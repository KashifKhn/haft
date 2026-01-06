package doctor

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/afero"
)

type Checker struct {
	fs          afero.Fs
	projectPath string
	buildFile   string
	buildTool   string
	content     string
}

func NewChecker(fs afero.Fs, projectPath string) *Checker {
	return &Checker{
		fs:          fs,
		projectPath: projectPath,
	}
}

func (c *Checker) detectBuildTool() {
	pomPath := filepath.Join(c.projectPath, "pom.xml")
	gradlePath := filepath.Join(c.projectPath, "build.gradle")
	gradleKtsPath := filepath.Join(c.projectPath, "build.gradle.kts")

	if exists, _ := afero.Exists(c.fs, pomPath); exists {
		c.buildFile = pomPath
		c.buildTool = "Maven"
		content, _ := afero.ReadFile(c.fs, pomPath)
		c.content = string(content)
	} else if exists, _ := afero.Exists(c.fs, gradleKtsPath); exists {
		c.buildFile = gradleKtsPath
		c.buildTool = "Gradle (Kotlin)"
		content, _ := afero.ReadFile(c.fs, gradleKtsPath)
		c.content = string(content)
	} else if exists, _ := afero.Exists(c.fs, gradlePath); exists {
		c.buildFile = gradlePath
		c.buildTool = "Gradle"
		content, _ := afero.ReadFile(c.fs, gradlePath)
		c.content = string(content)
	}
}

func (c *Checker) RunAllChecks() []CheckResult {
	c.detectBuildTool()

	var results []CheckResult

	results = append(results, c.checkBuildFile())
	results = append(results, c.checkSpringBootParent())
	results = append(results, c.checkJavaVersion())
	results = append(results, c.checkSourceDirectory())
	results = append(results, c.checkTestDirectory())
	results = append(results, c.checkMainClass())
	results = append(results, c.checkConfigFile())
	results = append(results, c.checkHardcodedSecrets())
	results = append(results, c.checkH2Scope())
	results = append(results, c.checkDevToolsScope())
	results = append(results, c.checkTestDependencies())
	results = append(results, c.suggestActuator())
	results = append(results, c.suggestValidation())
	results = append(results, c.suggestOpenAPI())
	results = append(results, c.checkLombokConfig())

	return results
}

func (c *Checker) checkBuildFile() CheckResult {
	result := CheckResult{
		Name:     "build_file",
		Category: CategoryBuild,
		Severity: SeverityError,
	}

	if c.buildFile == "" {
		result.Passed = false
		result.Message = "No build file found"
		result.Details = "Expected pom.xml or build.gradle(.kts)"
		result.FixHint = "Run 'haft init' to create a new project"
		return result
	}

	result.Passed = true
	result.Message = "Build file exists (" + c.buildTool + ")"
	return result
}

func (c *Checker) checkSpringBootParent() CheckResult {
	result := CheckResult{
		Name:     "spring_boot_config",
		Category: CategoryBuild,
		Severity: SeverityError,
	}

	if c.content == "" {
		result.Passed = false
		result.Message = "Cannot check Spring Boot configuration"
		result.Details = "Build file not readable"
		return result
	}

	hasSpringBoot := strings.Contains(c.content, "spring-boot") ||
		strings.Contains(c.content, "org.springframework.boot")

	if !hasSpringBoot {
		result.Passed = false
		result.Message = "Spring Boot not configured"
		result.Details = "No Spring Boot parent or plugin found"
		result.FixHint = "Add spring-boot-starter-parent or Spring Boot Gradle plugin"
		return result
	}

	result.Passed = true
	result.Message = "Spring Boot configured"
	return result
}

func (c *Checker) checkJavaVersion() CheckResult {
	result := CheckResult{
		Name:     "java_version",
		Category: CategoryBuild,
		Severity: SeverityWarning,
	}

	if c.content == "" {
		result.Passed = false
		result.Message = "Cannot check Java version"
		return result
	}

	javaVersionPatterns := []string{
		`<java\.version>(\d+)</java\.version>`,
		`sourceCompatibility\s*=\s*['"]?(\d+)['"]?`,
		`JavaVersion\.VERSION_(\d+)`,
		`java\s*{\s*sourceCompatibility\s*=\s*JavaVersion\.VERSION_(\d+)`,
	}

	for _, pattern := range javaVersionPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(c.content)
		if len(matches) > 1 {
			version := matches[1]
			if version == "8" || version == "1.8" {
				result.Passed = true
				result.Severity = SeverityWarning
				result.Message = "Java 8 detected - consider upgrading"
				result.Details = "Java 8 is end-of-life. Java 17 or 21 recommended"
				result.FixHint = "Update java.version to 17 or 21"
				return result
			}
			result.Passed = true
			result.Message = "Java " + version + " configured"
			return result
		}
	}

	result.Passed = false
	result.Message = "Java version not specified"
	result.Details = "Explicit Java version recommended for reproducible builds"
	result.FixHint = "Add <java.version>17</java.version> to pom.xml properties"
	return result
}

func (c *Checker) checkSourceDirectory() CheckResult {
	result := CheckResult{
		Name:     "source_directory",
		Category: CategorySource,
		Severity: SeverityError,
	}

	javaPath := filepath.Join(c.projectPath, "src", "main", "java")
	kotlinPath := filepath.Join(c.projectPath, "src", "main", "kotlin")

	javaExists, _ := afero.DirExists(c.fs, javaPath)
	kotlinExists, _ := afero.DirExists(c.fs, kotlinPath)

	if !javaExists && !kotlinExists {
		result.Passed = false
		result.Message = "Source directory missing"
		result.Details = "Expected src/main/java or src/main/kotlin"
		result.FixHint = "Create src/main/java directory"
		return result
	}

	lang := "Java"
	if kotlinExists {
		lang = "Kotlin"
	}
	result.Passed = true
	result.Message = "Source directory exists (" + lang + ")"
	return result
}

func (c *Checker) checkTestDirectory() CheckResult {
	result := CheckResult{
		Name:     "test_directory",
		Category: CategorySource,
		Severity: SeverityWarning,
	}

	javaPath := filepath.Join(c.projectPath, "src", "test", "java")
	kotlinPath := filepath.Join(c.projectPath, "src", "test", "kotlin")

	javaExists, _ := afero.DirExists(c.fs, javaPath)
	kotlinExists, _ := afero.DirExists(c.fs, kotlinPath)

	if !javaExists && !kotlinExists {
		result.Passed = false
		result.Message = "Test directory missing"
		result.Details = "No src/test/java or src/test/kotlin found"
		result.FixHint = "Create src/test/java and add test classes"
		return result
	}

	result.Passed = true
	result.Message = "Test directory exists"
	return result
}

func (c *Checker) checkMainClass() CheckResult {
	result := CheckResult{
		Name:     "main_class",
		Category: CategorySource,
		Severity: SeverityError,
	}

	srcDirs := []string{
		filepath.Join(c.projectPath, "src", "main", "java"),
		filepath.Join(c.projectPath, "src", "main", "kotlin"),
	}

	for _, srcDir := range srcDirs {
		if exists, _ := afero.DirExists(c.fs, srcDir); !exists {
			continue
		}

		found := false
		_ = afero.Walk(c.fs, srcDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(path, ".java") && !strings.HasSuffix(path, ".kt") {
				return nil
			}
			content, err := afero.ReadFile(c.fs, path)
			if err != nil {
				return nil
			}
			if strings.Contains(string(content), "@SpringBootApplication") {
				found = true
				return filepath.SkipAll
			}
			return nil
		})

		if found {
			result.Passed = true
			result.Message = "Main class with @SpringBootApplication found"
			return result
		}
	}

	result.Passed = false
	result.Message = "No @SpringBootApplication class found"
	result.Details = "Every Spring Boot app needs a main class with @SpringBootApplication"
	result.FixHint = "Create Application.java with @SpringBootApplication annotation"
	return result
}

func (c *Checker) checkConfigFile() CheckResult {
	result := CheckResult{
		Name:     "config_file",
		Category: CategoryConfig,
		Severity: SeverityWarning,
	}

	configFiles := []string{
		"src/main/resources/application.yml",
		"src/main/resources/application.yaml",
		"src/main/resources/application.properties",
	}

	for _, cf := range configFiles {
		path := filepath.Join(c.projectPath, cf)
		if exists, _ := afero.Exists(c.fs, path); exists {
			result.Passed = true
			result.Message = "Configuration file found (" + filepath.Base(cf) + ")"
			return result
		}
	}

	result.Passed = false
	result.Message = "No configuration file found"
	result.Details = "Expected application.yml or application.properties"
	result.FixHint = "Create src/main/resources/application.yml"
	return result
}

func (c *Checker) checkHardcodedSecrets() CheckResult {
	result := CheckResult{
		Name:     "hardcoded_secrets",
		Category: CategorySecurity,
		Severity: SeverityError,
	}

	configFiles := []string{
		"src/main/resources/application.yml",
		"src/main/resources/application.yaml",
		"src/main/resources/application.properties",
	}

	secretPatterns := []struct {
		pattern string
		name    string
	}{
		{`password\s*[:=]\s*['""]?[^${\s][^'"\s]+`, "password"},
		{`secret\s*[:=]\s*['""]?[^${\s][^'"\s]+`, "secret"},
		{`api[_-]?key\s*[:=]\s*['""]?[^${\s][^'"\s]+`, "api-key"},
		{`private[_-]?key\s*[:=]\s*['""]?[^${\s][^'"\s]+`, "private-key"},
		{`access[_-]?token\s*[:=]\s*['""]?[^${\s][^'"\s]+`, "access-token"},
	}

	for _, cf := range configFiles {
		path := filepath.Join(c.projectPath, cf)
		content, err := afero.ReadFile(c.fs, path)
		if err != nil {
			continue
		}

		contentStr := strings.ToLower(string(content))
		for _, sp := range secretPatterns {
			re := regexp.MustCompile(`(?i)` + sp.pattern)
			if re.MatchString(contentStr) {
				if !strings.Contains(contentStr, "${") {
					result.Passed = false
					result.Message = "Potential hardcoded secrets detected"
					result.Details = "Found hardcoded " + sp.name + " in " + filepath.Base(cf)
					result.FixHint = "Use environment variables: ${DB_PASSWORD}"
					return result
				}
			}
		}
	}

	result.Passed = true
	result.Message = "No hardcoded secrets detected"
	return result
}

func (c *Checker) checkH2Scope() CheckResult {
	result := CheckResult{
		Name:     "h2_scope",
		Category: CategorySecurity,
		Severity: SeverityWarning,
	}

	if c.content == "" {
		result.Passed = true
		result.Message = "H2 database not detected"
		return result
	}

	hasH2 := strings.Contains(c.content, "com.h2database") || strings.Contains(c.content, "h2database:h2")

	if !hasH2 {
		result.Passed = true
		result.Message = "H2 database not used"
		return result
	}

	isTestScope := strings.Contains(c.content, "<scope>test</scope>") ||
		strings.Contains(c.content, "<scope>runtime</scope>") ||
		strings.Contains(c.content, "testImplementation") ||
		strings.Contains(c.content, "runtimeOnly")

	if !isTestScope {
		result.Passed = false
		result.Message = "H2 database in compile scope"
		result.Details = "H2 should be test or runtime scope only"
		result.FixHint = "Change H2 scope to <scope>runtime</scope> or testImplementation"
		return result
	}

	result.Passed = true
	result.Message = "H2 database correctly scoped"
	return result
}

func (c *Checker) checkDevToolsScope() CheckResult {
	result := CheckResult{
		Name:     "devtools_scope",
		Category: CategorySecurity,
		Severity: SeverityWarning,
	}

	if c.content == "" {
		result.Passed = true
		result.Message = "DevTools not detected"
		return result
	}

	hasDevTools := strings.Contains(c.content, "spring-boot-devtools")

	if !hasDevTools {
		result.Passed = true
		result.Message = "DevTools not used"
		return result
	}

	if c.buildTool == "Maven" {
		if !strings.Contains(c.content, "<optional>true</optional>") {
			result.Passed = false
			result.Message = "DevTools not marked as optional"
			result.Details = "DevTools should be optional to prevent inclusion in production"
			result.FixHint = "Add <optional>true</optional> to DevTools dependency"
			return result
		}
	}

	result.Passed = true
	result.Message = "DevTools correctly configured"
	return result
}

func (c *Checker) checkTestDependencies() CheckResult {
	result := CheckResult{
		Name:     "test_dependencies",
		Category: CategoryDependencies,
		Severity: SeverityWarning,
	}

	if c.content == "" {
		result.Passed = false
		result.Message = "Cannot check test dependencies"
		return result
	}

	hasTestStarter := strings.Contains(c.content, "spring-boot-starter-test")

	if !hasTestStarter {
		result.Passed = false
		result.Message = "No test dependencies found"
		result.Details = "spring-boot-starter-test not found"
		result.FixHint = "Run: haft add test"
		return result
	}

	result.Passed = true
	result.Message = "Test dependencies configured"
	return result
}

func (c *Checker) suggestActuator() CheckResult {
	result := CheckResult{
		Name:     "suggest_actuator",
		Category: CategoryDependencies,
		Severity: SeveritySuggestion,
	}

	if c.content == "" {
		result.Passed = true
		return result
	}

	hasActuator := strings.Contains(c.content, "spring-boot-starter-actuator")

	if hasActuator {
		result.Passed = true
		result.Message = "Actuator configured for monitoring"
		return result
	}

	result.Passed = false
	result.Message = "Consider adding Actuator"
	result.Details = "Actuator provides health checks, metrics, and monitoring endpoints"
	result.FixHint = "Run: haft add actuator"
	return result
}

func (c *Checker) suggestValidation() CheckResult {
	result := CheckResult{
		Name:     "suggest_validation",
		Category: CategoryDependencies,
		Severity: SeveritySuggestion,
	}

	if c.content == "" {
		result.Passed = true
		return result
	}

	hasValidation := strings.Contains(c.content, "spring-boot-starter-validation") ||
		strings.Contains(c.content, "hibernate-validator")

	if hasValidation {
		result.Passed = true
		result.Message = "Validation configured"
		return result
	}

	result.Passed = false
	result.Message = "Consider adding Validation"
	result.Details = "Bean validation helps ensure data integrity with @Valid annotations"
	result.FixHint = "Run: haft add validation"
	return result
}

func (c *Checker) suggestOpenAPI() CheckResult {
	result := CheckResult{
		Name:     "suggest_openapi",
		Category: CategoryDependencies,
		Severity: SeveritySuggestion,
	}

	if c.content == "" {
		result.Passed = true
		return result
	}

	hasOpenAPI := strings.Contains(c.content, "springdoc-openapi") ||
		strings.Contains(c.content, "springfox")

	if hasOpenAPI {
		result.Passed = true
		result.Message = "API documentation configured"
		return result
	}

	hasWeb := strings.Contains(c.content, "spring-boot-starter-web")
	if !hasWeb {
		result.Passed = true
		result.Message = "Not a web project"
		return result
	}

	result.Passed = false
	result.Message = "Consider adding OpenAPI documentation"
	result.Details = "springdoc-openapi provides automatic API documentation"
	result.FixHint = "Run: haft add openapi"
	return result
}

func (c *Checker) checkLombokConfig() CheckResult {
	result := CheckResult{
		Name:     "lombok_config",
		Category: CategoryBestPractice,
		Severity: SeverityInfo,
	}

	if c.content == "" {
		result.Passed = true
		return result
	}

	hasLombok := strings.Contains(c.content, "org.projectlombok") ||
		strings.Contains(c.content, "lombok")

	if !hasLombok {
		result.Passed = true
		result.Message = "Lombok not used"
		return result
	}

	lombokConfigPath := filepath.Join(c.projectPath, "lombok.config")
	if exists, _ := afero.Exists(c.fs, lombokConfigPath); exists {
		result.Passed = true
		result.Message = "Lombok configured with lombok.config"
		return result
	}

	result.Passed = true
	result.Severity = SeverityInfo
	result.Message = "Lombok used without lombok.config"
	result.Details = "Consider adding lombok.config for consistent behavior"
	result.FixHint = "Create lombok.config with: config.stopBubbling = true"
	return result
}
