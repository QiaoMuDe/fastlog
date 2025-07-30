/*
fastlog_test.go - FastLog核心功能测试文件
包含对日志记录器初始化、并发日志处理、日志格式、日志级别过滤和文件切割等功能的单元测试。
*/
package fastlog

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

	// 清理日志目录
	if _, err := os.Stat("logs"); err == nil {
		if err := os.RemoveAll("logs"); err != nil {
			fmt.Printf("清理日志目录失败: %v\n", err)
		}
	}

	// 恢复原始输出
	if !testing.Verbose() {
		os.Stdout = originalStdout
		os.Stderr = originalStderr
		_ = nullFile.Close()
	}

	os.Exit(exitCode)
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

// TestLogFormats 测试日志库支持的所有日志格式
func TestLogFormats(t *testing.T) {
	// 定义日志格式及其对应的文件名
	formats := map[LogFormatType]string{
		Detailed: "detailed.log",
		Bracket:  "bracket.log",
		Json:     "json.log",
		Threaded: "threaded.log",
		// 注意：对于Custom格式，日志库内部不进行格式化，需要在外部格式化后传入
		Custom: "custom.log",
	}

	// 为每种格式创建日志记录器并测试所有日志级别
	for format, filename := range formats {
		t.Run(fmt.Sprintf("Format_%s", filename), func(t *testing.T) {
			// 创建日志配置
			cfg := NewFastLogConfig("logs", filename)
			cfg.LogLevel = DEBUG
			cfg.LogFormat = format
			cfg.OutputToConsole = true // 禁用控制台输出，避免测试干扰

			// 创建日志记录器
			log, err := NewFastLog(cfg)
			if err != nil {
				t.Fatalf("创建日志记录器失败: %v", err)
			}
			defer func() { log.Close() }()

			// 测试所有日志级别
			// 对于Custom格式，需要外部格式化
			if format == Custom {
				log.Debugf("[自定义格式] [%s] %s", "DEBUG", "这是一条调试日志")
				log.Infof("[自定义格式] [%s] %s", "INFO", "这是一条信息日志")
				log.Successf("[自定义格式] [%s] %s", "SUCCESS", "这是一条成功日志")
				log.Warnf("[自定义格式] [%s] %s", "WARN", "这是一条警告日志")
				log.Errorf("[自定义格式] [%s] %s", "ERROR", "这是一条错误日志")
			} else {
				log.Debug("这是一条调试日志")
				log.Info("这是一条信息日志")
				log.Success("这是一条成功日志")
				log.Warn("这是一条警告日志")
				log.Error("这是一条错误日志")
			}

			// 给日志处理器一些时间来处理和写入日志
			time.Sleep(100 * time.Millisecond)
		})
	}
}
