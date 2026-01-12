package doctor

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupFs(t *testing.T) afero.Fs {
	return afero.NewMemMapFs()
}

func writeFile(t *testing.T, fs afero.Fs, path string, content string) {
	err := afero.WriteFile(fs, path, []byte(content), 0644)
	require.NoError(t, err)
}

func mkdirAll(t *testing.T, fs afero.Fs, path string) {
	err := fs.MkdirAll(path, 0755)
	require.NoError(t, err)
}

func TestCheckBuildFile_Maven(t *testing.T) {
	fs := setupFs(t)
	writeFile(t, fs, "/project/pom.xml", "<project></project>")

	checker := NewChecker(fs, "/project")
	checker.detectBuildTool()
	result := checker.checkBuildFile()

	assert.True(t, result.Passed)
	assert.Contains(t, result.Message, "Maven")
}

func TestCheckBuildFile_Gradle(t *testing.T) {
	fs := setupFs(t)
	writeFile(t, fs, "/project/build.gradle", "plugins {}")

	checker := NewChecker(fs, "/project")
	checker.detectBuildTool()
	result := checker.checkBuildFile()

	assert.True(t, result.Passed)
	assert.Contains(t, result.Message, "Gradle")
}

func TestCheckBuildFile_Missing(t *testing.T) {
	fs := setupFs(t)
	mkdirAll(t, fs, "/project")

	checker := NewChecker(fs, "/project")
	checker.detectBuildTool()
	result := checker.checkBuildFile()

	assert.False(t, result.Passed)
	assert.Equal(t, SeverityError, result.Severity)
}

func TestCheckSpringBootParent_Present(t *testing.T) {
	fs := setupFs(t)
	pom := `<project>
		<parent>
			<groupId>org.springframework.boot</groupId>
			<artifactId>spring-boot-starter-parent</artifactId>
		</parent>
	</project>`
	writeFile(t, fs, "/project/pom.xml", pom)

	checker := NewChecker(fs, "/project")
	checker.detectBuildTool()
	result := checker.checkSpringBootParent()

	assert.True(t, result.Passed)
}

func TestCheckSpringBootParent_Missing(t *testing.T) {
	fs := setupFs(t)
	pom := `<project><groupId>com.example</groupId></project>`
	writeFile(t, fs, "/project/pom.xml", pom)

	checker := NewChecker(fs, "/project")
	checker.detectBuildTool()
	result := checker.checkSpringBootParent()

	assert.False(t, result.Passed)
	assert.Equal(t, SeverityError, result.Severity)
}

func TestCheckJavaVersion_17(t *testing.T) {
	fs := setupFs(t)
	pom := `<project>
		<properties>
			<java.version>17</java.version>
		</properties>
	</project>`
	writeFile(t, fs, "/project/pom.xml", pom)

	checker := NewChecker(fs, "/project")
	checker.detectBuildTool()
	result := checker.checkJavaVersion()

	assert.True(t, result.Passed)
	assert.Contains(t, result.Message, "17")
}

func TestCheckJavaVersion_8Warning(t *testing.T) {
	fs := setupFs(t)
	pom := `<project>
		<properties>
			<java.version>8</java.version>
		</properties>
	</project>`
	writeFile(t, fs, "/project/pom.xml", pom)

	checker := NewChecker(fs, "/project")
	checker.detectBuildTool()
	result := checker.checkJavaVersion()

	assert.True(t, result.Passed)
	assert.Equal(t, SeverityWarning, result.Severity)
	assert.Contains(t, result.Message, "upgrading")
}

func TestCheckJavaVersion_Missing(t *testing.T) {
	fs := setupFs(t)
	pom := `<project></project>`
	writeFile(t, fs, "/project/pom.xml", pom)

	checker := NewChecker(fs, "/project")
	checker.detectBuildTool()
	result := checker.checkJavaVersion()

	assert.False(t, result.Passed)
	assert.Contains(t, result.Message, "not specified")
}

func TestCheckSourceDirectory_Exists(t *testing.T) {
	fs := setupFs(t)
	mkdirAll(t, fs, "/project/src/main/java")

	checker := NewChecker(fs, "/project")
	result := checker.checkSourceDirectory()

	assert.True(t, result.Passed)
	assert.Contains(t, result.Message, "Java")
}

func TestCheckSourceDirectory_Kotlin(t *testing.T) {
	fs := setupFs(t)
	mkdirAll(t, fs, "/project/src/main/kotlin")

	checker := NewChecker(fs, "/project")
	result := checker.checkSourceDirectory()

	assert.True(t, result.Passed)
	assert.Contains(t, result.Message, "Kotlin")
}

func TestCheckSourceDirectory_Missing(t *testing.T) {
	fs := setupFs(t)
	mkdirAll(t, fs, "/project")

	checker := NewChecker(fs, "/project")
	result := checker.checkSourceDirectory()

	assert.False(t, result.Passed)
	assert.Equal(t, SeverityError, result.Severity)
}

func TestCheckTestDirectory_Exists(t *testing.T) {
	fs := setupFs(t)
	mkdirAll(t, fs, "/project/src/test/java")

	checker := NewChecker(fs, "/project")
	result := checker.checkTestDirectory()

	assert.True(t, result.Passed)
}

func TestCheckTestDirectory_Missing(t *testing.T) {
	fs := setupFs(t)
	mkdirAll(t, fs, "/project")

	checker := NewChecker(fs, "/project")
	result := checker.checkTestDirectory()

	assert.False(t, result.Passed)
	assert.Equal(t, SeverityWarning, result.Severity)
}

func TestCheckMainClass_Found(t *testing.T) {
	fs := setupFs(t)
	mkdirAll(t, fs, "/project/src/main/java/com/example")
	mainClass := `package com.example;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class Application {
    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }
}`
	writeFile(t, fs, "/project/src/main/java/com/example/Application.java", mainClass)

	checker := NewChecker(fs, "/project")
	result := checker.checkMainClass()

	assert.True(t, result.Passed)
	assert.Contains(t, result.Message, "@SpringBootApplication")
}

func TestCheckMainClass_NotFound(t *testing.T) {
	fs := setupFs(t)
	mkdirAll(t, fs, "/project/src/main/java/com/example")
	writeFile(t, fs, "/project/src/main/java/com/example/Service.java", "class Service {}")

	checker := NewChecker(fs, "/project")
	result := checker.checkMainClass()

	assert.False(t, result.Passed)
	assert.Equal(t, SeverityError, result.Severity)
}

func TestCheckConfigFile_YmlExists(t *testing.T) {
	fs := setupFs(t)
	mkdirAll(t, fs, "/project/src/main/resources")
	writeFile(t, fs, "/project/src/main/resources/application.yml", "server:\n  port: 8080")

	checker := NewChecker(fs, "/project")
	result := checker.checkConfigFile()

	assert.True(t, result.Passed)
	assert.Contains(t, result.Message, "application.yml")
}

func TestCheckConfigFile_PropertiesExists(t *testing.T) {
	fs := setupFs(t)
	mkdirAll(t, fs, "/project/src/main/resources")
	writeFile(t, fs, "/project/src/main/resources/application.properties", "server.port=8080")

	checker := NewChecker(fs, "/project")
	result := checker.checkConfigFile()

	assert.True(t, result.Passed)
	assert.Contains(t, result.Message, "application.properties")
}

func TestCheckConfigFile_Missing(t *testing.T) {
	fs := setupFs(t)
	mkdirAll(t, fs, "/project/src/main/resources")

	checker := NewChecker(fs, "/project")
	result := checker.checkConfigFile()

	assert.False(t, result.Passed)
	assert.Equal(t, SeverityWarning, result.Severity)
}

func TestCheckHardcodedSecrets_Found(t *testing.T) {
	fs := setupFs(t)
	mkdirAll(t, fs, "/project/src/main/resources")
	config := `spring:
  datasource:
    password: mysecretpassword123`
	writeFile(t, fs, "/project/src/main/resources/application.yml", config)

	checker := NewChecker(fs, "/project")
	result := checker.checkHardcodedSecrets()

	assert.False(t, result.Passed)
	assert.Equal(t, SeverityError, result.Severity)
	assert.Contains(t, result.Message, "secrets")
}

func TestCheckHardcodedSecrets_EnvVariable(t *testing.T) {
	fs := setupFs(t)
	mkdirAll(t, fs, "/project/src/main/resources")
	config := `spring:
  datasource:
    password: ${DB_PASSWORD}`
	writeFile(t, fs, "/project/src/main/resources/application.yml", config)

	checker := NewChecker(fs, "/project")
	result := checker.checkHardcodedSecrets()

	assert.True(t, result.Passed)
}

func TestCheckH2Scope_CompileScope(t *testing.T) {
	fs := setupFs(t)
	pom := `<project>
		<dependencies>
			<dependency>
				<groupId>com.h2database</groupId>
				<artifactId>h2</artifactId>
			</dependency>
		</dependencies>
	</project>`
	writeFile(t, fs, "/project/pom.xml", pom)

	checker := NewChecker(fs, "/project")
	checker.detectBuildTool()
	result := checker.checkH2Scope()

	assert.False(t, result.Passed)
	assert.Equal(t, SeverityWarning, result.Severity)
}

func TestCheckH2Scope_RuntimeScope(t *testing.T) {
	fs := setupFs(t)
	pom := `<project>
		<dependencies>
			<dependency>
				<groupId>com.h2database</groupId>
				<artifactId>h2</artifactId>
				<scope>runtime</scope>
			</dependency>
		</dependencies>
	</project>`
	writeFile(t, fs, "/project/pom.xml", pom)

	checker := NewChecker(fs, "/project")
	checker.detectBuildTool()
	result := checker.checkH2Scope()

	assert.True(t, result.Passed)
}

func TestCheckTestDependencies_Present(t *testing.T) {
	fs := setupFs(t)
	pom := `<project>
		<dependencies>
			<dependency>
				<groupId>org.springframework.boot</groupId>
				<artifactId>spring-boot-starter-test</artifactId>
				<scope>test</scope>
			</dependency>
		</dependencies>
	</project>`
	writeFile(t, fs, "/project/pom.xml", pom)

	checker := NewChecker(fs, "/project")
	checker.detectBuildTool()
	result := checker.checkTestDependencies()

	assert.True(t, result.Passed)
}

func TestCheckTestDependencies_Missing(t *testing.T) {
	fs := setupFs(t)
	pom := `<project>
		<dependencies>
			<dependency>
				<groupId>org.springframework.boot</groupId>
				<artifactId>spring-boot-starter-web</artifactId>
			</dependency>
		</dependencies>
	</project>`
	writeFile(t, fs, "/project/pom.xml", pom)

	checker := NewChecker(fs, "/project")
	checker.detectBuildTool()
	result := checker.checkTestDependencies()

	assert.False(t, result.Passed)
	assert.Equal(t, SeverityWarning, result.Severity)
}

func TestSuggestActuator_Missing(t *testing.T) {
	fs := setupFs(t)
	pom := `<project>
		<dependencies>
			<dependency>
				<groupId>org.springframework.boot</groupId>
				<artifactId>spring-boot-starter-web</artifactId>
			</dependency>
		</dependencies>
	</project>`
	writeFile(t, fs, "/project/pom.xml", pom)

	checker := NewChecker(fs, "/project")
	checker.detectBuildTool()
	result := checker.suggestActuator()

	assert.False(t, result.Passed)
	assert.Equal(t, SeveritySuggestion, result.Severity)
	assert.Contains(t, result.FixHint, "haft add actuator")
}

func TestSuggestActuator_Present(t *testing.T) {
	fs := setupFs(t)
	pom := `<project>
		<dependencies>
			<dependency>
				<groupId>org.springframework.boot</groupId>
				<artifactId>spring-boot-starter-actuator</artifactId>
			</dependency>
		</dependencies>
	</project>`
	writeFile(t, fs, "/project/pom.xml", pom)

	checker := NewChecker(fs, "/project")
	checker.detectBuildTool()
	result := checker.suggestActuator()

	assert.True(t, result.Passed)
}

func TestRunAllChecks(t *testing.T) {
	fs := setupFs(t)

	pom := `<project>
		<parent>
			<groupId>org.springframework.boot</groupId>
			<artifactId>spring-boot-starter-parent</artifactId>
			<version>3.2.0</version>
		</parent>
		<properties>
			<java.version>17</java.version>
		</properties>
		<dependencies>
			<dependency>
				<groupId>org.springframework.boot</groupId>
				<artifactId>spring-boot-starter-web</artifactId>
			</dependency>
			<dependency>
				<groupId>org.springframework.boot</groupId>
				<artifactId>spring-boot-starter-test</artifactId>
				<scope>test</scope>
			</dependency>
		</dependencies>
	</project>`
	writeFile(t, fs, "/project/pom.xml", pom)

	mkdirAll(t, fs, "/project/src/main/java/com/example")
	mkdirAll(t, fs, "/project/src/test/java")
	mkdirAll(t, fs, "/project/src/main/resources")

	mainClass := `@SpringBootApplication public class App {}`
	writeFile(t, fs, "/project/src/main/java/com/example/App.java", mainClass)
	writeFile(t, fs, "/project/src/main/resources/application.yml", "server:\n  port: 8080")

	checker := NewChecker(fs, "/project")
	results := checker.RunAllChecks()

	assert.NotEmpty(t, results)
	assert.GreaterOrEqual(t, len(results), 10)
}

func TestReportCalculateCounts(t *testing.T) {
	report := &Report{
		Results: []CheckResult{
			{Passed: true},
			{Passed: true},
			{Passed: false, Severity: SeverityError},
			{Passed: false, Severity: SeverityWarning},
			{Passed: false, Severity: SeverityWarning},
			{Passed: false, Severity: SeveritySuggestion},
		},
	}

	report.CalculateCounts()

	assert.Equal(t, 2, report.PassedCount)
	assert.Equal(t, 1, report.ErrorCount)
	assert.Equal(t, 2, report.WarningCount)
	assert.Equal(t, 1, report.SuggestionCount)
}

func TestReportHasIssues(t *testing.T) {
	report := &Report{ErrorCount: 1}
	assert.True(t, report.HasIssues())

	report = &Report{ErrorCount: 0}
	assert.False(t, report.HasIssues())
}

func TestReportHasWarnings(t *testing.T) {
	report := &Report{WarningCount: 1}
	assert.True(t, report.HasWarnings())

	report = &Report{WarningCount: 0}
	assert.False(t, report.HasWarnings())
}

func TestDoctorCommand_Structure(t *testing.T) {
	cmd := NewCommand()

	assert.Equal(t, "doctor", cmd.Use)
	assert.Contains(t, cmd.Aliases, "doc")
	assert.Contains(t, cmd.Aliases, "check")
	assert.Contains(t, cmd.Aliases, "health")
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
}

func TestDoctorCommand_Flags(t *testing.T) {
	cmd := NewCommand()

	jsonFlag := cmd.Flags().Lookup("json")
	assert.NotNil(t, jsonFlag)
	assert.Equal(t, "false", jsonFlag.DefValue)

	strictFlag := cmd.Flags().Lookup("strict")
	assert.NotNil(t, strictFlag)
	assert.Equal(t, "false", strictFlag.DefValue)

	categoryFlag := cmd.Flags().Lookup("category")
	assert.NotNil(t, categoryFlag)
	assert.Empty(t, categoryFlag.DefValue)
}

func TestFormatReport_JSON(t *testing.T) {
	report := &Report{
		ProjectPath: "/test",
		ProjectName: "test-project",
		BuildTool:   "Maven",
		Results: []CheckResult{
			{Name: "test", Passed: true, Message: "Test passed"},
		},
		PassedCount: 1,
	}

	output := FormatReport(report, Options{JSON: true})

	assert.Contains(t, output, `"project_name": "test-project"`)
	assert.Contains(t, output, `"build_tool": "Maven"`)
	assert.Contains(t, output, `"passed": true`)
}

func TestFormatReport_Text(t *testing.T) {
	report := &Report{
		ProjectPath: "/test",
		ProjectName: "test-project",
		BuildTool:   "Maven",
		Results: []CheckResult{
			{Name: "passed_check", Passed: true, Message: "Check passed"},
			{Name: "error_check", Passed: false, Severity: SeverityError, Message: "Error found"},
			{Name: "warning_check", Passed: false, Severity: SeverityWarning, Message: "Warning found"},
		},
		PassedCount:  1,
		ErrorCount:   1,
		WarningCount: 1,
	}

	output := FormatReport(report, Options{JSON: false})

	assert.Contains(t, output, "Haft Doctor")
	assert.Contains(t, output, "test-project")
	assert.Contains(t, output, "Maven")
	assert.Contains(t, output, "Check passed")
	assert.Contains(t, output, "Error found")
	assert.Contains(t, output, "Warning found")
}

func TestFilterByCategory(t *testing.T) {
	results := []CheckResult{
		{Name: "build1", Category: CategoryBuild},
		{Name: "build2", Category: CategoryBuild},
		{Name: "security1", Category: CategorySecurity},
		{Name: "config1", Category: CategoryConfig},
	}

	buildResults := filterByCategory(results, CategoryBuild)
	assert.Len(t, buildResults, 2)

	securityResults := filterByCategory(results, CategorySecurity)
	assert.Len(t, securityResults, 1)
}

func TestCheckDockerfile_NotPresent(t *testing.T) {
	fs := setupFs(t)
	mkdirAll(t, fs, "/project")

	checker := NewChecker(fs, "/project")
	result := checker.checkDockerfile()

	assert.True(t, result.Passed)
	assert.Contains(t, result.Message, "No Dockerfile")
}

func TestCheckDockerfile_MultiStageOptimized(t *testing.T) {
	fs := setupFs(t)
	dockerfile := `FROM maven:3.9-eclipse-temurin-17 AS build
WORKDIR /app
COPY pom.xml .
RUN mvn dependency:go-offline
COPY src ./src
RUN mvn clean package -DskipTests

FROM eclipse-temurin:17-jre-alpine
WORKDIR /app
COPY --from=build /app/target/*.jar app.jar
USER 1000
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s CMD wget -qO- http://localhost:8080/actuator/health || exit 1
ENTRYPOINT ["java", "-jar", "app.jar"]`
	writeFile(t, fs, "/project/Dockerfile", dockerfile)
	writeFile(t, fs, "/project/.dockerignore", ".git\ntarget/\n")

	checker := NewChecker(fs, "/project")
	result := checker.checkDockerfile()

	assert.True(t, result.Passed)
	assert.Contains(t, result.Message, "multi-stage")
}

func TestCheckDockerfile_MissingBestPractices(t *testing.T) {
	fs := setupFs(t)
	dockerfile := `FROM eclipse-temurin:latest
COPY . .
RUN mvn clean package
ENTRYPOINT ["java", "-jar", "target/app.jar"]`
	writeFile(t, fs, "/project/Dockerfile", dockerfile)

	checker := NewChecker(fs, "/project")
	result := checker.checkDockerfile()

	assert.False(t, result.Passed)
	assert.Equal(t, SeverityWarning, result.Severity)
	assert.Contains(t, result.Details, ":latest")
}

func TestCheckDockerfile_NoUserInstruction(t *testing.T) {
	fs := setupFs(t)
	dockerfile := `FROM eclipse-temurin:17-jre
WORKDIR /app
COPY target/*.jar app.jar
ENTRYPOINT ["java", "-jar", "app.jar"]`
	writeFile(t, fs, "/project/Dockerfile", dockerfile)
	writeFile(t, fs, "/project/.dockerignore", ".git\n")

	checker := NewChecker(fs, "/project")
	result := checker.checkDockerfile()

	assert.False(t, result.Passed)
	assert.Contains(t, result.Details, "EXPOSE")
}

func TestCheckDockerCompose_NotPresent(t *testing.T) {
	fs := setupFs(t)
	mkdirAll(t, fs, "/project")

	checker := NewChecker(fs, "/project")
	result := checker.checkDockerCompose()

	assert.True(t, result.Passed)
	assert.Contains(t, result.Message, "No Docker Compose")
}

func TestCheckDockerCompose_Valid(t *testing.T) {
	fs := setupFs(t)
	compose := `version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_PASSWORD=${DB_PASSWORD}
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:8080/actuator/health"]
      interval: 30s
    restart: unless-stopped
    depends_on:
      - db
  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_PASSWORD=${DB_PASSWORD}
    restart: unless-stopped`
	writeFile(t, fs, "/project/docker-compose.yml", compose)

	checker := NewChecker(fs, "/project")
	result := checker.checkDockerCompose()

	assert.True(t, result.Passed)
	assert.Contains(t, result.Message, "Docker Compose configured")
}

func TestCheckDockerCompose_HardcodedPassword(t *testing.T) {
	fs := setupFs(t)
	compose := `version: '3.8'
services:
  db:
    image: postgres:15
    environment:
      POSTGRES_PASSWORD: mysecretpassword`
	writeFile(t, fs, "/project/docker-compose.yml", compose)

	checker := NewChecker(fs, "/project")
	result := checker.checkDockerCompose()

	assert.False(t, result.Passed)
	assert.Equal(t, SeverityWarning, result.Severity)
	assert.Contains(t, result.Details, "hardcoded")
}

func TestCheckDockerignore_NotPresentWithDockerfile(t *testing.T) {
	fs := setupFs(t)
	dockerfile := `FROM eclipse-temurin:17-jre
COPY target/*.jar app.jar`
	writeFile(t, fs, "/project/Dockerfile", dockerfile)

	checker := NewChecker(fs, "/project")
	result := checker.checkDockerignore()

	assert.False(t, result.Passed)
	assert.Equal(t, SeverityWarning, result.Severity)
	assert.Contains(t, result.Message, "No .dockerignore")
}

func TestCheckDockerignore_Valid(t *testing.T) {
	fs := setupFs(t)
	writeFile(t, fs, "/project/Dockerfile", "FROM eclipse-temurin:17")
	dockerignore := `.git
target/
build/
.idea/
*.log
.env`
	writeFile(t, fs, "/project/.dockerignore", dockerignore)

	checker := NewChecker(fs, "/project")
	result := checker.checkDockerignore()

	assert.True(t, result.Passed)
	assert.Contains(t, result.Message, "configured")
}

func TestCheckDockerignore_MissingPatterns(t *testing.T) {
	fs := setupFs(t)
	writeFile(t, fs, "/project/Dockerfile", "FROM eclipse-temurin:17")
	dockerignore := `.idea/
*.log`
	writeFile(t, fs, "/project/.dockerignore", dockerignore)

	checker := NewChecker(fs, "/project")
	result := checker.checkDockerignore()

	assert.True(t, result.Passed)
	assert.Equal(t, SeverityInfo, result.Severity)
	assert.Contains(t, result.Details, ".git")
}

func TestSuggestDocker_Missing(t *testing.T) {
	fs := setupFs(t)
	writeFile(t, fs, "/project/pom.xml", "<project></project>")

	checker := NewChecker(fs, "/project")
	result := checker.suggestDocker()

	assert.False(t, result.Passed)
	assert.Equal(t, SeveritySuggestion, result.Severity)
	assert.Contains(t, result.FixHint, "haft dockerize")
}

func TestSuggestDocker_Present(t *testing.T) {
	fs := setupFs(t)
	writeFile(t, fs, "/project/pom.xml", "<project></project>")
	writeFile(t, fs, "/project/Dockerfile", "FROM eclipse-temurin:17")

	checker := NewChecker(fs, "/project")
	result := checker.suggestDocker()

	assert.True(t, result.Passed)
}

func TestFilterByCategory_Docker(t *testing.T) {
	results := []CheckResult{
		{Name: "dockerfile", Category: CategoryDocker},
		{Name: "compose", Category: CategoryDocker},
		{Name: "security", Category: CategorySecurity},
	}

	dockerResults := filterByCategory(results, CategoryDocker)
	assert.Len(t, dockerResults, 2)
}
