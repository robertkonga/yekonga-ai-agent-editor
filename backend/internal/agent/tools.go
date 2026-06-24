package agent

import "encoding/json"

type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  Schema `json:"parameters"`
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

// GetAvailableTools returns the list of tools available to the local LLM agent
func GetAvailableTools() []Tool {
	return []Tool{
		{
			Name:        "run_shell",
			Description: "Run a bash shell command in the secure Docker sandbox. Use this to create files, read files, or run scripts. Output will be returned as text.",
			Parameters: Schema{
				Type: "object",
				Properties: map[string]Property{
					"command": {
						Type:        "string",
						Description: "The bash command to run.",
					},
				},
				Required: []string{"command"},
			},
		},
		{
			Name:        "write_file",
			Description: "Write content to a file in the workspace.",
			Parameters: Schema{
				Type: "object",
				Properties: map[string]Property{
					"path": {
						Type:        "string",
						Description: "The absolute or relative path to the file inside the workspace (/workspace/...).",
					},
					"content": {
						Type:        "string",
						Description: "The exact content to write.",
					},
				},
				Required: []string{"path", "content"},
			},
		},
	}
}
