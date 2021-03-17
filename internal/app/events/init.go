package events

import "github.com/meir/Sweetheart/internal/pkg/bot"

func Initialize(sweetheart *bot.DiscordBot) {
	sweetheart.AddHandler(Ready(sweetheart))
	sweetheart.AddHandler(JoinedGuild(sweetheart))
	sweetheart.AddHandler(LeftGuild(sweetheart))

	// sweetheart.AddHandler(Message(sweetheart))
}
