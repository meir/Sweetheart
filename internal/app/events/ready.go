package events

import "github.com/bwmarrin/discordgo"

func Ready(discord *discordgo.Session, ready *discordgo.Ready) {
	activity := &discordgo.Activity{
		Name: "OHOHOHOHOHOHOHO",
		Type: discordgo.ActivityTypeCustom,
	}

	discord.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{activity},
	})
}
