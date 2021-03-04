package commands

import (
	"github.com/meir/Sweetheart/internal/pkg/commandeer"
	"github.com/meir/Sweetheart/internal/pkg/webhook"
)

func feedback(meta commandeer.Meta, command string, arguments []string) bool {
	webhook.Feedback(meta.Meta, *meta.Message)
	meta.Session.ChannelMessageSend(meta.Message.ChannelID, "Sent feedback to developer server! thank you :)")
	return true
}
