package agent

import (
	"context"
)

// LLMProvider defines the interface for interacting with different LLM models.
type LLMProvider interface {
	Complete(ctx context.Context, system string, history []ChatMessage, tools []Tool) (*LLMResponse, error)
}

// LLMResponse is a unified response format for all providers.
type LLMResponse struct {
	Content   string         `json:"content"`
	ToolCalls []ContentBlock `json:"tool_calls,omitempty"`
}
