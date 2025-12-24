package components

import (
	"github.com/KashifKhn/haft/internal/tui/styles"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type SpinnerModel struct {
	spinner spinner.Model
	message string
	done    bool
	err     error
}

type SpinnerConfig struct {
	Message string
}

func NewSpinner(cfg SpinnerConfig) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.Focused

	return SpinnerModel{
		spinner: s,
		message: cfg.Message,
	}
}

func (m SpinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m SpinnerModel) Update(msg tea.Msg) (SpinnerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case SpinnerDoneMsg:
		m.done = true
		m.err = msg.Err
		return m, nil
	}

	return m, nil
}

func (m SpinnerModel) View() string {
	if m.done {
		if m.err != nil {
			return styles.CrossMark + " " + styles.RenderError(m.err.Error())
		}
		return styles.CheckMark + " " + styles.RenderSuccess(m.message)
	}
	return m.spinner.View() + " " + m.message
}

func (m SpinnerModel) Done() bool {
	return m.done
}

func (m SpinnerModel) Err() error {
	return m.err
}

func (m *SpinnerModel) SetMessage(msg string) {
	m.message = msg
}

type SpinnerDoneMsg struct {
	Err error
}

func SpinnerComplete() tea.Msg {
	return SpinnerDoneMsg{Err: nil}
}

func SpinnerFailed(err error) tea.Msg {
	return SpinnerDoneMsg{Err: err}
}
