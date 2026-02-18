package logger

import (
	"sync"
)

const (
	defaultLogLevel         = "debug" //  "debug | info | warn | error"
	defaultTeamsLogLevel    = "warn"
	defaultSlackAPILogLevel = "warn"
	defaultLogChannelSize   = 2000
)

var logsChannel chan [][]byte
var logsWaitGroup = new(sync.WaitGroup)
var logsDone chan bool

type LogConfig struct {
	EnvPrefix   string
	Formatter   string // "text" | "json"
	LogLevel    string //  "debug | info | warn | error"
	LogChanSize int
	TeamsLogCfg TeamsLogCfg
	SlackAPICfg SlackAPICfg
	LogChanCfg  LogChanCfg
}

type TeamsLogCfg struct {
	Enabled  bool
	Endpoint string // Endpoint for your Teams hook
	LogLevel string //  "debug | info | warn | error | fatal"
}

type SlackAPICfg struct {
	Enabled   bool
	Token     string // Slack Bot User OAuth Token (xoxb-...)
	Channel   string // Channel ID (e.g., C086K...)
	LogLevel  string // "debug | info | warn | error | fatal"
	UseBlocks bool   // Whether to use rich block formatting
}

// LogChanCfg configures the LogChan hook which sends logrus-text-formatted
// messages to a caller-provided string channel.
type LogChanCfg struct {
	Enabled  bool
	Ch       chan string // caller-provided channel to receive log messages
	LogLevel string     // "debug | info | warn | error | fatal"
}
