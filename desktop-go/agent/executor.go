package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (a *Agent) executeTool(name string, rawInput json.RawMessage) (string, error) {
	// Helper to resolve paths relative to active workspace if needed
	resolvePath := func(p string) string {
		if filepath.IsAbs(p) {
			return p
		}
		if a.ActivePath != "" {
			return filepath.Join(a.ActivePath, p)
		}
		return p
	}

	switch name {
	case "read_file":
		var in struct {
			Path string `json:"path"`
		}
		json.Unmarshal(rawInput, &in)
		path := resolvePath(in.Path)
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(data), nil

	case "list_directory":
		var in struct {
			Path string `json:"path"`
		}
		json.Unmarshal(rawInput, &in)
		path := resolvePath(in.Path)
		if path == "" {
			path = "."
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			return "", err
		}
		var lines []string
		for _, e := range entries {
			prefix := "  "
			if e.IsDir() {
				prefix = "📁 "
			}
			lines = append(lines, prefix+e.Name())
		}
		return strings.Join(lines, "\n"), nil

	case "write_file":
		var in struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		}
		json.Unmarshal(rawInput, &in)
		path := resolvePath(in.Path)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return "", err
		}
		if err := os.WriteFile(path, []byte(in.Content), 0644); err != nil {
			return "", err
		}
		return fmt.Sprintf("written: %s", path), nil

	case "search_in_files":
		var in struct {
			Query string `json:"query"`
			Root  string `json:"root"`
		}
		json.Unmarshal(rawInput, &in)
		root := resolvePath(in.Root)
		// return searchFiles(root, in.Query)
		return fmt.Sprintf("Search feature for root '%s' is not implemented yet", root), nil

	case "initialize_project":
		var in struct {
			Prompt   string `json:"prompt"`
			RootPath string `json:"root_path"`
		}
		json.Unmarshal(rawInput, &in)
		path := resolvePath(in.RootPath)
		err := a.GenerateProject(in.Prompt, path, "")
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Project initialized at: %s", path), nil

	case "generate_dataschema":
		var in struct {
			Description string `json:"description"`
		}
		json.Unmarshal(rawInput, &in)
		schema, err := a.GenerateSchema(in.Description)
		if err != nil {
			return "", err
		}
		return schema, nil

	case "generate_crud":
		var in struct {
			SchemaJSON string `json:"schema_json"`
			TargetPath string `json:"target_path"`
		}
		json.Unmarshal(rawInput, &in)
		path := resolvePath(in.TargetPath)
		// Logic for CRUD generation would go here
		// For now, return a placeholder
		return fmt.Sprintf("CRUD components generation started for target: %s", path), nil

	case "read_template":
		var in struct {
			Name string `json:"name"`
		}
		json.Unmarshal(rawInput, &in)
		content := a.getSystemInstruction(in.Name)
		if content == "" {
			return "", fmt.Errorf("template not found: %s", in.Name)
		}
		return content, nil

	case "list_templates":
		entries, err := assets.ReadDir("templates")
		if err != nil {
			return "", err
		}
		var names []string
		for _, e := range entries {
			if !e.IsDir() {
				names = append(names, strings.TrimSuffix(e.Name(), ".md"))
			}
		}
		return strings.Join(names, ", "), nil

	case "get_active_workspace_path":
		if a.ActivePath == "" {
			return "No active workspace path set. Please use an absolute path or initialize a project.", nil
		}
		return a.ActivePath, nil

	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}
