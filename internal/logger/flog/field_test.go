package flog

import (
	"errors"
	"testing"
	"time"
)

// TestStringField 测试字符串字段
func TestStringField(t *testing.T) {
	// 正常情况
	field := String("name", "test")
	if field == nil {
		t.Fatal("Expected non-nil field")
	}
	if field.Key() != "name" {
		t.Errorf("Expected key 'name', got '%s'", field.Key())
	}
	if field.Type() != StringType {
		t.Errorf("Expected StringType, got %v", field.Type())
	}
	if field.Value() != "test" {
		t.Errorf("Expected value 'test', got '%s'", field.Value())
	}

	// 空key情况
	emptyField := String("", "test")
	if emptyField != nil {
		t.Error("Expected nil for empty key")
	}
}

// TestIntField 测试整数字段
func TestIntField(t *testing.T) {
	field := Int("count", 42)
	if field == nil {
		t.Fatal("Expected non-nil field")
	}
	if field.Key() != "count" {
		t.Errorf("Expected key 'count', got '%s'", field.Key())
	}
	if field.Type() != IntType {
		t.Errorf("Expected IntType, got %v", field.Type())
	}
	if field.Value() != "42" {
		t.Errorf("Expected value '42', got '%s'", field.Value())
	}

	// 空key情况
	emptyField := Int("", 42)
	if emptyField != nil {
		t.Error("Expected nil for empty key")
	}
}

// TestFloat64Field 测试浮点数字段
func TestFloat64Field(t *testing.T) {
	field := Float64("price", 19.99)
	if field == nil {
		t.Fatal("Expected non-nil field")
	}
	if field.Key() != "price" {
		t.Errorf("Expected key 'price', got '%s'", field.Key())
	}
	if field.Type() != Float64Type {
		t.Errorf("Expected Float64Type, got %v", field.Type())
	}
	if field.Value() != "19.99" {
		t.Errorf("Expected value '19.99', got '%s'", field.Value())
	}
}

// TestBoolField 测试布尔字段
func TestBoolField(t *testing.T) {
	field := Bool("enabled", true)
	if field == nil {
		t.Fatal("Expected non-nil field")
	}
	if field.Key() != "enabled" {
		t.Errorf("Expected key 'enabled', got '%s'", field.Key())
	}
	if field.Type() != BoolType {
		t.Errorf("Expected BoolType, got %v", field.Type())
	}
	if field.Value() != "true" {
		t.Errorf("Expected value 'true', got '%s'", field.Value())
	}
}

// TestTimeField 测试时间字段
func TestTimeField(t *testing.T) {
	now := time.Now()
	field := Time("timestamp", now)
	if field == nil {
		t.Fatal("Expected non-nil field")
	}
	if field.Key() != "timestamp" {
		t.Errorf("Expected key 'timestamp', got '%s'", field.Key())
	}
	if field.Type() != TimeType {
		t.Errorf("Expected TimeType, got %v", field.Type())
	}
	// 时间格式应该符合RFC3339
	if field.Value() == "" {
		t.Error("Expected non-empty time value")
	}
}

// TestDurationField 测试持续时间段
func TestDurationField(t *testing.T) {
	duration := 5 * time.Second
	field := Duration("timeout", duration)
	if field == nil {
		t.Fatal("Expected non-nil field")
	}
	if field.Key() != "timeout" {
		t.Errorf("Expected key 'timeout', got '%s'", field.Key())
	}
	if field.Type() != DurationType {
		t.Errorf("Expected DurationType, got %v", field.Type())
	}
	if field.Value() != "5s" {
		t.Errorf("Expected value '5s', got '%s'", field.Value())
	}
}

// TestUint64Field 测试无符号64位整数字段
func TestUint64Field(t *testing.T) {
	field := Uint64("filesize", 1024)
	if field == nil {
		t.Fatal("Expected non-nil field")
	}
	if field.Key() != "filesize" {
		t.Errorf("Expected key 'filesize', got '%s'", field.Key())
	}
	if field.Type() != Uint64Type {
		t.Errorf("Expected Uint64Type, got %v", field.Type())
	}
	if field.Value() != "1024" {
		t.Errorf("Expected value '1024', got '%s'", field.Value())
	}
}

// TestErrorField 测试错误字段
func TestErrorField(t *testing.T) {
	err := errors.New("test error")
	field := Error("error", err)
	if field == nil {
		t.Fatal("Expected non-nil field")
	}
	if field.Key() != "error" {
		t.Errorf("Expected key 'error', got '%s'", field.Key())
	}
	if field.Type() != ErrorType {
		t.Errorf("Expected ErrorType, got %v", field.Type())
	}
	if field.Value() != "test error" {
		t.Errorf("Expected value 'test error', got '%s'", field.Value())
	}

	// 测试nil error
	nilField := Error("error", nil)
	if nilField == nil {
		t.Fatal("Expected non-nil field for nil error")
	}
	if nilField.Value() != "" {
		t.Errorf("Expected empty value for nil error, got '%s'", nilField.Value())
	}
}