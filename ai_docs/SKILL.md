---
name: structured-error-logger
description: "logger is a structured error logger that integrates well with github.com/rohanthewiz/serr. logger.LogErr will write out all the context stored in the SErr."
---

# Logger Package

A structured logging package for Go that wraps `logrus` with enhanced support for `github.com/rohanthewiz/serr` structured errors.

## Installation

```go
import "github.com/rohanthewiz/logger"
```

## Initialization

Initialize the logger at application startup:

```go
logger.InitLog(logger.LogConfig{
    Formatter: "json",  // "json" or "text"
    LogLevel:  "debug", // "debug" | "info" | "warn" | "error"
})
defer logger.CloseLog() // Ensure all async logs are flushed
```

## Core Functions

### Basic Logging with Key-Value Pairs

The `Log` function accepts a level, message, and optional key-value pairs:

```go
logger.Log("info", "Conveying some info", "attribute1", "value1", "attribute2", "value2")
// => {"attribute1":"value1","attribute2":"value2","level":"info","msg":"Conveying some info","time":"..."}
```

### Convenience Wrappers

Use wrapper functions for cleaner code:

```go
logger.Info("Simple info message", "key1", "value1", "key2", "value2")
logger.Debug("Simple debug message", "key1", "value1")
logger.Warn("Simple warning message", "key1", "value1")
logger.Error("Simple error message", "key1", "value1")
```

### Formatted Logging

Printf-style formatting is available:

```go
logger.F("Number: %d String: %s", 42, "test") // log info-level message
logger.InfoF("Info message with %s and %d", "string", 123)
logger.DebugF("Debug message with %.2f", 3.14159)
logger.WarnF("Warning: %s occurred %d times", "event", 5)
logger.ErrorF("Error code: %x", 255)
```

## Error Logging with SErr Integration

The key feature of this package is its awareness of `github.com/rohanthewiz/serr` structured errors. When logging errors, use `LogErr` or its alias `Err`:

### Standard Error Logging

```go
err := errors.New("This is the original error")

// Log a standard error - message will be err.Error()
logger.LogErr(err)
// => {"level":"error","msg":"This is the original error","time":"..."}

// With additional context
logger.LogErr(err, "info", "message2", "key1", "value1", "key2", "value2")
// => {"info":"message2","key1":"value1","key2":"value2","level":"error","msg":"This is the original error","time":"..."}

// Shorter form using Err
logger.Err(err, "message2", "key1", "value1")
```

### SErr Structured Error Logging

When logging a `serr.SErr`, all embedded context is automatically extracted and included in the log output:

```go
ser := serr.New("This is the original error", "errFld1", "errVal1", "errFld2", "errVal2")

// All SErr fields are extracted and logged
logger.Err(ser)
// => {"errFld1":"errVal1","errFld2":"errVal2","error":"This is the original error",
//     "function":"package.Function","level":"error","location":"file.go:40",
//     "msg":"This is the original error","time":"..."}

// Add additional context at log time
logger.Err(ser, "error", "Err: from my point of view")
// => {"errFld1":"errVal1","errFld2":"errVal2","error":"Err: from my point of view",
//     "function":"package.Function","level":"error","location":"file.go:40",
//     "msg":"This is the original error","time":"..."}

// Combine SErr fields with additional key-value pairs
logger.Err(ser, "info", "message2", "key1", "value1", "key2", "value2")
```

### Key Behavior Notes

- `LogErr` and `Err` automatically append caller context (function name and location)
- All attributes stored in a `serr.SErr` are unpacked into log fields
- Multiple values for the same attribute are joined with ` -> `
- `UserMsg` and `UserMsgSeverity` keys are excluded (reserved for UI)
- Logging a `nil` error logs an info message about the nil error instead

## Async Logging

For non-blocking logging:

```go
logger.LogAsync("info", "An Async message")
logger.LogAsync("error", "An Async error message", "key1", "value1")

// Shorter form
logger.Async("info", "An Async message")
```

## Configuration Options

```go
type LogConfig struct {
    EnvPrefix   string      // Prefix for all log messages
    Formatter   string      // "text" | "json"
    LogLevel    string      // "debug" | "info" | "warn" | "error"
    LogChanSize int         // Buffer size for async logs (default: 2000)
    TeamsLogCfg TeamsLogCfg // Microsoft Teams integration
    SlackAPICfg SlackAPICfg // Slack integration
}
```

### Microsoft Teams Integration

```go
logger.InitLog(logger.LogConfig{
    Formatter: "json",
    LogLevel:  "debug",
    TeamsLogCfg: logger.TeamsLogCfg{
        Enabled:  true,
        Endpoint: "https://your-teams-webhook-url",
        LogLevel: "warn", // Only log warn and above to Teams
    },
})
```

### Slack Integration

```go
logger.InitLog(logger.LogConfig{
    Formatter: "json",
    LogLevel:  "debug",
    SlackAPICfg: logger.SlackAPICfg{
        Enabled:   true,
        Token:     "xoxb-your-token",
        Channel:   "C086K...",
        LogLevel:  "warn",
        UseBlocks: true, // Rich block formatting
    },
})
```

## Log Levels

Available via `logger.LogLevel`:

```go
logger.LogLevel.Debug // "debug"
logger.LogLevel.Info  // "info"
logger.LogLevel.Warn  // "warn"
logger.LogLevel.Error // "error"
```

Or use string constants:

```go
logger.StrLevelDebug // "debug"
logger.StrLevelInfo  // "info"
logger.StrLevelWarn  // "warn"
logger.StrLevelError // "error"
```

## Utility Functions

### Stack Trace

Print the current stack trace for debugging:

```go
logger.PrintStackTrace()
```

### Location Helper

Add location context to log messages:

```go
logger.Log("error", "Some error occurred", "location", serr.FunctionLoc(serr.FrameLevels.FrameLevel1))
// => {"level":"error","location":"package/file.go:25","msg":"Some error occurred","time":"..."}
```

## Best Practices

1. Always call `defer logger.CloseLog()` after `InitLog()` to ensure async logs are flushed
2. Use `logger.LogErr(err)` or `logger.Err(err)` for error logging to capture full context
3. Prefer `serr.Wrap(err)` to wrap errors before logging for maximum context
4. Use structured key-value pairs instead of string formatting for searchable logs
5. Set appropriate log levels for different environments (debug for dev, info/warn for prod)
