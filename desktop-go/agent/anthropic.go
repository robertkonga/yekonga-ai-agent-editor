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
	anthropicAPI     = "https://api.anthropic.com/v1/messages"
	anthropicVersion = "2023-06-01"
	anthropicModel   = "claude-3-5-sonnet-20241022"
	anthropicTimeout = 120 * time.Second
)

type AnthropicProvider struct {
	apiKey     string
	httpClient *http.Client
}

func NewAnthropicProvider(apiKey string) *AnthropicProvider {
	return &AnthropicProvider{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: anthropicTimeout},
	}
}

func (p *AnthropicProvider) Complete(ctx context.Context, system string, history []ChatMessage, tools []Tool) (*LLMResponse, error) {
	messages := make([]llmMessage, len(history))
	for i, msg := range history {
		content := ""
		if s, ok := msg.Content.(string); ok {
			content = s
		} else {
			// Handle tool results/calls if needed for Anthropic
			data, _ := json.Marshal(msg.Content)
			content = string(data)
		}
		messages[i] = llmMessage{Role: msg.Role, Content: content}
	}

	reqBody := llmRequest{
		Model:     anthropicModel,
		MaxTokens: 8192,
		System:    system,
		Messages:  messages,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, anthropicAPI, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", anthropicVersion)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request: %w", err)
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Anthropic API error %d: %s", resp.StatusCode, string(rawBody))
	}

	var llmResp llmResponse
	if err := json.Unmarshal(rawBody, &llmResp); err != nil {
		return nil, fmt.Errorf("parse response JSON: %w", err)
	}

	response := &LLMResponse{}
	for _, block := range llmResp.Content {
		if block.Type == "text" {
			response.Content += block.Text
		}
		// Handle tool_use blocks if we extend Anthropic for general agent use
	}

	return response, nil
}

// ── Anthropic Internal Types ────────────────────────────────────────────────

type llmRequest struct {
	Model     string       `json:"model"`
	MaxTokens int          `json:"max_tokens"`
	System    string       `json:"system"`
	Messages  []llmMessage `json:"messages"`
}

type llmMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type llmResponse struct {
	Content []llmContentBlock `json:"content"`
	Error   *llmError         `json:"error,omitempty"`
	Usage   llmUsage          `json:"usage"`
}

type llmContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type llmError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type llmUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}
