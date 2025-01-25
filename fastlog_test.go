package fastlog_test

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gitee.com/MM-Q/fastlog"
)

// TestLogger 测试日志记录器
func TestLogger(t *testing.T) {
	// 配置日志
	config := fastlog.NewConfig("applogs", "app.log")
	config.LogLevel = fastlog.Debug
	config.EnableLogRotation = true // 启用日志轮转
	config.EnableCompression = true // 启用日志压缩
	config.CompressionFormat = "zip"
	config.LogFormat = fastlog.Threaded
	config.LogMaxSize = "1kb"    // 设置日志文件大小限制为1KB，便于测试轮转
	config.RotationInterval = 1  // 设置日志轮转间隔为1秒
	config.LogRetentionCount = 2 // 设置日志保留数量为2

	// 初始化日志记录器
	logger, err := fastlog.NewLogger(config)
	if err != nil {
		t.Fatalf("初始化日志记录器失败: %v", err)
	}
	defer logger.Close()

	// 持续运行5分钟，每秒打印20条随机日志
	startTime := time.Now()
	jd := 0
	for time.Since(startTime) < 1*time.Minute {
		for i := 0; i < 50; i++ {
			level := rand.Intn(5)                    // 随机生成0-4的整数，对应Debug、Info、Warn、Error、Success
			message := fmt.Sprintf("随机日志消息 #%d", jd) // 随机生成日志消息
			switch level {
			case 0:
				logger.Debug(message)
			case 1:
				logger.Info(message)
			case 2:
				logger.Warn(message)
			case 3:
				logger.Error(message)
			case 4:
				logger.Success(message)
			}
			jd++
		}
		time.Sleep(1 * time.Second)
	}

	// 等待1秒，确保日志写入完成
	time.Sleep(1 * time.Second)

	// 检查日志文件是否存在
	if _, err := os.Stat("applogs/app.log"); os.IsNotExist(err) {
		t.Fatalf("日志文件未生成: %v", err)
	}

	// 等待所有协程完成
	time.Sleep(2 * time.Second) // 确保所有协程都已退出

	// 检查日志文件数量是否符合预期
	files, err := os.ReadDir("applogs")
	if err != nil {
		t.Fatalf("读取日志目录失败: %v", err)
	}

	// 计算符合条件的日志文件数量
	logFileCount := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) == ".log" || filepath.Ext(file.Name()) == ".zip" {
			logFileCount++
		}
	}

	// 验证日志文件数量是否符合预期
	expectedCount := config.LogRetentionCount + 1 // 当前日志文件 + 保留的日志文件数量
	if logFileCount != expectedCount {
		t.Fatalf("日志文件数量不符合预期，期望 %d 个，实际 %d 个", expectedCount, logFileCount)
	}

	// 等待1秒，确保日志压缩完成
	time.Sleep(10 * time.Second)

	// 清理日志文件
	os.RemoveAll("applogs")
	// 检查是否清理成功
	if _, err := os.Stat("applogs"); !os.IsNotExist(err) {
		t.Fatalf("日志文件清理失败: %v", err)
	}

	t.Log("日志测试通过")
}
