package game

import (
	"github.com/nsf/termbox-go"
	"math"
	"math/rand"
)

// ObstacleType represents the type of obstacle
type ObstacleType int

const (
	CactusType ObstacleType = iota
	GroupCactusType
	BirdType
	BigBirdType
)

// IObstacle defines the interface for all obstacle types
type IObstacle interface {
	Update()
	Draw()
	GetPosition() (float64, int)
	SetPosition(x float64, y int)
	Reset()
	GetSprite() Sprite
}

// BaseObstacle contains common properties and methods for all obstacles
type BaseObstacle struct {
	posX        float64
	y           int
	animFrame   int
	animCounter int
}

// Update moves the obstacle and updates animation
func (b *BaseObstacle) Update() {
	b.posX -= obstacleSpeed
	b.updateAnimation()
}

// GetPosition returns the current position of the obstacle
func (b *BaseObstacle) GetPosition() (float64, int) {
	return b.posX, b.y
}

// SetPosition sets the position of the obstacle
func (b *BaseObstacle) SetPosition(x float64, y int) {
	b.posX = x
	b.y = y
}

// updateAnimation advances obstacle animation frames
func (b *BaseObstacle) updateAnimation() {
	b.animCounter++
	if b.animCounter >= animPeriod {
		b.animCounter = 0
		b.animFrame = (b.animFrame + 1) % b.getFrameCount()
	}
}

// Cactus represents a cactus obstacle
type Cactus struct {
	BaseObstacle
}

// NewCactus creates a new cactus obstacle
func NewCactus() *Cactus {
	c := &Cactus{}
	c.Reset()
	return c
}

// Reset resets the cactus position and animation
func (c *Cactus) Reset() {
	c.posX = float64(width - 1)
	c.y = height - 2
	c.animFrame = 0
	c.animCounter = 0
}

// Draw renders the cactus on screen
func (c *Cactus) Draw() {
	sprite := obstacleFrames[c.animFrame]
	h := len(sprite)
	startY := c.y - (h - 1)
	x := int(math.Round(c.posX))
	sprite.Draw(x, startY, termbox.ColorRed, termbox.ColorDefault)
}

// GetSprite returns the current sprite for collision detection
func (c *Cactus) GetSprite() Sprite {
	return obstacleFrames[c.animFrame]
}

// getFrameCount returns the number of animation frames
func (c *BaseObstacle) getFrameCount() int {
	return len(obstacleFrames)
}

// Bird represents a small bird obstacle
type Bird struct {
	BaseObstacle
}

// NewBird creates a new bird obstacle
func NewBird() *Bird {
	b := &Bird{}
	b.Reset()
	return b
}

// Reset resets the bird position and animation
func (b *Bird) Reset() {
	b.posX = float64(width - 1)
	b.y = birdFlightRow
	b.animFrame = 0
	b.animCounter = 0
}

// Draw renders the bird on screen
func (b *Bird) Draw() {
	sprite := birdFrames[b.animFrame]
	h := len(sprite)
	startY := b.y - (h - 1)
	x := int(math.Round(b.posX))
	sprite.Draw(x, startY, termbox.ColorYellow, termbox.ColorDefault)
}

// GetSprite returns the current sprite for collision detection
func (b *Bird) GetSprite() Sprite {
	return birdFrames[b.animFrame]
}

// getFrameCount returns the number of animation frames
func (b *Bird) getFrameCount() int {
	return len(birdFrames)
}

// BigBird represents a large bird obstacle
type BigBird struct {
	BaseObstacle
}

// NewBigBird creates a new big bird obstacle
func NewBigBird() *BigBird {
	b := &BigBird{}
	b.Reset()
	return b
}

// Reset resets the big bird position and animation
func (b *BigBird) Reset() {
	b.posX = float64(width - 1)
	b.y = bigBirdFlightRow
	b.animFrame = 0
	b.animCounter = 0
}

// Draw renders the big bird on screen
func (b *BigBird) Draw() {
	sprite := bigBirdFrames[b.animFrame]
	h := len(sprite)
	startY := b.y - (h - 1)
	x := int(math.Round(b.posX))
	sprite.Draw(x, startY, termbox.ColorMagenta, termbox.ColorDefault)
}

// GetSprite returns the current sprite for collision detection
func (b *BigBird) GetSprite() Sprite {
	return bigBirdFrames[b.animFrame]
}

// getFrameCount returns the number of animation frames
func (b *BigBird) getFrameCount() int {
	return len(bigBirdFrames)
}

// ObstacleManager manages the creation and updating of obstacles
type ObstacleManager struct {
	obstacles    []IObstacle // 存储多个障碍物
	nextGapTimer int         // 下一个障碍物生成前的计时器
	minGap       int         // 当前阶段的最小间距
	maxGap       int         // 当前阶段的最大间距
	currentStage int         // 当前游戏阶段
}

// NewObstacleManager creates a new obstacle manager
func NewObstacleManager() *ObstacleManager {
	om := &ObstacleManager{
		obstacles:    make([]IObstacle, 0, 5), // 预分配5个障碍物的空间
		minGap:       stageConfigs[0].MinGap,
		maxGap:       stageConfigs[0].MaxGap,
		nextGapTimer: 0,
		currentStage: 0,
	}
	// 生成第一个障碍物
	om.generateNewObstacle()
	return om
}

// Update updates all obstacles and generates new ones if needed
func (om *ObstacleManager) Update() {
	// 更新所有现有障碍物
	for i := 0; i < len(om.obstacles); i++ {
		om.obstacles[i].Update()

		// 如果障碍物已经完全移出屏幕左侧，从列表中移除
		x, _ := om.obstacles[i].GetPosition()
		if x < -10 { // 使用-10确保障碍物完全离开屏幕
			// 移除障碍物（通过将最后一个元素移到当前位置，然后缩小切片）
			om.obstacles[i] = om.obstacles[len(om.obstacles)-1]
			om.obstacles = om.obstacles[:len(om.obstacles)-1]
			i-- // 调整索引，因为我们刚刚替换了当前元素
		}
	}

	// 计时器逻辑，决定何时生成新障碍物
	if om.nextGapTimer > 0 {
		om.nextGapTimer--
		if om.nextGapTimer <= 0 {
			om.generateNewObstacle()
		}
	} else if len(om.obstacles) == 0 {
		// 如果没有障碍物，立即生成一个
		om.generateNewObstacle()
	}
}

// UpdateStageGaps updates the gap parameters based on current stage
func (om *ObstacleManager) UpdateStageGaps(minGap, maxGap int, stageIndex int) {
	om.minGap = minGap
	om.maxGap = maxGap
	om.currentStage = stageIndex
}

// Draw renders all obstacles
func (om *ObstacleManager) Draw() {
	for _, obstacle := range om.obstacles {
		obstacle.Draw()
	}
}

// GetObstacles returns all active obstacles for collision detection
func (om *ObstacleManager) GetObstacles() []IObstacle {
	return om.obstacles
}

// generateNewObstacle creates a new obstacle based on probabilities
func (om *ObstacleManager) generateNewObstacle() {
	var newObstacle IObstacle

	// 根据当前游戏阶段动态调整障碍物生成逻辑
	maxObstaclesOnScreen := 3
	if om.currentStage <= 1 { // 初始阶段（0-100分）
		maxObstaclesOnScreen = 1 // 前期最多只有1个障碍物在屏幕上
	} else if om.currentStage <= 3 { // 中期阶段（100-600分）
		maxObstaclesOnScreen = 2 // 中期最多有2个障碍物
	}

	// 检查屏幕上已有的障碍物数量，如果超过阈值，延迟生成新障碍物
	if len(om.obstacles) >= maxObstaclesOnScreen {
		// 增加间隔时间，确保障碍物不会过于密集
		om.nextGapTimer = om.maxGap
		return
	}

	r := rand.Float64()
	if r < bigBirdProbability {
		newObstacle = NewBigBird()
	} else if r < bigBirdProbability+birdProbability {
		newObstacle = NewBird()
	} else if r < bigBirdProbability+birdProbability+groupCactusProbability {
		newObstacle = NewGroupCactus()
	} else {
		newObstacle = NewCactus()
	}

	// 添加到障碍物列表
	om.obstacles = append(om.obstacles, newObstacle)

	// 设置下一个障碍物的生成间隔
	// 根据游戏阶段调整间隔生成策略
	var gapMultiplier float64 = 1.0
	if om.currentStage <= 1 {
		// 初始阶段，使用更大的间隔
		gapMultiplier = 1.5
	} else if om.currentStage <= 3 {
		// 中期阶段，使用稍大的间隔
		gapMultiplier = 1.2
	}

	// 使用加权随机，更倾向于选择较大的间隔值
	gapRange := om.maxGap - om.minGap + 1
	weightedGap := om.minGap + int(float64(gapRange)*rand.Float64()*rand.Float64())

	// 应用阶段乘数
	weightedGap = int(float64(weightedGap) * gapMultiplier)

	// 确保不超过最大间隔
	if weightedGap > int(float64(om.maxGap)*1.5) {
		weightedGap = int(float64(om.maxGap) * 1.5)
	}

	om.nextGapTimer = weightedGap
}

// GroupCactus represents a group of connected cacti obstacle
type GroupCactus struct {
	BaseObstacle
}

// NewGroupCactus creates a new group cactus obstacle
func NewGroupCactus() *GroupCactus {
	c := &GroupCactus{}
	c.Reset()
	return c
}

// Reset resets the group cactus position and animation
func (c *GroupCactus) Reset() {
	c.posX = float64(width - 1)
	c.y = height - 2
	c.animFrame = 0
	c.animCounter = 0
}

// Draw renders the group cactus on screen
func (c *GroupCactus) Draw() {
	sprite := groupCactusFrames[c.animFrame]
	h := len(sprite)
	startY := c.y - (h - 1)
	x := int(math.Round(c.posX))
	sprite.Draw(x, startY, termbox.ColorRed, termbox.ColorDefault)
}

// GetSprite returns the current sprite for collision detection
func (c *GroupCactus) GetSprite() Sprite {
	return groupCactusFrames[c.animFrame]
}

// getFrameCount returns the number of animation frames
func (c *GroupCactus) getFrameCount() int {
	return len(groupCactusFrames)
}
