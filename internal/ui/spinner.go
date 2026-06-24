package ui

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"
)

// Spinner wraps a Bubble Tea spinner for long-running brew commands
// (update/upgrade). Log lines stream above the spinner; the label can be
// updated in flight. On non-TTY stdout it degrades to plain printing so
// piped output stays clean.
type Spinner struct {
	program *tea.Program
	done    chan struct{}
	tty     bool
}

type (
	labelMsg string
	logMsg   string
	quitMsg  struct{}
)

type spinnerModel struct {
	spinner spinner.Model
	label   string
}

func newSpinnerModel(label string) spinnerModel {
	s := spinner.New(
		spinner.WithSpinner(spinner.MiniDot),
		spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#FBB040"))),
	)
	return spinnerModel{spinner: s, label: label}
}

func (m spinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case labelMsg:
		m.label = string(msg)
		return m, nil
	case logMsg:
		return m, tea.Println(string(msg))
	case quitMsg:
		return m, tea.Quit
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m spinnerModel) View() tea.View {
	return tea.NewView(fmt.Sprintf("%s %s", m.spinner.View(), m.label))
}

func NewSpinner(label string) *Spinner {
	tty := term.IsTerminal(os.Stdout.Fd())
	s := &Spinner{
		done: make(chan struct{}),
		tty:  tty,
	}
	if tty {
		s.program = tea.NewProgram(newSpinnerModel(label))
	} else {
		fmt.Println(label)
	}
	return s
}

func (s *Spinner) Start() {
	if !s.tty {
		close(s.done)
		return
	}
	go func() {
		_, _ = s.program.Run()
		close(s.done)
	}()
}

func (s *Spinner) Stop() {
	if s.tty {
		s.program.Send(quitMsg{})
	}
	<-s.done
}

func (s *Spinner) Log(line string) {
	if !s.tty {
		fmt.Println(line)
		return
	}
	s.program.Send(logMsg(line))
}

func (s *Spinner) SetLabel(label string) {
	if !s.tty {
		return
	}
	s.program.Send(labelMsg(label))
}
