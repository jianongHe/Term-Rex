package game

import (
	"github.com/nsf/termbox-go"
	"math/rand"
)

// Ground decoration types
type GroundDecoration struct {
	x    int
	char rune
}

// Collection of ground decorations
var groundDecorations []GroundDecoration

// Initialize ground decorations
func InitGroundDecorations() {
	groundDecorations = make([]GroundDecoration, 0)

	// Add random decorations across the ground
	for x := 0; x < width; x += 2 + rand.Intn(5) {
		// Choose a decoration character
		var char rune
		switch rand.Intn(5) {
		case 0:
			char = '.'
		case 1:
			char = ','
		case 2:
			char = '\''
		case 3:
			char = '`'
		case 4:
			char = '-'
		}

		groundDecorations = append(groundDecorations, GroundDecoration{
			x:    x,
			char: char,
		})
	}
}

// ClearScreen clears the terminal
func ClearScreen() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

// DrawGround draws the ground line with decorations
func DrawGround() {
	// Draw the main ground line
	for x := 0; x < width; x++ {
		termbox.SetCell(x, height-1, '_', termbox.ColorWhite, termbox.ColorDefault)
	}

	// Draw decorations below the ground
	for _, decoration := range groundDecorations {
		if decoration.x < width {
			termbox.SetCell(decoration.x, height, decoration.char, termbox.ColorWhite, termbox.ColorDefault)
		}
	}
}

// drawGroundPartial draws ground between current Game boundaries with decorations
func (g *Game) drawGroundPartial() {
	// Draw the main ground line
	for x := g.groundStart; x <= g.groundEnd; x++ {
		termbox.SetCell(x, height-1, '_', termbox.ColorWhite, termbox.ColorDefault)
	}

	// Draw decorations below the ground
	for _, decoration := range groundDecorations {
		if decoration.x >= g.groundStart && decoration.x <= g.groundEnd {
			termbox.SetCell(decoration.x, height, decoration.char, termbox.ColorWhite, termbox.ColorDefault)
		}
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

// PrintCenterAt prints a message centered horizontally at the specified row
func PrintCenterAt(msg string, row int) {
	x := (width - len(msg)) / 2
	for i, c := range msg {
		termbox.SetCell(x+i, row, c, termbox.ColorWhite, termbox.ColorDefault)
	}
}

// PrintAt prints a message at the specified coordinates.
func PrintAt(x, y int, msg string) {
	for i, ch := range msg {
		termbox.SetCell(x+i, y, ch, termbox.ColorWhite, termbox.ColorDefault)
	}
}
