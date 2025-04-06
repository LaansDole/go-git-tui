package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Config holds styles for the UI components
type Config struct {
	AppStyle     lipgloss.Style
	TitleStyle   lipgloss.Style
	ListStyle    lipgloss.Style
	DiffStyle    lipgloss.Style
	StatusBar    lipgloss.Style
	HelpStyle    lipgloss.Style
	AddedStyle   lipgloss.Style
	DeletedStyle lipgloss.Style
	InfoStyle    lipgloss.Style
	DividerStyle lipgloss.Style // Style for the vertical divider between panes
}

// New initializes the style configuration with sensible defaults
func New() Config {
	return Config{
		AppStyle:   lipgloss.NewStyle().Margin(1, 2),
		TitleStyle: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")),
		ListStyle:  lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1),
		// Fixed dimensions for diff style to prevent content overflow
		DiffStyle:    lipgloss.NewStyle().Padding(0, 1),
		StatusBar:    lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
		HelpStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
		AddedStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("10")), // Green
		DeletedStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("9")),  // Red
		InfoStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("12")), // Blue
	}
}
