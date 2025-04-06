package add

import (
	"github.com/charmbracelet/lipgloss"
)

// StyleConfig holds styles for the UI components
type StyleConfig struct {
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

// NewStyleConfig initializes the style configuration with sensible defaults
func NewStyleConfig() StyleConfig {
	return StyleConfig{
		AppStyle:   lipgloss.NewStyle().Margin(0, 1), // Reduced margin
		TitleStyle: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")),
		ListStyle:  lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 0), // Removed padding
		// Fixed dimensions for diff style to prevent content overflow
		DiffStyle:    lipgloss.NewStyle().Padding(0, 0), // Removed padding
		StatusBar:    lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
		HelpStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
		AddedStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("10")), // Green
		DeletedStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("9")),  // Red
		InfoStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("12")), // Blue
		DividerStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("240")), // Add divider style
	}
}
