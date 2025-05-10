package game

import (
	"github.com/nsf/termbox-go"
	"math/rand"
)

// Obstacle moves across the screen
type Obstacle struct {
	X, Y   int
	isBird bool
}

func NewObstacle() *Obstacle {
	o := &Obstacle{}
	o.reset()
	return o
}

func (o *Obstacle) reset() {
	o.X = width - 1
	if rand.Float64() < birdProbability {
		o.isBird = true
		// place bird just above cactus level
		o.Y = height - 3
	} else {
		o.isBird = false
		o.Y = height - 2
	}
}

func (o *Obstacle) Update() {
	o.X -= obstacleSpeed
	if o.X < 0 {
		o.reset()
	}
}

func (o *Obstacle) Draw() {
	var sprite Sprite
	var fg termbox.Attribute
	if o.isBird {
		sprite = birdSprite
		fg = termbox.ColorYellow
	} else {
		sprite = obstacleSprite
		fg = termbox.ColorRed
	}
	h := len(sprite)
	startY := o.Y - (h - 1)
	sprite.Draw(o.X, startY, fg, termbox.ColorDefault)
}
