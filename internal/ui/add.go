package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/LaansDole/go-git-tui/internal/git"
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
	
	// Use a maximum width for file paths to prevent overflow
	path := truncatePath(i.path, 60, 30, 27)
	
	return prefix + statusFormatted + " " + path
}

func (i FileItem) Description() string {
	// Return empty description since we're now showing status in the title
	return ""
}

func (i FileItem) FilterValue() string { return i.path }

// Layout constants
const (
	// Layout ratios
	// ListRatio is the percentage of width allocated to the file list
	ListRatio = 40
	// DiffRatio is the percentage of width allocated to the diff view
	DiffRatio = 60
	
	// Spacing constants for viewport layout
	TitleSpaceReserved    = 1 // Space for title row
	StatsSpaceReserved    = 1 // Space for diff stats
	MessageSpaceReserved  = 1 // Space for status messages
	PaddingSpace          = 2 // Extra padding between elements
	SafetyMargin          = 1 // Extra margin to prevent overflow
	DividerWidth          = 1 // Width of the vertical divider
)

// AddModel for the file selection application
type AddModel struct {
	list           list.Model           // File list component
	selected       map[int]bool         // Tracks selected files by index
	quitting       bool                 // Whether the application is exiting
	diffViewport   viewport.Model       // Viewport for showing diff content
	currentDiff    *git.DiffResult      // Currently loaded diff
	currentFile    string               // Currently selected file path
	width, height  int                  // Terminal dimensions
	gitService     *git.DefaultGitService // Git service for operations
	ready          bool                 // Whether the UI is initialized
	styleConfig    StyleConfig          // Styling configuration
	loadingDiff    bool                 // Whether we're loading a diff
	message        string               // Status message to display
	messageTimeout int                  // Countdown for message display
}

// StyleConfig holds styles for the UI components
type StyleConfig struct {
	AppStyle     lipgloss.Style
	TitleStyle   lipgloss.Style
	ListStyle    lipgloss.Style
	DiffStyle    lipgloss.Style
	StatusBar    lipgloss.Style
	HelpStyle    lipgloss.Style
	AddedStyle   lipgloss.Style
	DeletedStyle lipgloss.Style
	InfoStyle    lipgloss.Style
	DividerStyle lipgloss.Style // Style for the vertical divider between panes
}

// Initialize the style configuration with sensible defaults
// Initialize the style configuration with sensible defaults and fixed dimensions
// Create a new style configuration with preset styles for all UI components
func newStyleConfig() StyleConfig {
	return StyleConfig{
		AppStyle:   lipgloss.NewStyle().Margin(1, 2),
		TitleStyle: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")),
		ListStyle:  lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1),
		// Fixed dimensions for diff style to prevent content overflow
		DiffStyle:  lipgloss.NewStyle().Padding(0, 1),
		StatusBar:  lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
		HelpStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
		AddedStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("10")), // Green
		DeletedStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("9")), // Red
		InfoStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("12")), // Blue
	}
}

// Initialize the AddModel with default values
func initialAddModel() AddModel {
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

	// Create a viewport for the diff display with fixed dimensions
	// This ensures it won't overflow into other UI elements
	diffViewport := viewport.New(0, 0)
	// Set initial viewport height to a reasonable value
	// This will be updated on window resize
	diffViewport.Height = 20
	diffViewport.Style = lipgloss.NewStyle().MaxHeight(20)
	// Enable mouse wheel scrolling
	diffViewport.MouseWheelEnabled = true
	diffViewport.YPosition = 0
	diffViewport.SetContent("Loading...")

	return AddModel{
		list:         l,
		selected:     make(map[int]bool),
		quitting:     false,
		diffViewport: diffViewport,
		gitService:   gitService,
		styleConfig:  newStyleConfig(),
		loadingDiff:  false,
		message:      "",
	}
}

func (m AddModel) Init() tea.Cmd {
	// Return the initial command to set up the UI
	return nil
}

// showDiff loads and displays the diff for the selected file
func (m *AddModel) showDiff(filePath string) tea.Cmd {
	return func() tea.Msg {
		if m.gitService == nil {
			return nil
		}

		diff, err := m.gitService.GetFileDiff(filePath)
		if err != nil {
			return errMsg{err}
		}

		return diffLoadedMsg{diff}
	}
}

// formatDiffContent formats the diff content with syntax highlighting
func (m *AddModel) formatDiffContent(diff *git.DiffResult) string {
	if diff == nil {
		return "No diff available"
	}

	if diff.IsBinary {
		return "Binary file differences not shown"
	}

	var result strings.Builder
	lines := strings.Split(diff.Content, "\n")
	
	// Calculate max line width based on viewport width
	maxWidth := m.diffViewport.Width - 2 // Account for padding
	if maxWidth < 20 { // Ensure minimum readable width
		maxWidth = 20
	}
	
	for _, line := range lines {
		if len(line) > 0 {
			// Truncate very long lines to prevent overflow using helper function
			displayLine := truncateText(line, maxWidth-5, "...")
			
			prefix := line[0:1]
			switch prefix {
			case "+":
				result.WriteString(m.styleConfig.AddedStyle.Render(displayLine) + "\n")
			case "-":
				result.WriteString(m.styleConfig.DeletedStyle.Render(displayLine) + "\n")
			default:
				result.WriteString(displayLine + "\n")
			}
		} else {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// getDiffStats returns a formatted string with diff statistics
func (m *AddModel) getDiffStats(diff *git.DiffResult) string {
	// Calculate total changes (sum of added and deleted lines)
	if diff == nil || (diff.Stats.Added == 0 && diff.Stats.Deleted == 0) {
		return "No changes"
	}

	return fmt.Sprintf(
		"%d insertions(+), %d deletions(-)",
		diff.Stats.Added,
		diff.Stats.Deleted,
	)
}

// Custom message types
type errMsg struct{ error }
type diffLoadedMsg struct{ diff *git.DiffResult }
type tickMsg struct{}
type stagingCompleteMsg struct{ files []string }

func (m AddModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.ready = true

		// Get frame sizes from styles
		appFrameH, appFrameV := m.styleConfig.AppStyle.GetFrameSize()
		listFrameH, listFrameV := m.styleConfig.ListStyle.GetFrameSize()
		diffFrameH, diffFrameV := m.styleConfig.DiffStyle.GetFrameSize()

		// Calculate sizes for side-by-side layout
		availableWidth := msg.Width - appFrameH
		availableHeight := msg.Height - appFrameV - 4 // Reserve space for title, status, help
		
		// Split width 40% for list, 60% for diff (like in gadd.sh)
		listWidth := availableWidth * ListRatio / 100
		diffWidth := availableWidth - listWidth
		
		// Calculate optimal component sizes with fixed dimensions and clear separation
		// Using the global space allocation constants
		
		// Reserve space for header elements
		reservedVerticalSpace := TitleSpaceReserved + StatsSpaceReserved + MessageSpaceReserved + PaddingSpace + SafetyMargin
		
		// Set fixed dimensions for file list
		m.list.SetSize(listWidth - listFrameH - DividerWidth, availableHeight - listFrameV)
		
		// Set fixed dimensions for diff viewport that won't overflow
		m.diffViewport.Width = diffWidth - diffFrameH
		m.diffViewport.Height = availableHeight - diffFrameV - reservedVerticalSpace
		
		// Completely disable viewport overlapping with overflow handling
		m.styleConfig.DiffStyle = m.styleConfig.DiffStyle.Copy().Height(m.diffViewport.Height).MaxHeight(m.diffViewport.Height)
		
		// Create a vertical divider style
		m.styleConfig.DividerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Width(DividerWidth)

		// If there's at least one item, load the diff for the first item
		if len(m.list.Items()) > 0 {
			if i, ok := m.list.SelectedItem().(FileItem); ok && m.currentFile == "" {
				m.currentFile = i.path
				m.loadingDiff = true
				return m, m.showDiff(i.path)
			}
		}

	case diffLoadedMsg:
		m.currentDiff = msg.diff
		content := m.formatDiffContent(msg.diff)
		m.diffViewport.SetContent(content)
		m.loadingDiff = false
		m.message = ""
		m.messageTimeout = 0
		return m, nil

	case tickMsg:
		// Handle message timeout
		if m.messageTimeout > 0 {
			m.messageTimeout--
			if m.messageTimeout == 0 {
				m.message = ""
			}
			return m, tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg {
				return tickMsg{}
			})
		}
		return m, nil

	case stagingCompleteMsg:
		// Display a confirmation message showing which files were staged
		m.message = fmt.Sprintf("%d files staged: %s", len(msg.files), strings.Join(msg.files, ", "))
		m.messageTimeout = 20  // Show message for a longer time
		
		// Set a delayed quit command to show the message before quitting
		return m, tea.Sequence(
			tea.Tick(1500*time.Millisecond, func(time.Time) tea.Msg { return tickMsg{} }),
			tea.Quit,
		)

	case errMsg:
		// Display error message and quit
		return m, tea.Sequence(
			tea.Printf("Error: %v", msg.error),
			tea.Quit,
		)

	case tea.KeyMsg:
		// Handle key commands
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "tab", " ":
			// Toggle selection for the current file
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

			// Move to the next item if not at the end of the list when using tab
			if idx < len(m.list.Items())-1 && msg.String() == "tab" {
				m.list.Select(idx + 1)
				// Load diff for the next item automatically
				if nxt, ok := m.list.Items()[idx+1].(FileItem); ok {
					m.currentFile = nxt.path
					m.loadingDiff = true
					m.message = "Loading diff..."
					m.messageTimeout = 10
					return m, m.showDiff(nxt.path)
				}
			}

			return m, nil
			
		case "enter":
			// Confirm staging with Enter
			// Don't quit immediately - show which files were staged first
			return m, m.confirmStaging()
			
		case "j":
			// Scroll down one line
			m.diffViewport.LineDown(1)
			return m, nil
			
		case "k":
			// Scroll up one line
			m.diffViewport.LineUp(1)
			return m, nil
			
		case "d":
			// Scroll down half page (like vim ctrl+d)
			m.diffViewport.HalfViewDown()
			return m, nil
			
		case "u":
			// Scroll up half page (like vim ctrl+u)
			m.diffViewport.HalfViewUp()
			return m, nil
			
		case "g":
			// Scroll to top (like vim)
			m.diffViewport.GotoTop()
			return m, nil
			
		case "G":
			// Scroll to bottom (like vim)
			m.diffViewport.GotoBottom()
			return m, nil
		}
	}

	// Handle navigation keys for separate viewports
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Diff viewport navigation keys (already handled above)
		case "j", "k", "u", "d", "g", "G":
			// These keys only affect the diff viewport and have been handled above
			// Don't pass them to the list at all
			break
			
		// File list navigation keys
		case "up", "down", "home", "end", "pgup", "pgdown":
			// These keys only affect the file list navigation
			var prevIndex = m.list.Index()
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
			
			// If the cursor moved and we have items, update the diff view
			if prevIndex != m.list.Index() && len(m.list.Items()) > 0 {
				if item, ok := m.list.SelectedItem().(FileItem); ok {
					m.currentFile = item.path
					m.loadingDiff = true
					m.message = "Loading diff..."
					m.messageTimeout = 10
					cmds = append(cmds, m.showDiff(item.path))
				}
			}
			
		// All other keys (selection, quit, etc.)
		default:
			// Pass other keys to list
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		}
		
	default:
		// Pass non-key messages (window resize, etc.) to list
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}
	
	// Tick for message timeout
	if m.messageTimeout > 0 {
		cmds = append(cmds, tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg {
			return tickMsg{}
		}))
	}
	
	return m, tea.Batch(cmds...)
}

// confirmStaging stages the selected files and exits
func (m *AddModel) confirmStaging() tea.Cmd {
	return func() tea.Msg {
		// Get all selected items
		var selectedPaths []string
		for i, item := range m.list.Items() {
			if m.selected[i] {
				if fileItem, ok := item.(FileItem); ok {
					selectedPaths = append(selectedPaths, fileItem.path)
				}
			}
		}

		// If nothing selected but we're in diff view, select the current file
		if len(selectedPaths) == 0 && m.currentFile != "" {
			selectedPaths = []string{m.currentFile}
		}

		// If still nothing to stage, just quit
		if len(selectedPaths) == 0 {
			return nil
		}

		// Stage the files
		gitService, err := git.NewGitService()
		if err != nil {
			return errMsg{err}
		}

		err = gitService.Stage(selectedPaths)
		if err != nil {
			return errMsg{err}
		}

		// Return a successful staging message
		return stagingCompleteMsg{files: selectedPaths}
	}
}

// Helper function to truncate text with proper bounds checking
func truncateText(text string, maxLength int, ellipsis string) string {
	// If text is already short enough, return it unchanged
	if len(text) <= maxLength || maxLength <= 0 {
		return text
	}
	
	// Ensure we have at least one character plus ellipsis
	truncateAt := maxLength
	if truncateAt <= len(ellipsis) {
		truncateAt = 1
	} else {
		truncateAt -= len(ellipsis)
	}
	
	// Make sure we don't exceed the string length
	if truncateAt > len(text) {
		truncateAt = len(text)
	}
	
	// Return truncated text with ellipsis
	return text[:truncateAt] + ellipsis
}

// Helper function for truncating paths with middle ellipsis
func truncatePath(path string, maxLength int, prefixChars int, suffixChars int) string {
	// If path is already short enough, return it unchanged
	if len(path) <= maxLength {
		return path
	}
	
	// Ensure we have reasonable prefix and suffix sizes
	const ellipsis = "..."
	
	// Make sure prefix + suffix doesn't exceed path length
	if prefixChars + suffixChars > len(path) {
		// Adjust proportionally
		total := len(path)
		prefixChars = total / 2
		suffixChars = total - prefixChars - len(ellipsis)
		if suffixChars < 0 {
			suffixChars = 0
		}
	}
	
	// Construct truncated path
	if len(path) > prefixChars + suffixChars + len(ellipsis) {
		return path[:prefixChars] + ellipsis + path[len(path)-suffixChars:]
	}
	
	// Fallback if calculations are off
	return path
}

func (m AddModel) View() string {
	if m.quitting {
		return ""
	}

	if !m.ready {
		return "Loading git repository..."
	}

	// Show message when there are no files to stage
	if len(m.list.Items()) == 0 {
		return m.styleConfig.InfoStyle.Render("No files to stage. Working tree clean.")
	}

	// Always show side-by-side view
	// Add a title and status line
	titleText := m.styleConfig.TitleStyle.Render("Go Git TUI - Stage Files")
	
	// Count selected items
	selectedCount := 0
	for _, selected := range m.selected {
		if selected {
			selectedCount++
		}
	}
	
	statusText := m.styleConfig.StatusBar.Render(
		fmt.Sprintf("%d files, %d selected", len(m.list.Items()), selectedCount))

	helpText := m.styleConfig.HelpStyle.Render(
		"↑↓: Navigate Files • j/k: Scroll Diff • Tab: Select • u/d: Page Up/Down • g/G: Top/Bottom • Enter: Confirm • q: Quit")

	// Prepare diff view
	diffTitle := "Diff"
	diffStats := ""
	if m.currentFile != "" {
		diffTitle = fmt.Sprintf("Diff for %s", m.currentFile)
		if m.currentDiff != nil {
			diffStats = m.getDiffStats(m.currentDiff)
		}
	}

	// Wrap in styles
	diffTitleText := m.styleConfig.TitleStyle.Render(diffTitle)
	diffStatsText := ""
	if diffStats != "" {
		diffStatsText = m.styleConfig.InfoStyle.Render(diffStats)
	}

	// Message handling moved to the bottom of the layout
	// for better visibility across the entire UI

	// Prepare the content for the diff viewport
	diffContent := ""
	if m.loadingDiff {
		diffContent = "Loading diff..."
	} else if m.currentDiff == nil {
		diffContent = "Select a file to view diff"
	} else {
		// Get the viewport content directly
		diffContent = m.diffViewport.View()
	}

	// Build the split view
	listView := m.styleConfig.ListStyle.Render(m.list.View())
	
	// Calculate available width for the diff panel
	diffWidth := m.width / 2
	
	// Truncate title if too long using helper function
	if len(diffTitle) > diffWidth - 5 && diffWidth > 15 {
		diffTitle = truncateText(diffTitle, diffWidth-10, "...")
		diffTitleText = m.styleConfig.TitleStyle.Render(diffTitle)
	}
	
	// Build the diff panel with fixed dimensions and strict height constraint
	diffPanel := lipgloss.JoinVertical(
		lipgloss.Left,
		diffTitleText,
		diffStatsText,
		// Apply style with fixed height and width to ensure containment
		// Use Height method to constrain the diff content
		m.styleConfig.DiffStyle.Height(m.diffViewport.Height).Width(diffWidth - 4).Render(diffContent),
	)

	// Create the vertical divider with appropriate height
	dividerHeight := m.diffViewport.Height + TitleSpaceReserved + StatsSpaceReserved + MessageSpaceReserved
	verticalDivider := strings.Repeat("│\n", dividerHeight)
	dividerView := m.styleConfig.DividerStyle.Render(verticalDivider)
	
	// Combine in a side-by-side layout with strict separation
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		listView,
		dividerView,
		diffPanel,
	)

	// Add title, status, message (for staging confirmations), and help below
	
	// Create a properly styled message text that spans the full width
	messageDisplay := ""
	if m.message != "" {
		// Use a style that stands out better for important confirmation messages
		messageStyle := m.styleConfig.InfoStyle.Copy().Padding(0, 1).Bold(true).Width(m.width)
		messageDisplay = messageStyle.Render(m.message)
	}
	
	content = lipgloss.JoinVertical(
		lipgloss.Left,
		titleText,
		statusText,
		content,
		messageDisplay,
		helpText,
	)

	return m.styleConfig.AppStyle.Render(content)
}

func RunAddUI() error {
	p := tea.NewProgram(initialAddModel())
	_, err := p.Run()
	return err
}
