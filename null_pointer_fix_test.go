/*
null_pointer_fix_test.go - 空指针修复验证测试文件
验证对关键路径空指针检查修复的有效性，确保在各种异常情况下系统不会panic，
包含高并发稳定性测试、资源耗尽处理测试和对象池安全性验证。
*/
package fastlog

import (
	"sync"
	"testing"
	"time"
)

// TestNullPointerProtection 测试空指针保护机制
func TestNullPointerProtection(t *testing.T) {
	t.Run("nil FastLog instance", func(t *testing.T) {
		var logger *FastLog = nil

		// 测试所有公共API方法都能安全处理nil实例
		logger.Info("test message")
		logger.Debug("test message")
		logger.Warn("test message")
		logger.Error("test message")
		logger.Success("test message")

		logger.Infof("test %s", "message")
		logger.Debugf("test %s", "message")
		logger.Warnf("test %s", "message")
		logger.Errorf("test %s", "message")
		logger.Successf("test %s", "message")

		// 如果程序能执行到这里，说明空指针保护有效
		t.Log("空指针保护测试通过")
	})

	t.Run("empty format string", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "test.log")
		cfg.OutputToConsole = false
		cfg.OutputToFile = false

		logger, err := NewFastLog(cfg)
		if err != nil {
			t.Fatalf("创建日志实例失败: %v", err)
		}
		defer logger.Close()

		// 测试空格式字符串的处理
		logger.Infof("")
		logger.Debugf("")
		logger.Warnf("")
		logger.Errorf("")
		logger.Successf("")

		t.Log("空格式字符串保护测试通过")
	})
}

// TestObjectPoolSafety 测试对象池安全性
func TestObjectPoolSafety(t *testing.T) {
	t.Run("concurrent object pool access", func(t *testing.T) {
		const goroutineCount = 100
		const operationsPerGoroutine = 1000

		var wg sync.WaitGroup
		wg.Add(goroutineCount)

		// 并发访问对象池
		for i := 0; i < goroutineCount; i++ {
			go func() {
				defer wg.Done()

				for j := 0; j < operationsPerGoroutine; j++ {
					// 获取对象
					msg := getLogMsg()
					if msg == nil {
						t.Errorf("getLogMsg() 返回了 nil")
						return
					}

					// 使用对象
					msg.Message = "test message"
					msg.Level = INFO
					msg.Timestamp = "2023-01-01 12:00:00"

					// 归还对象
					putLogMsg(msg)
				}
			}()
		}

		wg.Wait()
		t.Log("对象池并发安全测试通过")
	})

	t.Run("object pool nil handling", func(t *testing.T) {
		// 测试putLogMsg对nil的处理
		putLogMsg(nil) // 不应该panic

		// 测试getLogMsg永远不返回nil
		for i := 0; i < 100; i++ {
			msg := getLogMsg()
			if msg == nil {
				t.Fatal("getLogMsg() 不应该返回 nil")
			}
			putLogMsg(msg)
		}

		t.Log("对象池nil处理测试通过")
	})
}

// TestHighConcurrencyStability 测试高并发场景下的稳定性
func TestHighConcurrencyStability(t *testing.T) {
	cfg := NewFastLogConfig("logs", "stability_test.log")
	cfg.OutputToConsole = false
	cfg.OutputToFile = true
	cfg.ChanIntSize = 1000 // 较小的通道容量，更容易触发背压

	logger, err := NewFastLog(cfg)
	if err != nil {
		t.Fatalf("创建日志实例失败: %v", err)
	}
	defer logger.Close()

	const goroutineCount = 50
	const messagesPerGoroutine = 1000

	var wg sync.WaitGroup
	wg.Add(goroutineCount)

	// 启动多个goroutine并发写入日志
	for i := 0; i < goroutineCount; i++ {
		go func(id int) {
			defer wg.Done()

			for j := 0; j < messagesPerGoroutine; j++ {
				// 随机选择不同的日志方法
				switch j % 5 {
				case 0:
					logger.Info("并发测试消息", id, j)
				case 1:
					logger.Debug("并发测试消息", id, j)
				case 2:
					logger.Warn("并发测试消息", id, j)
				case 3:
					logger.Error("并发测试消息", id, j)
				case 4:
					logger.Success("并发测试消息", id, j)
				}
			}
		}(i)
	}

	wg.Wait()

	// 等待日志处理完成
	time.Sleep(500 * time.Millisecond)

	t.Logf("高并发稳定性测试完成: %d个goroutine，每个写入%d条日志",
		goroutineCount, messagesPerGoroutine)
}

// TestResourceExhaustionHandling 测试资源耗尽场景的处理
func TestResourceExhaustionHandling(t *testing.T) {
	cfg := NewFastLogConfig("logs", "resource_test.log")
	cfg.OutputToConsole = false
	cfg.OutputToFile = true
	cfg.ChanIntSize = 10 // 非常小的通道容量

	logger, err := NewFastLog(cfg)
	if err != nil {
		t.Fatalf("创建日志实例失败: %v", err)
	}
	defer logger.Close()

	// 快速发送大量日志，触发背压和通道满的情况
	for i := 0; i < 1000; i++ {
		logger.Infof("资源耗尽测试消息 %d", i)
	}

	// 系统应该能够优雅处理，不会panic或死锁
	time.Sleep(100 * time.Millisecond)

	t.Log("资源耗尽处理测试通过")
}

// BenchmarkNullPointerCheck 基准测试空指针检查的性能影响
func BenchmarkNullPointerCheck(b *testing.B) {
	cfg := NewFastLogConfig("logs", "benchmark.log")
	cfg.OutputToConsole = false
	cfg.OutputToFile = false // 禁用所有输出，只测试检查逻辑

	logger, err := NewFastLog(cfg)
	if err != nil {
		b.Fatalf("创建日志实例失败: %v", err)
	}
	defer logger.Close()

	b.ResetTimer()

	// 测试添加空指针检查后的性能影响
	for i := 0; i < b.N; i++ {
		logger.Info("基准测试消息")
	}
}
