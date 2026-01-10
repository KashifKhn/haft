package doctor

type Severity string

const (
	SeverityError      Severity = "error"
	SeverityWarning    Severity = "warning"
	SeverityInfo       Severity = "info"
	SeveritySuggestion Severity = "suggestion"
)

type Category string

const (
	CategoryBuild        Category = "build"
	CategorySource       Category = "source"
	CategoryConfig       Category = "config"
	CategorySecurity     Category = "security"
	CategoryDependencies Category = "dependencies"
	CategoryBestPractice Category = "best-practice"
)

type CheckResult struct {
	Name     string   `json:"name"`
	Category Category `json:"category"`
	Passed   bool     `json:"passed"`
	Severity Severity `json:"severity"`
	Message  string   `json:"message"`
	Details  string   `json:"details,omitempty"`
	FixHint  string   `json:"fix_hint,omitempty"`
}

type Report struct {
	ProjectPath     string        `json:"project_path"`
	ProjectName     string        `json:"project_name"`
	BuildTool       string        `json:"build_tool"`
	SpringVersion   string        `json:"spring_version,omitempty"`
	JavaVersion     string        `json:"java_version,omitempty"`
	Results         []CheckResult `json:"results"`
	PassedCount     int           `json:"passed_count"`
	ErrorCount      int           `json:"error_count"`
	WarningCount    int           `json:"warning_count"`
	InfoCount       int           `json:"info_count"`
	SuggestionCount int           `json:"suggestion_count"`
}

type Options struct {
	JSON     bool
	Fix      bool
	Strict   bool
	Category string
}

func (r *Report) CalculateCounts() {
	r.PassedCount = 0
	r.ErrorCount = 0
	r.WarningCount = 0
	r.InfoCount = 0
	r.SuggestionCount = 0

	for _, result := range r.Results {
		if result.Passed {
			r.PassedCount++
			continue
		}
		switch result.Severity {
		case SeverityError:
			r.ErrorCount++
		case SeverityWarning:
			r.WarningCount++
		case SeverityInfo:
			r.InfoCount++
		case SeveritySuggestion:
			r.SuggestionCount++
		}
	}
}

func (r *Report) HasIssues() bool {
	return r.ErrorCount > 0
}

func (r *Report) HasWarnings() bool {
	return r.WarningCount > 0
}
