package commands

import (
	"fmt"

	"github.com/meir/Sweetheart/internal/pkg/commandeer"
	"github.com/meir/Sweetheart/internal/pkg/settings"
)

func version(meta commandeer.Meta, command string, arguments []string) bool {
	meta.Session.ChannelMessageSend(meta.Message.ChannelID, fmt.Sprintf("Running Sweetheart VH-%v", meta.Settings[settings.VERSION]))
	return true
}
