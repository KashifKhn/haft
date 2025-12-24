package wizard

import (
	"fmt"
	"strings"

	"github.com/KashifKhn/haft/internal/tui/components"
	"github.com/KashifKhn/haft/internal/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

type StepType int

const (
	StepTypeInput StepType = iota
	StepTypeSelect
	StepTypeMultiSelect
	StepTypeConfirm
)

type Step interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Step, tea.Cmd)
	View() string
	Value() any
	Validate() error
	Submitted() bool
}

type TextInputStep struct {
	model components.TextInputModel
}

func NewTextInputStep(cfg components.TextInputConfig) *TextInputStep {
	return &TextInputStep{model: components.NewTextInput(cfg)}
}

func (s *TextInputStep) Init() tea.Cmd {
	return s.model.Init()
}

func (s *TextInputStep) Update(msg tea.Msg) (Step, tea.Cmd) {
	updated, cmd := s.model.Update(msg)
	s.model = updated
	return s, cmd
}

func (s *TextInputStep) View() string {
	return s.model.View()
}

func (s *TextInputStep) Value() any {
	return s.model.Value()
}

func (s *TextInputStep) Validate() error {
	return s.model.Validate()
}

func (s *TextInputStep) Submitted() bool {
	return s.model.Submitted()
}

type SelectStep struct {
	model components.SelectModel
}

func NewSelectStep(cfg components.SelectConfig) *SelectStep {
	return &SelectStep{model: components.NewSelect(cfg)}
}

func (s *SelectStep) Init() tea.Cmd {
	return s.model.Init()
}

func (s *SelectStep) Update(msg tea.Msg) (Step, tea.Cmd) {
	updated, cmd := s.model.Update(msg)
	s.model = updated
	return s, cmd
}

func (s *SelectStep) View() string {
	return s.model.View()
}

func (s *SelectStep) Value() any {
	return s.model.Value()
}

func (s *SelectStep) Validate() error {
	return nil
}

func (s *SelectStep) Submitted() bool {
	return s.model.Submitted()
}

type MultiSelectStep struct {
	model components.MultiSelectModel
}

func NewMultiSelectStep(cfg components.MultiSelectConfig) *MultiSelectStep {
	return &MultiSelectStep{model: components.NewMultiSelect(cfg)}
}

func (s *MultiSelectStep) Init() tea.Cmd {
	return s.model.Init()
}

func (s *MultiSelectStep) Update(msg tea.Msg) (Step, tea.Cmd) {
	updated, cmd := s.model.Update(msg)
	s.model = updated
	return s, cmd
}

func (s *MultiSelectStep) View() string {
	return s.model.View()
}

func (s *MultiSelectStep) Value() any {
	return s.model.Values()
}

func (s *MultiSelectStep) Validate() error {
	return s.model.Validate()
}

func (s *MultiSelectStep) Submitted() bool {
	return s.model.Submitted()
}

type WizardModel struct {
	title       string
	steps       []Step
	currentStep int
	cancelled   bool
	completed   bool
	values      map[string]any
	stepKeys    []string
}

type WizardConfig struct {
	Title    string
	Steps    []Step
	StepKeys []string
}

func New(cfg WizardConfig) WizardModel {
	return WizardModel{
		title:       cfg.Title,
		steps:       cfg.Steps,
		stepKeys:    cfg.StepKeys,
		currentStep: 0,
		values:      make(map[string]any),
	}
}

func (m WizardModel) Init() tea.Cmd {
	if len(m.steps) > 0 {
		return m.steps[0].Init()
	}
	return nil
}

func (m WizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.cancelled = true
			return m, tea.Quit
		}
	}

	if m.currentStep >= len(m.steps) {
		return m, nil
	}

	updatedStep, cmd := m.steps[m.currentStep].Update(msg)
	m.steps[m.currentStep] = updatedStep

	if updatedStep.Submitted() {
		if m.currentStep < len(m.stepKeys) {
			m.values[m.stepKeys[m.currentStep]] = updatedStep.Value()
		}

		if m.currentStep < len(m.steps)-1 {
			m.currentStep++
			return m, m.steps[m.currentStep].Init()
		} else {
			m.completed = true
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m WizardModel) View() string {
	var b strings.Builder

	if m.title != "" {
		b.WriteString(styles.RenderTitle(m.title))
		b.WriteString("\n")
	}

	b.WriteString(m.renderProgress())
	b.WriteString("\n\n")

	if m.currentStep < len(m.steps) {
		b.WriteString(m.steps[m.currentStep].View())
	}

	return b.String()
}

func (m WizardModel) renderProgress() string {
	var parts []string

	for i := range m.steps {
		if i < m.currentStep {
			parts = append(parts, styles.CompletedStep.Render(fmt.Sprintf("(%d)", i+1)))
		} else if i == m.currentStep {
			parts = append(parts, styles.ActiveStep.Render(fmt.Sprintf("[%d]", i+1)))
		} else {
			parts = append(parts, styles.StepIndicator.Render(fmt.Sprintf("(%d)", i+1)))
		}
	}

	return strings.Join(parts, " → ")
}

func (m WizardModel) Values() map[string]any {
	return m.values
}

func (m WizardModel) Value(key string) any {
	return m.values[key]
}

func (m WizardModel) StringValue(key string) string {
	if v, ok := m.values[key].(string); ok {
		return v
	}
	return ""
}

func (m WizardModel) StringSliceValue(key string) []string {
	if v, ok := m.values[key].([]string); ok {
		return v
	}
	return nil
}

func (m WizardModel) Cancelled() bool {
	return m.cancelled
}

func (m WizardModel) Completed() bool {
	return m.completed
}

func (m WizardModel) CurrentStepIndex() int {
	return m.currentStep
}

func (m WizardModel) TotalSteps() int {
	return len(m.steps)
}
