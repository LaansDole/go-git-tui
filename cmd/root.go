package cmd

import (
	"fmt"
	"os"

	"go-git-tui/internal/git"
	"go-git-tui/internal/ui"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// Version information - these will be overridden during build
var (
	// Version is the application version from the git tag
	Version = "dev"
	// Commit is the git commit hash
	Commit = "none"
	// BuildDate is the date the binary was built
	BuildDate = "unknown"
)

var (
	// Global flags
	verbose bool
	// Output format for version command
	versionFormat string

	// Commands
	rootCmd = &cobra.Command{
		Use:     "git-tui",
		Short:   "A Git TUI application",
		Long:    `A terminal user interface for Git operations built with go-git and Cobra CLI.`,
		Version: Version,
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

// Version command to display version information
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Display detailed version information about the git-tui application.`,
	Run: func(cmd *cobra.Command, args []string) {
		switch versionFormat {
		case "json":
			fmt.Printf(`{"version":"%s","commit":"%s","buildDate":"%s"}`+"\n", Version, Commit, BuildDate)
		case "short":
			fmt.Printf("git-tui v%s\n", Version)
		default: // full
			fmt.Printf("git-tui v%s (commit: %s, built: %s)\n", Version, Commit, BuildDate)
		}
	},
}

// Generate documentation for the CLI commands
var docsCmd = &cobra.Command{
	Use:    "generate-docs",
	Hidden: true,
	Short:  "Generate markdown documentation",
	Long:   `Generate markdown documentation for all git-tui commands.`,
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

	// Version command flags
	versionCmd.Flags().StringVarP(&versionFormat, "format", "f", "full", "Output format (full, short, json)")

	// Add subcommands
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(docsCmd)

	// Enable bash completion
	rootCmd.CompletionOptions.DisableDefaultCmd = false
}
