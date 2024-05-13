package logger

import "github.com/sirupsen/logrus"

type loggerLevels struct {
	Debug, Info, Warn, Error string
}

// Consumable log levels (convenience var)
var LogLevel = loggerLevels{
	Error: "error",
	Warn:  "warn",
	Info:  "info",
	Debug: "debug",
}

const (
	StrLevelError = "error"
	StrLevelWarn  = "warn"
	StrLevelInfo  = "info"
	StrLevelDebug = "debug"
)

// Internal only - keep private
var logrusLevels = map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
	"fatal": logrus.FatalLevel,
}
