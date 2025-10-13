package types

// SimpleLogMsg 简化的日志消息结构（用于JsonSimple格式）
type SimpleLogMsg struct {
	Timestamp string   `json:"time"`    // 预格式化的时间字符串
	Level     LogLevel `json:"level"`   // 日志级别
	Message   string   `json:"message"` // 日志消息
}

// LogFormatType 日志格式选项
//
// 格式:
//   - Detailed: 详细格式
//   - Json: json格式
//   - JsonSimple: json简化格式(无文件信息)
//   - Simple: 简约格式(无文件信息)
//   - Structured: 结构化格式
//   - BasicStructured: 基础结构化格式(无文件信息)
//   - SimpleTimestamp: 简单时间格式(无文件信息)
//   - Custom: 自定义格式(无文件信息)
type LogFormatType int

// 日志格式选项
const (
	Detailed        LogFormatType = iota // 详细格式
	Json                                 // json格式
	JsonSimple                           // json简化格式(无文件信息)
	Simple                               // 简约格式(无文件信息)
	Structured                           // 结构化格式
	BasicStructured                      // 基础结构化格式(无文件信息)
	SimpleTimestamp                      // 简单时间格式(无文件信息)
	Custom                               // 自定义格式(无文件信息)
)

// FileInfoRequiredFormats 需要处理文件信息的日志格式集合
// 如果日志格式需要文件信息(文件名、函数名、行号)，则将其添加到该集合中
var FileInfoRequiredFormats = map[LogFormatType]struct{}{
	Json:       {},
	Detailed:   {},
	Structured: {},
}

// 定义日志格式
var LogFormatMap = map[LogFormatType]string{
	Json:            `{"time":"%s","level":"%s","file":"%s","function":"%s","line":"%d","message":"%s"}`, // Json格式
	JsonSimple:      `{"time":"%s","level":"%s","message":"%s"}`,                                         // Json简化格式(无文件信息)
	Detailed:        `%s | %-6s | %s:%s:%d - %s`,                                                         // 详细格式
	Simple:          `%s | %-6s | %s`,                                                                    // 简约格式                                                                                                // 自定义格式
	Structured:      `T:%s|L:%-6s|F:%s:%s:%d|M:%s`,                                                       // 结构化格式
	BasicStructured: `T:%s|L:%-6s|M:%s`,                                                                  // 基础结构化格式(无文件信息)
	SimpleTimestamp: `%s %s %s`,                                                                          // 简单时间格式
}
