# Go Git TUI

This project provides a terminal user interface for Git operations, specifically focusing on staging files and committing changes using fuzzy selection. It consists of two main commands: `gadd` for staging files and `gcommit` for committing changes.

## Features

- **Fuzzy File Selection**: Quickly stage files using fuzzy matching.
- **Interactive Commit Messages**: Easily specify commit types and messages through prompts.

## Project Structure

```
go-git-tui
├── cmd
│   ├── gadd
│   │   └── main.go       # Main entry for the gadd command
│   └── gcommit
│       └── main.go       # Main entry for the gcommit command
├── internal
│   ├── ui
│   │   ├── common.go     # Common UI functions
│   │   ├── add.go        # UI logic for gadd command
│   │   └── commit.go     # UI logic for gcommit command
│   └── git
│       └── commands.go   # Git command interactions
├── go.mod                 # Go module configuration
├── go.sum                 # Module dependency checksums
├── Makefile               # Build instructions
└── README.md              # Project documentation
```

## Installation

1. Clone the repository:
    ```shell
    git clone https://github.com/AnDoLeLongANZ/go-git-tui.git
    cd go-git-tui
    ```

2. Ensure you have Go installed on your machine.

3. Install the required dependencies:
    ```shell
    go mod tidy
    ```

## Usage

### gadd

To stage files using the `gadd` command, run:
```shell
go run cmd/gadd/main.go
```
This will open a fuzzy finder interface to select files to stage.

### gcommit

To commit staged changes using the `gcommit` command, run:
```shell
go run cmd/gcommit/main.go
```
You will be prompted to enter a commit type and message.

### Makefile Utilities

The Makefile provides several utilities to streamline development and usage:

- **Build**: Compile the project binaries.
    ```shell
    make build
    ```

- **Clean**: Remove compiled binaries and other generated files.
    ```shell
    make clean
    ```

- **Run**: Execute the main application.
    ```shell
    make run
    ```

- **Test**: Run the test suite.
    ```shell
    make test
    ```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.