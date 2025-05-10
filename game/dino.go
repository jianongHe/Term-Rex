package game

import "github.com/nsf/termbox-go"

// Dino represents the player character with smooth jump physics
type Dino struct {
	X          int
	posY       float64
	velY       float64
	hangFrames int
}

// NewDino creates a new Dino at the ground position
func NewDino() *Dino {
	return &Dino{X: 2, posY: float64(height - 2)}
}

// Update advances the dino's position with smooth jump and hang time
func (d *Dino) Update() {
	// Only update when in air or hangFrames > 0
	if d.posY < float64(height-2) || d.velY != 0 || d.hangFrames > 0 {
		// calculate next velocity
		nextVel := d.velY + gravity
		// detect apex crossing and start hang
		if d.velY < 0 && nextVel >= 0 {
			d.hangFrames = hangDuration
			d.velY = 0
		} else if d.hangFrames > 0 {
			// hang time: stay at apex
			d.hangFrames--
		} else {
			// normal physics update
			d.posY += d.velY
			d.velY = nextVel
		}
		// landing on ground
		if d.posY >= float64(height-2) {
			d.posY = float64(height - 2)
			d.velY = 0
			d.hangFrames = 0
		}
	}
}

// Jump initiates an upward velocity if on the ground
func (d *Dino) Jump() {
	if d.posY == float64(height-2) {
		d.velY = jumpVelocity
		d.hangFrames = 0
	}
}

// Draw renders the dino sprite at its current position
func (d *Dino) Draw() {
	h := len(dinoSprite)
	y := int(d.posY)
	startY := y - (h - 1)
	dinoSprite.Draw(d.X, startY, termbox.ColorGreen, termbox.ColorDefault)
}
