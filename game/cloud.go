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
}

// CloudManager manages multiple clouds in the sky
type CloudManager struct {
	clouds    []*Cloud
	maxClouds int
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

// NewCloudManager creates a new cloud manager with initial clouds
func NewCloudManager() *CloudManager {
	cm := &CloudManager{
		maxClouds: 3 + rand.Intn(3), // 3-5 clouds
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
			buffer := cloudWidth + 5 // Add buffer space between clouds
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

		// Create cloud at this position
		cm.clouds = append(cm.clouds, &Cloud{
			posX:      float64(startPos),
			y:         2 + rand.Intn(4), // Random height in the sky (rows 2-5)
			width:     cloudWidth,
			speed:     0.2 + rand.Float64()*0.3, // Random speed between 0.2-0.5
			cloudType: cloudType,
		})
	}

	return cm
}

// createNewCloud creates a new cloud at the right edge with proper spacing
func (cm *CloudManager) createNewCloud() *Cloud {
	cloudType := rand.Intn(len(cloudSprites))
	cloudWidth := len(cloudSprites[cloudType][0])

	// Check if there's already a cloud near the right edge
	rightEdgeBuffer := 15 // Buffer space to check near right edge
	hasNearbyCloud := false

	for _, cloud := range cm.clouds {
		if cloud.x > width-rightEdgeBuffer {
			hasNearbyCloud = true
			break
		}
	}

	// If there's a cloud near the right edge, add extra spacing
	extraSpace := 0
	if hasNearbyCloud {
		extraSpace = 10 + rand.Intn(15) // Add 10-25 extra spaces
	}

	return &Cloud{
		posX:      float64(width + extraSpace),
		y:         2 + rand.Intn(4), // Random height in the sky (rows 2-5)
		width:     cloudWidth,
		speed:     0.2 + rand.Float64()*0.3, // Random speed between 0.2-0.5
		cloudType: cloudType,
	}
}

// Update moves all clouds and cycles them when they move off-screen
func (cm *CloudManager) Update() {
	for i, cloud := range cm.clouds {
		cloud.posX -= cloud.speed
		cloud.x = int(cloud.posX)

		// If cloud has moved completely off-screen to the left
		if cloud.x+cloud.width < 0 {
			// Create a new cloud at the right edge of the screen with proper spacing
			cm.clouds[i] = cm.createNewCloud()
		}
	}
}

// Draw renders all clouds on the screen
func (cm *CloudManager) Draw() {
	for _, cloud := range cm.clouds {
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
