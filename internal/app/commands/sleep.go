package commands

import (
	"bytes"
	"image/png"

	"github.com/meir/Sweetheart/internal/pkg/commandeer"
	"github.com/meir/Sweetheart/internal/pkg/logging"
)

func sleep(meta commandeer.Meta, command string, arguments []string) bool {
	image := meta.DialogueGenerator.GenerateDialogue("Sleep is the most important meal of the day!", meta.DialogueGenerator.NormalFont, 500, 140)
	var buf bytes.Buffer
	err := png.Encode(&buf, image)
	if err != nil {
		logging.Warn("Failed to encode buffer to png", err)
		return false
	}
	_, err = meta.Session.ChannelFileSend(meta.Message.ChannelID, "file.png", bytes.NewReader(buf.Bytes()))
	if err != nil {
		logging.Warn("Failed to send message", err)
		return false
	}
	return true
}
