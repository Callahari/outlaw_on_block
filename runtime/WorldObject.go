package runtime

import "outlaw_on_block/animation"

type WorldObject interface {
	GetType() string
	GetAnimation() *animation.Animation
	GetPosition() struct{ X, Y float64 }
}
