package game

import (
	"errors"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"log"
	"math"
	"outlaw_on_block/car"
	"outlaw_on_block/player"
	"outlaw_on_block/runtime"
	"outlaw_on_block/tiles"
)

const (
	GameScene_Menu GameScene = iota
	GameScene_Play
	GameScene_Editor
)

type GameScene int

type Game struct {
	Player   *player.Player
	Cars     []*car.Car
	TilesMap []*tiles.Tile
	Scene    GameScene
	Menu     struct {
		PlayTriggered   bool
		EditorTriggered bool
		ExitTriggered   bool
	}
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		currentMousePosX, currentMousePosY := ebiten.CursorPosition()
		log.Println(currentMousePosX, currentMousePosY)
	}
	switch g.Scene {
	case GameScene_Menu:
		currentMousePosX, currentMousePosY := ebiten.CursorPosition()
		pointer := image.Rect(currentMousePosX, currentMousePosY, currentMousePosX+1, currentMousePosY+1)
		PlayButtonRect := image.Rect(30, 140, 380, 185)
		EditorButtonRect := image.Rect(30, 190, 260, 230)
		ExitButtonRect := image.Rect(30, 240, 245, 280)

		if pointer.Overlaps(PlayButtonRect) {
			g.Menu.PlayTriggered = true
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				g.Scene = GameScene_Play
			}
		} else {
			g.Menu.PlayTriggered = false
		}

		if pointer.Overlaps(EditorButtonRect) {
			g.Menu.EditorTriggered = true
		} else {
			g.Menu.EditorTriggered = false
		}

		if pointer.Overlaps(ExitButtonRect) {
			g.Menu.ExitTriggered = true
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				return errors.New("Exit")
			}
		} else {
			g.Menu.ExitTriggered = false
		}
	case GameScene_Play:
		//player Update
		g.Player.Update()
		//car update
		for _, c := range g.Cars {
			c.Update()
		}
	case GameScene_Editor:
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.Scene {
	case GameScene_Menu:
		runtime.DrawString("Outlaw on Block", 4, 10, 10, false, screen)
		runtime.DrawString("Auf gehts", 2, 10, 70, g.Menu.PlayTriggered, screen)
		runtime.DrawString("Editor", 2, 10, 95, g.Menu.EditorTriggered, screen)
		runtime.DrawString("Weg hier", 2, 10, 120, g.Menu.ExitTriggered, screen)
	case GameScene_Play:
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
	case GameScene_Editor:
	}

	msg := fmt.Sprintf(`TPS: %0.2f
FPS: %0.2f
`, ebiten.ActualTPS(), ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
