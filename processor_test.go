package fastlog

import (
	"bytes"
	"os"
	"sync"
	"testing"
	"time"
)

// TestHandleLog_BufferThreshold 测试缓冲区达到阈值时的自动刷新功能
func TestHandleLog_BufferThreshold(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	cfg := NewFastLogConfig(tempDir, "threshold.log")
	cfg.SetConsoleOnly(false)
	cfg.SetFlushInterval(1 * time.Hour) // 禁用自动刷新

	log, _ := NewFastLog(cfg)
	defer func() { _ = log.Close() }()

	// 计算触发阈值所需的日志大小 (230400字节)
	thresholdSize := flushThreshold
	msgSize := 1024 // 每条日志约1KB
	messagesNeeded := thresholdSize/msgSize + 1

	// 创建测试消息
	largeMsg := generateLargeString(msgSize)

	// 并发写入消息直到触发阈值
	var wg sync.WaitGroup
	for i := 0; i < messagesNeeded; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Info(largeMsg)
		}()
	}
	wg.Wait()

	// 记录初始索引
	log.fileBufferMu.Lock()
	initialIdx := log.fileBufferIdx.Load()
	log.fileBufferMu.Unlock()

	// 等待短暂时间确保刷新完成
	time.Sleep(100 * time.Millisecond)

	// 验证缓冲区已刷新
	log.fileBufferMu.Lock()
	currentIdx := log.fileBufferIdx.Load()
	bufSize := log.fileBuffers[initialIdx].Len()
	log.fileBufferMu.Unlock()

	if currentIdx == initialIdx {
		t.Error("达到阈值后缓冲区索引应切换")
	}
	if bufSize >= thresholdSize {
		t.Errorf("刷新后缓冲区大小应小于阈值，实际大小: %d bytes", bufSize)
	}
}

// TestProcessLogs_ChannelHandling 测试日志通道处理功能
func TestProcessLogs_ChannelHandling(t *testing.T) {
	cfg := NewFastLogConfig("", "")
	cfg.SetConsoleOnly(true)
	log, _ := NewFastLog(cfg)
	defer func() { _ = log.Close() }()

	// 发送测试消息
	log.Infof("test channel message")

	// 等待消息处理
	time.Sleep(100 * time.Millisecond)

	// 验证消息已处理
	log.consoleBufferMu.Lock()
	bufContent := log.consoleBuffers[log.consoleBufferIdx.Load()].String()
	log.consoleBufferMu.Unlock()

	if bufContent == "" {
		t.Error("日志通道处理失败，缓冲区内容为空")
	}
}

// TestHandleLog_ConcurrentProcessing 测试并发日志处理
func TestHandleLog_ConcurrentProcessing(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	cfg := NewFastLogConfig(tempDir, "concurrent.log")
	cfg.SetConsoleOnly(false)
	cfg.SetLogLevel(DEBUG)       // 确保DEBUG级别日志能输出
	cfg.SetPrintToConsole(false) // 只写入文件

	log, err := NewFastLog(cfg)
	if err != nil {
		t.Fatalf("创建日志记录器失败: %v", err)
	}

	// 并发写入1000条日志
	var wg sync.WaitGroup
	logCount := 1000
	for i := 0; i < logCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			log.Infof("concurrent log message %d", id) // 使用INFO级别确保输出
		}(i)
	}
	wg.Wait()

	// 等待所有日志写入完成
	time.Sleep(200 * time.Millisecond)

	_ = log.Close()

	// 再次等待确保文件写入完成
	time.Sleep(100 * time.Millisecond)

	// 验证日志总数
	content, err := os.ReadFile(log.logFilePath)
	if err != nil {
		t.Fatalf("读取日志文件失败: %v", err)
	}

	// 计算包含"concurrent log message"的行数
	lines := bytes.Split(content, []byte{'\n'})
	validLines := 0
	for _, line := range lines {
		if bytes.Contains(line, []byte("concurrent log message")) {
			validLines++
		}
	}

	if validLines != logCount {
		t.Errorf("并发处理应生成%d条日志, 实际: %d条", logCount, validLines)
		t.Logf("文件内容长度: %d bytes", len(content))
	}
}

// generateLargeString 创建指定大小的测试字符串
func generateLargeString(size int) string {
	buf := make([]byte, size)
	for i := range buf {
		buf[i] = 'a'
	}
	return string(buf)
}
