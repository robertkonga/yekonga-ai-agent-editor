package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"yekonga-builder/types"
)

// ── Constants ────────────────────────────────────────────────────────────────

const (
	anthropicAPI     = "https://api.anthropic.com/v1/messages"
	anthropicVersion = "2023-06-01"
	llmModel         = "claude-sonnet-4-20250514"
	llmMaxTokens     = 8096
	httpTimeout      = 120 * time.Second
)

const scaffoldSystemPrompt = `You are a project scaffolding agent embedded in a code editor.
When given a project description, respond ONLY with a valid JSON object — no markdown fences,
no explanation, no preamble. The JSON must match this exact schema:

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

// ── Internal Anthropic types ──────────────────────────────────────────────────

type llmRequest struct {
	Model     string       `json:"model"`
	MaxTokens int          `json:"max_tokens"`
	System    string       `json:"system"`
	Messages  []llmMessage `json:"messages"`
}

type llmMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type llmResponse struct {
	Content []llmContentBlock `json:"content"`
	Error   *llmError         `json:"error,omitempty"`
	// Usage is included for logging / future token-budget enforcement
	Usage llmUsage `json:"usage"`
}

type llmContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type llmError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type llmUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// ── LLM client ───────────────────────────────────────────────────────────────

type llmClient struct {
	apiKey     string
	httpClient *http.Client
}

func newLLMClient(apiKey string) *llmClient {
	return &llmClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: httpTimeout,
		},
	}
}

// complete sends a single request to the Anthropic API and returns the raw text response.
func (c *llmClient) complete(ctx context.Context, system, userPrompt string) (string, error) {
	reqBody := llmRequest{
		Model:     llmModel,
		MaxTokens: llmMaxTokens,
		System:    system,
		Messages: []llmMessage{
			{Role: "user", Content: userPrompt},
		},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, anthropicAPI, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", anthropicVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", fmt.Errorf("LLM request timed out after %s", httpTimeout)
		}
		return "", fmt.Errorf("HTTP request: %w", err)
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response body: %w", err)
	}

	// Non-2xx: surface the API error message
	if resp.StatusCode != http.StatusOK {
		var apiErr llmResponse
		if jsonErr := json.Unmarshal(rawBody, &apiErr); jsonErr == nil && apiErr.Error != nil {
			return "", fmt.Errorf("Anthropic API error %d: %s — %s",
				resp.StatusCode, apiErr.Error.Type, apiErr.Error.Message)
		}
		return "", fmt.Errorf("Anthropic API error %d: %s", resp.StatusCode, string(rawBody))
	}

	var llmResp llmResponse
	if err := json.Unmarshal(rawBody, &llmResp); err != nil {
		return "", fmt.Errorf("parse response JSON: %w", err)
	}
	if len(llmResp.Content) == 0 {
		return "", errors.New("LLM returned empty content")
	}

	// Collect all text blocks (there is usually just one)
	var sb strings.Builder
	for _, block := range llmResp.Content {
		if block.Type == "text" {
			sb.WriteString(block.Text)
		}
	}
	return sb.String(), nil
}

// ── Plan parsing ─────────────────────────────────────────────────────────────

// parsePlan extracts and validates the types.ScaffoldPlan from the raw LLM text.
// The model occasionally wraps the JSON in markdown fences — we strip those.
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

// ── Wails-facing App method ───────────────────────────────────────────────────

// GenerateProject is the Wails-bound method called from the Vue frontend.
// It calls the LLM, parses the scaffold plan, writes every file to rootPath,
// and emits real-time progress events.
func (a *Agent) GenerateProject(userPrompt string, rootPath string, extraConventions string) error {
	client := newLLMClient(a.ApiKey)

	// Optionally inject project-specific conventions on top of the base prompt
	system := scaffoldSystemPrompt
	if strings.TrimSpace(extraConventions) != "" {
		system += "\n\nAdditional conventions for this project:\n" + extraConventions
	}

	a.Emit(types.ScaffoldProgress{File: "Contacting LLM…"})

	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()

	raw, err := client.complete(ctx, system, userPrompt)
	if err != nil {
		a.Emit(types.ScaffoldProgress{Error: err.Error(), Done: true})
		return err
	}

	a.Emit(types.ScaffoldProgress{File: "Parsing plan…"})

	plan, err := parsePlan(raw)
	if err != nil {
		a.Emit(types.ScaffoldProgress{Error: err.Error(), Done: true})
		return err
	}

	total := len(plan.Files)

	for i, f := range plan.Files {
		dest := filepath.Join(rootPath, filepath.FromSlash(f.Path))

		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			msg := fmt.Sprintf("create dir for %s: %v", f.Path, err)
			a.Emit(types.ScaffoldProgress{Error: msg, Done: true})
			return errors.New(msg)
		}

		// Write file — fail if it already exists to avoid silent overwrites
		file, err := os.OpenFile(dest, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err != nil {
			if os.IsExist(err) {
				msg := fmt.Sprintf("file already exists, skipping: %s", f.Path)
				a.Emit(types.ScaffoldProgress{File: f.Path, Index: i + 1, Total: total})
				// Non-fatal: log and continue
				fmt.Println("[scaffold]", msg)
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
