// Package main is a simple test app that exercises logger's LogChanHook.
// It initialises the logger with a caller-provided string channel and
// consumes every formatted log line that the hook delivers into that channel.
package main

import (
	"fmt"
	"sync"

	logger "github.com/rohanthewiz/logger"
)

func main() {
	// Create a buffered channel to receive formatted log messages.
	// The buffer prevents the non-blocking hook from dropping messages
	// during the short window before the consumer goroutine starts.
	logCh := make(chan string, 100)

	// WaitGroup to ensure the consumer drains the channel before exit.
	var wg sync.WaitGroup

	// Consume messages from the channel in a separate goroutine and
	// print each one so the output is visible when running the script.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range logCh {
			fmt.Print("[LogChanHook] ", msg)
		}
	}()

	// Initialise the logger with the LogChanHook enabled.
	// Only warn-and-above messages will be forwarded to the hook so we
	// can show that level filtering works alongside unrestricted stdout logging.
	logger.InitLog(logger.LogConfig{
		Formatter: "text",
		LogLevel:  "debug", // overall log level
		LogChanCfg: logger.LogChanCfg{
			Enabled:  true,
			Ch:       logCh,
			LogLevel: "warn", // hook only receives warn / error / fatal / panic
		},
	})

	fmt.Println("--- Logging messages at various levels ---")

	// debug and info: visible on stdout but NOT sent to the hook
	logger.Log("debug", "Debug message – hook should NOT receive this")
	logger.Log("info", "Info message – hook should NOT receive this",
		"key1", "value1")

	// warn and above: visible on stdout AND forwarded to the hook
	logger.Log("warn", "Warn message – hook SHOULD receive this",
		"component", "test_app")
	logger.Log("error", "Error message – hook SHOULD receive this",
		"component", "test_app", "detail", "something went wrong")

	// CloseLog flushes the async logger, then we close the hook channel
	// so the consumer goroutine can exit cleanly.
	logger.CloseLog()
	close(logCh)

	// Wait for the consumer to finish printing all messages.
	wg.Wait()

	fmt.Println("--- Done ---")
}
