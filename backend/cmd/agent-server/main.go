package main

import (
	"context"
	"log"
	
	"github.com/yekonga/ai-agent/internal/sandbox"
	"github.com/yekonga/ai-agent/internal/llm"
	"github.com/yekonga/ai-agent/internal/agent"
	"github.com/yekonga/ai-agent/internal/api"
)

func main() {
	log.Println("Starting Agent Engine Server...")
	ctx := context.Background()

	// Initialize dependencies
	sandboxManager, err := sandbox.NewSandboxManager()
	if err != nil {
		log.Fatalf("Failed to initialize docker sandbox: %v", err)
	}

	// Spin up a workspace container (using base ubuntu for now)
	containerID, err := sandboxManager.StartSandbox(ctx, "ubuntu:22.04", "/tmp/workspace") // Tmp dir for prototype
	if err != nil {
		log.Fatalf("Failed to start docker sandbox: %v", err)
	}
	defer sandboxManager.StopSandbox(ctx, containerID)

	llmClient := llm.NewOllamaClient("", "")
	agentEngine := agent.NewEngine(llmClient, sandboxManager, containerID)

	server := api.NewServer(":8080", agentEngine)
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
