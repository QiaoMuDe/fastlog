package flog

import (
	"errors"
	"testing"
	"time"

	"gitee.com/MM-Q/fastlog/internal/config"
	"gitee.com/MM-Q/fastlog/internal/types"
)

func TestQM(t *testing.T) {
	log := NewFlog(config.ConsoleConfig())
	defer log.Close()

	log.Info("测试")
	log.Warn("测试")
	log.Error("测试")
	log.Debug("测试")

}

// TestLog 测试使用flog库记录不同级别的日志
func TestLog(t *testing.T) {
	// 创建开发环境配置
	cfg := config.DevConfig("example_logs", "app.log")

	// 也可以自定义配置
	cfg.LogLevel = types.DEBUG
	cfg.LogFormat = types.Json
	cfg.CallerInfo = true
	cfg.OutputToFile = true

	// 创建flog实例
	logger := NewFlog(cfg)
	defer logger.Close()

	// 记录不同级别的日志
	logger.Debug("这是一条调试日志")
	logger.Info("这是一条信息日志")
	logger.Warn("这是一条警告日志")
	logger.Error("这是一条错误日志")

	// 使用字段记录结构化日志
	logger.Info("用户登录",
		String("user_id", "12345"),
		String("username", "alice"),
		Int("age", 25),
		Bool("is_admin", false))

	// 记录带错误信息的日志
	err := errors.New("数据库连接失败")
	logger.Error("处理用户请求时发生错误",
		String("request_id", "req-12345"),
		String("url", "/api/users"),
		Error("error", err))

	// 记录带时间信息的日志
	logger.Info("任务执行完成",
		String("task_name", "data_backup"),
		Duration("duration", 5*time.Second),
		Int("records_processed", 1000))

	// 记录数值型字段
	logger.Info("系统状态监控",
		Int("cpu_usage", 75),
		Float64("memory_usage", 85.5),
		Uint64("disk_space", 1024000))

	// 使用不同日志格式
	jsonCfg := config.NewFastLogConfig("example_logs", "json.log")
	jsonCfg.LogFormat = types.Json
	jsonCfg.CallerInfo = true

	jsonLogger := NewFlog(jsonCfg)
	defer jsonLogger.Close()

	jsonLogger.Info("JSON格式日志示例",
		String("service", "user-service"),
		Int("status_code", 200))
}

// TestFlogUsage 测试flog库的基本使用方法
func TestFlogUsage(t *testing.T) {
	// 创建测试配置 - 使用控制台模式避免文件写入
	cfg := config.ConsoleConfig()
	cfg.LogFormat = types.Json

	// 创建flog实例
	logger := NewFlog(cfg)
	if logger == nil {
		t.Fatal("Failed to create flog instance")
	}
	defer logger.Close()

	// 测试不同级别的日志记录
	t.Run("BasicLogLevels", func(t *testing.T) {
		logger.Debug("这是一条调试日志")
		logger.Info("这是一条信息日志")
		logger.Warn("这是一条警告日志")
		logger.Error("这是一条错误日志")
	})

	// 测试带字段的日志记录
	t.Run("LogsWithFields", func(t *testing.T) {
		logger.Info("用户登录",
			String("user_id", "12345"),
			String("username", "alice"),
			Int("age", 25),
			Bool("is_admin", false))
	})

	// 测试带错误信息的日志记录
	t.Run("LogsWithError", func(t *testing.T) {
		err := errors.New("数据库连接失败")
		logger.Error("处理用户请求时发生错误",
			String("request_id", "req-12345"),
			String("url", "/api/users"),
			Error("error", err))
	})

	// 测试带时间信息的日志记录
	t.Run("LogsWithTimeInfo", func(t *testing.T) {
		logger.Info("任务执行完成",
			String("task_name", "data_backup"),
			Duration("duration", 5*time.Second),
			Int("records_processed", 1000))
	})

	// 测试数值型字段
	t.Run("LogsWithNumericFields", func(t *testing.T) {
		logger.Info("系统状态监控",
			Int("cpu_usage", 75),
			Float64("memory_usage", 85.5),
			Uint64("disk_space", 1024000))
	})
}

// TestFlogWithDifferentFormats 测试不同日志格式
func TestFlogWithDifferentFormats(t *testing.T) {
	formats := []types.LogFormatType{
		types.Def,
		types.Json,
		types.Structured,
		types.Timestamp,
	}

	for _, format := range formats {
		t.Run(format.String(), func(t *testing.T) {
			// 创建测试配置
			cfg := config.ConsoleConfig()
			cfg.LogFormat = format

			// 创建flog实例
			logger := NewFlog(cfg)
			if logger == nil {
				t.Fatal("Failed to create flog instance")
			}
			defer logger.Close()

			// 记录带字段的日志
			logger.Info("测试不同日志格式",
				String("service", "test-service"),
				Int("status_code", 200),
				String("user_id", "12345"))
		})
	}
}

// TestFlogConfigurations 测试不同的配置选项
func TestFlogConfigurations(t *testing.T) {
	// 测试启用调用者信息
	t.Run("WithCallerInfo", func(t *testing.T) {
		cfg := config.ConsoleConfig()
		cfg.CallerInfo = true

		logger := NewFlog(cfg)
		if logger == nil {
			t.Fatal("Failed to create flog instance")
		}
		defer logger.Close()

		logger.Info("启用调用者信息的日志")
	})

	// 测试不同日志级别
	t.Run("DifferentLogLevels", func(t *testing.T) {
		levels := []types.LogLevel{
			types.DEBUG,
			types.INFO,
			types.WARN,
			types.ERROR,
		}

		for _, level := range levels {
			t.Run(level.String(), func(t *testing.T) {
				cfg := config.ConsoleConfig()
				cfg.LogLevel = level

				logger := NewFlog(cfg)
				if logger == nil {
					t.Fatal("Failed to create flog instance")
				}
				defer logger.Close()

				// 尝试记录所有级别的日志
				logger.Debug("调试信息")
				logger.Info("一般信息")
				logger.Warn("警告信息")
				logger.Error("错误信息")
			})
		}
	})
}

// BenchmarkFlogInfo 简单的性能基准测试
func BenchmarkFlogInfo(b *testing.B) {
	// 创建测试配置
	cfg := config.ConsoleConfig()
	cfg.OutputToConsole = false // 关闭控制台输出以提高性能测试准确性

	// 创建flog实例
	logger := NewFlog(cfg)
	if logger == nil {
		b.Fatal("Failed to create flog instance")
	}
	defer logger.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("性能测试日志")
	}
}

// BenchmarkFlogInfoWithFields 带字段的性能基准测试
func BenchmarkFlogInfoWithFields(b *testing.B) {
	// 创建测试配置
	cfg := config.ConsoleConfig()
	cfg.OutputToConsole = false // 关闭控制台输出以提高性能测试准确性

	// 创建flog实例
	logger := NewFlog(cfg)
	if logger == nil {
		b.Fatal("Failed to create flog instance")
	}
	defer logger.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("带字段的性能测试日志",
			String("key1", "value1"),
			Int("key2", 123),
			Bool("key3", true))
	}
}
