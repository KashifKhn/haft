package remove

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
		Use:   "remove [dependency...]",
		Short: "Remove dependencies from your project",
		Long: `Remove dependencies from an existing Spring Boot project.

The remove command modifies your pom.xml to remove dependencies. It supports:
  - Interactive mode: haft remove (opens picker with current dependencies)
  - By artifact: haft remove lombok
  - By coordinates: haft remove org.projectlombok:lombok
  - Multiple: haft remove lombok validation jpa`,
		Example: `  # Interactive dependency picker
  haft remove

  # Remove by artifact name
  haft remove lombok
  haft remove spring-boot-starter-web

  # Remove by coordinates
  haft remove org.projectlombok:lombok

  # Remove multiple
  haft remove lombok validation h2`,
		Aliases: []string{"rm"},
		Args:    cobra.ArbitraryArgs,
		RunE:    runRemove,
	}

	return cmd
}

func runRemove(cmd *cobra.Command, args []string) error {
	log := logger.Default()

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

	if project.Dependencies == nil || len(project.Dependencies.Dependency) == 0 {
		log.Info("No dependencies found in pom.xml")
		return nil
	}

	if len(args) == 0 {
		return runInteractivePicker(cmd, parser, pomPath, project)
	}

	removedCount := 0
	notFoundCount := 0

	for _, arg := range args {
		groupId, artifactId := resolveInput(arg, project)

		if groupId == "" || artifactId == "" {
			log.Warning("Dependency not found", "input", arg)
			notFoundCount++
			continue
		}

		if parser.RemoveDependency(project, groupId, artifactId) {
			log.Success("Removed", "dependency", fmt.Sprintf("%s:%s", groupId, artifactId))
			removedCount++
		} else {
			log.Warning("Dependency not found", "input", arg)
			notFoundCount++
		}
	}

	if removedCount == 0 {
		if notFoundCount > 0 {
			log.Info("No dependencies removed (none found)")
		}
		return nil
	}

	if err := parser.Write(pomPath, project); err != nil {
		return fmt.Errorf("could not write pom.xml: %w", err)
	}

	log.Success(fmt.Sprintf("Removed %d dependencies from pom.xml", removedCount))
	return nil
}

func resolveInput(input string, project *maven.Project) (string, string) {
	if strings.Contains(input, ":") {
		parts := strings.Split(input, ":")
		if len(parts) >= 2 {
			return parts[0], parts[1]
		}
		return "", ""
	}

	for _, dep := range project.Dependencies.Dependency {
		if dep.ArtifactId == input {
			return dep.GroupId, dep.ArtifactId
		}
		if strings.HasSuffix(dep.ArtifactId, input) {
			return dep.GroupId, dep.ArtifactId
		}
	}

	return "", ""
}

func runInteractivePicker(cmd *cobra.Command, parser *maven.Parser, pomPath string, project *maven.Project) error {
	selected, err := RunRemovePicker(project.Dependencies.Dependency)
	if err != nil {
		return err
	}

	if len(selected) == 0 {
		return nil
	}

	log := logger.Default()
	removedCount := 0

	for _, dep := range selected {
		if parser.RemoveDependency(project, dep.GroupId, dep.ArtifactId) {
			log.Success("Removed", "dependency", fmt.Sprintf("%s:%s", dep.GroupId, dep.ArtifactId))
			removedCount++
		}
	}

	if removedCount == 0 {
		return nil
	}

	if err := parser.Write(pomPath, project); err != nil {
		return fmt.Errorf("could not write pom.xml: %w", err)
	}

	log.Success(fmt.Sprintf("Removed %d dependencies from pom.xml", removedCount))
	return nil
}
