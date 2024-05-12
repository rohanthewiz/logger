package logger

import (
	"strings"

	"github.com/rohanthewiz/serr"
	"github.com/sirupsen/logrus"
)

const (
	errorKey    = "error"
	msgKey      = "msg"
	messagesKey = "messages"
	prefixKey   = "prefix"
)

// Logging for structured errors (SErr)
// Take an err (preferrably a SErr), and optional attributes as key value pairs
// *Note* `msg` key is added automatically by Logrus
// Example:
//
//		er := serr.New("Just testing an error", "attribute1", "value1", "attribute2", "value2")
//	logger.LogErr(er, "'attribute3", "value3", "attribute4", "value4")
//	 logger.LogErr(er)
//	 logger.LogErr(er, "msg", "error message")
func LogErr(err error, keyValPairs ...string) {
	if err == nil {
		Log(LogLevel.Info, "In LogErr Not logging a nil err", "called from", serr.FunctionLoc(serr.FrameLevels.FrameLevel1))
		return
	}
	//
	// var errs []string // accumulate "error" fields, to include the inner error
	// var msgs []string // for accumulating "msg" fields which become like a description of the error

	// // Add error string from original error
	// if er := err.Error(); er != "" {
	// 	errs = []string{er}
	// }

	flds := logrus.Fields{}

	// If error is structured error, get key vals
	if ser, ok := err.(serr.SErr); ok {
		// Add any additional attributes
		ser.AppendKeyValPairs(keyValPairs...)

		// Get all attributes from the error
		for key, val := range ser.FieldsMap() {
			if key != "" {
				switch strings.ToLower(key) {
				case strings.ToLower(serr.UserMsgKey):
					continue // that one is for UI only
				case strings.ToLower(serr.UserMsgSeverityKey):
					continue // that one is for UI only
				case prefixKey:
					logPrefix = val
					continue
				default:
					flds[key] = val
				}
			}
		}

		/*	  Seems like logrus already takes care of this	// move any `msg` to new key `msgs`
		if val, ok := flds[msgKey]; !ok {
			flds[messagesKey] = val
		}
		*/
		if _, ok := flds[errorKey]; !ok {
			flds[errorKey] = ser.Error()
		}

	} else { // not an SErr
		key := ""
		for i, str := range keyValPairs {
			if i%2 == 0 { // even position is a key
				key = str
			} else {
				flds[key] = str
			}
		}

		// Fixup / Validate
		if len(keyValPairs)%2 != 0 {
			logrus.Warn("It is best to use pairs of key values with the LogErr function")
		}

	}

	logrus.WithFields(flds).Error(logPrefix + err.Error())
}
