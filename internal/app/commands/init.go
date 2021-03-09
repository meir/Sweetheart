package commands

import (
	"fmt"

	"github.com/meir/Sweetheart/internal/pkg/bot"
	"github.com/meir/Sweetheart/internal/pkg/commandeer"
	"github.com/meir/Sweetheart/internal/pkg/settings"
)

func Initialize(sweetheart *bot.DiscordBot) {
	sweetheart.Commandeer.Apply("version", version, commandeer.Arguments{Any: true}, "", "Shows the current Sweetheart version running")
	sweetheart.Commandeer.Apply("sleep", sleep, commandeer.Arguments{Any: true}, "", "Sleep is the most important meal of the day!")
	sweetheart.Commandeer.Apply("feedback", feedback, commandeer.Arguments{Min: 1}, "[message]", "Give your feedback or ideas about me!")
	sweetheart.Commandeer.Apply("wiserock", wiserock, commandeer.Arguments{Min: 0, Max: 1}, "(good/ok/bad)", "A wise rock giving wise advice (wiseness of advice may vary)")
	sweetheart.Commandeer.Apply("status", status, commandeer.Arguments{Any: true}, "", "Shows the status of the bot, http-server and database.")
	sweetheart.Commandeer.Apply("about", about, commandeer.Arguments{Min: 0, Max: 1}, "(user mention)", "Shows a profile of a user, also includes their time right now :)")
	sweetheart.Commandeer.Apply("commands", commands, commandeer.Arguments{Any: true}, "", "This!")

	sweetheart.Commandeer.FailedArguments = failedArguments
}

func failedArguments(meta commandeer.Meta, command string, arguments []string) bool {
	meta.Session.ChannelMessageSend(meta.Message.ChannelID, fmt.Sprintf("Cant really use the arguments, you can use the command like this: `%v%v %v`", meta.Settings[settings.PREFIX], command, meta.Usage))
	return true
}
