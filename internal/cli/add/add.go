package add

import (
	"fmt"
	"os"
	"strings"

	"github.com/KashifKhn/haft/internal/logger"
	"github.com/KashifKhn/haft/internal/maven"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <dependency> [dependencies...]",
		Short: "Add dependencies to your project",
		Long: `Add dependencies to an existing Spring Boot project.

The add command modifies your pom.xml to add new dependencies. It supports:
  - Shortcuts: haft add lombok, haft add jpa
  - Maven coordinates: haft add org.example:my-lib
  - With version: haft add org.example:my-lib:1.0.0

Dependencies are auto-detected from the catalog or parsed as Maven coordinates.`,
		Example: `  # Add using shortcuts
  haft add lombok
  haft add jpa validation

  # Add using Maven coordinates
  haft add org.mapstruct:mapstruct:1.5.5.Final

  # Add with specific scope
  haft add h2 --scope test

  # List available shortcuts
  haft add --list`,
		Args: cobra.ArbitraryArgs,
		RunE: runAdd,
	}

	cmd.Flags().String("scope", "", "Dependency scope (compile, runtime, test, provided)")
	cmd.Flags().String("version", "", "Override dependency version")
	cmd.Flags().Bool("list", false, "List available dependency shortcuts")

	return cmd
}

func runAdd(cmd *cobra.Command, args []string) error {
	log := logger.Default()

	listFlag, _ := cmd.Flags().GetBool("list")
	if listFlag {
		printCatalog()
		return nil
	}

	if len(args) == 0 {
		return fmt.Errorf("at least one dependency is required\n\nUse 'haft add --list' to see available shortcuts")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	parser := maven.NewParser()
	pomPath, err := parser.FindPomXml(cwd)
	if err != nil {
		return fmt.Errorf("could not find pom.xml: %w", err)
	}

	project, err := parser.Parse(pomPath)
	if err != nil {
		return fmt.Errorf("could not parse pom.xml: %w", err)
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

			if parser.HasDependency(project, dep.GroupId, dep.ArtifactId) {
				log.Warning("Skipped (already exists)", "dependency", formatDependency(dep))
				skippedCount++
				continue
			}

			parser.AddDependency(project, dep)
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

	if err := parser.Write(pomPath, project); err != nil {
		return fmt.Errorf("could not write pom.xml: %w", err)
	}

	log.Success(fmt.Sprintf("Added %d dependencies to pom.xml", addedCount))
	return nil
}

func resolveDependency(input string) ([]maven.Dependency, string, error) {
	if entry, ok := GetCatalogEntry(input); ok {
		return entry.Dependencies, entry.Name, nil
	}

	parts := strings.Split(input, ":")
	switch len(parts) {
	case 2:
		return []maven.Dependency{{GroupId: parts[0], ArtifactId: parts[1]}}, "", nil
	case 3:
		return []maven.Dependency{{GroupId: parts[0], ArtifactId: parts[1], Version: parts[2]}}, "", nil
	default:
		return nil, "", fmt.Errorf("unknown shortcut '%s'. Use 'haft add --list' or specify as groupId:artifactId", input)
	}
}

func formatDependency(dep maven.Dependency) string {
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
		"Developer Tools",
		"Testing",
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

	fmt.Println("Usage: haft add <shortcut>")
	fmt.Println("       haft add <groupId:artifactId>")
	fmt.Println("       haft add <groupId:artifactId:version>")
}
