package cmd

import (
	"fmt"
	"os"

	"go-git-tui/internal/ui"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-tui",
	Short: "A Git TUI application",
	Long:  `A terminal user interface for Git operations built in Go`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default behavior when no subcommand is specified
		if err := cmd.Help(); err != nil {
			fmt.Fprintf(os.Stderr, "Error showing help: %v\n", err)
			os.Exit(1)
		}
	},
}

// Add commands for gadd and gcommit
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Stage files interactively",
	Long: `Stage files using a terminal UI for faster interaction.
  
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

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Create commits interactively",
	Long: `Create Git commits with an interactive TUI.
  
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
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(commitCmd)
}
