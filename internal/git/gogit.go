package git

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

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
