// gemini_client.go
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

// ── Request types ────────────────────────────────────────────────────────────

type geminiRequest struct {
	SystemInstruction *geminiContent  `json:"system_instruction,omitempty"` // ← instructions go here
	Contents          []geminiContent `json:"contents"`                     // ← conversation history
	Tools             []geminiTool    `json:"tools,omitempty"`
	GenerationConfig  *geminiConfig   `json:"generationConfig,omitempty"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"` // "user" | "model"  (not "assistant")
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

// ── Response types ───────────────────────────────────────────────────────────

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

// toGeminiContents converts flat ChatMessage history into Gemini's contents format
func toGeminiContents(history []ChatMessage) []geminiContent {
	contents := make([]geminiContent, len(history))
	for i, msg := range history {
		text, _ := msg.Content.(string)

		contents[i] = geminiContent{
			Role: msg.Role,
			Parts: []geminiPart{
				{Text: text},
			},
		}
	}
	return contents
}

type geminiClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewGeminiClient(apiKey string) *geminiClient {
	return &geminiClient{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: geminiTimeout},
	}
}

func (c *geminiClient) complete(
	ctx context.Context,
	system string, // instructions
	history []ChatMessage, // user/model turns
	tools []Tool,
) (*geminiResponse, error) {

	req := geminiRequest{
		// System instruction — separate from contents, same as Anthropic's "system"
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

	url := fmt.Sprintf(geminiAPI, geminiModel, c.apiKey)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
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

	return &result, nil
}
