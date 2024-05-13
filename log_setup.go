package logger

import (
	"strings"

	"github.com/johntdyer/slackrus"
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

	// Important - must set formatter -- should not be nil
	if logCfg.Formatter == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

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

	// Slack Log
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

	if ll, ok := logrusLevels[strings.ToLower(logCfg.LogLevel)]; ok {
		logrus.SetLevel(ll)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}
