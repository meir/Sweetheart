package dialogue

import (
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"path"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/meir/Sweetheart/internal/pkg/settings"
)

type DialogueGenerator struct {
	NormalFont *truetype.Font
	ScaryFont  *truetype.Font
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

func loadFont(path string) *truetype.Font {
	fontBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}
	return f
}

const FONTSIZE = 40

func (d *DialogueGenerator) GenerateDialogue(text string, font *truetype.Font, width int, height int) *image.RGBA {
	fg, bg := image.White, image.Black
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	draw.Draw(rgba, rgba.Bounds().Inset(2), fg, image.ZP, draw.Src)
	draw.Draw(rgba, rgba.Bounds().Inset(6), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetFont(font)
	c.SetFontSize(FONTSIZE)
	c.SetClip(rgba.Bounds().Inset(15))
	c.SetDst(rgba)
	c.SetSrc(fg)

	pt := freetype.Pt(20, 20+int(c.PointToFixed(FONTSIZE)>>6))
	for _, char := range text {
		_, err := c.DrawString(string(char), pt)
		if err != nil {
			log.Println(err)
			return rgba
		}
		pt.X += c.PointToFixed(FONTSIZE / 3)
		if pt.X >= c.PointToFixed(float64(rgba.Bounds().Inset(15).Max.X)) {
			pt.X = c.PointToFixed(20)
			pt.Y += c.PointToFixed(FONTSIZE / 3)
		}
	}

	return rgba
}
