package commands

import (
	"bytes"
	"fmt"
	"image/png"

	"github.com/fogleman/gg"
	"github.com/meir/Sweetheart/internal/pkg/commandeer"
	"github.com/meir/Sweetheart/internal/pkg/logging"
	"github.com/meir/Sweetheart/internal/pkg/settings"
)

func status(meta commandeer.Meta, command string, arguments []string) bool {
	lines := map[string]bool{}
	lines[fmt.Sprintf("Sweetheart VH-%v", meta.Settings[settings.VERSION])] = true
	width := 500
	height := 20 * len(lines)
	dc := gg.NewContext(width, height)

	dc.SetRGB(0, 0, 0)
	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.Fill()
	dc.SetRGB(1, 1, 1)
	dc.DrawRectangle(2, 2, float64(width)-4, float64(height)-4)
	dc.Fill()
	dc.SetRGB(0, 0, 0)
	dc.DrawRectangle(6, 6, float64(width)-12, float64(height)-12)
	dc.Fill()

	dc.SetRGB(1, 1, 1)
	dc.SetFontFace(meta.DialogueGenerator.NormalFont)
	var i float64 = 0
	for k, v := range lines {
		dc.DrawStringWrapped(k, 20, 5+(20*i), 0, 0, float64(width)-20, 1, gg.AlignLeft)
		dc.Fill()
		if v {
			dc.SetRGB255(110, 255, 161)
			dc.DrawCircle(10, 10+(20*i), 5)
			dc.Fill()
		} else {
			dc.SetRGB255(255, 148, 138)
			dc.DrawCircle(10, 10+(20*i), 5)
			dc.SetLineWidth(2)
			dc.Stroke()
		}
		i++
	}

	var buf bytes.Buffer
	err := png.Encode(&buf, dc.Image())
	if err != nil {
		logging.Warn("Failed to encode buffer into png", err)
	}

	_, err = meta.Session.ChannelFileSend(meta.Message.ChannelID, "sweetheart-status.png", bytes.NewReader(buf.Bytes()))
	if err != nil {
		logging.Warn("Failed to send message", err)
		return false
	}
	return true
}
