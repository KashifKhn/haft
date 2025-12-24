package main

import (
	"fmt"
	"os"

	"github.com/KashifKhn/haft/internal/tui/components"
	"github.com/KashifKhn/haft/internal/tui/wizard"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run examples/tui_demo/main.go <component>")
		fmt.Println("")
		fmt.Println("Components:")
		fmt.Println("  textinput   - Test text input component")
		fmt.Println("  select      - Test select component")
		fmt.Println("  multiselect - Test multi-select component")
		fmt.Println("  wizard      - Test full wizard flow")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "textinput":
		runTextInput()
	case "select":
		runSelect()
	case "multiselect":
		runMultiSelect()
	case "wizard":
		runWizard()
	default:
		fmt.Printf("Unknown component: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func runTextInput() {
	model := components.NewTextInput(components.TextInputConfig{
		Label:       "Project Name",
		Placeholder: "my-spring-app",
		Required:    true,
	})

	wrapper := &textInputWrapper{model: model}
	p := tea.NewProgram(wrapper)

	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if w, ok := finalModel.(*textInputWrapper); ok {
		fmt.Printf("\nYou entered: %s\n", w.model.Value())
	}
}

type textInputWrapper struct {
	model components.TextInputModel
}

func (w *textInputWrapper) Init() tea.Cmd {
	return w.model.Init()
}

func (w *textInputWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "ctrl+c" {
		return w, tea.Quit
	}

	updated, cmd := w.model.Update(msg)
	w.model = updated

	if w.model.Submitted() {
		return w, tea.Quit
	}

	return w, cmd
}

func (w *textInputWrapper) View() string {
	return w.model.View() + "\n\n(Press Enter to submit, Ctrl+C to quit)"
}

func runSelect() {
	model := components.NewSelect(components.SelectConfig{
		Label: "Select Build Tool",
		Items: []components.SelectItem{
			{Label: "Maven", Value: "maven"},
			{Label: "Gradle (Groovy)", Value: "gradle"},
			{Label: "Gradle (Kotlin)", Value: "gradle-kotlin"},
		},
	})

	wrapper := &selectWrapper{model: model}
	p := tea.NewProgram(wrapper)

	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if w, ok := finalModel.(*selectWrapper); ok {
		fmt.Printf("\nYou selected: %s\n", w.model.Value())
	}
}

type selectWrapper struct {
	model components.SelectModel
}

func (w *selectWrapper) Init() tea.Cmd {
	return w.model.Init()
}

func (w *selectWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "ctrl+c" {
		return w, tea.Quit
	}

	updated, cmd := w.model.Update(msg)
	w.model = updated

	if w.model.Submitted() {
		return w, tea.Quit
	}

	return w, cmd
}

func (w *selectWrapper) View() string {
	return w.model.View()
}

func runMultiSelect() {
	model := components.NewMultiSelect(components.MultiSelectConfig{
		Label: "Select Dependencies",
		Items: []components.MultiSelectItem{
			{Label: "Spring Web", Value: "web"},
			{Label: "Spring Data JPA", Value: "jpa"},
			{Label: "Spring Security", Value: "security"},
			{Label: "Lombok", Value: "lombok"},
			{Label: "Spring Validation", Value: "validation"},
		},
		Required:  true,
		MinSelect: 1,
	})

	wrapper := &multiSelectWrapper{model: model}
	p := tea.NewProgram(wrapper)

	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if w, ok := finalModel.(*multiSelectWrapper); ok {
		fmt.Printf("\nYou selected: %v\n", w.model.Values())
	}
}

type multiSelectWrapper struct {
	model components.MultiSelectModel
}

func (w *multiSelectWrapper) Init() tea.Cmd {
	return w.model.Init()
}

func (w *multiSelectWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "ctrl+c" {
		return w, tea.Quit
	}

	updated, cmd := w.model.Update(msg)
	w.model = updated

	if w.model.Submitted() {
		return w, tea.Quit
	}

	return w, cmd
}

func (w *multiSelectWrapper) View() string {
	return w.model.View()
}

func runWizard() {
	steps := []wizard.Step{
		wizard.NewTextInputStep(components.TextInputConfig{
			Label:       "Project Name",
			Placeholder: "my-spring-app",
			Required:    true,
		}),
		wizard.NewSelectStep(components.SelectConfig{
			Label: "Build Tool",
			Items: []components.SelectItem{
				{Label: "Maven", Value: "maven"},
				{Label: "Gradle", Value: "gradle"},
			},
		}),
		wizard.NewSelectStep(components.SelectConfig{
			Label: "Java Version",
			Items: []components.SelectItem{
				{Label: "Java 21 (LTS)", Value: "21"},
				{Label: "Java 17 (LTS)", Value: "17"},
				{Label: "Java 11 (LTS)", Value: "11"},
			},
		}),
		wizard.NewMultiSelectStep(components.MultiSelectConfig{
			Label: "Dependencies",
			Items: []components.MultiSelectItem{
				{Label: "Spring Web", Value: "web"},
				{Label: "Spring Data JPA", Value: "jpa"},
				{Label: "Lombok", Value: "lombok"},
				{Label: "Validation", Value: "validation"},
			},
		}),
	}

	w := wizard.New(wizard.WizardConfig{
		Title:    "Create New Spring Boot Project",
		Steps:    steps,
		StepKeys: []string{"name", "buildTool", "javaVersion", "dependencies"},
	})

	p := tea.NewProgram(w)

	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if wiz, ok := finalModel.(wizard.WizardModel); ok {
		if wiz.Cancelled() {
			fmt.Println("\nWizard cancelled")
			return
		}

		fmt.Println("\n✓ Project Configuration:")
		fmt.Printf("  Name:         %s\n", wiz.StringValue("name"))
		fmt.Printf("  Build Tool:   %s\n", wiz.StringValue("buildTool"))
		fmt.Printf("  Java Version: %s\n", wiz.StringValue("javaVersion"))
		fmt.Printf("  Dependencies: %v\n", wiz.StringSliceValue("dependencies"))
	}
}
