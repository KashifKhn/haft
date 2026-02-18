package doctor

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39"))

	passedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39"))

	suggestionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("141"))

	mutedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("255"))
)

func FormatReport(report *Report, opts Options) string {
	if opts.JSON {
		return formatJSON(report)
	}
	return formatText(report, opts)
}

func formatJSON(report *Report) string {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(data)
}

func formatText(report *Report, opts Options) string {
	var sb strings.Builder

	sb.WriteString("\n")
	sb.WriteString(titleStyle.Render("🏥 Haft Doctor - Project Health Check"))
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("=", 45))
	sb.WriteString("\n\n")

	fmt.Fprintf(&sb, "Project: %s\n", report.ProjectName)
	fmt.Fprintf(&sb, "Path: %s\n", report.ProjectPath)
	if report.BuildTool != "" {
		fmt.Fprintf(&sb, "Build Tool: %s\n", report.BuildTool)
	}
	if report.JavaVersion != "" {
		fmt.Fprintf(&sb, "Java: %s\n", report.JavaVersion)
	}
	if report.SpringVersion != "" {
		fmt.Fprintf(&sb, "Spring Boot: %s\n", report.SpringVersion)
	}
	sb.WriteString("\n")

	passed := filterResults(report.Results, true, "")
	errors := filterResults(report.Results, false, string(SeverityError))
	warnings := filterResults(report.Results, false, string(SeverityWarning))
	infos := filterResults(report.Results, false, string(SeverityInfo))
	suggestions := filterResults(report.Results, false, string(SeveritySuggestion))

	if len(passed) > 0 {
		sb.WriteString(headerStyle.Render("Passed Checks:"))
		sb.WriteString("\n")
		for _, r := range passed {
			fmt.Fprintf(&sb, "  %s %s\n",
				passedStyle.Render("✓"),
				r.Message,
			)
		}
		sb.WriteString("\n")
	}

	if len(errors) > 0 {
		sb.WriteString(headerStyle.Render("Issues:"))
		sb.WriteString("\n")
		for _, r := range errors {
			fmt.Fprintf(&sb, "  %s %s\n",
				errorStyle.Render("✗"),
				errorStyle.Render(r.Message),
			)
			if r.Details != "" {
				fmt.Fprintf(&sb, "    %s\n", mutedStyle.Render(r.Details))
			}
			if r.FixHint != "" {
				fmt.Fprintf(&sb, "    %s %s\n", mutedStyle.Render("→"), r.FixHint)
			}
		}
		sb.WriteString("\n")
	}

	if len(warnings) > 0 {
		sb.WriteString(headerStyle.Render("Warnings:"))
		sb.WriteString("\n")
		for _, r := range warnings {
			fmt.Fprintf(&sb, "  %s %s\n",
				warningStyle.Render("⚠"),
				warningStyle.Render(r.Message),
			)
			if r.Details != "" {
				fmt.Fprintf(&sb, "    %s\n", mutedStyle.Render(r.Details))
			}
			if r.FixHint != "" {
				fmt.Fprintf(&sb, "    %s %s\n", mutedStyle.Render("→"), r.FixHint)
			}
		}
		sb.WriteString("\n")
	}

	if len(infos) > 0 {
		sb.WriteString(headerStyle.Render("Info:"))
		sb.WriteString("\n")
		for _, r := range infos {
			fmt.Fprintf(&sb, "  %s %s\n",
				infoStyle.Render("ℹ"),
				r.Message,
			)
			if r.Details != "" {
				fmt.Fprintf(&sb, "    %s\n", mutedStyle.Render(r.Details))
			}
		}
		sb.WriteString("\n")
	}

	if len(suggestions) > 0 {
		sb.WriteString(headerStyle.Render("Suggestions:"))
		sb.WriteString("\n")
		for _, r := range suggestions {
			fmt.Fprintf(&sb, "  %s %s\n",
				suggestionStyle.Render("💡"),
				r.Message,
			)
			if r.Details != "" {
				fmt.Fprintf(&sb, "    %s\n", mutedStyle.Render(r.Details))
			}
			if r.FixHint != "" {
				fmt.Fprintf(&sb, "    %s %s\n", mutedStyle.Render("→"), r.FixHint)
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString(strings.Repeat("-", 45))
	sb.WriteString("\n")
	fmt.Fprintf(&sb, "Summary: %s passed",
		passedStyle.Render(fmt.Sprintf("%d", report.PassedCount)))

	if report.ErrorCount > 0 {
		fmt.Fprintf(&sb, ", %s",
			errorStyle.Render(fmt.Sprintf("%d errors", report.ErrorCount)))
	}
	if report.WarningCount > 0 {
		fmt.Fprintf(&sb, ", %s",
			warningStyle.Render(fmt.Sprintf("%d warnings", report.WarningCount)))
	}
	if report.SuggestionCount > 0 {
		fmt.Fprintf(&sb, ", %s",
			suggestionStyle.Render(fmt.Sprintf("%d suggestions", report.SuggestionCount)))
	}
	sb.WriteString("\n")

	return sb.String()
}

func filterResults(results []CheckResult, passed bool, severity string) []CheckResult {
	var filtered []CheckResult
	for _, r := range results {
		if passed && r.Passed {
			filtered = append(filtered, r)
		} else if !passed && !r.Passed && (severity == "" || string(r.Severity) == severity) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
