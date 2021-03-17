package events

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/Sweetheart/internal/pkg/bot"
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

		ranking := bson.M{}
		ranking["global"] = bson.M{
			"$add": []interface{}{
				"global",
				1,
			},
		}
		ranking[msg.GuildID] = bson.M{
			"$add": []interface{}{
				msg.GuildID,
				1,
			},
		}

		_, err = collection.UpdateOne(context.Background(), bson.M{
			"ID": msg.Author.ID,
		}, bson.M{
			"$set": bson.M{
				"username":      msg.Author.Username,
				"avatar":        msg.Author.Avatar,
				"discriminator": msg.Author.Discriminator,
				"ranks":         ranking,
			},
			"$setOnInsert": bson.M{
				"id":      msg.Author.ID,
				"profile": nil,
			},
		}, &options.UpdateOptions{
			Upsert: &upsert,
		})
		if err != nil {
			logging.Warn(err)
		}
	}
}
