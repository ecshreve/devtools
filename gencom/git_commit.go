package gencom

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/log"
)

// Commit describes a git commit.
type Commit struct {
	Type   string `json:"type"`
	Scope  string `json:"scope"`
	Desc   string `json:"desc"`
	Body   string `json:"body"`
	Footer string `json:"footer"`
}

func (c Commit) String() string {
	out := c.Type
	if c.Scope != "" {
		out = fmt.Sprintf("%s(%s)", out, c.Scope)
	}
	out += ": " + c.Desc

	if c.Body != "" {
		out += "\n\n" + c.Body
	}

	if c.Footer != "" {
		out += "\n\n" + c.Footer
	}

	return out
}

func GenComStruct() Commit {
	return Commit{
		Type:  "feat",
		Scope: "gencom",
		Desc:  "Add bubbletea interface",
		Body: `Added bubbletea interface to gencom. This will allow for a more 
interactive experience when generating commit messages.

- Added bubbletea interface
- Added a new state machine to handle the bubbletea interface
- Added form to handle user input`,
		Footer: "",
	}
}

// Execute runs the git commit command.
func Execute(c *Commit) (string, error) {
	log.Info("Execute")
	msg := c.String()
	args := []string{"commit", "-m", msg}

	cmd := exec.Command("git", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Debug("Execute", "cmd", fmt.Sprintf("git %v", args))
	return fmt.Sprintf("git %v", args), cmd.Run()
}
