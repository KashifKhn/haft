package maven

import (
	"testing"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_ParseBytes(t *testing.T) {
	parser := NewParser()

	pomXml := `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0">
    <modelVersion>4.0.0</modelVersion>
    <parent>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-parent</artifactId>
        <version>3.4.1</version>
    </parent>
    <groupId>com.example</groupId>
    <artifactId>my-app</artifactId>
    <version>0.0.1-SNAPSHOT</version>
    <name>My App</name>
    <properties>
        <java.version>21</java.version>
    </properties>
    <dependencies>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-web</artifactId>
        </dependency>
        <dependency>
            <groupId>org.projectlombok</groupId>
            <artifactId>lombok</artifactId>
            <scope>provided</scope>
        </dependency>
    </dependencies>
</project>`

	project, err := parser.ParseBytes([]byte(pomXml))

	require.NoError(t, err)
	assert.Equal(t, "com.example", project.GroupId)
	assert.Equal(t, "my-app", project.ArtifactId)
	assert.Equal(t, "0.0.1-SNAPSHOT", project.Version)
	assert.Equal(t, "My App", project.Name)
	assert.Equal(t, "21", project.JavaVersion)
	assert.Len(t, project.Dependencies, 2)
}

func TestParser_HasDependency(t *testing.T) {
	parser := NewParser()
	project := &buildtool.Project{
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-web"},
			{GroupId: "org.projectlombok", ArtifactId: "lombok"},
		},
	}

	assert.True(t, parser.HasDependency(project, "org.projectlombok", "lombok"))
	assert.True(t, parser.HasDependency(project, "org.springframework.boot", "spring-boot-starter-web"))
	assert.False(t, parser.HasDependency(project, "org.mapstruct", "mapstruct"))
}

func TestParser_AddDependency(t *testing.T) {
	parser := NewParser()
	project := &buildtool.Project{}

	parser.AddDependency(project, buildtool.Dependency{
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-web",
	})

	assert.Len(t, project.Dependencies, 1)

	parser.AddDependency(project, buildtool.Dependency{
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-web",
	})
	assert.Len(t, project.Dependencies, 1)

	parser.AddDependency(project, buildtool.Dependency{
		GroupId:    "org.projectlombok",
		ArtifactId: "lombok",
		Scope:      "provided",
	})
	assert.Len(t, project.Dependencies, 2)
}

func TestParser_RemoveDependency(t *testing.T) {
	parser := NewParser()
	project := &buildtool.Project{
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-web"},
			{GroupId: "org.projectlombok", ArtifactId: "lombok"},
		},
	}

	removed := parser.RemoveDependency(project, "org.projectlombok", "lombok")

	assert.True(t, removed)
	assert.Len(t, project.Dependencies, 1)
	assert.False(t, parser.HasDependency(project, "org.projectlombok", "lombok"))

	removed = parser.RemoveDependency(project, "org.mapstruct", "mapstruct")
	assert.False(t, removed)
}

func TestParser_HelperMethods(t *testing.T) {
	parser := NewParser()
	project := &buildtool.Project{
		GroupId:           "com.example",
		ArtifactId:        "my-app",
		JavaVersion:       "21",
		SpringBootVersion: "3.4.1",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.projectlombok", ArtifactId: "lombok"},
			{GroupId: "org.mapstruct", ArtifactId: "mapstruct"},
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-jpa"},
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-web"},
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-validation"},
		},
	}

	assert.True(t, parser.HasLombok(project))
	assert.True(t, parser.HasMapStruct(project))
	assert.True(t, parser.HasSpringDataJpa(project))
	assert.True(t, parser.HasSpringWeb(project))
	assert.True(t, parser.HasValidation(project))
	assert.Equal(t, "21", parser.GetJavaVersion(project))
	assert.Equal(t, "3.4.1", parser.GetSpringBootVersion(project))
	assert.Equal(t, "com.example.myapp", parser.GetBasePackage(project))
}

func TestParser_Parse_FromFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	parser := NewParserWithFs(fs)

	pomContent := `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0">
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.test</groupId>
    <artifactId>test-app</artifactId>
    <version>1.0.0</version>
</project>`

	require.NoError(t, afero.WriteFile(fs, "/project/pom.xml", []byte(pomContent), 0644))

	project, err := parser.Parse("/project/pom.xml")

	require.NoError(t, err)
	assert.Equal(t, "com.test", project.GroupId)
	assert.Equal(t, "test-app", project.ArtifactId)
}

func TestParser_Write(t *testing.T) {
	fs := afero.NewMemMapFs()
	parser := NewParserWithFs(fs)

	project := &buildtool.Project{
		GroupId:    "com.example",
		ArtifactId: "my-app",
		Version:    "0.0.1-SNAPSHOT",
		BuildTool:  buildtool.Maven,
		Raw: &MavenProject{
			ModelVersion: "4.0.0",
			GroupId:      "com.example",
			ArtifactId:   "my-app",
			Version:      "0.0.1-SNAPSHOT",
		},
	}

	err := parser.Write("/project/pom.xml", project)

	require.NoError(t, err)

	data, err := afero.ReadFile(fs, "/project/pom.xml")
	require.NoError(t, err)
	assert.Contains(t, string(data), "com.example")
	assert.Contains(t, string(data), "my-app")
}

func TestParser_FindPomXml(t *testing.T) {
	fs := afero.NewMemMapFs()
	parser := NewParserWithFs(fs)

	require.NoError(t, fs.MkdirAll("/project/src/main/java", 0755))
	require.NoError(t, afero.WriteFile(fs, "/project/pom.xml", []byte("<project/>"), 0644))

	pomPath, err := parser.FindPomXml("/project/src/main/java")

	require.NoError(t, err)
	assert.Equal(t, "/project/pom.xml", pomPath)
}

func TestParser_FindPomXml_NotFound(t *testing.T) {
	fs := afero.NewMemMapFs()
	parser := NewParserWithFs(fs)

	require.NoError(t, fs.MkdirAll("/project/src", 0755))

	_, err := parser.FindPomXml("/project/src")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pom.xml not found")
}

func TestParser_Parse_FileNotFound(t *testing.T) {
	fs := afero.NewMemMapFs()
	parser := NewParserWithFs(fs)

	_, err := parser.Parse("/nonexistent/pom.xml")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read pom.xml")
}

func TestParser_ParseBytes_InvalidXml(t *testing.T) {
	parser := NewParser()

	_, err := parser.ParseBytes([]byte("not valid xml"))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse pom.xml")
}

func TestParser_GetDependency(t *testing.T) {
	parser := NewParser()
	project := &buildtool.Project{
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-web", Version: "3.4.0"},
			{GroupId: "org.projectlombok", ArtifactId: "lombok", Scope: "provided"},
		},
	}

	dep := parser.GetDependency(project, "org.projectlombok", "lombok")
	assert.NotNil(t, dep)
	assert.Equal(t, "provided", dep.Scope)

	dep = parser.GetDependency(project, "org.springframework.boot", "spring-boot-starter-web")
	assert.NotNil(t, dep)
	assert.Equal(t, "3.4.0", dep.Version)

	dep = parser.GetDependency(project, "nonexistent", "dep")
	assert.Nil(t, dep)
}

func TestParser_GetDependency_NilDependencies(t *testing.T) {
	parser := NewParser()
	project := &buildtool.Project{}

	dep := parser.GetDependency(project, "org.projectlombok", "lombok")

	assert.Nil(t, dep)
}

func TestParser_HasDependency_NilDependencies(t *testing.T) {
	parser := NewParser()
	project := &buildtool.Project{}

	assert.False(t, parser.HasDependency(project, "org.projectlombok", "lombok"))
}

func TestParser_RemoveDependency_NilDependencies(t *testing.T) {
	parser := NewParser()
	project := &buildtool.Project{}

	removed := parser.RemoveDependency(project, "org.projectlombok", "lombok")

	assert.False(t, removed)
}

func TestParser_GetJavaVersion_Empty(t *testing.T) {
	parser := NewParser()
	project := &buildtool.Project{}

	version := parser.GetJavaVersion(project)

	assert.Equal(t, "", version)
}

func TestParser_GetSpringBootVersion_Empty(t *testing.T) {
	parser := NewParser()
	project := &buildtool.Project{}

	version := parser.GetSpringBootVersion(project)

	assert.Equal(t, "", version)
}

func TestParser_Marshal(t *testing.T) {
	parser := NewParser()
	mavenProject := &MavenProject{
		ModelVersion: "4.0.0",
		GroupId:      "com.example",
		ArtifactId:   "test-app",
		Version:      "1.0.0",
	}

	data, err := parser.Marshal(mavenProject)

	require.NoError(t, err)
	assert.Contains(t, string(data), "<?xml version")
	assert.Contains(t, string(data), "http://maven.apache.org/POM/4.0.0")
	assert.Contains(t, string(data), "com.example")
}

func TestParser_Type(t *testing.T) {
	parser := NewParser()
	assert.Equal(t, buildtool.Maven, parser.Type())
}

func TestParser_LegacyMethods(t *testing.T) {
	parser := NewParser()
	mavenProject := &MavenProject{
		GroupId:    "com.example",
		ArtifactId: "my-app",
		Parent: &Parent{
			GroupId:    "org.springframework.boot",
			ArtifactId: "spring-boot-starter-parent",
			Version:    "3.4.1",
		},
		Properties: &Properties{
			JavaVersion: "21",
		},
		Dependencies: &Dependencies{
			Dependency: []Dependency{
				{GroupId: "org.projectlombok", ArtifactId: "lombok"},
			},
		},
	}

	assert.True(t, parser.HasDependencyLegacy(mavenProject, "org.projectlombok", "lombok"))
	assert.True(t, parser.HasLombokLegacy(mavenProject))
	assert.Equal(t, "21", parser.GetJavaVersionLegacy(mavenProject))
	assert.Equal(t, "3.4.1", parser.GetSpringBootVersionLegacy(mavenProject))
	assert.Equal(t, "com.example.myapp", parser.GetBasePackageLegacy(mavenProject))

	parser.AddDependencyLegacy(mavenProject, Dependency{GroupId: "org.mapstruct", ArtifactId: "mapstruct"})
	assert.True(t, parser.HasMapStructLegacy(mavenProject))

	parser.RemoveDependencyLegacy(mavenProject, "org.mapstruct", "mapstruct")
	assert.False(t, parser.HasMapStructLegacy(mavenProject))
}
