package fastlog

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestFatal 测试Fatal方法
func TestFatal(t *testing.T) {
	// 定义测试用例名称，用于环境变量标识
	const testName = "TestFatal"

	// 子进程模式：执行实际的Fatal调用
	if os.Getenv("TEST_MODE") == testName {
		config := NewFastLogConfig("test-logs", "fatal_test.log")
		if err := os.MkdirAll(config.GetLogDirName(), 0755); err != nil {
			panic(err)
		}
		log, err := NewFastLog(config)
		if err != nil {
			panic(err)
		}
		log.Fatal("fatal_test message")
		return
	}

	// 主进程模式：启动子进程并检查结果
	cmd := exec.Command(os.Args[0], "-test.run=^"+testName+"$", "-test.v")
	cmd.Env = append(os.Environ(), "TEST_MODE="+testName)
	output, err := cmd.CombinedOutput()

	// 检查是否以预期的退出码1退出
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitCode := exitErr.ExitCode(); exitCode != 1 {
			t.Errorf("期望退出码1，实际得到%d\n输出: %s", exitCode, output)
		}
	} else if err != nil {
		t.Fatalf("测试命令执行失败: %v\n输出: %s", err, output)
	} else {
		t.Error("Fatal方法没有导致程序退出")
	}

	// 验证日志文件是否正确创建并包含预期内容
	logPath := filepath.Join("test-logs", "fatal_test.log")
	defer func() {
		// 确保测试后清理文件
		_ = os.Remove(logPath)
		_ = os.Remove("test-logs")
	}()

	if _, statErr := os.Stat(logPath); os.IsNotExist(statErr) {
		t.Errorf("日志文件未创建: %s\n子进程输出: %s", logPath, output)
		return
	}

	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Errorf("无法读取日志文件: %v", err)
	} else if !strings.Contains(string(content), "fatal_test message") {
		t.Errorf("日志内容不包含预期消息，内容: %s", content)
	}
}

// TestFatalf 测试Fatalf方法
func TestFatalf(t *testing.T) {
	// 定义测试用例名称，用于环境变量标识
	const testName = "TestFatalf"

	// 子进程模式：执行实际的Fatalf调用
	if os.Getenv("TEST_MODE") == testName {
		config := NewFastLogConfig("test-logs", "fatalf_test.log")
		if err := os.MkdirAll(config.GetLogDirName(), 0755); err != nil {
			panic(err)
		}
		log, err := NewFastLog(config)
		if err != nil {
			panic(err)
		}
		log.Fatalf("fatalf_test %s message", "formatted")
		return
	}

	// 主进程模式：启动子进程并检查结果
	cmd := exec.Command(os.Args[0], "-test.run=^"+testName+"$", "-test.v")
	cmd.Env = append(os.Environ(), "TEST_MODE="+testName)
	output, err := cmd.CombinedOutput()

	// 检查是否以预期的退出码1退出
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitCode := exitErr.ExitCode(); exitCode != 1 {
			t.Errorf("期望退出码1，实际得到%d\n输出: %s", exitCode, output)
		}
	} else if err != nil {
		t.Fatalf("测试命令执行失败: %v\n输出: %s", err, output)
	} else {
		t.Error("Fatalf方法没有导致程序退出")
	}

	// 验证日志文件是否正确创建并包含预期内容
	logPath := filepath.Join("test-logs", "fatalf_test.log")
	defer func() {
		// 确保测试后清理文件
		_ = os.Remove(logPath)
		_ = os.Remove("test-logs")
	}()

	if _, statErr := os.Stat(logPath); os.IsNotExist(statErr) {
		t.Errorf("日志文件未创建: %s\n子进程输出: %s", logPath, output)
		return
	}

	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Errorf("无法读取日志文件: %v", err)
	} else if !strings.Contains(string(content), "fatalf_test formatted message") {
		t.Errorf("日志内容不包含预期消息，内容: %s", content)
	}
}
