package add

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/LaansDole/go-git-tui/internal/git"
)

// DelayedDiffMsg is sent after the navigation debounce period to trigger a diff load
type DelayedDiffMsg struct{ 
    FilePath string
    RequestID int64 // Unique ID to track requests and avoid race conditions
}

// DelayedResizeMsg is sent after the resize debounce period to finalize a resize operation
type DelayedResizeMsg struct{ 
    Width, Height int
    RequestID int64 // Unique ID to track requests
}

// Update handles state changes based on messages - implements tea.Model interface
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Safety check - avoid processing messages if we're quitting
	if m.Quitting {
		return m, tea.Quit
	}

	// Limited message processing while loading diff
	if m.LoadingDiff {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "q", "ctrl+c", "esc":
				m.Quitting = true
				m.LoadingDiff = false
				return m, tea.Quit
			default:
				// Allow only escape keys while loading diff
				return m, nil
			}
		}
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Skip resize processing if we're quitting
		if m.Quitting {
			return m, nil
		}
		
		// Handle window resize with debouncing
		currentTime := time.Now()
		m.lastResizeTime = currentTime
		
		// Only set resize in progress if not already in progress
		if !m.resizeInProgress {
			m.resizeInProgress = true
		}

		// Update dimensions immediately but debounce the full resize operation
		m.Width, m.Height = msg.Width, msg.Height
		
		// Generate unique ID for this resize request
		requestID := time.Now().UnixNano()
		
		// Schedule a delayed resize operation
		return m, tea.Tick(m.resizeDebounce, func(t time.Time) tea.Msg {
			// Only proceed if this is the most recent resize request
			if t.Sub(m.lastResizeTime) >= m.resizeDebounce {
				return DelayedResizeMsg{
					Width: msg.Width, 
					Height: msg.Height,
					RequestID: requestID,
				}
			}
			return nil
		})

	case DiffLoadedMsg:
		// Handle diff loaded message
		return m.handleDiffLoaded(msg)

	case TickMsg:
		// Handle message timeout
		return m.handleTick()

	case DelayedDiffMsg:
		// Only load the diff if we're not currently navigating and not quitting
		// Also ensure we're not in the middle of a resize operation
		if !m.isNavigating && !m.Quitting && !m.resizeInProgress && msg.FilePath == m.CurrentFile {
			// Safely trigger the diff loading
			m.LoadingDiff = true
			return m, m.ShowDiff(msg.FilePath)
		}
		return m, nil
		
	case DelayedResizeMsg:
		// Skip if we're quitting
		if m.Quitting {
			m.resizeInProgress = false
			return m, nil
		}
		
		// Process the resize only if not loading diff
		if !m.LoadingDiff {
			// Handle a delayed resize event after debouncing
			defer func() { m.resizeInProgress = false }()
			return m.handleWindowResize(tea.WindowSizeMsg{Width: msg.Width, Height: msg.Height})
		}
		
		// If loading diff, mark resize as not in progress but don't process yet
		m.resizeInProgress = false
		return m, nil

	case StagingCompleteMsg:
		// Handle staging complete message
		return m.handleStagingComplete(msg)

	case ErrMsg:
		// Display error message and quit
		return m, tea.Sequence(
			tea.Printf("Error: %v", msg.error),
			tea.Quit,
		)

	case tea.KeyMsg:
		// Handle various keyboard commands
		switch msg.String() {
		// Explicit handling for exit keys
		case "q", "ctrl+c", "esc":
			// Quit the application
			m.Quitting = true
			return m, tea.Quit

		// Explicit handling for selection keys
		case "tab", " ":
			// Toggle selection for the current file
			return m.handleSelectionToggle(msg)

		case "enter":
			// Confirm staging with Enter
			return m, m.ConfirmStaging()

		// Explicit handling for diff viewport navigation
		case "j":
			m.DiffViewport.LineDown(1)
			return m, nil

		case "k":
			m.DiffViewport.LineUp(1)
			return m, nil

		case "d":
			// Toggle diff display for current file
			if i, ok := m.List.SelectedItem().(FileItem); ok {
				// If current diff is already showing this file, clear it
				if m.CurrentFile == i.Path && m.CurrentDiff != nil {
					m.CurrentDiff = nil
					m.DiffViewport.SetContent("")
					m.Message = "Diff hidden"
					m.MessageTimeout = 5
					return m, nil
				} else {
					// Otherwise, show diff for this file
					m.CurrentFile = i.Path
					m.LoadingDiff = true
					m.Message = "Loading diff..."
					m.MessageTimeout = 10
					return m, m.ShowDiff(i.Path)
				}
			}
			return m, nil

		case "g":
			// Scroll to top (like vim)
			m.DiffViewport.GotoTop()
			return m, nil

		case "G":
			// Scroll to bottom (like vim)
			m.DiffViewport.GotoBottom()
			return m, nil

		// Explicit handling for left/right arrow keys to prevent unexpected exits
		case "left", "right":
			// Intentionally ignore these keys to prevent unexpected exits
			return m, nil
		}
	}

	// Handle navigation keys for separate viewports
	return m.handleNavigationKeys(msg)
}

// handleWindowResize handles window resize messages
func (m *Model) handleWindowResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	// Store dimensions
	m.Width, m.Height = msg.Width, msg.Height
	m.Ready = true

	// Get frame sizes from styles
	appFrameH, appFrameV := m.StyleConfig.AppStyle.GetFrameSize()
	listFrameH, listFrameV := m.StyleConfig.ListStyle.GetFrameSize()
	diffFrameH, diffFrameV := m.StyleConfig.DiffStyle.GetFrameSize()

	// Calculate sizes with more compact layout
	availableWidth := m.Width - appFrameH 
	availableHeight := m.Height - appFrameV - 3 // Reduced reservation for title, status, help

	// Split width 40% for list, 60% for diff
	listWidth := availableWidth * ListRatio / 100
	diffWidth := availableWidth - listWidth - DividerWidth - 1 // Add extra space buffer

	// Reserve less vertical space for more compact layout
	reservedVerticalSpace := TitleSpaceReserved + StatsSpaceReserved + MessageSpaceReserved + 1

	// Update file list dimensions
	m.List.SetSize(listWidth-listFrameH, availableHeight-listFrameV)

	// Calculate new viewport dimensions
	newViewportWidth := diffWidth - diffFrameH
	newViewportHeight := availableHeight - diffFrameV - reservedVerticalSpace
	
	// Only fully recreate the viewport if dimensions have changed
	if m.DiffViewport.Width != newViewportWidth || m.DiffViewport.Height != newViewportHeight {
		// Safe viewport recreation with multiple safety checks
		if m.isNavigating || m.LoadingDiff {
			// Delay viewport recreation if we're in the middle of another operation
			// This prevents crashes from race conditions
			return m, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
				return DelayedResizeMsg{
					Width: m.Width,
					Height: m.Height,
					RequestID: time.Now().UnixNano(),
				}
			})
		}
		
		// Save current state
		currentYPosition := m.DiffViewport.YPosition
		var content string
		
		// Get content to restore
		if m.CurrentDiff != nil {
			content = m.FormatDiffContent(m.CurrentDiff)
		} else {
			content = m.DiffViewport.View()
		}
		
		// Completely recreate the viewport to avoid artifacts
		m.DiffViewport = viewport.New(newViewportWidth, newViewportHeight)
		m.DiffViewport.MouseWheelEnabled = false
		m.DiffViewport.SetContent(content)
		
		// Restore scroll position if possible
		if currentYPosition > 0 && content != "" {
			m.DiffViewport.YPosition = currentYPosition
			// Ensure we don't scroll beyond content
			if m.DiffViewport.YPosition > m.DiffViewport.TotalLineCount() - m.DiffViewport.Height {
				if m.DiffViewport.TotalLineCount() - m.DiffViewport.Height > 0 {
					m.DiffViewport.YPosition = m.DiffViewport.TotalLineCount() - m.DiffViewport.Height
				} else {
					m.DiffViewport.YPosition = 0
				}
			}
		}
	}

	// Set fixed height for diff style to prevent overflow
	m.StyleConfig.DiffStyle = m.StyleConfig.DiffStyle.Copy().
		Height(m.DiffViewport.Height).
		MaxHeight(m.DiffViewport.Height).
		Width(m.DiffViewport.Width).
		MaxWidth(m.DiffViewport.Width)
	
	// Reset the divider style with the proper color and size
	m.StyleConfig.DividerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Width(DividerWidth)

	// If there's at least one item, set the current file on first load
	if len(m.List.Items()) > 0 && m.CurrentFile == "" {
		if i, ok := m.List.SelectedItem().(FileItem); ok {
			m.CurrentFile = i.Path
			m.DiffViewport.SetContent("Select a file and press TAB to view the diff")
		}
	}

	// Return a command to perform another frame draw after a slight delay
	// to ensure any artifacts are cleaned up
	return m, tea.Tick(50*time.Millisecond, func(_ time.Time) tea.Msg {
		// This causes a redraw without changing state
		return nil
	})
}

// handleDiffLoaded handles when a diff is loaded
func (m *Model) handleDiffLoaded(msg DiffLoadedMsg) (tea.Model, tea.Cmd) {
	if m.Quitting {
		return m, nil
	}

	m.CurrentDiff = msg.Diff
	content := m.FormatDiffContent(msg.Diff)
	m.DiffViewport.SetContent(content)
	m.LoadingDiff = false
	m.Message = ""
	m.MessageTimeout = 0
	return m, nil
}

// handleTick handles tick messages for message timeout
func (m *Model) handleTick() (tea.Model, tea.Cmd) {
	if m.MessageTimeout > 0 {
		m.MessageTimeout--
		if m.MessageTimeout == 0 {
			m.Message = ""
		}
		return m, tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg {
			return TickMsg{}
		})
	}
	return m, nil
}

// handleStagingComplete handles when staging is complete
func (m *Model) handleStagingComplete(msg StagingCompleteMsg) (tea.Model, tea.Cmd) {
	m.Message = fmt.Sprintf("%d files staged: %s", len(msg.Files), strings.Join(msg.Files, ", "))
	m.MessageTimeout = 20

	return m, tea.Sequence(
		tea.Tick(1500*time.Millisecond, func(time.Time) tea.Msg { return TickMsg{} }),
		tea.Quit,
	)
}

// handleSelectionToggle handles toggling selection of files
func (m *Model) handleSelectionToggle(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	idx := m.List.Index()
	m.Selected[idx] = !m.Selected[idx]

	// Update the item in the list
	if i, ok := m.List.SelectedItem().(FileItem); ok {
		items := []list.Item{}
		for j, item := range m.List.Items() {
			if j == idx {
				i.IsSelected = m.Selected[idx]
				items = append(items, i)
			} else {
				items = append(items, item)
			}
		}
		m.List.SetItems(items)
	}

	if _, ok := m.List.SelectedItem().(FileItem); ok && msg.String() == "tab" {
		isLastFile := m.List.Index() == len(m.List.Items())-1

		if !isLastFile {
			m.List.Select(m.List.Index() + 1)
		}

		return m, nil
	}

	return m, nil
}

// handleNavigationKeys handles navigation keys for both viewports
func (m *Model) handleNavigationKeys(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Safety checks - avoid processing if program is in certain states
	if m.Quitting || m.resizeInProgress {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "k", "g", "G", "left", "right", "up", "down":
			break

		case "w", "s", "home", "end", "pgup", "pgdown":
			if len(m.List.Items()) == 0 {
				break
			}

			var prevIndex = m.List.Index()

			isLastFile := m.List.Index() == len(m.List.Items())-1

			var keyToSend tea.KeyMsg
			switch msg.String() {
			case "w":
				if m.List.Index() <= 0 {
					break
				}
				keyToSend = tea.KeyMsg{Type: tea.KeyUp}
			case "s":
				if isLastFile {
					break
				}
				keyToSend = tea.KeyMsg{Type: tea.KeyDown}
			default:
				keyToSend = msg
			}

			var cmd tea.Cmd
			if keyToSend.Type != 0 {
				m.List, cmd = m.List.Update(keyToSend)
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			}

			if prevIndex != m.List.Index() && len(m.List.Items()) > 0 {
				if item, ok := m.List.SelectedItem().(FileItem); ok {
					m.CurrentFile = item.Path

					// Handle w/s navigation debouncing
					if msg.String() == "w" || msg.String() == "s" {
						// Set navigating state
						m.isNavigating = true
						m.lastNavTime = time.Now()

						// Generate unique request ID for this navigation request
						requestID := time.Now().UnixNano()

						// Start a debounced diff loading with unique ID
						cmds = append(cmds, tea.Tick(m.navDebounceTime, func(t time.Time) tea.Msg {
							// Only load the diff if no other navigation has happened since this timer started
							if time.Since(m.lastNavTime) >= m.navDebounceTime && !m.Quitting && !m.resizeInProgress {
								m.isNavigating = false
								return DelayedDiffMsg{
									FilePath: item.Path,
									RequestID: requestID,
								}
							}
							return nil
						}))
					}
				}
			}
		default:
			var cmd tea.Cmd
			m.List, cmd = m.List.Update(msg)
			cmds = append(cmds, cmd)
		}

	default:
		var cmd tea.Cmd
		m.List, cmd = m.List.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Tick for message timeout
	if m.MessageTimeout > 0 {
		cmds = append(cmds, tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg {
			return TickMsg{}
		}))
	}

	if len(cmds) > 0 {
		return m, tea.Batch(cmds...)
	}
	return m, nil
}

// ShowDiff loads and displays the diff for the selected file
func (m *Model) ShowDiff(filePath string) tea.Cmd {
	return func() tea.Msg {
		// Multiple safety checks to avoid race conditions
		if m.Quitting {
			m.LoadingDiff = false
			return nil
		}

		if m.GitService == nil {
			m.LoadingDiff = false
			return nil
		}

		// Use mutex to ensure exclusive access to diff operations
		m.diffMutex.Lock()
		defer m.diffMutex.Unlock()

		// Double-check all conditions that could cause conflicts before proceeding
		if m.isNavigating || m.resizeInProgress {
			m.LoadingDiff = false
			return nil
		}

		if filePath != m.CurrentFile {
			m.LoadingDiff = false
			return nil
		}

		elapsed := time.Since(m.lastDiffTime)
		if elapsed < m.minDiffDelay {
			time.Sleep(m.minDiffDelay - elapsed)
		}

		m.lastDiffTime = time.Now()

		diff, err := m.GitService.GetFileDiff(filePath)
		if err != nil {
			return ErrMsg{err}
		}

		return DiffLoadedMsg{Diff: diff}
	}
}

// FormatDiffContent formats the diff content with syntax highlighting
func (m *Model) FormatDiffContent(diff *git.DiffResult) string {
	if diff == nil {
		return "No diff available"
	}

	if diff.IsBinary {
		return "Binary file differences not shown"
	}

	var result strings.Builder
	lines := strings.Split(diff.Content, "\n")

	maxWidth := max(m.DiffViewport.Width-2, 20)

	for _, line := range lines {
		if len(line) > 0 {
			displayLine := m.truncateText(line, maxWidth-5, "...")

			prefix := line[0:1]
			switch prefix {
			case "+":
				result.WriteString(m.StyleConfig.AddedStyle.Render(displayLine) + "\n")
			case "-":
				result.WriteString(m.StyleConfig.DeletedStyle.Render(displayLine) + "\n")
			default:
				result.WriteString(displayLine + "\n")
			}
		} else {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// GetDiffStats returns a formatted string with diff statistics
func (m *Model) GetDiffStats(diff *git.DiffResult) string {
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

// ConfirmStaging stages the selected files and exits
func (m *Model) ConfirmStaging() tea.Cmd {
	return func() tea.Msg {
		m.Quitting = true
		m.LoadingDiff = false

		var selectedPaths []string
		for i, item := range m.List.Items() {
			if m.Selected[i] {
				if fileItem, ok := item.(FileItem); ok {
					selectedPaths = append(selectedPaths, fileItem.Path)
				}
			}
		}

		if len(selectedPaths) == 0 && m.CurrentFile != "" {
			selectedPaths = []string{m.CurrentFile}
		}

		if len(selectedPaths) == 0 {
			return nil
		}

		gitService, err := git.NewGitService()
		if err != nil {
			return ErrMsg{err}
		}

		err = gitService.Stage(selectedPaths)
		if err != nil {
			return ErrMsg{err}
		}

		return StagingCompleteMsg{Files: selectedPaths}
	}
}

func (m *Model) truncateText(text string, maxLength int, ellipsis string) string {
	if len(text) <= maxLength || maxLength <= 0 {
		return text
	}

	truncateAt := maxLength
	if truncateAt <= len(ellipsis) {
		truncateAt = 1
	} else {
		truncateAt -= len(ellipsis)
	}

	if truncateAt > len(text) {
		truncateAt = len(text)
	}

	return text[:truncateAt] + ellipsis
}
