/*
graceful_shutdown_test.go - 优雅关闭机制测试文件
测试FastLog的Close()方法在各种情况下的正确性，
验证优雅关闭、超时处理和资源清理的完整性。
*/
package fastlog

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"
)

// TestGracefulShutdown 测试优雅关闭机制
func TestGracefulShutdown(t *testing.T) {
	t.Run("正常关闭流程", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "shutdown_test.log")
		cfg.OutputToConsole = false

		logger, err := NewFastLog(cfg)
		if err != nil {
			t.Fatalf("创建日志实例失败: %v", err)
		}

		// 写入一些日志
		for i := 0; i < 100; i++ {
			logger.Infof("关闭测试消息 %d", i)
		}

		// 记录关闭开始时间
		startTime := time.Now()

		// 执行关闭
		logger.Close()

		// 记录关闭耗时
		shutdownDuration := time.Since(startTime)
		t.Logf("关闭耗时: %v", shutdownDuration)

		// 验证关闭后不能再写入日志（不应该panic）
		logger.Info("关闭后的消息")

		// 验证日志文件存在且有内容
		logFile := filepath.Join("logs", "shutdown_test.log")
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			t.Error("关闭后日志文件不存在")
		}
	})

	t.Run("重复关闭测试", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "repeat_close_test.log")
		cfg.OutputToConsole = false

		logger, err := NewFastLog(cfg)
		if err != nil {
			t.Fatalf("创建日志实例失败: %v", err)
		}

		logger.Info("重复关闭测试")

		// 第一次关闭
		logger.Close()

		// 重复关闭应该是安全的（幂等性）
		logger.Close()
		logger.Close()

		// 不应该panic
		t.Log("重复关闭测试通过")
	})

	t.Run("并发关闭测试", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "concurrent_close_test.log")
		cfg.OutputToConsole = false

		logger, err := NewFastLog(cfg)
		if err != nil {
			t.Fatalf("创建日志实例失败: %v", err)
		}

		logger.Info("并发关闭测试")

		// 启动多个goroutine同时关闭
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				logger.Close()
			}(i)
		}

		wg.Wait()
		t.Log("并发关闭测试通过")
	})
}

// TestShutdownWithActiveWriting 测试关闭过程中的活跃写入
func TestShutdownWithActiveWriting(t *testing.T) {
	cfg := NewFastLogConfig("logs", "active_writing_test.log")
	cfg.OutputToConsole = false
	cfg.ChanIntSize = 1000

	logger, err := NewFastLog(cfg)
	if err != nil {
		t.Fatalf("创建日志实例失败: %v", err)
	}

	// 启动后台写入
	var wg sync.WaitGroup
	stopWriting := make(chan bool)

	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 0
		for {
			select {
			case <-stopWriting:
				return
			default:
				logger.Infof("活跃写入测试消息 %d", i)
				i++
				time.Sleep(1 * time.Millisecond)
			}
		}
	}()

	// 让写入运行一段时间
	time.Sleep(100 * time.Millisecond)

	// 停止写入并关闭
	close(stopWriting)
	logger.Close()

	wg.Wait()
	t.Log("活跃写入关闭测试通过")
}

// TestShutdownTimeout 测试关闭超时处理
func TestShutdownTimeout(t *testing.T) {
	cfg := NewFastLogConfig("logs", "timeout_test.log")
	cfg.OutputToConsole = false
	cfg.ChanIntSize = 10000              // 大通道容量
	cfg.FlushInterval = 10 * time.Second // 长刷新间隔

	logger, err := NewFastLog(cfg)
	if err != nil {
		t.Fatalf("创建日志实例失败: %v", err)
	}

	// 填充大量日志到通道中
	for i := 0; i < 5000; i++ {
		logger.Infof("超时测试消息 %d", i)
	}

	// 测试关闭超时
	startTime := time.Now()
	logger.Close()
	shutdownDuration := time.Since(startTime)

	t.Logf("关闭耗时: %v", shutdownDuration)

	// 验证关闭在合理时间内完成
	if shutdownDuration > 30*time.Second {
		t.Errorf("关闭耗时过长: %v", shutdownDuration)
	}
}

// TestResourceCleanup 测试资源清理
func TestResourceCleanup(t *testing.T) {
	t.Run("goroutine清理验证", func(t *testing.T) {
		initialGoroutines := runtime.NumGoroutine()

		cfg := NewFastLogConfig("logs", "cleanup_test.log")
		cfg.OutputToConsole = false

		logger, err := NewFastLog(cfg)
		if err != nil {
			t.Fatalf("创建日志实例失败: %v", err)
		}

		logger.Info("资源清理测试")
		time.Sleep(100 * time.Millisecond)

		logger.Close()
		time.Sleep(100 * time.Millisecond) // 等待goroutine退出

		finalGoroutines := runtime.NumGoroutine()

		// 允许一定的goroutine数量波动
		if finalGoroutines > initialGoroutines+2 {
			t.Errorf("可能存在goroutine泄漏: 初始=%d, 最终=%d", initialGoroutines, finalGoroutines)
		}
	})

	t.Run("文件句柄清理验证", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "handle_cleanup_test.log")
		cfg.OutputToConsole = false

		logger, err := NewFastLog(cfg)
		if err != nil {
			t.Fatalf("创建日志实例失败: %v", err)
		}

		logger.Info("文件句柄清理测试")
		time.Sleep(100 * time.Millisecond)

		logger.Close()

		// 验证日志文件可以被删除（说明文件句柄已关闭）
		logFile := filepath.Join("logs", "handle_cleanup_test.log")
		time.Sleep(100 * time.Millisecond) // 等待文件句柄释放

		err = os.Remove(logFile)
		if err != nil {
			t.Errorf("无法删除日志文件，可能存在文件句柄泄漏: %v", err)
		}
	})
}

// BenchmarkShutdownPerformance 关闭性能基准测试
func BenchmarkShutdownPerformance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cfg := NewFastLogConfig("logs", "shutdown_bench.log")
		cfg.OutputToConsole = false

		logger, err := NewFastLog(cfg)
		if err != nil {
			b.Fatalf("创建日志实例失败: %v", err)
		}

		// 写入一些日志
		for j := 0; j < 10; j++ {
			logger.Infof("基准测试消息 %d", j)
		}

		// 测量关闭时间
		logger.Close()
	}
}
