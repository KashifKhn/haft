package docker

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/KashifKhn/haft/internal/generator"
	"github.com/KashifKhn/haft/internal/logger"
	"github.com/KashifKhn/haft/internal/output"
	"github.com/KashifKhn/haft/internal/tui/components"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type DatabaseInfo struct {
	Type       string
	Image      string
	Port       int
	EnvPrefix  string
	SpringURL  string
	EnvVars    map[string]string
	VolumePath string
}

var databaseDrivers = map[string]DatabaseInfo{
	"postgresql": {
		Type:       "postgres",
		Image:      "postgres:16-alpine",
		Port:       5432,
		EnvPrefix:  "POSTGRES",
		SpringURL:  "jdbc:postgresql://postgres:5432/${APP_NAME}",
		VolumePath: "/var/lib/postgresql/data",
		EnvVars: map[string]string{
			"POSTGRES_DB":       "${APP_NAME}",
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
		},
	},
	"mysql-connector-j": {
		Type:       "mysql",
		Image:      "mysql:8",
		Port:       3306,
		EnvPrefix:  "MYSQL",
		SpringURL:  "jdbc:mysql://mysql:3306/${APP_NAME}",
		VolumePath: "/var/lib/mysql",
		EnvVars: map[string]string{
			"MYSQL_DATABASE":      "${APP_NAME}",
			"MYSQL_ROOT_PASSWORD": "root",
			"MYSQL_USER":          "mysql",
			"MYSQL_PASSWORD":      "mysql",
		},
	},
	"mysql": {
		Type:       "mysql",
		Image:      "mysql:8",
		Port:       3306,
		EnvPrefix:  "MYSQL",
		SpringURL:  "jdbc:mysql://mysql:3306/${APP_NAME}",
		VolumePath: "/var/lib/mysql",
		EnvVars: map[string]string{
			"MYSQL_DATABASE":      "${APP_NAME}",
			"MYSQL_ROOT_PASSWORD": "root",
			"MYSQL_USER":          "mysql",
			"MYSQL_PASSWORD":      "mysql",
		},
	},
	"mariadb": {
		Type:       "mariadb",
		Image:      "mariadb:11",
		Port:       3306,
		EnvPrefix:  "MARIADB",
		SpringURL:  "jdbc:mariadb://mariadb:3306/${APP_NAME}",
		VolumePath: "/var/lib/mysql",
		EnvVars: map[string]string{
			"MARIADB_DATABASE":      "${APP_NAME}",
			"MARIADB_ROOT_PASSWORD": "root",
			"MARIADB_USER":          "mariadb",
			"MARIADB_PASSWORD":      "mariadb",
		},
	},
	"data-mongodb": {
		Type:       "mongodb",
		Image:      "mongo:7",
		Port:       27017,
		EnvPrefix:  "MONGO",
		SpringURL:  "mongodb://mongodb:27017/${APP_NAME}",
		VolumePath: "/data/db",
		EnvVars: map[string]string{
			"MONGO_INITDB_DATABASE": "${APP_NAME}",
		},
	},
	"mongodb": {
		Type:       "mongodb",
		Image:      "mongo:7",
		Port:       27017,
		EnvPrefix:  "MONGO",
		SpringURL:  "mongodb://mongodb:27017/${APP_NAME}",
		VolumePath: "/data/db",
		EnvVars: map[string]string{
			"MONGO_INITDB_DATABASE": "${APP_NAME}",
		},
	},
	"data-redis": {
		Type:       "redis",
		Image:      "redis:7-alpine",
		Port:       6379,
		EnvPrefix:  "REDIS",
		SpringURL:  "redis://redis:6379",
		VolumePath: "/data",
		EnvVars:    map[string]string{},
	},
	"data-cassandra": {
		Type:       "cassandra",
		Image:      "cassandra:4",
		Port:       9042,
		EnvPrefix:  "CASSANDRA",
		SpringURL:  "cassandra://cassandra:9042/${APP_NAME}",
		VolumePath: "/var/lib/cassandra",
		EnvVars: map[string]string{
			"CASSANDRA_CLUSTER_NAME": "${APP_NAME}-cluster",
		},
	},
}

var databaseChoices = []struct {
	Label       string
	Value       string
	Description string
}{
	{"PostgreSQL", "postgresql", "Recommended for most applications"},
	{"MySQL", "mysql", "Popular relational database"},
	{"MariaDB", "mariadb", "MySQL-compatible database"},
	{"MongoDB", "mongodb", "Document-oriented NoSQL database"},
	{"Redis", "redis", "In-memory data store (caching)"},
	{"None", "none", "No database service"},
}

type DockerConfig struct {
	AppName         string
	JavaVersion     string
	Port            int
	BuildTool       buildtool.Type
	HasWrapper      bool
	DatabaseType    string
	DatabaseInfo    *DatabaseInfo
	GenerateCompose bool
}

type DockerOutput struct {
	Generated []string `json:"generated"`
	Skipped   []string `json:"skipped"`
	Config    struct {
		AppName      string `json:"appName"`
		JavaVersion  string `json:"javaVersion"`
		Port         int    `json:"port"`
		BuildTool    string `json:"buildTool"`
		DatabaseType string `json:"databaseType,omitempty"`
	} `json:"config"`
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dockerize",
		Aliases: []string{"docker"},
		Short:   "Generate Docker configuration files",
		Long: `Generate optimized Docker configuration for your Spring Boot project.

This command generates:
  - Dockerfile: Multi-stage build for optimized image size
  - docker-compose.yml: Service orchestration with database
  - .dockerignore: Exclude unnecessary files from build context

The command automatically detects:
  - Build tool (Maven/Gradle) for appropriate Dockerfile
  - Java version for correct base image
  - Database dependencies for docker-compose services
  - Application name and port from configuration

Database services supported:
  - PostgreSQL (postgres:16-alpine)
  - MySQL (mysql:8)
  - MariaDB (mariadb:11)
  - MongoDB (mongo:7)
  - Redis (redis:7-alpine)
  - Cassandra (cassandra:4)`,
		Example: `  # Generate Docker files with auto-detection
  haft dockerize

  # Generate only Dockerfile (no docker-compose)
  haft dockerize --no-compose

  # Specify database type explicitly
  haft dockerize --db postgresql

  # Override application port
  haft dockerize --port 9000

  # Override Java version
  haft dockerize --java 21

  # Non-interactive mode
  haft dockerize --no-interactive

  # Output as JSON
  haft dockerize --json`,
		RunE: runDockerize,
	}

	cmd.Flags().IntP("port", "p", 0, "Application port (default: auto-detect or 8080)")
	cmd.Flags().StringP("java", "j", "", "Java version (default: auto-detect from build file)")
	cmd.Flags().String("db", "", "Database type (postgresql, mysql, mariadb, mongodb, redis, none)")
	cmd.Flags().Bool("no-compose", false, "Skip docker-compose.yml generation")
	cmd.Flags().Bool("no-interactive", false, "Skip interactive prompts")
	cmd.Flags().Bool("json", false, "Output result as JSON")

	return cmd
}

func runDockerize(cmd *cobra.Command, args []string) error {
	noInteractive, _ := cmd.Flags().GetBool("no-interactive")
	noCompose, _ := cmd.Flags().GetBool("no-compose")
	jsonOutput, _ := cmd.Flags().GetBool("json")
	portFlag, _ := cmd.Flags().GetInt("port")
	javaFlag, _ := cmd.Flags().GetString("java")
	dbFlag, _ := cmd.Flags().GetString("db")

	log := logger.Default()
	fs := afero.NewOsFs()

	cwd, err := os.Getwd()
	if err != nil {
		if jsonOutput {
			return output.Error("DIRECTORY_ERROR", "Could not get current directory", err.Error())
		}
		return fmt.Errorf("could not get current directory: %w", err)
	}

	result, err := buildtool.Detect(cwd, fs)
	if err != nil {
		if jsonOutput {
			return output.Error("BUILD_FILE_ERROR", "No build file found", err.Error())
		}
		return fmt.Errorf("no build file found: %w. Run this command from a Maven or Gradle project root", err)
	}

	project, err := result.Parser.Parse(result.FilePath)
	if err != nil {
		if jsonOutput {
			return output.Error("PARSE_ERROR", "Could not parse build file", err.Error())
		}
		return fmt.Errorf("could not parse build file: %w", err)
	}

	cfg := DockerConfig{
		AppName:         detectAppName(project, cwd),
		JavaVersion:     detectJavaVersion(project, result.Parser, javaFlag),
		Port:            detectPort(cwd, portFlag),
		BuildTool:       result.BuildTool,
		HasWrapper:      hasWrapper(cwd, result.BuildTool),
		GenerateCompose: !noCompose,
	}

	detectedDB := detectDatabaseFromDependencies(project.Dependencies)
	hasJPA := hasJPADependency(project.Dependencies)

	if dbFlag != "" {
		if dbFlag == "none" {
			cfg.DatabaseType = ""
			cfg.DatabaseInfo = nil
		} else if info, ok := databaseDrivers[dbFlag]; ok {
			cfg.DatabaseType = dbFlag
			cfg.DatabaseInfo = &info
		} else {
			if jsonOutput {
				return output.Error("INVALID_DB", "Invalid database type", "Valid options: postgresql, mysql, mariadb, mongodb, redis, none")
			}
			return fmt.Errorf("invalid database type: %s. Valid options: postgresql, mysql, mariadb, mongodb, redis, none", dbFlag)
		}
	} else if detectedDB != "" {
		info := databaseDrivers[detectedDB]
		cfg.DatabaseType = detectedDB
		cfg.DatabaseInfo = &info
	} else if hasJPA && !noInteractive {
		selectedDB, err := runDatabasePicker()
		if err != nil {
			if strings.Contains(err.Error(), "cancelled") {
				if jsonOutput {
					return output.Error("CANCELLED", "Operation cancelled by user")
				}
				return fmt.Errorf("operation cancelled")
			}
			if jsonOutput {
				return output.Error("PICKER_ERROR", "Database picker failed", err.Error())
			}
			return err
		}
		if selectedDB != "none" && selectedDB != "" {
			info := databaseDrivers[selectedDB]
			cfg.DatabaseType = selectedDB
			cfg.DatabaseInfo = &info
		}
	}

	if !jsonOutput {
		log.Info("Detected configuration",
			"buildTool", cfg.BuildTool.DisplayName(),
			"java", cfg.JavaVersion,
			"port", cfg.Port,
			"app", cfg.AppName,
		)
		if cfg.DatabaseType != "" {
			log.Info("Database", "type", cfg.DatabaseInfo.Type, "image", cfg.DatabaseInfo.Image)
		}
	}

	return generateDockerFiles(cfg, cwd, jsonOutput)
}

func generateDockerFiles(cfg DockerConfig, cwd string, jsonOutput bool) error {
	log := logger.Default()
	fs := afero.NewOsFs()
	engine := generator.NewEngineWithLoader(fs, cwd)

	generated := []string{}
	skipped := []string{}

	data := buildDockerTemplateData(cfg)

	dockerfilePath := filepath.Join(cwd, "Dockerfile")
	if engine.FileExists(dockerfilePath) {
		if !jsonOutput {
			log.Warning("File exists, skipping", "file", "Dockerfile")
		}
		skipped = append(skipped, "Dockerfile")
	} else {
		templateName := "docker/Dockerfile.maven.tmpl"
		if cfg.BuildTool.IsGradle() {
			templateName = "docker/Dockerfile.gradle.tmpl"
		}

		if err := engine.RenderAndWrite(templateName, dockerfilePath, data); err != nil {
			if jsonOutput {
				return output.Error("GENERATION_ERROR", "Failed to generate Dockerfile", err.Error())
			}
			return fmt.Errorf("failed to generate Dockerfile: %w", err)
		}
		if !jsonOutput {
			log.Success("Created", "file", "Dockerfile")
		}
		generated = append(generated, "Dockerfile")
	}

	dockerignorePath := filepath.Join(cwd, ".dockerignore")
	if engine.FileExists(dockerignorePath) {
		if !jsonOutput {
			log.Warning("File exists, skipping", "file", ".dockerignore")
		}
		skipped = append(skipped, ".dockerignore")
	} else {
		if err := engine.RenderAndWrite("docker/dockerignore.tmpl", dockerignorePath, data); err != nil {
			if jsonOutput {
				return output.Error("GENERATION_ERROR", "Failed to generate .dockerignore", err.Error())
			}
			return fmt.Errorf("failed to generate .dockerignore: %w", err)
		}
		if !jsonOutput {
			log.Success("Created", "file", ".dockerignore")
		}
		generated = append(generated, ".dockerignore")
	}

	if cfg.GenerateCompose {
		composePath := filepath.Join(cwd, "docker-compose.yml")
		if engine.FileExists(composePath) {
			if !jsonOutput {
				log.Warning("File exists, skipping", "file", "docker-compose.yml")
			}
			skipped = append(skipped, "docker-compose.yml")
		} else {
			if err := engine.RenderAndWrite("docker/docker-compose.yml.tmpl", composePath, data); err != nil {
				if jsonOutput {
					return output.Error("GENERATION_ERROR", "Failed to generate docker-compose.yml", err.Error())
				}
				return fmt.Errorf("failed to generate docker-compose.yml: %w", err)
			}
			if !jsonOutput {
				log.Success("Created", "file", "docker-compose.yml")
			}
			generated = append(generated, "docker-compose.yml")
		}
	}

	if !jsonOutput {
		if len(generated) > 0 {
			log.Success(fmt.Sprintf("Generated %d Docker file(s)", len(generated)))
		}
		if len(skipped) > 0 {
			log.Info(fmt.Sprintf("Skipped %d existing file(s)", len(skipped)))
		}

		fmt.Println()
		log.Info("Next steps:")
		fmt.Println("  1. Review and customize the generated files")
		fmt.Println("  2. Build the image: docker build -t " + cfg.AppName + " .")
		if cfg.GenerateCompose {
			fmt.Println("  3. Start services: docker compose up -d")
		}
	}

	if jsonOutput {
		out := DockerOutput{
			Generated: generated,
			Skipped:   skipped,
		}
		out.Config.AppName = cfg.AppName
		out.Config.JavaVersion = cfg.JavaVersion
		out.Config.Port = cfg.Port
		out.Config.BuildTool = string(cfg.BuildTool)
		if cfg.DatabaseType != "" {
			out.Config.DatabaseType = cfg.DatabaseInfo.Type
		}
		return output.Success(out)
	}

	return nil
}

func buildDockerTemplateData(cfg DockerConfig) map[string]any {
	data := map[string]any{
		"AppName":       cfg.AppName,
		"AppNameLower":  strings.ToLower(cfg.AppName),
		"JavaVersion":   cfg.JavaVersion,
		"Port":          cfg.Port,
		"BuildTool":     string(cfg.BuildTool),
		"HasWrapper":    cfg.HasWrapper,
		"IsGradle":      cfg.BuildTool.IsGradle(),
		"IsMaven":       cfg.BuildTool == buildtool.Maven,
		"HasDatabase":   cfg.DatabaseInfo != nil,
		"DatabaseType":  "",
		"DatabaseImage": "",
		"DatabasePort":  0,
		"DatabaseEnv":   map[string]string{},
		"VolumeName":    "",
		"VolumePath":    "",
		"SpringURL":     "",
	}

	if cfg.DatabaseInfo != nil {
		appNameLower := strings.ToLower(cfg.AppName)
		appNameLower = regexp.MustCompile(`[^a-z0-9]`).ReplaceAllString(appNameLower, "")
		if appNameLower == "" {
			appNameLower = "app"
		}

		envVars := make(map[string]string)
		for k, v := range cfg.DatabaseInfo.EnvVars {
			envVars[k] = strings.ReplaceAll(v, "${APP_NAME}", appNameLower)
		}

		data["DatabaseType"] = cfg.DatabaseInfo.Type
		data["DatabaseImage"] = cfg.DatabaseInfo.Image
		data["DatabasePort"] = cfg.DatabaseInfo.Port
		data["DatabaseEnv"] = envVars
		data["VolumeName"] = cfg.DatabaseInfo.Type + "_data"
		data["VolumePath"] = cfg.DatabaseInfo.VolumePath
		data["SpringURL"] = strings.ReplaceAll(cfg.DatabaseInfo.SpringURL, "${APP_NAME}", appNameLower)
	}

	return data
}

func detectAppName(project *buildtool.Project, cwd string) string {
	if project.ArtifactId != "" {
		return project.ArtifactId
	}
	if project.Name != "" {
		return strings.ToLower(strings.ReplaceAll(project.Name, " ", "-"))
	}
	return filepath.Base(cwd)
}

func detectJavaVersion(project *buildtool.Project, parser buildtool.Parser, override string) string {
	if override != "" {
		return normalizeJavaVersion(override)
	}

	if v := parser.GetJavaVersion(project); v != "" {
		return normalizeJavaVersion(v)
	}

	return "21"
}

func normalizeJavaVersion(v string) string {
	v = strings.TrimPrefix(v, "1.")
	v = strings.TrimSuffix(v, ".0")

	switch v {
	case "8", "11", "17", "21", "22", "23":
		return v
	default:
		if strings.HasPrefix(v, "1") && len(v) == 2 {
			return v
		}
		return "21"
	}
}

func detectPort(cwd string, override int) int {
	if override > 0 {
		return override
	}

	propsFiles := []string{
		filepath.Join(cwd, "src", "main", "resources", "application.yml"),
		filepath.Join(cwd, "src", "main", "resources", "application.yaml"),
		filepath.Join(cwd, "src", "main", "resources", "application.properties"),
	}

	for _, path := range propsFiles {
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)

			if strings.HasPrefix(line, "server.port=") {
				portStr := strings.TrimPrefix(line, "server.port=")
				portStr = strings.TrimSpace(portStr)
				var port int
				if _, err := fmt.Sscanf(portStr, "%d", &port); err == nil && port > 0 {
					return port
				}
			}

			if strings.Contains(line, "port:") && !strings.HasPrefix(line, "#") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					portStr := strings.TrimSpace(parts[1])
					var port int
					if _, err := fmt.Sscanf(portStr, "%d", &port); err == nil && port > 0 {
						return port
					}
				}
			}
		}
	}

	return 8080
}

func hasWrapper(cwd string, buildTool buildtool.Type) bool {
	var wrapperFile string
	if buildTool.IsGradle() {
		wrapperFile = filepath.Join(cwd, "gradlew")
	} else {
		wrapperFile = filepath.Join(cwd, "mvnw")
	}

	_, err := os.Stat(wrapperFile)
	return err == nil
}

func detectDatabaseFromDependencies(deps []buildtool.Dependency) string {
	for _, dep := range deps {
		key := dep.ArtifactId
		if _, ok := databaseDrivers[key]; ok {
			return key
		}

		combined := dep.GroupId + ":" + dep.ArtifactId
		if strings.Contains(combined, "postgresql") {
			return "postgresql"
		}
		if strings.Contains(combined, "mysql") {
			return "mysql"
		}
		if strings.Contains(combined, "mariadb") {
			return "mariadb"
		}
		if strings.Contains(combined, "mongodb") || strings.Contains(combined, "data-mongodb") {
			return "data-mongodb"
		}
		if strings.Contains(combined, "redis") || strings.Contains(combined, "data-redis") {
			return "data-redis"
		}
		if strings.Contains(combined, "cassandra") {
			return "data-cassandra"
		}
	}
	return ""
}

func hasJPADependency(deps []buildtool.Dependency) bool {
	for _, dep := range deps {
		if strings.Contains(dep.ArtifactId, "data-jpa") ||
			strings.Contains(dep.ArtifactId, "spring-data-jpa") ||
			strings.Contains(dep.ArtifactId, "hibernate") {
			return true
		}
	}
	return false
}

type selectWrapper struct {
	model components.SelectModel
}

func (w selectWrapper) Init() tea.Cmd {
	return w.model.Init()
}

func (w selectWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return w, tea.Quit
		}
	}

	updated, cmd := w.model.Update(msg)
	w.model = updated

	if w.model.Submitted() || w.model.GoBack() {
		return w, tea.Quit
	}

	return w, cmd
}

func (w selectWrapper) View() string {
	return w.model.View()
}

func runDatabasePicker() (string, error) {
	items := make([]components.SelectItem, len(databaseChoices))
	for i, choice := range databaseChoices {
		items[i] = components.SelectItem{
			Label:       choice.Label,
			Value:       choice.Value,
			Description: choice.Description,
		}
	}

	model := components.NewSelect(components.SelectConfig{
		Label: "JPA detected but no database driver found. Select a database for docker-compose:",
		Items: items,
	})

	wrapper := selectWrapper{model: model}
	p := tea.NewProgram(wrapper)
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	result := finalModel.(selectWrapper)
	if result.model.GoBack() {
		return "", fmt.Errorf("wizard cancelled")
	}

	return result.model.Value(), nil
}
