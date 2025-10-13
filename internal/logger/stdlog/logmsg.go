package stdlog

import (
	"sync"

	"gitee.com/MM-Q/fastlog/internal/types"
)

// logMsg 结构体用于封装日志消息
type logMsg struct {
	Timestamp string         `json:"time"`     // 预格式化的时间字符串
	Level     types.LogLevel `json:"level"`    // 日志级别
	FileName  string         `json:"file"`     // 文件名
	FuncName  string         `json:"function"` // 调用函数名
	Line      uint16         `json:"line"`     // 行号
	Message   string         `json:"message"`  // 日志消息
}

// logMsgPool 是一个日志消息对象池
var logMsgPool = sync.Pool{
	New: func() interface{} {
		return &logMsg{}
	},
}

// getLogMsg 获取日志消息对象，使用安全的类型断言
//
// 返回：
//   - *logMsg: 日志消息对象指针，保证非nil
//   - 注意：返回的对象总是可以安全地传递给putLogMsg
func getLogMsg() *logMsg {
	// 尝试从对象池获取对象并进行类型断言
	if msg, ok := logMsgPool.Get().(*logMsg); ok {
		return msg
	}

	// 创建新的对象
	return &logMsg{}
}

// putLogMsg 归还日志消息对象
//
// 参数：
//   - msg: 日志消息对象指针
//   - 注意：该函数可以安全地处理任何来源的logMsg对象，
//     包括从getLogMsg获取的对象和通过new/&logMsg{}创建的对象
func putLogMsg(msg *logMsg) {
	// 安全检查：防止空指针
	if msg == nil {
		return
	}

	// 使用零值重置，确保完全清理所有字段
	// 这种方式比逐个字段清理更安全，不会遗漏任何字段
	*msg = logMsg{}

	// 归还对象到池中
	logMsgPool.Put(msg)
}
