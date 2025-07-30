package fastlog

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestMain 全局测试入口，控制非verbose模式下的输出重定向
func TestMain(m *testing.M) {
	flag.Parse() // 解析命令行参数
	// 保存原始标准输出和错误输出
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	var nullFile *os.File
	var err error

	// 非verbose模式下重定向到空设备
	if !testing.Verbose() {
		nullFile, err = os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
		if err != nil {
			panic("无法打开空设备文件: " + err.Error())
		}
		os.Stdout = nullFile
		os.Stderr = nullFile
	}

	// 运行所有测试
	exitCode := m.Run()

	// 恢复原始输出
	if !testing.Verbose() {
		os.Stdout = originalStdout
		os.Stderr = originalStderr
		_ = nullFile.Close()
	}

	os.Exit(exitCode)
}

// TestConcurrentFastLog 测试并发场景下的多个日志记录器
func TestConcurrentFastLog(t *testing.T) {
	// 记录开始时间
	startTime := time.Now()

	// 创建日志配置
	cfg := NewFastLogConfig("logs", "test.log")

	// 创建日志记录器
	log, err := NewFastLog(cfg)
	if err != nil {
		t.Fatalf("创建日志记录器失败: %v", err)
	}

	// 持续时间为3秒
	duration := 3
	// 每秒生成100条日志
	rate := 10

	defer func() {
		log.Close()
		// 计算总耗时并打印
		totalDuration := time.Since(startTime)
		fmt.Printf("=== 并发日志测试结果 ===\n")
		fmt.Printf("测试配置: 持续时间 %d秒 | 目标速率 %d条/秒\n", duration, rate)
		fmt.Printf("预期生成: %d条日志\n", duration*rate)
		fmt.Printf("实际耗时: %.3fs (%.2fms)\n", totalDuration.Seconds(), float64(totalDuration.Nanoseconds())/1e6)
		fmt.Printf("========================\n")
	}()

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
func randomLog(log *FastLog, duration int, rate int, t *testing.T) {
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
