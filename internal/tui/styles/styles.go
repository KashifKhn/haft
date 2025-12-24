package styles

import "github.com/charmbracelet/lipgloss"

var (
	Primary   = lipgloss.Color("12")
	Secondary = lipgloss.Color("5")
	Success   = lipgloss.Color("10")
	Warning   = lipgloss.Color("11")
	Error     = lipgloss.Color("9")
	Muted     = lipgloss.Color("8")
	White     = lipgloss.Color("15")
	Black     = lipgloss.Color("0")
)

var (
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(Primary).
		MarginBottom(1)

	Subtitle = lipgloss.NewStyle().
			Foreground(Muted).
			MarginBottom(1)

	Focused = lipgloss.NewStyle().
		Foreground(Primary).
		Bold(true)

	Blurred = lipgloss.NewStyle().
		Foreground(Muted)

	Cursor = lipgloss.NewStyle().
		Foreground(Primary)

	SuccessText = lipgloss.NewStyle().
			Foreground(Success)

	ErrorText = lipgloss.NewStyle().
			Foreground(Error)

	WarningText = lipgloss.NewStyle().
			Foreground(Warning)

	HelpText = lipgloss.NewStyle().
			Foreground(Muted).
			MarginTop(1)

	Selected = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true)

	Unselected = lipgloss.NewStyle().
			Foreground(White)

	ActiveItem = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)

	InactiveItem = lipgloss.NewStyle().
			Foreground(Muted)

	Border = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Primary).
		Padding(1, 2)

	InputField = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(Muted).
			Padding(0, 1)

	FocusedInput = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(Primary).
			Padding(0, 1)

	ProgressBar = lipgloss.NewStyle().
			Foreground(Primary)

	StepIndicator = lipgloss.NewStyle().
			Foreground(Muted)

	ActiveStep = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)

	CompletedStep = lipgloss.NewStyle().
			Foreground(Success)
)

var (
	CheckMark = SuccessText.Render("✓")
	CrossMark = ErrorText.Render("✗")
	Arrow     = Focused.Render("→")
	Bullet    = Blurred.Render("•")
)

var (
	Subtle = lipgloss.NewStyle().
		Foreground(Muted)

	CategoryStyle = lipgloss.NewStyle().
			Foreground(Secondary).
			Bold(true).
			MarginTop(1)
)

func RenderTitle(text string) string {
	return Title.Render(text)
}

func RenderSubtitle(text string) string {
	return Subtitle.Render(text)
}

func RenderSuccess(text string) string {
	return SuccessText.Render(text)
}

func RenderError(text string) string {
	return ErrorText.Render(text)
}

func RenderWarning(text string) string {
	return WarningText.Render(text)
}

func RenderHelp(text string) string {
	return HelpText.Render(text)
}
