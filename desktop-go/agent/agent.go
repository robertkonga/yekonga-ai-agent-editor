package agent

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"yekonga-builder/console"
	"yekonga-builder/types"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:templates
var assets embed.FS

type Agent struct {
	ApiKey    string
	ctx       *context.Context
	llmClient geminiClient
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content any    `json:"content"` // string OR []ContentBlock
}

type ContentBlock struct {
	Type      string          `json:"type"`
	Text      string          `json:"text,omitempty"`
	ID        string          `json:"id,omitempty"`          // tool_use
	Name      string          `json:"name,omitempty"`        // tool_use
	Input     json.RawMessage `json:"input,omitempty"`       // tool_use
	ToolUseID string          `json:"tool_use_id,omitempty"` // tool_result
	Content   string          `json:"content,omitempty"`     // tool_result
}

func NewAgent(apiKey string, ctx *context.Context) *Agent {
	return &Agent{
		ApiKey: apiKey,
		ctx:    ctx,
	}
}

func (a *Agent) Emit(p types.ScaffoldProgress) {
	runtime.EventsEmit(*a.ctx, "scaffold:progress", p)
}

func (a *Agent) getSystemInstruction(name string) string {
	value, err := assets.ReadFile(fmt.Sprintf("template/%s.%s", name, "md"))

	if err != nil {
		console.Log(err.Error())
	}

	return string(value)
}

// AgentChat is bound to Wails — called from Vue when user sends a message.
// history is the full conversation so far (sent from frontend).
func (a *Agent) AgentChat(userMessage string, history []ChatMessage) error {
	// Append the new user message to history
	history = append(history, ChatMessage{
		Role:    "user",
		Content: userMessage,
	})

	// Agentic loop — keeps going until Claude stops calling tools
	for {
		// Collect all content blocks from the response
		var textParts []string
		var toolCalls []ContentBlock

		// resp, err := a.llmClient.complete(
		// 	context.Background(),
		// 	agentSystemPrompt,
		// 	history,
		// 	agentTools,
		// )
		// if err != nil {
		// 	runtime.EventsEmit(*a.ctx, "agent:error", err.Error())
		// 	return err
		// }

		// for _, block := range resp.Candidates[0].Content.Parts {
		// 	switch block.Text {
		// 	case "text":
		// 		if strings.TrimSpace(block.Text) != "" {
		// 			textParts = append(textParts, block.Text)
		// 		}
		// 	case "tool_use":
		// 		toolCalls = append(toolCalls, block)
		// 	}
		// }

		// Stream any text to frontend immediately
		if len(textParts) > 0 {
			runtime.EventsEmit(*a.ctx, "agent:message", strings.Join(textParts, ""))
		}

		// No tool calls → Claude is done, exit loop
		if len(toolCalls) == 0 {
			runtime.EventsEmit(*a.ctx, "agent:done", nil)
			return nil
		}

		// Add Claude's response (with tool_use blocks) to history
		history = append(history, ChatMessage{
			Role: "assistant",
			// Content: resp.Content,
		})

		// Execute each tool and collect results
		var toolResults []ContentBlock
		for _, call := range toolCalls {
			// Tell the frontend which tool is running
			runtime.EventsEmit(*a.ctx, "agent:tool", map[string]string{
				"name":  call.Name,
				"input": string(call.Input),
			})

			result, toolErr := a.executeTool(call.Name, call.Input)
			if toolErr != nil {
				result = fmt.Sprintf("error: %s", toolErr.Error())
			}

			toolResults = append(toolResults, ContentBlock{
				Type:      "tool_result",
				ToolUseID: call.ID,
				Content:   result,
			})
		}

		// Add tool results to history and loop again
		history = append(history, ChatMessage{
			Role:    "user",
			Content: toolResults,
		})
	}
}

const agentSystemPrompt = `You are an AI coding assistant embedded in a code editor.
You have tools to read files, list directories, write files, and search code.
Always read relevant files before answering questions about code.
When asked to make changes, read the file first, then write the updated version.`
