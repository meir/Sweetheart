package sweetheart

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/Sweetheart/internal/pkg/bot"
	"github.com/meir/Sweetheart/internal/pkg/commandeer"
)

func Sweetheart() {
	sweetheart, err := bot.NewBot(bot.GatherSettings())
	if err != nil {
		panic(err)
	}

	sweetheart.Identify.Intents = discordgo.IntentsGuildMessages
	sweetheart.Commandeer.Start(sweetheart.Session)

	sweetheart.Commandeer.Apply("version", func(session *discordgo.Session, user *discordgo.User, command string, arguments []string, raw *discordgo.Message) bool {
		session.ChannelMessageSend(raw.ChannelID, fmt.Sprintf("Running Sweetheart VH-%v", os.Getenv("VERSION")))
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
