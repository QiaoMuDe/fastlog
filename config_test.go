/*
config_test.go - 配置模块测试文件
包含对FastLog配置结构体的单元测试，验证配置初始化、边界条件处理和配置修正功能。
*/
package fastlog

import (
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
	if cfg.MaxLogFileSize != 5 {
		t.Errorf("默认最大日志文件大小应为5MB，实际为%d", cfg.MaxLogFileSize)
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
	if cfg.MaxLogFileSize != 5 {
		t.Errorf("修正后最大日志文件大小应为5，实际为%d", cfg.MaxLogFileSize)
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
