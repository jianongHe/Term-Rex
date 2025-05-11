# Term-Rex: Command-Line Dino Runner

[![Release](https://img.shields.io/github/v/release/jianongHe/term-rex)](https://github.com/jianongHe/term-rex/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/jianongHe/term-rex)](https://goreportcard.com/report/github.com/jianongHe/term-rex)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A terminal-based dinosaur runner game inspired by the Chrome offline game, implemented in Golang using the termbox-go library.

![Term-Rex Game Demo](https://github.com/jianongHe/term-rex/raw/main/assets/demo.gif)

## Installation

### macOS / Linux

```bash
# Install with Homebrew
brew tap jianongHe/tap
brew install term-rex
```

### Windows

```powershell
# Install with Scoop
scoop bucket add jianongHe https://github.com/jianongHe/scoop-bucket.git
scoop install term-rex
```

### Direct Download

Download the pre-built binary for your system from the [releases page](https://github.com/jianongHe/term-rex/releases/latest).

### Building from Source

```bash
# Clone the repository
git clone https://github.com/jianongHe/term-rex.git
cd term-rex

# Build the game
go build -o term-rex
```

## Game Controls

| Key | Action |
|-----|--------|
| <kbd>Space</kbd> / <kbd>↑</kbd> | Jump |
| <kbd>↓</kbd> | Duck |
| <kbd>p</kbd> | Pause/Resume |
| <kbd>m</kbd> | Toggle sound |
| <kbd>r</kbd> | Restart (after game over) |
| <kbd>q</kbd> / <kbd>Esc</kbd> | Quit |

## Game Features

- **Challenging Gameplay**: Navigate your dinosaur through a desert landscape, avoiding obstacles
- **Multiple Obstacle Types**: Dodge cacti, small birds, and large birds with unique movement patterns
- **Progressive Difficulty**: Game speed and obstacle frequency increase as your score grows
- **Visual Effects**: Animated clouds, day/night cycle, and ground decorations
- **Sound Effects**: Optional audio feedback for jumping, collisions, and milestones
- **High Score Tracking**: Your best score is saved between game sessions

## Game Mechanics

### Obstacles

- **Single Cacti**: Ground-level obstacles that require jumping over
- **Group Cacti**: Wider ground obstacles that require precise jumping
- **Small Birds**: Mid-air obstacles that require ducking under
- **Big Birds**: Higher flying, larger obstacles that require careful timing

### Difficulty Progression

The game features multiple difficulty stages based on your score:

1. **Beginner** (0-100 points): Slow speed, mostly single cacti
2. **Intermediate** (100-300 points): Medium speed, introduction of birds
3. **Advanced** (300-600 points): Faster speed, more birds and group cacti
4. **Expert** (600-1000 points): High speed, frequent birds and big birds
5. **Master** (1000+ points): Maximum difficulty with all obstacle types

## System Requirements

- Any terminal with Unicode support
- Minimum terminal size: 80×24 characters
- Audio support (optional, for sound effects)

## Uninstallation

### Homebrew (macOS and Linux)

```bash
brew uninstall term-rex
brew untap jianongHe/tap  # Optional: remove the tap repository
```

### Scoop (Windows)

```powershell
scoop uninstall term-rex
scoop bucket rm jianongHe  # Optional: remove the bucket
```

## Contributing

Contributions are welcome! Feel free to submit issues or pull requests.

## Credits

- Inspired by Chrome's T-Rex Runner game
- Built with [termbox-go](https://github.com/nsf/termbox-go) for terminal rendering
- Sound effects powered by [beep](https://github.com/faiface/beep) library

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
