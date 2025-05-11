package game

import (
	"github.com/nsf/termbox-go"
	"time"
)

// —— 可调参数 ——

// 默认屏幕宽度（会被 SetWidth 覆盖）
var width = 80

// 最大有效游戏宽度（超过这个宽度，障碍物不会从更远处生成）
const maxEffectiveWidth = 120

// 固定游戏高度（行数）
const height = 15

// 帧率（FPS）
const fps = 60 // 从24提高到60，使游戏更流畅

// 原始帧率（用于计算速度调整因子）
const originalFps = 24

// 速度调整因子（保持与原始24FPS相同的游戏速度）
const speedFactor = float64(originalFps) / float64(fps)

// 跳跃高度（行数）
const jumpHeight = 5

// 跳跃持续帧数（根据新帧率调整）
const jumpDuration = fps / 4

// initial jump velocity (calculated to reach jumpHeight in jumpDuration frames)
var jumpVelocity = -2 * float64(jumpHeight) / float64(jumpDuration)

// gravity acceleration (calculated to bring velocity back to zero over jumpDuration frames)
var gravity = -jumpVelocity / float64(jumpDuration)

// hang time at apex in frames
const hangDuration = 2

// 障碍物每帧移动的格数（根据速度因子调整）
var obstacleSpeed float64 = 1.0 * speedFactor

// ground extension speed in cells per frame（根据速度因子调整）
var groundExtendSpeed float64 = 3 * speedFactor

// initial ground length in cells (total width)
const initialGroundLength = 24

// duck hold duration in frames
var duckHoldDuration = fps/2 + 1

// 每帧间隔
var tickDuration = time.Second / time.Duration(fps)

// 音效相关配置
const (
	AudioEnabled   = true // 默认启用音效
	ScoreMilestone = 100  // 每得到100分播放一次得分音效
)

// 分数闪烁相关配置
const (
	ScoreBlinkDuration = 1500 * time.Millisecond // 分数闪烁持续时间（1.5秒）
	ScoreBlinkInterval = 200 * time.Millisecond  // 分数闪烁间隔（200毫秒）
)

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

// frames between animation switches (adjusted for new FPS)
const animPeriod = fps / 12

// Animation frames for single cactus obstacles
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

// Animation frames for group cactus obstacles
var groupCactusFrames = []Sprite{
	{
		"    |  ",
		"/|\\/|\\",
		" |  |",
	},
	{
		"    |  ",
		"\\|//|\\",
		" |  |",
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
	ScoreThreshold  int     // minimum score to enter this stage
	Speed           float64 // obstacleSpeed for this stage
	BirdProb        float64 // birdProbability for this stage
	BigBirdProb     float64 // bigBirdProbability for this stage
	GroupCactusProb float64 // groupCactusProbability for this stage
	MinGap          int     // 障碍物之间的最小间距（屏幕单位）
	MaxGap          int     // 障碍物之间的最大间距（屏幕单位）
}

// stageConfigs lists the stages in ascending order of score threshold.
var stageConfigs = []StageConfig{
	{ScoreThreshold: 0, Speed: 1.2, BirdProb: 0.05, BigBirdProb: 0.01, GroupCactusProb: 0.10, MinGap: 60, MaxGap: 90},
	{ScoreThreshold: 100, Speed: 1.4, BirdProb: 0.10, BigBirdProb: 0.03, GroupCactusProb: 0.15, MinGap: 55, MaxGap: 85},
	{ScoreThreshold: 300, Speed: 1.6, BirdProb: 0.15, BigBirdProb: 0.05, GroupCactusProb: 0.20, MinGap: 50, MaxGap: 80},
	{ScoreThreshold: 600, Speed: 1.8, BirdProb: 0.20, BigBirdProb: 0.08, GroupCactusProb: 0.25, MinGap: 45, MaxGap: 75},
	{ScoreThreshold: 1000, Speed: 2.0, BirdProb: 0.25, BigBirdProb: 0.10, GroupCactusProb: 0.30, MinGap: 40, MaxGap: 70},
	{ScoreThreshold: 1500, Speed: 2.3, BirdProb: 0.30, BigBirdProb: 0.15, GroupCactusProb: 0.30, MinGap: 35, MaxGap: 65},
	{ScoreThreshold: 2000, Speed: 2.6, BirdProb: 0.40, BigBirdProb: 0.18, GroupCactusProb: 0.35, MinGap: 30, MaxGap: 60},
	{ScoreThreshold: 2500, Speed: 2.8, BirdProb: 0.45, BigBirdProb: 0.20, GroupCactusProb: 0.35, MinGap: 28, MaxGap: 55},
	{ScoreThreshold: 3000, Speed: 3.0, BirdProb: 0.50, BigBirdProb: 0.25, GroupCactusProb: 0.40, MinGap: 25, MaxGap: 50},
}

// duration of smooth transition between stages
var stageTransitionDuration = 3000 * time.Millisecond

// Key bindings
const (
	KeyJump    = termbox.KeySpace     // jump action
	KeyJumpAlt = termbox.KeyArrowUp   // alternate jump action
	KeyDuck    = termbox.KeyArrowDown // duck action
	KeyQuit    = termbox.KeyEsc       // quit action
	KeyRelease = termbox.KeyArrowDown // 用于检测下键释放（实际上是同一个键）
)

// Character key bindings
const (
	KeyQuitRune = 'q' // alternate quit
	// restart game
	KeyRestartRune = 'r'
)

// group cactus appearance probability (default value, will be overridden by stage config)
var groupCactusProbability = 0.25

// 障碍物组合配置
type ObstacleCombination struct {
	Type1    ObstacleType // 第一个障碍物类型
	Type2    ObstacleType // 第二个障碍物类型（可选）
	Gap      int          // 两个障碍物之间的间距（如果有第二个障碍物）
	HasCombo bool         // 是否是组合障碍物
}

// 障碍物组合概率（随着难度增加）
var obstacleComboProbability = 0.0
