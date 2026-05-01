package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	Green  = lipgloss.Color("#34d399")
	Red    = lipgloss.Color("#f87171")
	Yellow = lipgloss.Color("#fbbf24")
	Gray   = lipgloss.Color("#9ca3af")
	White  = lipgloss.Color("#f9fafb")
	Dim    = lipgloss.Color("#6b7280")

	// Text styles
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(White)

	SuccessIcon = lipgloss.NewStyle().
			Foreground(Green).
			SetString("✓")

	FailureIcon = lipgloss.NewStyle().
			Foreground(Red).
			SetString("✗")

	ErrorIcon = lipgloss.NewStyle().
			Foreground(Yellow).
			SetString("!")

	CheckName = lipgloss.NewStyle().
			Foreground(White)

	Message = lipgloss.NewStyle().
		Foreground(Gray)

	Hint = lipgloss.NewStyle().
		Foreground(Dim).
		Italic(true)

	FixLabel = lipgloss.NewStyle().
			Foreground(Yellow).
			Bold(true)

	FixCommand = lipgloss.NewStyle().
			Foreground(Gray)

	SummaryPassed = lipgloss.NewStyle().
			Foreground(Green)

	SummaryFailed = lipgloss.NewStyle().
			Foreground(Red)

	SummaryDot = lipgloss.NewStyle().
			Foreground(Dim).
			SetString("•")
)
