package main

import (
	"flag"
	"fmt"
	"os"

	"go-git-tui/internal/ui"
)

func showHelp() {
	fmt.Println("Usage: gadd [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -h, --help      Show this help message and exit")
	fmt.Println("")
	fmt.Println("Description:")
	fmt.Println("  This program allows you to stage files using a terminal UI for faster interaction.")
	fmt.Println("")
	fmt.Println("User Manual:")
	fmt.Println("  - Use ARROW KEYS (UP/DOWN) to move")
	fmt.Println("  - TAB to select files")
	fmt.Println("  - ENTER once you have done selecting the files you want to add")
}

func main() {
	help := flag.Bool("help", false, "Show help")
	h := flag.Bool("h", false, "Show help (shorthand)")
	flag.Parse()

	if *help || *h {
		showHelp()
		return
	}

	if err := ui.RunAddUI(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
