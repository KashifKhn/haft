package components

import (
	"github.com/KashifKhn/haft/internal/tui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TextInputModel struct {
	textInput      textinput.Model
	label          string
	placeholder    string
	helpText       string
	required       bool
	validator      func(string) error
	dynamicDefault func(map[string]any) string
	err            error
	submitted      bool
	goBack         bool
}

type TextInputConfig struct {
	Label          string
	Placeholder    string
	Default        string
	HelpText       string
	Required       bool
	CharLimit      int
	Width          int
	Validator      func(string) error
	DynamicDefault func(map[string]any) string
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

	if cfg.Default != "" {
		ti.SetValue(cfg.Default)
	}

	return TextInputModel{
		textInput:      ti,
		label:          cfg.Label,
		placeholder:    cfg.Placeholder,
		helpText:       cfg.HelpText,
		required:       cfg.Required,
		validator:      cfg.Validator,
		dynamicDefault: cfg.DynamicDefault,
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
			m.goBack = true
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

	if m.helpText != "" {
		view += "\n" + styles.RenderHelp(m.helpText)
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

func (m TextInputModel) GoBack() bool {
	return m.goBack
}

func (m *TextInputModel) Reset() {
	m.submitted = false
	m.goBack = false
	m.err = nil
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

func (m *TextInputModel) ApplyDynamicDefault(values map[string]any) {
	if m.dynamicDefault != nil && m.textInput.Value() == "" {
		if defaultVal := m.dynamicDefault(values); defaultVal != "" {
			m.textInput.SetValue(defaultVal)
		}
	}
}

func (m TextInputModel) Focused() bool {
	return m.textInput.Focused()
}
