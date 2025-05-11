# Changelog

All notable changes to Term-Rex will be documented in this file.

## [0.1.8] - 2025-05-11

### Added
- Enhanced cloud system with improved distribution and movement patterns
- Optimized obstacle generation algorithm for better gameplay balance

### Changed
- Refined big bird obstacle behavior for more challenging gameplay
- Improved visual consistency across different terminal sizes
- Removed debug logs for cleaner production code

### Fixed
- Fixed potential race condition in obstacle management
- Improved performance on systems with limited resources

## [0.1.7] - 2025-05-11

### Changed
- Removed unused Makefile as CI uses GoReleaser directly
- Simplified project structure

## [0.1.6] - 2025-05-11

### Fixed
- Fixed installation script to use correct version number
- Updated package.json version to match application version
- Synchronized version numbers across all project files

## [0.1.5] - 2025-05-11

### Fixed
- Fixed CI/CD pipeline to prevent dirty state in release process
- Improved build stability for cross-platform releases

## [0.1.4] - 2025-05-11

### Added
- Added big bird obstacle for increased gameplay variety
- Implemented cloud system with different speeds and positions
- Added ground decorations for improved visual appeal

### Changed
- Refactored obstacle system using object-oriented design
- Improved game scene with layered visual elements
- Optimized collision detection for different obstacle types

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
