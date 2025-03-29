package git

import (
	"os"
	"path/filepath"
)

// GitService is an interface for git operations
type GitService interface {
	Status() ([]GitFile, error)
	Stage(paths []string) error
	Commit(commitType, message string) error
}

// DefaultGitService provides an implementation of GitService interface
type DefaultGitService struct {
	repo *GitRepository
}

func NewGitService() (*DefaultGitService, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	repoPath, err := findGitRepository(cwd)
	if err != nil {
		return nil, err
	}

	repo, err := NewGitRepository(repoPath)
	if err != nil {
		return nil, err
	}

	return &DefaultGitService{
		repo: repo,
	}, nil
}

func (s *DefaultGitService) Status() ([]GitFile, error) {
	files, err := s.repo.Status()
	if err != nil {
		// Fall back to exec implementation if go-git fails
		return GetStatus()
	}
	return files, nil
}

func (s *DefaultGitService) Stage(paths []string) error {
	err := s.repo.Stage(paths)
	if err != nil {
		// Fall back to exec implementation if go-git fails
		return StageFiles(paths)
	}
	return nil
}

func (s *DefaultGitService) Commit(commitType, message string) error {
	err := s.repo.Commit(commitType, message)
	if err != nil {
		// Fall back to exec implementation if go-git fails
		return Commit(commitType, message)
	}
	return nil
}

func findGitRepository(startPath string) (string, error) {
	path := startPath
	for {
		if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
			return path, nil
		}

		parentPath := filepath.Dir(path)
		if parentPath == path {
			return "", os.ErrNotExist
		}
		path = parentPath
	}
}
