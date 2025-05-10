package game

import (
	"github.com/nsf/termbox-go"
	"time"
)

var (
	width  int
	height = 10
	fps    = 24
)

// Game holds all state
type Game struct {
	dino     *Dino
	obstacle *Obstacle
	ticker   *time.Ticker
	events   chan termbox.Event
}

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
		ticker:   time.NewTicker(time.Second / fps),
		events:   events,
	}
}

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

func (g *Game) update() {
	g.dino.Update()
	g.obstacle.Update()
	if g.obstacle.X == g.dino.X && g.obstacle.Y == g.dino.Y {
		g.gameOver()
	}
}

func (g *Game) draw() {
	ClearScreen()
	DrawGround()
	g.dino.Draw()
	g.obstacle.Draw()
	termbox.Flush()
}

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
