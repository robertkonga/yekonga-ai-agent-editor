package agent

import (
	"context"
	"log"
	"os"

	"google.golang.org/genai"
)

func ClientAgent(contents string) {
	ctx := context.Background()

	instructionBytes, err := os.ReadFile("instruction.md")
	if err != nil {
		log.Fatal(err)
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Models.GenerateContent(
		ctx,
		"gemini-3.5-flash",
		genai.Text(contents),
		&genai.GenerateContentConfig{
			Temperature:      0.1,
			ResponseMIMEType: "application/json", // Crucial for reliable parsing
			SystemInstruction: &genai.Content{
				Parts: []*genai.Part{
					genai.NewPartFromText(string(instructionBytes)),
				},
			},
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp.Text())
}
