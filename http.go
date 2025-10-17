package fastlog

import (
	"net/http"
	"sync"
	"time"
)

// logWriterPool 是一个用于 logWriter 实例的对象池
var logWriterPool = sync.Pool{
	New: func() interface{} {
		return &logWriter{}
	},
}

// getLogWriter 从池中获取一个 logWriter 实例
func getLogWriter(w http.ResponseWriter) *logWriter {
	lw := logWriterPool.Get().(*logWriter)
	lw.ResponseWriter = w
	lw.statusCode = http.StatusOK
	return lw
}

// putLogWriter 将一个 logWriter 实例归还到池中
func putLogWriter(lw *logWriter) {
	// 清理引用，防止内存泄漏
	lw.ResponseWriter = nil
	logWriterPool.Put(lw)
}

// logWriter 是一个包装器，用于捕获http.ResponseWriter写入的状态码
type logWriter struct {
	http.ResponseWriter     // 原始的 ResponseWriter
	statusCode          int // 捕获的状态码
}

// WriteHeader 捕获状态码并调用原始的 WriteHeader
func (lw *logWriter) WriteHeader(code int) {
	lw.statusCode = code
	lw.ResponseWriter.WriteHeader(code)
}

// LogRequest 日志中间件，用于记录HTTP请求日志
//
// 参数:
//   - log: 日志实例
//   - next: 下一个处理器
//
// 返回:
//   - http.Handler: 中间件处理后的处理器
func LogRequest(log *FLog, next http.Handler) http.Handler {
	if log == nil || next == nil {
		log.Fatal("fastlog: invalid arguments")
	}

	// 用HandlerFunc适配匿名函数，满足Handler接口
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 日志前置操作：记录请求开始时间
		startTime := time.Now()

		// 从对象池获取 logWriter 实例
		lw := getLogWriter(w)
		defer putLogWriter(lw) // 确保请求处理完毕后将实例归还池中

		// 调用原处理器（核心：执行业务逻辑）
		next.ServeHTTP(lw, r)

		// 日志后置操作：计算耗时并打印日志
		duration := time.Since(startTime)

		// 打印HTTP日志
		log.InfoFields("[HTTP LOG]",
			String("method", r.Method),            // 请求方法
			String("path", r.URL.Path),            // 请求路径
			Int("status", lw.statusCode),          // HTTP状态码
			Duration("duration", duration),        // 处理耗时
			String("remote_addr", r.RemoteAddr),   // 客户端IP
			String("user_agent", r.UserAgent()),   // User-Agent
			Int64("content_len", r.ContentLength), // 请求体大小
		)
	})
}
