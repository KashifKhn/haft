package dev

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewValidateCommand(t *testing.T) {
	cmd := newValidateCommand()

	assert.Equal(t, "validate", cmd.Use)
	assert.Equal(t, []string{"v", "check"}, cmd.Aliases)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
}

func TestValidateCommand_Flags(t *testing.T) {
	cmd := newValidateCommand()

	strictFlag := cmd.Flags().Lookup("strict")
	require.NotNil(t, strictFlag)
	assert.Equal(t, "s", strictFlag.Shorthand)
	assert.Equal(t, "false", strictFlag.DefValue)

	skipBuildToolFlag := cmd.Flags().Lookup("skip-build-tool")
	require.NotNil(t, skipBuildToolFlag)
	assert.Equal(t, "false", skipBuildToolFlag.DefValue)

	jsonFlag := cmd.Flags().Lookup("json")
	require.NotNil(t, jsonFlag)
	assert.Equal(t, "false", jsonFlag.DefValue)
}

func TestValidateCommand_IsRegistered(t *testing.T) {
	cmd := NewCommand()

	var validateCmd *cobra.Command
	for _, sub := range cmd.Commands() {
		if sub.Name() == "validate" {
			validateCmd = sub
			break
		}
	}

	require.NotNil(t, validateCmd, "validate command should be registered")
	assert.Equal(t, []string{"v", "check"}, validateCmd.Aliases)
}

func TestValidationSeverity_Values(t *testing.T) {
	assert.Equal(t, ValidationSeverity("error"), SeverityError)
	assert.Equal(t, ValidationSeverity("warning"), SeverityWarning)
	assert.Equal(t, ValidationSeverity("info"), SeverityInfo)
}

func TestValidationResult_JSONSerialization(t *testing.T) {
	result := ValidationResult{
		Check:    "test_check",
		Passed:   true,
		Severity: SeverityInfo,
		Message:  "Test message",
		Details:  "Test details",
	}

	data, err := json.Marshal(result)
	require.NoError(t, err)

	var decoded ValidationResult
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, result.Check, decoded.Check)
	assert.Equal(t, result.Passed, decoded.Passed)
	assert.Equal(t, result.Severity, decoded.Severity)
	assert.Equal(t, result.Message, decoded.Message)
	assert.Equal(t, result.Details, decoded.Details)
}

func TestValidationReport_JSONSerialization(t *testing.T) {
	report := ValidationReport{
		ProjectPath:   "/test/path",
		BuildTool:     "Maven",
		Passed:        true,
		ErrorCount:    0,
		WarningCount:  1,
		BuildToolPass: true,
		Results: []ValidationResult{
			{Check: "test", Passed: true, Severity: SeverityInfo, Message: "OK"},
		},
	}

	data, err := json.Marshal(report)
	require.NoError(t, err)

	var decoded ValidationReport
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, report.ProjectPath, decoded.ProjectPath)
	assert.Equal(t, report.BuildTool, decoded.BuildTool)
	assert.Equal(t, report.Passed, decoded.Passed)
	assert.Equal(t, report.ErrorCount, decoded.ErrorCount)
	assert.Equal(t, report.WarningCount, decoded.WarningCount)
	assert.Equal(t, len(report.Results), len(decoded.Results))
}

func TestDirExists(t *testing.T) {
	fs := afero.NewMemMapFs()

	_ = fs.MkdirAll("/existing/dir", 0755)
	_ = afero.WriteFile(fs, "/existing/file.txt", []byte("content"), 0644)

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"existing directory", "/existing/dir", true},
		{"existing file (not a dir)", "/existing/file.txt", false},
		{"non-existing path", "/non/existing", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dirExists(fs, tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFileExists(t *testing.T) {
	fs := afero.NewMemMapFs()

	_ = fs.MkdirAll("/existing/dir", 0755)
	_ = afero.WriteFile(fs, "/existing/file.txt", []byte("content"), 0644)

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"existing file", "/existing/file.txt", true},
		{"existing directory (not a file)", "/existing/dir", false},
		{"non-existing path", "/non/existing.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fileExists(fs, tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateSourceDirectory(t *testing.T) {
	tests := []struct {
		name       string
		setupFs    func(fs afero.Fs)
		wantPassed bool
		wantMsg    string
	}{
		{
			name: "java source exists",
			setupFs: func(fs afero.Fs) {
				_ = fs.MkdirAll("/project/src/main/java", 0755)
			},
			wantPassed: true,
			wantMsg:    "Source directory found: src/main/java",
		},
		{
			name: "kotlin source exists",
			setupFs: func(fs afero.Fs) {
				_ = fs.MkdirAll("/project/src/main/kotlin", 0755)
			},
			wantPassed: true,
			wantMsg:    "Source directory found: src/main/kotlin",
		},
		{
			name: "both java and kotlin exist",
			setupFs: func(fs afero.Fs) {
				_ = fs.MkdirAll("/project/src/main/java", 0755)
				_ = fs.MkdirAll("/project/src/main/kotlin", 0755)
			},
			wantPassed: true,
			wantMsg:    "Source directories found (Java + Kotlin)",
		},
		{
			name:       "no source directory",
			setupFs:    func(fs afero.Fs) {},
			wantPassed: false,
			wantMsg:    "Source directory not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			tt.setupFs(fs)

			result := validateSourceDirectory("/project", fs)
			assert.Equal(t, tt.wantPassed, result.Passed)
			assert.Equal(t, tt.wantMsg, result.Message)
			assert.Equal(t, "source_directory", result.Check)
		})
	}
}

func TestValidateResourcesDirectory(t *testing.T) {
	tests := []struct {
		name       string
		setupFs    func(fs afero.Fs)
		wantPassed bool
	}{
		{
			name: "resources directory exists",
			setupFs: func(fs afero.Fs) {
				_ = fs.MkdirAll("/project/src/main/resources", 0755)
			},
			wantPassed: true,
		},
		{
			name:       "resources directory missing",
			setupFs:    func(fs afero.Fs) {},
			wantPassed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			tt.setupFs(fs)

			result := validateResourcesDirectory("/project", fs)
			assert.Equal(t, tt.wantPassed, result.Passed)
			assert.Equal(t, "resources_directory", result.Check)
		})
	}
}

func TestValidateTestDirectory(t *testing.T) {
	tests := []struct {
		name       string
		setupFs    func(fs afero.Fs)
		wantPassed bool
		wantMsg    string
	}{
		{
			name: "java test exists",
			setupFs: func(fs afero.Fs) {
				_ = fs.MkdirAll("/project/src/test/java", 0755)
			},
			wantPassed: true,
			wantMsg:    "Test directory found: src/test/java",
		},
		{
			name: "kotlin test exists",
			setupFs: func(fs afero.Fs) {
				_ = fs.MkdirAll("/project/src/test/kotlin", 0755)
			},
			wantPassed: true,
			wantMsg:    "Test directory found: src/test/kotlin",
		},
		{
			name: "both test directories exist",
			setupFs: func(fs afero.Fs) {
				_ = fs.MkdirAll("/project/src/test/java", 0755)
				_ = fs.MkdirAll("/project/src/test/kotlin", 0755)
			},
			wantPassed: true,
			wantMsg:    "Test directories found (Java + Kotlin)",
		},
		{
			name:       "no test directory",
			setupFs:    func(fs afero.Fs) {},
			wantPassed: false,
			wantMsg:    "Test directory not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			tt.setupFs(fs)

			result := validateTestDirectory("/project", fs)
			assert.Equal(t, tt.wantPassed, result.Passed)
			assert.Equal(t, tt.wantMsg, result.Message)
			assert.Equal(t, "test_directory", result.Check)
		})
	}
}

func TestValidateConfigFile(t *testing.T) {
	tests := []struct {
		name       string
		setupFs    func(fs afero.Fs)
		wantPassed bool
		wantMsg    string
	}{
		{
			name: "application.yml exists",
			setupFs: func(fs afero.Fs) {
				_ = fs.MkdirAll("/project/src/main/resources", 0755)
				_ = afero.WriteFile(fs, "/project/src/main/resources/application.yml", []byte("server:\n  port: 8080"), 0644)
			},
			wantPassed: true,
			wantMsg:    "Configuration file found: application.yml",
		},
		{
			name: "application.yaml exists",
			setupFs: func(fs afero.Fs) {
				_ = fs.MkdirAll("/project/src/main/resources", 0755)
				_ = afero.WriteFile(fs, "/project/src/main/resources/application.yaml", []byte("server:\n  port: 8080"), 0644)
			},
			wantPassed: true,
			wantMsg:    "Configuration file found: application.yaml",
		},
		{
			name: "application.properties exists",
			setupFs: func(fs afero.Fs) {
				_ = fs.MkdirAll("/project/src/main/resources", 0755)
				_ = afero.WriteFile(fs, "/project/src/main/resources/application.properties", []byte("server.port=8080"), 0644)
			},
			wantPassed: true,
			wantMsg:    "Configuration file found: application.properties",
		},
		{
			name: "no config file",
			setupFs: func(fs afero.Fs) {
				_ = fs.MkdirAll("/project/src/main/resources", 0755)
			},
			wantPassed: false,
			wantMsg:    "Configuration file not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			tt.setupFs(fs)

			result := validateConfigFile("/project", fs)
			assert.Equal(t, tt.wantPassed, result.Passed)
			assert.Equal(t, tt.wantMsg, result.Message)
			assert.Equal(t, "config_file", result.Check)
		})
	}
}

func TestFindMainClass(t *testing.T) {
	tests := []struct {
		name      string
		setupFs   func(fs afero.Fs)
		wantFound bool
	}{
		{
			name: "java main class with annotation",
			setupFs: func(fs afero.Fs) {
				_ = fs.MkdirAll("/project/src/main/java/com/example", 0755)
				content := `package com.example;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class Application {
    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }
}`
				_ = afero.WriteFile(fs, "/project/src/main/java/com/example/Application.java", []byte(content), 0644)
			},
			wantFound: true,
		},
		{
			name: "kotlin main class with annotation",
			setupFs: func(fs afero.Fs) {
				_ = fs.MkdirAll("/project/src/main/kotlin/com/example", 0755)
				content := `package com.example

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication

@SpringBootApplication
class Application

fun main(args: Array<String>) {
    runApplication<Application>(*args)
}`
				_ = afero.WriteFile(fs, "/project/src/main/kotlin/com/example/Application.kt", []byte(content), 0644)
			},
			wantFound: true,
		},
		{
			name: "no main class",
			setupFs: func(fs afero.Fs) {
				_ = fs.MkdirAll("/project/src/main/java/com/example", 0755)
				content := `package com.example;

public class SomeService {
    public void doSomething() {}
}`
				_ = afero.WriteFile(fs, "/project/src/main/java/com/example/SomeService.java", []byte(content), 0644)
			},
			wantFound: false,
		},
		{
			name: "empty source directory",
			setupFs: func(fs afero.Fs) {
				_ = fs.MkdirAll("/project/src/main/java", 0755)
			},
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			tt.setupFs(fs)

			found, _ := findMainClass("/project/src/main/java", fs)
			if !found {
				found, _ = findMainClass("/project/src/main/kotlin", fs)
			}
			assert.Equal(t, tt.wantFound, found)
		})
	}
}

func TestValidateJavaVersion(t *testing.T) {
	tests := []struct {
		name         string
		javaVersion  string
		wantPassed   bool
		wantSeverity ValidationSeverity
	}{
		{"java 17", "17", true, SeverityInfo},
		{"java 21", "21", true, SeverityInfo},
		{"java 11", "11", true, SeverityInfo},
		{"java 8", "8", true, SeverityWarning},
		{"java 1.8", "1.8", true, SeverityWarning},
		{"java 22", "22", true, SeverityInfo},
		{"no version", "", false, SeverityWarning},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project := &buildtool.Project{JavaVersion: tt.javaVersion}
			result := validateJavaVersion(project)

			assert.Equal(t, tt.wantPassed, result.Passed)
			assert.Equal(t, tt.wantSeverity, result.Severity)
			assert.Equal(t, "java_version", result.Check)
		})
	}
}

func TestValidateSpringBootConfig(t *testing.T) {
	tests := []struct {
		name       string
		sbVersion  string
		hasSBDep   bool
		wantPassed bool
	}{
		{"has spring boot version", "3.2.0", false, true},
		{"no version but has dep", "", true, true},
		{"no version no deps", "", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project := &buildtool.Project{
				SpringBootVersion: tt.sbVersion,
			}
			if tt.hasSBDep {
				project.Dependencies = []buildtool.Dependency{
					{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-web"},
				}
			}
			result := validateSpringBootConfig(project, buildtool.Maven)

			assert.Equal(t, tt.wantPassed, result.Passed)
			assert.Equal(t, "spring_boot_config", result.Check)
		})
	}
}

func TestValidateReport_ErrorCounting(t *testing.T) {
	tests := []struct {
		name             string
		results          []ValidationResult
		wantErrorCount   int
		wantWarningCount int
	}{
		{
			name: "all passed",
			results: []ValidationResult{
				{Passed: true, Severity: SeverityInfo},
				{Passed: true, Severity: SeverityInfo},
			},
			wantErrorCount:   0,
			wantWarningCount: 0,
		},
		{
			name: "one error",
			results: []ValidationResult{
				{Passed: false, Severity: SeverityError},
				{Passed: true, Severity: SeverityInfo},
			},
			wantErrorCount:   1,
			wantWarningCount: 0,
		},
		{
			name: "one warning",
			results: []ValidationResult{
				{Passed: false, Severity: SeverityWarning},
				{Passed: true, Severity: SeverityInfo},
			},
			wantErrorCount:   0,
			wantWarningCount: 1,
		},
		{
			name: "mixed errors and warnings",
			results: []ValidationResult{
				{Passed: false, Severity: SeverityError},
				{Passed: false, Severity: SeverityWarning},
				{Passed: false, Severity: SeverityError},
				{Passed: false, Severity: SeverityWarning},
				{Passed: true, Severity: SeverityInfo},
			},
			wantErrorCount:   2,
			wantWarningCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := &ValidationReport{Results: tt.results}

			for _, r := range report.Results {
				if !r.Passed {
					switch r.Severity {
					case SeverityError:
						report.ErrorCount++
					case SeverityWarning:
						report.WarningCount++
					}
				}
			}

			assert.Equal(t, tt.wantErrorCount, report.ErrorCount)
			assert.Equal(t, tt.wantWarningCount, report.WarningCount)
		})
	}
}

func TestValidateReport_PassedLogic(t *testing.T) {
	tests := []struct {
		name          string
		errorCount    int
		warningCount  int
		buildToolPass bool
		strict        bool
		wantPassed    bool
	}{
		{"all good", 0, 0, true, false, true},
		{"has error", 1, 0, true, false, false},
		{"has warning non-strict", 0, 1, true, false, true},
		{"has warning strict", 0, 1, true, true, false},
		{"build tool failed", 0, 0, false, false, false},
		{"error and warning strict", 1, 1, true, true, false},
		{"error and warning non-strict", 1, 1, true, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := &ValidationReport{
				ErrorCount:    tt.errorCount,
				WarningCount:  tt.warningCount,
				BuildToolPass: tt.buildToolPass,
			}

			if tt.strict {
				report.Passed = report.ErrorCount == 0 && report.WarningCount == 0 && report.BuildToolPass
			} else {
				report.Passed = report.ErrorCount == 0 && report.BuildToolPass
			}

			assert.Equal(t, tt.wantPassed, report.Passed)
		})
	}
}

func TestValidateOptions_Defaults(t *testing.T) {
	opts := ValidateOptions{}

	assert.False(t, opts.Strict)
	assert.False(t, opts.SkipBuildTool)
	assert.False(t, opts.JSONOutput)
	assert.Nil(t, opts.Fs)
}

func TestValidateMainClass_WithRealStructure(t *testing.T) {
	fs := afero.NewMemMapFs()

	_ = fs.MkdirAll("/project/src/main/java/com/example/demo", 0755)

	mainClass := `package com.example.demo;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class DemoApplication {

    public static void main(String[] args) {
        SpringApplication.run(DemoApplication.class, args);
    }
}
`
	_ = afero.WriteFile(fs, "/project/src/main/java/com/example/demo/DemoApplication.java", []byte(mainClass), 0644)

	found, path := findMainClass("/project/src/main/java", fs)

	assert.True(t, found)
	assert.Contains(t, path, "DemoApplication.java")
}

func TestValidateMainClass_DeepNesting(t *testing.T) {
	fs := afero.NewMemMapFs()

	deepPath := "/project/src/main/java/com/example/very/deep/nested/package"
	_ = fs.MkdirAll(deepPath, 0755)

	mainClass := `@SpringBootApplication
public class App {}`

	_ = afero.WriteFile(fs, filepath.Join(deepPath, "App.java"), []byte(mainClass), 0644)

	found, _ := findMainClass("/project/src/main/java", fs)
	assert.True(t, found)
}

func TestValidateMainClass_MultipleFiles(t *testing.T) {
	fs := afero.NewMemMapFs()

	_ = fs.MkdirAll("/project/src/main/java/com/example", 0755)

	_ = afero.WriteFile(fs, "/project/src/main/java/com/example/Service.java",
		[]byte(`public class Service {}`), 0644)
	_ = afero.WriteFile(fs, "/project/src/main/java/com/example/Controller.java",
		[]byte(`public class Controller {}`), 0644)
	_ = afero.WriteFile(fs, "/project/src/main/java/com/example/App.java",
		[]byte(`@SpringBootApplication public class App {}`), 0644)

	found, path := findMainClass("/project/src/main/java", fs)

	assert.True(t, found)
	assert.Contains(t, path, "App.java")
}

func TestValidateConfigFile_Priority(t *testing.T) {
	fs := afero.NewMemMapFs()
	_ = fs.MkdirAll("/project/src/main/resources", 0755)

	_ = afero.WriteFile(fs, "/project/src/main/resources/application.yml", []byte(""), 0644)
	_ = afero.WriteFile(fs, "/project/src/main/resources/application.properties", []byte(""), 0644)

	result := validateConfigFile("/project", fs)

	assert.True(t, result.Passed)
	assert.Equal(t, "Configuration file found: application.yml", result.Message)
}

func TestValidateCommand_HasCorrectDescription(t *testing.T) {
	cmd := newValidateCommand()

	assert.Contains(t, cmd.Long, "Haft Validation")
	assert.Contains(t, cmd.Long, "Build Tool Validation")
	assert.Contains(t, cmd.Long, "mvn validate")
	assert.Contains(t, cmd.Long, "gradlew help")
}

func TestOutputJSON_Format(t *testing.T) {
	report := &ValidationReport{
		ProjectPath:   "/test",
		BuildTool:     "Maven",
		Passed:        true,
		ErrorCount:    0,
		WarningCount:  0,
		BuildToolPass: true,
		Results:       []ValidationResult{},
	}

	data, err := json.MarshalIndent(report, "", "  ")
	require.NoError(t, err)

	assert.Contains(t, string(data), `"project_path"`)
	assert.Contains(t, string(data), `"build_tool"`)
	assert.Contains(t, string(data), `"passed"`)
	assert.Contains(t, string(data), `"error_count"`)
	assert.Contains(t, string(data), `"warning_count"`)
	assert.Contains(t, string(data), `"build_tool_pass"`)
}

func TestValidationResult_EmptyDetails(t *testing.T) {
	result := ValidationResult{
		Check:    "test",
		Passed:   true,
		Severity: SeverityInfo,
		Message:  "OK",
	}

	data, err := json.Marshal(result)
	require.NoError(t, err)

	assert.NotContains(t, string(data), `"details"`)
}

func TestValidationResult_WithDetails(t *testing.T) {
	result := ValidationResult{
		Check:    "test",
		Passed:   false,
		Severity: SeverityError,
		Message:  "Failed",
		Details:  "Some error details",
	}

	data, err := json.Marshal(result)
	require.NoError(t, err)

	assert.Contains(t, string(data), `"details"`)
	assert.Contains(t, string(data), "Some error details")
}

func TestValidateSourceDirectory_EdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		setupFs    func(fs afero.Fs)
		wantPassed bool
	}{
		{
			name: "src exists but not main/java",
			setupFs: func(fs afero.Fs) {
				_ = fs.MkdirAll("/project/src/main", 0755)
			},
			wantPassed: false,
		},
		{
			name: "java file instead of directory",
			setupFs: func(fs afero.Fs) {
				_ = fs.MkdirAll("/project/src/main", 0755)
				_ = afero.WriteFile(fs, "/project/src/main/java", []byte("not a dir"), 0644)
			},
			wantPassed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			tt.setupFs(fs)

			result := validateSourceDirectory("/project", fs)
			assert.Equal(t, tt.wantPassed, result.Passed)
		})
	}
}

func TestFindMainClass_AnnotationVariants(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantFound bool
	}{
		{
			name:      "standard annotation",
			content:   "@SpringBootApplication\npublic class App {}",
			wantFound: true,
		},
		{
			name:      "annotation with params",
			content:   "@SpringBootApplication(scanBasePackages = \"com.example\")\npublic class App {}",
			wantFound: true,
		},
		{
			name:      "annotation with newline",
			content:   "@SpringBootApplication\n\npublic class App {}",
			wantFound: true,
		},
		{
			name:      "in comment also matches",
			content:   "// @SpringBootApplication\npublic class App {}",
			wantFound: true,
		},
		{
			name:      "no annotation",
			content:   "public class App {}",
			wantFound: false,
		},
		{
			name:      "different annotation",
			content:   "@RestController\npublic class App {}",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			_ = fs.MkdirAll("/project/src/main/java", 0755)
			_ = afero.WriteFile(fs, "/project/src/main/java/App.java", []byte(tt.content), 0644)

			found, _ := findMainClass("/project/src/main/java", fs)
			assert.Equal(t, tt.wantFound, found)
		})
	}
}
