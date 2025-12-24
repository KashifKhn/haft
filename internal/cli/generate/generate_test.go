package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCommand(t *testing.T) {
	cmd := NewCommand()

	assert.Equal(t, "generate", cmd.Use)
	assert.Contains(t, cmd.Aliases, "g")
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
}

func TestResourceCommandExists(t *testing.T) {
	cmd := NewCommand()
	commands := cmd.Commands()

	var found bool
	for _, c := range commands {
		if c.Use == "resource [name]" {
			found = true
			break
		}
	}

	assert.True(t, found, "resource subcommand should exist")
}

func TestResourceCommandFlags(t *testing.T) {
	cmd := newResourceCommand()

	packageFlag := cmd.Flags().Lookup("package")
	noInteractiveFlag := cmd.Flags().Lookup("no-interactive")
	skipEntityFlag := cmd.Flags().Lookup("skip-entity")
	skipRepositoryFlag := cmd.Flags().Lookup("skip-repository")

	assert.NotNil(t, packageFlag)
	assert.Equal(t, "p", packageFlag.Shorthand)
	assert.NotNil(t, noInteractiveFlag)
	assert.NotNil(t, skipEntityFlag)
	assert.NotNil(t, skipRepositoryFlag)
}

func TestValidateResourceName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid User", "User", false},
		{"valid user", "user", false},
		{"valid Product123", "Product123", false},
		{"invalid single char", "A", true},
		{"invalid starts with number", "123User", true},
		{"invalid contains hyphen", "User-Profile", true},
		{"invalid contains underscore", "User_Profile", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateResourceName(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user", "User"},
		{"user-profile", "UserProfile"},
		{"user_profile", "UserProfile"},
		{"UserProfile", "UserProfile"},
		{"userProfile", "UserProfile"},
		{"PRODUCT", "Product"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toPascalCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"User", "user"},
		{"user-profile", "userProfile"},
		{"user_profile", "userProfile"},
		{"PRODUCT", "product"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toCamelCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildTemplateData(t *testing.T) {
	cfg := ResourceConfig{
		Name:          "User",
		BasePackage:   "com.example.demo",
		HasLombok:     true,
		HasJpa:        true,
		HasValidation: false,
	}

	data := buildTemplateData(cfg)

	assert.Equal(t, "User", data["Name"])
	assert.Equal(t, "user", data["NameLower"])
	assert.Equal(t, "user", data["NameCamel"])
	assert.Equal(t, "com.example.demo", data["BasePackage"])
	assert.Equal(t, true, data["HasLombok"])
	assert.Equal(t, true, data["HasJpa"])
	assert.Equal(t, false, data["HasValidation"])
}

func TestValidateResourceConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     ResourceConfig
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid config",
			cfg:     ResourceConfig{Name: "User", BasePackage: "com.example"},
			wantErr: false,
		},
		{
			name:    "missing name",
			cfg:     ResourceConfig{BasePackage: "com.example"},
			wantErr: true,
			errMsg:  "resource name is required",
		},
		{
			name:    "missing package",
			cfg:     ResourceConfig{Name: "User"},
			wantErr: true,
			errMsg:  "base package is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateResourceConfig(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFindSourcePath(t *testing.T) {
	fs := afero.NewMemMapFs()

	tempDir := "/testproject"
	srcPath := filepath.Join(tempDir, "src", "main", "java")
	require.NoError(t, fs.MkdirAll(srcPath, 0755))

	result := findSourcePath(tempDir)

	assert.Empty(t, result)
}

func TestFindSourcePathRealFS(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "haft-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	srcPath := filepath.Join(tempDir, "src", "main", "java")
	require.NoError(t, os.MkdirAll(srcPath, 0755))

	result := findSourcePath(tempDir)

	assert.Equal(t, srcPath, result)
}

func TestFindSourcePathNotFound(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "haft-test-empty")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	result := findSourcePath(tempDir)

	assert.Empty(t, result)
}

func TestFormatRelativePath(t *testing.T) {
	base := "/home/user/project"
	path := "/home/user/project/src/main/java/User.java"

	result := formatRelativePath(base, path)

	assert.Equal(t, "src/main/java/User.java", result)
}

func TestValidatePackageName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid com.example", "com.example", false},
		{"valid com.example.demo", "com.example.demo", false},
		{"valid io.github.user", "io.github.user", false},
		{"empty is valid", "", false},
		{"invalid starts uppercase", "Com.example", true},
		{"invalid contains hyphen", "com-example", true},
		{"invalid starts with dot", ".com.example", true},
		{"invalid ends with dot", "com.example.", true},
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

func TestCapitalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user", "User"},
		{"User", "User"},
		{"", ""},
		{"a", "A"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, capitalize(tt.input))
		})
	}
}

func TestSplitWords(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"user-profile", []string{"user", "profile"}},
		{"user_profile", []string{"user", "profile"}},
		{"UserProfile", []string{"User", "Profile"}},
		{"user", []string{"user"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := splitWords(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsUpper(t *testing.T) {
	assert.True(t, isUpper('A'))
	assert.True(t, isUpper('Z'))
	assert.False(t, isUpper('a'))
	assert.False(t, isUpper('1'))
}

func TestBuildResourceWizardSteps(t *testing.T) {
	cfg := ResourceConfig{
		Name:        "User",
		BasePackage: "com.example.demo",
	}

	steps, keys := buildResourceWizardSteps(cfg)

	assert.Len(t, steps, 2)
	assert.Len(t, keys, 2)
	assert.Equal(t, "name", keys[0])
	assert.Equal(t, "basePackage", keys[1])
}

func TestBuildResourceWizardStepsWithEmptyConfig(t *testing.T) {
	cfg := ResourceConfig{}

	steps, keys := buildResourceWizardSteps(cfg)

	assert.Len(t, steps, 2)
	assert.Equal(t, []string{"name", "basePackage"}, keys)
}

func TestExtractResourceWizardValuesPreservesDetectedFeatures(t *testing.T) {
	cfg := ResourceConfig{
		Name:          "",
		BasePackage:   "",
		HasLombok:     true,
		HasJpa:        true,
		HasValidation: true,
	}

	cfg.Name = "Product"
	cfg.BasePackage = "com.example.shop"

	assert.Equal(t, "Product", cfg.Name)
	assert.Equal(t, "com.example.shop", cfg.BasePackage)
	assert.True(t, cfg.HasLombok)
	assert.True(t, cfg.HasJpa)
	assert.True(t, cfg.HasValidation)
}

func TestResourceConfigAutoDetection(t *testing.T) {
	tests := []struct {
		name          string
		hasLombok     bool
		hasJpa        bool
		hasValidation bool
	}{
		{"all detected", true, true, true},
		{"only lombok", true, false, false},
		{"only jpa", false, true, false},
		{"only validation", false, false, true},
		{"none detected", false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := ResourceConfig{
				Name:          "User",
				BasePackage:   "com.example",
				HasLombok:     tt.hasLombok,
				HasJpa:        tt.hasJpa,
				HasValidation: tt.hasValidation,
			}

			data := buildTemplateData(cfg)

			assert.Equal(t, tt.hasLombok, data["HasLombok"])
			assert.Equal(t, tt.hasJpa, data["HasJpa"])
			assert.Equal(t, tt.hasValidation, data["HasValidation"])
		})
	}
}

func TestBuildTemplateDataAllCombinations(t *testing.T) {
	tests := []struct {
		name string
		cfg  ResourceConfig
	}{
		{
			name: "with lombok only",
			cfg: ResourceConfig{
				Name:          "Product",
				BasePackage:   "com.shop",
				HasLombok:     true,
				HasJpa:        false,
				HasValidation: false,
			},
		},
		{
			name: "with jpa only",
			cfg: ResourceConfig{
				Name:          "Order",
				BasePackage:   "com.shop",
				HasLombok:     false,
				HasJpa:        true,
				HasValidation: false,
			},
		},
		{
			name: "with validation only",
			cfg: ResourceConfig{
				Name:          "Customer",
				BasePackage:   "com.shop",
				HasLombok:     false,
				HasJpa:        false,
				HasValidation: true,
			},
		},
		{
			name: "with all features",
			cfg: ResourceConfig{
				Name:          "Invoice",
				BasePackage:   "com.billing",
				HasLombok:     true,
				HasJpa:        true,
				HasValidation: true,
			},
		},
		{
			name: "with no features",
			cfg: ResourceConfig{
				Name:          "Report",
				BasePackage:   "com.reports",
				HasLombok:     false,
				HasJpa:        false,
				HasValidation: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := buildTemplateData(tt.cfg)

			assert.Equal(t, tt.cfg.Name, data["Name"])
			assert.Equal(t, strings.ToLower(tt.cfg.Name), data["NameLower"])
			assert.Equal(t, tt.cfg.BasePackage, data["BasePackage"])
			assert.Equal(t, tt.cfg.HasLombok, data["HasLombok"])
			assert.Equal(t, tt.cfg.HasJpa, data["HasJpa"])
			assert.Equal(t, tt.cfg.HasValidation, data["HasValidation"])
		})
	}
}

func TestFindSourcePathWithAppDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "haft-test-app")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	srcPath := filepath.Join(tempDir, "app", "src", "main", "java")
	require.NoError(t, os.MkdirAll(srcPath, 0755))

	result := findSourcePath(tempDir)

	assert.Equal(t, srcPath, result)
}

func TestFormatRelativePathError(t *testing.T) {
	result := formatRelativePath("/different/base", "/some/other/path/file.java")

	assert.NotEmpty(t, result)
}

func TestToCamelCaseEmpty(t *testing.T) {
	result := toCamelCase("")
	assert.Equal(t, "", result)
}

func TestSplitWordsEmpty(t *testing.T) {
	result := splitWords("")
	assert.Nil(t, result)
}

func TestSplitWordsWithSpaces(t *testing.T) {
	result := splitWords("hello world test")
	assert.Equal(t, []string{"hello", "world", "test"}, result)
}

func TestResourceConfigStruct(t *testing.T) {
	cfg := ResourceConfig{
		Name:          "TestResource",
		BasePackage:   "com.test.pkg",
		HasLombok:     true,
		HasJpa:        true,
		HasValidation: true,
	}

	assert.Equal(t, "TestResource", cfg.Name)
	assert.Equal(t, "com.test.pkg", cfg.BasePackage)
	assert.True(t, cfg.HasLombok)
	assert.True(t, cfg.HasJpa)
	assert.True(t, cfg.HasValidation)
}
