package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"yekonga-builder/agent/provider"
	agentTypes "yekonga-builder/agent/types"
	"yekonga-builder/types"
)

const scaffoldSystemPrompt = `You are a project scaffolding agent embedded in a code editor.
When given a project description, respond ONLY with a valid JSON object — no markdown fences,
no explanation, no preamble. 

The JSON must match this exact schema:
{
  "name": "project-name",
  "description": "one-line summary",
  "files": [
    {
      "path": "relative/path/to/file",
      "content": "complete file content"
    }
  ]
}

Rules:
- All paths must be relative and use forward slashes.
- Never use ".." in any path.
- Every file must contain complete, working content — no placeholders or TODOs.
- Include all config files required to run the project (package.json, tsconfig, go.mod, etc.).
- Follow the user's specified framework, language, and conventions exactly.`

// parsePlan extracts and validates the types.ScaffoldPlan from the raw LLM text.
func parsePlan(raw string) (*types.ScaffoldPlan, error) {
	cleaned := strings.TrimSpace(raw)

	// Strip optional ```json ... ``` fences
	if strings.HasPrefix(cleaned, "```") {
		start := strings.Index(cleaned, "\n")
		end := strings.LastIndex(cleaned, "```")
		if start != -1 && end > start {
			cleaned = strings.TrimSpace(cleaned[start+1 : end])
		}
	}

	var plan types.ScaffoldPlan
	if err := json.Unmarshal([]byte(cleaned), &plan); err != nil {
		return nil, fmt.Errorf("parse scaffold JSON: %w\nraw output: %.400s", err, cleaned)
	}

	if err := validatePlan(&plan); err != nil {
		return nil, err
	}
	return &plan, nil
}

// validatePlan rejects plans that could cause harm (path traversal, empty content, etc.)
func validatePlan(plan *types.ScaffoldPlan) error {
	if strings.TrimSpace(plan.Name) == "" {
		return errors.New("plan has no project name")
	}
	if len(plan.Files) == 0 {
		return errors.New("plan contains no files")
	}
	seen := make(map[string]bool, len(plan.Files))
	for i, f := range plan.Files {
		if strings.TrimSpace(f.Path) == "" {
			return fmt.Errorf("file[%d] has an empty path", i)
		}
		// Block path traversal
		clean := filepath.ToSlash(filepath.Clean(f.Path))
		if strings.HasPrefix(clean, "..") || filepath.IsAbs(f.Path) {
			return fmt.Errorf("file[%d] has unsafe path: %q", i, f.Path)
		}
		if seen[clean] {
			return fmt.Errorf("duplicate file path in plan: %q", clean)
		}
		seen[clean] = true
	}
	return nil
}

// GenerateProject is the Wails-bound method called from the Vue frontend.
func (a *Agent) GenerateProject(userPrompt string, rootPath string, extraConventions string) error {
	// Default to Anthropic for scaffolding as it's typically better at large JSON
	provider := provider.NewAnthropicProvider(a.ApiKey, "")

	system := scaffoldSystemPrompt
	if strings.TrimSpace(extraConventions) != "" {
		system += "\n\nAdditional conventions for this project:\n" + extraConventions
	}

	a.Emit(types.ScaffoldProgress{File: "Contacting LLM…"})

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	resp, err := provider.Complete(ctx, system, []agentTypes.ChatMessage{{Role: "user", Content: userPrompt}}, nil)
	if err != nil {
		a.Emit(types.ScaffoldProgress{Error: err.Error(), Done: true})
		return err
	}

	a.Emit(types.ScaffoldProgress{File: "Parsing plan…"})

	plan, err := parsePlan(resp.Content)
	if err != nil {
		a.Emit(types.ScaffoldProgress{Error: err.Error(), Done: true})
		return err
	}

	total := len(plan.Files)
	for i, f := range plan.Files {
		dest := filepath.Join(rootPath, filepath.FromSlash(f.Path))

		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			msg := fmt.Sprintf("create dir for %s: %v", f.Path, err)
			a.Emit(types.ScaffoldProgress{Error: msg, Done: true})
			return errors.New(msg)
		}

		file, err := os.OpenFile(dest, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err != nil {
			if os.IsExist(err) {
				a.Emit(types.ScaffoldProgress{File: f.Path, Index: i + 1, Total: total})
				continue
			}
			msg := fmt.Sprintf("create file %s: %v", f.Path, err)
			a.Emit(types.ScaffoldProgress{Error: msg, Done: true})
			return errors.New(msg)
		}

		_, writeErr := file.WriteString(f.Content)
		file.Close()
		if writeErr != nil {
			msg := fmt.Sprintf("write file %s: %v", f.Path, writeErr)
			a.Emit(types.ScaffoldProgress{Error: msg, Done: true})
			return errors.New(msg)
		}

		a.Emit(types.ScaffoldProgress{
			File:  f.Path,
			Index: i + 1,
			Total: total,
		})
	}

	a.Emit(types.ScaffoldProgress{Done: true, Total: total, Index: total})
	return nil
}
