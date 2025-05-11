package game

import (
	"github.com/nsf/termbox-go"
	"math/rand"
)

// Ground decoration types
type GroundDecoration struct {
	x    float64 // 使用浮点数以支持平滑移动
	char rune
}

// Collection of ground decorations
var groundDecorations []GroundDecoration

// Initialize ground decorations
func InitGroundDecorations() {
	groundDecorations = make([]GroundDecoration, 0)

	// Add random decorations across the ground
	for x := 0; x < width*2; x += 2 + rand.Intn(5) { // 生成更多装饰，以便滚动时有足够的装饰
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
			x:    float64(x),
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
		intX := int(decoration.x) % (width * 2) // 使用取模运算使装饰在屏幕范围内循环
		if intX < width {
			termbox.SetCell(intX, height, decoration.char, termbox.ColorWhite, termbox.ColorDefault)
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
		intX := int(decoration.x) % (width * 2) // 使用取模运算使装饰在屏幕范围内循环
		if intX >= g.groundStart && intX <= g.groundEnd {
			termbox.SetCell(intX, height, decoration.char, termbox.ColorWhite, termbox.ColorDefault)
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
