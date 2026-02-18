package components

import (
	"fmt"
	"strings"

	"github.com/KashifKhn/haft/internal/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

type SelectItem struct {
	Label       string
	Value       string
	Description string
}

type SelectModel struct {
	items     []SelectItem
	cursor    int
	selected  int
	label     string
	helpText  string
	submitted bool
	goBack    bool
}

type SelectConfig struct {
	Label    string
	Items    []SelectItem
	HelpText string
	Default  string
}

func NewSelect(cfg SelectConfig) SelectModel {
	cursor := 0
	if cfg.Default != "" {
		for i, item := range cfg.Items {
			if item.Value == cfg.Default {
				cursor = i
				break
			}
		}
	}
	return SelectModel{
		items:    cfg.Items,
		label:    cfg.Label,
		helpText: cfg.HelpText,
		cursor:   cursor,
		selected: -1,
	}
}

func (m SelectModel) Init() tea.Cmd {
	return nil
}

func (m SelectModel) Update(msg tea.Msg) (SelectModel, tea.Cmd) {
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
		case "enter", " ":
			m.selected = m.cursor
			m.submitted = true
			return m, nil
		case "esc":
			m.goBack = true
			return m, nil
		}
	}

	return m, nil
}

func (m SelectModel) View() string {
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

		var line string
		if m.cursor == i {
			line = styles.ActiveItem.Render(item.Label)
		} else {
			line = styles.InactiveItem.Render(item.Label)
		}

		fmt.Fprintf(&b, "%s%s\n", cursor, line)

		if item.Description != "" && m.cursor == i {
			fmt.Fprintf(&b, "     %s\n", styles.Subtle.Render(item.Description))
		}
	}

	if m.helpText != "" {
		b.WriteString("\n" + styles.RenderHelp(m.helpText))
	} else {
		b.WriteString(styles.RenderHelp("↑/↓: navigate • enter: select"))
	}

	return b.String()
}

func (m SelectModel) Value() string {
	if m.selected >= 0 && m.selected < len(m.items) {
		return m.items[m.selected].Value
	}
	return ""
}

func (m SelectModel) SelectedItem() (SelectItem, bool) {
	if m.selected >= 0 && m.selected < len(m.items) {
		return m.items[m.selected], true
	}
	return SelectItem{}, false
}

func (m SelectModel) Submitted() bool {
	return m.submitted
}

func (m SelectModel) GoBack() bool {
	return m.goBack
}

func (m *SelectModel) Reset() {
	m.submitted = false
	m.goBack = false
}

func (m SelectModel) SelectedIndex() int {
	return m.selected
}
