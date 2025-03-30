package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// GitFile represents a file in git status
type GitFile struct {
	Status string
	Path   string
}

// NOTE: This file contains fallback implementations using native git commands.
// These functions should ONLY be used when the go-git implementation fails.
// All direct usage should go through the GitService interface.

// GetStatus is a fallback implementation that uses the git command-line tool.
// It parses the output of "git status --porcelain" to get the status of files in the repository.
// This should only be used when the go-git implementation fails.
func GetStatus() ([]GitFile, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("fallback git status failed: %w", err)
	}

	files := []GitFile{}
	for _, line := range strings.Split(string(output), "\n") {
		if len(line) > 3 {
			status := line[0:2]
			path := strings.TrimSpace(line[2:])
			if path != "" {
				files = append(files, GitFile{
					Status: status,
					Path:   path,
				})
			}
		}
	}

	return files, nil
}

// StageFiles is a fallback implementation that uses the git command-line tool.
// It stages the specified files using "git add".
// This should only be used when the go-git implementation fails.
func StageFiles(paths []string) error {
	if len(paths) == 0 {
		return nil
	}

	args := append([]string{"add", "--"}, paths...)
	cmd := exec.Command("git", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("fallback git add failed: %w\nOutput: %s", err, output)
	}

	return nil
}

// Commit is a fallback implementation that uses the git command-line tool.
// It creates a commit with the specified type and message.
// This should only be used when the go-git implementation fails.
func Commit(commitType, message string) error {
	if commitType == "" || message == "" {
		return fmt.Errorf("commit type and message cannot be empty")
	}

	fullMessage := fmt.Sprintf("%s: %s", commitType, message)

	cmd := exec.Command("git", "commit", "-m", fullMessage)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("fallback git commit failed: %w\nOutput: %s", err, output)
	}

	return nil
}
