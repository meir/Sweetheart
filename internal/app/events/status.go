package events

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/Sweetheart/internal/pkg/bot"
	"github.com/meir/Sweetheart/internal/pkg/settings"
)

func Ready(sweetheart *bot.DiscordBot) func(session *discordgo.Session, guild *discordgo.Ready) {
	return func(session *discordgo.Session, guild *discordgo.Ready) {
		updateStatus(session)
	}
}

func JoinedGuild(sweetheart *bot.DiscordBot) func(session *discordgo.Session, guild *discordgo.GuildCreate) {
	return func(session *discordgo.Session, guild *discordgo.GuildCreate) {
		updateStatus(session)
		roles, err := session.GuildRoles(guild.ID)
		if err != nil {
			panic(err)
		}
		println(len(roles))
		for _, role := range roles {
			println(role.Name)
			if role.Name == session.State.User.Username {
				color, err := strconv.Atoi(sweetheart.Meta.Settings[settings.ROLE_COLOR])
				if err != nil {
					panic(err)
				}
				_, err = session.GuildRoleEdit(guild.ID, role.ID, role.Name, color, role.Hoist, role.Permissions, role.Mentionable)
				if err != nil {
					panic(err)
				}
			}
		}
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

	session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{activity},
		Status:     string(discordgo.StatusOnline),
	})
}
