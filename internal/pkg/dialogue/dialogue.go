package dialogue

import (
	"image"
	"image/color"
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

func (d *DialogueGenerator) GenerateDialogue(text string, font *truetype.Font, width int, height int) *image.RGBA {
	fg, bg := image.Black, image.White
	ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}
	rgba := image.NewRGBA(image.Rect(0, 0, 640, 480))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetFont(font)
	c.SetFontSize(20)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)

	for i := 0; i < 200; i++ {
		rgba.Set(10, 10+i, ruler)
		rgba.Set(10+i, 10, ruler)
	}

	pt := freetype.Pt(10, 10+int(c.PointToFixed(20)>>6))
	for _, s := range text {
		_, err := c.DrawString(string(s), pt)
		if err != nil {
			log.Println(err)
			return rgba
		}
		pt.Y += c.PointToFixed(20 * 5)
	}

	return rgba
}
