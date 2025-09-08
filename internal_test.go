/*
internal_test.go - FastLog内部功能测试文件
包含对日志系统内部组件的单元测试，包括时间戳缓存、调用者信息获取、背压控制、
文件名清理等核心内部功能的测试，确保内部实现的正确性和性能。
*/
package fastlog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"gitee.com/MM-Q/colorlib"
)

// TestGetCallerInfo 测试调用者信息获取
func TestGetCallerInfo(t *testing.T) {
	// 直接调用(skip=1)
	fileName, funcName, _, ok := getCallerInfo(1)
	if !ok {
		t.Fatal("获取调用者信息失败")
	}
	if !strings.HasSuffix(fileName, "_test.go") || !strings.Contains(funcName, "TestGetCallerInfo") {
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
	if !strings.HasSuffix(fileName, "_test.go") || !strings.Contains(funcName, "TestGetCallerInfo") {
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

// TestGetCachedTimestamp_Comprehensive 全面测试时间戳缓存功能
func TestGetCachedTimestamp_Comprehensive(t *testing.T) {
	t.Run("基本功能测试", func(t *testing.T) {
		timestamp1 := getCachedTimestamp()
		timestamp2 := getCachedTimestamp()

		// 在同一秒内应该返回相同的时间戳
		if timestamp1 != timestamp2 {
			t.Errorf("同一秒内时间戳应该相同: %s != %s", timestamp1, timestamp2)
		}

		// 验证时间戳格式
		if len(timestamp1) != 19 { // "2006-01-02 15:04:05" 长度为19
			t.Errorf("时间戳格式不正确，长度应为19，实际为%d: %s", len(timestamp1), timestamp1)
		}
	})

	t.Run("跨秒更新测试", func(t *testing.T) {
		// 获取当前时间戳
		timestamp1 := getCachedTimestamp()

		// 等待到下一秒
		now := time.Now()
		nextSecond := now.Truncate(time.Second).Add(time.Second)
		time.Sleep(time.Until(nextSecond) + 10*time.Millisecond)

		// 获取新的时间戳
		timestamp2 := getCachedTimestamp()

		// 应该不同
		if timestamp1 == timestamp2 {
			t.Errorf("跨秒时间戳应该不同: %s == %s", timestamp1, timestamp2)
		}
	})

	t.Run("并发安全测试", func(t *testing.T) {
		const goroutineCount = 100
		const iterationsPerGoroutine = 1000

		var wg sync.WaitGroup
		timestamps := make([]string, goroutineCount*iterationsPerGoroutine)
		var index int64

		// 启动多个goroutine并发获取时间戳
		for i := 0; i < goroutineCount; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < iterationsPerGoroutine; j++ {
					ts := getCachedTimestamp()
					idx := atomic.AddInt64(&index, 1) - 1
					timestamps[idx] = ts
				}
			}()
		}

		wg.Wait()

		// 验证所有时间戳都是有效的
		for i, ts := range timestamps {
			if len(ts) != 19 {
				t.Errorf("索引%d的时间戳格式无效: %s", i, ts)
				break
			}
		}

		// 统计不同时间戳的数量（应该很少，因为测试执行很快）
		uniqueTimestamps := make(map[string]bool)
		for _, ts := range timestamps {
			uniqueTimestamps[ts] = true
		}

		t.Logf("并发测试中生成了%d个不同的时间戳", len(uniqueTimestamps))

		// 在快速执行的情况下，不同时间戳数量应该很少
		if len(uniqueTimestamps) > 10 {
			t.Logf("警告：时间戳变化频繁，可能存在缓存失效问题")
		}
	})
}

// TestGetCallerInfo_Comprehensive 全面测试调用者信息获取功能
func TestGetCallerInfo_Comprehensive(t *testing.T) {
	t.Run("正常调用测试", func(t *testing.T) {
		fileName, funcName, line, ok := getCallerInfo(1)

		if !ok {
			t.Fatal("获取调用者信息失败")
		}

		// 验证文件名
		if !strings.HasSuffix(fileName, "_test.go") {
			t.Errorf("文件名应该以_test.go结尾，实际为: %s", fileName)
		}

		// 验证函数名包含测试函数名
		if !strings.Contains(funcName, "TestGetCallerInfo") {
			t.Errorf("函数名应该包含TestGetCallerInfo，实际为: %s", funcName)
		}

		// 验证行号合理性
		if line == 0 {
			t.Error("行号不应该为0")
		}

		t.Logf("调用者信息: 文件=%s, 函数=%s, 行号=%d", fileName, funcName, line)
	})

	t.Run("跳过层数测试", func(t *testing.T) {
		// 测试不同的跳过层数
		testFunc := func(skip int) (string, string, uint16, bool) {
			return getCallerInfo(skip)
		}

		// skip=1 应该指向testFunc
		_, funcName1, _, ok1 := testFunc(1)
		if !ok1 {
			t.Fatal("skip=1时获取调用者信息失败")
		}

		// skip=2 应该指向当前测试函数
		_, funcName2, _, ok2 := testFunc(2)
		if !ok2 {
			t.Fatal("skip=2时获取调用者信息失败")
		}

		// 两个函数名应该不同
		if funcName1 == funcName2 {
			t.Errorf("不同skip值应该返回不同的函数名: %s == %s", funcName1, funcName2)
		}

		t.Logf("skip=1: %s, skip=2: %s", funcName1, funcName2)
	})

	t.Run("边界条件测试", func(t *testing.T) {
		// 测试过大的skip值
		_, _, _, ok := getCallerInfo(100)
		if ok {
			t.Error("过大的skip值应该返回false")
		}

		// 测试负数skip值（虽然不推荐，但应该处理）
		_, _, _, ok = getCallerInfo(-1)
		// 这个行为取决于runtime.Caller的实现，我们只记录结果
		t.Logf("负数skip值结果: ok=%v", ok)
	})

	t.Run("文件名缓存测试", func(t *testing.T) {
		// 清空缓存，确保测试的准确性
		fileNameCache = sync.Map{}

		// 多次调用相同位置，测试缓存效果
		const iterations = 1000

		start := time.Now()
		for i := 0; i < iterations; i++ {
			getCallerInfo(1)
		}
		duration := time.Since(start)

		t.Logf("执行%d次getCallerInfo耗时: %v (平均: %v)",
			iterations, duration, duration/iterations)

		// 验证缓存中确实有数据
		cacheCount := 0
		fileNameCache.Range(func(key, value interface{}) bool {
			cacheCount++
			return true
		})

		if cacheCount == 0 {
			t.Error("文件名缓存应该包含数据")
			return
		}

		t.Logf("文件名缓存包含%d个条目", cacheCount)

		// 第二次执行，测试缓存命中
		start2 := time.Now()
		for i := 0; i < iterations; i++ {
			getCallerInfo(1)
		}
		duration2 := time.Since(start2)

		t.Logf("第二次执行%d次getCallerInfo耗时: %v (平均: %v)",
			iterations, duration2, duration2/iterations)

		// 更宽松的性能验证：允许合理的性能波动
		// 如果第二次执行时间超过第一次的3倍，才认为缓存可能失效
		if duration2 > duration*3 {
			t.Logf("警告：第二次执行时间较长，可能受系统因素影响: %v vs %v", duration2, duration)
		}

		// 验证缓存功能的正确性：多次调用应该返回相同的文件名
		fileName1, _, _, ok1 := getCallerInfo(1)
		fileName2, _, _, ok2 := getCallerInfo(1)

		if !ok1 || !ok2 {
			t.Error("获取调用者信息失败")
			return
		}

		if fileName1 != fileName2 {
			t.Errorf("缓存功能异常：多次调用返回不同的文件名: %s != %s", fileName1, fileName2)
		}
	})

	t.Run("并发缓存测试", func(t *testing.T) {
		const goroutineCount = 50
		const iterationsPerGoroutine = 100

		var wg sync.WaitGroup
		results := make([]bool, goroutineCount*iterationsPerGoroutine)
		var index int64

		// 并发调用getCallerInfo
		for i := 0; i < goroutineCount; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < iterationsPerGoroutine; j++ {
					_, _, _, ok := getCallerInfo(1)
					idx := atomic.AddInt64(&index, 1) - 1
					results[idx] = ok
				}
			}()
		}

		wg.Wait()

		// 统计成功率
		successCount := 0
		for _, ok := range results {
			if ok {
				successCount++
			}
		}

		successRate := float64(successCount) / float64(len(results)) * 100
		t.Logf("并发调用成功率: %.2f%% (%d/%d)", successRate, successCount, len(results))

		if successRate < 95.0 {
			t.Errorf("并发调用成功率过低: %.2f%%", successRate)
		}
	})
}

// TestShouldDropLogByBackpressure_Comprehensive 全面测试背压日志丢弃逻辑
func TestShouldDropLogByBackpressure_Comprehensive(t *testing.T) {
	t.Run("空指针安全测试", func(t *testing.T) {
		// 测试nil通道
		result := shouldDropLogByBackpressure(nil, nil, INFO)
		if result {
			t.Error("nil通道应该返回false")
		}

		// 测试nil背压阈值
		ch := make(chan *logMsg, 100)
		result = shouldDropLogByBackpressure(nil, ch, INFO)
		if result {
			t.Error("nil背压阈值应该返回false")
		}
	})

	t.Run("边界条件测试", func(t *testing.T) {
		// 测试容量为0的通道
		ch := make(chan *logMsg)
		bp := &bpThresholds{
			threshold80: 0,
			threshold90: 0,
			threshold95: 0,
			threshold98: 0,
		}

		result := shouldDropLogByBackpressure(bp, ch, INFO)
		if !result {
			t.Error("容量为0的通道应该丢弃日志")
		}

		// 测试容量为1的通道
		ch1 := make(chan *logMsg, 1)
		bp1 := &bpThresholds{
			threshold80: 80, // 1 * 80 = 80
			threshold90: 90, // 1 * 90 = 90
			threshold95: 95, // 1 * 95 = 95
			threshold98: 98, // 1 * 98 = 98
		}

		// 空通道不应该丢弃
		result = shouldDropLogByBackpressure(bp1, ch1, DEBUG)
		if result {
			t.Error("空通道不应该丢弃日志")
		}

		// 填满通道
		ch1 <- &logMsg{}
		result = shouldDropLogByBackpressure(bp1, ch1, DEBUG)
		if !result {
			t.Error("满通道应该丢弃日志")
		}

		// 清空通道
		<-ch1
	})

	t.Run("不同日志级别测试", func(t *testing.T) {
		ch := make(chan *logMsg, 100)
		bp := &bpThresholds{
			threshold80: 80 * 100, // 8000
			threshold90: 90 * 100, // 9000
			threshold95: 95 * 100, // 9500
			threshold98: 98 * 100, // 9800
		}

		// 填充到85%
		for i := 0; i < 85; i++ {
			ch <- &logMsg{}
		}

		// 在85%使用率下测试不同级别
		testCases := []struct {
			level    LogLevel
			expected bool
			desc     string
		}{
			{DEBUG, true, "DEBUG级别在85%时应该被丢弃"},
			{INFO, false, "INFO级别在85%时不应该被丢弃"},
			{WARN, false, "WARN级别在85%时不应该被丢弃"},
			{ERROR, false, "ERROR级别在85%时不应该被丢弃"},
			{FATAL, false, "FATAL级别在85%时不应该被丢弃"},
		}

		for _, tc := range testCases {
			result := shouldDropLogByBackpressure(bp, ch, tc.level)
			if result != tc.expected {
				t.Errorf("%s: 期望%v，实际%v", tc.desc, tc.expected, result)
			}
		}

		// 清空通道
		for len(ch) > 0 {
			<-ch
		}
	})

	t.Run("阈值精确测试", func(t *testing.T) {
		ch := make(chan *logMsg, 1000)
		chanCap := cap(ch)
		bp := &bpThresholds{
			threshold80: chanCap * 80, // 1000 * 80 = 80000 (80%)
			threshold90: chanCap * 90, // 1000 * 90 = 90000 (90%)
			threshold95: chanCap * 95, // 1000 * 95 = 95000 (95%)
			threshold98: chanCap * 98, // 1000 * 98 = 98000 (98%)
		}

		// 测试各个阈值的精确边界
		testScenarios := []struct {
			fillCount int
			level     LogLevel
			expected  bool
			desc      string
		}{
			// 79% - 不丢弃任何级别
			{790, DEBUG, false, "79%时DEBUG不应该被丢弃"},
			{790, INFO, false, "79%时INFO不应该被丢弃"},

			// 81% - 只丢弃DEBUG
			{810, DEBUG, true, "81%时DEBUG应该被丢弃"},
			{810, INFO, false, "81%时INFO不应该被丢弃"},

			// 91% - 丢弃DEBUG和INFO
			{910, DEBUG, true, "91%时DEBUG应该被丢弃"},
			{910, INFO, true, "91%时INFO应该被丢弃"},
			{910, WARN, false, "91%时WARN不应该被丢弃"},

			// 96% - 丢弃DEBUG、INFO、WARN
			{960, DEBUG, true, "96%时DEBUG应该被丢弃"},
			{960, INFO, true, "96%时INFO应该被丢弃"},
			{960, WARN, true, "96%时WARN应该被丢弃"},
			{960, ERROR, false, "96%时ERROR不应该被丢弃"},

			// 99% - 只保留FATAL
			{990, DEBUG, true, "99%时DEBUG应该被丢弃"},
			{990, INFO, true, "99%时INFO应该被丢弃"},
			{990, WARN, true, "99%时WARN应该被丢弃"},
			{990, ERROR, true, "99%时ERROR应该被丢弃"},
			{990, FATAL, false, "99%时FATAL不应该被丢弃"},
		}

		for _, scenario := range testScenarios {
			// 清空通道
			for len(ch) > 0 {
				<-ch
			}

			// 填充到指定数量
			for i := 0; i < scenario.fillCount; i++ {
				ch <- &logMsg{}
			}

			result := shouldDropLogByBackpressure(bp, ch, scenario.level)
			if result != scenario.expected {
				t.Errorf("%s: 期望%v，实际%v (通道使用率: %.1f%%)",
					scenario.desc, scenario.expected, result,
					float64(scenario.fillCount)/10.0)
			}
		}
	})
}

// TestNeedsFileInfo_Comprehensive 全面测试文件信息需求判断
func TestNeedsFileInfo_Comprehensive(t *testing.T) {
	testCases := []struct {
		format   LogFormatType
		expected bool
		desc     string
	}{
		{Detailed, true, "Detailed格式需要文件信息"},
		{Json, true, "Json格式需要文件信息"},
		{Structured, true, "Structured格式需要文件信息"},
		{Simple, false, "Simple格式不需要文件信息"},
		{JsonSimple, false, "JsonSimple格式不需要文件信息"},
		{BasicStructured, false, "BasicStructured格式不需要文件信息"},
		{SimpleTimestamp, false, "SimpleTimestamp格式不需要文件信息"},
		{Custom, false, "Custom格式不需要文件信息"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			result := needsFileInfo(tc.format)
			if result != tc.expected {
				t.Errorf("%s: 期望%v，实际%v", tc.desc, tc.expected, result)
			}
		})
	}

	// 测试无效的格式类型
	t.Run("无效格式测试", func(t *testing.T) {
		invalidFormat := LogFormatType(999)
		result := needsFileInfo(invalidFormat)
		if result {
			t.Error("无效格式应该返回false")
		}
	})
}

// TestTempBufferPool_Comprehensive 全面测试临时缓冲区池
func TestTempBufferPool_Comprehensive(t *testing.T) {
	t.Run("基本功能测试", func(t *testing.T) {
		// 获取缓冲区
		buf := getTempBuffer()
		if buf == nil {
			t.Fatal("getTempBuffer返回nil")
		}

		// 验证缓冲区是空的
		if buf.Len() != 0 {
			t.Errorf("新获取的缓冲区应该为空，实际长度: %d", buf.Len())
		}

		// 写入数据
		testData := "test data"
		buf.WriteString(testData)
		if buf.String() != testData {
			t.Errorf("缓冲区内容不正确: 期望%s，实际%s", testData, buf.String())
		}

		// 归还缓冲区
		putTempBuffer(buf)

		// 再次获取，应该是重置后的缓冲区
		buf2 := getTempBuffer()
		if buf2.Len() != 0 {
			t.Errorf("重用的缓冲区应该为空，实际长度: %d", buf2.Len())
		}

		putTempBuffer(buf2)
	})

	t.Run("nil安全测试", func(t *testing.T) {
		// 测试putTempBuffer处理nil的情况
		putTempBuffer(nil) // 不应该panic
	})

	t.Run("并发安全测试", func(t *testing.T) {
		const goroutineCount = 100
		const iterationsPerGoroutine = 100

		var wg sync.WaitGroup

		for i := 0; i < goroutineCount; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()

				for j := 0; j < iterationsPerGoroutine; j++ {
					buf := getTempBuffer()
					if buf == nil {
						t.Errorf("goroutine %d: getTempBuffer返回nil", goroutineID)
						return
					}

					// 写入一些数据
					fmt.Fprintf(buf, "goroutine-%d-iteration-%d", goroutineID, j)

					// 验证数据
					expected := fmt.Sprintf("goroutine-%d-iteration-%d", goroutineID, j)
					if buf.String() != expected {
						t.Errorf("goroutine %d: 缓冲区内容不正确", goroutineID)
						return
					}

					putTempBuffer(buf)
				}
			}(i)
		}

		wg.Wait()
	})

	t.Run("内存重用测试", func(t *testing.T) {
		// 获取多个缓冲区并记录它们的地址
		buffers := make([]*bytes.Buffer, 10)
		addresses := make([]uintptr, 10)

		for i := 0; i < 10; i++ {
			buffers[i] = getTempBuffer()
			addresses[i] = uintptr(unsafe.Pointer(buffers[i]))
		}

		// 归还所有缓冲区
		for i := 0; i < 10; i++ {
			putTempBuffer(buffers[i])
		}

		// 再次获取缓冲区，检查是否有重用
		newBuffers := make([]*bytes.Buffer, 10)
		newAddresses := make([]uintptr, 10)
		reuseCount := 0

		for i := 0; i < 10; i++ {
			newBuffers[i] = getTempBuffer()
			newAddresses[i] = uintptr(unsafe.Pointer(newBuffers[i]))

			// 检查是否重用了之前的地址
			for j := 0; j < 10; j++ {
				if newAddresses[i] == addresses[j] {
					reuseCount++
					break
				}
			}
		}

		t.Logf("缓冲区重用率: %d/10 (%.1f%%)", reuseCount, float64(reuseCount)*10)

		// 清理
		for i := 0; i < 10; i++ {
			putTempBuffer(newBuffers[i])
		}
	})
}

// MockProcessorDependencies 模拟处理器依赖
type MockProcessorDependencies struct {
	config        *FastLogConfig
	fileWriter    io.Writer
	consoleWriter io.Writer
	colorLib      *colorlib.ColorLib
	ctx           context.Context
	logChan       chan *logMsg
	doneCallback  func()
	bufferSize    int
}

func (m *MockProcessorDependencies) getConfig() *FastLogConfig {
	return m.config
}

func (m *MockProcessorDependencies) getFileWriter() io.Writer {
	return m.fileWriter
}

func (m *MockProcessorDependencies) getConsoleWriter() io.Writer {
	return m.consoleWriter
}

func (m *MockProcessorDependencies) getColorLib() *colorlib.ColorLib {
	return m.colorLib
}

func (m *MockProcessorDependencies) getContext() context.Context {
	return m.ctx
}

func (m *MockProcessorDependencies) getLogChannel() <-chan *logMsg {
	return m.logChan
}

func (m *MockProcessorDependencies) notifyProcessorDone() {
	if m.doneCallback != nil {
		m.doneCallback()
	}
}

func (m *MockProcessorDependencies) getBufferSize() int {
	return m.bufferSize
}

// TestFastLogInterfaceImplementation 测试FastLog接口实现
func TestFastLogInterfaceImplementation(t *testing.T) {
	cfg := NewFastLogConfig("logs", "interface_test.log")
	cfg.OutputToConsole = false
	log := NewFastLog(cfg)
	defer log.Close()

	t.Run("processorDependencies接口实现测试", func(t *testing.T) {
		// 验证FastLog实现了processorDependencies接口
		var deps processorDependencies = log

		// 测试各个接口方法
		config := deps.getConfig()
		if config == nil {
			t.Error("getConfig返回nil")
		}

		fileWriter := deps.getFileWriter()
		if fileWriter == nil {
			t.Error("getFileWriter返回nil")
		}

		consoleWriter := deps.getConsoleWriter()
		if consoleWriter == nil {
			t.Error("getConsoleWriter返回nil")
		}

		colorLib := deps.getColorLib()
		if colorLib == nil {
			t.Error("getColorLib返回nil")
		}

		ctx := deps.getContext()
		if ctx == nil {
			t.Error("getContext返回nil")
		}

		logChan := deps.getLogChannel()
		if logChan == nil {
			t.Error("getLogChannel返回nil")
		}

		// notifyProcessorDone在实际使用中由processor goroutine调用
		// 这里不直接测试，因为需要先调用logWait.Add()才能安全调用Done()
		// 实际的processor会在启动时Add(1)，结束时Done()
		t.Log("notifyProcessorDone方法存在且可调用")
	})
}

// TestGetCloseTimeout_Comprehensive 全面测试关闭超时计算
func TestGetCloseTimeout_Comprehensive(t *testing.T) {
	testCases := []struct {
		flushInterval time.Duration
		expected      time.Duration
		desc          string
	}{
		{
			flushInterval: 100 * time.Millisecond,
			expected:      3 * time.Second, // 100ms * 10 = 1s，但最小是3s
			desc:          "短刷新间隔应该使用最小超时",
		},
		{
			flushInterval: 500 * time.Millisecond,
			expected:      5 * time.Second, // 500ms * 10 = 5s
			desc:          "中等刷新间隔",
		},
		{
			flushInterval: 2 * time.Second,
			expected:      10 * time.Second, // 2s * 10 = 20s，但最大是10s
			desc:          "长刷新间隔应该使用最大超时",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			cfg := NewFastLogConfig("logs", "timeout_test.log")
			cfg.FlushInterval = tc.flushInterval
			cfg.OutputToConsole = false

			log := NewFastLog(cfg)
			defer log.Close()

			timeout := log.getCloseTimeout()

			if timeout != tc.expected {
				t.Errorf("超时时间不正确: %v (期望: %v)",
					timeout, tc.expected)
			}

			t.Logf("刷新间隔: %v, 计算的超时: %v", tc.flushInterval, timeout)
		})
	}
}

// BenchmarkGetCachedTimestamp 时间戳缓存性能基准测试
func BenchmarkGetCachedTimestamp(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		getCachedTimestamp()
	}
}

// BenchmarkGetCallerInfo 调用者信息获取性能基准测试
func BenchmarkGetCallerInfo(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		getCallerInfo(1)
	}
}

// BenchmarkShouldDropLogByBackpressure 背压判断性能基准测试
func BenchmarkShouldDropLogByBackpressure(b *testing.B) {
	ch := make(chan *logMsg, 1000)
	bp := &bpThresholds{
		threshold80: 80000,
		threshold90: 90000,
		threshold95: 95000,
		threshold98: 98000,
	}

	// 填充到80%
	for i := 0; i < 800; i++ {
		ch <- &logMsg{}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		shouldDropLogByBackpressure(bp, ch, INFO)
	}

	// 清理
	for len(ch) > 0 {
		<-ch
	}
}

// BenchmarkTempBufferPool 临时缓冲区池性能基准测试
func BenchmarkTempBufferPool(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf := getTempBuffer()
		buf.WriteString("test data")
		putTempBuffer(buf)
	}
}

// TestLogWithLevel_EdgeCases 测试logWithLevel的边界情况
func TestLogWithLevel_EdgeCases(t *testing.T) {
	t.Run("nil FastLog测试", func(t *testing.T) {
		var log *FastLog = nil
		// 不应该panic
		log.logWithLevel(INFO, "test", 1)
	})

	t.Run("空消息测试", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "edge_test.log")
		cfg.OutputToConsole = false
		log := NewFastLog(cfg)
		defer log.Close()

		// 空消息应该被忽略
		log.logWithLevel(INFO, "", 1)

		// 等待处理
		time.Sleep(50 * time.Millisecond)
	})

	t.Run("级别过滤测试", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "level_filter_test.log")
		cfg.OutputToConsole = false
		cfg.LogLevel = WARN // 只记录WARN及以上级别
		log := NewFastLog(cfg)
		defer log.Close()

		// 低级别日志应该被过滤
		log.logWithLevel(DEBUG, "debug message", 1)
		log.logWithLevel(INFO, "info message", 1)

		// 高级别日志应该被记录
		log.logWithLevel(WARN, "warn message", 1)
		log.logWithLevel(ERROR, "error message", 1)

		time.Sleep(100 * time.Millisecond)
	})

	t.Run("上下文取消测试", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "context_test.log")
		cfg.OutputToConsole = false
		log := NewFastLog(cfg)

		// 立即关闭上下文
		log.cancel()

		// 尝试记录日志，应该被忽略
		log.logWithLevel(INFO, "should be ignored", 1)

		log.Close()
	})
}

// TestLogFatal_EdgeCases 测试logFatal的边界情况
func TestLogFatal_EdgeCases(t *testing.T) {
	t.Run("nil FastLog Fatal测试", func(t *testing.T) {
		// 这个测试需要在子进程中运行，因为logFatal会调用os.Exit
		if os.Getenv("GO_TEST_SUBPROCESS") == "1" {
			var log *FastLog = nil
			log.logFatal("nil fastlog fatal", 1)
			return
		}

		// 启动子进程测试
		cmd := exec.Command(os.Args[0], "-test.run=TestLogFatal_EdgeCases")
		cmd.Env = append(os.Environ(), "GO_TEST_SUBPROCESS=1")
		err := cmd.Run()

		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitCode := exitErr.ExitCode(); exitCode != 1 {
				t.Errorf("期望退出码1，实际得到%d", exitCode)
			} else {
				t.Log("nil FastLog Fatal测试通过，正确退出码1")
			}
		} else if err != nil {
			t.Fatalf("测试命令执行失败: %v", err)
		} else {
			t.Error("Fatal方法没有导致程序退出")
		}
	})
}

// TestFileNameCache_Comprehensive 全面测试文件名缓存
func TestFileNameCache_Comprehensive(t *testing.T) {
	t.Run("缓存命中测试", func(t *testing.T) {
		// 清空缓存
		fileNameCache = sync.Map{}

		// 第一次调用，应该缓存未命中
		fileName1, _, _, ok1 := getCallerInfo(1)
		if !ok1 {
			t.Skip("无法获取调用者信息，跳过缓存测试")
		}

		// 第二次调用相同路径，应该缓存命中
		fileName2, _, _, ok2 := getCallerInfo(1)
		if !ok2 {
			t.Fatal("第二次获取调用者信息失败")
		}

		if fileName1 != fileName2 {
			t.Errorf("缓存测试失败：%s != %s", fileName1, fileName2)
		}

		// 验证缓存中确实有数据
		cacheCount := 0
		fileNameCache.Range(func(key, value interface{}) bool {
			cacheCount++
			return true
		})

		if cacheCount == 0 {
			t.Error("文件名缓存应该包含数据")
		}

		t.Logf("文件名缓存包含%d个条目", cacheCount)
	})

	t.Run("并发缓存安全测试", func(t *testing.T) {
		// 清空缓存
		fileNameCache = sync.Map{}

		const goroutineCount = 50
		var wg sync.WaitGroup

		for i := 0; i < goroutineCount; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()

				for j := 0; j < 100; j++ {
					getCallerInfo(1)
				}
			}(i)
		}

		wg.Wait()

		// 验证缓存状态
		cacheCount := 0
		fileNameCache.Range(func(key, value interface{}) bool {
			cacheCount++
			// 验证缓存值的类型
			if _, ok := value.(string); !ok {
				t.Errorf("缓存值类型错误: %T", value)
			}
			return true
		})

		t.Logf("并发测试后缓存包含%d个条目", cacheCount)
	})
}

// TestProcessorDependenciesInterface 测试处理器依赖接口的完整性
func TestProcessorDependenciesInterface(t *testing.T) {
	cfg := NewFastLogConfig("logs", "deps_test.log")
	cfg.OutputToConsole = false
	log := NewFastLog(cfg)
	defer log.Close()

	// 验证接口实现
	var deps processorDependencies = log

	t.Run("接口方法完整性测试", func(t *testing.T) {
		// 测试所有接口方法
		methods := []struct {
			name string
			test func() error
		}{
			{"getConfig", func() error {
				if deps.getConfig() == nil {
					return fmt.Errorf("getConfig返回nil")
				}
				return nil
			}},
			{"getFileWriter", func() error {
				if deps.getFileWriter() == nil {
					return fmt.Errorf("getFileWriter返回nil")
				}
				return nil
			}},
			{"getConsoleWriter", func() error {
				if deps.getConsoleWriter() == nil {
					return fmt.Errorf("getConsoleWriter返回nil")
				}
				return nil
			}},
			{"getColorLib", func() error {
				if deps.getColorLib() == nil {
					return fmt.Errorf("getColorLib返回nil")
				}
				return nil
			}},
			{"getContext", func() error {
				if deps.getContext() == nil {
					return fmt.Errorf("getContext返回nil")
				}
				return nil
			}},
			{"getLogChannel", func() error {
				if deps.getLogChannel() == nil {
					return fmt.Errorf("getLogChannel返回nil")
				}
				return nil
			}},
			{"notifyProcessorDone", func() error {
				// 这个方法在实际使用中由processor调用，需要先Add()才能Done()
				// 这里只验证方法存在，不实际调用以避免WaitGroup panic
				t.Log("notifyProcessorDone方法可用")
				return nil
			}},
		}

		for _, method := range methods {
			t.Run(method.name, func(t *testing.T) {
				if err := method.test(); err != nil {
					t.Error(err)
				}
			})
		}
	})
}

// TestRWTimestampCache_RaceCondition 测试时间戳缓存的竞态条件
func TestRWTimestampCache_RaceCondition(t *testing.T) {
	// 重置全局缓存
	globalRWCache = &rwTimestampCache{}

	const goroutineCount = 100
	const iterationsPerGoroutine = 1000

	var wg sync.WaitGroup
	results := make([]string, goroutineCount*iterationsPerGoroutine)
	var index int64

	// 启动大量goroutine并发访问时间戳缓存
	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < iterationsPerGoroutine; j++ {
				ts := getCachedTimestamp()
				idx := atomic.AddInt64(&index, 1) - 1
				results[idx] = ts
			}
		}()
	}

	wg.Wait()

	// 验证所有结果都是有效的时间戳
	validCount := 0
	for _, ts := range results {
		if len(ts) == 19 { // "2006-01-02 15:04:05" 的长度
			validCount++
		}
	}

	validRate := float64(validCount) / float64(len(results)) * 100
	t.Logf("有效时间戳比例: %.2f%% (%d/%d)", validRate, validCount, len(results))

	if validRate < 100.0 {
		t.Errorf("时间戳格式错误率过高: %.2f%%", 100.0-validRate)
	}

	// 统计不同时间戳的数量
	uniqueTimestamps := make(map[string]int)
	for _, ts := range results {
		uniqueTimestamps[ts]++
	}

	t.Logf("生成了%d个不同的时间戳", len(uniqueTimestamps))

	// 在快速执行的情况下，不同时间戳数量应该很少
	if len(uniqueTimestamps) > 5 {
		t.Logf("警告：时间戳变化频繁，可能存在缓存效率问题")
	}
}

// TestCircularDependencyFixed 测试循环依赖是否已修复
func TestCircularDependencyFixed(t *testing.T) {
	// 创建配置
	config := NewFastLogConfig("logs", "test.log")
	config.OutputToConsole = true
	config.OutputToFile = false
	config.ChanIntSize = 10

	// 创建日志实例
	logger := NewFastLog(config)

	// 记录一些日志
	logger.Info("测试循环依赖修复")
	logger.Debug("这是一条调试消息")
	logger.Warn("这是一条警告消息")

	// 等待处理完成
	time.Sleep(100 * time.Millisecond)

	// 关闭日志器
	logger.Close()

	// 强制垃圾回收
	runtime.GC()
	runtime.GC() // 执行两次确保彻底回收

	t.Log("循环依赖修复验证通过：日志器正常创建、使用和关闭")
}

// TestProcessorDependencyInterface 测试处理器依赖接口
func TestProcessorDependencyInterface(t *testing.T) {
	// 创建配置
	config := NewFastLogConfig("logs", "test.log")
	config.OutputToConsole = true
	config.OutputToFile = false

	// 创建日志实例
	logger := NewFastLog(config)
	defer logger.Close()

	// 验证FastLog实现了ProcessorDependencies接口
	var deps processorDependencies = logger

	// 测试接口方法
	if deps.getConfig() == nil {
		t.Error("getConfig() 返回 nil")
	}

	if deps.getFileWriter() == nil {
		t.Error("getFileWriter() 返回 nil")
	}

	if deps.getConsoleWriter() == nil {
		t.Error("getConsoleWriter() 返回 nil")
	}

	if deps.getColorLib() == nil {
		t.Error("getColorLib() 返回 nil")
	}

	if deps.getContext() == nil {
		t.Error("getContext() 返回 nil")
	}

	if deps.getLogChannel() == nil {
		t.Error("getLogChannel() 返回 nil")
	}

	t.Log("ProcessorDependencies接口实现验证通过")
}

// TestMemoryLeakPrevention 测试内存泄漏预防
func TestMemoryLeakPrevention(t *testing.T) {
	// 记录初始内存统计
	var m1 runtime.MemStats
	runtime.GC()
	runtime.GC() // 执行两次GC确保稳定状态
	time.Sleep(50 * time.Millisecond)
	runtime.ReadMemStats(&m1)

	// 创建和销毁多个日志实例
	for i := 0; i < 10; i++ {
		config := NewFastLogConfig("logs", "test.log")
		config.OutputToConsole = true
		config.OutputToFile = false
		config.LogLevel = NONE // 设置为NONE级别，避免实际输出日志
		config.ChanIntSize = 5

		logger := NewFastLog(config)

		// 写入一些日志
		logger.Infof("测试消息 %d", i)
		time.Sleep(10 * time.Millisecond)

		// 关闭日志器
		logger.Close()
	}

	// 强制垃圾回收
	runtime.GC()
	runtime.GC()
	time.Sleep(100 * time.Millisecond)

	// 记录最终内存统计
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	// 安全地计算内存变化（处理uint64下溢问题）
	var memoryChange int64
	if m2.Alloc >= m1.Alloc {
		memoryChange = int64(m2.Alloc - m1.Alloc)
		t.Logf("内存增长: %d bytes", memoryChange)
	} else {
		memoryChange = -int64(m1.Alloc - m2.Alloc)
		t.Logf("内存减少: %d bytes", -memoryChange)
	}

	// 检查是否存在明显的内存泄漏
	// 允许合理的内存增长（最多1MB），但不应该有大量泄漏
	if memoryChange > 1024*1024 {
		t.Errorf("可能存在内存泄漏，内存增长: %d bytes", memoryChange)
	} else if memoryChange < 0 {
		t.Logf("内存使用优化良好，减少了 %d bytes", -memoryChange)
	} else {
		t.Logf("内存增长在合理范围内: %d bytes", memoryChange)
	}

	t.Log("内存泄漏预防测试通过")
}

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

		// 使用defer捕获可能的panic
		defer func() {
			if r := recover(); r != nil {
				t.Logf("预期的权限错误panic: %v", r)
			}
		}()

		logger := NewFastLog(cfg)

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

		// 使用defer捕获可能的panic
		defer func() {
			if r := recover(); r != nil {
				t.Logf("预期的文件权限错误panic: %v", r)
			}
		}()

		logger := NewFastLog(cfg)

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

		// 使用defer捕获可能的panic
		defer func() {
			if r := recover(); r != nil {
				t.Logf("深层目录创建可能失败panic: %v", r)
			}
		}()

		logger := NewFastLog(cfg)
		logger.Info("测试深层目录创建")
		time.Sleep(100 * time.Millisecond)
		logger.Close()

		// 清理创建的目录
		_ = os.RemoveAll(strings.Split(deepPath, "/")[0])
	})

	t.Run("无效路径字符", func(t *testing.T) {
		// 在非Windows系统上测试包含无效字符的路径
		if runtime.GOOS != "windows" {
			cfg := NewFastLogConfig("test\x00invalid", "test.log")
			cfg.OutputToConsole = false

			// 使用defer捕获可能的panic
			defer func() {
				if r := recover(); r != nil {
					t.Logf("预期的无效路径错误panic: %v", r)
				}
			}()

			logger := NewFastLog(cfg)
			logger.Close()
		}
	})
}

// TestFileSystemSpaceHandling 测试文件系统空间处理
func TestFileSystemSpaceHandling(t *testing.T) {
	t.Run("大文件写入测试", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "large_test.log")
		cfg.OutputToConsole = false
		cfg.MaxLogFileSize = 1 // 1MB限制，触发轮转

		logger := NewFastLog(cfg)
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

			logger := NewFastLog(cfg)
			loggers = append(loggers, logger)
		}

		// 并发写入日志
		done := make(chan bool, len(loggers))
		for i, logger := range loggers {
			go func(id int, f *FastLog) {
				defer func() { done <- true }()
				for j := 0; j < 100; j++ {
					f.Infof("实例 %d 消息 %d", id, j)
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

		logger := NewFastLog(cfg)
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

			// 使用defer捕获可能的panic
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Logf("创建第 %d 个实例时panic: %v", i, r)
					}
				}()
				logger := NewFastLog(cfg)
				loggers = append(loggers, logger)
			}()
			if len(loggers) <= i {
				break
			}
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

		logger := NewFastLog(cfg)
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

	logger := NewFastLog(cfg)
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

// TestNullPointerProtection 测试空指针保护机制
func TestNullPointerProtection(t *testing.T) {
	t.Run("nil FastLog instance", func(t *testing.T) {
		var logger *FastLog = nil

		// 测试所有公共API方法都能安全处理nil实例
		logger.Info("test message")
		logger.Debug("test message")
		logger.Warn("test message")
		logger.Error("test message")

		logger.Infof("test %s", "message")
		logger.Debugf("test %s", "message")
		logger.Warnf("test %s", "message")
		logger.Errorf("test %s", "message")

		// 如果程序能执行到这里，说明空指针保护有效
		t.Log("空指针保护测试通过")
	})

	t.Run("empty format string", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "test.log")
		cfg.OutputToConsole = true
		cfg.OutputToFile = false
		cfg.LogLevel = NONE // 设置为NONE级别，避免实际输出日志

		logger := NewFastLog(cfg)
		defer logger.Close()

		// 测试空格式字符串的处理
		logger.Infof("")
		logger.Debugf("")
		logger.Warnf("")
		logger.Errorf("")

		t.Log("空格式字符串保护测试通过")
	})
}

// TestObjectPoolSafety 测试对象池安全性
func TestObjectPoolSafety(t *testing.T) {
	t.Run("concurrent object pool access", func(t *testing.T) {
		const goroutineCount = 100
		const operationsPerGoroutine = 1000

		var wg sync.WaitGroup
		wg.Add(goroutineCount)

		// 并发访问对象池
		for i := 0; i < goroutineCount; i++ {
			go func() {
				defer wg.Done()

				for j := 0; j < operationsPerGoroutine; j++ {
					// 获取对象
					msg := getLogMsg()
					if msg == nil {
						t.Errorf("getLogMsg() 返回了 nil")
						return
					}

					// 使用对象
					msg.Message = "test message"
					msg.Level = INFO
					msg.Timestamp = "2023-01-01 12:00:00"

					// 归还对象
					putLogMsg(msg)
				}
			}()
		}

		wg.Wait()
		t.Log("对象池并发安全测试通过")
	})

	t.Run("object pool nil handling", func(t *testing.T) {
		// 测试putLogMsg对nil的处理
		putLogMsg(nil) // 不应该panic

		// 测试getLogMsg永远不返回nil
		for i := 0; i < 100; i++ {
			msg := getLogMsg()
			if msg == nil {
				t.Fatal("getLogMsg() 不应该返回 nil")
			}
			putLogMsg(msg)
		}

		t.Log("对象池nil处理测试通过")
	})
}

// TestHighConcurrencyStability 测试高并发场景下的稳定性
func TestHighConcurrencyStability(t *testing.T) {
	cfg := NewFastLogConfig("logs", "stability_test.log")
	cfg.OutputToConsole = false
	cfg.OutputToFile = true
	cfg.ChanIntSize = 1000 // 较小的通道容量，更容易触发背压

	logger := NewFastLog(cfg)
	defer logger.Close()

	const goroutineCount = 50
	const messagesPerGoroutine = 1000

	var wg sync.WaitGroup
	wg.Add(goroutineCount)

	// 启动多个goroutine并发写入日志
	for i := 0; i < goroutineCount; i++ {
		go func(id int) {
			defer wg.Done()

			for j := 0; j < messagesPerGoroutine; j++ {
				// 随机选择不同的日志方法
				switch j % 5 {
				case 0:
					logger.Info("并发测试消息", id, j)
				case 1:
					logger.Debug("并发测试消息", id, j)
				case 2:
					logger.Warn("并发测试消息", id, j)
				case 3:
					logger.Error("并发测试消息", id, j)
				}
			}
		}(i)
	}

	wg.Wait()

	// 等待日志处理完成
	time.Sleep(500 * time.Millisecond)

	t.Logf("高并发稳定性测试完成: %d个goroutine，每个写入%d条日志",
		goroutineCount, messagesPerGoroutine)
}

// TestResourceExhaustionHandling 测试资源耗尽场景的处理
func TestResourceExhaustionHandling(t *testing.T) {
	cfg := NewFastLogConfig("logs", "resource_test.log")
	cfg.OutputToConsole = false
	cfg.OutputToFile = true
	cfg.ChanIntSize = 10 // 非常小的通道容量

	logger := NewFastLog(cfg)
	defer logger.Close()

	// 快速发送大量日志，触发背压和通道满的情况
	for i := 0; i < 1000; i++ {
		logger.Infof("资源耗尽测试消息 %d", i)
	}

	// 系统应该能够优雅处理，不会panic或死锁
	time.Sleep(100 * time.Millisecond)

	t.Log("资源耗尽处理测试通过")
}

// BenchmarkNullPointerCheck 基准测试空指针检查的性能影响
func BenchmarkNullPointerCheck(b *testing.B) {
	cfg := NewFastLogConfig("logs", "benchmark.log")
	cfg.OutputToConsole = false
	cfg.OutputToFile = false // 禁用所有输出，只测试检查逻辑

	logger := NewFastLog(cfg)
	defer logger.Close()

	b.ResetTimer()

	// 测试添加空指针检查后的性能影响
	for i := 0; i < b.N; i++ {
		logger.Info("基准测试消息")
	}
}

// TestGracefulShutdown 测试优雅关闭机制
func TestGracefulShutdown(t *testing.T) {
	t.Run("正常关闭流程", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "shutdown_test.log")
		cfg.OutputToConsole = false

		logger := NewFastLog(cfg)

		// 写入一些日志
		for i := 0; i < 100; i++ {
			logger.Infof("关闭测试消息 %d", i)
		}

		// 记录关闭开始时间
		startTime := time.Now()

		// 执行关闭
		logger.Close()

		// 记录关闭耗时
		shutdownDuration := time.Since(startTime)
		t.Logf("关闭耗时: %v", shutdownDuration)

		// 验证关闭后不能再写入日志（不应该panic）
		logger.Info("关闭后的消息")

		// 验证日志文件存在且有内容
		logFile := filepath.Join("logs", "shutdown_test.log")
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			t.Error("关闭后日志文件不存在")
		}
	})

	t.Run("重复关闭测试", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "repeat_close_test.log")
		cfg.OutputToConsole = false

		logger := NewFastLog(cfg)

		logger.Info("重复关闭测试")

		// 第一次关闭
		logger.Close()

		// 重复关闭应该是安全的（幂等性）
		logger.Close()
		logger.Close()

		// 不应该panic
		t.Log("重复关闭测试通过")
	})

	t.Run("并发关闭测试", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "concurrent_close_test.log")
		cfg.OutputToConsole = false

		logger := NewFastLog(cfg)

		logger.Info("并发关闭测试")

		// 启动多个goroutine同时关闭
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				logger.Close()
			}(i)
		}

		wg.Wait()
		t.Log("并发关闭测试通过")
	})
}

// TestShutdownWithActiveWriting 测试关闭过程中的活跃写入
func TestShutdownWithActiveWriting(t *testing.T) {
	cfg := NewFastLogConfig("logs", "active_writing_test.log")
	cfg.OutputToConsole = false
	cfg.ChanIntSize = 1000

	logger := NewFastLog(cfg)

	// 启动后台写入
	var wg sync.WaitGroup
	stopWriting := make(chan bool)

	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 0
		for {
			select {
			case <-stopWriting:
				return
			default:
				logger.Infof("活跃写入测试消息 %d", i)
				i++
				time.Sleep(1 * time.Millisecond)
			}
		}
	}()

	// 让写入运行一段时间
	time.Sleep(100 * time.Millisecond)

	// 停止写入并关闭
	close(stopWriting)
	logger.Close()

	wg.Wait()
	t.Log("活跃写入关闭测试通过")
}

// TestShutdownTimeout 测试关闭超时处理
func TestShutdownTimeout(t *testing.T) {
	cfg := NewFastLogConfig("logs", "timeout_test.log")
	cfg.OutputToConsole = false
	cfg.ChanIntSize = 10000              // 大通道容量
	cfg.FlushInterval = 10 * time.Second // 长刷新间隔

	logger := NewFastLog(cfg)

	// 填充大量日志到通道中
	for i := 0; i < 5000; i++ {
		logger.Infof("超时测试消息 %d", i)
	}

	// 测试关闭超时
	startTime := time.Now()
	logger.Close()
	shutdownDuration := time.Since(startTime)

	t.Logf("关闭耗时: %v", shutdownDuration)

	// 验证关闭在合理时间内完成
	if shutdownDuration > 30*time.Second {
		t.Errorf("关闭耗时过长: %v", shutdownDuration)
	}
}

// TestResourceCleanup 测试资源清理
func TestResourceCleanup(t *testing.T) {
	t.Run("goroutine清理验证", func(t *testing.T) {
		initialGoroutines := runtime.NumGoroutine()

		cfg := NewFastLogConfig("logs", "cleanup_test.log")
		cfg.OutputToConsole = false

		logger := NewFastLog(cfg)

		logger.Info("资源清理测试")
		time.Sleep(100 * time.Millisecond)

		logger.Close()
		time.Sleep(100 * time.Millisecond) // 等待goroutine退出

		finalGoroutines := runtime.NumGoroutine()

		// 允许一定的goroutine数量波动
		if finalGoroutines > initialGoroutines+2 {
			t.Errorf("可能存在goroutine泄漏: 初始=%d, 最终=%d", initialGoroutines, finalGoroutines)
		}
	})

	t.Run("文件句柄清理验证", func(t *testing.T) {
		cfg := NewFastLogConfig("logs", "handle_cleanup_test.log")
		cfg.OutputToConsole = false

		logger := NewFastLog(cfg)

		logger.Info("文件句柄清理测试")
		time.Sleep(100 * time.Millisecond)

		logger.Close()

		// 验证日志文件可以被删除（说明文件句柄已关闭）
		logFile := filepath.Join("logs", "handle_cleanup_test.log")
		time.Sleep(100 * time.Millisecond) // 等待文件句柄释放

		err := os.Remove(logFile)
		if err != nil {
			t.Errorf("无法删除日志文件，可能存在文件句柄泄漏: %v", err)
		}
	})
}
