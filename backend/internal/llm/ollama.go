package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type OllamaClient struct {
	BaseURL    string
	Model      string
	HTTPClient *http.Client
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ChatResponseChunk struct {
	Model     string  `json:"model"`
	Message   Message `json:"message"`
	Done      bool    `json:"done"`
}

func NewOllamaClient(baseURL, model string) *OllamaClient {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	if model == "" {
		model = "llama3" // Default, can be qwen2.5-coder or mistral
	}
	return &OllamaClient{
		BaseURL:    baseURL,
		Model:      model,
		HTTPClient: &http.Client{},
	}
}

func (c *OllamaClient) StreamChat(ctx context.Context, messages []Message, onChunk func(string) error) error {
	reqBody := ChatRequest{
		Model:    c.Model,
		Messages: messages,
		Stream:   true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/api/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama returned status code: %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Bytes()
		var chunk ChatResponseChunk
		if err := json.Unmarshal(line, &chunk); err != nil {
			return fmt.Errorf("failed to unmarshal chunk: %w", err)
		}
		
		if chunk.Message.Content != "" {
			if err := onChunk(chunk.Message.Content); err != nil {
				return err // allow caller to abort streaming
			}
		}

		if chunk.Done {
			break
		}
	}

	return scanner.Err()
}
