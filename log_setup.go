package logger

import (
	"fmt"
	"strings"

	"github.com/johntdyer/slackrus"
	"github.com/rohanthewiz/logger/slack_api"
	"github.com/rohanthewiz/logger/teams_log"
	"github.com/sirupsen/logrus"
)

var logPrefix string

func InitLog(logCfg LogConfig) {
	initLogrus(logCfg)

	logChanSize := logCfg.LogChanSize
	if logChanSize == 0 {
		logChanSize = defaultLogChannelSize
	}

	logsChannel = make(chan [][]byte, logChanSize)
	logsDone = make(chan bool)

	go pollForLogs(logsChannel, logsDone) // start the listener
}

func CloseLog() {
	logsWaitGroup.Wait()
	// Close the channel so nothing else can be added and the log poller knows to start wrapping up
	close(logsChannel)
	<-logsDone // wait for *all* log processing to complete
	logrus.Info("Logs gracefully shutdown")
}

func initLogrus(logCfg LogConfig) {
	logPrefix = logCfg.EnvPrefix

	if logCfg.LogLevel == "" {
		logCfg.LogLevel = defaultLogLevel
	}

	SetLogFormat(logCfg.Formatter)

	SetLogLevel(logCfg.LogLevel)

	// HOOKS

	// Teams Log
	if logCfg.TeamsLogCfg.Enabled {
		// Pass config down to local package
		teams_log.SetTeamsCfg(teams_log.TeamsCfg{
			Enabled:     logCfg.TeamsLogCfg.Enabled,
			LogEndpoint: logCfg.TeamsLogCfg.Endpoint,
			LogLevel:    logCfg.TeamsLogCfg.LogLevel,
		})

		if logCfg.TeamsLogCfg.LogLevel == "" {
			logCfg.TeamsLogCfg.LogLevel = defaultTeamsLogLevel
		}

		logrus.AddHook(&teams_log.TeamsLogHook{
			URL:            logCfg.TeamsLogCfg.Endpoint,
			AcceptedLevels: teams_log.AllowedLevels(logrusLevels[strings.ToLower(logCfg.TeamsLogCfg.LogLevel)]),
		})
	}

	// Slack Log via WebHook - deprecated, prefer to use Slack API instead
	if logCfg.SlackrusCfg.Enabled {
		if logCfg.SlackrusCfg.LogLevel == "" {
			logCfg.SlackrusCfg.LogLevel = defaultSlackrusLogLevel
		}

		logrus.AddHook(&slackrus.SlackrusHook{
			HookURL:        logCfg.SlackrusCfg.Endpoint,
			AcceptedLevels: slackrus.LevelThreshold(logrusLevels[strings.ToLower(logCfg.SlackrusCfg.LogLevel)]),
			IconEmoji:      ":computer:",
		})
	}

	// Slack API Log
	if logCfg.SlackAPICfg.Enabled {
		if logCfg.SlackAPICfg.LogLevel == "" {
			logCfg.SlackAPICfg.LogLevel = defaultSlackAPILogLevel
		}

		// Convert string log level to logrus level
		acceptedLevel := logrusLevels[strings.ToLower(logCfg.SlackAPICfg.LogLevel)]
		acceptedLevels := AllowedLevels(acceptedLevel)
		fmt.Println("Slack API acceptedLevels:", acceptedLevels)

		hook := slack_api.NewSlackAPIHook(
			logCfg.SlackAPICfg.Token,
			logCfg.SlackAPICfg.Channel,
			acceptedLevels,
			logCfg.SlackAPICfg.UseBlocks,
		)
		logrus.AddHook(hook)
	}
}

// SetLogFormat sets the log format with "json" for JSON, otherwise text
func SetLogFormat(format string) {
	format = strings.ToLower(format)
	if format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{})
	}
}

// SetLogLevel sets the log level, defaulting to info
// logLevel can be "debug | info | warn | error"
func SetLogLevel(logLevel string) {
	logLevel = strings.ToLower(logLevel)

	if logLevel == "warning" {
		logLevel = "warn"
	}

	logrusLevel := logrus.InfoLevel

	if ll, ok := logrusLevels[logLevel]; ok {
		logrusLevel = ll
	}

	logrus.SetLevel(logrusLevel)
}
