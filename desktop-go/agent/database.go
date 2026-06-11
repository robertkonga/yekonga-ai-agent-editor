package agent

import (
	"context"
	"encoding/json"
	"fmt"
)

func (a *Agent) GenerateSchema(description string) (string, error) {
	schemaSystemPrompt := a.getSystemInstruction("instruction")
	client := newLLMClient(a.ApiKey)

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
