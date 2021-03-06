package sweetheart

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/Sweetheart/internal/app/commands"
	"github.com/meir/Sweetheart/internal/app/events"
	"github.com/meir/Sweetheart/internal/pkg/bot"
	"github.com/meir/Sweetheart/internal/pkg/settings"
)

func Sweetheart() {
	sweetheart, err := bot.NewBot(settings.GatherSettings())
	if err != nil {
		panic(err)
	}

	sweetheart.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds
	sweetheart.Commandeer.Start(sweetheart.Session) // initialize commandeer

	// initialize commands and event listeners
	commands.Initialize(sweetheart)
	events.Initialize(sweetheart)

	err = sweetheart.Open()
	if err != nil {
		panic(err)
	}

	sweetheart.Meta.Status[fmt.Sprintf("Sweetheart VH-%v", sweetheart.Meta.Settings[settings.VERSION])] = true

	// wait for kill signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	sweetheart.Close()
}
