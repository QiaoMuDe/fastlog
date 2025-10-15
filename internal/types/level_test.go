package types

import (
	"testing"
)

// TestShouldLog 测试ShouldLog函数的各种情况
func TestShouldLog(t *testing.T) {
	// 测试用例结构
	testCases := []struct {
		name     string
		logLevel LogLevel // 当前要记录的日志级别
		minLevel LogLevel // 配置的最低记录级别
		expected bool     // 期望的结果
	}{
		// DEBUG配置级别测试
		{"DEBUG配置下记录DEBUG级别", DEBUG_Mask, DEBUG, true},
		{"DEBUG配置下记录INFO级别", INFO_Mask, DEBUG, true},
		{"DEBUG配置下记录WARN级别", WARN_Mask, DEBUG, true},
		{"DEBUG配置下记录ERROR级别", ERROR_Mask, DEBUG, true},
		{"DEBUG配置下记录FATAL级别", FATAL_Mask, DEBUG, true},

		// INFO配置级别测试
		{"INFO配置下记录DEBUG级别", DEBUG_Mask, INFO, false},
		{"INFO配置下记录INFO级别", INFO_Mask, INFO, true},
		{"INFO配置下记录WARN级别", WARN_Mask, INFO, true},
		{"INFO配置下记录ERROR级别", ERROR_Mask, INFO, true},
		{"INFO配置下记录FATAL级别", FATAL_Mask, INFO, true},

		// WARN配置级别测试
		{"WARN配置下记录DEBUG级别", DEBUG_Mask, WARN, false},
		{"WARN配置下记录INFO级别", INFO_Mask, WARN, false},
		{"WARN配置下记录WARN级别", WARN_Mask, WARN, true},
		{"WARN配置下记录ERROR级别", ERROR_Mask, WARN, true},
		{"WARN配置下记录FATAL级别", FATAL_Mask, WARN, true},

		// ERROR配置级别测试
		{"ERROR配置下记录DEBUG级别", DEBUG_Mask, ERROR, false},
		{"ERROR配置下记录INFO级别", INFO_Mask, ERROR, false},
		{"ERROR配置下记录WARN级别", WARN_Mask, ERROR, false},
		{"ERROR配置下记录ERROR级别", ERROR_Mask, ERROR, true},
		{"ERROR配置下记录FATAL级别", FATAL_Mask, ERROR, true},

		// FATAL配置级别测试
		{"FATAL配置下记录DEBUG级别", DEBUG_Mask, FATAL, false},
		{"FATAL配置下记录INFO级别", INFO_Mask, FATAL, false},
		{"FATAL配置下记录WARN级别", WARN_Mask, FATAL, false},
		{"FATAL配置下记录ERROR级别", ERROR_Mask, FATAL, false},
		{"FATAL配置下记录FATAL级别", FATAL_Mask, FATAL, true},

		// NONE配置级别测试
		{"NONE配置下记录DEBUG级别", DEBUG_Mask, NONE, false},
		{"NONE配置下记录INFO级别", INFO_Mask, NONE, false},
		{"NONE配置下记录WARN级别", WARN_Mask, NONE, false},
		{"NONE配置下记录ERROR级别", ERROR_Mask, NONE, false},
		{"NONE配置下记录FATAL级别", FATAL_Mask, NONE, false},
	}

	// 执行测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ShouldLog(tc.logLevel, tc.minLevel)
			if result != tc.expected {
				t.Errorf("ShouldLog(%v, %v) = %v; expected %v", tc.logLevel, tc.minLevel, result, tc.expected)
			}
		})
	}
}

// TestLogLevelString 测试日志级别的字符串表示
func TestLogLevelString(t *testing.T) {
	testCases := []struct {
		level    LogLevel
		expected string
	}{
		{DEBUG, "DEBUG"},
		{INFO, "INFO"},
		{WARN, "WARN"},
		{ERROR, "ERROR"},
		{FATAL, "FATAL"},
		{NONE, "NONE"},
		{LogLevel(99), "UNKNOWN"}, // 无效的日志级别
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result := tc.level.String()
			if result != tc.expected {
				t.Errorf("LogLevel(%v).String() = %s; expected %s", tc.level, result, tc.expected)
			}
		})
	}
}

// TestLogLevelToPaddedString 测试带填充的日志级别字符串表示
func TestLogLevelToPaddedString(t *testing.T) {
	testCases := []struct {
		level    LogLevel
		expected string
	}{
		{DEBUG, "DEBUG"},
		{INFO, "INFO "},
		{WARN, "WARN "},
		{ERROR, "ERROR"},
		{FATAL, "FATAL"},
		{NONE, "NONE "},
		{LogLevel(99), "UNKNOWN"}, // 无效的日志级别
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result := LogLevelToPaddedString(tc.level)
			if result != tc.expected {
				t.Errorf("LogLevelToPaddedString(%v) = %s; expected %s", tc.level, result, tc.expected)
			}
		})
	}
}

// TestLogLevelMarshalJSON 测试日志级别的JSON序列化
func TestLogLevelMarshalJSON(t *testing.T) {
	testCases := []struct {
		level    LogLevel
		expected string
	}{
		{DEBUG, `"DEBUG"`},
		{INFO, `"INFO"`},
		{WARN, `"WARN"`},
		{ERROR, `"ERROR"`},
		{FATAL, `"FATAL"`},
		{NONE, `"NONE"`},
		{LogLevel(99), `"UNKNOWN"`}, // 无效的日志级别
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result, err := tc.level.MarshalJSON()
			if err != nil {
				t.Errorf("LogLevel(%v).MarshalJSON() returned error: %v", tc.level, err)
			}
			if string(result) != tc.expected {
				t.Errorf("LogLevel(%v).MarshalJSON() = %s; expected %s", tc.level, string(result), tc.expected)
			}
		})
	}
}
