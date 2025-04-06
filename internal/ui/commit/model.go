package commit

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the commit UI state
type Model struct {
	Step          int // 0 = select type, 1 = enter message, 2 = confirm
	TypeList      list.Model
	MessageInput  textinput.Model
	SelectedType  string
	CommitMessage string
	Quitting      bool
	Err           error
}

// New initializes a new commit model
func New() Model {
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
			TypeTitle:       "chores",
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

	typeList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	typeList.Title = "Select Commit Type"

	// Setup message input
	ti := textinput.New()
	ti.Placeholder = "Enter commit message"
	ti.CharLimit = 100
	ti.Width = 50

	return Model{
		Step:         0,
		TypeList:     typeList,
		MessageInput: ti,
		Quitting:     false,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}
