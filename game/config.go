package game

import (
	"github.com/nsf/termbox-go"
	"time"
)

// —— 可调参数 ——

// 默认屏幕宽度（会被 SetWidth 覆盖）
var width = 80

// 固定游戏高度（行数）
const height = 15

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

// ground extension speed in cells per frame
const groundExtendSpeed = 3

// initial ground length in cells (total width)
const initialGroundLength = 6

// duck hold duration in frames
const duckHoldDuration = fps - 11

// 每帧间隔
var tickDuration = time.Second / time.Duration(fps)

// ASCII art sprite for the dinosaur

var dinoSprite = Sprite{
	" /oo\\  ",
	"/|  |\\ ",
	"  ||   ",
	" /  \\  ",
}

// ASCII art sprite for the crouching dinosaur
var dinoDuckSprite = Sprite{
	"  _  ",
	" /_\\ ",
	"/___\\",
}

var obstacleSprite = Sprite{
	" | ",
	"/|\\",
	" | ",
}

// bird appearance probability
const birdProbability = 0.3

// ASCII art sprite for the bird
var birdSprite = Sprite{
	" <o> ",
}

// Key bindings
const (
	KeyJump    = termbox.KeySpace     // jump action
	KeyJumpAlt = termbox.KeyArrowUp   // alternate jump action
	KeyDuck    = termbox.KeyArrowDown // duck action
	KeyQuit    = termbox.KeyEsc       // quit action
)

// Character key bindings
const (
	KeyQuitRune = 'q' // alternate quit
	// restart game
	KeyRestartRune = 'r'
)
