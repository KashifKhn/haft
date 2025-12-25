package gradle

import (
	"testing"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_Parse_GroovyDSL(t *testing.T) {
	fs := afero.NewMemMapFs()
	parser := NewParserWithFs(fs, false)

	buildGradle := `plugins {
    id 'java'
    id 'org.springframework.boot' version '3.4.1'
    id 'io.spring.dependency-management' version '1.1.7'
}

group = 'com.example'
version = '0.0.1-SNAPSHOT'

java {
    sourceCompatibility = '21'
}

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-web'
    compileOnly 'org.projectlombok:lombok'
    annotationProcessor 'org.projectlombok:lombok'
    testImplementation 'org.springframework.boot:spring-boot-starter-test'
}
`
	settingsGradle := `rootProject.name = 'my-app'`

	require.NoError(t, fs.MkdirAll("/project", 0755))
	require.NoError(t, afero.WriteFile(fs, "/project/build.gradle", []byte(buildGradle), 0644))
	require.NoError(t, afero.WriteFile(fs, "/project/settings.gradle", []byte(settingsGradle), 0644))

	project, err := parser.Parse("/project/build.gradle")

	require.NoError(t, err)
	assert.Equal(t, "com.example", project.GroupId)
	assert.Equal(t, "0.0.1-SNAPSHOT", project.Version)
	assert.Equal(t, "my-app", project.ArtifactId)
	assert.Equal(t, "3.4.1", project.SpringBootVersion)
	assert.Len(t, project.Dependencies, 3)
}

func TestParser_Parse_KotlinDSL(t *testing.T) {
	fs := afero.NewMemMapFs()
	parser := NewParserWithFs(fs, true)

	buildGradleKts := `plugins {
    java
    id("org.springframework.boot") version "3.4.1"
    id("io.spring.dependency-management") version "1.1.7"
}

group = "com.example"
version = "0.0.1-SNAPSHOT"

java {
    toolchain {
        languageVersion.set(JavaLanguageVersion.of(21))
    }
}

dependencies {
    implementation("org.springframework.boot:spring-boot-starter-web")
    compileOnly("org.projectlombok:lombok")
    annotationProcessor("org.projectlombok:lombok")
    testImplementation("org.springframework.boot:spring-boot-starter-test")
}
`
	settingsGradleKts := `rootProject.name = "kotlin-app"`

	require.NoError(t, fs.MkdirAll("/project", 0755))
	require.NoError(t, afero.WriteFile(fs, "/project/build.gradle.kts", []byte(buildGradleKts), 0644))
	require.NoError(t, afero.WriteFile(fs, "/project/settings.gradle.kts", []byte(settingsGradleKts), 0644))

	project, err := parser.Parse("/project/build.gradle.kts")

	require.NoError(t, err)
	assert.Equal(t, "com.example", project.GroupId)
	assert.Equal(t, "0.0.1-SNAPSHOT", project.Version)
	assert.Equal(t, "kotlin-app", project.ArtifactId)
	assert.Equal(t, "21", project.JavaVersion)
	assert.Equal(t, "3.4.1", project.SpringBootVersion)
	assert.Len(t, project.Dependencies, 3)
}

func TestParser_HasDependency(t *testing.T) {
	parser := NewParser(false)
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

func TestParser_AddDependency_Groovy(t *testing.T) {
	fs := afero.NewMemMapFs()
	parser := NewParserWithFs(fs, false)

	buildGradle := `plugins {
    id 'java'
}

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-web'
}
`
	require.NoError(t, fs.MkdirAll("/project", 0755))
	require.NoError(t, afero.WriteFile(fs, "/project/build.gradle", []byte(buildGradle), 0644))

	project, err := parser.Parse("/project/build.gradle")
	require.NoError(t, err)

	parser.AddDependency(project, buildtool.Dependency{
		GroupId:    "org.projectlombok",
		ArtifactId: "lombok",
		Scope:      "provided",
	})

	assert.Len(t, project.Dependencies, 2)
	assert.True(t, parser.HasDependency(project, "org.projectlombok", "lombok"))

	gradleProject := project.Raw.(*GradleProject)
	assert.Contains(t, gradleProject.Content, "compileOnly 'org.projectlombok:lombok'")
}

func TestParser_AddDependency_Kotlin(t *testing.T) {
	fs := afero.NewMemMapFs()
	parser := NewParserWithFs(fs, true)

	buildGradleKts := `plugins {
    java
}

dependencies {
    implementation("org.springframework.boot:spring-boot-starter-web")
}
`
	require.NoError(t, fs.MkdirAll("/project", 0755))
	require.NoError(t, afero.WriteFile(fs, "/project/build.gradle.kts", []byte(buildGradleKts), 0644))

	project, err := parser.Parse("/project/build.gradle.kts")
	require.NoError(t, err)

	parser.AddDependency(project, buildtool.Dependency{
		GroupId:    "org.projectlombok",
		ArtifactId: "lombok",
		Scope:      "provided",
	})

	assert.Len(t, project.Dependencies, 2)

	gradleProject := project.Raw.(*GradleProject)
	assert.Contains(t, gradleProject.Content, `compileOnly("org.projectlombok:lombok")`)
}

func TestParser_AddDependency_NoDuplicates(t *testing.T) {
	parser := NewParser(false)
	project := &buildtool.Project{
		Raw: &GradleProject{Content: "dependencies {}"},
	}

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
}

func TestParser_RemoveDependency_Groovy(t *testing.T) {
	fs := afero.NewMemMapFs()
	parser := NewParserWithFs(fs, false)

	buildGradle := `dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-web'
    compileOnly 'org.projectlombok:lombok'
}
`
	require.NoError(t, fs.MkdirAll("/project", 0755))
	require.NoError(t, afero.WriteFile(fs, "/project/build.gradle", []byte(buildGradle), 0644))

	project, err := parser.Parse("/project/build.gradle")
	require.NoError(t, err)
	assert.Len(t, project.Dependencies, 2)

	removed := parser.RemoveDependency(project, "org.projectlombok", "lombok")

	assert.True(t, removed)
	assert.Len(t, project.Dependencies, 1)
	assert.False(t, parser.HasDependency(project, "org.projectlombok", "lombok"))

	gradleProject := project.Raw.(*GradleProject)
	assert.NotContains(t, gradleProject.Content, "lombok")
}

func TestParser_RemoveDependency_Kotlin(t *testing.T) {
	fs := afero.NewMemMapFs()
	parser := NewParserWithFs(fs, true)

	buildGradleKts := `dependencies {
    implementation("org.springframework.boot:spring-boot-starter-web")
    compileOnly("org.projectlombok:lombok")
}
`
	require.NoError(t, fs.MkdirAll("/project", 0755))
	require.NoError(t, afero.WriteFile(fs, "/project/build.gradle.kts", []byte(buildGradleKts), 0644))

	project, err := parser.Parse("/project/build.gradle.kts")
	require.NoError(t, err)

	removed := parser.RemoveDependency(project, "org.projectlombok", "lombok")

	assert.True(t, removed)
	assert.False(t, parser.HasDependency(project, "org.projectlombok", "lombok"))
}

func TestParser_RemoveDependency_NotFound(t *testing.T) {
	parser := NewParser(false)
	project := &buildtool.Project{
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-web"},
		},
		Raw: &GradleProject{Content: "dependencies {}"},
	}

	removed := parser.RemoveDependency(project, "nonexistent", "dep")

	assert.False(t, removed)
	assert.Len(t, project.Dependencies, 1)
}

func TestParser_GetDependency(t *testing.T) {
	parser := NewParser(false)
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

func TestParser_GetDependencies(t *testing.T) {
	parser := NewParser(false)
	project := &buildtool.Project{
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-web"},
			{GroupId: "org.projectlombok", ArtifactId: "lombok"},
		},
	}

	deps := parser.GetDependencies(project)

	assert.Len(t, deps, 2)
}

func TestParser_HelperMethods(t *testing.T) {
	parser := NewParser(false)
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

func TestParser_Type(t *testing.T) {
	groovyParser := NewParser(false)
	assert.Equal(t, buildtool.Gradle, groovyParser.Type())

	kotlinParser := NewParser(true)
	assert.Equal(t, buildtool.GradleKotln, kotlinParser.Type())
}

func TestParser_FindBuildFile_Groovy(t *testing.T) {
	fs := afero.NewMemMapFs()
	parser := NewParserWithFs(fs, false)

	require.NoError(t, fs.MkdirAll("/project/src/main/java", 0755))
	require.NoError(t, afero.WriteFile(fs, "/project/build.gradle", []byte("plugins {}"), 0644))

	buildPath, err := parser.FindBuildFile("/project/src/main/java")

	require.NoError(t, err)
	assert.Equal(t, "/project/build.gradle", buildPath)
}

func TestParser_FindBuildFile_Kotlin(t *testing.T) {
	fs := afero.NewMemMapFs()
	parser := NewParserWithFs(fs, true)

	require.NoError(t, fs.MkdirAll("/project/src/main/java", 0755))
	require.NoError(t, afero.WriteFile(fs, "/project/build.gradle.kts", []byte("plugins {}"), 0644))

	buildPath, err := parser.FindBuildFile("/project/src/main/java")

	require.NoError(t, err)
	assert.Equal(t, "/project/build.gradle.kts", buildPath)
}

func TestParser_FindBuildFile_NotFound(t *testing.T) {
	fs := afero.NewMemMapFs()
	parser := NewParserWithFs(fs, false)

	require.NoError(t, fs.MkdirAll("/project/src", 0755))

	_, err := parser.FindBuildFile("/project/src")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "build.gradle not found")
}

func TestParser_Write(t *testing.T) {
	fs := afero.NewMemMapFs()
	parser := NewParserWithFs(fs, false)

	require.NoError(t, fs.MkdirAll("/project", 0755))

	buildGradle := `plugins {
    id 'java'
}

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-web'
}
`
	require.NoError(t, afero.WriteFile(fs, "/project/build.gradle", []byte(buildGradle), 0644))

	project, err := parser.Parse("/project/build.gradle")
	require.NoError(t, err)

	parser.AddDependency(project, buildtool.Dependency{
		GroupId:    "org.projectlombok",
		ArtifactId: "lombok",
	})

	err = parser.Write("/project/build.gradle", project)
	require.NoError(t, err)

	data, err := afero.ReadFile(fs, "/project/build.gradle")
	require.NoError(t, err)
	assert.Contains(t, string(data), "lombok")
}

func TestParser_Parse_FileNotFound(t *testing.T) {
	fs := afero.NewMemMapFs()
	parser := NewParserWithFs(fs, false)

	_, err := parser.Parse("/nonexistent/build.gradle")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read gradle build file")
}

func TestParser_ExtractJavaVersion_Variations(t *testing.T) {
	testCases := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "sourceCompatibility string",
			content:  `sourceCompatibility = '21'`,
			expected: "21",
		},
		{
			name:     "sourceCompatibility number",
			content:  `sourceCompatibility = 17`,
			expected: "17",
		},
		{
			name:     "JavaVersion enum",
			content:  `sourceCompatibility = JavaVersion.VERSION_21`,
			expected: "21",
		},
		{
			name:     "java.sourceCompatibility",
			content:  `java.sourceCompatibility = JavaVersion.VERSION_17`,
			expected: "17",
		},
		{
			name:     "toolchain",
			content:  `toolchain { languageVersion.set(JavaLanguageVersion.of(21)) }`,
			expected: "21",
		},
		{
			name:     "jvmToolchain",
			content:  `jvmToolchain(17)`,
			expected: "17",
		},
	}

	parser := NewParser(false)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			version := parser.extractJavaVersion(tc.content, false)
			assert.Equal(t, tc.expected, version)
		})
	}
}

func TestParser_ExtractSpringBootVersion(t *testing.T) {
	testCases := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "Groovy DSL",
			content:  `id 'org.springframework.boot' version '3.4.1'`,
			expected: "3.4.1",
		},
		{
			name:     "Kotlin DSL",
			content:  `id("org.springframework.boot") version "3.4.1"`,
			expected: "3.4.1",
		},
		{
			name:     "Variable",
			content:  `springBootVersion = '3.4.0'`,
			expected: "3.4.0",
		},
	}

	parser := NewParser(false)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			version := parser.extractSpringBootVersion(tc.content, false)
			assert.Equal(t, tc.expected, version)
		})
	}
}

func TestParser_DependencyWithVersion(t *testing.T) {
	fs := afero.NewMemMapFs()
	parser := NewParserWithFs(fs, false)

	buildGradle := `dependencies {
    implementation 'com.google.guava:guava:32.1.2-jre'
    runtimeOnly 'org.postgresql:postgresql:42.6.0'
}
`
	require.NoError(t, fs.MkdirAll("/project", 0755))
	require.NoError(t, afero.WriteFile(fs, "/project/build.gradle", []byte(buildGradle), 0644))

	project, err := parser.Parse("/project/build.gradle")

	require.NoError(t, err)
	assert.Len(t, project.Dependencies, 2)

	guavaDep := parser.GetDependency(project, "com.google.guava", "guava")
	assert.NotNil(t, guavaDep)
	assert.Equal(t, "32.1.2-jre", guavaDep.Version)

	pgDep := parser.GetDependency(project, "org.postgresql", "postgresql")
	assert.NotNil(t, pgDep)
	assert.Equal(t, "runtime", pgDep.Scope)
}

func TestParser_ScopeMapping(t *testing.T) {
	parser := NewParser(false)

	testCases := []struct {
		gradleScope string
		mavenScope  string
	}{
		{"implementation", "compile"},
		{"compileOnly", "provided"},
		{"runtimeOnly", "runtime"},
		{"testImplementation", "test"},
		{"testCompileOnly", "test"},
		{"testRuntimeOnly", "test"},
		{"annotationProcessor", "provided"},
	}

	for _, tc := range testCases {
		t.Run(tc.gradleScope, func(t *testing.T) {
			scope := parser.mapGradleScope(tc.gradleScope)
			assert.Equal(t, tc.mavenScope, scope)
		})
	}
}

func TestParser_ReverseScopeMapping(t *testing.T) {
	parser := NewParser(false)

	testCases := []struct {
		mavenScope  string
		gradleScope string
	}{
		{"compile", "implementation"},
		{"", "implementation"},
		{"provided", "compileOnly"},
		{"runtime", "runtimeOnly"},
		{"test", "testImplementation"},
	}

	for _, tc := range testCases {
		t.Run(tc.mavenScope, func(t *testing.T) {
			scope := parser.mapScopeToGradle(tc.mavenScope, false)
			assert.Equal(t, tc.gradleScope, scope)
		})
	}
}

func TestFindGradleBuildFile(t *testing.T) {
	fs := afero.NewMemMapFs()

	require.NoError(t, fs.MkdirAll("/project/src/main/java", 0755))
	require.NoError(t, afero.WriteFile(fs, "/project/build.gradle", []byte("plugins {}"), 0644))

	path, isKotlin, err := FindGradleBuildFile("/project/src/main/java", fs)

	require.NoError(t, err)
	assert.Equal(t, "/project/build.gradle", path)
	assert.False(t, isKotlin)
}

func TestFindGradleBuildFile_KotlinDSL(t *testing.T) {
	fs := afero.NewMemMapFs()

	require.NoError(t, fs.MkdirAll("/project/src/main/java", 0755))
	require.NoError(t, afero.WriteFile(fs, "/project/build.gradle.kts", []byte("plugins {}"), 0644))

	path, isKotlin, err := FindGradleBuildFile("/project/src/main/java", fs)

	require.NoError(t, err)
	assert.Equal(t, "/project/build.gradle.kts", path)
	assert.True(t, isKotlin)
}

func TestFindGradleBuildFile_NotFound(t *testing.T) {
	fs := afero.NewMemMapFs()

	require.NoError(t, fs.MkdirAll("/project/src", 0755))

	_, _, err := FindGradleBuildFile("/project/src", fs)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no gradle build file found")
}

func TestParser_AddDependency_NoDependenciesBlock(t *testing.T) {
	fs := afero.NewMemMapFs()
	parser := NewParserWithFs(fs, false)

	buildGradle := `plugins {
    id 'java'
}

group = 'com.example'
`
	require.NoError(t, fs.MkdirAll("/project", 0755))
	require.NoError(t, afero.WriteFile(fs, "/project/build.gradle", []byte(buildGradle), 0644))

	project, err := parser.Parse("/project/build.gradle")
	require.NoError(t, err)

	parser.AddDependency(project, buildtool.Dependency{
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-web",
	})

	gradleProject := project.Raw.(*GradleProject)
	assert.Contains(t, gradleProject.Content, "dependencies {")
	assert.Contains(t, gradleProject.Content, "spring-boot-starter-web")
}

func TestParser_KaptScope(t *testing.T) {
	fs := afero.NewMemMapFs()
	parser := NewParserWithFs(fs, true)

	buildGradleKts := `dependencies {
    kapt("org.mapstruct:mapstruct-processor:1.5.5.Final")
}
`
	require.NoError(t, fs.MkdirAll("/project", 0755))
	require.NoError(t, afero.WriteFile(fs, "/project/build.gradle.kts", []byte(buildGradleKts), 0644))

	project, err := parser.Parse("/project/build.gradle.kts")

	require.NoError(t, err)
	assert.Len(t, project.Dependencies, 1)
	assert.Equal(t, "provided", project.Dependencies[0].Scope)
}

func TestParser_Write_InvalidProject(t *testing.T) {
	fs := afero.NewMemMapFs()
	parser := NewParserWithFs(fs, false)

	project := &buildtool.Project{
		Raw: nil,
	}

	err := parser.Write("/project/build.gradle", project)

	assert.Error(t, err)
}
