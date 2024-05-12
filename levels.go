package logger

import "github.com/sirupsen/logrus"

type loggerLevels struct {
	Debug, Info, Warn, Error string
}

// Consumable log levels (convenience var)
var LogLevel = loggerLevels{
	Debug: "debug",
	Info:  "info",
	Warn:  "warn",
	Error: "error",
}

// Even more convenient
const (
	Error = "error"
	Warn  = "warn"
	Info  = "info"
	Debug = "debug"
)

// Internal only - keep private
var logrusLevels = map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
	"fatal": logrus.FatalLevel,
}
