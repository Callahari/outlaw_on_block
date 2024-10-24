package editor

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"log"
	"math"
	"os"
	"outlaw_on_block/runtime"
	"outlaw_on_block/tiles"
	"outlaw_on_block/ui"
	"path/filepath"
)

type Editor struct {
	Tiles      []*tiles.Tile
	startTile  int
	ArrowRight string
	ArrowLeft  string
}

func NewEditor() *Editor {
	e := &Editor{}
	e.startTile = 0
	e.ArrowRight = "green"
	e.ArrowLeft = "green"
	_ = filepath.Walk("/home/callahari/Code/node-io.dev/outlaw_on_block/raw/gta2_tiles/out2", func(path string, info os.FileInfo, err error) error {
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
			t.TileImage = ebiten.NewImageFromImage(img)

			e.Tiles = append(e.Tiles, t)
		}
		return nil
	})
	return e
}

func (e *Editor) Update() error {
	currentMousePosX, currentMousePosY := ebiten.CursorPosition()
	cursorTrigger := image.Rect(currentMousePosX, currentMousePosY, currentMousePosX+1, currentMousePosY+1)

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

	return nil
}

func (e *Editor) Draw(screen *ebiten.Image) {
	stringLength := (len("OoB Editor") * 20) / 2
	runtime.DrawString("OoB Editor", 3, ((1920/3)/2)-stringLength, 10, false, screen)

	cnt := 0
	row := 0
	for idx, t := range e.Tiles {
		if idx < e.startTile {
			continue
		}
		cnt++
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(cnt*64), float64(row*64))
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
func (e *Editor) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
