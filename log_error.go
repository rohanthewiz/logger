package logger

import (
	"fmt"
	"time"
	"strings"
	"github.com/rohanthewiz/serr"
	"github.com/sirupsen/logrus"
)

// Special logging for errors and structured errors (github.com/rohanthewiz/serr)
func LogErr(err error, mesg...string) {
	msgs := []string{}  // accumulate "msg" fields
	errs := []string{}  // accumulate "error" fields

	// Add standard logging fields
	flds := logrus.Fields{"level": "error"}
	if seq, ok := flds["sequence"]; !ok || seq == "" {  // set a sequence if not already set
		flds["sequence"] = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	if app, ok := flds["app"]; !ok || app == "" {  // and do both "env" and "app" together
		flds["app"] = logOptions.AppName
		flds["environment"] = logOptions.Environment
	}

	// Add mesg if supplied
	if len(mesg) == 1 {  // if single item, add it to mesg
		msgs = []string{mesg[0]}
	}
	// If multiple mesgs supplied wrap the error with them
	if len(mesg) > 1 {
		err = serr.Wrap(err, mesg...)
	}

	// Add error string from original error
	if er := err.Error(); er != "" {
		errs = []string{er}
	}

	// If error is structured error, get key vals
	if ser, ok := err.(serr.SErr); ok {
		for key, val := range ser.FieldsMap() {
			if key != "" {
				switch strings.ToLower(key) {
				case "error":
					errs = append(errs, val)
				case "msg":
					msgs = append(msgs, val)
				default:
					flds[key] = val
				}
			}
		}
	}
	// message is required by logrus so use the original error string if msgs empty
	if len(msgs) == 0 {
		msgs = []string{err.Error()}
	}
	// Populate the "error" field
	if len(errs) > 0 {
		flds["error"] = strings.Join(errs, " - ")
	}
	// Log it
	logrus.WithFields(flds).Error(strings.Join(msgs, " - "))
}
