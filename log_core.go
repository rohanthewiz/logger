package logger

import (
	"fmt"
	"time"
	"strings"
	"github.com/sirupsen/logrus"
)

// Our general logging functions
// Level can be one of "debug", "info", "warn", "error", "fatal"
// `args` is a list of argument pairs
// `msg` can be empty if `event` key and value present
// `bin` can transport binary data

// Use this method for synchronously logging strings
func Log(level, msg string, args... string) {
	LogCore(level, msg, nil, args...)
}

// This is generally a landing point for async logging which may include a binary third argument
func LogBinary(level string, msg string, bin []byte, args... []byte) {
	str_args := []string{}
	for _, arg := range args {
		str_args = append(str_args, string(arg))
	}
	LogCore(level, msg, &bin, str_args...)
}

// Build flds and msg for logrus
func LogCore(level, msg string, bin *[]byte, args... string) {
	lvl := strings.ToLower(level)
	flds := logrus.Fields{"level": lvl}
	if bin != nil {
		flds["bin"] = bin
	}

	// Gather the other keys and values
	key := ""
	keys := []string{}
	for i, arg := range args {
		if i % 2 == 0 {  // arg is a key
			key = arg
			keys = append(keys, key)
		} else {
			flds[string(key)] = arg
		}
	}
	// Fixup / Validate
	if len(args) % 2 != 0 {
		logrus.Warn("Even number of meta arguments required to Log function. Odd argument will be paired with a blank")
	}
	if seq, ok := flds["sequence"]; !ok || seq == "" {  // set a sequence if not already set
		flds["sequence"] = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	// msg can be empty if event present - msg will be set to event value
	if msg == "" && flds["event"] != "" {
		msg = flds["event"].(string)
	}

	// If this is an external message relay, add our app name and env
	if app, ok := flds["app"]; !ok || app == "" {  // and do both "env" and "app" together
		flds["app"] = logOptions.AppName
		flds["environment"] = logOptions.Environment
	}

	// Call the logger
	lg := logrus.WithFields(flds)
	switch lvl {
	case "debug":
		lg.Debug(msg)
	case "info":
		lg.Info(msg)
	case "warn":
		lg.Warn(msg)
	case "error":
		lg.Error(msg)  // Log error, but don't quit
	case "fatal":
		lg.Fatal(msg)  // Calls os.Exit(1) after logging
	}
}
