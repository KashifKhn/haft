package init

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToArtifactId(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"my-app", "my-app"},
		{"MyApp", "myapp"},
		{"My App", "my-app"},
		{"my_app", "my-app"},
		{"My--App", "my-app"},
		{"  my app  ", "my-app"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toArtifactId(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"my-app", "MyApp"},
		{"user-profile", "UserProfile"},
		{"hello_world", "HelloWorld"},
		{"simple", "Simple"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toPascalCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateProjectName(t *testing.T) {
	validNames := []string{"myapp", "my-app", "my_app", "MyApp123"}
	for _, name := range validNames {
		t.Run("valid_"+name, func(t *testing.T) {
			err := validateProjectName(name)
			assert.NoError(t, err)
		})
	}

	invalidNames := []string{"a", "123app", "-myapp", "_myapp"}
	for _, name := range invalidNames {
		t.Run("invalid_"+name, func(t *testing.T) {
			err := validateProjectName(name)
			assert.Error(t, err)
		})
	}
}

func TestValidateGroupId(t *testing.T) {
	validIds := []string{"com.example", "org.springframework", "io.github.user"}
	for _, id := range validIds {
		t.Run("valid_"+id, func(t *testing.T) {
			err := validateGroupId(id)
			assert.NoError(t, err)
		})
	}

	invalidIds := []string{"Com.Example", "com..example", ".com", "com-example"}
	for _, id := range invalidIds {
		t.Run("invalid_"+id, func(t *testing.T) {
			err := validateGroupId(id)
			assert.Error(t, err)
		})
	}
}

func TestValidateArtifactId(t *testing.T) {
	validIds := []string{"myapp", "my-app", "demo123"}
	for _, id := range validIds {
		t.Run("valid_"+id, func(t *testing.T) {
			err := validateArtifactId(id)
			assert.NoError(t, err)
		})
	}

	invalidIds := []string{"a", "MyApp", "my_app", "-myapp"}
	for _, id := range invalidIds {
		t.Run("invalid_"+id, func(t *testing.T) {
			err := validateArtifactId(id)
			assert.Error(t, err)
		})
	}
}

func TestContains(t *testing.T) {
	slice := []string{"a", "b", "c"}

	assert.True(t, contains(slice, "a"))
	assert.True(t, contains(slice, "b"))
	assert.False(t, contains(slice, "d"))
	assert.False(t, contains(nil, "a"))
}

func TestBuildDependencies(t *testing.T) {
	deps := buildDependencies([]string{"web", "data-jpa", "lombok"})

	assert.Len(t, deps, 4)

	var hasWeb, hasJpa, hasLombok, hasH2 bool
	for _, d := range deps {
		if d.ArtifactId == "spring-boot-starter-web" {
			hasWeb = true
		}
		if d.ArtifactId == "spring-boot-starter-data-jpa" {
			hasJpa = true
		}
		if d.ArtifactId == "lombok" {
			hasLombok = true
			assert.Equal(t, "provided", d.Scope)
		}
		if d.ArtifactId == "h2" {
			hasH2 = true
			assert.Equal(t, "runtime", d.Scope)
		}
	}

	assert.True(t, hasWeb)
	assert.True(t, hasJpa)
	assert.True(t, hasLombok)
	assert.True(t, hasH2)
}

func TestBuildDependenciesWithExplicitDb(t *testing.T) {
	deps := buildDependencies([]string{"data-jpa", "postgresql"})

	assert.Len(t, deps, 2)

	var hasJpa, hasPostgres, hasH2 bool
	for _, d := range deps {
		if d.ArtifactId == "spring-boot-starter-data-jpa" {
			hasJpa = true
		}
		if d.ArtifactId == "postgresql" {
			hasPostgres = true
		}
		if d.ArtifactId == "h2" {
			hasH2 = true
		}
	}

	assert.True(t, hasJpa)
	assert.True(t, hasPostgres)
	assert.False(t, hasH2)
}

func TestValidateConfig(t *testing.T) {
	valid := ProjectConfig{
		Name:       "test",
		GroupId:    "com.example",
		ArtifactId: "test",
	}
	assert.NoError(t, validateConfig(valid))

	noName := ProjectConfig{GroupId: "com.example", ArtifactId: "test"}
	assert.Error(t, validateConfig(noName))

	noGroup := ProjectConfig{Name: "test", ArtifactId: "test"}
	assert.Error(t, validateConfig(noGroup))

	noArtifact := ProjectConfig{Name: "test", GroupId: "com.example"}
	assert.Error(t, validateConfig(noArtifact))
}

func TestNormalizeDependencies(t *testing.T) {
	depsWithJpaNoDb := []string{"web", "data-jpa"}
	normalized := normalizeDependencies(depsWithJpaNoDb)
	assert.Contains(t, normalized, "h2")

	depsWithJpaAndDb := []string{"web", "data-jpa", "postgresql"}
	normalized = normalizeDependencies(depsWithJpaAndDb)
	assert.NotContains(t, normalized, "h2")
}

func TestApplyDefaultsConfigFormat(t *testing.T) {
	cfg := &ProjectConfig{
		Name:       "test",
		GroupId:    "com.example",
		ArtifactId: "test",
	}

	err := applyDefaults(cfg)

	assert.NoError(t, err)
	assert.Equal(t, "yaml", cfg.ConfigFormat)
}

func TestBuildWizardStepsConfigFormatDefault(t *testing.T) {
	cfg := ProjectConfig{}

	steps, keys := buildWizardSteps(cfg)

	configFormatIdx := -1
	for i, key := range keys {
		if key == "configFormat" {
			configFormatIdx = i
			break
		}
	}

	assert.NotEqual(t, -1, configFormatIdx, "configFormat step should exist")
	assert.True(t, configFormatIdx < len(steps), "configFormat step index should be valid")
}

func TestApplyDefaultsAllFields(t *testing.T) {
	cfg := &ProjectConfig{
		Name:       "myapp",
		GroupId:    "com.example",
		ArtifactId: "myapp",
	}

	err := applyDefaults(cfg)

	assert.NoError(t, err)
	assert.Equal(t, "21", cfg.JavaVersion)
	assert.Equal(t, "3.4.0", cfg.SpringBootVersion)
	assert.Equal(t, "maven", cfg.BuildTool)
	assert.Equal(t, "jar", cfg.Packaging)
	assert.Equal(t, "yaml", cfg.ConfigFormat)
	assert.Equal(t, "com.example.myapp", cfg.PackageName)
	assert.Contains(t, cfg.Description, "myapp")
}

func TestApplyDefaultsPreservesExistingValues(t *testing.T) {
	cfg := &ProjectConfig{
		Name:              "myapp",
		GroupId:           "com.example",
		ArtifactId:        "myapp",
		JavaVersion:       "17",
		SpringBootVersion: "3.3.0",
		BuildTool:         "gradle",
		Packaging:         "war",
		ConfigFormat:      "properties",
		PackageName:       "com.custom.pkg",
		Description:       "Custom description",
	}

	err := applyDefaults(cfg)

	assert.NoError(t, err)
	assert.Equal(t, "17", cfg.JavaVersion)
	assert.Equal(t, "3.3.0", cfg.SpringBootVersion)
	assert.Equal(t, "gradle", cfg.BuildTool)
	assert.Equal(t, "war", cfg.Packaging)
	assert.Equal(t, "properties", cfg.ConfigFormat)
	assert.Equal(t, "com.custom.pkg", cfg.PackageName)
	assert.Equal(t, "Custom description", cfg.Description)
}

func TestValidatePackageName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid com.example", "com.example", false},
		{"valid com.example.demo", "com.example.demo", false},
		{"empty is valid", "", false},
		{"invalid uppercase", "Com.example", true},
		{"invalid hyphen", "com-example", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePackageName(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
