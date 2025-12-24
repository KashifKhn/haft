package maven

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/afero"
)

type Parser struct {
	fs afero.Fs
}

func NewParser() *Parser {
	return &Parser{fs: afero.NewOsFs()}
}

func NewParserWithFs(fs afero.Fs) *Parser {
	return &Parser{fs: fs}
}

func (p *Parser) Parse(path string) (*Project, error) {
	data, err := afero.ReadFile(p.fs, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read pom.xml: %w", err)
	}
	return p.ParseBytes(data)
}

func (p *Parser) ParseBytes(data []byte) (*Project, error) {
	var project Project
	if err := xml.Unmarshal(data, &project); err != nil {
		return nil, fmt.Errorf("failed to parse pom.xml: %w", err)
	}
	return &project, nil
}

func (p *Parser) Write(path string, project *Project) error {
	data, err := p.Marshal(project)
	if err != nil {
		return err
	}
	return afero.WriteFile(p.fs, path, data, 0644)
}

func (p *Parser) Marshal(project *Project) ([]byte, error) {
	project.Xmlns = "http://maven.apache.org/POM/4.0.0"
	project.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	project.SchemaLocation = "http://maven.apache.org/POM/4.0.0 https://maven.apache.org/xsd/maven-4.0.0.xsd"

	data, err := xml.MarshalIndent(project, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal pom.xml: %w", err)
	}

	return append([]byte(xml.Header), data...), nil
}

func (p *Parser) HasDependency(project *Project, groupId, artifactId string) bool {
	if project.Dependencies == nil {
		return false
	}
	for _, dep := range project.Dependencies.Dependency {
		if dep.GroupId == groupId && dep.ArtifactId == artifactId {
			return true
		}
	}
	return false
}

func (p *Parser) AddDependency(project *Project, dep Dependency) {
	if project.Dependencies == nil {
		project.Dependencies = &Dependencies{}
	}
	if !p.HasDependency(project, dep.GroupId, dep.ArtifactId) {
		project.Dependencies.Dependency = append(project.Dependencies.Dependency, dep)
	}
}

func (p *Parser) RemoveDependency(project *Project, groupId, artifactId string) bool {
	if project.Dependencies == nil {
		return false
	}
	for i, dep := range project.Dependencies.Dependency {
		if dep.GroupId == groupId && dep.ArtifactId == artifactId {
			project.Dependencies.Dependency = append(
				project.Dependencies.Dependency[:i],
				project.Dependencies.Dependency[i+1:]...,
			)
			return true
		}
	}
	return false
}

func (p *Parser) GetDependency(project *Project, groupId, artifactId string) *Dependency {
	if project.Dependencies == nil {
		return nil
	}
	for i, dep := range project.Dependencies.Dependency {
		if dep.GroupId == groupId && dep.ArtifactId == artifactId {
			return &project.Dependencies.Dependency[i]
		}
	}
	return nil
}

func (p *Parser) HasLombok(project *Project) bool {
	return p.HasDependency(project, "org.projectlombok", "lombok")
}

func (p *Parser) HasMapStruct(project *Project) bool {
	return p.HasDependency(project, "org.mapstruct", "mapstruct")
}

func (p *Parser) HasSpringDataJpa(project *Project) bool {
	return p.HasDependency(project, "org.springframework.boot", "spring-boot-starter-data-jpa")
}

func (p *Parser) HasSpringWeb(project *Project) bool {
	return p.HasDependency(project, "org.springframework.boot", "spring-boot-starter-web")
}

func (p *Parser) HasValidation(project *Project) bool {
	return p.HasDependency(project, "org.springframework.boot", "spring-boot-starter-validation")
}

func (p *Parser) GetJavaVersion(project *Project) string {
	if project.Properties != nil && project.Properties.JavaVersion != "" {
		return project.Properties.JavaVersion
	}
	return ""
}

func (p *Parser) GetSpringBootVersion(project *Project) string {
	if project.Parent != nil &&
		project.Parent.GroupId == "org.springframework.boot" &&
		project.Parent.ArtifactId == "spring-boot-starter-parent" {
		return project.Parent.Version
	}
	return ""
}

func (p *Parser) GetBasePackage(project *Project) string {
	groupId := project.GroupId
	artifactId := strings.ReplaceAll(project.ArtifactId, "-", "")
	return groupId + "." + artifactId
}

func (p *Parser) FindPomXml(startDir string) (string, error) {
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
