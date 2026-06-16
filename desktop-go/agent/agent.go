package agent

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"yekonga-builder/console"
	"yekonga-builder/types"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:templates
var assets embed.FS

type Agent struct {
	ApiKey         string
	ActivePath     string
	ctx            *context.Context
	sessionManager *SessionManager
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
	configDir, _ := os.UserConfigDir()
	storageDir := filepath.Join(configDir, "YekongaEditor", "sessions")

	sm, _ := NewSessionManager(storageDir)

	return &Agent{
		ApiKey:         apiKey,
		ctx:            ctx,
		sessionManager: sm,
	}
}

func (a *Agent) Emit(p types.ScaffoldProgress) {
	runtime.EventsEmit(*a.ctx, "scaffold:progress", p)
}

func (a *Agent) getSystemInstruction(name string) string {
	value, err := assets.ReadFile(fmt.Sprintf("templates/%s.%s", name, "md"))

	if err != nil {
		console.Log(err.Error())
	}

	return string(value)
}

func (a *Agent) ListSessions() ([]*Session, error) {
	return a.sessionManager.ListSessions()
}

func (a *Agent) GetSession(id string) (*Session, error) {
	return a.sessionManager.GetSession(id)
}

// AgentChat is bound to Wails — called from Vue when user sends a message.
func (a *Agent) AgentChat(sessionID string, userMessage string, providerName string) error {
	// 1. Get or create session
	session, err := a.sessionManager.GetSession(sessionID)
	if err != nil {
		return err
	}
	if session == nil {
		session = &Session{
			ID:          sessionID,
			Provider:    providerName,
			History:     []ChatMessage{},
			LastUpdated: time.Now(),
		}
	}

	// 2. Select provider
	var provider LLMProvider
	if session.Provider == "anthropic" {
		provider = NewAnthropicProvider(a.ApiKey)
	} else {
		provider = NewGeminiProvider(a.ApiKey)
	}

	// 3. Append user message
	session.History = append(session.History, ChatMessage{
		Role:    "user",
		Content: userMessage,
	})

	// 4. Agentic loop
	for {
		resp, err := provider.Complete(
			context.Background(),
			agentSystemPrompt,
			session.History,
			agentTools,
		)
		if err != nil {
			runtime.EventsEmit(*a.ctx, "agent:error", err.Error())
			return err
		}

		// Stream text to frontend
		if resp.Content != "" {
			runtime.EventsEmit(*a.ctx, "agent:message", resp.Content)
		}

		// Add assistant response to history
		assistantMsg := ChatMessage{
			Role:    "assistant",
			Content: resp.Content,
		}

		if len(resp.ToolCalls) > 0 {
			assistantMsg.Content = resp.ToolCalls
		}
		session.History = append(session.History, assistantMsg)

		// No tool calls → done
		if len(resp.ToolCalls) == 0 {
			a.sessionManager.SaveSession(session)
			runtime.EventsEmit(*a.ctx, "agent:done", nil)
			return nil
		}

		// Execute tools
		var toolResults []ContentBlock
		for _, call := range resp.ToolCalls {
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
				Name:      call.Name,
				ToolUseID: call.ID,
				Content:   result,
			})
		}

		// Add results to history and continue loop
		session.History = append(session.History, ChatMessage{
			Role:    "user",
			Content: toolResults,
		})

		a.sessionManager.SaveSession(session)
	}
}

const agentSystemPrompt = `You are an AI coding assistant embedded in a code editor.
You have tools to read files, list directories, write files, and search code.
Always read relevant files before answering questions about code.
When asked to make changes, read the file first, then write the updated version.`
