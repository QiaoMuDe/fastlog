package flog

import (
	"strings"
	"testing"
	"time"

	"gitee.com/MM-Q/fastlog/internal/config"
	"gitee.com/MM-Q/fastlog/internal/types"
)

// TestBuildLogFormats 测试buildLog函数在不同格式下的输出是否符合预期
func TestBuildLogFormats(t *testing.T) {
	// 创建测试配置
	cfg := config.NewFastLogConfig("test_logs", "test.log")

	// 创建测试时间戳
	testTime := time.Now().Format("2006-01-02 15:04:05")

	// 创建测试字段
	fields := []*Field{
		String("user_id", "12345"),
		Int("age", 25),
		String("action", "login"),
	}

	// 使用NewEntry构造函数创建测试Entry
	entry := NewEntry(true, types.INFO, "用户登录成功", fields...)

	// 定义支持的日志格式（目前仅支持Json和JsonSimple）
	formats := []struct {
		formatType  types.LogFormatType
		name        string
		contains    []string // 期望包含的字符串片段
		notContains []string // 期望不包含的字符串片段
		caller      bool     // 是否包含caller信息
	}{
		{
			formatType:  types.Json,
			name:        "JSON格式",
			caller:      true,
			contains:    []string{"{", `"time":"` + testTime, `"level":"INFO"`, `"msg":"用户登录成功"`, `"user_id":"12345"`, `"age":"25"`, `"action":"login"`, "}"},
			notContains: []string{},
		},
		{
			formatType:  types.Json,
			name:        "Json格式(无caller)",
			caller:      false,
			contains:    []string{"{", `"time":"` + testTime, `"level":"INFO"`, `"msg":"用户登录成功"`, `"user_id":"12345"`, `"age":"25"`, `"action":"login"`, "}"},
			notContains: []string{},
		},
	}

	// 遍历所有格式进行测试
	for _, format := range formats {
		t.Run(format.name, func(t *testing.T) {
			// 设置当前测试的日志格式
			cfg.LogFormat = format.formatType

			// 设置是否包含caller信息
			cfg.CallerInfo = format.caller

			// 调用buildLog函数
			result := buildLog(cfg, entry)
			resultStr := string(result)

			t.Logf("格式: %s", format.name)
			t.Logf("输出: %s", resultStr)

			// 检查是否包含期望的字符串片段
			for _, expected := range format.contains {
				if !strings.Contains(resultStr, expected) {
					t.Errorf("格式 %s 的输出中未包含期望的字符串: %s", format.name, expected)
				}
			}

			// 检查是否不包含不应该包含的字符串片段
			for _, notExpected := range format.notContains {
				if strings.Contains(resultStr, notExpected) {
					t.Errorf("格式 %s 的输出中不应该包含字符串: %s", format.name, notExpected)
				}
			}

			// 特殊检查：确保JSON格式是合法的
			if format.formatType == types.Json {
				if !strings.HasPrefix(resultStr, "{") || !strings.HasSuffix(resultStr, "}") {
					t.Errorf("格式 %s 的输出不是合法的JSON格式", format.name)
				}
			}
		})
	}
}

// TestBuildLogEdgeCases 测试buildLog函数的边界情况
func TestBuildLogEdgeCases(t *testing.T) {
	cfg := config.NewFastLogConfig("test_logs", "test.log")

	// 测试空配置
	t.Run("空配置", func(t *testing.T) {
		result := buildLog(nil, &Entry{})
		if len(result) != 0 {
			t.Errorf("空配置应该返回空字节数组，实际返回: %s", string(result))
		}
	})

	// 测试空Entry
	t.Run("空Entry", func(t *testing.T) {
		result := buildLog(cfg, nil)
		if len(result) != 0 {
			t.Errorf("空Entry应该返回空字节数组，实际返回: %s", string(result))
		}
	})

	// 测试无字段的情况
	t.Run("无字段", func(t *testing.T) {
		entry := NewEntry(true, types.INFO, "测试消息")

		// 测试Json格式
		cfg.LogFormat = types.Json
		cfg.CallerInfo = true
		result := buildLog(cfg, entry)
		resultStr := string(result)

		t.Logf("无字段JSON输出: %s", resultStr)

		// 应该包含基本字段，但不包含额外的字段
		expectedParts := []string{"time", "level", "msg"}
		for _, part := range expectedParts {
			if !strings.Contains(resultStr, `"`+part+`":"`) {
				t.Errorf("无字段JSON中缺少基本字段: %s", part)
			}
		}

		// 检查是否包含caller字段（因为NewEntry(true, ...)会包含caller信息）
		if !strings.Contains(resultStr, `"caller"`) {
			t.Errorf("无字段JSON中应该包含caller字段")
		}
	})

	// 测试包含特殊字符的消息
	t.Run("特殊字符消息", func(t *testing.T) {
		specialFields := []*Field{
			String("path", `/usr/local/bin/app`),
			String("regex", `\d{3}-\d{3}-\d{4}`),
		}
		entry := NewEntry(true, types.ERROR, `错误消息包含"引号"和\反斜杠`, specialFields...)

		cfg.LogFormat = types.Json
		result := buildLog(cfg, entry)
		resultStr := string(result)

		t.Logf("特殊字符JSON输出: %s", resultStr)

		// 检查特殊字符是否正确转义
		if !strings.Contains(resultStr, `错误消息包含\"引号\"和\\反斜杠`) {
			t.Errorf("JSON中的特殊字符未正确转义")
		}
	})

	// 测试不同日志级别
	t.Run("不同日志级别", func(t *testing.T) {
		levels := []struct {
			level    types.LogLevel
			levelStr string
		}{
			{types.DEBUG, "DEBUG"},
			{types.INFO, "INFO"},
			{types.WARN, "WARN"},
			{types.ERROR, "ERROR"},
			{types.FATAL, "FATAL"},
		}

		for _, level := range levels {
			t.Run(level.levelStr, func(t *testing.T) {
				entry := NewEntry(true, level.level, "测试消息")

				cfg.LogFormat = types.Json
				result := buildLog(cfg, entry)
				resultStr := string(result)

				if !strings.Contains(resultStr, level.levelStr) {
					t.Errorf("日志级别 %s 未正确显示在输出中", level.levelStr)
				}
			})
		}
	})
}

// TestBuildLogPerformance 简单的性能测试
func TestBuildLogPerformance(t *testing.T) {
	cfg := config.NewFastLogConfig("test_logs", "test.log")
	cfg.LogFormat = types.Json

	performanceFields := []*Field{
		String("user_id", "12345"),
		Int("request_count", 1000),
		Float64("response_time", 123.456),
		Bool("success", true),
	}
	entry := NewEntry(true, types.INFO, "性能测试消息", performanceFields...)

	// 预热
	for i := 0; i < 100; i++ {
		_ = buildLog(cfg, entry)
	}

	// 正式测试
	start := time.Now()
	iterations := 10000
	for i := 0; i < iterations; i++ {
		_ = buildLog(cfg, entry)
	}
	duration := time.Since(start)

	t.Logf("性能测试: %d 次调用，耗时: %v，平均每次: %v", iterations, duration, duration/time.Duration(iterations))

	// 基本性能要求：10000次调用应该在1秒内完成（宽松要求）
	if duration > time.Second {
		t.Errorf("性能测试失败: 10000次调用耗时 %v，超过1秒", duration)
	}
}
