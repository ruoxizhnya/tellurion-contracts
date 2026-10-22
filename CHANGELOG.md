# CHANGELOG — tellurion-contracts

> 本仓的版本账本。**每个发版一节**，标题格式固定 `## [vX.Y.Z] - YYYY-MM-DD`（`tools/check_release_matrix` 按它和 README 兼容矩阵、`git tag` 三方对拍）。
> 未发版的改动记在 `## [Unreleased]` 下，发版时把那一节的内容搬进新版本节。
> 规则口径见 [`README.md`](README.md) 的「版本规则」（源头：`tellurion:docs/TASKS.md` QT-19 ③）。
> breaking 的条目必须写成「破坏性：<旧 → 新>」，光写「改了签名」不算留账。

## [Unreleased]

尚未发版（**零 tag**；QT-05 ③-① 本仓侧落码完成，四个库包 / 429 个导出符号）。首个版本将是 `v0.1.0`，由 `tellurion:docs/TASKS.md` QT-07 ③-S2-a 在 S2 打；在此之前使用方只能 pin `main`，**无兼容承诺**。

### 新增（MINOR 级：加包 / 加导出符号）

- **QT-05 ③-①（2026-10-22，本仓侧已完成）**：新增 **`fees/`** 包（13 个导出符号，已登记进 `docs/API.md` §4）——自 `tellurion:pkg/fees` **迁入（代码逐字）**。搬运与改注释分两步、各自取证：① 拷贝落地时 `ashare.go` / `ashare_test.go` 与 core 侧 **sha256 逐字节相同**（证明搬运无损）；② 随后只改注释（跨仓路径改 `repo:path`；另订正 `Validate` 注释里与实现不符的区间 `(0, 0.05)` → `[0, 0.05)` 和「被 OrderManager 调用」的失实描述），**代码行零差异**（用「剥掉全部注释行再 diff」验证，两个文件退出码均为 0）。落地后 `ashare.go` 比 core 原文多 12 行，多出的全是注释（包注释 + `Validate` 区间/调用方描述的订正）。
  - 迁入面：7 个费率常量（`Default*` / `FixedSlippageRate`）、`StampTaxCutDate`、`StampTaxRateFor`（按成交日取档的**唯一**比较点）、`AShareFees` + `DefaultAShareFees` / `ApplyDefaults` / `Validate`。
  - **裁决 `docs/API.md` §8 #3**：整包迁入，包括三个「看起来像行为」的函数。依据是实测——`tellurion:pkg/fees` **零算术**，算钱在 `tellurion:pkg/backtest/tracker` 与 `tellurion:pkg/live/mock_trader.go`；把 `StampTaxRateFor` 留在 core 会拆散「削减日只在一处比较」这条不变量。
  - 第三方依赖：**零新增**（`fees` 的测试只用标准库）。本仓 `go.mod` 未变。
  - ⚠️ 登记在案的形状代价：`StampTaxCutDate` 是**导出 `var`**（`time.Time` 做不了 const），进程内可被导入方改写；本次按 QT-05 ④「不改任何业务逻辑」没动它，下游只读。

- **QT-05 ③-①（2026-10-22，本仓侧已完成）**：新增 **`domain/market/`** 包（46 个导出符号，已登记进 `docs/API.md` §4）——自 `tellurion:pkg/domain/market` **迁入（代码逐字）**。7 个文件落地时 sha256 逐字节相同；随后只改注释里的跨仓路径为 `repo:path` 形式，剥掉全部注释行再 diff，**七个文件退出码全 0**（代码行零差异）。
  - 迁入面：行情值对象（`OHLCV` / `Stock` / `Fundamental` / `FundamentalData` / `Split` / `Dividend` / `IndexDaily` / `IndexConstituent`）+ 频率词表（`Frequency` / `FrequencySpec`）+ 时间语义与 PIT 原语（`TimeSemantics` / `AsiaShanghai` / `DailyVisibleAt` / `IntradayVisibleAt` / `OHLCV.VisibleBy`）+ `Provider` 接口（5 个方法）。
  - 稳定性：字面量词表（`Frequency` / `FrequencySpec` / `TimeSemantics` 及各自常量）标 **frozen**——取值与 nt-sidecar `sidecar/mapping.py` 逐字对齐，改字面量 = 同时改 Python 侧；`OHLCV` / `Fundamental` 的字段集也标 **frozen**（**SPEC 互锁**：core 的 `pkg/domain/spec_drift_test.go` 与 `tellurion:docs/SPEC.md` 双向比对，那条测试**刻意留 core**，core 侧改依赖后改为 import 本包）。
  - 第三方依赖：**零新增**（本包只用标准库，测试也不用 testify）。本仓 `go.mod` 仍只有 module 行。

- **QT-05 ③-①（2026-10-22，本仓侧已完成）**：新增 **`domain/`** 包（176 个导出符号，已登记进 `docs/API.md` §4）——自 `tellurion:pkg/domain` 迁入（23 文件：16 非测试 + 7 测试）。**这不是逐字等价迁**，刻意改动只有三类，且逐项留痕：
  - ① import 换成本仓路径（`domain/market`、`fees`）；
  - ② **删掉 5 个零引用的部署配置类型**：`Config` / `DatabaseConfig` / `RedisConfig` / `ServiceConfig` / `TushareConfig`——core 与 agent 实测 **0 引用**（连 core 自己都不用），`Config` 还带 `Password` 字段。裁决记录：[`docs/adr/0002-config-族的边界处置.md`](docs/adr/0002-config-族的边界处置.md)；core 侧随切换依赖那次 commit 一并删除（0 调用点）。**破坏性：`domain.Config` / `domain.DatabaseConfig` / `domain.RedisConfig` / `domain.ServiceConfig` / `domain.TushareConfig` → 移除**（本仓不提供；core 侧删除后彻底消失——今天无使用方，故无实际迁移成本）。
  - ③ 随之删掉只为它们存在、且断言恒真的 `TestConfig_Types`。
  - **等价性证据**：剥注释 + 归一化 import 后，20/23 文件与 core 完全相同；`execution_test.go` 只差 gofmt 的 import 顺序；`types.go` / `types_test.go` 只差被删的那两段（−31 / −24 行）。
  - 稳定性：领域词汇（信号/订单/持仓/组合/公司行为/因子名册/回测结果）标 stable；core-only 数据族标 experimental + 进 QT-05 ③-⑨ 回收清单（2026-10-22 重测：119 个根名，agent import 38 个，0 命中 81）。
  - **第三方依赖**：`github.com/stretchr/testify v1.11.1`（**仅测试**）——`domain` 的 `_test.go` 随迁带入，版本与 core 对齐；`go.mod` 从此有 `require`（另有 3 个 testify 的间接依赖）。这不是新增决策：AGENTS §4 的第三方表在 C1/C2 时已预登记。
  - ⚠️ 记录在案：`FactorDefinition.Validate` / `FactorRegistration.Validate` / `FactorMiningRecord.Validate` **会返回错误**（fail-closed），但它们是包内 `fmt.Errorf`，**不是** `AppError`（本包不依赖 kit，也不依赖尚不存在的三类哨兵）。消费方按「`error != nil` ⇒ 拒绝登记 / 拒绝入库」处理，别按错误类型分支。

- **QT-05 ③-①（2026-10-22，本仓侧收官）**：新增 **`engine/`** 包（**包名仍是 `contracts`**，194 个导出符号，已登记进 `docs/API.md` §4）——自 `tellurion:pkg/backtest/contracts` **逐字迁入**（20 个文件：10 非测试 + 10 测试）。刻意的代码改动**只有 import 路径**（`domain` / `fees` / 测试包自引用）：剥注释 + 归一化 import 后 17/20 文件完全相同，3 个测试文件只差 gofmt 的 import 顺序，**0 处真实代码差异**。
  - 迁入面：四契约（`BacktestRequest` / `BacktestResponse` / `SignalWeight` / `EngineRunner` + `BacktestRunner`）、内核注册与市场路由、成交口径（`FillTiming` + 三个 `ErrFillTiming*`）、canonical 编码与 provenance 三 hash、`CheckValidity` 与无效 Code、**三类哨兵与 abi 常量**（同包）。
  - `TradingConfig` / `PriceLimitConfig` 按 [`adr/0002`](docs/adr/0002-config-族的边界处置.md) ③「冻结面伴生名」保留原名。
  - **裁决 `docs/API.md` §8 #1**：`errors/` 与 `abi/` **不独立成包**（[`adr/0003`](docs/adr/0003-哨兵与-abi-不独立成包.md)）——实测 `contracts.Err*` 只有 10 个消费文件、`contracts.Host*` 只有 4 个，且每个文件都在同时用同包的其它符号，拆包没有调用点能少 import 一个包。代价如实登记：`engine/` 194 个符号、占本仓 429 的 45%。
  - **留在 core**：`resultspec_test.go`（读 core 的 `docs/SPEC.md` 做文档互锁；core 侧切换依赖后改为 import 本包）。
  - 第三方依赖：零新增（本包测试沿用已登记的 testify）。

### 变更（PATCH 级：文档与仓内自查工具，不改导出面）

- **QT-05 ③-① 自审**（2026-10-22）：订正文档失真——ADR-0002 文件名死链、`FactorDefinition` / `FactorRegistration` / `FactorMiningRecord` 的方法类别与校验描述、`EngineRunner.RunBacktest` 返回值（指针）、`ParseStrategyRef` 的小写约束、`fees.AShareFees.Validate` 注释里的区间与调用方描述；新增 [`adr/0004`](docs/adr/0004-值对象上的纯函数行为随类型走.md)（「零实现」的四条判据）并补进 AGENTS / ARCHITECTURE / API 索引。代码改动**仅注释**（`doc.go` / `engine/contracts.go` / `domain/types.go` / `domain/execution.go` / `fees/ashare.go` / `fees/ashare_test.go`），导出面零变化（429/429）。
- **CT-02**（2026-10-22）：README 补「版本规则表 + 兼容矩阵 + 发版流程 + 拉取前提（`GOPRIVATE`）」；新建本 `CHANGELOG.md`；新建护栏 `tools/check_release_matrix`（矩阵 ↔ CHANGELOG ↔ tag 三件套不一致即红）并接进 CI。导出面零变化。
- **CT-01**（2026-10-21）：新建 `docs/API.md`（导出面清单：稳定性标注 + 错误语义 + 冻结面节号）与 `tools/check_api_surface`（AST 扫描，新增导出符号未登记即红），接进 CI 步骤「导出面核对」。导出面零变化 —— 登记区今天为空。
