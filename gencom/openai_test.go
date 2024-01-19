package main

import "testing"

func TestGenerateCommitMessage(t *testing.T) {
	openai := MockOpenAIClient{}
	diff := "mocked git diff output"
	commitMessage, err := openai.GenerateCommitMessage(diff)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	expectedMessage := "mocked commit message"
	if commitMessage != expectedMessage {
		t.Errorf("Expected %s, got %s", expectedMessage, commitMessage)
	}
}
