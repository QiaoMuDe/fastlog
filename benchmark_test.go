/*
benchmark_test.go - FastLog性能基准测试文件
包含对背压控制函数、数据类型性能对比等关键功能的基准测试，
用于评估和优化FastLog在高负载场景下的性能表现。
*/
package fastlog

import (
	"testing"
)

// BenchmarkShouldDropLogByBackpressureOriginal 测试背压函数的性能（原始版本）
func BenchmarkShouldDropLogByBackpressureOriginal(b *testing.B) {
	// 创建一个测试通道
	logChan := make(chan *logMsg, 1000)
	bp := &bpThresholds{
		threshold80: 80,
		threshold90: 90,
		threshold95: 95,
		threshold98: 98,
	}

	// 填充一些数据模拟实际使用
	for i := 0; i < 5000; i++ {
		logChan <- &logMsg{Level: INFO, Message: "test"}
	}

	b.ResetTimer()

	// 基准测试
	for i := 0; i < b.N; i++ {
		shouldDropLogByBackpressure(bp, logChan, INFO)
	}
}

// BenchmarkShouldDropLogByBackpressure_HighLoad 测试高负载场景
func BenchmarkShouldDropLogByBackpressure_HighLoad(b *testing.B) {
	// 创建一个接近满载的通道
	logChan := make(chan *logMsg, 1000)
	bp := &bpThresholds{
		threshold80: 80,
		threshold90: 90,
		threshold95: 95,
		threshold98: 98,
	}

	// 填充95%的数据
	for i := 0; i < 9500; i++ {
		logChan <- &logMsg{Level: INFO, Message: "test"}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		shouldDropLogByBackpressure(bp, logChan, DEBUG)
	}
}

// BenchmarkShouldDropLogByBackpressure_EmptyChannel 测试空通道场景
func BenchmarkShouldDropLogByBackpressure_EmptyChannel(b *testing.B) {
	logChan := make(chan *logMsg, 1000)
	bp := &bpThresholds{
		threshold80: 80,
		threshold90: 90,
		threshold95: 95,
		threshold98: 98,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		shouldDropLogByBackpressure(bp, logChan, INFO)
	}
}

// BenchmarkInt64VsInt 对比int64和int的性能差异
func BenchmarkInt64VsInt(b *testing.B) {
	chanLen := 5000
	chanCap := 10000

	b.Run("Int64Calculation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var usage int64
			usage = (int64(chanLen) * 100) / int64(chanCap)
			if usage > 100 {
				usage = 100
			}
			_ = usage
		}
	})

	b.Run("IntCalculation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var usage int
			usage = (chanLen * 100) / chanCap
			if usage > 100 {
				usage = 100
			}
			_ = usage
		}
	})
}

// BenchmarkShutdownPerformance 关闭性能基准测试
func BenchmarkShutdownPerformance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cfg := NewFastLogConfig("logs", "shutdown_bench.log")
		cfg.OutputToConsole = false

		logger := NewFastLog(cfg)

		// 写入一些日志
		for j := 0; j < 10; j++ {
			logger.Infof("基准测试消息 %d", j)
		}

		// 测量关闭时间
		logger.Close()
	}
}
