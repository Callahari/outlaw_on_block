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
	"outlaw_on_block/assetManager"
	"outlaw_on_block/car"
	"outlaw_on_block/editor"
	"outlaw_on_block/player"
	"outlaw_on_block/runtime"
	"outlaw_on_block/tiles"
)

const (
	GameScene_Menu GameScene = iota
	GameScene_Play
	GameScene_Editor
	GameScene_AssetManager
)

type GameScene int

type Game struct {
	Player   *player.Player
	Cars     []*car.Car
	TilesMap []*tiles.Tile
	Scene    GameScene
	Menu     struct {
		PlayTriggered         bool
		EditorTriggered       bool
		ExitTriggered         bool
		AssetManagerTriggered bool
	}
	Editor        *editor.Editor
	AssertManager *assetManager.AssetManager
}

func (g *Game) Update() error {
	//Debug exit
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) && ebiten.IsKeyPressed(ebiten.KeyShift) {
		return errors.New("Force Debug super Power Exit")
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		currentMousePosX, currentMousePosY := ebiten.CursorPosition()
		log.Println(currentMousePosX, currentMousePosY)
	}
	switch g.Scene {
	case GameScene_Menu:
		currentMousePosX, currentMousePosY := ebiten.CursorPosition()
		pointer := image.Rect(currentMousePosX, currentMousePosY, currentMousePosX+1, currentMousePosY+1)
		PlayButtonRect := image.Rect(12, 76, 89, 93)
		AssetManagerButtonRect := image.Rect(12, 100, 134, 118)
		EditorButtonRect := image.Rect(12, 123, 62, 143)
		ExitButtonRect := image.Rect(12, 151, 91, 167)

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
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				if g.Editor == nil {
					g.Editor = editor.NewEditor()
				}
				g.Scene = GameScene_Editor
			}
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
		if pointer.Overlaps(AssetManagerButtonRect) {
			g.Menu.AssetManagerTriggered = true
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				if g.AssertManager == nil {
					g.AssertManager = &assetManager.AssetManager{}
				}
				g.Scene = GameScene_AssetManager
			}
		} else {
			g.Menu.AssetManagerTriggered = false
		}
	case GameScene_Play:
		//player Update
		g.Player.Update()
		//car update
		for _, c := range g.Cars {
			c.Update()
		}
	case GameScene_Editor:
		g.Editor.Update()
	case GameScene_AssetManager:
		g.AssertManager.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.Scene {
	case GameScene_Menu:
		runtime.DrawString("Outlaw on Block", 52, 10, 10, false, screen)
		runtime.DrawString("Auf gehts", 24, 10, 70, g.Menu.PlayTriggered, screen)
		runtime.DrawString("Asset Manager", 24, 10, 95, g.Menu.AssetManagerTriggered, screen)
		runtime.DrawString("Editor", 24, 10, 120, g.Menu.EditorTriggered, screen)
		runtime.DrawString("Weg hier !", 24, 10, 145, g.Menu.ExitTriggered, screen)
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
		g.Editor.Draw(screen)
	case GameScene_AssetManager:
		g.AssertManager.Draw(screen)
	}

	msg := fmt.Sprintf(`TPS: %0.2f
FPS: %0.2f
`, ebiten.ActualTPS(), ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
