package ui

import "github.com/hajimehoshi/ebiten/v2"

const (
	BUTTON_NORMAL = iota
	BUTTON_HOVER
	BUTTON_PRESSED
)

type ButtonState int

type Button struct {
	Image    *ebiten.Image
	Icon     *ebiten.Image
	State    ButtonState
	Label    string
	Position struct {
		X float64
		Y float64
	}
	Scale struct {
		X float64
		Y float64
	}
	OnClick func(map[string]interface{})
}
