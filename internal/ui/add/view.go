package add

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the UI - implements tea.Model interface
func (m *Model) View() string {
	if m.Quitting {
		return m.Message
	}

	if !m.Ready {
		return "Loading git repository..."
	}

	// Show message when there are no files to stage
	if len(m.List.Items()) == 0 {
		return m.StyleConfig.InfoStyle.Render("No files to stage. Working tree clean.")
	}

	titleText := m.StyleConfig.TitleStyle.Render("Go Git TUI - Stage Files")

	selectedCount := 0
	for _, selected := range m.Selected {
		if selected {
			selectedCount++
		}
	}

	statusText := m.StyleConfig.StatusBar.Render(
		fmt.Sprintf("%d files, %d selected", len(m.List.Items()), selectedCount))

	helpText := m.StyleConfig.HelpStyle.Render(
		"w/s: Navigate Files • j/k: Scroll Diff • Tab: Select • Enter: Confirm • q: Quit")

	diffTitle := "Diff"
	diffStats := ""
	if m.CurrentFile != "" {
		diffTitle = fmt.Sprintf("Diff for %s", m.CurrentFile)
		if m.CurrentDiff != nil {
			diffStats = m.GetDiffStats(m.CurrentDiff)
		}
	}

	diffTitleText := m.StyleConfig.TitleStyle.Render(diffTitle)
	diffStatsText := ""
	if diffStats != "" {
		diffStatsText = m.StyleConfig.InfoStyle.Render(diffStats)
	}

	diffContent := ""
	if m.LoadingDiff {
		diffContent = "Loading diff..."
	} else if m.CurrentDiff == nil {
		diffContent = "Press TAB to select a file and press 'd' to view diff"
	} else {
		diffContent = m.DiffViewport.View()
	}

	listView := m.StyleConfig.ListStyle.Render(m.List.View())

	diffWidth := m.Width / 2

	if len(diffTitle) > diffWidth-5 && diffWidth > 15 {
		diffTitle = m.truncateText(diffTitle, diffWidth-10, "...")
		diffTitleText = m.StyleConfig.TitleStyle.Render(diffTitle)
	}

	// Reduce spacing in diff panel and ensure content stays within bounds
	diffPanel := lipgloss.JoinVertical(
		lipgloss.Left,
		diffTitleText,
		diffStatsText,
		m.StyleConfig.DiffStyle.Height(m.DiffViewport.Height).Width(diffWidth-2).Render(diffContent),
	)

	// Create the vertical divider with exact height to match content
	dividerHeight := m.DiffViewport.Height
	if dividerHeight < 0 {
		dividerHeight = 0
	}
	verticalDivider := strings.Repeat("│\n", dividerHeight)
	dividerView := m.StyleConfig.DividerStyle.Render(verticalDivider)

	// Join horizontally with tight spacing
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		listView,
		dividerView,
		diffPanel,
	)

	messageDisplay := ""
	if m.Message != "" {
		messageStyle := m.StyleConfig.InfoStyle.Copy().Padding(0, 1).Bold(true).Width(m.Width)
		messageDisplay = messageStyle.Render(m.Message)
	}

	content = lipgloss.JoinVertical(
		lipgloss.Left,
		titleText,
		statusText,
		content,
		messageDisplay,
		helpText,
	)

	return m.StyleConfig.AppStyle.Render(content)
}
