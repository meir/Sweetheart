package sweetheart

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/Sweetheart/internal/pkg/bot"
	"github.com/meir/Sweetheart/internal/pkg/commandeer"
	"github.com/meir/Sweetheart/internal/pkg/settings"
)

func Sweetheart() {
	sweetheart, err := bot.NewBot(settings.GatherSettings())
	if err != nil {
		panic(err)
	}

	sweetheart.Identify.Intents = discordgo.IntentsGuildMessages
	sweetheart.Commandeer.Start(sweetheart.Session)

	sweetheart.Commandeer.Apply("version", func(meta commandeer.Meta, command string, arguments []string) bool {
		meta.Session.ChannelMessageSend(meta.Message.ChannelID, fmt.Sprintf("Running Sweetheart VH-%v", meta.Settings[settings.VERSION]))
		return true
	}, commandeer.Arguments{Any: true})

	err = sweetheart.Open()
	if err != nil {
		panic(err)
	}

	println("Sweetheart is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	sweetheart.Close()
}
