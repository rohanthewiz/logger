package logger

import (
	"strings"

	"github.com/rohanthewiz/serr"
	"github.com/sirupsen/logrus"
)

const (
	errorKey  = "error"
	prefixKey = "prefix"
)

// Logging for structured errors (SErr)
// Take an err (preferrably a SErr), and optional attributes as key value pairs
// *Note* `msg` key is added automatically by Logrus
// Example:
//
//		er := serr.New("Just testing an error", "attribute1", "value1", "attribute2", "value2")
//	logger.LogErr(er, "'attribute3", "value3", "attribute4", "value4")
//
// see the tests for more examples
func LogErr(err error, keyValPairs ...string) {
	if err == nil {
		Log(LogLevel.Info, "In LogErr Not logging a nil err", "called from", serr.FunctionLoc(serr.FrameLevels.FrameLevel1))
		return
	}

	flds := logrus.Fields{}

	ser, ok := err.(serr.SErr)
	if !ok { // make into SErr just for the sake of logging
		ser = serr.NewSerrNoContext(err)
	}

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

		if _, ok := flds[errorKey]; !ok {
			flds[errorKey] = ser.Error()
		}
	}

	logrus.WithFields(flds).Error(logPrefix + err.Error())
}

// Err is a convenience wrapper for LogErr
func Err(err error, keyValPairs ...string) {
	LogErr(err, keyValPairs...)
}
