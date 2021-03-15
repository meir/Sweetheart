package commands

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/Sweetheart/internal/pkg/commandeer"
	"github.com/meir/Sweetheart/internal/pkg/logging"
	"github.com/meir/Sweetheart/internal/pkg/settings"
)

func boop(meta commandeer.Meta, command string, arguments []string) bool {
	if meta.IsAdmin(meta.Message.Author.ID) {
		return true
	}
	images := []string{
		"perfectheart.png",
		"sweetheart.png",
		"mutantheart.png",
	}

	i, err := strconv.Atoi(arguments[0])
	if err != nil {
		return false
	}

	if !(i >= 0 && i < len(images)) {
		return false
	}

	img, err := ioutil.ReadFile(path.Join(meta.Settings[settings.ASSETS], "images", images[i]))
	if err != nil {
		logging.Warn("Failed to read file", err)
		return false
	}
	t := http.DetectContentType(img)
	im := base64.StdEncoding.EncodeToString(img)

	user, err := meta.Session.UserUpdate("", "", "", fmt.Sprintf("data:%s;base64,%s", t, im), "")
	if err != nil || user == nil {
		logging.Warn("Failed to update profile picture", err)
		return false
	}

	rolecolor, err := strconv.Atoi(meta.Settings[settings.ROLE_COLOR])
	if err != nil {
		logging.Warn("Failed to convert role color to int", err)
		return false
	}

	_, err = meta.Session.ChannelMessageSendEmbed(meta.Message.ChannelID, &discordgo.MessageEmbed{
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: user.AvatarURL("256"),
		},
		Title:       "Boop!",
		Description: fmt.Sprintf("Updated to %v", images[i]),
		Color:       rolecolor,
		Type:        discordgo.EmbedTypeRich,
	})
	if err != nil {
		logging.Warn("Failed to send message", err)
		return false
	}
	return true
}
