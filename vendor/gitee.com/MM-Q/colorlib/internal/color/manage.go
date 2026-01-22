// Package color 提供了颜色管理功能。
// 该文件实现了 ColorManager 结构体，用于管理 ANSI 颜色代码和颜色名称之间的映射关系，
// 支持标准颜色和亮色系列，提供颜色代码与名称的双向转换功能。
package color

// ColorManager 颜色管理
type ColorManager struct {
	codeToNameMap map[int]string // 颜色代码到名称
	nameToCodeMap map[string]int // 颜色名称到代码
}

// NewColorManager 创建颜色管理器
//
// 返回值:
//   - *ColorManager: 颜色管理器实例
func NewColorManager() *ColorManager {
	// 创建颜色管理器
	cm := &ColorManager{
		codeToNameMap: codeToNameMap, // 颜色代码到名称
		nameToCodeMap: nameToCodeMap, // 颜色名称到代码
	}

	return cm
}

const (
	// 标准颜色 (30-37)
	Black   = 30 // 黑色
	Red     = 31 // 红色
	Green   = 32 // 绿色
	Yellow  = 33 // 黄色
	Blue    = 34 // 蓝色
	Magenta = 35 // 品红色
	Cyan    = 36 // 青色
	White   = 37 // 白色
	Gray    = 90 // 灰色

	// 亮色 (90-97) - 统一用 Bright 前缀
	BrightRed     = 91 // 亮红色
	BrightGreen   = 92 // 亮绿色
	BrightYellow  = 93 // 亮黄色
	BrightBlue    = 94 // 亮蓝色
	BrightMagenta = 95 // 亮品红色
	BrightCyan    = 96 // 亮青色
	BrightWhite   = 97 // 亮白色
)

// 在 color/manage.go 中添加
var colorCodeToStringMap = map[int]string{
	// 标准颜色 (30-37)
	Black:   "30", // 黑色
	Red:     "31", // 红色
	Green:   "32", // 绿色
	Yellow:  "33", // 黄色
	Blue:    "34", // 蓝色
	Magenta: "35", // 品红色
	Cyan:    "36", // 青色
	White:   "37", // 白色
	Gray:    "90", // 灰色

	// 亮色 (90-97)
	BrightRed:     "91", // 亮红色
	BrightGreen:   "92", // 亮绿色
	BrightYellow:  "93", // 亮黄色
	BrightBlue:    "94", // 亮蓝色
	BrightMagenta: "95", // 亮品红色
	BrightCyan:    "96", // 亮青色
	BrightWhite:   "97", // 亮白色
}

// 添加一个获取颜色代码字符串的方法
func (cm *ColorManager) GetColorCodeString(code int) (string, bool) {
	str, ok := colorCodeToStringMap[code]
	return str, ok
}

// 内部映射表 - 颜色代码到名称
var codeToNameMap = map[int]string{
	// 标准颜色
	Black:   "black",   // 黑色
	Red:     "red",     // 红色
	Green:   "green",   // 绿色
	Yellow:  "yellow",  // 黄色
	Blue:    "blue",    // 蓝色
	Magenta: "magenta", // 品红色
	Cyan:    "cyan",    // 青色
	White:   "white",   // 白色
	Gray:    "gray",    // 灰色

	// 亮色
	BrightRed:     "bright_red",     // 亮红色
	BrightGreen:   "bright_green",   // 亮绿色
	BrightYellow:  "bright_yellow",  // 亮黄色
	BrightBlue:    "bright_blue",    // 亮蓝色
	BrightMagenta: "bright_magenta", // 亮品红色
	BrightCyan:    "bright_cyan",    // 亮青色
	BrightWhite:   "bright_white",   // 亮白色
}

// 内部映射表 - 名称到颜色代码
var nameToCodeMap = map[string]int{
	// 标准颜色
	"black":   Black,   // 黑色
	"red":     Red,     // 红色
	"green":   Green,   // 绿色
	"yellow":  Yellow,  // 黄色
	"blue":    Blue,    // 蓝色
	"magenta": Magenta, // 品红色
	"cyan":    Cyan,    // 青色
	"white":   White,   // 白色
	"gray":    Gray,    // 灰色

	// 亮色
	"bright_red":     BrightRed,     // 亮红色
	"bright_green":   BrightGreen,   // 亮绿色
	"bright_yellow":  BrightYellow,  // 亮黄色
	"bright_blue":    BrightBlue,    // 亮蓝色
	"bright_magenta": BrightMagenta, // 亮品红色
	"bright_cyan":    BrightCyan,    // 亮青色
	"bright_white":   BrightWhite,   // 亮白色
}

// GetColorName 根据颜色代码获取颜色名称
//
// 参数:
//   - code: 颜色代码
//
// 返回值:
//   - name: 颜色名称
func (cm *ColorManager) GetColorName(code int) (string, bool) {
	name, ok := cm.codeToNameMap[code]
	return name, ok
}

// GetColorCode 根据颜色名称获取颜色代码
//
// 参数:
//   - name: 颜色名称
//
// 返回值:
//   - code: 颜色代码
func (cm *ColorManager) GetColorCode(name string) (int, bool) {
	code, ok := cm.nameToCodeMap[name]
	return code, ok
}

// IsColorName 判断是否为颜色名称
//
// 参数:
//   - name: 颜色名称
//
// 返回值:
//   - ok: 是否为颜色名称
func (cm *ColorManager) IsColorName(name string) bool {
	_, ok := cm.nameToCodeMap[name]
	return ok
}

// IsColorCode 判断是否为颜色代码
//
// 参数:
//   - code: 颜色代码
//
// 返回值:
//   - ok: 是否为颜色代码
func (cm *ColorManager) IsColorCode(code int) bool {
	_, ok := cm.codeToNameMap[code]
	return ok
}
