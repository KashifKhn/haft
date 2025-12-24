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
