package add

import (
	"fmt"
	"sort"
	"strings"

	"github.com/KashifKhn/haft/internal/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

type pickerItem struct {
	alias       string
	name        string
	description string
	category    string
	selected    bool
}

type PickerModel struct {
	items         []pickerItem
	filtered      []int
	cursor        int
	searchQuery   string
	submitted     bool
	cancelled     bool
	viewportStart int
	viewportSize  int
}

func NewPicker() PickerModel {
	items := buildPickerItems()
	m := PickerModel{
		items:        items,
		viewportSize: 15,
	}
	m.resetFilter()
	return m
}

func buildPickerItems() []pickerItem {
	var items []pickerItem

	categories := GetCatalogByCategory()
	categoryOrder := []string{
		"Web", "SQL", "NoSQL", "Security", "Messaging",
		"I/O", "Template Engines", "Ops", "Developer Tools", "Testing",
	}

	for _, category := range categoryOrder {
		aliases, ok := categories[category]
		if !ok {
			continue
		}
		sort.Strings(aliases)
		for _, alias := range aliases {
			entry, _ := GetCatalogEntry(alias)
			items = append(items, pickerItem{
				alias:       alias,
				name:        entry.Name,
				description: entry.Description,
				category:    entry.Category,
			})
		}
	}

	return items
}

func (m *PickerModel) resetFilter() {
	m.filtered = make([]int, len(m.items))
	for i := range m.items {
		m.filtered[i] = i
	}
}

func (m *PickerModel) applyFilter() {
	if m.searchQuery == "" {
		m.resetFilter()
		return
	}

	m.filtered = nil
	query := strings.ToLower(m.searchQuery)

	for i, item := range m.items {
		if strings.Contains(strings.ToLower(item.alias), query) ||
			strings.Contains(strings.ToLower(item.name), query) ||
			strings.Contains(strings.ToLower(item.description), query) ||
			strings.Contains(strings.ToLower(item.category), query) {
			m.filtered = append(m.filtered, i)
		}
	}

	m.cursor = 0
	m.viewportStart = 0
}

func (m PickerModel) Init() tea.Cmd {
	return nil
}

func (m PickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *PickerModel) moveCursor(delta int) {
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

func (m *PickerModel) adjustViewport() {
	if m.cursor < m.viewportStart {
		m.viewportStart = m.cursor
	} else if m.cursor >= m.viewportStart+m.viewportSize {
		m.viewportStart = m.cursor - m.viewportSize + 1
	}
}

func (m *PickerModel) toggleCurrent() {
	if m.cursor >= 0 && m.cursor < len(m.filtered) {
		idx := m.filtered[m.cursor]
		m.items[idx].selected = !m.items[idx].selected
	}
}

func (m *PickerModel) selectAllVisible() {
	for _, idx := range m.filtered {
		m.items[idx].selected = true
	}
}

func (m *PickerModel) selectNone() {
	for i := range m.items {
		m.items[i].selected = false
	}
}

func (m PickerModel) countSelected() int {
	count := 0
	for _, item := range m.items {
		if item.selected {
			count++
		}
	}
	return count
}

func (m PickerModel) View() string {
	var b strings.Builder

	b.WriteString(styles.RenderTitle("Add Dependencies"))
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
		b.WriteString(styles.Subtle.Render(fmt.Sprintf("  %d results\n\n", len(m.filtered))))
	}

	end := m.viewportStart + m.viewportSize
	if end > len(m.filtered) {
		end = len(m.filtered)
	}

	currentCategory := ""
	for i := m.viewportStart; i < end; i++ {
		idx := m.filtered[i]
		item := m.items[idx]

		if item.category != currentCategory && m.searchQuery == "" {
			currentCategory = item.category
			b.WriteString(styles.CategoryStyle.Render(currentCategory) + "\n")
		}

		cursor := "  "
		if m.cursor == i {
			cursor = styles.Arrow + " "
		}

		checkbox := "[ ]"
		if item.selected {
			checkbox = "[" + styles.CheckMark + "]"
		}

		var nameStyle string
		if m.cursor == i {
			nameStyle = styles.ActiveItem.Render(item.name)
		} else if item.selected {
			nameStyle = styles.Selected.Render(item.name)
		} else {
			nameStyle = styles.InactiveItem.Render(item.name)
		}

		aliasText := styles.Subtle.Render(fmt.Sprintf("(%s)", item.alias))
		b.WriteString(fmt.Sprintf("%s%s %s %s\n", cursor, checkbox, nameStyle, aliasText))

		if m.cursor == i && item.description != "" {
			b.WriteString(fmt.Sprintf("       %s\n", styles.Subtle.Render(item.description)))
		}
	}

	selectedCount := m.countSelected()
	b.WriteString(fmt.Sprintf("\n%s\n", styles.Subtle.Render(fmt.Sprintf("Selected: %d", selectedCount))))

	b.WriteString(styles.RenderHelp("↑/↓: navigate • space: toggle • a: all • n: none • enter: add • esc: cancel"))

	return b.String()
}

func (m PickerModel) SelectedAliases() []string {
	var aliases []string
	for _, item := range m.items {
		if item.selected {
			aliases = append(aliases, item.alias)
		}
	}
	return aliases
}

func (m PickerModel) Submitted() bool {
	return m.submitted
}

func (m PickerModel) Cancelled() bool {
	return m.cancelled
}

func RunPicker() ([]string, error) {
	model := NewPicker()
	p := tea.NewProgram(model)

	finalModel, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run picker: %w", err)
	}

	m := finalModel.(PickerModel)
	if m.Cancelled() {
		return nil, nil
	}

	return m.SelectedAliases(), nil
}
