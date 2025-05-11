# Term-Rex Installation Guide

## Installation Methods

### Direct Download

Download the binary for your platform from [GitHub Releases](https://github.com/jianongHe/Term-Rex/releases).

### Using Docker

```bash
# Pull the image
docker pull jianonghe/term-rex:latest

# Run the game
docker run -it --rm jianonghe/term-rex:latest
```

### Using npm

```bash
# Global installation
npm install -g term-rex

# Run the game
term-rex
```

### Using Homebrew (macOS)

```bash
# Add tap
brew tap jianonghe/term-rex

# Install
brew install term-rex

# Run the game
term-rex
```

## Building from Source

If you want to build Term-Rex from source, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/jianongHe/Term-Rex.git
   cd Term-Rex
   ```

2. Build the project:
   ```bash
   ./build.sh
   ```

3. Run the game:
   ```bash
   ./bin/term-rex
   ```

## System Requirements

- Terminal with Unicode support
- Audio support (optional, for sound effects)

## Troubleshooting

If you encounter issues running the game, try these solutions:

1. Ensure your terminal supports Unicode characters
2. Check if necessary audio libraries are installed (e.g., ALSA)
3. Try running the game in different terminal emulators

If the problem persists, please report the issue on [GitHub Issues](https://github.com/jianongHe/Term-Rex/issues).
