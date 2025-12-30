package maven

import (
	"encoding/xml"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/spf13/afero"
)

func init() {
	buildtool.Register(buildtool.Maven, func(fs afero.Fs) buildtool.Parser {
		return NewParserWithFs(fs)
	})
}

type Parser struct {
	fs afero.Fs
}

func NewParser() *Parser {
	return &Parser{fs: afero.NewOsFs()}
}

func NewParserWithFs(fs afero.Fs) *Parser {
	return &Parser{fs: fs}
}

func (p *Parser) Type() buildtool.Type {
	return buildtool.Maven
}

func (p *Parser) FindBuildFile(startDir string) (string, error) {
	current := startDir
	for {
		pomPath := current + "/pom.xml"
		if _, err := p.fs.Stat(pomPath); err == nil {
			return pomPath, nil
		}

		parent := current[:strings.LastIndex(current, string(os.PathSeparator))]
		if parent == current || parent == "" {
			break
		}
		current = parent
	}
	return "", fmt.Errorf("pom.xml not found in %s or any parent directory", startDir)
}

func (p *Parser) Parse(path string) (*buildtool.Project, error) {
	data, err := afero.ReadFile(p.fs, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read pom.xml: %w", err)
	}
	return p.ParseBytes(data)
}

func (p *Parser) ParseBytes(data []byte) (*buildtool.Project, error) {
	var mavenProject MavenProject
	if err := xml.Unmarshal(data, &mavenProject); err != nil {
		return nil, fmt.Errorf("failed to parse pom.xml: %w", err)
	}

	project := &buildtool.Project{
		GroupId:     mavenProject.GroupId,
		ArtifactId:  mavenProject.ArtifactId,
		Version:     mavenProject.Version,
		Name:        mavenProject.Name,
		Description: mavenProject.Description,
		Packaging:   mavenProject.Packaging,
		BuildTool:   buildtool.Maven,
		Raw:         &mavenProject,
	}

	if mavenProject.Properties != nil {
		project.JavaVersion = mavenProject.Properties.JavaVersion
	}

	if mavenProject.Parent != nil &&
		mavenProject.Parent.GroupId == "org.springframework.boot" &&
		mavenProject.Parent.ArtifactId == "spring-boot-starter-parent" {
		project.SpringBootVersion = mavenProject.Parent.Version
	}

	if mavenProject.Dependencies != nil {
		for _, dep := range mavenProject.Dependencies.Dependency {
			project.Dependencies = append(project.Dependencies, buildtool.Dependency{
				GroupId:    dep.GroupId,
				ArtifactId: dep.ArtifactId,
				Version:    dep.Version,
				Scope:      dep.Scope,
				Optional:   dep.Optional == "true",
				Type:       dep.Type,
				Classifier: dep.Classifier,
			})
		}
	}

	return project, nil
}

func (p *Parser) Write(path string, project *buildtool.Project) error {
	exists, err := afero.Exists(p.fs, path)
	if err != nil {
		return fmt.Errorf("failed to check file: %w", err)
	}

	if exists {
		return p.WriteMinimal(path, project)
	}

	return p.WriteFull(path, project)
}

func (p *Parser) WriteFull(path string, project *buildtool.Project) error {
	mavenProject := p.toMavenProject(project)
	data, err := p.Marshal(mavenProject)
	if err != nil {
		return err
	}
	return afero.WriteFile(p.fs, path, data, 0644)
}

func (p *Parser) WriteMinimal(path string, project *buildtool.Project) error {
	originalData, err := afero.ReadFile(p.fs, path)
	if err != nil {
		return fmt.Errorf("failed to read pom.xml: %w", err)
	}

	content := string(originalData)
	mavenProject := p.getMavenProject(project)

	var originalDeps *Dependencies
	var originalMaven MavenProject
	if err := xml.Unmarshal(originalData, &originalMaven); err == nil {
		originalDeps = originalMaven.Dependencies
	}

	var depsToAdd []Dependency
	if mavenProject != nil && mavenProject.Dependencies != nil {
		for _, dep := range mavenProject.Dependencies.Dependency {
			if !p.hasDependencyInList(originalDeps, dep.GroupId, dep.ArtifactId) {
				depsToAdd = append(depsToAdd, dep)
			}
		}
	}

	var depsToRemove []Dependency
	if originalDeps != nil {
		for _, dep := range originalDeps.Dependency {
			if mavenProject == nil || mavenProject.Dependencies == nil {
				depsToRemove = append(depsToRemove, dep)
			} else if !p.hasDependencyInList(mavenProject.Dependencies, dep.GroupId, dep.ArtifactId) {
				depsToRemove = append(depsToRemove, dep)
			}
		}
	}

	if len(depsToAdd) == 0 && len(depsToRemove) == 0 {
		return nil
	}

	for _, dep := range depsToRemove {
		content = p.removeDependencyFromContent(content, dep.GroupId, dep.ArtifactId)
	}

	if len(depsToAdd) > 0 {
		content = p.insertDependencies(content, depsToAdd)
	}

	return afero.WriteFile(p.fs, path, []byte(content), 0644)
}

func (p *Parser) hasDependencyInList(deps *Dependencies, groupId, artifactId string) bool {
	if deps == nil {
		return false
	}
	for _, dep := range deps.Dependency {
		if dep.GroupId == groupId && dep.ArtifactId == artifactId {
			return true
		}
	}
	return false
}

func (p *Parser) removeDependencyFromContent(content, groupId, artifactId string) string {
	pattern := fmt.Sprintf(`(?s)<dependency>\s*<groupId>%s</groupId>\s*<artifactId>%s</artifactId>.*?</dependency>`,
		regexp.QuoteMeta(groupId), regexp.QuoteMeta(artifactId))
	re := regexp.MustCompile(pattern)

	loc := re.FindStringIndex(content)
	if loc == nil {
		return content
	}

	start := loc[0]
	end := loc[1]

	leadingNewline := start
	for leadingNewline > 0 && (content[leadingNewline-1] == ' ' || content[leadingNewline-1] == '\t') {
		leadingNewline--
	}
	if leadingNewline > 0 && content[leadingNewline-1] == '\n' {
		leadingNewline--
	}

	trailingNewline := end
	for trailingNewline < len(content) && (content[trailingNewline] == ' ' || content[trailingNewline] == '\t') {
		trailingNewline++
	}
	if trailingNewline < len(content) && content[trailingNewline] == '\n' {
		trailingNewline++
	}

	return content[:leadingNewline] + content[trailingNewline:]
}

func (p *Parser) insertDependencies(content string, deps []Dependency) string {
	depsEndIdx := strings.LastIndex(content, "</dependencies>")
	if depsEndIdx == -1 {
		return p.insertDependenciesSection(content, deps)
	}

	indent := p.detectIndent(content, depsEndIdx)
	var depXML strings.Builder
	for _, dep := range deps {
		depXML.WriteString(p.formatDependency(dep, indent))
	}

	return content[:depsEndIdx] + depXML.String() + content[depsEndIdx:]
}

func (p *Parser) insertDependenciesSection(content string, deps []Dependency) string {
	insertPoints := []string{"</properties>", "</parent>", "</version>", "</artifactId>"}

	for _, point := range insertPoints {
		idx := strings.LastIndex(content, point)
		if idx != -1 {
			insertPos := idx + len(point)
			indent := p.detectBaseIndent(content)

			var section strings.Builder
			section.WriteString("\n\n")
			section.WriteString(indent)
			section.WriteString("<dependencies>")
			for _, dep := range deps {
				section.WriteString(p.formatDependency(dep, indent+indent))
			}
			section.WriteString("\n")
			section.WriteString(indent)
			section.WriteString("</dependencies>")

			return content[:insertPos] + section.String() + content[insertPos:]
		}
	}

	return content
}

func (p *Parser) detectIndent(content string, pos int) string {
	lineStart := strings.LastIndex(content[:pos], "\n")
	if lineStart == -1 {
		return "        "
	}

	line := content[lineStart+1 : pos]
	indent := ""
	for _, ch := range line {
		if ch == ' ' || ch == '\t' {
			indent += string(ch)
		} else {
			break
		}
	}

	if indent == "" {
		return "        "
	}
	return indent
}

func (p *Parser) detectBaseIndent(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")
		if strings.HasPrefix(trimmed, "<groupId>") || strings.HasPrefix(trimmed, "<artifactId>") {
			indent := line[:len(line)-len(trimmed)]
			if indent != "" {
				return indent
			}
		}
	}
	return "    "
}

func (p *Parser) formatDependency(dep Dependency, indent string) string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(indent)
	b.WriteString("<dependency>\n")
	b.WriteString(indent)
	b.WriteString("    <groupId>")
	b.WriteString(dep.GroupId)
	b.WriteString("</groupId>\n")
	b.WriteString(indent)
	b.WriteString("    <artifactId>")
	b.WriteString(dep.ArtifactId)
	b.WriteString("</artifactId>\n")

	if dep.Version != "" {
		b.WriteString(indent)
		b.WriteString("    <version>")
		b.WriteString(dep.Version)
		b.WriteString("</version>\n")
	}

	if dep.Scope != "" {
		b.WriteString(indent)
		b.WriteString("    <scope>")
		b.WriteString(dep.Scope)
		b.WriteString("</scope>\n")
	}

	if dep.Optional == "true" {
		b.WriteString(indent)
		b.WriteString("    <optional>true</optional>\n")
	}

	if dep.Type != "" {
		b.WriteString(indent)
		b.WriteString("    <type>")
		b.WriteString(dep.Type)
		b.WriteString("</type>\n")
	}

	if dep.Classifier != "" {
		b.WriteString(indent)
		b.WriteString("    <classifier>")
		b.WriteString(dep.Classifier)
		b.WriteString("</classifier>\n")
	}

	b.WriteString(indent)
	b.WriteString("</dependency>")

	return b.String()
}

func (p *Parser) Marshal(mavenProject *MavenProject) ([]byte, error) {
	mavenProject.Xmlns = "http://maven.apache.org/POM/4.0.0"
	mavenProject.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	mavenProject.SchemaLocation = "http://maven.apache.org/POM/4.0.0 https://maven.apache.org/xsd/maven-4.0.0.xsd"

	data, err := xml.MarshalIndent(mavenProject, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal pom.xml: %w", err)
	}

	return append([]byte(xml.Header), data...), nil
}

func (p *Parser) toMavenProject(project *buildtool.Project) *MavenProject {
	if raw, ok := project.Raw.(*MavenProject); ok {
		return raw
	}

	mavenProject := &MavenProject{
		ModelVersion: "4.0.0",
		GroupId:      project.GroupId,
		ArtifactId:   project.ArtifactId,
		Version:      project.Version,
		Name:         project.Name,
		Description:  project.Description,
		Packaging:    project.Packaging,
	}

	if project.JavaVersion != "" {
		mavenProject.Properties = &Properties{JavaVersion: project.JavaVersion}
	}

	if len(project.Dependencies) > 0 {
		mavenProject.Dependencies = &Dependencies{}
		for _, dep := range project.Dependencies {
			optional := ""
			if dep.Optional {
				optional = "true"
			}
			mavenProject.Dependencies.Dependency = append(mavenProject.Dependencies.Dependency, Dependency{
				GroupId:    dep.GroupId,
				ArtifactId: dep.ArtifactId,
				Version:    dep.Version,
				Scope:      dep.Scope,
				Optional:   optional,
				Type:       dep.Type,
				Classifier: dep.Classifier,
			})
		}
	}

	return mavenProject
}

func (p *Parser) getMavenProject(project *buildtool.Project) *MavenProject {
	if raw, ok := project.Raw.(*MavenProject); ok {
		return raw
	}
	return p.toMavenProject(project)
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

	if mavenProject := p.getMavenProject(project); mavenProject != nil {
		if mavenProject.Dependencies == nil {
			mavenProject.Dependencies = &Dependencies{}
		}
		optional := ""
		if dep.Optional {
			optional = "true"
		}
		mavenProject.Dependencies.Dependency = append(mavenProject.Dependencies.Dependency, Dependency{
			GroupId:    dep.GroupId,
			ArtifactId: dep.ArtifactId,
			Version:    dep.Version,
			Scope:      dep.Scope,
			Optional:   optional,
			Type:       dep.Type,
			Classifier: dep.Classifier,
		})
	}
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

	if mavenProject := p.getMavenProject(project); mavenProject != nil && mavenProject.Dependencies != nil {
		var newMavenDeps []Dependency
		for _, dep := range mavenProject.Dependencies.Dependency {
			if dep.GroupId == groupId && dep.ArtifactId == artifactId {
				continue
			}
			newMavenDeps = append(newMavenDeps, dep)
		}
		mavenProject.Dependencies.Dependency = newMavenDeps
	}

	return found
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

func (p *Parser) FindPomXml(startDir string) (string, error) {
	return p.FindBuildFile(startDir)
}

func (p *Parser) ParseLegacy(path string) (*MavenProject, error) {
	data, err := afero.ReadFile(p.fs, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read pom.xml: %w", err)
	}
	var mavenProject MavenProject
	if err := xml.Unmarshal(data, &mavenProject); err != nil {
		return nil, fmt.Errorf("failed to parse pom.xml: %w", err)
	}
	return &mavenProject, nil
}

func (p *Parser) WriteLegacy(path string, mavenProject *MavenProject) error {
	data, err := p.Marshal(mavenProject)
	if err != nil {
		return err
	}
	return afero.WriteFile(p.fs, path, data, 0644)
}

func (p *Parser) HasDependencyLegacy(mavenProject *MavenProject, groupId, artifactId string) bool {
	if mavenProject.Dependencies == nil {
		return false
	}
	for _, dep := range mavenProject.Dependencies.Dependency {
		if dep.GroupId == groupId && dep.ArtifactId == artifactId {
			return true
		}
	}
	return false
}

func (p *Parser) AddDependencyLegacy(mavenProject *MavenProject, dep Dependency) {
	if mavenProject.Dependencies == nil {
		mavenProject.Dependencies = &Dependencies{}
	}
	if !p.HasDependencyLegacy(mavenProject, dep.GroupId, dep.ArtifactId) {
		mavenProject.Dependencies.Dependency = append(mavenProject.Dependencies.Dependency, dep)
	}
}

func (p *Parser) RemoveDependencyLegacy(mavenProject *MavenProject, groupId, artifactId string) bool {
	if mavenProject.Dependencies == nil {
		return false
	}
	for i, dep := range mavenProject.Dependencies.Dependency {
		if dep.GroupId == groupId && dep.ArtifactId == artifactId {
			mavenProject.Dependencies.Dependency = append(
				mavenProject.Dependencies.Dependency[:i],
				mavenProject.Dependencies.Dependency[i+1:]...,
			)
			return true
		}
	}
	return false
}

func (p *Parser) GetDependencyLegacy(mavenProject *MavenProject, groupId, artifactId string) *Dependency {
	if mavenProject.Dependencies == nil {
		return nil
	}
	for i, dep := range mavenProject.Dependencies.Dependency {
		if dep.GroupId == groupId && dep.ArtifactId == artifactId {
			return &mavenProject.Dependencies.Dependency[i]
		}
	}
	return nil
}

func (p *Parser) HasLombokLegacy(mavenProject *MavenProject) bool {
	return p.HasDependencyLegacy(mavenProject, "org.projectlombok", "lombok")
}

func (p *Parser) HasMapStructLegacy(mavenProject *MavenProject) bool {
	return p.HasDependencyLegacy(mavenProject, "org.mapstruct", "mapstruct")
}

func (p *Parser) HasSpringDataJpaLegacy(mavenProject *MavenProject) bool {
	return p.HasDependencyLegacy(mavenProject, "org.springframework.boot", "spring-boot-starter-data-jpa")
}

func (p *Parser) HasSpringWebLegacy(mavenProject *MavenProject) bool {
	return p.HasDependencyLegacy(mavenProject, "org.springframework.boot", "spring-boot-starter-web")
}

func (p *Parser) HasValidationLegacy(mavenProject *MavenProject) bool {
	return p.HasDependencyLegacy(mavenProject, "org.springframework.boot", "spring-boot-starter-validation")
}

func (p *Parser) GetJavaVersionLegacy(mavenProject *MavenProject) string {
	if mavenProject.Properties != nil && mavenProject.Properties.JavaVersion != "" {
		return mavenProject.Properties.JavaVersion
	}
	return ""
}

func (p *Parser) GetSpringBootVersionLegacy(mavenProject *MavenProject) string {
	if mavenProject.Parent != nil &&
		mavenProject.Parent.GroupId == "org.springframework.boot" &&
		mavenProject.Parent.ArtifactId == "spring-boot-starter-parent" {
		return mavenProject.Parent.Version
	}
	return ""
}

func (p *Parser) GetBasePackageLegacy(mavenProject *MavenProject) string {
	artifactId := strings.ReplaceAll(mavenProject.ArtifactId, "-", "")
	return mavenProject.GroupId + "." + artifactId
}
