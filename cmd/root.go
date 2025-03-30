package cmd

import (
	"fmt"
	"os"

	"github.com/LaansDole/go-git-tui/internal/ui"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// Version number hardcoded for simplicity
const Version = "v1.0.5-stable"

var (
	// Global flags
	verbose bool

	// Commands
	rootCmd = &cobra.Command{
		Use:   "go-git-tui",
		Short: "A Git TUI application",
		Long:  `A terminal user interface for Git operations built with go-git and Cobra CLI.`,
		// Don't set Version here to avoid the auto-generated --version flag
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
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// Version command to display version information
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Display the version of the go-git-tui application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("go-git-tui %s\n", Version)
	},
}

// Generate documentation for the CLI commands
var docsCmd = &cobra.Command{
	Use:    "generate-docs",
	Hidden: true,
	Short:  "Generate markdown documentation",
	Long:   `Generate markdown documentation for all go-git-tui commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		outDir := "./docs/"
		if len(args) > 0 {
			outDir = args[0]
		}

		if err := os.MkdirAll(outDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating docs directory: %v\n", err)
			os.Exit(1)
		}

		if err := doc.GenMarkdownTree(rootCmd, outDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating docs: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Documentation generated in %s\n", outDir)
	},
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// Add subcommands
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(docsCmd)

	// Enable bash completion
	rootCmd.CompletionOptions.DisableDefaultCmd = false
}
