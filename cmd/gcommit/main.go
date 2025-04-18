package main

import (
	"fmt"
	"os"

	"github.com/LaansDole/go-git-tui/internal/ui"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "gcommit",
		Short: "Interactive TUI for creating Git commits",
		Long: `This program provides an interactive TUI for creating Git commits.
It allows you to select commit types and enter commit messages.

User Manual:
  - Follow the on-screen prompts to select a commit type and enter a commit message
  - Use ARROW KEYS (UP/DOWN) to navigate options
  - ENTER to confirm selections`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := ui.StartCommitTUI(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
