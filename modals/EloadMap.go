package modals

import (
	"bytes"
	"encoding/json"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"log"
	"os"
	"outlaw_on_block/player"
	"outlaw_on_block/runtime"
	"outlaw_on_block/tiles"
	"path/filepath"
	"strings"
)

type EloadMap struct {
	Name       string
	TileMap    []tiles.Tile
	Closed     bool
	HoverFile  string
	SavedFiles []string
	PlayerObject *player.Player
}

func NewEloadMapModal(name string, tileMap []tiles.Tile) *EloadMap {
	e := &EloadMap{}
	e.Name = name
	e.TileMap = tileMap
	e.Closed = false

	//Lookup saved Files
	_ = filepath.Walk(runtime.OOBUserDir+"/OOB/save", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() && filepath.Ext(path) == ".json" {
			strip := strings.Split(path, "/manifest.json")[0]
			e.SavedFiles = append(e.SavedFiles, strip)
		}
		return nil
	})
	return e
}

func (this *EloadMap) GetPlayerObject() *player.Player {
	return this.PlayerObject
}

func (this *EloadMap) GetTileMap() []tiles.Tile {
	return this.TileMap
}

func (this *EloadMap) IsClosed() bool {
	return this.Closed
}

func (this *EloadMap) GetName() string {
	return this.Name
}

func (this *EloadMap) Update() error {
	cursorPosX, cursorPosY := ebiten.CursorPosition()
	cursorTriggerRect := image.Rect(cursorPosX, cursorPosY, cursorPosX+1, cursorPosY+1)

	relCoods := struct {
		X float32
		Y float32
	}{
		X: (1920 / 2) - 400,
		Y: (1080 / 2) - 300,
	}

	for idx, path := range this.SavedFiles {
		saveTriggerRect := image.Rect((int(relCoods.X) + 32), (int(relCoods.Y)+50)+(idx*32), (int(relCoods.X)+32)+300, (int(relCoods.Y)+50)+(idx*32)+32)
		if cursorTriggerRect.In(saveTriggerRect) {
			this.HoverFile = path
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				//1. Read manifest
				loadMap := []tiles.Tile{}
				manifestAsBytes, _ := os.ReadFile(path + "/manifest.json")
				json.Unmarshal(manifestAsBytes, &loadMap)
				//2.Recreate ebiten images
				for _, tile := range loadMap {
					imgAsByte, err := os.ReadFile(path + "/img/" + tile.ID + ".png")
					if err != nil {
						log.Printf("err: %v", err)
						continue
					}
					rawImg, _, _ := image.Decode(bytes.NewReader(imgAsByte))
					tile.TileImage = ebiten.NewImageFromImage(rawImg)
					this.TileMap = append(this.TileMap, tile)
				}
				this.Closed = true
				return nil
			}
		}
	}

	return nil
}

func (this *EloadMap) Draw(screen *ebiten.Image) {
	relCoods := struct {
		X float32
		Y float32
	}{
		X: (1920 / 2) - 400,
		Y: (1080 / 2) - 300,
	}

	vector.DrawFilledRect(screen, relCoods.X, relCoods.Y, 800, 600, color.RGBA{128, 128, 128, 255}, false)
	runtime.DrawString("Load Map", runtime.FONT_NORMAL, (int(relCoods.X)+400)-(5*20), int(relCoods.Y)+8, screen, &runtime.OOBFontOptions{
		Colors: struct {
			Normal color.Color
			Hover  color.Color
			Active color.Color
		}{
			Normal: color.RGBA{255, 255, 255, 255},
		},
		Size: 52,
	})

	for idx, path := range this.SavedFiles {
		if path == this.HoverFile {
			runtime.DrawString(path, runtime.FONT_NORMAL, (int(relCoods.X) + 32), (int(relCoods.Y)+50)+(idx*32), screen, nil)
		} else {
			runtime.DrawString(path, runtime.FONT_NORMAL, (int(relCoods.X) + 32), (int(relCoods.Y)+50)+(idx*32), screen, nil)
		}
	}
}

func (this *EloadMap) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
