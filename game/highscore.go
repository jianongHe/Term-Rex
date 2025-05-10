package game

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const highScoreFileName = ".term-rex-highscore"

// SaveHighScore 将最高分保存到文件中
func SaveHighScore(score int) error {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("无法获取用户主目录: %v", err)
	}

	// 创建高分文件路径
	highScorePath := filepath.Join(homeDir, highScoreFileName)

	// 将分数转换为字符串并写入文件
	return ioutil.WriteFile(highScorePath, []byte(strconv.Itoa(score)), 0644)
}

// LoadHighScore 从文件中加载最高分
func LoadHighScore() (int, error) {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return 0, fmt.Errorf("无法获取用户主目录: %v", err)
	}

	// 创建高分文件路径
	highScorePath := filepath.Join(homeDir, highScoreFileName)

	// 检查文件是否存在
	if _, err := os.Stat(highScorePath); os.IsNotExist(err) {
		// 文件不存在，返回0分
		return 0, nil
	}

	// 读取文件内容
	data, err := ioutil.ReadFile(highScorePath)
	if err != nil {
		return 0, fmt.Errorf("无法读取高分文件: %v", err)
	}

	// 将内容转换为整数
	scoreStr := strings.TrimSpace(string(data))
	score, err := strconv.Atoi(scoreStr)
	if err != nil {
		return 0, fmt.Errorf("无法解析高分: %v", err)
	}

	return score, nil
}
