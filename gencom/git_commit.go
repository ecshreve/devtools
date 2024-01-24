package gencom

import "fmt"

type Commit struct {
	Type             string
	Scope            string
	Desc             string
	Body             string
	Footer           string
	DoesWantToCommit bool
}

// Message returns the commit message.
func (c Commit) MessageString() *string {
	out := c.Type
	if c.Scope != "" {
		out = fmt.Sprintf("%s(%s)", out, c.Scope)
	}
	out += ": " + c.Desc

	return &out
}

func (c Commit) String() string {
	out := *c.MessageString()

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
