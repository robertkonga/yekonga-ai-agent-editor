package agent

import "github.com/yekonga/ai-agent/internal/llm"

type Memory struct {
	systemPrompt string
	messages     []llm.Message
	maxTokens    int
}

func NewMemory(systemPrompt string) *Memory {
	return &Memory{
		systemPrompt: systemPrompt,
		messages: []llm.Message{
			{Role: "system", Content: systemPrompt},
		},
		maxTokens: 8192,
	}
}

func (m *Memory) AddUserMessage(content string) {
	m.messages = append(m.messages, llm.Message{Role: "user", Content: content})
}

func (m *Memory) AddAssistantMessage(content string) {
	m.messages = append(m.messages, llm.Message{Role: "assistant", Content: content})
}

func (m *Memory) AddToolResult(toolName, result string) {
	m.messages = append(m.messages, llm.Message{
		Role: "user", 
		Content: "Tool Result for " + toolName + ":\n" + result,
	})
}

func (m *Memory) GetMessages() []llm.Message {
	return m.messages
}

func (m *Memory) Clear() {
	m.messages = []llm.Message{
		{Role: "system", Content: m.systemPrompt},
	}
}
