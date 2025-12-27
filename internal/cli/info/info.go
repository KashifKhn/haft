package info

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/KashifKhn/haft/internal/stats"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	valueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	countStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
)

func NewCommand() *cobra.Command {
	var jsonOutput bool
	var showLoc bool

	cmd := &cobra.Command{
		Use:   "info",
		Short: "Show project information",
		Long: `Display information about the current Spring Boot project.

Shows project metadata, build configuration, Spring Boot version,
Java version, and dependency summary.`,
		Example: `  # Show project info
  haft info

  # Output as JSON
  haft info --json

  # Include lines of code summary
  haft info --loc`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInfo(jsonOutput, showLoc)
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	cmd.Flags().BoolVar(&showLoc, "loc", false, "Show lines of code summary")

	return cmd
}

func runInfo(jsonOutput bool, showLoc bool) error {
	fs := afero.NewOsFs()
	result, err := buildtool.DetectWithCwd(fs)
	if err != nil {
		return fmt.Errorf("not a Spring Boot project: %w", err)
	}

	project, err := result.Parser.Parse(result.FilePath)
	if err != nil {
		return fmt.Errorf("failed to parse build file: %w", err)
	}

	var locStats *stats.ProjectStats
	if showLoc {
		cwd, _ := os.Getwd()
		locStats, _ = stats.CountProjectQuick(cwd)
	}

	if jsonOutput {
		return printJSON(project, result, locStats)
	}

	return printFormatted(project, result, locStats)
}

func printFormatted(project *buildtool.Project, result *buildtool.DetectionResult, locStats *stats.ProjectStats) error {
	cwd, _ := os.Getwd()
	projectName := filepath.Base(cwd)

	fmt.Println()
	fmt.Println(titleStyle.Render("  Project Information"))
	fmt.Println(strings.Repeat("─", 50))

	printRow("Name", projectName)
	printRow("Group ID", project.GroupId)
	printRow("Artifact ID", project.ArtifactId)
	printRow("Version", project.Version)
	if project.Description != "" {
		printRow("Description", project.Description)
	}

	fmt.Println()
	fmt.Println(titleStyle.Render("  Build Configuration"))
	fmt.Println(strings.Repeat("─", 50))

	printRow("Build Tool", result.BuildTool.DisplayName())
	printRow("Build File", filepath.Base(result.FilePath))
	printRow("Java Version", project.JavaVersion)
	printRow("Spring Boot", project.SpringBootVersion)
	if project.Packaging != "" {
		printRow("Packaging", project.Packaging)
	}

	fmt.Println()
	fmt.Println(titleStyle.Render("  Dependencies"))
	fmt.Println(strings.Repeat("─", 50))

	deps := project.Dependencies
	printRow("Total", countStyle.Render(fmt.Sprintf("%d", len(deps))))

	starters := countByPrefix(deps, "spring-boot-starter")
	if starters > 0 {
		printRow("Spring Starters", countStyle.Render(fmt.Sprintf("%d", starters)))
	}

	springDeps := countByPrefix(deps, "spring")
	if springDeps > 0 {
		printRow("Spring Libraries", countStyle.Render(fmt.Sprintf("%d", springDeps-starters)))
	}

	testDeps := countByScope(deps, "test")
	if testDeps > 0 {
		printRow("Test Dependencies", countStyle.Render(fmt.Sprintf("%d", testDeps)))
	}

	fmt.Println()
	fmt.Println(titleStyle.Render("  Key Dependencies"))
	fmt.Println(strings.Repeat("─", 50))

	keyDeps := getKeyDependencies(project, result.Parser)
	for _, kd := range keyDeps {
		printRow(kd.name, kd.status)
	}

	if locStats != nil {
		fmt.Println()
		fmt.Println(titleStyle.Render("  Code Statistics"))
		fmt.Println(strings.Repeat("─", 50))

		printRow("Total Files", countStyle.Render(fmt.Sprintf("%d", locStats.TotalFiles)))
		printRow("Lines of Code", countStyle.Render(fmt.Sprintf("%d", locStats.TotalCode)))
		printRow("Comments", countStyle.Render(fmt.Sprintf("%d", locStats.TotalComments)))
		printRow("Blank Lines", countStyle.Render(fmt.Sprintf("%d", locStats.TotalBlanks)))
	}

	fmt.Println()

	return nil
}

type keyDep struct {
	name   string
	status string
}

func getKeyDependencies(project *buildtool.Project, parser buildtool.Parser) []keyDep {
	var deps []keyDep

	checkStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	crossStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	check := checkStyle.Render("✓")
	cross := crossStyle.Render("–")

	if parser.HasSpringWeb(project) {
		deps = append(deps, keyDep{"Spring Web", check})
	} else {
		deps = append(deps, keyDep{"Spring Web", cross})
	}

	if parser.HasSpringDataJpa(project) {
		deps = append(deps, keyDep{"Spring Data JPA", check})
	} else {
		deps = append(deps, keyDep{"Spring Data JPA", cross})
	}

	if parser.HasLombok(project) {
		deps = append(deps, keyDep{"Lombok", check})
	} else {
		deps = append(deps, keyDep{"Lombok", cross})
	}

	if parser.HasValidation(project) {
		deps = append(deps, keyDep{"Validation", check})
	} else {
		deps = append(deps, keyDep{"Validation", cross})
	}

	if parser.HasMapStruct(project) {
		deps = append(deps, keyDep{"MapStruct", check})
	} else {
		deps = append(deps, keyDep{"MapStruct", cross})
	}

	if parser.HasDependency(project, "org.springframework.boot", "spring-boot-starter-security") {
		deps = append(deps, keyDep{"Spring Security", check})
	} else {
		deps = append(deps, keyDep{"Spring Security", cross})
	}

	return deps
}

func printRow(label, value string) {
	fmt.Printf("  %s %s\n",
		labelStyle.Render(fmt.Sprintf("%-18s", label+":")),
		valueStyle.Render(value))
}

func countByPrefix(deps []buildtool.Dependency, prefix string) int {
	count := 0
	for _, d := range deps {
		if strings.HasPrefix(d.ArtifactId, prefix) {
			count++
		}
	}
	return count
}

func countByScope(deps []buildtool.Dependency, scope string) int {
	count := 0
	for _, d := range deps {
		if strings.EqualFold(d.Scope, scope) {
			count++
		}
	}
	return count
}

func printJSON(project *buildtool.Project, result *buildtool.DetectionResult, locStats *stats.ProjectStats) error {
	cwd, _ := os.Getwd()
	projectName := filepath.Base(cwd)

	locJSON := ""
	if locStats != nil {
		locJSON = fmt.Sprintf(`,
  "codeStats": {
    "totalFiles": %d,
    "linesOfCode": %d,
    "comments": %d,
    "blankLines": %d
  }`,
			locStats.TotalFiles,
			locStats.TotalCode,
			locStats.TotalComments,
			locStats.TotalBlanks)
	}

	fmt.Printf(`{
  "name": "%s",
  "groupId": "%s",
  "artifactId": "%s",
  "version": "%s",
  "description": "%s",
  "buildTool": "%s",
  "buildFile": "%s",
  "javaVersion": "%s",
  "springBootVersion": "%s",
  "packaging": "%s",
  "dependencyCount": %d%s
}
`,
		projectName,
		project.GroupId,
		project.ArtifactId,
		project.Version,
		project.Description,
		result.BuildTool,
		filepath.Base(result.FilePath),
		project.JavaVersion,
		project.SpringBootVersion,
		project.Packaging,
		len(project.Dependencies),
		locJSON)

	return nil
}
