package generator

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEngine(t *testing.T) {
	fs := afero.NewMemMapFs()
	engine := NewEngine(fs)

	assert.NotNil(t, engine)
	assert.NotNil(t, engine.fs)
	assert.NotNil(t, engine.funcMap)
}

func TestRenderString(t *testing.T) {
	fs := afero.NewMemMapFs()
	engine := NewEngine(fs)

	tests := []struct {
		name     string
		template string
		data     any
		expected string
	}{
		{
			name:     "simple variable",
			template: "Hello {{.Name}}!",
			data:     map[string]string{"Name": "World"},
			expected: "Hello World!",
		},
		{
			name:     "lower function",
			template: "{{lower .Name}}",
			data:     map[string]string{"Name": "HELLO"},
			expected: "hello",
		},
		{
			name:     "upper function",
			template: "{{upper .Name}}",
			data:     map[string]string{"Name": "hello"},
			expected: "HELLO",
		},
		{
			name:     "capitalize function",
			template: "{{capitalize .Name}}",
			data:     map[string]string{"Name": "hello"},
			expected: "Hello",
		},
		{
			name:     "camelCase function",
			template: "{{camelCase .Name}}",
			data:     map[string]string{"Name": "user-profile"},
			expected: "userProfile",
		},
		{
			name:     "pascalCase function",
			template: "{{pascalCase .Name}}",
			data:     map[string]string{"Name": "user-profile"},
			expected: "UserProfile",
		},
		{
			name:     "snakeCase function",
			template: "{{snakeCase .Name}}",
			data:     map[string]string{"Name": "UserProfile"},
			expected: "user_profile",
		},
		{
			name:     "kebabCase function",
			template: "{{kebabCase .Name}}",
			data:     map[string]string{"Name": "UserProfile"},
			expected: "user-profile",
		},
		{
			name:     "plural function",
			template: "{{plural .Name}}",
			data:     map[string]string{"Name": "user"},
			expected: "users",
		},
		{
			name:     "singular function",
			template: "{{singular .Name}}",
			data:     map[string]string{"Name": "users"},
			expected: "user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.RenderString(tt.template, tt.data)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWriteFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	engine := NewEngine(fs)

	err := engine.WriteFile("/test/dir/file.txt", "content")
	require.NoError(t, err)

	content, err := afero.ReadFile(fs, "/test/dir/file.txt")
	require.NoError(t, err)
	assert.Equal(t, "content", string(content))
}

func TestFileExists(t *testing.T) {
	fs := afero.NewMemMapFs()
	engine := NewEngine(fs)

	assert.False(t, engine.FileExists("/nonexistent"))

	err := engine.WriteFile("/exists.txt", "content")
	require.NoError(t, err)
	assert.True(t, engine.FileExists("/exists.txt"))
}

func TestRenderTemplate(t *testing.T) {
	fs := afero.NewMemMapFs()
	engine := NewEngine(fs)

	data := map[string]any{
		"BasePackage":     "com.example.demo",
		"ApplicationName": "Demo",
	}

	result, err := engine.RenderTemplate("project/Application.java.tmpl", data)
	require.NoError(t, err)
	assert.Contains(t, result, "package com.example.demo")
	assert.Contains(t, result, "class DemoApplication")
	assert.Contains(t, result, "@SpringBootApplication")
}

func TestRenderAndWrite(t *testing.T) {
	fs := afero.NewMemMapFs()
	engine := NewEngine(fs)

	data := map[string]any{
		"BasePackage":     "com.example.demo",
		"ApplicationName": "Demo",
	}

	err := engine.RenderAndWrite("project/Application.java.tmpl", "/output/Application.java", data)
	require.NoError(t, err)

	content, err := afero.ReadFile(fs, "/output/Application.java")
	require.NoError(t, err)
	assert.Contains(t, string(content), "package com.example.demo")
}

func TestListTemplates(t *testing.T) {
	fs := afero.NewMemMapFs()
	engine := NewEngine(fs)

	templates, err := engine.ListTemplates("project")
	require.NoError(t, err)
	assert.NotEmpty(t, templates)
	assert.Contains(t, templates, "project/Application.java.tmpl")
	assert.Contains(t, templates, "project/pom.xml.tmpl")
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"a", "A"},
		{"hello", "Hello"},
		{"Hello", "Hello"},
		{"HELLO", "HELLO"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, capitalize(tt.input))
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user", "user"},
		{"user-profile", "userProfile"},
		{"user_profile", "userProfile"},
		{"UserProfile", "userProfile"},
		{"USER_PROFILE", "userProfile"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, toCamelCase(tt.input))
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user", "User"},
		{"user-profile", "UserProfile"},
		{"user_profile", "UserProfile"},
		{"userProfile", "UserProfile"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, toPascalCase(tt.input))
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user", "user"},
		{"UserProfile", "user_profile"},
		{"user-profile", "user_profile"},
		{"userProfile", "user_profile"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, toSnakeCase(tt.input))
		})
	}
}

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user", "user"},
		{"UserProfile", "user-profile"},
		{"user_profile", "user-profile"},
		{"userProfile", "user-profile"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, toKebabCase(tt.input))
		})
	}
}

func TestPluralize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"user", "users"},
		{"category", "categories"},
		{"box", "boxes"},
		{"church", "churches"},
		{"bush", "bushes"},
		{"person", "people"},
		{"child", "children"},
		{"leaf", "leaves"},
		{"Person", "People"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, pluralize(tt.input))
		})
	}
}

func TestSingularize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"users", "user"},
		{"categories", "category"},
		{"boxes", "box"},
		{"churches", "church"},
		{"people", "person"},
		{"children", "child"},
		{"leaves", "leaf"},
		{"People", "Person"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, singularize(tt.input))
		})
	}
}

func TestToTitleCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"hello", "Hello"},
		{"hello world", "Hello World"},
		{"HELLO WORLD", "Hello World"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, toTitleCase(tt.input))
		})
	}
}
