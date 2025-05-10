package game

import "github.com/nsf/termbox-go"

// Dino represents the player character with smooth jump physics
type Dino struct {
	X           int
	posY        float64
	velY        float64
	hangFrames  int
	duckFrames  int
	drawingDuck bool
}

// NewDino creates a new Dino at the ground position
func NewDino() *Dino {
	return &Dino{X: 2, posY: float64(height - 2)}
}

// Update advances the dino's position with smooth jump and hang time
func (d *Dino) Update() {
	// handle duck hold
	if d.duckFrames > 0 {
		d.duckFrames--
	}
	d.drawingDuck = d.duckFrames > 0

	if d.shouldUpdate() {
		d.applyPhysics()
		d.checkLanding()
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
// Draw renders the dino sprite at its current position
func (d *Dino) Draw() {
	var sprite Sprite
	// if airborne, always use standing sprite
	if int(d.posY) != height-2 {
		sprite = dinoSprite
	} else if d.drawingDuck {
		sprite = dinoDuckSprite
	} else {
		sprite = dinoSprite
	}
	h := len(sprite)
	y := int(d.posY)
	startY := y - (h - 1)
	sprite.Draw(d.X, startY, termbox.ColorGreen, termbox.ColorDefault)
}

// shouldUpdate returns true if the dino is airborne or hanging
func (d *Dino) shouldUpdate() bool {
	return d.posY < float64(height-2) || d.velY != 0 || d.hangFrames > 0
}

// applyPhysics updates velocity and position, handling apex hang
func (d *Dino) applyPhysics() {
	nextVel := d.velY + gravity
	if d.isApex(nextVel) {
		d.startHang()
	} else if d.isHanging() {
		d.hangFrames--
	} else {
		d.posY += d.velY
		d.velY = nextVel
	}
}

// isApex returns true when dino transitions from rising to falling
func (d *Dino) isApex(nextVel float64) bool {
	return d.velY < 0 && nextVel >= 0
}

// startHang sets up hang frames at apex
func (d *Dino) startHang() {
	d.hangFrames = hangDuration
	d.velY = 0
}

// isHanging returns true if the dino is in hang time
func (d *Dino) isHanging() bool {
	return d.hangFrames > 0
}

// checkLanding resets the dino when landing on ground
func (d *Dino) checkLanding() {
	if d.posY >= float64(height-2) {
		d.posY = float64(height - 2)
		d.velY = 0
		d.hangFrames = 0
	}
}
