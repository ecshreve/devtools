package main

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/sashabaranov/go-openai"
)

type OpenAIInterface interface {
	GenerateCommitMessage(diff string) (string, error)
}

type OpenAIClient struct {
	Client openai.Client
}

func (o OpenAIClient) GenerateCommitMessage(diff string) (string, error) {
	log.Info("OpenAIClient.GenerateCommitMessage")
	const promptTemplate = `Analyze the following git diff of a codebase and generate
a concise informative commit message. Focus on the intention  behind the changes and 
the impact on the project. Generate the message in JSON format with the following details:
- DESC: A concise description containing 40 characters or less.
- BODY: A detailed explanation of the changes suitable for the body of a git commit message.
- PRETTY: A pretty-printed version of the BODY, formatted as a bulleted list.
- TYPE: Classify this change as one of the [fix, feat, test, docs, refactor, chore].
- SCOPE: Provide one word describing the area of the codebase most affected.

Git Diff:
%s

Format the response as a JSON dictionary:
{
  "DESC": "<description>",
  "BODY": "<body>",
  "TYPE": "<type>",
  "SCOPE": "<scope>",
	"PRETTY": "<pretty>",
}
`
	p := fmt.Sprintf(promptTemplate, diff)

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		Seed:      nil,
		MaxTokens: 1024,
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

	return resp.Choices[0].Message.Content, nil
}

type MockOpenAIClient struct{}

func (mo MockOpenAIClient) GenerateCommitMessage(diff string) (string, error) {
	log.Info("MockOpenAIClient.GenerateCommitMessage")
	return "mocked commit message", nil
}
