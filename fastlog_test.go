package fastlog_test

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
	cfg.IsLocalTime = false

	// 创建日志记录器
	log, err := fastlog.NewFastLog(cfg)
	if err != nil {
		t.Fatalf("创建日志记录器失败: %v", err)
	}
	defer log.Close()

	// 持续时间为3秒
	duration := 3
	// 每秒生成100条日志
	rate := 100

	// 启动随机日志函数
	randomLog(log, duration, rate, t)
}

// generateLogs 生成指定数量的模拟日志到通道中
func generateLogs(logMethodsNoFormat []func(v ...any), logMethodsWithFormat []func(format string, v ...interface{}), totalLoops int) <-chan func() {
	// 创建一个通道用于传递日志，缓冲区大小为总循环数的2倍
	logChan := make(chan func(), totalLoops*2)

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
func randomLog(log *fastlog.FastLog, duration int, rate int, t *testing.T) {
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
	// 由于 generateLogs 仅返回一个通道，这里只使用一个变量接收
	logChan := generateLogs(logMethodsNoFormat, logMethodsWithFormat, totalLoops)

	// 使用WaitGroup同步并发操作
	var wg sync.WaitGroup
	wg.Add(totalLoops)

	// 按指定速率从通道中取出日志并打印
	go func() {
		for logFunc := range logChan {
			go func(f func()) {
				defer wg.Done()
				f()
			}(logFunc)
			time.Sleep(interval)
		}
	}()

	// 等待所有日志生成完成
	wg.Wait()

	// 验证日志文件内容
	content, err := os.ReadFile(filepath.Join("logs", "test.log"))
	if err != nil {
		t.Fatalf("读取日志文件失败: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	validLines := 0
	for _, line := range lines {
		if strings.Contains(line, "这是一个测试日志") {
			validLines++
		}
	}
}

// TestCustomFormat 测试自定义日志格式
func TestCustomFormat(t *testing.T) {
	// 创建日志配置
	cfg := fastlog.NewFastLogConfig("logs", "custom.log")
	cfg.LogLevel = fastlog.DEBUG
	cfg.LogFormat = fastlog.Custom
	cfg.PrintToConsole = true

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
	cfg.NoColor = false

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

func TestNewFastLog(t *testing.T) {
	config := fastlog.NewFastLogConfig("logs", "test.log")

	log, err := fastlog.NewFastLog(config)
	if err != nil {
		t.Fatalf("创建FastLog实例失败: %v", err)
	}
	log.Info("测试日志")
	log.Warn("测试日志")
	log.Error("测试日志")

	if err := log.Close(); err != nil {
		t.Fatalf("关闭FastLog实例失败: %v", err)
	}
}

// TestCleanupLogs 测试完成后清理日志目录
func TestCleanupLogs(t *testing.T) {
	if _, err := os.Stat("logs"); err == nil {
		if err := os.RemoveAll("logs"); err != nil {
			t.Fatalf("删除logs目录失败: %v", err)
		}
	}
}
