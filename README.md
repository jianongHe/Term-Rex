# Term-Rex: Command-Line Dino Runner

A terminal-based dinosaur runner game inspired by the Chrome offline game, implemented in Golang using the termbox-go library.

## Game Features

- Control a dinosaur running through a desert, avoiding obstacles by jumping or ducking
- Visually rich environment with animated clouds and ground decorations
- Progressive difficulty with multiple game stages based on score
- Three types of obstacles with unique movement patterns:
  - Single cacti (ground level)
  - Group cacti (wider ground obstacles)
  - Small birds (mid-air obstacles)
  - Big birds (higher flying, larger obstacles)
- Day/night cycle that changes every 60 seconds
- Score tracking with milestone celebrations (blinking score)
- High score persistence between game sessions
- Sound effects for jumping, collision, and scoring milestones
- Smooth stage transitions with gradually increasing difficulty

## Controls

- **Space** or **Up Arrow**: Jump
- **Down Arrow**: Duck (crouch to avoid flying obstacles)
- **p**: Pause/Unpause the game
- **q** or **Esc**: Quit the game
- **Ctrl+C**: Quit the game (standard terminal interrupt)
- **r**: Restart after game over
- **m**: Toggle sound effects on/off

## Requirements

- Go (1.16 or higher recommended)
- Terminal with Unicode support
- Audio support (optional, for sound effects)

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/Term-Rex.git
cd Term-Rex

# Install dependencies
go get github.com/nsf/termbox-go
go get github.com/faiface/beep
```

## How to Run

```bash
go run main.go
```

## Game Rules

- The dinosaur automatically runs forward
- Obstacles (cacti, group cacti, small birds, and big birds) move from right to left
- Press space or up arrow to jump over ground obstacles
- Press down arrow to duck under flying obstacles
- Jumping has a fixed height and duration (no double jumping)
- Colliding with any obstacle ends the game
- Your score increases the longer you survive
- The game gets progressively harder over time with:
  - Increased obstacle speed
  - Higher probability of more challenging obstacles
  - Shorter intervals between obstacles

## Game Stages

The game features multiple difficulty stages based on your score:

1. **Beginner** (0-100 points): Slow speed, mostly single cacti
2. **Intermediate** (100-300 points): Medium speed, introduction of birds
3. **Advanced** (300-600 points): Faster speed, more birds and group cacti
4. **Expert** (600-1000 points): High speed, frequent birds and big birds
5. **Master** (1000+ points): Maximum difficulty with all obstacle types

## Visual Elements

- **Animated Dinosaur**: Running animation with different frames for standing and ducking
- **Cloud System**: Multiple cloud shapes floating at different heights and speeds
- **Ground Decorations**: Random decorative elements below the ground line
- **Obstacle Variety**: Four distinct obstacle types with unique animations
- **Score Display**: Dynamic score counter with blinking effect at milestones

## Code Structure

The game is built using object-oriented design with the following main components:

- **Game**: Main game controller that manages the game loop and state
- **Dino**: Represents the dinosaur character with jumping and ducking mechanics
- **ObstacleManager**: Manages the creation and updating of obstacles
- **CloudManager**: Controls the cloud system for visual enhancement
- **AudioManager**: Handles sound effects for game events
- **StageManager**: Controls difficulty progression based on score

## Implementation Details

- Fixed 24 FPS update loop for consistent gameplay
- Non-blocking input handling for responsive controls
- Interface-based design for obstacle types enabling easy extension
- Terminal display is properly restored after exit
- High score tracking between sessions using file persistence
- Sound effects using the beep library with MP3 and WAV support
- Smooth difficulty transitions between game stages

## Credits

- Inspired by Chrome's T-Rex Runner game
- Built with [termbox-go](https://github.com/nsf/termbox-go) for terminal rendering
- Sound effects powered by [beep](https://github.com/faiface/beep) library
