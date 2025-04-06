package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GitRepositoryInterface defines the operations that can be performed on a git repository
type GitRepositoryInterface interface {
	Status() ([]GitFile, error)
	Stage(paths []string) error
	Commit(commitType, message string) error
	GetCurrentBranch() (string, error)
	GetFileDiff(filePath string) (*DiffResult, error)
}

// GitRepository represents a repository managed by go-git
type GitRepository struct {
	repo *git.Repository
	path string
}

// NewGitRepository creates a new GitRepository instance
func NewGitRepository(path string) (*GitRepository, error) {
	// Open the repository at the given path
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository: %w", err)
	}

	return &GitRepository{
		repo: repo,
		path: path,
	}, nil
}

// Status gets the repository status
func (g *GitRepository) Status() ([]GitFile, error) {
	if g.repo == nil {
		return nil, errors.New("repository not initialized")
	}

	wt, err := g.repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	status, err := wt.Status()
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	files := make([]GitFile, 0, len(status))
	for file, fileStatus := range status {
		statusStr := ""

		if fileStatus.Staging == git.Untracked && fileStatus.Worktree == git.Untracked {
			statusStr = "??"
		} else {
			staging := " "
			switch fileStatus.Staging {
			case git.Added:
				staging = "A"
			case git.Modified:
				staging = "M"
			case git.Deleted:
				staging = "D"
			case git.Renamed:
				staging = "R"
			}

			worktree := " "
			switch fileStatus.Worktree {
			case git.Added:
				worktree = "A"
			case git.Modified:
				worktree = "M"
			case git.Deleted:
				worktree = "D"
			case git.Untracked:
				worktree = "?"
			}

			statusStr = staging + worktree
		}

		files = append(files, GitFile{
			Status: statusStr,
			Path:   file,
		})
	}

	return files, nil
}

// Stage adds files to the staging area
func (g *GitRepository) Stage(paths []string) error {
	if g.repo == nil {
		return errors.New("repository not initialized")
	}

	wt, err := g.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	for _, path := range paths {
		_, err := wt.Add(path)
		if err != nil {
			return fmt.Errorf("failed to stage %s: %w", path, err)
		}
	}

	return nil
}

// Commit creates a new commit with the given message
func (g *GitRepository) Commit(commitType, message string) error {
	if g.repo == nil {
		return errors.New("repository not initialized")
	}

	wt, err := g.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Format the commit message
	fullMessage := fmt.Sprintf("%s: %s", commitType, message)

	// Get the configuration to retrieve the user's name and email
	config, err := g.repo.Config()
	if err != nil {
		return fmt.Errorf("failed to get git config: %w", err)
	}

	// Create the commit
	commit, err := wt.Commit(fullMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  config.User.Name,
			Email: config.User.Email,
			When:  time.Now(),
		},
	})

	if err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	// Get the commit object to verify
	_, err = g.repo.CommitObject(commit)
	if err != nil {
		return fmt.Errorf("failed to get commit object: %w", err)
	}

	return nil
}

// GetCurrentBranch returns the current branch name
func (g *GitRepository) GetCurrentBranch() (string, error) {
	if g.repo == nil {
		return "", errors.New("repository not initialized")
	}

	head, err := g.repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	if head.Name().IsBranch() {
		return head.Name().Short(), nil
	}

	return head.Hash().String()[:7], nil
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
