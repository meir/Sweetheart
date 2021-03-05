package commands

import (
	"fmt"

	"github.com/meir/Sweetheart/internal/pkg/bot"
	"github.com/meir/Sweetheart/internal/pkg/commandeer"
	"github.com/meir/Sweetheart/internal/pkg/settings"
)

func Initialize(sweetheart *bot.DiscordBot) {
	sweetheart.Commandeer.Apply("version", version, commandeer.Arguments{Any: true}, "")
	sweetheart.Commandeer.Apply("sleep", sleep, commandeer.Arguments{Any: true}, "")
	sweetheart.Commandeer.Apply("feedback", feedback, commandeer.Arguments{Min: 1}, "[message]")
	sweetheart.Commandeer.Apply("wiserock", wiserock, commandeer.Arguments{Min: 0, Max: 1}, "(good/ok/bad)")

	sweetheart.Commandeer.FailedArguments = failedArguments
}

func failedArguments(meta commandeer.Meta, command string, arguments []string) bool {
	meta.Session.ChannelMessageSend(meta.Message.ChannelID, fmt.Sprintf("Cant really use the arguments, you can use the command like this: `%v%v %v`", meta.Settings[settings.PREFIX], command, meta.Usage))
	return true
}
