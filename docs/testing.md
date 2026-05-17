# fastlog 测试计划

> 规划时间：2026-05-15

---

## 测试策略

按依赖关系从底层到上层逐层测试，先测无副作用的纯函数，再测需要 mock 的组件。

---

## 1. `fastlog.go` — Level 体系 | P0 先测

**测试要点**：
- `Level.String()` — 6 个级别返回各自名称
- `Level.Enabled()` — **历史 Bug 回归**：INFO 级别下，WARN/ERROR 放行，DEBUG 抑制
- `ParseLevel()` — 大小写不敏感解析正常字符串，非法字符串返回 INFO+error
- `AllLevels()` — 返回 6 个级别且顺序正确

---

## 2. `field.go` — Field 字段系统 | P0 次之

**测试要点**：
- 12 个构造函数返回的 Key/Type/Value 正确（一个类型造一个典型值即可）
- `Error(nil)` → StringVal="<nil>"
- 9 个取值方法在类型匹配时返回值正确，类型不匹配时返回零值
- `Stack()` 返回 Key="stack" 且值非空

---

## 3. `sampler.go` — Sampler 采样器 | P1 中段

**测试要点**：
- 前 N 条放行、第 N+1 条抑制（initial=3）
- 窗口过期后计数器重置，重新放行（tick=50ms，sleep 60ms 验证）
- 不同 level 独立计数（DEBUG 抑制不影响 INFO）
- `thereafter=0` 边界：前 N 条放行后永久抑制

---

## 4. `formatter.go` — 5 种格式化器 | P1 中段

**测试要点**：
- 每种格式输出模板正确（带/不带 Caller、带/不带 Fields）
- JSON 末尾有 `\n`，且是合法 JSON
- 空 Entry（Message="" 或 Fields=nil）不 panic

---

## 5. `writer.go` — 写入器 | P2 后段

**测试要点**：
- ColorWriter：NoColor=true 原样输出；NoColor=false 检测级别后着色输出；未识别级别原样输出
- MultiWriter：所有 writer 都收到数据；一个 writer 失败返回 error；Close 关闭所有；空列表不 panic

---

## 6. `logger.go` — Logger 核心 | P3 最后

**测试要点**：
- 6 个 Option 函数各自生效
- 级别过滤正确（INFO 级别下 Debug 不写入，Warn 写入）
- 采样器集成后被抑制不走 writer
- Sync/Close 委托给底层 writer

---

## 执行顺序

```
P0 ─── fastlog.go → field.go      ← 纯函数，无依赖
P1 ─── sampler.go → formatter.go  ← 核心逻辑
P2 ─── writer.go               ← 需要 mock Writer
P3 ─── logger.go               ← 最复杂，依赖所有组件
```

## 测试工具

- 标准库 `testing`，表格驱动测试
- `bytes.Buffer` 做 mock writer
- 不引入第三方测试框架
