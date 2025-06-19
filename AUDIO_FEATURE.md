# Audio System Feature

## Overview
The Term-Rex game now includes a fully functional audio system that provides sound effects for various game events. The audio system is cross-platform and uses both custom audio files and system sounds to enhance the gaming experience.

## Features

### Sound Effects
- **Jump Sound**: Custom MP3 file (`assets/sounds/jump.mp3`) - played when the dinosaur jumps
- **Drop Sound**: Custom MP3 file (`assets/sounds/drop.mp3`) - played when the dinosaur ducks/drops quickly
- **Collision Sound**: System sound - played when the dinosaur hits an obstacle
- **Score Sound**: System sound - played when reaching score milestones

### Controls
- **'m' key**: Toggle audio on/off during gameplay
- Audio state is displayed in the game UI

### Hybrid Audio Approach
The system uses a hybrid approach combining custom audio files with system sounds:

#### Custom Audio Files (Jump & Drop)
- **Jump**: Uses `assets/sounds/jump.mp3` for a more immersive jumping experience
- **Drop**: Uses `assets/sounds/drop.mp3` for ducking/dropping actions
- Files are played using platform-specific audio players

#### System Sounds (Collision & Score)
- **Collision**: Uses system alert sounds for immediate feedback
- **Score**: Uses system notification sounds for achievement feedback

### Cross-Platform Support
The audio system adapts to different operating systems:

#### macOS
- **Custom files**: Uses `afplay` command to play MP3 files
- **System sounds**: Uses built-in system sounds (Sosumi.aiff, Glass.aiff)

#### Linux
- **Custom files**: Tries multiple audio backends in order:
  - PulseAudio (`paplay`)
  - ALSA (`aplay`)
  - mpg123 (`mpg123`)
  - FFmpeg (`ffplay`)
  - VLC (`cvlc`)
- **System sounds**: Uses ALSA system sounds or beep commands

#### Windows
- **Custom files**: Uses PowerShell Media.SoundPlayer or Windows Media Player
- **System sounds**: Uses `rundll32` with system message beeps

#### Other Systems
- Falls back to terminal bell character (`\a`)

## File Structure
```
assets/
└── sounds/
    ├── jump.mp3      # Custom jump sound (used)
    ├── drop.mp3      # Custom drop sound (used)
    ├── collison.mp3  # Available but not used (uses system sound)
    └── score.mp3     # Available but not used (uses system sound)
```

## Configuration
Audio is enabled by default but can be controlled via:
- Config setting: `AudioEnabled = true` in `config.go`
- Runtime toggle: Press 'm' key during gameplay
- Programmatic control: `GetAudioManager().SetEnabled(bool)`
- Custom sounds directory: `GetAudioManager().SetSoundsDirectory(string)`

## Implementation Details

### Architecture
- **Singleton Pattern**: Single `AudioManager` instance manages all audio
- **Hybrid Approach**: Custom files for jump/drop, system sounds for collision/score
- **Non-blocking**: Sound commands run asynchronously to avoid game lag
- **Graceful Degradation**: Falls back to simpler sounds if advanced features unavailable

### Error Handling
- Commands that fail silently fall back to simpler alternatives
- Missing audio files fall back to system sounds or terminal bell
- No crashes if audio systems are unavailable
- Game continues normally even if all audio fails

### Performance
- Minimal overhead when audio is disabled
- Asynchronous sound execution prevents game stuttering
- No audio file preloading - files are played on demand
- Efficient path resolution for audio files

## Usage
The audio system is automatically initialized when the game starts. Players can:
1. Enjoy enhanced sound effects during gameplay (custom jump/drop sounds)
2. Toggle audio on/off using the 'm' key
3. See audio status in the game interface

The system automatically detects and uses the best available audio method for the current platform, providing a rich audio experience with custom sounds for key actions while maintaining system integration for feedback sounds.
