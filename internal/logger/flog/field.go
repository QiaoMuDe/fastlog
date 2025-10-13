package flog

import "strconv"

// FieldType 定义了日志字段的类型
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
)

// Field  定义了日志字段的结构体
type Field struct {
	key       string    // 字段键名
	fieldType FieldType // 字段类型
	value     string    // 字符串值（默认存储格式，用于快速输出）
}

// 访问方法
func (f *Field) Key() string {
	return f.key
}

func (f *Field) Type() FieldType {
	return f.fieldType
}

func (f *Field) Value() string {
	return f.value
}

// String 添加字符串字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为字符串。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针, 用于链式调用。
func String(key string, value string) *Field {
	if key != "" {
		return &Field{
			key:       key,
			fieldType: StringType,
			value:     value,
		}
	}
	return nil
}

// Int 添加整数字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为整数。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针, 用于链式调用。
func Int(key string, value int) *Field {
	if key != "" {
		return &Field{
			key:       key,
			fieldType: IntType,
			value:     strconv.Itoa(value),
		}
	}
	return nil
}
