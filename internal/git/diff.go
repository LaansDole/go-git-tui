package git

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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

// GetFileDiff returns the diff content for a specific file
func (g *GitRepository) GetFileDiff(filePath string) (*DiffResult, error) {
	if g.repo == nil {
		return nil, fmt.Errorf("repository not initialized")
	}

	wt, err := g.repo.Worktree()
	if err != nil {
		return nil, err
	}

	// Get the working tree status for the file
	status, err := wt.Status()
	if err != nil {
		return nil, err
	}

	fileStatus, ok := status[filePath]
	if !ok {
		return nil, fmt.Errorf("file %s not found in status", filePath)
	}

	// Check if this is an untracked file
	if fileStatus.Worktree == git.Untracked {
		// For untracked files, return all content as added
		return g.diffForNewFile(filePath)
	}

	// For tracked files, get diff between HEAD and working tree
	return g.diffBetweenHeadAndWorktree(filePath)
}

// diffForNewFile creates a diff for untracked files
func (g *GitRepository) diffForNewFile(filePath string) (*DiffResult, error) {
	fullPath := filepath.Join(g.path, filePath)
	
	// Read the file content
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	// Check if it's a binary file
	if isBinary(content) {
		return &DiffResult{
			Path:     filePath,
			IsBinary: true,
			Content:  "[Binary file]",
			Stats: DiffStats{
				Added:    1,
				Deleted:  0,
				Modified: 0,
			},
		}, nil
	}

	// Create a diff with all lines added
	lines := strings.Split(string(content), "\n")
	
	var diffContent strings.Builder
	for _, line := range lines {
		diffContent.WriteString("+ " + line + "\n")
	}

	return &DiffResult{
		Path:     filePath,
		IsBinary: false,
		Content:  diffContent.String(),
		Stats: DiffStats{
			Added:    len(lines),
			Deleted:  0,
			Modified: 0,
		},
	}, nil
}

// diffBetweenHeadAndWorktree returns diff between HEAD and working tree
func (g *GitRepository) diffBetweenHeadAndWorktree(filePath string) (*DiffResult, error) {
	// Get the HEAD commit
	headRef, err := g.repo.Head()
	if err != nil {
		// If no HEAD exists (e.g., new repo), treat file as new
		if err == plumbing.ErrReferenceNotFound {
			return g.diffForNewFile(filePath)
		}
		return nil, err
	}

	headCommit, err := g.repo.CommitObject(headRef.Hash())
	if err != nil {
		return nil, err
	}

	// Get the file in HEAD
	headFile, err := headCommit.File(filePath)
	
	// If file doesn't exist in HEAD, treat as new file
	if err == object.ErrFileNotFound {
		return g.diffForNewFile(filePath)
	} else if err != nil {
		return nil, err
	}

	// Get file contents from HEAD
	headContents, err := headFile.Contents()
	if err != nil {
		return nil, err
	}

	// Get current file contents
	fullPath := filepath.Join(g.path, filePath)
	currentContents, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	// Check if file is binary
	if isBinary(currentContents) {
		return &DiffResult{
			Path:     filePath,
			IsBinary: true,
			Content:  "[Binary file]",
			Stats: DiffStats{
				Modified: 1,
			},
		}, nil
	}

	// Generate diff between HEAD and working copy
	diffContent, stats := generateDiff(headContents, string(currentContents))
	
	return &DiffResult{
		Path:     filePath,
		IsBinary: false,
		Content:  diffContent,
		Stats:    stats,
	}, nil
}

// generateDiff creates a readable diff with +/- indicators and context lines
func generateDiff(oldContent, newContent string) (string, DiffStats) {
	// Split content into lines
	oldLines := strings.Split(oldContent, "\n")
	newLines := strings.Split(newContent, "\n")

	// Format for display
	var sb strings.Builder
	stats := DiffStats{}

	// We'll use a more comprehensive approach that shows context
	// First, identify common lines and changed lines
	oldMap := make(map[string][]int) // maps content to line numbers
	for i, line := range oldLines {
		oldMap[line] = append(oldMap[line], i)
	}

	newMap := make(map[string][]int)
	for i, line := range newLines {
		newMap[line] = append(newMap[line], i)
	}

	// Track which lines we've processed
	oldProcessed := make([]bool, len(oldLines))
	newProcessed := make([]bool, len(newLines))

	// Find matching (unchanged) lines first to provide context
	type Match struct {
		OldIdx, NewIdx int
	}

	matches := []Match{}
	for lineContent, oldIndices := range oldMap {
		if newIndices, exists := newMap[lineContent]; exists {
			// For simplicity, we'll match the first occurrences
			// A more sophisticated approach would use longest common subsequence
			for i := 0; i < min(len(oldIndices), len(newIndices)); i++ {
				if !oldProcessed[oldIndices[i]] && !newProcessed[newIndices[i]] {
					matches = append(matches, Match{oldIndices[i], newIndices[i]})
					oldProcessed[oldIndices[i]] = true
					newProcessed[newIndices[i]] = true
				}
			}
		}
	}

	// Now build the output with context
	// We'll include some context lines around changes
	// First mark added/deleted lines
	for i, processed := range oldProcessed {
		if !processed {
			sb.WriteString("- " + oldLines[i] + "\n")
			stats.Deleted++
		}
	}

	for i, processed := range newProcessed {
		if !processed {
			sb.WriteString("+ " + newLines[i] + "\n")
			stats.Added++
		}
	}

	// Add some context lines 
	for i, line := range newLines {
		if newProcessed[i] {
			// This is a common line, show it as context
			sb.WriteString("  " + line + "\n")
		}
	}

	// Sort the output by line number for better readability
	// This is a simplified approach - a real diff would use an LCS algorithm
	lines := strings.Split(sb.String(), "\n")
	var result strings.Builder
	for _, line := range lines {
		if line != "" {
			result.WriteString(line + "\n")
		}
	}

	// Count modified lines for stats
	stats.Modified = stats.Added + stats.Deleted

	return result.String(), stats
}

// Helper function to get the minimum of two values
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// isBinary checks if a file is likely to be binary by looking for null bytes
func isBinary(content []byte) bool {
	if len(content) > 512 {
		content = content[:512]
	}
	return bytes.IndexByte(content, 0) != -1
}

// GetDiff is a convenience method on the service
func (s *DefaultGitService) GetFileDiff(path string) (*DiffResult, error) {
	return s.repo.GetFileDiff(path)
}
