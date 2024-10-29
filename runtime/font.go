package runtime

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"log"
	"outlaw_on_block/res"
)

const (
	FONT_NORMAL = iota
	FONT_HOVER
	FONT_ACTIVE
)

var (
	OobFont *text.GoTextFaceSource
)

type FontStatus int

type OOBFontOptions struct {
	Colors struct {
		Normal color.Color
		Hover  color.Color
		Active color.Color
	}
	Size float64
}

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(res.OOBFont))
	if err != nil {
		log.Fatal(err)
	}
	OobFont = s
}

func DrawString(s string, status FontStatus, posX, posY int, screen *ebiten.Image, opt *OOBFontOptions) {
	if opt == nil {
		//Use Default settings
		opt = &OOBFontOptions{
			Colors: struct {
				Normal color.Color
				Hover  color.Color
				Active color.Color
			}{
				Normal: color.RGBA{255, 255, 255, 255},
				Hover:  color.RGBA{128, 128, 128, 255},
				Active: color.RGBA{0, 0, 128, 255},
			},
			Size: 11,
		}
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(posX), float64(posY))
	switch status {
	case FONT_NORMAL:
		op.ColorScale.ScaleWithColor(opt.Colors.Normal)
	case FONT_HOVER:
		op.ColorScale.ScaleWithColor(opt.Colors.Hover)
	case FONT_ACTIVE:
		op.ColorScale.ScaleWithColor(opt.Colors.Active)
	default:
		op.ColorScale.ScaleWithColor(opt.Colors.Normal)
	}
	text.Draw(screen, s, &text.GoTextFace{
		Source: OobFont,
		Size:   opt.Size,
	}, op)
}
