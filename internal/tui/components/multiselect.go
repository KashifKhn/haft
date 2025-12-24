package components

import (
	"fmt"
	"strings"

	"github.com/KashifKhn/haft/internal/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

type MultiSelectItem struct {
	Label    string
	Value    string
	Selected bool
}

type MultiSelectModel struct {
	items     []MultiSelectItem
	cursor    int
	label     string
	submitted bool
	goBack    bool
	required  bool
	minSelect int
	maxSelect int
	err       error
}

type MultiSelectConfig struct {
	Label     string
	Items     []MultiSelectItem
	Required  bool
	MinSelect int
	MaxSelect int
}

func NewMultiSelect(cfg MultiSelectConfig) MultiSelectModel {
	return MultiSelectModel{
		items:     cfg.Items,
		label:     cfg.Label,
		cursor:    0,
		required:  cfg.Required,
		minSelect: cfg.MinSelect,
		maxSelect: cfg.MaxSelect,
	}
}

func (m MultiSelectModel) Init() tea.Cmd {
	return nil
}

func (m MultiSelectModel) Update(msg tea.Msg) (MultiSelectModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case " ", "x":
			m.toggleCurrent()
			m.err = nil
		case "enter":
			if err := m.Validate(); err != nil {
				m.err = err
				return m, nil
			}
			m.submitted = true
			return m, nil
		case "esc":
			m.goBack = true
			return m, nil
		case "a":
			m.selectAll()
		case "n":
			m.selectNone()
		}
	}

	return m, nil
}

func (m *MultiSelectModel) toggleCurrent() {
	if m.cursor >= 0 && m.cursor < len(m.items) {
		currentSelected := m.countSelected()
		if m.items[m.cursor].Selected {
			m.items[m.cursor].Selected = false
		} else if m.maxSelect == 0 || currentSelected < m.maxSelect {
			m.items[m.cursor].Selected = true
		}
	}
}

func (m *MultiSelectModel) selectAll() {
	for i := range m.items {
		if m.maxSelect > 0 && m.countSelected() >= m.maxSelect {
			break
		}
		m.items[i].Selected = true
	}
}

func (m *MultiSelectModel) selectNone() {
	for i := range m.items {
		m.items[i].Selected = false
	}
}

func (m MultiSelectModel) countSelected() int {
	count := 0
	for _, item := range m.items {
		if item.Selected {
			count++
		}
	}
	return count
}

func (m MultiSelectModel) View() string {
	var b strings.Builder

	if m.label != "" {
		b.WriteString(styles.RenderTitle(m.label))
		b.WriteString("\n")
	}

	for i, item := range m.items {
		cursor := "  "
		if m.cursor == i {
			cursor = styles.Arrow + " "
		}

		checkbox := "[ ]"
		if item.Selected {
			checkbox = "[" + styles.CheckMark + "]"
		}

		var line string
		if m.cursor == i {
			line = styles.ActiveItem.Render(item.Label)
		} else if item.Selected {
			line = styles.Selected.Render(item.Label)
		} else {
			line = styles.InactiveItem.Render(item.Label)
		}

		b.WriteString(fmt.Sprintf("%s%s %s\n", cursor, checkbox, line))
	}

	if m.err != nil {
		b.WriteString("\n" + styles.RenderError(m.err.Error()) + "\n")
	}

	b.WriteString(styles.RenderHelp("↑/↓: navigate • space: toggle • a: all • n: none • enter: confirm"))

	return b.String()
}

func (m MultiSelectModel) Values() []string {
	var values []string
	for _, item := range m.items {
		if item.Selected {
			values = append(values, item.Value)
		}
	}
	return values
}

func (m MultiSelectModel) SelectedItems() []MultiSelectItem {
	var selected []MultiSelectItem
	for _, item := range m.items {
		if item.Selected {
			selected = append(selected, item)
		}
	}
	return selected
}

func (m MultiSelectModel) Validate() error {
	selectedCount := m.countSelected()

	if m.required && selectedCount == 0 {
		return errRequired
	}

	if m.minSelect > 0 && selectedCount < m.minSelect {
		return fmt.Errorf("select at least %d items", m.minSelect)
	}

	if m.maxSelect > 0 && selectedCount > m.maxSelect {
		return fmt.Errorf("select at most %d items", m.maxSelect)
	}

	return nil
}

func (m MultiSelectModel) Submitted() bool {
	return m.submitted
}

func (m MultiSelectModel) GoBack() bool {
	return m.goBack
}

func (m *MultiSelectModel) Reset() {
	m.submitted = false
	m.goBack = false
	m.err = nil
}
