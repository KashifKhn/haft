package buildtool

import "github.com/spf13/afero"

type Type string

const (
	Maven       Type = "maven"
	Gradle      Type = "gradle"
	GradleKotln Type = "gradle-kotlin"
)

type Dependency struct {
	GroupId    string
	ArtifactId string
	Version    string
	Scope      string
	Optional   bool
	Type       string
	Classifier string
}

type Project struct {
	GroupId           string
	ArtifactId        string
	Version           string
	Name              string
	Description       string
	JavaVersion       string
	SpringBootVersion string
	Packaging         string
	Dependencies      []Dependency
	BuildTool         Type
	Raw               any
}

type Parser interface {
	FindBuildFile(startDir string) (string, error)

	Parse(path string) (*Project, error)
	Write(path string, project *Project) error

	HasDependency(project *Project, groupId, artifactId string) bool
	AddDependency(project *Project, dep Dependency)
	RemoveDependency(project *Project, groupId, artifactId string) bool
	GetDependencies(project *Project) []Dependency
	GetDependency(project *Project, groupId, artifactId string) *Dependency

	GetJavaVersion(project *Project) string
	GetSpringBootVersion(project *Project) string
	GetBasePackage(project *Project) string

	HasLombok(project *Project) bool
	HasMapStruct(project *Project) bool
	HasSpringDataJpa(project *Project) bool
	HasSpringWeb(project *Project) bool
	HasValidation(project *Project) bool

	Type() Type
}

type ParserFactory func(fs afero.Fs) Parser

var registry = make(map[Type]ParserFactory)

func Register(t Type, factory ParserFactory) {
	registry[t] = factory
}

func GetParser(t Type, fs afero.Fs) Parser {
	if factory, ok := registry[t]; ok {
		return factory(fs)
	}
	return nil
}

func GetRegisteredTypes() []Type {
	var types []Type
	for t := range registry {
		types = append(types, t)
	}
	return types
}
