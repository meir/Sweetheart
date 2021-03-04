package commands

import (
	"image/png"
	"io"

	"github.com/meir/Sweetheart/internal/pkg/commandeer"
)

func sleep(meta commandeer.Meta, command string, arguments []string) bool {
	image := meta.DialogueGenerator.GenerateDialogue("Sleep is the most important meal of the day!", meta.DialogueGenerator.NormalFont, 500, 200)
	reader, writer := io.Pipe()
	err := png.Encode(writer, image)
	if err != nil {
		panic(err)
	}
	meta.Session.ChannelFileSend(meta.Message.ChannelID, "file.png", reader)
	return true
}
