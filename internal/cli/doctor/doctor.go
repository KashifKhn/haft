package doctor

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var jsonOutput bool
	var strict bool
	var category string

	cmd := &cobra.Command{
		Use:     "doctor",
		Aliases: []string{"doc", "check", "health"},
		Short:   "Check project health and best practices",
		Long: `Analyze your Spring Boot project for issues, warnings, and suggestions.

The doctor command performs comprehensive health checks including:
  - Build configuration (pom.xml/build.gradle)
  - Spring Boot setup
  - Source code structure
  - Configuration files
  - Security best practices
  - Dependency recommendations
  - Docker configuration (Dockerfile, docker-compose, .dockerignore)

It helps identify problems early and suggests improvements.`,
		Example: `  # Run full health check
  haft doctor

  # Output as JSON (for CI/CD)
  haft doctor --json

  # Strict mode (exit 1 on warnings)
  haft doctor --strict

  # Filter by category
  haft doctor --category security
  haft doctor --category docker`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDoctor(Options{
				JSON:     jsonOutput,
				Strict:   strict,
				Category: category,
			})
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output results as JSON")
	cmd.Flags().BoolVar(&strict, "strict", false, "Exit with code 1 on any warning")
	cmd.Flags().StringVar(&category, "category", "", "Filter by category (build, source, config, security, dependencies, best-practice, docker)")

	return cmd
}

func runDoctor(opts Options) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	fs := afero.NewOsFs()
	report := RunDoctorChecks(fs, cwd, opts)

	output := FormatReport(report, opts)
	fmt.Print(output)

	if report.HasIssues() {
		return fmt.Errorf("health check failed with %d errors", report.ErrorCount)
	}

	if opts.Strict && report.HasWarnings() {
		return fmt.Errorf("health check failed with %d warnings (strict mode)", report.WarningCount)
	}

	return nil
}

func RunDoctorChecks(fs afero.Fs, projectPath string, opts Options) *Report {
	checker := NewChecker(fs, projectPath)
	results := checker.RunAllChecks()

	if opts.Category != "" {
		results = filterByCategory(results, Category(opts.Category))
	}

	report := &Report{
		ProjectPath: projectPath,
		ProjectName: filepath.Base(projectPath),
		BuildTool:   checker.buildTool,
		Results:     results,
	}

	report.CalculateCounts()
	return report
}

func filterByCategory(results []CheckResult, category Category) []CheckResult {
	var filtered []CheckResult
	for _, r := range results {
		if r.Category == category {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
