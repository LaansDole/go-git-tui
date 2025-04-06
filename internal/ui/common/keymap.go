package common

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

// KeyMap defines the key bindings for the UI
type KeyMap struct {
	Select key.Binding
	Cancel key.Binding
}

// NewKeyMap creates a new KeyMap with default key bindings
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

// RenderText is a helper function to render text with a style
func RenderText(text string, style lipgloss.Style) string {
	return style.Render(text)
}
