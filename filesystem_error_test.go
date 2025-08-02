/*
filesystem_error_test.go - 文件系统错误处理测试文件
测试FastLog在各种文件系统异常情况下的处理能力，
包括权限不足、磁盘空间不足、文件锁定等场景。
*/
package fastlog

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
