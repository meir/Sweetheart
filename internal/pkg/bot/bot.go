package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/Sweetheart/internal/pkg/commandeer"
	"github.com/meir/Sweetheart/internal/pkg/dialogue"
	"github.com/meir/Sweetheart/internal/pkg/meta"
	"github.com/meir/Sweetheart/internal/pkg/settings"
)

type DiscordBot struct {
	*discordgo.Session

	Commandeer *commandeer.Commandeer
	Meta       *meta.Meta
}

func NewBot(st map[settings.BotSetting]string) (*DiscordBot, error) {
	if _, ok := st[settings.TOKEN]; !ok {
		return nil, fmt.Errorf("requires bot token in settings")
	}

	c, err := discordgo.New(fmt.Sprintf("Bot %v", strings.ReplaceAll(st[settings.TOKEN], "\n", "")))
	if err != nil {
		return nil, err
	}

	meta := &meta.Meta{
		c,
		st,
		dialogue.NewDialogueGenerator(st),
	}

	return &DiscordBot{
		c,
		commandeer.NewCommandeer(st[settings.PREFIX], meta),
		meta,
	}, nil
}
