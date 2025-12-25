package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
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

func TestAllSubcommandsExist(t *testing.T) {
	cmd := NewCommand()
	commands := cmd.Commands()

	expectedCommands := map[string]bool{
		"resource [name]":   false,
		"controller [name]": false,
		"service [name]":    false,
		"repository [name]": false,
		"entity [name]":     false,
		"dto [name]":        false,
	}

	for _, c := range commands {
		if _, ok := expectedCommands[c.Use]; ok {
			expectedCommands[c.Use] = true
		}
	}

	for name, found := range expectedCommands {
		assert.True(t, found, "%s subcommand should exist", name)
	}
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

func TestControllerCommandFlags(t *testing.T) {
	cmd := newControllerCommand()

	assert.Equal(t, "controller [name]", cmd.Use)
	assert.Contains(t, cmd.Aliases, "co")

	packageFlag := cmd.Flags().Lookup("package")
	noInteractiveFlag := cmd.Flags().Lookup("no-interactive")

	assert.NotNil(t, packageFlag)
	assert.Equal(t, "p", packageFlag.Shorthand)
	assert.NotNil(t, noInteractiveFlag)
}

func TestServiceCommandFlags(t *testing.T) {
	cmd := newServiceCommand()

	assert.Equal(t, "service [name]", cmd.Use)
	assert.Contains(t, cmd.Aliases, "s")

	packageFlag := cmd.Flags().Lookup("package")
	noInteractiveFlag := cmd.Flags().Lookup("no-interactive")

	assert.NotNil(t, packageFlag)
	assert.NotNil(t, noInteractiveFlag)
}

func TestRepositoryCommandFlags(t *testing.T) {
	cmd := newRepositoryCommand()

	assert.Equal(t, "repository [name]", cmd.Use)
	assert.Contains(t, cmd.Aliases, "repo")

	packageFlag := cmd.Flags().Lookup("package")
	noInteractiveFlag := cmd.Flags().Lookup("no-interactive")

	assert.NotNil(t, packageFlag)
	assert.NotNil(t, noInteractiveFlag)
}

func TestEntityCommandFlags(t *testing.T) {
	cmd := newEntityCommand()

	assert.Equal(t, "entity [name]", cmd.Use)
	assert.Contains(t, cmd.Aliases, "e")

	packageFlag := cmd.Flags().Lookup("package")
	noInteractiveFlag := cmd.Flags().Lookup("no-interactive")

	assert.NotNil(t, packageFlag)
	assert.NotNil(t, noInteractiveFlag)
}

func TestDtoCommandFlags(t *testing.T) {
	cmd := newDtoCommand()

	assert.Equal(t, "dto [name]", cmd.Use)

	packageFlag := cmd.Flags().Lookup("package")
	noInteractiveFlag := cmd.Flags().Lookup("no-interactive")
	requestOnlyFlag := cmd.Flags().Lookup("request-only")
	responseOnlyFlag := cmd.Flags().Lookup("response-only")

	assert.NotNil(t, packageFlag)
	assert.NotNil(t, noInteractiveFlag)
	assert.NotNil(t, requestOnlyFlag)
	assert.NotNil(t, responseOnlyFlag)
}

func TestValidateComponentName(t *testing.T) {
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
			err := ValidateComponentName(tt.input)
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
			result := ToPascalCase(tt.input)
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
			result := ToCamelCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildTemplateData(t *testing.T) {
	cfg := ComponentConfig{
		Name:          "User",
		BasePackage:   "com.example.demo",
		HasLombok:     true,
		HasJpa:        true,
		HasValidation: false,
	}

	data := BuildTemplateData(cfg)

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

func TestValidateComponentConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     ComponentConfig
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid config",
			cfg:     ComponentConfig{Name: "User", BasePackage: "com.example"},
			wantErr: false,
		},
		{
			name:    "missing name",
			cfg:     ComponentConfig{BasePackage: "com.example"},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name:    "missing package",
			cfg:     ComponentConfig{Name: "User"},
			wantErr: true,
			errMsg:  "base package is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateComponentConfig(tt.cfg)
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

	result := FindSourcePath(tempDir)

	assert.Empty(t, result)
}

func TestFindSourcePathRealFS(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "haft-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	srcPath := filepath.Join(tempDir, "src", "main", "java")
	require.NoError(t, os.MkdirAll(srcPath, 0755))

	result := FindSourcePath(tempDir)

	assert.Equal(t, srcPath, result)
}

func TestFindSourcePathNotFound(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "haft-test-empty")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	result := FindSourcePath(tempDir)

	assert.Empty(t, result)
}

func TestFormatRelativePath(t *testing.T) {
	base := "/home/user/project"
	path := "/home/user/project/src/main/java/User.java"

	result := FormatRelativePath(base, path)

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
			err := ValidatePackageName(tt.input)
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
			assert.Equal(t, tt.expected, Capitalize(tt.input))
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
			result := SplitWords(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsUpper(t *testing.T) {
	assert.True(t, IsUpper('A'))
	assert.True(t, IsUpper('Z'))
	assert.False(t, IsUpper('a'))
	assert.False(t, IsUpper('1'))
}

func TestBuildComponentWizardSteps(t *testing.T) {
	cfg := ComponentConfig{
		Name:        "User",
		BasePackage: "com.example.demo",
	}

	steps, keys := buildComponentWizardSteps(cfg, "Controller")

	assert.Len(t, steps, 2)
	assert.Len(t, keys, 2)
	assert.Equal(t, "name", keys[0])
	assert.Equal(t, "basePackage", keys[1])
}

func TestBuildComponentWizardStepsWithEmptyConfig(t *testing.T) {
	cfg := ComponentConfig{}

	steps, keys := buildComponentWizardSteps(cfg, "Service")

	assert.Len(t, steps, 2)
	assert.Equal(t, []string{"name", "basePackage"}, keys)
}

func TestComponentConfigPreservesDetectedFeatures(t *testing.T) {
	cfg := ComponentConfig{
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

func TestComponentConfigAutoDetection(t *testing.T) {
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
			cfg := ComponentConfig{
				Name:          "User",
				BasePackage:   "com.example",
				HasLombok:     tt.hasLombok,
				HasJpa:        tt.hasJpa,
				HasValidation: tt.hasValidation,
			}

			data := BuildTemplateData(cfg)

			assert.Equal(t, tt.hasLombok, data["HasLombok"])
			assert.Equal(t, tt.hasJpa, data["HasJpa"])
			assert.Equal(t, tt.hasValidation, data["HasValidation"])
		})
	}
}

func TestBuildTemplateDataAllCombinations(t *testing.T) {
	tests := []struct {
		name string
		cfg  ComponentConfig
	}{
		{
			name: "with lombok only",
			cfg: ComponentConfig{
				Name:          "Product",
				BasePackage:   "com.shop",
				HasLombok:     true,
				HasJpa:        false,
				HasValidation: false,
			},
		},
		{
			name: "with jpa only",
			cfg: ComponentConfig{
				Name:          "Order",
				BasePackage:   "com.shop",
				HasLombok:     false,
				HasJpa:        true,
				HasValidation: false,
			},
		},
		{
			name: "with validation only",
			cfg: ComponentConfig{
				Name:          "Customer",
				BasePackage:   "com.shop",
				HasLombok:     false,
				HasJpa:        false,
				HasValidation: true,
			},
		},
		{
			name: "with all features",
			cfg: ComponentConfig{
				Name:          "Invoice",
				BasePackage:   "com.billing",
				HasLombok:     true,
				HasJpa:        true,
				HasValidation: true,
			},
		},
		{
			name: "with no features",
			cfg: ComponentConfig{
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
			data := BuildTemplateData(tt.cfg)

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

	result := FindSourcePath(tempDir)

	assert.Equal(t, srcPath, result)
}

func TestFormatRelativePathError(t *testing.T) {
	result := FormatRelativePath("/different/base", "/some/other/path/file.java")

	assert.NotEmpty(t, result)
}

func TestToCamelCaseEmpty(t *testing.T) {
	result := ToCamelCase("")
	assert.Equal(t, "", result)
}

func TestSplitWordsEmpty(t *testing.T) {
	result := SplitWords("")
	assert.Nil(t, result)
}

func TestSplitWordsWithSpaces(t *testing.T) {
	result := SplitWords("hello world test")
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

func TestComponentConfigStruct(t *testing.T) {
	cfg := ComponentConfig{
		Name:          "TestComponent",
		BasePackage:   "com.test.pkg",
		HasLombok:     true,
		HasJpa:        true,
		HasValidation: true,
	}

	assert.Equal(t, "TestComponent", cfg.Name)
	assert.Equal(t, "com.test.pkg", cfg.BasePackage)
	assert.True(t, cfg.HasLombok)
	assert.True(t, cfg.HasJpa)
	assert.True(t, cfg.HasValidation)
}

func TestBuildResourceTemplateData(t *testing.T) {
	cfg := ResourceConfig{
		Name:          "Product",
		BasePackage:   "com.shop.api",
		HasLombok:     true,
		HasJpa:        true,
		HasValidation: true,
	}

	data := buildResourceTemplateData(cfg)

	assert.Equal(t, "Product", data["Name"])
	assert.Equal(t, "product", data["NameLower"])
	assert.Equal(t, "product", data["NameCamel"])
	assert.Equal(t, "com.shop.api", data["BasePackage"])
	assert.Equal(t, true, data["HasLombok"])
	assert.Equal(t, true, data["HasJpa"])
	assert.Equal(t, true, data["HasValidation"])
}

func TestBuildResourceTemplateDataWithMultiWordName(t *testing.T) {
	cfg := ResourceConfig{
		Name:          "UserProfile",
		BasePackage:   "com.example",
		HasLombok:     false,
		HasJpa:        false,
		HasValidation: false,
	}

	data := buildResourceTemplateData(cfg)

	assert.Equal(t, "UserProfile", data["Name"])
	assert.Equal(t, "userprofile", data["NameLower"])
	assert.Equal(t, "userProfile", data["NameCamel"])
}

func TestCommandAliases(t *testing.T) {
	tests := []struct {
		cmdFunc       func() *cobra.Command
		expectedAlias string
	}{
		{newResourceCommand, "r"},
		{newControllerCommand, "co"},
		{newServiceCommand, "s"},
		{newRepositoryCommand, "repo"},
		{newEntityCommand, "e"},
	}

	for _, tt := range tests {
		cmd := tt.cmdFunc()
		assert.Contains(t, cmd.Aliases, tt.expectedAlias, "Command %s should have alias %s", cmd.Use, tt.expectedAlias)
	}
}

func TestCommandExamples(t *testing.T) {
	commands := []func() *cobra.Command{
		newResourceCommand,
		newControllerCommand,
		newServiceCommand,
		newRepositoryCommand,
		newEntityCommand,
		newDtoCommand,
	}

	for _, cmdFunc := range commands {
		cmd := cmdFunc()
		assert.NotEmpty(t, cmd.Example, "Command %s should have examples", cmd.Use)
		assert.NotEmpty(t, cmd.Long, "Command %s should have long description", cmd.Use)
		assert.NotEmpty(t, cmd.Short, "Command %s should have short description", cmd.Use)
	}
}

func TestValidateComponentNameEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"exactly 2 chars", "Ab", false},
		{"long name", "VeryLongResourceNameWithManyCharacters", false},
		{"only numbers after letter", "A123456", false},
		{"empty string", "", true},
		{"special chars", "User@Name", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateComponentName(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePackageNameEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"single segment", "com", false},
		{"with numbers", "com.example123", false},
		{"three segments", "org.apache.commons", false},
		{"double dots", "com..example", true},
		{"uppercase in middle", "com.Example.demo", true},
		{"starts with number", "1com.example", true},
		{"segment starts with number", "com.1example", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePackageName(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestToPascalCaseEdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"a", "A"},
		{"ABC", "Abc"},
		{"my-super-long-name", "MySuperLongName"},
		{"my_super_long_name", "MySuperLongName"},
		{"mySuper_long-Name", "MySuperLongName"},
		{"123", "123"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToPascalCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToCamelCaseEdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"A", "a"},
		{"ABC", "abc"},
		{"MySuper_long-Name", "mySuperLongName"},
		{"XMLParser", "xmlparser"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToCamelCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSplitWordsEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"multiple hyphens", "a-b-c-d", []string{"a", "b", "c", "d"}},
		{"multiple underscores", "a_b_c_d", []string{"a", "b", "c", "d"}},
		{"mixed separators", "a-b_c", []string{"a", "b", "c"}},
		{"consecutive uppercase", "XMLParser", []string{"XMLParser"}},
		{"single uppercase", "A", []string{"A"}},
		{"trailing separator", "user-", []string{"user"}},
		{"leading separator", "-user", []string{"user"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitWords(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildComponentWizardStepsAllTypes(t *testing.T) {
	componentTypes := []string{"Controller", "Service", "Repository", "Entity", "DTO", "Resource"}

	for _, compType := range componentTypes {
		t.Run(compType, func(t *testing.T) {
			cfg := ComponentConfig{Name: "Test", BasePackage: "com.test"}
			steps, keys := buildComponentWizardSteps(cfg, compType)

			assert.Len(t, steps, 2)
			assert.Len(t, keys, 2)
			assert.Equal(t, "name", keys[0])
			assert.Equal(t, "basePackage", keys[1])
		})
	}
}

func TestFormatRelativePathSamePath(t *testing.T) {
	base := "/home/user/project"
	path := "/home/user/project"

	result := FormatRelativePath(base, path)

	assert.Equal(t, ".", result)
}

func TestFormatRelativePathChildPath(t *testing.T) {
	base := "/home/user/project"
	path := "/home/user/project/src"

	result := FormatRelativePath(base, path)

	assert.Equal(t, "src", result)
}

func TestIsUpperEdgeCases(t *testing.T) {
	tests := []struct {
		input    rune
		expected bool
	}{
		{'A', true},
		{'Z', true},
		{'M', true},
		{'a', false},
		{'z', false},
		{'0', false},
		{'9', false},
		{' ', false},
		{'@', false},
		{'[', false},
		{'`', false},
	}

	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			assert.Equal(t, tt.expected, IsUpper(tt.input))
		})
	}
}

func TestResourceConfigToComponentConfig(t *testing.T) {
	resCfg := ResourceConfig{
		Name:          "User",
		BasePackage:   "com.example",
		HasLombok:     true,
		HasJpa:        true,
		HasValidation: true,
	}

	compCfg := ComponentConfig{
		Name:          resCfg.Name,
		BasePackage:   resCfg.BasePackage,
		HasLombok:     resCfg.HasLombok,
		HasJpa:        resCfg.HasJpa,
		HasValidation: resCfg.HasValidation,
	}

	assert.Equal(t, resCfg.Name, compCfg.Name)
	assert.Equal(t, resCfg.BasePackage, compCfg.BasePackage)
	assert.Equal(t, resCfg.HasLombok, compCfg.HasLombok)
	assert.Equal(t, resCfg.HasJpa, compCfg.HasJpa)
	assert.Equal(t, resCfg.HasValidation, compCfg.HasValidation)
}

func TestSubcommandCount(t *testing.T) {
	cmd := NewCommand()
	assert.Equal(t, 6, len(cmd.Commands()), "Should have 6 subcommands: resource, controller, service, repository, entity, dto")
}

func TestGenerateCommandHasNoRunE(t *testing.T) {
	cmd := NewCommand()
	assert.Nil(t, cmd.RunE, "Parent generate command should not have RunE")
	assert.Nil(t, cmd.Run, "Parent generate command should not have Run")
}
