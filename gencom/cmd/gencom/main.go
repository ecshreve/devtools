package main

import (
	"fmt"
	"gencom"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

const maxWidth = 160

var (
	red               = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo            = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green             = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
	shortCircuit      = false
	useCommit         = false
	CommitBuilderChan = make(chan gencom.Commit)
	CommitReadyChan   = make(chan struct{})
	DoneChan          = make(chan struct{})
	toCommit          *gencom.Commit
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

type Model struct {
	lg     *lipgloss.Renderer
	styles *Styles
	form   *huh.Form
	comm   *gencom.Commit
	width  int
}

func NewModel() Model {
	log.Info("NewModel")
	m := Model{width: maxWidth}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)

	newCommit := toCommit
	m.comm = newCommit
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Skip editing?").
				Value(&shortCircuit),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Type").
				Inline(true).
				Value(&newCommit.Type).
				Suggestions([]string{"feat", "fix", "docs", "test", "refactor"}),
			huh.NewInput().
				Title("Scope").
				CharLimit(16).
				Inline(true).
				Value(&newCommit.Scope),
			huh.NewInput().
				Title("Desc").
				Inline(true).
				Value(&newCommit.Desc).
				Validate(func(s string) error {
					if len(s) < 10 {
						return fmt.Errorf("summary must be at least 10 characters")
					}

					if len(s) > 48-len(newCommit.Type)-len(newCommit.Scope) {
						return fmt.Errorf("summary line must be less than 50 characters")
					}

					return nil
				}),
			huh.NewText().
				Value(&newCommit.Body).
				Title("Body").
				Lines(4).
				Validate(func(s string) error {
					for _, l := range strings.Split(s, "\n") {
						if len(l) > 72 {
							return fmt.Errorf("body line length must be less than 72 characters")
						}
					}
					return nil
				}),
		).WithHideFunc(func() bool { return shortCircuit }),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Ready to commit").
				Value(&useCommit),
		),
	).WithWidth(80).WithShowErrors(false).WithShowHelp(false)

	return m
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
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
		// Commit the changes
		log.Debug("committing", "commit", m.comm, "useCommit", useCommit)

		return s.Status.Copy().Margin(0, 1).Padding(1, 2).Width(80).Render(m.comm.String()) + "\n\n" + "---"

	default:
		// Form (left side)
		v := strings.TrimSuffix(m.form.View(), "\n\n")
		form := m.lg.NewStyle().Margin(1, 0).Render(v)

		// Status (right side)
		var status string
		{

			const statusWidth = 74
			statusMarginLeft := m.width - statusWidth - lipgloss.Width(form) - s.Status.GetMarginRight()
			status = s.Status.Copy().
				Height(10).
				Width(statusWidth).
				MarginLeft(statusMarginLeft).
				Render(s.StatusHeader.Render("|------------------------------------------- 50 >|" + "\n\n" +
					m.comm.String() + "\n\n" +
					"|----------------------------------------------------------------- 72 >|"))
		}

		errors := m.form.Errors()
		header := m.appBoundaryView("Commit Message Helper")
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

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)

	if os.Getenv("GENCOM_ENV") == "dev" {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			log.Fatal("Error opening log file", "err", err)
		}
		log.SetOutput(f)
		defer f.Close()
	}

	wrk := gencom.NewWorker()
	log.Print("Starting worker")
	err := wrk.Run()
	if err != nil {
		log.Fatal("Error running worker", "err", err)
	}
	log.Print("Worker finished", "commit", wrk.CommitData)

	toCommit = wrk.CommitData
	_, err = tea.NewProgram(NewModel()).Run()
	if err != nil {
		log.Fatal("Error running program", "err", err)
	}

	if useCommit {
		log.Info("Running git commit")
		_, err := gencom.Execute(toCommit)
		if err != nil {
			log.Error("Error running git commit", "err", err)
		}
	} else {
		log.Info("Skipping commit", "useCommit", useCommit, "shortCircuit", shortCircuit)
	}

	log.Info("Done")
}
