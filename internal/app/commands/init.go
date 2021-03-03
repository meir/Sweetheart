package commands

import (
	"github.com/meir/Sweetheart/internal/pkg/bot"
	"github.com/meir/Sweetheart/internal/pkg/commandeer"
)

func Initialize(sweetheart *bot.DiscordBot) {
	sweetheart.Commandeer.Apply("version", version, commandeer.Arguments{Any: true})
}
