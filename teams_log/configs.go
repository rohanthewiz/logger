package teams_log

var teamsConfig TeamsCfg

type TeamsCfg struct {
	Enabled     bool
	LogEndpoint string
	LogLevel    string
}

func GetTeamsCfg() TeamsCfg {
	return teamsConfig
}

func SetTeamsCfg(cfg TeamsCfg) {
	teamsConfig = cfg
}
