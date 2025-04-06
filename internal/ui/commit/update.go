package commit

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/LaansDole/go-git-tui/internal/git"
)

// Update handles events and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.Quitting = true
			return m, tea.Quit

		case "enter":
			if m.Step == 0 {
				// Get the selected commit type
				if i, ok := m.TypeList.SelectedItem().(CommitTypeItem); ok {
					m.SelectedType = i.TypeTitle
					m.Step = 1
					m.MessageInput.Focus()
					return m, textinput.Blink
				}
			} else if m.Step == 1 {
				// Get the commit message
				m.CommitMessage = m.MessageInput.Value()
				if m.CommitMessage == "" {
					return m, nil // Don't proceed with empty message
				}
				m.Step = 2

				// Use git service through interface for better decoupling
				var gitService git.GitService
				gitService, err := git.NewGitService()
				if err != nil {
					m.Err = err
					return m, tea.Quit
				}

				err = gitService.Commit(m.SelectedType, m.CommitMessage)
				if err != nil {
					m.Err = err
					return m, tea.Quit
				}

				// Show confirmation and exit after a key press
				return m, nil
			} else if m.Step == 2 {
				// Exit after confirmation
				return m, tea.Quit
			}

		case "esc":
			if m.Step == 1 {
				// Go back to type selection
				m.Step = 0
				m.MessageInput.Blur()
				return m, nil
			}
		}

		// If in step 2 (confirmation), any key exits
		if m.Step == 2 {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		// Adjust the list height and width
		if m.Step == 0 {
			m.TypeList.SetSize(msg.Width, msg.Height-4)
		}
	}

	// Handle updates for the current step components
	if m.Step == 0 {
		m.TypeList, cmd = m.TypeList.Update(msg)
		return m, cmd
	} else if m.Step == 1 {
		m.MessageInput, cmd = m.MessageInput.Update(msg)
		return m, cmd
	}

	return m, nil
}
