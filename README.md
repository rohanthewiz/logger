## Logging package

### Usage

```go
import "github.com/rohanthewiz/logger

func ExerciseLogging() {
    // We can log a standard error, the message will be err.Error()
    err := errors.New("I'm making this error up :-)")
    logger.LogErr(err)
      //=> ERRO[0000] I'm making this error up :-)     error="I'm making this error up :-)"
    // We can log with a custom message  
    logger.LogErr(err, "Custom message here")
      //=> ERRO[0000] Custom message here   error="I'm making this error up :-)"
    
    // We can log an SErr wrapped error
    err2 := serr.Wrap(err, "cat", "aight", "dogs", "I dunno")
    logger.LogErr(err2, "Animals, do we really need them? Yes!!")
       //=> ERRO[0000] Animals, do we really need them? Yes!! cat=aight dogs="I dunno" error="I'm making this error up :-)"
}
```
