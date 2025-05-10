package game

import (
	"github.com/nsf/termbox-go"
	"math"
	"math/rand"
)

// ObstacleType represents the type of obstacle
type ObstacleType int

const (
	Cactus ObstacleType = iota
	Bird
	BigBird
)

// Obstacle moves across the screen
type Obstacle struct {
	posX         float64
	Y            int
	obstacleType ObstacleType
	animFrame    int
	animCounter  int
}

func NewObstacle() *Obstacle {
	o := &Obstacle{}
	o.reset()
	o.animFrame = 0
	o.animCounter = 0
	return o
}

func (o *Obstacle) reset() {
	o.posX = float64(width - 1)

	// Determine obstacle type based on probabilities
	r := rand.Float64()
	if r < bigBirdProbability {
		o.obstacleType = BigBird
		// place big bird at configurable flight height
		o.Y = bigBirdFlightRow
	} else if r < bigBirdProbability+birdProbability {
		o.obstacleType = Bird
		// place bird at configurable flight height
		o.Y = birdFlightRow
	} else {
		o.obstacleType = Cactus
		o.Y = height - 2
	}

	o.animFrame = 0
	o.animCounter = 0
}

func (o *Obstacle) Update() {
	o.posX -= obstacleSpeed
	if o.posX < 0 {
		o.reset()
	}
	o.updateAnimation()
}

func (o *Obstacle) Draw() {
	var sprite Sprite
	var fg termbox.Attribute

	switch o.obstacleType {
	case Bird:
		sprite = birdFrames[o.animFrame]
		fg = termbox.ColorYellow
	case BigBird:
		sprite = bigBirdFrames[o.animFrame]
		fg = termbox.ColorMagenta
	default: // Cactus
		sprite = obstacleFrames[o.animFrame]
		fg = termbox.ColorRed
	}

	h := len(sprite)
	startY := o.Y - (h - 1)
	x := int(math.Round(o.posX))
	sprite.Draw(x, startY, fg, termbox.ColorDefault)
}

// updateAnimation advances obstacle animation frames
func (o *Obstacle) updateAnimation() {
	o.animCounter++
	if o.animCounter >= animPeriod {
		o.animCounter = 0
		var frameCount int

		switch o.obstacleType {
		case Bird:
			frameCount = len(birdFrames)
		case BigBird:
			frameCount = len(bigBirdFrames)
		default: // Cactus
			frameCount = len(obstacleFrames)
		}

		o.animFrame = (o.animFrame + 1) % frameCount
	}
}
