package game

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
)

// 音效类型常量
const (
	SoundJump      = "jump"
	SoundCollision = "collision"
	SoundScore     = "score"
	SoundDrop      = "drop" // 新增：快速下降音效
)

// AudioManager 管理游戏音效
type AudioManager struct {
	enabled   bool
	soundsDir string // Path to sounds directory
}

var (
	audioManager *AudioManager
)

// GetAudioManager 返回单例的音频管理器
func GetAudioManager() *AudioManager {
	if audioManager == nil {
		audioManager = &AudioManager{
			enabled:   AudioEnabled, // Use config value
			soundsDir: "assets/sounds",
		}
	}
	return audioManager
}

// Initialize 初始化音频系统并加载所有音效
func (am *AudioManager) Initialize() error {
	// For terminal-based audio, we don't need to preload files
	// We'll play them on demand
	return nil
}

// PlaySound 播放指定的音效
func (am *AudioManager) PlaySound(name string) {
	if !am.enabled {
		return
	}

	// Use custom audio files for jump and drop, system sounds for others
	switch name {
	case SoundJump:
		am.playCustomSound("jump.mp3")
	case SoundDrop:
		am.playCustomSound("drop.mp3")
	case SoundCollision:
		am.playSystemSound("collision")
	case SoundScore:
		am.playSystemSound("score")
	}
}

// playCustomSound plays a custom audio file from the sounds directory
func (am *AudioManager) playCustomSound(filename string) {
	soundPath := filepath.Join(am.soundsDir, filename)

	switch runtime.GOOS {
	case "darwin": // macOS
		exec.Command("afplay", soundPath).Start()
	case "linux":
		am.tryLinuxAudioCommands(soundPath)
	case "windows":
		// For Windows, try different media players
		am.tryWindowsAudioCommands(soundPath)
	default:
		// Fallback to terminal bell for unsupported systems
		fmt.Print("\a")
	}
}

// tryLinuxAudioCommands tries different Linux audio players for custom files
func (am *AudioManager) tryLinuxAudioCommands(soundPath string) {
	commands := []string{
		"paplay " + soundPath,
		"aplay " + soundPath,
		"mpg123 -q " + soundPath,
		"ffplay -nodisp -autoexit -v quiet " + soundPath,
		"cvlc --play-and-exit --intf dummy " + soundPath,
	}

	for _, cmdStr := range commands {
		if am.executeCommand(cmdStr) {
			return
		}
	}

	// Fallback to terminal bell
	fmt.Print("\a")
}

// tryWindowsAudioCommands tries different Windows audio players for custom files
func (am *AudioManager) tryWindowsAudioCommands(soundPath string) {
	commands := []string{
		"powershell -c \"(New-Object Media.SoundPlayer '" + soundPath + "').PlaySync()\"",
		"start /min wmplayer /close " + soundPath,
	}

	for _, cmdStr := range commands {
		if am.executeCommand(cmdStr) {
			return
		}
	}

	// Fallback to system beep
	exec.Command("rundll32", "user32.dll,MessageBeep", "0x00000040").Start()
}

// playSystemSound plays a system sound based on the OS
func (am *AudioManager) playSystemSound(soundType string) {
	switch runtime.GOOS {
	case "darwin": // macOS
		am.playMacSound(soundType)
	case "linux":
		am.playLinuxSound(soundType)
	case "windows":
		am.playWindowsSound(soundType)
	default:
		// Fallback to terminal bell
		fmt.Print("\a")
	}
}

// playMacSound plays sounds on macOS using afplay or say commands
func (am *AudioManager) playMacSound(soundType string) {
	switch soundType {
	case "collision":
		// Use system sound for collision
		exec.Command("afplay", "/System/Library/Sounds/Sosumi.aiff").Start()
	case "score":
		// Use system sound for score
		exec.Command("afplay", "/System/Library/Sounds/Glass.aiff").Start()
	default:
		// Fallback to terminal bell
		fmt.Print("\a")
	}
}

// playLinuxSound plays sounds on Linux using paplay, aplay, or beep
func (am *AudioManager) playLinuxSound(soundType string) {
	// Try different Linux sound approaches
	switch soundType {
	case "collision":
		am.tryLinuxSystemCommands([]string{
			"paplay /usr/share/sounds/alsa/Side_Left.wav",
			"aplay /usr/share/sounds/alsa/Side_Left.wav",
			"beep -f 200 -l 300",
		})
	case "score":
		am.tryLinuxSystemCommands([]string{
			"paplay /usr/share/sounds/alsa/Front_Right.wav",
			"aplay /usr/share/sounds/alsa/Front_Right.wav",
			"beep -f 1000 -l 200",
		})
	default:
		fmt.Print("\a")
	}
}

// playWindowsSound plays sounds on Windows
func (am *AudioManager) playWindowsSound(soundType string) {
	switch soundType {
	case "collision":
		exec.Command("rundll32", "user32.dll,MessageBeep", "0x00000010").Start()
	case "score":
		exec.Command("rundll32", "user32.dll,MessageBeep", "0x00000000").Start()
	default:
		fmt.Print("\a")
	}
}

// tryLinuxSystemCommands tries multiple Linux sound commands until one works
func (am *AudioManager) tryLinuxSystemCommands(commands []string) {
	for _, cmdStr := range commands {
		if am.executeCommand(cmdStr) {
			return
		}
	}
	// Fallback to terminal bell
	fmt.Print("\a")
}

// executeCommand executes a shell command and returns true if successful
func (am *AudioManager) executeCommand(cmdStr string) bool {
	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Stdout = nil
	cmd.Stderr = nil
	err := cmd.Start()
	return err == nil
}

// SetEnabled 启用或禁用音效
func (am *AudioManager) SetEnabled(enabled bool) {
	am.enabled = enabled
}

// IsEnabled 返回音效是否启用
func (am *AudioManager) IsEnabled() bool {
	return am.enabled
}

// ToggleEnabled 切换音效启用状态
func (am *AudioManager) ToggleEnabled() {
	am.enabled = !am.enabled
}

// SetSoundsDirectory sets the directory where sound files are located
func (am *AudioManager) SetSoundsDirectory(dir string) {
	am.soundsDir = dir
}

// checkSoundAvailability checks if sound commands are available on the system
func (am *AudioManager) checkSoundAvailability() bool {
	switch runtime.GOOS {
	case "darwin":
		// Check if afplay is available
		_, err := exec.LookPath("afplay")
		return err == nil
	case "linux":
		// Check if any Linux sound command is available
		commands := []string{"paplay", "aplay", "beep", "mpg123", "ffplay", "cvlc"}
		for _, cmd := range commands {
			if _, err := exec.LookPath(cmd); err == nil {
				return true
			}
		}
		return false
	case "windows":
		// PowerShell and wmplayer should be available on Windows
		return true
	default:
		// For other systems, assume terminal bell is available
		return true
	}
}
