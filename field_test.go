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
