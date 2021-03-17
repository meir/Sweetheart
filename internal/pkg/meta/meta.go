package meta

import (
	"context"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/Sweetheart/internal/pkg/dialogue"
	"github.com/meir/Sweetheart/internal/pkg/settings"
	"go.mongodb.org/mongo-driver/mongo"
)

type Meta struct {
	Session           *discordgo.Session
	Settings          map[settings.BotSetting]string
	DialogueGenerator *dialogue.DialogueGenerator
	Database          *mongo.Client
	Status            map[string]bool
	Commands          map[string]string
}

func (m *Meta) IsAdmin(id string) bool {
	admins := strings.Split(m.Settings[settings.ADMINS], ",")
	for _, v := range admins {
		if id == v {
			return true
		}
	}
	return false
}

func (meta *Meta) GetCollection(name string) (*mongo.Collection, error) {
	database := meta.Database.Database("sweetheart")
	col := database.Collection(name)
	if col == nil {
		err := database.CreateCollection(context.Background(), name)
		if err != nil {
			return nil, err
		}
		col = database.Collection(name)
	}
	return col, nil
}
