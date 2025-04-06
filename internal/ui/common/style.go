package common

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

// NewStyleConfig creates a new style configuration with preset styles for all UI components
func NewStyleConfig() StyleConfig {
	return StyleConfig{
		AppStyle:     lipgloss.NewStyle().Margin(1, 2),
		TitleStyle:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")),
		ListStyle:    lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1),
		DiffStyle:    lipgloss.NewStyle().Padding(0, 1),
		StatusBar:    lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
		HelpStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
		AddedStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("10")), // Green
		DeletedStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("9")),  // Red
		InfoStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("12")), // Blue
		DividerStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),  // Gray
	}
}

// Layout constants
const (
	// Layout ratios
	ListRatio = 40 // Percentage of width allocated to the file list
	DiffRatio = 60 // Percentage of width allocated to the diff view

	// Spacing constants for viewport layout
	TitleSpaceReserved   = 1 // Space for title row
	StatsSpaceReserved   = 1 // Space for diff stats
	MessageSpaceReserved = 1 // Space for status messages
	PaddingSpace         = 2 // Extra padding between elements
	SafetyMargin         = 1 // Extra margin to prevent overflow
	DividerWidth         = 1 // Width of the vertical divider
)
