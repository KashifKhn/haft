package add

import (
	"testing"

	"github.com/KashifKhn/haft/internal/maven"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCatalogEntry(t *testing.T) {
	tests := []struct {
		alias    string
		expected bool
	}{
		{"lombok", true},
		{"jpa", true},
		{"web", true},
		{"validation", true},
		{"postgresql", true},
		{"unknown", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.alias, func(t *testing.T) {
			entry, ok := GetCatalogEntry(tt.alias)
			assert.Equal(t, tt.expected, ok)
			if ok {
				assert.NotEmpty(t, entry.Name)
				assert.NotEmpty(t, entry.Description)
				assert.NotEmpty(t, entry.Category)
				assert.NotEmpty(t, entry.Dependencies)
			}
		})
	}
}

func TestGetAllAliases(t *testing.T) {
	aliases := GetAllAliases()
	assert.NotEmpty(t, aliases)
	assert.Contains(t, aliases, "lombok")
	assert.Contains(t, aliases, "jpa")
	assert.Contains(t, aliases, "web")
}

func TestGetCatalogByCategory(t *testing.T) {
	categories := GetCatalogByCategory()
	assert.NotEmpty(t, categories)
	assert.Contains(t, categories, "Web")
	assert.Contains(t, categories, "SQL")
	assert.Contains(t, categories, "Developer Tools")
}

func TestSearchCatalog(t *testing.T) {
	tests := []struct {
		query       string
		expectFound bool
	}{
		{"lombok", true},
		{"database", false},
		{"JPA", true},
		{"web", true},
		{"nonexistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			results := SearchCatalog(tt.query)
			if tt.expectFound {
				assert.NotEmpty(t, results)
			}
		})
	}
}

func TestResolveDependency(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantDeps  int
		wantName  string
		wantError bool
	}{
		{
			name:     "shortcut lombok",
			input:    "lombok",
			wantDeps: 1,
			wantName: "Lombok",
		},
		{
			name:     "shortcut mapstruct with multiple deps",
			input:    "mapstruct",
			wantDeps: 2,
			wantName: "MapStruct",
		},
		{
			name:      "invalid shortcut",
			input:     "unknownshortcut",
			wantError: true,
		},
		{
			name:      "invalid format too many colons",
			input:     "org:example:my:lib",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps, name, err := resolveDependency(tt.input)
			if tt.wantError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Len(t, deps, tt.wantDeps)
			assert.Equal(t, tt.wantName, name)
		})
	}
}

func TestFormatDependency(t *testing.T) {
	tests := []struct {
		dep      maven.Dependency
		expected string
	}{
		{
			dep:      maven.Dependency{GroupId: "org.example", ArtifactId: "my-lib"},
			expected: "org.example:my-lib",
		},
		{
			dep:      maven.Dependency{GroupId: "org.example", ArtifactId: "my-lib", Version: "1.0.0"},
			expected: "org.example:my-lib:1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatDependency(tt.dep)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCatalogEntriesHaveRequiredFields(t *testing.T) {
	for alias, entry := range dependencyCatalog {
		t.Run(alias, func(t *testing.T) {
			assert.NotEmpty(t, entry.Name, "Name should not be empty for %s", alias)
			assert.NotEmpty(t, entry.Description, "Description should not be empty for %s", alias)
			assert.NotEmpty(t, entry.Category, "Category should not be empty for %s", alias)
			assert.NotEmpty(t, entry.Dependencies, "Dependencies should not be empty for %s", alias)

			for _, dep := range entry.Dependencies {
				assert.NotEmpty(t, dep.GroupId, "GroupId should not be empty for %s", alias)
				assert.NotEmpty(t, dep.ArtifactId, "ArtifactId should not be empty for %s", alias)
			}
		})
	}
}

func TestLombokDependencyHasCorrectScope(t *testing.T) {
	entry, ok := GetCatalogEntry("lombok")
	require.True(t, ok)
	require.Len(t, entry.Dependencies, 1)
	assert.Equal(t, "provided", entry.Dependencies[0].Scope)
}

func TestDevtoolsHasRuntimeScope(t *testing.T) {
	entry, ok := GetCatalogEntry("devtools")
	require.True(t, ok)
	require.Len(t, entry.Dependencies, 1)
	assert.Equal(t, "runtime", entry.Dependencies[0].Scope)
	assert.Equal(t, "true", entry.Dependencies[0].Optional)
}

func TestDatabaseDriversHaveRuntimeScope(t *testing.T) {
	drivers := []string{"postgresql", "mysql", "mariadb", "h2"}

	for _, driver := range drivers {
		t.Run(driver, func(t *testing.T) {
			entry, ok := GetCatalogEntry(driver)
			require.True(t, ok)
			require.NotEmpty(t, entry.Dependencies)
			assert.Equal(t, "runtime", entry.Dependencies[0].Scope)
		})
	}
}

func TestToLower(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HELLO", "hello"},
		{"Hello World", "hello world"},
		{"already lowercase", "already lowercase"},
		{"MixedCase123", "mixedcase123"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toLower(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		s        string
		substr   string
		expected bool
	}{
		{"hello world", "world", true},
		{"hello world", "WORLD", false},
		{"hello", "hello", true},
		{"hello", "hello world", false},
		{"", "", true},
		{"hello", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.s+"_"+tt.substr, func(t *testing.T) {
			result := contains(tt.s, tt.substr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewCommandHasRequiredFields(t *testing.T) {
	cmd := NewCommand()

	assert.Equal(t, "add [dependency...]", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
	assert.NotNil(t, cmd.RunE)
}

func TestNewCommandHasFlags(t *testing.T) {
	cmd := NewCommand()

	scopeFlag := cmd.Flag("scope")
	assert.NotNil(t, scopeFlag)
	assert.Equal(t, "string", scopeFlag.Value.Type())

	versionFlag := cmd.Flag("version")
	assert.NotNil(t, versionFlag)
	assert.Equal(t, "string", versionFlag.Value.Type())

	listFlag := cmd.Flag("list")
	assert.NotNil(t, listFlag)
	assert.Equal(t, "bool", listFlag.Value.Type())

	browseFlag := cmd.Flag("browse")
	assert.NotNil(t, browseFlag)
	assert.Equal(t, "bool", browseFlag.Value.Type())
}

func TestNewMavenClient(t *testing.T) {
	client := NewMavenClient()
	assert.NotNil(t, client)
	assert.NotNil(t, client.httpClient)
	assert.NotEmpty(t, client.baseURL)
}

func TestJWTDependencyHasMultipleArtifacts(t *testing.T) {
	entry, ok := GetCatalogEntry("jwt")
	require.True(t, ok)
	require.Len(t, entry.Dependencies, 3)

	assert.Equal(t, "io.jsonwebtoken", entry.Dependencies[0].GroupId)
	assert.Equal(t, "jjwt-api", entry.Dependencies[0].ArtifactId)

	assert.Equal(t, "runtime", entry.Dependencies[1].Scope)
	assert.Equal(t, "runtime", entry.Dependencies[2].Scope)
}

func TestNewCatalogEntries(t *testing.T) {
	newEntries := []string{
		"jwt", "commons-lang", "commons-io", "guava", "modelmapper",
		"jackson-datatype", "feign", "resilience4j", "micrometer",
		"cassandra", "neo4j", "freemarker", "mustache",
		"security-test", "mockito", "restdocs", "hateoas", "data-rest",
	}

	for _, alias := range newEntries {
		t.Run(alias, func(t *testing.T) {
			entry, ok := GetCatalogEntry(alias)
			assert.True(t, ok, "Entry should exist for %s", alias)
			assert.NotEmpty(t, entry.Name)
			assert.NotEmpty(t, entry.Dependencies)
		})
	}
}

func TestPickerModel(t *testing.T) {
	picker := NewPicker()

	assert.NotEmpty(t, picker.items)
	assert.NotEmpty(t, picker.filtered)
	assert.Equal(t, 15, picker.viewportSize)
	assert.False(t, picker.submitted)
	assert.False(t, picker.cancelled)
}

func TestPickerFilter(t *testing.T) {
	picker := NewPicker()
	originalCount := len(picker.filtered)

	picker.searchQuery = "lombok"
	picker.applyFilter()
	assert.Less(t, len(picker.filtered), originalCount)
	assert.Greater(t, len(picker.filtered), 0)

	picker.searchQuery = ""
	picker.applyFilter()
	assert.Equal(t, originalCount, len(picker.filtered))
}

func TestPickerSelectToggle(t *testing.T) {
	picker := NewPicker()

	assert.Equal(t, 0, picker.countSelected())

	picker.toggleCurrent()
	assert.Equal(t, 1, picker.countSelected())

	picker.toggleCurrent()
	assert.Equal(t, 0, picker.countSelected())
}

func TestPickerSelectAllVisible(t *testing.T) {
	picker := NewPicker()

	picker.selectAllVisible()
	assert.Equal(t, len(picker.items), picker.countSelected())
}

func TestPickerSelectNone(t *testing.T) {
	picker := NewPicker()

	picker.selectAllVisible()
	assert.Greater(t, picker.countSelected(), 0)

	picker.selectNone()
	assert.Equal(t, 0, picker.countSelected())
}

func TestPickerSelectedAliases(t *testing.T) {
	picker := NewPicker()

	aliases := picker.SelectedAliases()
	assert.Empty(t, aliases)

	picker.toggleCurrent()
	aliases = picker.SelectedAliases()
	assert.Len(t, aliases, 1)
}

func TestBuildBrowserCategories(t *testing.T) {
	categories := buildBrowserCategories()

	assert.NotEmpty(t, categories)

	foundWeb := false
	for _, cat := range categories {
		if cat.Name == "Web" {
			foundWeb = true
			assert.NotEmpty(t, cat.Dependencies)
		}
	}
	assert.True(t, foundWeb, "Web category should exist")
}

func TestCommandUseUpdatedForInteractive(t *testing.T) {
	cmd := NewCommand()
	assert.Equal(t, "add [dependency...]", cmd.Use)
}
