package commands

import (
	"image/png"
	"io"

	"github.com/meir/Sweetheart/internal/pkg/commandeer"
)

func sleep(meta commandeer.Meta, command string, arguments []string) bool {
	println("A")
	image := meta.DialogueGenerator.GenerateDialogue("Sleep is the most important meal of the day!", meta.DialogueGenerator.NormalFont, 500, 200)
	println("B")
	reader, writer := io.Pipe()
	println("C")
	err := png.Encode(writer, image)
	println("D")
	writer.Close()
	defer reader.Close()
	if err != nil {
		panic(err)
	}
	println("E")
	meta.Session.ChannelFileSend(meta.Message.ChannelID, "file.png", reader)
	println("F")
	return true
}
