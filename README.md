## Logging package

### Usage

```go
package main
import (
    "errors"
    . "github.com/rohanthewiz/logger
    "github.com/rohanthewiz/serr
)

func main() {
	InitLog(LogOptions{
		AppName: "My App",
		Environment: "development",
		Format: "text",
		Level: "info",
		InfoPath: "log/info.log",
		ErrorPath: "log/error.log",
	})
	defer CloseLog()

    Log("info", "Conveying some info", "attribute1", "value1", "attribute2", "value2")
    Log("error", "Some error occurred", "attribute1", "value1", "attribute2", "value2")

	err := errors.New("This is the original error")

	// We can log a standard error, the message will be err.Error()
	LogErr(err)
		//=> ERRO[0000] This is the original error	error="This is the original error"

	// Single argument after err becomes part of logrus message
	LogErr(err, "Custom message here")
		//=> ERRO[0000] Custom message here   error="I'm making this error up :-)"

	// Multiple arguments after err are treated as a key, value list and will wrap the error
	LogErr(err, "key1", "value1", "key2", "value2")
		//=> ERRO[0000] This is the original error error="This is the original error" key1=value1 key2=value2

	// Multiple arguments after err are treated as a key, value list and will wrap the error
	LogErr(err, "error", "Error Code: ABCDE321", "msg", "This is a critical error", "key2", "value2")
		//=> ERRO[0000] This is a critical error  error="This is the original error - Error Code: ABCDE321" key2=value2

	err2 := serr.Wrap(err, "Gosh! We got an error!")
	LogErr(err2)
		//=> ERRO[0000] Gosh! We got an error!	error="This is the original error"
	LogErr(err2, "Additional info")
		//=> ERRO[0000] Additional info - Gosh! We got an error!	error="This is the original error"

	// We can log an SErr wrapped error
	err3 := serr.Wrap(err, "cat", "aight", "dogs", "I dunno")
	LogErr(err3, "Animals, do we really need them? Yes!!")
		//=> ERRO[0000] Animals, do we really need them? Yes!! cat=aight dogs="I dunno" error="I'm making this error up :-)"
}
```
