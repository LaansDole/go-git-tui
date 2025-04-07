package commit

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// View renders the current state of the model
func (m Model) View() string {
	// Show error if any
	if m.Err != nil {
		return m.StyleConfig.ErrorStyle.Render(fmt.Sprintf("Error: %v\nPress any key to exit", m.Err))
	}

	if m.Quitting {
		return "Exiting..."
	}

	// Different help text based on current step
	var helpText string
	if m.Step == 0 {
		helpText = m.StyleConfig.HelpStyle.Render("w/s: Navigate Types • Tab: Select Type • q: Quit")
	} else if m.Step == 1 {
		helpText = m.StyleConfig.HelpStyle.Render("Enter: Commit • Esc: Back • q: Quit")
	} else {
		helpText = m.StyleConfig.HelpStyle.Render("a: Amend Commit • Enter/q: Exit")
	}

	var content string
	switch m.Step {
	case 0:
		// Show commit type selection list in a compact format
		content = lipgloss.JoinVertical(lipgloss.Left, m.TypeList.View())

	case 1:
		// Show commit message input with compact styling
		messageTitle := m.StyleConfig.SubTitleStyle.Render("Enter commit message:")
		messageType := m.StyleConfig.InfoStyle.Render(fmt.Sprintf("Type: %s", m.SelectedType))
		messageInput := m.StyleConfig.InputStyle.Render(m.MessageInput.View())
		instructions := m.StyleConfig.InfoStyle.Render("Press Enter to commit or Esc to go back")

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			messageTitle,
			messageType,
			messageInput,
			instructions,
		)

	case 2:
		// Show confirmation with compact styling
		successTitle := m.StyleConfig.SubTitleStyle.Render("Commit successfully created:")
		commitDetails := m.StyleConfig.SuccessStyle.Render(
			fmt.Sprintf("Type: %s\nMessage: %s", m.SelectedType, m.CommitMessage),
		)
		exitInstructions := m.StyleConfig.InfoStyle.Render("Press any key to exit or 'a' to amend commit message")

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			successTitle,
			commitDetails,
			exitInstructions,
		)
	}

	// Combine all elements with consistent padding and spacing
	verticalLayout := []string{}

	// Add content and help text
	verticalLayout = append(verticalLayout, content, helpText)

	return m.StyleConfig.AppStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			verticalLayout...,
		),
	)
}
