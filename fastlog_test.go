package fastlog

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestLogger(t *testing.T) {
	// 创建临时目录用于存储日志文件
	tmpDir, err := ioutil.TempDir("", "logtest")
	if err != nil {
		t.Fatalf("无法创建临时目录: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建日志配置
	cfg := NewConfig(tmpDir, "test.log")
	cfg.LogLevel = Debug

	// 创建日志记录器
	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("创建日志记录器失败: %v", err)
	}

	// 测试 Info 日志
	logger.Info("这是一条 Info 日志")

	// 测试 Warn 日志
	logger.Warn("这是一条 Warn 日志")

	// 测试 Error 日志
	logger.Error("这是一条 Error 日志")

	// 测试 Success 日志
	logger.Success("这是一条 Success 日志")

	// 测试 Debug 日志
	logger.Debug("这是一条 Debug 日志")

	// 测试支持格式化的 Info 日志
	logger.Infof("这是一条格式化的 Info 日志: %s", "参数")

	// 测试支持格式化的 Warn 日志
	logger.Warnf("这是一条格式化的 Warn 日志: %s", "参数")

	// 测试支持格式化的 Error 日志
	logger.Errorf("这是一条格式化的 Error 日志: %s", "参数")

	// 测试支持格式化的 Success 日志
	logger.Successf("这是一条格式化的 Success 日志: %s", "参数")

	// 测试支持格式化的 Debug 日志
	logger.Debugf("这是一条格式化的 Debug 日志: %s", "参数")

	// 关闭日志记录器
	logger.Close()

	// 检查日志文件是否存在
	logPath := filepath.Join(tmpDir, "test.log")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Errorf("日志文件不存在")
	}

	// 读取日志文件内容
	logContent, err := ioutil.ReadFile(logPath)
	if err != nil {
		t.Fatalf("读取日志文件失败: %v", err)
	}

	// 检查日志内容是否包含预期的日志信息
	expectedLogs := []string{
		"这是一条 Info 日志",
		"这是一条 Warn 日志",
		"这是一条 Error 日志",
		"这是一条 Success 日志",
		"这是一条 Debug 日志",
		"这是一条格式化的 Info 日志: 参数",
		"这是一条格式化的 Warn 日志: 参数",
		"这是一条格式化的 Error 日志: 参数",
		"这是一条格式化的 Success 日志: 参数",
		"这是一条格式化的 Debug 日志: 参数",
	}

	for _, expectedLog := range expectedLogs {
		if !bytes.Contains(logContent, []byte(expectedLog)) {
			t.Errorf("日志文件中未找到预期的日志信息: %s", expectedLog)
		}
	}
}
