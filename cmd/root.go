package cmd

import (
	"fmt"
	"os"

	"github.com/LaansDole/go-git-tui/internal/ui/add"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

const version = "v1.0.6"

var (
	verbose bool

	rootCmd = &cobra.Command{
		Use:   "go-git-tui",
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

	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Stage files interactively",
		Long: `Stage files using a terminal UI for faster interaction.
  
User Manual:
  - Use ARROW KEYS (UP/DOWN) to move
  - TAB to select files
  - ENTER once you have done selecting the files you want to add`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := add.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	commitCmd = &cobra.Command{
		Use:   "commit",
		Short: "Create commits interactively",
		Long: `Create Git commits with an interactive TUI.
  
User Manual:
  - Follow the on-screen prompts to select a commit type and enter a commit message
  - Use ARROW KEYS (UP/DOWN) to navigate options
  - ENTER to confirm selections`,
		Run: func(cmd *cobra.Command, args []string) {
			// Directly return a not-yet-implemented error without going through the UI package
			// This avoids the need to load non-existent files
			fmt.Fprintf(os.Stderr, "Error: commit UI not yet implemented\n")
			os.Exit(1)
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Display the version of the go-git-tui application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("go-git-tui %s\n", version)
	},
}

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
