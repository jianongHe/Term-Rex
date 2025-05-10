package game

import (
	"math"
)

// checkCollision detects if the dino has collided with an obstacle
func (g *Game) checkCollision() bool {
	// 检查与所有障碍物的碰撞
	for _, obstacle := range g.obstacleManager.GetObstacles() {
		if checkSingleCollision(g.dino, obstacle) {
			return true
		}
	}
	return false
}

// checkSingleCollision checks collision between dino and a single obstacle
func checkSingleCollision(dino *Dino, obstacle IObstacle) bool {
	// get dino sprite based on state
	var dinoSprite Sprite
	if dino.IsDucking() {
		dinoSprite = dinoDuckFrames[dino.animFrame]
	} else {
		dinoSprite = dinoStandFrames[dino.animFrame]
	}

	// get obstacle sprite
	obstacleSprite := obstacle.GetSprite()

	// get positions
	dinoX := dino.X
	dinoY := dino.GetY() - len(dinoSprite) + 1 // adjust for sprite height
	obstacleX, obstacleY := obstacle.GetPosition()
	obstacleY = obstacleY - len(obstacleSprite) + 1 // adjust for sprite height

	// convert to integer position
	obstacleXInt := int(math.Round(obstacleX))

	// check for overlap in x and y dimensions
	dinoWidth := getMaxWidth(dinoSprite)
	obstacleWidth := getMaxWidth(obstacleSprite)
	dinoHeight := len(dinoSprite)
	obstacleHeight := len(obstacleSprite)

	// check for overlap in x dimension
	if dinoX+dinoWidth <= obstacleXInt || dinoX >= obstacleXInt+obstacleWidth {
		return false
	}

	// check for overlap in y dimension
	if dinoY+dinoHeight <= obstacleY || dinoY >= obstacleY+obstacleHeight {
		return false
	}

	// check for character-level collision
	for dy := 0; dy < dinoHeight; dy++ {
		for dx := 0; dx < dinoWidth; dx++ {
			dinoRow := dy
			dinoCol := dx
			if dinoRow >= len(dinoSprite) || dinoCol >= len(dinoSprite[dinoRow]) {
				continue
			}

			// 确保所有类型都转换为int，避免类型不匹配错误
			obstacleRow := dy + dinoY - obstacleY
			obstacleCol := dx + dinoX - int(math.Round(obstacleX))

			if obstacleRow < 0 || obstacleRow >= len(obstacleSprite) ||
				obstacleCol < 0 || obstacleCol >= len(obstacleSprite[obstacleRow]) {
				continue
			}

			dinoChar := dinoSprite[dinoRow][dinoCol]
			obstacleChar := obstacleSprite[obstacleRow][obstacleCol]
			if dinoChar != ' ' && obstacleChar != ' ' {
				return true
			}
		}
	}

	return false
}

// getMaxWidth returns the maximum width of a sprite
func getMaxWidth(s Sprite) int {
	maxWidth := 0
	for _, row := range s {
		if len(row) > maxWidth {
			maxWidth = len(row)
		}
	}
	return maxWidth
}
