package dev

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/KashifKhn/haft/internal/logger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type ValidationSeverity string

const (
	SeverityError   ValidationSeverity = "error"
	SeverityWarning ValidationSeverity = "warning"
	SeverityInfo    ValidationSeverity = "info"
)

type ValidationResult struct {
	Check    string             `json:"check"`
	Passed   bool               `json:"passed"`
	Severity ValidationSeverity `json:"severity"`
	Message  string             `json:"message"`
	Details  string             `json:"details,omitempty"`
}

type ValidationReport struct {
	ProjectPath   string             `json:"project_path"`
	BuildTool     string             `json:"build_tool"`
	Passed        bool               `json:"passed"`
	ErrorCount    int                `json:"error_count"`
	WarningCount  int                `json:"warning_count"`
	Results       []ValidationResult `json:"results"`
	BuildToolPass bool               `json:"build_tool_pass"`
}

type ValidateOptions struct {
	Strict        bool
	SkipBuildTool bool
	JSONOutput    bool
	Fs            afero.Fs
}

func newValidateCommand() *cobra.Command {
	var strict bool
	var skipBuildTool bool
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:     "validate",
		Aliases: []string{"v", "check"},
		Short:   "Validate project configuration and structure",
		Long: `Validate the Spring Boot project configuration and structure.

This command performs two types of validation:

1. Haft Validation (Custom Checks):
   - Build file exists and is valid (pom.xml/build.gradle)
   - Spring Boot parent/plugin configured
   - Java version specified
   - Source directory structure exists
   - Main application class with @SpringBootApplication
   - Configuration files present (application.properties/yml)
   - Required Spring Boot starters present

2. Build Tool Validation:
   - Runs 'mvn validate' for Maven projects
   - Runs './gradlew help' for Gradle projects
   - Verifies build tool configuration is correct

Use --skip-build-tool to only run Haft's custom validation checks.`,
		Example: `  # Validate project (both custom and build tool validation)
  haft dev validate

  # Validate with strict mode (warnings become errors)
  haft dev validate --strict

  # Skip build tool validation (only run Haft checks)
  haft dev validate --skip-build-tool

  # Output as JSON for CI pipelines
  haft dev validate --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := ValidateOptions{
				Strict:        strict,
				SkipBuildTool: skipBuildTool,
				JSONOutput:    jsonOutput,
				Fs:            afero.NewOsFs(),
			}
			return runValidate(opts)
		},
	}

	cmd.Flags().BoolVarP(&strict, "strict", "s", false, "Treat warnings as errors")
	cmd.Flags().BoolVar(&skipBuildTool, "skip-build-tool", false, "Skip build tool validation (mvn validate/gradle)")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output results as JSON")

	return cmd
}

func runValidate(opts ValidateOptions) error {
	report, err := RunValidation(opts)
	if err != nil {
		return err
	}

	if opts.JSONOutput {
		return outputJSON(report)
	}

	return outputText(report, opts.Strict)
}

func RunValidation(opts ValidateOptions) (*ValidationReport, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("could not get current directory: %w", err)
	}

	if opts.Fs == nil {
		opts.Fs = afero.NewOsFs()
	}

	result, err := buildtool.Detect(cwd, opts.Fs)
	if err != nil {
		return &ValidationReport{
			ProjectPath: cwd,
			Passed:      false,
			ErrorCount:  1,
			Results: []ValidationResult{{
				Check:    "build_file",
				Passed:   false,
				Severity: SeverityError,
				Message:  "No build file found",
				Details:  "Could not find pom.xml or build.gradle in current directory or parents",
			}},
		}, nil
	}

	report := &ValidationReport{
		ProjectPath: cwd,
		BuildTool:   result.BuildTool.DisplayName(),
		Results:     make([]ValidationResult, 0),
	}

	runCustomValidation(cwd, result, opts.Fs, report)

	if !opts.SkipBuildTool {
		runBuildToolValidation(result, report)
	} else {
		report.BuildToolPass = true
	}

	report.ErrorCount = 0
	report.WarningCount = 0
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

	if opts.Strict {
		report.Passed = report.ErrorCount == 0 && report.WarningCount == 0 && report.BuildToolPass
	} else {
		report.Passed = report.ErrorCount == 0 && report.BuildToolPass
	}

	return report, nil
}

func runCustomValidation(projectPath string, detection *buildtool.DetectionResult, fs afero.Fs, report *ValidationReport) {
	report.Results = append(report.Results, validateBuildFile(detection))

	project, err := detection.Parser.Parse(detection.FilePath)
	if err != nil {
		report.Results = append(report.Results, ValidationResult{
			Check:    "build_file_parse",
			Passed:   false,
			Severity: SeverityError,
			Message:  "Failed to parse build file",
			Details:  err.Error(),
		})
		return
	}

	report.Results = append(report.Results, ValidationResult{
		Check:    "build_file_parse",
		Passed:   true,
		Severity: SeverityInfo,
		Message:  "Build file parsed successfully",
	})

	report.Results = append(report.Results, validateSpringBootConfig(project, detection.BuildTool))
	report.Results = append(report.Results, validateJavaVersion(project))
	report.Results = append(report.Results, validateSourceDirectory(projectPath, fs))
	report.Results = append(report.Results, validateResourcesDirectory(projectPath, fs))
	report.Results = append(report.Results, validateMainClass(projectPath, project, fs))
	report.Results = append(report.Results, validateConfigFile(projectPath, fs))
	report.Results = append(report.Results, validateSpringBootStarter(project, detection.Parser))
	report.Results = append(report.Results, validateTestDirectory(projectPath, fs))
}

func validateBuildFile(detection *buildtool.DetectionResult) ValidationResult {
	return ValidationResult{
		Check:    "build_file",
		Passed:   true,
		Severity: SeverityInfo,
		Message:  fmt.Sprintf("Build file found: %s", filepath.Base(detection.FilePath)),
		Details:  detection.FilePath,
	}
}

func validateSpringBootConfig(project *buildtool.Project, buildTool buildtool.Type) ValidationResult {
	if project.SpringBootVersion != "" {
		return ValidationResult{
			Check:    "spring_boot_config",
			Passed:   true,
			Severity: SeverityInfo,
			Message:  fmt.Sprintf("Spring Boot %s configured", project.SpringBootVersion),
		}
	}

	hasSpringBootDep := false
	for _, dep := range project.Dependencies {
		if dep.GroupId == "org.springframework.boot" {
			hasSpringBootDep = true
			break
		}
	}

	if hasSpringBootDep {
		return ValidationResult{
			Check:    "spring_boot_config",
			Passed:   true,
			Severity: SeverityWarning,
			Message:  "Spring Boot dependencies found but version not detected from parent",
		}
	}

	return ValidationResult{
		Check:    "spring_boot_config",
		Passed:   false,
		Severity: SeverityWarning,
		Message:  "Spring Boot parent/plugin not configured",
		Details:  "Add spring-boot-starter-parent as parent (Maven) or spring-boot plugin (Gradle)",
	}
}

func validateJavaVersion(project *buildtool.Project) ValidationResult {
	if project.JavaVersion != "" {
		version := project.JavaVersion
		if version == "17" || version == "21" || version == "11" {
			return ValidationResult{
				Check:    "java_version",
				Passed:   true,
				Severity: SeverityInfo,
				Message:  fmt.Sprintf("Java version %s configured", version),
			}
		}
		if version == "8" || version == "1.8" {
			return ValidationResult{
				Check:    "java_version",
				Passed:   true,
				Severity: SeverityWarning,
				Message:  fmt.Sprintf("Java version %s configured (consider upgrading to 17 or 21)", version),
			}
		}
		return ValidationResult{
			Check:    "java_version",
			Passed:   true,
			Severity: SeverityInfo,
			Message:  fmt.Sprintf("Java version %s configured", version),
		}
	}

	return ValidationResult{
		Check:    "java_version",
		Passed:   false,
		Severity: SeverityWarning,
		Message:  "Java version not specified",
		Details:  "Set java.version property in build file",
	}
}

func validateSourceDirectory(projectPath string, fs afero.Fs) ValidationResult {
	srcMainJava := filepath.Join(projectPath, "src", "main", "java")
	srcMainKotlin := filepath.Join(projectPath, "src", "main", "kotlin")

	javaExists := dirExists(fs, srcMainJava)
	kotlinExists := dirExists(fs, srcMainKotlin)

	if javaExists && kotlinExists {
		return ValidationResult{
			Check:    "source_directory",
			Passed:   true,
			Severity: SeverityInfo,
			Message:  "Source directories found (Java + Kotlin)",
		}
	}

	if javaExists {
		return ValidationResult{
			Check:    "source_directory",
			Passed:   true,
			Severity: SeverityInfo,
			Message:  "Source directory found: src/main/java",
		}
	}

	if kotlinExists {
		return ValidationResult{
			Check:    "source_directory",
			Passed:   true,
			Severity: SeverityInfo,
			Message:  "Source directory found: src/main/kotlin",
		}
	}

	return ValidationResult{
		Check:    "source_directory",
		Passed:   false,
		Severity: SeverityError,
		Message:  "Source directory not found",
		Details:  "Expected src/main/java or src/main/kotlin directory",
	}
}

func validateResourcesDirectory(projectPath string, fs afero.Fs) ValidationResult {
	resourcesDir := filepath.Join(projectPath, "src", "main", "resources")

	if dirExists(fs, resourcesDir) {
		return ValidationResult{
			Check:    "resources_directory",
			Passed:   true,
			Severity: SeverityInfo,
			Message:  "Resources directory found: src/main/resources",
		}
	}

	return ValidationResult{
		Check:    "resources_directory",
		Passed:   false,
		Severity: SeverityWarning,
		Message:  "Resources directory not found",
		Details:  "Expected src/main/resources directory",
	}
}

func validateMainClass(projectPath string, project *buildtool.Project, fs afero.Fs) ValidationResult {
	srcDirs := []string{
		filepath.Join(projectPath, "src", "main", "java"),
		filepath.Join(projectPath, "src", "main", "kotlin"),
	}

	for _, srcDir := range srcDirs {
		if !dirExists(fs, srcDir) {
			continue
		}

		found, path := findMainClass(srcDir, fs)
		if found {
			return ValidationResult{
				Check:    "main_class",
				Passed:   true,
				Severity: SeverityInfo,
				Message:  "Main application class found",
				Details:  path,
			}
		}
	}

	return ValidationResult{
		Check:    "main_class",
		Passed:   false,
		Severity: SeverityWarning,
		Message:  "Main application class not found",
		Details:  "Expected a class with @SpringBootApplication annotation",
	}
}

func findMainClass(srcDir string, fs afero.Fs) (bool, string) {
	var foundPath string
	found := false

	_ = afero.Walk(fs, srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".java") && !strings.HasSuffix(path, ".kt") {
			return nil
		}

		content, err := afero.ReadFile(fs, path)
		if err != nil {
			return nil
		}

		if strings.Contains(string(content), "@SpringBootApplication") {
			found = true
			foundPath = path
			return filepath.SkipAll
		}

		return nil
	})

	return found, foundPath
}

func validateConfigFile(projectPath string, fs afero.Fs) ValidationResult {
	resourcesDir := filepath.Join(projectPath, "src", "main", "resources")

	configFiles := []string{
		"application.yml",
		"application.yaml",
		"application.properties",
	}

	for _, file := range configFiles {
		configPath := filepath.Join(resourcesDir, file)
		if fileExists(fs, configPath) {
			return ValidationResult{
				Check:    "config_file",
				Passed:   true,
				Severity: SeverityInfo,
				Message:  fmt.Sprintf("Configuration file found: %s", file),
			}
		}
	}

	return ValidationResult{
		Check:    "config_file",
		Passed:   false,
		Severity: SeverityWarning,
		Message:  "Configuration file not found",
		Details:  "Expected application.properties or application.yml in src/main/resources",
	}
}

func validateSpringBootStarter(project *buildtool.Project, parser buildtool.Parser) ValidationResult {
	starterDeps := []struct {
		groupId    string
		artifactId string
		name       string
	}{
		{"org.springframework.boot", "spring-boot-starter", "spring-boot-starter"},
		{"org.springframework.boot", "spring-boot-starter-web", "spring-boot-starter-web"},
		{"org.springframework.boot", "spring-boot-starter-webflux", "spring-boot-starter-webflux"},
		{"org.springframework.boot", "spring-boot-starter-data-jpa", "spring-boot-starter-data-jpa"},
		{"org.springframework.boot", "spring-boot-starter-data-mongodb", "spring-boot-starter-data-mongodb"},
		{"org.springframework.boot", "spring-boot-starter-security", "spring-boot-starter-security"},
		{"org.springframework.boot", "spring-boot-starter-actuator", "spring-boot-starter-actuator"},
		{"org.springframework.boot", "spring-boot-starter-test", "spring-boot-starter-test"},
	}

	var foundStarters []string
	for _, starter := range starterDeps {
		if parser.HasDependency(project, starter.groupId, starter.artifactId) {
			foundStarters = append(foundStarters, starter.name)
		}
	}

	if len(foundStarters) > 0 {
		return ValidationResult{
			Check:    "spring_boot_starter",
			Passed:   true,
			Severity: SeverityInfo,
			Message:  fmt.Sprintf("Found %d Spring Boot starters", len(foundStarters)),
			Details:  strings.Join(foundStarters, ", "),
		}
	}

	return ValidationResult{
		Check:    "spring_boot_starter",
		Passed:   false,
		Severity: SeverityWarning,
		Message:  "No Spring Boot starters found",
		Details:  "Add at least one spring-boot-starter dependency",
	}
}

func validateTestDirectory(projectPath string, fs afero.Fs) ValidationResult {
	testJava := filepath.Join(projectPath, "src", "test", "java")
	testKotlin := filepath.Join(projectPath, "src", "test", "kotlin")

	javaExists := dirExists(fs, testJava)
	kotlinExists := dirExists(fs, testKotlin)

	if javaExists || kotlinExists {
		var msg string
		if javaExists && kotlinExists {
			msg = "Test directories found (Java + Kotlin)"
		} else if javaExists {
			msg = "Test directory found: src/test/java"
		} else {
			msg = "Test directory found: src/test/kotlin"
		}
		return ValidationResult{
			Check:    "test_directory",
			Passed:   true,
			Severity: SeverityInfo,
			Message:  msg,
		}
	}

	return ValidationResult{
		Check:    "test_directory",
		Passed:   false,
		Severity: SeverityWarning,
		Message:  "Test directory not found",
		Details:  "Expected src/test/java or src/test/kotlin directory",
	}
}

func runBuildToolValidation(detection *buildtool.DetectionResult, report *ValidationReport) {
	var executable string
	var args []string

	switch detection.BuildTool {
	case buildtool.Maven:
		executable = getMavenExecutable()
		args = []string{"validate", "-q"}
	case buildtool.Gradle, buildtool.GradleKotln:
		executable = getGradleExecutable()
		args = []string{"help", "-q"}
	}

	err := executeCommandSilent(executable, args)
	if err != nil {
		report.BuildToolPass = false
		report.Results = append(report.Results, ValidationResult{
			Check:    "build_tool_validation",
			Passed:   false,
			Severity: SeverityError,
			Message:  fmt.Sprintf("%s validation failed", detection.BuildTool.DisplayName()),
			Details:  err.Error(),
		})
		return
	}

	report.BuildToolPass = true
	report.Results = append(report.Results, ValidationResult{
		Check:    "build_tool_validation",
		Passed:   true,
		Severity: SeverityInfo,
		Message:  fmt.Sprintf("%s validation passed", detection.BuildTool.DisplayName()),
	})
}

func executeCommandSilent(executable string, args []string) error {
	cmd := exec.Command(executable, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func dirExists(fs afero.Fs, path string) bool {
	info, err := fs.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func fileExists(fs afero.Fs, path string) bool {
	info, err := fs.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func outputJSON(report *ValidationReport) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

func outputText(report *ValidationReport, strict bool) error {
	if !report.Passed || report.ErrorCount > 0 || report.WarningCount > 0 {
		logger.Info("Validating project", "path", report.ProjectPath, "build-tool", report.BuildTool)
	} else {
		logger.Info("Validating project", "path", report.ProjectPath, "build-tool", report.BuildTool)
	}

	fmt.Println()

	for _, result := range report.Results {
		var icon string
		if result.Passed {
			icon = "\u2713"
		} else {
			switch result.Severity {
			case SeverityError:
				icon = "\u2717"
			case SeverityWarning:
				icon = "!"
			default:
				icon = "-"
			}
		}

		fmt.Printf("  %s %s\n", icon, result.Message)
		if !result.Passed && result.Details != "" {
			fmt.Printf("    %s\n", result.Details)
		}
	}

	fmt.Println()

	if report.Passed {
		logger.Info("Validation passed",
			"errors", report.ErrorCount,
			"warnings", report.WarningCount)
	} else {
		if strict && report.WarningCount > 0 {
			logger.Error("Validation failed (strict mode)",
				"errors", report.ErrorCount,
				"warnings", report.WarningCount)
		} else {
			logger.Error("Validation failed",
				"errors", report.ErrorCount,
				"warnings", report.WarningCount)
		}
		return fmt.Errorf("validation failed with %d errors", report.ErrorCount)
	}

	return nil
}
