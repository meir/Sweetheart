package commands

import (
	"fmt"
	"strings"

	"github.com/meir/Sweetheart/internal/pkg/commandeer"
	"github.com/meir/Sweetheart/internal/pkg/logging"
)

func servers(meta commandeer.Meta, command string, arguments []string) bool {
	if !meta.IsAdmin(meta.Message.Author.ID) {
		return true
	}
	servers := []string{}
	count := 0
	for _, v := range meta.Session.State.Guilds {
		line := fmt.Sprintf("%v - %v", v.Name, v.Description)
		servers = append(servers, line)
		count += len(line)
		if count >= 1750 {
			_, err := meta.Session.ChannelMessageSend(meta.Message.ChannelID, strings.Join(servers, "\n"))
			if err != nil {
				logging.Warn("Failed to send message", err)
				return false
			}
			servers = []string{}
			count = 0
		}
	}
	_, err := meta.Session.ChannelMessageSend(meta.Message.ChannelID, strings.Join(servers, "\n"))
	if err != nil {
		logging.Warn("Failed to send message", err)
		return false
	}
	return true
}
