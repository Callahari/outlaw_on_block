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
)

type EsaveMap struct {
	Name        string
	HoveredChar string
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
	inputRect := image.Rect(592, 281, 1311, 387)
	saveBtnTrigger := image.Rect(1265, 791, 1343, 814)

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
			_ = os.Mkdir("save/"+this.FileName, 0755)
			//2. Create inside img
			_ = os.Mkdir("save/"+this.FileName+"/img", 0755)
			//3. Save all tiles to <saveFolder>/img
			for _, img := range this.TileMap {
				f, _ := os.Create("save/" + this.FileName + "/img/" + img.ID + ".png")
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
			manifestFile, _ := os.Create("save/" + this.FileName + "/manifest.json")
			manifestFile.Write(manifestAsByte)
			this.Closed = true
			return nil
		}
	} else {
		this.HoverSave = false
	}

	if mouseCursorTrigger.In(inputRect) {
		row := 282
		cnt := 593
		for idx, char := range runtime.Alphabet {
			if idx == 18 {
				row += 51
				cnt = 593
			}
			charTrigger := image.Rect(cnt, row, cnt+41, row+51)
			if mouseCursorTrigger.In(charTrigger) {
				this.HoveredChar = char
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
					this.FileName += char
				}
			}
			cnt += 41
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
	runtime.DrawString("Save Map", 1, (int(relCoods.X)+400)-(5*20), int(relCoods.Y)+8, false, screen)

	//Draw Clickable Alphabet 18
	row := 0
	cnt := 0
	for _, char := range runtime.Alphabet {
		if cnt == 18 {
			row++
			cnt = 0
		}
		if char == this.HoveredChar {
			runtime.DrawString(char, 2, (int(relCoods.X)-270)+(cnt*20), (int(relCoods.Y)-100)+(row*32), true, screen)
		} else {
			runtime.DrawString(char, 2, (int(relCoods.X)-270)+(cnt*20), (int(relCoods.Y)-100)+(row*32), false, screen)
		}
		cnt++
	}

	fileNameString := "Filename: " + this.FileName + ".json"
	runtime.DrawString(fileNameString, 1, (int(relCoods.X) + 32), int(relCoods.Y)+550, false, screen)

	if this.HoverSave {
		runtime.DrawString("Save", 1, (int(relCoods.X) + 700), int(relCoods.Y)+550, true, screen)
	} else {
		runtime.DrawString("Save", 1, (int(relCoods.X) + 700), int(relCoods.Y)+550, false, screen)
	}

}

func (this *EsaveMap) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}