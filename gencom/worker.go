package gencom

import (
	"encoding/json"
	"os"

	"github.com/charmbracelet/log"
)

// Worker is a struct that describes the worker that
// generates the commit message.
type Worker struct {
	git    GitClient
	openai OpenAIInterface

	CommitData *Commit
}

// NewWorker creates a new instance of Worker.
func NewWorker() *Worker {
	checkRequiredCommands([]string{"git", "less"})
	checkOpenAIKey()

	var gc GitClient
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

// Run runs the worker.
func (w *Worker) Run() error {
	log.Info("Worker.Run")

	// Get the git diff
	diff, err := w.git.GetDiff()
	if err != nil {
		log.Error("Error getting git diff", "err", err)
		return err
	}

	// Process the diff
	log.Debug("before processing", "len", len(diff))
	proc := summarizeDiff(diff)
	log.Debug("after processing", "len", len(proc))

	// Generate the commit message
	commitMessage, err := w.openai.GenerateCommitMessage(proc)
	if err != nil {
		log.Error("Error generating commit message", "err", err, "commitMessage", commitMessage)
		return err
	}

	// Unmarshal the commit message
	var cmt Commit
	err = json.Unmarshal([]byte(commitMessage), &cmt)
	if err != nil {
		log.Error("Error unmarshalling commit", "err", err)
		return err
	}

	// Post-process the commit message
	cmt.Body = foldString(cmt.Body, 72)
	w.CommitData = &cmt
	return nil
}
