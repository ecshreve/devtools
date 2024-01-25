package gencom

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/log"
)

type Commit struct {
	Type   string
	Scope  string
	Desc   string
	Body   string
	Footer string
}

func (c Commit) Parts() (string, string, string) {
	out := c.Type
	if c.Scope != "" {
		out = fmt.Sprintf("%s(%s)", out, c.Scope)
	}
	out += ": " + c.Desc

	return out, c.Body, c.Footer
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

func Execute(c *Commit) (string, error) {
	log.Info("Execute")
	msg, body, _ := c.Parts()
	args := []string{"commit", "-m", msg}
	if body != "" {
		args = append(args, "-m", body)
	}

	cmd := exec.Command("git", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Debug("Execute", "cmd", fmt.Sprintf("git %v", args))
	return fmt.Sprintf("git %v", args), cmd.Run()
}
