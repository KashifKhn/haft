package initializr

type Metadata struct {
	Dependencies DependencyGroup `json:"dependencies"`
	Type         OptionGroup     `json:"type"`
	Packaging    OptionGroup     `json:"packaging"`
	JavaVersion  OptionGroup     `json:"javaVersion"`
	Language     OptionGroup     `json:"language"`
	BootVersion  OptionGroup     `json:"bootVersion"`
	GroupId      TextOption      `json:"groupId"`
	ArtifactId   TextOption      `json:"artifactId"`
	Version      TextOption      `json:"version"`
	Name         TextOption      `json:"name"`
	Description  TextOption      `json:"description"`
	PackageName  TextOption      `json:"packageName"`
}

type DependencyGroup struct {
	Type   string               `json:"type"`
	Values []DependencyCategory `json:"values"`
}

type DependencyCategory struct {
	Name   string       `json:"name"`
	Values []Dependency `json:"values"`
}

type Dependency struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	VersionRange string `json:"versionRange,omitempty"`
}

type OptionGroup struct {
	Type    string   `json:"type"`
	Default string   `json:"default"`
	Values  []Option `json:"values"`
}

type Option struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type TextOption struct {
	Type    string `json:"type"`
	Default string `json:"default"`
}
