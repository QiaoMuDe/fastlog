package fastlog

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gitee.com/MM-Q/fastlog/internal/types"
)

// TestWithLog 打印多组模拟 HTTP 请求的日志（无断言）
func TestWithLog(t *testing.T) {
	// 配置日志
	cfg := ConsoleConfig()
	cfg.LogFormat = types.Def // 使用默认/键值对格式，便于阅读
	cfg.Color = true          // 开启颜色，便于阅读
	cfg.CallerInfo = false
	log := NewFLog(cfg)

	// 定义多组用例
	cases := []struct {
		name       string
		method     string
		path       string
		statusCode int
		handler    http.HandlerFunc
	}{
		{
			name:       "GET OK",
			method:     "GET",
			path:       "/test/success",
			statusCode: http.StatusOK,
			handler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(10 * time.Millisecond)
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("Success"))
			},
		},
		{
			name:       "POST Client Error",
			method:     "POST",
			path:       "/test/client-error",
			statusCode: http.StatusBadRequest,
			handler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(5 * time.Millisecond)
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("Bad Request"))
			},
		},
		{
			name:       "PUT Server Error",
			method:     "PUT",
			path:       "/test/server-error",
			statusCode: http.StatusInternalServerError,
			handler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(15 * time.Millisecond)
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Internal Server Error"))
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// 每个用例单独的业务处理器并包上中间件
			middleware := LogRequest(log, tc.handler)

			// 构造请求并添加常见头部
			req := httptest.NewRequest(tc.method, tc.path, nil)
			req.Header.Set("User-Agent", "test-agent/1.0")
			req.Header.Set("X-Request-ID", fmt.Sprintf("req-%d", time.Now().UnixNano()))
			req.Header.Set("X-Forwarded-For", "203.0.113.10")
			rr := httptest.NewRecorder()

			// 在日志前后打印分隔，方便阅读
			t.Logf("===== START %s %s (%s) =====", tc.method, tc.path, tc.name)
			middleware.ServeHTTP(rr, req)
			t.Logf("===== END %s %s (%s) =====", tc.method, tc.path, tc.name)
		})
	}

	_ = log.Close()
}
