package fastlog

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// PathInfo 是一个结构体，用于封装路径的信息
type PathInfo struct {
	Path    string      // 路径
	Exists  bool        // 是否存在
	IsFile  bool        // 是否为文件
	IsDir   bool        // 是否为目录
	Size    int64       // 文件大小（字节）
	Mode    os.FileMode // 文件权限
	ModTime time.Time   // 文件修改时间
}

// 定义正则表达式，用于匹配日志级别
var re = regexp.MustCompile(`\b(INFO|WARN|ERROR|SUCCESS|DEBUG)\b`)

// checkPath 检查给定路径的信息
func checkPath(path string) (PathInfo, error) {
	// 创建一个 PathInfo 结构体
	var info PathInfo

	// 清理路径，确保没有多余的斜杠
	path = filepath.Clean(path)

	// 设置路径
	info.Path = path

	// 使用 os.Stat 获取文件状态
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果路径不存在, 则直接返回
			info.Exists = false
			return info, fmt.Errorf("路径 '%s' 不存在，请检查路径是否正确: %s", path, err)
		} else {
			return info, fmt.Errorf("无法访问路径 '%s': %s", path, err)
		}
	}

	// 路径存在，填充信息
	info.Exists = true                // 标记路径存在
	info.IsFile = !fileInfo.IsDir()   // 通过取反判断是否为文件，因为 IsDir 返回 false 表示是文件
	info.IsDir = fileInfo.IsDir()     // 直接使用 IsDir 方法判断是否为目录
	info.Size = fileInfo.Size()       // 获取文件大小
	info.Mode = fileInfo.Mode()       // 获取文件权限
	info.ModTime = fileInfo.ModTime() // 获取文件的最后修改时间

	// 返回路径信息结构体
	return info, nil
}

// getCallerInfo 获取调用者的信息
// 参数：
// skip - 跳过的调用层数（通常设置为1或2，具体取决于调用链的深度）
// 返回值：
// fileName - 调用者的文件名（不包含路径）
// functionName - 调用者的函数名
// line - 调用者的行号
// ok - 是否成功获取到调用者信息
func getCallerInfo(skip int) (fileName string, functionName string, line int, ok bool) {
	// 获取调用者信息，跳过指定的调用层数
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "???"
		line = 0
		functionName = "???"
		return
	}

	// 获取文件名（只保留文件名，不包含路径）
	fileName = filepath.Base(file)

	// 获取函数名
	function := runtime.FuncForPC(pc)
	if function != nil {
		functionName = function.Name()
	} else {
		functionName = "???"
	}

	return
}

// 获取当前 Goroutine 的 ID
func getGoroutineID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := bytes.Fields(buf[:n])[1]
	id, _ := strconv.ParseInt(string(idField), 10, 64)
	return id
}

// logLevelToString 将 LogLevel 转换为对应的字符串，并以大写形式返回
// 参数：
// level - 要转换的日志级别
// 返回值：
// string - 对应的日志级别字符串，如果 level 无效，则返回 "UNKNOWN"
func logLevelToString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case SUCCESS:
		return "SUCCESS"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case None:
		return "NONE"
	default:
		return "UNKNOWN"
	}
}

// addColor 根据日志级别添加颜色
func addColor(s string) string {
	// 使用正则表达式精确匹配日志级别，确保匹配的是独立的单词
	match := re.FindString(s)

	// 根据匹配到的日志级别添加颜色
	switch match {
	case "INFO":
		return CL.Sblue(s) // Blue
	case "WARN":
		return CL.Syellow(s) // Yellow
	case "ERROR":
		return CL.Sred(s) // Red
	case "SUCCESS":
		return CL.Sgreen(s) // Green
	case "DEBUG":
		return CL.Spurple(s) // Purple
	default:
		return s // 如果没有匹配到日志级别，返回原始字符串
	}
}

// formatLog 格式化日志消息。
func formatLog(f *FastLog, l *logMessage) string {
	// 定义一个变量，用于存储格式化后的日志消息。
	var logMsg string
	switch f.logFormat {
	// Json格式
	case Json:
		logMsg = fmt.Sprintf(
			`{"time":"%s","level":"%s","file":"%s","function":"%s","line":"%d", "thread":"%d","message":"%s"}`,
			l.timestamp.Format("2006-01-02 15:04:05"), logLevelToString(l.level), l.fileName, l.funcName, l.line, l.goroutineID, l.message,
		)
	// 详细格式
	case Detailed:
		// 按照指定格式输出日志，使用%-7s让日志级别左对齐且宽度为7个字符
		logMsg = fmt.Sprintf(
			"%s | %-7s | %s:%s:%d - %s",
			l.timestamp.Format("2006-01-02 15:04:05"), logLevelToString(l.level), l.fileName, l.funcName, l.line, l.message,
		)
	// 括号格式
	case Bracket:
		logMsg = fmt.Sprintf("[%s] %s", logLevelToString(l.level), l.message)
	// 协程格式
	case Threaded:
		logMsg = fmt.Sprintf(`%s | %-7s | [thread="%d"] %s`, l.timestamp.Format("2006-01-02 15:04:05"), logLevelToString(l.level), l.goroutineID, l.message)
	// 无法识别的日志格式选项
	default:
		logMsg = fmt.Sprintf("无法识别的日志格式选项: %v", f.logFormat)
	}

	// 返回格式化后的日志消息。
	return logMsg
}

// 发送缓冲区内容到文件和控制台
func sendBufferToFile(f *FastLog, logMsg *logMessage) error {
	// 加锁，确保并发安全
	f.sendMu.Lock()
	defer f.sendMu.Unlock()

	// 计算缓冲区大小
	BufferSize := f.bufferKbSize * 1024

	// 先检查文件缓冲区大小，如果达到阈值则立即进行写入
	if f.fileBuffer != nil && f.fileBuffer.Len() >= BufferSize {
		if !f.consoleOnly {
			if err := logBuffersWriter(f); err != nil {
				return fmt.Errorf("写入文件缓冲区内容到文件失败: %v", err)
			}
		}
	}

	// 初始化控制台字符串构建器
	var consoleBuilder strings.Builder

	// 初始化日志文件字符串构建器
	var fileBuilder strings.Builder

	// 格式化日志消息
	formattedLog := formatLog(f, logMsg)

	// 如果允许将日志输出到控制台，或者仅输出到控制台
	if f.consoleOnly || f.printToConsole {
		// 检查控制台缓冲区剩余空间
		remainingConsoleSpace := BufferSize - f.consoleBuffer.Len()
		if len(formattedLog) > remainingConsoleSpace {
			// 如果日志消息长度超过了控制台缓冲区剩余空间，触发写入操作
			if err := logBuffersWriter(f); err != nil {
				return fmt.Errorf("写入控制台缓冲区内容到控制台失败: %v", err)
			}
		}

		// 调用addColor方法给日志消息添加颜色
		formattedLog = addColor(formattedLog)

		// 给格式化后的日志消息添加换行符
		consoleBuilder.WriteString(formattedLog)
		consoleBuilder.WriteString("\n")
		formattedLog = consoleBuilder.String()

		// 将格式化后的日志消息写入到控制台缓冲区中
		if _, err := f.consoleBuffer.WriteString(formattedLog); err != nil {
			return fmt.Errorf("写入控制台缓冲区失败: %v", err)
		}
	}

	// 如果不是仅输出到控制台，则将日志消息写入到日志文件缓冲区中。
	if !f.consoleOnly {
		// 检查缓冲区剩余空间
		remainingSpace := BufferSize - f.fileBuffer.Len()
		if len(formattedLog) > remainingSpace {
			// 如果日志消息长度超过了缓冲区剩余空间，触发写入操作
			if err := logBuffersWriter(f); err != nil {
				return fmt.Errorf("写入文件缓冲区内容到文件失败: %v", err)
			}
		}

		// 给格式化后的日志消息添加换行符
		fileBuilder.WriteString(formattedLog)
		fileBuilder.WriteString("\n")
		formattedLog = fileBuilder.String()

		// 将格式化后的日志消息写入到日志文件缓冲区中
		if _, err := f.fileBuffer.WriteString(formattedLog); err != nil {
			return fmt.Errorf("写入文件缓冲区失败: %v", err)
		}
	}

	return nil
}
