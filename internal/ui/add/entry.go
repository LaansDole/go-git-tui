package add

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Run initializes and runs the add UI component
func Run() error {
	p := tea.NewProgram(New())
	_, err := p.Run()
	return err
}
