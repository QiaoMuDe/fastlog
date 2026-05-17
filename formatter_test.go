package fastlog

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func makeEntry(msg, caller string, fields ...Field) *Entry {
	return &Entry{
		Time:    time.Date(2026, 1, 15, 10, 30, 45, 0, time.UTC),
		Level:   INFO,
		Message: msg,
		Caller:  caller,
		Fields:  fields,
	}
}

func TestDefFormat(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		entry := makeEntry("hello", "")
		b, err := Def{}.Format(entry)
		if err != nil {
			t.Fatalf("Format() error = %v", err)
		}
		want := "2026-01-15T10:30:45Z | INFO   | hello\n"
		if got := string(b); got != want {
			t.Errorf("Def.Format() = %q, want %q", got, want)
		}
	})

	t.Run("with caller", func(t *testing.T) {
		entry := makeEntry("hello", "main.go:main:10")
		b, _ := Def{}.Format(entry)
		want := "2026-01-15T10:30:45Z | INFO   | main.go:main:10 - hello\n"
		if got := string(b); got != want {
			t.Errorf("Def.Format() with caller = %q, want %q", got, want)
		}
	})

	t.Run("with fields", func(t *testing.T) {
		entry := makeEntry("hello", "", String("k", "v"))
		b, _ := Def{}.Format(entry)
		want := "2026-01-15T10:30:45Z | INFO   | hello k=v\n"
		if got := string(b); got != want {
			t.Errorf("Def.Format() with fields = %q, want %q", got, want)
		}
	})

	t.Run("empty message", func(t *testing.T) {
		entry := makeEntry("", "")
		b, err := Def{}.Format(entry)
		if err != nil {
			t.Fatalf("Format() empty msg error = %v", err)
		}
		if !strings.HasSuffix(string(b), "\n") {
			t.Errorf("Def.Format() should end with newline")
		}
	})

	t.Run("nil fields", func(t *testing.T) {
		entry := makeEntry("hello", "")
		entry.Fields = nil
		b, err := Def{}.Format(entry)
		if err != nil {
			t.Fatalf("Format() nil fields error = %v", err)
		}
		want := "2026-01-15T10:30:45Z | INFO   | hello\n"
		if got := string(b); got != want {
			t.Errorf("Def.Format() nil fields = %q, want %q", got, want)
		}
	})
}

func TestJSONFormat(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		entry := makeEntry("hello", "")
		b, err := JSON{}.Format(entry)
		if err != nil {
			t.Fatalf("Format() error = %v", err)
		}

		// 确认末尾有换行
		if b[len(b)-1] != '\n' {
			t.Errorf("JSON format should end with newline")
		}

		// 反序列化验证字段
		var data map[string]interface{}
		if err := json.Unmarshal(b, &data); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}

		if data["level"] != "INFO" {
			t.Errorf("JSON level = %v, want INFO", data["level"])
		}
		if data["message"] != "hello" {
			t.Errorf("JSON message = %v, want hello", data["message"])
		}
		if _, ok := data["caller"]; ok {
			t.Errorf("JSON should not have caller field when empty")
		}
	})

	t.Run("with caller and fields", func(t *testing.T) {
		entry := makeEntry("hello", "main.go:main:10", String("key", "val"))
		b, _ := JSON{}.Format(entry)

		var data map[string]interface{}
		_ = json.Unmarshal(b, &data)

		if data["caller"] != "main.go:main:10" {
			t.Errorf("JSON caller = %v, want main.go:main:10", data["caller"])
		}
		if data["key"] != "val" {
			t.Errorf("JSON key = %v, want val", data["key"])
		}
		if data["time"] != "2026-01-15T10:30:45Z" {
			t.Errorf("JSON time = %v, want 2026-01-15T10:30:45Z", data["time"])
		}
	})
}

func TestTimestampFormat(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		entry := makeEntry("hello", "")
		b, _ := Timestamp{}.Format(entry)
		want := "2026-01-15T10:30:45Z INFO hello\n"
		if got := string(b); got != want {
			t.Errorf("Timestamp.Format() = %q, want %q", got, want)
		}
	})

	t.Run("with fields", func(t *testing.T) {
		entry := makeEntry("hello", "", String("k", "v"))
		b, _ := Timestamp{}.Format(entry)
		want := "2026-01-15T10:30:45Z INFO hello k=v\n"
		if got := string(b); got != want {
			t.Errorf("Timestamp.Format() with fields = %q, want %q", got, want)
		}
	})
}

func TestKVFormat(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		entry := makeEntry("hello", "")
		b, _ := KV{}.Format(entry)
		want := "time=2026-01-15T10:30:45Z level=INFO message=hello\n"
		if got := string(b); got != want {
			t.Errorf("KV.Format() = %q, want %q", got, want)
		}
	})

	t.Run("with caller", func(t *testing.T) {
		entry := makeEntry("hello", "main.go:main:10")
		b, _ := KV{}.Format(entry)
		want := "time=2026-01-15T10:30:45Z level=INFO message=hello caller=main.go:main:10\n"
		if got := string(b); got != want {
			t.Errorf("KV.Format() with caller = %q, want %q", got, want)
		}
	})

	t.Run("with fields", func(t *testing.T) {
		entry := makeEntry("hello", "", String("k", "v"))
		b, _ := KV{}.Format(entry)
		want := "time=2026-01-15T10:30:45Z level=INFO message=hello k=v\n"
		if got := string(b); got != want {
			t.Errorf("KV.Format() with fields = %q, want %q", got, want)
		}
	})
}

func TestLogFmtFormat(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		entry := makeEntry("hello", "")
		b, _ := LogFmt{}.Format(entry)
		want := "2026-01-15T10:30:45Z [INFO ] hello\n"
		if got := string(b); got != want {
			t.Errorf("LogFmt.Format() = %q, want %q", got, want)
		}
	})

	t.Run("with caller", func(t *testing.T) {
		entry := makeEntry("hello", "main.go:main:10")
		b, _ := LogFmt{}.Format(entry)
		want := "2026-01-15T10:30:45Z [INFO ] main.go:main:10 hello\n"
		if got := string(b); got != want {
			t.Errorf("LogFmt.Format() with caller = %q, want %q", got, want)
		}
	})

	t.Run("with fields", func(t *testing.T) {
		entry := makeEntry("hello", "", String("k", "v"))
		b, _ := LogFmt{}.Format(entry)
		want := "2026-01-15T10:30:45Z [INFO ] hello [k=v]\n"
		if got := string(b); got != want {
			t.Errorf("LogFmt.Format() with fields = %q, want %q", got, want)
		}
	})
}
