# Go Git TUI

This project provides a terminal user interface for Git operations, specifically focusing on staging files, git diff and committing changes using interactive selection.

## Roadmap

**TODOs:** 
- [x] Implement Cobra integration for CLI, combined with Bubbles/Charm TUI
- [x] Replace pure git commands with go-git dependency (primary implementation)
- [x] Add status command to view repository status
- [x] Add version command with multiple output formats
- [x] Add documentation generation and shell completion
- [x] Implement automatic version from Git tags
- [x] Refer to `gadd.sh` in `scripts/`, the tui should display the `git diff`
- [ ] A clearer UI to display that you have selected an option in `gcommit`
- [ ] Add better guidance in the TUI
- [x] Implement more comprehensive unit tests

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

### Option 2: Using Go Install

If you have Go installed, you can directly install the package using:

```shell
go install github.com/LaansDole/go-git-tui@latest
```

This will compile and install the `go-git-tui` binary to your $GOPATH/bin directory (or $HOME/go/bin if GOPATH is not set). 

To install the standalone commands:

```shell
# Create symlinks or wrappers for gadd and gcommit
echo -e '#!/bin/bash\nexec go-git-tui add "$@"' > $HOME/go/bin/gadd
echo -e '#!/bin/bash\nexec go-git-tui commit "$@"' > $HOME/go/bin/gcommit
chmod +x $HOME/go/bin/gadd $HOME/go/bin/gcommit
```

Make sure $GOPATH/bin (or $HOME/go/bin) is in your PATH.

## Usage

### As a Combined Tool

Use the main `go-git-tui` command with subcommands:

```shell
# Show status
go-git-tui status

# Interactive staging
go-git-tui add

# Interactive commit
go-git-tui commit

# Generate documentation
go-git-tui generate-docs

# Generate shell completion script
go-git-tui completion bash > ~/.bash_completion.d/go-git-tui

# Check version information
go-git-tui version

### As Standalone Commands

```shell
# Interactive staging with gadd
gadd

# Interactive commit with gcommit
gcommit
```

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
    make test-coverage
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