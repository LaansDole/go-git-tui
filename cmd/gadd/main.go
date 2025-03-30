package main

import (
	"fmt"
	"os"

	"github.com/LaansDole/go-git-tui/internal/ui"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "gadd",
		Short: "Interactive TUI for staging Git files",
		Long: `This program allows you to stage files using a terminal UI for faster interaction.

User Manual:
  - Use ARROW KEYS (UP/DOWN) to move
  - TAB to select files
  - ENTER once you have done selecting the files you want to add`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := ui.RunAddUI(); err != nil {
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
