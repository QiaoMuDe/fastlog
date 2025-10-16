package fastlog

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gitee.com/MM-Q/fastlog/internal/types"
)

// TestWithLog 测试 WithLog 中间件是否能正确记录 HTTP 请求日志
func TestWithLog(t *testing.T) {
	// 设置日志记录器，使用 ConsoleConfig 并将格式设置为 KVFmt 以便断言
	cfg := ConsoleConfig()
	cfg.LogFormat = types.Def // 使用键值对格式，方便检查
	cfg.Color = true          // 关闭颜色，避免 ANSI 转义字符干扰
	cfg.CallerInfo = false    // 关闭调用者信息
	log := NewFLog(cfg)

	// 创建模拟的 HTTP 处理器和请求
	// 模拟业务逻辑处理器
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 模拟业务耗时
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// 使用中间件包装处理器
	middleware := LogRequest(log, nextHandler)

	// 创建一个模拟请求
	req := httptest.NewRequest("GET", "/test/path", nil)
	rr := httptest.NewRecorder()

	// 执行中间件
	middleware.ServeHTTP(rr, req)
	_ = log.Close() // 关闭日志，确保缓冲区内容被刷新
}
