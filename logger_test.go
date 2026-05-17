package fastlog

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

type syncWriter struct {
	*bytes.Buffer
	syncCalled bool
	syncErr    error
}

func (s *syncWriter) Sync() error {
	s.syncCalled = true
	return s.syncErr
}

func (s *syncWriter) Close() error {
	return nil
}

func TestNewLogger(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		l := New(&Config{
			Level:         INFO,
			OutputConsole: true,
		})
		// New 在成功时总是返回非 nil, 失败时会 panic
		// 这里不需要检查 nil, staticcheck 会误报
		_ = l.config.Level // 确保可以正常访问
	})

	t.Run("nil config panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("New(nil) should panic")
			}
		}()
		New(nil)
	})

	t.Run("invalid config panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("New with invalid config should panic")
			}
		}()
		New(&Config{
			// 未设置任何输出, 会触发 Validate 错误
		})
	})

	t.Run("default values applied", func(t *testing.T) {
		l := New(&Config{
			OutputConsole: true,
		})
		if l.config.Level != INFO {
			t.Errorf("default level should be INFO, got %v", l.config.Level)
		}
		if l.config.Formatter == nil {
			t.Errorf("default formatter should not be nil")
		}
	})
}

func TestLoggerLevelFilter(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(&Config{
		Level:         INFO,
		OutputConsole: true,
		Formatter:     &testFormatter{buf: buf},
	})

	l.Debug("should be suppressed")
	if buf.Len() > 0 {
		t.Errorf("Debug at INFO level should be suppressed, got: %q", buf.String())
	}

	l.Info("should pass")
	if buf.Len() == 0 {
		t.Errorf("Info at INFO level should pass")
	}
	if !strings.Contains(buf.String(), "should pass") {
		t.Errorf("Info output should contain message, got: %q", buf.String())
	}
}

func TestLoggerSamplerSuppress(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(&Config{
		Level:             DEBUG,
		OutputConsole:     true,
		Formatter:         &testFormatter{buf: buf},
		SamplerTick:       time.Minute,
		SamplerInitial:    0, // initial=0 → 第1条放行, 之后永久抑制
		SamplerThereafter: 0,
	})

	l.Info("same message")
	before := buf.Len()

	l.Info("same message")
	if buf.Len() != before {
		t.Errorf("Sampler should suppress after initial=0, second message wrote %d bytes", buf.Len()-before)
	}
}

func TestLoggerSamplerAllow(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(&Config{
		Level:             DEBUG,
		OutputConsole:     true,
		Formatter:         &testFormatter{buf: buf},
		SamplerTick:       time.Minute,
		SamplerInitial:    1, // initial=1 → 第1条放行
		SamplerThereafter: 0,
	})

	l.Info("first message")
	if buf.Len() == 0 {
		t.Errorf("Sampler should allow first message when initial=1")
	}
	if !strings.Contains(buf.String(), "first message") {
		t.Errorf("Output should contain message, got: %q", buf.String())
	}
}

func TestLoggerSync(t *testing.T) {
	sw := &syncWriter{Buffer: &bytes.Buffer{}}
	l := New(&Config{
		Level:         INFO,
		OutputConsole: true,
	})
	// 替换写入器为 syncWriter
	l.writer = sw

	l.Info("test")
	if err := l.Sync(); err != nil {
		t.Errorf("Sync() error = %v", err)
	}
	if !sw.syncCalled {
		t.Errorf("Sync() should delegate to underlying writer")
	}
}

func TestLoggerClose(t *testing.T) {
	m := &mockWriteCloser{Buffer: &bytes.Buffer{}}
	l := New(&Config{
		Level:         INFO,
		OutputConsole: true,
	})
	// 替换写入器为 mock
	l.writer = m

	l.Info("test")
	if err := l.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
	if !m.closed {
		t.Errorf("Close() should close underlying writer")
	}
}

func TestLoggerLevelsWriteData(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(&Config{
		Level:         DEBUG,
		OutputConsole: true,
		Formatter:     &testFormatter{buf: buf},
	})

	l.Debug("debug msg")
	l.Info("info msg")
	l.Warn("warn msg")
	l.Error("error msg")

	output := buf.String()
	if !strings.Contains(output, "debug msg") {
		t.Errorf("Debug should produce output at DEBUG level")
	}
	if !strings.Contains(output, "info msg") {
		t.Errorf("Info should produce output at DEBUG level")
	}
	if !strings.Contains(output, "warn msg") {
		t.Errorf("Warn should produce output at DEBUG level")
	}
	if !strings.Contains(output, "error msg") {
		t.Errorf("Error should produce output at DEBUG level")
	}
}

func TestLoggerFatalPanic(t *testing.T) {
	t.Run("Panic recovers", func(t *testing.T) {
		buf := &bytes.Buffer{}
		l := New(&Config{
			Level:         INFO,
			OutputConsole: true,
			Formatter:     &testFormatter{buf: buf},
		})

		panicked := false
		func() {
			defer func() {
				if r := recover(); r != nil {
					panicked = true
				}
			}()
			l.Panic("panic msg")
		}()

		if !panicked {
			t.Errorf("Panic should panic")
		}
		if buf.Len() == 0 {
			t.Errorf("Panic should write before panic")
		}
	})
}

func TestLoggerFields(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(&Config{
		Level:         INFO,
		OutputConsole: true,
		Formatter:     &testFormatter{buf: buf},
		Fields:        []Field{String("app", "fastlog")},
	})

	l.Infow("user login", String("user", "alice"))

	output := buf.String()
	if !strings.Contains(output, "app=fastlog") {
		t.Errorf("Should contain preset field 'app=fastlog', got: %q", output)
	}
	if !strings.Contains(output, "user=alice") {
		t.Errorf("Should contain local field 'user=alice', got: %q", output)
	}
}

func TestLoggerCaller(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(&Config{
		Level:         INFO,
		OutputConsole: true,
		Formatter:     &testFormatter{buf: buf},
		Caller:        true,
	})

	l.Info("with caller")

	output := buf.String()
	if !strings.Contains(output, "logger_test.go") {
		t.Errorf("Should contain caller info, got: %q", output)
	}
}

// testFormatter 用于测试的格式化器
type testFormatter struct {
	buf *bytes.Buffer
}

func (f *testFormatter) Format(entry *Entry) ([]byte, error) {
	var parts []string
	parts = append(parts, entry.Level.String())
	parts = append(parts, entry.Message)
	for _, field := range entry.Fields {
		parts = append(parts, field.Key()+"="+field.Value())
	}
	if entry.Caller != "" {
		parts = append(parts, entry.Caller)
	}
	result := strings.Join(parts, " ") + "\n"
	f.buf.WriteString(result)
	return []byte(result), nil
}

// ======== 格式化日志方法测试 ========

func TestLoggerDebugf(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(&Config{Level: DEBUG, OutputConsole: true, Formatter: &testFormatter{buf: buf}})
	l.Debugf("hello %s %d", "world", 42)
	if !strings.Contains(buf.String(), "hello world 42") {
		t.Errorf("Debugf output = %q, want 'hello world 42'", buf.String())
	}
}

func TestLoggerInfof(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(&Config{Level: DEBUG, OutputConsole: true, Formatter: &testFormatter{buf: buf}})
	l.Infof("count %d", 99)
	if !strings.Contains(buf.String(), "count 99") {
		t.Errorf("Infof output = %q, want 'count 99'", buf.String())
	}
}

func TestLoggerWarnf(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(&Config{Level: DEBUG, OutputConsole: true, Formatter: &testFormatter{buf: buf}})
	l.Warnf("disk %s %.1f%%", "/dev/sda", 85.5)
	if !strings.Contains(buf.String(), "disk /dev/sda 85.5%") {
		t.Errorf("Warnf output = %q", buf.String())
	}
}

func TestLoggerErrorf(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(&Config{Level: DEBUG, OutputConsole: true, Formatter: &testFormatter{buf: buf}})
	l.Errorf("err %s", "timeout")
	if !strings.Contains(buf.String(), "err timeout") {
		t.Errorf("Errorf output = %q", buf.String())
	}
}

// ======== 结构化日志方法测试 ========

func TestLoggerPanicf(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(&Config{Level: DEBUG, OutputConsole: true, Formatter: &testFormatter{buf: buf}})
	func() {
		defer func() { _ = recover() }()
		l.Panicf("panic %s %d", "test", 1)
	}()
	if !strings.Contains(buf.String(), "panic test 1") {
		t.Errorf("Panicf should write before panic, got: %q", buf.String())
	}
}

func TestLoggerDebugw(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(&Config{Level: DEBUG, OutputConsole: true, Formatter: &testFormatter{buf: buf}})
	l.Debugw("debug", String("k", "v"))
	if !strings.Contains(buf.String(), "k=v") {
		t.Errorf("Debugw output = %q, want 'k=v'", buf.String())
	}
}

func TestLoggerWarnw(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(&Config{Level: DEBUG, OutputConsole: true, Formatter: &testFormatter{buf: buf}})
	l.Warnw("warn", Int("code", 500))
	if !strings.Contains(buf.String(), "code=500") {
		t.Errorf("Warnw output = %q, want 'code=500'", buf.String())
	}
}

func TestLoggerErrorw(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(&Config{Level: DEBUG, OutputConsole: true, Formatter: &testFormatter{buf: buf}})
	l.Errorw("error", String("db", "mysql"), Int("retry", 3))
	output := buf.String()
	if !strings.Contains(output, "db=mysql") || !strings.Contains(output, "retry=3") {
		t.Errorf("Errorw output = %q, want 'db=mysql' and 'retry=3'", output)
	}
}

// ======== 格式化器/写入器错误路径测试 ========

// errFormatter 返回错误的格式化器
type errFormatter struct{}

func (e errFormatter) Format(_ *Entry) ([]byte, error) {
	return nil, errFmt
}

var errFmt = fmt.Errorf("format error")

// errWriter 写入时返回错误
type errWriter struct{}

func (e errWriter) Write(_ []byte) (int, error) {
	return 0, fmt.Errorf("write error")
}

func (e errWriter) Close() error { return nil }

func TestLoggerFormatterError(t *testing.T) {
	// 格式化错误应该通过 stderr 输出, 不 panic
	l := New(&Config{Level: DEBUG, OutputConsole: true, Formatter: errFormatter{}})
	l.Info("test")
	// 走到这里就算通过, 不应该 panic
}

func TestLoggerWriterError(t *testing.T) {
	l := New(&Config{Level: DEBUG, OutputConsole: true, Formatter: &testFormatter{buf: &bytes.Buffer{}}})
	l.writer = &errWriter{}
	l.Info("test")
	// 写入错误不应该 panic
}

// ======== EntryPool 测试 ========

func TestEntryPoolReuse(t *testing.T) {
	e1 := GetEntry()
	e1.Message = "original message"
	e1.Caller = "test.go:1"
	e1.Time = time.Now()
	e1.Fields = append(e1.Fields, String("k", "v"))
	PutEntry(e1)

	// 再次获取, 字段应已被清空
	e2 := GetEntry()
	if e2.Message != "" {
		t.Errorf("EntryPool reuse: Message should be empty, got %q", e2.Message)
	}
	if e2.Caller != "" {
		t.Errorf("EntryPool reuse: Caller should be empty, got %q", e2.Caller)
	}
	if !e2.Time.IsZero() {
		t.Errorf("EntryPool reuse: Time should be zero, got %v", e2.Time)
	}
	if len(e2.Fields) != 0 {
		t.Errorf("EntryPool reuse: Fields should be empty, got %d", len(e2.Fields))
	}
	PutEntry(e2)
}

// ======== getCaller 测试 ========

func TestGetCaller(t *testing.T) {
	// 跳过 1 层: test 函数本身
	caller := getCaller(1)
	if caller == "" {
		t.Errorf("getCaller(1) should not be empty")
	}
	if !strings.Contains(caller, "logger_test.go") {
		t.Errorf("getCaller(1) should contain test file name, got: %q", caller)
	}
	if !strings.Contains(caller, "TestGetCaller") {
		t.Errorf("getCaller(1) should contain function name, got: %q", caller)
	}
}

func TestGetCallerSkipTooDeep(t *testing.T) {
	caller := getCaller(100)
	if caller != "?:?:0" {
		t.Errorf("getCaller(100) should return '?:?:0', got: %q", caller)
	}
}

// ======== 级别过滤边界测试 ========

func TestLoggerFatalSuppressed(t *testing.T) {
	// FATAL 级别只放行 FATAL 和 PANIC
	buf := &bytes.Buffer{}
	l := New(&Config{Level: FATAL, OutputConsole: true, Formatter: &testFormatter{buf: buf}})
	l.Debug("debug")
	l.Info("info")
	l.Warn("warn")
	l.Error("error")
	if buf.Len() > 0 {
		t.Errorf("All logs below FATAL should be suppressed, got: %q", buf.String())
	}
}

func TestLoggerDebugEnabledAll(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(&Config{Level: DEBUG, OutputConsole: true, Formatter: &testFormatter{buf: buf}})
	l.Debug("d")
	l.Info("i")
	l.Warn("w")
	l.Error("e")
	output := buf.String()
	for _, level := range []string{"DEBUG", "INFO", "WARN", "ERROR"} {
		if !strings.Contains(output, level) {
			t.Errorf("DEBUG level should output %s, got: %q", level, output)
		}
	}
}

// ======== 关闭/同步边界测试 ========

func TestLoggerCloseNilWriter(t *testing.T) {
	l := New(&Config{Level: INFO, OutputConsole: true})
	// 第一次 Close 正常
	if err := l.Close(); err != nil {
		t.Errorf("First Close() error = %v", err)
	}
}

func TestLoggerSyncNoSyncer(t *testing.T) {
	l := New(&Config{Level: INFO, OutputConsole: true})
	// ConsoleWriter 没有 Sync 方法, Sync() 应返回 nil
	if err := l.Sync(); err != nil {
		t.Errorf("Sync() on non-sync writer should return nil, got: %v", err)
	}
}

// ======== 采样边界测试 ========

func TestLoggerSamplerDisabled(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(&Config{Level: DEBUG, OutputConsole: true, Formatter: &testFormatter{buf: buf}})
	for i := 0; i < 100; i++ {
		l.Info("all messages pass")
	}
	// 未配置采样器, 100 条都应输出
	lines := strings.Count(buf.String(), "\n")
	if lines != 100 {
		t.Errorf("Without sampler, expected 100 lines, got %d", lines)
	}
}

// ======== 预设字段加日志字段合并测试 ========

func TestLoggerPresetAndLocalFieldsMerge(t *testing.T) {
	buf := &bytes.Buffer{}
	// 带多字段的预设
	l := New(&Config{
		Level:         DEBUG,
		OutputConsole: true,
		Formatter:     &testFormatter{buf: buf},
		Fields:        []Field{String("env", "prod"), String("app", "mysvc")},
	})
	l.Infow("request", String("method", "GET"))
	output := buf.String()
	if !strings.Contains(output, "env=prod") {
		t.Errorf("Should contain preset field 'env=prod'")
	}
	if !strings.Contains(output, "app=mysvc") {
		t.Errorf("Should contain preset field 'app=mysvc'")
	}
	if !strings.Contains(output, "method=GET") {
		t.Errorf("Should contain local field 'method=GET'")
	}
}
