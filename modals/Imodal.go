package modals

import (
	"github.com/hajimehoshi/ebiten/v2"
	"outlaw_on_block/player"
	"outlaw_on_block/tiles"
)

type IModal interface {
	GetName() string
	IsClosed() bool
	GetTileMap() []tiles.Tile
	GetPlayerObject() *player.Player
	Update() error
	Draw(screen *ebiten.Image)
	Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
}
