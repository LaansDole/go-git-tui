package commit

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// View renders the current state of the model
func (m Model) View() string {
	// Define styling
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	// Show error if any
	if m.Err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v\nPress any key to exit", m.Err))
	}

	if m.Quitting {
		return "Exiting..."
	}

	switch m.Step {
	case 0:
		// Show commit type selection list
		return titleStyle.Render("Select commit type:") + "\n\n" + m.TypeList.View()

	case 1:
		// Show commit message input
		return fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			titleStyle.Render("Enter commit message:"),
			m.MessageInput.View(),
			"Press Enter to commit or Esc to go back",
		)

	case 2:
		// Show confirmation
		commitStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
		return fmt.Sprintf(
			"%s\n\nCommit type: %s\nMessage: %s\n\n%s",
			titleStyle.Render("Commit successfully created:"),
			commitStyle.Render(m.SelectedType),
			commitStyle.Render(m.CommitMessage),
			"Press any key to exit",
		)
	}

	return "Unknown state"
}
