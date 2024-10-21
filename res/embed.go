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
	//go:embed font/font.png
	OOBFont []byte
	//go:embed font/font_hover.png
	OOBFontHover []byte
)
