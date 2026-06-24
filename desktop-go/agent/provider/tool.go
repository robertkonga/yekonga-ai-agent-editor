package provider

import "encoding/json"

// Tool matches the Anthropic API tool schema
type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema Schema `json:"input_schema"`
}

type Schema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required,omitempty"`
}

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type ToolCall struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}
