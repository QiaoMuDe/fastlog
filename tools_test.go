package fastlog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestCheckPath 测试路径检查功能
func TestCheckPath(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")

	// 创建测试文件
	_ = os.WriteFile(filePath, []byte("test"), 0644)

	// 测试存在的文件
	info, err := checkPath(filePath)
	if err != nil {
		t.Errorf("检查存在的文件失败: %v", err)
	}
	if !info.Exists || !info.IsFile || info.IsDir {
		t.Error("存在的文件属性判断错误")
	}

	// 测试存在的目录
	info, err = checkPath(tempDir)
	if err != nil {
		t.Errorf("检查存在的目录失败: %v", err)
	}
	if !info.Exists || !info.IsDir || info.IsFile {
		t.Error("存在的目录属性判断错误")
	}

	// 测试不存在的路径
	info, err = checkPath(filepath.Join(tempDir, "nonexistent"))
	if err == nil || info.Exists {
		t.Error("检查不存在的路径时应返回错误")
	}
}

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

// TestGetGoroutineID 测试协程ID获取
func TestGetGoroutineID(t *testing.T) {
	// 主协程ID
	mainID := getGoroutineID()
	if mainID <= 0 {
		t.Error("主协程ID应为正数")
	}

	// 测试不同协程ID
	var wg sync.WaitGroup
	ids := make([]int64, 5)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			ids[idx] = getGoroutineID()
			if ids[idx] <= 0 {
				t.Errorf("协程ID应为正数，实际为%d", ids[idx])
			}
		}(i)
	}
	wg.Wait()

	// 检查所有ID是否唯一
	idSet := make(map[int64]bool)
	for _, id := range ids {
		if idSet[id] {
			t.Errorf("发现重复的协程ID: %d", id)
		}
		idSet[id] = true
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
		level: SUCCESS, expected: "SUCCESS",
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

// TestAddColor 测试日志颜色添加
func TestAddColor(t *testing.T) {
	cfg := NewFastLogConfig("logs", "app.log")
	log, _ := NewFastLog(cfg)

	msg := &logMessage{
		level:   INFO,
		message: "test color",
	}

	colored := addColor(log, msg, "test color")
	if !strings.Contains(colored, "test color") {
		t.Error("彩色日志应包含原始消息")
	}

	// 测试禁用颜色
	cfg.NoColor = true
	log2, _ := NewFastLog(cfg)
	uncolored := addColor(log2, msg, "test color")
	// 禁用颜色后，返回的应该就是原始消息，不包含颜色代码
	if uncolored != "test color" {
		t.Error("禁用颜色后应返回原始消息")
	}
}

// TestFormatLog 测试日志格式化
func TestFormatLog(t *testing.T) {
	msg := &logMessage{
		timestamp:   time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		level:       INFO,
		message:     "test format",
		fileName:    "test.go",
		funcName:    "TestFunc",
		line:        42,
		goroutineID: 123,
	}

	// 测试详细格式 - 创建第一个实例
	cfg1 := NewFastLogConfig("", "")
	cfg1.LogFormat = Detailed
	cfg1.OutputToConsole = false // 禁用控制台输出，避免启动处理器
	cfg1.OutputToFile = false    // 禁用文件输出，避免启动处理器
	log1, _ := NewFastLog(cfg1)
	defer log1.Close()

	detailed := formatLog(log1, msg)
	if !strings.Contains(detailed, "2023-01-01 12:00:00") || !strings.Contains(detailed, "INFO") {
		t.Error("详细格式日志不完整")
	}

	// 测试JSON格式 - 创建第二个实例
	cfg2 := NewFastLogConfig("", "")
	cfg2.LogFormat = Json
	cfg2.OutputToConsole = false // 禁用控制台输出，避免启动处理器
	cfg2.OutputToFile = false    // 禁用文件输出，避免启动处理器
	log2, _ := NewFastLog(cfg2)
	defer log2.Close()

	jsonLog := formatLog(log2, msg)
	if !strings.Contains(jsonLog, "\"level\":\"INFO\"") || !strings.Contains(jsonLog, "\"message\":\"test format\"") {
		t.Error("JSON格式日志不完整")
	}

	fmt.Println(jsonLog)
}

// TestBackpressure 测试背压功能
func TestBackpressure(t *testing.T) {
	// 创建一个小容量的通道来模拟背压情况
	testChan := make(chan *logMessage, 10) // 容量为10的通道

	// 测试用例数据
	testCases := []struct {
		name         string
		channelFill  int      // 通道填充数量
		level        LogLevel // 测试的日志级别
		expectedDrop bool     // 是否应该被丢弃
		description  string   // 测试描述
	}{
		// 正常情况（0-60%）
		{"正常_DEBUG", 0, DEBUG, false, "通道空闲时，DEBUG日志应该保留"},
		{"正常_INFO", 5, INFO, false, "通道50%时，INFO日志应该保留"},
		{"正常_WARN", 6, WARN, false, "通道60%时，WARN日志应该保留"},

		// 70%背压（丢弃DEBUG）
		{"70%背压_DEBUG", 7, DEBUG, true, "通道70%时，DEBUG日志应该被丢弃"},
		{"70%背压_INFO", 7, INFO, false, "通道70%时，INFO日志应该保留"},
		{"70%背压_WARN", 7, WARN, false, "通道70%时，WARN日志应该保留"},

		// 80%背压（只保留WARN及以上）
		{"80%背压_DEBUG", 8, DEBUG, true, "通道80%时，DEBUG日志应该被丢弃"},
		{"80%背压_INFO", 8, INFO, true, "通道80%时，INFO日志应该被丢弃"},
		{"80%背压_WARN", 8, WARN, false, "通道80%时，WARN日志应该保留"},
		{"80%背压_ERROR", 8, ERROR, false, "通道80%时，ERROR日志应该保留"},

		// 90%背压（只保留ERROR和FATAL）
		{"90%背压_DEBUG", 9, DEBUG, true, "通道90%时，DEBUG日志应该被丢弃"},
		{"90%背压_INFO", 9, INFO, true, "通道90%时，INFO日志应该被丢弃"},
		{"90%背压_WARN", 9, WARN, true, "通道90%时，WARN日志应该被丢弃"},
		{"90%背压_ERROR", 9, ERROR, false, "通道90%时，ERROR日志应该保留"},
		{"90%背压_FATAL", 9, FATAL, false, "通道90%时，FATAL日志应该保留"},

		// 边界情况
		{"满通道_ERROR", 10, ERROR, false, "通道满时，ERROR日志应该保留"},
		{"满通道_DEBUG", 10, DEBUG, true, "通道满时，DEBUG日志应该被丢弃"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 清空通道
			for len(testChan) > 0 {
				<-testChan
			}

			// 填充通道到指定数量
			for i := 0; i < tc.channelFill; i++ {
				testChan <- &logMessage{level: INFO, message: "填充消息"}
			}

			// 测试背压函数
			shouldDrop := shouldDropLogByBackpressure(testChan, tc.level)

			// 验证结果
			if shouldDrop != tc.expectedDrop {
				t.Errorf("测试失败: %s\n"+
					"通道使用率: %d/%d (%.0f%%)\n"+
					"日志级别: %s\n"+
					"预期丢弃: %v, 实际丢弃: %v\n"+
					"描述: %s",
					tc.name,
					len(testChan), cap(testChan), float64(len(testChan))*100/float64(cap(testChan)),
					logLevelToString(tc.level),
					tc.expectedDrop, shouldDrop,
					tc.description)
			} else {
				t.Logf("✓ %s - 通道使用率: %d%%, 级别: %s, 丢弃: %v",
					tc.name,
					len(testChan)*100/cap(testChan),
					logLevelToString(tc.level),
					shouldDrop)
			}
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

	log, err := NewFastLog(cfg)
	if err != nil {
		t.Fatalf("创建日志失败: %v", err)
	}
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
	testChan := make(chan *logMessage, 1000)

	// 填充通道到80%
	for i := 0; i < 800; i++ {
		testChan <- &logMessage{level: INFO}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		shouldDropLogByBackpressure(testChan, WARN)
	}
}
