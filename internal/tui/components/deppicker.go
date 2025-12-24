package components

import (
	"fmt"
	"strings"

	"github.com/KashifKhn/haft/internal/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

type DepCategory struct {
	Name         string
	Dependencies []DepItem
}

type DepItem struct {
	ID          string
	Name        string
	Description string
	Selected    bool
}

type DepPickerModel struct {
	categories      []DepCategory
	flatItems       []flatDepItem
	cursor          int
	label           string
	submitted       bool
	goBack          bool
	searchQuery     string
	searchMode      bool
	categoryFilter  int
	filteredIndices []int
	viewportStart   int
	viewportSize    int
	err             error
}

type flatDepItem struct {
	categoryIdx int
	depIdx      int
	isCategory  bool
}

type DepPickerConfig struct {
	Label      string
	Categories []DepCategory
}

func NewDepPicker(cfg DepPickerConfig) DepPickerModel {
	m := DepPickerModel{
		categories:     cfg.Categories,
		label:          cfg.Label,
		cursor:         0,
		categoryFilter: -1,
		viewportSize:   15,
	}
	m.buildFlatList()
	m.resetFilter()
	return m
}

func (m *DepPickerModel) buildFlatList() {
	m.flatItems = nil
	for catIdx, cat := range m.categories {
		m.flatItems = append(m.flatItems, flatDepItem{categoryIdx: catIdx, isCategory: true})
		for depIdx := range cat.Dependencies {
			m.flatItems = append(m.flatItems, flatDepItem{categoryIdx: catIdx, depIdx: depIdx, isCategory: false})
		}
	}
}

func (m *DepPickerModel) resetFilter() {
	m.filteredIndices = make([]int, len(m.flatItems))
	for i := range m.flatItems {
		m.filteredIndices[i] = i
	}
}

func (m *DepPickerModel) applyFilter() {
	m.filteredIndices = nil

	if m.categoryFilter >= 0 && m.categoryFilter < len(m.categories) {
		for i, item := range m.flatItems {
			if item.categoryIdx == m.categoryFilter {
				if m.searchQuery == "" {
					m.filteredIndices = append(m.filteredIndices, i)
				} else if !item.isCategory {
					dep := m.categories[item.categoryIdx].Dependencies[item.depIdx]
					query := strings.ToLower(m.searchQuery)
					if strings.Contains(strings.ToLower(dep.Name), query) ||
						strings.Contains(strings.ToLower(dep.Description), query) ||
						strings.Contains(strings.ToLower(dep.ID), query) {
						m.filteredIndices = append(m.filteredIndices, i)
					}
				}
			}
		}
		m.cursor = 0
		m.viewportStart = 0
		return
	}

	if m.searchQuery == "" {
		m.resetFilter()
		return
	}

	query := strings.ToLower(m.searchQuery)
	matchedCategories := make(map[int]bool)

	for i, item := range m.flatItems {
		if item.isCategory {
			continue
		}
		dep := m.categories[item.categoryIdx].Dependencies[item.depIdx]
		if strings.Contains(strings.ToLower(dep.Name), query) ||
			strings.Contains(strings.ToLower(dep.Description), query) ||
			strings.Contains(strings.ToLower(dep.ID), query) {
			matchedCategories[item.categoryIdx] = true
			m.filteredIndices = append(m.filteredIndices, i)
		}
	}

	var result []int
	for i, item := range m.flatItems {
		if item.isCategory && matchedCategories[item.categoryIdx] {
			result = append(result, i)
		} else if !item.isCategory {
			for _, fi := range m.filteredIndices {
				if fi == i {
					result = append(result, i)
					break
				}
			}
		}
	}
	m.filteredIndices = result
	m.cursor = 0
	m.viewportStart = 0
}

func (m DepPickerModel) Init() tea.Cmd {
	return nil
}

func (m DepPickerModel) Update(msg tea.Msg) (DepPickerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.searchMode {
			return m.handleSearchInput(msg)
		}
		return m.handleNavigation(msg)
	}
	return m, nil
}

func (m DepPickerModel) handleSearchInput(msg tea.KeyMsg) (DepPickerModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.searchMode = false
		m.searchQuery = ""
		m.resetFilter()
	case "enter":
		m.searchMode = false
	case "backspace":
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			m.applyFilter()
		}
	default:
		if len(msg.String()) == 1 {
			m.searchQuery += msg.String()
			m.applyFilter()
		}
	}
	return m, nil
}

func (m DepPickerModel) handleNavigation(msg tea.KeyMsg) (DepPickerModel, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.moveCursor(-1)
	case "down", "j":
		m.moveCursor(1)
	case "pgup":
		m.moveCursor(-m.viewportSize)
	case "pgdown":
		m.moveCursor(m.viewportSize)
	case "home", "g":
		m.cursor = 0
		m.viewportStart = 0
	case "end", "G":
		m.cursor = len(m.filteredIndices) - 1
		m.adjustViewport()
	case " ", "x":
		m.toggleCurrent()
	case "/":
		m.searchMode = true
		m.searchQuery = ""
	case "tab":
		m.nextCategory()
	case "shift+tab":
		m.prevCategory()
	case "enter":
		m.submitted = true
	case "esc":
		m.goBack = true
	case "a":
		m.selectAllVisible()
	case "n":
		m.selectNone()
	case "c":
		m.clearCategoryFilter()
	}
	return m, nil
}

func (m *DepPickerModel) moveCursor(delta int) {
	newCursor := m.cursor + delta
	for newCursor >= 0 && newCursor < len(m.filteredIndices) {
		item := m.flatItems[m.filteredIndices[newCursor]]
		if !item.isCategory {
			break
		}
		if delta > 0 {
			newCursor++
		} else {
			newCursor--
		}
	}

	if newCursor < 0 {
		newCursor = 0
		for newCursor < len(m.filteredIndices) && m.flatItems[m.filteredIndices[newCursor]].isCategory {
			newCursor++
		}
	}
	if newCursor >= len(m.filteredIndices) {
		newCursor = len(m.filteredIndices) - 1
		for newCursor >= 0 && m.flatItems[m.filteredIndices[newCursor]].isCategory {
			newCursor--
		}
	}

	if newCursor >= 0 && newCursor < len(m.filteredIndices) {
		m.cursor = newCursor
		m.adjustViewport()
	}
}

func (m *DepPickerModel) adjustViewport() {
	if m.cursor < m.viewportStart {
		m.viewportStart = m.cursor
	} else if m.cursor >= m.viewportStart+m.viewportSize {
		m.viewportStart = m.cursor - m.viewportSize + 1
	}
}

func (m *DepPickerModel) toggleCurrent() {
	if m.cursor < 0 || m.cursor >= len(m.filteredIndices) {
		return
	}
	flatIdx := m.filteredIndices[m.cursor]
	item := m.flatItems[flatIdx]
	if item.isCategory {
		return
	}
	m.categories[item.categoryIdx].Dependencies[item.depIdx].Selected =
		!m.categories[item.categoryIdx].Dependencies[item.depIdx].Selected
}

func (m *DepPickerModel) selectAllVisible() {
	for _, flatIdx := range m.filteredIndices {
		item := m.flatItems[flatIdx]
		if !item.isCategory {
			m.categories[item.categoryIdx].Dependencies[item.depIdx].Selected = true
		}
	}
}

func (m *DepPickerModel) selectNone() {
	for catIdx := range m.categories {
		for depIdx := range m.categories[catIdx].Dependencies {
			m.categories[catIdx].Dependencies[depIdx].Selected = false
		}
	}
}

func (m *DepPickerModel) nextCategory() {
	m.categoryFilter++
	if m.categoryFilter >= len(m.categories) {
		m.categoryFilter = -1
	}
	m.applyFilter()
}

func (m *DepPickerModel) prevCategory() {
	m.categoryFilter--
	if m.categoryFilter < -1 {
		m.categoryFilter = len(m.categories) - 1
	}
	m.applyFilter()
}

func (m *DepPickerModel) clearCategoryFilter() {
	m.categoryFilter = -1
	m.searchQuery = ""
	m.resetFilter()
}

func (m DepPickerModel) renderCategoryTabs() string {
	var tabs []string

	allLabel := "All"
	if m.categoryFilter == -1 {
		allLabel = "[All]"
		tabs = append(tabs, styles.Focused.Render(allLabel))
	} else {
		tabs = append(tabs, styles.Subtle.Render(allLabel))
	}

	for i, cat := range m.categories {
		shortName := cat.Name
		if len(shortName) > 12 {
			shortName = shortName[:10] + ".."
		}
		if i == m.categoryFilter {
			tabs = append(tabs, styles.Focused.Render("["+shortName+"]"))
		} else {
			tabs = append(tabs, styles.Subtle.Render(shortName))
		}
	}

	return strings.Join(tabs, " | ")
}

func (m DepPickerModel) View() string {
	var b strings.Builder

	if m.label != "" {
		b.WriteString(styles.RenderTitle(m.label))
		b.WriteString("\n")
	}

	b.WriteString(m.renderCategoryTabs())
	b.WriteString("\n")

	if m.searchMode {
		b.WriteString(styles.Focused.Render("Search: ") + m.searchQuery + "▌\n\n")
	} else if m.searchQuery != "" {
		b.WriteString(styles.Subtle.Render(fmt.Sprintf("Filter: %s (%d results)\n\n", m.searchQuery, m.countFilteredDeps())))
	}

	end := m.viewportStart + m.viewportSize
	if end > len(m.filteredIndices) {
		end = len(m.filteredIndices)
	}

	for i := m.viewportStart; i < end; i++ {
		flatIdx := m.filteredIndices[i]
		item := m.flatItems[flatIdx]

		if item.isCategory {
			if m.categoryFilter < 0 {
				catName := m.categories[item.categoryIdx].Name
				b.WriteString("\n" + styles.CategoryStyle.Render(catName) + "\n")
			}
			continue
		}

		dep := m.categories[item.categoryIdx].Dependencies[item.depIdx]
		cursor := "  "
		if m.cursor == i {
			cursor = styles.Arrow + " "
		}

		checkbox := "[ ]"
		if dep.Selected {
			checkbox = "[" + styles.CheckMark + "]"
		}

		var nameStyle, descStyle string
		if m.cursor == i {
			nameStyle = styles.ActiveItem.Render(dep.Name)
			descStyle = styles.Subtle.Render(dep.Description)
		} else if dep.Selected {
			nameStyle = styles.Selected.Render(dep.Name)
			descStyle = styles.Subtle.Render(dep.Description)
		} else {
			nameStyle = styles.InactiveItem.Render(dep.Name)
			descStyle = styles.Subtle.Render(dep.Description)
		}

		b.WriteString(fmt.Sprintf("%s%s %s\n", cursor, checkbox, nameStyle))
		if dep.Description != "" && m.cursor == i {
			b.WriteString(fmt.Sprintf("       %s\n", descStyle))
		}
	}

	selectedCount := m.countSelected()
	b.WriteString(fmt.Sprintf("\n%s\n", styles.Subtle.Render(fmt.Sprintf("Selected: %d", selectedCount))))

	if m.err != nil {
		b.WriteString(styles.RenderError(m.err.Error()) + "\n")
	}

	if m.searchMode {
		b.WriteString(styles.RenderHelp("type to search • enter: apply • esc: cancel"))
	} else {
		b.WriteString(styles.RenderHelp("↑/↓: navigate • space: toggle • tab: category • /: search • c: clear • enter: confirm"))
	}

	return b.String()
}

func (m DepPickerModel) countSelected() int {
	count := 0
	for _, cat := range m.categories {
		for _, dep := range cat.Dependencies {
			if dep.Selected {
				count++
			}
		}
	}
	return count
}

func (m DepPickerModel) countFilteredDeps() int {
	count := 0
	for _, flatIdx := range m.filteredIndices {
		if !m.flatItems[flatIdx].isCategory {
			count++
		}
	}
	return count
}

func (m DepPickerModel) Values() []string {
	var values []string
	for _, cat := range m.categories {
		for _, dep := range cat.Dependencies {
			if dep.Selected {
				values = append(values, dep.ID)
			}
		}
	}
	return values
}

func (m DepPickerModel) SelectedItems() []DepItem {
	var selected []DepItem
	for _, cat := range m.categories {
		for _, dep := range cat.Dependencies {
			if dep.Selected {
				selected = append(selected, dep)
			}
		}
	}
	return selected
}

func (m DepPickerModel) Submitted() bool {
	return m.submitted
}

func (m DepPickerModel) Validate() error {
	return nil
}

func (m DepPickerModel) GoBack() bool {
	return m.goBack
}

func (m *DepPickerModel) Reset() {
	m.submitted = false
	m.goBack = false
	m.err = nil
}
