package dialogue

import (
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/meir/Sweetheart/internal/pkg/settings"
	"golang.org/x/image/math/fixed"
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

const FONTSIZE = 35

func (d *DialogueGenerator) GenerateDialogue(text string, font *truetype.Font, width int, height int) *image.RGBA {
	fg, bg := image.White, image.Black
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	draw.Draw(rgba, rgba.Bounds().Inset(2), fg, image.ZP, draw.Src)
	draw.Draw(rgba, rgba.Bounds().Inset(6), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetFont(font)
	c.SetFontSize(FONTSIZE)
	c.SetClip(rgba.Bounds().Inset(10))
	c.SetDst(rgba)
	c.SetSrc(fg)

	pt := freetype.Pt(20, 15+int(c.PointToFixed(FONTSIZE/2)>>6))
	for _, word := range strings.Split(text, " ") {
		length, err := c.DrawString(word, pt)
		if err != nil {
			log.Println(err)
			return rgba
		}
		pt.X += d.getSize(" ", font) + length.X
		if pt.X+d.getSize(word, font) >= c.PointToFixed(float64(rgba.Bounds().Inset(15).Max.X)) {
			pt.X = c.PointToFixed(20)
			pt.Y += c.PointToFixed(FONTSIZE / 1.25)
		}
	}

	return rgba
}

func (d *DialogueGenerator) getSize(str string, font *truetype.Font) fixed.Int26_6 {
	var size fixed.Int26_6
	for _, char := range str {
		hmet := font.HMetric(FONTSIZE, font.Index(char))
		size += hmet.AdvanceWidth
	}
	return size
}
