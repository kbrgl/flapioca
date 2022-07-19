package internal

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
	ViewStyle     = lipgloss.NewStyle().MarginTop(1).MarginBottom(2)
	ViewportStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	GameOverStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)
