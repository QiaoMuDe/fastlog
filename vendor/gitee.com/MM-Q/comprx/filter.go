// Package comprx 提供文件过滤功能，支持从忽略文件加载排除模式。
//
// 该文件提供了从忽略文件（如 .gitignore、.comprxignore）加载排除模式的功能。
// 支持标准的忽略文件格式，包括注释行、空行处理和 glob 模式匹配。
//
// 主要功能：
//   - 从忽略文件加载排除模式
//   - 支持注释行和空行处理
//   - 自动去重排除模式
//   - 支持 glob 模式匹配
//   - 提供文件不存在时的容错处理
//
// 支持的忽略文件格式：
//   - 每行一个模式
//   - # 开头的注释行
//   - 空行自动忽略
//   - 支持标准 glob 通配符
//
// 使用示例：
//
//	// 加载忽略文件，文件不存在会报错
//	patterns, err := comprx.LoadExcludeFromFile(".gitignore")
//
//	// 加载忽略文件，文件不存在返回空列表
//	patterns, err := comprx.LoadExcludeFromFileOrEmpty(".comprxignore")
package comprx

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadExcludeFromFile 从忽略文件加载排除模式
//
// 参数:
//   - ignoreFilePath: 忽略文件路径（如 ".comprxignore", ".gitignore"）
//
// 返回:
//   - []string: 排除模式列表（已去重）
//   - error: 错误信息
//
// 支持的文件格式:
//   - 每行一个模式
//   - 支持 # 开头的注释行
//   - 自动忽略空行
//   - 支持 glob 模式匹配
//   - 自动去除重复模式
//
// 使用示例:
//
//	patterns, err := comprx.LoadExcludeFromFile(".comprxignore")
func LoadExcludeFromFile(ignoreFilePath string) ([]string, error) {
	// 参数验证
	if ignoreFilePath == "" {
		return nil, fmt.Errorf("忽略文件路径不能为空")
	}

	// 获取绝对路径用于错误报告
	absPath, err := filepath.Abs(ignoreFilePath)
	if err != nil {
		return nil, fmt.Errorf("获取文件绝对路径失败 '%s': %w", ignoreFilePath, err)
	}

	file, err := os.Open(ignoreFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("忽略文件不存在: %s", absPath)
		}
		return nil, fmt.Errorf("打开忽略文件失败 '%s': %w", absPath, err)
	}
	defer func() { _ = file.Close() }()

	// 预分配切片容量 - 获取文件大小估算行数
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败 '%s': %w", absPath, err)
	}

	// 估算行数：假设平均每行20字符，预分配容量避免频繁扩容
	estimatedLines := int(stat.Size()/20) + 10 // 额外预留10行
	if estimatedLines < 16 {
		estimatedLines = 16 // 最小预分配16行
	}
	patterns := make([]string, 0, estimatedLines)

	// 用于去重的map
	seen := make(map[string]bool)

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释行
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 验证模式是否有效
		if _, err := filepath.Match(line, "test"); err != nil {
			return nil, fmt.Errorf("文件 '%s' 第 %d 行包含无效的 glob 模式 '%s': %w",
				filepath.Base(absPath), lineNum, line, err)
		}

		// 去重处理：只添加未见过的模式
		if !seen[line] {
			seen[line] = true
			patterns = append(patterns, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取忽略文件失败 '%s': %w", absPath, err)
	}

	return patterns, nil
}

// LoadExcludeFromFileOrEmpty 从忽略文件加载排除模式，文件不存在时返回空列表
//
// 参数:
//   - ignoreFilePath: 忽略文件路径
//
// 返回:
//   - []string: 排除模式列表，文件不存在时返回空列表
//   - error: 错误信息（文件不存在不算错误）
//
// 使用示例:
//
//	patterns, err := comprx.LoadExcludeFromFileOrEmpty(".comprxignore")
func LoadExcludeFromFileOrEmpty(ignoreFilePath string) ([]string, error) {
	// 参数验证
	if ignoreFilePath == "" {
		return nil, fmt.Errorf("忽略文件路径不能为空")
	}

	// 直接检查文件是否存在，避免包装错误的问题
	if _, err := os.Stat(ignoreFilePath); err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil // 文件不存在返回空列表，不是错误
		}
		return nil, fmt.Errorf("检查文件状态失败 '%s': %w", ignoreFilePath, err)
	}

	// 文件存在，调用正常的加载函数
	return LoadExcludeFromFile(ignoreFilePath)
}
