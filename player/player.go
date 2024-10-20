package player

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"log"
	"math"
	"outlaw_on_block/animation"
	"outlaw_on_block/car"
	"outlaw_on_block/res"
	"outlaw_on_block/runtime"
)

const (
	PlayerDown = iota
	PlayerUp
	PlayerRight
	PlayerLeft
)

var (
	ebitenImage *ebiten.Image
)

type PlayerDirection int

type Player struct {
	Position             struct{ X, Y float64 }
	Direction            PlayerDirection
	Show                 bool
	Speed                float64
	Type                 string
	FirstCollisionIgnore bool
	InFrontOfCar         *car.Car
	InCar                *car.Car
	Animation            *animation.Animation
}

func init() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(res.PlayerSprite))
	if err != nil {
		log.Fatal(err)
	}
	origEbitenImage := ebiten.NewImageFromImage(img)

	s := origEbitenImage.Bounds().Size()
	ebitenImage = ebiten.NewImage(s.X, s.Y)

	op := &ebiten.DrawImageOptions{}
	//op.ColorScale.ScaleAlpha(0.5)
	ebitenImage.DrawImage(origEbitenImage, op)
}

func NewPlayer() *Player {
	p := &Player{
		Direction: PlayerDown,
		Show:      true,
		Position:  struct{ X, Y float64 }{100, 100},
		Speed:     5.5,
		Type:      "player",
	}
	s := []*ebiten.Image{}
	for y := range ebitenImage.Bounds().Size().Y / 16 {
		//for y := range 2 {

		newSpriteImage := ebiten.NewImage(18, 16)
		op := &ebiten.DrawImageOptions{}
		newSpriteImage.DrawImage(ebitenImage.SubImage(image.Rect(0, y*16, 18, (y+1)*16)).(*ebiten.Image), op)
		//newSpriteImage.DrawImage(ebitenImage, op)
		s = append(s, newSpriteImage)

		log.Printf("y: %d, bounds: %v,offsetY: %f", y, newSpriteImage.Bounds(), float64(y)*16)
		//s = append(s, ebitenImage)
	}
	p.Animation = &animation.Animation{}
	p.Animation.AnimationName = "walk"
	p.Animation.Sprites = map[string][]*ebiten.Image{"walk": s}

	return p
}
func (p *Player) GetType() string {
	return p.Type
}
func (p *Player) GetAnimation() *animation.Animation {
	return p.Animation
}
func (p *Player) GetPosition() struct{ X, Y float64 } {
	return p.Position
}
func (p *Player) move() {
	p.Animation.UpdateCounter++
	if p.Animation.UpdateCounter%8 == 0 {
		oldPosition := p.Position
		p.Animation.SpriteIdx++
		switch p.Direction {
		case 0:
			{
				p.Position.Y += p.Speed
			}
		case 1:
			{
				p.Position.Y -= p.Speed
			}
		case 2:
			{
				p.Position.X += p.Speed
			}
		case 3:
			{
				p.Position.X -= p.Speed
			}
		}
		if p.InCar != nil {
			p.InCar.Position = p.Position
		}
		if p.detectCollision() {
			p.Position = oldPosition
		}
		if p.Animation.SpriteIdx >= len(p.Animation.Sprites[p.Animation.AnimationName])-1 {
			p.Animation.SpriteIdx = 0
		}
	}
}
func (p *Player) detectCollision() bool {
	if !p.Show {
		return false
	}
	currentPosBounds := image.Rect(
		int(p.Position.X+(p.Speed/2)),
		int(p.Position.Y+(p.Speed/2)),
		p.GetAnimation().Sprites[p.Animation.AnimationName][p.Animation.SpriteIdx].Bounds().Dx()+int(p.Position.X+(p.Speed/2)),
		p.GetAnimation().Sprites[p.Animation.AnimationName][p.Animation.SpriteIdx].Bounds().Dy()+int(p.Position.Y+(p.Speed/2)))
	for _, WorldObject := range runtime.WorldCollisionObjects {

		WCOSprite := WorldObject.GetAnimation().Sprites[WorldObject.GetAnimation().AnimationName][WorldObject.GetAnimation().SpriteIdx]
		WCORect := image.Rect(
			int(WorldObject.GetPosition().X),
			int(WorldObject.GetPosition().Y),
			int(int(WorldObject.GetPosition().X)+WCOSprite.Bounds().Dx()),
			int(int(WorldObject.GetPosition().Y)+WCOSprite.Bounds().Dy()),
		)

		rotatedRect := runtime.RotateRect(WCORect, float64(90%360)*2*math.Pi/360)

		if currentPosBounds.Overlaps(rotatedRect) {
			log.Printf("collision detected")
			log.Printf("currentPosBounds: %v", WCOSprite.Bounds())
			if WorldObject.GetType() == "car" {
				p.InFrontOfCar = WorldObject.(*car.Car)
			}
			return true
		}
	}
	p.InFrontOfCar = nil
	return false
}

func (p *Player) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		p.Direction = PlayerUp
		p.move()
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		p.Direction = PlayerDown
		p.move()
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		p.Direction = PlayerRight
		p.move()
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		p.Direction = PlayerLeft
		p.move()
	} else {
		p.Animation.UpdateCounter = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		if p.InFrontOfCar != nil {
			log.Println("enter the car")
			p.Show = false
			p.InCar = p.InFrontOfCar
			p.InFrontOfCar = nil
		} else {
			if p.InCar != nil {
				log.Println("Get out of the car")
				p.Show = true
				p.InCar = nil
				p.InFrontOfCar = p.InFrontOfCar
			} else {
				log.Println("Not in the front of a car")
			}

		}
	}
	return nil
}

func (p *Player) Draw(screen *ebiten.Image) {
	if !p.Show {
		return
	}
	currentSprite := p.Animation.Sprites[p.Animation.AnimationName][p.Animation.SpriteIdx]
	w, h := float64(currentSprite.Bounds().Size().X), float64(currentSprite.Bounds().Size().Y)
	op := &ebiten.DrawImageOptions{}
	switch p.Direction {
	case 0:
		{
		}
	case 1:
		{
			op.GeoM.Translate(-w/2, -h/2)
			op.GeoM.Rotate(math.Pi)
			op.GeoM.Translate(w/2, h/2)
		}
	case 2:
		{
			op.GeoM.Translate(-w/2, -h/2)
			op.GeoM.Rotate(float64(-90%360) * 2 * math.Pi / 360)
			op.GeoM.Translate(w/2, h/2)
		}
	case 3:
		{
			op.GeoM.Translate(-w/2, -h/2)
			op.GeoM.Rotate(float64(90%360) * 2 * math.Pi / 360)
			op.GeoM.Translate(w/2, h/2)
		}

	}
	op.GeoM.Translate(float64(p.Position.X), float64(p.Position.Y))
	screen.DrawImage(currentSprite, op)
}

func (g *Player) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
