package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alessio/shellescape"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

const maxWidth = 160

var (
	red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
)

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().
		Foreground(indigo).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(indigo).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(green).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.Copy().
		Foreground(red)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

type state int

const (
	statusNormal state = iota
	stateDone
)

type Model struct {
	state  state
	lg     *lipgloss.Renderer
	styles *Styles
	form   *huh.Form
	comm   *Commit
	width  int
}

// commit commits the changes to git
func commit(msg string) (string, error) {
	args := append([]string{"commit", "-m", msg}, os.Args[1:]...)
	cmd := exec.Command("git", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return fmt.Sprintf("git %v", shellescape.QuoteCommand(args)), cmd.Run()
}

func NewModel() Model {
	m := Model{width: maxWidth}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)

	newCommit := GenComStruct()
	m.comm = &newCommit

	shortCircuit := false
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Ready to commit").
				Inline(true).
				Affirmative("Yes!").
				Negative("Nope.").
				Value(&shortCircuit),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Type").
				Options(typeOptions...).
				Value(&newCommit.Type),
			huh.NewInput().
				Title("Scope").
				CharLimit(16).
				Value(&newCommit.Scope),
			huh.NewInput().
				Description(fmt.Sprint(func() int { return 48 - len(newCommit.Type) - len(newCommit.Scope) }())).
				Value(&newCommit.Desc).
				Title("Desc").
				CharLimit(func() int { return 48 - len(newCommit.Type) - len(newCommit.Scope) }()).
				Validate(func(s string) error {
					if len(s) < 10 {
						return fmt.Errorf("summary must be at least 10 characters")
					}

					if len(s) > 48-len(newCommit.Type)-len(newCommit.Scope) {
						return fmt.Errorf("summary must be less than 48 characters")
					}

					return nil
				}),
			huh.NewText().
				Value(&newCommit.Body).
				Title("Body").
				Lines(8).
				Validate(func(s string) error {
					if len(s) < 10 {
						return fmt.Errorf("body must be at least 10 characters")
					}

					for _, l := range strings.Split(s, "\n") {
						if len(l) > 72 {
							return fmt.Errorf("body line length must be less than 72 characters")
						}
					}
					return nil
				}),
			huh.NewConfirm().
				Title("Ready to commit").
				Inline(true).
				Affirmative("Yes!").
				Negative("Nope.").
				Value(&newCommit.DoesWantToCommit),
		),
	).WithWidth(80).WithShowErrors(false).WithShowHelp(false)

	if newCommit.DoesWantToCommit {
		commitMsg := shellescape.Quote(newCommit.String())

		cmd, err := commit(commitMsg)
		if err != nil {
			log.Error("error committing", "err", err, "cmd", cmd)
			os.Exit(1)
		}
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, maxWidth) - m.styles.Base.GetHorizontalFrameSize()
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd

	// Process the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		// Quit when the form is done.
		cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := m.styles

	switch m.form.State {
	case huh.StateCompleted:
		return s.Status.Copy().Margin(0, 1).Padding(1, 2).Width(48).Render(m.form.GetString("desc")) + "\n\n"
	default:

		// Form (left side)
		v := strings.TrimSuffix(m.form.View(), "\n\n")
		form := m.lg.NewStyle().Margin(1, 0).Render(v)

		// Status (right side)
		var status string
		{

			const statusWidth = 80
			statusMarginLeft := m.width - statusWidth - lipgloss.Width(form) - s.Status.GetMarginRight()
			status = s.Status.Copy().
				Height(lipgloss.Height(m.comm.String())).
				Width(statusWidth).
				MarginLeft(statusMarginLeft).
				Render(s.StatusHeader.Render("Current Build") + "\n" +
					"|============================================ 50 >|" + "\n" +
					m.comm.String() + "\n" +
					"|================================================================= 72 >|")

		}

		errors := m.form.Errors()
		header := m.appBoundaryView("Charm Employment Application")
		if len(errors) > 0 {
			header = m.appErrorBoundaryView(m.errorView())
		}
		body := lipgloss.JoinHorizontal(lipgloss.Top, form, status)

		footer := m.appBoundaryView(m.form.Help().ShortHelpView(m.form.KeyBinds()))
		if len(errors) > 0 {
			footer = m.appErrorBoundaryView("")
		}

		return s.Base.Render(header + "\n" + body + "\n\n" + footer)
	}
}

func (m Model) errorView() string {
	var s string
	for _, err := range m.form.Errors() {
		s += err.Error()
	}
	return s
}

func (m Model) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(indigo),
	)
}

func (m Model) appErrorBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.ErrorHeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(red),
	)
}

func main() {
	_, err := tea.NewProgram(NewModel()).Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}
}

var typeOptions = []huh.Option[string]{
	huh.NewOption("feat", "feat"),
	huh.NewOption("fix", "fix"),
	huh.NewOption("docs", "docs"),
	huh.NewOption("test", "test"),
	huh.NewOption("refactor", "refactor"),
}

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
