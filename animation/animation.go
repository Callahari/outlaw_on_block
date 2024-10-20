package animation

import "github.com/hajimehoshi/ebiten/v2"

type Animation struct {
	Sprites       map[string][]*ebiten.Image
	AnimationName string
	UpdateCounter int
	SpriteIdx     int
}
