package add

import (
	"fmt"
	"os"
	"strings"

	"github.com/KashifKhn/haft/internal/buildtool"
	_ "github.com/KashifKhn/haft/internal/gradle"
	"github.com/KashifKhn/haft/internal/logger"
	_ "github.com/KashifKhn/haft/internal/maven"
	"github.com/KashifKhn/haft/internal/output"
	"github.com/KashifKhn/haft/internal/tui/components"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [dependency...]",
		Short: "Add dependencies to your project",
		Long: `Add dependencies to an existing Spring Boot project.

The add command modifies your build file (pom.xml or build.gradle) to add new dependencies. It supports:
  - Interactive mode: haft add (opens search picker)
  - Browse mode: haft add --browse (category browser)
  - Shortcuts: haft add lombok, haft add jpa
  - Maven coordinates: haft add org.example:my-lib
  - With version: haft add org.example:my-lib:1.0.0

Dependencies are auto-detected from the catalog or verified against Maven Central.`,
		Example: `  # Interactive search picker
  haft add

  # Browse by category
  haft add --browse

  # Add using shortcuts
  haft add lombok
  haft add jpa validation

  # Add using Maven coordinates
  haft add org.mapstruct:mapstruct:1.5.5.Final

  # Add with specific scope
  haft add h2 --scope test

  # List available shortcuts
  haft add --list

  # Output catalog as JSON
  haft add --list --json`,
		Args: cobra.ArbitraryArgs,
		RunE: runAdd,
	}

	cmd.Flags().String("scope", "", "Dependency scope (compile, runtime, test, provided)")
	cmd.Flags().String("version", "", "Override dependency version")
	cmd.Flags().Bool("list", false, "List available dependency shortcuts")
	cmd.Flags().BoolP("browse", "b", false, "Browse dependencies by category")
	cmd.Flags().Bool("no-interactive", false, "Skip interactive picker (requires dependency argument)")
	cmd.Flags().Bool("json", false, "Output as JSON (use with --list)")

	return cmd
}

func runAdd(cmd *cobra.Command, args []string) error {
	log := logger.Default()

	listFlag, _ := cmd.Flags().GetBool("list")
	jsonFlag, _ := cmd.Flags().GetBool("json")

	if listFlag {
		if jsonFlag {
			return printCatalogJSON()
		}
		printCatalog()
		return nil
	}

	browseFlag, _ := cmd.Flags().GetBool("browse")
	noInteractive, _ := cmd.Flags().GetBool("no-interactive")

	if len(args) == 0 && !browseFlag {
		if noInteractive {
			return fmt.Errorf("dependency argument required when using --no-interactive")
		}
		return runInteractivePicker(cmd)
	}

	if browseFlag {
		return runBrowser(cmd)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	fs := afero.NewOsFs()
	result, err := buildtool.Detect(cwd, fs)
	if err != nil {
		return fmt.Errorf("could not find build file: %w", err)
	}

	project, err := result.Parser.Parse(result.FilePath)
	if err != nil {
		return fmt.Errorf("could not parse %s: %w", result.FilePath, err)
	}

	scopeOverride, _ := cmd.Flags().GetString("scope")
	versionOverride, _ := cmd.Flags().GetString("version")

	addedCount := 0
	skippedCount := 0

	for _, arg := range args {
		deps, entryName, err := resolveDependency(arg)
		if err != nil {
			log.Error("Invalid dependency", "input", arg, "error", err.Error())
			continue
		}

		for _, dep := range deps {
			if scopeOverride != "" {
				dep.Scope = scopeOverride
			}
			if versionOverride != "" {
				dep.Version = versionOverride
			}

			if result.Parser.HasDependency(project, dep.GroupId, dep.ArtifactId) {
				log.Warning("Skipped (already exists)", "dependency", formatDependency(dep))
				skippedCount++
				continue
			}

			result.Parser.AddDependency(project, dep)
			if entryName != "" {
				log.Success("Added", "dependency", entryName, "artifact", dep.ArtifactId)
			} else {
				log.Success("Added", "dependency", formatDependency(dep))
			}
			addedCount++
		}
	}

	if addedCount == 0 {
		if skippedCount > 0 {
			log.Info("No new dependencies added (all already exist)")
		}
		return nil
	}

	if err := result.Parser.Write(result.FilePath, project); err != nil {
		return fmt.Errorf("could not write %s: %w", result.FilePath, err)
	}

	log.Success(fmt.Sprintf("Added %d dependencies to %s", addedCount, buildtool.GetBuildFileName(result.BuildTool)))
	return nil
}

func runInteractivePicker(cmd *cobra.Command) error {
	aliases, err := RunPicker()
	if err != nil {
		return err
	}

	if len(aliases) == 0 {
		return nil
	}

	return runAdd(cmd, aliases)
}

func runBrowser(cmd *cobra.Command) error {
	categories := buildBrowserCategories()
	model := components.NewDepPicker(components.DepPickerConfig{
		Label:      "Browse Dependencies",
		Categories: categories,
	})

	p := tea.NewProgram(browserWrapper{model: model})
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run browser: %w", err)
	}

	wrapper := finalModel.(browserWrapper)
	if wrapper.model.GoBack() || !wrapper.model.Submitted() {
		return nil
	}

	aliases := wrapper.model.Values()
	if len(aliases) == 0 {
		return nil
	}

	return runAdd(cmd, aliases)
}

type browserWrapper struct {
	model components.DepPickerModel
}

func (w browserWrapper) Init() tea.Cmd {
	return w.model.Init()
}

func (w browserWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return w, tea.Quit
		}
	}

	newModel, cmd := w.model.Update(msg)
	w.model = newModel

	if w.model.Submitted() || w.model.GoBack() {
		return w, tea.Quit
	}

	return w, cmd
}

func (w browserWrapper) View() string {
	return w.model.View()
}

func buildBrowserCategories() []components.DepCategory {
	catalogCategories := GetCatalogByCategory()
	categoryOrder := []string{
		"Web", "SQL", "NoSQL", "Security", "Messaging",
		"I/O", "Template Engines", "Ops", "Observability",
		"AI", "Cloud", "Notifications", "Payments", "Search",
		"Utilities", "Workflow", "Developer Tools", "Testing",
		"Maps", "Media", "Fintech", "Social", "Data",
		"Feature Flags", "Microservices", "Integration", "IoT",
		"DevOps", "Quality", "Caching", "Content", "Networking",
		"API", "Scheduling", "Logging",
	}

	var categories []components.DepCategory
	for _, catName := range categoryOrder {
		aliases, ok := catalogCategories[catName]
		if !ok {
			continue
		}

		var deps []components.DepItem
		for _, alias := range aliases {
			entry, _ := GetCatalogEntry(alias)
			deps = append(deps, components.DepItem{
				ID:          alias,
				Name:        entry.Name,
				Description: entry.Description,
			})
		}

		categories = append(categories, components.DepCategory{
			Name:         catName,
			Dependencies: deps,
		})
	}

	return categories
}

func resolveDependency(input string) ([]buildtool.Dependency, string, error) {
	if entry, ok := GetCatalogEntry(input); ok {
		return entry.Dependencies, entry.Name, nil
	}

	parts := strings.Split(input, ":")
	switch len(parts) {
	case 2:
		return verifyAndResolve(parts[0], parts[1], "")
	case 3:
		return verifyAndResolve(parts[0], parts[1], parts[2])
	default:
		return nil, "", fmt.Errorf("unknown shortcut '%s'. Use 'haft add --list' or specify as groupId:artifactId", input)
	}
}

func verifyAndResolve(groupId, artifactId, version string) ([]buildtool.Dependency, string, error) {
	client := NewMavenClient()
	artifact, err := client.VerifyDependency(groupId, artifactId)

	if err != nil {
		log := logger.Default()
		log.Warning("Could not verify dependency on Maven Central", "error", err.Error())
		return []buildtool.Dependency{{GroupId: groupId, ArtifactId: artifactId, Version: version}}, "", nil
	}

	if artifact == nil {
		return nil, "", fmt.Errorf("dependency '%s:%s' not found on Maven Central", groupId, artifactId)
	}

	if version == "" && artifact.LatestVersion != "" {
		version = artifact.LatestVersion
	}

	return []buildtool.Dependency{{GroupId: groupId, ArtifactId: artifactId, Version: version}}, "", nil
}

func formatDependency(dep buildtool.Dependency) string {
	if dep.Version != "" {
		return fmt.Sprintf("%s:%s:%s", dep.GroupId, dep.ArtifactId, dep.Version)
	}
	return fmt.Sprintf("%s:%s", dep.GroupId, dep.ArtifactId)
}

func printCatalog() {
	log := logger.Default()
	categories := GetCatalogByCategory()

	categoryOrder := []string{
		"Web",
		"SQL",
		"NoSQL",
		"Security",
		"Messaging",
		"I/O",
		"Template Engines",
		"Ops",
		"Observability",
		"AI",
		"Cloud",
		"Notifications",
		"Payments",
		"Search",
		"Utilities",
		"Workflow",
		"Developer Tools",
		"Testing",
		"Maps",
		"Media",
		"Fintech",
		"Social",
		"Data",
		"Feature Flags",
		"Microservices",
		"Integration",
		"IoT",
		"DevOps",
		"Quality",
		"Caching",
		"Content",
		"Networking",
		"API",
		"Scheduling",
		"Logging",
	}

	fmt.Println()
	fmt.Println("Available dependency shortcuts:")
	fmt.Println()

	for _, category := range categoryOrder {
		aliases, ok := categories[category]
		if !ok {
			continue
		}

		log.Info(category)
		for _, alias := range aliases {
			entry, _ := GetCatalogEntry(alias)
			fmt.Printf("  %-25s %s\n", alias, entry.Description)
		}
		fmt.Println()
	}

	fmt.Println("Usage: haft add                      (interactive picker)")
	fmt.Println("       haft add --browse             (browse by category)")
	fmt.Println("       haft add <shortcut>")
	fmt.Println("       haft add <groupId:artifactId>")
	fmt.Println("       haft add <groupId:artifactId:version>")
}

func printCatalogJSON() error {
	categories := GetCatalogByCategory()
	categoryOrder := []string{
		"Web", "SQL", "NoSQL", "Security", "Messaging",
		"I/O", "Template Engines", "Ops", "Observability",
		"AI", "Cloud", "Notifications", "Payments", "Search",
		"Utilities", "Workflow", "Developer Tools", "Testing",
		"Maps", "Media", "Fintech", "Social", "Data",
		"Feature Flags", "Microservices", "Integration", "IoT",
		"DevOps", "Quality", "Caching", "Content", "Networking",
		"API", "Scheduling", "Logging",
	}

	var outputCategories []output.CatalogCategory
	totalDeps := 0

	for _, catName := range categoryOrder {
		aliases, ok := categories[catName]
		if !ok {
			continue
		}

		var items []output.CatalogItem
		for _, alias := range aliases {
			entry, _ := GetCatalogEntry(alias)
			item := output.CatalogItem{
				Shortcut:    alias,
				Name:        entry.Name,
				Description: entry.Description,
			}
			if len(entry.Dependencies) > 0 {
				item.GroupID = entry.Dependencies[0].GroupId
				item.ArtifactID = entry.Dependencies[0].ArtifactId
			}
			items = append(items, item)
			totalDeps++
		}

		outputCategories = append(outputCategories, output.CatalogCategory{
			Name:         catName,
			Dependencies: items,
		})
	}

	return output.Success(output.CatalogOutput{
		Categories: outputCategories,
		Total:      totalDeps,
	})
}
