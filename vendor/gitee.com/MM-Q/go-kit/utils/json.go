package utils

// hexTable 用于将字节转换为 JSON 转义序列中的十六进制字符
var hexTable = "0123456789abcdef"

// needsEsc 判断是否需要转义 JSON 字符串中的字符
//
// 参数：
//   - c: 待判断的字符
//
// 返回：
//   - 是否需要转义
func needsEsc(c byte) bool {
	return c < 0x20 || c == '"' || c == '\\'
}

// QuoteBytes 将输入字节切片转义为合法 JSON 字符串字面量。
//
// 转义规则：
//  1. 7 个缩写控制字符 => \" \\ \b \f \n \r \t
//  2. 其余 0x00–0x1F 统一写成 \u00XX
//  3. 无转义时直接原串返回，零额外分配
//
// 参数：
//   - raw: 待转义的原始字节切片
//
// 返回：
//   - 转义后的 JSON 字节串
func QuoteBytes(raw []byte) []byte {
	// 1. 先统计需转义字符数量，同时算最大可能长度
	var cnt int
	for _, c := range raw {
		if needsEsc(c) {
			cnt++
		}
	}
	if cnt == 0 {
		// 无转义可直接返回原切片（只读场景安全）
		return raw
	}

	// 2. 一次性分配足够大缓冲区
	out := make([]byte, 0, len(raw)+cnt*6) // 最坏每字符+6
	for _, c := range raw {
		switch c {
		case '"':
			out = append(out, '\\', '"')
		case '\\':
			out = append(out, '\\', '\\')
		case '\b':
			out = append(out, '\\', 'b')
		case '\f':
			out = append(out, '\\', 'f')
		case '\n':
			out = append(out, '\\', 'n')
		case '\r':
			out = append(out, '\\', 'r')
		case '\t':
			out = append(out, '\\', 't')
		default:
			if c < 0x20 {
				out = append(out, '\\', 'u', '0', '0', hexTable[c>>4], hexTable[c&0xF])
			} else {
				out = append(out, c)
			}
		}
	}

	return out
}

// QuoteString 将输入字符串转义为合法 JSON 字符串字面量
//
// 转义规则：
//  1. 7 个缩写控制字符 => \" \\ \b \f \n \r \t
//  2. 其余 0x00–0x1F 统一写成 \u00XX
//  3. 无转义时直接原串返回，零额外分配
//
// 参数：
//   - raw: 待转义的原始字符串
//
// 返回：
//   - 转义后的 JSON 字符串
func QuoteString(raw string) string {
	if raw == "" {
		return ""
	}

	// 复用 []byte 路径，安全但多一次拷贝
	return string(QuoteBytes([]byte(raw)))
}
