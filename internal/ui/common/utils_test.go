package common

import (
	"testing"
)

// TestTruncateText tests the text truncation helper function
func TestTruncateText(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		maxLength int
		ellipsis  string
		want      string
	}{
		{
			name:      "GIVEN string shorter than max length THEN returns original string",
			text:      "short",
			maxLength: 10,
			ellipsis:  "...",
			want:      "short",
		},
		{
			name:      "GIVEN string longer than max length THEN returns truncated string with ellipsis",
			text:      "this is a long string",
			maxLength: 10,
			ellipsis:  "...",
			want:      "this is...",
		},
		{
			name:      "GIVEN max length equal to text length THEN returns original string",
			text:      "exactly",
			maxLength: 7,
			ellipsis:  "...",
			want:      "exactly",
		},
		{
			name:      "GIVEN max length shorter than ellipsis THEN returns only part of ellipsis",
			text:      "some text",
			maxLength: 2,
			ellipsis:  "...",
			want:      "...",
		},
		{
			name:      "GIVEN empty string THEN returns empty string",
			text:      "",
			maxLength: 5,
			ellipsis:  "...",
			want:      "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := TruncateText(tc.text, tc.maxLength, tc.ellipsis)
			if got != tc.want {
				t.Errorf("truncateText() = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestTruncatePath tests the path truncation helper function
func TestTruncatePath(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		maxLength   int
		prefixChars int
		suffixChars int
		want        string
	}{
		{
			name:        "GIVEN path shorter than max length THEN returns original path",
			path:        "short/path",
			maxLength:   20,
			prefixChars: 5,
			suffixChars: 5,
			want:        "short/path",
		},
		{
			name:        "GIVEN long path THEN returns truncated path with middle ellipsis",
			path:        "some/very/long/path/to/file.txt",
			maxLength:   20,
			prefixChars: 10,
			suffixChars: 7,
			want:        "some/very/...o/file.txt",
		},
		{
			name:        "GIVEN max length too small for format THEN returns simple truncation",
			path:        "path/to/file.txt",
			maxLength:   5,
			prefixChars: 10,
			suffixChars: 10,
			want:        "pa...",
		},
		{
			name:        "GIVEN max length equal to path length THEN returns original path",
			path:        "exactly",
			maxLength:   7,
			prefixChars: 3,
			suffixChars: 3,
			want:        "exactly",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := TruncatePath(tc.path, tc.maxLength, tc.prefixChars, tc.suffixChars)
			if got != tc.want {
				t.Errorf("truncatePath() = %v, want %v", got, tc.want)
			}
		})
	}
}
