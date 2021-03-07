package commands

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/Sweetheart/internal/pkg/commandeer"
	"github.com/meir/Sweetheart/internal/pkg/logging"
	"github.com/meir/Sweetheart/internal/pkg/settings"
)

func feedback(meta commandeer.Meta, command string, arguments []string) bool {
	if meta.Settings[settings.FEEDBACK_WEBHOOK] == "" {
		logging.Warn("No feedback webhook has been set; could not perform feedback command.")
		return false
	}

	selfIcon := meta.Session.State.User.AvatarURL("256")
	self := meta.Session.State.User.Username
	userIcon := meta.Message.Author.AvatarURL("256")
	user := meta.Message.Author.Username

	guild, err := meta.Session.Guild(meta.Message.GuildID)
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
		Description: strings.Join(arguments, " "),
		Author: &discordgo.MessageEmbedAuthor{
			Name:    user,
			IconURL: userIcon,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    guildName,
			IconURL: guildIcon,
		},
		Timestamp: string(meta.Message.Timestamp),
	}

	webhook := discordgo.WebhookParams{
		Username:  self,
		AvatarURL: selfIcon,
		Embeds:    []*discordgo.MessageEmbed{embed},
	}

	data, err := json.Marshal(webhook)
	if err != nil {
		logging.Warn("Failed to marshal webhook params to json", err)
		return false
	}

	r, err := http.Post(meta.Settings[settings.FEEDBACK_WEBHOOK], "application/json", bytes.NewReader(data))
	if err != nil {
		logging.Warn("Failed to call feedback webhook", err)
		return false
	}
	defer r.Body.Close()
	_, err = meta.Session.ChannelMessageSend(meta.Message.ChannelID, "Sent feedback to developer server! thank you :)")
	if err != nil {
		logging.Warn("Failed to send message", err)
		return false
	}
	return true
}
