package game

import (
	"errors"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"image/color"
	"log"
	"math"
	"outlaw_on_block/assetManager"
	"outlaw_on_block/car"
	"outlaw_on_block/editor"
	"outlaw_on_block/modals"
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

var (
	menuItemTextOp *runtime.OOBFontOptions
)

type GameScene int

type Game struct {
	Player   *player.Player
	Cars     []*car.Car
	TilesMap []tiles.Tile
	Scene    GameScene
	Modal    modals.IModal
	Menu     struct {
		PlayTriggered         runtime.FontStatus
		EditorTriggered       runtime.FontStatus
		ExitTriggered         runtime.FontStatus
		AssetManagerTriggered runtime.FontStatus
	}
	Editor        *editor.Editor
	AssertManager *assetManager.AssetManager
}

func init() {
	menuItemTextOp = &runtime.OOBFontOptions{
		Colors: struct {
			Normal color.Color
			Hover  color.Color
			Active color.Color
		}{
			Normal: color.RGBA{255, 255, 255, 255},
			Hover:  color.RGBA{128, 128, 128, 255},
		},
		Size: 24,
	}
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

	if g.Modal != nil {
		if g.Modal.IsClosed() {
			if g.Modal.GetTileMap() != nil {
				g.TilesMap = g.Modal.GetTileMap()
				g.Player = player.NewPlayer()
				g.Player.Position.X = g.Modal.GetPlayerObject().Position.X
				g.Player.Position.Y = g.Modal.GetPlayerObject().Position.Y
				for idx, t := range g.TilesMap {
					if t.Name == "playerPos" {
						g.TilesMap = append(g.TilesMap[:idx], g.TilesMap[idx+1:]...)
					}
				}
				g.Scene = GameScene_Play
			}
			g.Modal = nil
			return nil
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.Modal = nil
			return nil
		}
		return g.Modal.Update()
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
			g.Menu.PlayTriggered = runtime.FONT_HOVER
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				//g.Scene = GameScene_Play
				g.Modal = modals.NewEloadMapModal("", g.TilesMap)
			}
		} else {
			g.Menu.PlayTriggered = runtime.FONT_NORMAL
		}

		if pointer.Overlaps(EditorButtonRect) {
			g.Menu.EditorTriggered = runtime.FONT_HOVER
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				if g.Editor == nil {
					g.Editor = editor.NewEditor()
				}
				g.Scene = GameScene_Editor
			}
		} else {
			g.Menu.EditorTriggered = runtime.FONT_NORMAL
		}

		if pointer.Overlaps(ExitButtonRect) {
			g.Menu.ExitTriggered = runtime.FONT_HOVER
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				return errors.New("Exit")
			}
		} else {
			g.Menu.ExitTriggered = runtime.FONT_NORMAL
		}
		if pointer.Overlaps(AssetManagerButtonRect) {
			g.Menu.AssetManagerTriggered = runtime.FONT_HOVER
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				if g.AssertManager == nil {
					g.AssertManager = &assetManager.AssetManager{}
				}
				g.Scene = GameScene_AssetManager
			}
		} else {
			g.Menu.AssetManagerTriggered = runtime.FONT_NORMAL
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
		runtime.DrawString("Outlaw on Block", runtime.FONT_NORMAL, 10, 10, screen, &runtime.OOBFontOptions{
			Colors: struct {
				Normal color.Color
				Hover  color.Color
				Active color.Color
			}{
				Normal: color.RGBA{255, 255, 255, 255},
			},
			Size: 52,
		})
		runtime.DrawString("Auf gehts", g.Menu.PlayTriggered, 10, 70, screen, menuItemTextOp)
		runtime.DrawString("Asset Manager", g.Menu.AssetManagerTriggered, 10, 95, screen, menuItemTextOp)
		runtime.DrawString("Editor", g.Menu.EditorTriggered, 10, 120, screen, menuItemTextOp)
		runtime.DrawString("Weg hier !", g.Menu.ExitTriggered, 10, 145, screen, menuItemTextOp)
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

	//Draw Modal
	if g.Modal != nil {
		g.Modal.Draw(screen)
	}

	msg := fmt.Sprintf(`TPS: %0.2f
FPS: %0.2f
`, ebiten.ActualTPS(), ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
