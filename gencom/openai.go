package gencom

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/sashabaranov/go-openai"
)

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
	const promptTemplate = `Analyze the following git diff of a codebase and generate
a concise informative commit message. Focus on the intention behind the changes and 
the impact on the project.

Generate output with the following details:
- DESC: A concise summary containing no more than 34 characters.
- BODY: A detailed explanation of the changes suitable for the body of a git commit message.
- TYPE: Classification of this set of changes as one of [fix, feat, test, docs, refactor, chore].
- SCOPE: Single token identifying the area of the codebase most affected. SCOPE should never include spaces or punctuation.

Git Diff:
%s

Format the response as a JSON dictionary:
{
  "desc": "<DESC>",
  "body": "<BODY>",
  "type": "<TYPE>",
  "scope": "<SCOPE>",
}
`
	p := fmt.Sprintf(promptTemplate, diff)

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT40613,
		Seed:      nil,
		MaxTokens: 2400,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: p,
			},
		},
	}

	resp, err := o.Client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", err
	}
	log.Debug("OpenAIClient.GenerateCommitMessage", "resp", resp)

	return resp.Choices[0].Message.Content, nil
}

type MockOpenAIClient struct{}

func (mo MockOpenAIClient) GenerateCommitMessage(diff string) (string, error) {
	log.Info("MockOpenAIClient.GenerateCommitMessage")
	return "mocked commit message", nil
}
