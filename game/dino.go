package game

import "github.com/nsf/termbox-go"

// Dino represents the player character with smooth jump physics
type Dino struct {
	X                int
	posY             float64
	velY             float64
	hangFrames       int
	animFrame        int
	animCounter      int
	duckFrames       int
	isFastDropping   bool // 标记是否正在快速下降
	isDownKeyPressed bool // 新增：标记下键是否被按住
}

// NewDino creates a new Dino at the ground position
func NewDino() *Dino {
	return &Dino{
		X:                2,
		posY:             float64(height - 2),
		isDownKeyPressed: false,
		isFastDropping:   false,
	}
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

	// 如果恐龙在地面上且下键被按住，持续刷新蹲下状态
	// 注意：这里不需要检查d.isFastDropping，因为落地时已经处理了
	if d.posY == float64(height-2) && d.isDownKeyPressed {
		d.Duck()
	}
}

// Jump initiates an upward velocity if on the ground
func (d *Dino) Jump() {
	if d.posY == float64(height-2) {
		d.velY = jumpVelocity
		d.hangFrames = 0
		d.isFastDropping = false

		// 播放跳跃音效
		GetAudioManager().PlaySound(SoundJump)
	}
}

// FastDrop initiates a fast downward velocity if in the air
func (d *Dino) FastDrop() {
	// 只有在空中才能快速下降
	if d.posY < float64(height-2) {
		// 设置一个较大的向下速度，比重力加速度更快
		d.velY = -jumpVelocity * 0.8 // 使用跳跃速度的80%作为下降速度
		d.hangFrames = 0             // 取消任何悬停时间
		d.isFastDropping = true
		d.isDownKeyPressed = true // 确保下键状态被设置为按住

		// 播放快速下降音效
		GetAudioManager().PlaySound(SoundDrop) // 使用专门的下降音效
	}
}

// Draw renders the dino sprite at its current position with animation
func (d *Dino) Draw() {
	var sprite Sprite
	onGround := int(d.posY) == height-2
	if !onGround {
		// 如果在快速下降，可以使用不同的精灵图（可选）
		if d.isFastDropping {
			sprite = dinoDuckFrames[0] // 使用蹲下的精灵图表示快速下降
		} else {
			sprite = dinoStandFrames[0]
		}
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

	// 如果正在快速下降，使用更大的重力加速度
	if d.isFastDropping {
		nextVel = d.velY + gravity*1.5
	}

	if d.isApex(nextVel) && !d.isFastDropping {
		d.startHang()
	} else if d.isHanging() {
		d.hangFrames--
	} else {
		// 使用更小的增量来实现更平滑的移动
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

		// 如果是从快速下降状态落地，立即进入蹲下状态
		if d.isFastDropping {
			d.Duck() // 立即蹲下
			// 只有在落地后才重置快速下降状态
			d.isFastDropping = false
		} else if d.isDownKeyPressed {
			// 如果不是快速下降但下键被按住，也立即蹲下
			d.Duck()
		}
	}
}

// Duck initiates ducking state for the specified duration
func (d *Dino) Duck() {
	d.duckFrames = duckHoldDuration
}

// IsDucking returns true if the dino is currently in ducking state
func (d *Dino) IsDucking() bool {
	return d.duckFrames > 0
}

// GetY returns the current Y position of the dino
func (d *Dino) GetY() int {
	return int(d.posY)
}
