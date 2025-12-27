package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/KashifKhn/haft/internal/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCommand(t *testing.T) {
	cmd := NewCommand()

	assert.Equal(t, "template", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
}

func TestNewCommandHasSubcommands(t *testing.T) {
	cmd := NewCommand()

	subcommands := cmd.Commands()
	assert.Len(t, subcommands, 3)

	subcommandNames := make([]string, len(subcommands))
	for i, sub := range subcommands {
		subcommandNames[i] = sub.Use
	}

	assert.Contains(t, subcommandNames, "init")
	assert.Contains(t, subcommandNames, "list")
	assert.Contains(t, subcommandNames, "validate [template-path]")
}

func TestNewInitCommand(t *testing.T) {
	cmd := newInitCommand()

	assert.Equal(t, "init", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
}

func TestNewInitCommandFlags(t *testing.T) {
	cmd := newInitCommand()

	categoryFlag := cmd.Flags().Lookup("category")
	assert.NotNil(t, categoryFlag)
	assert.Equal(t, "c", categoryFlag.Shorthand)

	globalFlag := cmd.Flags().Lookup("global")
	assert.NotNil(t, globalFlag)
	assert.Equal(t, "g", globalFlag.Shorthand)

	forceFlag := cmd.Flags().Lookup("force")
	assert.NotNil(t, forceFlag)
	assert.Equal(t, "f", forceFlag.Shorthand)
}

func TestNewListCommand(t *testing.T) {
	cmd := newListCommand()

	assert.Equal(t, "list", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
}

func TestNewListCommandFlags(t *testing.T) {
	cmd := newListCommand()

	customFlag := cmd.Flags().Lookup("custom")
	assert.NotNil(t, customFlag)

	categoryFlag := cmd.Flags().Lookup("category")
	assert.NotNil(t, categoryFlag)
	assert.Equal(t, "c", categoryFlag.Shorthand)

	pathsFlag := cmd.Flags().Lookup("paths")
	assert.NotNil(t, pathsFlag)
}

func TestNewValidateCommand(t *testing.T) {
	cmd := newValidateCommand()

	assert.Equal(t, "validate [template-path]", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
}

func TestNewValidateCommandFlags(t *testing.T) {
	cmd := newValidateCommand()

	varsFlag := cmd.Flags().Lookup("vars")
	assert.NotNil(t, varsFlag)

	conditionsFlag := cmd.Flags().Lookup("conditions")
	assert.NotNil(t, conditionsFlag)
}

func TestFilterByCategory(t *testing.T) {
	templates := []generator.TemplateInfo{
		{Name: "resource/Controller.java.tmpl", Source: generator.SourceEmbedded},
		{Name: "resource/Service.java.tmpl", Source: generator.SourceEmbedded},
		{Name: "test/ControllerTest.java.tmpl", Source: generator.SourceEmbedded},
		{Name: "project/Application.java.tmpl", Source: generator.SourceEmbedded},
	}

	filtered := filterByCategory(templates, "resource")
	assert.Len(t, filtered, 2)
	for _, tmpl := range filtered {
		assert.True(t, tmpl.Name[:9] == "resource/")
	}

	filtered = filterByCategory(templates, "test")
	assert.Len(t, filtered, 1)
	assert.Equal(t, "test/ControllerTest.java.tmpl", filtered[0].Name)

	filtered = filterByCategory(templates, "nonexistent")
	assert.Len(t, filtered, 0)
}

func TestFilterCustomOnly(t *testing.T) {
	templates := []generator.TemplateInfo{
		{Name: "resource/Controller.java.tmpl", Source: generator.SourceEmbedded},
		{Name: "resource/Service.java.tmpl", Source: generator.SourceProject},
		{Name: "test/ControllerTest.java.tmpl", Source: generator.SourceGlobal},
		{Name: "project/Application.java.tmpl", Source: generator.SourceEmbedded},
	}

	filtered := filterCustomOnly(templates)
	assert.Len(t, filtered, 2)

	for _, tmpl := range filtered {
		assert.NotEqual(t, generator.SourceEmbedded, tmpl.Source)
	}
}

func TestGroupByCategory(t *testing.T) {
	templates := []generator.TemplateInfo{
		{Name: "resource/Controller.java.tmpl", Source: generator.SourceEmbedded},
		{Name: "resource/Service.java.tmpl", Source: generator.SourceEmbedded},
		{Name: "test/ControllerTest.java.tmpl", Source: generator.SourceEmbedded},
		{Name: "project/Application.java.tmpl", Source: generator.SourceEmbedded},
	}

	grouped := groupByCategory(templates)

	assert.Len(t, grouped["resource"], 2)
	assert.Len(t, grouped["test"], 1)
	assert.Len(t, grouped["project"], 1)
}

func TestGroupByCategoryWithOther(t *testing.T) {
	templates := []generator.TemplateInfo{
		{Name: "custom/MyTemplate.java.tmpl", Source: generator.SourceProject},
	}

	grouped := groupByCategory(templates)
	assert.Len(t, grouped["custom"], 1)
}

func TestCountBySource(t *testing.T) {
	templates := []generator.TemplateInfo{
		{Name: "t1", Source: generator.SourceEmbedded},
		{Name: "t2", Source: generator.SourceEmbedded},
		{Name: "t3", Source: generator.SourceProject},
		{Name: "t4", Source: generator.SourceGlobal},
		{Name: "t5", Source: generator.SourceGlobal},
		{Name: "t6", Source: generator.SourceGlobal},
	}

	assert.Equal(t, 2, countBySource(templates, generator.SourceEmbedded))
	assert.Equal(t, 1, countBySource(templates, generator.SourceProject))
	assert.Equal(t, 3, countBySource(templates, generator.SourceGlobal))
}

func TestCountBySourceEmpty(t *testing.T) {
	var templates []generator.TemplateInfo

	assert.Equal(t, 0, countBySource(templates, generator.SourceEmbedded))
	assert.Equal(t, 0, countBySource(templates, generator.SourceProject))
	assert.Equal(t, 0, countBySource(templates, generator.SourceGlobal))
}

func TestTemplateInfoStruct(t *testing.T) {
	info := generator.TemplateInfo{
		Name:   "resource/Controller.java.tmpl",
		Path:   "/path/to/template",
		Source: generator.SourceProject,
	}

	assert.Equal(t, "resource/Controller.java.tmpl", info.Name)
	assert.Equal(t, "/path/to/template", info.Path)
	assert.Equal(t, generator.SourceProject, info.Source)
}

func TestTemplateSourceConstants(t *testing.T) {
	assert.Equal(t, "embedded", generator.SourceEmbedded.String())
	assert.Equal(t, "project", generator.SourceProject.String())
	assert.Equal(t, "global", generator.SourceGlobal.String())
}

func TestStyleVariablesExist(t *testing.T) {
	assert.NotNil(t, headerStyle)
	assert.NotNil(t, sourceStyle)
	assert.NotNil(t, projectStyle)
	assert.NotNil(t, globalStyle)
	assert.NotNil(t, pathStyle)
}

func TestValidateStyleVariablesExist(t *testing.T) {
	assert.NotNil(t, errorStyle)
	assert.NotNil(t, warningStyle)
	assert.NotNil(t, successStyle)
	assert.NotNil(t, lineStyle)
	assert.NotNil(t, infoStyle)
}

func TestFilterByCategoryEmptyInput(t *testing.T) {
	var templates []generator.TemplateInfo
	filtered := filterByCategory(templates, "resource")
	assert.Len(t, filtered, 0)
}

func TestFilterCustomOnlyEmptyInput(t *testing.T) {
	var templates []generator.TemplateInfo
	filtered := filterCustomOnly(templates)
	assert.Len(t, filtered, 0)
}

func TestGroupByCategoryEmptyInput(t *testing.T) {
	var templates []generator.TemplateInfo
	grouped := groupByCategory(templates)
	assert.Len(t, grouped, 0)
}

func TestFilterByCategoryAllEmbedded(t *testing.T) {
	templates := []generator.TemplateInfo{
		{Name: "resource/Controller.java.tmpl", Source: generator.SourceEmbedded},
		{Name: "resource/Service.java.tmpl", Source: generator.SourceEmbedded},
	}

	filtered := filterCustomOnly(templates)
	assert.Len(t, filtered, 0)
}

func TestGroupByCategoryMultipleCategories(t *testing.T) {
	templates := []generator.TemplateInfo{
		{Name: "resource/A.tmpl", Source: generator.SourceEmbedded},
		{Name: "resource/B.tmpl", Source: generator.SourceEmbedded},
		{Name: "test/A.tmpl", Source: generator.SourceEmbedded},
		{Name: "test/B.tmpl", Source: generator.SourceEmbedded},
		{Name: "project/A.tmpl", Source: generator.SourceEmbedded},
		{Name: "config/A.tmpl", Source: generator.SourceEmbedded},
		{Name: "exception/A.tmpl", Source: generator.SourceEmbedded},
	}

	grouped := groupByCategory(templates)

	assert.Len(t, grouped["resource"], 2)
	assert.Len(t, grouped["test"], 2)
	assert.Len(t, grouped["project"], 1)
	assert.Len(t, grouped["config"], 1)
	assert.Len(t, grouped["exception"], 1)
}

func TestPrintTemplateInfoDoesNotPanic(t *testing.T) {
	info := generator.TemplateInfo{
		Name:   "resource/Controller.java.tmpl",
		Path:   "/path/to/template",
		Source: generator.SourceProject,
	}

	assert.NotPanics(t, func() {
		printTemplateInfo(info, false)
	})

	assert.NotPanics(t, func() {
		printTemplateInfo(info, true)
	})
}

func TestPrintTemplateInfoAllSources(t *testing.T) {
	sources := []generator.TemplateSource{
		generator.SourceEmbedded,
		generator.SourceProject,
		generator.SourceGlobal,
	}

	for _, source := range sources {
		info := generator.TemplateInfo{
			Name:   "test/Template.java.tmpl",
			Path:   "/path/to/template",
			Source: source,
		}

		assert.NotPanics(t, func() {
			printTemplateInfo(info, false)
			printTemplateInfo(info, true)
		})
	}
}

func TestPrintAvailableVariablesDoesNotPanic(t *testing.T) {
	assert.NotPanics(t, func() {
		printAvailableVariables()
	})
}

func TestPrintAvailableConditionsDoesNotPanic(t *testing.T) {
	assert.NotPanics(t, func() {
		printAvailableConditions()
	})
}

func TestRunListNoError(t *testing.T) {
	err := runList(false, "", false)
	assert.NoError(t, err)
}

func TestRunListWithCustomOnly(t *testing.T) {
	err := runList(true, "", false)
	assert.NoError(t, err)
}

func TestRunListWithCategory(t *testing.T) {
	err := runList(false, "resource", false)
	assert.NoError(t, err)
}

func TestRunListWithCategoryAndPaths(t *testing.T) {
	err := runList(false, "resource", true)
	assert.NoError(t, err)
}

func TestRunListWithAllOptions(t *testing.T) {
	err := runList(true, "test", true)
	assert.NoError(t, err)
}

func TestRunValidateNoTemplates(t *testing.T) {
	err := runValidate(nil)
	assert.NoError(t, err)
}

func TestRunValidateNonexistentPath(t *testing.T) {
	err := runValidate([]string{"/nonexistent/path/to/template.tmpl"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path not found")
}

func TestValidateCommandVarsFlag(t *testing.T) {
	cmd := newValidateCommand()
	cmd.SetArgs([]string{"--vars"})

	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestValidateCommandConditionsFlag(t *testing.T) {
	cmd := newValidateCommand()
	cmd.SetArgs([]string{"--conditions"})

	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestInitCommandInvalidCategory(t *testing.T) {
	err := runInit("nonexistent_category", false, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid category")
}

func TestRunValidateWithTempDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "haft-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	validTemplate := `package ${BasePackage}.controller;

import org.springframework.web.bind.annotation.RestController;

@RestController
public class ${Name}Controller {
}
`
	templatePath := filepath.Join(tmpDir, "Controller.java.tmpl")
	err = os.WriteFile(templatePath, []byte(validTemplate), 0644)
	require.NoError(t, err)

	err = runValidate([]string{templatePath})
	assert.NoError(t, err)
}

func TestRunValidateWithInvalidTemplate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "haft-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	invalidTemplate := `package ${BasePackage.controller;

// @if HasLombok
@Data
// missing @endif
`
	templatePath := filepath.Join(tmpDir, "Invalid.java.tmpl")
	err = os.WriteFile(templatePath, []byte(invalidTemplate), 0644)
	require.NoError(t, err)

	err = runValidate([]string{templatePath})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestRunValidateWithDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "haft-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	validTemplate := `package ${BasePackage}.service;

public class ${Name}Service {
}
`
	templatePath := filepath.Join(tmpDir, "Service.java.tmpl")
	err = os.WriteFile(templatePath, []byte(validTemplate), 0644)
	require.NoError(t, err)

	err = runValidate([]string{tmpDir})
	assert.NoError(t, err)
}

func TestRunValidateEmptyDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "haft-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	err = runValidate([]string{tmpDir})
	assert.NoError(t, err)
}

func TestRunValidateWithWarnings(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "haft-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	templateWithUnknown := `package ${BasePackage}.service;
import ${UnknownVariable};

public class ${Name}Service {
}
`
	templatePath := filepath.Join(tmpDir, "ServiceWarning.java.tmpl")
	err = os.WriteFile(templatePath, []byte(templateWithUnknown), 0644)
	require.NoError(t, err)

	err = runValidate([]string{templatePath})
	assert.NoError(t, err)
}

func TestListCommandExecution(t *testing.T) {
	cmd := newListCommand()
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestListCommandWithCustomFlag(t *testing.T) {
	cmd := newListCommand()
	cmd.SetArgs([]string{"--custom"})

	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestListCommandWithCategoryFlag(t *testing.T) {
	cmd := newListCommand()
	cmd.SetArgs([]string{"--category", "resource"})

	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestListCommandWithPathsFlag(t *testing.T) {
	cmd := newListCommand()
	cmd.SetArgs([]string{"--paths"})

	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestFilterByCategoryPrefix(t *testing.T) {
	templates := []generator.TemplateInfo{
		{Name: "resource/layered/Controller.java.tmpl", Source: generator.SourceEmbedded},
		{Name: "resource/feature/Controller.java.tmpl", Source: generator.SourceEmbedded},
		{Name: "resourceful/Other.java.tmpl", Source: generator.SourceEmbedded},
	}

	filtered := filterByCategory(templates, "resource")
	assert.Len(t, filtered, 2)

	for _, tmpl := range filtered {
		assert.True(t, strings.HasPrefix(tmpl.Name, "resource/"))
	}
}

func TestCountBySourceMixed(t *testing.T) {
	templates := []generator.TemplateInfo{
		{Name: "t1", Source: generator.SourceEmbedded},
		{Name: "t2", Source: generator.SourceProject},
		{Name: "t3", Source: generator.SourceEmbedded},
		{Name: "t4", Source: generator.SourceProject},
		{Name: "t5", Source: generator.SourceEmbedded},
	}

	assert.Equal(t, 3, countBySource(templates, generator.SourceEmbedded))
	assert.Equal(t, 2, countBySource(templates, generator.SourceProject))
	assert.Equal(t, 0, countBySource(templates, generator.SourceGlobal))
}

func TestGroupByCategoryNested(t *testing.T) {
	templates := []generator.TemplateInfo{
		{Name: "resource/layered/Controller.java.tmpl", Source: generator.SourceEmbedded},
		{Name: "resource/layered/Service.java.tmpl", Source: generator.SourceEmbedded},
		{Name: "resource/feature/Controller.java.tmpl", Source: generator.SourceEmbedded},
	}

	grouped := groupByCategory(templates)
	assert.Len(t, grouped["resource"], 3)
}

func TestInitCommandCategoryResource(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "haft-init-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	err = runInit("resource", false, false)
	assert.NoError(t, err)
}

func TestInitCommandAllTemplates(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "haft-init-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	err = runInit("", false, false)
	assert.NoError(t, err)
}

func TestInitCommandForce(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "haft-init-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	err = runInit("resource", false, false)
	require.NoError(t, err)

	err = runInit("resource", false, true)
	assert.NoError(t, err)
}

func TestInitCommandSkipExisting(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "haft-init-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	err = runInit("resource", false, false)
	require.NoError(t, err)

	err = runInit("resource", false, false)
	assert.NoError(t, err)
}
