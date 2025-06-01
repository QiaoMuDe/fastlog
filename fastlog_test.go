package fastlog_test

import (
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gitee.com/MM-Q/fastlog/v2"
)

// TestConcurrentFastLog 测试并发场景下的多个日志记录器
func TestConcurrentFastLog(t *testing.T) {
	// 创建日志配置
	cfg := fastlog.NewFastLogConfig("logs", "test.log")
	cfg.LogLevel = fastlog.DEBUG
	cfg.LogFormat = fastlog.Simple
	cfg.IsLocalTime = true // 使用本地时间

	// 检查日志文件是否存在，如果存在则清空
	// if _, err := os.Stat(filepath.Join("logs", "test.log")); err == nil {
	// 	if err := os.Truncate(filepath.Join("logs", "test.log"), 0); err != nil {
	// 		t.Fatalf("清空日志文件失败: %v", err)
	// 	}
	// }

	// 创建日志记录器
	log, err := fastlog.NewFastLog(cfg)
	if err != nil {
		t.Fatalf("创建日志记录器失败: %v", err)
	}
	defer log.Close()

	// 持续时间为5秒
	duration := 10
	// 每秒生成10条日志
	rate := 1000

	// 启动随机日志函数
	randomLog(log, duration, rate)
}

// generateLogs 生成指定数量的模拟日志到通道中
func generateLogs(logMethodsNoFormat []func(v ...any), logMethodsWithFormat []func(format string, v ...interface{}), totalLoops int) <-chan func() {
	// 创建一个通道用于传递日志
	logChan := make(chan func(), totalLoops)

	// 创建一个随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	go func() {
		for i := 0; i < totalLoops; i++ {
			if r.Intn(2) == 0 {
				method := logMethodsNoFormat[r.Intn(len(logMethodsNoFormat))]
				logChan <- func() { method("这是一个测试日志") }
			} else {
				method := logMethodsWithFormat[r.Intn(len(logMethodsWithFormat))]
				logChan <- func() { method("这是一个测试日志: %s", "test") }
			}
		}
		close(logChan)
	}()

	return logChan
}

// randomLog 该函数用于按指定速率从通道中取出日志并打印
func randomLog(log *fastlog.FastLog, duration int, rate int) {
	// 定义无格式化日志方法的切片
	logMethodsNoFormat := []func(v ...any){
		log.Info,
		log.Warn,
		log.Error,
		log.Debug,
		log.Success,
	}
	// 定义格式化日志方法的切片
	logMethodsWithFormat := []func(format string, v ...interface{}){
		log.Infof,
		log.Warnf,
		log.Errorf,
		log.Debugf,
		log.Successf,
	}

	// 计算总循环次数
	totalLoops := duration * rate
	// 计算每次输出日志的间隔时间
	interval := time.Second / time.Duration(rate)

	// 生成日志到通道
	logChan := generateLogs(logMethodsNoFormat, logMethodsWithFormat, totalLoops)

	// 按指定速率从通道中取出日志并打印
	for logFunc := range logChan {
		logFunc()
		time.Sleep(interval)
	}
}

// TestCustomFormat 测试自定义日志格式
func TestCustomFormat(t *testing.T) {
	// 创建日志配置
	cfg := fastlog.NewFastLogConfig("logs", "custom.log")
	cfg.LogLevel = fastlog.DEBUG
	cfg.LogFormat = fastlog.Custom

	// 检查日志文件是否存在，如果存在则清空
	if _, err := os.Stat(filepath.Join("logs", "custom.log")); err == nil {
		if err := os.Truncate(filepath.Join("logs", "custom.log"), 0); err != nil {
			t.Fatalf("清空日志文件失败: %v", err)
		}
	}

	// 创建日志记录器
	log, err := fastlog.NewFastLog(cfg)
	if err != nil {
		t.Fatalf("创建日志记录器失败: %v", err)
	}
	defer log.Close()

	// 定义web应用程序日志格式
	webAppLogFormat := `%s [%s] %s %s %s %d %d %s %s %dms`

	// 模拟web应用程序日志记录
	log.Errorf(webAppLogFormat, "127.0.0.1", "2023-10-01 12:00:00", "GET", "/index.html", "HTTP/1.1", 200, 1234, "Mozilla/5.0", "en-US", 100)
	log.Infof(webAppLogFormat, "127.0.0.1", "2023-10-01 12:00:01", "POST", "/login", "HTTP/1.1", 401, 0, "Mozilla/5.0", "en-US", 200)
	log.Successf(webAppLogFormat, "127.0.0.1", "2023-10-01 12:00:02", "GET", "/profile", "HTTP/1.1", 200, 5678, "Mozilla/5.0", "en-US", 300)
	log.Warnf(webAppLogFormat, "127.0.0.1", "2023-10-01 12:00:03", "PUT", "/settings", "HTTP/1.1", 500, 0, "Mozilla/5.0", "en-US", 500)
	log.Debugf(webAppLogFormat, "127.0.0.1", "2023-10-01 12:00:04", "DELETE", "/logout", "HTTP/1.1", 200, 0, "Mozilla/5.0", "en-US", 50)
}

// TestNoColor 测试无颜色日志
func TestNoColor(t *testing.T) {
	// 创建日志配置
	cfg := fastlog.NewFastLogConfig("logs", "nocolor.log")
	cfg.LogLevel = fastlog.DEBUG
	cfg.NoColor = true // 禁用终端颜色

	// 检查日志文件是否存在，如果存在则清空
	if _, err := os.Stat(filepath.Join("logs", "custom.log")); err == nil {
		if err := os.Truncate(filepath.Join("logs", "custom.log"), 0); err != nil {
			t.Fatalf("清空日志文件失败: %v", err)
		}
	}

	// 创建日志记录器
	log, err := fastlog.NewFastLog(cfg)
	if err != nil {
		t.Fatalf("创建日志记录器失败: %v", err)
	}
	defer log.Close()

	// 打印测试日志
	log.Info("测试无颜色日志")
	log.Warn("测试无颜色日志")
	log.Error("测试无颜色日志")
	log.Debug("测试无颜色日志")
	log.Success("测试无颜色日志")
}

// TestNoBold 测试无加粗日志
func TestNoBold(t *testing.T) {
	// 创建日志配置
	cfg := fastlog.NewFastLogConfig("logs", "nobold.log")
	cfg.LogLevel = fastlog.DEBUG
	cfg.NoBold = true // 禁用终端字体加粗

	// 检查日志文件是否存在，如果存在则清空
	if _, err := os.Stat(filepath.Join("logs", "custom.log")); err == nil {
		if err := os.Truncate(filepath.Join("logs", "custom.log"), 0); err != nil {
			t.Fatalf("清空日志文件失败: %v", err)
		}
	}

	// 创建日志记录器
	log, err := fastlog.NewFastLog(cfg)
	if err != nil {
		t.Fatalf("创建日志记录器失败: %v", err)
	}
	defer log.Close()

	// 打印测试日志
	log.Info("测试无加粗日志")
	log.Warn("测试无加粗日志")
	log.Error("测试无加粗日志")
	log.Debug("测试无加粗日志")
	log.Success("测试无加粗日志")
}

func TestRmLogs(t *testing.T) {
	// 检查当前目录下是否存在logs目录
	if _, err := os.Stat("logs"); err == nil {
		// 删除logs目录及其下的所有文件和子目录
		if err := os.RemoveAll("logs"); err != nil {
			t.Fatalf("删除logs目录失败: %v", err)
		}
	}
}
