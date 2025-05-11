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
	// 使用有效宽度而不是实际宽度
	effectiveWidth := math.Min(float64(width), float64(maxEffectiveWidth))
	c.posX = effectiveWidth
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
	// 使用有效宽度而不是实际宽度
	effectiveWidth := math.Min(float64(width), float64(maxEffectiveWidth))
	b.posX = effectiveWidth
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
	// 使用有效宽度而不是实际宽度
	effectiveWidth := math.Min(float64(width), float64(maxEffectiveWidth))
	b.posX = effectiveWidth
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
	} else {
		// 检查最右边的障碍物是否已经进入屏幕足够距离
		rightmostX := float64(-100) // 初始值设为屏幕外
		for _, obs := range om.obstacles {
			x, _ := obs.GetPosition()
			if x > rightmostX {
				rightmostX = x
			}
		}

		// 如果最右边的障碍物已经进入屏幕一定距离，可以考虑生成新障碍物
		// 这个距离是基于有效屏幕宽度动态计算的
		effectiveWidth := math.Min(float64(width), float64(maxEffectiveWidth))
		entryThreshold := effectiveWidth * 0.7 // 有效宽度的70%
		if rightmostX < entryThreshold {
			// 有一定概率立即生成新障碍物，而不等待计时器
			// 概率随着屏幕宽度增加而增加，但基于有效宽度
			effectiveWidthFactor := effectiveWidth / 80.0
			spawnChance := 0.1 * math.Min(effectiveWidthFactor, 3.0) // 最高30%概率
			if rand.Float64() < spawnChance {
				om.generateNewObstacle()
			}
		}
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

	// Select obstacle type based on probabilities
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

	// Add to obstacle list
	om.obstacles = append(om.obstacles, newObstacle)

	// Calculate gap for next obstacle
	// The gap is measured in frames (how many update cycles before generating the next obstacle)

	// Base gap multiplier - higher value means larger gaps between obstacles
	// Increasing from 0.5 to 0.8 to make obstacles less frequent
	var gapMultiplier float64 = 0.8

	// Adjust gap based on screen width
	// For wider screens, we need proportionally larger gaps
	effectiveWidth := math.Min(float64(width), float64(maxEffectiveWidth))
	effectiveWidthFactor := effectiveWidth / 80.0

	// For wider screens, increase the gap proportionally
	// This ensures consistent difficulty regardless of screen width
	if effectiveWidthFactor > 1.0 {
		// Linear scaling for wider screens (instead of reducing gaps)
		gapMultiplier *= effectiveWidthFactor * 0.8
	}

	// Select a random gap value between min and max for current stage
	gapRange := om.maxGap - om.minGap + 1
	baseGap := om.minGap + rand.Intn(gapRange)

	// Apply the multiplier to get final gap
	finalGap := int(float64(baseGap) * gapMultiplier)

	// Ensure minimum reasonable gap
	minAllowedGap := 10 // Increased from 3 to 10 for better spacing
	if finalGap < minAllowedGap {
		finalGap = minAllowedGap
	}

	// Set timer for next obstacle generation
	om.nextGapTimer = finalGap
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
	// 使用有效宽度而不是实际宽度
	effectiveWidth := math.Min(float64(width), float64(maxEffectiveWidth))
	c.posX = effectiveWidth
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
