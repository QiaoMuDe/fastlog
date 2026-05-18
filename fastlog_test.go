package fastlog

import (
	"bytes"
	"testing"
)

func TestLevelString(t *testing.T) {
	tests := []struct {
		name string
		l    Level
		want string
	}{
		{"DEBUG", DEBUG, "DEBUG"},
		{"INFO", INFO, "INFO"},
		{"WARN", WARN, "WARN"},
		{"ERROR", ERROR, "ERROR"},
		{"FATAL", FATAL, "FATAL"},
		{"PANIC", PANIC, "PANIC"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.String(); got != tt.want {
				t.Errorf("Level.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLevelEnabled(t *testing.T) {
	tests := []struct {
		name  string
		level Level
		lvl   Level
		want  bool
	}{
		{"INFO enabled INFO", INFO, INFO, true},
		{"INFO enabled WARN", INFO, WARN, true},
		{"INFO enabled ERROR", INFO, ERROR, true},
		{"INFO enabled FATAL", INFO, FATAL, true},
		{"INFO enabled PANIC", INFO, PANIC, true},
		{"INFO not enabled DEBUG", INFO, DEBUG, false},

		{"DEBUG enabled all", DEBUG, PANIC, true},
		{"PANIC not enabled FATAL", PANIC, FATAL, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.Enabled(tt.lvl); got != tt.want {
				t.Errorf("Level(%d).Enabled(%d) = %v, want %v", tt.level, tt.lvl, got, tt.want)
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    Level
		wantErr bool
	}{
		{"lowercase", "debug", DEBUG, false},
		{"uppercase", "DEBUG", DEBUG, false},
		{"mixed case", "Debug", DEBUG, false},
		{"info", "INFO", INFO, false},
		{"warn", "WARN", WARN, false},
		{"error", "ERROR", ERROR, false},
		{"fatal", "FATAL", FATAL, false},
		{"panic", "PANIC", PANIC, false},
		{"unknown string", "unknown", INFO, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLevel(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLevel(%q) error = %v, wantErr %v", tt.s, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestAllLevels(t *testing.T) {
	levels := AllLevels()
	want := []Level{DEBUG, INFO, WARN, ERROR, FATAL, PANIC}
	if len(levels) != len(want) {
		t.Errorf("AllLevels() length = %d, want %d", len(levels), len(want))
		return
	}
	for i := range want {
		if levels[i] != want[i] {
			t.Errorf("AllLevels()[%d] = %v, want %v", i, levels[i], want[i])
		}
	}
}

// TestLoggerDynamicLevel 测试动态级别调整功能
func TestLoggerDynamicLevel(t *testing.T) {
	t.Run("initial level from config", func(t *testing.T) {
		l := New(&Config{
			Level:         WARN,
			OutputConsole: true,
			TimeFormat:    DefaultTimeFormat,
		})

		if l.Level() != WARN {
			t.Errorf("initial level should be WARN, got %v", l.Level())
		}
	})

	t.Run("set level changes behavior", func(t *testing.T) {
		buf := &bytes.Buffer{}
		l := New(&Config{
			Level:         INFO,
			OutputConsole: true,
			TimeFormat:    DefaultTimeFormat,
			Formatter:     &testFormatter{buf: buf},
		})

		// INFO 级别可以输出 INFO 日志
		buf.Reset()
		l.Info("info message")
		if buf.Len() == 0 {
			t.Errorf("INFO should pass at INFO level")
		}

		// 提升到 WARN 级别
		l.SetLevel(WARN)
		if l.Level() != WARN {
			t.Errorf("level should be WARN after SetLevel, got %v", l.Level())
		}

		// WARN 级别抑制 INFO 日志
		buf.Reset()
		l.Info("should be suppressed")
		if buf.Len() > 0 {
			t.Errorf("INFO should be suppressed at WARN level, got: %q", buf.String())
		}

		// WARN 级别允许 WARN 日志
		buf.Reset()
		l.Warn("warn message")
		if buf.Len() == 0 {
			t.Errorf("WARN should pass at WARN level")
		}
	})

	t.Run("set level to DEBUG enables all", func(t *testing.T) {
		buf := &bytes.Buffer{}
		l := New(&Config{
			Level:         ERROR,
			OutputConsole: true,
			TimeFormat:    DefaultTimeFormat,
			Formatter:     &testFormatter{buf: buf},
		})

		// ERROR 级别抑制 INFO
		buf.Reset()
		l.Info("should be suppressed")
		if buf.Len() > 0 {
			t.Errorf("INFO should be suppressed at ERROR level")
		}

		// 降级到 DEBUG
		l.SetLevel(DEBUG)

		// DEBUG 级别允许所有日志
		buf.Reset()
		l.Debug("debug message")
		if buf.Len() == 0 {
			t.Errorf("DEBUG should pass at DEBUG level")
		}

		buf.Reset()
		l.Info("info message")
		if buf.Len() == 0 {
			t.Errorf("INFO should pass at DEBUG level")
		}
	})

	t.Run("multiple level changes", func(t *testing.T) {
		l := New(&Config{
			Level:         INFO,
			OutputConsole: true,
			TimeFormat:    DefaultTimeFormat,
		})

		// 多次切换级别
		l.SetLevel(DEBUG)
		if l.Level() != DEBUG {
			t.Errorf("level should be DEBUG")
		}

		l.SetLevel(WARN)
		if l.Level() != WARN {
			t.Errorf("level should be WARN")
		}

		l.SetLevel(ERROR)
		if l.Level() != ERROR {
			t.Errorf("level should be ERROR")
		}

		l.SetLevel(INFO)
		if l.Level() != INFO {
			t.Errorf("level should be INFO")
		}
	})

	t.Run("set level concurrent safe", func(t *testing.T) {
		l := New(&Config{
			Level:         INFO,
			OutputConsole: true,
			TimeFormat:    DefaultTimeFormat,
		})

		// 并发设置级别和读取级别
		done := make(chan bool, 10)

		// 5 个 goroutine 设置不同级别
		for i := 0; i < 5; i++ {
			go func(idx int) {
				levels := []Level{DEBUG, INFO, WARN, ERROR, FATAL}
				l.SetLevel(levels[idx%5])
				done <- true
			}(i)
		}

		// 5 个 goroutine 读取级别
		for i := 0; i < 5; i++ {
			go func() {
				_ = l.Level()
				done <- true
			}()
		}

		// 等待所有 goroutine 完成
		for i := 0; i < 10; i++ {
			<-done
		}

		// 如果并发不安全，上面的代码会触发 data race
		// go test -race 会检测到
	})

	t.Run("level affects all log methods", func(t *testing.T) {
		buf := &bytes.Buffer{}
		l := New(&Config{
			Level:         PANIC,
			OutputConsole: true,
			TimeFormat:    DefaultTimeFormat,
			Formatter:     &testFormatter{buf: buf},
		})

		// PANIC 级别抑制所有其他级别
		buf.Reset()
		l.Debug("debug")
		l.Info("info")
		l.Warn("warn")
		l.Error("error")
		if buf.Len() > 0 {
			t.Errorf("all levels should be suppressed at PANIC level")
		}

		// 降级到 DEBUG
		l.SetLevel(DEBUG)

		// DEBUG 级别允许所有级别
		l.Debug("debug2")
		l.Info("info2")
		l.Warn("warn2")
		l.Error("error2")
		// 注意：不测试 FATAL 和 PANIC，因为它们会退出/崩溃
	})
}
