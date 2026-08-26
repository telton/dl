# dl

A minimal devlog CLI for quick daily progress notes.

## What it does

- Interactive REPL: type a line, hit Enter, it gets logged with a timestamp.
- Pipe mode: `echo "did the thing" | dl` for one-off entries.
- Appends entries to a daily log file (`YYYY-MM-DD.log`) with `HH:MM:SS` timestamps.
- Stores files in `~/.local/share/dl/` by default (configurable).

## Features

- Simple REPL: Enter submits the line
- Timestamped entries in log format
- TOML configuration with XDG directory support
- Single static binary, trivial cross-compilation
- Pipe mode for shell integration

## Installation

### Download a Release

Download the latest release for your platform from the [releases page](https://github.com/telton/dl/releases) and place it somewhere in your `PATH`.

```bash
# Example: Linux x86_64
curl -L -o dl "https://github.com/telton/dl/releases/latest/download/dl_linux_amd64"
chmod +x dl
mv dl /usr/local/bin/
```

### Install with Go

Requires Go 1.27 or later.

```bash
go install github.com/telton/dl@latest
```

### Build from Source

```bash
git clone https://github.com/telton/dl.git
cd dl
go build -o dl .
```

### Nix

This repo includes a flake for development.

```bash
nix develop    # Enter shell with Go tooling
nix build      # Build the package (if flake outputs are configured)
```

## Usage

### REPL Mode

```bash
./dl
```

```
[14:32:01] fixed the race condition in pool.go
[14:45:22] started CLI tool scaffolding

dl — ctrl+c to quit
> fixed the race condition in pool.go
> started CLI tool scaffolding
```

Hit Enter to submit. `ctrl+c` to quit.

### Pipe Mode

```bash
echo "fixed the race condition in pool.go" | ./dl
```

### Configuration

The config file lives at `~/.config/dl/config.toml` by default, or specify `--config /path/to/config.toml`.

Generated on first run if absent:

```toml
data_dir = "/home/user/.local/share/dl"
```

## Output Format

Each day gets its own file in the data directory:

```
[14:32:01] fixed the race condition in the connection pool
[14:45:22] started CLI tool scaffolding
```

Simple log format — readable in any editor or pager.

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with race detection
go test -race ./...
```

### Linting

```bash
golangci-lint run ./...
```

### Cross-Compilation

```bash
# Linux x86_64
GOOS=linux GOARCH=amd64 go build -o dl-linux-amd64 .

# macOS ARM64
GOOS=darwin GOARCH=arm64 go build -o dl-darwin-arm64 .
```

