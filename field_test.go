package fastlog

import (
	"errors"
	"testing"
	"time"
)

func TestFieldString(t *testing.T) {
	f := String("key", "value")
	if f.Key() != "key" {
		t.Errorf("Field.Key() = %q, want %q", f.Key(), "key")
	}
	if f.Type() != StringType {
		t.Errorf("Field.Type() = %v, want %v", f.Type(), StringType)
	}
	if f.Value() != "value" {
		t.Errorf("Field.Value() = %q, want %q", f.Value(), "value")
	}
}

func TestFieldInt(t *testing.T) {
	f := Int("count", 42)
	if f.Key() != "count" {
		t.Errorf("Field.Key() = %q, want %q", f.Key(), "count")
	}
	if f.Type() != IntType {
		t.Errorf("Field.Type() = %v, want %v", f.Type(), IntType)
	}
	if f.Value() != "42" {
		t.Errorf("Field.Value() = %q, want %q", f.Value(), "42")
	}
}

func TestFieldInt64(t *testing.T) {
	f := Int64("big", 1<<62)
	if f.Type() != Int64Type {
		t.Errorf("Field.Type() = %v, want %v", f.Type(), Int64Type)
	}
	if f.Value() != "4611686018427387904" {
		t.Errorf("Field.Value() = %q, want %q", f.Value(), "4611686018427387904")
	}
}

func TestFieldUint(t *testing.T) {
	f := Uint("count", 42)
	if f.Type() != UintType {
		t.Errorf("Field.Type() = %v, want %v", f.Type(), UintType)
	}
	if f.Value() != "42" {
		t.Errorf("Field.Value() = %q, want %q", f.Value(), "42")
	}
}

func TestFieldUint64(t *testing.T) {
	f := Uint64("big", 1<<63)
	if f.Type() != Uint64Type {
		t.Errorf("Field.Type() = %v, want %v", f.Type(), Uint64Type)
	}
	if f.Value() != "9223372036854775808" {
		t.Errorf("Field.Value() = %q, want %q", f.Value(), "9223372036854775808")
	}
}

func TestFieldFloat64(t *testing.T) {
	f := Float64("pi", 3.14)
	if f.Type() != Float64Type {
		t.Errorf("Field.Type() = %v, want %v", f.Type(), Float64Type)
	}
	if f.Value() != "3.14" {
		t.Errorf("Field.Value() = %q, want %q", f.Value(), "3.14")
	}
}

func TestFieldBool(t *testing.T) {
	f := Bool("flag", true)
	if f.Type() != BoolType {
		t.Errorf("Field.Type() = %v, want %v", f.Type(), BoolType)
	}
	if f.Value() != "true" {
		t.Errorf("Field.Value() = %q, want %q", f.Value(), "true")
	}
}

func TestFieldTime(t *testing.T) {
	now := time.Date(2025, 1, 15, 10, 30, 45, 0, time.UTC)
	f := Time("created", now)
	if f.Type() != TimeType {
		t.Errorf("Field.Type() = %v, want %v", f.Type(), TimeType)
	}
	if f.Value() != "2025-01-15T10:30:45Z" {
		t.Errorf("Field.Value() = %q, want %q", f.Value(), "2025-01-15T10:30:45Z")
	}
}

func TestFieldDuration(t *testing.T) {
	d := 5 * time.Second
	f := Duration("dur", d)
	if f.Type() != DurationType {
		t.Errorf("Field.Type() = %v, want %v", f.Type(), DurationType)
	}
	if f.Value() != "5s" {
		t.Errorf("Field.Value() = %q, want %q", f.Value(), "5s")
	}
}

func TestFieldError(t *testing.T) {
	t.Run("non-nil error", func(t *testing.T) {
		f := Error(errors.New("something wrong"))
		if f.Key() != "error" {
			t.Errorf("Field.Key() = %q, want %q", f.Key(), "error")
		}
		if f.Type() != ErrorType {
			t.Errorf("Field.Type() = %v, want %v", f.Type(), ErrorType)
		}
		if f.Value() != "something wrong" {
			t.Errorf("Field.Value() = %q, want %q", f.Value(), "something wrong")
		}
	})

	t.Run("nil error", func(t *testing.T) {
		f := Error(nil)
		if f.Key() != "error" {
			t.Errorf("Field.Key() = %q, want %q", f.Key(), "error")
		}
		if f.Type() != StringType {
			t.Errorf("Field.Type() = %v, want %v", f.Type(), StringType)
		}
		if f.Value() != "<nil>" {
			t.Errorf("Field.Value() = %q, want %q", f.Value(), "<nil>")
		}
	})
}

func TestFieldErr(t *testing.T) {
	f := Err("db_error", errors.New("timeout"))
	if f.Key() != "db_error" {
		t.Errorf("Field.Key() = %q, want %q", f.Key(), "db_error")
	}
	if f.Type() != ErrorType {
		t.Errorf("Field.Type() = %v, want %v", f.Type(), ErrorType)
	}
	if f.Value() != "timeout" {
		t.Errorf("Field.Value() = %q, want %q", f.Value(), "timeout")
	}
}

func TestFieldAny(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		f := Any("data", "hello")
		if f.Type() != AnyType {
			t.Errorf("Field.Type() = %v, want %v", f.Type(), AnyType)
		}
		if f.Value() != "hello" {
			t.Errorf("Field.Value() = %q, want %q", f.Value(), "hello")
		}
	})

	t.Run("int", func(t *testing.T) {
		f := Any("num", 42)
		if f.Value() != "42" {
			t.Errorf("Field.Value() = %q, want %q", f.Value(), "42")
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Any("empty", nil)
		if f.Value() != "null" {
			t.Errorf("Field.Value() = %q, want %q", f.Value(), "null")
		}
	})
}

func TestFieldStack(t *testing.T) {
	f := Stack()
	if f.Key() != "stack" {
		t.Errorf("Field.Key() = %q, want %q", f.Key(), "stack")
	}
	if f.Type() != StringType {
		t.Errorf("Field.Type() = %v, want %v", f.Type(), StringType)
	}
	if f.Value() == "" {
		t.Error("Field.Value() should not be empty for Stack()")
	}
}

// ======== toInterface 测试 ========

func TestFieldToInterface(t *testing.T) {
	t.Run("StringType", func(t *testing.T) {
		f := String("k", "hello")
		if got := f.toInterface(); got != "hello" {
			t.Errorf("toInterface() = %v, want %v", got, "hello")
		}
	})

	t.Run("IntType", func(t *testing.T) {
		f := Int("k", -42)
		if got := f.toInterface(); got != int64(-42) {
			t.Errorf("toInterface() = %v, want %v", got, int64(-42))
		}
	})

	t.Run("Int64Type", func(t *testing.T) {
		f := Int64("k", 1<<62)
		if got := f.toInterface(); got != int64(1<<62) {
			t.Errorf("toInterface() = %v", got)
		}
	})

	t.Run("UintType", func(t *testing.T) {
		f := Uint("k", 99)
		if got := f.toInterface(); got != uint64(99) {
			t.Errorf("toInterface() = %v, want %v", got, uint64(99))
		}
	})

	t.Run("Uint64Type", func(t *testing.T) {
		f := Uint64("k", 1<<63)
		if got := f.toInterface(); got != uint64(1<<63) {
			t.Errorf("toInterface() = %v", got)
		}
	})

	t.Run("Float64Type", func(t *testing.T) {
		f := Float64("k", 3.14)
		if got := f.toInterface(); got != float64(3.14) {
			t.Errorf("toInterface() = %v, want %v", got, 3.14)
		}
	})

	t.Run("BoolType", func(t *testing.T) {
		f := Bool("k", true)
		if got := f.toInterface(); got != true {
			t.Errorf("toInterface() = %v, want true", got)
		}
	})

	t.Run("TimeType", func(t *testing.T) {
		now := time.Date(2025, 1, 15, 10, 30, 45, 0, time.UTC)
		f := Time("k", now)
		if got := f.toInterface(); got != "2025-01-15T10:30:45Z" {
			t.Errorf("toInterface() = %v, want RFC3339 string", got)
		}
	})

	t.Run("DurationType", func(t *testing.T) {
		f := Duration("k", 5*time.Second)
		if got := f.toInterface(); got != "5s" {
			t.Errorf("toInterface() = %v, want '5s'", got)
		}
	})

	t.Run("ErrorType", func(t *testing.T) {
		f := Error(errNoServer)
		if got := f.toInterface(); got != "no server" {
			t.Errorf("toInterface() = %v, want 'no server'", got)
		}
	})

	t.Run("AnyType", func(t *testing.T) {
		f := Any("k", map[string]int{"a": 1})
		if got := f.toInterface(); got == nil {
			t.Errorf("toInterface() for AnyType should return original value")
		}
	})

	t.Run("UnknownType", func(t *testing.T) {
		f := Field{key: "k", typ: UnknownType}
		if got := f.toInterface(); got != nil {
			t.Errorf("toInterface() for UnknownType = %v, want nil", got)
		}
	})
}

var errNoServer = errors.New("no server")

// ======== Any 更多类型测试 ========

func TestFieldAnyExtraTypes(t *testing.T) {
	t.Run("int8", func(t *testing.T) {
		if got := Any("k", int8(8)).Value(); got != "8" {
			t.Errorf("Any(int8) Value = %q, want '8'", got)
		}
	})

	t.Run("int16", func(t *testing.T) {
		if got := Any("k", int16(16)).Value(); got != "16" {
			t.Errorf("Any(int16) Value = %q, want '16'", got)
		}
	})

	t.Run("int32", func(t *testing.T) {
		if got := Any("k", int32(32)).Value(); got != "32" {
			t.Errorf("Any(int32) Value = %q, want '32'", got)
		}
	})

	t.Run("uint8", func(t *testing.T) {
		if got := Any("k", uint8(8)).Value(); got != "8" {
			t.Errorf("Any(uint8) Value = %q, want '8'", got)
		}
	})

	t.Run("uint16", func(t *testing.T) {
		if got := Any("k", uint16(16)).Value(); got != "16" {
			t.Errorf("Any(uint16) Value = %q, want '16'", got)
		}
	})

	t.Run("uint32", func(t *testing.T) {
		if got := Any("k", uint32(32)).Value(); got != "32" {
			t.Errorf("Any(uint32) Value = %q, want '32'", got)
		}
	})

	t.Run("float32", func(t *testing.T) {
		if got := Any("k", float32(3.14)).Value(); got != "3.14" {
			t.Errorf("Any(float32) Value = %q, want '3.14'", got)
		}
	})

	t.Run("error", func(t *testing.T) {
		if got := Any("k", errNoServer).Value(); got != "no server" {
			t.Errorf("Any(error) Value = %q, want 'no server'", got)
		}
	})

	t.Run("time.Duration", func(t *testing.T) {
		if got := Any("k", 3*time.Minute).Value(); got != "3m0s" {
			t.Errorf("Any(Duration) Value = %q, want '3m0s'", got)
		}
	})

	t.Run("nil", func(t *testing.T) {
		if got := Any("k", nil).Value(); got != "null" {
			t.Errorf("Any(nil) Value = %q, want 'null'", got)
		}
	})

	t.Run("unrecognized struct type", func(t *testing.T) {
		type custom struct{ x int }
		if got := Any("k", custom{x: 1}).Value(); got != "" {
			t.Errorf("Any(struct{}) Value = %q, want ''", got)
		}
	})
}

// ======== 边界值测试 ========

func TestFieldValueBoundary(t *testing.T) {
	t.Run("negative int", func(t *testing.T) {
		if got := Int("k", -1).Value(); got != "-1" {
			t.Errorf("Int(-1) Value = %q, want '-1'", got)
		}
	})

	t.Run("max uint64", func(t *testing.T) {
		if got := Uint64("k", 1<<64-1).Value(); got != "18446744073709551615" {
			t.Errorf("Uint64(max) Value = %q", got)
		}
	})

	t.Run("zero time", func(t *testing.T) {
		if got := Time("k", time.Time{}).Value(); got != "0001-01-01T00:00:00Z" {
			t.Errorf("Time(zero) Value = %q", got)
		}
	})

	t.Run("zero duration", func(t *testing.T) {
		if got := Duration("k", 0).Value(); got != "0s" {
			t.Errorf("Duration(0) Value = %q, want '0s'", got)
		}
	})

	t.Run("large float", func(t *testing.T) {
		if got := Float64("k", 1e20).Value(); got != "1e+20" {
			t.Errorf("Float64(1e20) Value = %q, want '1e+20'", got)
		}
	})
}

// ======== UnknownType 返回空字符串 ========

func TestFieldUnknownType(t *testing.T) {
	f := Field{key: "unknown", typ: UnknownType}
	if got := f.Value(); got != "" {
		t.Errorf("UnknownType Value = %q, want ''", got)
	}
}

// ======== Err(nil) 自定义键名 ========

func TestFieldErrNilCustomKey(t *testing.T) {
	f := Err("db_error", nil)
	if f.Key() != "db_error" {
		t.Errorf("Err(nil) Key = %q, want 'db_error'", f.Key())
	}
	if f.Type() != StringType {
		t.Errorf("Err(nil) Type = %v, want StringType", f.Type())
	}
	if f.Value() != "<nil>" {
		t.Errorf("Err(nil) Value = %q, want '<nil>'", f.Value())
	}
}

// ======== 空键名 ========

func TestFieldEmptyKey(t *testing.T) {
	f := String("", "value")
	if f.Key() != "" {
		t.Errorf("Empty key should be empty, got %q", f.Key())
	}
	if f.Value() != "value" {
		t.Errorf("Empty key field Value = %q, want 'value'", f.Value())
	}
}
