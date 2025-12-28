package output

import (
	"encoding/json"
	"fmt"
	"os"
)

type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *MetaInfo   `json:"meta,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type MetaInfo struct {
	Command   string `json:"command,omitempty"`
	Duration  string `json:"duration,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

type ProjectInfo struct {
	Name              string            `json:"name"`
	GroupID           string            `json:"groupId"`
	ArtifactID        string            `json:"artifactId"`
	Version           string            `json:"version"`
	Description       string            `json:"description,omitempty"`
	BuildTool         string            `json:"buildTool"`
	BuildFile         string            `json:"buildFile"`
	JavaVersion       string            `json:"javaVersion"`
	SpringBootVersion string            `json:"springBootVersion"`
	Packaging         string            `json:"packaging,omitempty"`
	BasePackage       string            `json:"basePackage,omitempty"`
	Dependencies      *DependencyInfo   `json:"dependencies"`
	Features          *FeatureInfo      `json:"features"`
	CodeStats         *CodeStatsInfo    `json:"codeStats,omitempty"`
	Architecture      *ArchitectureInfo `json:"architecture,omitempty"`
}

type DependencyInfo struct {
	Total          int                  `json:"total"`
	SpringStarters int                  `json:"springStarters"`
	SpringLibs     int                  `json:"springLibraries"`
	TestDeps       int                  `json:"testDependencies"`
	List           []DependencyListItem `json:"list,omitempty"`
}

type DependencyListItem struct {
	GroupID    string `json:"groupId"`
	ArtifactID string `json:"artifactId"`
	Version    string `json:"version,omitempty"`
	Scope      string `json:"scope,omitempty"`
}

type FeatureInfo struct {
	HasWeb        bool `json:"hasWeb"`
	HasJPA        bool `json:"hasJpa"`
	HasLombok     bool `json:"hasLombok"`
	HasValidation bool `json:"hasValidation"`
	HasMapStruct  bool `json:"hasMapStruct"`
	HasSecurity   bool `json:"hasSecurity"`
	HasActuator   bool `json:"hasActuator"`
	HasDevTools   bool `json:"hasDevTools"`
}

type ArchitectureInfo struct {
	Pattern      string `json:"pattern"`
	FeatureStyle string `json:"featureStyle,omitempty"`
	DTONaming    string `json:"dtoNaming,omitempty"`
	IDType       string `json:"idType,omitempty"`
	MapperType   string `json:"mapperType,omitempty"`
}

type CodeStatsInfo struct {
	TotalFiles      int64   `json:"totalFiles"`
	LinesOfCode     int64   `json:"linesOfCode"`
	Comments        int64   `json:"comments"`
	BlankLines      int64   `json:"blankLines"`
	TotalBytes      int64   `json:"totalBytes,omitempty"`
	TotalLines      int64   `json:"totalLines,omitempty"`
	Complexity      int64   `json:"complexity,omitempty"`
	EstimatedCost   float64 `json:"estimatedCost,omitempty"`
	EstimatedMonths float64 `json:"estimatedMonths,omitempty"`
	EstimatedPeople float64 `json:"estimatedPeople,omitempty"`
}

type LanguageStats struct {
	Name       string `json:"name"`
	Files      int64  `json:"files"`
	Lines      int64  `json:"lines"`
	Code       int64  `json:"code"`
	Comments   int64  `json:"comments"`
	Blanks     int64  `json:"blanks"`
	Complexity int64  `json:"complexity,omitempty"`
}

type StatsOutput struct {
	Languages []LanguageStats `json:"languages"`
	Summary   CodeStatsInfo   `json:"summary"`
}

type RouteInfo struct {
	Method     string `json:"method"`
	Path       string `json:"path"`
	Controller string `json:"controller"`
	Handler    string `json:"handler"`
	File       string `json:"file"`
	Line       int    `json:"line"`
}

type RoutesOutput struct {
	Routes []RouteInfo `json:"routes"`
	Total  int         `json:"total"`
}

type CatalogCategory struct {
	Name         string        `json:"name"`
	Dependencies []CatalogItem `json:"dependencies"`
}

type CatalogItem struct {
	Shortcut    string `json:"shortcut"`
	Name        string `json:"name"`
	Description string `json:"description"`
	GroupID     string `json:"groupId,omitempty"`
	ArtifactID  string `json:"artifactId,omitempty"`
}

type CatalogOutput struct {
	Categories []CatalogCategory `json:"categories"`
	Total      int               `json:"total"`
}

type GenerateResult struct {
	Type      string   `json:"type"`
	Name      string   `json:"name"`
	Generated []string `json:"generated"`
	Skipped   []string `json:"skipped,omitempty"`
	Errors    []string `json:"errors,omitempty"`
}

type GenerateOutput struct {
	Results        []GenerateResult `json:"results"`
	TotalGenerated int              `json:"totalGenerated"`
	TotalSkipped   int              `json:"totalSkipped"`
}

type AddRemoveResult struct {
	Action  string   `json:"action"`
	Added   []string `json:"added,omitempty"`
	Removed []string `json:"removed,omitempty"`
	Skipped []string `json:"skipped,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

type GeneratorType struct {
	Name        string   `json:"name"`
	Alias       string   `json:"alias,omitempty"`
	Description string   `json:"description"`
	Options     []string `json:"options,omitempty"`
}

type GeneratorTypesOutput struct {
	Types []GeneratorType `json:"types"`
}

type SecurityType struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SecurityTypesOutput struct {
	Types []SecurityType `json:"types"`
}

type InitOutput struct {
	ProjectName       string   `json:"projectName"`
	GroupID           string   `json:"groupId"`
	ArtifactID        string   `json:"artifactId"`
	PackageName       string   `json:"packageName"`
	JavaVersion       string   `json:"javaVersion"`
	SpringBootVersion string   `json:"springBootVersion"`
	BuildTool         string   `json:"buildTool"`
	Packaging         string   `json:"packaging"`
	Dependencies      []string `json:"dependencies"`
	ProjectDir        string   `json:"projectDir"`
	FilesCreated      []string `json:"filesCreated"`
	GitInitialized    bool     `json:"gitInitialized"`
}

type TemplateInfo struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Source   string `json:"source"`
	Path     string `json:"path,omitempty"`
}

type TemplateListOutput struct {
	Templates     []TemplateInfo `json:"templates"`
	Total         int            `json:"total"`
	ProjectCount  int            `json:"projectCount"`
	GlobalCount   int            `json:"globalCount"`
	EmbeddedCount int            `json:"embeddedCount"`
}

type DevCommandOutput struct {
	Command   string `json:"command"`
	BuildTool string `json:"buildTool"`
	ExitCode  int    `json:"exitCode"`
	Output    string `json:"output,omitempty"`
}

type TemplateInitOutput struct {
	TargetDir string   `json:"targetDir"`
	Scope     string   `json:"scope"`
	Category  string   `json:"category,omitempty"`
	Copied    []string `json:"copied"`
	Skipped   []string `json:"skipped"`
}

type TemplateValidationError struct {
	Line    int    `json:"line"`
	Message string `json:"message"`
}

type TemplateValidationResult struct {
	Path     string                    `json:"path"`
	Valid    bool                      `json:"valid"`
	Errors   []TemplateValidationError `json:"errors,omitempty"`
	Warnings []TemplateValidationError `json:"warnings,omitempty"`
}

type TemplateValidateOutput struct {
	Results      []TemplateValidationResult `json:"results"`
	TotalFiles   int                        `json:"totalFiles"`
	ValidCount   int                        `json:"validCount"`
	ErrorCount   int                        `json:"errorCount"`
	WarningCount int                        `json:"warningCount"`
}

func JSON(v interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}

func Success(data interface{}) error {
	return JSON(Response{
		Success: true,
		Data:    data,
	})
}

func Error(code, message string, details ...string) error {
	errInfo := &ErrorInfo{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		errInfo.Details = details[0]
	}
	return JSON(Response{
		Success: false,
		Error:   errInfo,
	})
}

func ErrorWithExit(code, message string, exitCode int) {
	_ = Error(code, message)
	os.Exit(exitCode)
}

func Print(format Format, textFn func() error, data interface{}) error {
	if format == FormatJSON {
		return Success(data)
	}
	return textFn()
}

func PrintRaw(format Format, textFn func() error, jsonFn func() error) error {
	if format == FormatJSON {
		return jsonFn()
	}
	return textFn()
}

func MustJSON(v interface{}) {
	if err := JSON(v); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}
