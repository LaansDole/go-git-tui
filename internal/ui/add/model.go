package add

import (
	"fmt"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/LaansDole/go-git-tui/internal/git"
	"github.com/LaansDole/go-git-tui/internal/ui/common"
)

// FileItem represents a git status file item
type FileItem struct {
	Status     string
	Path       string
	IsSelected bool
}

// Title implements the list.Item interface
func (i FileItem) Title() string {
	prefix := "  "
	if i.IsSelected {
		prefix = "✓ "
	}

	// Color code different statuses
	statusStyle := lipgloss.NewStyle()
	switch i.Status {
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
	statusFormatted := statusStyle.Render(fmt.Sprintf("[%s]", i.Status))

	// Use a maximum width for file paths to prevent overflow
	path := common.TruncatePath(i.Path, 60, 30, 27)

	return prefix + statusFormatted + " " + path
}

// Description implements the list.Item interface
func (i FileItem) Description() string { return "" }

// FilterValue implements the list.Item interface
func (i FileItem) FilterValue() string { return i.Path }

// Layout constants for the add UI
const (
	// Layout ratios
	// ListRatio is the percentage of width allocated to the file list
	ListRatio = 40
	// DiffRatio is the percentage of width allocated to the diff view
	DiffRatio = 60

	// Spacing constants for viewport layout
	TitleSpaceReserved   = 1 // Space for title row
	StatsSpaceReserved   = 1 // Space for diff stats
	MessageSpaceReserved = 1 // Space for status messages
	PaddingSpace         = 2 // Extra padding between elements
	SafetyMargin         = 1 // Extra margin to prevent overflow
	DividerWidth         = 1 // Width of the vertical divider
)

// Custom message types
type ErrMsg struct{ error }
type DiffLoadedMsg struct{ Diff *git.DiffResult }
type TickMsg struct{}
type StagingCompleteMsg struct{ Files []string }

// Model represents the main UI model for the git add component
type Model struct {
	// UI Components
	List         list.Model
	DiffViewport viewport.Model

	// State
	Selected       map[int]bool
	Quitting       bool
	CurrentDiff    *git.DiffResult
	CurrentFile    string
	Width          int
	Height         int
	Ready          bool
	LoadingDiff    bool
	Message        string
	MessageTimeout int

	// Dependencies
	GitService  *git.DefaultGitService
	StyleConfig StyleConfig

	// Concurrency control
	diffMutex       sync.Mutex
	lastDiffTime    time.Time
	minDiffDelay    time.Duration
	lastNavTime     time.Time
	navDebounceTime time.Duration
	isNavigating    bool
	lastResizeTime  time.Time
}

// New initializes a new instance of the add UI model
func New() *Model {
	items := []list.Item{}
	var gitService *git.DefaultGitService

	// Get git status using internal/git package
	gitServiceTemp, err := git.NewGitService()
	// Only proceed to get status if service is initialized successfully
	if err == nil {
		gitService = gitServiceTemp
		files, err := gitService.Status()
		if err == nil {
			for _, file := range files {
				items = append(items, FileItem{
					Status:     file.Status,
					Path:       file.Path,
					IsSelected: false,
				})
			}
		}
	}

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
	l.SetShowHelp(false) // Disable the built-in help text since we have our own custom help below
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

	// Create a viewport for the diff display with fixed dimensions
	// This ensures it won't overflow into other UI elements
	diffViewport := viewport.New(0, 0)
	// Set initial viewport height to a reasonable value
	// This will be updated on window resize
	diffViewport.Height = 20
	diffViewport.Style = lipgloss.NewStyle().MaxHeight(20)
	// Disable mouse wheel scrolling to prevent viewport glitches
	diffViewport.MouseWheelEnabled = false
	diffViewport.YPosition = 0
	diffViewport.SetContent("Loading...")

	return &Model{
		List:            l,
		Selected:        make(map[int]bool),
		Quitting:        false,
		DiffViewport:    diffViewport,
		GitService:      gitService,
		StyleConfig:     NewStyleConfig(),
		LoadingDiff:     false,
		lastDiffTime:    time.Now(),
		lastNavTime:     time.Now(),
		minDiffDelay:    250 * time.Millisecond, // Minimum time between diff loads
		navDebounceTime: 500 * time.Millisecond, // Time to wait after navigation stops before loading diff
		lastResizeTime:  time.Now(),
	}
}

// Init initializes the model - implements tea.Model interface
func (m *Model) Init() tea.Cmd {
	return nil
}
