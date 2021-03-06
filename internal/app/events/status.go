package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/Sweetheart/internal/pkg/bot"
	"github.com/meir/Sweetheart/internal/pkg/logging"
)

func Ready(sweetheart *bot.DiscordBot) func(session *discordgo.Session, guild *discordgo.Ready) {
	return func(session *discordgo.Session, guild *discordgo.Ready) {
		updateStatus(session)
	}
}

func JoinedGuild(sweetheart *bot.DiscordBot) func(session *discordgo.Session, guild *discordgo.GuildCreate) {
	return func(session *discordgo.Session, guild *discordgo.GuildCreate) {
		updateStatus(session)
	}
}

func LeftGuild(sweetheart *bot.DiscordBot) func(session *discordgo.Session, guild *discordgo.GuildDelete) {
	return func(session *discordgo.Session, guild *discordgo.GuildDelete) {
		updateStatus(session)
	}
}

func updateStatus(session *discordgo.Session) {
	activity := &discordgo.Activity{
		Name: fmt.Sprintf("%v servers", len(session.State.Guilds)),
		Type: discordgo.ActivityTypeListening,
	}

	err := session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{activity},
		Status:     string(discordgo.StatusOnline),
	})
	if err != nil {
		logging.Warn("Failed to update bot status to amount of servers", err)
		return
	}
}
