package webhook

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/Sweetheart/internal/pkg/meta"
	"github.com/meir/Sweetheart/internal/pkg/settings"
)

func Feedback(meta *meta.Meta, message discordgo.Message, feedback string) {
	if meta.Settings[settings.FEEDBACK_WEBHOOK] == "" {
		// TODO: logging
		return
	}

	selfIcon := meta.Session.State.User.AvatarURL("256")
	self := meta.Session.State.User.Username
	userIcon := message.Author.AvatarURL("256")
	user := message.Author.Username

	guild, err := meta.Session.Guild(message.GuildID)
	if err != nil {
		guild = nil
	}

	guildIcon := ""
	guildName := "private"

	if guild != nil {
		guildIcon = guild.IconURL()
		guildName = guild.Name
	}

	embed := &discordgo.MessageEmbed{
		Description: feedback,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    user,
			IconURL: userIcon,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    guildName,
			IconURL: guildIcon,
		},
		Timestamp: string(message.Timestamp),
	}

	webhook := discordgo.WebhookParams{
		Username:  self,
		AvatarURL: selfIcon,
		Embeds:    []*discordgo.MessageEmbed{embed},
	}

	data, err := json.Marshal(webhook)
	if err != nil {
		// TODO: logging
		return
	}

	_, err = http.Post(meta.Settings[settings.FEEDBACK_WEBHOOK], "application/json", bytes.NewReader(data))
	if err != nil {
		// TODO: logging
		return
	}
}
