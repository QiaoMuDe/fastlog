// Package comprx 提供内存中的压缩和解压缩功能。
//
// 该文件提供了 GZIP 和 ZLIB 格式的内存压缩和流式压缩功能。
// 支持字节数组、字符串和流式数据的压缩与解压缩操作。
//
// 主要功能：
//   - GZIP 内存压缩：字节数组和字符串的压缩解压
//   - GZIP 流式压缩：支持 io.Reader 和 io.Writer 接口
//   - ZLIB 内存压缩：字节数组和字符串的压缩解压
//   - ZLIB 流式压缩：支持 io.Reader 和 io.Writer 接口
//   - 支持自定义压缩等级
//
// 使用示例：
//
//	// GZIP 压缩字符串
//	compressed, err := comprx.GzipString("hello world")
//
//	// ZLIB 解压字节数据
//	decompressed, err := comprx.UnzlibBytes(compressedData)
package comprx

import (
	"io"

	"gitee.com/MM-Q/comprx/internal/cxgzip"
	"gitee.com/MM-Q/comprx/internal/cxzlib"
)

// ==================== gzip内存压缩API ====================

// GzipBytes 压缩字节数据（使用默认压缩等级）
//
// 参数:
//   - data: 要压缩的字节数据
//
// 返回:
//   - []byte: 压缩后的数据
//   - error: 错误信息
//
// 使用示例:
//
//	compressed, err := GzipBytes([]byte("hello world"))
func GzipBytes(data []byte) ([]byte, error) {
	return cxgzip.CompressBytes(data, CompressionLevelDefault)
}

// GzipBytesWithLevel 压缩字节数据（指定压缩等级）
//
// 参数:
//   - data: 要压缩的字节数据
//   - level: 压缩级别
//
// 返回:
//   - []byte: 压缩后的数据
//   - error: 错误信息
//
// 使用示例:
//
//	compressed, err := GzipBytesWithLevel([]byte("hello world"), CompressionLevelBest)
func GzipBytesWithLevel(data []byte, level CompressionLevel) ([]byte, error) {
	return cxgzip.CompressBytes(data, level)
}

// UngzipBytes 解压字节数据
//
// 参数:
//   - compressedData: 压缩的字节数据
//
// 返回:
//   - []byte: 解压后的数据
//   - error: 错误信息
//
// 使用示例:
//
//	decompressed, err := UngzipBytes(compressedData)
func UngzipBytes(compressedData []byte) ([]byte, error) {
	return cxgzip.DecompressBytes(compressedData)
}

// GzipString 压缩字符串（使用默认压缩等级）
//
// 参数:
//   - text: 要压缩的字符串
//
// 返回:
//   - []byte: 压缩后的数据
//   - error: 错误信息
//
// 使用示例:
//
//	compressed, err := GzipString("hello world")
func GzipString(text string) ([]byte, error) {
	return cxgzip.CompressString(text, CompressionLevelDefault)
}

// GzipStringWithLevel 压缩字符串（指定压缩等级）
//
// 参数:
//   - text: 要压缩的字符串
//   - level: 压缩级别
//
// 返回:
//   - []byte: 压缩后的数据
//   - error: 错误信息
//
// 使用示例:
//
//	compressed, err := GzipStringWithLevel("hello world", CompressionLevelBest)
func GzipStringWithLevel(text string, level CompressionLevel) ([]byte, error) {
	return cxgzip.CompressString(text, level)
}

// UngzipString 解压为字符串
//
// 参数:
//   - compressedData: 压缩的字节数据
//
// 返回:
//   - string: 解压后的字符串
//   - error: 错误信息
//
// 使用示例:
//
//	text, err := UngzipString(compressedData)
func UngzipString(compressedData []byte) (string, error) {
	return cxgzip.DecompressString(compressedData)
}

// ==================== gzip流式压缩API ====================

// GzipStream 流式压缩数据（使用默认压缩等级）
//
// 参数:
//   - dst: 目标写入器
//   - src: 源读取器
//
// 返回:
//   - error: 错误信息
//
// 使用示例:
//
//	file, _ := os.Open("input.txt")
//	defer file.Close()
//
//	var buf bytes.Buffer
//	err := GzipStream(&buf, file)
func GzipStream(dst io.Writer, src io.Reader) error {
	return cxgzip.CompressStream(dst, src, CompressionLevelDefault)
}

// GzipStreamWithLevel 流式压缩数据（指定压缩等级）
//
// 参数:
//   - dst: 目标写入器
//   - src: 源读取器
//   - level: 压缩级别
//
// 返回:
//   - error: 错误信息
//
// 使用示例:
//
//	file, _ := os.Open("input.txt")
//	defer file.Close()
//
//	output, _ := os.Create("output.gz")
//	defer output.Close()
//
//	err := GzipStreamWithLevel(output, file, CompressionLevelBest)
func GzipStreamWithLevel(dst io.Writer, src io.Reader, level CompressionLevel) error {
	return cxgzip.CompressStream(dst, src, level)
}

// UngzipStream 流式解压数据
//
// 参数:
//   - dst: 目标写入器
//   - src: 源读取器（压缩数据）
//
// 返回:
//   - error: 错误信息
//
// 使用示例:
//
//	compressedFile, _ := os.Open("input.gz")
//	defer compressedFile.Close()
//
//	output, _ := os.Create("output.txt")
//	defer output.Close()
//
//	err := UngzipStream(output, compressedFile)
func UngzipStream(dst io.Writer, src io.Reader) error {
	return cxgzip.DecompressStream(dst, src)
}

// ==================== ZLIB 内存压缩API ====================

// ZlibBytes 压缩字节数据（使用默认压缩等级）
//
// 参数:
//   - data: 要压缩的字节数据
//
// 返回:
//   - []byte: 压缩后的数据
//   - error: 错误信息
//
// 使用示例:
//
//	compressed, err := ZlibBytes([]byte("hello world"))
func ZlibBytes(data []byte) ([]byte, error) {
	return cxzlib.CompressBytes(data, CompressionLevelDefault)
}

// ZlibBytesWithLevel 压缩字节数据（指定压缩等级）
//
// 参数:
//   - data: 要压缩的字节数据
//   - level: 压缩级别
//
// 返回:
//   - []byte: 压缩后的数据
//   - error: 错误信息
//
// 使用示例:
//
//	compressed, err := ZlibBytesWithLevel([]byte("hello world"), CompressionLevelBest)
func ZlibBytesWithLevel(data []byte, level CompressionLevel) ([]byte, error) {
	return cxzlib.CompressBytes(data, level)
}

// UnzlibBytes 解压字节数据
//
// 参数:
//   - compressedData: 压缩的字节数据
//
// 返回:
//   - []byte: 解压后的数据
//   - error: 错误信息
//
// 使用示例:
//
//	decompressed, err := UnzlibBytes(compressedData)
func UnzlibBytes(compressedData []byte) ([]byte, error) {
	return cxzlib.DecompressBytes(compressedData)
}

// ZlibString 压缩字符串（使用默认压缩等级）
//
// 参数:
//   - text: 要压缩的字符串
//
// 返回:
//   - []byte: 压缩后的数据
//   - error: 错误信息
//
// 使用示例:
//
//	compressed, err := ZlibString("hello world")
func ZlibString(text string) ([]byte, error) {
	return cxzlib.CompressString(text, CompressionLevelDefault)
}

// ZlibStringWithLevel 压缩字符串（指定压缩等级）
//
// 参数:
//   - text: 要压缩的字符串
//   - level: 压缩级别
//
// 返回:
//   - []byte: 压缩后的数据
//   - error: 错误信息
//
// 使用示例:
//
//	compressed, err := ZlibStringWithLevel("hello world", CompressionLevelBest)
func ZlibStringWithLevel(text string, level CompressionLevel) ([]byte, error) {
	return cxzlib.CompressString(text, level)
}

// UnzlibString 解压为字符串
//
// 参数:
//   - compressedData: 压缩的字节数据
//
// 返回:
//   - string: 解压后的字符串
//   - error: 错误信息
//
// 使用示例:
//
//	text, err := UnzlibString(compressedData)
func UnzlibString(compressedData []byte) (string, error) {
	return cxzlib.DecompressString(compressedData)
}

// ==================== ZLIB 流式压缩API ====================

// ZlibStream 流式压缩数据（使用默认压缩等级）
//
// 参数:
//   - dst: 目标写入器
//   - src: 源读取器
//
// 返回:
//   - error: 错误信息
//
// 使用示例:
//
//	file, _ := os.Open("input.txt")
//	defer file.Close()
//
//	var buf bytes.Buffer
//	err := ZlibStream(&buf, file)
func ZlibStream(dst io.Writer, src io.Reader) error {
	return cxzlib.CompressStream(dst, src, CompressionLevelDefault)
}

// ZlibStreamWithLevel 流式压缩数据（指定压缩等级）
//
// 参数:
//   - dst: 目标写入器
//   - src: 源读取器
//   - level: 压缩级别
//
// 返回:
//   - error: 错误信息
//
// 使用示例:
//
//	file, _ := os.Open("input.txt")
//	defer file.Close()
//
//	output, _ := os.Create("output.zlib")
//	defer output.Close()
//
//	err := ZlibStreamWithLevel(output, file, CompressionLevelBest)
func ZlibStreamWithLevel(dst io.Writer, src io.Reader, level CompressionLevel) error {
	return cxzlib.CompressStream(dst, src, level)
}

// UnzlibStream 流式解压数据
//
// 参数:
//   - dst: 目标写入器
//   - src: 源读取器（压缩数据）
//
// 返回:
//   - error: 错误信息
//
// 使用示例:
//
//	compressedFile, _ := os.Open("input.zlib")
//	defer compressedFile.Close()
//
//	output, _ := os.Create("output.txt")
//	defer output.Close()
//
//	err := UnzlibStream(output, compressedFile)
func UnzlibStream(dst io.Writer, src io.Reader) error {
	return cxzlib.DecompressStream(dst, src)
}
