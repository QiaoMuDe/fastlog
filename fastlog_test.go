package fastlog_test

import (
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gitee.com/MM-Q/fastlog"
)

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

// TestConcurrentFastLog 测试并发场景下的多个日志记录器
func TestConcurrentFastLog(t *testing.T) {
	// 创建日志配置
	cfg := fastlog.NewFastLogConfig("logs", "test.log")
	cfg.LogLevel = fastlog.DEBUG
	cfg.LogFormat = fastlog.Simple

	// 检查日志文件是否存在，如果存在则清空
	if _, err := os.Stat(filepath.Join("logs", "test.log")); err == nil {
		if err := os.Truncate(filepath.Join("logs", "test.log"), 0); err != nil {
			t.Fatalf("清空日志文件失败: %v", err)
		}
	}

	// 创建日志记录器
	log, err := fastlog.NewFastLog(cfg)
	if err != nil {
		t.Fatalf("创建日志记录器失败: %v", err)
	}
	defer log.Close()

	// 持续时间为5秒
	duration := 5
	// 每秒生成10条日志
	rate := 10

	// 启动随机日志函数
	randomLog(log, duration, rate)
}
