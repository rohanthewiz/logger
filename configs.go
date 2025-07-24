package logger

import (
	"sync"
)

const (
	defaultLogLevel         = "debug" //  "debug | info | warn | error"
	defaultSlackrusLogLevel = "warn"
	defaultTeamsLogLevel    = "warn"
	defaultSlackAPILogLevel = "warn"
	defaultLogChannelSize   = 2000
)

var logsChannel chan [][]byte
var logsWaitGroup = new(sync.WaitGroup)
var logsDone chan bool

type SlackrusCfg struct {
	Enabled  bool
	Endpoint string // Endpoint for your Slack hook
	LogLevel string //  "debug | info | warn | error"
}

type LogConfig struct {
	EnvPrefix   string
	Formatter   string // "text" | "json"
	LogLevel    string //  "debug | info | warn | error"
	LogChanSize int
	SlackrusCfg SlackrusCfg
	TeamsLogCfg TeamsLogCfg
	SlackAPICfg SlackAPICfg
}

type TeamsLogCfg struct {
	Enabled  bool
	Endpoint string // Endpoint for your Teams hook
	LogLevel string //  "debug | info | warn | error | fatal"
}

type SlackAPICfg struct {
	Enabled   bool
	Token     string // Slack Bot User OAuth Token (xoxb-...)
	Channel   string // Channel ID (e.g., C086KQQUGCW)
	LogLevel  string // "debug | info | warn | error | fatal"
	UseBlocks bool   // Whether to use rich block formatting
}
