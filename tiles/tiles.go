package tiles

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type Tile struct {
	ID           string
	Name         string
	TileImage    *ebiten.Image
	Pos          struct{ X, Y int }
	Rotation     int
	CollisionMap *image.Rectangle
}
