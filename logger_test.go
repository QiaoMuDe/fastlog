package fastlog

import (
	"bytes"
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
