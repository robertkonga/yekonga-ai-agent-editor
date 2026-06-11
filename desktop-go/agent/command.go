package agent

import (
	"context"
	"log"

	"google.golang.org/genai"
)

func (a *Agent) ClientAgent(contents string) {
	ctx := context.Background()

	instructionString := a.getSystemInstruction("instruction")

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: a.ApiKey,
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
					genai.NewPartFromText(instructionString),
				},
			},
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp.Text())
}
