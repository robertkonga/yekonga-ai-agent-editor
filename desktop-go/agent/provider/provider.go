package provider

import (
	"context"
	"yekonga-builder/agent/types"
)

// LLMProvider defines the interface for interacting with different LLM models.
type LLMProvider interface {
	Complete(ctx context.Context, system string, history []types.ChatMessage, tools []Tool) (*LLMResponse, error)
}

// LLMResponse is a unified response format for all providers.
type LLMResponse struct {
	Content   string               `json:"content"`
	ToolCalls []types.ContentBlock `json:"tool_calls,omitempty"`
}
