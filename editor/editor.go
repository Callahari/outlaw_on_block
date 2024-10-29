package editor

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"outlaw_on_block/modals"
	"outlaw_on_block/player"
	"outlaw_on_block/runtime"
	"outlaw_on_block/tiles"
	"outlaw_on_block/ui"
	"path/filepath"
)

var (
	payerMarker *ebiten.Image
)

type Editor struct {
	Tiles                 []*tiles.Tile
	startTile             int
	ArrowRight            string
	ArrowLeft             string
	Selected              *tiles.Tile
	FineJustage           bool
	MapItems              []tiles.Tile
	Modal                 modals.IModal
	Selection             []tiles.Tile
	CopyMode              bool
	DebugMode             bool
	PlayerObject          *player.Player
	SetPlayerPosBtnStatus runtime.FontStatus
	Camera                struct {
		Position struct {
			X float64
			Y float64
		}
		ScrollSpeed float64
	}
}

func init() {
	t := ebiten.NewImage(64, 64)
	t.Fill(color.Transparent)
	m := ebiten.NewImage(21, 21)
	m.Fill(color.RGBA{255, 0, 0, 255})
	mOp := &ebiten.DrawImageOptions{}
	mOp.GeoM.Translate(21, 21)
	t.DrawImage(m, mOp)

	payerMarker = t
}
func NewEditor() *Editor {
	e := &Editor{}
	e.startTile = 0
	e.ArrowRight = "green"
	e.ArrowLeft = "green"
	e.Camera.ScrollSpeed = 5
	e.PlayerObject = &player.Player{}
	e.Selection = make([]tiles.Tile, 0)
	_ = filepath.Walk("/home/callahari/Code/node-io.dev/outlaw_on_block/raw/gta2_tiles", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode().IsRegular() && filepath.Ext(path) == ".png" {
			// PNG-Datei gefunden, lese die Datei ein
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
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		e.DebugMode = !e.DebugMode
	}
	if e.Modal != nil {
		if e.Modal.IsClosed() {
			if e.Modal.GetTileMap() != nil {
				e.MapItems = e.Modal.GetTileMap()
				e.PlayerObject = e.Modal.GetPlayerObject()
			}
			e.Modal = nil
			return nil
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			e.Modal = nil
			return nil
		}
		return e.Modal.Update()
	}
	runtime.ViewPort.X = 1654/2 + e.Camera.Position.X
	runtime.ViewPort.Y = 979/2 + e.Camera.Position.Y
	currentMousePosX, currentMousePosY := ebiten.CursorPosition()
	mouseMapOffsetX := currentMousePosX - 232
	mouseMapOffsetY := currentMousePosY - 58

	relX := mouseMapOffsetX - int(e.Camera.Position.X) - 598
	relY := mouseMapOffsetY - int(e.Camera.Position.Y) - 425
	cursorTrigger := image.Rect(currentMousePosX, currentMousePosY, currentMousePosX+1, currentMousePosY+1)
	relCursorTrigger := image.Rect(relX, relY, relX+1, relY+1)
	// 	vector.DrawFilledRect(screen, 261, 96, 1654, 979, color.RGBA{0, 255, 0, 2}, true)
	mapRect := image.Rect(261, 96, 1654+261, 979+96)
	setPlayerPosBtnTrigger := image.Rect(300, 15, 457, 34)

	//SetPlayerPos Button interaction
	if cursorTrigger.In(setPlayerPosBtnTrigger) {
		if e.SetPlayerPosBtnStatus != runtime.FONT_ACTIVE {
			e.SetPlayerPosBtnStatus = runtime.FONT_HOVER
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if e.SetPlayerPosBtnStatus == runtime.FONT_ACTIVE {
				e.SetPlayerPosBtnStatus = runtime.FONT_HOVER
			} else {
				e.SetPlayerPosBtnStatus = runtime.FONT_ACTIVE
				e.Selected = nil
				for idx, t := range e.MapItems {
					if t.Name == "playerPos" {
						e.MapItems = append(e.MapItems[:idx], e.MapItems[idx+1:]...)
					}
				}
			}
		}
	} else {
		if e.SetPlayerPosBtnStatus != runtime.FONT_ACTIVE {
			e.SetPlayerPosBtnStatus = runtime.FONT_NORMAL
		}
	}

	//Try to enter Copy mode
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyC) {
		if len(e.Selection) > 0 {
			e.CopyMode = !e.CopyMode
		}
	}

	//Click on Save 	runtime.DrawString("Save Map", 1, 1700, 10, false, screen)
	saveBtnRect := image.Rect(1700, 10, 1860, 32)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && cursorTrigger.In(saveBtnRect) {
		if e.MapItems == nil || len(e.MapItems) == 0 {
			log.Println("TileMap is empty, nothing to save.")
		} else {
			m := &modals.EsaveMap{Name: "Foo", TileMap: e.MapItems, PlayerObject: e.PlayerObject}
			e.Modal = m
		}
	}
	loadBtnRect := image.Rect(1707, 38, 1857, 58)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && cursorTrigger.In(loadBtnRect) {
		m := modals.NewEloadMapModal("eladap", nil)
		e.Modal = m
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
			if e.SetPlayerPosBtnStatus == runtime.FONT_ACTIVE {
				t := tiles.Tile{
					ID:        uuid.NewString(),
					Name:      "playerPos",
					TileImage: payerMarker,
					Pos: struct{ X, Y int }{
						X: relX - 32,
						Y: relY - 32,
					},
					Rotation:     0,
					CollisionMap: nil,
				}
				e.MapItems = append(e.MapItems, t)
				e.SetPlayerPosBtnStatus = runtime.FONT_NORMAL
				e.PlayerObject.Position.X = float64(relX - 32)
				e.PlayerObject.Position.Y = float64(relY - 32)
				return nil
			}
			//Check if Teile already pleased.
			isPlaced := false
			tileIdx := -1
			tile := tiles.Tile{}
			for idx, t := range e.MapItems {
				trigger := image.Rect(t.Pos.X, t.Pos.Y, t.Pos.X+64, t.Pos.Y+64)
				if relCursorTrigger.In(trigger) {
					isPlaced = true
					tileIdx = idx
					tile = t
					log.Printf("Tile already placed: %v\n", t.Name)
					break
				}
			}
			if ebiten.IsKeyPressed(ebiten.KeyControl) {
				//Add Tile to Selection
				if isPlaced {
					onSelection := false
					for idx, t := range e.Selection {
						trigger := image.Rect(t.Pos.X, t.Pos.Y, t.Pos.X+64, t.Pos.Y+64)
						if relCursorTrigger.In(trigger) {
							e.Selection = append(e.Selection[:idx], e.Selection[idx+1:]...)
							onSelection = true
							break
						}
					}
					if !onSelection {
						e.Selection = append(e.Selection, tile)
					}
				}
				return nil
			}
			log.Printf("isPlaced: %v\n", isPlaced)
			if e.Selected != nil {
				log.Printf("rel.X: %v; rel.Y: %v\n", relX, relY)
				log.Printf("tile.X: %v; tile.Y: %v\n", relX/64, relY/64)
				e.Selected.Pos.X = int(relX/64) * 64
				e.Selected.Pos.Y = int(relY/64) * 64

				e.MapItems = append(e.MapItems, *e.Selected)
				//e.Selected = nil

				return nil
			} else {
				if isPlaced {
					e.MapItems = append(e.MapItems[:tileIdx], e.MapItems[tileIdx+1:]...)
				}
			}
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
func (e *Editor) drawTileMenu(screen *ebiten.Image) {
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
}
func (e *Editor) Draw(screen *ebiten.Image) {
	currentMousePosX, currentMousePosY := ebiten.CursorPosition()
	f := &text.GoTextFace{
		Source: runtime.OobFont,
		Size:   52,
	}
	textW, _ := text.Measure("OoB Editor", f, 1)
	runtime.DrawString("OoB Editor", runtime.FONT_NORMAL, int((1920/2)-(textW/2)), 10, screen, &runtime.OOBFontOptions{
		Colors: struct {
			Normal color.Color
			Hover  color.Color
			Active color.Color
		}{
			Normal: color.RGBA{255, 255, 255, 255},
		},
		Size: 52,
	})
	runtime.DrawString("Set Player Position", e.SetPlayerPosBtnStatus, 300, 10, screen, &runtime.OOBFontOptions{
		Colors: struct {
			Normal color.Color
			Hover  color.Color
			Active color.Color
		}{
			Normal: color.White,
			Hover:  color.RGBA{200, 200, 200, 255},
			Active: color.RGBA{0, 0, 128, 255},
		},
		Size: 24,
	})
	//Draw map grid 64x64 px
	for y := range 1024 {
		vector.StrokeLine(screen, float32(runtime.ViewPort.X), float32(runtime.ViewPort.Y+(float64(y)*64)), float32(runtime.ViewPort.X+65536), float32(runtime.ViewPort.Y+(float64(y)*64)), 1, color.RGBA{0, 255, 0, 255}, false)
	}
	for x := range 1024 {
		vector.StrokeLine(screen, float32(runtime.ViewPort.X+(float64(x)*64)), float32(runtime.ViewPort.Y), float32(runtime.ViewPort.X+(float64(x)*64)), float32(runtime.ViewPort.Y+65536), 1, color.RGBA{0, 255, 0, 255}, false)
	}
	//Draw Selected marker
	if len(e.Selection) > 0 {
		for _, s := range e.Selection {
			vector.DrawFilledRect(screen, float32(s.Pos.X)+float32(runtime.ViewPort.X)-1, float32(s.Pos.Y)+float32(runtime.ViewPort.Y)-1, 66, 66, color.RGBA{255, 0, 255, 255}, false)
		}
	}
	//draw player marker object
	if e.SetPlayerPosBtnStatus == runtime.FONT_ACTIVE {
		tOp := &ebiten.DrawImageOptions{}
		tOp.GeoM.Translate(float64(currentMousePosX)-32, float64(currentMousePosY)-32)
		screen.DrawImage(payerMarker, tOp)
	}
	//Draw if item selected
	if e.Selected != nil && !e.FineJustage {
		w, h := float64(e.Selected.TileImage.Bounds().Size().X), float64(e.Selected.TileImage.Bounds().Size().Y)

		op := &ebiten.DrawImageOptions{}

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
	btnFontOp := &runtime.OOBFontOptions{
		Colors: struct {
			Normal color.Color
			Hover  color.Color
			Active color.Color
		}{
			Normal: color.RGBA{255, 255, 255, 255},
			Hover:  color.RGBA{200, 200, 200, 255},
		},
		Size: 16,
	}
	runtime.DrawString("Save Map", runtime.FONT_NORMAL, 1700, 10, screen, btnFontOp)
	runtime.DrawString("Load Map", runtime.FONT_NORMAL, 1700, 35, screen, btnFontOp)
	//map
	//vector.DrawFilledRect(screen, 261, 96, 1654, 979, color.RGBA{0, 255, 0, 2}, true)

	e.drawTileMenu(screen)

	//Draw Modal
	if e.Modal != nil {
		e.Modal.Draw(screen)
	}

	if e.DebugMode {
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
		/*msg = fmt.Sprintf("selection: %v", e.Selection)
		ebitenutil.DebugPrintAt(screen, msg, 0, 76)
		msg = fmt.Sprintf("copyMode: %v", e.CopyMode)
		ebitenutil.DebugPrintAt(screen, msg, 0, 87)*/

		yPos := 76
		for idx, m := range e.MapItems {
			msg := fmt.Sprintf("MapItem: %v; Pos: %v", m, m.Pos)
			ebitenutil.DebugPrintAt(screen, msg, 0, yPos+idx)
			yPos += 11
		}
	}
}
func (e *Editor) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
