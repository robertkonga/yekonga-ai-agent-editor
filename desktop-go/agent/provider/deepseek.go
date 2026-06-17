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
	deepSeekBaseURL = "https://api.deepseek.com/v1"
	deepSeekTimeout = 120 * time.Second
)

// Available models:
//
//	"deepseek-chat"    — DeepSeek-V3, general-purpose (recommended)
//	"deepseek-reasoner" — DeepSeek-R1, chain-of-thought reasoning
const deepSeekDefaultModel = "deepseek-chat"

type DeepSeekProvider struct {
	apiKey     string
	modelName  string
	httpClient *http.Client
}

func NewDeepSeekProvider(apiKey string, modelName string) *DeepSeekProvider {
	return &DeepSeekProvider{
		apiKey:     apiKey,
		modelName:  modelName,
		httpClient: &http.Client{Timeout: deepSeekTimeout},
	}
}

func NewDeepSeekProviderWithModel(apiKey, model string) *DeepSeekProvider {
	p := NewDeepSeekProvider(apiKey, model)
	p.modelName = model
	return p
}

func (p *DeepSeekProvider) Complete(ctx context.Context, system string, history []types.ChatMessage, tools []Tool) (*LLMResponse, error) {
	messages := toDeepSeekMessages(system, history)

	req := deepSeekRequest{
		Model:       p.modelName,
		Messages:    messages,
		MaxTokens:   8192,
		Temperature: 0.2,
	}

	if len(tools) > 0 {
		req.Tools = toDeepSeekTools(tools)
		req.ToolChoice = "auto"
	}

	payload, err := json.Marshal(req)
	console.Error("deepseek.payload", string(payload))
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := deepSeekBaseURL + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr deepSeekErrorResponse
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Error.Message != "" {
			return nil, fmt.Errorf("DeepSeek API %d %s: %s", resp.StatusCode, apiErr.Error.Type, apiErr.Error.Message)
		}
		return nil, fmt.Errorf("DeepSeek API %d: %s", resp.StatusCode, string(body))
	}

	var result deepSeekResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("DeepSeek returned no choices")
	}

	choice := result.Choices[0]
	response := &LLMResponse{}

	// Reasoning content (deepseek-reasoner only) — prepend as context if present.
	if choice.Message.ReasoningContent != "" {
		response.Content = "<think>\n" + choice.Message.ReasoningContent + "\n</think>\n"
	}
	if choice.Message.Content != "" {
		response.Content += choice.Message.Content
	}

	for _, tc := range choice.Message.ToolCalls {
		var args map[string]any
		json.Unmarshal([]byte(tc.Function.Arguments), &args) //nolint:errcheck
		input, _ := json.Marshal(args)
		response.ToolCalls = append(response.ToolCalls, types.ContentBlock{
			Type:  "tool_use",
			ID:    tc.ID,
			Name:  tc.Function.Name,
			Input: input,
		})
	}

	return response, nil
}

// ── DeepSeek Internal Types ─────────────────────────────────────────────────

type deepSeekRequest struct {
	Model       string            `json:"model"`
	Messages    []deepSeekMessage `json:"messages"`
	Tools       []deepSeekTool    `json:"tools,omitempty"`
	ToolChoice  string            `json:"tool_choice,omitempty"`
	MaxTokens   int               `json:"max_tokens,omitempty"`
	Temperature float64           `json:"temperature,omitempty"`
	Stream      bool              `json:"stream"`
}

type deepSeekMessage struct {
	Role             string             `json:"role"`
	Content          string             `json:"content"`
	ToolCalls        []deepSeekToolCall `json:"tool_calls,omitempty"`
	ToolCallID       string             `json:"tool_call_id,omitempty"`
	ReasoningContent string             `json:"reasoning_content,omitempty"`
}

type deepSeekToolCall struct {
	ID       string                   `json:"id"`
	Type     string                   `json:"type"` // always "function"
	Function deepSeekToolCallFunction `json:"function"`
}

type deepSeekToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // raw JSON string (OpenAI convention)
}

type deepSeekTool struct {
	Type     string               `json:"type"` // always "function"
	Function deepSeekToolFunction `json:"function"`
}

type deepSeekToolFunction struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  Schema `json:"parameters"`
}

type deepSeekResponse struct {
	ID      string           `json:"id"`
	Model   string           `json:"model"`
	Choices []deepSeekChoice `json:"choices"`
	Usage   deepSeekUsage    `json:"usage"`
}

type deepSeekChoice struct {
	Index        int             `json:"index"`
	Message      deepSeekMessage `json:"message"`
	FinishReason string          `json:"finish_reason"`
}

type deepSeekUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type deepSeekErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// ── Conversion Helpers ──────────────────────────────────────────────────────

func toDeepSeekTools(tools []Tool) []deepSeekTool {
	result := make([]deepSeekTool, len(tools))
	for i, t := range tools {
		result[i] = deepSeekTool{
			Type: "function",
			Function: deepSeekToolFunction{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.InputSchema,
			},
		}
	}
	return result
}

// toDeepSeekMessages converts the agent's ChatMessage history into the flat
// message list that DeepSeek's OpenAI-compatible endpoint expects.
//
// Role mapping:
//   - system prompt  → role "system"
//   - user string    → role "user"
//   - assistant text → role "assistant", content = text
//   - tool_use block → role "assistant", tool_calls populated
//   - tool_result    → role "tool",      tool_call_id + content
func toDeepSeekMessages(system string, history []types.ChatMessage) []deepSeekMessage {
	msgs := make([]deepSeekMessage, 0, len(history)+1)

	if system != "" {
		msgs = append(msgs, deepSeekMessage{Role: "system", Content: system})
	}

	for _, msg := range history {
		role := msg.Role

		switch v := msg.Content.(type) {
		case string:
			msgs = append(msgs, deepSeekMessage{Role: role, Content: v})

		case []types.ContentBlock:
			msgs = append(msgs, blocksToDeepSeekMessages(role, v)...)

		case []any:
			data, _ := json.Marshal(v)
			var blocks []types.ContentBlock
			if err := json.Unmarshal(data, &blocks); err == nil {
				msgs = append(msgs, blocksToDeepSeekMessages(role, blocks)...)
			}
		}
	}

	return msgs
}

func blocksToDeepSeekMessages(role string, blocks []types.ContentBlock) []deepSeekMessage {
	var msgs []deepSeekMessage

	var toolCalls []deepSeekToolCall
	var textContent string

	for _, block := range blocks {
		switch block.Type {
		case "tool_use":
			// DeepSeek expects arguments as a raw JSON string, not a parsed object.
			args := string(block.Input)
			if args == "" {
				args = "{}"
			}
			toolCalls = append(toolCalls, deepSeekToolCall{
				ID:   block.ID,
				Type: "function",
				Function: deepSeekToolCallFunction{
					Name:      block.Name,
					Arguments: args,
				},
			})

		case "tool_result":
			// Must follow the paired assistant tool_call message.
			// Flush any buffered tool calls first.
			if len(toolCalls) > 0 {
				msgs = append(msgs, deepSeekMessage{
					Role:      "assistant",
					ToolCalls: toolCalls,
				})
				toolCalls = nil
			}
			resultText := fmt.Sprintf("%v", block.Content)
			msgs = append(msgs, deepSeekMessage{
				Role:       "tool",
				ToolCallID: block.ID,
				Content:    resultText,
			})

		case "text":
			textContent += block.Content
		}
	}

	// Flush tool calls (no tool_result followed them in this batch).
	if len(toolCalls) > 0 {
		msgs = append(msgs, deepSeekMessage{
			Role:      "assistant",
			ToolCalls: toolCalls,
		})
	}

	// Flush plain text.
	if textContent != "" {
		msgs = append(msgs, deepSeekMessage{Role: role, Content: textContent})
	}

	return msgs
}

// stripThinkTags removes <think>…</think> blocks that deepseek-reasoner
// prepends to its output, returning only the final answer text.
// Useful if you don't want to surface chain-of-thought to users.
func stripThinkTags(s string) string {
	const open, close = "<think>\n", "\n</think>\n"
	for {
		start := indexStr(s, open)
		if start == -1 {
			break
		}
		end := indexStr(s[start:], close)
		if end == -1 {
			break
		}
		s = s[:start] + s[start+end+len(close):]
	}
	return s
}

func indexStr(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

// Ensure DeepSeekProvider satisfies the same provider interface at compile time.
// Uncomment once your Provider interface is defined in this package:
// var _ Provider = (*DeepSeekProvider)(nil)

// Example usage:
//
//	provider := agent.NewDeepSeekProvider(os.Getenv("DEEPSEEK_API_KEY"))
//	// or for the reasoner model:
//	provider := agent.NewDeepSeekProviderWithModel(os.Getenv("DEEPSEEK_API_KEY"), "deepseek-reasoner")
//
//	resp, err := provider.Complete(ctx, systemPrompt, history, tools)
