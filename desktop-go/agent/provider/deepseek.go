package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"yekonga-builder/agent/types"
	// "yekonga-builder/console"
)

const (
	deepSeekBaseURL = "https://api.deepseek.com/v1"
	deepSeekTimeout = 120 * time.Second
)

// Available models:
//
//	"deepseek-chat"     — DeepSeek-V3, general-purpose (recommended)
//	"deepseek-reasoner" — DeepSeek-R1, chain-of-thought reasoning
const deepSeekDefaultModel = "deepseek-chat"

type DeepSeekProvider struct {
	apiKey     string
	modelName  string
	httpClient *http.Client
}

func NewDeepSeekProvider(apiKey string, modelName string) *DeepSeekProvider {
	if modelName == "" {
		modelName = deepSeekDefaultModel
	}
	return &DeepSeekProvider{
		apiKey:     apiKey,
		modelName:  modelName,
		httpClient: &http.Client{Timeout: deepSeekTimeout},
	}
}

// NewDeepSeekProviderWithModel is an alias kept for backward compatibility.
func NewDeepSeekProviderWithModel(apiKey, model string) *DeepSeekProvider {
	return NewDeepSeekProvider(apiKey, model)
}

func (p *DeepSeekProvider) Complete(ctx context.Context, system string, history []types.ChatMessage, tools []Tool) (*LLMResponse, error) {
	messages, err := toDeepSeekMessages(system, history)
	if err != nil {
		return nil, fmt.Errorf("build messages: %w", err)
	}

	req := deepSeekRequest{
		Model:       p.modelName,
		Messages:    messages,
		MaxTokens:   8192,
		Temperature: 0.2,
		Stream:      false,
	}

	if len(tools) > 0 {
		req.Tools = toDeepSeekTools(tools)
		req.ToolChoice = "auto"
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	// console.Error("deepseek.payload", string(payload))

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
	// console.Error("deepseek.response", string(body))

	choice := result.Choices[0]
	response := &LLMResponse{}

	// Reasoning content (deepseek-reasoner only) — surface as <think> wrapper.
	if choice.Message.ReasoningContent != "" {
		response.Content = "<think>\n" + choice.Message.ReasoningContent + "\n</think>\n"
	}
	if choice.Message.Content != nil && *choice.Message.Content != "" {
		response.Content += *choice.Message.Content
	}

	for _, tc := range choice.Message.ToolCalls {
		var args map[string]any
		_ = json.Unmarshal([]byte(tc.Function.Arguments), &args)
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

// ── Internal Types ──────────────────────────────────────────────────────────

type deepSeekRequest struct {
	Model       string            `json:"model"`
	Messages    []deepSeekMessage `json:"messages"`
	Tools       []deepSeekTool    `json:"tools,omitempty"`
	ToolChoice  string            `json:"tool_choice,omitempty"`
	MaxTokens   int               `json:"max_tokens,omitempty"`
	Temperature float64           `json:"temperature,omitempty"`
	Stream      bool              `json:"stream"`
}

// deepSeekMessage represents a single turn in the conversation.
//
// Content is a *string (pointer) so that assistant tool-call messages can
// serialize as `"content": null` rather than `"content": ""`.  DeepSeek
// rejects a non-null empty string when tool_calls is also present.
type deepSeekMessage struct {
	Role             string             `json:"role"`
	Content          *string            `json:"content"` // null for assistant tool-call turns
	ToolCalls        []deepSeekToolCall `json:"tool_calls,omitempty"`
	ToolCallID       string             `json:"tool_call_id,omitempty"` // role="tool" only
	ReasoningContent string             `json:"reasoning_content,omitempty"`
}

type deepSeekToolCall struct {
	ID       string                   `json:"id"`
	Type     string                   `json:"type"` // always "function"
	Function deepSeekToolCallFunction `json:"function"`
}

type deepSeekToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // raw JSON string
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

// strPtr returns a pointer to s. Used to produce non-null JSON strings.
func strPtr(s string) *string { return &s }

// toDeepSeekMessages converts the agent history into the flat message list
// that DeepSeek's OpenAI-compatible endpoint expects.
//
// The critical invariant DeepSeek enforces:
//
//  1. An assistant message that contains tool_calls must have content = null.
//  2. Each tool result must be a separate message with role "tool" and a
//     tool_call_id that exactly matches the id of the tool_use block that
//     preceded it in the assistant message.
//  3. tool_use and tool_result blocks live in different ChatMessages in the
//     agent history:
//     - assistant message → []ContentBlock of type "tool_use" (+ optional "text")
//     - user message      → []ContentBlock of type "tool_result"
//
// This function never mixes tool_use and tool_result in the same pass.
func toDeepSeekMessages(system string, history []types.ChatMessage) ([]deepSeekMessage, error) {
	msgs := make([]deepSeekMessage, 0, len(history)+1)

	if system != "" {
		msgs = append(msgs, deepSeekMessage{Role: "system", Content: strPtr(system)})
	}

	for i, msg := range history {
		converted, err := chatMessageToDeepSeek(msg)
		if err != nil {
			return nil, fmt.Errorf("message[%d]: %w", i, err)
		}
		msgs = append(msgs, converted...)
	}

	return msgs, nil
}

// chatMessageToDeepSeek converts a single ChatMessage to one or more
// deepSeekMessages. A single ChatMessage can expand to multiple API messages
// when it contains a mix of tool_use calls (one assistant message) plus
// tool_result responses (one "tool" message per result).
func chatMessageToDeepSeek(msg types.ChatMessage) ([]deepSeekMessage, error) {
	switch v := msg.Content.(type) {

	case string:
		return []deepSeekMessage{{Role: msg.Role, Content: strPtr(v)}}, nil

	case []types.ContentBlock:
		return blocksToDeepSeekMessages(msg.Role, v)

	case []any:
		// History round-tripped through JSON loses concrete types; re-hydrate.
		data, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("re-marshal []any: %w", err)
		}
		var blocks []types.ContentBlock
		if err := json.Unmarshal(data, &blocks); err != nil {
			return nil, fmt.Errorf("unmarshal []ContentBlock: %w", err)
		}
		return blocksToDeepSeekMessages(msg.Role, blocks)

	default:
		return nil, fmt.Errorf("unsupported content type %T", msg.Content)
	}
}

// blocksToDeepSeekMessages maps a slice of ContentBlocks to DeepSeek messages.
//
// Block layout rules by role:
//
//	role="assistant"
//	  "text"     → assistant message with content string
//	  "tool_use" → assistant message with tool_calls array, content = null
//	               (text and tool_calls are merged into ONE assistant message)
//
//	role="user"
//	  "tool_result" → one "tool" message per result, tool_call_id = block.ID
//	  "text"        → user message with content string
func blocksToDeepSeekMessages(role string, blocks []types.ContentBlock) ([]deepSeekMessage, error) {
	var msgs []deepSeekMessage

	// ── Assistant turn ──────────────────────────────────────────────────────
	if role == "assistant" {
		var toolCalls []deepSeekToolCall
		var textParts []string

		for _, block := range blocks {
			switch block.Type {
			case "tool_use":
				args := string(block.Input)
				if args == "" || args == "null" {
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

			case "text":
				if block.Content != "" {
					textParts = append(textParts, block.Content)
				}

			// tool_result should never appear in an assistant message; skip
			// gracefully rather than producing a malformed API call.
			case "tool_result":
				return nil, fmt.Errorf("tool_result block found in assistant message (block.ID=%q); check agentic loop history construction", block.ID)
			}
		}

		// Build a single assistant message.
		// DeepSeek rules:
		//   • If tool_calls present → content MUST be null (not "")
		//   • If no tool_calls     → content is the text (may be "")
		m := deepSeekMessage{Role: "assistant"}
		if len(toolCalls) > 0 {
			m.ToolCalls = toolCalls
			m.Content = nil // explicit null
		} else {
			text := strings.Join(textParts, "\n")
			m.Content = strPtr(text)
		}
		msgs = append(msgs, m)
		return msgs, nil
	}

	// ── User turn ───────────────────────────────────────────────────────────
	// tool_result blocks become role="tool" messages.
	// Plain text blocks become role="user" messages.
	for _, block := range blocks {
		switch block.Type {
		case "tool_result":
			// block.ID here is the tool_call_id — it must match the id that was
			// assigned to the tool_use block in the preceding assistant message.
			content := extractToolResultContent(block.Content)
			toolCallID := block.ID
			if toolCallID == "" {
				toolCallID = block.ToolUseID
			}
			msgs = append(msgs, deepSeekMessage{
				Role:       "tool",
				ToolCallID: toolCallID,
				Content:    strPtr(content),
			})

		case "text":
			if block.Content != "" {
				msgs = append(msgs, deepSeekMessage{
					Role:    "user",
					Content: strPtr(block.Content),
				})
			}
		}
	}

	return msgs, nil
}

// extractToolResultContent safely converts the tool_result Content field
// (which is typed as any in ContentBlock) to a plain string.
func extractToolResultContent(raw any) string {
	if raw == nil {
		return ""
	}
	switch v := raw.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		// Structured content (e.g. []ContentBlock) — marshal to JSON string.
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", raw)
		}
		return string(b)
	}
}

// ── Utility ─────────────────────────────────────────────────────────────────

// StripThinkTags removes <think>…</think> blocks that deepseek-reasoner
// prepends to its output, returning only the final answer text.
func StripThinkTags(s string) string {
	const open, close = "<think>\n", "\n</think>\n"
	for {
		start := strings.Index(s, open)
		if start == -1 {
			break
		}
		end := strings.Index(s[start+len(open):], close)
		if end == -1 {
			break
		}
		s = s[:start] + s[start+len(open)+end+len(close):]
	}
	return s
}

// Ensure DeepSeekProvider satisfies the Provider interface at compile time.
// Uncomment once Provider is defined in this package:
// var _ Provider = (*DeepSeekProvider)(nil)
