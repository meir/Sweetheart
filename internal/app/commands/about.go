package commands

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/Sweetheart/internal/pkg/commandeer"
	"github.com/meir/Sweetheart/internal/pkg/logging"
	"github.com/meir/Sweetheart/internal/pkg/webserver"
	"go.mongodb.org/mongo-driver/bson"
)

var reg = regexp.MustCompile(`(?m)^<@([0-9]{1,})>$`)

func about(meta commandeer.Meta, command string, arguments []string) bool {
	if len(arguments) == 0 {
		_, err := meta.Session.ChannelMessageSend(meta.Message.ChannelID, fmt.Sprintf("The about command is used to seek profiles of other people\nYou can make your own profile on https://sweetheart.flamingo.dev/"))
		if err != nil {
			logging.Warn("Failed to send message", err)
			return false
		}
		return true
	}

	id := arguments[0]

	if !reg.MatchString(id) {
		return false
	}
	id = strings.ReplaceAll(id, "<", "")
	id = strings.ReplaceAll(id, "@", "")
	id = strings.ReplaceAll(id, ">", "")

	collection := meta.Database.Database("sweetheart").Collection("users")
	res := collection.FindOne(context.Background(), bson.M{
		"id": id,
	})
	if res == nil {
		_, err := meta.Session.ChannelMessageSend(meta.Message.ChannelID, fmt.Sprintf("Could not find the profile of this user\nYou can make your own profile on https://sweetheart.flamingo.dev/"))
		if err != nil {
			logging.Warn("Failed to send message", err)
			return false
		}
		return true
	}

	var details webserver.DiscordDetails
	err := res.Decode(&details)
	if err != nil {
		logging.Warn("could not parse mongodb query to DiscordDetails.", err)
		return false
	}

	image := fmt.Sprintf("https://cdn.discordapp.com/avatars/%v/%v", details.ID, details.Avatar)
	t := time.Now().UTC()
	t.Add(time.Duration(int64(time.Minute) * int64(details.Profile.Timezone)))

	socials := ""
	for _, v := range details.Profile.Socials {
		socials += fmt.Sprintf("> __%v:__ %v\n", v.Name, v.Handle)
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("%v#%v", details.Username, details.Discriminator),
			IconURL: image,
		},
		Color: details.Profile.FavoriteColor,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "About me",
				Value: details.Profile.Description,
			},
			{
				Name:  "Socials",
				Value: socials,
			},
			{
				Name:  "Sexual Orientation",
				Value: fmt.Sprintf("__Sexuality:__ %v\n__Gender:__ %v\n__Pronouns:__ %v", details.Profile.Sexuality, details.Profile.Gender, details.Profile.Pronouns),
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%v %v • %v", details.Profile.Country, webserver.Countries[details.Profile.Country], t.Format("15:04")),
		},
	}

	_, err = meta.Session.ChannelMessageSendEmbed(meta.Message.ChannelID, embed)
	if err != nil {
		logging.Warn("Failed to send message", err)
		return false
	}
	return true
}
