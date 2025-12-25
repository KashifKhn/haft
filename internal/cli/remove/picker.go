package remove

import (
	"fmt"
	"strings"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/KashifKhn/haft/internal/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

type removePickerModel struct {
	deps          []buildtool.Dependency
	cursor        int
	selected      map[int]bool
	submitted     bool
	cancelled     bool
	viewportStart int
	viewportSize  int
	searchQuery   string
	filtered      []int
}

func newRemovePickerModel(deps []buildtool.Dependency) removePickerModel {
	m := removePickerModel{
		deps:         deps,
		selected:     make(map[int]bool),
		viewportSize: 15,
	}
	m.resetFilter()
	return m
}

func (m *removePickerModel) resetFilter() {
	m.filtered = make([]int, len(m.deps))
	for i := range m.deps {
		m.filtered[i] = i
	}
}

func (m *removePickerModel) applyFilter() {
	if m.searchQuery == "" {
		m.resetFilter()
		return
	}

	m.filtered = nil
	query := strings.ToLower(m.searchQuery)

	for i, dep := range m.deps {
		if strings.Contains(strings.ToLower(dep.ArtifactId), query) ||
			strings.Contains(strings.ToLower(dep.GroupId), query) {
			m.filtered = append(m.filtered, i)
		}
	}

	m.cursor = 0
	m.viewportStart = 0
}

func (m removePickerModel) Init() tea.Cmd {
	return nil
}

func (m removePickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.moveCursor(-1)
		case "down", "j":
			m.moveCursor(1)
		case "pgup":
			m.moveCursor(-m.viewportSize)
		case "pgdown":
			m.moveCursor(m.viewportSize)
		case " ", "x":
			m.toggleCurrent()
		case "enter":
			if m.countSelected() > 0 {
				m.submitted = true
				return m, tea.Quit
			}
		case "esc", "q", "ctrl+c":
			m.cancelled = true
			return m, tea.Quit
		case "a":
			m.selectAllVisible()
		case "n":
			m.selectNone()
		case "backspace":
			if len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				m.applyFilter()
			}
		default:
			if len(msg.String()) == 1 && msg.String() >= " " {
				m.searchQuery += msg.String()
				m.applyFilter()
			}
		}
	}
	return m, nil
}

func (m *removePickerModel) moveCursor(delta int) {
	if len(m.filtered) == 0 {
		return
	}

	m.cursor += delta
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor >= len(m.filtered) {
		m.cursor = len(m.filtered) - 1
	}
	m.adjustViewport()
}

func (m *removePickerModel) adjustViewport() {
	if m.cursor < m.viewportStart {
		m.viewportStart = m.cursor
	} else if m.cursor >= m.viewportStart+m.viewportSize {
		m.viewportStart = m.cursor - m.viewportSize + 1
	}
}

func (m *removePickerModel) toggleCurrent() {
	if m.cursor >= 0 && m.cursor < len(m.filtered) {
		idx := m.filtered[m.cursor]
		m.selected[idx] = !m.selected[idx]
	}
}

func (m *removePickerModel) selectAllVisible() {
	for _, idx := range m.filtered {
		m.selected[idx] = true
	}
}

func (m *removePickerModel) selectNone() {
	m.selected = make(map[int]bool)
}

func (m removePickerModel) countSelected() int {
	count := 0
	for _, sel := range m.selected {
		if sel {
			count++
		}
	}
	return count
}

func (m removePickerModel) View() string {
	var b strings.Builder

	b.WriteString(styles.RenderTitle("Remove Dependencies"))
	b.WriteString("\n")

	b.WriteString(styles.Focused.Render("Search: "))
	if m.searchQuery == "" {
		b.WriteString(styles.Subtle.Render("type to filter..."))
	} else {
		b.WriteString(m.searchQuery)
	}
	b.WriteString("▌\n")

	if len(m.filtered) == 0 {
		b.WriteString(styles.Subtle.Render("\n  No matching dependencies found\n"))
	} else {
		b.WriteString(styles.Subtle.Render(fmt.Sprintf("  %d dependencies\n\n", len(m.filtered))))
	}

	end := m.viewportStart + m.viewportSize
	if end > len(m.filtered) {
		end = len(m.filtered)
	}

	for i := m.viewportStart; i < end; i++ {
		idx := m.filtered[i]
		dep := m.deps[idx]

		cursor := "  "
		if m.cursor == i {
			cursor = styles.Arrow + " "
		}

		checkbox := "[ ]"
		if m.selected[idx] {
			checkbox = "[" + styles.CheckMark + "]"
		}

		artifactStyle := styles.InactiveItem.Render(dep.ArtifactId)
		if m.cursor == i {
			artifactStyle = styles.ActiveItem.Render(dep.ArtifactId)
		} else if m.selected[idx] {
			artifactStyle = styles.Selected.Render(dep.ArtifactId)
		}

		groupStyle := styles.Subtle.Render(dep.GroupId)
		b.WriteString(fmt.Sprintf("%s%s %s\n", cursor, checkbox, artifactStyle))
		b.WriteString(fmt.Sprintf("       %s", groupStyle))

		if dep.Version != "" {
			b.WriteString(styles.Subtle.Render(fmt.Sprintf(":%s", dep.Version)))
		}
		if dep.Scope != "" {
			b.WriteString(styles.Subtle.Render(fmt.Sprintf(" (%s)", dep.Scope)))
		}
		b.WriteString("\n")
	}

	selectedCount := m.countSelected()
	b.WriteString(fmt.Sprintf("\n%s\n", styles.Subtle.Render(fmt.Sprintf("Selected for removal: %d", selectedCount))))

	if selectedCount > 0 {
		b.WriteString(styles.WarningText.Render("Press Enter to remove selected dependencies\n"))
	}

	b.WriteString(styles.RenderHelp("↑/↓: navigate • space: toggle • a: all • n: none • enter: remove • esc: cancel"))

	return b.String()
}

func (m removePickerModel) selectedDeps() []buildtool.Dependency {
	var result []buildtool.Dependency
	for idx, sel := range m.selected {
		if sel && idx < len(m.deps) {
			result = append(result, m.deps[idx])
		}
	}
	return result
}

func RunRemovePicker(deps []buildtool.Dependency) ([]buildtool.Dependency, error) {
	if len(deps) == 0 {
		return nil, nil
	}

	model := newRemovePickerModel(deps)
	p := tea.NewProgram(model)

	finalModel, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run picker: %w", err)
	}

	m := finalModel.(removePickerModel)
	if m.cancelled {
		return nil, nil
	}

	return m.selectedDeps(), nil
}
