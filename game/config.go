package game

import (
	"github.com/nsf/termbox-go"
	"time"
)

// —— 可调参数 ——

// 默认屏幕宽度（会被 SetWidth 覆盖）
var width = 80

// 固定游戏高度（行数）
const height = 10

// 帧率（FPS）
const fps = 24

// 跳跃高度（行数）
const jumpHeight = 4

// 跳跃持续帧数
const jumpDuration = fps / 4

// initial jump velocity (calculated to reach jumpHeight in jumpDuration frames)
var jumpVelocity = -2 * float64(jumpHeight) / float64(jumpDuration)

// gravity acceleration (calculated to bring velocity back to zero over jumpDuration frames)
var gravity = -jumpVelocity / float64(jumpDuration)

// hang time at apex in frames
const hangDuration = 2

// 障碍物每帧移动的格数
const obstacleSpeed = 1

// 每帧间隔
var tickDuration = time.Second / time.Duration(fps)

// ASCII art sprite for the dinosaur
var dinoSprite = Sprite{
	"  _  ",
	" /_\\ ",
	"/___\\",
}

var obstacleSprite = Sprite{
	" | ",
	"/|\\",
	" | ",
}

// Key bindings
const (
	KeyJump = termbox.KeySpace // jump action
	KeyQuit = termbox.KeyEsc   // quit action
)

// Character key bindings
const (
	KeyQuitRune = 'q' // alternate quit
)
