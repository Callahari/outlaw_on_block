package player

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"image/color"
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

	DefaultSpeed          = 5.5
	DefaultBackwartsSpeed = 2.5
)

var (
	ebitenImage *ebiten.Image
)

type PlayerDirection int

// Player represents the player character in the game.
type Player struct {
	Position             struct{ X, Y float64 }
	Direction            PlayerDirection
	Rotation             int
	Show                 bool
	Speed                float64
	BackwartsSpeed       float64
	Type                 string
	UpdateDivider        int
	FirstCollisionIgnore bool
	InFrontOfCar         *car.Car
	InCar                *car.Car
	Animation            *animation.Animation
}

func init() {
	img, _, err := image.Decode(bytes.NewReader(res.PlayerSprite))
	if err != nil {
		log.Fatal(err)
	}
	origEbitenImage := ebiten.NewImageFromImage(img)
	s := origEbitenImage.Bounds().Size()
	ebitenImage = ebiten.NewImage(s.X, s.Y)
	op := &ebiten.DrawImageOptions{}
	ebitenImage.DrawImage(origEbitenImage, op)
}

func NewPlayer() *Player {
	p := &Player{
		Direction:     PlayerDown,
		Show:          true,
		Position:      struct{ X, Y float64 }{100, 100},
		Speed:         5.5,
		Type:          "player",
		UpdateDivider: 5,
	}
	s := []*ebiten.Image{}
	for y := 0; y < ebitenImage.Bounds().Size().Y/16; y++ {
		newSpriteImage := ebiten.NewImage(18, 16)
		op := &ebiten.DrawImageOptions{}
		newSpriteImage.DrawImage(ebitenImage.SubImage(image.Rect(0, y*16, 18, (y+1)*16)).(*ebiten.Image), op)
		s = append(s, newSpriteImage)
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
	if p.Animation.UpdateCounter%p.UpdateDivider == 0 {
		if p.InCar != nil {
			if p.InCar.Accelerating {
				switch p.InCar.Direction {
				case car.CarForward:
					p.InCar.CurrentSpeed = p.InCar.CurrentSpeed + car.Velocity
					if p.InCar.CurrentSpeed > car.MaxSpeed {
						p.InCar.CurrentSpeed = car.MaxSpeed
					}
				case car.CarBackwards:
					p.InCar.CurrentSpeed -= car.Velocity
					if p.InCar.CurrentSpeed < -(car.MaxSpeed / 2) {
						p.InCar.CurrentSpeed = -(car.MaxSpeed / 2)
					}
				}

			} else {
				if p.InCar.CurrentSpeed > 0 {
					p.InCar.CurrentSpeed -= car.Velocity
					if p.InCar.CurrentSpeed < 0 {
						p.InCar.CurrentSpeed = 0
						p.InCar.Direction = car.CarNeutral
					}
				} else {
					p.InCar.CurrentSpeed += car.Velocity
					if p.InCar.CurrentSpeed > 0 {
						p.InCar.CurrentSpeed = 0
						p.InCar.Direction = car.CarNeutral
					}
				}
			}
		}
		oldPosition := p.Position
		p.Animation.SpriteIdx++
		switch p.Direction {
		case PlayerUp:
			if p.InCar != nil {
				p.InCar.IsMoving = true
			}
			radians := float64(p.Rotation+90) * (math.Pi / 180)
			p.Position.X += p.Speed * math.Cos(radians)
			p.Position.Y += p.Speed * math.Sin(radians)
		case PlayerDown:
			if p.InCar != nil {
				if p.InCar.CurrentSpeed > 0 {
					p.InCar.CurrentSpeed -= car.Velocity * 2
				}
			}
			radians := float64(p.Rotation+90) * (math.Pi / 180)
			p.Position.X += p.Speed * math.Cos(radians)
			p.Position.Y += p.Speed * math.Sin(radians)

		case PlayerRight:
			p.Rotation += 1
			if p.Rotation > 359 {
				p.Rotation = 0
			}
		case PlayerLeft:
			p.Rotation -= 1
			if p.Rotation > 359 {
				p.Rotation = 0
			}
		}
		if p.InCar != nil {
			p.InCar.Position = p.Position
			p.InCar.Rotation = p.Rotation
			p.Speed = p.InCar.CurrentSpeed
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
	playerRect := image.Rect(
		int(p.Position.X),
		int(p.Position.Y),
		int(p.Position.X+18),
		int(p.Position.Y+16))
	playerRectTrans := playerRect.Add(image.Point{int(-runtime.ViewPort.X), int(-runtime.ViewPort.Y)})
	for _, WorldObject := range runtime.WorldCollisionObjects {
		WCOSprite := WorldObject.GetAnimation().Sprites[WorldObject.GetAnimation().AnimationName][WorldObject.GetAnimation().SpriteIdx]
		carRect := image.Rect(
			int(WorldObject.GetPosition().X),
			int(WorldObject.GetPosition().Y),
			int(WorldObject.GetPosition().X)+WCOSprite.Bounds().Dx(),
			int(WorldObject.GetPosition().Y)+WCOSprite.Bounds().Dy())

		carRectTrans := carRect.Add(image.Point{int(-runtime.ViewPort.X), int(-runtime.ViewPort.Y)})
		rotatedRect := runtime.RotateRect(carRectTrans, float64(WorldObject.GetRotation()%360)*2*math.Pi/360)
		if playerRectTrans.Overlaps(rotatedRect) {
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
	runtime.ViewPort.X = 1920/2 - p.Position.X
	runtime.ViewPort.Y = 1080/2 - p.Position.Y
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		p.Direction = PlayerUp
		if p.InCar != nil {
			p.InCar.Accelerating = true
			if p.InCar.Direction == car.CarNeutral {
				p.InCar.Direction = car.CarForward
			} else if p.InCar.Direction == car.CarBackwards {
				p.InCar.Direction = car.CarForward
			}
		} else {
			p.Speed = DefaultSpeed
		}
		p.move()
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		p.Direction = PlayerDown
		if p.InCar != nil {
			p.InCar.Accelerating = true
			if p.InCar.Direction == car.CarNeutral {
				p.InCar.Direction = car.CarBackwards
			} else if p.InCar.Direction == car.CarForward {
				p.InCar.Direction = car.CarBackwards
			}
		} else {
			p.Speed = -DefaultBackwartsSpeed
		}
		p.move()
	} else {
		p.Animation.UpdateCounter = 0
		if p.InCar != nil {
			p.InCar.Accelerating = false
			p.move()
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		if p.InCar != nil {
			if p.InCar.CurrentSpeed != 0 {
				p.Rotation += 3
			}
		} else {
			p.Rotation += 3
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		if p.InCar != nil {
			if p.InCar.CurrentSpeed != 0 {
				p.Rotation -= 3
			}
		} else {
			p.Rotation -= 3
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		if p.InFrontOfCar != nil {
			p.Show = false
			p.InCar = p.InFrontOfCar
			p.InFrontOfCar = nil
			p.Rotation = p.InCar.Rotation
			p.Speed = p.InCar.CurrentSpeed
			p.BackwartsSpeed = p.InCar.BackwardsSpeed
			p.UpdateDivider = 1
			log.Println("enter the car")
		} else {
			if p.InCar != nil {
				p.Show = true
				p.InCar = nil
				p.InFrontOfCar = p.InFrontOfCar
				p.Speed = DefaultSpeed
				p.BackwartsSpeed = DefaultBackwartsSpeed
				p.UpdateDivider = 5
				log.Println("get out of the car")
			}
		}
	}
	return nil
}

func (p *Player) Draw(screen *ebiten.Image) {
	inCarStr := "false"
	carDirection := 0
	if p.InCar != nil {
		inCarStr = "true"
		carDirection = int(p.InCar.Direction)
	}
	msg := fmt.Sprintf(`
inCar: %s
Speed: %f
car.Direction: %d
`, inCarStr, p.Speed, carDirection)
	ebitenutil.DebugPrintAt(screen, msg, 0, 30)
	if !p.Show {
		return
	}

	currentSprite := p.Animation.Sprites[p.Animation.AnimationName][p.Animation.SpriteIdx]
	w, h := float64(currentSprite.Bounds().Size().X), float64(currentSprite.Bounds().Size().Y)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Rotate(float64(p.Rotation%360) * 2 * math.Pi / 360)
	op.GeoM.Translate(w/2, h/2)
	//op.GeoM.Translate(p.Position.X-cameraPos.X, p.Position.Y-cameraPos.Y)
	op.GeoM.Translate(1920/2, 1080/2)
	screen.DrawImage(currentSprite, op)

	////////////////DEBUG SECTION
	for _, WorldObject := range runtime.WorldCollisionObjects {
		WCOSprite := WorldObject.GetAnimation().Sprites[WorldObject.GetAnimation().AnimationName][WorldObject.GetAnimation().SpriteIdx]
		WCORect := image.Rect(
			int(WorldObject.GetPosition().X)-int(runtime.ViewPort.X),
			int(WorldObject.GetPosition().Y)-int(runtime.ViewPort.Y),
			int(WorldObject.GetPosition().X)+WCOSprite.Bounds().Dx()-int(runtime.ViewPort.X),
			int(WorldObject.GetPosition().Y)+WCOSprite.Bounds().Dy()-int(runtime.ViewPort.Y),
		)
		//rotatedRect := runtime.RotateRect(WCORect, float64(WorldObject.GetRotation()%360)*2*math.Pi/360)
		ni := ebiten.NewImage(WCOSprite.Bounds().Dx(), WCOSprite.Bounds().Dy())
		ni.Fill(color.RGBA{255, 0, 255, 128})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(100, 100)
		screen.DrawImage(ni.SubImage(WCORect).(*ebiten.Image), op)
	}

}

func (g *Player) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
