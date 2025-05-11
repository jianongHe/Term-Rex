package game

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"strings"
	"time"
)

// SetWidth updates game width based on terminal size
func SetWidth(w int) {
	width = w
	// Reinitialize ground decorations when width changes
	InitGroundDecorations()
}

// Game holds all state
type Game struct {
	dino                     *Dino
	obstacleManager          *ObstacleManager
	downKeyHeld              bool // 添加一个字段来跟踪下键状态
	cloudManager             *CloudManager
	ticker                   *time.Ticker
	events                   chan termbox.Event
	score                    int
	highestScore             int
	groundStart              int
	groundEnd                int
	started                  bool
	groundExtending          bool
	collided                 bool // indicates collision occurred
	stageIndexActive         int
	stageIndexTarget         int
	stageTransitionStart     time.Time
	scoreBlinking            bool      // 标记分数是否正在闪烁
	scoreBlinkStart          time.Time // 分数开始闪烁的时间
	scoreBlinkVisible        bool      // 控制分数闪烁的显示/隐藏状态
	lastBlinkToggle          time.Time // 上次闪烁状态切换的时间
	groundSpecialCharCounter int       // 用于控制特殊地面字符的添加频率
}

// NewGame initializes and returns a new Game
func NewGame() *Game {
	events := make(chan termbox.Event)
	go func() {
		for {
			events <- termbox.PollEvent()
		}
	}()
	// initialize player
	d := NewDino()
	// calculate initial ground boundaries
	half := initialGroundLength / 2
	gs := d.X - half
	if gs < 0 {
		gs = 0
	}
	ge := d.X + half
	if ge > width-1 {
		ge = width - 1
	}

	// Initialize ground decorations
	InitGroundDecorations()

	// 加载历史最高分
	highScore, err := LoadHighScore()
	if err != nil {
		// 如果加载失败，使用默认值0
		fmt.Println("无法加载最高分:", err)
		highScore = 0
	}

	return &Game{
		dino:                     d,
		obstacleManager:          NewObstacleManager(),
		cloudManager:             NewCloudManager(),
		ticker:                   time.NewTicker(tickDuration),
		events:                   events,
		score:                    0,
		highestScore:             highScore,
		groundStart:              gs,
		groundEnd:                ge,
		started:                  false,
		groundExtending:          false,
		downKeyHeld:              false,
		stageIndexActive:         0,
		stageIndexTarget:         0,
		stageTransitionStart:     time.Time{},
		scoreBlinking:            false,
		scoreBlinkStart:          time.Time{},
		scoreBlinkVisible:        true,
		lastBlinkToggle:          time.Time{},
		groundSpecialCharCounter: 0,
	}
}

// drawStartScreen renders the initial start prompt and partial ground
func (g *Game) drawStartScreen() {
	// Draw clouds first (always across the entire sky)
	g.cloudManager.Draw()

	// Draw partial ground
	g.drawGroundPartial()

	// Draw the dinosaur at its starting position
	g.dino.Draw()

	PrintCenter("Press Space or Up Arrow to Start")

	// 显示音效控制提示
	//soundMsg := "Press 'm' to toggle sound"
	//if !GetAudioManager().IsEnabled() {
	//	soundMsg = "Sound OFF - Press 'm' to enable"
	//}
	//PrintCenterAt(soundMsg, height/2+2)
}

// drawGameScene renders the full game scene after start
func (g *Game) drawGameScene() {
	// Draw clouds first (always across the entire sky)
	g.cloudManager.Draw()

	// ground
	if g.groundExtending {
		g.drawGroundPartial()
	} else {
		DrawGround()
	}

	// dino
	g.dino.Draw()

	// obstacle
	g.obstacleManager.Draw()
}

// draw renders the current game state
func (g *Game) draw() {
	ClearScreen()

	// 处理分数闪烁逻辑
	if g.scoreBlinking {
		elapsed := time.Since(g.scoreBlinkStart)

		// 检查是否需要结束闪烁
		if elapsed >= ScoreBlinkDuration {
			g.scoreBlinking = false
			g.scoreBlinkVisible = true
		} else {
			// 检查是否需要切换闪烁状态
			if time.Since(g.lastBlinkToggle) >= ScoreBlinkInterval {
				g.scoreBlinkVisible = !g.scoreBlinkVisible
				g.lastBlinkToggle = time.Now()
			}
		}
	}

	// score and quit hint
	if g.scoreBlinking && !g.scoreBlinkVisible {
		// 闪烁状态下，用空格替换分数的每一位，保持原有位数
		scoreStr := fmt.Sprintf("%d", g.score)
		blankScore := strings.Repeat(" ", len(scoreStr))
		PrintAt(0, 0, fmt.Sprintf("Score: %s  (Q to quit)", blankScore))
	} else {
		PrintAt(0, 0, fmt.Sprintf("Score: %d  (Q to quit)", g.score))
	}

	// 显示音效状态
	//soundStatus := "Sound: ON (m)"
	//if !GetAudioManager().IsEnabled() {
	//	soundStatus = "Sound: OFF (m)"
	//}
	//PrintAt(width/2-len(soundStatus)/2, 0, soundStatus)

	// 始终显示最高分，即使是0
	hsText := fmt.Sprintf("High: %d", g.highestScore)
	x := width - len(hsText)
	PrintAt(x, 0, hsText)
	if !g.started {
		g.drawStartScreen()
		termbox.Flush()
		return
	}
	// main game view
	g.drawGameScene()
	termbox.Flush()
}

// Reset resets the game state for a new game
func (g *Game) Reset() {
	// 重置恐龙
	g.dino = NewDino()

	// 重置障碍物管理器
	g.obstacleManager = NewObstacleManager()

	// 重置云朵管理器
	g.cloudManager = NewCloudManager()

	// 重置分数
	g.score = 0

	// 重置游戏状态
	g.started = true
	g.groundExtending = true
	g.collided = false
	g.downKeyHeld = false

	// 重置阶段
	g.stageIndexActive = 0
	g.stageIndexTarget = 0
	g.stageTransitionStart = time.Time{}

	// 重置分数闪烁状态
	g.scoreBlinking = false
	g.scoreBlinkStart = time.Time{}
	g.scoreBlinkVisible = true
	g.lastBlinkToggle = time.Time{}
}

// TogglePause toggles the game's paused state
func (g *Game) TogglePause() {
	// 暂停/继续游戏的逻辑
	// 这里可以添加暂停功能的实现
}

// Run starts the game loop
func (g *Game) Run() {
	// 初始化音频系统
	audioManager := GetAudioManager()
	// 如果音频初始化失败，记录警告但继续游戏
	if !audioManager.IsEnabled() {
		fmt.Println("Warning: Audio system initialization failed. Game will run without sound.")
	}

	lastScoreMilestone := 0

	// 用于跟踪下键状态的变量
	lastKeyPressTime := time.Now()
	keyCheckInterval := 100 * time.Millisecond

	for range g.ticker.C {
		// 定期检查是否有按键事件
		// 如果一段时间内没有收到下键的按键事件，则认为下键已释放
		if g.downKeyHeld && time.Since(lastKeyPressTime) > keyCheckInterval {
			// 检查是否有新的按键事件
			select {
			case ev := <-g.events:
				if ev.Type == termbox.EventKey && ev.Key == KeyDuck {
					// 如果是下键，更新最后按键时间
					lastKeyPressTime = time.Now()
				} else {
					// 如果是其他键或非按键事件，处理它
					if !g.handleEvent(ev) {
						return
					}
				}
			default:
				// 如果没有新的按键事件，认为下键已释放
				g.downKeyHeld = false
				g.dino.isDownKeyPressed = false
			}
		} else {
			// 正常处理按键事件
			select {
			case ev := <-g.events:
				if ev.Type == termbox.EventKey && ev.Key == KeyDuck {
					// 如果是下键，更新最后按键时间
					lastKeyPressTime = time.Now()
				}
				if !g.handleEvent(ev) {
					return
				}
			default:
			}
		}
		g.update()
		if g.collided {
			g.draw()
			g.gameOver()
			// clear collision flag and restart loop
			g.collided = false
			continue
		}
		if g.started {
			// 只有在分数不闪烁时才增加分数
			if !g.scoreBlinking {
				g.score++

				// 每得到100分播放一次得分音效
				if g.score/ScoreMilestone > lastScoreMilestone {
					lastScoreMilestone = g.score / ScoreMilestone
					GetAudioManager().PlaySound(SoundScore)
				}
			}
		}
		g.draw()
	}
}
