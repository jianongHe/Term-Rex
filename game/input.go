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
			// 跳跃时重置下键状态
			g.downKeyHeld = false
			g.dino.isDownKeyPressed = false
		case KeyDuck:
			if g.started {
				// 设置下键被按住的状态
				g.downKeyHeld = true
				g.dino.isDownKeyPressed = true

				if int(g.dino.posY) == height-2 {
					// 在地面上按下键时蹲下
					g.dino.Duck()
				} else {
					// 在空中按下键时快速下降
					g.dino.FastDrop()
				}
			}
		case KeyQuit:
			return false
		case termbox.KeyCtrlC: // 添加对 Ctrl+C 的处理
			return false
		default:
			// 如果按下了其他键，认为下键已释放
			if g.downKeyHeld {
				g.downKeyHeld = false
				g.dino.isDownKeyPressed = false
			}
		}

		// 处理字符键
		switch ev.Ch {
		case KeyQuitRune:
			return false
		case 'm': // 音效开关
			am := GetAudioManager()
			am.SetEnabled(!am.IsEnabled())
			// 保持蹲下状态，如果当前正在蹲下
			if g.started && int(g.dino.posY) == height-2 && g.dino.IsDucking() {
				g.dino.Duck()
			}
		case KeyRestartRune: // 重新开始游戏
			if g.collided {
				g.Reset()
			}
		case 'p': // 暂停/继续游戏
			if g.started && !g.collided {
				g.TogglePause()
			}
		default:
			// 如果按下了其他字符键，认为下键已释放
			if g.downKeyHeld {
				g.downKeyHeld = false
				g.dino.isDownKeyPressed = false
			}
		}
	}
	return true
}
