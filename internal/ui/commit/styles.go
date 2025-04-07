package commit

import (
	"github.com/charmbracelet/lipgloss"
)

// StyleConfig holds style definitions for the commit UI
type StyleConfig struct {
	AppStyle      lipgloss.Style
	TitleStyle    lipgloss.Style
	SubTitleStyle lipgloss.Style
	ListStyle     lipgloss.Style
	InputStyle    lipgloss.Style
	ErrorStyle    lipgloss.Style
	SuccessStyle  lipgloss.Style
	InfoStyle     lipgloss.Style
	HelpStyle     lipgloss.Style
	StatusBar     lipgloss.Style
}

// NewStyleConfig creates and initializes style configuration
func NewStyleConfig() StyleConfig {
	appStyle := lipgloss.NewStyle().
		Padding(1, 2, 1, 2)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		MarginBottom(1)

	subTitleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63")).
		MarginBottom(1)

	listStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		PaddingRight(1)

	inputStyle := lipgloss.NewStyle().
		PaddingLeft(1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		MarginTop(1).
		MarginBottom(1).
		Width(60)

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Bold(true)

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Bold(true).
		Padding(1, 0)

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("14")).
		Italic(true)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(1)

	statusBar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("8")).
		Bold(true).
		Padding(0, 1).
		MarginBottom(1)

	return StyleConfig{
		AppStyle:      appStyle,
		TitleStyle:    titleStyle,
		SubTitleStyle: subTitleStyle,
		ListStyle:     listStyle,
		InputStyle:    inputStyle,
		ErrorStyle:    errorStyle,
		SuccessStyle:  successStyle,
		InfoStyle:     infoStyle,
		HelpStyle:     helpStyle,
		StatusBar:     statusBar,
	}
}
