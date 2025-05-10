package main

import (
	"github.com/jianongHe/term-rex/game"
	"github.com/nsf/termbox-go"
)

func main() {
	// initialize terminal
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	// dynamic width based on terminal size
	w, _ := termbox.Size()
	game.SetWidth(w)

	// start game
	g := game.NewGame()
	g.Run()
}
