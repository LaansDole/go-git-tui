package git

import (
	"bytes"
	"strings"
)

// DiffResult contains the diff content of a file
type DiffResult struct {
	Path     string
	IsBinary bool
	Content  string
	Stats    DiffStats
}

// DiffStats contains statistics about a diff
type DiffStats struct {
	Added    int
	Deleted  int
	Modified int
}

// generateDiff creates a readable diff with +/- indicators and context lines
func generateDiff(oldContent, newContent string) (string, DiffStats) {
	const (
		addedLinePrefix     = "+ "
		deletedLinePrefix   = "- "
		unchangedLinePrefix = "  "
	)

	oldLines := strings.Split(oldContent, "\n")
	newLines := strings.Split(newContent, "\n")
	stats := DiffStats{}
	var sb strings.Builder

	// Find matching lines using longest common subsequence approach
	matchingLines := findLCS(oldLines, newLines)

	oldIdx, newIdx := 0, 0

	// Process each matching line and the lines between matches
	for _, match := range matchingLines {
		// Add deleted lines (those present in old but not in new)
		processDeletions(&sb, oldLines, oldIdx, match.OldIndex, deletedLinePrefix, &stats)
		oldIdx = match.OldIndex

		// Add added lines (those present in new but not in old)
		processAdditions(&sb, newLines, newIdx, match.NewIndex, addedLinePrefix, &stats)
		newIdx = match.NewIndex

		// Add the matching line as a context line
		appendLine(&sb, unchangedLinePrefix, newLines[newIdx])
		oldIdx++
		newIdx++
	}

	// Handle any remaining lines after the last match
	processDeletions(&sb, oldLines, oldIdx, len(oldLines), deletedLinePrefix, &stats)
	processAdditions(&sb, newLines, newIdx, len(newLines), addedLinePrefix, &stats)

	// Count modified lines as the sum of added and deleted
	stats.Modified = stats.Added + stats.Deleted

	return sb.String(), stats
}

// processDeletions adds deleted lines to the diff output and updates stats
func processDeletions(sb *strings.Builder, lines []string, fromIdx, toIdx int, prefix string, stats *DiffStats) {
	for i := fromIdx; i < toIdx; i++ {
		appendLine(sb, prefix, lines[i])
		stats.Deleted++
	}
}

// processAdditions adds new lines to the diff output and updates stats
func processAdditions(sb *strings.Builder, lines []string, fromIdx, toIdx int, prefix string, stats *DiffStats) {
	for i := fromIdx; i < toIdx; i++ {
		appendLine(sb, prefix, lines[i])
		stats.Added++
	}
}

// appendLine adds a line with the given prefix to the string builder
func appendLine(sb *strings.Builder, prefix, line string) {
	sb.WriteString(prefix + line + "\n")
}

// LineMatch represents a matching line in both old and new content
type LineMatch struct {
	OldIndex int
	NewIndex int
}

// findLCS finds the longest common subsequence of lines
func findLCS(oldLines, newLines []string) []LineMatch {
	// Create a map of line content to positions in the new content
	newLineMap := make(map[string][]int)
	for i, line := range newLines {
		newLineMap[line] = append(newLineMap[line], i)
	}

	// Find matching lines using a greedy approach
	var matches []LineMatch
	lastNewIndex := -1

	for oldIndex, line := range oldLines {
		if positions, ok := newLineMap[line]; ok {
			// Find the next valid position after lastNewIndex
			for _, newIndex := range positions {
				if newIndex > lastNewIndex {
					matches = append(matches, LineMatch{oldIndex, newIndex})
					lastNewIndex = newIndex
					break
				}
			}
		}
	}

	return matches
}

// isBinary checks if a file is likely to be binary by looking for null bytes
func isBinary(content []byte) bool {
	if len(content) > 512 {
		content = content[:512]
	}
	return bytes.IndexByte(content, 0) != -1
}
