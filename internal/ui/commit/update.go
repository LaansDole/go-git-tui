package commit

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/LaansDole/go-git-tui/internal/git"
)

// Update handles events and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Skip processing if quitting
	if m.Quitting {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case errMsg:
		// Handle error messages
		m.Err = msg.err
		return m, nil

	case commitSuccessMsg:
		// Handle successful commit
		// Already in step 2 (confirmation)
		return m, nil

	case tea.WindowSizeMsg:
		// Store dimensions and mark as ready
		m.Width, m.Height = msg.Width, msg.Height
		m.Ready = true

		// Set list size with proper margins for compact layout
		listHeight := msg.Height - 7 // Reserve space for title, help, and margins
		m.TypeList.SetSize(msg.Width-4, listHeight)

		// Adjust input field width based on window size
		m.MessageInput.Width = msg.Width - 10

		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.Quitting = true
			return m, tea.Quit

		// w/s navigation for type selection
		case "w":
			if m.Step == 0 {
				currentIndex := m.TypeList.Index()
				if currentIndex > 0 {
					m.TypeList.Select(currentIndex - 1)
				}
				return m, nil
			}

		case "s":
			if m.Step == 0 {
				currentIndex := m.TypeList.Index()
				if currentIndex < len(m.TypeList.Items())-1 {
					m.TypeList.Select(currentIndex + 1)
				}
				return m, nil
			}

		// Use tab to select commit type and proceed to message input
		case "tab":
			if m.Step == 0 {
				currentIndex := m.TypeList.Index()
				if currentIndex >= 0 && currentIndex < len(m.TypeList.Items()) {
					if i, ok := m.TypeList.SelectedItem().(CommitTypeItem); ok {
						// Set the selected type and proceed to message input
						m.SelectedIndex = currentIndex
						m.SelectedType = i.TypeTitle
						m.Step = 1
						m.MessageInput.Focus()
						return m, textinput.Blink
					}
				}
			}
			return m, nil

		case "enter":
			if m.Step == 1 {
				// Get the commit message
				m.CommitMessage = m.MessageInput.Value()
				if m.CommitMessage == "" {
					return m, nil // Don't proceed with empty message
				}
				m.Step = 2

				// Run the commit operation asynchronously
				return m, m.performCommit()

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

		// Handle special keys in step 2 (confirmation)
		if m.Step == 2 {
			switch msg.String() {
			case "a":
				// Amend commit - go back to message input with current message
				m.Step = 1
				m.MessageInput.SetValue(m.CommitMessage)
				m.MessageInput.Focus()
				return m, textinput.Blink
				
			default:
				// Any other key exits
				return m, tea.Quit
			}
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

// performCommit handles the git commit operation
func (m Model) performCommit() tea.Cmd {
	return func() tea.Msg {
		// Use git service for the commit operation
		gitService, err := git.NewGitService()
		if err != nil {
			return errMsg{err}
		}

		err = gitService.Commit(m.SelectedType, m.CommitMessage)
		if err != nil {
			return errMsg{err}
		}

		return commitSuccessMsg{}
	}
}

// Custom message types
type errMsg struct{ err error }
type commitSuccessMsg struct{}
