package commands

import (
	"fmt"
	"strings"

	"github.com/meir/Sweetheart/internal/pkg/commandeer"
	"github.com/meir/Sweetheart/internal/pkg/logging"
)

func commands(meta commandeer.Meta, command string, arguments []string) bool {

	long := 0
	for k := range meta.Commands {
		if long < len(k) {
			long = len(k)
		}
	}

	output := []string{}
	for k, v := range meta.Commands {
		spaces := len(k) - long
		output = append(output, fmt.Sprintf("%v%v : %v", k, strings.Repeat(" ", spaces), v))
	}

	_, err := meta.Session.ChannelMessageSend(meta.Message.ChannelID, fmt.Sprintf("```\n-- Commands --\n\n%v\n```", strings.Join(output, "\n")))
	if err != nil {
		logging.Warn("Failed to send message", err)
		return false
	}
	return true
}
