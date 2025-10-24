package logger

import (
	"fmt"
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
func LogErr(err error, keyValPairs ...any) {
	logErrCore(err, keyValPairs...)
}

// Err is a convenience wrapper for LogErr
// We have to duplicate the function body or use a common core so as to keep error framelevels consistent
func Err(err error, keyValPairs ...any) {
	logErrCore(err, keyValPairs...)
}

// logErrCore is the common core for logging errors.
// This exists so as to keep framelevels consistent among calling functions
func logErrCore(err error, keyValPairs ...any) {
	if err == nil {
		Log(LogLevel.Info, "In LogErr Not logging a nil err", "called from",
			serr.FunctionLoc(serr.FrameLevels.FrameLevel3))
		return
	}

	flds := logrus.Fields{}

	ser, ok := err.(serr.SErr)
	if !ok { // make into SErr just for the sake of logging
		ser = serr.NewSerrNoContext(err)
	}

	// Add current location context so we don't have to wrap errors at the point of logging
	ser.AppendCallerContext(serr.FrameLevels.FrameLevel4)

	// Add any additional attributes
	ser.AppendAttributes(keyValPairs...)

	// Get all attributes from the error
	for key, anyArr := range ser.FieldsMapOfSliceOfAny() {
		strArr := make([]string, 0, len(anyArr))

		for _, a := range anyArr {
			strArr = append(strArr, fmt.Sprintf("%v", a))
		}

		strVal := strings.Join(strArr, " -> ")

		if key != "" {
			switch strings.ToLower(key) {
			case strings.ToLower(serr.UserMsgKey):
				continue // that one is for UI only
			case strings.ToLower(serr.UserMsgSeverityKey):
				continue // that one is for UI only
			case prefixKey:
				logPrefix = strVal
				continue
			default:
				flds[key] = strVal
			}
		}

		if _, ok := flds[errorKey]; !ok {
			flds[errorKey] = ser.Error()
		}
	}

	logrus.WithFields(flds).Error(logPrefix + err.Error())
}
