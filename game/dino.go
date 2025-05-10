package game

import "github.com/nsf/termbox-go"

// Dino represents the player character with smooth jump physics
type Dino struct {
	X           int
	posY        float64
	velY        float64
	hangFrames  int
	animFrame   int
	animCounter int
	duckFrames  int
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

	if d.shouldUpdate() {
		d.applyPhysics()
		d.checkLanding()
	}
	d.updateAnimation()
}

// Jump initiates an upward velocity if on the ground
func (d *Dino) Jump() {
	if d.posY == float64(height-2) {
		d.velY = jumpVelocity
		d.hangFrames = 0
	}
}

// Draw renders the dino sprite at its current position
// Draw renders the dino sprite at its current position with animation
func (d *Dino) Draw() {
	var sprite Sprite
	onGround := int(d.posY) == height-2
	if !onGround {
		sprite = dinoStandFrames[0]
	} else if d.duckFrames > 0 {
		sprite = dinoDuckFrames[d.animFrame]
	} else {
		sprite = dinoStandFrames[d.animFrame]
	}
	h := len(sprite)
	y := int(d.posY)
	startY := y - (h - 1)
	sprite.Draw(d.X, startY, termbox.ColorGreen, termbox.ColorDefault)
}

// updateAnimation advances animation frames
func (d *Dino) updateAnimation() {
	d.animCounter++
	if d.animCounter >= animPeriod {
		d.animCounter = 0
		d.animFrame = (d.animFrame + 1) % 2
	}
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
