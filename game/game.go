package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"math"
	"outlaw_on_block/car"
	"outlaw_on_block/player"
	"outlaw_on_block/runtime"
	"outlaw_on_block/tiles"
)

type Game struct {
	Player   *player.Player
	Cars     []*car.Car
	TilesMap []*tiles.Tile
}

func (g *Game) Update() error {
	//player Update
	g.Player.Update()
	//car update
	for _, c := range g.Cars {
		c.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	//DrawTiles
	for _, t := range g.TilesMap {
		w, h := float64(t.TileImage.Bounds().Size().X), float64(t.TileImage.Bounds().Size().Y)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-w/2, -h/2)
		op.GeoM.Rotate(float64(t.Rotation%360.0) * 2 * math.Pi / 360)
		op.GeoM.Translate(w/2, h/2)
		op.GeoM.Translate(float64(t.Pos.X)+runtime.ViewPort.X, float64(t.Pos.Y)+runtime.ViewPort.Y)
		screen.DrawImage(t.TileImage, op)
	}
	//Draw Player
	g.Player.Draw(screen)
	//Draw Cars
	for _, c := range g.Cars {
		c.Draw(screen)
	}
	msg := fmt.Sprintf(`TPS: %0.2f
FPS: %0.2f
`, ebiten.ActualTPS(), ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
