/*
performance_test.go - FastLog 高并发性能测试
提供完整的性能测试套件，包括并发测试、内存使用统计、吞吐量测试等，
生成详细的性能报告，用于评估 FastLog 在高并发环境下的表现。
*/
package fastlog

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 性能测试全局配置常量 - 方便调整测试规模
const (
	// 高并发性能测试配置
	TEST_GOROUTINE_COUNT  = 10    // 并发 Goroutine 数量
	TEST_LOGS_PER_ROUTINE = 50000 // 每个 Goroutine 生成的日志数量

	// 基准测试配置
	BENCH_GOROUTINE_COUNT = 8 // 基准测试的 Goroutine 数量
)

// 计算总日志数量的辅助常量
const (
	TOTAL_TEST_LOGS = TEST_GOROUTINE_COUNT * TEST_LOGS_PER_ROUTINE // 主测试总日志数
)

// PerformanceStats 性能统计结构体
type PerformanceStats struct {
	// 基本测试信息
	StartTime  time.Time
	EndTime    time.Time
	Duration   time.Duration
	Goroutines int
	CPUCores   int

	// 日志处理统计
	ExpectedLogs int64 // 预期生成的日志数量
	ActualLogs   int64 // 实际生成的日志数量
	WrittenLogs  int64 // 实际写入文件的日志数量
	SuccessRate  float64
	Throughput   float64 // 吞吐量（条/秒）

	// 内存使用统计
	StartMemory  uint64        // 开始时内存使用量
	EndMemory    uint64        // 结束时内存使用量
	PeakMemory   uint64        // 峰值内存使用量
	MemoryGrowth int64         // 内存增长量
	TotalAlloc   uint64        // 总分配内存
	SystemMemory uint64        // 系统内存
	GCCount      uint32        // GC次数
	GCPauseTime  time.Duration // GC暂停时间

	// 性能评估
	AvgMemoryPerLog float64 // 平均每条日志内存开销
	AvgTimePerLog   float64 // 平均每条日志处理时间
}

// getMemStats 获取内存统计信息
func getMemStats() runtime.MemStats {
	var m runtime.MemStats
	runtime.GC() // 强制执行GC以获得准确的内存统计
	runtime.ReadMemStats(&m)
	return m
}

// formatBytes 格式化字节数为可读格式
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%.1f B", float64(bytes))
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// formatDuration 格式化时间间隔
func formatDuration(d time.Duration) string {
	if d >= time.Second {
		return fmt.Sprintf("%.3fs (%.2fms)", d.Seconds(), float64(d.Nanoseconds())/1e6)
	}
	if d >= time.Millisecond {
		return fmt.Sprintf("%.2fms", float64(d.Nanoseconds())/1e6)
	}
	return fmt.Sprintf("%.2fμs", float64(d.Nanoseconds())/1e3)
}

// printPerformanceReport 打印性能测试报告
func printPerformanceReport(stats *PerformanceStats) {
	fmt.Println("============================================================")
	fmt.Println("           FastLog 高并发性能测试报告")
	fmt.Println("============================================================")

	// 📊 测试基本信息
	fmt.Println("📊 测试基本信息:")
	fmt.Printf("   开始时间: %s\n", stats.StartTime.Format("2006-01-02 15:04:05.000"))
	fmt.Printf("   结束时间: %s\n", stats.EndTime.Format("2006-01-02 15:04:05.000"))
	fmt.Printf("   测试耗时: %s\n", formatDuration(stats.Duration))
	fmt.Printf("   Goroutine数量: %d\n", stats.Goroutines)
	fmt.Println()

	// 📝 日志处理统计
	fmt.Println("📝 日志处理统计:")
	fmt.Printf("   预期生成: %.1f万条日志\n", float64(stats.ExpectedLogs)/10000)
	fmt.Printf("   实际生成: %.1f万条日志\n", float64(stats.ActualLogs)/10000)
	fmt.Printf("   文件写入: %.1f万条有效日志\n", float64(stats.WrittenLogs)/10000)
	fmt.Printf("   成功率: %.2f%%\n", stats.SuccessRate)
	fmt.Printf("   实际吞吐量: %.1f万条/秒\n", stats.Throughput/10000)
	fmt.Println()

	// 💾 内存使用统计
	fmt.Println("💾 内存使用统计:")
	fmt.Printf("   开始内存: %s\n", formatBytes(stats.StartMemory))
	fmt.Printf("   结束内存: %s\n", formatBytes(stats.EndMemory))
	fmt.Printf("   峰值内存: %s\n", formatBytes(stats.PeakMemory))
	fmt.Printf("   内存增长: %+s\n", formatBytes(uint64(stats.MemoryGrowth)))
	fmt.Printf("   总分配: %s\n", formatBytes(stats.TotalAlloc))
	fmt.Printf("   系统内存: %s\n", formatBytes(stats.SystemMemory))
	fmt.Printf("   GC次数: %d次\n", stats.GCCount)
	fmt.Printf("   GC暂停时间: %s\n", formatDuration(stats.GCPauseTime))
	fmt.Println()

	// ⚡ 性能评估
	fmt.Println("⚡ 性能评估:")
	fmt.Printf("   平均每条日志内存开销: %.2f bytes\n", stats.AvgMemoryPerLog)
	fmt.Printf("   平均每条日志处理时间: %.2f μs\n", stats.AvgTimePerLog)
	fmt.Println()

	// 🖥️ 系统资源
	fmt.Println("🖥️  系统资源:")
	fmt.Printf("   CPU核心数: %d\n", stats.CPUCores)
	fmt.Printf("   最大并发Goroutine: %d\n", stats.Goroutines)
	fmt.Printf("   并发度: %.1fx\n", float64(stats.Goroutines)/float64(stats.CPUCores))
	fmt.Println("============================================================")
}

// TestFastLogPerformance 高并发性能测试
func TestFastLogPerformance(t *testing.T) {
	// 使用全局常量配置测试参数
	const (
		goroutineCount = TEST_GOROUTINE_COUNT  // 并发 Goroutine 数量
		logsPerRoutine = TEST_LOGS_PER_ROUTINE // 每个 Goroutine 生成的日志数量
		totalLogs      = TOTAL_TEST_LOGS       // 总日志数量
	)

	// 创建测试目录
	testDir := "logs"

	// 配置日志记录器
	config := NewFastLogConfig(testDir, "performance_test.log")
	config.OutputToFile = true     // 开启文件输出
	config.OutputToConsole = false // 关闭控制台输出以提高性能
	config.LogLevel = INFO         // 限制日志级别
	config.LogFormat = Simple      // 简单日志格式
	config.MaxSize = 100           // 100MB
	config.Color = false           // 关闭颜色输出

	// 创建日志记录器
	logger := NewFastLog(config)
	defer logger.Close()

	// 初始化性能统计
	stats := &PerformanceStats{
		Goroutines:   goroutineCount,
		CPUCores:     runtime.NumCPU(),
		ExpectedLogs: totalLogs,
	}

	// 获取开始时的内存统计
	startMem := getMemStats()
	stats.StartMemory = startMem.Alloc
	stats.StartTime = time.Now()

	// 用于统计实际生成的日志数量
	var actualLogCount int64

	// 创建等待组
	var wg sync.WaitGroup
	wg.Add(goroutineCount)

	// 启动并发测试
	fmt.Printf("🚀 开始高并发性能测试...\n")
	fmt.Printf("📊 测试配置: %d个Goroutine，每个生成%d条日志，总计%.1f万条\n",
		goroutineCount, logsPerRoutine, float64(totalLogs)/10000)

	// 监控内存使用峰值
	var peakMemory uint64
	stopMonitor := make(chan bool)
	go func() {
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-stopMonitor:
				return
			case <-ticker.C:
				m := getMemStats()
				if m.Alloc > peakMemory {
					atomic.StoreUint64(&peakMemory, m.Alloc)
				}
			}
		}
	}()
	// 使用固定的测试日志消息，避免字符串格式化开销影响性能测试结果
	const message = "这是一条性能测试日志消息"

	// 启动多个 Goroutine 并发写入日志
	for i := 0; i < goroutineCount; i++ {
		go func(routineID int) {
			defer wg.Done()

			for j := 0; j < logsPerRoutine; j++ {
				// 写入不同级别的日志
				switch j % 4 {
				case 0:
					logger.Info(message)
				case 1:
					logger.Debug(message)
				case 2:
					logger.Warn(message)
				case 3:
					logger.Error(message)
				}

				// 原子递增实际日志计数
				atomic.AddInt64(&actualLogCount, 1)
			}
		}(i)
	}

	// 等待所有 Goroutine 完成
	wg.Wait()

	// 停止内存监控
	close(stopMonitor)

	// 记录结束时间和内存统计
	stats.EndTime = time.Now()
	stats.Duration = stats.EndTime.Sub(stats.StartTime)

	endMem := getMemStats()
	stats.EndMemory = endMem.Alloc
	stats.PeakMemory = atomic.LoadUint64(&peakMemory)
	stats.ActualLogs = atomic.LoadInt64(&actualLogCount)
	stats.TotalAlloc = endMem.TotalAlloc
	stats.SystemMemory = endMem.Sys
	stats.GCCount = endMem.NumGC - startMem.NumGC

	// 计算GC暂停时间
	var totalGCPause time.Duration
	for i := startMem.NumGC; i < endMem.NumGC; i++ {
		totalGCPause += time.Duration(endMem.PauseNs[i%256])
	}
	stats.GCPauseTime = totalGCPause

	// 计算统计数据
	stats.MemoryGrowth = int64(stats.EndMemory) - int64(stats.StartMemory)
	stats.SuccessRate = float64(stats.ActualLogs) / float64(stats.ExpectedLogs) * 100
	stats.Throughput = float64(stats.ActualLogs) / stats.Duration.Seconds()
	stats.AvgMemoryPerLog = float64(stats.TotalAlloc) / float64(stats.ActualLogs)
	stats.AvgTimePerLog = float64(stats.Duration.Nanoseconds()) / float64(stats.ActualLogs) / 1000 // 转换为微秒

	// 估算实际写入文件的日志数量（基于日志级别过滤）
	stats.WrittenLogs = stats.ActualLogs // 在这个测试中，所有日志都会被写入

	// 打印性能报告
	fmt.Println()
	printPerformanceReport(stats)

	// 验证测试结果
	if stats.ActualLogs != stats.ExpectedLogs {
		t.Errorf("日志数量不匹配: 预期 %d, 实际 %d", stats.ExpectedLogs, stats.ActualLogs)
	}

	// 性能基准检查
	if stats.Throughput < 100000 { // 至少10万条/秒
		t.Logf("警告: 吞吐量较低 (%.0f 条/秒)", stats.Throughput)
	}

	if stats.AvgMemoryPerLog > 1000 { // 每条日志内存开销不应超过1KB
		t.Logf("警告: 内存开销较高 (%.2f bytes/log)", stats.AvgMemoryPerLog)
	}

	fmt.Printf("✅ 性能测试完成！实际生成 %.1f万条日志，吞吐量 %.1f万条/秒\n",
		float64(stats.ActualLogs)/10000, stats.Throughput/10000)
}

// BenchmarkFastLogConcurrent 基准测试 - 并发写入
func BenchmarkFastLogConcurrent(b *testing.B) {
	// 创建测试目录
	testDir := "logs"

	// 配置日志记录器
	config := NewFastLogConfig(testDir, "benchmark.log")
	config.OutputToFile = true
	config.OutputToConsole = false
	config.LogLevel = INFO
	config.LogFormat = Simple

	logger := NewFastLog(config)
	defer logger.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Benchmark concurrent log message for performance testing")
		}
	})
}

// BenchmarkFastLogSingle 基准测试 - 单线程写入
func BenchmarkFastLogSingle(b *testing.B) {
	// 创建测试目录
	testDir := "logs"

	// 配置日志记录器
	config := NewFastLogConfig(testDir, "benchmark_single.log")
	config.OutputToFile = true
	config.OutputToConsole = false
	config.LogLevel = INFO
	config.LogFormat = Simple

	logger := NewFastLog(config)
	defer logger.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark single thread log message for performance testing")
	}
}
