package fastlog

import (
	"runtime/debug"
	"strconv"
	"time"
)

// FieldType 表示字段类型
type FieldType int8

// 字段类型常量
const (
	UnknownType  FieldType = iota // 未知类型
	StringType                    // 字符串类型
	IntType                       // 整数类型
	Int64Type                     // 64位整数类型
	UintType                      // 无符号整数类型
	Uint64Type                    // 64位无符号整数类型
	Float64Type                   // 浮点数类型
	BoolType                      // 布尔类型
	TimeType                      // 时间类型
	DurationType                  // 时间持续类型
	ErrorType                     // 错误类型
	AnyType                       // any类型
)

// Field 表示一个键值对字段, 包含所有可能的类型
type Field struct {
	key       string        // 字段键
	typ       FieldType     // 字段类型
	stringVal string        // 字符串值
	intVal    int64         // 整数值
	uintVal   uint64        // 无符号整数值
	floatVal  float64       // 浮点数值
	boolVal   bool          // 布尔值
	timeVal   time.Time     // 时间值
	duration  time.Duration // 时间持续值
	iface     interface{}   // any类型值
}

// Key 返回字段键
//
// 返回:
//   - string: 字段键名
func (f Field) Key() string {
	return f.key
}

// Type 返回字段类型
//
// 返回:
//   - FieldType: 字段类型
func (f Field) Type() FieldType {
	return f.typ
}

// Value 将字段值转换为字符串返回
//
// 根据字段类型将值格式化为字符串:
//   - StringType/ErrorType: 直接返回字符串值
//   - IntType/Int64Type: 转为 10 进制字符串
//   - UintType/Uint64Type: 转为 10 进制无符号字符串
//   - Float64Type: 转为浮点数字符串
//   - BoolType: 转为 "true" 或 "false"
//   - TimeType: 转为 RFC3339 格式时间字符串
//   - DurationType: 转为持续时间字符串 (如 "1h30m")
//   - AnyType: 使用 fmt.Sprintf("%v") 格式化
//   - 其他类型: 返回空字符串
//
// 返回:
//   - string: 字段值的字符串表示
func (f Field) Value() string {
	switch f.typ {
	case StringType, ErrorType:
		return f.stringVal

	case IntType, Int64Type:
		return strconv.FormatInt(f.intVal, 10)

	case UintType, Uint64Type:
		return strconv.FormatUint(f.uintVal, 10)

	case Float64Type:
		return strconv.FormatFloat(f.floatVal, 'g', -1, 64)

	case BoolType:
		return strconv.FormatBool(f.boolVal)

	case TimeType:
		return f.timeVal.Format(DefaultTimeFormat)

	case DurationType:
		return f.duration.String()

	case AnyType:
		return f.anyString()

	default:
		return ""
	}
}

// anyString 将 iface 字段值转为字符串
//
// 返回:
//   - string: 任意类型值的字符串表示
func (f Field) anyString() string {
	if f.iface == nil {
		return "null"
	}
	switch val := f.iface.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)

	case int8:
		return strconv.FormatInt(int64(val), 10)

	case int16:
		return strconv.FormatInt(int64(val), 10)

	case int32:
		return strconv.FormatInt(int64(val), 10)

	case int64:
		return strconv.FormatInt(val, 10)

	case uint:
		return strconv.FormatUint(uint64(val), 10)

	case uint8:
		return strconv.FormatUint(uint64(val), 10)

	case uint16:
		return strconv.FormatUint(uint64(val), 10)

	case uint32:
		return strconv.FormatUint(uint64(val), 10)

	case uint64:
		return strconv.FormatUint(val, 10)

	case float32:
		return strconv.FormatFloat(float64(val), 'g', -1, 32)

	case float64:
		return strconv.FormatFloat(val, 'g', -1, 64)

	case bool:
		return strconv.FormatBool(val)

	case time.Time:
		return val.Format(DefaultTimeFormat)

	case time.Duration:
		return val.String()

	case error:
		return val.Error()

	default:
		return ""
	}
}

// toInterface 将字段值转换为 interface{}
// 用于 JSON 格式化器
func (f Field) toInterface() interface{} {
	switch f.typ {
	case StringType, ErrorType:
		return f.stringVal

	case IntType, Int64Type:
		return f.intVal

	case UintType, Uint64Type:
		return f.uintVal

	case Float64Type:
		return f.floatVal

	case BoolType:
		return f.boolVal

	case TimeType:
		return f.timeVal.Format(DefaultTimeFormat)

	case DurationType:
		return f.duration.String()

	case AnyType:
		return f.iface

	default:
		return nil
	}
}

// Format 将字段格式化为 key=value 形式
//
// 返回:
//   - string: 格式化后的字段字符串，格式为 "key=value"
//
// 示例:
//
//	field := fastlog.String("user", "admin")
//	result := field.Format() // 返回 "user=admin"
func (f Field) Format() string {
	return f.key + "=" + f.Value()
}

// Stack 创建一个堆栈字段
//
// 返回:
//   - Field: 字段实例, 键名为 "stack", 值为当前堆栈信息
func Stack() Field {
	return Field{key: "stack", typ: StringType, stringVal: string(debug.Stack())}
}

// String 创建一个字符串字段
//
// 参数:
//   - key: 字段键
//   - val: 字符串值
//
// 返回:
//   - Field: 字段实例
func String(key, val string) Field {
	return Field{key: key, typ: StringType, stringVal: val}
}

// Int 创建一个 int 字段
//
// 参数:
//   - key: 字段键
//   - val: 整数值
//
// 返回:
//   - Field: 字段实例
func Int(key string, val int) Field {
	return Field{key: key, typ: IntType, intVal: int64(val)}
}

// Int64 创建一个 int64 字段
//
// 参数:
//   - key: 字段键
//   - val: 64位整数值
//
// 返回:
//   - Field: 字段实例
func Int64(key string, val int64) Field {
	return Field{key: key, typ: Int64Type, intVal: val}
}

// Uint 创建一个 uint 字段
//
// 参数:
//   - key: 字段键
//   - val: 无符号整数值
//
// 返回:
//   - Field: 字段实例
func Uint(key string, val uint) Field {
	return Field{key: key, typ: UintType, uintVal: uint64(val)}
}

// Uint64 创建一个 uint64 字段
//
// 参数:
//   - key: 字段键
//   - val: 64位无符号整数值
//
// 返回:
//   - Field: 字段实例
func Uint64(key string, val uint64) Field {
	return Field{key: key, typ: Uint64Type, uintVal: val}
}

// Float64 创建一个 float64 字段
//
// 参数:
//   - key: 字段键
//   - val: 浮点数值
//
// 返回:
//   - Field: 字段实例
func Float64(key string, val float64) Field {
	return Field{key: key, typ: Float64Type, floatVal: val}
}

// Bool 创建一个 bool 字段
//
// 参数:
//   - key: 字段键
//   - val: 布尔值
//
// 返回:
//   - Field: 字段实例
func Bool(key string, val bool) Field {
	return Field{key: key, typ: BoolType, boolVal: val}
}

// Time 创建一个 time.Time 字段
//
// 参数:
//   - key: 字段键
//   - val: 时间值
//
// 返回:
//   - Field: 字段实例
func Time(key string, val time.Time) Field {
	return Field{key: key, typ: TimeType, timeVal: val}
}

// Duration 创建一个 time.Duration 字段
//
// 参数:
//   - key: 字段键
//   - val: 时间持续值
//
// 返回:
//   - Field: 字段实例
func Duration(key string, val time.Duration) Field {
	return Field{key: key, typ: DurationType, duration: val}
}

// Error 创建一个 error 字段
//
// 参数:
//   - err: 错误值, nil 时返回 "<nil>" 字符串
//
// 返回:
//   - Field: 字段实例, 键名为 "error"
func Error(err error) Field {
	if err == nil {
		return Field{key: "error", typ: StringType, stringVal: "<nil>"}
	}
	return Field{key: "error", typ: ErrorType, stringVal: err.Error()}
}

// Err 创建一个自定义键名的 error 字段
//
// 参数:
//   - key: 字段键名
//   - err: 错误值, nil 时返回 "<nil>" 字符串
//
// 返回:
//   - Field: 字段实例
//
// 示例:
//
//	logger.Infow("操作失败",
//	    fastlog.Err("db_error", dbErr),
//	    fastlog.Err("cache_error", cacheErr),
//	)
func Err(key string, err error) Field {
	if err == nil {
		return Field{key: key, typ: StringType, stringVal: "<nil>"}
	}
	return Field{key: key, typ: ErrorType, stringVal: err.Error()}
}

// Any 创建一个任意类型字段
//
// 参数:
//   - key: 字段键
//   - val: 任意类型的值
//
// 返回:
//   - Field: 字段实例
func Any(key string, val interface{}) Field {
	return Field{key: key, typ: AnyType, iface: val}
}
