package stats

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/KashifKhn/haft/internal/stats"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	labelStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	valueStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	langStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	numberStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	cocomoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
)

func NewCommand() *cobra.Command {
	var jsonOutput bool
	var showCocomo bool

	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Show code statistics",
		Long: `Display code statistics for the current Spring Boot project.

Uses SCC (Sloc Cloc and Code) to count lines of code, comments, blanks,
and complexity for each language in the project.`,
		Example: `  # Show code statistics
  haft stats

  # Output as JSON
  haft stats --json

  # Include COCOMO cost estimates
  haft stats --cocomo`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStats(jsonOutput, showCocomo)
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	cmd.Flags().BoolVar(&showCocomo, "cocomo", false, "Show COCOMO cost estimates")

	return cmd
}

func runStats(jsonOutput bool, showCocomo bool) error {
	fs := afero.NewOsFs()
	_, err := buildtool.DetectWithCwd(fs)
	if err != nil {
		return fmt.Errorf("not a Spring Boot project: %w", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	projectStats, err := stats.CountProject(cwd)
	if err != nil {
		return fmt.Errorf("failed to count code: %w", err)
	}

	if jsonOutput {
		return printJSON(projectStats, showCocomo)
	}

	return printFormatted(projectStats, showCocomo)
}

func printFormatted(s *stats.ProjectStats, showCocomo bool) error {
	fmt.Println()
	fmt.Println(headerStyle.Render("  Code Statistics"))
	fmt.Println(strings.Repeat("─", 85))

	fmt.Printf("  %-20s %10s %10s %10s %10s %10s\n",
		labelStyle.Render("Language"),
		labelStyle.Render("Files"),
		labelStyle.Render("Lines"),
		labelStyle.Render("Code"),
		labelStyle.Render("Comments"),
		labelStyle.Render("Blanks"))
	fmt.Println(strings.Repeat("─", 85))

	for _, lang := range s.Languages {
		fmt.Printf("  %-20s %10s %10s %10s %10s %10s\n",
			langStyle.Render(truncate(lang.Name, 20)),
			numberStyle.Render(formatNumber(lang.Files)),
			valueStyle.Render(formatNumber(lang.Lines)),
			numberStyle.Render(formatNumber(lang.Code)),
			valueStyle.Render(formatNumber(lang.Comments)),
			valueStyle.Render(formatNumber(lang.Blanks)))
	}

	fmt.Println(strings.Repeat("─", 85))
	fmt.Printf("  %-20s %10s %10s %10s %10s %10s\n",
		headerStyle.Render("Total"),
		numberStyle.Render(formatNumber(s.TotalFiles)),
		valueStyle.Render(formatNumber(s.TotalLines)),
		numberStyle.Render(formatNumber(s.TotalCode)),
		valueStyle.Render(formatNumber(s.TotalComments)),
		valueStyle.Render(formatNumber(s.TotalBlanks)))

	if showCocomo && s.EstimatedCost > 0 {
		fmt.Println()
		fmt.Println(headerStyle.Render("  COCOMO Estimates"))
		fmt.Println(strings.Repeat("─", 85))
		fmt.Printf("  %s %s\n",
			labelStyle.Render(fmt.Sprintf("%-18s", "Estimated Cost:")),
			cocomoStyle.Render(fmt.Sprintf("$%s", formatNumber(int64(s.EstimatedCost)))))
		fmt.Printf("  %s %s\n",
			labelStyle.Render(fmt.Sprintf("%-18s", "Schedule Effort:")),
			cocomoStyle.Render(fmt.Sprintf("%.2f months", s.EstimatedMonths)))
		fmt.Printf("  %s %s\n",
			labelStyle.Render(fmt.Sprintf("%-18s", "People Required:")),
			cocomoStyle.Render(fmt.Sprintf("%.2f", s.EstimatedPeople)))
	}

	fmt.Println()
	fmt.Printf("  %s %s\n",
		labelStyle.Render("Processed:"),
		valueStyle.Render(formatBytes(s.TotalBytes)))
	fmt.Println()

	return nil
}

func printJSON(s *stats.ProjectStats, showCocomo bool) error {
	type jsonLang struct {
		Name       string `json:"name"`
		Files      int64  `json:"files"`
		Lines      int64  `json:"lines"`
		Code       int64  `json:"code"`
		Comments   int64  `json:"comments"`
		Blanks     int64  `json:"blanks"`
		Complexity int64  `json:"complexity,omitempty"`
	}

	type jsonOutput struct {
		Languages       []jsonLang `json:"languages"`
		TotalFiles      int64      `json:"totalFiles"`
		TotalLines      int64      `json:"totalLines"`
		TotalCode       int64      `json:"totalCode"`
		TotalComments   int64      `json:"totalComments"`
		TotalBlanks     int64      `json:"totalBlanks"`
		TotalBytes      int64      `json:"totalBytes"`
		EstimatedCost   float64    `json:"estimatedCost,omitempty"`
		EstimatedMonths float64    `json:"estimatedMonths,omitempty"`
		EstimatedPeople float64    `json:"estimatedPeople,omitempty"`
	}

	output := jsonOutput{
		Languages:     make([]jsonLang, len(s.Languages)),
		TotalFiles:    s.TotalFiles,
		TotalLines:    s.TotalLines,
		TotalCode:     s.TotalCode,
		TotalComments: s.TotalComments,
		TotalBlanks:   s.TotalBlanks,
		TotalBytes:    s.TotalBytes,
	}

	for i, lang := range s.Languages {
		output.Languages[i] = jsonLang{
			Name:       lang.Name,
			Files:      lang.Files,
			Lines:      lang.Lines,
			Code:       lang.Code,
			Comments:   lang.Comments,
			Blanks:     lang.Blanks,
			Complexity: lang.Complexity,
		}
	}

	if showCocomo {
		output.EstimatedCost = s.EstimatedCost
		output.EstimatedMonths = s.EstimatedMonths
		output.EstimatedPeople = s.EstimatedPeople
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}

func formatNumber(n int64) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	if n < 1000000 {
		return fmt.Sprintf("%d,%03d", n/1000, n%1000)
	}
	return fmt.Sprintf("%d,%03d,%03d", n/1000000, (n/1000)%1000, n%1000)
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d bytes", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-1] + "…"
}
