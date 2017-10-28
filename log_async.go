package logger

import (
	"fmt"
	"time"
)

// Queue up logs on logsChannel. Level, message, and all arguments are strings
func LogAsync(level string, msg string, args... string) {
	LogAsyncBin(level, msg, nil, args...)
}

// The transport here allows for a binary third argument, all other arguments are strings
func LogAsyncBin(level string, msg string, bin *[]byte, args ...string) {
	blank_bytes := []byte{}
	if bin == nil {
		bin = &blank_bytes
	}
	args_slice := [][]byte{ []byte(level), []byte(msg), *bin } // first section of arguments
	for _, arg := range args {  // get the rest of the arguments
		args_slice = append(args_slice, []byte(arg))
	}
	// Lock in a sequence attribute here before the async call
	args_slice = append(args_slice, []byte("seq"))
	args_slice = append(args_slice, []byte(fmt.Sprintf("%d", time.Now().UnixNano())))

	logsWaitGroup.Add(1)  // track the number of log senders
	go func() {
		logsChannel <- args_slice // send to the channel.
		logsWaitGroup.Done()  // one less log sender
	}()
}

// Poll the LogsChannel for incoming messages of [][]byte
// Arguments are the receive only logs channel and send only done channel
func pollForLogs(logs_channel <-chan [][]byte, done chan <- bool) {
	defer func() {
		// Flush any hooks here
		done <- true  // signal caller when we are done
	}()
	for {
		select {  // Select can multiplex cases reading from multiple channels
		case attrs, ok := <- logs_channel:  // we will block till there is a message on the channel
			if !ok { // the channel is closed *and* empty, so wrap up
				return
			} else {
				LogBinary(string(attrs[0]), string(attrs[1]), attrs[2], attrs[3:]...) // receive the item and call Log()
			}
			// we don't timeout. Logs run for the life of the app
		}
	}
}