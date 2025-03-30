package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestRenderText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		style    lipgloss.Style
		wantSame bool // true if output should be the same as input
	}{
		{
			name:     "GIVEN text with red foreground style THEN styled text is returned",
			text:     "TestText",
			style:    lipgloss.NewStyle().Foreground(lipgloss.Color("1")),
			wantSame: false,
		},
		{
			name:     "GIVEN text with blue background style THEN styled text is returned",
			text:     "TestText",
			style:    lipgloss.NewStyle().Background(lipgloss.Color("4")),
			wantSame: false,
		},
		{
			name:     "GIVEN text with bold style THEN styled text is returned",
			text:     "TestText",
			style:    lipgloss.NewStyle().Bold(true),
			wantSame: false,
		},
		{
			name:     "GIVEN text with no style THEN original text is returned",
			text:     "TestText",
			style:    lipgloss.NewStyle(),
			wantSame: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Get direct rendering result for comparison
			directResult := tc.style.Render(tc.text)

			// Get result from our function
			result := RenderText(tc.text, tc.style)

			// Debug logging removed to keep test output clean

			// Check if results match direct rendering
			if result != directResult {
				t.Errorf("RenderText(%q, style) = %q, want %q", tc.text, result, directResult)
			}

			// Check if styling was applied when expected
			if tc.wantSame && result != tc.text {
				t.Errorf("Expected unstyled text to remain unchanged, got %q, want %q", result, tc.text)
			}

			if !tc.wantSame && result == tc.text && directResult != tc.text {
				t.Errorf("Expected styled text to be different from input, but got same text %q", result)
			}
		})
	}
}

func TestNewKeyMap(t *testing.T) {
	tests := []struct {
		name           string
		wantSelectKeys []string
		wantCancelKeys []string
	}{
		{
			name:           "GIVEN keymap creation THEN default keys are properly configured",
			wantSelectKeys: []string{"enter"},
			wantCancelKeys: []string{"ctrl+c", "esc"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			keyMap := NewKeyMap()

			// Check select keys
			if len(keyMap.Select.Keys()) != len(tc.wantSelectKeys) {
				t.Errorf("NewKeyMap().Select has %d keys, want %d", len(keyMap.Select.Keys()), len(tc.wantSelectKeys))
			}

			for _, wantKey := range tc.wantSelectKeys {
				found := false
				for _, gotKey := range keyMap.Select.Keys() {
					if gotKey == wantKey {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("NewKeyMap().Select missing key %q", wantKey)
				}
			}

			// Check cancel keys
			if len(keyMap.Cancel.Keys()) != len(tc.wantCancelKeys) {
				t.Errorf("NewKeyMap().Cancel has %d keys, want %d", len(keyMap.Cancel.Keys()), len(tc.wantCancelKeys))
			}

			for _, wantKey := range tc.wantCancelKeys {
				found := false
				for _, gotKey := range keyMap.Cancel.Keys() {
					if gotKey == wantKey {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("NewKeyMap().Cancel missing key %q", wantKey)
				}
			}
		})
	}
}
