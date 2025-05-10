package game

import (
	"github.com/nsf/termbox-go"
	"math/rand"
)

// Cloud represents a decorative cloud in the sky
type Cloud struct {
	x         int
	y         int
	width     int
	speed     float64
	posX      float64
	cloudType int
	color     termbox.Attribute
}

// CloudManager manages multiple clouds in the sky
type CloudManager struct {
	clouds    []*Cloud
	maxClouds int
	colors    []termbox.Attribute
}

// NewCloudManager creates a new cloud manager with initial clouds
func NewCloudManager() *CloudManager {
	cm := &CloudManager{
		maxClouds: cloudMinCount + rand.Intn(cloudMaxCount-cloudMinCount+1),
		colors:    []termbox.Attribute{termbox.ColorWhite},
	}

	// Create initial clouds with good spacing
	availablePositions := make([]bool, width)
	for i := 0; i < width; i++ {
		availablePositions[i] = true
	}

	for i := 0; i < cm.maxClouds; i++ {
		// Find a position that doesn't overlap with existing clouds
		cloudType := rand.Intn(len(cloudSprites))
		cloudWidth := len(cloudSprites[cloudType][0])

		// Try to find a position with enough space
		var startPos int
		for attempts := 0; attempts < 20; attempts++ { // Limit attempts to avoid infinite loop
			startPos = rand.Intn(width)
			hasSpace := true

			// Check if there's enough space (with buffer)
			buffer := cloudWidth + cloudBufferSpace
			for j := startPos - buffer/2; j < startPos+buffer/2; j++ {
				if j >= 0 && j < width && !availablePositions[j] {
					hasSpace = false
					break
				}
			}

			if hasSpace {
				break
			}
		}

		// Mark this position as used
		for j := startPos - cloudWidth/2; j < startPos+cloudWidth/2; j++ {
			if j >= 0 && j < width {
				availablePositions[j] = false
			}
		}

		// Create cloud at this position with random color
		cm.clouds = append(cm.clouds, &Cloud{
			posX:      float64(startPos),
			x:         startPos, // Explicitly set x to match posX initially
			y:         cloudMinHeight + rand.Intn(cloudMaxHeight-cloudMinHeight+1),
			width:     cloudWidth,
			speed:     cloudMinSpeed + rand.Float64()*(cloudMaxSpeed-cloudMinSpeed),
			cloudType: cloudType,
			color:     cm.colors[rand.Intn(len(cm.colors))],
		})
	}

	return cm
}

// createNewCloud creates a new cloud at the right edge with proper spacing
func (cm *CloudManager) createNewCloud() *Cloud {
	cloudType := rand.Intn(len(cloudSprites))
	cloudWidth := len(cloudSprites[cloudType][0])

	// Check if there's already a cloud near the right edge
	hasNearbyCloud := false

	for _, cloud := range cm.clouds {
		if cloud.x > width-cloudRightEdgeBuffer {
			hasNearbyCloud = true
			break
		}
	}

	// If there's a cloud near the right edge, add extra spacing
	extraSpace := 0
	if hasNearbyCloud {
		extraSpace = cloudMinExtraSpace + rand.Intn(cloudMaxExtraSpace-cloudMinExtraSpace+1)
	}

	// Position the new cloud completely off-screen to the right
	newX := width + extraSpace

	return &Cloud{
		posX:      float64(newX),
		x:         newX, // Explicitly set x to match posX initially
		y:         cloudMinHeight + rand.Intn(cloudMaxHeight-cloudMinHeight+1),
		width:     cloudWidth,
		speed:     (cloudMinSpeed + rand.Float64()*(cloudMaxSpeed-cloudMinSpeed)),
		cloudType: cloudType,
		color:     cm.colors[rand.Intn(len(cm.colors))],
	}
}

// Update moves all clouds and cycles them when they move off-screen
func (cm *CloudManager) Update() {
	for i, cloud := range cm.clouds {
		// 使用速度因子调整云朵移动速度，保持与原始帧率下相同的视觉速度
		cloud.posX -= cloud.speed * speedFactor
		cloud.x = int(cloud.posX)

		// If cloud has moved completely off-screen to the left
		if cloud.x+cloud.width < -5 { // Add extra buffer to ensure it's fully off-screen
			// Create a new cloud at the right edge of the screen with proper spacing
			cm.clouds[i] = cm.createNewCloud()
		}
	}
}

// Draw renders all clouds on the screen
func (cm *CloudManager) Draw() {
	for _, cloud := range cm.clouds {
		// Skip drawing if the cloud is completely off-screen
		if cloud.x+cloud.width < 0 || cloud.x > width {
			continue
		}

		// Draw all clouds regardless of game state or ground extension
		sprite := cloudSprites[cloud.cloudType]
		for y, line := range sprite {
			for x, ch := range line {
				// Only draw non-space characters that are within screen bounds
				if ch != ' ' && cloud.x+x >= 0 && cloud.x+x < width {
					termbox.SetCell(cloud.x+x, cloud.y+y, ch, termbox.ColorWhite, termbox.ColorDefault)
				}
			}
		}
	}
}
