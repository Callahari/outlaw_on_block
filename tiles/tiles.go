package tiles

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type Tile struct {
	TileImage    *ebiten.Image
	Pos          struct{ X, Y int }
	Rotation     int
	CollisionMap *image.Rectangle
}
