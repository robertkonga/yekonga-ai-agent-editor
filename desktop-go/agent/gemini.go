package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	geminiAPI     = "https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s"
	geminiModel   = "gemini-2.0-flash"
	geminiTimeout = 120 * time.Second
)

type GeminiProvider struct {
	apiKey     string
	httpClient *http.Client
}

func NewGeminiProvider(apiKey string) *GeminiProvider {
	return &GeminiProvider{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: geminiTimeout},
	}
}

func (p *GeminiProvider) Complete(ctx context.Context, system string, history []ChatMessage, tools []Tool) (*LLMResponse, error) {
	req := geminiRequest{
		SystemInstruction: &geminiContent{
			Parts: []geminiPart{{Text: system}},
		},
		Contents: toGeminiContents(history),
		GenerationConfig: &geminiConfig{
			MaxOutputTokens: 8192,
			Temperature:     0.2,
		},
	}

	if len(tools) > 0 {
		req.Tools = []geminiTool{
			{FunctionDeclarations: tools},
		}
	}

	payload, err := json.Marshal(req)
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
			response.ToolCalls = append(response.ToolCalls, ContentBlock{
				Type:  "tool_use",
				ID:    fmt.Sprintf("call_%d", time.Now().UnixNano()), // Gemini doesn't always provide IDs in v1beta
				Name:  part.FunctionCall.Name,
				Input: input,
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
}

type functionCall struct {
	Name string         `json:"name"`
	Args map[string]any `json:"args"`
}

type functionResponse struct {
	Name     string         `json:"name"`
	Response map[string]any `json:"response"`
}

type geminiTool struct {
	FunctionDeclarations []Tool `json:"function_declarations"`
}

type geminiConfig struct {
	MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
	Temperature     float64 `json:"temperature,omitempty"`
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

func toGeminiContents(history []ChatMessage) []geminiContent {
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
		case []ContentBlock:
			for _, block := range v {
				if block.Type == "tool_result" {
					parts = append(parts, geminiPart{
						FunctionResponse: &functionResponse{
							Name:     block.Name, // We need to ensure Name is preserved
							Response: map[string]any{"result": block.Content},
						},
					})
				} else if block.Type == "tool_use" {
					var args map[string]any
					json.Unmarshal(block.Input, &args)
					parts = append(parts, geminiPart{
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
			var blocks []ContentBlock
			if err := json.Unmarshal(data, &blocks); err == nil {
				for _, block := range blocks {
					if block.Type == "tool_result" {
						parts = append(parts, geminiPart{
							FunctionResponse: &functionResponse{
								Name:     block.Name,
								Response: map[string]any{"result": block.Content},
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
