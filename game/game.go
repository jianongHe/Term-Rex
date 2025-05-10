package game

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"math"
	"time"
)

// SetWidth updates game width based on terminal size
func SetWidth(w int) {
	width = w
}

// Game holds all state
type Game struct {
	dino                 *Dino
	obstacleManager      *ObstacleManager
	ticker               *time.Ticker
	events               chan termbox.Event
	score                int
	highestScore         int
	groundStart          int
	groundEnd            int
	started              bool
	groundExtending      bool
	collided             bool // indicates collision occurred
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
		obstacleManager:      NewObstacleManager(),
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

// drawStartScreen renders the initial start prompt and partial ground
func (g *Game) drawStartScreen() {
	g.drawGroundPartial()
	// draw the dinosaur at its starting position
	g.dino.Draw()
	PrintCenter("Press Space or Up Arrow to Start")
}

// drawGameScene renders the full game scene after start
func (g *Game) drawGameScene() {
	// ground
	if g.groundExtending {
		g.drawGroundPartial()
	} else {
		DrawGround()
	}
	// dino
	g.dino.Draw()
	// obstacle
	obstacle := g.obstacleManager.GetCurrentObstacle()
	x, _ := obstacle.GetPosition()
	xPos := int(math.Round(x))
	if xPos >= g.groundStart && xPos <= g.groundEnd {
		g.obstacleManager.Draw()
	}
}

// draw renders the current game state
func (g *Game) draw() {
	ClearScreen()
	// score and quit hint
	PrintAt(0, 0, fmt.Sprintf("Score: %d  (q to quit)", g.score))
	if g.highestScore > 0 {
		hsText := fmt.Sprintf("High: %d", g.highestScore)
		x := width - len(hsText)
		PrintAt(x, 0, hsText)
	}
	if !g.started {
		g.drawStartScreen()
		termbox.Flush()
		return
	}
	// main game view
	g.drawGameScene()
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
		if g.collided {
			g.draw()
			g.gameOver()
			// clear collision flag and restart loop
			g.collided = false
			continue
		}
		if g.started {
			g.score++
		}
		g.draw()
	}
}
