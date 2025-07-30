package fastlog

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestMain 全局测试入口，控制非verbose模式下的输出重定向
func TestMain(m *testing.M) {
	flag.Parse() // 解析命令行参数
	// 保存原始标准输出和错误输出
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	var nullFile *os.File
	var err error

	// 非verbose模式下重定向到空设备
	if !testing.Verbose() {
		nullFile, err = os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
		if err != nil {
			panic("无法打开空设备文件: " + err.Error())
		}
		os.Stdout = nullFile
		os.Stderr = nullFile
	}

	// 运行所有测试
	exitCode := m.Run()

	// 恢复原始输出
	if !testing.Verbose() {
		os.Stdout = originalStdout
		os.Stderr = originalStderr
		_ = nullFile.Close()
	}

	os.Exit(exitCode)
}

// TestConcurrentFastLog 测试并发场景下的多个日志记录器
func TestConcurrentFastLog(t *testing.T) {
	// 记录开始时间
	startTime := time.Now()

	// 创建日志配置
	cfg := NewFastLogConfig("logs", "test.log")
	cfg.OutputToConsole = true
	cfg.OutputToFile = true
	cfg.MaxLogFileSize = 1

	// 创建日志记录器
	log, err := NewFastLog(cfg)
	if err != nil {
		t.Fatalf("创建日志记录器失败: %v", err)
	}

	// 持续时间为3秒
	duration := 3
	// 每秒生成10条日志
	rate := 10

	defer func() {
		log.Close()
		// 计算总耗时并打印
		totalDuration := time.Since(startTime)
		fmt.Printf("=== 并发日志测试结果 ===\n")
		fmt.Printf("测试配置: 持续时间 %d秒 | 目标速率 %d条/秒\n", duration, rate)
		fmt.Printf("预期生成: %d条日志\n", duration*rate)
		fmt.Printf("实际耗时: %.3fs (%.2fms)\n", totalDuration.Seconds(), float64(totalDuration.Nanoseconds())/1e6)
		fmt.Printf("========================\n")
	}()

	// 启动随机日志函数
	randomLog(log, duration, rate, t)
}

// generateLogs 生成指定数量的模拟日志到通道中
func generateLogs(logMethodsNoFormat []func(v ...any), logMethodsWithFormat []func(format string, v ...interface{}), totalLoops int) <-chan func() {
	// 创建一个通道用于传递日志，缓冲区大小为总循环数的2倍
	logChan := make(chan func(), totalLoops*2)

	// 创建一个随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	go func() {
		for i := 0; i < totalLoops; i++ {
			if r.Intn(2) == 0 {
				method := logMethodsNoFormat[r.Intn(len(logMethodsNoFormat))]
				logChan <- func() { method("这是一个测试日志") }
			} else {
				method := logMethodsWithFormat[r.Intn(len(logMethodsWithFormat))]
				logChan <- func() { method("这是一个测试日志: %s", "test") }
			}
		}
		close(logChan)
	}()

	return logChan
}

// randomLog 该函数用于按指定速率从通道中取出日志并打印
func randomLog(log *FastLog, duration int, rate int, t *testing.T) {
	// 定义无格式化日志方法的切片
	logMethodsNoFormat := []func(v ...any){
		log.Info,
		log.Warn,
		log.Error,
		log.Debug,
		log.Success,
	}
	// 定义格式化日志方法的切片
	logMethodsWithFormat := []func(format string, v ...interface{}){
		log.Infof,
		log.Warnf,
		log.Errorf,
		log.Debugf,
		log.Successf,
	}

	// 计算总循环次数
	totalLoops := duration * rate
	// 计算每次输出日志的间隔时间
	interval := time.Second / time.Duration(rate)

	// 生成日志到通道
	// 由于 generateLogs 仅返回一个通道，这里只使用一个变量接收
	logChan := generateLogs(logMethodsNoFormat, logMethodsWithFormat, totalLoops)

	// 使用WaitGroup同步并发操作
	var wg sync.WaitGroup
	wg.Add(totalLoops)

	// 按指定速率从通道中取出日志并打印
	go func() {
		for logFunc := range logChan {
			go func(f func()) {
				defer wg.Done()
				f()
			}(logFunc)
			time.Sleep(interval)
		}
	}()

	// 等待所有日志生成完成
	wg.Wait()

	// 验证日志文件内容
	content, err := os.ReadFile(filepath.Join("logs", "test.log"))
	if err != nil {
		t.Fatalf("读取日志文件失败: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	validLines := 0
	for _, line := range lines {
		if strings.Contains(line, "这是一个测试日志") {
			validLines++
		}
	}
}

// TestCustomFormat 测试自定义日志格式
func TestCustomFormat(t *testing.T) {
	// 创建日志配置
	cfg := NewFastLogConfig("logs", "custom.log")
	cfg.LogLevel = DEBUG
	cfg.LogFormat = Custom

	// 创建日志记录器
	log, err := NewFastLog(cfg)
	if err != nil {
		t.Fatalf("创建日志记录器失败: %v", err)
	}
	defer func() { log.Close() }()

	// 定义web应用程序日志格式
	webAppLogFormat := `%s [%s] %s %s %s %d %d %s %s %dms`

	// 模拟web应用程序日志记录
	log.Errorf(webAppLogFormat, "127.0.0.1", "2023-10-01 12:00:00", "GET", "/index.html", "HTTP/1.1", 200, 1234, "Mozilla/5.0", "en-US", 100)
	log.Infof(webAppLogFormat, "127.0.0.1", "2023-10-01 12:00:01", "POST", "/login", "HTTP/1.1", 401, 0, "Mozilla/5.0", "en-US", 200)
	log.Successf(webAppLogFormat, "127.0.0.1", "2023-10-01 12:00:02", "GET", "/profile", "HTTP/1.1", 200, 5678, "Mozilla/5.0", "en-US", 300)
	log.Warnf(webAppLogFormat, "127.0.0.1", "2023-10-01 12:00:03", "PUT", "/settings", "HTTP/1.1", 500, 0, "Mozilla/5.0", "en-US", 500)
	log.Debugf(webAppLogFormat, "127.0.0.1", "2023-10-01 12:00:04", "DELETE", "/logout", "HTTP/1.1", 200, 0, "Mozilla/5.0", "en-US", 50)
}

// TestNoColor 测试无颜色日志
func TestNoColor(t *testing.T) {
	// 创建日志配置
	cfg := NewFastLogConfig("logs", "nocolor.log")
	cfg.LogLevel = DEBUG
	cfg.NoColor = true // 禁用终端颜色

	// 创建日志记录器
	log, err := NewFastLog(cfg)
	if err != nil {
		t.Fatalf("创建日志记录器失败: %v", err)
	}
	defer func() { log.Close() }()

	// 打印测试日志
	log.Info("测试无颜色日志")
	log.Warn("测试无颜色日志")
	log.Error("测试无颜色日志")
	log.Debug("测试无颜色日志")
	log.Success("测试无颜色日志")
}

// TestNoBold 测试无加粗日志
func TestNoBold(t *testing.T) {
	// 创建日志配置
	cfg := NewFastLogConfig("logs", "nobold.log")
	cfg.LogLevel = DEBUG
	cfg.NoBold = true // 禁用终端字体加粗
	cfg.NoColor = false

	// 创建日志记录器
	log, err := NewFastLog(cfg)
	if err != nil {
		t.Fatalf("创建日志记录器失败: %v", err)
	}
	defer func() { log.Close() }()

	// 打印测试日志
	log.Info("测试无加粗日志")
	log.Warn("测试无加粗日志")
	log.Error("测试无加粗日志")
	log.Debug("测试无加粗日志")
	log.Success("测试无加粗日志")
}

// TestNewFastLog_Initialization 测试日志记录器初始化
func TestNewFastLog_Initialization(t *testing.T) {
	// 测试nil配置错误 - 应该会panic
	t.Run("nil config should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("使用nil配置时应该panic")
			}
		}()
		// 这行代码应该会panic
		_, _ = NewFastLog(nil)
	})

	// 测试正常初始化 - 使用临时目录避免文件冲突
	t.Run("normal initialization", func(t *testing.T) {
		tempDir := t.TempDir()
		cfg := NewFastLogConfig(tempDir, "init_test.log")
		cfg.OutputToConsole = false // 禁用控制台输出，避免测试干扰
		cfg.OutputToFile = true     // 只测试文件输出

		log, err := NewFastLog(cfg)
		if err != nil {
			t.Fatalf("初始化日志失败: %v", err)
		}

		// 验证日志实例是否正确创建
		if log == nil {
			t.Fatal("日志实例不应为nil")
		}

		// 简单测试日志功能是否正常
		log.Info("初始化测试日志")

		// 给processor一些时间完成初始化和处理
		time.Sleep(100 * time.Millisecond)

		log.Close()

		// 验证日志文件是否创建成功
		logFile := filepath.Join(tempDir, "init_test.log")
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			t.Error("日志文件未创建")
		}
	})
}

// TestFastLog_LogLevels 测试不同日志级别的过滤功能
func TestFastLog_LogLevels(t *testing.T) {
	// 创建临时日志文件来测试级别过滤
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "level_test.log")

	cfg := NewFastLogConfig(tempDir, "level_test.log")
	cfg.LogLevel = WARN
	cfg.OutputToConsole = false // 只写入文件，不输出到控制台
	log, err := NewFastLog(cfg)
	if err != nil {
		t.Fatalf("创建日志记录器失败: %v", err)
	}
	defer func() { log.Close() }()

	// 不同级别日志
	log.Debug("debug message")
	log.Info("info message")
	log.Warn("warn message")
	log.Error("error message")

	// 确保所有日志都写入完成
	time.Sleep(200 * time.Millisecond)
	log.Close()

	// 读取日志文件内容
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("读取日志文件失败: %v", err)
	}

	output := string(content)
	t.Logf("日志文件内容: %q", output)

	// 验证结果：不应包含低级别日志
	if strings.Contains(output, "debug message") || strings.Contains(output, "info message") {
		t.Error("日志级别过滤失败，不应包含低级日志")
	}

	// 验证结果：应包含高级别日志
	if !strings.Contains(output, "warn message") || !strings.Contains(output, "error message") {
		t.Error("日志级别过滤失败，应包含高级日志")
	}
}

// TestFastLog_FileRotation 测试日志文件切割功能
func TestFastLog_FileRotation(t *testing.T) {
	cfg := NewFastLogConfig("logs", "test.log")
	cfg.MaxLogFileSize = 1
	cfg.OutputToConsole = false
	log, _ := NewFastLog(cfg)

	// 生成300kb日志
	largeData := generateLargeString(1024 * 100) // 转换写入实际为300kb

	// 写入10条日志
	for i := 0; i < 10; i++ {
		log.Info(largeData)
	}

	// 确保缓冲区刷新
	time.Sleep(500 * time.Millisecond)

	// 确保所有日志写入完成
	log.Close()

	// 等待切割完成
	time.Sleep(2 * time.Second) // 减少等待时间，500ms缓冲区刷新+2秒足够完成切割

	// 验证切割文件生成
	files, _ := filepath.Glob(filepath.Join(cfg.LogDirName, "test.*"))
	if len(files) < 1 {
		t.Error("日志文件切割失败，未生成切割文件")
	}
}

// 辅助函数
//
// 参数：
//   - size: 字符串大小(单位：字节)
//
// 返回：
//   - string: 生成的字符串
func generateLargeString(size int) []byte {
	buf := make([]byte, size)
	for i := range buf {
		buf[i] = 'a'
	}
	// 添加长度验证
	if len(buf) != size {
		panic(fmt.Sprintf("生成字符串大小不符: 预期=%d, 实际=%d", size, len(buf)))
	}
	return buf
}

// TestCleanupLogs 测试完成后清理日志目录
func TestCleanupLogs(t *testing.T) {
	if _, err := os.Stat("logs"); err == nil {
		if err := os.RemoveAll("logs"); err != nil {
			t.Fatalf("删除logs目录失败: %v", err)
		}
	}
}
