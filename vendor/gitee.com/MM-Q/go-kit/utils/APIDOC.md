# Package utils

`import "gitee.com/MM-Q/go-kit/utils"`

工具类封装，提供 JSON 转义、字节格式化等常用辅助函数。

## FUNCTIONS

### func QuoteString

```go
func QuoteString(raw string) string
```

QuoteString 将输入字符串转义为合法 JSON 字符串字面量。

**参数：**

- `raw` – 待转义的原始字符串

**返回：**

- 转义后的 JSON 字符串

---

### func QuoteBytes

```go
func QuoteBytes(raw []byte) []byte
```

QuoteBytes 将输入字节切片转义为合法 JSON 字符串字面量。

**转义规则：**

1. 7 个缩写控制字符 => `\" \\ \b \f \n \r \t`  
2. 其余 `0x00–0x1F` 统一写成 `\u00XX`  
3. 无转义时直接原串返回，零额外分配

**参数：**

- `raw` – 待转义的原始字节切片

**返回：**

- 转义后的 JSON 字节串

---

### func FormatBytes

```go
func FormatBytes(bytes int64) string
```

FormatBytes 将字节数转换为人类可读的带单位的字符串，用于将字节数格式化为易读的存储单位格式，支持 B 到 PB 的转换。

**参数：**

- `bytes` – 字节数（`int64` 类型）

**返回：**

- `string` – 格式化后的字符串，如 `"1.23 KB"`、`"456.78 MB"`、`"2.34 GB"` 等