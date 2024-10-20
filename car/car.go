package car

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"outlaw_on_block/animation"
	"outlaw_on_block/runtime"
)

const (
	CarDown = iota
	CarUp
	CarLeft
	CarRight
)

type CarDirection int

type Car struct {
	Type      string
	Position  struct{ X, Y float64 }
	Direction CarDirection
	Animation *animation.Animation
}

func NewCar(spriteImage *ebiten.Image) *Car {
	c := &Car{}
	c.Type = "car"
	c.Position.X = 150
	c.Position.Y = 150
	c.Direction = CarLeft
	c.Animation = &animation.Animation{}
	c.Animation.AnimationName = "normal"
	c.Animation.Sprites = make(map[string][]*ebiten.Image)
	c.Animation.Sprites["normal"] = []*ebiten.Image{spriteImage}

	runtime.WorldCollisionObjects = append(runtime.WorldCollisionObjects, c)
	return c

}
func (c *Car) GetType() string {
	return c.Type
}
func (c *Car) GetAnimation() *animation.Animation {
	return c.Animation
}
func (c *Car) GetPosition() struct{ X, Y float64 } {
	return c.Position
}
func (c *Car) Update() error {
	return nil
}

func (c *Car) Draw(screen *ebiten.Image) {
	currentSprite := c.Animation.Sprites[c.Animation.AnimationName][c.Animation.SpriteIdx]
	w, h := float64(currentSprite.Bounds().Dx()), float64(currentSprite.Bounds().Dy())
	op := &ebiten.DrawImageOptions{}
	switch c.Direction {
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
	op.GeoM.Translate(float64(c.Position.X), float64(c.Position.Y))
	screen.DrawImage(currentSprite, op)
}

func (c *Car) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
