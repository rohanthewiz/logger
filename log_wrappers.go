package logger

import (
	"fmt"
)

func Warn(msg string, args ...any) {
	Log(StrLevelWarn, msg, strArrayFromAnyArgs(args...)...)
}

func Info(msg string, args ...any) {
	Log(StrLevelInfo, msg, strArrayFromAnyArgs(args...)...)
}

func Debug(msg string, args ...any) {
	Log(StrLevelDebug, msg, strArrayFromAnyArgs(args...)...)
}

// Error will create a new error based on msg.
// It is better to use LogErr(err, ...) if you are logging an existing error,
// as that will include any error stack trace from serr
func Error(msg string, args ...any) {
	Log(StrLevelError, msg, strArrayFromAnyArgs(args...)...)
}

// F is a convenience function for log Info using a formatted string
func F(format string, args ...any) {
	Log(StrLevelInfo, fmt.Sprintf(format, args...))
}

// InfoF is a convenience function for log Info using a formatted string
func InfoF(format string, args ...any) {
	Log(StrLevelInfo, fmt.Sprintf(format, args...))
}

// DebugF is a convenience function for log Debug using a formatted string
func DebugF(format string, args ...any) {
	Log(StrLevelDebug, fmt.Sprintf(format, args...))
}

// WarnF is a convenience function for log Warn using a formatted string
func WarnF(format string, args ...any) {
	Log(StrLevelWarn, fmt.Sprintf(format, args...))
}

// ErrorF is a convenience function for log Error using a formatted string
func ErrorF(format string, args ...any) {
	Log(StrLevelError, fmt.Sprintf(format, args...))
}
