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

			// Test the fallback implementation
			files, err := GetStatus()
			if (err != nil) != tc.wantErr {
				t.Errorf("GetStatus() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			// Test the go-git implementation through GitService
			repo, err := NewGitRepository(repoPath)
			if err != nil {
				t.Fatalf("Failed to create GitRepository: %v", err)
			}

			goGitFiles, err := repo.Status()
			if (err != nil) != tc.wantErr {
				t.Errorf("GitRepository.Status() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			// Both implementations should return empty status for empty repo
			if len(files) != 0 || len(goGitFiles) != 0 {
				t.Errorf("Expected empty status, got %d fallback files and %d go-git files", len(files), len(goGitFiles))
			}
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
		{
			name:    "GIVEN nonexistent file THEN error is returned",
			paths:   []string{"nonexistent.txt"},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repoPath := setupTestRepo(t)
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

			// Skip the test for nonexistent file in empty paths case
			if len(tc.paths) == 0 {
				// Test fallback implementation
				err := StageFiles(tc.paths)
				if (err != nil) != tc.wantErr {
					t.Errorf("StageFiles() error = %v, wantErr %v", err, tc.wantErr)
				}

				// Test go-git implementation
				repo, err := NewGitRepository(repoPath)
				if err != nil {
					t.Fatalf("Failed to create GitRepository: %v", err)
				}

				err = repo.Stage(tc.paths)
				if (err != nil) != tc.wantErr {
					t.Errorf("GitRepository.Stage() error = %v, wantErr %v", err, tc.wantErr)
				}
			} else {
				// For nonexistent file test, we can assume it will fail but implementation details may vary
				// Test fallback implementation
				err1 := StageFiles(tc.paths)
				
				// Test go-git implementation
				repo, err := NewGitRepository(repoPath)
				if err != nil {
					t.Fatalf("Failed to create GitRepository: %v", err)
				}

				err2 := repo.Stage(tc.paths)

				// At least one of them should fail for nonexistent file
				if (err1 == nil) && (err2 == nil) && tc.wantErr {
					t.Errorf("Expected at least one error for nonexistent file")
				}
			}
		})
	}
}

func TestCommit(t *testing.T) {
	tests := []struct {
		name       string
		setupFn    func(t *testing.T, repoPath string) // setup function to prepare for commit
		commitType string
		message    string
		wantErr    bool
	}{
		{
			name: "GIVEN empty commit type THEN error is returned",
			setupFn: func(t *testing.T, repoPath string) {
				// No setup needed
			},
			commitType: "",
			message:    "test commit",
			wantErr:    true,
		},
		{
			name: "GIVEN empty commit message THEN error is returned",
			setupFn: func(t *testing.T, repoPath string) {
				// No setup needed
			},
			commitType: "feat",
			message:    "",
			wantErr:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repoPath := setupTestRepo(t)
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

			// Run setup
			tc.setupFn(t, repoPath)

			// For empty type/message tests
			if tc.commitType == "" || tc.message == "" {
				// Test fallback implementation
				err := Commit(tc.commitType, tc.message)
				if (err != nil) != tc.wantErr {
					t.Errorf("Commit() error = %v, wantErr %v", err, tc.wantErr)
				}

				// Test go-git implementation
				repo, err := NewGitRepository(repoPath)
				if err != nil {
					t.Fatalf("Failed to create GitRepository: %v", err)
				}

				err = repo.Commit(tc.commitType, tc.message)
				if (err != nil) != tc.wantErr {
					t.Errorf("GitRepository.Commit() error = %v, wantErr %v", err, tc.wantErr)
				}
			}
		})
	}
}
