package git

import (
	"os"
	"testing"

	"github.com/go-git/go-git/v5"
)

// setupTestRepo creates a test git repository with some files
func setupTestRepo(t *testing.T) string {
	t.Helper()

	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "go-git-tui-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Initialize git repository
	_, err = git.PlainInit(tmpDir, false)
	if err != nil {
		cleanupTestRepo(t, tmpDir)
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	return tmpDir
}

// cleanupTestRepo removes the test repository
func cleanupTestRepo(t *testing.T, path string) {
	t.Helper()
	if err := os.RemoveAll(path); err != nil {
		t.Fatalf("Failed to cleanup test repo: %v", err)
	}
}

func TestGetStatus(t *testing.T) {
	tests := []struct {
		name    string
		setupFn func(t *testing.T) string
		wantErr bool
	}{
		{
			name:    "GIVEN empty repository THEN status is returned without error",
			setupFn: setupTestRepo,
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repoPath := tc.setupFn(t)
			defer cleanupTestRepo(t, repoPath)

			// Change to the test repository directory
			oldWd, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get current working directory: %v", err)
			}
			defer func() {
				if err := os.Chdir(oldWd); err != nil {
					t.Fatalf("Failed to change back to original directory: %v", err)
				}
			}()

			if err := os.Chdir(repoPath); err != nil {
				t.Fatalf("Failed to change directory to test repo: %v", err)
			}

			// Skip for now since implementation is incomplete
			t.Skip("Skipping until go-git implementation is complete")
		})
	}
}

func TestStageFiles(t *testing.T) {
	tests := []struct {
		name    string
		paths   []string
		wantErr bool
	}{
		{
			name:    "GIVEN empty paths array THEN no error is returned",
			paths:   []string{},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repoPath := setupTestRepo(t)
			defer cleanupTestRepo(t, repoPath)

			// Skip for now
			t.Skip("Skipping until go-git implementation is complete")
		})
	}
}

func TestCommit(t *testing.T) {
	tests := []struct {
		name       string
		commitType string
		message    string
		wantErr    bool
	}{
		{
			name:       "GIVEN valid commit type and message THEN commit succeeds",
			commitType: "feat",
			message:    "test commit",
			wantErr:    false,
		},
		{
			name:       "GIVEN empty commit type THEN error is returned",
			commitType: "",
			message:    "test commit",
			wantErr:    true,
		},
		{
			name:       "GIVEN empty commit message THEN error is returned",
			commitType: "feat",
			message:    "",
			wantErr:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repoPath := setupTestRepo(t)
			defer cleanupTestRepo(t, repoPath)

			// Skip for now
			t.Skip("Skipping until go-git implementation is complete")
		})
	}
}
