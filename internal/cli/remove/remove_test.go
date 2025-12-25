package remove

import (
	"testing"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	cmd := NewCommand()

	assert.Equal(t, "remove [dependency...]", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
	assert.Contains(t, cmd.Aliases, "rm")
	assert.NotNil(t, cmd.RunE)
}

func TestResolveInput(t *testing.T) {
	project := &buildtool.Project{
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.projectlombok", ArtifactId: "lombok"},
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-web"},
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-jpa"},
		},
	}

	tests := []struct {
		name          string
		input         string
		expectedGroup string
		expectedArtif string
	}{
		{
			name:          "full coordinates",
			input:         "org.projectlombok:lombok",
			expectedGroup: "org.projectlombok",
			expectedArtif: "lombok",
		},
		{
			name:          "artifact only exact match",
			input:         "lombok",
			expectedGroup: "org.projectlombok",
			expectedArtif: "lombok",
		},
		{
			name:          "artifact only suffix match",
			input:         "spring-boot-starter-web",
			expectedGroup: "org.springframework.boot",
			expectedArtif: "spring-boot-starter-web",
		},
		{
			name:          "not found",
			input:         "nonexistent",
			expectedGroup: "",
			expectedArtif: "",
		},
		{
			name:          "partial suffix match",
			input:         "jpa",
			expectedGroup: "org.springframework.boot",
			expectedArtif: "spring-boot-starter-jpa",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groupId, artifactId := resolveInput(tt.input, project)
			assert.Equal(t, tt.expectedGroup, groupId)
			assert.Equal(t, tt.expectedArtif, artifactId)
		})
	}
}

func TestRemovePickerModel(t *testing.T) {
	deps := []buildtool.Dependency{
		{GroupId: "org.projectlombok", ArtifactId: "lombok"},
		{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-web"},
	}

	model := newRemovePickerModel(deps)

	assert.Len(t, model.deps, 2)
	assert.Len(t, model.filtered, 2)
	assert.Equal(t, 0, model.countSelected())
	assert.False(t, model.submitted)
	assert.False(t, model.cancelled)
}

func TestRemovePickerToggle(t *testing.T) {
	deps := []buildtool.Dependency{
		{GroupId: "org.projectlombok", ArtifactId: "lombok"},
	}

	model := newRemovePickerModel(deps)

	assert.Equal(t, 0, model.countSelected())

	model.toggleCurrent()
	assert.Equal(t, 1, model.countSelected())

	model.toggleCurrent()
	assert.Equal(t, 0, model.countSelected())
}

func TestRemovePickerSelectAll(t *testing.T) {
	deps := []buildtool.Dependency{
		{GroupId: "org.projectlombok", ArtifactId: "lombok"},
		{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-web"},
	}

	model := newRemovePickerModel(deps)

	model.selectAllVisible()
	assert.Equal(t, 2, model.countSelected())
}

func TestRemovePickerSelectNone(t *testing.T) {
	deps := []buildtool.Dependency{
		{GroupId: "org.projectlombok", ArtifactId: "lombok"},
	}

	model := newRemovePickerModel(deps)

	model.selectAllVisible()
	assert.Equal(t, 1, model.countSelected())

	model.selectNone()
	assert.Equal(t, 0, model.countSelected())
}

func TestRemovePickerFilter(t *testing.T) {
	deps := []buildtool.Dependency{
		{GroupId: "org.projectlombok", ArtifactId: "lombok"},
		{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-web"},
	}

	model := newRemovePickerModel(deps)

	model.searchQuery = "lombok"
	model.applyFilter()
	assert.Len(t, model.filtered, 1)

	model.searchQuery = ""
	model.applyFilter()
	assert.Len(t, model.filtered, 2)
}

func TestRemovePickerSelectedDeps(t *testing.T) {
	deps := []buildtool.Dependency{
		{GroupId: "org.projectlombok", ArtifactId: "lombok"},
		{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-web"},
	}

	model := newRemovePickerModel(deps)

	model.toggleCurrent()
	selected := model.selectedDeps()
	assert.Len(t, selected, 1)
	assert.Equal(t, "lombok", selected[0].ArtifactId)
}

func TestRemovePickerMoveCursor(t *testing.T) {
	deps := []buildtool.Dependency{
		{GroupId: "org.projectlombok", ArtifactId: "lombok"},
		{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-web"},
		{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-jpa"},
	}

	model := newRemovePickerModel(deps)

	assert.Equal(t, 0, model.cursor)

	model.moveCursor(1)
	assert.Equal(t, 1, model.cursor)

	model.moveCursor(1)
	assert.Equal(t, 2, model.cursor)

	model.moveCursor(1)
	assert.Equal(t, 2, model.cursor)

	model.moveCursor(-10)
	assert.Equal(t, 0, model.cursor)
}
