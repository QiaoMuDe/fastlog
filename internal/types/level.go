package types

// LogLevel 定义为位掩码类型，每一位代表一个日志级别
type LogLevel uint8

// 定义日志级别的位掩码
const (
	DEBUG_Mask LogLevel = 1 << iota // 1 表示Debug级别
	INFO_Mask                       // 2  表示Info级别
	WARN_Mask                       // 4  表示Warn级别
	ERROR_Mask                      // 8  表示Error级别
	FATAL_Mask                      // 16 表示Fatal级别
	NONE_Mask  LogLevel = 0         // 0 表示不启用任何日志

	// 预定义的日志级别组合
	DEBUG LogLevel = DEBUG_Mask | INFO_Mask | WARN_Mask | ERROR_Mask | FATAL_Mask // Debug及以上级别
	INFO  LogLevel = INFO_Mask | WARN_Mask | ERROR_Mask | FATAL_Mask              // Info及以上级别
	WARN  LogLevel = WARN_Mask | ERROR_Mask | FATAL_Mask                          // Warn及以上级别
	ERROR LogLevel = ERROR_Mask | FATAL_Mask                                      // Error及以上级别
	FATAL LogLevel = FATAL_Mask                                                   // Fatal级别
	NONE  LogLevel = NONE_Mask
)

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

// ShouldLog 检查是否应该记录指定级别的日志
// 使用位运算优化日志级别比较，提高判断性能
//
// 参数：
//   - logLevel: 当前要记录的日志级别（基本级别，如DEBUG_Mask）
//   - minLevel: 配置的最低记录级别（组合级别，如DEBUG, INFO等）
//
// 返回值：
//   - bool: 如果当前级别应该记录，则返回 true；否则返回 false
func ShouldLog(logLevel, minLevel LogLevel) bool {
	// NONE级别表示不记录任何日志
	if minLevel == NONE {
		return false
	}

	// 使用位运算进行级别判断
	// 如果最低记录级别包含当前日志级别，则应该记录日志
	return minLevel&logLevel != 0
}

// String 将 LogLevel 转换为字符串
func (l LogLevel) String() string {
	// 使用预构建的映射表进行O(1)查询(不带填充，适用于JSON)
	if str, exists := LogLevelStringMap[l]; exists {
		return str
	}
	return "UNKNOWN"
}

// MarshalJSON 将 LogLevel 转换为 JSON 字符串
//
// 返回值：
//   - []byte: 包含日志级别的 JSON 字符串（带双引号）
//   - error: 如果转换过程中发生错误，返回非 nil 错误；否则返回 nil
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
