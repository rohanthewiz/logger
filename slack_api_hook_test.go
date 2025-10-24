package logger

import (
	"testing"

	"github.com/rohanthewiz/serr"
	"github.com/sirupsen/logrus"
)

// TestSlackAPIHook demonstrates usage of the Slack API hook
// To run this test with actual Slack integration:
//  1. Set environment variables:
//     export SLACK_API_TOKEN="xoxb-your-bot-token"
//     export SLACK_CHANNEL="C086KQQUGCW"
//  2. Run: go test -v -run TestSlackAPIHook
func TestSlackAPIHook(t *testing.T) {
	// Skip if no Slack token is provided
	token := "bogus-token"    // Replace with os.Getenv("SLACK_API_TOKEN") for real testing
	channel := "test-channel" // Replace with os.Getenv("SLACK_CHANNEL") for real testing

	if token == "test-token" {
		t.Skip("Skipping Slack API test - no real token provided")
	}

	// Test configuration with simple messages
	t.Run("SimpleMessages", func(t *testing.T) {
		// Initialize logger with Slack API hook for simple messages
		InitLog(LogConfig{
			Formatter: "text",
			LogLevel:  "debug",
			SlackAPICfg: SlackAPICfg{
				Enabled:   true,
				Token:     token,
				Channel:   channel,
				LogLevel:  "info",
				UseBlocks: false, // Simple messages
			},
		})
		defer CloseLog()

		// Test various log levels
		Debug("This is a debug message", "debug_key", "debug_value")
		Info("Payment processed successfully", "service", "payment-gateway", "environment", "staging")
		Warn("API response time is high", "latency_ms", 1500)
		Error("Failed to connect to database", "error_type", "DatabaseConnectionError")

		// Test with structured error from serr
		err := serr.New("connection timeout",
			"host", "db.example.com",
			"port", "5432",
			"retry_count", "3")
		LogErr(err, "Database connection failed after retries")
	})

	// Test configuration with block messages
	t.Run("BlockMessages", func(t *testing.T) {
		// Re-initialize with block formatting
		logrus.StandardLogger().Hooks = make(logrus.LevelHooks) // Clear existing hooks

		InitLog(LogConfig{
			Formatter: "json",
			LogLevel:  "debug",
			SlackAPICfg: SlackAPICfg{
				Enabled:   true,
				Token:     token,
				Channel:   channel,
				LogLevel:  "warn",
				UseBlocks: true, // Rich block messages
			},
		})
		defer CloseLog()

		// Test warning with multiple fields
		Warn(
			"service", "authentication-service",
			"environment", "production",
			"component", "jwt-validator",
			"user_id", "12345",
			"ip_address", "192.168.1.100",
			"Token validation failed - invalid signature",
		)

		// Test error with stack trace
		Error(
			"service", "payment-gateway",
			"environment", "production",
			"error_type", "PaymentProcessingError",
			"transaction_id", "TXN-789456",
			"amount", 99.99,
			"currency", "USD",
			"stack_trace", "goroutine 1 [running]:\nmain.processPayment(0xc000010200)\n\t/app/payment.go:45 +0x123\nmain.main()\n\t/app/main.go:20 +0x45",
			"log_url", "https://logs.example.com/search?id=ERR-789456",
			"incident_id", "INC-123",
			"Payment processing failed due to gateway timeout",
		)

		// Test structured error with rich context
		err := serr.New("database query failed",
			"query", "SELECT * FROM users WHERE id = ?",
			"table", "users",
			"timeout_seconds", "30",
			"error_code", "TIMEOUT_ERROR")

		LogErr(err,
			"service", "user-service",
			"environment", "production",
			"module", "user-repository",
			"Failed to fetch user data",
		)
	})
}

// Example function showing how to configure Slack API hook in a real application
func ExampleInitLog_slackAPI() {
	// Example configuration for production use
	config := LogConfig{
		Formatter: "json",
		LogLevel:  "info",
		SlackAPICfg: SlackAPICfg{
			Enabled:   true,
			Token:     "xoxb-your-bot-token", // Bot User OAuth Token from Slack App
			Channel:   "C086KQQUGCW",         // Channel ID where logs should be sent
			LogLevel:  "error",               // Only send errors and above to Slack
			UseBlocks: true,                  // Use rich formatting for better readability
		},
	}

	InitLog(config)
	defer CloseLog()

	// Your application code here
	Info("Application started successfully")
	Error("service", "api", "Database connection lost")
}
