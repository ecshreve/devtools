package gencom

import (
	"context"
	"fmt"
	"os"

	_ "embed"

	"github.com/charmbracelet/log"
	"github.com/sashabaranov/go-openai"
)

//go:embed templates/system_prompt.tmpl
var systemPrompt string

//go:embed templates/user_prompt.tmpl
var userPrompt string

type OpenAIInterface interface {
	GenerateCommitMessage(diff string) (string, error)
}

type OpenAIClient struct {
	Client openai.Client
}

func NewOpenAIClient() OpenAIInterface {
	return OpenAIClient{
		Client: *openai.NewClient(os.Getenv("OPENAI_API_KEY")),
	}
}

func (o OpenAIClient) GenerateCommitMessage(diff string) (string, error) {
	log.Info("OpenAIClient.GenerateCommitMessage")

	fullUserPrompt := fmt.Sprintf("%s\n%s", userPrompt, diff)
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT4TurboPreview,
		Seed:      nil,
		MaxTokens: 2400,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fullUserPrompt,
			},
		},
	}

	resp, err := o.Client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", err
	}
	log.Debug("OpenAIClient.GenerateCommitMessage", "resp", resp)
	log.Debug("OpenAIClient.GenerateCommitMessage", "resp.Choices[0].Message.Content", resp.Choices[0].Message.Content)

	return resp.Choices[0].Message.Content, nil
}

type MockOpenAIClient struct{}

func (mo MockOpenAIClient) GenerateCommitMessage(diff string) (string, error) {
	log.Info("MockOpenAIClient.GenerateCommitMessage")
	return "mocked commit message", nil
}
