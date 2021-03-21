package events

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/Sweetheart/internal/pkg/bot"
	"github.com/meir/Sweetheart/internal/pkg/data"
	"github.com/meir/Sweetheart/internal/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Message(sweetheart *bot.DiscordBot) func(session *discordgo.Session, guild *discordgo.MessageCreate) {
	return func(session *discordgo.Session, msg *discordgo.MessageCreate) {
		upsert := true
		collection, err := sweetheart.Meta.GetCollection("users")
		if err != nil {
			logging.Warn(err)
			return
		}

		_, err = collection.UpdateOne(context.Background(), bson.M{
			"id": msg.Author.ID,
		}, bson.M{
			"$set": bson.M{
				"username":      msg.Author.Username,
				"avatar":        msg.Author.Avatar,
				"discriminator": msg.Author.Discriminator,
			},
			"$setOnInsert": bson.M{
				"id":      msg.Author.ID,
				"profile": data.User{},
			},
			"$inc": bson.M{
				"ranks.global":                       1,
				fmt.Sprintf("ranks.%v", msg.GuildID): 1,
			},
		}, &options.UpdateOptions{
			Upsert: &upsert,
		})
		if err != nil {
			logging.Warn(err)
		}
	}
}
