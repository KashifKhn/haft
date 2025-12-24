package components

import (
	"github.com/KashifKhn/haft/internal/tui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TextInputModel struct {
	textInput   textinput.Model
	label       string
	placeholder string
	required    bool
	validator   func(string) error
	err         error
	submitted   bool
}

type TextInputConfig struct {
	Label       string
	Placeholder string
	Required    bool
	CharLimit   int
	Width       int
	Validator   func(string) error
}

func NewTextInput(cfg TextInputConfig) TextInputModel {
	ti := textinput.New()
	ti.Placeholder = cfg.Placeholder
	ti.Focus()
	ti.CharLimit = cfg.CharLimit
	if ti.CharLimit == 0 {
		ti.CharLimit = 256
	}
	ti.Width = cfg.Width
	if ti.Width == 0 {
		ti.Width = 40
	}
	ti.PromptStyle = styles.Focused
	ti.TextStyle = styles.Focused
	ti.Cursor.Style = styles.Cursor

	return TextInputModel{
		textInput:   ti,
		label:       cfg.Label,
		placeholder: cfg.Placeholder,
		required:    cfg.Required,
		validator:   cfg.Validator,
	}
}

func (m TextInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TextInputModel) Update(msg tea.Msg) (TextInputModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if err := m.Validate(); err != nil {
				m.err = err
				return m, nil
			}
			m.submitted = true
			return m, nil
		case tea.KeyEsc:
			m.submitted = true
			return m, nil
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m TextInputModel) View() string {
	var view string

	if m.label != "" {
		view = styles.RenderTitle(m.label) + "\n"
	}

	view += m.textInput.View()

	if m.err != nil {
		view += "\n" + styles.RenderError(m.err.Error())
	}

	return view
}

func (m TextInputModel) Value() string {
	return m.textInput.Value()
}

func (m TextInputModel) Validate() error {
	value := m.textInput.Value()

	if m.required && value == "" {
		return errRequired
	}

	if m.validator != nil {
		return m.validator(value)
	}

	return nil
}

func (m TextInputModel) Submitted() bool {
	return m.submitted
}

func (m *TextInputModel) Focus() tea.Cmd {
	return m.textInput.Focus()
}

func (m *TextInputModel) Blur() {
	m.textInput.Blur()
}

func (m *TextInputModel) SetValue(value string) {
	m.textInput.SetValue(value)
}

func (m TextInputModel) Focused() bool {
	return m.textInput.Focused()
}
