# Changelog

All notable changes to Term-Rex will be documented in this file.

## [0.1.3] - 2025-05-11

### Added
- Added pause functionality with 'P' key
- Improved pause UI with clear visual indicators
- Updated README to reflect pause feature

### Changed
- Optimized game loop to handle pause state properly
- Clouds continue to move during pause for visual effect
- Refactored input handling for better code organization

## [0.1.2] - 2025-05-11

### Added
- Added npm package support for easier installation via npm
- Updated installation scripts to match current version

## [0.1.1] - 2025-05-11

### Changed
- Moved dinosaur position 2 cells to the right for better gameplay experience
- Slowed down score accumulation rate for a more gradual difficulty curve
- Optimized README documentation for clearer installation and usage instructions

## [0.1.0] - 2025-05-11

### Added
- Initial release of Term-Rex
- Terminal-based dinosaur runner game inspired by Chrome's offline game
- Features:
  - Control a dinosaur running through a desert, avoiding obstacles
  - Three types of obstacles: cacti, small birds, and big birds
  - Day/night cycle that changes every 60 seconds
  - Score tracking with milestone celebrations
  - High score persistence between game sessions
  - Sound effects for jumping, collision, and scoring milestones
  - Multiple game stages with increasing difficulty
  - Visual elements including animated clouds and ground decorations
- Controls:
  - Space/Up Arrow: Jump
  - Down Arrow: Duck
  - p: Pause/Unpause
  - q/Esc: Quit
  - r: Restart after game over
  - m: Toggle sound effects
