package gencom

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/log"
)

type Worker struct {
	git    GitInterface
	openai OpenAIInterface
}

func NewWorker() *Worker {
	checkRequiredCommands([]string{"git", "less"})
	checkOpenAIKey()

	var gc GitInterface
	if os.Getenv("MOCK_GIT_DIFF") == "true" {
		gc = MockGitCommand{}
	} else {
		gc = NewGitCommand()
	}

	var oc OpenAIInterface
	if os.Getenv("MOCK_OPENAI") == "true" {
		oc = MockOpenAIClient{}
	} else {
		oc = NewOpenAIClient()
	}

	return &Worker{
		git:    gc,
		openai: oc,
	}
}

func (w *Worker) Run() *Commit {
	log.Info("Worker.Run")
	// Get the git diff
	diff, err := w.git.GetDiff()
	if err != nil {
		log.Error("Error getting git diff", "err", err)
		os.Exit(1)
	}
	log.Debug("before processing", "len", len(diff))

	proc := summarizeDiff(diff)
	log.Debug("after processing", "len", len(proc))

	// Generate the commit message
	commitMessage, err := w.openai.GenerateCommitMessage(proc)
	if err != nil {
		log.Error("Error generating commit message", "err", err)
		os.Exit(1)
	}

	var cmt Commit
	err = json.Unmarshal([]byte(commitMessage), &cmt)
	if err != nil {
		log.Error("Error unmarshalling commit", "err", err)
		return nil
	}

	cmt.Body = foldString(cmt.Body, 72)
	return &cmt
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
