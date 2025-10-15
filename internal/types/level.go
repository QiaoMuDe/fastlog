package types

// 预构建的日志级别到字符串的映射表（不带填充，用于JSON序列化）
var LogLevelStringMap = map[LogLevel]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
	NONE:  "NONE",
}

// 预构建的日志级别到字符串的映射表（带填充，用于文本格式化）
var LogLevelPaddedStringMap = map[LogLevel]string{
	DEBUG: "DEBUG", // 5个字符(预填充空格)
	INFO:  "INFO ", // 5个字符
	WARN:  "WARN ", // 5个字符
	ERROR: "ERROR", // 5个字符
	FATAL: "FATAL", // 5个字符
	NONE:  "NONE ", // 5个字符
}

// 日志级别枚举
type LogLevel uint8

func (l LogLevel) String() string {
	// 使用预构建的映射表进行O(1)查询(不带填充，适用于JSON)
	if str, exists := LogLevelStringMap[l]; exists {
		return str
	}
	return "UNKNOWN"
}

// 将日志级别转换为字符串
func (l LogLevel) MarshalJSON() ([]byte, error) {
	return []byte(`"` + l.String() + `"`), nil
}

// 定义日志级别
const (
	DEBUG LogLevel = 10  // 调试级别
	INFO  LogLevel = 20  // 信息级别
	WARN  LogLevel = 30  // 警告级别
	ERROR LogLevel = 40  // 错误级别
	FATAL LogLevel = 50  // 致命级别
	NONE  LogLevel = 255 // 无日志级别
)

// LogLevelToPaddedString 将 LogLevel 转换为带填充的字符串（用于文本格式化）
//
// 参数：
//   - level: 要转换的日志级别
//
// 返回值：
//   - string: 对应的带填充的日志级别字符串（7个字符），如果 level 无效, 则返回 "UNKNOWN"
func LogLevelToPaddedString(level LogLevel) string {
	// 使用预构建的带填充映射表进行O(1)查询（适用于文本格式）
	if str, exists := LogLevelPaddedStringMap[level]; exists {
		return str
	}
	return "UNKNOWN"
}
