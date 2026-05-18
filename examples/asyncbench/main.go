package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	fastlog "gitee.com/MM-Q/fastlog"
)

func main() {
	fmt.Println("FastLog Async 异步轮转/压缩性能测试")
	fmt.Println(repeat("=", 55))

	// 清理上次测试的日志文件
	_ = os.RemoveAll("logs")

	// ============================================================
	// 配置说明:
	//   使用 Prod 配置 (Async: true, Compress: true) 演示高频日志
	//   写入场景下, 后台异步轮转和压缩的效果。
	//
	//   Async: true    — 日志轮转和压缩在后台 goroutine 执行,
	//                     不阻塞写入路径, 适合高频场景。
	//   MaxSize: 1     — 设置为 1MB, 快速触发轮转。
	//   Compress: true — 轮转后的旧文件自动压缩为 .gz。
	// ============================================================
	cfg := fastlog.Prod("logs/asyncbench.log")
	cfg.Level = fastlog.INFO
	cfg.MaxSize = 1           // 1MB 触发轮转 (单位: MB, 降低阈值确保测试中多次触发)
	cfg.MaxFiles = 10         // 保留 10 个历史文件
	cfg.MaxAge = 1            // 保留 1 天 (避免测试期间被清理)
	cfg.Compress = true       // 启用压缩
	cfg.Async = true          // 异步轮转/压缩 (Prod 默认已启用)
	cfg.Caller = true         // 记录调用者信息 (观察日志格式)
	cfg.RotateByDay = false   // 关闭按天轮转, 只按大小轮转
	cfg.DateDirLayout = false // 关闭日期目录, 所有文件平铺展示

	logger := fastlog.New(cfg)
	defer func() { _ = logger.Close() }()

	// 阶段一: 高频写入 — 大量结构化日志触发多次轮转
	// 使用单 goroutine 顺序写入, 避免 Windows 文件锁冲突导致 async 轮转失败
	highFrequencyWrite(logger, 2_000_000, 1)

	// 阶段二: 观察写入后, Close 前文件状态 (看看哪些已轮转/压缩)
	fmt.Println("\n📦 Close 前文件状态:")
	showLogDir("logs")

	// 关闭 Logger: 触发缓冲刷新 + 写入器关闭
	_ = logger.Close()

	// 等待异步轮转/压缩完成 (后台 goroutine 需要时间处理)
	fmt.Println("\n⏳ 等待异步轮转/压缩完成...")
	time.Sleep(5 * time.Second)

	// 阶段三: 查看 Close + 异步处理后的最终文件状态
	fmt.Println("\n📦 Close 后文件状态:")
	showLogDir("logs")

	// 阶段四: 统计摘要
	showSummary("logs")

	// 清理
	_ = os.RemoveAll("logs")
}

func highFrequencyWrite(logger *fastlog.Logger, total int, concurrency int) {
	fmt.Printf("\n⚡ 高频写入: %d 并发 × %d 条 = %d 条结构化日志\n",
		concurrency, total/concurrency, total)
	fmt.Println(repeat("-", 55))

	start := time.Now()

	var wg sync.WaitGroup
	ch := make(chan int, total)
	for i := 0; i < total; i++ {
		ch <- i
	}
	close(ch)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for seq := range ch {
				writeLogLine(logger, seq)
			}
		}()
	}
	wg.Wait()

	elapsed := time.Since(start)
	throughput := float64(total) / elapsed.Seconds()

	fmt.Printf("  耗时: %v\n", elapsed.Round(time.Millisecond))
	fmt.Printf("  吞吐量: %.0f 条/秒\n", throughput)
}

func writeLogLine(logger *fastlog.Logger, seq int) {
	// 轮流写入不同业务类型的日志, 模拟真实高频场景
	switch seq % 4 {
	case 0:
		logger.Infow("订单创建",
			fastlog.Int("order_id", seq),
			fastlog.Int("user_id", seq%10000),
			fastlog.String("product", fmt.Sprintf("sku_%d", seq%500)),
			fastlog.Float64("amount", float64(seq%10000)/100),
			fastlog.Int("quantity", seq%10+1),
			fastlog.String("status", "pending"),
		)
	case 1:
		logger.Infow("支付回调",
			fastlog.Int("order_id", seq),
			fastlog.String("channel", []string{"alipay", "wechat", "card"}[seq%3]),
			fastlog.String("trade_no", fmt.Sprintf("T%015d", seq)),
			fastlog.Float64("amount", float64(seq%50000)/100),
			fastlog.Bool("success", seq%10 != 0),
			fastlog.Int("retry_count", seq%3),
		)
	case 2:
		logger.Warnw("库存预警",
			fastlog.String("sku", fmt.Sprintf("sku_%d", seq%500)),
			fastlog.Int("remaining", seq%20),
			fastlog.Int("threshold", 50),
			fastlog.String("warehouse", []string{"华东", "华南", "华北", "西南"}[seq%4]),
			fastlog.Int("seq", seq),
		)
	case 3:
		logger.Infow("用户行为",
			fastlog.Int("user_id", seq%10000),
			fastlog.String("action", []string{"view", "click", "add_cart", "favorite"}[seq%4]),
			fastlog.String("page", fmt.Sprintf("/product/%d", seq%500)),
			fastlog.Duration("stay", time.Duration(seq%30000)*time.Millisecond),
			fastlog.String("device", []string{"mobile", "pc", "pad"}[seq%3]),
			fastlog.Int("seq", seq),
		)
	}
}

func showLogDir(dir string) {
	type fileInfo struct {
		name string
		size int64
	}

	var files []fileInfo
	var totalSize int64

	// 递归扫描目录 (DateDirLayout 会将轮转文件放入日期子目录)
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(dir, path)
		files = append(files, fileInfo{name: rel, size: info.Size()})
		totalSize += info.Size()
		return nil
	})

	if len(files) == 0 {
		fmt.Println("  (空目录)")
		return
	}

	// 按文件名排序
	sort.Slice(files, func(i, j int) bool {
		return files[i].name < files[j].name
	})

	for _, f := range files {
		var mark string
		if strings.HasSuffix(f.name, ".gz") {
			mark = "�"
		} else if strings.HasSuffix(f.name, ".log") {
			// 判断是否是当前活动文件
			if filepath.Base(f.name) == "asyncbench.log" {
				mark = "📄"
			} else {
				mark = "🗜"
			}
		} else {
			mark = "🗜"
		}
		fmt.Printf("  %s %-40s %s\n", mark, f.name, humanBytes(f.size))
	}
	fmt.Printf("  %s %-40s %s\n", "", "----------------------------------------", "--------")
	fmt.Printf("  %s %-40s %s (%d 文件)\n", "", "合计", humanBytes(totalSize), len(files))
}

func showSummary(dir string) {
	fmt.Printf("\n%s\n", repeat("=", 55))
	fmt.Println("异步轮转/压缩摘要")
	fmt.Println(repeat("=", 55))

	// 递归扫描所有文件
	var allFiles []os.FileInfo
	var allPaths []string
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		allFiles = append(allFiles, info)
		allPaths = append(allPaths, path)
		return nil
	})

	if len(allFiles) == 0 {
		fmt.Println("  无法读取日志目录")
		return
	}

	var currentSize int64
	var archiveCount int
	var archiveSize int64
	var compressedCount int
	var compressedSize int64

	for i, info := range allFiles {
		name := allPaths[i]

		if strings.HasSuffix(name, ".gz") {
			compressedCount++
			compressedSize += info.Size()
		} else if strings.HasSuffix(name, "asyncbench.log") {
			currentSize += info.Size()
		} else {
			archiveCount++
			archiveSize += info.Size()
		}
	}

	fmt.Printf("  当前活动文件:     %s\n", humanBytes(currentSize))
	fmt.Printf("  已轮转(未压缩):   %d 个, 共 %s\n", archiveCount, humanBytes(archiveSize))
	fmt.Printf("  已压缩:           %d 个, 共 %s\n", compressedCount, humanBytes(compressedSize))
	fmt.Printf("  总文件数:         %d\n", len(allFiles))

	totalLogData := currentSize + archiveSize + compressedSize
	fmt.Printf("  总原始数据:       %s\n", humanBytes(totalLogData))
	if compressedSize > 0 {
		fmt.Printf("  压缩后节省:        %.1f%%\n",
			(1-float64(compressedSize)/float64(totalLogData))*100)
	}

	fmt.Printf("\n✅ Async 异步模式说明:\n")
	fmt.Printf("   轮转和压缩操作在后台 goroutine 中异步执行,\n")
	fmt.Printf("   不阻塞日志写入路径, 保证高频场景下的写入性能。\n")
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
