package main

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"outlaw_on_block/car"
	"outlaw_on_block/game"
	"outlaw_on_block/player"
	"outlaw_on_block/res"
	"outlaw_on_block/tiles"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Sprites (Ebitengine Demo)")
	ebiten.SetWindowResizable(true)

	timg, _, _ := image.Decode(bytes.NewReader(res.Tile0Sprite))
	timg1, _, _ := image.Decode(bytes.NewReader(res.Tile276Sprite))

	tMap := []*tiles.Tile{}
	t := &tiles.Tile{
		TileImage:    ebiten.NewImageFromImage(timg),
		Pos:          struct{ X, Y int }{100, 100},
		Rotation:     0,
		CollisionMap: nil,
	}
	tMap = append(tMap, t)
	t = &tiles.Tile{
		TileImage:    ebiten.NewImageFromImage(timg),
		Pos:          struct{ X, Y int }{100, 164},
		Rotation:     0,
		CollisionMap: nil,
	}
	tMap = append(tMap, t)

	t = &tiles.Tile{
		TileImage:    ebiten.NewImageFromImage(timg1),
		Pos:          struct{ X, Y int }{100, 228},
		Rotation:     0,
		CollisionMap: nil,
	}
	tMap = append(tMap, t)
	t = &tiles.Tile{
		TileImage:    ebiten.NewImageFromImage(timg1),
		Pos:          struct{ X, Y int }{100, 292},
		Rotation:     180,
		CollisionMap: nil,
	}
	tMap = append(tMap, t)

	game := &game.Game{}
	p := player.NewPlayer()
	game.Player = p
	game.TilesMap = tMap

	cars := []*car.Car{}
	taxiSrc, _, _ := image.Decode(bytes.NewReader(res.Car50Sprite))
	taxi := ebiten.NewImageFromImage(taxiSrc)
	c := car.NewCar(taxi)
	c.Rotation = 90
	cars = append(cars, c)

	game.Cars = cars
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
