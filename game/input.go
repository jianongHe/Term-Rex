package game

import "github.com/nsf/termbox-go"

// handleEvent processes a single input event.
func (g *Game) handleEvent(ev termbox.Event) bool {
	if ev.Type == termbox.EventKey {
		switch ev.Key {
		case KeyJump, KeyJumpAlt:
			if !g.started {
				g.started = true
				g.groundExtending = true
			}
			g.dino.Jump()
			// cancel duck when jumping
			g.dino.duckFrames = 0
		case KeyDuck:
			// only allow ducking when on ground
			if g.started && int(g.dino.posY) == height-2 {
				g.dino.duckFrames = duckHoldDuration
			}
		case KeyQuit:
			return false
		case termbox.KeyCtrlC: // 添加对 Ctrl+C 的处理
			return false
		}

		// 处理字符键
		switch ev.Ch {
		case KeyQuitRune:
			return false
		case 'm': // 音效开关
			am := GetAudioManager()
			am.SetEnabled(!am.IsEnabled())
		}
	}
	return true
}
