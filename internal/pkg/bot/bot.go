package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/Sweetheart/internal/pkg/commandeer"
)

type DiscordBot struct {
	*discordgo.Session

	Commandeer *commandeer.Commandeer
	Settings   map[BotSetting]string
}

func NewBot(settings map[BotSetting]string) (*DiscordBot, error) {
	if _, ok := settings[TOKEN]; !ok {
		return nil, fmt.Errorf("requires bot token in settings")
	}

	c, err := discordgo.New(fmt.Sprintf("Bot %v", settings[TOKEN]))
	if err != nil {
		return nil, err
	}

	return &DiscordBot{
		c,
		commandeer.NewCommandeer(settings[PREFIX]),
		settings,
	}, nil
}
