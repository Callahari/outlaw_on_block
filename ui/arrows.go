package ui

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"outlaw_on_block/res"
)

var (
	Arrows     map[string]*ebiten.Image
	ArrowNames = []string{"green", "marine", "black", "blue", "yellow", "orange", "red"}
)

func init() {
	fontSpriteAsByte, _, err := image.Decode(bytes.NewReader(res.UIArrows))
	if err != nil {
		log.Printf("%v", err)
		return
	}
	newEbitenImage := ebiten.NewImageFromImage(fontSpriteAsByte)
	Arrows = make(map[string]*ebiten.Image)
	for x := 0; x < newEbitenImage.Bounds().Size().X/11; x++ {
		newSpriteImage := ebiten.NewImage(11, 18)
		op := &ebiten.DrawImageOptions{}
		newSpriteImage.DrawImage(newEbitenImage.SubImage(image.Rect(x*11, 0, (x+1)*11, 18)).(*ebiten.Image), op)
		Arrows[ArrowNames[x]] = newSpriteImage
	}
}

func GetArrow(name string) *ebiten.Image {
	if _, ok := Arrows[name]; !ok {
		log.Printf("Arrow %s not found", name)
		return nil
	}
	return Arrows[name]
}
