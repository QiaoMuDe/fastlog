/*
config_test.go - 日志配置模块测试文件
包含对FastLog配置结构体的单元测试，验证配置初始化、边界条件处理和配置验证功能，
确保配置参数在各种异常情况下都能正确验证和设置默认值。
*/
package fastlog

import (
	"strings"
	"testing"
	"time"
)

// ========================================================================
// 基础配置创建测试
// ========================================================================

// TestNewFastLogConfig 测试基础配置创建
func TestNewFastLogConfig(t *testing.T) {
	cfg := NewFastLogConfig("logs", "test.log")

	// 验证基本字段设置
	if cfg.LogDirName != "logs" {
		t.Errorf("LogDirName expected 'logs', got '%s'", cfg.LogDirName)
	}
	if cfg.LogFileName != "test.log" {
		t.Errorf("LogFileName expected 'test.log', got '%s'", cfg.LogFileName)
	}
	if !cfg.OutputToConsole {
		t.Error("OutputToConsole should be true by default")
	}
	if !cfg.OutputToFile {
		t.Error("OutputToFile should be true by default")
	}
	if cfg.LogLevel != INFO {
		t.Errorf("LogLevel expected %d (INFO), got %d", INFO, cfg.LogLevel)
	}
	if cfg.ChanIntSize != defaultChanSize {
		t.Errorf("ChanIntSize expected %d, got %d", defaultChanSize, cfg.ChanIntSize)
	}
	if cfg.FlushInterval != normalFlushInterval {
		t.Errorf("FlushInterval expected %v, got %v", normalFlushInterval, cfg.FlushInterval)
	}
	if cfg.MaxSize != defaultMaxFileSize {
		t.Errorf("MaxSize expected %d, got %d", defaultMaxFileSize, cfg.MaxSize)
	}
}

// ========================================================================
// 预设配置模式测试
// ========================================================================

// TestDevConfig 测试开发模式配置
func TestDevConfig(t *testing.T) {
	cfg := DevConfig("dev_logs", "dev.log")

	if cfg.LogLevel != DEBUG {
		t.Errorf("Dev mode LogLevel expected %d (DEBUG), got %d", DEBUG, cfg.LogLevel)
	}
	if cfg.FlushInterval != fastFlushInterval {
		t.Errorf("Dev mode FlushInterval expected %v, got %v", fastFlushInterval, cfg.FlushInterval)
	}
	if cfg.LogFormat != Detailed {
		t.Errorf("Dev mode LogFormat expected %d (Detailed), got %d", Detailed, cfg.LogFormat)
	}
	if cfg.MaxAge != developmentMaxAge {
		t.Errorf("Dev mode MaxAge expected %d, got %d", developmentMaxAge, cfg.MaxAge)
	}
	if !cfg.OutputToConsole || !cfg.OutputToFile {
		t.Error("Dev mode should enable both console and file output")
	}
}

// TestProdConfig 测试生产模式配置
func TestProdConfig(t *testing.T) {
	cfg := ProdConfig("prod_logs", "prod.log")

	if cfg.OutputToConsole {
		t.Error("Prod mode should disable console output")
	}
	if !cfg.OutputToFile {
		t.Error("Prod mode should enable file output")
	}
	if cfg.ChanIntSize != largeChanSize {
		t.Errorf("Prod mode ChanIntSize expected %d, got %d", largeChanSize, cfg.ChanIntSize)
	}
	if cfg.FlushInterval != slowFlushInterval {
		t.Errorf("Prod mode FlushInterval expected %v, got %v", slowFlushInterval, cfg.FlushInterval)
	}
	if cfg.LogFormat != Json {
		t.Errorf("Prod mode LogFormat expected %d (Json), got %d", Json, cfg.LogFormat)
	}
	if cfg.Color || cfg.Bold {
		t.Error("Prod mode should disable color and bold")
	}
	if !cfg.Compress {
		t.Error("Prod mode should enable compression")
	}
}

// TestConsoleConfig 测试控制台模式配置
func TestConsoleConfig(t *testing.T) {
	cfg := ConsoleConfig()

	if !cfg.OutputToConsole {
		t.Error("Console mode should enable console output")
	}
	if cfg.OutputToFile {
		t.Error("Console mode should disable file output")
	}
	if cfg.ChanIntSize != smallChanSize {
		t.Errorf("Console mode ChanIntSize expected %d, got %d", smallChanSize, cfg.ChanIntSize)
	}
	if cfg.FlushInterval != fastFlushInterval {
		t.Errorf("Console mode FlushInterval expected %v, got %v", fastFlushInterval, cfg.FlushInterval)
	}
}

// TestFileConfig 测试文件模式配置
func TestFileConfig(t *testing.T) {
	cfg := FileConfig("file_logs", "file.log")

	if cfg.OutputToConsole {
		t.Error("File mode should disable console output")
	}
	if !cfg.OutputToFile {
		t.Error("File mode should enable file output")
	}
	if cfg.LogFormat != BasicStructured {
		t.Errorf("File mode LogFormat expected %d (BasicStructured), got %d", BasicStructured, cfg.LogFormat)
	}
	if cfg.Color || cfg.Bold {
		t.Error("File mode should disable color and bold")
	}
}

// ========================================================================
// 配置验证测试 - 正常情况
// ========================================================================

// TestValidateConfig_ValidConfigs 测试有效配置通过验证
func TestValidateConfig_ValidConfigs(t *testing.T) {
	testCases := []struct {
		name   string
		config *FastLogConfig
	}{
		{
			name: "默认配置",
			config: &FastLogConfig{
				OutputToConsole: true,
				OutputToFile:    true,
				LogDirName:      "logs",
				LogFileName:     "test.log",
			},
		},
		{
			name: "仅控制台输出",
			config: &FastLogConfig{
				OutputToConsole: true,
				OutputToFile:    false,
			},
		},
		{
			name: "自定义有效值",
			config: &FastLogConfig{
				OutputToConsole: true,
				OutputToFile:    true,
				LogDirName:      "custom_logs",
				LogFileName:     "custom.log",
				ChanIntSize:     5000,
				FlushInterval:   100 * time.Millisecond,
				BatchSize:       500,
				LogLevel:        DEBUG,
				LogFormat:       Json,
				MaxSize:         100,
				MaxAge:          30,
				MaxFiles:        10,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Valid config should not panic: %v", r)
				}
			}()
			tc.config.validateConfig()
		})
	}
}

// ========================================================================
// 配置验证测试 - 错误情况
// ========================================================================

// TestValidateConfig_InvalidConfigs 测试无效配置触发panic
func TestValidateConfig_InvalidConfigs(t *testing.T) {
	testCases := []struct {
		name        string
		setupConfig func() *FastLogConfig
		expectedMsg string
	}{
		{
			name: "配置对象为nil",
			setupConfig: func() *FastLogConfig {
				return nil
			},
			expectedMsg: "cannot be nil",
		},
		{
			name: "未启用任何输出方式",
			setupConfig: func() *FastLogConfig {
				return &FastLogConfig{
					OutputToConsole: false,
					OutputToFile:    false,
				}
			},
			expectedMsg: "at least one output method must be enabled",
		},
		{
			name: "ChanIntSize为负数",
			setupConfig: func() *FastLogConfig {
				cfg := NewFastLogConfig("logs", "test.log")
				cfg.ChanIntSize = -1
				return cfg
			},
			expectedMsg: "cannot be negative",
		},
		{
			name: "ChanIntSize超过最大值",
			setupConfig: func() *FastLogConfig {
				cfg := NewFastLogConfig("logs", "test.log")
				cfg.ChanIntSize = maxChanSize + 1
				return cfg
			},
			expectedMsg: "exceeds maximum",
		},
		{
			name: "MaxLogFileSize超过最大值",
			setupConfig: func() *FastLogConfig {
				cfg := NewFastLogConfig("logs", "test.log")
				cfg.MaxSize = maxSingleFileSize + 1
				return cfg
			},
			expectedMsg: "exceeds maximum",
		},
		{
			name: "LogDirName包含路径遍历",
			setupConfig: func() *FastLogConfig {
				return &FastLogConfig{
					OutputToConsole: false,
					OutputToFile:    true,
					LogDirName:      "../../../etc",
					LogFileName:     "test.log",
				}
			},
			expectedMsg: "path traversal",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := tc.setupConfig()

			defer func() {
				if r := recover(); r != nil {
					panicStr := r.(string)
					if !strings.Contains(panicStr, tc.expectedMsg) {
						t.Errorf("Expected panic message to contain '%s', got '%s'", tc.expectedMsg, panicStr)
					}
				} else {
					t.Error("Expected panic but none occurred")
				}
			}()

			cfg.validateConfig()
		})
	}
}

// ========================================================================
// 默认值设置测试
// ========================================================================

// TestValidateConfig_DefaultValues 测试默认值设置
func TestValidateConfig_DefaultValues(t *testing.T) {
	cfg := &FastLogConfig{
		OutputToConsole: true,
		OutputToFile:    true,
		LogDirName:      "logs",
		LogFileName:     "test.log",
		// 其他字段保持零值，测试默认值设置
	}

	cfg.validateConfig()

	// 验证默认值是否正确设置
	if cfg.ChanIntSize != defaultChanSize {
		t.Errorf("ChanIntSize default expected %d, got %d", defaultChanSize, cfg.ChanIntSize)
	}
	if cfg.FlushInterval != normalFlushInterval {
		t.Errorf("FlushInterval default expected %v, got %v", normalFlushInterval, cfg.FlushInterval)
	}
	if cfg.BatchSize != defaultBatchSize {
		t.Errorf("BatchSize default expected %d, got %d", defaultBatchSize, cfg.BatchSize)
	}
	if cfg.LogLevel != INFO {
		t.Errorf("LogLevel default expected %d (INFO), got %d", INFO, cfg.LogLevel)
	}
	if cfg.LogFormat != Simple {
		t.Errorf("LogFormat default expected %d (Simple), got %d", Simple, cfg.LogFormat)
	}
	if cfg.MaxSize != defaultMaxFileSize {
		t.Errorf("MaxSize default expected %d, got %d", defaultMaxFileSize, cfg.MaxSize)
	}
}
