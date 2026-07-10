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
	ollamaBaseURL = "http://core.tz:11434"
	ollamaTimeout = 10 * time.Minute
)

// Supported models:	"qwen3", "qwen3:8b", "qwen3:14b", "qwen3:32b",
//
//	"deepseek-r1", "deepseek-r1:7b", "deepseek-r1:14b", "deepseek-r1:32b"
const ollamaModel = "qwen3"

type OllamaProvider struct {
	modelName  string
	baseURL    string
	httpClient *http.Client
}

// ── Ollama Internal Types ───────────────────────────────────────────────────

type ollamaChatRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Tools    []ollamaTool    `json:"tools,omitempty"`
	Stream   bool            `json:"stream"`
	Options  ollamaOptions   `json:"options,omitempty"`
}

type ollamaOptions struct {
	Temperature float64 `json:"temperature,omitempty"`
	NumPredict  int     `json:"num_predict,omitempty"`
}

type ollamaMessage struct {
	Role      string           `json:"role"`
	Content   string           `json:"content"`
	Reasoning string           `json:"reasoning"`
	ToolCalls []ollamaToolCall `json:"tool_calls,omitempty"`
}

type ollamaToolCall struct {
	Function ollamaToolCallFunction `json:"function"`
}

type ollamaToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ollamaTool struct {
	Type     string             `json:"type"` // always "function"
	Function ollamaToolFunction `json:"function"`
}

type ollamaToolFunction struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  Schema `json:"parameters"`
}

// type ollamaChatResponse struct {
// 	Model   string        `json:"model"`
// 	Message ollamaMessage `json:"message"`
// 	Done    bool          `json:"done"`
// }

type ollamaChatResponse struct {
	Model   string `json:"model"`
	Choices []struct {
		Message      ollamaMessage `json:"message"`
		FinishReason string        `json:"finish_reason"`
	} `json:"choices"`
}

func NewOllamaProvider(modelName string) *OllamaProvider {
	if modelName == "" {
		modelName = ollamaModel
	}
	return &OllamaProvider{
		modelName:  modelName,
		baseURL:    ollamaBaseURL,
		httpClient: &http.Client{
			// Timeout: ollamaTimeout,
		},
	}
}

// NewOllamaProviderWithURL allows overriding the Ollama host (e.g. remote instance).
func NewOllamaProviderWithURL(model, baseURL string) *OllamaProvider {
	p := NewOllamaProvider(model)
	p.baseURL = baseURL
	return p
}

func (p *OllamaProvider) Complete(ctx context.Context, system string, history []types.ChatMessage, tools []Tool) (*LLMResponse, error) {
	messages := toOllamaMessages(system, history)

	req := ollamaChatRequest{
		Model:    p.modelName,
		Messages: messages,
		Stream:   false,
		Options: ollamaOptions{
			Temperature: 0.2,
			NumPredict:  8192,
		},
	}

	if len(tools) > 0 {
		req.Tools = toOllamaTools(tools)
	}

	payload, err := json.Marshal(req)
	// console.Error("ollama.payload", string(payload))
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	// url := fmt.Sprintf("%s/api/generate", p.baseURL)
	url := fmt.Sprintf("%s/v1/chat/completions", p.baseURL)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url, bytes.NewReader(payload),
	)
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

	console.Error("ollama.response", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama API %d: %s", resp.StatusCode, string(body))
	}

	var result ollamaChatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		console.Error("body", string(body))
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	choice := result.Choices[0]
	response := &LLMResponse{}

	// Reasoning content (deepseek-reasoner only) — surface as <think> wrapper.
	if choice.Message.Reasoning != "" {
		response.Content = "<think>\n" + choice.Message.Reasoning + "\n</think>\n"
	}
	if choice.Message.Content != "" {
		response.Content += choice.Message.Content
	}

	for _, tc := range choice.Message.ToolCalls {
		var args map[string]any
		_ = json.Unmarshal([]byte(tc.Function.Arguments), &args)
		input, _ := json.Marshal(args)
		response.ToolCalls = append(response.ToolCalls, types.ContentBlock{
			Type:  "tool_use",
			ID:    fmt.Sprintf("call_%d", time.Now().UnixNano()),
			Name:  tc.Function.Name,
			Input: input,
		})
	}
	return response, nil
}

// ── Conversion Helpers ──────────────────────────────────────────────────────

func toOllamaTools(tools []Tool) []ollamaTool {
	result := make([]ollamaTool, len(tools))
	for i, t := range tools {
		result[i] = ollamaTool{
			Type: "function",
			Function: ollamaToolFunction{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.InputSchema,
			},
		}
	}
	return result
}

// toOllamaMessages converts the agent's shared ChatMessage history into the
// flat message list that Ollama's /api/chat endpoint expects.
//
// Mapping:
//   - system prompt  → role "system"
//   - assistant text → role "assistant", Content = text
//   - tool_use block → role "assistant", ToolCalls populated
//   - tool_result    → role "tool",      Content = serialised result
//   - user string    → role "user"
func toOllamaMessages(system string, history []types.ChatMessage) []ollamaMessage {
	msgs := make([]ollamaMessage, 0, len(history)+1)

	if system != "" {
		msgs = append(msgs, ollamaMessage{Role: "system", Content: system})
	}

	for _, msg := range history {
		role := msg.Role

		switch v := msg.Content.(type) {

		case string:
			msgs = append(msgs, ollamaMessage{Role: role, Content: v})

		case []types.ContentBlock:
			msgs = append(msgs, blocksToOllamaMessages(role, v)...)

		case []any:
			// Deserialised from JSON — re-encode then decode into []ContentBlock.
			data, _ := json.Marshal(v)
			var blocks []types.ContentBlock
			if err := json.Unmarshal(data, &blocks); err == nil {
				msgs = append(msgs, blocksToOllamaMessages(role, blocks)...)
			}
		}
	}

	return msgs
}

func blocksToOllamaMessages(role string, blocks []types.ContentBlock) []ollamaMessage {
	var msgs []ollamaMessage

	var toolCalls []ollamaToolCall
	var textContent string

	for _, block := range blocks {
		switch block.Type {
		case "tool_use":
			args := string(block.Input)
			if args == "" || args == "null" {
				args = "{}"
			}

			toolCalls = append(toolCalls, ollamaToolCall{
				Function: ollamaToolCallFunction{
					Name:      block.Name,
					Arguments: args,
				},
			})

		case "tool_result":
			// Each tool result becomes its own "tool" role message.
			resultText := fmt.Sprintf("%v", block.Content)
			msgs = append(msgs, ollamaMessage{Role: "tool", Content: resultText})

		case "text":
			textContent += block.Content
		}
	}

	// Flush any accumulated tool calls as a single assistant message.
	if len(toolCalls) > 0 {
		msgs = append(msgs, ollamaMessage{
			Role:      "assistant",
			ToolCalls: toolCalls,
		})
	}

	// Flush plain text.
	if textContent != "" {
		msgs = append(msgs, ollamaMessage{Role: role, Content: textContent})
	}

	return msgs
}
