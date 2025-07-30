package fastlog

import (
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
	cfg.SetNoColor(true)
	log2, _ := NewFastLog(cfg)
	uncolored := addColor(log2, msg, "test color")
	// 禁用颜色后，返回的应该就是原始消息，不包含颜色代码
	if uncolored != "test color" {
		t.Error("禁用颜色后应返回原始消息")
	}
}

// TestFormatLog 测试日志格式化
func TestFormatLog(t *testing.T) {
	cfg := NewFastLogConfig("", "")
	log, _ := NewFastLog(cfg)

	msg := &logMessage{
		timestamp:   time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		level:       INFO,
		message:     "test format",
		fileName:    "test.go",
		funcName:    "TestFunc",
		line:        42,
		goroutineID: 123,
	}

	// 测试详细格式
	cfg.SetLogFormat(Detailed)
	detailed := formatLog(log, msg)
	if !strings.Contains(detailed, "2023-01-01 12:00:00") || !strings.Contains(detailed, "INFO") {
		t.Error("详细格式日志不完整")
	}

	// 测试JSON格式
	cfg.SetLogFormat(Json)
	jsonLog := formatLog(log, msg)
	if !strings.Contains(jsonLog, "\"level\":\"INFO\"") || !strings.Contains(jsonLog, "\"message\":\"test format\"") {
		t.Error("JSON格式日志不完整")
	}
}
