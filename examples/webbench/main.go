package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	fastlog "gitee.com/MM-Q/fastlog"
)

func main() {
	fmt.Println("FastLog Web 高并发日志写入性能测试")
	fmt.Println(repeat("=", 50))

	// 清理上次测试的日志文件
	_ = os.RemoveAll("logs")

	// 配置说明:
	//   使用 NewConfig 而非 Prod, 避免 Async 异步模式带来的缓冲延迟,
	//   确保写入后数据立即可见, 方便 reportStats 读取日志文件和大小。
	cfg := fastlog.NewConfig("logs/webbench.log")
	cfg.Level = fastlog.INFO
	cfg.OutputConsole = false

	logger := fastlog.New(cfg)

	// 模拟 Web 请求日志
	runWebBench(logger, 1000, 10)

	// 模拟真实场景: 混合日志类型 (增加操作量让耗时可测量)
	runMixedScenario(logger, 5000, 10)

	// 先关闭 Logger 确保所有日志刷入磁盘, 再统计结果
	_ = logger.Close()
	reportStats("logs")

	// 最后清理测试文件
	_ = os.RemoveAll("logs")
}

func runWebBench(logger *fastlog.Logger, totalRequests, concurrency int) {
	fmt.Printf("\n场景一: Web 请求并发写入 (%d 并发, %d 请求)\n", concurrency, totalRequests)
	fmt.Println(repeat("-", 50))

	start := time.Now()

	var wg sync.WaitGroup
	ch := make(chan int, totalRequests)
	for i := 0; i < totalRequests; i++ {
		ch <- i
	}
	close(ch)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for seq := range ch {
				simulateWebRequest(logger, seq)
			}
		}()
	}
	wg.Wait()

	elapsed := time.Since(start)
	throughput := float64(totalRequests) / elapsed.Seconds()

	fmt.Printf("  耗时: %v\n", elapsed.Round(time.Millisecond))
	fmt.Printf("  吞吐量: %.0f 请求/秒\n", throughput)
}

func simulateWebRequest(logger *fastlog.Logger, seq int) {
	// 模拟 HTTP 请求处理耗时
	time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)

	methods := []string{"GET", "POST", "PUT", "DELETE"}
	paths := []string{"/api/users", "/api/orders", "/api/products", "/api/auth", "/health"}
	statusCodes := []int{200, 201, 200, 204, 200, 200, 301, 400, 401, 403, 404, 500}
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
		"curl/8.0.1",
		"PostmanRuntime/7.36.0",
		"Go-http-client/2.0",
	}

	method := methods[rand.Intn(len(methods))]
	path := paths[rand.Intn(len(paths))]
	status := statusCodes[rand.Intn(len(statusCodes))]
	ua := userAgents[rand.Intn(len(userAgents))]
	duration := time.Duration(rand.Intn(200)+1) * time.Millisecond
	bodyBytes := rand.Intn(4096) + 64

	// 模拟不同的日志级别
	switch {
	case status >= 500:
		logger.Errorw("服务器内部错误",
			fastlog.String("method", method),
			fastlog.String("path", path),
			fastlog.Int("status", status),
			fastlog.String("user_agent", ua),
			fastlog.Duration("latency", duration),
			fastlog.Int("bytes", bodyBytes),
			fastlog.Int("request_id", seq),
			fastlog.String("error", "internal server error"),
		)
	case status >= 400:
		logger.Warnw("客户端请求异常",
			fastlog.String("method", method),
			fastlog.String("path", path),
			fastlog.Int("status", status),
			fastlog.String("user_agent", ua),
			fastlog.Duration("latency", duration),
			fastlog.Int("bytes", bodyBytes),
			fastlog.Int("request_id", seq),
		)
	default:
		logger.Infow("请求处理完成",
			fastlog.String("method", method),
			fastlog.String("path", path),
			fastlog.Int("status", status),
			fastlog.String("user_agent", ua),
			fastlog.Duration("latency", duration),
			fastlog.Int("bytes", bodyBytes),
			fastlog.Int("request_id", seq),
		)
	}
}

func runMixedScenario(logger *fastlog.Logger, totalOps, concurrency int) {
	fmt.Printf("\n场景二: 混合业务日志 (%d 并发, %d 条操作)\n", concurrency, totalOps)
	fmt.Println(repeat("-", 50))

	start := time.Now()

	var wg sync.WaitGroup
	ch := make(chan int, totalOps)
	for i := 0; i < totalOps; i++ {
		ch <- i
	}
	close(ch)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for seq := range ch {
				simulateBusinessOp(logger, seq)
			}
		}()
	}
	wg.Wait()

	elapsed := time.Since(start)
	throughput := float64(totalOps) / elapsed.Seconds()

	fmt.Printf("  耗时: %v\n", elapsed.Round(time.Millisecond))
	fmt.Printf("  吞吐量: %.0f 操作/秒\n", throughput)
}

func simulateBusinessOp(logger *fastlog.Logger, seq int) {
	// 模拟业务处理耗时
	time.Sleep(time.Duration(rand.Intn(3)) * time.Millisecond)

	// 随机模拟一种业务操作
	op := rand.Intn(4)
	userID := rand.Intn(10000) + 1

	switch op {
	case 0:
		// 用户登录
		logger.Infow("用户登录",
			fastlog.Int("user_id", userID),
			fastlog.String("username", fmt.Sprintf("user_%d", userID)),
			fastlog.String("ip", fmt.Sprintf("192.168.%d.%d", rand.Intn(255), rand.Intn(255))),
			fastlog.Bool("success", true),
			fastlog.Int("seq", seq),
		)
	case 1:
		// 数据库查询
		dur := time.Duration(rand.Intn(100)+1) * time.Millisecond
		logger.Infow("数据库查询",
			fastlog.String("query", "SELECT * FROM orders WHERE user_id = ?"),
			fastlog.Duration("elapsed", dur),
			fastlog.Int("rows", rand.Intn(50)),
			fastlog.Int("user_id", userID),
			fastlog.Int("seq", seq),
		)
	case 2:
		// 缓存操作
		logger.Infow("缓存操作",
			fastlog.String("action", "set"),
			fastlog.String("key", fmt.Sprintf("user:%d:profile", userID)),
			fastlog.Int("ttl", 3600),
			fastlog.Bool("hit_cache", rand.Intn(2) == 1),
			fastlog.Int("seq", seq),
		)
	case 3:
		// 外部 API 调用
		dur := time.Duration(rand.Intn(500)+10) * time.Millisecond
		status := rand.Intn(2)
		if status == 0 {
			logger.Infow("外部 API 调用",
				fastlog.String("service", "payment-gateway"),
				fastlog.String("endpoint", "/v1/charge"),
				fastlog.Duration("elapsed", dur),
				fastlog.Int("status", 200),
				fastlog.Int("seq", seq),
			)
		} else {
			logger.Warnw("外部 API 超时",
				fastlog.String("service", "sms-provider"),
				fastlog.String("endpoint", "/v1/send"),
				fastlog.Duration("elapsed", dur),
				fastlog.Int("status", 504),
				fastlog.Int("retry", 3),
				fastlog.Int("seq", seq),
			)
		}
	}
}

func reportStats(logDir string) {
	fmt.Printf("\n%s\n", repeat("=", 50))
	fmt.Println("测试完成!")

	entries, err := os.ReadDir(logDir)
	if err != nil {
		fmt.Printf("  无法读取日志目录: %v\n", err)
		return
	}

	var totalSize int64
	var fileCount int
	for _, entry := range entries {
		if !entry.IsDir() {
			info, err := entry.Info()
			if err == nil {
				totalSize += info.Size()
				fileCount++
			}
		}
	}

	fmt.Printf("  日志文件数: %d\n", fileCount)
	fmt.Printf("  总写入大小: %s\n", humanBytes(totalSize))
}

func humanBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func repeat(ch string, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = ch[0]
	}
	return string(b)
}
