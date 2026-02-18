## Logging package

### Usage

```go
package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rohanthewiz/logger"
)

func main() {
	logger.InitLog(logger.LogConfig{
		Formatter: "json", // or "text"
		LogLevel:  "debug",
		/*		SlackrusCfg: logger.SlackrusCfg{
				Enabled: false,
				Endpoint: "<Endpoint>",
				LogLevel: "<LogLevel>",
			}, */
		/*		SlackAPICfg: logger.SlackAPICfg{
				Enabled:   true,
				Token:     "xoxb-your-bot-token", // Bot User OAuth Token
				Channel:   "C086KQQUGCW",         // Channel ID
				LogLevel:  "error",               // Minimum level to send to Slack
				UseBlocks: true,                  // Use rich block formatting
			}, */
	})
	defer logger.CloseLog()

	logger.Log("info", "Conveying some info", "attribute1", "value1", "attribute2", "value2")
	// => {"attribute1":"value1","attribute2":"value2","level":"info","msg":"Conveying some info","time":"2024-05-13T13:18:43-05:00"}

	logger.Log("error", "Some error occurred", "attribute1", "value1", "attribute2", "value2")
	// => {"attribute1":"value1","attribute2":"value2","level":"error","msg":"Some error occurred","time":"2024-05-13T13:18:43-05:00"}

	fmt.Println("\n--- Using the convenience functions", strings.Repeat("-", 50))

	logger.Info("Conveying some info", "attribute1", "value1", "attribute2", "value2")
	// => {"attribute1":"value1","attribute2":"value2","level":"info","msg":"Conveying some info","time":"2024-05-13T13:18:43-05:00"}

	logger.Debug("something interesting happened here")
	// => {"level":"debug","msg":"something interesting happened here","time":"2024-05-13T13:49:13-05:00"}

	fmt.Println("\n--- Using LogErr to unpack a structured error", strings.Repeat("--", 50))

	err := errors.New("this is the original error")
	logger.LogErr(err)
	// => {"level":"error","msg":"this is the original error","time":"2024-05-13T13:18:43-05:00"}

	// Log with a single message
	logger.LogErr(err, "message")
	// => {"error":"this is the original error","fields.msg":"message","level":"error","msg":"this is the original error","time":"2024-05-13T13:18:43-05:00"}

	// Multiple arguments after message are treated as a key, value list and will wrap the error
	// If an odd number of fields are provided, the first field is treated as the message.
	logger.Err(err, "message", "key1", "value1", "key2", "value2")
	// => {"error":"this is the original error","fields.msg":"message","key1":"value1","key2":"value2","level":"error","msg":"this is the original error","time":"2024-05-13T13:18:43-05:00"

	// See log_test.go for more examples

	// [At close] {"level":"info","msg":"Logs gracefully shutdown","time":"2024-05-13T14:24:40-05:00"}
}

```

### Slack API Hook

The logger package now includes a Slack API hook that sends log messages directly to Slack using the Web API. This provides more flexibility than webhook-based approaches.

#### Features

- **Simple Messages**: Send basic formatted log messages to Slack
- **Rich Block Formatting**: Create visually appealing messages with headers, fields, and action buttons
- **Structured Data Support**: Automatically formats log fields into organized sections
- **Error Context**: Special handling for error logs with stack traces and debugging links
- **Asynchronous Sending**: Non-blocking log dispatch to maintain application performance

#### Configuration

```go
logger.InitLog(logger.LogConfig{
	Formatter: "json",
	LogLevel:  "info",
	SlackAPICfg: logger.SlackAPICfg{
		Enabled:   true,
		Token:     "xoxb-your-bot-token", // Bot User OAuth Token from Slack App
		Channel:   "C086KQQUGCW",         // Channel ID where logs should be sent
		LogLevel:  "error",               // Minimum level to send to Slack
		UseBlocks: true,                  // Enable rich block formatting
	},
})
```

#### Usage Examples

```go
// Simple error logging
logger.Error("service", "payment-api", "environment", "production", 
	"Payment processing failed")

// Rich error with debugging information
logger.Error(
	"service", "auth-service",
	"environment", "production",
	"error_type", "TokenValidationError",
	"user_id", "12345",
	"stack_trace", "goroutine 1 [running]:\nmain.validateToken(...)",
	"log_url", "https://logs.example.com/search?id=ERR-123",
	"incident_id", "INC-456",
	"JWT token validation failed",
)
```

### LogChan Hook

The LogChan hook sends logrus-text-formatted log messages to a caller-provided `chan string`. This lets you route logs to any custom consumer — a UI, a database writer, a remote forwarder, etc. — without coupling to a specific transport.

#### Configuration

```go
// Create a buffered channel to receive log messages
logCh := make(chan string, 100)

logger.InitLog(logger.LogConfig{
	Formatter: "text",
	LogLevel:  "debug",
	LogChanCfg: logger.LogChanCfg{
		Enabled:  true,
		Ch:       logCh,
		LogLevel: "warn", // Only send warn and above to the channel
	},
})
defer logger.CloseLog()
```

#### Consuming Log Messages

```go
// Start a goroutine to process incoming log lines
go func() {
	for msg := range logCh {
		fmt.Print("received: ", msg) // each msg is a full text-formatted log line
	}
}()

logger.Warn("disk usage high", "percent", "92")
logger.Error("connection lost", "host", "db-primary")
```

The hook performs a non-blocking send — if the channel buffer is full, messages are dropped (with a warning to stdout) rather than blocking the logging goroutine.

#### Slack App Setup

1. Create a Slack App at https://api.slack.com/apps
2. Add the `chat:write` OAuth scope to your Bot Token
3. Install the app to your workspace
4. Copy the Bot User OAuth Token (starts with `xoxb-`)
5. Invite the bot to the desired channel `/invite @bot_name`
6. Get the Channel ID where you want logs sent
