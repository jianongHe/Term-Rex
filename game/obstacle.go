package game

import "github.com/nsf/termbox-go"

// Obstacle moves across the screen
type Obstacle struct {
	X, Y int
}

func NewObstacle() *Obstacle {
	return &Obstacle{X: width - 1, Y: height - 2}
}

func (o *Obstacle) Update() {
	o.X--
	if o.X < 0 {
		o.X = width - 1
	}
}

var obstacleSprite = Sprite{
	" | ",
	"/|\\",
	" | ",
}

func (o *Obstacle) Draw() {
	h := len(obstacleSprite)
	startY := o.Y - (h - 1)
	obstacleSprite.Draw(o.X, startY, termbox.ColorRed, termbox.ColorDefault)
}
