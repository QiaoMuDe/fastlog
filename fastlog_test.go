package gitee.com/MM-Q/fastlog.git

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestLoggerConfig 测试 LoggerConfig 的默认配置生成是否正确
func TestLoggerConfig(t *testing.T) {
	expectedConfig := LoggerConfig{
		LogDirName:     "logs",
		LogFileName:    "app.log",
		LogPath:        filepath.Join("logs", "app.log"), // 使用 filepath.Join 生成路径
		PrintToConsole: true,
		LogLevel:       Info,
		ChanIntSize:    1000,
		BufferKbSize:   1024,
	}
	actualConfig := DefaultConfig("logs", "app.log")
	if actualConfig != expectedConfig {
		t.Errorf("DefaultConfig() = %v, want %v", actualConfig, expectedConfig)
	}
}

// TestNewLogger 测试 NewLogger 是否能正确创建 Logger 实例
func TestNewLogger(t *testing.T) {
	config := DefaultConfig("logs", "test.log")
	logger, err := NewLogger(config)
	if err != nil {
		t.Errorf("NewLogger() error = %v", err)
	}
	if logger == nil {
		t.Errorf("NewLogger() = nil, want non-nil Logger")
	}
	// 关闭 Logger 以清理资源
	logger.Close()
	// 删除测试日志文件
	os.Remove(config.LogPath)
}

// TestLogLevels 测试不同日志级别的记录是否正确
func TestLogLevels(t *testing.T) {
	config := DefaultConfig("logs", "test.log")
	config.LogLevel = Debug
	logger, err := NewLogger(config)
	if err != nil {
		t.Errorf("NewLogger() error = %v", err)
	}
	// 记录不同级别的日志
	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Warn("This is a warn message")
	logger.Error("This is an error message")
	logger.Success("This is a success message")
	// 等待一段时间，确保日志写入完成
	time.Sleep(1 * time.Second)
	// 关闭 Logger
	logger.Close()
	// 读取日志文件内容进行验证
	data, err := os.ReadFile(config.LogPath)
	if err != nil {
		t.Errorf("ReadFile() error = %v", err)
	}
	logContent := string(data)
	expectedLevels := []string{"DEBUG", "INFO", "WARN", "ERROR", "SUCCESS"}
	for _, level := range expectedLevels {
		if !strings.Contains(logContent, level) {
			t.Errorf("Log file does not contain %s level log", level)
		}
	}
	// 删除测试日志文件
	os.Remove(config.LogPath)
}

// TestLogLevelFiltering 测试日志级别过滤功能
func TestLogLevelFiltering(t *testing.T) {
	config := DefaultConfig("logs", "test.log")
	config.LogLevel = Warn
	logger, err := NewLogger(config)
	if err != nil {
		t.Errorf("NewLogger() error = %v", err)
	}
	// 记录不同级别的日志
	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Warn("This is a warn message")
	logger.Error("This is an error message")
	logger.Success("This is a success message")
	// 等待一段时间，确保日志写入完成
	time.Sleep(2 * time.Second) // 增加等待时间
	// 关闭 Logger
	logger.Close()
	// 读取日志文件内容进行验证
	data, err := os.ReadFile(config.LogPath)
	if err != nil {
		t.Errorf("ReadFile() error = %v", err)
	}
	logContent := string(data)
	expectedLevels := []string{"WARN", "ERROR", "SUCCESS"}
	unexpectedLevels := []string{"DEBUG", "INFO"}
	for _, level := range expectedLevels {
		if !strings.Contains(logContent, level) {
			t.Errorf("Log file does not contain %s level log", level)
		}
	}
	for _, level := range unexpectedLevels {
		if strings.Contains(logContent, level) {
			t.Errorf("Log file contains unexpected %s level log", level)
		}
	}
	// 删除测试日志文件
	os.Remove(config.LogPath)
}

// TestLoggerClose 测试 Logger 的 Close 方法是否能正确关闭资源
func TestLoggerClose(t *testing.T) {
	config := DefaultConfig("logs", "test.log")
	logger, err := NewLogger(config)
	if err != nil {
		t.Errorf("NewLogger() error = %v", err)
	}
	// 关闭 Logger
	logger.Close()
	// 尝试再次写入日志，验证是否失败
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Recovered from panic: %v", r)
			}
		}()
		logger.Info("This is a test message after closing logger")
	}()
	// 等待一段时间，确保 goroutine 运行完成
	time.Sleep(1 * time.Second)
	// 读取日志文件内容进行验证
	data, err := os.ReadFile(config.LogPath)
	if err != nil {
		t.Errorf("ReadFile() error = %v", err)
	}
	logContent := string(data)
	if strings.Contains(logContent, "This is a test message after closing logger") {
		t.Errorf("Logger did not close properly, log was written after closing")
	}
	// 删除测试日志文件
	os.Remove(config.LogPath)
}
