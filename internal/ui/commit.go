package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/LaansDole/go-git-tui/internal/git"
)

// CommitTypeItem represents a commit type option
type CommitTypeItem struct {
	title       string
	description string
}

// Implementation for list.Item interface
func (i CommitTypeItem) Title() string       { return i.title }
func (i CommitTypeItem) Description() string { return i.description }
func (i CommitTypeItem) FilterValue() string { return i.title }

// CommitModel represents the commit UI state
type CommitModel struct {
	step          int // 0 = select type, 1 = enter message, 2 = confirm
	typeList      list.Model
	messageInput  textinput.Model
	selectedType  string
	commitMessage string
	quitting      bool
	err           error
}

func initialCommitModel() CommitModel {
	// Setup type selection list
	items := []list.Item{
		CommitTypeItem{title: "feat", description: "A new feature"},
		CommitTypeItem{title: "fix", description: "A bug fix"},
		CommitTypeItem{title: "docs", description: "Documentation changes"},
		CommitTypeItem{title: "chores", description: "Chores and maintenance tasks"},
	}

	typeList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	typeList.Title = "Select Commit Type"

	// Setup message input
	ti := textinput.New()
	ti.Placeholder = "Enter commit message"
	ti.CharLimit = 100
	ti.Width = 50

	return CommitModel{
		step:         0,
		typeList:     typeList,
		messageInput: ti,
		quitting:     false,
	}
}

func (m CommitModel) Init() tea.Cmd {
	return nil
}

func (m CommitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			if m.step == 0 {
				// Get the selected commit type
				if i, ok := m.typeList.SelectedItem().(CommitTypeItem); ok {
					m.selectedType = i.title
					m.step = 1
					m.messageInput.Focus()
					return m, textinput.Blink
				}
			} else if m.step == 1 {
				// Get the commit message
				m.commitMessage = m.messageInput.Value()
				if m.commitMessage == "" {
					return m, nil // Don't proceed with empty message
				}
				m.step = 2

				gitService, err := git.NewGitService()
				if err != nil {
					m.err = err
					return m, tea.Quit
				}

				err = gitService.Commit(m.selectedType, m.commitMessage)
				if err != nil {
					m.err = err
					return m, tea.Quit
				}
				m.quitting = true
				return m, tea.Sequence(
					tea.Printf("Commit successful: %s: %s", m.selectedType, m.commitMessage),
					tea.Quit,
				)
			}
		}

	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		if m.step == 0 {
			m.typeList.SetSize(msg.Width-h, msg.Height-v)
		}
	}

	// Handle updates based on current step
	if m.step == 0 {
		m.typeList, cmd = m.typeList.Update(msg)
		return m, cmd
	} else if m.step == 1 {
		m.messageInput, cmd = m.messageInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m CommitModel) View() string {
	if m.quitting {
		return ""
	}

	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	if m.step == 0 {
		return lipgloss.NewStyle().Margin(1, 2).Render(m.typeList.View())
	} else if m.step == 1 {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("170")).
				Render(fmt.Sprintf("Commit Type: %s", m.selectedType)),
			"",
			"Enter commit message:",
			m.messageInput.View(),
			"",
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Render("Press Enter to commit"),
		)
	}

	return ""
}

func RunCommitUI() error {
	p := tea.NewProgram(initialCommitModel())
	_, err := p.Run()
	return err
}
