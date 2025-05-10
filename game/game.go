package game

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
	"time"
)

// checkCollision returns true if the dino and obstacle sprites overlap.
func (g *Game) checkCollision() bool {
	// Dino bounds
	dW := len(dinoSprite[0])
	dH := len(dinoSprite)
	dX0 := g.dino.X
	// current Dino vertical position
	y := int(g.dino.posY)
	dY0 := y - (dH - 1)
	dX1 := g.dino.X + dW - 1
	dY1 := y

	// Obstacle bounds
	oW := len(obstacleSprite[0])
	oH := len(obstacleSprite)
	oX0 := g.obstacle.X
	oY0 := g.obstacle.Y - (oH - 1)
	oX1 := g.obstacle.X + oW - 1
	oY1 := g.obstacle.Y

	// Check for intersection
	return !(dX1 < oX0 || oX1 < dX0 || dY1 < oY0 || oY1 < dY0)
}

// SetWidth updates game width based on terminal size
func SetWidth(w int) {
	width = w
}

// Game holds all state
type Game struct {
	dino     *Dino
	obstacle *Obstacle
	ticker   *time.Ticker
	events   chan termbox.Event
	score    int
}

// NewGame initializes and returns a new Game
func NewGame() *Game {
	events := make(chan termbox.Event)
	go func() {
		for {
			events <- termbox.PollEvent()
		}
	}()
	return &Game{
		dino:     NewDino(),
		obstacle: NewObstacle(),
		ticker:   time.NewTicker(tickDuration),
		events:   events,
		score:    0,
	}
}

// handleEvent processes a single input event
func (g *Game) handleEvent(ev termbox.Event) bool {
	if ev.Type == termbox.EventKey {
		switch ev.Key {
		case KeyJump:
			g.dino.Jump()
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
	g.obstacle.Update()
	if g.checkCollision() {
		g.gameOver()
	}
}

// draw renders the current game state
func (g *Game) draw() {
	ClearScreen()
	// display score and quit hint
	PrintAt(0, 0, fmt.Sprintf("Score: %d  (q to quit)", g.score))
	DrawGround()
	g.dino.Draw()
	g.obstacle.Draw()
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
		g.score++
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
