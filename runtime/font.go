package runtime

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"log"
	"outlaw_on_block/res"
)

var (
	OobFont *text.GoTextFaceSource
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(res.OOBFont))
	if err != nil {
		log.Fatal(err)
	}
	OobFont = s
}

func DrawString(s string, size, posX, posY int, hover bool, screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(posX), float64(posY))
	if hover {
		op.ColorScale.ScaleWithColor(color.RGBA{128, 128, 128, 255})
	} else {
		op.ColorScale.ScaleWithColor(color.White)
	}
	text.Draw(screen, s, &text.GoTextFace{
		Source: OobFont,
		Size:   float64(size),
	}, op)
}
