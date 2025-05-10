package game

import "github.com/nsf/termbox-go"

func ClearScreen() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func DrawGround() {
	for x := 0; x < width; x++ {
		termbox.SetCell(x, height-1, '_', termbox.ColorWhite, termbox.ColorDefault)
	}
}

func PrintCenter(msg string) {
	x := (width - len(msg)) / 2
	y := height / 2
	for i, c := range msg {
		termbox.SetCell(x+i, y, c, termbox.ColorWhite, termbox.ColorDefault)
	}
}
