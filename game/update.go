package game

import (
	"os"

	"github.com/nsf/termbox-go"
	"time"
)

// update updates game state
func (g *Game) update() {
	g.dino.Update()
	if g.started {
		g.applyStage()
		g.obstacle.Update()
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

// gameOver displays game over screen and waits for restart or quit
func (g *Game) gameOver() {
	PrintCenter("GAME OVER (r to retry, q to quit)")
	termbox.Flush()
	for {
		ev := <-g.events
		if ev.Type == termbox.EventKey {
			if ev.Ch == KeyRestartRune {
				// update highest score
				if g.score > g.highestScore {
					g.highestScore = g.score
				}
				// reset game state
				g.dino = NewDino()
				g.obstacle = NewObstacle()
				g.score = 0
				// reset stage progression and parameters
				g.stageIndexActive = 0
				g.stageIndexTarget = 0
				g.stageTransitionStart = time.Time{}
				obstacleSpeed = stageConfigs[0].Speed
				birdProbability = stageConfigs[0].BirdProb
				bigBirdProbability = stageConfigs[0].BigBirdProb
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
		g.groundStart -= groundExtendSpeed
		if g.groundStart < 0 {
			g.groundStart = 0
		}
	}
	if g.groundEnd < width-1 {
		g.groundEnd += groundExtendSpeed
		if g.groundEnd > width-1 {
			g.groundEnd = width - 1
		}
	}
}
