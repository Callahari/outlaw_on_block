package main

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"os"
	"outlaw_on_block/car"
	"outlaw_on_block/game"
	game2 "outlaw_on_block/game"
	"outlaw_on_block/player"
	"outlaw_on_block/res"
	"outlaw_on_block/runtime"
	"outlaw_on_block/tiles"
)

const (
	screenWidth  = 1920
	screenHeight = 1080
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Outlaw on Block")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

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
	game.Scene = game2.GameScene_Menu
	game.Editor = nil

	//set UserDirectory
	runtime.OOBUserDir, _ = os.UserHomeDir()
	//
	os.MkdirAll(runtime.OOBUserDir+"/OOB/images", 0744)

	//Check if assert.db exist , if not create
	_, err := os.Stat(runtime.OOBUserDir + "/OOB/assets.db")
	if os.IsNotExist(err) {
		file, err := os.Create(runtime.OOBUserDir + "/OOB/assets.db")
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		file.Close()
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
