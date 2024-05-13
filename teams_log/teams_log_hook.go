package teams_log

import (
	"fmt"

	"github.com/rohanthewiz/serr"
	"github.com/sirupsen/logrus"
)

var logIcons = map[logrus.Level]string{
	logrus.DebugLevel: "https://d2kk8pyj1kjlmo.cloudfront.net/icons/notepad_32.png",
	logrus.InfoLevel:  "https://d2kk8pyj1kjlmo.cloudfront.net/icons/note_32.png",
	logrus.WarnLevel:  "https://d2kk8pyj1kjlmo.cloudfront.net/icons/flash_32.png",
	logrus.ErrorLevel: "https://d2kk8pyj1kjlmo.cloudfront.net/icons/error_32.png",
	logrus.FatalLevel: "https://d2kk8pyj1kjlmo.cloudfront.net/icons/dead_scrn_32.png",
	logrus.PanicLevel: "https://d2kk8pyj1kjlmo.cloudfront.net/icons/dead_scrn_32.png",
}

type TeamsLogHook struct {
	AcceptedLevels []logrus.Level
	URL            string
	Disabled       bool
}

var allLevels = []logrus.Level{
	logrus.DebugLevel,
	logrus.InfoLevel,
	logrus.WarnLevel,
	logrus.ErrorLevel,
	logrus.FatalLevel,
	logrus.PanicLevel,
}

// Levels sets which levels to send to Teams
// This method is required for logrus hooks
func (th *TeamsLogHook) Levels() []logrus.Level {
	if th.AcceptedLevels == nil {
		return allLevels
	}
	return th.AcceptedLevels
}

// levelThreshold - Returns every logging level above and including the given parameter.
func AllowedLevels(lvl logrus.Level) []logrus.Level {
	for i := range allLevels {
		if allLevels[i] == lvl {
			return allLevels[i:]
		}
	}
	return []logrus.Level{}
}

func (th TeamsLogHook) Fire(le *logrus.Entry) (err error) {
	if th.Disabled {
		return nil
	}

	mc := MessageCard{
		Type:    messageCardType,
		Context: messageCardContext,
		Summary: "Log",
	}

	sec := Section{
		ActivityTitle: le.Message,
		// ActivitySubtitle: // le.Time.Format("2006-01-02 15:04 MST"),
		ActivityImage: logIcons[le.Level],
	}

	for k, v := range le.Data {
		val, ok := v.(string)
		if !ok {
			continue
		}

		switch k {
		case "msg":
			sec.ActivityTitle = val
		case "error":
			sec.ActivityText = "`" + val + "`" // quiet markdown formatting

		default:
			sec.Facts = append(sec.Facts, Fact{Name: k, Value: "`" + val + "`"})
		}
	}

	mc.Sections = []Section{sec}

	err = SendLog(mc, th.URL)
	if err != nil {
		ser, ok := err.(serr.SErr)
		if ok {
			fmt.Println(ser.String())
		} else {
			fmt.Println(err)
		}
	}

	return
}
