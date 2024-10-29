package modals

import (
	"encoding/json"
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"outlaw_on_block/runtime"
	"outlaw_on_block/tiles"
	"outlaw_on_block/ui"
	"strings"
)

type EsaveMap struct {
	Name        string
	InputStatus runtime.FontStatus
	FileName    string
	HoverSave   bool
	TileMap     []tiles.Tile
	Closed      bool
}

func (this *EsaveMap) GetTileMap() []tiles.Tile {
	return this.TileMap
}

func (this *EsaveMap) IsClosed() bool {
	return this.Closed
}

func (this *EsaveMap) GetName() string {
	return this.Name
}

func (this *EsaveMap) Update() error {
	mouseCursorPosX, mouseCursorPosY := ebiten.CursorPosition()
	mouseCursorTrigger := image.Rect(mouseCursorPosX, mouseCursorPosY, mouseCursorPosX+1, mouseCursorPosY+1)
	inputRect := image.Rect(675, 342, 982, 372)
	saveBtnTrigger := image.Rect(1260, 792, 1285, 807)

	if mouseCursorTrigger.In(saveBtnTrigger) {
		this.HoverSave = true
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if len(this.FileName) <= 0 {
				//Nope .... not filenam was enterd
				this.Closed = true
				log.Println("Nope .... not filename was entered")
				return errors.New("Filename was empty")
			}
			//1. Create folder named saveName
			_ = os.MkdirAll(runtime.OOBUserDir+"/OOB/save/"+this.FileName+"/img", 0755)
			//2. Create inside img
			//_ = os.Mkdir(runtime.OOBUserDir+"/OOB/save/"+this.FileName+"/img", 0755)
			//3. Save all tiles to <saveFolder>/img
			for _, img := range this.TileMap {
				f, _ := os.Create(runtime.OOBUserDir + "/OOB/save/" + this.FileName + "/img/" + img.ID + ".png")
				defer f.Close()
				png.Encode(f, img.TileImage)
			}
			//4. Write manifestfile
			manifest := []tiles.Tile{}
			for _, img := range this.TileMap {
				img.TileImage = nil
				manifest = append(manifest, img)
			}
			manifestAsByte, err := json.Marshal(manifest)
			if err != nil {
				this.Closed = true
				return err
			}
			manifestFile, _ := os.Create(runtime.OOBUserDir + "/OOB/save/" + this.FileName + "/manifest.json")
			manifestFile.Write(manifestAsByte)
			this.Closed = true
			return nil
		}
	} else {
		this.HoverSave = false
	}

	if mouseCursorTrigger.In(inputRect) {
		if this.InputStatus != runtime.FONT_ACTIVE {
			this.InputStatus = runtime.FONT_HOVER
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if this.InputStatus == runtime.FONT_ACTIVE {
				this.InputStatus = runtime.FONT_HOVER
				log.Println("Set inactive")
			} else {
				this.InputStatus = runtime.FONT_ACTIVE
				log.Println("Set active")
			}
		}
	} else {
		if this.InputStatus != runtime.FONT_ACTIVE {
			this.InputStatus = runtime.FONT_NORMAL
		}
	}
	if this.InputStatus == runtime.FONT_ACTIVE {
		if keys := inpututil.AppendJustPressedKeys(nil); len(keys) > 0 {
			for _, key := range keys {
				if len(key.String()) == 1 {
					if !ebiten.IsKeyPressed(ebiten.KeyShift) {
						this.FileName += strings.ToLower(key.String())
					} else {
						this.FileName += strings.ToUpper(key.String())
					}
				} else {
					if key == ebiten.KeyBackspace {
						this.FileName = this.FileName[:len(this.FileName)-1]
					} else if key == ebiten.KeyMinus {
						if ebiten.IsKeyPressed(ebiten.KeyShift) {
							this.FileName += "_"
						} else {
							this.FileName += "-"
						}
					} else if key == ebiten.KeyEnter {
						this.InputStatus = runtime.FONT_NORMAL
					} else if strings.Contains(key.String(), "Digit") {
						s := strings.Split(key.String(), "Digit")
						this.FileName += s[1]
					}
				}
			}
		}
	}
	return nil
}

func (this *EsaveMap) Draw(screen *ebiten.Image) {
	relCoods := struct {
		X float32
		Y float32
	}{
		X: (1920 / 2) - 400,
		Y: (1080 / 2) - 300,
	}

	vector.DrawFilledRect(screen, relCoods.X, relCoods.Y, 800, 600, color.RGBA{128, 128, 128, 255}, false)
	runtime.DrawString("Save Map", runtime.FONT_NORMAL, (int(relCoods.X)+400)-(5*20), int(relCoods.Y)+8, screen, &runtime.OOBFontOptions{
		Colors: struct {
			Normal color.Color
			Hover  color.Color
			Active color.Color
		}{
			Normal: color.RGBA{255, 255, 255, 255},
		},
		Size: 52,
	})

	runtime.DrawString("Filename:", runtime.FONT_NORMAL, (int(relCoods.X) + 32), int(relCoods.Y)+100, screen, &runtime.OOBFontOptions{
		Colors: struct {
			Normal color.Color
			Hover  color.Color
			Active color.Color
		}{
			Normal: color.RGBA{255, 255, 255, 255},
			Hover:  color.RGBA{200, 200, 200, 255},
			Active: color.RGBA{0, 0, 128, 255},
		},
		Size: 24,
	})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(5, 1)
	op.GeoM.Translate((float64(relCoods.X) + 108), float64(relCoods.Y)+100)
	screen.DrawImage(ui.NewInputUI(), op)
	runtime.DrawString(this.FileName+".json", this.InputStatus, (int(relCoods.X) + 118), int(relCoods.Y)+100, screen, &runtime.OOBFontOptions{
		Colors: struct {
			Normal color.Color
			Hover  color.Color
			Active color.Color
		}{
			Normal: color.RGBA{128, 128, 128, 255},
			Hover:  color.RGBA{200, 200, 200, 255},
			Active: color.RGBA{0, 0, 128, 255},
		},
		Size: 24,
	})

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
	if this.HoverSave {
		runtime.DrawString("Save", runtime.FONT_HOVER, (int(relCoods.X) + 700), int(relCoods.Y)+550, screen, btnFontOp)
	} else {
		runtime.DrawString("Save", runtime.FONT_NORMAL, (int(relCoods.X) + 700), int(relCoods.Y)+550, screen, btnFontOp)
	}

}

func (this *EsaveMap) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
