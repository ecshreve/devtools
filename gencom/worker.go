package gencom

import (
	"encoding/json"
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
		return nil
	}
	log.Debug("before processing", "len", len(diff))

	proc := summarizeDiff(diff)
	log.Debug("after processing", "len", len(proc))

	// Generate the commit message
	commitMessage, err := w.openai.GenerateCommitMessage(proc)
	if err != nil {
		log.Error("Error generating commit message", "err", err)
		return nil
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
			log.Fatal("command is not installed.", "cmd", cmd)
		}
	}
}

func checkOpenAIKey() {
	if os.Getenv("OPENAI_API_KEY") == "" {
		log.Fatal("OPENAI_API_KEY environment variable is not set.")
	}
}
