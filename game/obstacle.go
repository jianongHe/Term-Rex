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

func (o *Obstacle) Draw() {
	termbox.SetCell(o.X, o.Y, '|', termbox.ColorRed, termbox.ColorDefault)
}
