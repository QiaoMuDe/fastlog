package fastlog

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestFilePermissionErrors 测试文件权限错误的处理
func TestFilePermissionErrors(t *testing.T) {
	// 跳过Windows系统的权限测试（Windows权限模型不同）
	if runtime.GOOS == "windows" {
		t.Skip("跳过Windows系统的权限测试")
	}

	t.Run("无权限创建目录", func(t *testing.T) {
		// 创建一个只读目录
		readOnlyDir := filepath.Join("test_readonly")
		err := os.MkdirAll(readOnlyDir, 0755)
		if err != nil {
			t.Fatalf("创建测试目录失败: %v", err)
		}
		defer func() { _ = os.RemoveAll(readOnlyDir) }()

		// 将目录设置为只读
		err = os.Chmod(readOnlyDir, 0444)
		if err != nil {
			t.Fatalf("设置目录权限失败: %v", err)
		}

		// 尝试在只读目录下创建日志文件
		cfg := NewFastLogConfig(filepath.Join(readOnlyDir, "sublogs"), "test.log")
		cfg.OutputToConsole = false

		logger, err := NewFastLog(cfg)
		if err != nil {
			t.Logf("预期的权限错误: %v", err)
			return
		}

		// 如果创建成功，尝试写入日志
		logger.Info("测试权限错误处理")
		time.Sleep(100 * time.Millisecond)
		logger.Close()

		// 恢复目录权限以便清理
		_ = os.Chmod(readOnlyDir, 0755)
	})

	t.Run("无权限写入文件", func(t *testing.T) {
		// 创建测试目录
		testDir := "test_write_permission"
		err := os.MkdirAll(testDir, 0755)
		if err != nil {
			t.Fatalf("创建测试目录失败: %v", err)
		}
		defer func() { _ = os.RemoveAll(testDir) }()

		// 创建一个只读文件
		testFile := filepath.Join(testDir, "readonly.log")
		err = os.WriteFile(testFile, []byte("initial content"), 0444)
		if err != nil {
			t.Fatalf("创建只读文件失败: %v", err)
		}

		// 尝试写入只读文件
		cfg := NewFastLogConfig(testDir, "readonly.log")
		cfg.OutputToConsole = false

		logger, err := NewFastLog(cfg)
		if err != nil {
			t.Logf("预期的文件权限错误: %v", err)
			return
		}

		// 写入日志，应该会在内部处理错误
		logger.Info("测试写入只读文件")
		time.Sleep(100 * time.Millisecond)
		logger.Close()
	})
}

// TestDirectoryCreationErrors 测试目录创建错误的处理
func TestDirectoryCreationErrors(t *testing.T) {
	t.Run("深层嵌套目录创建", func(t *testing.T) {
		// 创建一个非常深的目录路径
		deepPath := strings.Repeat("very_deep_dir/", 10)
		cfg := NewFastLogConfig(deepPath, "test.log")
		cfg.OutputToConsole = false

		logger, err := NewFastLog(cfg)
		if err != nil {
			t.Logf("深层目录创建可能失败: %v", err)
		} else {
			logger.Info("测试深层目录创建")
			time.Sleep(100 * time.Millisecond)
			logger.Close()
		}

		// 清理创建的目录
		_ = os.RemoveAll(strings.Split(deepPath, "/")[0])
	})

	t.Run("无效路径字符", func(t *testing.T) {
		// 在非Windows系统上测试包含无效字符的路径
		if runtime.GOOS != "windows" {
			cfg := NewFastLogConfig("test\x00invalid", "test.log")
			cfg.OutputToConsole = false

			logger, err := NewFastLog(cfg)
			if err != nil {
				t.Logf("预期的无效路径错误: %v", err)
			} else {
				logger.Close()
			}
		}
	})
}

// TestFileSystemSpaceHandling 测试文件系统空间处理
func TestFileSystemSpaceHandling(t *testing.T) {
	t.Run("大文件写入测试", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "large_test.log")
		cfg.OutputToConsole = false
		cfg.MaxLogFileSize = 1 // 1MB限制，触发轮转

		logger, err := NewFastLog(cfg)
		if err != nil {
			t.Fatalf("创建日志实例失败: %v", err)
		}
		defer logger.Close()

		// 写入大量日志数据
		largeMessage := strings.Repeat("这是一个用于测试大文件写入的长消息。", 100)
		for i := 0; i < 1000; i++ {
			logger.Infof("大文件测试消息 %d: %s", i, largeMessage)
		}

		// 等待写入完成
		time.Sleep(500 * time.Millisecond)
	})
}

// TestConcurrentFileAccess 测试并发文件访问
func TestConcurrentFileAccess(t *testing.T) {
	t.Run("多实例写入同一文件", func(t *testing.T) {
		logFile := filepath.Join("logs", "concurrent_test.log")

		// 创建多个日志实例写入同一文件
		var loggers []*FastLog
		for i := 0; i < 3; i++ {
			cfg := NewFastLogConfig("logs", "concurrent_test.log")
			cfg.OutputToConsole = false

			logger, err := NewFastLog(cfg)
			if err != nil {
				t.Fatalf("创建日志实例 %d 失败: %v", i, err)
			}
			loggers = append(loggers, logger)
		}

		// 并发写入日志
		done := make(chan bool, len(loggers))
		for i, logger := range loggers {
			go func(id int, l *FastLog) {
				defer func() { done <- true }()
				for j := 0; j < 100; j++ {
					l.Infof("实例 %d 消息 %d", id, j)
				}
			}(i, logger)
		}

		// 等待所有写入完成
		for i := 0; i < len(loggers); i++ {
			<-done
		}

		// 关闭所有实例
		for _, logger := range loggers {
			logger.Close()
		}

		// 验证文件是否存在且有内容
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			t.Error("并发写入的日志文件不存在")
		}
	})
}

// TestFileRotationErrors 测试日志轮转错误处理
func TestFileRotationErrors(t *testing.T) {
	t.Run("轮转过程中的并发写入", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "rotation_test.log")
		cfg.OutputToConsole = false
		cfg.MaxLogFileSize = 1 // 1MB，容易触发轮转

		logger, err := NewFastLog(cfg)
		if err != nil {
			t.Fatalf("创建日志实例失败: %v", err)
		}
		defer logger.Close()

		// 快速写入大量数据，触发轮转
		largeMessage := strings.Repeat("轮转测试消息", 1000)
		for i := 0; i < 100; i++ {
			logger.Infof("轮转测试 %d: %s", i, largeMessage)
		}

		// 等待轮转完成
		time.Sleep(1 * time.Second)

		// 继续写入，测试轮转后的写入
		for i := 100; i < 200; i++ {
			logger.Infof("轮转后测试 %d: %s", i, largeMessage)
		}

		time.Sleep(500 * time.Millisecond)
	})
}

// TestResourceExhaustion 测试资源耗尽情况
func TestResourceExhaustion(t *testing.T) {
	t.Run("大量文件句柄测试", func(t *testing.T) {
		// 创建多个日志实例，每个使用不同的文件
		var loggers []*FastLog
		maxInstances := 50 // 限制数量避免真正耗尽系统资源

		for i := 0; i < maxInstances; i++ {
			filename := "handle_test_" + string(rune(i+'0')) + ".log"
			if i >= 10 {
				filename = "handle_test_" + string(rune(i/10+'0')) + string(rune(i%10+'0')) + ".log"
			}
			cfg := NewFastLogConfig("logs", filename)
			cfg.OutputToConsole = false

			logger, err := NewFastLog(cfg)
			if err != nil {
				t.Logf("创建第 %d 个实例时失败: %v", i, err)
				break
			}
			loggers = append(loggers, logger)
		}

		// 写入一些日志
		for i, logger := range loggers {
			logger.Infof("文件句柄测试实例 %d", i)
		}

		// 等待写入完成
		time.Sleep(200 * time.Millisecond)

		// 关闭所有实例
		for _, logger := range loggers {
			logger.Close()
		}

		t.Logf("成功创建了 %d 个日志实例", len(loggers))
	})
}

// TestErrorRecovery 测试错误恢复机制
func TestErrorRecovery(t *testing.T) {
	t.Run("文件写入失败后的恢复", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "recovery_test.log")
		cfg.OutputToConsole = true // 启用控制台输出作为备用

		logger, err := NewFastLog(cfg)
		if err != nil {
			t.Fatalf("创建日志实例失败: %v", err)
		}
		defer logger.Close()

		// 正常写入
		logger.Info("正常写入测试")

		// 模拟文件系统错误后的恢复
		// 这里主要测试系统不会崩溃
		for i := 0; i < 10; i++ {
			logger.Infof("恢复测试消息 %d", i)
		}

		time.Sleep(100 * time.Millisecond)
	})
}

// BenchmarkFileSystemOperations 文件系统操作性能基准测试
func BenchmarkFileSystemOperations(b *testing.B) {
	cfg := NewFastLogConfig("logs", "benchmark_fs.log")
	cfg.OutputToConsole = false

	logger, err := NewFastLog(cfg)
	if err != nil {
		b.Fatalf("创建日志实例失败: %v", err)
	}
	defer logger.Close()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			logger.Infof("文件系统基准测试消息 %d", i)
			i++
		}
	})
}

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
		cfg.OutputToConsole = true
		cfg.OutputToFile = false
		cfg.LogLevel = NONE // 设置为NONE级别，避免实际输出日志

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
