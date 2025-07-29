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
	cfg := NewFastLogConfig(t.TempDir(), "concurrent.log")
	cfg.SetConsoleOnly(false)
	log, _ := NewFastLog(cfg)
	defer func() { _ = log.Close() }()

	// 并发写入1000条日志
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			log.Debugf("concurrent log message %d", id)
		}(i)
	}
	wg.Wait()

	// 强制刷新剩余内容
	log.flushBufferNow()

	// 验证日志总数
	content, _ := os.ReadFile(cfg.logDirName + "/concurrent.log")
	lines := bytes.Count(content, []byte{'\n'})
	if lines != 1000 {
		t.Errorf("并发处理应生成1000条日志，实际: %d条", lines)
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
