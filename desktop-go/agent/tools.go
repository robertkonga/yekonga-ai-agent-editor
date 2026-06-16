package agent

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

var agentTools = []Tool{
	{
		Name:        "read_file",
		Description: "Read the full content of a file by its absolute path",
		InputSchema: Schema{
			Type: "object",
			Properties: map[string]Property{
				"path": {
					Type:        "string",
					Description: "Absolute path to the file",
				},
			},
			Required: []string{"path"},
		},
	},
	{
		Name:        "list_directory",
		Description: "List all files and folders inside a directory",
		InputSchema: Schema{
			Type: "object",
			Properties: map[string]Property{
				"path": {
					Type:        "string",
					Description: "Absolute path to the directory",
				},
			},
			Required: []string{"path"},
		},
	},
	{
		Name:        "write_file",
		Description: "Write or overwrite a file with new content",
		InputSchema: Schema{
			Type: "object",
			Properties: map[string]Property{
				"path":    {Type: "string"},
				"content": {Type: "string"},
			},
			Required: []string{"path", "content"},
		},
	},
	{
		Name:        "search_in_files",
		Description: "Search for a string or pattern across all files in the workspace",
		InputSchema: Schema{
			Type: "object",
			Properties: map[string]Property{
				"query": {Type: "string"},
				"root":  {Type: "string", Description: "Root directory to search from"},
			},
			Required: []string{"query", "root"},
		},
	},
	{
		Name:        "run_shell",
		Description: "Run a bash shell command in the secure Docker sandbox. Use this to create files, read files, or run scripts. Output will be returned as text.",
		InputSchema: Schema{
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
}
