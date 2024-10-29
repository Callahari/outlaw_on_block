package ui

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"outlaw_on_block/res"
)

func NewInputUI() *ebiten.Image {
	rawimg, _, _ := image.Decode(bytes.NewReader(res.UIFlatInputField01a))
	return ebiten.NewImageFromImage(rawimg)
}
