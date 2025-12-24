package initializr

import (
	_ "embed"
	"encoding/json"
)

//go:embed metadata.json
var metadataJSON []byte

var cachedMetadata *Metadata

func LoadMetadata() (*Metadata, error) {
	if cachedMetadata != nil {
		return cachedMetadata, nil
	}

	var m Metadata
	if err := json.Unmarshal(metadataJSON, &m); err != nil {
		return nil, err
	}
	cachedMetadata = &m
	return &m, nil
}

func GetDependencyCategories() ([]DependencyCategory, error) {
	m, err := LoadMetadata()
	if err != nil {
		return nil, err
	}
	return m.Dependencies.Values, nil
}

func GetJavaVersions() ([]Option, error) {
	m, err := LoadMetadata()
	if err != nil {
		return nil, err
	}
	return m.JavaVersion.Values, nil
}

func GetBootVersions() ([]Option, error) {
	m, err := LoadMetadata()
	if err != nil {
		return nil, err
	}
	return m.BootVersion.Values, nil
}

func GetPackagingOptions() ([]Option, error) {
	m, err := LoadMetadata()
	if err != nil {
		return nil, err
	}
	return m.Packaging.Values, nil
}

func GetBuildTypes() ([]Option, error) {
	m, err := LoadMetadata()
	if err != nil {
		return nil, err
	}
	return m.Type.Values, nil
}

func GetDefaults() (groupId, artifactId, name, description, packageName, version string, err error) {
	m, err := LoadMetadata()
	if err != nil {
		return "", "", "", "", "", "", err
	}
	return m.GroupId.Default, m.ArtifactId.Default, m.Name.Default,
		m.Description.Default, m.PackageName.Default, m.Version.Default, nil
}
