package commands

import (
	"strings"

	"github.com/meir/Sweetheart/internal/pkg/commandeer"
	"github.com/meir/Sweetheart/internal/pkg/webhook"
)

func feedback(meta commandeer.Meta, command string, arguments []string) bool {
	webhook.Feedback(meta.Meta, *meta.Message, strings.Join(arguments, " "))
	meta.Session.ChannelMessageSend(meta.Message.ChannelID, "Sent feedback to developer server! thank you :)")
	return true
}
