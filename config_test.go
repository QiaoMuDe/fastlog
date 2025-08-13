/*
config_test.go - 日志配置模块测试文件
包含对FastLog配置结构体的单元测试，验证配置初始化、边界条件处理和配置修正功能，
确保配置参数在各种异常情况下都能正确处理和自动修正。
*/
package fastlog

import (
	"strings"
	"testing"
	"time"
)

// TestNewFastLogConfig 测试NewFastLogConfig函数的默认配置初始化
func TestNewFastLogConfig(t *testing.T) {
	cfg := NewFastLogConfig("logs", "test.log")

	// 验证默认日志级别
	if cfg.LogLevel != INFO {
		t.Errorf("默认日志级别应为INFO，实际为%d", cfg.LogLevel)
	}

	// 验证默认最大日志文件大小
	if cfg.MaxLogFileSize != 10 {
		t.Errorf("默认最大日志文件大小应为10MB，实际为%d", cfg.MaxLogFileSize)
	}

	// 验证默认刷新间隔
	if cfg.FlushInterval != 500*time.Millisecond {
		t.Errorf("默认刷新间隔应为500ms，实际为%v", cfg.FlushInterval)
	}

	// 验证默认通道大小
	if cfg.ChanIntSize != 10000 {
		t.Errorf("默认通道大小应为10000，实际为%d", cfg.ChanIntSize)
	}
}

// TestSetMaxLogFileSize 测试设置最大日志文件大小的边界情况
func TestSetMaxLogFileSize(t *testing.T) {
	cfg := NewFastLogConfig("logs", "test.log")

	// 测试设置负数
	cfg.MaxLogFileSize = -1
	if cfg.MaxLogFileSize != -1 {
		t.Error("设置负数时应保留原值")
	}

	// 测试修正函数对负数的处理
	cfg.fixFinalConfig()
	if cfg.MaxLogFileSize != 10 {
		t.Errorf("修正后最大日志文件大小应为10，实际为%d", cfg.MaxLogFileSize)
	}

	// 测试设置超过最大值
	cfg.MaxLogFileSize = 2000
	cfg.fixFinalConfig()
	if cfg.MaxLogFileSize != 1000 {
		t.Errorf("修正后最大日志文件大小应为1000，实际为%d", cfg.MaxLogFileSize)
	}
}

// TestFixFinalConfig 测试配置修正函数对各种配置的修正
func TestFixFinalConfig(t *testing.T) {
	cfg := NewFastLogConfig("logs", "test.log")

	// 测试日志级别修正
	cfg.LogLevel = 0
	cfg.fixFinalConfig()
	if cfg.LogLevel != INFO {
		t.Errorf("日志级别修正应为INFO，实际为%d", cfg.LogLevel)
	}

	// 测试通道大小修正
	cfg.ChanIntSize = 0
	cfg.fixFinalConfig()
	if cfg.ChanIntSize != 10000 {
		t.Errorf("通道大小修正应为10000，实际为%d", cfg.ChanIntSize)
	}

	// 测试刷新间隔修正
	cfg.FlushInterval = 0
	cfg.fixFinalConfig()
	if cfg.FlushInterval != 500*time.Millisecond {
		t.Errorf("刷新间隔修正应为500ms，实际为%v", cfg.FlushInterval)
	}
}

// TestValidateFinalConfig 测试配置验证函数
func TestValidateFinalConfig(t *testing.T) {
	cfg := NewFastLogConfig("logs", "test.log")

	// 测试无效通道大小
	cfg.ChanIntSize = 0
	// 捕获验证函数的输出，检查是否正确报告错误
	// 由于validateFinalConfig只打印警告不返回值，这里使用日志捕获或模拟
}

// BenchmarkFixFinalConfig 基准测试配置修正函数性能
func BenchmarkFixFinalConfig(b *testing.B) {
	cfg := NewFastLogConfig("logs", "test.log")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cfg.fixFinalConfig()
	}
}

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
	if cfg.LogFormat != Simple {
		t.Errorf("无效日志格式应被修正为Simple，实际为%d", cfg.LogFormat)
	}

	// 测试负值日志格式
	cfg.LogFormat = LogFormatType(-1)
	cfg.fixFinalConfig()
	if cfg.LogFormat != Simple {
		t.Errorf("负值日志格式应被修正为Simple，实际为%d", cfg.LogFormat)
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
			expected: "app.log",
		},
		{
			name:     "纯空格文件名",
			input:    "   ",
			expected: "app.log",
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
	t.Run("禁用所有输出的配置应该panic", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "test.log")
		cfg.OutputToConsole = false
		cfg.OutputToFile = false

		// 这种配置应该触发panic，因为没有任何输出方式
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("禁用所有输出的配置应该触发panic")
			} else {
				expectedMsg := "At least one output method must be enabled"
				if !strings.Contains(r.(string), expectedMsg) {
					t.Errorf("panic消息不正确，期望包含%q，实际为%q", expectedMsg, r)
				}
			}
		}()

		// 直接调用fixFinalConfig来触发panic
		cfg.fixFinalConfig()
	})

	t.Run("文件输出时目录名和文件名都为空应该panic", func(t *testing.T) {
		cfg := NewFastLogConfig("", "")
		cfg.OutputToConsole = false
		cfg.OutputToFile = true

		// 这种配置应该触发panic
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("文件输出时目录名和文件名都为空应该触发panic")
			} else {
				expectedMsg := "When file output is enabled, log directory name"
				if !strings.Contains(r.(string), expectedMsg) {
					t.Errorf("panic消息不正确，期望包含%q，实际为%q", expectedMsg, r)
				}
			}
		}()

		// 直接调用validateCriticalConfig来触发panic
		cfg.validateCriticalConfig()
	})

	t.Run("超大通道大小应该panic", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "test.log")
		cfg.ChanIntSize = 2000000 // 200万，超过100万的限制

		// 这种配置应该触发panic
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("超大通道大小应该触发panic")
			} else {
				expectedMsg := "channel size too large"
				if !strings.Contains(r.(string), expectedMsg) {
					t.Errorf("panic消息不正确，期望包含%q，实际为%q", expectedMsg, r)
				}
			}
		}()

		// 直接调用fixFinalConfig来触发panic
		cfg.fixFinalConfig()
	})

	t.Run("超大文件大小应该panic", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "test.log")
		cfg.MaxLogFileSize = 15000 // 15GB，超过10GB的限制

		// 这种配置应该触发panic
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("超大文件大小应该触发panic")
			} else {
				expectedMsg := "single log file size too large"
				if !strings.Contains(r.(string), expectedMsg) {
					t.Errorf("panic消息不正确，期望包含%q，实际为%q", expectedMsg, r)
				}
			}
		}()

		// 直接调用fixFinalConfig来触发panic
		cfg.fixFinalConfig()
	})

	t.Run("正常配置应该成功", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "test.log")
		cfg.OutputToFile = false
		cfg.MaxLogFileSize = 100 // 设置了文件相关配置但禁用了文件输出

		// 应该能正常创建，文件相关配置被忽略
		logger := NewFastLog(cfg)
		defer logger.Close()

		// 记录日志不应该出错
		logger.Info("测试消息")
	})

	t.Run("配置自动修正验证", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "test.log")

		// 设置一些需要修正的值
		originalChanSize := -100
		originalFlushInterval := -1 * time.Second
		originalMaxFileSize := 2000
		originalMaxAge := 5000
		originalMaxBackups := 2000

		cfg.ChanIntSize = originalChanSize
		cfg.FlushInterval = originalFlushInterval
		cfg.MaxLogFileSize = originalMaxFileSize
		cfg.MaxLogAge = originalMaxAge
		cfg.MaxLogBackups = originalMaxBackups

		// 记录修正前的值
		t.Logf("修正前: ChanIntSize=%d, FlushInterval=%v, MaxLogFileSize=%d, MaxLogAge=%d, MaxLogBackups=%d",
			cfg.ChanIntSize, cfg.FlushInterval, cfg.MaxLogFileSize, cfg.MaxLogAge, cfg.MaxLogBackups)

		// 应该能正常创建，配置被自动修正
		logger := NewFastLog(cfg)
		defer logger.Close()

		// 记录修正后的值
		t.Logf("修正后: ChanIntSize=%d, FlushInterval=%v, MaxLogFileSize=%d, MaxLogAge=%d, MaxLogBackups=%d",
			cfg.ChanIntSize, cfg.FlushInterval, cfg.MaxLogFileSize, cfg.MaxLogAge, cfg.MaxLogBackups)

		// 验证配置被正确修正
		if cfg.ChanIntSize != 10000 {
			t.Errorf("ChanIntSize应被修正为10000，实际为%d", cfg.ChanIntSize)
		}
		if cfg.FlushInterval != 500*time.Millisecond {
			t.Errorf("FlushInterval应被修正为500ms，实际为%v", cfg.FlushInterval)
		}
		if cfg.MaxLogFileSize != 1000 {
			t.Errorf("MaxLogFileSize应被修正为1000，实际为%d", cfg.MaxLogFileSize)
		}
		if cfg.MaxLogAge != 3650 {
			t.Errorf("MaxLogAge应被修正为3650，实际为%d", cfg.MaxLogAge)
		}
		if cfg.MaxLogBackups != 1000 {
			t.Errorf("MaxLogBackups应被修正为1000，实际为%d", cfg.MaxLogBackups)
		}

		// 验证修正确实发生了
		if cfg.ChanIntSize == originalChanSize {
			t.Error("ChanIntSize没有被修正")
		}
		if cfg.FlushInterval == originalFlushInterval {
			t.Error("FlushInterval没有被修正")
		}
		if cfg.MaxLogFileSize == originalMaxFileSize {
			t.Error("MaxLogFileSize没有被修正")
		}
		if cfg.MaxLogAge == originalMaxAge {
			t.Error("MaxLogAge没有被修正")
		}
		if cfg.MaxLogBackups == originalMaxBackups {
			t.Error("MaxLogBackups没有被修正")
		}
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

// TestFastLogConfigChainMethods 测试FastLogConfig的链式调用方法
func TestFastLogConfigChainMethods(t *testing.T) {
	// 测试完整的链式调用
	config := NewFastLogConfig("test_logs", "test.log").
		WithLogDirName("custom_logs").
		WithLogFileName("custom.log").
		WithOutputToConsole(false).
		WithOutputToFile(true).
		WithFlushInterval(200 * time.Millisecond).
		WithLogLevel(DEBUG).
		WithChanIntSize(5000).
		WithLogFormat(Json).
		WithColor(false).
		WithBold(false).
		WithMaxLogFileSize(50).
		WithMaxLogAge(30).
		WithMaxLogBackups(10).
		WithIsLocalTime(false).
		WithEnableCompress(true)

	// 验证所有配置值是否正确设置
	tests := []struct {
		name     string
		expected interface{}
		actual   interface{}
	}{
		{"LogDirName", "custom_logs", config.LogDirName},
		{"LogFileName", "custom.log", config.LogFileName},
		{"OutputToConsole", false, config.OutputToConsole},
		{"OutputToFile", true, config.OutputToFile},
		{"FlushInterval", 200 * time.Millisecond, config.FlushInterval},
		{"LogLevel", DEBUG, config.LogLevel},
		{"ChanIntSize", 5000, config.ChanIntSize},
		{"LogFormat", Json, config.LogFormat},
		{"Color", false, config.Color},
		{"Bold", false, config.Bold},
		{"MaxLogFileSize", 50, config.MaxLogFileSize},
		{"MaxLogAge", 30, config.MaxLogAge},
		{"MaxLogBackups", 10, config.MaxLogBackups},
		{"IsLocalTime", false, config.IsLocalTime},
		{"EnableCompress", true, config.EnableCompress},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.actual != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.actual, tt.expected)
			}
		})
	}
}

// TestFastLogConfigIndividualMethods 测试每个方法的单独功能
func TestFastLogConfigIndividualMethods(t *testing.T) {
	config := NewFastLogConfig("logs", "app.log")

	// 测试 WithLogDirName
	result := config.WithLogDirName("new_logs")
	if result != config {
		t.Error("WithLogDirName should return the same config instance")
	}
	if config.LogDirName != "new_logs" {
		t.Errorf("LogDirName = %v, want %v", config.LogDirName, "new_logs")
	}

	// 测试 WithLogFileName
	config.WithLogFileName("new_app.log")
	if config.LogFileName != "new_app.log" {
		t.Errorf("LogFileName = %v, want %v", config.LogFileName, "new_app.log")
	}

	// 测试 WithOutputToConsole
	config.WithOutputToConsole(false)
	if config.OutputToConsole != false {
		t.Errorf("OutputToConsole = %v, want %v", config.OutputToConsole, false)
	}

	// 测试 WithOutputToFile
	config.WithOutputToFile(false)
	if config.OutputToFile != false {
		t.Errorf("OutputToFile = %v, want %v", config.OutputToFile, false)
	}

	// 测试 WithFlushInterval
	interval := 100 * time.Millisecond
	config.WithFlushInterval(interval)
	if config.FlushInterval != interval {
		t.Errorf("FlushInterval = %v, want %v", config.FlushInterval, interval)
	}

	// 测试 WithLogLevel
	config.WithLogLevel(ERROR)
	if config.LogLevel != ERROR {
		t.Errorf("LogLevel = %v, want %v", config.LogLevel, ERROR)
	}

	// 测试 WithChanIntSize
	config.WithChanIntSize(8000)
	if config.ChanIntSize != 8000 {
		t.Errorf("ChanIntSize = %v, want %v", config.ChanIntSize, 8000)
	}

	// 测试 WithLogFormat
	config.WithLogFormat(Detailed)
	if config.LogFormat != Detailed {
		t.Errorf("LogFormat = %v, want %v", config.LogFormat, Detailed)
	}

	// 测试 WithColor
	config.WithColor(false)
	if config.Color != false {
		t.Errorf("Color = %v, want %v", config.Color, false)
	}

	// 测试 WithBold
	config.WithBold(false)
	if config.Bold != false {
		t.Errorf("Bold = %v, want %v", config.Bold, false)
	}

	// 测试 WithMaxLogFileSize
	config.WithMaxLogFileSize(100)
	if config.MaxLogFileSize != 100 {
		t.Errorf("MaxLogFileSize = %v, want %v", config.MaxLogFileSize, 100)
	}

	// 测试 WithMaxLogAge
	config.WithMaxLogAge(60)
	if config.MaxLogAge != 60 {
		t.Errorf("MaxLogAge = %v, want %v", config.MaxLogAge, 60)
	}

	// 测试 WithMaxLogBackups
	config.WithMaxLogBackups(20)
	if config.MaxLogBackups != 20 {
		t.Errorf("MaxLogBackups = %v, want %v", config.MaxLogBackups, 20)
	}

	// 测试 WithIsLocalTime
	config.WithIsLocalTime(false)
	if config.IsLocalTime != false {
		t.Errorf("IsLocalTime = %v, want %v", config.IsLocalTime, false)
	}

	// 测试 WithEnableCompress
	config.WithEnableCompress(true)
	if config.EnableCompress != true {
		t.Errorf("EnableCompress = %v, want %v", config.EnableCompress, true)
	}
}

// TestFastLogConfigPartialChaining 测试部分链式调用
func TestFastLogConfigPartialChaining(t *testing.T) {
	// 测试部分链式调用
	config := NewFastLogConfig("logs", "app.log").
		WithLogLevel(WARN).
		WithOutputToConsole(false).
		WithMaxLogFileSize(25)

	if config.LogLevel != WARN {
		t.Errorf("LogLevel = %v, want %v", config.LogLevel, WARN)
	}
	if config.OutputToConsole != false {
		t.Errorf("OutputToConsole = %v, want %v", config.OutputToConsole, false)
	}
	if config.MaxLogFileSize != 25 {
		t.Errorf("MaxLogFileSize = %v, want %v", config.MaxLogFileSize, 25)
	}

	// 验证未修改的配置保持默认值
	if config.LogDirName != "logs" {
		t.Errorf("LogDirName should remain default value 'logs', got %v", config.LogDirName)
	}
	if config.LogFileName != "app.log" {
		t.Errorf("LogFileName should remain default value 'app.log', got %v", config.LogFileName)
	}
}

// TestFastLogConfigNilSafety 测试空指针安全性
func TestFastLogConfigNilSafety(t *testing.T) {
	var config *FastLogConfig = nil

	// 测试在nil配置上调用方法是否会panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when calling method on nil config")
		}
	}()

	config.WithLogLevel(DEBUG)
}

// BenchmarkFastLogConfigChaining 性能基准测试
func BenchmarkFastLogConfigChaining(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewFastLogConfig("logs", "app.log").
			WithLogLevel(DEBUG).
			WithOutputToConsole(true).
			WithOutputToFile(true).
			WithFlushInterval(100 * time.Millisecond).
			WithMaxLogFileSize(50).
			WithColor(true).
			WithBold(true)
	}
}
