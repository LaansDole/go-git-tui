package ui

import (
	"github.com/LaansDole/go-git-tui/internal/ui/add"
	"github.com/LaansDole/go-git-tui/internal/ui/commit"
)

// StartAddTUI runs the add UI application with terminal UI
func StartAddTUI() error {
	return add.Run()
}

// StartCommitTUI runs the commit UI application with terminal UI
func StartCommitTUI() error {
	return commit.Run()
}
