package main

import (
	"flag"
	"fmt"
	"os"

	"go-git-tui/internal/ui"
)

func showHelp() {
	fmt.Println("Usage: gcommit [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -h, --help      Show this help message and exit")
	fmt.Println("")
	fmt.Println("Description:")
	fmt.Println("  This program provides an interactive TUI for creating Git commits.")
	fmt.Println("  It allows you to select commit types and enter commit messages.")
	fmt.Println("")
	fmt.Println("User Manual:")
	fmt.Println("  - Follow the on-screen prompts to select a commit type and enter a commit message")
	fmt.Println("  - Use ARROW KEYS (UP/DOWN) to navigate options")
	fmt.Println("  - ENTER to confirm selections")
}

func main() {
	help := flag.Bool("help", false, "Show help")
	h := flag.Bool("h", false, "Show help (shorthand)")
	flag.Parse()

	if *help || *h {
		showHelp()
		return
	}

	if err := ui.RunCommitUI(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
