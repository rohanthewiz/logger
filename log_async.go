package logger

// LogAsync adds a log message asynchronously to the queue (logsChannel)
// level can be one of "debug", "info", "warn", "error"
// msg is required and should be a simple description of the event we are logging
// args are (optional) pairs of of attribute, values
// Example: LogAsync("info", "Downloading a file", "filename", "export.xlsx")
func LogAsync(level, msg string, args ...string) {
	argsSlice := [][]byte{[]byte(level), []byte(msg)} // first section of arguments
	for _, arg := range args {                        // get the rest of the arguments
		argsSlice = append(argsSlice, []byte(arg))
	}

	logsWaitGroup.Add(1) // track the number of log senders
	go func() {
		logsChannel <- argsSlice // send to the channel.
		logsWaitGroup.Done()     // one less log sender
	}()
}

// Poll the LogsChannel for incoming messages of [][]byte
func pollForLogs(logsChan <-chan [][]byte, done chan<- bool) {
	defer func() {
		done <- true // signal caller when we are done // Flush any hooks here
	}()
	for {
		select { // select blocks until there is a message on a monitored channel
		case attrs, ok := <-logsChan:
			if !ok { // the channel is closed *and* empty, so wrap up
				return
			} else {
				logBytes(string(attrs[0]), string(attrs[1]), attrs[2:]...) // receive the item and call Log()
			}
		}
	}
}
