package fastlog

import (
	"encoding/json"
	"fmt"
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

// ======== JSON 所有字段类型测试 ========

func TestJSONFormatAllFieldTypes(t *testing.T) {
	entry := makeEntry("test", "",
		String("str", "hello"),
		Int("int", -42),
		Int64("int64", 1<<62),
		Uint("uint", 99),
		Uint64("uint64", 1<<63),
		Float64("float", 3.14),
		Bool("bool", true),
		Time("time", time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)),
		Duration("dur", 5*time.Second),
		Error(errFmtTest),
		Any("any", "anyval"),
	)

	b, err := JSON{}.Format(entry)
	if err != nil {
		t.Fatalf("JSON.Format() error = %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	tests := []struct {
		key  string
		want interface{}
	}{
		{"str", "hello"},
		{"int", float64(-42)},
		{"int64", float64(1 << 62)},
		{"uint", float64(99)},
		{"uint64", float64(1 << 63)},
		{"float", 3.14},
		{"bool", true},
		{"time", "2025-06-01T12:00:00Z"},
		{"dur", "5s"},
		{"error", "test error"},
		{"any", "anyval"},
	}

	for _, tc := range tests {
		t.Run(tc.key, func(t *testing.T) {
			if data[tc.key] != tc.want {
				t.Errorf("JSON field %q = %v (%T), want %v (%T)", tc.key, data[tc.key], data[tc.key], tc.want, tc.want)
			}
		})
	}
}

// ======== formatField 函数测试 ========

var errFmtTest = fmt.Errorf("test error")

func TestFormatField(t *testing.T) {
	t.Run("string field", func(t *testing.T) {
		f := String("key", "value")
		if got := formatField(f); got != "key=value" {
			t.Errorf("formatField() = %q, want 'key=value'", got)
		}
	})

	t.Run("int field", func(t *testing.T) {
		if got := formatField(Int("n", 42)); got != "n=42" {
			t.Errorf("formatField(Int) = %q", got)
		}
	})

	t.Run("bool field", func(t *testing.T) {
		if got := formatField(Bool("flag", true)); got != "flag=true" {
			t.Errorf("formatField(Bool) = %q", got)
		}
	})
}

// ======== 未知级别测试 ========

func TestDefFormatUnknownLevel(t *testing.T) {
	entry := &Entry{
		Time:    time.Date(2026, 1, 15, 10, 30, 45, 0, time.UTC),
		Level:   Level(99),
		Message: "test",
	}
	b, err := Def{}.Format(entry)
	if err != nil {
		t.Fatalf("Def.Format() unknown level error = %v", err)
	}
	if !strings.Contains(string(b), "Level(99)") {
		t.Errorf("Def.Format() with unknown level should contain 'Level(99)', got: %q", string(b))
	}
}

func TestJSONFormatUnknownLevel(t *testing.T) {
	entry := &Entry{
		Time:    time.Date(2026, 1, 15, 10, 30, 45, 0, time.UTC),
		Level:   Level(99),
		Message: "test",
	}
	b, err := JSON{}.Format(entry)
	if err != nil {
		t.Fatalf("JSON.Format() unknown level error = %v", err)
	}
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if data["level"] != "Level(99)" {
		t.Errorf("JSON unknown level = %v, want 'Level(99)'", data["level"])
	}
}

// ======== 多字段格式测试 ========

func TestDefFormatMultipleFields(t *testing.T) {
	entry := makeEntry("multi", "",
		String("a", "1"),
		String("b", "2"),
		String("c", "3"),
	)
	b, _ := Def{}.Format(entry)
	want := "2026-01-15T10:30:45Z | INFO   | multi a=1, b=2, c=3\n"
	if got := string(b); got != want {
		t.Errorf("Def.Format() multiple fields = %q, want %q", got, want)
	}
}

func TestKVFormatMultipleFields(t *testing.T) {
	entry := makeEntry("multi", "",
		String("a", "1"),
		String("b", "2"),
	)
	b, _ := KV{}.Format(entry)
	want := "time=2026-01-15T10:30:45Z level=INFO message=multi a=1 b=2\n"
	if got := string(b); got != want {
		t.Errorf("KV.Format() multiple fields = %q, want %q", got, want)
	}
}

// ======== JSON nil/空字段 ========

func TestJSONFormatNilFields(t *testing.T) {
	entry := makeEntry("hello", "")
	entry.Fields = nil
	b, err := JSON{}.Format(entry)
	if err != nil {
		t.Fatalf("JSON.Format() nil fields error = %v", err)
	}
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if data["message"] != "hello" {
		t.Errorf("JSON nil fields message = %v", data["message"])
	}
}

func TestJSONFormatEmptyFields(t *testing.T) {
	entry := makeEntry("hello", "")
	entry.Fields = []Field{}
	b, err := JSON{}.Format(entry)
	if err != nil {
		t.Fatalf("JSON.Format() empty fields error = %v", err)
	}
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if data["message"] != "hello" {
		t.Errorf("JSON empty fields message = %v", data["message"])
	}
}
