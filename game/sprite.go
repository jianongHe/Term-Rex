package game

import "github.com/nsf/termbox-go"

// Sprite 是一组字符串，表示多行 ASCII 艺术图
type Sprite []string

// Draw 在 (x,y) 处逐字符绘制非空格字符
func (s Sprite) Draw(x, y int, fg, bg termbox.Attribute) {
	for row, line := range s {
		for col, ch := range line {
			if ch != ' ' {
				termbox.SetCell(x+col, y+row, ch, fg, bg)
			}
		}
	}
}
