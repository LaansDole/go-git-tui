package commit

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Run initializes and runs the commit UI component in a fullscreen terminal view
func Run() error {
	p := tea.NewProgram(
		New(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := p.Run()
	return err
}
