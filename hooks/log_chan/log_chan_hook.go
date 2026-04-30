package log_chan

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// LogChanHook is a logrus hook that formats log entries as text
// and sends them to a caller-provided string channel. This enables
// consumers to receive structured log output without coupling to
// a specific transport (e.g. Teams, Slack) — the caller decides
// what to do with each formatted log line.
type LogChanHook struct {
	Ch             chan string      // destination channel for formatted log messages
	AcceptedLevels []logrus.Level  // levels that trigger this hook; nil means all levels
	Disabled       bool            // allows the hook to be temporarily silenced
	formatter      logrus.Formatter // text formatter used to serialize log entries
}

// allLevels enumerates every logrus severity so we have a default
// when the caller does not restrict accepted levels.
var allLevels = []logrus.Level{
	logrus.TraceLevel,
	logrus.DebugLevel,
	logrus.InfoLevel,
	logrus.WarnLevel,
	logrus.ErrorLevel,
	logrus.FatalLevel,
	logrus.PanicLevel,
}

// NewLogChanHook creates a LogChanHook that writes logrus-text-formatted
// messages into ch. Pass nil for acceptedLevels to receive all levels.
func NewLogChanHook(ch chan string, acceptedLevels []logrus.Level) *LogChanHook {
	return &LogChanHook{
		Ch:             ch,
		AcceptedLevels: acceptedLevels,
		formatter:      &logrus.TextFormatter{DisableColors: true},
	}
}

// Levels returns the set of log levels this hook responds to.
// Required by the logrus.Hook interface.
func (h *LogChanHook) Levels() []logrus.Level {
	if h.AcceptedLevels == nil {
		return allLevels
	}
	return h.AcceptedLevels
}

// AllowedLevels returns every logrus level at or above the given level.
// "Above" means more severe — e.g. AllowedLevels(WarnLevel) returns
// [Warn, Error, Fatal, Panic]. This mirrors the helper in teams_log.
func AllowedLevels(lvl logrus.Level) []logrus.Level {
	for i := range allLevels {
		if allLevels[i] == lvl {
			return allLevels[i:]
		}
	}
	return []logrus.Level{}
}

// Fire formats the log entry as logrus text and sends it to the channel.
// A non-blocking send is used so a slow or full channel does not block
// the logging goroutine — messages are dropped with a stderr warning
// if the channel cannot accept them immediately.
// Required by the logrus.Hook interface.
func (h *LogChanHook) Fire(entry *logrus.Entry) error {
	if h.Disabled {
		return nil
	}

	// Format the entry using the logrus text formatter
	formatted, err := h.formatter.Format(entry)
	if err != nil {
		fmt.Printf("log_chan: failed to format log entry: %v\n", err)
		return nil // don't propagate formatter errors to logrus
	}

	// Non-blocking send: drop the message rather than stall the caller
	select {
	case h.Ch <- string(formatted):
	default:
		fmt.Println("log_chan: channel full, dropping log message")
	}

	return nil
}
