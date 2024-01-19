package main

type OpenAIInterface interface {
	GenerateCommitMessage(diff string) (string, error)
}

type OpenAIClient struct{}

func (o OpenAIClient) GenerateCommitMessage(diff string) (string, error) {
	// Real implementation
	return "real commit message", nil
}

type MockOpenAIClient struct{}

func (mo MockOpenAIClient) GenerateCommitMessage(diff string) (string, error) {
	return "mocked commit message", nil
}
