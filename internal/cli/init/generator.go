package init

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/KashifKhn/haft/internal/config"
	"github.com/KashifKhn/haft/internal/generator"
	"github.com/KashifKhn/haft/internal/output"
	"github.com/KashifKhn/haft/internal/tui/styles"
	"github.com/spf13/afero"
)

func generateProject(cfg ProjectConfig, projectDir string, jsonOutput bool) error {
	fs := afero.NewOsFs()

	if exists, _ := afero.DirExists(fs, projectDir); exists {
		if jsonOutput {
			return output.Error("PROJECT_EXISTS", fmt.Sprintf("directory %s already exists", projectDir))
		}
		return fmt.Errorf("directory %s already exists", projectDir)
	}

	var filesCreated []string

	if !jsonOutput {
		fmt.Println()
		fmt.Printf("  Creating project %s...\n", styles.Focused.Render(cfg.Name))
		fmt.Println()
	}

	engine := generator.NewEngine(fs)

	if !jsonOutput {
		fmt.Printf("  %s Creating directory structure\n", styles.CheckMark)
	}
	if err := createDirectoryStructure(engine, projectDir, cfg); err != nil {
		if jsonOutput {
			return output.Error("DIR_STRUCTURE_FAILED", "failed to create directory structure", err.Error())
		}
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	if !jsonOutput {
		fmt.Printf("  %s Generating project files\n", styles.CheckMark)
	}
	generatedFiles, err := generateProjectFiles(engine, projectDir, cfg)
	if err != nil {
		if jsonOutput {
			return output.Error("PROJECT_FILES_FAILED", "failed to generate project files", err.Error())
		}
		return fmt.Errorf("failed to generate project files: %w", err)
	}
	filesCreated = append(filesCreated, generatedFiles...)

	if cfg.BuildTool == "maven" {
		if !jsonOutput {
			fmt.Printf("  %s Adding Maven wrapper\n", styles.CheckMark)
		}
		wrapperFiles, err := copyMavenWrapper(engine, projectDir)
		if err != nil {
			if jsonOutput {
				return output.Error("MAVEN_WRAPPER_FAILED", "failed to copy Maven wrapper", err.Error())
			}
			return fmt.Errorf("failed to copy Maven wrapper: %w", err)
		}
		filesCreated = append(filesCreated, wrapperFiles...)
	} else if cfg.BuildTool == "gradle" || cfg.BuildTool == "gradle-kotlin" {
		if !jsonOutput {
			fmt.Printf("  %s Adding Gradle wrapper\n", styles.CheckMark)
		}
		wrapperFiles, err := copyGradleWrapper(engine, projectDir)
		if err != nil {
			if jsonOutput {
				return output.Error("GRADLE_WRAPPER_FAILED", "failed to copy Gradle wrapper", err.Error())
			}
			return fmt.Errorf("failed to copy Gradle wrapper: %w", err)
		}
		filesCreated = append(filesCreated, wrapperFiles...)
	}

	if !jsonOutput {
		fmt.Printf("  %s Writing configuration\n", styles.CheckMark)
	}
	if err := writeHaftConfig(fs, projectDir, cfg); err != nil {
		if jsonOutput {
			return output.Error("CONFIG_FAILED", "failed to write .haft.yaml", err.Error())
		}
		return fmt.Errorf("failed to write .haft.yaml: %w", err)
	}
	filesCreated = append(filesCreated, filepath.Join(projectDir, ".haft.yaml"))

	gitInitialized := false
	if cfg.InitGit {
		if !jsonOutput {
			fmt.Printf("  %s Initializing git repository\n", styles.CheckMark)
		}
		if err := initGitRepository(projectDir); err != nil {
			if !jsonOutput {
				fmt.Printf("  %s Failed to initialize git: %v\n", styles.ErrorText.Render("✗"), err)
			}
		} else {
			gitInitialized = true
		}
	}

	if jsonOutput {
		return output.Success(output.InitOutput{
			ProjectName:       cfg.Name,
			GroupID:           cfg.GroupId,
			ArtifactID:        cfg.ArtifactId,
			PackageName:       cfg.PackageName,
			JavaVersion:       cfg.JavaVersion,
			SpringBootVersion: cfg.SpringBootVersion,
			BuildTool:         cfg.BuildTool,
			Packaging:         cfg.Packaging,
			Dependencies:      cfg.Dependencies,
			ProjectDir:        projectDir,
			FilesCreated:      filesCreated,
			GitInitialized:    gitInitialized,
		})
	}

	fmt.Println()
	fmt.Printf("  %s Project created successfully!\n", styles.SuccessText.Render("✓"))
	fmt.Println()
	fmt.Println(styles.HelpText.Render("  Next steps:"))
	fmt.Println()
	fmt.Printf("    %s cd %s\n", styles.Arrow, styles.Focused.Render(projectDir))
	if cfg.BuildTool == "maven" {
		fmt.Printf("    %s ./mvnw spring-boot:run\n", styles.Arrow)
	} else {
		fmt.Printf("    %s ./gradlew bootRun\n", styles.Arrow)
	}
	fmt.Println()

	return nil
}

func createDirectoryStructure(engine *generator.Engine, projectDir string, cfg ProjectConfig) error {
	basePackagePath := strings.ReplaceAll(cfg.PackageName, ".", string(os.PathSeparator))

	dirs := []string{
		filepath.Join(projectDir, "src", "main", "java", basePackagePath),
		filepath.Join(projectDir, "src", "main", "java", basePackagePath, "controller"),
		filepath.Join(projectDir, "src", "main", "java", basePackagePath, "service"),
		filepath.Join(projectDir, "src", "main", "java", basePackagePath, "service", "impl"),
		filepath.Join(projectDir, "src", "main", "java", basePackagePath, "repository"),
		filepath.Join(projectDir, "src", "main", "java", basePackagePath, "entity"),
		filepath.Join(projectDir, "src", "main", "java", basePackagePath, "dto"),
		filepath.Join(projectDir, "src", "main", "java", basePackagePath, "mapper"),
		filepath.Join(projectDir, "src", "main", "java", basePackagePath, "exception"),
		filepath.Join(projectDir, "src", "main", "java", basePackagePath, "config"),
		filepath.Join(projectDir, "src", "main", "resources"),
		filepath.Join(projectDir, "src", "test", "java", basePackagePath),
		filepath.Join(projectDir, "src", "test", "resources"),
	}

	for _, dir := range dirs {
		if err := engine.GetFS().MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

func generateProjectFiles(engine *generator.Engine, projectDir string, cfg ProjectConfig) ([]string, error) {
	var filesCreated []string
	basePackagePath := strings.ReplaceAll(cfg.PackageName, ".", string(os.PathSeparator))
	applicationName := toPascalCase(cfg.ArtifactId)

	data := map[string]any{
		"Name":              cfg.Name,
		"GroupId":           cfg.GroupId,
		"ArtifactId":        cfg.ArtifactId,
		"Version":           "0.0.1-SNAPSHOT",
		"Description":       cfg.Description,
		"JavaVersion":       cfg.JavaVersion,
		"SpringBootVersion": cfg.SpringBootVersion,
		"BasePackage":       cfg.PackageName,
		"ApplicationName":   applicationName,
		"Packaging":         cfg.Packaging,
		"Dependencies":      buildDependencies(cfg.Dependencies),
		"HasLombok":         contains(cfg.Dependencies, "lombok"),
		"HasJpa":            contains(cfg.Dependencies, "data-jpa") || contains(cfg.Dependencies, "jpa"),
		"HasWeb":            contains(cfg.Dependencies, "web"),
		"HasSecurity":       contains(cfg.Dependencies, "security"),
		"HasValidation":     contains(cfg.Dependencies, "validation"),
	}

	if cfg.BuildTool == "maven" {
		filePath := filepath.Join(projectDir, "pom.xml")
		if err := engine.RenderAndWrite(
			"project/pom.xml.tmpl",
			filePath,
			data,
		); err != nil {
			return nil, err
		}
		filesCreated = append(filesCreated, filePath)
	} else if cfg.BuildTool == "gradle" {
		buildFile := filepath.Join(projectDir, "build.gradle")
		if err := engine.RenderAndWrite(
			"project/build.gradle.tmpl",
			buildFile,
			data,
		); err != nil {
			return nil, err
		}
		filesCreated = append(filesCreated, buildFile)

		settingsFile := filepath.Join(projectDir, "settings.gradle")
		if err := engine.RenderAndWrite(
			"project/settings.gradle.tmpl",
			settingsFile,
			data,
		); err != nil {
			return nil, err
		}
		filesCreated = append(filesCreated, settingsFile)
	} else if cfg.BuildTool == "gradle-kotlin" {
		buildFile := filepath.Join(projectDir, "build.gradle.kts")
		if err := engine.RenderAndWrite(
			"project/build.gradle.kts.tmpl",
			buildFile,
			data,
		); err != nil {
			return nil, err
		}
		filesCreated = append(filesCreated, buildFile)

		settingsFile := filepath.Join(projectDir, "settings.gradle.kts")
		if err := engine.RenderAndWrite(
			"project/settings.gradle.kts.tmpl",
			settingsFile,
			data,
		); err != nil {
			return nil, err
		}
		filesCreated = append(filesCreated, settingsFile)
	}

	appFile := filepath.Join(projectDir, "src", "main", "java", basePackagePath, applicationName+"Application.java")
	if err := engine.RenderAndWrite(
		"project/Application.java.tmpl",
		appFile,
		data,
	); err != nil {
		return nil, err
	}
	filesCreated = append(filesCreated, appFile)

	configFile := "application.properties"
	configTemplate := "project/application.properties.tmpl"
	if cfg.ConfigFormat == "yaml" {
		configFile = "application.yml"
		configTemplate = "project/application.yml.tmpl"
	}

	configPath := filepath.Join(projectDir, "src", "main", "resources", configFile)
	if err := engine.RenderAndWrite(
		configTemplate,
		configPath,
		data,
	); err != nil {
		return nil, err
	}
	filesCreated = append(filesCreated, configPath)

	testFile := filepath.Join(projectDir, "src", "test", "java", basePackagePath, applicationName+"ApplicationTests.java")
	if err := engine.RenderAndWrite(
		"project/ApplicationTests.java.tmpl",
		testFile,
		data,
	); err != nil {
		return nil, err
	}
	filesCreated = append(filesCreated, testFile)

	gitignoreContent, err := engine.ReadTemplateFile("project/gitignore.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to read gitignore template: %w", err)
	}
	gitignorePath := filepath.Join(projectDir, ".gitignore")
	if err := engine.WriteFileWithPerm(gitignorePath, gitignoreContent, 0644); err != nil {
		return nil, fmt.Errorf("failed to write .gitignore: %w", err)
	}
	filesCreated = append(filesCreated, gitignorePath)

	return filesCreated, nil
}

func writeHaftConfig(fs afero.Fs, projectDir string, cfg ProjectConfig) error {
	projectCfg := config.ProjectConfig{
		Version: "1",
		Project: config.ProjectSettings{
			Name:        cfg.Name,
			Group:       cfg.GroupId,
			Artifact:    cfg.ArtifactId,
			Description: cfg.Description,
			Package:     cfg.PackageName,
		},
		Spring: config.SpringSettings{
			Version: cfg.SpringBootVersion,
		},
		Java: config.JavaSettings{
			Version: cfg.JavaVersion,
		},
		Build: config.BuildSettings{
			Tool: cfg.BuildTool,
		},
		Architecture: config.ArchSettings{
			Style: cfg.Architecture,
		},
		Database: config.DatabaseSettings{
			Type: "h2",
		},
		Generators: config.GeneratorSettings{
			DTO: config.DTOSettings{
				Style: "record",
			},
			Tests: config.TestSettings{
				Enabled: true,
			},
		},
	}

	homeDir, _ := os.UserHomeDir()
	cm := config.NewConfigManager(fs, projectDir, homeDir)
	return cm.SaveProjectConfig(&projectCfg)
}

func copyMavenWrapper(engine *generator.Engine, projectDir string) ([]string, error) {
	var filesCreated []string

	wrapperDir := filepath.Join(projectDir, ".mvn", "wrapper")
	if err := engine.GetFS().MkdirAll(wrapperDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create .mvn/wrapper directory: %w", err)
	}

	propsContent, err := engine.ReadTemplateFile("wrapper/maven-wrapper.properties")
	if err != nil {
		return nil, fmt.Errorf("failed to read maven-wrapper.properties: %w", err)
	}
	propsPath := filepath.Join(wrapperDir, "maven-wrapper.properties")
	if err := engine.WriteFileWithPerm(propsPath, propsContent, 0644); err != nil {
		return nil, fmt.Errorf("failed to write maven-wrapper.properties: %w", err)
	}
	filesCreated = append(filesCreated, propsPath)

	mvnwContent, err := engine.ReadTemplateFile("wrapper/mvnw")
	if err != nil {
		return nil, fmt.Errorf("failed to read mvnw: %w", err)
	}
	mvnwPath := filepath.Join(projectDir, "mvnw")
	if err := engine.WriteFileWithPerm(mvnwPath, mvnwContent, 0755); err != nil {
		return nil, fmt.Errorf("failed to write mvnw: %w", err)
	}
	filesCreated = append(filesCreated, mvnwPath)

	mvnwCmdContent, err := engine.ReadTemplateFile("wrapper/mvnw.cmd")
	if err != nil {
		return nil, fmt.Errorf("failed to read mvnw.cmd: %w", err)
	}
	mvnwCmdPath := filepath.Join(projectDir, "mvnw.cmd")
	if err := engine.WriteFileWithPerm(mvnwCmdPath, mvnwCmdContent, 0644); err != nil {
		return nil, fmt.Errorf("failed to write mvnw.cmd: %w", err)
	}
	filesCreated = append(filesCreated, mvnwCmdPath)

	return filesCreated, nil
}

func copyGradleWrapper(engine *generator.Engine, projectDir string) ([]string, error) {
	var filesCreated []string

	wrapperDir := filepath.Join(projectDir, "gradle", "wrapper")
	if err := engine.GetFS().MkdirAll(wrapperDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create gradle/wrapper directory: %w", err)
	}

	propsContent, err := engine.ReadTemplateFile("wrapper/gradle-wrapper.properties")
	if err != nil {
		return nil, fmt.Errorf("failed to read gradle-wrapper.properties: %w", err)
	}
	propsPath := filepath.Join(wrapperDir, "gradle-wrapper.properties")
	if err := engine.WriteFileWithPerm(propsPath, propsContent, 0644); err != nil {
		return nil, fmt.Errorf("failed to write gradle-wrapper.properties: %w", err)
	}
	filesCreated = append(filesCreated, propsPath)

	gradlewContent, err := engine.ReadTemplateFile("wrapper/gradlew")
	if err != nil {
		return nil, fmt.Errorf("failed to read gradlew: %w", err)
	}
	gradlewPath := filepath.Join(projectDir, "gradlew")
	if err := engine.WriteFileWithPerm(gradlewPath, gradlewContent, 0755); err != nil {
		return nil, fmt.Errorf("failed to write gradlew: %w", err)
	}
	filesCreated = append(filesCreated, gradlewPath)

	gradlewBatContent, err := engine.ReadTemplateFile("wrapper/gradlew.bat")
	if err != nil {
		return nil, fmt.Errorf("failed to read gradlew.bat: %w", err)
	}
	gradlewBatPath := filepath.Join(projectDir, "gradlew.bat")
	if err := engine.WriteFileWithPerm(gradlewBatPath, gradlewBatContent, 0644); err != nil {
		return nil, fmt.Errorf("failed to write gradlew.bat: %w", err)
	}
	filesCreated = append(filesCreated, gradlewBatPath)

	return filesCreated, nil
}

func toPascalCase(s string) string {
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, "")
}

func initGitRepository(projectDir string) error {
	gitInit := exec.Command("git", "init")
	gitInit.Dir = projectDir
	if err := gitInit.Run(); err != nil {
		return fmt.Errorf("git init failed: %w", err)
	}

	gitAdd := exec.Command("git", "add", ".")
	gitAdd.Dir = projectDir
	if err := gitAdd.Run(); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}

	gitCommit := exec.Command("git", "commit", "-m", "Initial commit")
	gitCommit.Dir = projectDir
	if err := gitCommit.Run(); err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}

	return nil
}
