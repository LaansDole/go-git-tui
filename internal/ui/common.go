package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

// Style for the UI components
var (
	primaryStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	secondaryStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

// KeyMap defines the key bindings for the UI
type KeyMap struct {
	Select key.Binding
	Cancel key.Binding
}

// NewKeyMap initializes the key bindings
func NewKeyMap() KeyMap {
	return KeyMap{
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("ctrl+c", "esc"),
			key.WithHelp("ctrl+c/esc", "cancel"),
		),
	}
}

// RenderText renders styled text
func RenderText(text string, style lipgloss.Style) string {
	return style.Render(text)
}
