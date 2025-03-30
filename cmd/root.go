package cmd

import (
	"fmt"
	"os"

	"go-git-tui/internal/git"
	"go-git-tui/internal/ui"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	verbose bool
	
	// Commands
	rootCmd = &cobra.Command{
		Use:   "git-tui",
		Short: "A Git TUI application",
		Long:  `A terminal user interface for Git operations built with go-git and Cobra CLI.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Default behavior when no subcommand is specified
			if err := cmd.Help(); err != nil {
				fmt.Fprintf(os.Stderr, "Error showing help: %v\n", err)
				os.Exit(1)
			}
		},
	}

	// Add command for staging files interactively
	addCmd = &cobra.Command{
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

	// Commit command for creating commits interactively
	commitCmd = &cobra.Command{
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

	// Status command to display repository status
	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Show repository status",
		Long:  `Display the current status of the Git repository using go-git.`,
		Run: func(cmd *cobra.Command, args []string) {
			gitService, err := git.NewGitService()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error initializing git service: %v\n", err)
				os.Exit(1)
			}
			
			files, err := gitService.Status()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting status: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Println("Git Repository Status:")
			if len(files) == 0 {
				fmt.Println("Working directory clean")
			} else {
				for _, file := range files {
					fmt.Printf("%s %s\n", file.Status, file.Path)
				}
			}
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	
	// Add subcommands
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(statusCmd)
}
