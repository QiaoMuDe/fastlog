package types

import "sync"

// LogMsg 结构体用于封装日志消息
type LogMsg struct {
	Timestamp string   `json:"time"`     // 预格式化的时间字符串
	Level     LogLevel `json:"level"`    // 日志级别
	FileName  string   `json:"file"`     // 文件名
	FuncName  string   `json:"function"` // 调用函数名
	Line      uint16   `json:"line"`     // 行号
	Message   string   `json:"message"`  // 日志消息
}

// LogMsgPool 是一个日志消息对象池
var LogMsgPool = sync.Pool{
	New: func() interface{} {
		return &LogMsg{}
	},
}

// getLogMsg 获取日志消息对象，使用安全的类型断言
//
// 返回：
//   - *logMsg: 日志消息对象指针，保证非nil
//   - 注意：返回的对象总是可以安全地传递给putLogMsg
func GetLogMsg() *LogMsg {
	// 尝试从对象池获取对象并进行类型断言
	if msg, ok := LogMsgPool.Get().(*LogMsg); ok {
		return msg
	}

	// 创建新的对象
	return &LogMsg{}
}

// putLogMsg 归还日志消息对象
//
// 参数：
//   - msg: 日志消息对象指针
//   - 注意：该函数可以安全地处理任何来源的logMsg对象，
//     包括从getLogMsg获取的对象和通过new/&logMsg{}创建的对象
func PutLogMsg(msg *LogMsg) {
	// 安全检查：防止空指针
	if msg == nil {
		return
	}

	// 使用零值重置，确保完全清理所有字段
	// 这种方式比逐个字段清理更安全，不会遗漏任何字段
	*msg = LogMsg{}

	// 归还对象到池中
	LogMsgPool.Put(msg)
}
