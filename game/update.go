package game

import (
	"os"

	"github.com/nsf/termbox-go"
	"time"
)

// update updates game state
func (g *Game) update() {
	// 如果下键被按住，确保恐龙知道这一点
	if g.downKeyHeld {
		g.dino.isDownKeyPressed = true
	}

	g.dino.Update()
	g.cloudManager.Update() // Update clouds regardless of game state

	if g.started {
		g.applyStage()
		g.obstacleManager.Update()

		// 即使地面已经完全扩展，也要更新地面装饰的位置
		g.updateGroundDecorations()

		if g.checkCollision() {
			g.collided = true
		}
		if g.groundExtending {
			g.updateGround()
			// stop extending once ground fully spans screen
			if g.groundStart == 0 && g.groundEnd == width-1 {
				g.groundExtending = false
			}
		}
	}
}

// updateGroundDecorations 更新地面装饰的位置，使其随着游戏进行而移动
func (g *Game) updateGroundDecorations() {
	if !g.collided {
		// 移动所有地面装饰，速度与障碍物相同
		for i := range groundDecorations {
			groundDecorations[i].x -= obstacleSpeed

			// 如果装饰移出了屏幕左侧，将其移到屏幕右侧重新出现
			if groundDecorations[i].x < -5 {
				groundDecorations[i].x += float64(width * 2)
			}
		}
	}
}

// gameOver displays game over screen and waits for restart or quit
func (g *Game) gameOver() {
	// 播放碰撞音效
	GetAudioManager().PlaySound(SoundCollision)

	// 更新最高分并保存
	if g.score > g.highestScore {
		g.highestScore = g.score
		// 保存最高分到文件
		err := SaveHighScore(g.highestScore)
		if err != nil {
			// 如果保存失败，只在控制台打印错误，不影响游戏
			// 这里不能直接显示在游戏界面，因为会干扰游戏结束画面
			// fmt.Println("无法保存最高分:", err)
		}
	}

	PrintCenter("GAME OVER (r to retry, q to quit)")

	// 显示音效控制提示
	soundMsg := "Press 'm' to toggle sound"
	if !GetAudioManager().IsEnabled() {
		soundMsg = "Sound OFF - Press 'm' to enable"
	}
	PrintCenterAt(soundMsg, height/2+2)

	termbox.Flush()
	for {
		ev := <-g.events
		if ev.Type == termbox.EventKey {
			if ev.Ch == KeyRestartRune {
				// reset game state
				g.dino = NewDino()
				g.obstacleManager = NewObstacleManager()
				// Don't reset clouds, just let them continue
				g.score = 0
				// reset stage progression and parameters
				g.stageIndexActive = 0
				g.stageIndexTarget = 0
				g.stageTransitionStart = time.Time{}
				obstacleSpeed = stageConfigs[0].Speed * speedFactor

				// 重置障碍物管理器
				g.obstacleManager = NewObstacleManager()
				return
			}
			if ev.Key == KeyQuit || ev.Ch == KeyQuitRune {
				termbox.Close()
				os.Exit(0)
			}
		}
	}
}

// updateGround expands the ground boundaries until filling the screen.
func (g *Game) updateGround() {
	if g.groundStart > 0 {
		// 使用向下取整的方式将浮点数转换为整数
		moveAmount := int(groundExtendSpeed)
		if moveAmount < 1 {
			moveAmount = 1 // 确保至少移动1个单位
		}
		g.groundStart -= moveAmount
		if g.groundStart < 0 {
			g.groundStart = 0
		}
	}
	if g.groundEnd < width-1 {
		// 使用向下取整的方式将浮点数转换为整数
		moveAmount := int(groundExtendSpeed)
		if moveAmount < 1 {
			moveAmount = 1 // 确保至少移动1个单位
		}
		g.groundEnd += moveAmount
		if g.groundEnd > width-1 {
			g.groundEnd = width - 1
		}
	}
}
