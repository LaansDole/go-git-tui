package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"go-git-tui/internal/git"
)

// FileItem represents a git status file item
type FileItem struct {
	status     string
	path       string
	isSelected bool
}

// Implementation for list.Item interface
func (i FileItem) Title() string {
	prefix := "  "
	if i.isSelected {
		prefix = "✓ "
	}

	// Color code different statuses
	statusStyle := lipgloss.NewStyle()
	switch i.status {
	case "M ":
		statusStyle = statusStyle.Foreground(lipgloss.Color("3")) // Yellow for modified
	case "A ":
		statusStyle = statusStyle.Foreground(lipgloss.Color("2")) // Green for added
	case "D ":
		statusStyle = statusStyle.Foreground(lipgloss.Color("1")) // Red for deleted
	case "??":
		statusStyle = statusStyle.Foreground(lipgloss.Color("4")) // Blue for untracked
	}

	// Format the status in brackets next to the file path
	statusFormatted := statusStyle.Render(fmt.Sprintf("[%s]", i.status))
	return prefix + statusFormatted + " " + i.path
}

func (i FileItem) Description() string {
	// Return empty description since we're now showing status in the title
	return ""
}

func (i FileItem) FilterValue() string { return i.path }

// AddModel for the file selection application
type AddModel struct {
	list     list.Model
	selected map[int]bool
	quitting bool
}

func initialAddModel() AddModel {
	items := []list.Item{}

	// Get git status using internal/git package
	gitService, err := git.NewGitService()
	// Only proceed to get status if service is initialized successfully
	if err == nil {
		files, err := gitService.Status()
		if err == nil {
			for _, file := range files {
				items = append(items, FileItem{
					status:     file.Status,
					path:       file.Path,
					isSelected: false,
				})
			}
		}
	}
	// If there was an error with git service or status, we continue with empty items list

	// Create a custom delegate with more compact spacing
	delegate := list.NewDefaultDelegate()

	// Customize the delegate styles for more compact display
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("170")).
		Margin(0, 0)

	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("240")).
		Margin(0, 0)

	// Reduce padding and margins for all items
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Padding(0, 0).
		Margin(0, 0)

	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.
		Padding(0, 0).
		Margin(0, 0)

	// Set spacing between items to 0
	delegate.SetSpacing(0)

	l := list.New(items, delegate, 0, 0)
	l.Title = "Git Files"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "select"),
			),
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "confirm"),
			),
		}
	}

	return AddModel{
		list:     l,
		selected: make(map[int]bool),
		quitting: false,
	}
}

func (m AddModel) Init() tea.Cmd {
	return nil
}

func (m AddModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "tab":
			idx := m.list.Index()
			m.selected[idx] = !m.selected[idx]

			// Update the item in the list
			if i, ok := m.list.SelectedItem().(FileItem); ok {
				items := []list.Item{}
				for j, item := range m.list.Items() {
					if j == idx {
						i.isSelected = m.selected[idx]
						items = append(items, i)
					} else {
						items = append(items, item)
					}
				}
				m.list.SetItems(items)
			}

			// Move to the next item if not at the end of the list
			if idx < len(m.list.Items())-1 {
				m.list.Select(idx + 1)
			}

			return m, nil

		case "enter":
			// Get all selected items
			var selectedPaths []string
			for i, item := range m.list.Items() {
				if m.selected[i] {
					if fileItem, ok := item.(FileItem); ok {
						selectedPaths = append(selectedPaths, fileItem.path)
					}
				}
			}

			if len(selectedPaths) == 0 {
				return m, tea.Quit
			}

			// Stage files using git package
			err := git.StageFiles(selectedPaths)
			if err != nil {
				return m, tea.Quit
			}

			m.quitting = true
			return m, tea.Sequence(
				tea.Printf("Files staged: %s", strings.Join(selectedPaths, ", ")),
				tea.Quit,
			)
		}

	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m AddModel) View() string {
	if m.quitting {
		return ""
	}

	helpText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("↑/↓: Navigate • Tab: Select/Deselect • Enter: Confirm • q: Quit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.list.View(),
		helpText,
	)
}

func RunAddUI() error {
	p := tea.NewProgram(initialAddModel())
	_, err := p.Run()
	return err
}
