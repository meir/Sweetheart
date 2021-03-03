package events

import "github.com/meir/Sweetheart/internal/pkg/bot"

func Initialize(sweetheart *bot.DiscordBot) {
	sweetheart.AddHandler(Ready)
	sweetheart.AddHandler(JoinedGuild)
	sweetheart.AddHandler(LeftGuild)
}
