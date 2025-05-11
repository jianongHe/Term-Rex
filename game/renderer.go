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

// 地面线字符类型
type GroundLineChar struct {
	x    float64 // 使用浮点数以支持平滑移动
	char rune
}

// Collection of ground decorations
var groundDecorations []GroundDecoration

// 地面线字符集合
var groundLineChars []GroundLineChar

// Initialize ground decorations
func InitGroundDecorations() {
	groundDecorations = make([]GroundDecoration, 0)
	groundLineChars = make([]GroundLineChar, 0)

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

	// 初始化地面线字符 - 以较大间隔放置特殊字符
	// 首先用下划线填充整个地面
	for x := 0; x < width*2; x++ {
		groundLineChars = append(groundLineChars, GroundLineChar{
			x:    float64(x),
			char: '_', // 默认全部使用下划线
		})
	}

	// 然后在较大间隔处放置特殊字符
	minInterval := 200              // 最小间隔距离
	nextSpecialPos := rand.Intn(50) // 第一个特殊字符的位置（随机起点）

	for nextSpecialPos < width*2 {
		// 选择一个特殊字符
		var specialChar rune
		switch rand.Intn(4) {
		case 0:
			specialChar = '='
		case 1:
			specialChar = '~'
		case 2:
			specialChar = '-'
		case 3:
			specialChar = '^'
		}

		// 在特定位置放置特殊字符
		if nextSpecialPos < len(groundLineChars) {
			groundLineChars[nextSpecialPos].char = specialChar
		}

		// 计算下一个特殊字符的位置
		// 最小间隔为minInterval，再加上一些随机变化
		nextSpecialPos += minInterval + rand.Intn(100)
	}
}

// ClearScreen clears the terminal
func ClearScreen() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

// DrawGround draws the ground line with decorations
func DrawGround() {
	// Draw the main ground line using varied characters
	for x := 0; x < width; x++ {
		// 查找对应位置的地面线字符
		found := false
		for _, lineChar := range groundLineChars {
			intX := int(lineChar.x) % (width * 2)
			if intX == x {
				termbox.SetCell(x, height-1, lineChar.char, termbox.ColorWhite, termbox.ColorDefault)
				found = true
				break
			}
		}

		// 如果没有找到对应的特殊字符，使用默认的下划线
		if !found {
			termbox.SetCell(x, height-1, '_', termbox.ColorWhite, termbox.ColorDefault)
		}
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
	// Draw the main ground line using varied characters
	for x := g.groundStart; x <= g.groundEnd; x++ {
		// 查找对应位置的地面线字符
		found := false
		for _, lineChar := range groundLineChars {
			intX := int(lineChar.x) % (width * 2)
			if intX == x {
				termbox.SetCell(x, height-1, lineChar.char, termbox.ColorWhite, termbox.ColorDefault)
				found = true
				break
			}
		}

		// 如果没有找到对应的特殊字符，使用默认的下划线
		if !found {
			termbox.SetCell(x, height-1, '_', termbox.ColorWhite, termbox.ColorDefault)
		}
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
