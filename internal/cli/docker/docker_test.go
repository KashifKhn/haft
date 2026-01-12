package docker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KashifKhn/haft/internal/buildtool"
	_ "github.com/KashifKhn/haft/internal/gradle"
	_ "github.com/KashifKhn/haft/internal/maven"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCommand(t *testing.T) {
	cmd := NewCommand()

	assert.Equal(t, "dockerize", cmd.Use)
	assert.Contains(t, cmd.Aliases, "docker")
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
}

func TestNewCommandFlags(t *testing.T) {
	cmd := NewCommand()

	portFlag := cmd.Flags().Lookup("port")
	assert.NotNil(t, portFlag)
	assert.Equal(t, "p", portFlag.Shorthand)

	javaFlag := cmd.Flags().Lookup("java")
	assert.NotNil(t, javaFlag)
	assert.Equal(t, "j", javaFlag.Shorthand)

	dbFlag := cmd.Flags().Lookup("db")
	assert.NotNil(t, dbFlag)

	noComposeFlag := cmd.Flags().Lookup("no-compose")
	assert.NotNil(t, noComposeFlag)

	noInteractiveFlag := cmd.Flags().Lookup("no-interactive")
	assert.NotNil(t, noInteractiveFlag)

	jsonFlag := cmd.Flags().Lookup("json")
	assert.NotNil(t, jsonFlag)
}

func TestDatabaseDriversDefinition(t *testing.T) {
	expectedDrivers := []string{
		"postgresql", "mysql-connector-j", "mysql", "mariadb",
		"data-mongodb", "mongodb", "data-redis", "data-cassandra",
	}

	for _, driver := range expectedDrivers {
		info, exists := databaseDrivers[driver]
		assert.True(t, exists, "driver %s should exist", driver)
		assert.NotEmpty(t, info.Type)
		assert.NotEmpty(t, info.Image)
		assert.Greater(t, info.Port, 0)
	}
}

func TestDatabaseChoicesDefinition(t *testing.T) {
	assert.Len(t, databaseChoices, 6)

	expectedChoices := []string{
		"PostgreSQL", "MySQL", "MariaDB", "MongoDB", "Redis", "None",
	}

	for i, expected := range expectedChoices {
		assert.Equal(t, expected, databaseChoices[i].Label)
		assert.NotEmpty(t, databaseChoices[i].Value)
		assert.NotEmpty(t, databaseChoices[i].Description)
	}
}

func TestNormalizeJavaVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"standard 8", "8", "8"},
		{"standard 11", "11", "11"},
		{"standard 17", "17", "17"},
		{"standard 21", "21", "21"},
		{"standard 22", "22", "22"},
		{"with prefix 1.8", "1.8", "8"},
		{"with suffix 17.0", "17.0", "17"},
		{"with prefix 1.11", "1.11", "11"},
		{"invalid version", "99", "21"},
		{"empty string", "", "21"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeJavaVersion(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDetectAppName(t *testing.T) {
	tests := []struct {
		name     string
		project  *buildtool.Project
		cwd      string
		expected string
	}{
		{
			name:     "from artifactId",
			project:  &buildtool.Project{ArtifactId: "my-app"},
			cwd:      "/some/path/project",
			expected: "my-app",
		},
		{
			name:     "from name",
			project:  &buildtool.Project{Name: "My Application"},
			cwd:      "/some/path/project",
			expected: "my-application",
		},
		{
			name:     "from directory",
			project:  &buildtool.Project{},
			cwd:      "/some/path/demo-project",
			expected: "demo-project",
		},
		{
			name:     "artifactId takes priority",
			project:  &buildtool.Project{ArtifactId: "artifact", Name: "Name"},
			cwd:      "/some/path/directory",
			expected: "artifact",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectAppName(tt.project, tt.cwd)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDetectPort(t *testing.T) {
	t.Run("override takes priority", func(t *testing.T) {
		result := detectPort("/nonexistent", 9000)
		assert.Equal(t, 9000, result)
	})

	t.Run("default when no config", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "test-port-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		result := detectPort(tmpDir, 0)
		assert.Equal(t, 8080, result)
	})

	t.Run("from properties file", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "test-port-props-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		resourcesPath := filepath.Join(tmpDir, "src", "main", "resources")
		require.NoError(t, os.MkdirAll(resourcesPath, 0755))

		propsContent := "server.port=9090\n"
		require.NoError(t, os.WriteFile(
			filepath.Join(resourcesPath, "application.properties"),
			[]byte(propsContent),
			0644,
		))

		result := detectPort(tmpDir, 0)
		assert.Equal(t, 9090, result)
	})

	t.Run("from yml file", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "test-port-yml-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		resourcesPath := filepath.Join(tmpDir, "src", "main", "resources")
		require.NoError(t, os.MkdirAll(resourcesPath, 0755))

		ymlContent := "server:\n  port: 8888\n"
		require.NoError(t, os.WriteFile(
			filepath.Join(resourcesPath, "application.yml"),
			[]byte(ymlContent),
			0644,
		))

		result := detectPort(tmpDir, 0)
		assert.Equal(t, 8888, result)
	})
}

func TestHasWrapper(t *testing.T) {
	t.Run("maven wrapper exists", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "test-wrapper-mvn-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		require.NoError(t, os.WriteFile(
			filepath.Join(tmpDir, "mvnw"),
			[]byte("#!/bin/bash\n"),
			0755,
		))

		result := hasWrapper(tmpDir, buildtool.Maven)
		assert.True(t, result)
	})

	t.Run("gradle wrapper exists", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "test-wrapper-gradle-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		require.NoError(t, os.WriteFile(
			filepath.Join(tmpDir, "gradlew"),
			[]byte("#!/bin/bash\n"),
			0755,
		))

		result := hasWrapper(tmpDir, buildtool.Gradle)
		assert.True(t, result)
	})

	t.Run("no wrapper", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "test-wrapper-none-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		result := hasWrapper(tmpDir, buildtool.Maven)
		assert.False(t, result)
	})
}

func TestDetectDatabaseFromDependencies(t *testing.T) {
	tests := []struct {
		name     string
		deps     []buildtool.Dependency
		expected string
	}{
		{
			name:     "empty dependencies",
			deps:     []buildtool.Dependency{},
			expected: "",
		},
		{
			name: "postgresql driver",
			deps: []buildtool.Dependency{
				{GroupId: "org.postgresql", ArtifactId: "postgresql"},
			},
			expected: "postgresql",
		},
		{
			name: "mysql connector j",
			deps: []buildtool.Dependency{
				{GroupId: "com.mysql", ArtifactId: "mysql-connector-j"},
			},
			expected: "mysql-connector-j",
		},
		{
			name: "mysql in groupId",
			deps: []buildtool.Dependency{
				{GroupId: "mysql", ArtifactId: "mysql-connector-java"},
			},
			expected: "mysql",
		},
		{
			name: "mariadb",
			deps: []buildtool.Dependency{
				{GroupId: "org.mariadb.jdbc", ArtifactId: "mariadb-java-client"},
			},
			expected: "mariadb",
		},
		{
			name: "mongodb spring data",
			deps: []buildtool.Dependency{
				{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-mongodb"},
			},
			expected: "data-mongodb",
		},
		{
			name: "redis spring data",
			deps: []buildtool.Dependency{
				{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-redis"},
			},
			expected: "data-redis",
		},
		{
			name: "cassandra",
			deps: []buildtool.Dependency{
				{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-cassandra"},
			},
			expected: "data-cassandra",
		},
		{
			name: "no database dependency",
			deps: []buildtool.Dependency{
				{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-web"},
				{GroupId: "org.projectlombok", ArtifactId: "lombok"},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectDatabaseFromDependencies(tt.deps)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHasJPADependency(t *testing.T) {
	tests := []struct {
		name     string
		deps     []buildtool.Dependency
		expected bool
	}{
		{
			name:     "empty dependencies",
			deps:     []buildtool.Dependency{},
			expected: false,
		},
		{
			name: "has data-jpa",
			deps: []buildtool.Dependency{
				{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-jpa"},
			},
			expected: true,
		},
		{
			name: "has hibernate",
			deps: []buildtool.Dependency{
				{GroupId: "org.hibernate", ArtifactId: "hibernate-core"},
			},
			expected: true,
		},
		{
			name: "no jpa dependency",
			deps: []buildtool.Dependency{
				{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-web"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasJPADependency(tt.deps)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildDockerTemplateData(t *testing.T) {
	t.Run("without database", func(t *testing.T) {
		cfg := DockerConfig{
			AppName:         "my-app",
			JavaVersion:     "21",
			Port:            8080,
			BuildTool:       buildtool.Maven,
			HasWrapper:      true,
			GenerateCompose: true,
			DatabaseInfo:    nil,
		}

		data := buildDockerTemplateData(cfg)

		assert.Equal(t, "my-app", data["AppName"])
		assert.Equal(t, "my-app", data["AppNameLower"])
		assert.Equal(t, "21", data["JavaVersion"])
		assert.Equal(t, 8080, data["Port"])
		assert.Equal(t, "maven", data["BuildTool"])
		assert.True(t, data["HasWrapper"].(bool))
		assert.False(t, data["IsGradle"].(bool))
		assert.True(t, data["IsMaven"].(bool))
		assert.False(t, data["HasDatabase"].(bool))
	})

	t.Run("with postgresql database", func(t *testing.T) {
		dbInfo := databaseDrivers["postgresql"]
		cfg := DockerConfig{
			AppName:         "My App",
			JavaVersion:     "17",
			Port:            9000,
			BuildTool:       buildtool.Gradle,
			HasWrapper:      false,
			GenerateCompose: true,
			DatabaseType:    "postgresql",
			DatabaseInfo:    &dbInfo,
		}

		data := buildDockerTemplateData(cfg)

		assert.Equal(t, "My App", data["AppName"])
		assert.Equal(t, "my app", data["AppNameLower"])
		assert.Equal(t, "17", data["JavaVersion"])
		assert.Equal(t, 9000, data["Port"])
		assert.Equal(t, "gradle", data["BuildTool"])
		assert.False(t, data["HasWrapper"].(bool))
		assert.True(t, data["IsGradle"].(bool))
		assert.False(t, data["IsMaven"].(bool))
		assert.True(t, data["HasDatabase"].(bool))
		assert.Equal(t, "postgres", data["DatabaseType"])
		assert.Equal(t, "postgres:16-alpine", data["DatabaseImage"])
		assert.Equal(t, 5432, data["DatabasePort"])
		assert.Equal(t, "postgres_data", data["VolumeName"])
	})
}

func TestDockerConfigStruct(t *testing.T) {
	cfg := DockerConfig{
		AppName:         "test-app",
		JavaVersion:     "21",
		Port:            8080,
		BuildTool:       buildtool.Maven,
		HasWrapper:      true,
		DatabaseType:    "postgresql",
		GenerateCompose: true,
	}

	assert.Equal(t, "test-app", cfg.AppName)
	assert.Equal(t, "21", cfg.JavaVersion)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, buildtool.Maven, cfg.BuildTool)
	assert.True(t, cfg.HasWrapper)
	assert.Equal(t, "postgresql", cfg.DatabaseType)
	assert.True(t, cfg.GenerateCompose)
}

func TestDockerOutputStruct(t *testing.T) {
	out := DockerOutput{
		Generated: []string{"Dockerfile", "docker-compose.yml"},
		Skipped:   []string{".dockerignore"},
	}
	out.Config.AppName = "my-app"
	out.Config.JavaVersion = "21"
	out.Config.Port = 8080
	out.Config.BuildTool = "maven"
	out.Config.DatabaseType = "postgres"

	assert.Len(t, out.Generated, 2)
	assert.Contains(t, out.Generated, "Dockerfile")
	assert.Len(t, out.Skipped, 1)
	assert.Equal(t, "my-app", out.Config.AppName)
	assert.Equal(t, "21", out.Config.JavaVersion)
	assert.Equal(t, 8080, out.Config.Port)
	assert.Equal(t, "maven", out.Config.BuildTool)
	assert.Equal(t, "postgres", out.Config.DatabaseType)
}

func TestDatabaseInfoStruct(t *testing.T) {
	info := DatabaseInfo{
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
	}

	assert.Equal(t, "postgres", info.Type)
	assert.Equal(t, "postgres:16-alpine", info.Image)
	assert.Equal(t, 5432, info.Port)
	assert.Equal(t, "POSTGRES", info.EnvPrefix)
	assert.Contains(t, info.SpringURL, "postgresql")
	assert.Len(t, info.EnvVars, 3)
}

func TestRunDockerizeNoBuildFile(t *testing.T) {
	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	tmpDir, err := os.MkdirTemp("", "test-docker-nobuild-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	require.NoError(t, os.Chdir(tmpDir))

	cmd := NewCommand()
	cmd.SetArgs([]string{"--no-interactive"})

	err = cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no build file found")
}

func TestRunDockerizeInvalidDB(t *testing.T) {
	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	tmpDir, err := os.MkdirTemp("", "test-docker-invaliddb-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	pomContent := `<?xml version="1.0" encoding="UTF-8"?>
<project>
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>demo</artifactId>
    <version>1.0.0</version>
</project>`
	require.NoError(t, os.WriteFile(
		filepath.Join(tmpDir, "pom.xml"),
		[]byte(pomContent),
		0644,
	))

	require.NoError(t, os.Chdir(tmpDir))

	cmd := NewCommand()
	cmd.SetArgs([]string{"--db", "invaliddb", "--no-interactive"})

	err = cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid database type")
}

func TestGenerateDockerFilesMaven(t *testing.T) {
	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	tmpDir, err := os.MkdirTemp("", "test-docker-maven-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	pomContent := `<?xml version="1.0" encoding="UTF-8"?>
<project>
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>demo</artifactId>
    <version>1.0.0</version>
    <properties>
        <java.version>21</java.version>
    </properties>
</project>`
	require.NoError(t, os.WriteFile(
		filepath.Join(tmpDir, "pom.xml"),
		[]byte(pomContent),
		0644,
	))

	require.NoError(t, os.WriteFile(
		filepath.Join(tmpDir, "mvnw"),
		[]byte("#!/bin/bash\n"),
		0755,
	))

	require.NoError(t, os.Chdir(tmpDir))

	cfg := DockerConfig{
		AppName:         "demo",
		JavaVersion:     "21",
		Port:            8080,
		BuildTool:       buildtool.Maven,
		HasWrapper:      true,
		GenerateCompose: true,
	}

	err = generateDockerFiles(cfg, tmpDir, false)
	require.NoError(t, err)

	assert.FileExists(t, filepath.Join(tmpDir, "Dockerfile"))
	assert.FileExists(t, filepath.Join(tmpDir, ".dockerignore"))
	assert.FileExists(t, filepath.Join(tmpDir, "docker-compose.yml"))

	dockerfileContent, err := os.ReadFile(filepath.Join(tmpDir, "Dockerfile"))
	require.NoError(t, err)
	assert.Contains(t, string(dockerfileContent), "eclipse-temurin:21")
	assert.Contains(t, string(dockerfileContent), "mvnw")
}

func TestGenerateDockerFilesGradle(t *testing.T) {
	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	tmpDir, err := os.MkdirTemp("", "test-docker-gradle-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	buildGradle := `plugins {
    id 'java'
    id 'org.springframework.boot' version '3.2.0'
}

java {
    sourceCompatibility = '17'
}`
	require.NoError(t, os.WriteFile(
		filepath.Join(tmpDir, "build.gradle"),
		[]byte(buildGradle),
		0644,
	))

	require.NoError(t, os.WriteFile(
		filepath.Join(tmpDir, "gradlew"),
		[]byte("#!/bin/bash\n"),
		0755,
	))

	require.NoError(t, os.Chdir(tmpDir))

	cfg := DockerConfig{
		AppName:         "gradle-app",
		JavaVersion:     "17",
		Port:            8080,
		BuildTool:       buildtool.Gradle,
		HasWrapper:      true,
		GenerateCompose: true,
	}

	err = generateDockerFiles(cfg, tmpDir, false)
	require.NoError(t, err)

	assert.FileExists(t, filepath.Join(tmpDir, "Dockerfile"))

	dockerfileContent, err := os.ReadFile(filepath.Join(tmpDir, "Dockerfile"))
	require.NoError(t, err)
	assert.Contains(t, string(dockerfileContent), "eclipse-temurin:17")
	assert.Contains(t, string(dockerfileContent), "gradlew")
	assert.Contains(t, string(dockerfileContent), "bootJar")
}

func TestGenerateDockerFilesWithDatabase(t *testing.T) {
	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	tmpDir, err := os.MkdirTemp("", "test-docker-db-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	pomContent := `<?xml version="1.0" encoding="UTF-8"?>
<project>
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>demo</artifactId>
    <version>1.0.0</version>
</project>`
	require.NoError(t, os.WriteFile(
		filepath.Join(tmpDir, "pom.xml"),
		[]byte(pomContent),
		0644,
	))

	require.NoError(t, os.Chdir(tmpDir))

	dbInfo := databaseDrivers["postgresql"]
	cfg := DockerConfig{
		AppName:         "demo",
		JavaVersion:     "21",
		Port:            8080,
		BuildTool:       buildtool.Maven,
		HasWrapper:      false,
		GenerateCompose: true,
		DatabaseType:    "postgresql",
		DatabaseInfo:    &dbInfo,
	}

	err = generateDockerFiles(cfg, tmpDir, false)
	require.NoError(t, err)

	composeContent, err := os.ReadFile(filepath.Join(tmpDir, "docker-compose.yml"))
	require.NoError(t, err)
	assert.Contains(t, string(composeContent), "postgres:16-alpine")
	assert.Contains(t, string(composeContent), "POSTGRES_DB")
	assert.Contains(t, string(composeContent), "depends_on")
}

func TestGenerateDockerFilesNoCompose(t *testing.T) {
	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	tmpDir, err := os.MkdirTemp("", "test-docker-nocompose-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	pomContent := `<?xml version="1.0" encoding="UTF-8"?>
<project>
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>demo</artifactId>
    <version>1.0.0</version>
</project>`
	require.NoError(t, os.WriteFile(
		filepath.Join(tmpDir, "pom.xml"),
		[]byte(pomContent),
		0644,
	))

	require.NoError(t, os.Chdir(tmpDir))

	cfg := DockerConfig{
		AppName:         "demo",
		JavaVersion:     "21",
		Port:            8080,
		BuildTool:       buildtool.Maven,
		HasWrapper:      false,
		GenerateCompose: false,
	}

	err = generateDockerFiles(cfg, tmpDir, false)
	require.NoError(t, err)

	assert.FileExists(t, filepath.Join(tmpDir, "Dockerfile"))
	assert.FileExists(t, filepath.Join(tmpDir, ".dockerignore"))
	assert.NoFileExists(t, filepath.Join(tmpDir, "docker-compose.yml"))
}

func TestGenerateDockerFilesSkipsExisting(t *testing.T) {
	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	tmpDir, err := os.MkdirTemp("", "test-docker-skip-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	pomContent := `<?xml version="1.0" encoding="UTF-8"?>
<project>
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>demo</artifactId>
    <version>1.0.0</version>
</project>`
	require.NoError(t, os.WriteFile(
		filepath.Join(tmpDir, "pom.xml"),
		[]byte(pomContent),
		0644,
	))

	existingContent := "existing dockerfile content"
	require.NoError(t, os.WriteFile(
		filepath.Join(tmpDir, "Dockerfile"),
		[]byte(existingContent),
		0644,
	))

	require.NoError(t, os.Chdir(tmpDir))

	cfg := DockerConfig{
		AppName:         "demo",
		JavaVersion:     "21",
		Port:            8080,
		BuildTool:       buildtool.Maven,
		HasWrapper:      false,
		GenerateCompose: true,
	}

	err = generateDockerFiles(cfg, tmpDir, false)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(tmpDir, "Dockerfile"))
	require.NoError(t, err)
	assert.Equal(t, existingContent, string(content))

	assert.FileExists(t, filepath.Join(tmpDir, ".dockerignore"))
	assert.FileExists(t, filepath.Join(tmpDir, "docker-compose.yml"))
}

func TestSelectWrapperInit(t *testing.T) {
	wrapper := selectWrapper{}
	cmd := wrapper.Init()
	assert.Nil(t, cmd)
}

func TestSelectWrapperView(t *testing.T) {
	wrapper := selectWrapper{}
	view := wrapper.View()
	assert.NotNil(t, view)
}

func TestDetectJavaVersionWithOverride(t *testing.T) {
	result := normalizeJavaVersion("17")
	assert.Equal(t, "17", result)
}

func TestDetectJavaVersionDefault(t *testing.T) {
	result := normalizeJavaVersion("")
	assert.Equal(t, "21", result)
}
