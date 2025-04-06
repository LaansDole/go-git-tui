package ui

import (
	"fmt"

	"github.com/LaansDole/go-git-tui/internal/ui/add"
)

// StartAddTUI runs the add UI application with terminal UI
func StartAddTUI() error {
	return add.Run()
}

// StartCommitTUI runs the commit UI application with terminal UI
func StartCommitTUI() error {
	// TODO: Implement commit UI component
	return fmt.Errorf("commit UI not yet implemented")
}
