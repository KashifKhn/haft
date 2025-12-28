package remove

import (
	"fmt"
	"os"
	"strings"

	"github.com/KashifKhn/haft/internal/buildtool"
	_ "github.com/KashifKhn/haft/internal/gradle"
	"github.com/KashifKhn/haft/internal/logger"
	_ "github.com/KashifKhn/haft/internal/maven"
	"github.com/KashifKhn/haft/internal/output"
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
  haft remove lombok validation h2

  # Output as JSON
  haft remove lombok --json`,
		Aliases: []string{"rm"},
		Args:    cobra.ArbitraryArgs,
		RunE:    runRemove,
	}

	cmd.Flags().Bool("no-interactive", false, "Skip interactive picker (requires dependency argument)")
	cmd.Flags().Bool("json", false, "Output result as JSON")

	return cmd
}

func runRemove(cmd *cobra.Command, args []string) error {
	log := logger.Default()
	jsonFlag, _ := cmd.Flags().GetBool("json")

	cwd, err := os.Getwd()
	if err != nil {
		if jsonFlag {
			return output.Error("CWD_ERROR", "could not get current directory", err.Error())
		}
		return err
	}

	fs := afero.NewOsFs()
	result, err := buildtool.Detect(cwd, fs)
	if err != nil {
		if jsonFlag {
			return output.Error("NO_BUILD_FILE", "could not find build file", err.Error())
		}
		return fmt.Errorf("could not find build file: %w", err)
	}

	project, err := result.Parser.Parse(result.FilePath)
	if err != nil {
		if jsonFlag {
			return output.Error("PARSE_ERROR", fmt.Sprintf("could not parse %s", result.FilePath), err.Error())
		}
		return fmt.Errorf("could not parse %s: %w", result.FilePath, err)
	}

	if len(project.Dependencies) == 0 {
		if jsonFlag {
			return output.Success(output.AddRemoveResult{
				Action:  "remove",
				Removed: []string{},
				Skipped: []string{},
			})
		}
		log.Info("No dependencies found in " + buildtool.GetBuildFileName(result.BuildTool))
		return nil
	}

	noInteractive, _ := cmd.Flags().GetBool("no-interactive")

	if len(args) == 0 {
		if noInteractive {
			if jsonFlag {
				return output.Error("MISSING_ARGUMENT", "dependency argument required when using --no-interactive")
			}
			return fmt.Errorf("dependency argument required when using --no-interactive")
		}
		return runInteractivePicker(result.Parser, result.FilePath, project)
	}

	var removed, notFound []string

	for _, arg := range args {
		groupId, artifactId := resolveInput(arg, project)

		if groupId == "" || artifactId == "" {
			if !jsonFlag {
				log.Warning("Dependency not found", "input", arg)
			}
			notFound = append(notFound, arg)
			continue
		}

		if result.Parser.RemoveDependency(project, groupId, artifactId) {
			if !jsonFlag {
				log.Success("Removed", "dependency", fmt.Sprintf("%s:%s", groupId, artifactId))
			}
			removed = append(removed, fmt.Sprintf("%s:%s", groupId, artifactId))
		} else {
			if !jsonFlag {
				log.Warning("Dependency not found", "input", arg)
			}
			notFound = append(notFound, arg)
		}
	}

	if len(removed) == 0 {
		if jsonFlag {
			return output.Success(output.AddRemoveResult{
				Action:  "remove",
				Removed: removed,
				Skipped: notFound,
			})
		}
		if len(notFound) > 0 {
			log.Info("No dependencies removed (none found)")
		}
		return nil
	}

	if err := result.Parser.Write(result.FilePath, project); err != nil {
		if jsonFlag {
			return output.Error("WRITE_ERROR", fmt.Sprintf("could not write %s", result.FilePath), err.Error())
		}
		return fmt.Errorf("could not write %s: %w", result.FilePath, err)
	}

	if jsonFlag {
		return output.Success(output.AddRemoveResult{
			Action:  "remove",
			Removed: removed,
			Skipped: notFound,
		})
	}

	log.Success(fmt.Sprintf("Removed %d dependencies from %s", len(removed), buildtool.GetBuildFileName(result.BuildTool)))
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
