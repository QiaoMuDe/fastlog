package flog

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
	Key   string
	Type  FieldType
	Value any
}

// String 添加字符串字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为字符串。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针, 用于链式调用。
func (f *Field) String(key string, value string) *Field {
	if f != nil && key != "" {
		f.Key = key
		f.Type = StringType
		f.Value = value
	}
	return f
}

// Int 添加整数字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为整数。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针, 用于链式调用。
func (f *Field) Int(key string, value int) *Field {
	if f != nil && key != "" {
		f.Key = key
		f.Type = IntType
		f.Value = value
	}
	return f
}
