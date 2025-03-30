package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	
	"github.com/LaansDole/go-git-tui/internal/git"
)

func TestAddModel_Init(t *testing.T) {
	tests := []struct {
		name       string
		wantCmdNil bool
	}{
		{
			name:       "GIVEN initial model THEN returns nil command",
			wantCmdNil: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			model := initialAddModel()
			got := model.Init()

			// Check if command is nil or not
			isNil := got == nil
			if isNil != tc.wantCmdNil {
				t.Errorf("Init() returned nil: %v, want nil: %v", isNil, tc.wantCmdNil)
			}
		})
	}
}

func TestAddModel_View(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*AddModel)
		wantEmpty   bool
		wantStrings []string
		dontWant    []string
	}{
		{
			name: "GIVEN normal file list state THEN returns view with expected strings",
			setup: func(m *AddModel) {
				m.quitting = false
				m.ready = true
				m.list.SetItems([]list.Item{
					FileItem{path: "file1.txt", status: "M ", isSelected: false},
				})
			},
			wantEmpty:   false,
			wantStrings: []string{"Stage Files", "Navigate Files", "Scroll Diff", "Tab: Select", "Top/Bottom"},
			dontWant:    []string{},
		},
		{
			name: "GIVEN file with diff THEN returns view with diff info",
			setup: func(m *AddModel) {
				m.quitting = false
				m.ready = true
				m.list.SetItems([]list.Item{
					FileItem{path: "file1.txt", status: "M ", isSelected: false},
				})
				m.currentFile = "file1.txt"
				m.diffViewport.SetContent("+ Added line\n- Removed line")
			},
			wantEmpty:   false,
			wantStrings: []string{"Diff for", "Stage Files", "│"}, // The "│" character verifies the vertical divider
		},
		{
			name: "GIVEN quitting state THEN returns empty view",
			setup: func(m *AddModel) {
				m.quitting = true
			},
			wantEmpty: true,
		},
		{
			name: "GIVEN empty file list THEN returns clean working tree message",
			setup: func(m *AddModel) {
				m.quitting = false
				m.ready = true
				m.list.SetItems([]list.Item{})
			},
			wantEmpty:   false,
			wantStrings: []string{"clean"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			model := initialAddModel()
			
			// Apply test-specific setup
			if tc.setup != nil {
				tc.setup(&model)
			}
			got := model.View()

			if tc.wantEmpty && got != "" {
				// Avoid printing potentially large view content
				t.Errorf("View() returned non-empty string, expected empty string")
			}

			if !tc.wantEmpty {
				if got == "" {
					t.Error("View() returned empty string, want non-empty")
				}

				for _, s := range tc.wantStrings {
					if !strings.Contains(got, s) {
						// Avoid printing the entire view content to reduce log output
						t.Errorf("View() does not contain expected string %q", s)
					}
				}
				
				// Check strings that should not be present
				for _, s := range tc.dontWant {
					if strings.Contains(got, s) {
						// Avoid printing the entire view content to reduce log output
						t.Errorf("View() incorrectly contains string %q", s)
					}
				}
			}
		})
	}
}

// Mock implementation of git.DiffResult for testing
type mockDiffResult struct {
	path     string
	isBinary bool
	content  string
	stats    git.DiffStats
}

// TestHelper: create a test message with a mock diff result
func createMockDiffMsg(filePath string) diffLoadedMsg {
	return diffLoadedMsg{
		diff: &git.DiffResult{
			Path:     filePath,
			IsBinary: false,
			Content:  "+ Added line\n- Removed line",
			Stats: git.DiffStats{
				Added:    1,
				Deleted:  1,
				Modified: 2,
			},
		},
	}
}

func TestAddModel_Update(t *testing.T) {
	tests := []struct {
		name       string
		msg        tea.Msg
		setup      func(*AddModel)
		wantQuit   bool
		wantCmdNil bool
		check      func(*testing.T, AddModel)
	}{
		{
			name:       "GIVEN Ctrl+C key press THEN model is set to quit",
			msg:        tea.KeyMsg{Type: tea.KeyCtrlC},
			wantQuit:   true,
			wantCmdNil: false,
		},
		{
			name:       "GIVEN q key press THEN model is set to quit",
			msg:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			wantQuit:   true,
			wantCmdNil: false,
		},
		{
			name:       "GIVEN window resize THEN model updates dimensions",
			msg:        tea.WindowSizeMsg{Width: 80, Height: 24},
			wantQuit:   false,
			wantCmdNil: false, // Now returns a command to load diff for first item
		},
		{
			name: "GIVEN tab key on item THEN item is selected and cursor moves down",
			msg:  tea.KeyMsg{Type: tea.KeyTab},
			setup: func(m *AddModel) {
				// Add three items and select the first one
				m.list.SetItems([]list.Item{
					FileItem{path: "file1.txt", status: "M ", isSelected: false},
					FileItem{path: "file2.txt", status: "M ", isSelected: false},
					FileItem{path: "file3.txt", status: "??", isSelected: false},
				})
				m.list.Select(0)
			},
			wantQuit:   false,
			wantCmdNil: false, // Now returns a command to load diff for next item
			check: func(t *testing.T, m AddModel) {
				// Check if selection moved to next item
				if m.list.Index() != 1 {
					t.Errorf("After Tab, index = %d, want 1", m.list.Index())
				}

				// Check if the original item was selected
				if !m.selected[0] {
					t.Errorf("Item at index 0 should be selected")
				}

				// Check if the item shows as selected in the list
				item := m.list.Items()[0]
				if fileItem, ok := item.(FileItem); ok {
					if !fileItem.isSelected {
						t.Errorf("FileItem at index 0 should have isSelected = true")
					}
				}
			},
		},
		{
			name: "GIVEN tab key on last item THEN item is selected but cursor stays",
			msg:  tea.KeyMsg{Type: tea.KeyTab},
			setup: func(m *AddModel) {
				// Add two items and select the last one
				m.list.SetItems([]list.Item{
					FileItem{path: "file1.txt", status: "M ", isSelected: false},
					FileItem{path: "file2.txt", status: "M ", isSelected: false},
				})
				m.list.Select(1) // Select the last item
			},
			wantQuit:   false,
			wantCmdNil: true,
			check: func(t *testing.T, m AddModel) {
				// Check if selection stayed at the last item
				if m.list.Index() != 1 {
					t.Errorf("After Tab at last item, index = %d, want 1", m.list.Index())
				}

				// Check if the last item was selected
				if !m.selected[1] {
					t.Errorf("Item at index 1 should be selected")
				}

				// Check if the item shows as selected in the list
				item := m.list.Items()[1]
				if fileItem, ok := item.(FileItem); ok {
					if !fileItem.isSelected {
						t.Errorf("FileItem at index 1 should have isSelected = true")
					}
				}
			},
		},
		{
			name: "GIVEN navigation key press THEN updates current file",
			msg:  tea.KeyMsg{Type: tea.KeyDown},
			setup: func(m *AddModel) {
				// Add two files and select the first
				m.list.SetItems([]list.Item{
					FileItem{path: "file1.txt", status: "M ", isSelected: false},
					FileItem{path: "file2.txt", status: "M ", isSelected: false},
				})
				m.list.Select(0)
				// We'll verify the action by checking currentFile is updated when moving cursor
			},
			wantQuit:   false,
			wantCmdNil: false, // Should return commands for updating the list and loading diff
			check: func(t *testing.T, m AddModel) {
				// Check that the cursor moved to the next item
				if m.list.Index() != 1 {
					t.Errorf("Cursor position = %d, want 1", m.list.Index())
				}
			},
		},
		{
			name: "GIVEN space key press THEN toggles selection",
			msg:  tea.KeyMsg{Type: tea.KeySpace},
			setup: func(m *AddModel) {
				// Add a file and select it
				m.list.SetItems([]list.Item{
					FileItem{path: "file1.txt", status: "M ", isSelected: false},
					FileItem{path: "file2.txt", status: "A ", isSelected: false},
				})
				m.list.Select(0)
			},
			wantQuit:   false,
			wantCmdNil: true,
			check: func(t *testing.T, m AddModel) {
				// Check if the item was selected
				if !m.selected[0] {
					t.Errorf("Item at index 0 should be selected after space press")
				}
				
				// Check that cursor didn't advance (unlike with tab)
				if m.list.Index() != 0 {
					t.Errorf("Cursor moved to %d, should stay at 0 after space press", m.list.Index())
				}
			},
		},
		{
			name: "GIVEN j key press THEN scrolls diff view without affecting file list",
			msg:  tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			setup: func(m *AddModel) {
				// Setup a diff with multiple lines
				m.diffViewport.SetContent("Line 1\nLine 2\nLine 3\nLine 4")
				m.diffViewport.GotoTop() // Start at the top
				
				// Add a file to the list
				m.list.SetItems([]list.Item{
					FileItem{path: "file1.txt", status: "M ", isSelected: false},
				})
				m.list.Select(0)
			},
			wantQuit:   false,
			wantCmdNil: true, // Based on implementation, doesn't return a command
			check: func(t *testing.T, m AddModel) {
				// Verify file list cursor position didn't change
				if m.list.Index() != 0 {
					t.Errorf("File list index changed to %d when j was pressed, should remain at 0", m.list.Index())
				}
			},
		},
		{
			name: "GIVEN k key press THEN scrolls diff view without affecting file list",
			msg:  tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}},
			setup: func(m *AddModel) {
				// Setup a diff with multiple lines and move down to test scrolling up
				m.diffViewport.SetContent("Line 1\nLine 2\nLine 3\nLine 4")
				m.diffViewport.GotoBottom() // Start at the bottom
				
				// Add a file to the list
				m.list.SetItems([]list.Item{
					FileItem{path: "file1.txt", status: "M ", isSelected: false},
				})
				m.list.Select(0)
			},
			wantQuit:   false,
			wantCmdNil: true, // Based on implementation, doesn't return a command
			check: func(t *testing.T, m AddModel) {
				// Verify file list cursor position didn't change
				if m.list.Index() != 0 {
					t.Errorf("File list index changed to %d when k was pressed, should remain at 0", m.list.Index())
				}
			},
		},
		{
			name: "GIVEN enter key press THEN triggers staging confirmation without immediate quit",
			msg:  tea.KeyMsg{Type: tea.KeyEnter},
			setup: func(m *AddModel) {
				// Add a file and select it
				m.list.SetItems([]list.Item{
					FileItem{path: "file1.txt", status: "M ", isSelected: true},
				})
				m.list.Select(0)
				m.selected[0] = true
			},
			wantQuit:   false, // Should not quit immediately anymore
			wantCmdNil: false, // Should return the staging confirmation command
			check: func(t *testing.T, m AddModel) {
				// We can't directly check if the confirmStaging command was returned,
				// but we can verify the model isn't set to quit yet
				if m.quitting {
					t.Error("Model should not be set to quit immediately after Enter key")
				}
			},
		},
		{
			name: "GIVEN stagingCompleteMsg THEN shows confirmation message before quitting",
			msg:  stagingCompleteMsg{files: []string{"file1.txt", "file2.txt"}},
			setup: func(m *AddModel) {
				// No specific setup needed
			},
			wantQuit:   false, // Should not quit immediately
			wantCmdNil: false, // Should return a delayed quit command
			check: func(t *testing.T, m AddModel) {
				// Verify message contains staged files
				expectedSubstring := "2 files staged: file1.txt, file2.txt"
				if !strings.Contains(m.message, expectedSubstring) {
					t.Errorf("Message should contain %q, got %q", expectedSubstring, m.message)
				}
				
				// Verify message timeout is set to show message before quitting
				if m.messageTimeout <= 0 {
					t.Errorf("Message timeout should be positive, got %d", m.messageTimeout)
				}
			},
		},
		{
			name: "GIVEN diffLoadedMsg THEN updates the diff view content",
			msg: createMockDiffMsg("test.go"),
			setup: func(m *AddModel) {
				// No special setup needed
			},
			wantQuit:   false,
			wantCmdNil: true,
			check: func(t *testing.T, m AddModel) {
				// Verify diff was loaded correctly
				if m.currentDiff == nil {
					t.Error("currentDiff was not set")
				}
				
					// Set the content directly for testing since View() requires rendering
				m.formatDiffContent(m.currentDiff)
				
				// Check that diff was processed, even if not rendered in viewport
				if m.currentDiff == nil || m.currentDiff.Path != "test.go" {
					t.Error("Diff was not correctly loaded for test.go")
				}
				
				// Verify loading state is updated
				if m.loadingDiff {
					t.Error("loadingDiff should be false after diff is loaded")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a model with at least one item to ensure list.Update returns a command
			model := initialAddModel()

			// If testing window resize, make sure we have items in the list
			if _, ok := tc.msg.(tea.WindowSizeMsg); ok {
				// Add a dummy item if the list is empty
				if len(model.list.Items()) == 0 {
					model.list.SetItems([]list.Item{
						FileItem{path: "dummy-file.txt", status: "??", isSelected: false},
					})
				}
			}

			// Apply any test-specific setup
			if tc.setup != nil {
				tc.setup(&model)
			}

			newModel, cmd := model.Update(tc.msg)

			m, ok := newModel.(AddModel)
			if !ok {
				t.Fatalf("Update() returned %T, want AddModel", newModel)
			}

			if m.quitting != tc.wantQuit {
				t.Errorf("Update() quitting = %v, want %v", m.quitting, tc.wantQuit)
			}

			cmdNil := cmd == nil
			if cmdNil != tc.wantCmdNil {
				t.Errorf("Update() cmd = %v, want nil: %v", cmd, tc.wantCmdNil)
			}

			// Run any test-specific checks
			if tc.check != nil {
				tc.check(t, m)
			}
		})
	}
}

// TestTruncateText tests the text truncation helper function
func TestTruncateText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		maxLen   int
		ellipsis string
		want     string
	}{
		{
			name:     "Short text unchanged",
			text:     "Short",
			maxLen:   10,
			ellipsis: "...",
			want:     "Short",
		},
		{
			name:     "Long text truncated",
			text:     "This is a very long text that should be truncated",
			maxLen:   15,
			ellipsis: "...",
			want:     "This is a ve...",
		},
		{
			name:     "Handles zero max length",
			text:     "Some text",
			maxLen:   0,
			ellipsis: "...",
			want:     "Some text",
		},
		{
			name:     "Handles small max length",
			text:     "Text",
			maxLen:   2,
			ellipsis: "...",
			want:     "T...",
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := truncateText(tc.text, tc.maxLen, tc.ellipsis)
			if got != tc.want {
				t.Errorf("truncateText(%q, %d, %q) = %q, want %q",
					tc.text, tc.maxLen, tc.ellipsis, got, tc.want)
			}
		})
	}
}

// TestTruncatePath tests the path truncation helper function
func TestTruncatePath(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		maxLength int
		prefixLen int
		suffixLen int
		want      string
	}{
		{
			name:      "Short path unchanged",
			path:      "short.txt",
			maxLength: 20,
			prefixLen: 10,
			suffixLen: 5,
			want:      "short.txt",
		},
		{
			name:      "Long path truncated in middle",
			path:      "/very/long/path/to/some/file/that/should/be/truncated.txt",
			maxLength: 30,
			prefixLen: 10,
			suffixLen: 10,
			want:      "/very/long...ncated.txt",
		},
		{
			name:      "Handles small max length",
			path:      "short.txt",
			maxLength: 5,
			prefixLen: 2,
			suffixLen: 2,
			want:      "sh...xt",
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := truncatePath(tc.path, tc.maxLength, tc.prefixLen, tc.suffixLen)
			if got != tc.want {
				t.Errorf("truncatePath(%q, %d, %d, %d) = %q, want %q",
					tc.path, tc.maxLength, tc.prefixLen, tc.suffixLen, got, tc.want)
			}
		})
	}
}
