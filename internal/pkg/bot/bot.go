package bot

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/Sweetheart/internal/pkg/commandeer"
	"github.com/meir/Sweetheart/internal/pkg/dialogue"
	"github.com/meir/Sweetheart/internal/pkg/logging"
	"github.com/meir/Sweetheart/internal/pkg/meta"
	"github.com/meir/Sweetheart/internal/pkg/settings"
	"github.com/meir/Sweetheart/internal/pkg/webserver"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DiscordBot struct {
	*discordgo.Session

	Commandeer *commandeer.Commandeer
	Webserver  *webserver.Webserver
	Meta       *meta.Meta
}

func NewBot(st map[settings.BotSetting]string) (*DiscordBot, error) {
	if _, ok := st[settings.TOKEN]; !ok {
		return nil, fmt.Errorf("requires bot token in settings")
	}

	c, err := discordgo.New(fmt.Sprintf("Bot %v", st[settings.TOKEN]))
	if err != nil {
		return nil, err
	}

	var database *mongo.Client = nil
	if url := st[settings.MONGODB_URL]; url != "" {
		logging.Debug("Trying to connect to database:", url)
		opts := options.Client()
		opts.ApplyURI(url)
		opts.SetMaxPoolSize(5)
		if database, err = mongo.Connect(context.Background(), opts); err != nil {
			panic(err)
		}

		err = database.Ping(context.Background(), nil)
		if err != nil {
			panic(err)
		}
	} else {
		logging.Debug("No database url has been given.")
	}

	meta := &meta.Meta{
		c,
		st,
		dialogue.NewDialogueGenerator(st),
		database,
		map[string]bool{},
		map[string]string{},
	}

	return &DiscordBot{
		c,
		commandeer.NewCommandeer(st[settings.PREFIX], meta),
		webserver.NewWebserver(meta),
		meta,
	}, nil
}
