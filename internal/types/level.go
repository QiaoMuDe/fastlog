package types

// LogLevel 定义为位掩码类型，每一位代表一个日志级别
type LogLevel uint8

// 定义日志级别的位掩码
const (
	DEBUG LogLevel = 1 << iota // 1
	INFO                       // 2
	WARN                       // 4
	ERROR                      // 8
	FATAL                      // 16
	NONE   LogLevel = 0        // 0 表示不启用任何日志
	ALL    LogLevel = DEBUG | INFO | WARN | ERROR | FATAL // 31 表示启用所有日志
)

// 预构建的日志级别到字符串的映射表（不带填充，用于JSON序列化）
var LogLevelStringMap = map[LogLevel]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
	NONE:  "NONE",
	ALL:   "ALL",
}

// 预构建的日志级别到字符串的映射表（带填充，用于文本格式化）
var LogLevelPaddedStringMap = map[LogLevel]string{
	DEBUG: "DEBUG", // 5个字符(预填充空格)
	INFO:  "INFO ", // 5个字符
	WARN:  "WARN ", // 5个字符
	ERROR: "ERROR", // 5个字符
	FATAL: "FATAL", // 5个字符
	NONE:  "NONE ", // 5个字符
	ALL:   "ALL  ", // 5个字符
}

// ShouldLog 检查是否应该记录指定级别的日志
// 使用位运算优化日志级别比较，提高判断性能
func ShouldLog(currentLevel, configLevel LogLevel) bool {
	// NONE级别表示不记录任何日志
	if configLevel == NONE {
		return false
	}
	
	// 使用位运算进行级别判断
	// 如果配置级别包含当前级别，则应该记录日志
	return configLevel&currentLevel != 0
}

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
