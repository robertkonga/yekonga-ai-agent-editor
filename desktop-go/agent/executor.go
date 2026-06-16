package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (a *Agent) executeTool(name string, rawInput json.RawMessage) (string, error) {
	switch name {
	case "read_file":
		var in struct {
			Path string `json:"path"`
		}
		json.Unmarshal(rawInput, &in)
		data, err := os.ReadFile(in.Path)
		if err != nil {
			return "", err
		}
		return string(data), nil

	case "list_directory":
		var in struct {
			Path string `json:"path"`
		}
		json.Unmarshal(rawInput, &in)
		entries, err := os.ReadDir(in.Path)
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
		if err := os.MkdirAll(filepath.Dir(in.Path), 0755); err != nil {
			return "", err
		}
		if err := os.WriteFile(in.Path, []byte(in.Content), 0644); err != nil {
			return "", err
		}
		return fmt.Sprintf("written: %s", in.Path), nil

	case "search_in_files":
		var in struct {
			Query string `json:"query"`
			Root  string `json:"root"`
		}
		json.Unmarshal(rawInput, &in)
		// return searchFiles(in.Root, in.Query)
		return "", nil

	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}
