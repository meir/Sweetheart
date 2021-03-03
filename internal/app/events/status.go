package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Ready(session *discordgo.Session, ready *discordgo.Ready) {
	updateStatus(session)
}

func JoinedGuild(session *discordgo.Session, guild *discordgo.GuildCreate) {
	updateStatus(session)
}

func LeftGuild(session *discordgo.Session, guild *discordgo.GuildCreate) {
	updateStatus(session)
}

func updateStatus(session *discordgo.Session) {
	activity := &discordgo.Activity{
		Name: fmt.Sprintf("%v servers", len(session.State.Guilds)),
		Type: discordgo.ActivityTypeListening,
	}

	session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{activity},
		Status:     string(discordgo.StatusOnline),
	})
}
