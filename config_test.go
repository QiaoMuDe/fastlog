package fastlog

import (
	"testing"
	"time"
)

// helper: 断言调用发生 panic
func mustPanic(t *testing.T, name string, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("%s: 预期发生 panic，但未发生", name)
		}
	}()
	fn()
}

// -----------------------------------------------------------------------------
// NewFastLogConfig 默认值测试
// -----------------------------------------------------------------------------
func TestNewFastLogConfig_Defaults(t *testing.T) {
	cfg := NewFastLogConfig("logs", "app.log")

	// 布尔与枚举默认值
	if !cfg.OutputToConsole {
		t.Error("默认应启用控制台输出")
	}
	if !cfg.OutputToFile {
		t.Error("默认应启用文件输出")
	}
	if cfg.LogLevel != INFO {
		t.Errorf("默认 LogLevel 应为 INFO，得到 %v", cfg.LogLevel)
	}
	if cfg.LogFormat != Simple {
		t.Errorf("默认 LogFormat 应为 Simple，得到 %v", cfg.LogFormat)
	}
	if !cfg.Color {
		t.Error("默认 Color 应为 true")
	}
	if !cfg.Bold {
		t.Error("默认 Bold 应为 true")
	}
	if !cfg.LocalTime {
		t.Error("默认 LocalTime 应为 true")
	}
	if cfg.Compress {
		t.Error("默认 Compress 应为 false")
	}

	// 其他字段基本合理性（不校验具体常量值）
	if cfg.LogDirName != "logs" {
		t.Errorf("LogDirName 期望 'logs'，得到 %q", cfg.LogDirName)
	}
	if cfg.LogFileName != "app.log" {
		t.Errorf("LogFileName 期望 'app.log'，得到 %q", cfg.LogFileName)
	}
	if cfg.MaxSize < 0 {
		t.Error("MaxSize 不应为负数")
	}
	if cfg.MaxBufferSize <= 0 {
		t.Error("MaxBufferSize 应为正数")
	}
	if cfg.MaxWriteCount <= 0 {
		t.Error("MaxWriteCount 应为正数")
	}
	if cfg.FlushInterval != defaultFlushInterval {
		t.Errorf("默认 FlushInterval 应为 %v，得到 %v", defaultFlushInterval, cfg.FlushInterval)
	}
}

// -----------------------------------------------------------------------------
// 便捷模式测试：Dev / Prod / Console
// -----------------------------------------------------------------------------
func TestDevConfig(t *testing.T) {
	cfg := DevConfig("logs", "dev.log")

	if cfg.LogLevel != DEBUG {
		t.Errorf("DevConfig: LogLevel 应为 DEBUG，得到 %v", cfg.LogLevel)
	}
	if cfg.LogFormat != Detailed {
		t.Errorf("DevConfig: LogFormat 应为 Detailed，得到 %v", cfg.LogFormat)
	}
	if cfg.MaxFiles != 5 {
		t.Errorf("DevConfig: MaxFiles 应为 5，得到 %d", cfg.MaxFiles)
	}
	if cfg.MaxAge != 7 {
		t.Errorf("DevConfig: MaxAge 应为 7，得到 %d", cfg.MaxAge)
	}
	// Dev 保持默认的输出设置（来自 NewFastLogConfig）
	if !cfg.OutputToConsole || !cfg.OutputToFile {
		t.Error("DevConfig: 默认应同时输出到控制台与文件")
	}
}

func TestProdConfig(t *testing.T) {
	cfg := ProdConfig("logs", "prod.log")

	if cfg.Compress != true {
		t.Error("ProdConfig: Compress 应启用")
	}
	if cfg.OutputToConsole != false {
		t.Error("ProdConfig: 应禁用控制台输出")
	}
	if cfg.MaxAge != 30 {
		t.Errorf("ProdConfig: MaxAge 应为 30，得到 %d", cfg.MaxAge)
	}
	if cfg.MaxFiles != 24 {
		t.Errorf("ProdConfig: MaxFiles 应为 24，得到 %d", cfg.MaxFiles)
	}
	// 保持文件输出启用（来源默认）
	if !cfg.OutputToFile {
		t.Error("ProdConfig: 应启用文件输出")
	}
}

func TestConsoleConfig(t *testing.T) {
	cfg := ConsoleConfig()

	if cfg.OutputToFile != false {
		t.Error("ConsoleConfig: 应禁用文件输出")
	}
	if cfg.LogLevel != DEBUG {
		t.Errorf("ConsoleConfig: LogLevel 应为 DEBUG，得到 %v", cfg.LogLevel)
	}
	// 目录与文件名在控制台模式可为空，但 validateConfig 不会校验（因为未输出到文件）
	if cfg.LogDirName != "" || cfg.LogFileName != "" {
		t.Errorf("ConsoleConfig: 目录与文件名应为空占位，得到 dir=%q file=%q", cfg.LogDirName, cfg.LogFileName)
	}
}

// -----------------------------------------------------------------------------
// validateConfig 异常场景测试
// -----------------------------------------------------------------------------
func TestValidateConfig_OutputTargetsRequired(t *testing.T) {
	cfg := NewFastLogConfig("logs", "app.log")
	cfg.OutputToConsole = false
	cfg.OutputToFile = false
	mustPanic(t, "输出目标至少启用一个", func() {
		cfg.validateConfig()
	})
}

func TestValidateConfig_InvalidLevels(t *testing.T) {
	cfg := NewFastLogConfig("logs", "app.log")

	// 非法级别（下界：小于 DEBUG）
	cfg.LogLevel = LogLevel(0)
	mustPanic(t, "非法 LogLevel（下界）", func() {
		cfg.validateConfig()
	})
}

func TestValidateConfig_InvalidFormats(t *testing.T) {
	cfg := NewFastLogConfig("logs", "app.log")

	// 非法格式（上界：大于 Custom）
	cfg.LogFormat = Custom + 1
	mustPanic(t, "非法 LogFormat（上界）", func() {
		cfg.validateConfig()
	})
}

func TestValidateConfig_FileOutputRequirements(t *testing.T) {
	cfg := NewFastLogConfig("logs", "app.log")
	cfg.OutputToFile = true // 确保触发文件相关校验

	// 空目录
	cfg.LogDirName = "   "
	mustPanic(t, "启用文件输出时目录不能为空", func() {
		cfg.validateConfig()
	})
	cfg.LogDirName = "logs"

	// 空文件名
	cfg.LogFileName = " "
	mustPanic(t, "启用文件输出时文件名不能为空", func() {
		cfg.validateConfig()
	})
	cfg.LogFileName = "app.log"

	// 路径穿越检测
	cfg.LogDirName = "../logs"
	mustPanic(t, "目录包含路径穿越", func() {
		cfg.validateConfig()
	})
	cfg.LogDirName = "logs"

	cfg.LogFileName = "app../log"
	mustPanic(t, "文件名包含路径穿越", func() {
		cfg.validateConfig()
	})
	cfg.LogFileName = "app.log"
}

func TestValidateConfig_ValidExample(t *testing.T) {
	cfg := NewFastLogConfig("logs", "app.log")
	// 一个典型的生产配置示例（应通过）
	cfg.OutputToConsole = false
	cfg.OutputToFile = true
	cfg.LogLevel = INFO
	cfg.LogFormat = Structured
	cfg.MaxSize = 100
	cfg.MaxAge = 30
	cfg.MaxFiles = 24
	cfg.LocalTime = true
	cfg.Compress = true
	cfg.MaxBufferSize = defaultMaxBufferSize
	cfg.MaxWriteCount = defaultMaxWriteCount
	cfg.FlushInterval = 1 * time.Second

	// 不应 panic
	cfg.validateConfig()
}
