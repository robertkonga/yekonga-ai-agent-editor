package agent

import "fmt"

type AgentRequest struct {
	Instructions string
	Tools        string
	Context      string
	History      string
	UserMessage  string
}

func BuildPrompt(req AgentRequest) string {
	return fmt.Sprintf(`
[SYSTEM INSTRUCTIONS]
%s

[AVAILABLE TOOLS]
%s

[PROJECT CONTEXT]
%s

[CONVERSATION HISTORY]
%s

[CURRENT USER REQUEST]
%s
`,
		req.Instructions,
		req.Tools,
		req.Context,
		req.History,
		req.UserMessage,
	)
}
