package settings

import (
	"os"
	"strings"
)

type BotSetting string

const (
	VERSION     BotSetting = "VERSION"
	TOKEN       BotSetting = "TOKEN"
	PREFIX      BotSetting = "PREFIX"
	ASSETS      BotSetting = "ASSETS"
	ROLE_COLOR  BotSetting = "ROLE_COLOR"
	MONGODB_URL BotSetting = "MONGODB_URL"
	PORT        BotSetting = "PORT"

	INVITE_URL BotSetting = "INVITE_URL"
	OAUTH_URL  BotSetting = "OAUTH_URL"

	CLIENT_ID     BotSetting = "CLIENT_ID"
	CLIENT_SECRET BotSetting = "CLIENT_SECRET"
	GRANT_TYPE    BotSetting = "GRANT_TYPE"
	SCOPE         BotSetting = "SCOPE"
	REDIRECT      BotSetting = "REDIRECT"

	ADMINS BotSetting = "ADMINS"

	FEEDBACK_WEBHOOK BotSetting = "FEEDBACK_WEBHOOK"
	DEBUG_WEBHOOK    BotSetting = "DEBUG_WEBHOOK"
	DEBUG            BotSetting = "DEBUG"
	MAINTENANCE      BotSetting = "MAINTENANCE"
)

func AllSettings() []BotSetting {
	return []BotSetting{
		VERSION, TOKEN, PREFIX, ASSETS, ROLE_COLOR, MONGODB_URL, PORT,
		INVITE_URL, OAUTH_URL,
		CLIENT_ID, CLIENT_SECRET, GRANT_TYPE, SCOPE, REDIRECT,
		ADMINS,
		FEEDBACK_WEBHOOK, DEBUG_WEBHOOK, DEBUG, MAINTENANCE,
	}
}

func GatherSettings() map[BotSetting]string {
	s := map[BotSetting]string{}
	for _, v := range AllSettings() {
		s[v] = strings.ReplaceAll(os.Getenv(string(v)), "\n", "")
		println(v, ":", s[v])
	}
	return s
}
