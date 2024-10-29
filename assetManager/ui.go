package assetManager

import (
	"bytes"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"outlaw_on_block/modals"
	"outlaw_on_block/res"
	"strings"
)

const (
	BUTTON = iota
	BUTTON_HOVER
	BUTTON_PRESSED
)

var (
	Button        *ebiten.Image
	ButtonHover   *ebiten.Image
	ButtonPressed *ebiten.Image
	InputField    *ebiten.Image
)

func init() {
	btnImage, _, _ := image.Decode(bytes.NewReader(res.UIFlatButton01a3))
	Button = ebiten.NewImageFromImage(btnImage)
	btnImage, _, _ = image.Decode(bytes.NewReader(res.UIFlatButton01a2))
	ButtonHover = ebiten.NewImageFromImage(btnImage)
	btnImage, _, _ = image.Decode(bytes.NewReader(res.UIFlatButton01a1))
	ButtonPressed = ebiten.NewImageFromImage(btnImage)

	inputImage, _, _ := image.Decode(bytes.NewReader(res.UIFlatInputField01a))
	InputField = ebiten.NewImageFromImage(inputImage)
}

type ButtonState int

type AssetManager struct {
	Debug  bool
	Search struct {
		SearchText string
		SearchMode bool
	}
	NewEntryBtnState ButtonState
	Modal            modals.IModal
}

func (this *AssetManager) Update() error {
	if this.Modal != nil {
		if this.Modal.IsClosed() {
			this.Modal = nil
			return nil
		}
		this.Modal.Update()
		return nil
	}
	cursorX, cursorY := ebiten.CursorPosition()
	cursorTrigger := image.Rect(cursorX, cursorY, cursorX+1, cursorY+1)

	//Toggle Debug
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		this.Debug = !this.Debug
	}

	//Enter Search mode
	searchFieldTrigger := image.Rect(8, 40, 1020-8, 72)
	if cursorTrigger.In(searchFieldTrigger) {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			this.Search.SearchMode = true
		}
	}
	if this.Search.SearchMode {
		if keys := inpututil.AppendJustPressedKeys(nil); len(keys) > 0 {
			for _, key := range keys {
				if key == ebiten.KeyBackspace {
					this.Search.SearchText = this.Search.SearchText[:len(this.Search.SearchText)-1]
				} else if key == ebiten.KeyEnter {
					this.Search.SearchMode = false
				} else if key == ebiten.KeyEscape {
					this.Search.SearchText = ""
					this.Search.SearchMode = false
				} else {
					if key == ebiten.KeySpace {
						this.Search.SearchText += " "
					} else if strings.Contains(key.String(), "Digit") {
						s := strings.Split(key.String(), "Digit")
						this.Search.SearchText += s[1]
					} else if len(key.String()) == 1 {
						this.Search.SearchText += key.String()
					}
				}
			}
		}
	}
	if cursorTrigger.In(image.Rect(14, 82, 323, 112)) {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			this.NewEntryBtnState = BUTTON_PRESSED
			this.Modal = modals.NewASNewRecord("NewASRecord")
		} else {
			this.NewEntryBtnState = BUTTON_HOVER
		}
	} else {
		this.NewEntryBtnState = BUTTON
	}
	return nil
}

func (this *AssetManager) Draw(screen *ebiten.Image) {
	//runtime.DrawString("Asset Manager", 1, 8, 8, false, screen)
	//Define Areas
	//search bar
	vector.DrawFilledRect(screen, 8, 40+32, 1920-16, 1, color.RGBA{128, 128, 128, 255}, false)
	//search input
	//runtime.DrawString("Suche: ", 1, 8, 45, false, screen)
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(16, 1)
	op.GeoM.Translate(128, 42)
	screen.DrawImage(InputField, &op)
	//runtime.DrawString(this.Search.SearchText, 1, 140, 45, false, screen)

	//New DB entry
	//vector.DrawFilledRect(screen, 8, 32+32+16, 1920-16, 32, color.RGBA{0, 128, 0, 255}, false)
	op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(8.7, 1)
	op.GeoM.Translate(8, 64+16)
	switch this.NewEntryBtnState {
	case BUTTON:
		screen.DrawImage(Button, &op)
	case BUTTON_HOVER:
		screen.DrawImage(ButtonHover, &op)
	case BUTTON_PRESSED:
		screen.DrawImage(ButtonPressed, &op)
	default:
		screen.DrawImage(Button, &op)
	}

	//runtime.DrawString("New Record", 1, 15, 87, false, screen)

	//Database table
	vector.DrawFilledRect(screen, 8, 32+32+16+32+8, 1920-16, 1020-(32+8+32), color.RGBA{0, 0, 128, 255}, false)
	//Database item
	vector.DrawFilledRect(screen, 16, 32+32+16+32+8+32+8, 1920-32, 32, color.RGBA{128, 128, 0, 255}, false)

	//Draw Modal if set
	if this.Modal != nil {
		vector.DrawFilledRect(screen, 0, 0, 1920, 1080, color.RGBA{0, 0, 0, 128}, false)
		this.Modal.Draw(screen)
	}

	//Debug
	if this.Debug {
		msg := fmt.Sprintf("Asset Manager: %v\n", spew.Sdump(this))
		ebitenutil.DebugPrintAt(screen, msg, 1, 32)
	}
}

func (this *AssetManager) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
