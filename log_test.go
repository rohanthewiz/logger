package logger

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/rohanthewiz/serr"
)

func TestLog(t *testing.T) {
	formatter := "text" // or json
	InitLog(LogConfig{
		Formatter: formatter,
		LogLevel:  "debug",
		SlackrusCfg: SlackrusCfg{
			Enabled: false,
			// Endpoint: sc.Endpoint,
			// LogLevel: sc.LogLevel,
		},
	})
	defer CloseLog()

	Log("info", "Conveying some info", "attribute1", "value1", "attribute2", "value2")
	// => {"attribute1":"value1","attribute2":"value2","level":"info","msg":"Conveying some info","time":"2024-05-11T19:30:09-05:00"}

	Log("error", "Some error occurred", "attribute1", "value1", "attribute2", "value2",
		"location", serr.FunctionLoc(serr.FrameLevels.FrameLevel1))
	// => {"attribute1":"value1","attribute2":"value2","level":"error","location":"logger/log_test.go:25","msg":"Some error occurred","time":"2024-05-11T19:30:09-05:00"

	// With regular error
	err := errors.New("This is the original error")

	// We can log a standard error, the message will be err.Error()
	LogErr(err)
	// => {"level":"error","msg":"This is the original error","time":"2024-05-11T19:30:09-05:00"}

	LogErr(err, "info", "message2", "key1", "value1", "key2", "value2")
	// => {"info":"message2","key1":"value1","key2":"value2","level":"error","msg":"This is the original error","time":"2024-05-11T19:30:09-05:00"}

	// Use the shorter form of LogErr
	Err(err, "message2", "key1", "value1", "key2", "value2")
	// => {"info":"message2","key1":"value1","key2":"value2","level":"error","msg":"This is the original error","time":"2024-05-11T19:30:09-05:00"}

	// Log a nil error
	LogErr(nil, "message")

	fmt.Println("------------------------------------------------------------------------------")
	// With SErr
	ser := serr.New("This is the original error", "err0Fld1", "errVal1", "errFld2", "errVal2")

	fmt.Println("\n--- ", serr.StringFromErr(ser), "---\n ")

	Err(ser)
	// {"errFld1":"errVal1","errFld2":"errVal2","error":"This is the original error","function":"rohanthewiz/logger.TestLog","level":"error","location":"logger/log_test.go:40","msg":"This is the original error","time":"2024-05-11T19:30:09-05:00"}

	Err(ser, "error", "Err: from my point of view")
	// {"errFld1":"errVal1","errFld2":"errVal2","error":"Err: from my point of view","function":"rohanthewiz/logger.TestLog","level":"error","location":"logger/log_test.go:40","msg":"This is the original error","time":"2024-05-11T19:30:09-05:00"}

	Err(ser, "info", "message2", "key1", "value1", "key2", "value2", "msg", "my message")
	// => {"errFld1":"errVal1","errFld2":"errVal2","error":"This is the original error","fields.msg":"my message","function":"rohanthewiz/logger.TestLog","info":"message2","key1":"value1","key2":"value2","level":"error","location":"logger/log_test.go:40","msg":"This is the original error","time":"2024-05-11T19:30:09-05:00"

	fmt.Println("\n--- Sending some logs async'ly", strings.Repeat("--", 40))
	LogAsync("info", "An Async message")
	LogAsync("error", "An Async error message")

	// See log_test.go for more examples

	// User printed stack trace
	PrintStackTrace()
}
