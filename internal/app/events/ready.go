package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Ready(discord *discordgo.Session, ready *discordgo.Ready) {
	activity := &discordgo.Activity{
		Name: fmt.Sprintf("%v servers", len(discord.State.Guilds)),
		Type: discordgo.ActivityTypeListening,
	}

	discord.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{activity},
		Status:     string(discordgo.StatusOnline),
	})
}
