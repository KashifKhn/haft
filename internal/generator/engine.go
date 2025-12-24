package generator

import (
	"bytes"
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/afero"
)

//go:embed all:templates
var templateFS embed.FS

type Engine struct {
	fs        afero.Fs
	templates *template.Template
	funcMap   template.FuncMap
}

func NewEngine(filesystem afero.Fs) *Engine {
	e := &Engine{
		fs:      filesystem,
		funcMap: defaultFuncMap(),
	}
	return e
}

func defaultFuncMap() template.FuncMap {
	return template.FuncMap{
		"lower":      strings.ToLower,
		"upper":      strings.ToUpper,
		"title":      toTitleCase,
		"capitalize": capitalize,
		"camelCase":  toCamelCase,
		"pascalCase": toPascalCase,
		"snakeCase":  toSnakeCase,
		"kebabCase":  toKebabCase,
		"plural":     pluralize,
		"singular":   singularize,
		"package":    toPackagePath,
	}
}

func (e *Engine) RenderTemplate(name string, data any) (string, error) {
	content, err := templateFS.ReadFile("templates/" + name)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New(name).Funcs(e.funcMap).Parse(string(content))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (e *Engine) RenderString(content string, data any) (string, error) {
	tmpl, err := template.New("inline").Funcs(e.funcMap).Parse(content)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (e *Engine) WriteFile(path string, content string) error {
	dir := filepath.Dir(path)
	if err := e.fs.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return afero.WriteFile(e.fs, path, []byte(content), 0644)
}

func (e *Engine) RenderAndWrite(templateName string, outputPath string, data any) error {
	content, err := e.RenderTemplate(templateName, data)
	if err != nil {
		return err
	}
	return e.WriteFile(outputPath, content)
}

func (e *Engine) FileExists(path string) bool {
	_, err := e.fs.Stat(path)
	return err == nil
}

func (e *Engine) CopyTemplateDir(templateDir string, outputDir string, data any) error {
	return fs.WalkDir(templateFS, "templates/"+templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath := strings.TrimPrefix(path, "templates/"+templateDir)
		if relPath == "" {
			return nil
		}
		relPath = strings.TrimPrefix(relPath, "/")

		outputPath := filepath.Join(outputDir, relPath)

		outputPath, err = e.RenderString(outputPath, data)
		if err != nil {
			return err
		}

		if d.IsDir() {
			return e.fs.MkdirAll(outputPath, 0755)
		}

		if strings.HasSuffix(path, ".tmpl") {
			outputPath = strings.TrimSuffix(outputPath, ".tmpl")
			content, err := templateFS.ReadFile(path)
			if err != nil {
				return err
			}
			rendered, err := e.RenderString(string(content), data)
			if err != nil {
				return err
			}
			return e.WriteFile(outputPath, rendered)
		}

		content, err := templateFS.ReadFile(path)
		if err != nil {
			return err
		}
		return e.WriteFile(outputPath, string(content))
	})
}

func (e *Engine) ListTemplates(dir string) ([]string, error) {
	var templates []string
	err := fs.WalkDir(templateFS, "templates/"+dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			relPath := strings.TrimPrefix(path, "templates/")
			templates = append(templates, relPath)
		}
		return nil
	})
	return templates, err
}

func (e *Engine) GetFS() afero.Fs {
	return e.fs
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func toTitleCase(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		words[i] = capitalize(strings.ToLower(word))
	}
	return strings.Join(words, " ")
}

func toCamelCase(s string) string {
	words := splitWords(s)
	if len(words) == 0 {
		return s
	}
	result := strings.ToLower(words[0])
	for _, word := range words[1:] {
		result += capitalize(strings.ToLower(word))
	}
	return result
}

func toPascalCase(s string) string {
	words := splitWords(s)
	var result string
	for _, word := range words {
		result += capitalize(strings.ToLower(word))
	}
	return result
}

func toSnakeCase(s string) string {
	words := splitWords(s)
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}
	return strings.Join(words, "_")
}

func toKebabCase(s string) string {
	words := splitWords(s)
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}
	return strings.Join(words, "-")
}

func splitWords(s string) []string {
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")

	var words []string
	var currentWord strings.Builder

	for i, r := range s {
		if r == ' ' {
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
			continue
		}

		if i > 0 && isUpper(r) && currentWord.Len() > 0 {
			lastChar := []rune(currentWord.String())[currentWord.Len()-1]
			if !isUpper(lastChar) {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
		}

		currentWord.WriteRune(r)
	}

	if currentWord.Len() > 0 {
		words = append(words, currentWord.String())
	}

	return words
}

func isUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func pluralize(s string) string {
	if len(s) == 0 {
		return s
	}

	irregulars := map[string]string{
		"person": "people",
		"child":  "children",
		"man":    "men",
		"woman":  "women",
		"foot":   "feet",
		"tooth":  "teeth",
		"goose":  "geese",
		"mouse":  "mice",
	}

	lower := strings.ToLower(s)
	if plural, ok := irregulars[lower]; ok {
		if isUpper(rune(s[0])) {
			return capitalize(plural)
		}
		return plural
	}

	if strings.HasSuffix(lower, "s") || strings.HasSuffix(lower, "x") ||
		strings.HasSuffix(lower, "z") || strings.HasSuffix(lower, "ch") ||
		strings.HasSuffix(lower, "sh") {
		return s + "es"
	}

	if strings.HasSuffix(lower, "y") && len(s) > 1 {
		prev := lower[len(lower)-2]
		if prev != 'a' && prev != 'e' && prev != 'i' && prev != 'o' && prev != 'u' {
			return s[:len(s)-1] + "ies"
		}
	}

	if strings.HasSuffix(lower, "f") {
		return s[:len(s)-1] + "ves"
	}

	if strings.HasSuffix(lower, "fe") {
		return s[:len(s)-2] + "ves"
	}

	return s + "s"
}

func singularize(s string) string {
	if len(s) == 0 {
		return s
	}

	irregulars := map[string]string{
		"people":   "person",
		"children": "child",
		"men":      "man",
		"women":    "woman",
		"feet":     "foot",
		"teeth":    "tooth",
		"geese":    "goose",
		"mice":     "mouse",
	}

	lower := strings.ToLower(s)
	if singular, ok := irregulars[lower]; ok {
		if isUpper(rune(s[0])) {
			return capitalize(singular)
		}
		return singular
	}

	if strings.HasSuffix(lower, "ies") && len(s) > 3 {
		return s[:len(s)-3] + "y"
	}

	if strings.HasSuffix(lower, "ves") {
		return s[:len(s)-3] + "f"
	}

	if strings.HasSuffix(lower, "es") {
		base := s[:len(s)-2]
		baseLower := strings.ToLower(base)
		if strings.HasSuffix(baseLower, "s") || strings.HasSuffix(baseLower, "x") ||
			strings.HasSuffix(baseLower, "z") || strings.HasSuffix(baseLower, "ch") ||
			strings.HasSuffix(baseLower, "sh") {
			return base
		}
	}

	if strings.HasSuffix(lower, "s") && len(s) > 1 {
		return s[:len(s)-1]
	}

	return s
}

func toPackagePath(s string) string {
	return strings.ReplaceAll(s, ".", string(os.PathSeparator))
}
