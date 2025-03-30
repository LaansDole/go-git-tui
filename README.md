# Go Git TUI

This project provides a terminal user interface for Git operations, specifically focusing on staging files and committing changes using interactive selection. It includes three main commands: `git-tui` (the main application), `gadd` for staging files, and `gcommit` for committing changes.

## Features

- **Interactive File Selection**: Quickly stage files with visual selection interface
- **Commit Type Selection**: Choose from predefined commit types (feat, fix, docs, chores)
- **Customizable Commit Messages**: Enter detailed commit messages with type prefixes
- **Repository Status Display**: View repository status with colored file indicators
- **Go-Git Integration**: Primary implementation using native Go Git library
- **Shell Command Fallback**: Automatic fallback to Git CLI when needed
- **Cobra CLI Framework**: Structured command-line interface
- **Terminal UI**: Rich terminal interface using Bubble Tea, Bubbles, and Lipgloss
- **Version Information**: Check application version with different output formats
- **Shell Completion**: Built-in support for bash/zsh/fish completion
- **Documentation Generation**: Auto-generated markdown documentation for all commands

## Roadmap

**TODOs:** 
- [x] Implement Cobra integration for CLI, combined with Bubbles/Charm TUI
- [x] Replace pure git commands with go-git dependency (primary implementation)
- [x] Add status command to view repository status
- [x] Add version command with multiple output formats
- [x] Add documentation generation and shell completion
- [x] Implement automatic version from Git tags
- [ ] Refer to `gadd.sh` in `scripts/`, the tui should display the `git diff`
- [ ] A clearer UI to display that you have selected an option in `gcommit`
- [ ] Add better guidance in the TUI
- [ ] Implement more comprehensive unit tests

## Project Structure

```
go-git-tui
├── cmd/                      # Command entry points
│   ├── gadd/
│   │   └── main.go          # Entry for standalone gadd command  
│   ├── gcommit/
│   │   └── main.go          # Entry for standalone gcommit command
│   ├── root.go              # Main CLI command definitions
│   └── root_test.go         # Tests for root commands
├── internal/
│   ├── ui/
│   │   ├── common.go        # Common UI functions and styles
│   │   ├── add.go           # UI logic for gadd command
│   │   ├── add_test.go      # Tests for add UI
│   │   ├── commit.go        # UI logic for gcommit command
│   │   └── common_test.go   # Tests for common UI components
│   └── git/
│       ├── commands.go      # Git shell command operations
│       ├── commands_test.go # Tests for git commands
│       ├── gogit.go         # Go-git library implementation
│       └── service.go       # Git service interface
├── scripts/                  # Original shell script implementations
│   ├── gadd.sh              # Original bash script for gadd
│   └── gcommit.sh           # Original bash script for gcommit
├── .build/                   # Build output directory
├── bin/                      # Executable wrappers
├── go.mod                    # Go module configuration
├── go.sum                    # Module dependency checksums
├── Makefile                  # Build automation
├── main.go                   # Main entry point for git-tui
└── README.md                 # Project documentation
```

## Installation

### Option 1: Building from Source

1. Clone the repository:
    ```shell
    git clone https://github.com/LaansDole/go-git-tui.git
    cd go-git-tui
    ```

2. Ensure you have Go installed (Go 1.20+ recommended).

3. Build the project:
    ```shell
    make build
    ```

4. Add the bin directory to your PATH or use the executables directly:
    ```shell
    export PATH=$PATH:$PWD/bin
    ```

### Option 2: Global Installation

Run the following command to install globally:

```shell
make install
```

This will install `git-tui`, `gadd`, and `gcommit` to your $GOPATH/bin directory.

### Option 3: Using Go Install

If you have Go installed, you can directly install the package using:

```shell
go install github.com/LaansDole/go-git-tui@latest
```

This will compile and install the `git-tui` binary to your $GOPATH/bin directory (or $HOME/go/bin if GOPATH is not set). 

To install the standalone commands:

```shell
# Create symlinks or wrappers for gadd and gcommit
echo -e '#!/bin/bash\nexec git-tui add "$@"' > $HOME/go/bin/gadd
echo -e '#!/bin/bash\nexec git-tui commit "$@"' > $HOME/go/bin/gcommit
chmod +x $HOME/go/bin/gadd $HOME/go/bin/gcommit
```

Make sure $GOPATH/bin (or $HOME/go/bin) is in your PATH.

## Usage

### As a Combined Tool

Use the main `git-tui` command with subcommands:

```shell
# Show status
git-tui status

# Interactive staging
git-tui add

# Interactive commit
git-tui commit

# Generate documentation
git-tui generate-docs

# Generate shell completion script
git-tui completion bash > ~/.bash_completion.d/git-tui

# Check version information
git-tui version

# Check version in short format
git-tui version --format=short

# Check version in JSON format
git-tui version --format=json
```

### As Standalone Commands

```shell
# Interactive staging with gadd
gadd

# Interactive commit with gcommit
gcommit
```

### TUI Usage Guide

#### gadd (Interactive Staging)
- Use **↑/↓ arrow keys** to navigate through files
- Press **Tab** to select/deselect files for staging
- Press **Enter** to confirm and stage selected files
- Press **q** to quit without staging

#### gcommit (Interactive Commit)
- Use **↑/↓ arrow keys** to select commit type
- Press **Enter** to confirm type selection
- Type your commit message
- Press **Enter** to create the commit
- Press **Ctrl+C** to cancel at any time

### Makefile Utilities

The Makefile provides several utilities to streamline development and usage:

- **Build**: Compile the project binaries
    ```shell
    make build
    ```

- **Clean**: Remove compiled binaries and generated files
    ```shell
    make clean
    ```

- **Run**: Execute the main application
    ```shell
    make run
    ```

- **Direct Command Execution**:
    ```shell
    make gadd
    make gcommit
    ```

- **Test**: Run the test suite
    ```shell
    make test
    ```

- **Quick Check**: Run formatting and linting
    ```shell
    make check
    ```

## Implementation Details

### Git Operations

The application uses a dual-implementation approach for Git operations:
1. **Primary**: Uses the go-git library for native Go Git operations
2. **Fallback**: Automatically falls back to executing Git shell commands if go-git operations fail

This ensures maximum compatibility while leveraging the benefits of a native Go implementation.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.