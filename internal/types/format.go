package types

// LogFormatType 日志格式选项
//
// 格式:
//   - Detailed: 详细格式
//   - Json: json格式
//   - Structured: 结构化格式
//   - Timestamp: 时间格式
//   - Custom: 自定义格式
type LogFormatType int

// String 将 LogFormatType 转换为对应的字符串
func (f LogFormatType) String() string {
	switch f {
	case Json:
		return "json"
	case Structured:
		return "structured"
	case Timestamp:
		return "timestamp"
	case Custom:
		return "custom"
	default:
		return "def"
	}
}

// 日志格式选项
const (
	Def        LogFormatType = iota // 默认格式
	Json                            // json格式
	Structured                      // 结构化格式
	Timestamp                       // 时间格式
	Custom                          // 自定义格式
)

// 仅供内部参考的日志格式
// // 定义日志格式
// var LogFormatMap = map[LogFormatType]string{
// 	Json:            `{"time":"%s","level":"%s","caller":"%s","message":"%s"}`, // Json格式
// 	JsonSimple:      `{"time":"%s","level":"%s","message":"%s"}`,               // Json简化格式(无文件信息)
// 	Detailed:        `%s | %-6s | %s - %s`,                                     // 详细格式
// 	Simple:          `%s | %-6s | %s`,                                          // 简约格式(无文件信息)                                                                                                 // 自定义格式
// 	Structured:      `T:%s|L:%-6s|C:%s|M:%s`,                                   // 结构化格式
// 	BasicStructured: `T:%s|L:%-6s|M:%s`,                                        // 基础结构化格式(无文件信息)
// 	Timestamp:       `%s %s %s`,                                                // 简单时间格式
// }
