package log_chan

import (
	"time"

	"github.com/sirupsen/logrus"
)

// LogEntry represents a log event sent on the channel.
type LogEntry struct {
	Level   string
	Message string
	Time    time.Time
	Data    map[string]string
}

// LogChanHook is a logrus hook that sends log entries to a provided Go channel.
type LogChanHook struct {
	Channel        chan LogEntry
	AcceptedLevels []logrus.Level
}

var allLevels = []logrus.Level{
	logrus.DebugLevel,
	logrus.InfoLevel,
	logrus.WarnLevel,
	logrus.ErrorLevel,
	logrus.FatalLevel,
	logrus.PanicLevel,
}

// NewLogChanHook creates a new LogChanHook that sends entries to the given channel.
func NewLogChanHook(ch chan LogEntry, acceptedLevels []logrus.Level) *LogChanHook {
	return &LogChanHook{
		Channel:        ch,
		AcceptedLevels: acceptedLevels,
	}
}

// Levels returns the log levels this hook should fire on.
func (h *LogChanHook) Levels() []logrus.Level {
	if h.AcceptedLevels == nil {
		return allLevels
	}
	return h.AcceptedLevels
}

// Fire is called by logrus when a log event fires. It sends the entry to the channel.
func (h *LogChanHook) Fire(entry *logrus.Entry) error {
	data := make(map[string]string, len(entry.Data))
	for k, v := range entry.Data {
		if s, ok := v.(string); ok {
			data[k] = s
		}
	}

	le := LogEntry{
		Level:   entry.Level.String(),
		Message: entry.Message,
		Time:    entry.Time,
		Data:    data,
	}

	select {
	case h.Channel <- le:
	default:
		// Channel full — drop the entry to avoid blocking the logger
	}

	return nil
}
