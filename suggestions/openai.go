package suggestions

import (
	"context"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

const SYSTEM_PROMPT = `
I want you to act as a senior software engineer. 
I will provide you with code and I want you to improve it, 
and it will be your job to come up with best practices and 
patterns to make this a robust application. 
DO NOT include any explanations in your responses.
Simply improve the code and make it production ready,
without explaining why you are doing what you are doing.
`

func openAiClient() *openai.Client {
	return openai.NewClient(os.Getenv("OPEN_AI_TOKEN"))
}

func suggestImprovements(client *openai.Client, fileContent string) (openai.ChatCompletionResponse, error) {
	return client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: os.Getenv("OPEN_AI_MODEL"),
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: SYSTEM_PROMPT,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fileContent,
				},
			},
		},
	)
}

func getAvailableModels(ctx context.Context, client *openai.Client) (openai.ModelsList, error) {
	return client.ListModels(ctx)
}

func displayAvailableModels(ctx context.Context, client *openai.Client, err error) {
	if err.(*openai.APIError).Type == "invalid_request_error" {
		models, err := getAvailableModels(ctx, client)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("===== AVAILABLE MODELS =====")
		for _, m := range models.Models {
			log.Println(m.ID)
		}
		log.Println("============================")
	}
}
