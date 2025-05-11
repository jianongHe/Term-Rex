package game

import (
	"os"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
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
	enabled      bool
	soundBuffers map[string]*beep.Buffer
	mutex        sync.Mutex
}

var (
	audioManager *AudioManager
	once         sync.Once
)

// GetAudioManager 返回单例的音频管理器
func GetAudioManager() *AudioManager {
	once.Do(func() {
		audioManager = &AudioManager{
			enabled:      true,
			soundBuffers: make(map[string]*beep.Buffer),
		}
		// 初始化音频系统
		err := audioManager.Initialize()
		if err != nil {
			// 如果初始化失败，禁用音效但不中断游戏
			audioManager.enabled = false
		}
	})
	return audioManager
}

// Initialize 初始化音频系统并加载所有音效
func (am *AudioManager) Initialize() error {
	// 初始化音频播放器
	sampleRate := beep.SampleRate(44100)
	err := speaker.Init(sampleRate, sampleRate.N(time.Millisecond*10))
	if err != nil {
		return err
	}

	// 尝试加载音效，但不要因为加载失败而中断游戏
	am.tryLoadSound(SoundJump, "assets/sounds/jump.mp3")
	am.tryLoadSound(SoundCollision, "assets/sounds/collision.wav")
	am.tryLoadSound(SoundScore, "assets/sounds/score.wav")
	am.tryLoadSound(SoundDrop, "assets/sounds/drop.mp3")

	// 如果没有成功加载任何音效，禁用音频系统
	if len(am.soundBuffers) == 0 {
		return nil
	}

	return nil
}

// tryLoadSound 尝试加载单个音效文件，如果失败则忽略错误
func (am *AudioManager) tryLoadSound(name, path string) {
	// 检查文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return
	}

	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	var streamer beep.StreamSeekCloser
	var format beep.Format

	// 根据文件扩展名选择解码器
	if len(path) > 4 && path[len(path)-4:] == ".mp3" {
		streamer, format, err = mp3.Decode(f)
	} else {
		streamer, format, err = wav.Decode(f)
	}

	if err != nil {
		return
	}
	defer streamer.Close()

	// 创建缓冲区以便重复播放
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)

	am.mutex.Lock()
	am.soundBuffers[name] = buffer
	am.mutex.Unlock()
}

// PlaySound 播放指定的音效
func (am *AudioManager) PlaySound(name string) {
	if !am.enabled {
		return
	}

	am.mutex.Lock()
	buffer, exists := am.soundBuffers[name]
	am.mutex.Unlock()

	if !exists {
		return
	}

	// 创建新的流以便播放
	streamer := buffer.Streamer(0, buffer.Len())
	speaker.Play(streamer)
}

// SetEnabled 启用或禁用音效
func (am *AudioManager) SetEnabled(enabled bool) {
	am.enabled = enabled
}

// IsEnabled 返回音效是否启用
func (am *AudioManager) IsEnabled() bool {
	return am.enabled
}
