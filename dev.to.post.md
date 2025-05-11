*This is a submission for the [Amazon Q Developer "Quack The Code" Challenge](https://dev.to/challenges/aws-amazon-q-v2025-04-30): That's Entertainment!*

# Term-Rex: How Amazon Q Developer Built My Terminal Dinosaur Game

## What I Built

With the incredible assistance of Amazon Q Developer, I created **Term-Rex** - a terminal-based dinosaur runner game inspired by Chrome's offline T-Rex game, but with enhanced features and gameplay. This retro-style arcade game runs directly in your terminal, bringing nostalgic gaming joy to developers everywhere.

What makes Term-Rex special is that **Amazon Q Developer handled approximately 90% of the development work**. From designing the game architecture to implementing complex game mechanics like collision detection, obstacle generation, and animation systems - Amazon Q was the driving force behind this project.

The game features:
- An animated dinosaur character that can jump and duck
- Three types of obstacles: cacti, small birds, and big birds
- Day/night cycle with changing visuals
- Cloud animations and ground decorations for visual depth
- Score tracking with milestone celebrations
- High score persistence between sessions
- Sound effects for jumping, collisions, and milestones
- Multiple game stages with increasing difficulty
- Cross-platform support (macOS, Linux, Windows)
- Docker containerization for universal play

## Demo

![Term-Rex Game Demo](https://raw.githubusercontent.com/jianongHe/Term-Rex/main/assets/demo.gif)

### How to Play

Term-Rex can be installed and played through multiple methods:

**Homebrew (macOS/Linux):**
```bash
brew install jianongHe/tap/term-rex
```

**NPM:**
```bash
npm install -g term-rex
```

**Scoop (Windows):**
```powershell
scoop bucket add jianongHe https://github.com/jianongHe/scoop-bucket.git
scoop install term-rex
```

**Docker:**
```bash
docker run -it --rm jianonghe/term-rex
```

**Game Controls:**
- Space/Up Arrow: Jump
- Down Arrow: Duck
- P: Pause/Resume
- M: Toggle sound
- R: Restart (after game over)
- Q/Esc: Quit

## Code Repository

<a href="https://github.com/jianongHe/Term-Rex">
  <div style="display: flex; align-items: center; margin-bottom: 16px;">
    <img src="https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png" width="32" height="32" style="margin-right: 8px;">
    <span style="font-weight: bold;">jianongHe/Term-Rex</span>
  </div>
</a>

## How I Used Amazon Q Developer

My journey with Term-Rex showcases the transformative power of Amazon Q Developer in game development. Here's how Amazon Q revolutionized my workflow:

### 1. Architecture Design & Project Setup

I started with a simple prompt: "Help me create a terminal-based dinosaur runner game in Go." Amazon Q immediately:
- Designed a comprehensive game architecture with proper separation of concerns
- Set up the project structure with all necessary files and modules
- Implemented the game loop with proper timing and frame rate control
- Created a clean, maintainable codebase following Go best practices

### 2. Game Mechanics Implementation

Amazon Q Developer implemented sophisticated game mechanics that would have taken me weeks to code:

**Obstacle System:**
```go
// Amazon Q designed this elegant object-oriented obstacle system
type IObstacle interface {
    Update(float64)
    Draw(*termbox.Buffer)
    GetX() float64
    GetWidth() int
    GetHeight() int
    GetRow() int
    IsVisible() bool
}

// Base obstacle with shared functionality
type BaseObstacle struct {
    x        float64
    width    int
    height   int
    row      int
    visible  bool
    frames   []string
    frameIdx int
}
```

**Animation System:**
Amazon Q created a flexible animation system that handles sprite transitions for the dinosaur and obstacles, making the game visually engaging.

**Collision Detection:**
The collision detection algorithm Amazon Q developed is both efficient and accurate, handling different hitboxes for various obstacle types and dinosaur states (running, jumping, ducking).

### 3. Visual Enhancements

When I asked Amazon Q to improve the game's visual appeal, it:
- Added a cloud system with different cloud types moving at varying speeds
- Implemented ground decorations for visual depth
- Created a day/night cycle that changes the game's color scheme
- Designed ASCII art for all game elements

### 4. Cross-Platform Distribution

Amazon Q helped me make Term-Rex available everywhere:
- Created installation scripts for Homebrew, NPM, and Scoop
- Set up GitHub Actions workflows for automated releases
- Implemented Docker containerization with multi-architecture support
- Ensured the game runs consistently across different terminal types

### 5. Continuous Improvement

Throughout development, Amazon Q:
- Refactored code for better performance and maintainability
- Fixed bugs and edge cases I hadn't even considered
- Suggested feature enhancements that made the game more engaging
- Helped synchronize version information across all project files

### 6. Documentation & Community Engagement

Amazon Q even helped with:
- Writing comprehensive documentation and README files
- Creating this very submission for the "Quack The Code" Challenge
- Generating engaging GIFs and screenshots for promotion

## The Amazon Q Advantage

What impressed me most about working with Amazon Q Developer was:

1. **Contextual Understanding**: Amazon Q understood the entire codebase and could make targeted improvements to specific components without breaking others.

2. **Problem-Solving**: When I encountered issues with CI/CD pipelines or version synchronization, Amazon Q quickly diagnosed and resolved them.

3. **Code Quality**: The code Amazon Q produced was clean, well-documented, and followed best practices - no spaghetti code or technical debt.

4. **Learning Opportunity**: Working with Amazon Q taught me advanced Go programming techniques and game development patterns I wouldn't have discovered on my own.

## Conclusion

Term-Rex demonstrates how Amazon Q Developer can transform game development. With minimal input from me, Amazon Q created a fully-featured, cross-platform game that's both entertaining and educational. The project showcases Amazon Q's ability to handle complex, creative coding tasks while maintaining high code quality.

If Amazon Q can build a complete game with sophisticated mechanics, animations, and cross-platform support, imagine what it can do for your next project!

<!-- ⚠️ By submitting this entry, you agree to receive communications from AWS regarding products, services, events, and special offers. You can unsubscribe at any time. Your information will be handled in accordance with [AWS's Privacy Policy](https://aws.amazon.com/privacy/). Additionally, your submission and project may be publicly featured on AWS's social media channels or related promotional materials. -->
