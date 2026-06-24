package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"yekonga-builder/agent/provider"
	"yekonga-builder/agent/types"
)

func (a *Agent) GenerateSchema(description string) (string, error) {
	schemaSystemPrompt := a.getSystemInstruction("database") // Template is named database.md
	provider := provider.NewGeminiProvider(a.ApiKey, "")

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	resp, err := provider.Complete(ctx, schemaSystemPrompt, []types.ChatMessage{{Role: "user", Content: description}}, nil)
	if err != nil {
		return "", err
	}

	// Validate it's parseable before returning to frontend
	var schema map[string]any
	if err := json.Unmarshal([]byte(resp.Content), &schema); err != nil {
		return "", fmt.Errorf("LLM returned invalid JSON: %w", err)
	}

	return resp.Content, nil
}
