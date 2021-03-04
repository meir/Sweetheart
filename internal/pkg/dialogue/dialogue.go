package dialogue

import (
	"image"
	"path"

	"github.com/fogleman/gg"
	"github.com/meir/Sweetheart/internal/pkg/settings"
	"golang.org/x/image/font"
)

type DialogueGenerator struct {
	NormalFont font.Face
	ScaryFont  font.Face
}

const NORMALFONT = "/fonts/omori_normal.ttf"
const SCARYFONT = "/fonts/omori_scary.ttf"

func NewDialogueGenerator(st map[settings.BotSetting]string) *DialogueGenerator {
	assetFolder := st[settings.ASSETS]
	normalPath := path.Join(assetFolder, NORMALFONT)
	scaryPath := path.Join(assetFolder, SCARYFONT)

	return &DialogueGenerator{
		NormalFont: loadFont(normalPath),
		ScaryFont:  loadFont(scaryPath),
	}
}

func loadFont(path string) font.Face {
	f, err := gg.LoadFontFace(path, 35)
	if err != nil {
		panic(err)
	}
	return f
}

func (d *DialogueGenerator) GenerateDialogue(text string, font font.Face, width int, height int) image.Image {
	dc := gg.NewContext(width, height)
	dc.SetRGB(0, 0, 0)
	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.SetRGB(255, 255, 255)
	dc.DrawRectangle(2, 2, float64(width)-2, float64(height)-2)
	dc.SetRGB(0, 0, 0)
	dc.DrawRectangle(6, 6, float64(width)-6, float64(height)-6)

	dc.SetFontFace(font)
	dc.SetRGB(255, 255, 255)
	dc.DrawStringWrapped(text, 20, 15, 0, 0, float64(width)-20, 1.5, gg.AlignLeft)

	return dc.Image()
}
