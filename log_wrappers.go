package logger

import "fmt"

// F is a convenience function for log Info using a formatted string
func F(format string, args ...any) {
	Log(StrLevelInfo, fmt.Sprintf(format, args...))
}

func InfoF(format string, args ...any) {
	Log(StrLevelInfo, fmt.Sprintf(format, args...))
}
func DebugF(format string, args ...any) {
	Log(StrLevelDebug, fmt.Sprintf(format, args...))
}
func WarnF(format string, args ...any) {
	Log(StrLevelWarn, fmt.Sprintf(format, args...))
}
func ErrorF(format string, args ...any) {
	Log(StrLevelError, fmt.Sprintf(format, args...))
}

func Warn(msg string, args ...string) {
	Log(StrLevelWarn, msg, args...)
}

func Info(msg string, args ...string) {
	Log(StrLevelInfo, msg, args...)
}

func Debug(msg string, args ...string) {
	Log(StrLevelDebug, msg, args...)
}

// Error will create a new error based on msg
// If we are logging an existing error, better to use LogErr(err, ...)
func Error(msg string, args ...string) {
	Log(StrLevelError, msg, args...)
}

type Attr struct {
	Key   any
	Value any
}

func (attr Attr) String() []string {
	key := fmt.Sprintf("%v", attr.Key)
	val := fmt.Sprintf("%v", attr.Value)
	return []string{key, val}
}

func InfoAttrs(msg string, attrs ...Attr) {
	var strArgs []string
	for _, attr := range attrs {
		strArgs = append(strArgs, attr.String()...)
	}
	Log(StrLevelInfo, msg, strArgs...)
}

func DebugAttrs(msg string, attrs ...Attr) {
	var strArgs []string
	for _, attr := range attrs {
		strArgs = append(strArgs, attr.String()...)
	}
	Log(StrLevelDebug, msg, strArgs...)
}

func WarnAttrs(msg string, attrs ...Attr) {
	var strArgs []string
	for _, attr := range attrs {
		strArgs = append(strArgs, attr.String()...)
	}
	Log(StrLevelWarn, msg, strArgs...)
}

func ErrorAttrs(msg string, attrs ...Attr) {
	var strArgs []string
	for _, attr := range attrs {
		strArgs = append(strArgs, attr.String()...)
	}
	Log(StrLevelError, msg, strArgs...)
}
