package types

// LogFormatType 日志格式选项
//
// 格式:
//   - Detailed: 详细格式
//   - Json: json格式
//   - Timestamp: 时间格式
//   - KVfmt: 键值对格式
//   - Custom: 自定义格式
type LogFormatType int

// String 将 LogFormatType 转换为对应的字符串
func (f LogFormatType) String() string {
	switch f {
	case Json:
		return "json"

	case Timestamp:
		return "timestamp"

	case Custom:
		return "custom"

	case KVfmt:
		return "kvfmt"

	case LogFmt:
		return "logfmt"

	default:
		return "def"
	}
}

// 日志格式选项
const (
	Def       LogFormatType = iota // 默认格式
	Json                           // json格式
	Timestamp                      // 时间格式
	KVfmt                          // KVfmt 键值对格式
	LogFmt                         // logfmt格式
	Custom                         // 自定义格式
)
