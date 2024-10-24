package runtime

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"outlaw_on_block/res"
	"strings"
)

var (
	OOBFontSprites      map[string]*ebiten.Image
	OOBFontHoverSprites map[string]*ebiten.Image
	Alphabet            = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
)

func init() {
	fontSpriteAsByte, _, err := image.Decode(bytes.NewReader(res.OOBFont))
	if err != nil {
		log.Printf("%v", err)
		return
	}
	newEbitenImage := ebiten.NewImageFromImage(fontSpriteAsByte)
	OOBFontSprites = make(map[string]*ebiten.Image)
	for y := 0; y < newEbitenImage.Bounds().Size().Y/20; y++ {
		newSpriteImage := ebiten.NewImage(25, 20)
		op := &ebiten.DrawImageOptions{}
		newSpriteImage.DrawImage(newEbitenImage.SubImage(image.Rect(0, y*20, 25, (y+1)*20)).(*ebiten.Image), op)
		OOBFontSprites[Alphabet[y]] = newSpriteImage
	}

	fontSpriteAsByte, _, err = image.Decode(bytes.NewReader(res.OOBFontHover))
	if err != nil {
		log.Printf("%v", err)
		return
	}
	newEbitenImage = ebiten.NewImageFromImage(fontSpriteAsByte)
	OOBFontHoverSprites = make(map[string]*ebiten.Image)
	for y := 0; y < newEbitenImage.Bounds().Size().Y/20; y++ {
		newSpriteImage := ebiten.NewImage(25, 20)
		op := &ebiten.DrawImageOptions{}
		newSpriteImage.DrawImage(newEbitenImage.SubImage(image.Rect(0, y*20, 25, (y+1)*20)).(*ebiten.Image), op)
		OOBFontHoverSprites[Alphabet[y]] = newSpriteImage
	}
}

func DrawString(s string, scale, posX, posY int, hover bool, screen *ebiten.Image) {
	for idx, char := range s {
		if string(char) != " " {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(posX+(idx*20)), float64(posY))
			op.GeoM.Scale(float64(scale), float64(scale))
			if hover {
				screen.DrawImage(OOBFontHoverSprites[strings.ToUpper(string(char))], op)
			} else {
				screen.DrawImage(OOBFontSprites[strings.ToUpper(string(char))], op)
			}
		}
	}
}
