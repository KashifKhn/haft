package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

const (
	ProjectTemplateDir = ".haft/templates"
	GlobalTemplateDir  = ".haft/templates"
)

type TemplateSource int

const (
	SourceEmbedded TemplateSource = iota
	SourceProject
	SourceGlobal
)

func (s TemplateSource) String() string {
	switch s {
	case SourceProject:
		return "project"
	case SourceGlobal:
		return "global"
	default:
		return "embedded"
	}
}

type TemplateLoader struct {
	fs          afero.Fs
	projectRoot string
	homeDir     string
}

type LoadedTemplate struct {
	Content []byte
	Source  TemplateSource
	Path    string
}

func NewTemplateLoader(fs afero.Fs, projectRoot string) *TemplateLoader {
	homeDir, _ := os.UserHomeDir()
	return &TemplateLoader{
		fs:          fs,
		projectRoot: projectRoot,
		homeDir:     homeDir,
	}
}

func NewTemplateLoaderWithHome(fs afero.Fs, projectRoot, homeDir string) *TemplateLoader {
	return &TemplateLoader{
		fs:          fs,
		projectRoot: projectRoot,
		homeDir:     homeDir,
	}
}

func (l *TemplateLoader) LoadTemplate(name string) (*LoadedTemplate, error) {
	if content, path, err := l.loadFromProject(name); err == nil {
		return &LoadedTemplate{Content: content, Source: SourceProject, Path: path}, nil
	}

	if content, path, err := l.loadFromGlobal(name); err == nil {
		return &LoadedTemplate{Content: content, Source: SourceGlobal, Path: path}, nil
	}

	content, err := l.loadFromEmbedded(name)
	if err != nil {
		return nil, fmt.Errorf("template not found: %s", name)
	}

	return &LoadedTemplate{Content: content, Source: SourceEmbedded, Path: "embedded:" + name}, nil
}

func (l *TemplateLoader) loadFromProject(name string) ([]byte, string, error) {
	if l.projectRoot == "" {
		return nil, "", fmt.Errorf("project root not set")
	}

	path := filepath.Join(l.projectRoot, ProjectTemplateDir, name)
	content, err := afero.ReadFile(l.fs, path)
	if err != nil {
		return nil, "", err
	}

	return content, path, nil
}

func (l *TemplateLoader) loadFromGlobal(name string) ([]byte, string, error) {
	if l.homeDir == "" {
		return nil, "", fmt.Errorf("home directory not set")
	}

	path := filepath.Join(l.homeDir, GlobalTemplateDir, name)
	content, err := afero.ReadFile(l.fs, path)
	if err != nil {
		return nil, "", err
	}

	return content, path, nil
}

func (l *TemplateLoader) loadFromEmbedded(name string) ([]byte, error) {
	return templateFS.ReadFile("templates/" + name)
}

func (l *TemplateLoader) TemplateExists(name string) bool {
	_, err := l.LoadTemplate(name)
	return err == nil
}

func (l *TemplateLoader) GetTemplateSource(name string) (TemplateSource, error) {
	loaded, err := l.LoadTemplate(name)
	if err != nil {
		return SourceEmbedded, err
	}
	return loaded.Source, nil
}

func (l *TemplateLoader) ListProjectTemplates() ([]string, error) {
	if l.projectRoot == "" {
		return nil, nil
	}

	templateDir := filepath.Join(l.projectRoot, ProjectTemplateDir)
	return l.listTemplatesInDir(templateDir)
}

func (l *TemplateLoader) ListGlobalTemplates() ([]string, error) {
	if l.homeDir == "" {
		return nil, nil
	}

	templateDir := filepath.Join(l.homeDir, GlobalTemplateDir)
	return l.listTemplatesInDir(templateDir)
}

func (l *TemplateLoader) listTemplatesInDir(dir string) ([]string, error) {
	exists, err := afero.DirExists(l.fs, dir)
	if err != nil || !exists {
		return nil, nil
	}

	var templates []string
	err = afero.Walk(l.fs, dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".tmpl" {
			relPath, _ := filepath.Rel(dir, path)
			templates = append(templates, relPath)
		}
		return nil
	})

	return templates, err
}

func (l *TemplateLoader) GetProjectTemplateDir() string {
	if l.projectRoot == "" {
		return ""
	}
	return filepath.Join(l.projectRoot, ProjectTemplateDir)
}

func (l *TemplateLoader) GetGlobalTemplateDir() string {
	if l.homeDir == "" {
		return ""
	}
	return filepath.Join(l.homeDir, GlobalTemplateDir)
}

func (l *TemplateLoader) ProjectTemplatesExist() bool {
	if l.projectRoot == "" {
		return false
	}
	exists, _ := afero.DirExists(l.fs, filepath.Join(l.projectRoot, ProjectTemplateDir))
	return exists
}

func (l *TemplateLoader) GlobalTemplatesExist() bool {
	if l.homeDir == "" {
		return false
	}
	exists, _ := afero.DirExists(l.fs, filepath.Join(l.homeDir, GlobalTemplateDir))
	return exists
}

func (l *TemplateLoader) CopyEmbeddedToProject(templateName string) error {
	content, err := l.loadFromEmbedded(templateName)
	if err != nil {
		return fmt.Errorf("embedded template not found: %s", templateName)
	}

	destPath := filepath.Join(l.projectRoot, ProjectTemplateDir, templateName)
	destDir := filepath.Dir(destPath)

	if err := l.fs.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return afero.WriteFile(l.fs, destPath, content, 0644)
}

func (l *TemplateLoader) CopyAllEmbeddedToProject(templateDir string) error {
	templates, err := ListEmbeddedTemplates(templateDir)
	if err != nil {
		return err
	}

	for _, tmpl := range templates {
		if err := l.CopyEmbeddedToProject(tmpl); err != nil {
			return err
		}
	}

	return nil
}

func ListEmbeddedTemplates(dir string) ([]string, error) {
	var templates []string

	entries, err := templateFS.ReadDir("templates/" + dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		name := entry.Name()
		path := dir + "/" + name

		if entry.IsDir() {
			subTemplates, err := ListEmbeddedTemplates(path)
			if err != nil {
				return nil, err
			}
			templates = append(templates, subTemplates...)
		} else if filepath.Ext(name) == ".tmpl" {
			templates = append(templates, path)
		}
	}

	return templates, nil
}

func ListAllEmbeddedTemplates() ([]string, error) {
	var allTemplates []string

	dirs := []string{"resource", "test", "project"}
	for _, dir := range dirs {
		templates, err := ListEmbeddedTemplates(dir)
		if err != nil {
			continue
		}
		allTemplates = append(allTemplates, templates...)
	}

	return allTemplates, nil
}

type TemplateInfo struct {
	Name   string
	Source TemplateSource
	Path   string
}

func (l *TemplateLoader) GetTemplateInfo(name string) (*TemplateInfo, error) {
	loaded, err := l.LoadTemplate(name)
	if err != nil {
		return nil, err
	}

	return &TemplateInfo{
		Name:   name,
		Source: loaded.Source,
		Path:   loaded.Path,
	}, nil
}

func (l *TemplateLoader) ListAllTemplatesWithSource() ([]TemplateInfo, error) {
	embeddedTemplates, err := ListAllEmbeddedTemplates()
	if err != nil {
		return nil, err
	}

	var result []TemplateInfo
	seen := make(map[string]bool)

	projectTemplates, _ := l.ListProjectTemplates()
	for _, tmpl := range projectTemplates {
		result = append(result, TemplateInfo{
			Name:   tmpl,
			Source: SourceProject,
			Path:   filepath.Join(l.GetProjectTemplateDir(), tmpl),
		})
		seen[tmpl] = true
	}

	globalTemplates, _ := l.ListGlobalTemplates()
	for _, tmpl := range globalTemplates {
		if !seen[tmpl] {
			result = append(result, TemplateInfo{
				Name:   tmpl,
				Source: SourceGlobal,
				Path:   filepath.Join(l.GetGlobalTemplateDir(), tmpl),
			})
			seen[tmpl] = true
		}
	}

	for _, tmpl := range embeddedTemplates {
		if !seen[tmpl] {
			result = append(result, TemplateInfo{
				Name:   tmpl,
				Source: SourceEmbedded,
				Path:   "embedded:" + tmpl,
			})
		}
	}

	return result, nil
}
