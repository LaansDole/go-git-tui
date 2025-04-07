package commit

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the commit UI state
type Model struct {
	Step          int // 0 = select type, 1 = enter message, 2 = confirm
	TypeList      list.Model
	MessageInput  textinput.Model
	SelectedIndex int
	SelectedType  string
	CommitMessage string
	Quitting      bool
	Width         int
	Height        int
	Ready         bool
	Err           error
	StyleConfig   StyleConfig
}

// New initializes a new commit model
func New() *Model {
	// Setup type selection list
	items := []list.Item{
		CommitTypeItem{
			TypeTitle:       "feat",
			TypeDescription: "A new feature",
		},
		CommitTypeItem{
			TypeTitle:       "fix",
			TypeDescription: "A bug fix",
		},
		CommitTypeItem{
			TypeTitle:       "docs",
			TypeDescription: "Documentation changes",
		},
		CommitTypeItem{
			TypeTitle:       "chore",
			TypeDescription: "Chores and maintenance tasks",
		},
		CommitTypeItem{
			TypeTitle:       "refactor",
			TypeDescription: "Code refactoring without functionality change",
		},
		CommitTypeItem{
			TypeTitle:       "test",
			TypeDescription: "Adding or fixing tests",
		},
		CommitTypeItem{
			TypeTitle:       "style",
			TypeDescription: "Code style/formatting changes",
		},
	}

	// Create a custom delegate with more compact styling
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("170")).
		Margin(0, 0)

	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("240")).
		Margin(0, 0)

	// Reduce padding for items
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Padding(0, 0).
		Margin(0, 0)

	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.
		Padding(0, 0).
		Margin(0, 0)

	// Set spacing between items to 0
	delegate.SetSpacing(0)

	typeList := list.New(items, delegate, 0, 0)
	typeList.Title = "Select commit type:" // We'll use our own title styling
	typeList.SetShowStatusBar(false)
	typeList.SetFilteringEnabled(false)
	typeList.SetShowHelp(false) // Use custom help instead

	// Setup message input with improved styling
	ti := textinput.New()
	ti.Placeholder = "Enter commit message"
	ti.CharLimit = 100
	ti.Width = 50
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))

	return &Model{
		Step:          0,
		TypeList:      typeList,
		SelectedIndex: -1, // No selection initially
		MessageInput:  ti,
		Quitting:      false,
		Ready:         false,
		StyleConfig:   NewStyleConfig(),
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}
