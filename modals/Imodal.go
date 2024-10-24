package modals

import "github.com/hajimehoshi/ebiten/v2"

type IModal interface {
	GetName() string
	IsClosed() bool
	Update() error
	Draw(screen *ebiten.Image)
	Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
}
