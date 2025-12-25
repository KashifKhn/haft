package gradle

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/spf13/afero"
)

func init() {
	buildtool.Register(buildtool.Gradle, func(fs afero.Fs) buildtool.Parser {
		return NewParserWithFs(fs, false)
	})
	buildtool.Register(buildtool.GradleKotln, func(fs afero.Fs) buildtool.Parser {
		return NewParserWithFs(fs, true)
	})
}

type Parser struct {
	fs       afero.Fs
	isKotlin bool
}

type GradleProject struct {
	Content      string
	FilePath     string
	IsKotlin     bool
	Group        string
	Version      string
	SourceCompat string
	TargetCompat string
	JavaVersion  string
}

func NewParser(isKotlin bool) *Parser {
	return &Parser{fs: afero.NewOsFs(), isKotlin: isKotlin}
}

func NewParserWithFs(fs afero.Fs, isKotlin bool) *Parser {
	return &Parser{fs: fs, isKotlin: isKotlin}
}

func (p *Parser) Type() buildtool.Type {
	if p.isKotlin {
		return buildtool.GradleKotln
	}
	return buildtool.Gradle
}

func (p *Parser) FindBuildFile(startDir string) (string, error) {
	current := startDir
	for {
		var buildFile string
		if p.isKotlin {
			buildFile = filepath.Join(current, "build.gradle.kts")
		} else {
			buildFile = filepath.Join(current, "build.gradle")
		}

		if _, err := p.fs.Stat(buildFile); err == nil {
			return buildFile, nil
		}

		parent := filepath.Dir(current)
		if parent == current || parent == "" {
			break
		}
		current = parent
	}

	fileName := "build.gradle"
	if p.isKotlin {
		fileName = "build.gradle.kts"
	}
	return "", fmt.Errorf("%s not found in %s or any parent directory", fileName, startDir)
}

func (p *Parser) Parse(path string) (*buildtool.Project, error) {
	data, err := afero.ReadFile(p.fs, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read gradle build file: %w", err)
	}

	content := string(data)
	isKotlin := strings.HasSuffix(path, ".kts")

	gradleProject := &GradleProject{
		Content:  content,
		FilePath: path,
		IsKotlin: isKotlin,
	}

	project := &buildtool.Project{
		BuildTool: p.Type(),
		Raw:       gradleProject,
	}

	project.GroupId = p.extractGroup(content, isKotlin)
	project.Version = p.extractVersion(content, isKotlin)
	project.ArtifactId = p.extractArtifactName(path)
	project.JavaVersion = p.extractJavaVersion(content, isKotlin)
	project.SpringBootVersion = p.extractSpringBootVersion(content, isKotlin)
	project.Dependencies = p.extractDependencies(content, isKotlin)

	gradleProject.Group = project.GroupId
	gradleProject.Version = project.Version
	gradleProject.JavaVersion = project.JavaVersion

	return project, nil
}

func (p *Parser) Write(path string, project *buildtool.Project) error {
	gradleProject := p.getGradleProject(project)
	if gradleProject == nil {
		return fmt.Errorf("invalid gradle project")
	}
	return afero.WriteFile(p.fs, path, []byte(gradleProject.Content), 0644)
}

func (p *Parser) getGradleProject(project *buildtool.Project) *GradleProject {
	if raw, ok := project.Raw.(*GradleProject); ok {
		return raw
	}
	return nil
}

func (p *Parser) extractGroup(content string, isKotlin bool) string {
	var patterns []*regexp.Regexp
	if isKotlin {
		patterns = []*regexp.Regexp{
			regexp.MustCompile(`group\s*=\s*"([^"]+)"`),
			regexp.MustCompile(`group\s*=\s*'([^']+)'`),
		}
	} else {
		patterns = []*regexp.Regexp{
			regexp.MustCompile(`group\s*=\s*['"]([^'"]+)['"]`),
			regexp.MustCompile(`group\s+['"]([^'"]+)['"]`),
		}
	}

	for _, pattern := range patterns {
		if match := pattern.FindStringSubmatch(content); len(match) > 1 {
			return match[1]
		}
	}
	return ""
}

func (p *Parser) extractVersion(content string, isKotlin bool) string {
	var patterns []*regexp.Regexp
	if isKotlin {
		patterns = []*regexp.Regexp{
			regexp.MustCompile(`(?m)^version\s*=\s*"([^"]+)"`),
			regexp.MustCompile(`(?m)^version\s*=\s*'([^']+)'`),
		}
	} else {
		patterns = []*regexp.Regexp{
			regexp.MustCompile(`(?m)^version\s*=\s*['"]([^'"]+)['"]`),
			regexp.MustCompile(`(?m)^version\s+['"]([^'"]+)['"]`),
		}
	}

	for _, pattern := range patterns {
		if match := pattern.FindStringSubmatch(content); len(match) > 1 {
			return match[1]
		}
	}
	return ""
}

func (p *Parser) extractArtifactName(path string) string {
	dir := filepath.Dir(path)
	settingsFile := filepath.Join(dir, "settings.gradle")
	if p.isKotlin {
		settingsFile = filepath.Join(dir, "settings.gradle.kts")
	}

	data, err := afero.ReadFile(p.fs, settingsFile)
	if err == nil {
		content := string(data)
		patterns := []*regexp.Regexp{
			regexp.MustCompile(`rootProject\.name\s*=\s*["']([^"']+)["']`),
			regexp.MustCompile(`rootProject\.name\s*=\s*"([^"]+)"`),
		}
		for _, pattern := range patterns {
			if match := pattern.FindStringSubmatch(content); len(match) > 1 {
				return match[1]
			}
		}
	}

	return filepath.Base(dir)
}

func (p *Parser) extractJavaVersion(content string, isKotlin bool) string {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`sourceCompatibility\s*=\s*['"]?(\d+)['"]?`),
		regexp.MustCompile(`sourceCompatibility\s*=\s*JavaVersion\.VERSION_(\d+)`),
		regexp.MustCompile(`java\.sourceCompatibility\s*=\s*JavaVersion\.VERSION_(\d+)`),
		regexp.MustCompile(`toolchain\s*\{[^}]*languageVersion\.set\(JavaLanguageVersion\.of\((\d+)\)\)`),
		regexp.MustCompile(`jvmToolchain\((\d+)\)`),
		regexp.MustCompile(`JavaLanguageVersion\.of\((\d+)\)`),
	}

	for _, pattern := range patterns {
		if match := pattern.FindStringSubmatch(content); len(match) > 1 {
			return match[1]
		}
	}
	return ""
}

func (p *Parser) extractSpringBootVersion(content string, isKotlin bool) string {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`id\s*\(?["']org\.springframework\.boot["']\)?\s*version\s*["']([^"']+)["']`),
		regexp.MustCompile(`org\.springframework\.boot["']?\)?\s*version\s*["']([^"']+)["']`),
		regexp.MustCompile(`springBootVersion\s*=\s*["']([^"']+)["']`),
	}

	for _, pattern := range patterns {
		if match := pattern.FindStringSubmatch(content); len(match) > 1 {
			return match[1]
		}
	}
	return ""
}

func (p *Parser) extractDependencies(content string, isKotlin bool) []buildtool.Dependency {
	var deps []buildtool.Dependency

	var patterns []*regexp.Regexp
	if isKotlin {
		patterns = []*regexp.Regexp{
			regexp.MustCompile(`(implementation|compileOnly|runtimeOnly|testImplementation|testCompileOnly|testRuntimeOnly|annotationProcessor|kapt)\s*\(\s*"([^"]+)"\s*\)`),
			regexp.MustCompile(`(implementation|compileOnly|runtimeOnly|testImplementation|testCompileOnly|testRuntimeOnly|annotationProcessor|kapt)\s*\(\s*'([^']+)'\s*\)`),
		}
	} else {
		patterns = []*regexp.Regexp{
			regexp.MustCompile(`(implementation|compileOnly|runtimeOnly|testImplementation|testCompileOnly|testRuntimeOnly|annotationProcessor)\s+['"]([^'"]+)['"]`),
			regexp.MustCompile(`(implementation|compileOnly|runtimeOnly|testImplementation|testCompileOnly|testRuntimeOnly|annotationProcessor)\s*\(\s*['"]([^'"]+)['"]\s*\)`),
		}
	}

	seen := make(map[string]bool)

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) < 3 {
				continue
			}
			scope := match[1]
			coords := match[2]

			if seen[coords] {
				continue
			}
			seen[coords] = true

			dep := p.parseCoordinates(coords)
			dep.Scope = p.mapGradleScope(scope)
			deps = append(deps, dep)
		}
	}

	return deps
}

func (p *Parser) parseCoordinates(coords string) buildtool.Dependency {
	parts := strings.Split(coords, ":")
	dep := buildtool.Dependency{}

	if len(parts) >= 1 {
		dep.GroupId = parts[0]
	}
	if len(parts) >= 2 {
		dep.ArtifactId = parts[1]
	}
	if len(parts) >= 3 {
		dep.Version = parts[2]
	}

	return dep
}

func (p *Parser) mapGradleScope(gradleScope string) string {
	switch gradleScope {
	case "implementation":
		return "compile"
	case "compileOnly":
		return "provided"
	case "runtimeOnly":
		return "runtime"
	case "testImplementation", "testCompileOnly":
		return "test"
	case "testRuntimeOnly":
		return "test"
	case "annotationProcessor", "kapt":
		return "provided"
	default:
		return "compile"
	}
}

func (p *Parser) mapScopeToGradle(mavenScope string, isKotlin bool) string {
	switch mavenScope {
	case "compile", "":
		return "implementation"
	case "provided":
		return "compileOnly"
	case "runtime":
		return "runtimeOnly"
	case "test":
		return "testImplementation"
	default:
		return "implementation"
	}
}

func (p *Parser) HasDependency(project *buildtool.Project, groupId, artifactId string) bool {
	for _, dep := range project.Dependencies {
		if dep.GroupId == groupId && dep.ArtifactId == artifactId {
			return true
		}
	}
	return false
}

func (p *Parser) AddDependency(project *buildtool.Project, dep buildtool.Dependency) {
	if p.HasDependency(project, dep.GroupId, dep.ArtifactId) {
		return
	}

	project.Dependencies = append(project.Dependencies, dep)

	gradleProject := p.getGradleProject(project)
	if gradleProject == nil {
		return
	}

	gradleScope := p.mapScopeToGradle(dep.Scope, gradleProject.IsKotlin)
	depLine := p.formatDependencyLine(dep, gradleScope, gradleProject.IsKotlin)

	gradleProject.Content = p.insertDependency(gradleProject.Content, depLine, gradleProject.IsKotlin)
}

func (p *Parser) formatDependencyLine(dep buildtool.Dependency, scope string, isKotlin bool) string {
	coords := dep.GroupId + ":" + dep.ArtifactId
	if dep.Version != "" {
		coords += ":" + dep.Version
	}

	if isKotlin {
		return fmt.Sprintf("\t%s(\"%s\")", scope, coords)
	}
	return fmt.Sprintf("\t%s '%s'", scope, coords)
}

func (p *Parser) insertDependency(content, depLine string, isKotlin bool) string {
	depBlockRegex := regexp.MustCompile(`(?s)(dependencies\s*\{)([^}]*)(\})`)

	if match := depBlockRegex.FindStringSubmatchIndex(content); match != nil {
		closingBracePos := match[5]
		return content[:closingBracePos] + depLine + "\n" + content[closingBracePos:]
	}

	return content + "\n\ndependencies {\n" + depLine + "\n}\n"
}

func (p *Parser) RemoveDependency(project *buildtool.Project, groupId, artifactId string) bool {
	found := false
	var newDeps []buildtool.Dependency
	for _, dep := range project.Dependencies {
		if dep.GroupId == groupId && dep.ArtifactId == artifactId {
			found = true
			continue
		}
		newDeps = append(newDeps, dep)
	}
	project.Dependencies = newDeps

	gradleProject := p.getGradleProject(project)
	if gradleProject == nil {
		return found
	}

	gradleProject.Content = p.removeDependencyFromContent(gradleProject.Content, groupId, artifactId, gradleProject.IsKotlin)

	return found
}

func (p *Parser) removeDependencyFromContent(content, groupId, artifactId string, isKotlin bool) string {
	coords := groupId + ":" + artifactId

	var patterns []*regexp.Regexp
	if isKotlin {
		patterns = []*regexp.Regexp{
			regexp.MustCompile(`(?m)^\s*(implementation|compileOnly|runtimeOnly|testImplementation|testCompileOnly|testRuntimeOnly|annotationProcessor|kapt)\s*\(\s*"` + regexp.QuoteMeta(coords) + `[^"]*"\s*\)\s*\n?`),
			regexp.MustCompile(`(?m)^\s*(implementation|compileOnly|runtimeOnly|testImplementation|testCompileOnly|testRuntimeOnly|annotationProcessor|kapt)\s*\(\s*'` + regexp.QuoteMeta(coords) + `[^']*'\s*\)\s*\n?`),
		}
	} else {
		patterns = []*regexp.Regexp{
			regexp.MustCompile(`(?m)^\s*(implementation|compileOnly|runtimeOnly|testImplementation|testCompileOnly|testRuntimeOnly|annotationProcessor)\s+['"]` + regexp.QuoteMeta(coords) + `[^'"]*['"]\s*\n?`),
			regexp.MustCompile(`(?m)^\s*(implementation|compileOnly|runtimeOnly|testImplementation|testCompileOnly|testRuntimeOnly|annotationProcessor)\s*\(\s*['"]` + regexp.QuoteMeta(coords) + `[^'"]*['"]\s*\)\s*\n?`),
		}
	}

	for _, pattern := range patterns {
		content = pattern.ReplaceAllString(content, "")
	}

	return content
}

func (p *Parser) GetDependencies(project *buildtool.Project) []buildtool.Dependency {
	return project.Dependencies
}

func (p *Parser) GetDependency(project *buildtool.Project, groupId, artifactId string) *buildtool.Dependency {
	for i, dep := range project.Dependencies {
		if dep.GroupId == groupId && dep.ArtifactId == artifactId {
			return &project.Dependencies[i]
		}
	}
	return nil
}

func (p *Parser) GetJavaVersion(project *buildtool.Project) string {
	return project.JavaVersion
}

func (p *Parser) GetSpringBootVersion(project *buildtool.Project) string {
	return project.SpringBootVersion
}

func (p *Parser) GetBasePackage(project *buildtool.Project) string {
	artifactId := strings.ReplaceAll(project.ArtifactId, "-", "")
	return project.GroupId + "." + artifactId
}

func (p *Parser) HasLombok(project *buildtool.Project) bool {
	return p.HasDependency(project, "org.projectlombok", "lombok")
}

func (p *Parser) HasMapStruct(project *buildtool.Project) bool {
	return p.HasDependency(project, "org.mapstruct", "mapstruct")
}

func (p *Parser) HasSpringDataJpa(project *buildtool.Project) bool {
	return p.HasDependency(project, "org.springframework.boot", "spring-boot-starter-data-jpa")
}

func (p *Parser) HasSpringWeb(project *buildtool.Project) bool {
	return p.HasDependency(project, "org.springframework.boot", "spring-boot-starter-web")
}

func (p *Parser) HasValidation(project *buildtool.Project) bool {
	return p.HasDependency(project, "org.springframework.boot", "spring-boot-starter-validation")
}

func FindGradleBuildFile(startDir string, fs afero.Fs) (string, bool, error) {
	current := startDir
	for {
		ktsPath := filepath.Join(current, "build.gradle.kts")
		if _, err := fs.Stat(ktsPath); err == nil {
			return ktsPath, true, nil
		}

		groovyPath := filepath.Join(current, "build.gradle")
		if _, err := fs.Stat(groovyPath); err == nil {
			return groovyPath, false, nil
		}

		parent := filepath.Dir(current)
		if parent == current || parent == "" {
			break
		}
		current = parent
	}

	return "", false, fmt.Errorf("no gradle build file found in %s or any parent directory", startDir)
}

func FindGradleBuildFileWithCwd(fs afero.Fs) (string, bool, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", false, fmt.Errorf("could not get current directory: %w", err)
	}
	return FindGradleBuildFile(cwd, fs)
}
