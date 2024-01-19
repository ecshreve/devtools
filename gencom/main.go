package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/log"
	"github.com/kr/pretty"
	"github.com/sashabaranov/go-openai"
)

type Worker struct {
	git    GitInterface
	openai OpenAIInterface
}

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	checkRequiredCommands([]string{"git", "less"})
	checkOpenAIKey()

	// Get the git diff
	diff, _ := MockGitCommand{}.GetDiff()
	log.Debug("before processing", "len", len(diff))

	proc := summarizeDiff(diff)
	log.Debug("after processing", "len", len(proc))
	pretty.Println(proc)
	// wrk := Worker{
	// 	git: MockGitCommand{},
	// 	openai: OpenAIClient{
	// 		Client: *openai.NewClient(os.Getenv("OPENAI_API_KEY")),
	// 	},
	// }

	// diff, err := wrk.git.GetDiff()
	// if err != nil {
	// 	log.Error("Error getting git diff:", err)
	// 	os.Exit(1)
	// }
	// log.Debug(strings.Join(strings.Split(diff[:100], "\n"), ""))

	// commitMessage, err := wrk.openai.GenerateCommitMessage(diff)
	// if err != nil {
	// 	log.Error(err)
	// 	os.Exit(1)
	// }
	// log.Debug(commitMessage)

}

func checkRequiredCommands(cmds []string) {
	// Check if cmds are installed
	for _, cmd := range cmds {
		_, err := exec.LookPath(cmd)
		if err != nil {
			log.Error("command is not installed.", "cmd", cmd)
			os.Exit(1)
		}
	}
}

func checkOpenAIKey() {
	if os.Getenv("OPENAI_API_KEY") == "" {
		fmt.Println("OPENAI_API_KEY environment variable is not set.")
		os.Exit(1)
	}
}

func getGitDiff() string {
	cmd := exec.Command("git", "diff", "--cached", "--unified=0")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getting git diff:", err)
		os.Exit(1)
	}
	return string(output)
}

func generateCommitMessage(diff string) string {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	// Prepare your prompt or use the diff directly
	prompt := "Your prompt based on the diff: " + diff

	// Call the OpenAI API (you might need to adjust parameters according to your needs)
	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo, // or your preferred model
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: prompt,
			},
		},
		MaxTokens: 100,
	})

	if err != nil {
		fmt.Println("Error calling OpenAI API:", err)
		os.Exit(1)
	}

	return resp.Choices[0].Message.Content
}
