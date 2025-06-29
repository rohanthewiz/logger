# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
This is a Go logging package that wraps the `logrus` library with additional features including:
- Structured logging with key-value pairs
- Error logging with structured error support via `github.com/rohanthewiz/serr`
- Asynchronous logging capabilities
- Slack and Microsoft Teams integration
- Convenience wrapper functions for common logging patterns

## Development Commands

### Build
```bash
go build ./...
```

### Run Tests
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -v -run TestLog
```

### Dependencies
```bash
# Update dependencies
go mod tidy

# Download dependencies
go mod download
```

## Code Architecture

### Core Components
- **log_core.go**: Main logging functionality and core Log() function
- **log_setup.go**: InitLog() and CloseLog() functions for initialization
- **log_wrappers.go**: Convenience functions (Info, Debug, Warn, Error, Fatal) that wrap Log()
- **log_error.go**: LogErr() and Err() functions for structured error logging
- **log_async.go**: Asynchronous logging implementation
- **levels.go**: Log level constants and mappings
- **configs.go**: Configuration structures (LogConfig, SlackrusCfg)

### Key Design Patterns
1. **Variadic Arguments**: Functions accept `...any` for flexible key-value pair logging
2. **Structured Errors**: Integration with `serr` package for rich error context
3. **Wrapper Pattern**: Convenience functions wrap the core Log() function
4. **Configuration**: Single LogConfig struct controls formatter, level, and integrations

### Testing Approach
- Uses standard Go testing package
- Tests demonstrate proper usage patterns
- Focus on testing various argument combinations and error scenarios

## Important Notes
- When modifying wrapper functions, ensure they properly handle variadic arguments
- Error logging functions (LogErr, Err) unpack structured errors from the `serr` package
- The package supports both JSON and text formatters
- Async logging must be properly closed with CloseLog() to ensure all logs are flushed