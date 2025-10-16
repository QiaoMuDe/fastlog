package fastlog

import (
	"strconv"
	"sync"
	"time"
)

// FieldType 定义了日志字段的类型(13个)
type FieldType uint8

const (
	StringType   FieldType = iota // 字符串字段
	IntType                       // 整数字段
	Int64Type                     // 64位整数字段
	Float64Type                   // 64位浮点数字段
	BoolType                      // 布尔字段
	TimeType                      // 时间字段
	DurationType                  // 间隔时间字段
	Uint64Type                    // 64位无符号整数字段
	Uint32Type                    // 32位无符号整数字段
	Uint16Type                    // 16位无符号整数字段
	Uint8Type                     // 8位无符号整数字段
	ErrorType                     // 错误字段
	UnknownType                   // 未知字段
)

// Field  定义了日志字段的结构体
type Field struct {
	key string    // 字段键名
	typ FieldType // 字段类型

	// union
	strVal      string        // 字符串值
	intVal      int           // 整数值
	int64Val    int64         // 64位整数值
	float64Val  float64       // 64位浮点数值
	boolVal     bool          // 布尔值
	timeVal     time.Time     // 时间值
	durationVal time.Duration // 间隔时间值
	uint64Val   uint64        // 64位无符号整数值
	uint32Val   uint32        // 32位无符号整数值
	uint16Val   uint16        // 16位无符号整数值
	uint8Val    uint8         // 8位无符号整数值
	errorVal    error         // 错误值
}

// Key 获取字段键名
//
// 返回值:
//   - string: 字段键名。
func (f *Field) Key() string {
	if f == nil {
		return ""
	}

	return f.key
}

// Type 获取字段类型
//
// 返回值:
//   - FieldType: 字段类型。
func (f *Field) Type() FieldType {
	if f == nil {
		return UnknownType
	}

	return f.typ
}

// Value 获取字段值
//
// 返回值:
//   - string: 字段值的字符串表示。
func (f *Field) Value() string {
	if f == nil {
		return ""
	}

	switch f.typ {
	case StringType:
		return f.strVal

	case IntType:
		return strconv.Itoa(f.intVal)

	case Int64Type:
		return strconv.FormatInt(f.int64Val, 10)

	case Float64Type:
		return strconv.FormatFloat(f.float64Val, 'f', -1, 64)

	case BoolType:
		return strconv.FormatBool(f.boolVal)

	case TimeType:
		return f.timeVal.Format(time.DateTime)

	case DurationType:
		return f.durationVal.String()

	case Uint64Type:
		return strconv.FormatUint(f.uint64Val, 10)

	case Uint32Type:
		return strconv.FormatUint(uint64(f.uint32Val), 10)

	case Uint16Type:
		return strconv.FormatUint(uint64(f.uint16Val), 10)

	case Uint8Type:
		return strconv.FormatUint(uint64(f.uint8Val), 10)

	case ErrorType:
		if f.errorVal != nil {
			return f.errorVal.Error()
		}
		return ""

	default:
		return ""
	}
}

// fieldPool Field对象池，用于重用Field实例
var fieldPool = sync.Pool{
	New: func() interface{} {
		return &Field{}
	},
}

// getField 从对象池获取Field实例
func getField() *Field {
	return fieldPool.Get().(*Field)
}

// putField 将Field实例归还到对象池
func putField(f *Field) {
	if f == nil {
		return
	}

	// 重置字段值
	*f = Field{}
	fieldPool.Put(f)
}

// String 添加字符串字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为字符串。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
func String(key string, value string) *Field {
	f := getField()
	f.key = key
	f.typ = StringType
	f.strVal = value
	return f
}

// Int 添加整数字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为整数。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
func Int(key string, value int) *Field {
	f := getField()
	f.key = key
	f.typ = IntType
	f.intVal = value
	return f
}

// Int64 添加64位整数字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为64位整数。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
func Int64(key string, value int64) *Field {
	f := getField()
	f.key = key
	f.typ = Int64Type
	f.int64Val = value
	return f
}

// Float64 添加64位浮点数字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为64位浮点数。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
func Float64(key string, value float64) *Field {
	f := getField()
	f.key = key
	f.typ = Float64Type
	f.float64Val = value
	return f
}

// Bool 添加布尔字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为布尔值。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
func Bool(key string, value bool) *Field {
	f := getField()
	f.key = key
	f.typ = BoolType
	f.boolVal = value
	return f
}

// Time 添加时间字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为time.Time类型。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
func Time(key string, value time.Time) *Field {
	f := getField()
	f.key = key
	f.typ = TimeType
	f.timeVal = value
	return f
}

// Duration 添加持续时间字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为time.Duration类型。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
func Duration(key string, value time.Duration) *Field {
	f := getField()
	f.key = key
	f.typ = DurationType
	f.durationVal = value
	return f
}

// Uint64 添加64位无符号整数字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为64位无符号整数。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
func Uint64(key string, value uint64) *Field {
	f := getField()
	f.key = key
	f.typ = Uint64Type
	f.uint64Val = value
	return f
}

// Uint32 添加32位无符号整数字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为32位无符号整数。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
func Uint32(key string, value uint32) *Field {
	f := getField()
	f.key = key
	f.typ = Uint32Type
	f.uint32Val = value
	return f
}

// Uint16 添加16位无符号整数字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为16位无符号整数。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
func Uint16(key string, value uint16) *Field {
	f := getField()
	f.key = key
	f.typ = Uint16Type
	f.uint16Val = value
	return f
}

// Uint8 添加8位无符号整数字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为8位无符号整数。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
func Uint8(key string, value uint8) *Field {
	f := getField()
	f.key = key
	f.typ = Uint8Type
	f.uint8Val = value
	return f
}

// Error 添加错误字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为error类型。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
func Error(key string, value error) *Field {
	f := getField()
	f.key = key
	f.typ = ErrorType
	f.errorVal = value
	return f
}
