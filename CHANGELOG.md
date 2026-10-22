# CHANGELOG — tellurion-contracts

> 本仓的版本账本。**每个发版一节**，标题格式固定 `## [vX.Y.Z] - YYYY-MM-DD`（`tools/check_release_matrix` 按它和 README 兼容矩阵、`git tag` 三方对拍）。
> 未发版的改动记在 `## [Unreleased]` 下，发版时把那一节的内容搬进新版本节。
> 规则口径见 [`README.md`](README.md) 的「版本规则」（源头：`tellurion:docs/TASKS.md` QT-19 ③）。
> breaking 的条目必须写成「破坏性：<旧 → 新>」，光写「改了签名」不算留账。

## [Unreleased]

尚未发版（**零 tag**；QT-05 ③-① 落码进行中，已有 1 个库包）。首个版本将是 `v0.1.0`，由 `tellurion:docs/TASKS.md` QT-07 ③-S2-a 在 S2 打；在此之前使用方只能 pin `main`，**无兼容承诺**。

### 新增（MINOR 级：加包 / 加导出符号）

- **QT-05 ③-①（2026-10-22，进行中）**：新增 **`fees/`** 包（13 个导出符号，已登记进 `docs/API.md` §4）——自 `tellurion:pkg/fees` **逐字迁入**。搬运与改注释分两步、各自取证：① 拷贝落地时 `ashare.go` / `ashare_test.go` 与 core 侧 **sha256 逐字节相同**（证明搬运无损）；② 随后只改注释里的跨仓路径为 `repo:path` 形式，**代码行零差异**（用「剥掉全部注释行再 diff」验证，两个文件退出码均为 0）。落地后 `ashare.go` 比 core 原文多 10 行，多出的全是包注释。
  - 迁入面：7 个费率常量（`Default*` / `FixedSlippageRate`）、`StampTaxCutDate`、`StampTaxRateFor`（按成交日取档的**唯一**比较点）、`AShareFees` + `DefaultAShareFees` / `ApplyDefaults` / `Validate`。
  - **裁决 `docs/API.md` §8 #3**：整包迁入，包括三个「看起来像行为」的函数。依据是实测——`tellurion:pkg/fees` **零算术**，算钱在 `tellurion:pkg/backtest/tracker` 与 `tellurion:pkg/live/mock_trader.go`；把 `StampTaxRateFor` 留在 core 会拆散「削减日只在一处比较」这条不变量。
  - 第三方依赖：**零新增**（`fees` 的测试只用标准库）。本仓 `go.mod` 未变。
  - ⚠️ 登记在案的形状代价：`StampTaxCutDate` 是**导出 `var`**（`time.Time` 做不了 const），进程内可被导入方改写；本次按 QT-05 ④「不改任何业务逻辑」没动它，下游只读。

### 变更（PATCH 级：文档与仓内自查工具，不改导出面）

- **CT-02**（2026-10-22）：README 补「版本规则表 + 兼容矩阵 + 发版流程 + 拉取前提（`GOPRIVATE`）」；新建本 `CHANGELOG.md`；新建护栏 `tools/check_release_matrix`（矩阵 ↔ CHANGELOG ↔ tag 三件套不一致即红）并接进 CI。导出面零变化。
- **CT-01**（2026-10-21）：新建 `docs/API.md`（导出面清单：稳定性标注 + 错误语义 + 冻结面节号）与 `tools/check_api_surface`（AST 扫描，新增导出符号未登记即红），接进 CI 步骤「导出面核对」。导出面零变化 —— 登记区今天为空。
