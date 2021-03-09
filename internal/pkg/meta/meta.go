package meta

import (
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
