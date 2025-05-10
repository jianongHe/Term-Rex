package game

import (
	"github.com/nsf/termbox-go"
	"math"
	"math/rand"
)

// ObstacleType represents the type of obstacle
type ObstacleType int

const (
	CactusType ObstacleType = iota
	BirdType
	BigBirdType
)

// IObstacle defines the interface for all obstacle types
type IObstacle interface {
	Update()
	Draw()
	GetPosition() (float64, int)
	SetPosition(x float64, y int)
	Reset()
	GetSprite() Sprite
}

// BaseObstacle contains common properties and methods for all obstacles
type BaseObstacle struct {
	posX        float64
	y           int
	animFrame   int
	animCounter int
}

// Update moves the obstacle and updates animation
func (b *BaseObstacle) Update() {
	b.posX -= obstacleSpeed
	b.updateAnimation()
}

// GetPosition returns the current position of the obstacle
func (b *BaseObstacle) GetPosition() (float64, int) {
	return b.posX, b.y
}

// SetPosition sets the position of the obstacle
func (b *BaseObstacle) SetPosition(x float64, y int) {
	b.posX = x
	b.y = y
}

// updateAnimation advances obstacle animation frames
func (b *BaseObstacle) updateAnimation() {
	b.animCounter++
	if b.animCounter >= animPeriod {
		b.animCounter = 0
		b.animFrame = (b.animFrame + 1) % b.getFrameCount()
	}
}

// Cactus represents a cactus obstacle
type Cactus struct {
	BaseObstacle
}

// NewCactus creates a new cactus obstacle
func NewCactus() *Cactus {
	c := &Cactus{}
	c.Reset()
	return c
}

// Reset resets the cactus position and animation
func (c *Cactus) Reset() {
	c.posX = float64(width - 1)
	c.y = height - 2
	c.animFrame = 0
	c.animCounter = 0
}

// Draw renders the cactus on screen
func (c *Cactus) Draw() {
	sprite := obstacleFrames[c.animFrame]
	h := len(sprite)
	startY := c.y - (h - 1)
	x := int(math.Round(c.posX))
	sprite.Draw(x, startY, termbox.ColorRed, termbox.ColorDefault)
}

// GetSprite returns the current sprite for collision detection
func (c *Cactus) GetSprite() Sprite {
	return obstacleFrames[c.animFrame]
}

// getFrameCount returns the number of animation frames
func (c *BaseObstacle) getFrameCount() int {
	return len(obstacleFrames)
}

// Bird represents a small bird obstacle
type Bird struct {
	BaseObstacle
}

// NewBird creates a new bird obstacle
func NewBird() *Bird {
	b := &Bird{}
	b.Reset()
	return b
}

// Reset resets the bird position and animation
func (b *Bird) Reset() {
	b.posX = float64(width - 1)
	b.y = birdFlightRow
	b.animFrame = 0
	b.animCounter = 0
}

// Draw renders the bird on screen
func (b *Bird) Draw() {
	sprite := birdFrames[b.animFrame]
	h := len(sprite)
	startY := b.y - (h - 1)
	x := int(math.Round(b.posX))
	sprite.Draw(x, startY, termbox.ColorYellow, termbox.ColorDefault)
}

// GetSprite returns the current sprite for collision detection
func (b *Bird) GetSprite() Sprite {
	return birdFrames[b.animFrame]
}

// getFrameCount returns the number of animation frames
func (b *Bird) getFrameCount() int {
	return len(birdFrames)
}

// BigBird represents a large bird obstacle
type BigBird struct {
	BaseObstacle
}

// NewBigBird creates a new big bird obstacle
func NewBigBird() *BigBird {
	b := &BigBird{}
	b.Reset()
	return b
}

// Reset resets the big bird position and animation
func (b *BigBird) Reset() {
	b.posX = float64(width - 1)
	b.y = bigBirdFlightRow
	b.animFrame = 0
	b.animCounter = 0
}

// Draw renders the big bird on screen
func (b *BigBird) Draw() {
	sprite := bigBirdFrames[b.animFrame]
	h := len(sprite)
	startY := b.y - (h - 1)
	x := int(math.Round(b.posX))
	sprite.Draw(x, startY, termbox.ColorMagenta, termbox.ColorDefault)
}

// GetSprite returns the current sprite for collision detection
func (b *BigBird) GetSprite() Sprite {
	return bigBirdFrames[b.animFrame]
}

// getFrameCount returns the number of animation frames
func (b *BigBird) getFrameCount() int {
	return len(bigBirdFrames)
}

// ObstacleManager manages the creation and updating of obstacles
type ObstacleManager struct {
	currentObstacle IObstacle
}

// NewObstacleManager creates a new obstacle manager
func NewObstacleManager() *ObstacleManager {
	om := &ObstacleManager{}
	om.generateNewObstacle()
	return om
}

// Update updates the current obstacle and generates a new one if needed
func (om *ObstacleManager) Update() {
	om.currentObstacle.Update()
	x, _ := om.currentObstacle.GetPosition()
	if x < 0 {
		om.generateNewObstacle()
	}
}

// Draw renders the current obstacle
func (om *ObstacleManager) Draw() {
	om.currentObstacle.Draw()
}

// GetCurrentObstacle returns the current active obstacle
func (om *ObstacleManager) GetCurrentObstacle() IObstacle {
	return om.currentObstacle
}

// generateNewObstacle creates a new obstacle based on probabilities
func (om *ObstacleManager) generateNewObstacle() {
	r := rand.Float64()
	if r < bigBirdProbability {
		om.currentObstacle = NewBigBird()
	} else if r < bigBirdProbability+birdProbability {
		om.currentObstacle = NewBird()
	} else {
		om.currentObstacle = NewCactus()
	}
}
