## Logging package

### Usage

```go
package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rohanthewiz/logger"
)

func main() {
	logger.InitLog(logger.LogConfig{
		Formatter: "json", // or "text"
		LogLevel:  "debug",
		/*		SlackrusCfg: logger.SlackrusCfg{
				Enabled: false,
				Endpoint: "<Endpoint>",
				LogLevel: "<LogLevel>",
			}, */
	})
	defer logger.CloseLog()

	logger.Log("info", "Conveying some info", "attribute1", "value1", "attribute2", "value2")
	// => {"attribute1":"value1","attribute2":"value2","level":"info","msg":"Conveying some info","time":"2024-05-13T13:18:43-05:00"}

	logger.Log("error", "Some error occurred", "attribute1", "value1", "attribute2", "value2")
	// => {"attribute1":"value1","attribute2":"value2","level":"error","msg":"Some error occurred","time":"2024-05-13T13:18:43-05:00"}

	fmt.Println("\n--- Using the convenience functions", strings.Repeat("-", 50))

	logger.Info("Conveying some info", "attribute1", "value1", "attribute2", "value2")
	// => {"attribute1":"value1","attribute2":"value2","level":"info","msg":"Conveying some info","time":"2024-05-13T13:18:43-05:00"}

	logger.Error("Some error occurred", "attribute1", "value1", "attribute2", "value2")
	// => {"attribute1":"value1","attribute2":"value2","level":"error","msg":"Some error occurred","time":"2024-05-13T13:18:43-05:00"}

	logger.Debug("something interesting happened here")
	// => {"level":"debug","msg":"something interesting happened here","time":"2024-05-13T13:49:13-05:00"}

	fmt.Println("\n--- Using LogErr to unpack a structured error", strings.Repeat("--", 50))

	err := errors.New("this is the original error")
	logger.LogErr(err)
	// => {"level":"error","msg":"this is the original error","time":"2024-05-13T13:18:43-05:00"}

	// Log with a single message
	logger.LogErr(err, "message")
	// => {"error":"this is the original error","fields.msg":"message","level":"error","msg":"this is the original error","time":"2024-05-13T13:18:43-05:00"}

	// Multiple arguments after message are treated as a key, value list and will wrap the error
	// If an odd number of fields are provided, the first field is treated as the message.
	logger.LogErr(err, "message", "key1", "value1", "key2", "value2")
	// => {"error":"this is the original error","fields.msg":"message","key1":"value1","key2":"value2","level":"error","msg":"this is the original error","time":"2024-05-13T13:18:43-05:00"

	// See log_test.go for more examples

	// [At close] {"level":"info","msg":"Logs gracefully shutdown","time":"2024-05-13T14:24:40-05:00"}
}

```
