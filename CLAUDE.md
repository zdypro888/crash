# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is a Go library that provides crash and panic handling utilities for cross-platform applications. The library redirects stderr (for panic traces) and logs to files with automatic rotation and compression, and provides flexible panic handling hooks.

## Core Architecture

The library uses **platform-specific build tags** to handle OS differences:

- `crash_unix.go` (darwin, linux): Uses `unix.Dup2` from `golang.org/x/sys/unix`
- `crash_windows.go`: Uses `windows.SetStdHandle` from `golang.org/x/sys/windows`
- `hook.go`: Platform-independent panic handler system
- `logger.go`: Log rotation using lumberjack

### Key Components

1. **Panic Redirection** (`InitPanicFile` / `InitPanicFileWithTee`):
   - `InitPanicFile`: Redirects stderr to a file only (original behavior)
   - `InitPanicFileWithTee`: Redirects stderr to both file and console (new feature)
   - Unix: Uses `unix.Dup2` with pipe and goroutine for tee functionality
   - Windows: Uses `windows.SetStdHandle` with pipe and goroutine for tee functionality
   - File is opened with `O_APPEND` to preserve previous crashes

2. **Panic Hook System** (`hook.go`):
   - `AddPanicHandler`: Register custom panic handlers
   - `Recover`: Catch panic and execute all registered handlers
   - `WrapMain`: Convenience wrapper for main function
   - `RecoverToFile`: Built-in handler that writes panic to file
   - Handlers receive both panic value and stack trace
   - Thread-safe handler registration

3. **Log Rotation** (`RedirectLog`):
   - Uses `gopkg.in/natefinch/lumberjack.v2` for automatic log rotation
   - Default settings: 100MB max size, 3 backups, 30-day retention, compression enabled
   - Redirects Go's standard `log` package output only (not fmt.Println)

## Development Commands

### Building
```bash
go build .
```

### Testing
```bash
go test ./...          # Run all tests
go test -v ./...       # Verbose output
go test -run TestName  # Run specific test
```

### Cross-platform Testing
```bash
GOOS=windows go build .  # Test Windows build
GOOS=linux go build .    # Test Linux build
GOOS=darwin go build .   # Test macOS build
```

### Dependency Management
```bash
go mod tidy        # Clean up dependencies
go mod download    # Download dependencies
go list -m all     # List all dependencies
```

## Important Implementation Notes

- **Modern Unix API**: Uses `golang.org/x/sys/unix.Dup2` instead of deprecated `syscall.Dup2`
- **No finalizers**: The code intentionally avoids `runtime.SetFinalizer` for file descriptors
- **File descriptor lifetime**: File descriptors remain open for the process lifetime (intentional for crash logging)
- **Build tags**: Always test both Unix and Windows implementations when modifying panic redirection logic
- **Tee implementation**: Uses pipe + goroutine + io.MultiWriter to simultaneously write to file and console
- **Hook safety**: Panic handlers are protected from their own panics
- **Go version**: Requires Go 1.24.0 or higher

## Usage Examples

### Basic panic file redirection
```go
crash.InitPanicFile("panic.log")
```

### Tee to both file and console
```go
crash.InitPanicFileWithTee("panic.log")
```

### Custom panic handling
```go
crash.AddPanicHandler(func(panicValue interface{}, stackTrace []byte) {
    // Send alert, upload to monitoring system, etc.
})
defer crash.Recover(true)  // rePanic=true to continue panic after handling
```

### Wrap main function
```go
func main() {
    crash.WrapMain(func() {
        // Your application code
    })
}
```

## Dependencies

- `golang.org/x/sys v0.38.0`: Unix and Windows system calls
- `gopkg.in/natefinch/lumberjack.v2 v2.2.1`: Log rotation and compression
