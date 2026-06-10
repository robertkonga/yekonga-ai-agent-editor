package agent

import (
	"context"
	"encoding/json"
	"fmt"
)

const schemaSystemPrompt = ``

func (a *App) GenerateSchema(description string) (string, error) {
	client := newLLMClient(a.config.AnthropicAPIKey)

	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()

	raw, err := client.complete(ctx, schemaSystemPrompt, description)
	if err != nil {
		return "", err
	}

	// Validate it's parseable before returning to frontend
	var schema map[string]any
	if err := json.Unmarshal([]byte(raw), &schema); err != nil {
		return "", fmt.Errorf("LLM returned invalid JSON: %w", err)
	}

	return raw, nil
}
