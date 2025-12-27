package generator

import (
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

var (
	simplePlaceholderRegex = regexp.MustCompile(`\$\{([a-zA-Z_][a-zA-Z0-9_]*)\}`)

	simpleVarMappings = map[string]string{
		"Name":        "{{.Name}}",
		"name":        "{{.NameLower}}",
		"nameLower":   "{{.NameLower}}",
		"nameCamel":   "{{.NameCamel}}",
		"nameSnake":   "{{snakeCase .Name}}",
		"nameKebab":   "{{kebabCase .Name}}",
		"namePlural":  "{{plural .NameLower}}",
		"NamePlural":  "{{plural .Name}}",
		"BasePackage": "{{.BasePackage}}",
		"basePackage": "{{.BasePackage}}",
		"Package":     "{{.Package}}",
		"package":     "{{.Package}}",
		"IDType":      "{{.IDType}}",
		"idType":      "{{.IDType}}",
		"TableName":   "{{.TableName}}",
		"tableName":   "{{.TableName}}",
	}

	conditionalStartRegex = regexp.MustCompile(`(?m)^\s*//\s*@if\s+(\w+)\s*$`)
	conditionalElseRegex  = regexp.MustCompile(`(?m)^\s*//\s*@else\s*$`)
	conditionalEndRegex   = regexp.MustCompile(`(?m)^\s*//\s*@endif\s*$`)

	conditionMappings = map[string]string{
		"HasLombok":     ".HasLombok",
		"HasJpa":        ".HasJpa",
		"HasValidation": ".HasValidation",
		"HasMapStruct":  ".HasMapStruct",
		"HasSwagger":    ".HasSwagger",
		"HasBaseEntity": ".HasBaseEntity",
		"UsesUUID":      "eq .IDType \"UUID\"",
		"UsesLong":      "eq .IDType \"Long\"",
	}
)

func PreprocessTemplate(content string) string {
	result := content
	result = preprocessSimplePlaceholders(result)
	result = preprocessConditionals(result)
	return result
}

func preprocessSimplePlaceholders(content string) string {
	return simplePlaceholderRegex.ReplaceAllStringFunc(content, func(match string) string {
		varName := simplePlaceholderRegex.FindStringSubmatch(match)[1]

		if goTemplate, ok := simpleVarMappings[varName]; ok {
			return goTemplate
		}

		if strings.HasPrefix(varName, "name") || strings.HasPrefix(varName, "Name") {
			return fmt.Sprintf("{{.%s}}", varName)
		}

		return fmt.Sprintf("{{.%s}}", varName)
	})
}

func preprocessConditionals(content string) string {
	result := conditionalStartRegex.ReplaceAllStringFunc(content, func(match string) string {
		submatch := conditionalStartRegex.FindStringSubmatch(match)
		if len(submatch) < 2 {
			return match
		}
		condition := submatch[1]

		if goCondition, ok := conditionMappings[condition]; ok {
			return fmt.Sprintf("{{if %s}}", goCondition)
		}

		return fmt.Sprintf("{{if .%s}}", condition)
	})

	result = conditionalElseRegex.ReplaceAllString(result, "{{else}}")
	result = conditionalEndRegex.ReplaceAllString(result, "{{end}}")

	return result
}

type ValidationError struct {
	Line    int
	Column  int
	Message string
}

type ValidationResult struct {
	Valid    bool
	Errors   []ValidationError
	Warnings []ValidationError
}

func ValidateTemplate(content string, templateName string) ValidationResult {
	result := ValidationResult{Valid: true}

	lines := strings.Split(content, "\n")

	unclosedPlaceholders := findUnclosedPlaceholders(lines)
	for _, err := range unclosedPlaceholders {
		result.Errors = append(result.Errors, err)
		result.Valid = false
	}

	conditionalErrors := validateConditionals(lines)
	for _, err := range conditionalErrors {
		result.Errors = append(result.Errors, err)
		result.Valid = false
	}

	unknownVars := findUnknownVariables(lines)
	for _, warn := range unknownVars {
		result.Warnings = append(result.Warnings, warn)
	}

	preprocessed := PreprocessTemplate(content)
	goTemplateErrors := validateGoTemplate(preprocessed, templateName)
	for _, err := range goTemplateErrors {
		result.Errors = append(result.Errors, err)
		result.Valid = false
	}

	return result
}

func findUnclosedPlaceholders(lines []string) []ValidationError {
	var errors []ValidationError

	for i, line := range lines {
		dollarBraceIdx := 0
		for {
			idx := strings.Index(line[dollarBraceIdx:], "${")
			if idx == -1 {
				break
			}
			startIdx := dollarBraceIdx + idx
			closeIdx := strings.Index(line[startIdx:], "}")
			if closeIdx == -1 {
				errors = append(errors, ValidationError{
					Line:    i + 1,
					Column:  startIdx + 1,
					Message: "unclosed placeholder - missing closing '}'",
				})
			}
			dollarBraceIdx = startIdx + 2
		}
	}

	return errors
}

func validateConditionals(lines []string) []ValidationError {
	var errors []ValidationError
	var stack []int

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if conditionalStartRegex.MatchString(line) {
			stack = append(stack, i+1)
		} else if conditionalEndRegex.MatchString(line) {
			if len(stack) == 0 {
				errors = append(errors, ValidationError{
					Line:    i + 1,
					Column:  1,
					Message: "@endif without matching @if",
				})
			} else {
				stack = stack[:len(stack)-1]
			}
		} else if conditionalElseRegex.MatchString(line) {
			if len(stack) == 0 {
				errors = append(errors, ValidationError{
					Line:    i + 1,
					Column:  1,
					Message: "@else without matching @if",
				})
			}
		}

		if strings.Contains(trimmed, "// @") || strings.HasPrefix(trimmed, "// @") {
			directive := strings.TrimPrefix(trimmed, "//")
			directive = strings.TrimSpace(directive)
			if strings.HasPrefix(directive, "@") {
				if !conditionalStartRegex.MatchString(line) &&
					!conditionalElseRegex.MatchString(line) &&
					!conditionalEndRegex.MatchString(line) {
					errors = append(errors, ValidationError{
						Line:    i + 1,
						Column:  strings.Index(line, "@") + 1,
						Message: fmt.Sprintf("unknown directive: %s (valid: @if, @else, @endif)", directive),
					})
				}
			}
		}
	}

	for _, lineNum := range stack {
		errors = append(errors, ValidationError{
			Line:    lineNum,
			Column:  1,
			Message: "@if without matching @endif",
		})
	}

	return errors
}

func findUnknownVariables(lines []string) []ValidationError {
	var warnings []ValidationError

	knownVars := map[string]bool{
		"Name": true, "name": true, "nameLower": true, "nameCamel": true,
		"nameSnake": true, "nameKebab": true, "namePlural": true, "NamePlural": true,
		"BasePackage": true, "basePackage": true, "Package": true, "package": true,
		"IDType": true, "idType": true, "TableName": true, "tableName": true,
		"HasLombok": true, "HasJpa": true, "HasValidation": true,
		"HasMapStruct": true, "HasSwagger": true, "HasBaseEntity": true,
		"NameCamel": true, "NameLower": true, "NameSnake": true, "NameKebab": true,
		"Year": true, "Date": true, "Author": true,
	}

	for i, line := range lines {
		matches := simplePlaceholderRegex.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				varName := match[1]
				if !knownVars[varName] {
					warnings = append(warnings, ValidationError{
						Line:    i + 1,
						Column:  strings.Index(line, "${"+varName+"}") + 1,
						Message: fmt.Sprintf("unknown variable '%s' - may not be available in all contexts", varName),
					})
				}
			}
		}
	}

	return warnings
}

func validateGoTemplate(content string, name string) []ValidationError {
	var errors []ValidationError

	funcMap := defaultFuncMap()
	_, err := template.New(name).Funcs(funcMap).Parse(content)
	if err != nil {
		errors = append(errors, ValidationError{
			Line:    1,
			Column:  1,
			Message: fmt.Sprintf("template syntax error: %v", err),
		})
	}

	return errors
}

type VariableInfo struct {
	Name        string
	Description string
	Example     string
}

type ConditionInfo struct {
	Name        string
	Description string
}

func GetAvailableVariables() []VariableInfo {
	return []VariableInfo{
		{Name: "${Name}", Description: "Resource name in PascalCase", Example: "User"},
		{Name: "${name}", Description: "Resource name in lowercase", Example: "user"},
		{Name: "${nameCamel}", Description: "Resource name in camelCase", Example: "user"},
		{Name: "${nameSnake}", Description: "Resource name in snake_case", Example: "user_name"},
		{Name: "${nameKebab}", Description: "Resource name in kebab-case", Example: "user-name"},
		{Name: "${namePlural}", Description: "Pluralized lowercase name", Example: "users"},
		{Name: "${NamePlural}", Description: "Pluralized PascalCase name", Example: "Users"},
		{Name: "${BasePackage}", Description: "Base package path", Example: "com.example.app"},
		{Name: "${Package}", Description: "Full package path for current file", Example: "com.example.app.controller"},
		{Name: "${IDType}", Description: "Entity ID type", Example: "Long or UUID"},
		{Name: "${TableName}", Description: "Database table name", Example: "users"},
	}
}

func GetAvailableConditions() []ConditionInfo {
	return []ConditionInfo{
		{Name: "HasLombok", Description: "True if Lombok is available in project"},
		{Name: "HasJpa", Description: "True if Spring Data JPA is available"},
		{Name: "HasValidation", Description: "True if Bean Validation is available"},
		{Name: "HasMapStruct", Description: "True if MapStruct is available"},
		{Name: "HasSwagger", Description: "True if Swagger/OpenAPI is available"},
		{Name: "HasBaseEntity", Description: "True if project has a base entity class"},
		{Name: "UsesUUID", Description: "True if entity uses UUID as ID type"},
		{Name: "UsesLong", Description: "True if entity uses Long as ID type"},
	}
}
