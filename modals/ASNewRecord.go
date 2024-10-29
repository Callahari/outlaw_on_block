package modals

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"log"
	"os"
	"outlaw_on_block/databases"
	"outlaw_on_block/res"
	"outlaw_on_block/runtime"
	"outlaw_on_block/tiles"
	"outlaw_on_block/ui"
	"strings"
)

const (
	INPUTNAME_NAME     = "name"
	INPUTNAME_TAGS     = "tags"
	INPUTNAME_FILENAME = "filename"
)

var (
	UIFrame            *ebiten.Image
	UIFrame03a         *ebiten.Image
	UIInput            *ebiten.Image
	UIButton           *ebiten.Image
	UIButtonHover      *ebiten.Image
	UIButtonPressed    *ebiten.Image
	UIButtonExtraCross *ebiten.Image
	UIButtonExtraCheck *ebiten.Image
)

func init() {
	frameImg, _, _ := image.Decode(bytes.NewReader(res.UIFlatFrame01a))
	UIFrame = ebiten.NewImageFromImage(frameImg)

	inputImg, _, _ := image.Decode(bytes.NewReader(res.UIFlatInputField01a))
	UIInput = ebiten.NewImageFromImage(inputImg)

	uiImg, _, _ := image.Decode(bytes.NewReader(res.UIFlatButton01a3))
	UIButton = ebiten.NewImageFromImage(uiImg)

	uiImg, _, _ = image.Decode(bytes.NewReader(res.UIFlatButton01a2))
	UIButtonHover = ebiten.NewImageFromImage(uiImg)

	uiImg, _, _ = image.Decode(bytes.NewReader(res.UIFlatButton01a1))
	UIButtonPressed = ebiten.NewImageFromImage(uiImg)

	uiImg, _, _ = image.Decode(bytes.NewReader(res.UIFlatButtonCross01a))
	UIButtonExtraCross = ebiten.NewImageFromImage(uiImg)

	uiImg, _, _ = image.Decode(bytes.NewReader(res.UIFlatButtonCheck01a))
	UIButtonExtraCheck = ebiten.NewImageFromImage(uiImg)

	uiImg, _, _ = image.Decode(bytes.NewReader(res.UIFlatFrame03a))
	UIFrame03a = ebiten.NewImageFromImage(uiImg)
}

type InputName string

type ASNewRecord struct {
	Debug       bool
	Name        string
	Closed      bool
	Buttons     []*ui.Button
	InputActive struct {
		CurrentActive InputName
		Name          struct {
			Active bool
			Input  string
		}
		Tags struct {
			Active bool
			Input  string
		}
		FileName struct {
			Active bool
			Input  string
		}
	}
}

func (this *ASNewRecord) GetTileMap() []tiles.Tile {
	return nil
}

func (this *ASNewRecord) IsClosed() bool {
	return this.Closed
}

func (this *ASNewRecord) GetName() string {
	return this.Name
}
func NewASNewRecord(name string) *ASNewRecord {
	n := &ASNewRecord{
		Name:   name,
		Closed: false,
	}

	cancelBtn := &ui.Button{
		Image: UIButton,
		State: ui.BUTTON_NORMAL,
		Label: "Cancel",
		Icon:  UIButtonExtraCross,
		Position: struct {
			X float64
			Y float64
		}{
			X: float64((UIFrame.Bounds().Max.X*10)/2) + 30,
			Y: float64((UIFrame.Bounds().Max.Y*10)/2) + float64((UIFrame.Bounds().Max.Y * 10)) - 68,
		},
		Scale: struct {
			X float64
			Y float64
		}{X: 5, Y: 1},
		OnClick: func(p map[string]interface{}) {
			modal := p["modal"].(*ASNewRecord)
			modal.Closed = true
			log.Println("Cancel")
		},
	}

	saveBtn := &ui.Button{
		Image: UIButton,
		State: ui.BUTTON_NORMAL,
		Label: "Save",
		Icon:  UIButtonExtraCheck,
		Position: struct {
			X float64
			Y float64
		}{
			X: float64((UIFrame.Bounds().Max.X*10)/2) + float64(UIFrame.Bounds().Max.X*10) - float64(UIButton.Bounds().Max.X*4) - 32,
			Y: float64((UIFrame.Bounds().Max.Y*10)/2) + float64((UIFrame.Bounds().Max.Y * 10)) - 68,
		},
		Scale: struct {
			X float64
			Y float64
		}{X: 4, Y: 1},
		OnClick: func(p map[string]interface{}) {
			/*
				1. Check Name, at least 3 chars
				2. Tags separated by period
				3. Check if input file exists an loadable
				4. Create Asset Record
			*/
			log.Println("Save")
			if len(p["modal"].(*ASNewRecord).InputActive.Name.Input) < 3 {
				log.Println("Name too short")
				return
			}

			tags := []string{}
			if len(p["modal"].(*ASNewRecord).InputActive.Tags.Input) > 0 {
				tags = strings.Split(p["modal"].(*ASNewRecord).InputActive.Tags.Input, ".")
			}

			_, err := os.Stat(runtime.OOBUserDir + "/OOB/images/" + p["modal"].(*ASNewRecord).InputActive.FileName.Input)
			if err != nil {
				log.Println("File Name too short")
				return
			}
			asDBFile, err := os.ReadFile(runtime.OOBUserDir + "/OOB/assets.db")
			if err != nil {
				log.Println("Error opening assets.db")
				return
			}
			assertDB := make([]databases.AssertDB, 0)
			_ = json.Unmarshal(asDBFile, &assertDB)

			newRecord := databases.AssertDB{}
			newRecord.ID = uuid.NewString()
			newRecord.Name = p["modal"].(*ASNewRecord).InputActive.Name.Input
			newRecord.Tags = tags
			newRecord.FileNames = p["modal"].(*ASNewRecord).InputActive.FileName.Input

			assertDB = append(assertDB, newRecord)

			asDBFile, err = json.MarshalIndent(assertDB, "", "  ")
			if err != nil {
				log.Println("Error writing assets.db")
				return
			}
			err = os.WriteFile(runtime.OOBUserDir+"/OOB/assets.db", asDBFile, 0644)
			if err != nil {
				log.Printf("Error writing assets.db: %v", err)
				return
			}

		},
	}
	n.Buttons = append(n.Buttons, saveBtn)
	n.Buttons = append(n.Buttons, cancelBtn)
	return n
}
func (this *ASNewRecord) getAnyInputKey() string {
	if keys := inpututil.AppendJustPressedKeys(nil); len(keys) > 0 {
		for _, key := range keys {
			log.Println(key)
			if key == ebiten.KeyBackspace || key == ebiten.KeyEnter || key == ebiten.KeyEscape {
				return key.String()
			} else {
				if key == ebiten.KeySpace {
					return " "
				} else if strings.Contains(key.String(), "Digit") {
					s := strings.Split(key.String(), "Digit")
					return s[1]
				} else if key == ebiten.KeyPeriod {
					return "."
				} else if key == ebiten.KeySlash {
					return "/"
				} else if len(key.String()) == 1 {
					return key.String()
				}
			}
		}
	}
	return ""
}
func (this *ASNewRecord) setAllInputsToInactive() {
	this.InputActive.Name.Active = false
	this.InputActive.Tags.Active = false
	this.InputActive.FileName.Active = false
	this.InputActive.CurrentActive = ""
}
func (this *ASNewRecord) Update() error {
	cursorX, cursorY := ebiten.CursorPosition()
	cursorTrigger := image.Rect(cursorX, cursorY, cursorX+1, cursorY+1)

	nameInputRect := image.Rect(699, 419, 1313, 440)
	tagsInputRect := image.Rect(697, 466, 1312, 488)
	fileNameInputRect := image.Rect(698, 516, 1313, 535)

	//Toggle Debug
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		this.Debug = !this.Debug
	}

	if cursorTrigger.In(nameInputRect) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if this.InputActive.Name.Active {
			this.InputActive.Name.Active = false
			this.InputActive.CurrentActive = ""
		} else {
			this.setAllInputsToInactive()
			this.InputActive.Name.Active = true
			this.InputActive.CurrentActive = INPUTNAME_NAME
		}
	}

	if cursorTrigger.In(tagsInputRect) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if this.InputActive.Tags.Active {
			this.InputActive.Tags.Active = false
			this.InputActive.CurrentActive = ""
		} else {
			this.setAllInputsToInactive()
			this.InputActive.Tags.Active = true
			this.InputActive.CurrentActive = INPUTNAME_TAGS
		}
	}

	if cursorTrigger.In(fileNameInputRect) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if this.InputActive.FileName.Active {
			this.InputActive.FileName.Active = false
			this.InputActive.CurrentActive = ""
		} else {
			this.setAllInputsToInactive()
			this.InputActive.FileName.Active = true
			this.InputActive.CurrentActive = INPUTNAME_FILENAME
		}
	}

	switch this.InputActive.CurrentActive {
	case INPUTNAME_NAME:
		switch this.getAnyInputKey() {
		case "Enter":
			this.setAllInputsToInactive()
		case "Escape":
			this.setAllInputsToInactive()
			this.InputActive.Name.Input = ""
		case "Backspace":
			this.InputActive.Name.Input = this.InputActive.Name.Input[:len(this.InputActive.Name.Input)-1]
		default:
			this.InputActive.Name.Input += this.getAnyInputKey()
		}
	case INPUTNAME_TAGS:
		switch this.getAnyInputKey() {
		case "Enter":
			this.setAllInputsToInactive()
		case "Escape":
			this.setAllInputsToInactive()
			this.InputActive.Tags.Input = ""
		case "Backspace":
			this.InputActive.Tags.Input = this.InputActive.Tags.Input[:len(this.InputActive.Tags.Input)-1]
		default:
			this.InputActive.Tags.Input += this.getAnyInputKey()
		}
	case INPUTNAME_FILENAME:
		switch this.getAnyInputKey() {
		case "Enter":
			this.setAllInputsToInactive()
		case "Escape":
			this.setAllInputsToInactive()
			this.InputActive.FileName.Input = ""
		case "Backspace":
			this.InputActive.FileName.Input = this.InputActive.FileName.Input[:len(this.InputActive.FileName.Input)-1]
		default:
			this.InputActive.FileName.Input += this.getAnyInputKey()
		}
	}

	for _, btn := range this.Buttons {
		if btn.Image == UIButton {
			btnTrigger := image.Rect(int(btn.Position.X), int(btn.Position.Y), (btn.Image.Bounds().Dx()*int(btn.Scale.X))+int(btn.Position.X), (btn.Image.Bounds().Dy()*int(btn.Scale.Y))+int(btn.Position.Y))
			if cursorTrigger.In(btnTrigger) {
				btn.State = ui.BUTTON_HOVER
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
					btn.State = ui.BUTTON_PRESSED
					cancelBtnOnClickParam := make(map[string]interface{})
					cancelBtnOnClickParam["modal"] = this
					btn.OnClick(cancelBtnOnClickParam)
				}
			} else {
				btn.State = ui.BUTTON_NORMAL
			}
		}
	}
	return nil
}

func (this *ASNewRecord) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(10, 10)
	newH := UIFrame.Bounds().Max.Y * 10
	newW := UIFrame.Bounds().Max.X * 10
	op.GeoM.Translate(float64((newW / 2)), float64(newH/2))
	screen.DrawImage(UIFrame, op)

	//runtime.DrawString("New Asset Record", 1, (newW/2)+25, (newH/2)+42, false, screen)
	vector.StrokeLine(screen, float32(newW/2), float32(newH/2)+67, float32((newW/2)+newW), float32((newH/2)+67), 1, color.RGBA{0, 0, 0, 255}, false)

	//runtime.DrawString("Name    :", 1, (newW/2)+25, (newH/2)+100, false, screen)
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(10, 1)
	op.GeoM.Translate(float64(686), float64((newH/2)+95))
	screen.DrawImage(UIInput, op)
	//DrwarInputstring
	//runtime.DrawString(this.InputActive.Name.Input, 1, 686+10, ((newH/2)+95)+5, false, screen)

	//runtime.DrawString("Tags    :", 1, (newW/2)+25, (newH/2)+148, false, screen)
	op.GeoM.Translate(0, 48)
	screen.DrawImage(UIInput, op)
	//.DrawString(this.InputActive.Tags.Input, 1, 686+10, ((newH/2)+95)+5+48, false, screen)

	//runtime.DrawString("Filename:", 1, (newW/2)+25, (newH/2)+196, false, screen)
	op.GeoM.Translate(0, 48)
	screen.DrawImage(UIInput, op)
	//runtime.DrawString(this.InputActive.FileName.Input, 1, 686+10, ((newH/2)+95)+5+48+48, false, screen)

	for _, btn := range this.Buttons {
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Scale(btn.Scale.X, btn.Scale.Y)
		op.GeoM.Translate(btn.Position.X, btn.Position.Y)
		switch btn.State {
		case ui.BUTTON_NORMAL:
			screen.DrawImage(UIButton, op)
		case ui.BUTTON_HOVER:
			screen.DrawImage(UIButtonHover, op)
		case ui.BUTTON_PRESSED:
			screen.DrawImage(UIButtonPressed, op)
		default:
			screen.DrawImage(UIButton, op)
		}
		if btn.Icon != nil {
			op = &ebiten.DrawImageOptions{}
			op.GeoM.Translate(
				(btn.Position.X+float64(btn.Image.Bounds().Max.X)*btn.Scale.X)-float64(btn.Icon.Bounds().Max.X)-8,
				btn.Position.Y+8)
			screen.DrawImage(btn.Icon, op)
		}
		//runtime.DrawString(btn.Label, 1, int(btn.Position.X+5), int(btn.Position.Y)+5, false, screen)
	}

	//Debug
	if this.Debug {
		msg := fmt.Sprintf("New Record Modal: %v\n", this)
		ebitenutil.DebugPrintAt(screen, msg, 1, 32)
	}
}

func (this *ASNewRecord) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
