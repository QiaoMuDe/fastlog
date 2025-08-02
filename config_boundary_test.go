/*
config_boundary_test.go - 配置边界值和异常情况测试文件
测试FastLog配置在各种边界值和异常情况下的处理能力，
验证配置修正函数的正确性和系统的健壮性。
*/
package fastlog

import (
	"strings"
	"testing"
	"time"
)

// TestConfigBoundaryValues 测试配置字段的边界值处理
func TestConfigBoundaryValues(t *testing.T) {
	t.Run("FlushInterval边界值测试", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "test.log")

		// 测试负值
		cfg.FlushInterval = -1 * time.Second
		cfg.fixFinalConfig()
		if cfg.FlushInterval != 500*time.Millisecond {
			t.Errorf("负值FlushInterval应被修正为500ms，实际为%v", cfg.FlushInterval)
		}

		// 测试零值
		cfg.FlushInterval = 0
		cfg.fixFinalConfig()
		if cfg.FlushInterval != 500*time.Millisecond {
			t.Errorf("零值FlushInterval应被修正为500ms，实际为%v", cfg.FlushInterval)
		}

		// 测试过小值
		cfg.FlushInterval = 5 * time.Millisecond
		cfg.fixFinalConfig()
		if cfg.FlushInterval != 10*time.Millisecond {
			t.Errorf("过小FlushInterval应被修正为10ms，实际为%v", cfg.FlushInterval)
		}

		// 测试过大值
		cfg.FlushInterval = 60 * time.Second
		cfg.fixFinalConfig()
		if cfg.FlushInterval != 30*time.Second {
			t.Errorf("过大FlushInterval应被修正为30s，实际为%v", cfg.FlushInterval)
		}
	})

	t.Run("ChanIntSize边界值测试", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "test.log")

		// 测试负值
		cfg.ChanIntSize = -100
		cfg.fixFinalConfig()
		if cfg.ChanIntSize != 10000 {
			t.Errorf("负值ChanIntSize应被修正为10000，实际为%d", cfg.ChanIntSize)
		}

		// 测试零值
		cfg.ChanIntSize = 0
		cfg.fixFinalConfig()
		if cfg.ChanIntSize != 10000 {
			t.Errorf("零值ChanIntSize应被修正为10000，实际为%d", cfg.ChanIntSize)
		}

		// 测试过大值
		cfg.ChanIntSize = 200000
		cfg.fixFinalConfig()
		if cfg.ChanIntSize != 100000 {
			t.Errorf("过大ChanIntSize应被修正为100000，实际为%d", cfg.ChanIntSize)
		}
	})

	t.Run("MaxLogFileSize边界值测试", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "test.log")

		// 测试负值
		cfg.MaxLogFileSize = -10
		cfg.fixFinalConfig()
		if cfg.MaxLogFileSize != 10 {
			t.Errorf("负值MaxLogFileSize应被修正为10，实际为%d", cfg.MaxLogFileSize)
		}

		// 测试零值
		cfg.MaxLogFileSize = 0
		cfg.fixFinalConfig()
		if cfg.MaxLogFileSize != 10 {
			t.Errorf("零值MaxLogFileSize应被修正为10，实际为%d", cfg.MaxLogFileSize)
		}

		// 测试过大值
		cfg.MaxLogFileSize = 2000
		cfg.fixFinalConfig()
		if cfg.MaxLogFileSize != 1000 {
			t.Errorf("过大MaxLogFileSize应被修正为1000，实际为%d", cfg.MaxLogFileSize)
		}
	})

	t.Run("MaxLogAge边界值测试", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "test.log")

		// 测试负值
		cfg.MaxLogAge = -30
		cfg.fixFinalConfig()
		if cfg.MaxLogAge != 0 {
			t.Errorf("负值MaxLogAge应被修正为0，实际为%d", cfg.MaxLogAge)
		}

		// 测试过大值
		cfg.MaxLogAge = 5000
		cfg.fixFinalConfig()
		if cfg.MaxLogAge != 3650 {
			t.Errorf("过大MaxLogAge应被修正为3650，实际为%d", cfg.MaxLogAge)
		}
	})

	t.Run("MaxLogBackups边界值测试", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "test.log")

		// 测试负值
		cfg.MaxLogBackups = -10
		cfg.fixFinalConfig()
		if cfg.MaxLogBackups != 0 {
			t.Errorf("负值MaxLogBackups应被修正为0，实际为%d", cfg.MaxLogBackups)
		}

		// 测试过大值
		cfg.MaxLogBackups = 2000
		cfg.fixFinalConfig()
		if cfg.MaxLogBackups != 1000 {
			t.Errorf("过大MaxLogBackups应被修正为1000，实际为%d", cfg.MaxLogBackups)
		}
	})
}

// TestInvalidLogLevel 测试无效日志级别的处理
func TestInvalidLogLevel(t *testing.T) {
	cfg := NewFastLogConfig("logs", "test.log")

	// 测试过小的日志级别 (小于DEBUG=10)
	cfg.LogLevel = 5
	cfg.fixFinalConfig()
	if cfg.LogLevel != INFO {
		t.Errorf("过小日志级别应被修正为INFO，实际为%d", cfg.LogLevel)
	}

	// 测试边界值 - DEBUG级别应该有效
	cfg.LogLevel = DEBUG
	cfg.fixFinalConfig()
	if cfg.LogLevel != DEBUG {
		t.Errorf("DEBUG级别应该保持不变，实际为%d", cfg.LogLevel)
	}

	// 测试边界值 - NONE级别应该有效
	cfg.LogLevel = NONE
	cfg.fixFinalConfig()
	if cfg.LogLevel != NONE {
		t.Errorf("NONE级别应该保持不变，实际为%d", cfg.LogLevel)
	}

	// 测试有效范围内的其他级别
	cfg.LogLevel = INFO
	cfg.fixFinalConfig()
	if cfg.LogLevel != INFO {
		t.Errorf("INFO级别应该保持不变，实际为%d", cfg.LogLevel)
	}
}

// TestInvalidLogFormat 测试无效日志格式的处理
func TestInvalidLogFormat(t *testing.T) {
	cfg := NewFastLogConfig("logs", "test.log")

	// 测试无效的日志格式
	cfg.LogFormat = LogFormatType(99)
	cfg.fixFinalConfig()
	if cfg.LogFormat != Detailed {
		t.Errorf("无效日志格式应被修正为Detailed，实际为%d", cfg.LogFormat)
	}

	// 测试负值日志格式
	cfg.LogFormat = LogFormatType(-1)
	cfg.fixFinalConfig()
	if cfg.LogFormat != Detailed {
		t.Errorf("负值日志格式应被修正为Detailed，实际为%d", cfg.LogFormat)
	}
}

// TestFileNameCleaning 测试文件名清理功能
func TestFileNameCleaning(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "包含非法字符的文件名",
			input:    "test<>:\"|?*.log",
			expected: "test_______.log", // 7个非法字符被替换为7个下划线
		},
		{
			name:     "包含路径遍历的文件名",
			input:    "../../../etc/passwd",
			expected: "passwd",
		},
		{
			name:     "空文件名",
			input:    "",
			expected: "default.log",
		},
		{
			name:     "纯空格文件名",
			input:    "   ",
			expected: "default.log",
		},
		{
			name:     "以点开头的文件名",
			input:    ".hidden.log",
			expected: "hidden.log",
		},
		{
			name:     "以点结尾的文件名",
			input:    "test.log.",
			expected: "test.log",
		},
		{
			name:     "包含多个连续斜杠",
			input:    "logs//test.log",
			expected: "logs\\test.log", // Windows系统使用反斜杠
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := cleanFileName(tc.input)
			if result != tc.expected {
				t.Errorf("文件名清理失败：输入=%q，期望=%q，实际=%q", tc.input, tc.expected, result)
			}
		})
	}
}

// TestLongFileName 测试超长文件名的处理
func TestLongFileName(t *testing.T) {
	// 创建一个超长的文件名（300字符）
	longName := strings.Repeat("a", 300) + ".log"

	cfg := NewFastLogConfig("logs", longName)
	cfg.fixFinalConfig()

	// 验证文件名被截断到合理长度
	if len(cfg.LogFileName) > 255 {
		t.Errorf("超长文件名应被截断，实际长度为%d", len(cfg.LogFileName))
	}

	// 验证仍然保留了.log扩展名
	if !strings.HasSuffix(cfg.LogFileName, ".log") {
		t.Errorf("截断后的文件名应保留扩展名，实际为%s", cfg.LogFileName)
	}
}

// TestEmptyDirectoryAndFileName 测试空目录名和文件名的处理
func TestEmptyDirectoryAndFileName(t *testing.T) {
	cfg := NewFastLogConfig("", "")
	cfg.fixFinalConfig()

	// 验证空目录名被修正
	if cfg.LogDirName != "logs" {
		t.Errorf("空目录名应被修正为'logs'，实际为%q", cfg.LogDirName)
	}

	// 验证空文件名被修正
	if cfg.LogFileName != "app.log" {
		t.Errorf("空文件名应被修正为'app.log'，实际为%q", cfg.LogFileName)
	}
}

// TestConfigConsistency 测试配置一致性
func TestConfigConsistency(t *testing.T) {
	t.Run("禁用所有输出的配置", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "test.log")
		cfg.OutputToConsole = false
		cfg.OutputToFile = false

		// 这种配置应该能正常创建日志实例，但不会有实际输出
		logger, err := NewFastLog(cfg)
		if err != nil {
			t.Fatalf("禁用所有输出的配置应该能正常创建日志实例：%v", err)
		}
		defer logger.Close()

		// 记录日志不应该出错
		logger.Info("测试消息")
	})

	t.Run("矛盾配置的处理", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "test.log")
		cfg.OutputToFile = false
		cfg.MaxLogFileSize = 100 // 设置了文件相关配置但禁用了文件输出

		// 应该能正常创建，文件相关配置被忽略
		logger, err := NewFastLog(cfg)
		if err != nil {
			t.Fatalf("矛盾配置应该能正常处理：%v", err)
		}
		defer logger.Close()
	})
}

// BenchmarkConfigFixing 配置修正性能基准测试
func BenchmarkConfigFixing(b *testing.B) {
	cfg := NewFastLogConfig("logs", "test.log")

	// 设置一些需要修正的值
	cfg.FlushInterval = -1 * time.Second
	cfg.ChanIntSize = -100
	cfg.MaxLogFileSize = -10
	cfg.LogLevel = 200 // 超过FATAL(60)但在uint8范围内

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cfg.fixFinalConfig()
	}
}
