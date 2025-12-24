package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestNewTextInput(t *testing.T) {
	cfg := TextInputConfig{
		Label:       "Name",
		Placeholder: "Enter your name",
		Required:    true,
	}
	model := NewTextInput(cfg)

	assert.Equal(t, "Name", model.label)
	assert.True(t, model.required)
	assert.False(t, model.Submitted())
	assert.Equal(t, "", model.Value())
}

func TestTextInputValidation(t *testing.T) {
	cfg := TextInputConfig{
		Label:    "Name",
		Required: true,
	}
	model := NewTextInput(cfg)

	err := model.Validate()
	assert.Error(t, err)
	assert.Equal(t, errRequired, err)

	model.SetValue("test")
	err = model.Validate()
	assert.NoError(t, err)
}

func TestTextInputCustomValidator(t *testing.T) {
	cfg := TextInputConfig{
		Label: "Email",
		Validator: func(s string) error {
			if len(s) < 5 {
				return errRequired
			}
			return nil
		},
	}
	model := NewTextInput(cfg)
	model.SetValue("ab")

	err := model.Validate()
	assert.Error(t, err)

	model.SetValue("test@example.com")
	err = model.Validate()
	assert.NoError(t, err)
}

func TestNewSelect(t *testing.T) {
	cfg := SelectConfig{
		Label: "Choose",
		Items: []SelectItem{
			{Label: "Option 1", Value: "opt1"},
			{Label: "Option 2", Value: "opt2"},
		},
	}
	model := NewSelect(cfg)

	assert.Equal(t, "Choose", model.label)
	assert.Len(t, model.items, 2)
	assert.Equal(t, 0, model.cursor)
	assert.Equal(t, -1, model.selected)
	assert.False(t, model.Submitted())
	assert.Equal(t, "", model.Value())
}

func TestNewSelectWithDefault(t *testing.T) {
	tests := []struct {
		name           string
		defaultValue   string
		expectedCursor int
	}{
		{
			name:           "default to second item",
			defaultValue:   "opt2",
			expectedCursor: 1,
		},
		{
			name:           "default to third item",
			defaultValue:   "opt3",
			expectedCursor: 2,
		},
		{
			name:           "default to first item",
			defaultValue:   "opt1",
			expectedCursor: 0,
		},
		{
			name:           "non-existent default falls back to first",
			defaultValue:   "nonexistent",
			expectedCursor: 0,
		},
		{
			name:           "empty default stays at first",
			defaultValue:   "",
			expectedCursor: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := SelectConfig{
				Label: "Choose",
				Items: []SelectItem{
					{Label: "Option 1", Value: "opt1"},
					{Label: "Option 2", Value: "opt2"},
					{Label: "Option 3", Value: "opt3"},
				},
				Default: tt.defaultValue,
			}
			model := NewSelect(cfg)

			assert.Equal(t, tt.expectedCursor, model.cursor)
			assert.Equal(t, -1, model.selected)
			assert.False(t, model.Submitted())
		})
	}
}

func TestSelectWithDefaultAndSubmit(t *testing.T) {
	cfg := SelectConfig{
		Label: "Config Format",
		Items: []SelectItem{
			{Label: "Properties", Value: "properties"},
			{Label: "YAML", Value: "yaml"},
		},
		Default: "yaml",
	}
	model := NewSelect(cfg)

	assert.Equal(t, 1, model.cursor)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	assert.True(t, model.Submitted())
	assert.Equal(t, "yaml", model.Value())
	assert.Equal(t, 1, model.SelectedIndex())
}

func TestSelectNavigationWithVimKeys(t *testing.T) {
	cfg := SelectConfig{
		Items: []SelectItem{
			{Label: "A", Value: "a"},
			{Label: "B", Value: "b"},
			{Label: "C", Value: "c"},
		},
	}
	model := NewSelect(cfg)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	assert.Equal(t, 1, model.cursor)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	assert.Equal(t, 2, model.cursor)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	assert.Equal(t, 1, model.cursor)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	assert.Equal(t, 0, model.cursor)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	assert.Equal(t, 0, model.cursor)
}

func TestSelectEscapeGoBack(t *testing.T) {
	cfg := SelectConfig{
		Items: []SelectItem{
			{Label: "A", Value: "a"},
		},
	}
	model := NewSelect(cfg)

	assert.False(t, model.GoBack())

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})

	assert.True(t, model.GoBack())
	assert.False(t, model.Submitted())
}

func TestSelectReset(t *testing.T) {
	cfg := SelectConfig{
		Items: []SelectItem{
			{Label: "A", Value: "a"},
		},
	}
	model := NewSelect(cfg)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.True(t, model.Submitted())

	model.Reset()

	assert.False(t, model.Submitted())
	assert.False(t, model.GoBack())
}

func TestSelectSelectedItem(t *testing.T) {
	cfg := SelectConfig{
		Items: []SelectItem{
			{Label: "Option A", Value: "a", Description: "First option"},
			{Label: "Option B", Value: "b", Description: "Second option"},
		},
	}
	model := NewSelect(cfg)

	_, ok := model.SelectedItem()
	assert.False(t, ok)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	item, ok := model.SelectedItem()
	assert.True(t, ok)
	assert.Equal(t, "Option B", item.Label)
	assert.Equal(t, "b", item.Value)
	assert.Equal(t, "Second option", item.Description)
}

func TestSelectView(t *testing.T) {
	cfg := SelectConfig{
		Label: "Choose Option",
		Items: []SelectItem{
			{Label: "Option A", Value: "a", Description: "First option"},
			{Label: "Option B", Value: "b"},
		},
		HelpText: "Select an option",
	}
	model := NewSelect(cfg)

	view := model.View()

	assert.Contains(t, view, "Choose Option")
	assert.Contains(t, view, "Option A")
	assert.Contains(t, view, "Option B")
	assert.Contains(t, view, "First option")
	assert.Contains(t, view, "Select an option")
}

func TestSelectViewDefaultHelp(t *testing.T) {
	cfg := SelectConfig{
		Items: []SelectItem{
			{Label: "A", Value: "a"},
		},
	}
	model := NewSelect(cfg)

	view := model.View()

	assert.Contains(t, view, "navigate")
	assert.Contains(t, view, "select")
}

func TestSelectInit(t *testing.T) {
	cfg := SelectConfig{
		Items: []SelectItem{
			{Label: "A", Value: "a"},
		},
	}
	model := NewSelect(cfg)

	cmd := model.Init()

	assert.Nil(t, cmd)
}

func TestSelectSpaceKeySelection(t *testing.T) {
	cfg := SelectConfig{
		Items: []SelectItem{
			{Label: "A", Value: "a"},
			{Label: "B", Value: "b"},
		},
	}
	model := NewSelect(cfg)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeySpace})

	assert.True(t, model.Submitted())
	assert.Equal(t, "b", model.Value())
}

func TestSelectNavigation(t *testing.T) {
	cfg := SelectConfig{
		Items: []SelectItem{
			{Label: "A", Value: "a"},
			{Label: "B", Value: "b"},
			{Label: "C", Value: "c"},
		},
	}
	model := NewSelect(cfg)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, 1, model.cursor)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, 2, model.cursor)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, 2, model.cursor)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	assert.Equal(t, 1, model.cursor)
}

func TestSelectSelection(t *testing.T) {
	cfg := SelectConfig{
		Items: []SelectItem{
			{Label: "A", Value: "a"},
			{Label: "B", Value: "b"},
		},
	}
	model := NewSelect(cfg)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	assert.True(t, model.Submitted())
	assert.Equal(t, "b", model.Value())
	assert.Equal(t, 1, model.SelectedIndex())
}

func TestNewMultiSelect(t *testing.T) {
	cfg := MultiSelectConfig{
		Label: "Choose",
		Items: []MultiSelectItem{
			{Label: "Option 1", Value: "opt1"},
			{Label: "Option 2", Value: "opt2"},
		},
		Required: true,
	}
	model := NewMultiSelect(cfg)

	assert.Equal(t, "Choose", model.label)
	assert.Len(t, model.items, 2)
	assert.True(t, model.required)
	assert.False(t, model.Submitted())
	assert.Empty(t, model.Values())
}

func TestMultiSelectToggle(t *testing.T) {
	cfg := MultiSelectConfig{
		Items: []MultiSelectItem{
			{Label: "A", Value: "a"},
			{Label: "B", Value: "b"},
		},
	}
	model := NewMultiSelect(cfg)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeySpace})
	assert.True(t, model.items[0].Selected)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeySpace})

	values := model.Values()
	assert.Len(t, values, 2)
	assert.Contains(t, values, "a")
	assert.Contains(t, values, "b")
}

func TestMultiSelectValidation(t *testing.T) {
	cfg := MultiSelectConfig{
		Items: []MultiSelectItem{
			{Label: "A", Value: "a"},
			{Label: "B", Value: "b"},
		},
		Required:  true,
		MinSelect: 1,
	}
	model := NewMultiSelect(cfg)

	err := model.Validate()
	assert.Error(t, err)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeySpace})
	err = model.Validate()
	assert.NoError(t, err)
}

func TestMultiSelectMaxSelect(t *testing.T) {
	cfg := MultiSelectConfig{
		Items: []MultiSelectItem{
			{Label: "A", Value: "a"},
			{Label: "B", Value: "b"},
			{Label: "C", Value: "c"},
		},
		MaxSelect: 2,
	}
	model := NewMultiSelect(cfg)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeySpace})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeySpace})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeySpace})

	values := model.Values()
	assert.Len(t, values, 2)
}

func TestNewSpinner(t *testing.T) {
	cfg := SpinnerConfig{
		Message: "Loading...",
	}
	model := NewSpinner(cfg)

	assert.Equal(t, "Loading...", model.message)
	assert.False(t, model.Done())
	assert.Nil(t, model.Err())
}

func TestSpinnerDone(t *testing.T) {
	cfg := SpinnerConfig{Message: "Processing..."}
	model := NewSpinner(cfg)

	model, _ = model.Update(SpinnerDoneMsg{Err: nil})
	assert.True(t, model.Done())
	assert.Nil(t, model.Err())
}

func TestSpinnerFailed(t *testing.T) {
	cfg := SpinnerConfig{Message: "Processing..."}
	model := NewSpinner(cfg)

	testErr := errRequired
	model, _ = model.Update(SpinnerDoneMsg{Err: testErr})
	assert.True(t, model.Done())
	assert.Equal(t, testErr, model.Err())
}

func TestTextInputUpdate(t *testing.T) {
	cfg := TextInputConfig{
		Label:       "Name",
		Placeholder: "Enter name",
	}
	model := NewTextInput(cfg)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.True(t, model.Submitted())
}

func TestTextInputEscapeGoBack(t *testing.T) {
	cfg := TextInputConfig{
		Label: "Name",
	}
	model := NewTextInput(cfg)

	assert.False(t, model.GoBack())

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})

	assert.True(t, model.GoBack())
}

func TestTextInputView(t *testing.T) {
	cfg := TextInputConfig{
		Label:    "Username",
		HelpText: "Enter your username",
		Required: true,
	}
	model := NewTextInput(cfg)

	view := model.View()

	assert.Contains(t, view, "Username")
}

func TestTextInputInit(t *testing.T) {
	cfg := TextInputConfig{
		Label: "Test",
	}
	model := NewTextInput(cfg)

	cmd := model.Init()

	assert.NotNil(t, cmd)
}

func TestTextInputReset(t *testing.T) {
	cfg := TextInputConfig{
		Label: "Test",
	}
	model := NewTextInput(cfg)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.True(t, model.Submitted())

	model.Reset()

	assert.False(t, model.Submitted())
	assert.False(t, model.GoBack())
}

func TestTextInputFocusBlur(t *testing.T) {
	cfg := TextInputConfig{
		Label: "Test",
	}
	model := NewTextInput(cfg)

	model.Blur()
	assert.False(t, model.Focused())

	model.Focus()
	assert.True(t, model.Focused())
}

func TestTextInputApplyDynamicDefault(t *testing.T) {
	cfg := TextInputConfig{
		Label: "Package",
		DynamicDefault: func(values map[string]any) string {
			group, _ := values["group"].(string)
			artifact, _ := values["artifact"].(string)
			return group + "." + artifact
		},
	}
	model := NewTextInput(cfg)

	values := map[string]any{
		"group":    "com.example",
		"artifact": "demo",
	}
	model.ApplyDynamicDefault(values)

	assert.Equal(t, "com.example.demo", model.Value())
}

func TestMultiSelectInit(t *testing.T) {
	cfg := MultiSelectConfig{
		Items: []MultiSelectItem{
			{Label: "A", Value: "a"},
		},
	}
	model := NewMultiSelect(cfg)

	cmd := model.Init()

	assert.Nil(t, cmd)
}

func TestMultiSelectView(t *testing.T) {
	cfg := MultiSelectConfig{
		Label: "Select Options",
		Items: []MultiSelectItem{
			{Label: "Option A", Value: "a"},
			{Label: "Option B", Value: "b"},
		},
	}
	model := NewMultiSelect(cfg)

	view := model.View()

	assert.Contains(t, view, "Select Options")
	assert.Contains(t, view, "Option A")
	assert.Contains(t, view, "Option B")
}

func TestMultiSelectEscapeGoBack(t *testing.T) {
	cfg := MultiSelectConfig{
		Items: []MultiSelectItem{
			{Label: "A", Value: "a"},
		},
	}
	model := NewMultiSelect(cfg)

	assert.False(t, model.GoBack())

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})

	assert.True(t, model.GoBack())
}

func TestMultiSelectReset(t *testing.T) {
	cfg := MultiSelectConfig{
		Items: []MultiSelectItem{
			{Label: "A", Value: "a"},
		},
	}
	model := NewMultiSelect(cfg)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.True(t, model.Submitted())

	model.Reset()

	assert.False(t, model.Submitted())
	assert.False(t, model.GoBack())
}

func TestMultiSelectVimKeys(t *testing.T) {
	cfg := MultiSelectConfig{
		Items: []MultiSelectItem{
			{Label: "A", Value: "a"},
			{Label: "B", Value: "b"},
			{Label: "C", Value: "c"},
		},
	}
	model := NewMultiSelect(cfg)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	assert.Equal(t, 1, model.cursor)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	assert.Equal(t, 0, model.cursor)
}

func TestMultiSelectSelectAll(t *testing.T) {
	cfg := MultiSelectConfig{
		Items: []MultiSelectItem{
			{Label: "A", Value: "a"},
			{Label: "B", Value: "b"},
		},
	}
	model := NewMultiSelect(cfg)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

	values := model.Values()
	assert.Len(t, values, 2)
}

func TestMultiSelectSelectNone(t *testing.T) {
	cfg := MultiSelectConfig{
		Items: []MultiSelectItem{
			{Label: "A", Value: "a", Selected: true},
			{Label: "B", Value: "b", Selected: true},
		},
	}
	model := NewMultiSelect(cfg)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

	values := model.Values()
	assert.Len(t, values, 0)
}

func TestMultiSelectSelectedItems(t *testing.T) {
	cfg := MultiSelectConfig{
		Items: []MultiSelectItem{
			{Label: "A", Value: "a"},
			{Label: "B", Value: "b"},
		},
	}
	model := NewMultiSelect(cfg)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeySpace})

	items := model.SelectedItems()
	assert.Len(t, items, 1)
	assert.Equal(t, "A", items[0].Label)
}

func TestSpinnerView(t *testing.T) {
	cfg := SpinnerConfig{
		Message: "Loading data...",
	}
	model := NewSpinner(cfg)

	view := model.View()

	assert.Contains(t, view, "Loading data...")
}

func TestSpinnerInit(t *testing.T) {
	cfg := SpinnerConfig{
		Message: "Loading...",
	}
	model := NewSpinner(cfg)

	cmd := model.Init()

	assert.NotNil(t, cmd)
}

func TestSpinnerSetMessage(t *testing.T) {
	cfg := SpinnerConfig{
		Message: "Initial",
	}
	model := NewSpinner(cfg)

	model.SetMessage("Updated message")

	view := model.View()
	assert.Contains(t, view, "Updated message")
}

func TestSpinnerComplete(t *testing.T) {
	msg := SpinnerComplete()

	doneMsg, ok := msg.(SpinnerDoneMsg)
	assert.True(t, ok)
	assert.Nil(t, doneMsg.Err)
}

func TestSpinnerFailedFunc(t *testing.T) {
	testErr := errRequired
	msg := SpinnerFailed(testErr)

	doneMsg, ok := msg.(SpinnerDoneMsg)
	assert.True(t, ok)
	assert.Equal(t, testErr, doneMsg.Err)
}
