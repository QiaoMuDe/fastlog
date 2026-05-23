package fastlog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestLevelRouterBasic 测试级别路由基础功能
func TestLevelRouterBasic(t *testing.T) {
	tmpDir := t.TempDir()

	// 创建配置，启用级别路由
	cfg := NewConfig(filepath.Join(tmpDir, "app.log"))
	cfg.Level = DEBUG
	cfg.Formatter = Simple{}
	cfg.OutputConsole = false // 禁用控制台输出，确保写入文件
	cfg.LevelRouter = true

	logger := New(cfg)
	defer func() { _ = logger.Close() }()

	// 写入不同级别的日志
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	// 同步并等待
	_ = logger.Sync()
	_ = logger.Close()
	time.Sleep(200 * time.Millisecond)

	// 验证全量日志文件包含所有级别
	allContent := readLogFile(t, filepath.Join(tmpDir, "app.log"))
	if !strings.Contains(allContent, "debug message") {
		t.Error("all.log should contain debug message")
	}
	if !strings.Contains(allContent, "info message") {
		t.Error("all.log should contain info message")
	}
	if !strings.Contains(allContent, "warn message") {
		t.Error("all.log should contain warn message")
	}
	if !strings.Contains(allContent, "error message") {
		t.Error("all.log should contain error message")
	}

	// 验证级别专属文件只包含对应级别
	debugContent := readLogFile(t, filepath.Join(tmpDir, "DEBUG.log"))
	if !strings.Contains(debugContent, "debug message") {
		t.Error("DEBUG.log should contain debug message")
	}
	if strings.Contains(debugContent, "info message") {
		t.Error("DEBUG.log should not contain info message")
	}

	infoContent := readLogFile(t, filepath.Join(tmpDir, "INFO.log"))
	if !strings.Contains(infoContent, "info message") {
		t.Error("INFO.log should contain info message")
	}
	if strings.Contains(infoContent, "debug message") {
		t.Error("INFO.log should not contain debug message")
	}

	errorContent := readLogFile(t, filepath.Join(tmpDir, "ERROR.log"))
	if !strings.Contains(errorContent, "error message") {
		t.Error("ERROR.log should contain error message")
	}
	if strings.Contains(errorContent, "info message") {
		t.Error("ERROR.log should not contain info message")
	}
}

// TestLevelRouterWithLevelFiltering 测试级别路由配合级别过滤
func TestLevelRouterWithLevelFiltering(t *testing.T) {
	tmpDir := t.TempDir()

	// 创建配置，Level=WARN，启用级别路由
	cfg := NewConfig(filepath.Join(tmpDir, "app.log"))
	cfg.Level = WARN
	cfg.Formatter = Simple{}
	cfg.OutputConsole = false // 禁用控制台输出
	cfg.LevelRouter = true

	logger := New(cfg)
	defer func() { _ = logger.Close() }()

	// 写入各级别日志
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")

	_ = logger.Sync()
	_ = logger.Close()
	time.Sleep(200 * time.Millisecond)

	// 验证全量日志只包含 WARN 及以上
	allContent := readLogFile(t, filepath.Join(tmpDir, "app.log"))
	if strings.Contains(allContent, "debug") {
		t.Error("all.log should not contain debug when level is WARN")
	}
	if strings.Contains(allContent, "info") {
		t.Error("all.log should not contain info when level is WARN")
	}
	if !strings.Contains(allContent, "warn") {
		t.Error("all.log should contain warn")
	}
	if !strings.Contains(allContent, "error") {
		t.Error("all.log should contain error")
	}

	// 验证 DEBUG 和 INFO 专属文件不存在（因为 Level=WARN）
	if _, err := os.Stat(filepath.Join(tmpDir, "DEBUG.log")); !os.IsNotExist(err) {
		t.Error("DEBUG.log should not exist when level is WARN")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "INFO.log")); !os.IsNotExist(err) {
		t.Error("INFO.log should not exist when level is WARN")
	}

	// 验证 WARN 和 ERROR 专属文件存在且内容正确
	warnContent := readLogFile(t, filepath.Join(tmpDir, "WARN.log"))
	if !strings.Contains(warnContent, "warn") {
		t.Error("WARN.log should contain warn")
	}

	errorContent := readLogFile(t, filepath.Join(tmpDir, "ERROR.log"))
	if !strings.Contains(errorContent, "error") {
		t.Error("ERROR.log should contain error")
	}
}

// TestLevelRouterPathConflict 测试路径冲突检测
func TestLevelRouterPathConflict(t *testing.T) {
	// 测试 LogPath 与级别文件冲突的情况
	tmpDir := t.TempDir()
	cfg := NewConfig(filepath.Join(tmpDir, "INFO.log")) // 与级别文件冲突
	cfg.LevelRouter = true

	err := cfg.Validate()
	if err == nil {
		t.Fatal("should detect path conflict")
	}
	if !strings.Contains(err.Error(), "conflicts") {
		t.Errorf("error message should mention conflict, got: %v", err)
	}
}

// TestLevelRouterNoFileOutput 测试未启用文件输出时的错误
func TestLevelRouterNoFileOutput(t *testing.T) {
	cfg := Console() // 纯控制台输出
	cfg.LevelRouter = true

	err := cfg.Validate()
	if err == nil {
		t.Error("should require file output for level router")
	}
}

// TestLevelRouterCallerInfo 测试 Caller 信息正确性
func TestLevelRouterCallerInfo(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := NewConfig(filepath.Join(tmpDir, "app.log"))
	cfg.Level = INFO
	cfg.Caller = true
	cfg.Formatter = Def{} // 使用默认格式，包含 Caller
	cfg.LevelRouter = true

	logger := New(cfg)
	defer func() { _ = logger.Close() }()

	// 写入日志
	logger.Info("test caller")

	_ = logger.Sync()
	_ = logger.Close()
	time.Sleep(200 * time.Millisecond)

	// 验证 Caller 信息指向当前测试文件，不是 logger.go
	content := readLogFile(t, filepath.Join(tmpDir, "app.log"))
	if !strings.Contains(content, "logger_levelrouter_test.go") {
		t.Errorf("caller should point to test file, content: %s", content)
	}

	// 验证不包含 logger.go
	if strings.Contains(content, "logger.go") {
		t.Error("caller should not point to logger.go")
	}
}

// TestLevelRouterWithProdConfig 测试生产环境配置
func TestLevelRouterWithProdConfig(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := Prod(filepath.Join(tmpDir, "app.log"))
	cfg.Formatter = Simple{}
	cfg.OutputConsole = false // 禁用控制台输出
	cfg.LevelRouter = true

	logger := New(cfg)
	defer func() { _ = logger.Close() }()

	// Prod 配置级别是 WARN
	logger.Info("info") // 应该被过滤
	logger.Warn("warn")
	logger.Error("error")

	_ = logger.Sync()
	_ = logger.Close()
	time.Sleep(200 * time.Millisecond)

	// 验证 WARN、ERROR 的专属文件存在（因为实际写入了这些级别的日志）
	levels := []string{"WARN", "ERROR"}
	for _, lvl := range levels {
		path := filepath.Join(tmpDir, lvl+".log")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("%s.log should exist", lvl)
		}
	}

	// FATAL 和 PANIC 的专属文件可能不存在（因为没有实际写入这些级别的日志）
	// 这是正常行为，因为文件可能在首次写入时才创建

	// 验证 DEBUG、INFO 专属文件不存在
	if _, err := os.Stat(filepath.Join(tmpDir, "DEBUG.log")); !os.IsNotExist(err) {
		t.Error("DEBUG.log should not exist in prod config")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "INFO.log")); !os.IsNotExist(err) {
		t.Error("INFO.log should not exist in prod config")
	}
}

// TestLevelRouterDisabled 测试未启用时不创建专属文件
func TestLevelRouterDisabled(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := NewConfig(filepath.Join(tmpDir, "app.log"))
	cfg.Level = DEBUG
	cfg.Formatter = Simple{}
	cfg.OutputConsole = false // 禁用控制台输出
	cfg.LevelRouter = false   // 不启用

	logger := New(cfg)
	defer func() { _ = logger.Close() }()

	logger.Debug("debug")
	logger.Info("info")

	_ = logger.Sync()
	_ = logger.Close()
	time.Sleep(200 * time.Millisecond)

	// 验证只有主文件，没有级别专属文件
	if _, err := os.Stat(filepath.Join(tmpDir, "app.log")); os.IsNotExist(err) {
		t.Error("app.log should exist")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "DEBUG.log")); !os.IsNotExist(err) {
		t.Error("DEBUG.log should not exist when level router is disabled")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "INFO.log")); !os.IsNotExist(err) {
		t.Error("INFO.log should not exist when level router is disabled")
	}
}

// readLogFile 读取日志文件内容，如果不存在返回空字符串
func readLogFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ""
		}
		t.Fatalf("failed to read file %s: %v", path, err)
	}
	return string(content)
}
