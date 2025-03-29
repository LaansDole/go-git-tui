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

func GetStatus() ([]GitFile, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git status: %w", err)
	}

	files := []GitFile{}
	for line := range strings.SplitSeq(string(output), "\n") {
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

func StageFiles(paths []string) error {
	if len(paths) == 0 {
		return nil
	}

	args := append([]string{"add", "--"}, paths...)
	cmd := exec.Command("git", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to stage files: %w\nOutput: %s", err, output)
	}

	return nil
}

func Commit(commitType, message string) error {
	if commitType == "" || message == "" {
		return fmt.Errorf("commit type and message cannot be empty")
	}

	fullMessage := fmt.Sprintf("%s: %s", commitType, message)

	cmd := exec.Command("git", "commit", "-m", fullMessage)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to commit: %w\nOutput: %s", err, output)
	}

	return nil
}
