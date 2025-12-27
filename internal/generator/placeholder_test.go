package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreprocessSimplePlaceholders(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple Name placeholder",
			input:    "public class ${Name}Controller {",
			expected: "public class {{.Name}}Controller {",
		},
		{
			name:     "lowercase name placeholder",
			input:    "private ${name}Service service;",
			expected: "private {{.NameLower}}Service service;",
		},
		{
			name:     "camelCase placeholder",
			input:    "private final ${nameCamel}Service;",
			expected: "private final {{.NameCamel}}Service;",
		},
		{
			name:     "BasePackage placeholder",
			input:    "package ${BasePackage}.controller;",
			expected: "package {{.BasePackage}}.controller;",
		},
		{
			name:     "multiple placeholders",
			input:    "public class ${Name}Controller extends ${Name}Base {",
			expected: "public class {{.Name}}Controller extends {{.Name}}Base {",
		},
		{
			name:     "pluralized placeholders",
			input:    "/api/${namePlural}",
			expected: "/api/{{plural .NameLower}}",
		},
		{
			name:     "snake case placeholder",
			input:    "@Table(name = \"${nameSnake}\")",
			expected: "@Table(name = \"{{snakeCase .Name}}\")",
		},
		{
			name:     "no placeholders",
			input:    "public class UserController {",
			expected: "public class UserController {",
		},
		{
			name:     "preserves Go template syntax",
			input:    "{{.Name}} and ${Name}",
			expected: "{{.Name}} and {{.Name}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := preprocessSimplePlaceholders(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPreprocessConditionals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple if condition",
			input:    "// @if HasLombok",
			expected: "{{if .HasLombok}}",
		},
		{
			name:     "endif",
			input:    "// @endif",
			expected: "{{end}}",
		},
		{
			name:     "else",
			input:    "// @else",
			expected: "{{else}}",
		},
		{
			name:     "if with indentation",
			input:    "    // @if HasValidation",
			expected: "{{if .HasValidation}}",
		},
		{
			name:     "complete block",
			input:    "// @if HasLombok\n@Data\n// @endif",
			expected: "{{if .HasLombok}}\n@Data\n{{end}}",
		},
		{
			name:     "with else block",
			input:    "// @if HasLombok\n@Data\n// @else\n// no lombok\n// @endif",
			expected: "{{if .HasLombok}}\n@Data\n{{else}}\n// no lombok\n{{end}}",
		},
		{
			name:     "UsesUUID maps to eq comparison",
			input:    "// @if UsesUUID",
			expected: "{{if eq .IDType \"UUID\"}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := preprocessConditionals(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPreprocessTemplate(t *testing.T) {
	input := `package ${BasePackage}.controller;

// @if HasLombok
import lombok.Data;
// @endif

public class ${Name}Controller {
    private final ${Name}Service ${nameCamel}Service;
}`

	result := PreprocessTemplate(input)

	assert.Contains(t, result, "{{.BasePackage}}")
	assert.Contains(t, result, "{{if .HasLombok}}")
	assert.Contains(t, result, "{{end}}")
	assert.Contains(t, result, "{{.Name}}Controller")
	assert.Contains(t, result, "{{.NameCamel}}Service")
	assert.NotContains(t, result, "${")
	assert.NotContains(t, result, "// @if")
}

func TestValidateTemplate_Valid(t *testing.T) {
	template := `package ${BasePackage}.controller;

// @if HasLombok
@Data
// @endif

public class ${Name}Controller {
    private final ${Name}Service ${nameCamel}Service;
}`

	result := ValidateTemplate(template, "test.tmpl")

	assert.True(t, result.Valid)
	assert.Empty(t, result.Errors)
}

func TestValidateTemplate_UnclosedPlaceholder(t *testing.T) {
	template := `public class ${Name Controller {`

	result := ValidateTemplate(template, "test.tmpl")

	assert.False(t, result.Valid)
	assert.NotEmpty(t, result.Errors)

	hasUnclosedError := false
	for _, err := range result.Errors {
		if err.Message == "unclosed placeholder - missing closing '}'" {
			hasUnclosedError = true
			break
		}
	}
	assert.True(t, hasUnclosedError, "should have unclosed placeholder error")
}

func TestValidateTemplate_UnmatchedIf(t *testing.T) {
	template := `// @if HasLombok
@Data
// missing @endif`

	result := ValidateTemplate(template, "test.tmpl")

	assert.False(t, result.Valid)

	hasIfError := false
	for _, err := range result.Errors {
		if err.Message == "@if without matching @endif" {
			hasIfError = true
			break
		}
	}
	assert.True(t, hasIfError, "should have unmatched @if error")
}

func TestValidateTemplate_UnmatchedEndif(t *testing.T) {
	template := `@Data
// @endif`

	result := ValidateTemplate(template, "test.tmpl")

	assert.False(t, result.Valid)

	hasEndifError := false
	for _, err := range result.Errors {
		if err.Message == "@endif without matching @if" {
			hasEndifError = true
			break
		}
	}
	assert.True(t, hasEndifError, "should have unmatched @endif error")
}

func TestValidateTemplate_UnknownVariable(t *testing.T) {
	template := `public class ${Name}Controller {
    private ${unknownVar} field;
}`

	result := ValidateTemplate(template, "test.tmpl")

	assert.True(t, result.Valid)
	assert.NotEmpty(t, result.Warnings)

	hasWarning := false
	for _, warn := range result.Warnings {
		if warn.Message == "unknown variable 'unknownVar' - may not be available in all contexts" {
			hasWarning = true
			break
		}
	}
	assert.True(t, hasWarning, "should have unknown variable warning")
}

func TestGetAvailableVariables(t *testing.T) {
	vars := GetAvailableVariables()

	assert.NotEmpty(t, vars)

	varNames := make(map[string]bool)
	for _, v := range vars {
		varNames[v.Name] = true
	}

	assert.True(t, varNames["${Name}"])
	assert.True(t, varNames["${name}"])
	assert.True(t, varNames["${BasePackage}"])
}

func TestGetAvailableConditions(t *testing.T) {
	conditions := GetAvailableConditions()

	assert.NotEmpty(t, conditions)

	condNames := make(map[string]bool)
	for _, c := range conditions {
		condNames[c.Name] = true
	}

	assert.True(t, condNames["HasLombok"])
	assert.True(t, condNames["HasJpa"])
	assert.True(t, condNames["HasValidation"])
}

func TestEngineRenderWithSimpleSyntax(t *testing.T) {
	engine := NewEngine(nil)

	template := `Hello ${Name}!`
	data := map[string]any{
		"Name": "World",
	}

	result, err := engine.RenderString(template, data)

	assert.NoError(t, err)
	assert.Equal(t, "Hello World!", result)
}

func TestEngineRenderWithConditionals(t *testing.T) {
	engine := NewEngine(nil)

	template := `// @if HasLombok
@Data
// @else
// no lombok
// @endif
class Test {}`

	dataWithLombok := map[string]any{
		"HasLombok": true,
	}

	result, err := engine.RenderString(template, dataWithLombok)
	assert.NoError(t, err)
	assert.Contains(t, result, "@Data")
	assert.NotContains(t, result, "no lombok")

	dataWithoutLombok := map[string]any{
		"HasLombok": false,
	}

	result, err = engine.RenderString(template, dataWithoutLombok)
	assert.NoError(t, err)
	assert.NotContains(t, result, "@Data")
	assert.Contains(t, result, "no lombok")
}

func TestEngineRenderMixedSyntax(t *testing.T) {
	engine := NewEngine(nil)

	template := `package ${BasePackage}.controller;
{{if .HasValidation}}import jakarta.validation.Valid;{{end}}

public class ${Name}Controller {
    private final ${Name}Service {{.NameCamel}}Service;
}`

	data := map[string]any{
		"Name":          "User",
		"NameCamel":     "user",
		"BasePackage":   "com.example",
		"HasValidation": true,
	}

	result, err := engine.RenderString(template, data)

	assert.NoError(t, err)
	assert.Contains(t, result, "package com.example.controller;")
	assert.Contains(t, result, "import jakarta.validation.Valid;")
	assert.Contains(t, result, "public class UserController")
	assert.Contains(t, result, "UserService userService")
}
