package commands

import (
	"context"
	"fmt"
	"sort"

	"github.com/meir/Sweetheart/internal/pkg/commandeer"
	"github.com/meir/Sweetheart/internal/pkg/data"
	"github.com/meir/Sweetheart/internal/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func rank(meta commandeer.Meta, command string, arguments []string) bool {
	collection, err := meta.GetCollection("users")
	if err != nil {
		logging.Warn(err)
		return false
	}

	res, err := collection.Find(context.Background(), bson.M{
		"ranks": bson.M{
			meta.Message.GuildID: bson.M{
				"$exists": true,
			},
		},
	}, options.Find().SetProjection(bson.D{
		{fmt.Sprintf("ranks.%v", meta.Message.GuildID), 1},
		{"id", 1},
		{"username", 1},
	}))

	if res.Err() != nil || err != nil {
		logging.Warn(res.Err(), err)
		return false
	}

	var users []data.DiscordDetails
	err = res.All(context.Background(), &users)
	if err != nil {
		logging.Warn(err)
		return false
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].Ranks[meta.Message.GuildID] > users[i].Ranks[meta.Message.GuildID]
	})

	rank := -1
	var exp uint64 = 0
	for k, v := range users {
		if v.ID == meta.Message.Author.ID {
			rank = k + 1
			exp = v.Ranks[meta.Message.GuildID]
			break
		}
	}

	_, err = meta.Session.ChannelMessageSend(meta.Message.ChannelID, fmt.Sprintf("Rank #%v EXP: %v", rank, exp))
	if err != nil {
		logging.Warn("Failed to send message", err)
		return false
	}
	return true
}
