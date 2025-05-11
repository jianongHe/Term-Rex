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
		"++    ++Q+++",
		" + +++++  w ",
		"  ++++++    ",
		"   |   |    ",
	},
	{
		"       ++++ ",
		" +    ++Q+++",
		" + +++++  w ",
		"  ++++++    ",
		"   /   /    ",
	},
}

// Animation frames for ducking Dino
var dinoDuckFrames = []Sprite{
	{
		"       ++++ ",
		" -    ++@+++",
		"  ++++++  w ",
		"   :   :    ",
	},
	{
		"       ++++ ",
		" +    ++@+++",
		"  ++++++  w ",
		"   ;   ;    ",
	},
}

// frames between animation switches (adjusted for new FPS)
const animPeriod = fps / 12

// ObstacleType represents the type of obstacle
type ObstacleType int

const (
	SingleCactusType ObstacleType = iota
	ShortCactusType
	GroupCactusType
	BirdType
	BigBirdType
)

// ObstacleFrames stores all obstacle animation frames by type
var ObstacleFrames = map[ObstacleType][]Sprite{
	SingleCactusType: {
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
	},
	ShortCactusType: {
		{
			"/:\\/:\\",
			" | |",
		},
		{
			"/|\\/|\\",
			" | |",
		},
	},
	GroupCactusType: {
		{
			"    |  ",
			"/|\\/|\\",
			" |  |",
		},
		{
			"    |  ",
			"\\|/\\|/",
			" |  |",
		},
	},
	BirdType: {
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
	},
	BigBirdType: {
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
	CactusProb     float64 // 仙人掌类别的总概率 (0-1.0)
	BirdProb       float64 // 鸟类别的总概率 (= 1.0 - CactusProb)

	// 仙人掌类别内部的概率分布 (这些值加起来应该等于1.0)
	SingleCactusRatio float64 // 单个仙人掌在仙人掌类别中的占比
	ShortCactusRatio  float64 // 矮仙人掌在仙人掌类别中的占比
	GroupCactusRatio  float64 // 组合仙人掌在仙人掌类别中的占比

	// 鸟类别内部的概率分布 (这些值加起来应该等于1.0)
	SmallBirdRatio float64 // 小鸟在鸟类别中的占比
	BigBirdRatio   float64 // 大鸟在鸟类别中的占比

	MinGap int // 障碍物之间的最小间距（屏幕单位）
	MaxGap int // 障碍物之间的最大间距（屏幕单位）
}

// stageConfigs lists the stages in ascending order of score threshold.
var stageConfigs = []StageConfig{
	{
		ScoreThreshold:    0,
		Speed:             1.4,
		CactusProb:        0.90, // 90% 仙人掌, 10% 鸟类
		SingleCactusRatio: 0.60, // 仙人掌类别内: 60% 单个仙人掌
		ShortCactusRatio:  0.25, // 仙人掌类别内: 25% 矮仙人掌
		GroupCactusRatio:  0.15, // 仙人掌类别内: 15% 组合仙人掌
		SmallBirdRatio:    0.90, // 鸟类别内: 90% 小鸟
		BigBirdRatio:      0.10, // 鸟类别内: 10% 大鸟
		MinGap:            80,
		MaxGap:            90,
	},
	{
		ScoreThreshold:    100,
		Speed:             1.8,
		CactusProb:        0.80, // 80% 仙人掌, 20% 鸟类
		SingleCactusRatio: 0.55, // 仙人掌类别内: 55% 单个仙人掌
		ShortCactusRatio:  0.25, // 仙人掌类别内: 25% 矮仙人掌
		GroupCactusRatio:  0.20, // 仙人掌类别内: 20% 组合仙人掌
		SmallBirdRatio:    0.85, // 鸟类别内: 85% 小鸟
		BigBirdRatio:      0.15, // 鸟类别内: 15% 大鸟
		MinGap:            70,
		MaxGap:            85,
	},
	{
		ScoreThreshold:    300,
		Speed:             2.2,
		CactusProb:        0.75, // 75% 仙人掌, 25% 鸟类
		SingleCactusRatio: 0.50, // 仙人掌类别内: 50% 单个仙人掌
		ShortCactusRatio:  0.25, // 仙人掌类别内: 25% 矮仙人掌
		GroupCactusRatio:  0.25, // 仙人掌类别内: 25% 组合仙人掌
		SmallBirdRatio:    0.80, // 鸟类别内: 80% 小鸟
		BigBirdRatio:      0.20, // 鸟类别内: 20% 大鸟
		MinGap:            60,
		MaxGap:            80,
	},
	{
		ScoreThreshold:    500,
		Speed:             2.8,
		CactusProb:        0.70, // 70% 仙人掌, 30% 鸟类
		SingleCactusRatio: 0.45, // 仙人掌类别内: 45% 单个仙人掌
		ShortCactusRatio:  0.20, // 仙人掌类别内: 20% 矮仙人掌
		GroupCactusRatio:  0.35, // 仙人掌类别内: 35% 组合仙人掌
		SmallBirdRatio:    0.75, // 鸟类别内: 75% 小鸟
		BigBirdRatio:      0.25, // 鸟类别内: 25% 大鸟
		MinGap:            50,
		MaxGap:            75,
	},
	{
		ScoreThreshold:    1000,
		Speed:             3,
		CactusProb:        0.65, // 65% 仙人掌, 35% 鸟类
		SingleCactusRatio: 0.40, // 仙人掌类别内: 40% 单个仙人掌
		ShortCactusRatio:  0.20, // 仙人掌类别内: 20% 矮仙人掌
		GroupCactusRatio:  0.40, // 仙人掌类别内: 40% 组合仙人掌
		SmallBirdRatio:    0.70, // 鸟类别内: 70% 小鸟
		BigBirdRatio:      0.30, // 鸟类别内: 30% 大鸟
		MinGap:            47,
		MaxGap:            70,
	},
	{
		ScoreThreshold:    1500,
		Speed:             3,
		CactusProb:        0.60, // 60% 仙人掌, 40% 鸟类
		SingleCactusRatio: 0.35, // 仙人掌类别内: 35% 单个仙人掌
		ShortCactusRatio:  0.15, // 仙人掌类别内: 15% 矮仙人掌
		GroupCactusRatio:  0.50, // 仙人掌类别内: 50% 组合仙人掌
		SmallBirdRatio:    0.65, // 鸟类别内: 65% 小鸟
		BigBirdRatio:      0.35, // 鸟类别内: 35% 大鸟
		MinGap:            40,
		MaxGap:            65,
	},
	{
		ScoreThreshold:    2000,
		Speed:             3,
		CactusProb:        0.60, // 60% 仙人掌, 40% 鸟类
		SingleCactusRatio: 0.35, // 仙人掌类别内: 35% 单个仙人掌
		ShortCactusRatio:  0.15, // 仙人掌类别内: 15% 矮仙人掌
		GroupCactusRatio:  0.50, // 仙人掌类别内: 50% 组合仙人掌
		SmallBirdRatio:    0.65, // 鸟类别内: 65% 小鸟
		BigBirdRatio:      0.35, // 鸟类别内: 35% 大鸟
		MinGap:            40,
		MaxGap:            60,
	},
	{
		ScoreThreshold:    2500,
		Speed:             3.2,
		CactusProb:        0.60, // 60% 仙人掌, 40% 鸟类
		SingleCactusRatio: 0.35, // 仙人掌类别内: 35% 单个仙人掌
		ShortCactusRatio:  0.15, // 仙人掌类别内: 15% 矮仙人掌
		GroupCactusRatio:  0.50, // 仙人掌类别内: 50% 组合仙人掌
		SmallBirdRatio:    0.65, // 鸟类别内: 65% 小鸟
		BigBirdRatio:      0.35, // 鸟类别内: 35% 大鸟
		MinGap:            40,
		MaxGap:            55,
	},
	{
		ScoreThreshold:    3000,
		Speed:             3.3,
		CactusProb:        0.60, // 60% 仙人掌, 40% 鸟类
		SingleCactusRatio: 0.35, // 仙人掌类别内: 35% 单个仙人掌
		ShortCactusRatio:  0.15, // 仙人掌类别内: 15% 矮仙人掌
		GroupCactusRatio:  0.50, // 仙人掌类别内: 50% 组合仙人掌
		SmallBirdRatio:    0.65, // 鸟类别内: 65% 小鸟
		BigBirdRatio:      0.35, // 鸟类别内: 35% 大鸟
		MinGap:            30,
		MaxGap:            50,
	},
	{
		ScoreThreshold:    6000,
		Speed:             3.5,
		CactusProb:        0.60, // 60% 仙人掌, 40% 鸟类
		SingleCactusRatio: 0.35, // 仙人掌类别内: 35% 单个仙人掌
		ShortCactusRatio:  0.15, // 仙人掌类别内: 15% 矮仙人掌
		GroupCactusRatio:  0.50, // 仙人掌类别内: 50% 组合仙人掌
		SmallBirdRatio:    0.65, // 鸟类别内: 65% 小鸟
		BigBirdRatio:      0.35, // 鸟类别内: 35% 大鸟
		MinGap:            25,
		MaxGap:            50,
	},
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
	KeyQuitRune    = 'q' // alternate quit
	KeyRestartRune = 'r' // restart game
	KeyPauseRune   = 'p' // pause/resume game
)

// 障碍物组合配置
type ObstacleCombination struct {
	Type1    ObstacleType // 第一个障碍物类型
	Type2    ObstacleType // 第二个障碍物类型（可选）
	Gap      int          // 两个障碍物之间的间距（如果有第二个障碍物）
	HasCombo bool         // 是否是组合障碍物
}

// 障碍物组合概率（随着难度增加）
var obstacleComboProbability = 0.0
