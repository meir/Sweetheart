package meta

import (
	"github.com/bwmarrin/discordgo"
	"github.com/meir/Sweetheart/internal/pkg/dialogue"
	"github.com/meir/Sweetheart/internal/pkg/settings"
)

type Meta struct {
	Session           *discordgo.Session
	Settings          map[settings.BotSetting]string
	DialogueGenerator *dialogue.DialogueGenerator
}
