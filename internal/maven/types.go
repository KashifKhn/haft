package maven

import "encoding/xml"

type MavenProject struct {
	XMLName              xml.Name              `xml:"project"`
	Xmlns                string                `xml:"xmlns,attr"`
	XmlnsXsi             string                `xml:"xmlns:xsi,attr"`
	SchemaLocation       string                `xml:"xsi:schemaLocation,attr"`
	ModelVersion         string                `xml:"modelVersion"`
	Parent               *Parent               `xml:"parent,omitempty"`
	GroupId              string                `xml:"groupId"`
	ArtifactId           string                `xml:"artifactId"`
	Version              string                `xml:"version"`
	Name                 string                `xml:"name,omitempty"`
	Description          string                `xml:"description,omitempty"`
	Packaging            string                `xml:"packaging,omitempty"`
	Properties           *Properties           `xml:"properties,omitempty"`
	Dependencies         *Dependencies         `xml:"dependencies,omitempty"`
	DependencyManagement *DependencyManagement `xml:"dependencyManagement,omitempty"`
	Build                *Build                `xml:"build,omitempty"`
}

type Parent struct {
	GroupId      string `xml:"groupId"`
	ArtifactId   string `xml:"artifactId"`
	Version      string `xml:"version"`
	RelativePath string `xml:"relativePath,omitempty"`
}

type Properties struct {
	JavaVersion    string          `xml:"java.version,omitempty"`
	SourceEncoding string          `xml:"project.build.sourceEncoding,omitempty"`
	Entries        []PropertyEntry `xml:"-"`
	Raw            []byte          `xml:"-"`
}

type PropertyEntry struct {
	Key   string
	Value string
}

type Dependencies struct {
	Dependency []Dependency `xml:"dependency"`
}

type Dependency struct {
	GroupId    string      `xml:"groupId"`
	ArtifactId string      `xml:"artifactId"`
	Version    string      `xml:"version,omitempty"`
	Scope      string      `xml:"scope,omitempty"`
	Optional   string      `xml:"optional,omitempty"`
	Type       string      `xml:"type,omitempty"`
	Classifier string      `xml:"classifier,omitempty"`
	Exclusions *Exclusions `xml:"exclusions,omitempty"`
}

type Exclusions struct {
	Exclusion []Exclusion `xml:"exclusion"`
}

type Exclusion struct {
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
}

type DependencyManagement struct {
	Dependencies *Dependencies `xml:"dependencies,omitempty"`
}

type Build struct {
	Plugins          *Plugins          `xml:"plugins,omitempty"`
	PluginManagement *PluginManagement `xml:"pluginManagement,omitempty"`
}

type Plugins struct {
	Plugin []Plugin `xml:"plugin"`
}

type PluginManagement struct {
	Plugins *Plugins `xml:"plugins,omitempty"`
}

type Plugin struct {
	GroupId       string         `xml:"groupId,omitempty"`
	ArtifactId    string         `xml:"artifactId"`
	Version       string         `xml:"version,omitempty"`
	Configuration *Configuration `xml:"configuration,omitempty"`
	Executions    *Executions    `xml:"executions,omitempty"`
}

type Configuration struct {
	Raw []byte `xml:",innerxml"`
}

type Executions struct {
	Execution []Execution `xml:"execution"`
}

type Execution struct {
	ID    string `xml:"id,omitempty"`
	Phase string `xml:"phase,omitempty"`
	Goals *Goals `xml:"goals,omitempty"`
}

type Goals struct {
	Goal []string `xml:"goal"`
}

type Project = MavenProject
