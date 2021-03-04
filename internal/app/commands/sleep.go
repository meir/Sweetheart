package commands

import (
	"bytes"
	"image/png"

	"github.com/meir/Sweetheart/internal/pkg/commandeer"
)

func sleep(meta commandeer.Meta, command string, arguments []string) bool {
	image := meta.DialogueGenerator.GenerateDialogue("Sleep is the most important meal of the day!", meta.DialogueGenerator.NormalFont, 500, 140)
	var buf bytes.Buffer
	err := png.Encode(&buf, image)
	if err != nil {
		panic(err)
	}
	meta.Session.ChannelFileSend(meta.Message.ChannelID, "file.png", bytes.NewReader(buf.Bytes()))
	return true
}
