package settings

import "os"

type BotSetting string

const (
	VERSION BotSetting = "VERSION"
	TOKEN   BotSetting = "TOKEN"
	PREFIX  BotSetting = "PREFIX"
)

func AllSettings() []BotSetting {
	return []BotSetting{VERSION, TOKEN, PREFIX}
}

func GatherSettings() map[BotSetting]string {
	s := map[BotSetting]string{}
	for _, v := range AllSettings() {
		s[v] = os.Getenv(string(v))
	}
	return s
}
