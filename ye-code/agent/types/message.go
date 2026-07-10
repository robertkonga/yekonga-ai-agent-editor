package types

import "encoding/json"

type ChatMessage struct {
	Role    string `json:"role"`
	Content any    `json:"content"` // string OR []ContentBlock
}

type ContentBlock struct {
	Type             string          `json:"type"`
	Text             string          `json:"text,omitempty"`
	ID               string          `json:"id,omitempty"`               // tool_use
	Name             string          `json:"name,omitempty"`             // tool_use
	Input            json.RawMessage `json:"input,omitempty"`            // tool_use
	ToolUseID        string          `json:"tool_use_id,omitempty"`      // tool_result
	Content          string          `json:"content,omitempty"`          // tool_result
	ThoughtSignature string          `json:"thoughtSignature,omitempty"` // ← add this

}
