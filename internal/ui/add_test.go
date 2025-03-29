package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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
		quitting    bool
		wantEmpty   bool
		wantStrings []string
	}{
		{
			name:        "GIVEN normal state THEN returns view with expected strings",
			quitting:    false,
			wantEmpty:   false,
			wantStrings: []string{"Git Files", "Navigate", "Select/Deselect"},
		},
		{
			name:      "GIVEN quitting state THEN returns empty view",
			quitting:  true,
			wantEmpty: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			model := initialAddModel()
			model.quitting = tc.quitting
			got := model.View()

			if tc.wantEmpty && got != "" {
				t.Errorf("View() = %q, want empty string", got)
			}

			if !tc.wantEmpty {
				if got == "" {
					t.Error("View() returned empty string, want non-empty")
				}

				for _, s := range tc.wantStrings {
					if !strings.Contains(got, s) {
						t.Errorf("View() = %q, want it to contain %q", got, s)
					}
				}
			}
		})
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
			wantCmdNil: true,
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
			wantCmdNil: true,
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
