package logger

import (
	"github.com/rohanthewiz/serr"
	"github.com/sirupsen/logrus"
	"strings"
)

// Log Error (possibly an SErr - structured error) with optional mesg
func LogErr(err error, mesg ...string) {
	var msgs []string  // accumulate "msg" fields
	var errs []string  // accumulate "errs" fields

	flds := logrus.Fields{}

	// Add mesg if supplied
	if len(mesg) > 0 {
		msgs = []string{mesg[0]}
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
