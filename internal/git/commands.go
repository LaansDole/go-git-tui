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

// GetStatus returns a list of files and their status from git status
func GetStatus() ([]GitFile, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git status: %w", err)
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

// StageFiles adds the specified files to git staging area
func StageFiles(paths []string) error {
	if len(paths) == 0 {
		return nil
	}

	args := append([]string{"add"}, paths...)
	cmd := exec.Command("git", args...)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to stage files: %w", err)
	}

	return nil
}

// Commit creates a commit with the given type and message
func Commit(commitType, message string) error {
	fullMessage := fmt.Sprintf("%s: %s", commitType, message)

	cmd := exec.Command("git", "commit", "-m", fullMessage)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}
