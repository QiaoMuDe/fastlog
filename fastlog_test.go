package fastlog

import (
	"os"
	"testing"
	"time"
)

// 测试日志轮转功能
func TestLogRotate(t *testing.T) {
	// 设置测试环境
	logDir := "./test_logs"
	logFile := "test.log"

	// 确保测试目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatalf("创建日志目录失败: %v", err)
	}

	// 清理旧的日志文件
	_ = os.RemoveAll(logDir)

	// 创建 Logger 配置
	config := NewConfig(logDir, logFile)
	config.EnableLogRotation = true
	config.LogMaxSize = "1KB"         // 设置日志文件大小限制为 1KB，便于测试
	config.LogRetentionCount = 2      // 设置保留的日志文件数量为 2
	config.RotationInterval = 1       // 设置轮转间隔为 1 分钟
	config.EnableCompression = true   // 启用日志压缩
	config.CompressionFormat = "gzip" // 设置压缩格式为 gzip

	// 创建 Logger
	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("创建日志记录器失败: %v", err)
	}
	defer logger.Close()

	// 写入足够的日志以触发轮转
	for i := 0; i < 1000; i++ {
		logger.Infof("这是第 %d 条测试日志消息", i)
	}

	// 等待轮转完成
	time.Sleep(2 * time.Minute)
	logger.Debug("等待轮转完成")

	// 检查日志文件数量
	files, err := os.ReadDir(logDir)
	if err != nil {
		t.Fatalf("读取日志目录失败: %v", err)
	}

	// 验证日志文件数量是否符合预期
	expectedCount := config.LogRetentionCount + 1 // 当前日志文件 + 保留的日志文件数量
	if len(files) != expectedCount {
		t.Errorf("期望的日志文件数量为 %d，但实际为 %d", expectedCount, len(files))
	}
}
