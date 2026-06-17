package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"yekonga-builder/agent/types"
	"yekonga-builder/console"
)

const (
	geminiAPI = "https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s"
	// geminiModel   = "gemini-2.5-flash"
	geminiModel   = "gemini-3.1-flash-lite" // gemini-3.1-pro | gemini-3.5-flash | gemini-3.1-flash-lite | gemini-3-flash | gemini-2.5-pro, gemini-2.5-flash, gemini-2.5-flash-lite
	geminiTimeout = 120 * time.Second
)

type GeminiProvider struct {
	apiKey     string
	modelName  string
	httpClient *http.Client
}

func NewGeminiProvider(apiKey string, modelName string) *GeminiProvider {
	return &GeminiProvider{
		apiKey:     apiKey,
		modelName:  modelName,
		httpClient: &http.Client{Timeout: geminiTimeout},
	}
}

func (p *GeminiProvider) Complete(ctx context.Context, system string, history []types.ChatMessage, tools []Tool) (*LLMResponse, error) {
	req := geminiRequest{
		SystemInstruction: &geminiContent{
			Parts: []geminiPart{{Text: system}},
		},
		Contents: toGeminiContents(history),
		GenerationConfig: &geminiConfig{
			MaxOutputTokens: 8192,
			Temperature:     0.2,
			ThinkingConfig: &geminiThinkingConfig{
				ThinkingBudget: 0, // 0 = disable thinking
			},
		},
	}

	if len(tools) > 0 {
		declarations := make([]geminiFunctionDeclaration, len(tools))
		for i, t := range tools {
			declarations[i] = geminiFunctionDeclaration{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.InputSchema,
			}
		}
		req.Tools = []geminiTool{
			{FunctionDeclarations: declarations},
		}
	}

	payload, err := json.Marshal(req)
	console.Error("gemini.payload", string(payload))
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf(geminiAPI, geminiModel, p.apiKey)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	console.Error("gemini.response", string(body))

	var result geminiResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	if result.Error != nil {
		return nil, fmt.Errorf("Gemini API %d %s: %s", result.Error.Code, result.Error.Status, result.Error.Message)
	}
	if len(result.Candidates) == 0 {
		return nil, fmt.Errorf("Gemini returned no candidates")
	}

	candidate := result.Candidates[0]
	response := &LLMResponse{}

	for _, part := range candidate.Content.Parts {
		if part.Text != "" {
			response.Content += part.Text
		}
		if part.FunctionCall != nil {
			input, _ := json.Marshal(part.FunctionCall.Args)
			response.ToolCalls = append(response.ToolCalls, types.ContentBlock{
				Type:             "tool_use",
				ID:               fmt.Sprintf("call_%d", time.Now().UnixNano()), // Gemini doesn't always provide IDs in v1beta
				Name:             part.FunctionCall.Name,
				Input:            input,
				ThoughtSignature: part.ThoughtSignature, // ← carry it
			})
		}
	}

	return response, nil
}

// ── Gemini Internal Types ───────────────────────────────────────────────────

type geminiRequest struct {
	SystemInstruction *geminiContent  `json:"system_instruction,omitempty"`
	Contents          []geminiContent `json:"contents"`
	Tools             []geminiTool    `json:"tools,omitempty"`
	GenerationConfig  *geminiConfig   `json:"generationConfig,omitempty"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text             string            `json:"text,omitempty"`
	FunctionCall     *functionCall     `json:"functionCall,omitempty"`
	FunctionResponse *functionResponse `json:"functionResponse,omitempty"`
	Thought          bool              `json:"thought,omitempty"`          // thinking text marker
	ThoughtSignature string            `json:"thoughtSignature,omitempty"` // ← add this
}

type functionCall struct {
	ID   string         `json:"id,omitempty"` // use this instead of generating your own
	Name string         `json:"name"`
	Args map[string]any `json:"args"`
}

type functionResponse struct {
	Name     string         `json:"name"`
	Response map[string]any `json:"response"`
}

type geminiFunctionDeclaration struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  Schema `json:"parameters"`
}

type geminiTool struct {
	FunctionDeclarations []geminiFunctionDeclaration `json:"function_declarations"`
}

type geminiConfig struct {
	MaxOutputTokens int                   `json:"maxOutputTokens,omitempty"`
	Temperature     float64               `json:"temperature,omitempty"`
	ThinkingConfig  *geminiThinkingConfig `json:"thinkingConfig"`
}

type geminiThinkingConfig struct {
	ThinkingBudget int `json:"thinkingBudget"`
}

type geminiResponse struct {
	Candidates []geminiCandidate `json:"candidates"`
	Error      *geminiError      `json:"error,omitempty"`
}

type geminiCandidate struct {
	Content      geminiContent `json:"content"`
	FinishReason string        `json:"finishReason"`
}

type geminiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func toGeminiContents(history []types.ChatMessage) []geminiContent {
	contents := make([]geminiContent, 0, len(history))
	for _, msg := range history {
		role := msg.Role
		if role == "assistant" {
			role = "model"
		}

		// Handle complex content (tool results or tool calls)
		var parts []geminiPart
		switch v := msg.Content.(type) {
		case string:
			parts = []geminiPart{{Text: v}}
		case []types.ContentBlock:
			for _, block := range v {
				switch block.Type {
				case "tool_result":
					parts = append(parts, geminiPart{
						FunctionResponse: &functionResponse{
							Name:     block.Name, // We need to ensure Name is preserved
							Response: map[string]any{"result": block.Content},
						},
					})
				case "tool_use":
					var args map[string]any
					json.Unmarshal(block.Input, &args)
					parts = append(parts, geminiPart{
						ThoughtSignature: block.ThoughtSignature, // ← echo back
						FunctionCall: &functionCall{
							Name: block.Name,
							Args: args,
						},
					})
				}
			}
		case []any:
			// If it's from JSON unmarshal, it might be []any
			data, _ := json.Marshal(v)
			var blocks []types.ContentBlock
			if err := json.Unmarshal(data, &blocks); err == nil {
				for _, block := range blocks {
					switch block.Type {
					case "tool_result":
						parts = append(parts, geminiPart{
							FunctionResponse: &functionResponse{
								Name:     block.Name,
								Response: map[string]any{"result": block.Content},
							},
						})
					case "tool_use":
						var args map[string]any
						json.Unmarshal(block.Input, &args)
						parts = append(parts, geminiPart{
							ThoughtSignature: block.ThoughtSignature, // ← echo back
							FunctionCall: &functionCall{
								Name: block.Name,
								Args: args,
							},
						})
					}
				}
			}
		}

		if len(parts) > 0 {
			contents = append(contents, geminiContent{
				Role:  role,
				Parts: parts,
			})
		}
	}
	return contents
}
