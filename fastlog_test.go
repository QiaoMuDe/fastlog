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

// TestFastLog 测试日志记录器
func TestFastLog(t *testing.T) {
	// 配置日志
	config := fastlog.NewFastLogConfig("applogs", "app.log")
	config.LogLevel = fastlog.DEBUG
	config.LogFormat = fastlog.Threaded

	// 初始化日志记录器
	logger, err := fastlog.NewFastLog(config)
	if err != nil {
		t.Fatalf("初始化日志记录器失败: %v", err)
	}
	defer logger.Close()

	// 持续运行10秒，每秒打印50条随机日志
	startTime := time.Now()
	jd := 0
	for time.Since(startTime) < 10*time.Second { // 持续运行10秒
		for i := 0; i < 50; i++ { // 每秒打印50条随机日志
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

	t.Log("日志测试通过")
}
