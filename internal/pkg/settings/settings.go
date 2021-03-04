package settings

import "os"

type BotSetting string

const (
	VERSION BotSetting = "VERSION"
	TOKEN   BotSetting = "TOKEN"
	PREFIX  BotSetting = "PREFIX"
	ASSETS  BotSetting = "ASSETS"

	FEEDBACK_WEBHOOK BotSetting = "FEEDBACK_WEBHOOK"
	DEBUG_WEBHOOK    BotSetting = "DEBUG_WEBHOOK"
	DEBUG            BotSetting = "DEBUG"
)

func AllSettings() []BotSetting {
	return []BotSetting{VERSION, TOKEN, PREFIX, ASSETS, FEEDBACK_WEBHOOK, DEBUG_WEBHOOK, DEBUG}
}

func GatherSettings() map[BotSetting]string {
	s := map[BotSetting]string{}
	for _, v := range AllSettings() {
		s[v] = os.Getenv(string(v))
	}
	return s
}
