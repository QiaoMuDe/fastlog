package fastlog_test

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"gitee.com/MM-Q/fastlog"
)

// TestConcurrentFastLog 测试并发场景下的多个日志记录器
func TestConcurrentFastLog(t *testing.T) {
	// 配置第一个日志记录器
	config1 := fastlog.NewFastLogConfig("applogs1", "app1.log")
	config1.LogLevel = fastlog.DEBUG
	config1.LogFormat = fastlog.Bracket

	// 初始化第一个日志记录器
	logger1, err := fastlog.NewFastLog(config1)
	if err != nil {
		t.Fatalf("初始化第一个日志记录器失败: %v", err)
	}
	defer logger1.Close()

	// 配置第二个日志记录器
	config2 := fastlog.NewFastLogConfig("applogs2", "app2.log")
	config2.LogLevel = fastlog.DEBUG
	config2.LogFormat = fastlog.Bracket

	// 初始化第二个日志记录器
	logger2, err := fastlog.NewFastLog(config2)
	if err != nil {
		t.Fatalf("初始化第二个日志记录器失败: %v", err)
	}
	defer logger2.Close()

	// 定义协程数量
	numGoroutines := 10

	// 定义等待组
	var wg sync.WaitGroup

	// 启动协程，分别使用两个日志记录器记录日志
	for i := 0; i < numGoroutines; i++ {
		wg.Add(2) // 每个协程对应两个日志记录器

		go func(logger *fastlog.FastLog, id int) {
			defer wg.Done()

			// 持续运行10秒，每秒打印50条随机日志
			startTime := time.Now()
			jd := 0
			for time.Since(startTime) < 3*time.Second { // 持续运行10秒
				for i := 0; i < 10; i++ { // 每秒打印50条随机日志
					level := rand.Intn(5)                                   // 随机生成0-4的整数，对应Debug、Info、Warn、Error、Success
					message := fmt.Sprintf("随机日志消息 #%d (协程ID: %d)", jd, id) // 随机生成日志消息
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
		}(logger1, i)

		go func(logger *fastlog.FastLog, id int) {
			defer wg.Done()

			// 持续运行10秒，每秒打印50条随机日志
			startTime := time.Now()
			jd := 0
			for time.Since(startTime) < 3*time.Second { // 持续运行10秒
				for i := 0; i < 10; i++ { // 每秒打印50条随机日志
					level := rand.Intn(5)                                   // 随机生成0-4的整数，对应Debug、Info、Warn、Error、Success
					message := fmt.Sprintf("随机日志消息 #%d (协程ID: %d)", jd, id) // 随机生成日志消息
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
		}(logger2, i)
	}

	// 等待所有协程完成
	wg.Wait()

	// 等待1秒，确保日志写入完成
	time.Sleep(1 * time.Second)

	// 检查第一个日志文件是否存在
	if _, err := os.Stat("applogs1/app1.log"); os.IsNotExist(err) {
		t.Fatalf("第一个日志文件未生成: %v", err)
	}

	// 检查第二个日志文件是否存在
	if _, err := os.Stat("applogs2/app2.log"); os.IsNotExist(err) {
		t.Fatalf("第二个日志文件未生成: %v", err)
	}

	// 等待所有协程完成
	time.Sleep(2 * time.Second) // 确保所有协程都已退出

	// 检查第一个日志文件数量是否符合预期
	files1, err := os.ReadDir("applogs1")
	if err != nil {
		t.Fatalf("读取第一个日志目录失败: %v", err)
	}

	// 计算符合条件的第一个日志文件数量
	logFileCount1 := 0
	for _, file := range files1 {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) == ".log" || filepath.Ext(file.Name()) == ".zip" {
			logFileCount1++
		}
	}

	// 检查第二个日志文件数量是否符合预期
	files2, err := os.ReadDir("applogs2")
	if err != nil {
		t.Fatalf("读取第二个日志目录失败: %v", err)
	}

	// 计算符合条件的第二个日志文件数量
	logFileCount2 := 0
	for _, file := range files2 {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) == ".log" || filepath.Ext(file.Name()) == ".zip" {
			logFileCount2++
		}
	}

	t.Log("并发多日志记录器测试通过")
}
