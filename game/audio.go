package game

// This is a stub implementation of the audio system
// All audio functionality is disabled to troubleshoot crashes

// 音效类型常量
const (
	SoundJump      = "jump"
	SoundCollision = "collision"
	SoundScore     = "score"
	SoundDrop      = "drop" // 新增：快速下降音效
)

// AudioManager 管理游戏音效
type AudioManager struct {
	enabled bool
}

var (
	audioManager *AudioManager
)

// GetAudioManager 返回单例的音频管理器
func GetAudioManager() *AudioManager {
	if audioManager == nil {
		audioManager = &AudioManager{
			enabled: false,
		}
	}
	return audioManager
}

// Initialize 初始化音频系统并加载所有音效
func (am *AudioManager) Initialize() error {
	// Do nothing - audio is disabled
	return nil
}

// PlaySound 播放指定的音效
func (am *AudioManager) PlaySound(name string) {
	// Do nothing - audio is disabled
}

// SetEnabled 启用或禁用音效
func (am *AudioManager) SetEnabled(enabled bool) {
	am.enabled = false // Always keep disabled
}

// IsEnabled 返回音效是否启用
func (am *AudioManager) IsEnabled() bool {
	return false // Always return false
}
