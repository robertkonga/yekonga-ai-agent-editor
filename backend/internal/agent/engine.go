package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/yekonga/ai-agent/internal/llm"
	"github.com/yekonga/ai-agent/internal/sandbox"
)

type Engine struct {
	llmClient   *llm.OllamaClient
	sandbox     *sandbox.SandboxManager
	containerID string
	memory      *Memory
}

func NewEngine(llmClient *llm.OllamaClient, sm *sandbox.SandboxManager, containerID string) *Engine {
	systemPrompt := `You are an autonomous coding agent. 
You can use tools to run commands in a secure Docker sandbox to accomplish tasks.
You must return tool calls strictly in JSON format inside <tool_call> tags if you want to use a tool.
Available tools:
`
	tools := GetAvailableTools()
	toolsJSON, _ := json.MarshalIndent(tools, "", "  ")
	systemPrompt += string(toolsJSON)

	return &Engine{
		llmClient:   llmClient,
		sandbox:     sm,
		containerID: containerID,
		memory:      NewMemory(systemPrompt),
	}
}

// RunLoop executes the Thought -> Action -> Observation loop
func (e *Engine) RunLoop(ctx context.Context, initialTask string, outStream func(string)) error {
	e.memory.AddUserMessage(initialTask)

	for i := 0; i < 15; i++ { // limit iterations to prevent infinite loops
		responseContent := ""
		
		err := e.llmClient.StreamChat(ctx, e.memory.GetMessages(), func(chunk string) error {
			responseContent += chunk
			outStream(chunk)
			return nil
		})

		if err != nil {
			return fmt.Errorf("llm error: %w", err)
		}

		e.memory.AddAssistantMessage(responseContent)

		// Naive parsing of <tool_call>...</tool_call>
		startIdx := strings.Index(responseContent, "<tool_call>")
		endIdx := strings.Index(responseContent, "</tool_call>")

		if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
			toolCallJSON := responseContent[startIdx+11 : endIdx]
			var tc ToolCall
			if err := json.Unmarshal([]byte(toolCallJSON), &tc); err != nil {
				e.memory.AddToolResult("error", "Invalid tool call JSON: "+err.Error())
				continue
			}

			// Execute tool
			result, err := e.executeTool(ctx, tc)
			if err != nil {
				e.memory.AddToolResult(tc.Name, "Failed: "+err.Error())
			} else {
				e.memory.AddToolResult(tc.Name, result)
			}
			outStream("\n--- Tool Execution ---\n" + tc.Name + "\nResult: " + result + "\n----------------------\n")
		} else {
			// No tool call, assume task is done or needs user input
			break
		}
	}
	return nil
}

func (e *Engine) executeTool(ctx context.Context, tc ToolCall) (string, error) {
	log.Printf("Executing tool %s with args %s", tc.Name, string(tc.Arguments))

	switch tc.Name {
	case "run_shell":
		var args struct {
			Command string `json:"command"`
		}
		if err := json.Unmarshal(tc.Arguments, &args); err != nil {
			return "", err
		}
		
		reader, err := e.sandbox.ExecuteCommand(ctx, e.containerID, []string{"sh", "-c", args.Command})
		if err != nil {
			return "", err
		}
		outBytes, err := io.ReadAll(reader)
		// Clean up docker stream multiplexing headers (naive approach for prototype) // In real prod, use stdcopy
		return string(outBytes), err

	case "write_file":
		// Simplified file write via shell command echo / cat
		var args struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		}
		if err := json.Unmarshal(tc.Arguments, &args); err != nil {
			return "", err
		}

		cmd := fmt.Sprintf("cat << 'EOF' > %s\n%s\nEOF", args.Path, args.Content)
		reader, err := e.sandbox.ExecuteCommand(ctx, e.containerID, []string{"sh", "-c", cmd})
		if err != nil {
			return "", err
		}
		outBytes, _ := io.ReadAll(reader)
		return "File created successfully. Output: " + string(outBytes), nil

	default:
		return "", fmt.Errorf("unknown tool: %s", tc.Name)
	}
}
