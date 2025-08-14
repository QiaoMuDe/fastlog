package fastlog

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestGetCallerInfo 测试调用者信息获取
func TestGetCallerInfo(t *testing.T) {
	// 直接调用(skip=1)
	fileName, funcName, _, ok := getCallerInfo(1)
	if !ok {
		t.Fatal("获取调用者信息失败")
	}
	if fileName != "tools_test.go" || !strings.HasSuffix(funcName, "TestGetCallerInfo") {
		t.Errorf("调用者信息不匹配: 文件=%s, 函数=%s", fileName, funcName)
	}

	// 间接调用(skip=2)
	testGetCallerInfoHelper(t)
}

func testGetCallerInfoHelper(t *testing.T) {
	fileName, funcName, _, ok := getCallerInfo(2)
	if !ok {
		t.Fatal("获取调用者信息失败")
	}
	if fileName != "tools_test.go" || !strings.HasSuffix(funcName, "TestGetCallerInfo") {
		t.Errorf("间接调用者信息不匹配: 文件=%s, 函数=%s", fileName, funcName)
	}
}

// TestLogLevelToString 测试日志级别转换
func TestLogLevelToString(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{{
		level: DEBUG, expected: "DEBUG",
	}, {
		level: INFO, expected: "INFO",
	}, {
		level: WARN, expected: "WARN",
	}, {
		level: ERROR, expected: "ERROR",
	}, {
		level: FATAL, expected: "FATAL",
	}, {
		level: NONE, expected: "NONE",
	}, {
		level: 99, expected: "UNKNOWN",
	}}

	for _, tt := range tests {
		result := logLevelToString(tt.level)
		if result != tt.expected {
			t.Errorf("日志级别%d转换错误，预期%s，实际%s", tt.level, tt.expected, result)
		}
	}
}

// TestBackpressure 测试背压功能
func TestBackpressure(t *testing.T) {
	// 测试用例数据
	testCases := []struct {
		name         string
		channelLen   int      // 通道当前长度
		channelCap   int      // 通道容量
		level        LogLevel // 测试的日志级别
		expectedDrop bool     // 是否应该被丢弃
		description  string   // 测试描述
	}{
		// 正常情况（0-79%）
		{"正常_DEBUG", 0, 10, DEBUG, false, "通道空闲时，DEBUG日志应该保留"},
		{"正常_INFO", 5, 10, INFO, false, "通道50%时，INFO日志应该保留"},
		{"正常_WARN", 7, 10, WARN, false, "通道70%时，WARN日志应该保留"},

		// 80%背压（只保留INFO及以上）
		{"80%背压_DEBUG", 8, 10, DEBUG, true, "通道80%时，DEBUG日志应该被丢弃"},
		{"80%背压_INFO", 8, 10, INFO, false, "通道80%时，INFO日志应该保留"},
		{"80%背压_WARN", 8, 10, WARN, false, "通道80%时，WARN日志应该保留"},
		{"80%背压_ERROR", 8, 10, ERROR, false, "通道80%时，ERROR日志应该保留"},

		// 90%背压（只保留WARN及以上）
		{"90%背压_DEBUG", 9, 10, DEBUG, true, "通道90%时，DEBUG日志应该被丢弃"},
		{"90%背压_INFO", 9, 10, INFO, true, "通道90%时，INFO日志应该被丢弃"},
		{"90%背压_WARN", 9, 10, WARN, false, "通道90%时，WARN日志应该保留"},
		{"90%背压_ERROR", 9, 10, ERROR, false, "通道90%时，ERROR日志应该保留"},
		{"90%背压_FATAL", 9, 10, FATAL, false, "通道90%时，FATAL日志应该保留"},

		// 95%背压（只保留ERROR及以上）
		{"95%背压_DEBUG", 95, 100, DEBUG, true, "通道95%时，DEBUG日志应该被丢弃"},
		{"95%背压_INFO", 95, 100, INFO, true, "通道95%时，INFO日志应该被丢弃"},
		{"95%背压_WARN", 95, 100, WARN, true, "通道95%时，WARN日志应该被丢弃"},
		{"95%背压_ERROR", 95, 100, ERROR, false, "通道95%时，ERROR日志应该保留"},
		{"95%背压_FATAL", 95, 100, FATAL, false, "通道95%时，FATAL日志应该保留"},

		// 98%背压（只保留FATAL）
		{"98%背压_DEBUG", 98, 100, DEBUG, true, "通道98%时，DEBUG日志应该被丢弃"},
		{"98%背压_INFO", 98, 100, INFO, true, "通道98%时，INFO日志应该被丢弃"},
		{"98%背压_WARN", 98, 100, WARN, true, "通道98%时，WARN日志应该被丢弃"},
		{"98%背压_ERROR", 98, 100, ERROR, true, "通道98%时，ERROR日志应该被丢弃"},
		{"98%背压_FATAL", 98, 100, FATAL, false, "通道98%时，FATAL日志应该保留"},

		// 边界情况 - 保守策略：通道满时丢弃所有日志避免阻塞
		{"满通道_ERROR", 10, 10, ERROR, true, "通道满时，ERROR日志也应该被丢弃（保守策略）"},
		{"满通道_DEBUG", 10, 10, DEBUG, true, "通道满时，DEBUG日志应该被丢弃"},
		{"满通道_FATAL", 10, 10, FATAL, true, "通道满时，即使FATAL日志也应该被丢弃（保守策略）"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建预计算的背压阈值
			bp := &bpThresholds{
				threshold80: tc.channelCap * 80,
				threshold90: tc.channelCap * 90,
				threshold95: tc.channelCap * 95,
				threshold98: tc.channelCap * 98,
			}

			// 创建测试通道
			logChan := make(chan *logMsg, tc.channelCap)

			// 填充通道到指定长度
			for i := 0; i < tc.channelLen; i++ {
				logChan <- &logMsg{Level: INFO, Message: "test"}
			}

			// 测试背压函数（使用新的函数签名）
			shouldDrop := shouldDropLogByBackpressure(bp, logChan, tc.level)

			// 验证结果
			actualLen := len(logChan)
			actualCap := cap(logChan)

			if shouldDrop != tc.expectedDrop {
				t.Errorf("测试失败: %s\n"+
					"通道使用率: %d/%d (%.0f%%)\n"+
					"日志级别: %s\n"+
					"预期丢弃: %v, 实际丢弃: %v\n"+
					"描述: %s",
					tc.name,
					actualLen, actualCap, float64(actualLen)*100/float64(actualCap),
					logLevelToString(tc.level),
					tc.expectedDrop, shouldDrop,
					tc.description)
			} else {
				t.Logf("✓ %s - 通道使用率: %d%%, 级别: %s, 丢弃: %v",
					tc.name,
					actualLen*100/actualCap,
					logLevelToString(tc.level),
					shouldDrop)
			}

			// 清理通道
			close(logChan)
		})
	}
}

// TestBackpressureIntegration 集成测试：测试实际日志系统中的背压
func TestBackpressureIntegration(t *testing.T) {
	// 创建小容量的日志配置
	cfg := NewFastLogConfig("", "")
	cfg.ChanIntSize = 10 // 小通道容量
	cfg.OutputToConsole = true
	cfg.LogLevel = DEBUG

	log := NewFastLog(cfg)
	defer log.Close()

	// 快速发送大量日志，触发背压
	for i := 0; i < 50; i++ {
		log.Debug(fmt.Sprintf("Debug消息 %d", i))
		log.Info(fmt.Sprintf("Info消息 %d", i))
		log.Warn(fmt.Sprintf("Warn消息 %d", i))
		log.Error(fmt.Sprintf("Error消息 %d", i))
	}

	// 等待处理完成
	time.Sleep(100 * time.Millisecond)

	// 验证通道没有阻塞（这里主要是确保程序没有死锁）
	t.Log("背压集成测试完成，系统没有阻塞")
}

// BenchmarkBackpressureFunction 性能测试：测试背压函数的性能
func BenchmarkBackpressureFunction(b *testing.B) {
	logChan := make(chan *logMsg, 1000)
	bp := &bpThresholds{
		threshold80: 80,
		threshold90: 90,
		threshold95: 95,
		threshold98: 98,
	}

	// 填充通道到80%
	for i := 0; i < 800; i++ {
		logChan <- &logMsg{Level: INFO}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		shouldDropLogByBackpressure(bp, logChan, WARN)
	}
}
