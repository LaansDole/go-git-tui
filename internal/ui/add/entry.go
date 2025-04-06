package add

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Run initializes and runs the add UI component in a fullscreen terminal view
func Run() error {
	p := tea.NewProgram(
		New(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := p.Run()
	return err
}
