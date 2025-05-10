package game

import "math"

// checkCollision returns true if the dino and obstacle sprites overlap on non-space chars.
func (g *Game) checkCollision() bool {
	// Determine Dino sprite and bounds
	var dSprite Sprite
	onGround := int(g.dino.posY) == height-2
	if !onGround {
		dSprite = dinoStandFrames[0]
	} else if g.dino.duckFrames > 0 {
		dSprite = dinoDuckFrames[g.dino.animFrame]
	} else {
		dSprite = dinoStandFrames[g.dino.animFrame]
	}
	dW := len(dSprite[0])
	dH := len(dSprite)
	dX0 := g.dino.X
	dY0 := int(g.dino.posY) - (dH - 1)

	// Determine Obstacle sprite and bounds
	var oSprite Sprite
	if g.obstacle.isBird {
		oSprite = birdFrames[g.obstacle.animFrame]
	} else {
		oSprite = obstacleFrames[g.obstacle.animFrame]
	}
	oW := len(oSprite[0])
	oH := len(oSprite)
	oX0 := int(math.Round(g.obstacle.posX))
	oY0 := g.obstacle.Y - (oH - 1)

	// Compute overlap rectangle
	xStart := max(dX0, oX0)
	xEnd := min(dX0+dW-1, oX0+oW-1)
	yStart := max(dY0, oY0)
	yEnd := min(dY0+dH-1, oY0+oH-1)

	// Check each overlapping cell: collision only if both chars are non-space
	for y := yStart; y <= yEnd; y++ {
		for x := xStart; x <= xEnd; x++ {
			dx := x - dX0
			dy := y - dY0
			ox := x - oX0
			oy := y - oY0
			if dSprite[dy][dx] != ' ' && oSprite[oy][ox] != ' ' {
				return true
			}
		}
	}
	return false
}

// helper min and max
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
