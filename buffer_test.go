package fastlog

import (
	"bytes"
	"os"
	"sync"
	"testing"
	"time"
)

// TestFlushBufferNow_FileBuffer 测试文件缓冲区立即刷新功能
func TestFlushBufferNow_FileBuffer(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	cfg := NewFastLogConfig(tempDir, "test.log")
	cfg.SetConsoleOnly(false)
	cfg.SetFlushInterval(1 * time.Hour) // 禁用自动刷新

	log, _ := NewFastLog(cfg)
	defer func() { _ = log.Close() }()

	// 手动写入测试数据到缓冲区
	log.fileBufferMu.Lock()
	log.fileBuffers[0].WriteString("test file buffer content\n")
	log.fileBufferMu.Unlock()

	// 执行立即刷新
	log.flushBufferNow()

	// 验证缓冲区已切换且内容已写入
	log.fileBufferMu.Lock()
	currentIdx := log.fileBufferIdx.Load()
	bufContent := log.fileBuffers[0].String()
	log.fileBufferMu.Unlock()

	if currentIdx != 1 {
		t.Error("刷新后缓冲区索引应切换为1")
	}
	if bufContent != "" {
		t.Errorf("刷新后原缓冲区应清空，实际内容: %s", bufContent)
	}

	// 验证文件内容
	filePath := tempDir + "/test.log"
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("读取测试日志文件失败: %v", err)
	}
	if string(content) != "test file buffer content\n" {
		t.Errorf("文件内容不匹配，实际内容: %s", string(content))
	}
}

// TestFlushBuffer_TimedFlush 测试定时刷新功能
func TestFlushBuffer_TimedFlush(t *testing.T) {
	// 创建临时目录和短间隔配置
	tempDir := t.TempDir()
	cfg := NewFastLogConfig(tempDir, "timed.log")
	cfg.SetConsoleOnly(false)
	cfg.SetFlushInterval(10 * time.Millisecond) // 极短刷新间隔

	log, _ := NewFastLog(cfg)
	defer func() { _ = log.Close() }()

	// 写入测试数据
	log.fileBufferMu.Lock()
	log.fileBuffers[0].WriteString("timed flush test\n")
	log.fileBufferMu.Unlock()

	// 等待定时刷新触发
	time.Sleep(50 * time.Millisecond)

	// 验证文件内容
	content, err := os.ReadFile(tempDir + "/timed.log")
	if err != nil {
		t.Fatalf("读取定时刷新测试文件失败: %v", err)
	}
	if string(content) != "timed flush test\n" {
		t.Errorf("定时刷新内容不匹配，实际内容: %s", string(content))
	}
}

// TestFlushBuffer_ConcurrentSafety 测试并发环境下的缓冲区安全
func TestFlushBuffer_ConcurrentSafety(t *testing.T) {
	tempDir := t.TempDir()
	cfg := NewFastLogConfig(tempDir, "concurrent.log")
	cfg.SetConsoleOnly(false)
	cfg.SetLogLevel(DEBUG)                       // 确保DEBUG级别日志能输出
	cfg.SetFlushInterval(100 * time.Millisecond) // 使用合理的刷新间隔
	cfg.SetPrintToConsole(false)                 // 只写入文件，避免控制台输出干扰

	log, err := NewFastLog(cfg)
	if err != nil {
		t.Fatalf("创建日志记录器失败: %v", err)
	}
	defer func() { _ = log.Close() }()

	// 并发写入测试
	var wg sync.WaitGroup
	logCount := 100
	for i := 0; i < logCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			log.Infof("concurrent test %d", id) // 使用INFO级别确保输出
		}(i)
	}
	wg.Wait()

	// 等待所有日志写入完成
	time.Sleep(200 * time.Millisecond)

	// 手动刷新确保所有内容写入文件
	log.flushBufferNow()

	// 再等待一下确保文件写入完成
	time.Sleep(50 * time.Millisecond)

	// 验证文件内容
	logFilePath := tempDir + "/concurrent.log"
	content, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("读取日志文件失败: %v", err)
	}

	// 计算包含"concurrent test"的行数
	lines := bytes.Split(content, []byte{'\n'})
	validLines := 0
	for _, line := range lines {
		if bytes.Contains(line, []byte("concurrent test")) {
			validLines++
		}
	}

	if validLines != logCount {
		t.Errorf("并发写入应产生%d行日志，实际: %d行", logCount, validLines)
		t.Logf("文件内容长度: %d bytes", len(content))
	}
}
