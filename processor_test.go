/*
processor_test.go - 日志处理器性能测试文件
包含对日志处理器高并发性能、内存使用、吞吐量等关键指标的综合测试，
提供详细的性能统计报告和基准测试，用于评估FastLog在生产环境中的表现。
*/
package fastlog

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 测试配置常量
const (
	TestWan         = 10000         // 用于快捷计算的标准单位(万)
	TestDuration    = 3             // 测试时长（秒）
	TestRate        = 100 * TestWan // 每秒生成多少条日志（降低到100万避免过度压力）
	TaskChannelSize = 100000        // 任务通道缓冲区大小
)

var (
	WorkerPoolSize = 12 // 工作池大小（goroutine数量）
)

// TestStats 测试统计信息结构体
type TestStats struct {
	StartTime        time.Time        // 测试开始时间
	EndTime          time.Time        // 测试结束时间
	Duration         time.Duration    // 测试持续时间
	ExpectedLogs     int64            // 预期生成的日志数量
	ActualLogs       int64            // 实际生成的日志数量
	ValidLogLines    int64            // 有效日志行数
	StartMemStats    runtime.MemStats // 开始时的内存统计
	EndMemStats      runtime.MemStats // 结束时的内存统计
	PeakMemStats     runtime.MemStats // 峰值内存统计
	GoroutineCount   int              // goroutine数量
	SuccessRate      float64          // 成功率
	ThroughputPerSec float64          // 实际吞吐量（条/秒）
}

// TestConcurrentFastLog 测试并发场景下的多个日志记录器（优化版本）
func TestConcurrentFastLog(t *testing.T) {
	// 初始化测试统计信息
	stats := &TestStats{
		StartTime: time.Now(),
	}

	// 强制垃圾回收，获取干净的初始内存状态
	runtime.GC()
	runtime.GC() // 执行两次确保彻底回收
	time.Sleep(50 * time.Millisecond)
	runtime.ReadMemStats(&stats.StartMemStats)

	// 创建日志配置
	cfg := NewFastLogConfig("logs", "test.log")
	cfg.OutputToConsole = false       // 控制台输出
	cfg.OutputToFile = true           // 文件输出
	cfg.MaxSize = 5                   // 设置日志文件最大大小为5MB
	cfg.LogFormat = Simple            // 设置日志格式
	cfg.ChanIntSize = TaskChannelSize // 增大通道容量以支持更高并发
	cfg.DisableBackpressure = false   // 禁用背压

	// 创建日志记录器
	log := NewFastLog(cfg)

	// 测试参数
	stats.ExpectedLogs = int64(TestDuration * TestRate)
	stats.GoroutineCount = WorkerPoolSize // 使用实际的工作池大小

	// 启动内存监控goroutine
	stopMonitoring := make(chan bool)
	go monitorMemoryUsage(stats, stopMonitoring)

	defer func() {
		// 停止内存监控
		close(stopMonitoring)

		// 等待通道中的日志被处理完成
		waitStart := time.Now()
		maxWaitTime := 10 * time.Second // 最大等待10秒

		// 等待通道中的日志数量降到合理范围（批处理大小的一半）
		for len(log.logChan) > cfg.BatchSize/2 && time.Since(waitStart) < maxWaitTime {
			time.Sleep(100 * time.Millisecond)
		}

		// 关闭日志器并等待处理完成
		log.Close()
		time.Sleep(500 * time.Millisecond) // 增加等待时间确保处理完成

		// 记录结束时间和内存状态
		stats.EndTime = time.Now()
		stats.Duration = stats.EndTime.Sub(stats.StartTime)

		// 强制垃圾回收后获取最终内存状态
		runtime.GC()
		runtime.GC()
		time.Sleep(50 * time.Millisecond)
		runtime.ReadMemStats(&stats.EndMemStats)

		// 计算统计数据
		stats.ActualLogs = stats.ExpectedLogs // 在实际测试中会被更新
		stats.SuccessRate = float64(stats.ValidLogLines) / float64(stats.ExpectedLogs) * 100
		stats.ThroughputPerSec = float64(stats.ActualLogs) / stats.Duration.Seconds()

		// 打印详细统计结果
		stats.PrintDetailedStats()
	}()

	// 启动高并发随机日志函数
	actualLogs := highConcurrencyRandomLogWithStats(log, TestDuration, TestRate, stats, t)
	stats.ActualLogs = actualLogs
}

// PrintDetailedStats 打印详细的测试统计结果
func (s *TestStats) PrintDetailedStats() {
	separator := strings.Repeat("=", 60)
	fmt.Printf("\n%s\n", separator)
	fmt.Printf("           FastLog 高并发性能测试报告\n")
	fmt.Printf("%s\n", separator)

	// 基本测试信息
	fmt.Printf("📊 测试基本信息:\n")
	fmt.Printf("   开始时间: %s\n", s.StartTime.Format("2006-01-02 15:04:05.000"))
	fmt.Printf("   结束时间: %s\n", s.EndTime.Format("2006-01-02 15:04:05.000"))
	fmt.Printf("   测试耗时: %.3fs (%.2fms)\n", s.Duration.Seconds(), float64(s.Duration.Nanoseconds())/1e6)
	fmt.Printf("   Goroutine数量: %d\n", s.GoroutineCount)

	// 日志处理统计
	fmt.Printf("\n📝 日志处理统计:\n")
	expectedStr := formatNumber(s.ExpectedLogs)
	actualStr := formatNumber(s.ActualLogs)
	validStr := formatNumber(s.ValidLogLines)
	throughputStr := formatNumber(int64(s.ThroughputPerSec))
	fmt.Printf("   预期生成: %s条日志\n", expectedStr)
	fmt.Printf("   实际生成: %s条日志\n", actualStr)
	fmt.Printf("   文件写入: %s条有效日志\n", validStr)
	fmt.Printf("   成功率: %.2f%%\n", s.SuccessRate)
	fmt.Printf("   实际吞吐量: %s条/秒\n", throughputStr)

	// 内存使用统计
	fmt.Printf("\n💾 内存使用统计:\n")
	startMemStr := formatBytes(s.StartMemStats.Alloc)
	endMemStr := formatBytes(s.EndMemStats.Alloc)
	peakMemStr := formatBytes(s.PeakMemStats.Alloc)
	totalAllocStr := formatBytes(s.EndMemStats.TotalAlloc)
	sysMemStr := formatBytes(s.EndMemStats.Sys)

	fmt.Printf("   开始内存: %s\n", startMemStr)
	fmt.Printf("   结束内存: %s\n", endMemStr)
	fmt.Printf("   峰值内存: %s\n", peakMemStr)

	memoryChange := int64(s.EndMemStats.Alloc) - int64(s.StartMemStats.Alloc)
	if memoryChange >= 0 {
		changeStr := formatBytes(uint64(memoryChange))
		fmt.Printf("   内存增长: +%s\n", changeStr)
	} else {
		changeStr := formatBytes(uint64(-memoryChange))
		fmt.Printf("   内存减少: -%s\n", changeStr)
	}

	fmt.Printf("   总分配: %s\n", totalAllocStr)
	fmt.Printf("   系统内存: %s\n", sysMemStr)
	fmt.Printf("   GC次数: %d次\n", s.EndMemStats.NumGC-s.StartMemStats.NumGC)
	fmt.Printf("   GC暂停时间: %.2fms\n", float64(s.EndMemStats.PauseTotalNs-s.StartMemStats.PauseTotalNs)/1e6)

	// 性能评估
	fmt.Printf("\n⚡ 性能评估:\n")
	memPerLog := float64(memoryChange) / float64(s.ActualLogs)
	if memPerLog > 0 {
		fmt.Printf("   平均每条日志内存开销: %.2f bytes\n", memPerLog)
	}
	fmt.Printf("   平均每条日志处理时间: %.2f μs\n", float64(s.Duration.Nanoseconds())/float64(s.ActualLogs)/1000)

	// 系统资源利用率
	fmt.Printf("\n🖥️  系统资源:\n")
	fmt.Printf("   CPU核心数: %d\n", runtime.NumCPU())
	fmt.Printf("   最大并发Goroutine: %d\n", s.GoroutineCount)
	fmt.Printf("   并发度: %.1fx\n", float64(s.GoroutineCount)/float64(runtime.NumCPU()))

	finalSeparator := strings.Repeat("=", 60)
	fmt.Printf("%s\n\n", finalSeparator)
}

// formatNumber 格式化数字，添加中文单位
func formatNumber(n int64) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	if n < 10000 {
		return fmt.Sprintf("%.1f千", float64(n)/1000)
	}
	if n < 100000000 {
		return fmt.Sprintf("%.1f万", float64(n)/10000)
	}
	return fmt.Sprintf("%.1f亿", float64(n)/100000000)
}

// formatBytes 格式化字节数
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// monitorMemoryUsage 监控内存使用情况，记录峰值
func monitorMemoryUsage(stats *TestStats, stop <-chan bool) {
	ticker := time.NewTicker(10 * time.Millisecond) // 每10ms检查一次
	defer ticker.Stop()

	var maxAlloc uint64 = 0

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			if m.Alloc > maxAlloc {
				maxAlloc = m.Alloc
				stats.PeakMemStats = m
			}
		}
	}
}

// LogTask 日志任务结构体
type LogTask struct {
	Index int
	Type  int // 0: 无格式化, 1: 格式化
}

// highConcurrencyRandomLogWithStats 高并发随机日志生成函数（优化版本 - 使用工作池）
// 使用固定数量的goroutine处理大量日志任务，避免创建过多goroutine
func highConcurrencyRandomLogWithStats(log *FastLog, duration int, rate int, stats *TestStats, t *testing.T) int64 {
	// 定义无格式化日志方法的切片
	logMethodsNoFormat := []func(v ...any){
		log.Info,
		log.Warn,
		log.Error,
		log.Debug,
	}
	// 定义格式化日志方法的切片
	logMethodsWithFormat := []func(format string, v ...interface{}){
		log.Infof,
		log.Warnf,
		log.Errorf,
		log.Debugf,
	}

	// 计算总任务数
	totalTasks := duration * rate

	// 创建任务通道
	taskChan := make(chan LogTask, TaskChannelSize)

	// 使用WaitGroup同步工作池
	var wg sync.WaitGroup

	// 记录实际发送的日志数量（使用原子操作保证并发安全）
	var actualLogsSent int64

	// 启动工作池
	for i := 0; i < WorkerPoolSize; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// 为每个worker创建独立的随机数生成器
			r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(workerID)))

			// 处理任务
			for task := range taskChan {
				// 随机选择日志方法类型
				if task.Type == 0 {
					// 随机选择无格式化日志方法
					method := logMethodsNoFormat[r.Intn(len(logMethodsNoFormat))]
					method("这是一个高并发测试日志", task.Index)
				} else {
					// 随机选择格式化日志方法
					method := logMethodsWithFormat[r.Intn(len(logMethodsWithFormat))]
					method("这是一个高并发测试日志: %s [%d]", "test", task.Index)
				}

				// 原子递增实际发送的日志数量
				atomic.AddInt64(&actualLogsSent, 1)
			}
		}(i)
	}

	// 生成任务并发送到任务通道
	go func() {
		defer close(taskChan)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		for i := 0; i < totalTasks; i++ {
			task := LogTask{
				Index: i,
				Type:  r.Intn(2), // 随机选择日志类型
			}
			taskChan <- task
		}
	}()

	// 等待所有工作完成
	wg.Wait()

	// 等待一段时间让日志进入通道
	time.Sleep(200 * time.Millisecond)

	// 🔍 验证所有日志文件内容（包括轮转文件）
	validLines := int64(0)
	logDir := "logs"

	// 读取日志目录中的所有文件
	files, err := filepath.Glob(filepath.Join(logDir, "test*.log"))
	if err != nil {
		t.Logf("读取日志目录失败: %v", err)
		stats.ValidLogLines = 0
		return actualLogsSent
	}

	// 统计所有日志文件中的有效行数
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Logf("读取日志文件 %s 失败: %v", file, err)
			continue
		}

		lines := strings.Split(string(content), "\n")
		fileValidLines := int64(0)
		for _, line := range lines {
			if strings.Contains(line, "这是一个高并发测试日志") {
				fileValidLines++
			}
		}
		validLines += fileValidLines
	}

	stats.ValidLogLines = validLines
	return actualLogsSent
}

// BenchmarkFastLog 高并发基准测试
func BenchmarkFastLog(b *testing.B) {
	// 创建日志配置
	cfg := NewFastLogConfig("logs", "benchmark.log")
	cfg.OutputToConsole = false // 基准测试中关闭控制台输出以减少I/O影响
	cfg.OutputToFile = true
	cfg.ChanIntSize = 100000

	// 创建日志记录器
	log := NewFastLog(cfg)
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
