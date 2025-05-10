# Term-Rex: Command-Line Dino Runner

A terminal-based dinosaur runner game inspired by the Chrome offline game, implemented in Golang using the curses library.

## Game Features

- Control a dinosaur running through a desert, avoiding obstacles by jumping
- Day/night cycle that changes every 60 seconds
- Progressive difficulty that increases over time
- Score tracking based on survival time
- High score persistence between game sessions
- Three types of obstacles: cacti, small birds, and big birds

## Controls

- **Space** or **Up Arrow**: Jump
- **p**: Pause/Unpause the game
- **q**: Quit the game
- **r**: Restart after game over

## Requirements

- Go

## How to Run

```bash
go run main.go
```

## Game Rules

- The dinosaur automatically runs forward
- Obstacles (cacti, small birds, and big birds) move from right to left
- Press space or up arrow to jump over obstacles
- Jumping has a fixed height and duration (no double jumping)
- Colliding with an obstacle ends the game
- Your score increases the longer you survive
- The game gets progressively harder over time

## Code Structure

The game is built using object-oriented design with the following main components:

- `Game`: Main game controller that manages the game loop and state
- `Player`: Represents the dinosaur character with jumping mechanics
- `Obstacle`: Represents obstacles that the player must avoid

## Implementation Details

- Fixed 24 FPS update loop
- Non-blocking input handling
- Terminal display is properly restored after exit
- High score tracking between sessions
