package game

import "github.com/nsf/termbox-go"

// ClearScreen clears the terminal
func ClearScreen() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

// DrawGround draws the ground line
func DrawGround() {
	for x := 0; x < width; x++ {
		termbox.SetCell(x, height-1, '_', termbox.ColorWhite, termbox.ColorDefault)
	}
}

// PrintCenter prints a message at center of screen
func PrintCenter(msg string) {
	x := (width - len(msg)) / 2
	y := height / 2
	for i, c := range msg {
		termbox.SetCell(x+i, y, c, termbox.ColorWhite, termbox.ColorDefault)
	}
}

// PrintAt prints a message at the specified coordinates.
func PrintAt(x, y int, msg string) {
	for i, ch := range msg {
		termbox.SetCell(x+i, y, ch, termbox.ColorWhite, termbox.ColorDefault)
	}
}
