package game

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"math"
	"os"
	"time"
)

// applyStage smoothly transitions parameters based on score threshold crossings.
func (g *Game) applyStage() {
	// determine target stage for current score
	target := 0
	for i := len(stageConfigs) - 1; i >= 0; i-- {
		if g.score >= stageConfigs[i].ScoreThreshold {
			target = i
			break
		}
	}
	// on first crossing, start transition
	if target != g.stageIndexTarget {
		g.stageIndexTarget = target
		g.stageTransitionStart = time.Now()
	}
	// if currently transitioning between two stages
	if g.stageIndexActive != g.stageIndexTarget {
		elapsed := time.Since(g.stageTransitionStart)
		frac := float64(elapsed) / float64(stageTransitionDuration)
		if frac >= 1 {
			// finish transition
			g.stageIndexActive = g.stageIndexTarget
			obstacleSpeed = stageConfigs[g.stageIndexActive].Speed
			birdProbability = stageConfigs[g.stageIndexActive].BirdProb
			g.stageTransitionStart = time.Time{}
		} else {
			// interpolate between active and target
			old := stageConfigs[g.stageIndexActive]
			next := stageConfigs[g.stageIndexTarget]
			speed := old.Speed + frac*(next.Speed-old.Speed)
			obstacleSpeed = speed
			birdProbability = old.BirdProb + frac*(next.BirdProb-old.BirdProb)
		}
	} else {
		// no transition: keep active stage values
		sc := stageConfigs[g.stageIndexActive]
		obstacleSpeed = sc.Speed
		birdProbability = sc.BirdProb
	}
}

// checkCollision returns true if the dino and obstacle sprites overlap pixel-perfect (non-space chars).
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
	dX1 := g.dino.X + dW - 1
	dY1 := int(g.dino.posY)

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
	oX1 := oX0 + oW - 1
	oY1 := g.obstacle.Y

	// Quick reject bounding boxes
	if dX1 < oX0 || oX1 < dX0 || dY1 < oY0 || oY1 < dY0 {
		return false
	}

	// Pixel-perfect collision: check overlapping cells
	xStart := max(dX0, oX0)
	xEnd := min(dX1, oX1)
	yStart := max(dY0, oY0)
	yEnd := min(dY1, oY1)

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

// SetWidth updates game width based on terminal size
func SetWidth(w int) {
	width = w
}

// Game holds all state
type Game struct {
	dino                 *Dino
	obstacle             *Obstacle
	ticker               *time.Ticker
	events               chan termbox.Event
	score                int
	groundStart          int
	groundEnd            int
	started              bool
	groundExtending      bool
	stageIndexActive     int
	stageIndexTarget     int
	stageTransitionStart time.Time
}

// NewGame initializes and returns a new Game
func NewGame() *Game {
	events := make(chan termbox.Event)
	go func() {
		for {
			events <- termbox.PollEvent()
		}
	}()
	// initialize player
	d := NewDino()
	// calculate initial ground boundaries
	half := initialGroundLength / 2
	gs := d.X - half
	if gs < 0 {
		gs = 0
	}
	ge := d.X + half
	if ge > width-1 {
		ge = width - 1
	}
	return &Game{
		dino:                 d,
		obstacle:             NewObstacle(),
		ticker:               time.NewTicker(tickDuration),
		events:               events,
		score:                0,
		groundStart:          gs,
		groundEnd:            ge,
		started:              false,
		groundExtending:      false,
		stageIndexActive:     0,
		stageIndexTarget:     0,
		stageTransitionStart: time.Time{},
	}
}

// handleEvent processes a single input event
func (g *Game) handleEvent(ev termbox.Event) bool {
	if ev.Type == termbox.EventKey {
		switch ev.Key {
		case KeyJump, KeyJumpAlt:
			if !g.started {
				g.started = true
				g.groundExtending = true
			}
			g.dino.Jump()
			// cancel duck when jumping
			g.dino.duckFrames = 0
		case KeyDuck:
			// only allow ducking when on ground
			if g.started && int(g.dino.posY) == height-2 {
				g.dino.duckFrames = duckHoldDuration
			}
		case KeyQuit:
			return false
		}
		if ev.Ch == KeyQuitRune {
			return false
		}
	}
	return true
}

// update updates game state
func (g *Game) update() {
	g.dino.Update()
	if g.started {
		g.applyStage()
		g.obstacle.Update()
		if g.checkCollision() {
			g.gameOver()
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

// draw renders the current game state
func (g *Game) draw() {
	ClearScreen()
	// display score and quit hint
	PrintAt(0, 0, fmt.Sprintf("Score: %d  (q to quit)", g.score))
	g.drawGround()
	g.dino.Draw()
	if !g.started {
		PrintCenter("Press Space or Up Arrow to Start")
		termbox.Flush()
		return
	}
	// only draw obstacle if it is on extended ground
	xPos := int(math.Round(g.obstacle.posX))
	if xPos >= g.groundStart && xPos <= g.groundEnd {
		g.obstacle.Draw()
	}
	termbox.Flush()
}

// Run starts the game loop
func (g *Game) Run() {
	for range g.ticker.C {
		select {
		case ev := <-g.events:
			if !g.handleEvent(ev) {
				return
			}
		default:
		}
		g.update()
		if g.started {
			g.score++
		}
		g.draw()
	}
}

// gameOver displays game over screen and waits for restart or quit
func (g *Game) gameOver() {
	// overlay Game Over and options
	PrintCenter("GAME OVER (r to retry, q to quit)")
	termbox.Flush()
	for {
		ev := <-g.events
		if ev.Type == termbox.EventKey {
			if ev.Ch == KeyRestartRune {
				// reset game state
				g.dino = NewDino()
				g.obstacle = NewObstacle()
				g.score = 0
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

// drawGround draws ground only between the current boundaries.
func (g *Game) drawGround() {
	for x := g.groundStart; x <= g.groundEnd; x++ {
		termbox.SetCell(x, height-1, '_', termbox.ColorWhite, termbox.ColorDefault)
	}
}
