# CHANGELOG — tellurion-contracts

> 本仓的版本账本。**每个发版一节**，标题格式固定 `## [vX.Y.Z] - YYYY-MM-DD`（`tools/check_release_matrix` 按它和 README 兼容矩阵、`git tag` 三方对拍）。
> 未发版的改动记在 `## [Unreleased]` 下，发版时把那一节的内容搬进新版本节。
> 规则口径见 [`README.md`](README.md) 的「版本规则」（源头：`tellurion:docs/TASKS.md` QT-19 ③）。
> breaking 的条目必须写成「破坏性：<旧 → 新>」，光写「改了签名」不算留账。

## [Unreleased]

（暂无）

## [v0.5.0] - 2026-10-24

MINOR 发版：QT-22 把 `version_hash` 的 canonical 集合收进契约（**只增不改**），导出面 443 → 447
（+4：1 个 frozen `const` + 1 个 `type` + 2 个方法），冻结基线 426 → 430 行（`git diff --numstat` 实测 +4 / −0）。
无破坏性变更；既有六条域标签与 `run_spec_hash` / `dataset_snapshot_hash` 的编码一字未动。

### 新增（MINOR 级：加导出符号）

- **QT-22 契约腿（2026-10-24）**：`engine/provenance.go` 增 `StrategyVersionInput`（六字段
  `StrategyID` / `Carrier` / `ABIVersion` / `Scope` / `Entry` / `Source`）+ `Hash()` / `MustHash()`；
  `engine/canonical.go` 增 `DomainStrategyVersionHash = "tellurion.strategy_version_hash.v1"`。
  - **集合已裁死，不许就地改**：`strategy_id + carrier + abi_version + scope + entry + source`。
    **不含** `capability`（只参与市场路由，不改变源码语义）；**不含** `params_defaults`
    （资产服务今天没有该字段，零生产者；要进必须升 `v2` 域标签并走双标签过渡）。
  - **不复用 `DomainParamsHash`**：`params_hash` 是生效参数指纹，`version_hash` 是源码版本身份；
    同域会让「只比参数」的查询误命中换了源码版本的资产。
  - **与 `tellurion-strategy` 对拍**：`StrategyVersionInput` 与
    `tellurion-strategy:internal/asset/hash.go` 的 `VersionInput` **逐字段、逐顺序一致**；
    `TestStrategyVersionInput_MatchesStrategyRepoEncoding` 逐行重建 strategy 侧编码并断言同 hex，
    `TestGoldenVector_StrategyVersionHash` 钉住固定向量
    `0e6ad53637ce43742804fe325d97b6799372fe03319f1020fc26f0d748d68412`（输入：`s1` / `wasm` /
    `on_bar/1` / `symbol` / `OnBar` / `x`）。
  - **未随本单落（如实）**：`tellurion-strategy` 侧把本地 `DomainStrategyVersionHash` / `VersionInput`
    改为本仓别名（资产服务收编腿归 G4 / G13 通知腿）；`params_defaults` 的 `v2` 立项条件见
    [`docs/design/qt22-version-hash-decision-material.md`](docs/design/qt22-version-hash-decision-material.md) §3。
    本仓只交付 canonical 定义、向量与登记，**不改 core/strategy 源码**。

## [v0.4.0] - 2026-10-24

MINOR 发版：两笔**只增不改**的契约腿（QT-20 参数投递入口 + QT-41b ④ 编译 / 资产错误族），导出面 439 → 443（+4 哨兵 `var`；两个结构体各 +1 可选字段，字段不计符号数），冻结基线 422 → 426 行（`git diff --numstat` 实测 +4 / −2：两条 struct 行因字段增加而改写）。无破坏性变更；`SignalWeight` / `ResultSpec` 冻结面一字未动。

### 新增（MINOR 级：加导出符号 / 加可选字段）

- **QT-20 契约腿（2026-10-24）**：`engine/contracts.go` 的 `BacktestRequest` 增两个可选字段——
  - `StrategyRef string`（`json:"strategy_ref,omitempty"`）：wire 形态与 `engine/provenance.go` 的 `ParseStrategyRef` / `StrategyRef.String()` **互逆**（`id` 或 `id@hash`），**不新增第二种引用形状**；`ValidateRunSpec` 只做形状校验（`errors.Is(err, ErrStrategyRefMalformed)` 可判），「该 id 是否存在、解析成哪个 hash」属资产服务查表。
  - `StrategyParams map[string]any`（`json:"strategy_params,omitempty"`）：**浅层对象**，值必须是标量（string / number / bool / null）；嵌套对象或数组在 `ValidateRunSpec` 拒绝并点名键。**未知键清单不复制进本仓**（J2：那属于 core 的策略 manifest，契约层只做形状 / 类型校验；合并 defaults ⊕ 覆写与取值域归执行侧）。
  - `domain/backtest.go` 的 `BacktestParams` 增 `Params map[string]any`（`json:"params,omitempty"`）：**实际生效参数**，字段名对齐 core `RunMeta.Params`；与 `FillTiming` / `Kernel` 同一条纪律——记录执行侧解析后的值，引擎不猜。
  - 存量行为不变：两个字段带 `omitempty`，`TestRunSpec_LegacyRequestSerializationUnchanged` 的逐字节 golden 仍绿。
  - **未随本单落（如实）**：请求字段的消费（生效参数合并、未知键 400、`RunMeta.Params` 投递、各入口接线）在 core **G1-QT-20**；本仓只交付形状与形状校验。
- **QT-41b ④ 错误族契约腿（2026-10-24）**：`engine/execution_errors.go` 增四个跨仓可识别哨兵——`ErrCompileFailed` / `ErrCompileUnavailable` / `ErrSourceTooLarge` / `ErrArtifactTooLarge`，名字与语义（含错误文案）逐字取自 `tellurion-strategy:internal/asset/errors.go`，`errors.Is` 可判、`%w` 保链、四条互不命中。迁入判据只有一条：**跨仓消费方必须能识别**（agent 生成闭环要把「编译失败」与「服务没配好」分开、core 取 artifact 要分「fetch 失败 / 超限 / 没有这个版本」、orchestration 要判可重试性），且不能为此 import 资产服务的实现包。**私有哨兵不迁**（`ErrStrategyNotFound` / `ErrStrategyExists` / `ErrVersionNotFound` / `ErrInvalidTransition` / `ErrJobNotFound` 等只看 HTTP status + wire code）。
  - **过渡态（如实）**：`tellurion-strategy:internal/asset` 今天仍自己定义这四个值；该仓升 pin 后改为本包别名（G13 的通知腿）。两边同名同义，不是两份真相。

## [v0.3.0] - 2026-10-23

MINOR 发版：QT-10 的契约腿（策略执行错误分级，**只增不改**），导出面 435 → 439（+4：3 个哨兵 `var` + 1 个 **frozen** `const`），冻结基线同步重生成（418 → 422 行，`git diff --numstat` 实测 **+4 / −0**）。无破坏性变更。

### 新增（MINOR 级：加导出符号）

- **QT-10 ③-① 契约腿（2026-10-23）**：新增 `engine/execution_errors.go` 的三类执行错误哨兵 **`ErrDataGap` / `ErrStrategyFault` / `ErrContractBreach`**（3 个 `var`，均 stable，`errors.Is` 可判 + `%w` 保链），与 `engine/validity.go` 的稳定 Code **`CodeDataGapSkipped`**（`const`，**frozen**，`data_gap_skipped`）。语义：数据缺口 → 跳过当天 + 写 `invalid_reasons`、run 继续、终态 `invalid`；策略故障（wasm trap / 非法信号 / 预算超限 / artifact 校验失败）与契约违规（前视 / 未知参数 / abi 不符）→ **终止 run**、`status=failed`。落点在本仓的理由是生产者（provider / wasm 宿主）与消费者（引擎 / job / AI 管线）分属不同模块，哨兵必须在两侧之间可 import 且不成环（`errors/` 不独立成包按 [`adr/0003`](docs/adr/0003-哨兵与-abi-不独立成包.md)）。
  - ⚠️ **`status` 仍是冻结三值**（`completed` / `invalid` / `failed`）：本单只落**映射规则**与一个新 Code，**没有也不引入第四值**（QT-10 ④、QT-63 ④ 两处明文）。钉子：`TestResultSpecStatusSetStaysThreeValued`；另有一条 `TestCheckValidityNeverEmitsDataGapSkipped` 证明新 Code 走既有通道、不由纯函数伪造产出。
  - **版本取 MINOR（`v0.3.0`）而非派工原文的 `v0.2.1`，理由登记在此**：本仓 README 版本规则表明写「新增导出符号 = MINOR」，同型先例是上面的 `v0.2.0`（+6 符号记 MINOR）；PATCH 那一行写的是「文档 / 注释 / 无行为修正」，套不上。派工另要求「兼容矩阵加 `v0.2.1` 行」——**实测不可执行**：矩阵按版本**序列** `vX.Y` 建行（QT-19 ③-③），注入那一行后 `go run ./tools/check_release_matrix` 当场红「矩阵里序列 v0.2 有多行」（EXIT=1，已逐字节还原）。两条合起来按规则表落到 `v0.3.0` + 新增 `v0.3` 序列行。护栏不核 bump 档位，所以这一条只能靠账来拦。
  - **未随本单落（如实）**：① **归类腿**（QT-10 ③-②：底层错误往哪一类包）与**状态映射腿**（③-③：`RunBacktest` 在哪条分支 `EndRun`）在 **core**，今天 core 的 provider 仍可能返回未分类 error——`docs/API.md` §3 与 `docs/api-schema.md` §1.3 的注已按此改成「对生产者是验收标准、对消费方才是现状」；② `session_terminated` 今天定义在 core `pkg/sandbox/wasm/replay.go:42`（QT-15 落沙箱层时自带），**未迁本仓**，收编与否归 core 侧裁决（登记见 `docs/API.md` §5.3）。

## [v0.2.0] - 2026-10-23

MINOR 发版：两个只增不改的契约腿随同一 tag 分发，导出面 429 → 435（+6：1 个 frozen 常量 + 2 类型 + 1 方法 + 2 哨兵错误），冻结基线同步重生成（412 → 418 行）。无破坏性变更。

### 新增（MINOR 级：加导出符号）

- **QT-16 契约腿（2026-10-23）**：新增 `engine/DomainRunIdentityHash`（const，**frozen**）——检查点身份复合域标签 `"tellurion.run_identity_hash.v1"`，覆盖 `{run_universe_hash, scope, shard_key, param_effective_hash}` 四分量（p0-round2-review N3 冻结设计 + w1-repo-ownership §2.5 裁决：常量住 contracts，计算函数留 core `pkg/strategy`，消费方只有 core）。
- **QT-62 契约腿（2026-10-23）**：新增宿主区 × 模块静态数据重叠判定族（5 个符号，全部 stable）——`StaticDataSegment` / `StaticDataImage` / `HostRegionDeclaration.CheckAgainstStaticData` / `ErrHostRegionOverlapsData` / `ErrHostRegionDataUnprovable`。权威定则：模块静态数据 = Data section（id=11）活动段并集、宿主区 = 三个导出 global 的声明，两者都出自 artifact 自身 ⇒ 相交即拒绝，不重排地址、不缩区间；`Unprovable` 非空必须 fail-closed（读不懂 ≠ 没风险）。零长段不算重叠，区间算术 uint64 不绕回。QT-01 的冻结面（符号名 / i32 类型 / 下界 / 对齐 / 上限）一字未动。

## [v0.1.1] - 2026-10-23

QT-09 的随仓交付物随 tag 分发（本节内容自 `[Unreleased]` 搬入）。**版本取 PATCH（`v0.1.0` → `v0.1.1`）**：本批全部是仓内自查工具 / 示例 / 语料，**导出面零变化**（仍 429/429，冻结基线 412 行），下游 import 的任何符号的形状与语义都没动。`conformance/` 虽是新包，但刻意零导出符号（见 `conformance/cases.go` 头注释），它不改变兼容面 —— 它改变的是**发布物**：`v0.1.0` 的 tag 树里 `conformance/` 命中 0 个文件，而 core（`tellurion:pkg/conformance`）与 agent（`tellurion-agent:internal/conformance`）两侧的腿按 pin 的版本读同一份 `cases.json` 且约定定位失败一律 `t.Fatalf`，所以下游升 tag 之前必须先把语料打进去（C1）。

### 变更（PATCH 级：仓内自查工具、示例与文档，不改导出面）

- **QT-09（2026-10-23，本仓侧）**：新增三条随仓护栏与一份冻结基线，导出面零变化（仍 429/429）：
  - ⑤ `tools/check_contract_shape`：`domain/` `engine/` `fees/` 出现导出的 `*Config`/`*Options`/`*Impl` 形状名、或任何非测试源文件 import 数据库/驱动包即红；扫不到 Go 包或三条契约根目录缺任一同样红。ADR-0002 的在册例外做成**显式白名单**（逐符号写理由）。
  - surface 冻结：`tools/check_surface_freeze` + `docs/api-surface.txt`（412 行形状指纹：函数签名 / 结构体字段 / 接口方法 / 底层类型 / const·var）。改签名、改字段、增删符号都红；`-update` 重生成。与 `check_api_surface` 组成「登记 + 基线」双门。
  - ④ `.github/PULL_REQUEST_TEMPLATE.md`（四个必填段落）+ `tools/check_pr_body`（只在 `pull_request` 事件跑，其他事件由 CI `if:` 显式跳过）。
  - ② `examples/marketprovider` / `examples/enginerunner`：两个跨仓 port 的最小可编译用法（`package main`，不计入导出面），清单见 `examples/README.md`。
  - 以上全部接进 `.github/workflows/ci.yml`；README 快速开始、AGENTS、docs/README 的变更触发表同步。
- **QT-09 ③（2026-10-23，本仓侧）**：新增 `conformance/` —— 跨仓契约一致性语料 `cases.json`（33 条 / 9 kind）+ 契约侧自检腿。同一组用例分别喂给契约仓、core（`tellurion:pkg/conformance`）与 agent（`tellurion-agent:internal/conformance`），任一侧默认值或语义漂移即红（历史漂移点：agent 曾把滑点写成 0.001、佣金写成 0.00025）。语料是唯一真相，三侧各自读同一份 JSON（不共享加载器，解析漂移也在比对范围内）。本包**零导出符号**，导出面仍 429/429、冻结基线仍 412 行。
- **工具单测（2026-10-23）**：`tools/check_contract_shape` / `check_surface_freeze` / `check_pr_body` 各补一组表驱动单测（33 个用例），逐条红腿独立断言 + 破坏自证；三个工具的 `options.out` 由 `*os.File` 改为 `io.Writer`（零行为变化）以便注入 `bytes.Buffer`。

## [v0.1.0] - 2026-10-23

S1 落码首版（QT-05 ③-①，本仓侧）：`fees/` + `domain/market/` + `domain/` + `engine/` 四个库包，
**429 个导出符号**全部登记进 `docs/API.md` §4。0.x 阶段按 README 版本规则允许 breaking，逐条见下。

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
