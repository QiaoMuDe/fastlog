// Package types 提供文件过滤功能的核心类型和接口定义。
//
// 该文件定义了文件过滤器接口和过滤选项结构体，用于在压缩和解压缩过程中
// 根据文件路径、大小等条件过滤文件。支持 glob 模式匹配、文件大小限制等功能。
//
// 主要类型：
//   - FileFilter: 文件过滤器接口
//   - FilterOptions: 过滤配置选项结构体
//
// 主要功能：
//   - 支持包含和排除模式的 glob 匹配
//   - 支持文件大小范围过滤
//   - 提供高性能的快速匹配算法
//   - 支持复杂的 glob 模式匹配
//   - 提供过滤条件验证功能
//
// 使用示例：
//
//	// 创建过滤选项
//	filter := &types.FilterOptions{
//	    Include: []string{"*.go", "*.md"},
//	    Exclude: []string{"*_test.go"},
//	    MaxSize: 10 * 1024 * 1024, // 10MB
//	}
//
//	// 检查文件是否应该跳过
//	shouldSkip := filter.ShouldSkipByParams("main.go", 1024, false)
package types

import (
	"fmt"
	"path/filepath"
	"strings"
)

// FilterOptions 过滤配置选项
//
// 用于指定压缩时或解压时的文件过滤条件：
//   - Include: 包含模式列表，支持 glob 语法，只有匹配的文件才会被处理
//   - Exclude: 排除模式列表，支持 glob 语法，匹配的文件会被跳过
//   - MaxSize: 最大文件大小限制（字节），0 表示无限制，超过此大小的文件会被跳过
//   - MinSize: 最小文件大小限制（字节），默认为 0，小于此大小的文件会被跳过
type FilterOptions struct {
	Include []string // 包含模式，支持 glob 语法，只处理匹配的文件
	Exclude []string // 排除模式，支持 glob 语法，跳过匹配的文件
	MaxSize int64    // 最大文件大小（字节），0 表示无限制
	MinSize int64    // 最小文件大小（字节），默认为 0
}

// ShouldSkipByParams 判断文件是否应该被跳过(通用方法，用于压缩和解压)
//
// 过滤逻辑:
//  1. 检查文件大小是否符合要求
//  2. 如果指定了包含模式，检查文件是否匹配包含模式
//  3. 检查文件是否匹配排除模式
//
// 参数:
//   - path: 文件路径
//   - size: 文件大小（字节）
//   - isDir: 是否为目录
//
// 返回:
//   - bool: true 表示应该跳过，false 表示应该处理
func (f *FilterOptions) ShouldSkipByParams(path string, size int64, isDir bool) bool {
	// 如果过滤器为空或没有过滤条件，不跳过任何文件
	if !HasFilterConditions(f) {
		return false
	}

	// 1. 检查文件大小 - 不符合大小要求的文件需要跳过
	if f.shouldSkipBySize(size, isDir) {
		return true
	}

	// 2. 检查包含模式（如果指定了包含模式）
	// 不匹配包含模式的文件应该被跳过
	if len(f.Include) > 0 {
		if !f.matchAnyPattern(f.Include, path) {
			return true
		}
	}

	// 3. 检查排除模式
	// 匹配排除模式的文件应该被跳过
	if len(f.Exclude) > 0 {
		if f.matchAnyPattern(f.Exclude, path) {
			return true
		}
	}

	// 通过所有检查，不应该被跳过
	return false
}

// shouldSkipBySize 检查文件是否因大小限制需要跳过
//
// 检查逻辑：
//   - 目录：不受大小限制，永远不跳过
//   - 文件：检查是否在设定的大小范围内
//
// 参数:
//   - size: 文件大小（字节）
//   - isDir: 是否为目录
//
// 返回:
//   - bool: true=需要跳过, false=不需要跳过
//
// 示例:
//
//	MinSize=1024, MaxSize=1048576 (1KB-1MB)
//	- 500字节的文件 → 返回true（太小，跳过）
//	- 2048字节的文件 → 返回false（合适，不跳过）
//	- 2MB的文件 → 返回true（太大，跳过）
func (f *FilterOptions) shouldSkipBySize(size int64, isDir bool) bool {
	// 目录不受大小限制，永远不跳过
	if isDir {
		return false
	}

	// 快速路径: 如果没有设置大小限制，直接返回不跳过
	if f.MinSize <= 0 && f.MaxSize <= 0 {
		return false
	}

	// 文件太小，跳过（例如：设置MinSize=1024，文件只有500字节）
	if f.MinSize > 0 && size < f.MinSize {
		return true
	}

	// 文件太大，跳过（例如：设置MaxSize=1MB，文件有2MB）
	if f.MaxSize > 0 && size > f.MaxSize {
		return true
	}

	// 文件大小在允许范围内，不跳过
	return false
}

// matchAnyPattern 检查路径是否匹配任一模式
//
// 参数:
//   - patterns: 模式列表
//   - path: 文件路径
//
// 返回:
//   - bool: true 表示匹配任一模式，false 表示不匹配任何模式
func (f *FilterOptions) matchAnyPattern(patterns []string, path string) bool {
	// 快速失败: 如果没有模式或路径为空，直接返回
	if len(patterns) == 0 || path == "" {
		return false
	}

	// 预计算路径信息，避免在循环中重复计算
	baseName := filepath.Base(path)
	slashPath := filepath.ToSlash(path)

	for _, pattern := range patterns {
		// 跳过空模式
		if pattern == "" {
			continue
		}

		// 统一标准化模式，避免在子函数中重复处理
		normalizedPattern := strings.ReplaceAll(pattern, "\\", "/")

		// 匹配模式
		if f.matchPattern(normalizedPattern, path, baseName, slashPath) {
			return true
		}
	}
	return false
}

// matchPattern 检查路径是否匹配指定模式
//
// 支持多种匹配方式:
//  1. 文件名匹配（如 *.go 匹配 main.go）
//  2. 完整路径匹配（如 src/*.go 匹配 src/main.go）
//  3. 目录匹配（如 vendor/ 匹配 vendor 目录）
//
// 参数:
//   - normalizedPattern: 已标准化的 glob 模式（路径分隔符已统一为正斜杠）
//   - path: 文件路径
//   - baseName: 预计算的文件名（可选，空字符串表示需要计算）
//   - slashPath: 预计算的标准化路径（可选，空字符串表示需要计算）
//
// 返回:
//   - bool: true 表示匹配，false 表示不匹配
func (f *FilterOptions) matchPattern(normalizedPattern, path, baseName, slashPath string) bool {
	// 懒加载: 只在需要时计算文件名和标准化路径
	if baseName == "" {
		baseName = filepath.Base(path)
	}
	if slashPath == "" {
		slashPath = filepath.ToSlash(path)
	}

	// 快速匹配: 处理常见的简单模式，避免复杂的glob解析
	if matched, handled := f.fastMatch(normalizedPattern, slashPath, baseName); handled {
		return matched
	}

	// 复杂匹配(仅在快速匹配失败时使用): 处理复杂的 glob 模式，如 "src/*"、"vendor/**" 等
	return f.complexMatch(normalizedPattern, path, baseName, slashPath)
}

// fastMatch 快速匹配常见的简单模式
//
// 针对80%的常见模式进行优化，避免复杂的glob解析，显著提升性能。
// 支持的快速模式包括：*.ext、精确匹配、prefix*、dirname/、*pattern* 等
//
// 参数:
//   - normalizedPattern: 已标准化的 glob 模式字符串（路径分隔符已统一为正斜杠）
//   - slashPath: 标准化的文件路径（使用正斜杠）
//   - baseName: 文件名（不含路径）
//
// 返回值说明:
//   - matched (第1个返回值): 匹配结果（仅当 handled=true 时有效）
//   - handled (第2个返回值): 是否能够快速处理
//
// 匹配结果:
//   - handled=true: 已完成快速匹配，直接使用 matched 结果
//   - handled=false: 无法快速处理，需要回退到复杂的 glob 匹配
//
// 使用示例:
//
//	matched, handled := f.fastMatch("*.go", "src/main.go", "main.go")
//	if handled {
//	    return matched  // 快速匹配完成，直接返回结果
//	}
//	否则继续使用复杂匹配...
func (f *FilterOptions) fastMatch(normalizedPattern, slashPath, baseName string) (matched bool, handled bool) {
	// 0. 处理空模式 - 空模式不匹配任何文件
	if normalizedPattern == "" {
		return false, true
	}

	// 1. 处理后缀匹配 (*.ext) - 最常见的模式
	if strings.HasPrefix(normalizedPattern, "*.") && len(normalizedPattern) > 2 {
		ext := normalizedPattern[2:] // 获取扩展名

		// 确保扩展名中没有其他通配符
		if !strings.ContainsAny(ext, "*?[]") {
			return strings.HasSuffix(baseName, ext), true
		}
	}

	// 2. 处理精确匹配 (无通配符) - 如 "node_modules", "vendor"
	// 注意：跳过以 '/' 结尾的模式，这些应该由目录匹配处理
	if !strings.ContainsAny(normalizedPattern, "*?[]") && normalizedPattern[len(normalizedPattern)-1] != '/' {
		// 检查文件名精确匹配或路径包含匹配
		return baseName == normalizedPattern || strings.Contains(slashPath, normalizedPattern), true
	}

	// 3. 处理前缀匹配 (prefix*) - 如 "test*", "vendor/*"
	if strings.HasSuffix(normalizedPattern, "*") && len(normalizedPattern) > 1 {
		prefix := normalizedPattern[:len(normalizedPattern)-1] // 获取前缀

		// 确保前缀中没有其他通配符
		if !strings.ContainsAny(prefix, "*?[]") {
			// 检查文件名前缀匹配
			if strings.HasPrefix(baseName, prefix) {
				return true, true
			}
			// 检查路径前缀匹配（如 "vendor/*" 匹配 "vendor/package.go"）
			if strings.HasPrefix(slashPath, prefix) {
				return true, true
			}
			// 检查路径中是否包含该前缀（如 "vendor/*" 匹配 "src/vendor/package.go"）
			if strings.Contains(slashPath, "/"+prefix) {
				return true, true
			}
			return false, true
		}
	}

	// 4. 处理目录匹配 (dirname/) - 如 "node_modules/" 或 "test/"
	if len(normalizedPattern) > 1 && normalizedPattern[len(normalizedPattern)-1] == '/' {
		dirName := normalizedPattern[:len(normalizedPattern)-1] // 获取目录名

		// 确保目录名中没有通配符
		if !strings.ContainsAny(dirName, "*?[]") {
			// 检查完整路径是否匹配目录名（如 "vendor/" 匹配 "vendor"）
			if slashPath == dirName {
				return true, true
			}
			// 检查文件名是否匹配目录名（如 "vendor/" 匹配 "vendor"）
			if baseName == dirName {
				return true, true
			}
			// 检查路径是否以该目录开头（如 "vendor/" 匹配 "vendor/package.go"）
			if strings.HasPrefix(slashPath, dirName+"/") {
				return true, true
			}
			// 检查路径中是否包含该目录（如 "vendor/" 匹配 "src/vendor/package.go"）
			if strings.Contains(slashPath, "/"+dirName+"/") {
				return true, true
			}
			return false, true
		}
	}

	// 5. 处理中间通配符 (*pattern*) - 如 "*test*"
	if strings.HasPrefix(normalizedPattern, "*") && strings.HasSuffix(normalizedPattern, "*") && len(normalizedPattern) > 2 {
		middle := normalizedPattern[1 : len(normalizedPattern)-1] // 获取中间部分

		// 确保中间部分没有其他通配符
		if !strings.ContainsAny(middle, "*?[]") {
			return strings.Contains(baseName, middle) || strings.Contains(slashPath, middle), true
		}
	}

	// 无法使用快速匹配，需要使用复杂匹配
	return false, false
}

// complexMatch 复杂模式匹配
// 处理复杂的glob模式，如包含多个通配符、字符类等
//
// 参数:
//   - normalizedPattern: 已标准化的 glob 模式（路径分隔符已统一为正斜杠）
//   - path: 完整文件路径
//   - baseName: 文件名
//   - slashPath: 标准化路径
//
// 返回:
//   - bool: 是否匹配
func (f *FilterOptions) complexMatch(normalizedPattern, path, baseName, slashPath string) bool {
	// 优先使用标准化路径，如果为空则先标准化
	if slashPath == "" {
		slashPath = filepath.ToSlash(path)
	}

	// 1. 尝试匹配文件名
	if matched, err := filepath.Match(normalizedPattern, baseName); err == nil && matched {
		return true
	}

	// 2. 尝试匹配标准化路径
	if matched, err := filepath.Match(normalizedPattern, slashPath); err == nil && matched {
		return true
	}

	// 3. 尝试匹配目录(处理以正斜杠结尾的模式)
	if len(normalizedPattern) > 0 && normalizedPattern[len(normalizedPattern)-1] == '/' {
		dirPattern := normalizedPattern[:len(normalizedPattern)-1] // 获取目录模式

		// 先匹配目录名
		if matched, err := filepath.Match(dirPattern, baseName); err == nil && matched {
			return true
		}

		// 再匹配标准化路径
		if matched, err := filepath.Match(dirPattern, slashPath); err == nil && matched {
			return true
		}
	}

	// 4. 处理路径中包含模式的情况
	// 优化：避免不必要的分割，先检查模式复杂度
	if strings.ContainsAny(normalizedPattern, "*?[]") {
		pathParts := strings.Split(slashPath, "/") // 分割标准化路径
		for _, part := range pathParts {
			if matched, err := filepath.Match(normalizedPattern, part); err == nil && matched {
				return true
			}
		}
	}

	return false
}

// Validate 验证过滤器选项
//
// 返回:
//   - error: 验证错误，如果验证通过则返回 nil
func (f *FilterOptions) Validate() error {
	// 验证文件大小范围
	if f.MinSize < 0 {
		return fmt.Errorf("最小文件大小不能为负数: %d", f.MinSize)
	}

	if f.MaxSize < 0 {
		return fmt.Errorf("最大文件大小不能为负数: %d", f.MaxSize)
	}

	if f.MinSize > 0 && f.MaxSize > 0 && f.MinSize > f.MaxSize {
		return fmt.Errorf("最小文件大小 (%d) 不能大于最大文件大小 (%d)", f.MinSize, f.MaxSize)
	}

	// 验证包含模式
	for _, pattern := range f.Include {
		if pattern == "" {
			return fmt.Errorf("包含模式不能为空字符串")
		}
	}

	// 验证排除模式
	for _, pattern := range f.Exclude {
		if pattern == "" {
			return fmt.Errorf("排除模式不能为空字符串")
		}
	}

	return nil
}

// HasFilterConditions 检查过滤器是否有任何过滤条件
//
// 参数:
//   - filter: 过滤配置选项
//
// 返回:
//   - bool: true 表示有过滤条件，false 表示没有
func HasFilterConditions(filter *FilterOptions) bool {
	if filter == nil {
		return false
	}

	return len(filter.Include) > 0 ||
		len(filter.Exclude) > 0 ||
		filter.MinSize > 0 ||
		filter.MaxSize > 0
}

// 注意：LoadExcludeFromFile 和 LoadExcludeFromFileOrEmpty 函数已移至主包 comprx
// 请使用 comprx.LoadExcludeFromFile 和 comprx.LoadExcludeFromFileOrEmpty
