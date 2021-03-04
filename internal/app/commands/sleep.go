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
	go png.Encode(writer, image)
	println("D")
	defer writer.Close()
	defer reader.Close()
	println("E")
	meta.Session.ChannelFileSend(meta.Message.ChannelID, "file.png", reader)
	println("F")
	return true
}
