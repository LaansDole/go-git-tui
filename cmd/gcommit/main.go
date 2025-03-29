package main

import (
	"fmt"
	"os"

	"go-git-tui/internal/ui"

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
			if err := ui.RunCommitUI(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	// Execute adds all child commands to the root command and sets flags appropriately.
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
