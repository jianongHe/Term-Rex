package game

import (
	"github.com/nsf/termbox-go"
	"time"
)

// checkCollision returns true if the dino and obstacle sprites overlap.
func (g *Game) checkCollision() bool {
	// Dino bounds
	dW := len(dinoSprite[0])
	dH := len(dinoSprite)
	dX0 := g.dino.X
	dY0 := g.dino.Y - (dH - 1)
	dX1 := g.dino.X + dW - 1
	dY1 := g.dino.Y

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

var (
	width  = 50
	height = 10
	fps    = 24
)

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
		ticker:   time.NewTicker(time.Second / time.Duration(fps)),
		events:   events,
	}
}

// handleEvent processes a single input event
func (g *Game) handleEvent(ev termbox.Event) bool {
	if ev.Type == termbox.EventKey {
		switch ev.Key {
		case termbox.KeySpace:
			g.dino.Jump()
		case termbox.KeyEsc:
			return false
		}
		if ev.Ch == 'q' {
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
		g.draw()
	}
}

// gameOver displays game over screen and waits for quit
func (g *Game) gameOver() {
	ClearScreen()
	PrintCenter("GAME OVER (press q)")
	termbox.Flush()
	for ev := range g.events {
		if ev.Type == termbox.EventKey && (ev.Key == termbox.KeyEsc || ev.Ch == 'q') {
			return
		}
	}
}
