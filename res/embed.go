package res

import (
	_ "embed"
)

var (
	//go:embed player/player.png
	PlayerSprite []byte

	//////TILES////
	//go:embed tiles/GTA2_TILE_0.png
	Tile0Sprite []byte
	//go:embed tiles/GTA2_TILE_276.png
	Tile276Sprite []byte

	/////CARS
	//go:embed car/GTA2_CAR_50.png
	Car50Sprite []byte

	//FONT
	//go:embed font/m6x11plus.ttf
	OOBFont []byte
	//go:embed font/font_hover.png
	OOBFontHover []byte

	//UI elements
	//go:embed ui/arrows.png
	UIArrows []byte
	//go:embed ui/flat/UI_Flat_Button01a_1.png
	UIFlatButton01a1 []byte
	//go:embed ui/flat/UI_Flat_Button01a_2.png
	UIFlatButton01a2 []byte
	//go:embed ui/flat/UI_Flat_Button01a_3.png
	UIFlatButton01a3 []byte
	//go:embed ui/flat/UI_Flat_InputField01a.png
	UIFlatInputField01a []byte
	//go:embed ui/flat/UI_Flat_Frame01a.png
	UIFlatFrame01a []byte
	//go:embed ui/flat/UI_Flat_ButtonCross01a.png
	UIFlatButtonCross01a []byte
	//go:embed ui/flat/UI_Flat_ButtonCheck01a.png
	UIFlatButtonCheck01a []byte
	//go:embed ui/flat/UI_Flat_Frame03a.png
	UIFlatFrame03a []byte
)
