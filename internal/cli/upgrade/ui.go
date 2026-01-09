package upgrade

import (
	"fmt"
	"strings"

	"github.com/KashifKhn/haft/internal/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

var (
	cyan    = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	bold    = lipgloss.NewStyle().Bold(true)
	dim     = lipgloss.NewStyle().Foreground(styles.Muted)
	success = lipgloss.NewStyle().Foreground(styles.Success)
	primary = lipgloss.NewStyle().Foreground(styles.Primary)
)

func printLogo() {
	chevron := cyan.Render("   ▄▄")
	cursor := cyan.Render("█")

	lines := []struct {
		left  string
		text  string
		right string
	}{
		{"   ▄▄", "  ██  ██   ████   ████████  ████████", "█"},
		{"  ▀▀▀▄", " ██  ██  ██  ██  ██           ██", "█"},
		{"     █", " ██████  ██████  █████        ██", "█"},
		{"  ▄▄▄▀", " ██  ██  ██  ██  ██           ██", "█"},
		{"   ▀▀", "  ██  ██  ██  ██  ██           ██", "█"},
	}

	fmt.Println()
	for _, line := range lines {
		chevron = cyan.Render(line.left)
		text := bold.Render(line.text)
		cursor = cyan.Render(line.right)
		fmt.Printf("%s %s      %s\n", chevron, text, cursor)
	}
	fmt.Println()
}

func printTagline() {
	fmt.Println(dim.Render("  The missing Spring Boot CLI"))
	fmt.Println(dim.Render("  scaffolding at terminal speed"))
	fmt.Println()
}

func printHeader(title string) {
	fmt.Println(bold.Render(title))
}

func printStep(symbol, text string) {
	fmt.Printf("%s %s\n", symbol, text)
}

func printStepInProgress(text string) {
	symbol := primary.Render("●")
	printStep(symbol, text)
}

func printStepComplete(text string) {
	symbol := success.Render("◇")
	printStep(symbol, text)
}

func printStepSuccess(text string) {
	symbol := success.Render("✓")
	printStep(symbol, text)
}

func printVersionTransition(from, to string) {
	fromStyled := dim.Render(from)
	arrow := dim.Render("→")
	toStyled := success.Render(to)
	printStepInProgress(fmt.Sprintf("From %s %s %s", fromStyled, arrow, toStyled))
}

func printMethod(method string) {
	printStepInProgress(fmt.Sprintf("Using method: %s", method))
}

func printDone() {
	fmt.Println(success.Render("Done"))
}

func printAlreadyLatest(version string) {
	printStepSuccess(fmt.Sprintf("Already on latest version (%s)", version))
}

func printUpdateAvailable(current, latest string) {
	fmt.Printf("  Current: %s\n", dim.Render(current))
	fmt.Printf("  Latest:  %s\n", success.Render(latest))
	fmt.Println()
	printStepSuccess("Update available! Run 'haft upgrade' to install.")
}

func printProgress(current, total int64, width int) string {
	if total <= 0 {
		return ""
	}

	percent := float64(current) / float64(total)
	filled := int(percent * float64(width))
	empty := width - filled

	bar := strings.Repeat("■", filled) + strings.Repeat("･", empty)
	return fmt.Sprintf("%s %3d%%", primary.Render(bar), int(percent*100))
}

func clearLine() {
	fmt.Print("\r\033[K")
}
