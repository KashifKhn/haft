package remove

import (
	"fmt"
	"os"
	"strings"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/KashifKhn/haft/internal/logger"
	_ "github.com/KashifKhn/haft/internal/maven"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [dependency...]",
		Short: "Remove dependencies from your project",
		Long: `Remove dependencies from an existing Spring Boot project.

The remove command modifies your build file (pom.xml or build.gradle) to remove dependencies. It supports:
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

	fs := afero.NewOsFs()
	result, err := buildtool.Detect(cwd, fs)
	if err != nil {
		return fmt.Errorf("could not find build file: %w", err)
	}

	project, err := result.Parser.Parse(result.FilePath)
	if err != nil {
		return fmt.Errorf("could not parse %s: %w", result.FilePath, err)
	}

	if len(project.Dependencies) == 0 {
		log.Info("No dependencies found in " + buildtool.GetBuildFileName(result.BuildTool))
		return nil
	}

	if len(args) == 0 {
		return runInteractivePicker(result.Parser, result.FilePath, project)
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

		if result.Parser.RemoveDependency(project, groupId, artifactId) {
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

	if err := result.Parser.Write(result.FilePath, project); err != nil {
		return fmt.Errorf("could not write %s: %w", result.FilePath, err)
	}

	log.Success(fmt.Sprintf("Removed %d dependencies from %s", removedCount, buildtool.GetBuildFileName(result.BuildTool)))
	return nil
}

func resolveInput(input string, project *buildtool.Project) (string, string) {
	if strings.Contains(input, ":") {
		parts := strings.Split(input, ":")
		if len(parts) >= 2 {
			return parts[0], parts[1]
		}
		return "", ""
	}

	for _, dep := range project.Dependencies {
		if dep.ArtifactId == input {
			return dep.GroupId, dep.ArtifactId
		}
		if strings.HasSuffix(dep.ArtifactId, input) {
			return dep.GroupId, dep.ArtifactId
		}
	}

	return "", ""
}

func runInteractivePicker(parser buildtool.Parser, filePath string, project *buildtool.Project) error {
	selected, err := RunRemovePicker(project.Dependencies)
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

	if err := parser.Write(filePath, project); err != nil {
		return fmt.Errorf("could not write build file: %w", err)
	}

	log.Success(fmt.Sprintf("Removed %d dependencies", removedCount))
	return nil
}
