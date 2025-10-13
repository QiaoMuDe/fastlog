package flog

import (
	"strconv"
	"time"
)

// FieldLogger 字段日志构建器
type FieldLogger struct {
	fields map[string]string // 自定义字段, 用于在日志中添加额外的上下文信息
}

// String 添加字符串字段，支持链式调用
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为字符串。
//
// 返回值:
//   - *FieldLogger: 一个指向 FieldLogger 实例的指针, 用于链式调用。
func (fl *FieldLogger) String(key string, value string) *FieldLogger {
	if fl != nil && fl.fields != nil && key != "" {
		fl.fields[key] = value
	}
	return fl
}

// Int 添加整数字段，支持链式调用
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为整数。
//
// 返回值:
//   - *FieldLogger: 一个指向 FieldLogger 实例的指针, 用于链式调用。
func (fl *FieldLogger) Int(key string, value int) *FieldLogger {
	if fl != nil && fl.fields != nil && key != "" {
		fl.fields[key] = strconv.Itoa(value)
	}
	return fl
}

// Int64 添加64位整数字段，支持链式调用
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为64位整数。
//
// 返回值:
//   - *FieldLogger: 一个指向 FieldLogger 实例的指针, 用于链式调用。
func (fl *FieldLogger) Int64(key string, value int64) *FieldLogger {
	if fl != nil && fl.fields != nil && key != "" {
		fl.fields[key] = strconv.FormatInt(value, 10)
	}
	return fl
}

// Float64 添加64位浮点数字段，支持链式调用
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为64位浮点数。
//
// 返回值:
//   - *FieldLogger: 一个指向 FieldLogger 实例的指针, 用于链式调用。
func (fl *FieldLogger) Float64(key string, value float64) *FieldLogger {
	if fl != nil && fl.fields != nil && key != "" {
		fl.fields[key] = strconv.FormatFloat(value, 'f', 6, 64)
	}
	return fl
}

// Bool 添加布尔字段，支持链式调用
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为布尔值。
//
// 返回值:
//   - *FieldLogger: 一个指向 FieldLogger 实例的指针, 用于链式调用。
func (fl *FieldLogger) Bool(key string, value bool) *FieldLogger {
	if fl != nil && fl.fields != nil && key != "" {
		fl.fields[key] = strconv.FormatBool(value)
	}
	return fl
}

// Uint64 添加64位无符号整数字段，支持链式调用
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为64位无符号整数。
//
// 返回值:
//   - *FieldLogger: 一个指向 FieldLogger 实例的指针, 用于链式调用。
func (fl *FieldLogger) Uint64(key string, value uint64) *FieldLogger {
	if fl != nil && fl.fields != nil && key != "" {
		fl.fields[key] = strconv.FormatUint(value, 10)
	}
	return fl
}

// Uint32 添加32位无符号整数字段，支持链式调用
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为32位无符号整数。
//
// 返回值:
//   - *FieldLogger: 一个指向 FieldLogger 实例的指针, 用于链式调用。
func (fl *FieldLogger) Uint32(key string, value uint32) *FieldLogger {
	if fl != nil && fl.fields != nil && key != "" {
		fl.fields[key] = strconv.FormatUint(uint64(value), 10)
	}
	return fl
}

// Uint16 添加16位无符号整数字段，支持链式调用
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为16位无符号整数。
//
// 返回值:
//   - *FieldLogger: 一个指向 FieldLogger 实例的指针, 用于链式调用。
func (fl *FieldLogger) Uint16(key string, value uint16) *FieldLogger {
	if fl != nil && fl.fields != nil && key != "" {
		fl.fields[key] = strconv.FormatUint(uint64(value), 10)
	}
	return fl
}

// Uint8 添加8位无符号整数字段，支持链式调用
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为8位无符号整数。
//
// 返回值:
//   - *FieldLogger: 一个指向 FieldLogger 实例的指针, 用于链式调用。
func (fl *FieldLogger) Uint8(key string, value uint8) *FieldLogger {
	if fl != nil && fl.fields != nil && key != "" {
		fl.fields[key] = strconv.FormatUint(uint64(value), 10)
	}
	return fl
}

// Duration 添加时间间隔字段，支持链式调用
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为时间间隔。
//
// 返回值:
//   - *FieldLogger: 一个指向 FieldLogger 实例的指针, 用于链式调用。
func (fl *FieldLogger) Duration(key string, value time.Duration) *FieldLogger {
	if fl != nil && fl.fields != nil && key != "" {
		fl.fields[key] = value.String()
	}
	return fl
}

// Error 添加错误字段，支持链式调用
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为错误。
//
// 返回值:
//   - *FieldLogger: 一个指向 FieldLogger 实例的指针, 用于链式调用。
func (fl *FieldLogger) Error(key string, value error) *FieldLogger {
	if fl != nil && fl.fields != nil && key != "" {
		fl.fields[key] = value.Error()
	}
	return fl
}
