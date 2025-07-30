package fastlog

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestConcurrentFastLog 测试并发场景下的多个日志记录器
func TestConcurrentFastLog(t *testing.T) {
	// 记录开始时间
	startTime := time.Now()

	// 创建日志配置
	cfg := NewFastLogConfig("logs", "test.log")
	cfg.OutputToConsole = true
	cfg.OutputToFile = true
	cfg.MaxLogFileSize = 5
	cfg.LogFormat = Simple
	cfg.ChanIntSize = 100000 // 增大通道容量以支持更高并发

	// 创建日志记录器
	log, err := NewFastLog(cfg)
	if err != nil {
		t.Fatalf("创建日志记录器失败: %v", err)
	}

	// 持续时间为3秒
	duration := 1
	// 每秒生成10000条日志
	rate := 1000000

	defer func() {
		log.Close()
		// 计算总耗时并打印
		totalDuration := time.Since(startTime)
		fmt.Printf("=== 并发日志测试结果 ===\n")
		fmt.Printf("测试配置: 持续时间 %d秒 | 目标速率 %d条/秒\n", duration, rate)
		fmt.Printf("预期生成: %d条日志\n", duration*rate)
		fmt.Printf("实际耗时: %.3fs (%.2fms)\n", totalDuration.Seconds(), float64(totalDuration.Nanoseconds())/1e6)
		fmt.Printf("实际速率: %.0f条/秒\n", float64(duration*rate)/totalDuration.Seconds())
		fmt.Printf("========================\n")
	}()

	// 启动高并发随机日志函数
	highConcurrencyRandomLog(log, duration, rate, t)
}

// highConcurrencyRandomLog 高并发随机日志生成函数
// 该函数真正并发地生成日志，而不是通过time.Sleep限制并发度
func highConcurrencyRandomLog(log *FastLog, duration int, rate int, t *testing.T) {
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

	// 使用WaitGroup同步并发操作
	var wg sync.WaitGroup
	wg.Add(totalLoops)

	// 并发生成所有日志
	for i := 0; i < totalLoops; i++ {
		go func() {
			defer wg.Done()
			// 为每个 goroutine 创建独立的随机数生成器
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			// 随机选择日志方法类型（无格式化或格式化）
			if r.Intn(2) == 0 {
				// 随机选择无格式化日志方法
				method := logMethodsNoFormat[r.Intn(len(logMethodsNoFormat))]
				method("这是一个高并发测试日志")
			} else {
				// 随机选择格式化日志方法
				method := logMethodsWithFormat[r.Intn(len(logMethodsWithFormat))]
				method("这是一个高并发测试日志: %s", "test")
			}
		}()
	}

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
		if strings.Contains(line, "这是一个高并发测试日志") {
			validLines++
		}
	}

	// 输出验证结果
	fmt.Printf("写入有效日志行数: %d\n", validLines)
}

// BenchmarkFastLog 高并发基准测试
func BenchmarkFastLog(b *testing.B) {
	// 创建日志配置
	cfg := NewFastLogConfig("logs", "benchmark.log")
	cfg.OutputToConsole = false // 基准测试中关闭控制台输出以减少I/O影响
	cfg.OutputToFile = true
	cfg.ChanIntSize = 100000

	// 创建日志记录器
	log, err := NewFastLog(cfg)
	if err != nil {
		b.Fatalf("创建日志记录器失败: %v", err)
	}

	defer log.Close()

	// 重置计时器
	b.ResetTimer()

	// 并发运行基准测试
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Info("基准测试日志消息")
		}
	})
}
