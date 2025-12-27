package generator

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTemplateLoader(t *testing.T) {
	fs := afero.NewMemMapFs()
	loader := NewTemplateLoader(fs, "/project")

	assert.NotNil(t, loader)
	assert.Equal(t, "/project", loader.projectRoot)
}

func TestNewTemplateLoaderWithHome(t *testing.T) {
	fs := afero.NewMemMapFs()
	loader := NewTemplateLoaderWithHome(fs, "/project", "/home/user")

	assert.NotNil(t, loader)
	assert.Equal(t, "/project", loader.projectRoot)
	assert.Equal(t, "/home/user", loader.homeDir)
}

func TestTemplateLoader_LoadFromProject(t *testing.T) {
	fs := afero.NewMemMapFs()
	projectRoot := "/project"
	templatePath := filepath.Join(projectRoot, ProjectTemplateDir, "resource/layered/Controller.java.tmpl")

	require.NoError(t, fs.MkdirAll(filepath.Dir(templatePath), 0755))
	require.NoError(t, afero.WriteFile(fs, templatePath, []byte("project template content"), 0644))

	loader := NewTemplateLoaderWithHome(fs, projectRoot, "/home/user")

	loaded, err := loader.LoadTemplate("resource/layered/Controller.java.tmpl")
	require.NoError(t, err)

	assert.Equal(t, SourceProject, loaded.Source)
	assert.Equal(t, "project template content", string(loaded.Content))
	assert.Contains(t, loaded.Path, ProjectTemplateDir)
}

func TestTemplateLoader_LoadFromGlobal(t *testing.T) {
	fs := afero.NewMemMapFs()
	projectRoot := "/project"
	homeDir := "/home/user"
	templatePath := filepath.Join(homeDir, GlobalTemplateDir, "resource/layered/Controller.java.tmpl")

	require.NoError(t, fs.MkdirAll(filepath.Dir(templatePath), 0755))
	require.NoError(t, afero.WriteFile(fs, templatePath, []byte("global template content"), 0644))

	loader := NewTemplateLoaderWithHome(fs, projectRoot, homeDir)

	loaded, err := loader.LoadTemplate("resource/layered/Controller.java.tmpl")
	require.NoError(t, err)

	assert.Equal(t, SourceGlobal, loaded.Source)
	assert.Equal(t, "global template content", string(loaded.Content))
}

func TestTemplateLoader_LoadFromEmbedded(t *testing.T) {
	fs := afero.NewMemMapFs()
	loader := NewTemplateLoaderWithHome(fs, "/project", "/home/user")

	loaded, err := loader.LoadTemplate("resource/layered/Controller.java.tmpl")
	require.NoError(t, err)

	assert.Equal(t, SourceEmbedded, loaded.Source)
	assert.NotEmpty(t, loaded.Content)
	assert.Equal(t, "embedded:resource/layered/Controller.java.tmpl", loaded.Path)
}

func TestTemplateLoader_PriorityOrder(t *testing.T) {
	fs := afero.NewMemMapFs()
	projectRoot := "/project"
	homeDir := "/home/user"
	templateName := "resource/layered/Controller.java.tmpl"

	projectPath := filepath.Join(projectRoot, ProjectTemplateDir, templateName)
	require.NoError(t, fs.MkdirAll(filepath.Dir(projectPath), 0755))
	require.NoError(t, afero.WriteFile(fs, projectPath, []byte("project template"), 0644))

	globalPath := filepath.Join(homeDir, GlobalTemplateDir, templateName)
	require.NoError(t, fs.MkdirAll(filepath.Dir(globalPath), 0755))
	require.NoError(t, afero.WriteFile(fs, globalPath, []byte("global template"), 0644))

	loader := NewTemplateLoaderWithHome(fs, projectRoot, homeDir)

	loaded, err := loader.LoadTemplate(templateName)
	require.NoError(t, err)

	assert.Equal(t, SourceProject, loaded.Source)
	assert.Equal(t, "project template", string(loaded.Content))
}

func TestTemplateLoader_FallbackToGlobalWhenNoProject(t *testing.T) {
	fs := afero.NewMemMapFs()
	projectRoot := "/project"
	homeDir := "/home/user"
	templateName := "resource/layered/Controller.java.tmpl"

	globalPath := filepath.Join(homeDir, GlobalTemplateDir, templateName)
	require.NoError(t, fs.MkdirAll(filepath.Dir(globalPath), 0755))
	require.NoError(t, afero.WriteFile(fs, globalPath, []byte("global template"), 0644))

	loader := NewTemplateLoaderWithHome(fs, projectRoot, homeDir)

	loaded, err := loader.LoadTemplate(templateName)
	require.NoError(t, err)

	assert.Equal(t, SourceGlobal, loaded.Source)
	assert.Equal(t, "global template", string(loaded.Content))
}

func TestTemplateLoader_NotFound(t *testing.T) {
	fs := afero.NewMemMapFs()
	loader := NewTemplateLoaderWithHome(fs, "/project", "/home/user")

	_, err := loader.LoadTemplate("nonexistent/template.java.tmpl")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "template not found")
}

func TestTemplateLoader_TemplateExists(t *testing.T) {
	fs := afero.NewMemMapFs()
	loader := NewTemplateLoaderWithHome(fs, "/project", "/home/user")

	assert.True(t, loader.TemplateExists("resource/layered/Controller.java.tmpl"))
	assert.False(t, loader.TemplateExists("nonexistent/template.java.tmpl"))
}

func TestTemplateLoader_GetTemplateSource(t *testing.T) {
	fs := afero.NewMemMapFs()
	projectRoot := "/project"
	templatePath := filepath.Join(projectRoot, ProjectTemplateDir, "resource/layered/Controller.java.tmpl")

	require.NoError(t, fs.MkdirAll(filepath.Dir(templatePath), 0755))
	require.NoError(t, afero.WriteFile(fs, templatePath, []byte("content"), 0644))

	loader := NewTemplateLoaderWithHome(fs, projectRoot, "/home/user")

	source, err := loader.GetTemplateSource("resource/layered/Controller.java.tmpl")
	require.NoError(t, err)
	assert.Equal(t, SourceProject, source)

	source, err = loader.GetTemplateSource("resource/layered/Service.java.tmpl")
	require.NoError(t, err)
	assert.Equal(t, SourceEmbedded, source)
}

func TestTemplateLoader_ListProjectTemplates(t *testing.T) {
	fs := afero.NewMemMapFs()
	projectRoot := "/project"
	templateDir := filepath.Join(projectRoot, ProjectTemplateDir)

	templates := []string{
		"resource/layered/Controller.java.tmpl",
		"resource/layered/Service.java.tmpl",
	}

	for _, tmpl := range templates {
		path := filepath.Join(templateDir, tmpl)
		require.NoError(t, fs.MkdirAll(filepath.Dir(path), 0755))
		require.NoError(t, afero.WriteFile(fs, path, []byte("content"), 0644))
	}

	loader := NewTemplateLoaderWithHome(fs, projectRoot, "/home/user")

	listed, err := loader.ListProjectTemplates()
	require.NoError(t, err)

	assert.Len(t, listed, 2)
	for _, tmpl := range templates {
		assert.Contains(t, listed, tmpl)
	}
}

func TestTemplateLoader_ListGlobalTemplates(t *testing.T) {
	fs := afero.NewMemMapFs()
	homeDir := "/home/user"
	templateDir := filepath.Join(homeDir, GlobalTemplateDir)

	templates := []string{
		"resource/layered/Controller.java.tmpl",
		"resource/layered/Service.java.tmpl",
	}

	for _, tmpl := range templates {
		path := filepath.Join(templateDir, tmpl)
		require.NoError(t, fs.MkdirAll(filepath.Dir(path), 0755))
		require.NoError(t, afero.WriteFile(fs, path, []byte("content"), 0644))
	}

	loader := NewTemplateLoaderWithHome(fs, "/project", homeDir)

	listed, err := loader.ListGlobalTemplates()
	require.NoError(t, err)

	assert.Len(t, listed, 2)
	for _, tmpl := range templates {
		assert.Contains(t, listed, tmpl)
	}
}

func TestTemplateLoader_GetProjectTemplateDir(t *testing.T) {
	fs := afero.NewMemMapFs()
	loader := NewTemplateLoaderWithHome(fs, "/project", "/home/user")

	assert.Equal(t, "/project/.haft/templates", loader.GetProjectTemplateDir())
}

func TestTemplateLoader_GetGlobalTemplateDir(t *testing.T) {
	fs := afero.NewMemMapFs()
	loader := NewTemplateLoaderWithHome(fs, "/project", "/home/user")

	assert.Equal(t, "/home/user/.haft/templates", loader.GetGlobalTemplateDir())
}

func TestTemplateLoader_ProjectTemplatesExist(t *testing.T) {
	fs := afero.NewMemMapFs()
	projectRoot := "/project"

	loader := NewTemplateLoaderWithHome(fs, projectRoot, "/home/user")
	assert.False(t, loader.ProjectTemplatesExist())

	require.NoError(t, fs.MkdirAll(filepath.Join(projectRoot, ProjectTemplateDir), 0755))
	assert.True(t, loader.ProjectTemplatesExist())
}

func TestTemplateLoader_GlobalTemplatesExist(t *testing.T) {
	fs := afero.NewMemMapFs()
	homeDir := "/home/user"

	loader := NewTemplateLoaderWithHome(fs, "/project", homeDir)
	assert.False(t, loader.GlobalTemplatesExist())

	require.NoError(t, fs.MkdirAll(filepath.Join(homeDir, GlobalTemplateDir), 0755))
	assert.True(t, loader.GlobalTemplatesExist())
}

func TestTemplateLoader_CopyEmbeddedToProject(t *testing.T) {
	fs := afero.NewMemMapFs()
	projectRoot := "/project"
	loader := NewTemplateLoaderWithHome(fs, projectRoot, "/home/user")

	err := loader.CopyEmbeddedToProject("resource/layered/Controller.java.tmpl")
	require.NoError(t, err)

	destPath := filepath.Join(projectRoot, ProjectTemplateDir, "resource/layered/Controller.java.tmpl")
	exists, err := afero.Exists(fs, destPath)
	require.NoError(t, err)
	assert.True(t, exists)

	content, err := afero.ReadFile(fs, destPath)
	require.NoError(t, err)
	assert.NotEmpty(t, content)
}

func TestTemplateLoader_CopyEmbeddedToProject_NonExistent(t *testing.T) {
	fs := afero.NewMemMapFs()
	loader := NewTemplateLoaderWithHome(fs, "/project", "/home/user")

	err := loader.CopyEmbeddedToProject("nonexistent/template.java.tmpl")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "embedded template not found")
}

func TestListEmbeddedTemplates(t *testing.T) {
	templates, err := ListEmbeddedTemplates("resource/layered")
	require.NoError(t, err)

	assert.NotEmpty(t, templates)
	assert.Contains(t, templates, "resource/layered/Controller.java.tmpl")
	assert.Contains(t, templates, "resource/layered/Service.java.tmpl")
}

func TestListAllEmbeddedTemplates(t *testing.T) {
	templates, err := ListAllEmbeddedTemplates()
	require.NoError(t, err)

	assert.NotEmpty(t, templates)

	hasResourceTemplate := false
	hasTestTemplate := false
	hasProjectTemplate := false

	for _, tmpl := range templates {
		if filepath.HasPrefix(tmpl, "resource/") {
			hasResourceTemplate = true
		}
		if filepath.HasPrefix(tmpl, "test/") {
			hasTestTemplate = true
		}
		if filepath.HasPrefix(tmpl, "project/") {
			hasProjectTemplate = true
		}
	}

	assert.True(t, hasResourceTemplate, "should have resource templates")
	assert.True(t, hasTestTemplate, "should have test templates")
	assert.True(t, hasProjectTemplate, "should have project templates")
}

func TestTemplateLoader_GetTemplateInfo(t *testing.T) {
	fs := afero.NewMemMapFs()
	projectRoot := "/project"
	templatePath := filepath.Join(projectRoot, ProjectTemplateDir, "resource/layered/Controller.java.tmpl")

	require.NoError(t, fs.MkdirAll(filepath.Dir(templatePath), 0755))
	require.NoError(t, afero.WriteFile(fs, templatePath, []byte("content"), 0644))

	loader := NewTemplateLoaderWithHome(fs, projectRoot, "/home/user")

	info, err := loader.GetTemplateInfo("resource/layered/Controller.java.tmpl")
	require.NoError(t, err)

	assert.Equal(t, "resource/layered/Controller.java.tmpl", info.Name)
	assert.Equal(t, SourceProject, info.Source)
	assert.Contains(t, info.Path, ProjectTemplateDir)
}

func TestTemplateLoader_ListAllTemplatesWithSource(t *testing.T) {
	fs := afero.NewMemMapFs()
	projectRoot := "/project"
	homeDir := "/home/user"

	projectPath := filepath.Join(projectRoot, ProjectTemplateDir, "resource/layered/Controller.java.tmpl")
	require.NoError(t, fs.MkdirAll(filepath.Dir(projectPath), 0755))
	require.NoError(t, afero.WriteFile(fs, projectPath, []byte("project"), 0644))

	globalPath := filepath.Join(homeDir, GlobalTemplateDir, "custom/Custom.java.tmpl")
	require.NoError(t, fs.MkdirAll(filepath.Dir(globalPath), 0755))
	require.NoError(t, afero.WriteFile(fs, globalPath, []byte("global"), 0644))

	loader := NewTemplateLoaderWithHome(fs, projectRoot, homeDir)

	templates, err := loader.ListAllTemplatesWithSource()
	require.NoError(t, err)

	assert.NotEmpty(t, templates)

	hasProject := false
	hasGlobal := false
	hasEmbedded := false

	for _, info := range templates {
		switch info.Source {
		case SourceProject:
			hasProject = true
		case SourceGlobal:
			hasGlobal = true
		case SourceEmbedded:
			hasEmbedded = true
		}
	}

	assert.True(t, hasProject, "should have project templates")
	assert.True(t, hasGlobal, "should have global templates")
	assert.True(t, hasEmbedded, "should have embedded templates")
}

func TestTemplateSource_String(t *testing.T) {
	assert.Equal(t, "embedded", SourceEmbedded.String())
	assert.Equal(t, "project", SourceProject.String())
	assert.Equal(t, "global", SourceGlobal.String())
}

func TestTemplateLoader_EmptyProjectRoot(t *testing.T) {
	fs := afero.NewMemMapFs()
	loader := NewTemplateLoaderWithHome(fs, "", "/home/user")

	assert.Empty(t, loader.GetProjectTemplateDir())
	assert.False(t, loader.ProjectTemplatesExist())

	templates, err := loader.ListProjectTemplates()
	assert.NoError(t, err)
	assert.Nil(t, templates)
}

func TestTemplateLoader_EmptyHomeDir(t *testing.T) {
	fs := afero.NewMemMapFs()
	loader := NewTemplateLoaderWithHome(fs, "/project", "")

	assert.Empty(t, loader.GetGlobalTemplateDir())
	assert.False(t, loader.GlobalTemplatesExist())

	templates, err := loader.ListGlobalTemplates()
	assert.NoError(t, err)
	assert.Nil(t, templates)
}
