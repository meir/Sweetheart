package commands

import (
	"github.com/meir/Sweetheart/internal/pkg/commandeer"
	"github.com/meir/Sweetheart/internal/pkg/logging"
)

func commands(meta commandeer.Meta, command string, arguments []string) bool {
	_, err := meta.Session.ChannelMessageSend(meta.Message.ChannelID, "See all the commands on https://sweetheart.flamingo.dev/commands")
	if err != nil {
		logging.Warn("Failed to send message", err)
		return false
	}
	return true
}
