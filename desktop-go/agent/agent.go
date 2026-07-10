package agent

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"yekonga-builder/agent/provider"
	"yekonga-builder/agent/types"
	"yekonga-builder/console"
	fileTypes "yekonga-builder/types"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:templates
var assets embed.FS

type ApiKeys struct {
	ApiKey          string
	AnthropicApiKey string
	GeminiApiKey    string
	DeepseekApiKey  string
}

type Agent struct {
	ApiKey          string
	AnthropicApiKey string
	GeminiApiKey    string
	DeepseekApiKey  string
	ActivePath      string
	OllamaHost      string
	ctx             *context.Context
	sessionManager  *SessionManager
}

func NewAgent(apiKeys ApiKeys, ollamaHost string, ctx *context.Context) *Agent {
	configDir, _ := os.UserConfigDir()
	storageDir := filepath.Join(configDir, "YekongaEditor", "sessions")

	sm, _ := NewSessionManager(storageDir)

	return &Agent{
		ApiKey:          apiKeys.ApiKey,
		GeminiApiKey:    apiKeys.GeminiApiKey,
		AnthropicApiKey: apiKeys.AnthropicApiKey,
		DeepseekApiKey:  apiKeys.DeepseekApiKey,
		OllamaHost:      ollamaHost,
		ctx:             ctx,
		sessionManager:  sm,
	}
}

func (a *Agent) Emit(p fileTypes.ScaffoldProgress) {
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

func (a *Agent) ListWorkspaceSessions(workspace string) ([]*Session, error) {
	return a.sessionManager.ListWorkspaceSessions(workspace)
}

func (a *Agent) GetSession(id string) (*Session, error) {
	return a.sessionManager.GetSession(id)
}

// AgentChat is bound to Wails — called from Vue when user sends a message.
func (a *Agent) AgentChat(sessionID string, userMessage string, providerName string, modelName string) error {
	// 1. Get or create session
	executor := NewToolExecutor(a)
	session, err := a.sessionManager.GetSession(sessionID)
	if err != nil {
		return err
	}
	if session == nil {
		session = &Session{
			ID:          sessionID,
			Provider:    providerName,
			History:     []types.ChatMessage{},
			LastUpdated: time.Now(),
			Workspace:   a.ActivePath,
		}
	}

	console.Log(providerName, modelName, userMessage)

	// 2. Select selectedProvider
	var selectedProvider provider.LLMProvider
	switch providerName {
	case "anthropic":
		selectedProvider = provider.NewAnthropicProvider(a.AnthropicApiKey, modelName)
	case "gemini":
		selectedProvider = provider.NewGeminiProvider(a.GeminiApiKey, modelName)
	case "deepseek":
		selectedProvider = provider.NewDeepSeekProvider(a.DeepseekApiKey, modelName)
	default:
		selectedProvider = provider.NewOllamaProviderWithURL(modelName, a.OllamaHost)
	}

	console.Log("create.selectedProvider", selectedProvider)
	// 3. Append user message
	session.History = append(session.History, types.ChatMessage{
		Role:    "user",
		Content: userMessage,
	})

	// 4. Agentic loop
	for {
		resp, err := selectedProvider.Complete(
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
			console.Log("agent:message", resp.Content)
			console.Log("agent:message", resp.ToolCalls)
			runtime.EventsEmit(*a.ctx, "agent:message", resp.Content)
		}

		// Add assistant response to history
		assistantMsg := types.ChatMessage{
			Role: "assistant",
		}

		// if len(resp.ToolCalls) > 0 {
		// 	// assistantMsg.Content = resp.ToolCalls
		// }
		if len(resp.ToolCalls) > 0 {
			var blocks []types.ContentBlock
			if resp.Content != "" {
				blocks = append(blocks, types.ContentBlock{
					Type:    "text",
					Content: resp.Content,
				})
			}
			blocks = append(blocks, resp.ToolCalls...)
			assistantMsg.Content = blocks
		} else {
			assistantMsg.Content = resp.Content
		}
		session.History = append(session.History, assistantMsg)

		// No tool calls → done
		if len(resp.ToolCalls) == 0 {
			a.sessionManager.SaveSession(session)
			runtime.EventsEmit(*a.ctx, "agent:done", nil)
			return nil
		}

		// Execute tools
		var toolResults []types.ContentBlock
		for _, call := range resp.ToolCalls {
			runtime.EventsEmit(*a.ctx, "agent:tool", map[string]string{
				"name":  call.Name,
				"input": string(call.Input),
			})

			result, toolErr := executor.Execute(call.Name, call.Input)
			if toolErr != nil {
				result = fmt.Sprintf("error: %s", toolErr.Error())
			}

			toolResults = append(toolResults, types.ContentBlock{
				Type: "tool_result",
				Name: call.Name,
				ID:   call.ID,
				// ToolUseID: call.ID,
				Content: result,
			})
		}

		// Add results to history and continue loop
		session.History = append(session.History, types.ChatMessage{
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
