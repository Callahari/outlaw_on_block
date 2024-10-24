package editor

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"outlaw_on_block/modals"
	"outlaw_on_block/runtime"
	"outlaw_on_block/tiles"
	"outlaw_on_block/ui"
	"path/filepath"
)

type Editor struct {
	Tiles       []*tiles.Tile
	startTile   int
	ArrowRight  string
	ArrowLeft   string
	Selected    *tiles.Tile
	FineJustage bool
	MapItems    []tiles.Tile
	Modal       modals.IModal
	Camera      struct {
		Position struct {
			X float64
			Y float64
		}
		ScrollSpeed float64
	}
}

func NewEditor() *Editor {
	e := &Editor{}
	e.startTile = 0
	e.ArrowRight = "green"
	e.ArrowLeft = "green"
	e.Camera.ScrollSpeed = 5
	_ = filepath.Walk("/home/callahari/Code/node-io.dev/outlaw_on_block/raw/gta2_tiles", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode().IsRegular() && filepath.Ext(path) == ".png" {
			// PNG-Datei gefunden, lese die Datei ein
			log.Println("Found file: ", path)
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			img, _, err := image.Decode(file)
			if err != nil {
				return err
			}
			t := &tiles.Tile{}
			t.Name = filepath.Base(path)
			t.ID = uuid.New().String()
			t.TileImage = ebiten.NewImageFromImage(img)

			e.Tiles = append(e.Tiles, t)
		}
		return nil
	})
	return e
}

func (e *Editor) Update() error {
	if e.Modal != nil {
		if e.Modal.IsClosed() {
			e.Modal = nil
			return nil
		}
		return e.Modal.Update()
	}
	runtime.ViewPort.X = 1654/2 + e.Camera.Position.X
	runtime.ViewPort.Y = 979/2 + e.Camera.Position.Y
	currentMousePosX, currentMousePosY := ebiten.CursorPosition()
	cursorTrigger := image.Rect(currentMousePosX, currentMousePosY, currentMousePosX+1, currentMousePosY+1)
	// 	vector.DrawFilledRect(screen, 261, 96, 1654, 979, color.RGBA{0, 255, 0, 2}, true)
	mapRect := image.Rect(261, 96, 1654+261, 979+96)

	//Click on Save 	runtime.DrawString("Save Map", 1, 1700, 10, false, screen)
	savBtnRect := image.Rect(1700, 10, 1860, 32)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && cursorTrigger.In(savBtnRect) {
		if e.MapItems == nil || len(e.MapItems) == 0 {
			log.Println("TileMap is empty, nothing to save.")
		} else {
			m := &modals.InputSaveMap{Name: "Foo", TileMap: e.MapItems}
			e.Modal = m
		}
	}

	//Scroll map
	if ebiten.IsKeyPressed(ebiten.KeyUp) && !ebiten.IsKeyPressed(ebiten.KeyControl) {
		e.Camera.Position.Y += e.Camera.ScrollSpeed
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) && !ebiten.IsKeyPressed(ebiten.KeyControl) {
		e.Camera.Position.Y -= e.Camera.ScrollSpeed
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) && !ebiten.IsKeyPressed(ebiten.KeyControl) {
		e.Camera.Position.X += e.Camera.ScrollSpeed
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) && !ebiten.IsKeyPressed(ebiten.KeyControl) {
		e.Camera.Position.X -= e.Camera.ScrollSpeed
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyComma) {
		e.Camera.ScrollSpeed -= 0.1
	} else if inpututil.IsKeyJustPressed(ebiten.KeyPeriod) {
		e.Camera.ScrollSpeed += 0.1
	}

	arrowRightTrigger := image.Rect(50, 1010, 105, 1050)
	if cursorTrigger.In(arrowRightTrigger) {
		e.ArrowRight = "orange"
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			e.ArrowLeft = "red"
			e.startTile -= 15 * 3
			if e.startTile < 0 {
				e.startTile = 0
			}
		} else {
			e.ArrowLeft = "orange"
		}
	} else {
		e.ArrowRight = "green"
	}

	arrowLeftTrigger := image.Rect(215, 1005, 270, 1044)
	if cursorTrigger.In(arrowLeftTrigger) {
		e.ArrowLeft = "orange"
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			e.ArrowLeft = "red"
			e.startTile += 15 * 3
			if e.startTile > len(e.Tiles) {
				e.startTile = len(e.Tiles)
			}
		} else {
			e.ArrowLeft = "orange"
		}
	} else {
		e.ArrowLeft = "green"
	}
	//Disable Fine Justage Mode
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyShift) {
		if e.FineJustage {
			e.FineJustage = false
			e.Selected = nil
		}
	}
	//Fine Justage Mode
	if ebiten.IsKeyPressed(ebiten.KeyControl) && e.FineJustage {
		change := false
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) && ebiten.IsKeyPressed(ebiten.KeyControl) {
			e.Selected.Pos.Y -= 1
			change = true
		} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) && ebiten.IsKeyPressed(ebiten.KeyControl) {
			e.Selected.Pos.Y += 1
			change = true
		} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && ebiten.IsKeyPressed(ebiten.KeyControl) {
			e.Selected.Pos.X -= 1
			change = true
		} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) && ebiten.IsKeyPressed(ebiten.KeyControl) {
			e.Selected.Pos.X += 1
			change = true
		}

		if change {
			for k, v := range e.MapItems {
				if v.ID == e.Selected.ID {
					e.MapItems = append(e.MapItems[:k], e.MapItems[k+1:]...)
					e.MapItems = append(e.MapItems, *e.Selected)
				}
			}
		}
	}
	//Click on map
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if cursorTrigger.In(mapRect) {
			if ebiten.IsKeyPressed(ebiten.KeyControl) {
				log.Println("Mouse Clicked on map with control")
				for _, t := range e.MapItems {
					triggerRect := image.Rect(t.Pos.X+int(runtime.ViewPort.X), t.Pos.Y+int(runtime.ViewPort.Y), (t.Pos.X+int(runtime.ViewPort.X))+64, (t.Pos.Y+int(runtime.ViewPort.Y))+64)
					if cursorTrigger.In(triggerRect) {
						if e.Selected == nil {
							log.Println("Fine Justage from Tile")
							e.Selected = &t
							e.FineJustage = true
							return nil
						}
					}
				}
				return nil
			}
			log.Println("Mouse Clicked on map")
			if e.Selected != nil {
				mouseMapOffsetX := currentMousePosX - 261
				mouseMapOffsetY := currentMousePosY - 96
				log.Printf("mouseMapOffsetX: %v; mouseMapOffsetY: %v\n", mouseMapOffsetX, mouseMapOffsetY)
				mapTile := e.Selected
				mapTile.Pos.X = mouseMapOffsetX - int(e.Camera.Position.X) - 598
				mapTile.Pos.Y = mouseMapOffsetY - int(e.Camera.Position.Y) - 425
				e.MapItems = append(e.MapItems, *mapTile)
				log.Printf("place MapItem: %v; Pos: %v\n", mapTile, mapTile.Pos)
				log.Printf("cursor.pos: %v %v\n", currentMousePosX, currentMousePosY)
				log.Printf("runtime.ViewPort: %v\n", runtime.ViewPort)
				log.Printf("camera.pos: %v\n", e.Camera.Position)
				e.Selected = nil
				return nil
			} else {
				//Grep again
				for idx, t := range e.MapItems {
					triggerRect := image.Rect(t.Pos.X+int(runtime.ViewPort.X), t.Pos.Y+int(runtime.ViewPort.Y), (t.Pos.X+int(runtime.ViewPort.X))+64, (t.Pos.Y+int(runtime.ViewPort.Y))+64)
					if cursorTrigger.In(triggerRect) {
						e.Selected = &t
						e.MapItems = append(e.MapItems[:idx], e.MapItems[idx+1:]...)
						break
					}
				}
			}
			return nil
		}
	}
	//Rotate tile
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		e.Selected.Rotation += 90
	}
	///SELECT item from Menu
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		if e.Selected != nil {
			e.Selected = nil
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		log.Println("Mouse Clicked")
		cnt := 0
		row := 0
		topOffset := 8
		targetRect := image.Rectangle{}
		currentMousePosX, currentMousePosY = ebiten.CursorPosition()
		for idx, t := range e.Tiles {
			if idx < e.startTile {
				continue
			}
			cnt++
			if row == 1 {
				targetRect = image.Rect(cnt*64, (row*64)+topOffset, (cnt*64)+64, (row*64)+64+topOffset)
			} else {
				targetRect = image.Rect(cnt*64, row*64, (cnt*64)+64, (row*64)+64)
			}
			if cursorTrigger.In(targetRect) {
				log.Printf("idx: %v; cnt: %v; row: %v;img: %s \n", idx, cnt, row, t.Name)
				e.Selected = t
				break
			}

			if cnt == 3 {
				cnt = 0
				row++
				if row == 15 {
					break
				}
			}
		}
	}
	return nil
}

func (e *Editor) Draw(screen *ebiten.Image) {
	stringLength := (len("OoB Editor") * 20) / 2
	runtime.DrawString("OoB Editor", 3, ((1920/3)/2)-stringLength, 10, false, screen)

	cnt := 0
	row := 0
	topOffset := 8
	for idx, t := range e.Tiles {
		if idx < e.startTile {
			continue
		}
		cnt++
		op := &ebiten.DrawImageOptions{}
		if row == 1 {
			op.GeoM.Translate(float64(cnt*64), float64(row*64))
			vector.DrawFilledRect(screen, float32(cnt*64), float32(row*64), 64, 64, color.RGBA{0, 255, 0, 64}, true)

		} else {
			op.GeoM.Translate(float64(cnt*64), float64((row*64)+topOffset))
			vector.DrawFilledRect(screen, float32(cnt*64), float32((row*64)+topOffset), 64, 64, color.RGBA{0, 255, 0, 64}, true)

		}
		screen.DrawImage(t.TileImage, op)
		if cnt == 3 {
			cnt = 0
			row++
			if row == 15 {
				break
			}
		}
	}

	//
	aR := ui.GetArrow(e.ArrowRight)
	aL := ui.GetArrow(e.ArrowLeft)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(-aR.Bounds().Dx()/2), float64(-aR.Bounds().Dy()/2))
	op.GeoM.Rotate(float64(-90%360) * 2 * math.Pi / 360)
	op.GeoM.Translate(float64(aR.Bounds().Dx()/2), float64(aR.Bounds().Dy()/2))
	op.GeoM.Scale(3, 3)
	op.GeoM.Translate(64, 1000)
	screen.DrawImage(aR, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(-aL.Bounds().Dx()/2), float64(-aL.Bounds().Dy()/2))
	op.GeoM.Rotate(float64(90%360) * 2 * math.Pi / 360)
	op.GeoM.Translate(float64(aL.Bounds().Dx()/2), float64(aL.Bounds().Dy()/2))
	op.GeoM.Scale(3, 3)
	op.GeoM.Translate(228, 995)
	screen.DrawImage(aL, op)

	//Draw if item selected
	if e.Selected != nil && !e.FineJustage {
		w, h := float64(e.Selected.TileImage.Bounds().Size().X), float64(e.Selected.TileImage.Bounds().Size().Y)

		op := &ebiten.DrawImageOptions{}
		currentMousePosX, currentMousePosY := ebiten.CursorPosition()

		op.GeoM.Translate(-w/2, -h/2)
		op.GeoM.Rotate(float64(e.Selected.Rotation%360.0) * 2 * math.Pi / 360)
		op.GeoM.Translate(w/2, h/2)

		op.GeoM.Translate(float64(currentMousePosX)-32, float64(currentMousePosY)-32)
		op.ColorScale.ScaleAlpha(0.5)
		screen.DrawImage(e.Selected.TileImage, op)
	}
	//Draw MapTiles
	for _, m := range e.MapItems {
		w, h := float64(m.TileImage.Bounds().Size().X), float64(m.TileImage.Bounds().Size().Y)

		op := &ebiten.DrawImageOptions{}

		op.GeoM.Translate(-w/2, -h/2)
		op.GeoM.Rotate(float64(m.Rotation%360.0) * 2 * math.Pi / 360)
		op.GeoM.Translate(w/2, h/2)

		op.GeoM.Translate(float64(m.Pos.X)+runtime.ViewPort.X, float64(m.Pos.Y)+runtime.ViewPort.Y)

		screen.DrawImage(m.TileImage, op)
	}

	//Draw Save Map
	runtime.DrawString("Save Map", 1, 1700, 10, false, screen)
	runtime.DrawString("Load Map", 1, 1700, 35, false, screen)
	//map
	//vector.DrawFilledRect(screen, 261, 96, 1654, 979, color.RGBA{0, 255, 0, 2}, true)

	//Draw Modal
	if e.Modal != nil {
		e.Modal.Draw(screen)
	}

	//Debug
	msg := fmt.Sprintf("camera.pos: %v; camera.speed: %.2f", e.Camera.Position, e.Camera.ScrollSpeed)
	ebitenutil.DebugPrintAt(screen, msg, 0, 32)
	msg = fmt.Sprintf("runtime.viewport: %v", runtime.ViewPort)
	ebitenutil.DebugPrintAt(screen, msg, 0, 43)
	mX, mY := ebiten.CursorPosition()
	msg = fmt.Sprintf("cursor.pos: %v %v", mX, mY)
	ebitenutil.DebugPrintAt(screen, msg, 0, 54)
	msg = fmt.Sprintf("fineMode: %v", e.FineJustage)
	ebitenutil.DebugPrintAt(screen, msg, 0, 65)

	yPos := 76
	for idx, m := range e.MapItems {
		msg := fmt.Sprintf("MapItem: %v; Pos: %v", m, m.Pos)
		ebitenutil.DebugPrintAt(screen, msg, 0, yPos+idx)
		yPos += 11
	}
}
func (e *Editor) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
