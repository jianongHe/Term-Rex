package game

import (
	"github.com/nsf/termbox-go"
	"math/rand"
)

// Obstacle moves across the screen
type Obstacle struct {
	X, Y        int
	isBird      bool
	animFrame   int
	animCounter int
}

func NewObstacle() *Obstacle {
	o := &Obstacle{}
	o.reset()
	o.animFrame = 0
	o.animCounter = 0
	return o
}

func (o *Obstacle) reset() {
	o.X = width - 1
	if rand.Float64() < birdProbability {
		o.isBird = true
		// place bird at configurable flight height
		o.Y = birdFlightRow
	} else {
		o.isBird = false
		o.Y = height - 2
	}
	o.animFrame = 0
	o.animCounter = 0
}

func (o *Obstacle) Update() {
	o.X -= obstacleSpeed
	if o.X < 0 {
		o.reset()
	}
	o.updateAnimation()
}

func (o *Obstacle) Draw() {
	var sprite Sprite
	var fg termbox.Attribute
	if o.isBird {
		sprite = birdFrames[o.animFrame]
		fg = termbox.ColorYellow
	} else {
		sprite = obstacleFrames[o.animFrame]
		fg = termbox.ColorRed
	}
	h := len(sprite)
	startY := o.Y - (h - 1)
	sprite.Draw(o.X, startY, fg, termbox.ColorDefault)
}

// updateAnimation advances obstacle animation frames
func (o *Obstacle) updateAnimation() {
	o.animCounter++
	if o.animCounter >= animPeriod {
		o.animCounter = 0
		var frameCount int
		if o.isBird {
			frameCount = len(birdFrames)
		} else {
			frameCount = len(obstacleFrames)
		}
		o.animFrame = (o.animFrame + 1) % frameCount
	}
}
