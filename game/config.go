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
const jumpHeight = 5

// 跳跃持续帧数
const jumpDuration = fps / 4

// initial jump velocity (calculated to reach jumpHeight in jumpDuration frames)
var jumpVelocity = -2 * float64(jumpHeight) / float64(jumpDuration)

// gravity acceleration (calculated to bring velocity back to zero over jumpDuration frames)
var gravity = -jumpVelocity / float64(jumpDuration)

// hang time at apex in frames
const hangDuration = 2

// 障碍物每帧移动的格数
var obstacleSpeed float64 = 1.0

// ground extension speed in cells per frame
const groundExtendSpeed = 3

// initial ground length in cells (total width)
const initialGroundLength = 24

// duck hold duration in frames
const duckHoldDuration = fps - 11

// 每帧间隔
var tickDuration = time.Second / time.Duration(fps)

// Animation frames for standing Dino
var dinoStandFrames = []Sprite{
	{
		"       ++++ ",
		"++    ++@++ ",
		" + +++++  + ",
		"  ++++++    ",
		"   |   |    ",
	},
	{
		"       ++++ ",
		" +    ++@++ ",
		" + +++++  + ",
		"  ++++++    ",
		"   /   /    ",
	},
}

// Animation frames for ducking Dino
var dinoDuckFrames = []Sprite{
	{
		"       ++++ ",
		" -    ++@++ ",
		"  ++++++    ",
		"   :   :    ",
	},
	{
		"       ++++ ",
		" +    ++@++ ",
		"  ++++++    ",
		"   ;   ;    ",
	},
}

// frames between animation switches
const animPeriod = fps / 6

// Animation frames for cactus obstacles
var obstacleFrames = []Sprite{
	{
		" | ",
		"/|\\",
		" | ",
	},
	{
		" | ",
		"\\|/",
		" | ",
	},
}

// Animation frames for flying birds
var birdFrames = []Sprite{
	{
		" |   ",
		"<o=- ",
		" |   ",
	},
	{
		" /   ",
		"<O=- ",
		" \\   ",
	},
}

// Animation frames for big birds
var bigBirdFrames = []Sprite{
	{
		"  /\\    ",
		" /  \\   ",
		"<ooo=-- ",
		" \\__/   ",
	},
	{
		"  /\\    ",
		" /  \\   ",
		"<OOO=-- ",
		" \\__/   ",
	},
}

// Cloud sprites with different shapes
var cloudSprites = []Sprite{
	{
		"   .--.    ",
		" .(    ).  ",
		"(___.__)  ",
	},
	{
		"  .-.      ",
		" (   ).    ",
		"(___(__)  ",
	},
	{
		"    .--.   ",
		".-(    )-. ",
		"(________) ",
	},
}

// bird appearance probability
var birdProbability = 0.3

// big bird appearance probability
var bigBirdProbability = 0.15

// bird flight height (row index) above bottom of screen
const birdFlightRow = 9

// big bird flight height (row index) above bottom of screen
const bigBirdFlightRow = 7

// —— 云朵配置参数 ——

// 云朵最小高度（行号，从上往下计数）
const cloudMinHeight = 1

// 云朵最大高度（行号，从上往下计数）
const cloudMaxHeight = 3

// 云朵最小移动速度
const cloudMinSpeed = 0.2

// 云朵最大移动速度
const cloudMaxSpeed = 0.5

// 初始云朵最小数量
const cloudMinCount = 3

// 初始云朵最大数量
const cloudMaxCount = 3

// 云朵右侧边缘缓冲区大小
const cloudRightEdgeBuffer = 15

// 云朵最小额外间距
const cloudMinExtraSpace = 10

// 云朵最大额外间距
const cloudMaxExtraSpace = 25

// 云朵之间的最小缓冲区大小
const cloudBufferSpace = 5

// StageConfig defines dynamic game parameters per stage based on score.
type StageConfig struct {
	ScoreThreshold int     // minimum score to enter this stage
	Speed          float64 // obstacleSpeed for this stage
	BirdProb       float64 // birdProbability for this stage
	BigBirdProb    float64 // bigBirdProbability for this stage
}

// stageConfigs lists the stages in ascending order of score threshold.
var stageConfigs = []StageConfig{
	{ScoreThreshold: 0, Speed: 1.4, BirdProb: 0.1, BigBirdProb: 0.05},
	{ScoreThreshold: 300, Speed: 2.0, BirdProb: 0.3, BigBirdProb: 0.1},
	{ScoreThreshold: 1000, Speed: 3.0, BirdProb: 0.5, BigBirdProb: 0.2},
}

// duration of smooth transition between stages
var stageTransitionDuration = 3000 * time.Millisecond

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
