package commands

import (
	"context"
	"fmt"

	"github.com/meir/Sweetheart/internal/pkg/commandeer"
	"github.com/meir/Sweetheart/internal/pkg/data"
	"github.com/meir/Sweetheart/internal/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
)

func rank(meta commandeer.Meta, command string, arguments []string) bool {
	collection, err := meta.GetCollection("users")
	if err != nil {
		logging.Warn(err)
		return false
	}

	res := collection.FindOne(context.Background(), bson.M{
		"id": meta.Message.Author.ID,
	})

	if res.Err() != nil {
		logging.Warn(res.Err())
		return false
	}

	var user data.DiscordDetails
	err = res.Decode(&user)
	if err != nil {
		logging.Warn(err)
		return false
	}

	_, err = meta.Session.ChannelMessageSend(meta.Message.ChannelID, fmt.Sprintf("EXP: %v", user.Ranks[meta.Message.GuildID]))
	if err != nil {
		logging.Warn("Failed to send message", err)
		return false
	}
	return true
}
