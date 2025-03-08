package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"go-git-tui/internal/git"
)

var (
	// Styles for the UI components
	titleStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	selectedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	unselectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
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
	return prefix + i.path
}

func (i FileItem) Description() string {
	statusMap := map[string]string{
		"M ": "Modified",
		"A ": "Added",
		"D ": "Deleted",
		"R ": "Renamed",
		"C ": "Copied",
		"U ": "Updated",
		"??": "Untracked",
	}

	desc, ok := statusMap[i.status]
	if !ok {
		desc = i.status
	}
	return desc
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
	files, err := git.GetStatus()
	if err == nil {
		for _, file := range files {
			items = append(items, FileItem{
				status:     file.Status,
				path:       file.Path,
				isSelected: false,
			})
		}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("170"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(lipgloss.Color("240"))

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
