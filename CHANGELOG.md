# CHANGELOG — tellurion-contracts

> 本仓的版本账本。**每个发版一节**，标题格式固定 `## [vX.Y.Z] - YYYY-MM-DD`（`tools/check_release_matrix` 按它和 README 兼容矩阵、`git tag` 三方对拍）。
> 未发版的改动记在 `## [Unreleased]` 下，发版时把那一节的内容搬进新版本节。
> 规则口径见 [`README.md`](README.md) 的「版本规则」（源头：`tellurion:docs/TASKS.md` QT-19 ③）。
> breaking 的条目必须写成「破坏性：<旧 → 新>」，光写「改了签名」不算留账。

## [Unreleased]

### 破坏性（CT-06-B，2026-11-10）

- **破坏性：`domain/` 根包里 63 个「只有 core 一个编译使用方」的登记符号 → 删除本仓导出，实类型迁入 `tellurion:pkg/domain`；另 3 个 0 使用方形状（`Sector` / `StockSector` / `StockIndustry`）与 `StrategyConfig` 直接删除、不在 core 重建。**
  依据：补裁 #12（contracts C4：本仓只留 ≥2 使用方的形状）+ 11 仓 grep / Go AST 复核；逐名处置见 `docs/API.md` §6.1 的「CT-06-B 已执行」块。
  - **迁回 core**（实类型物理搬家；market 别名在 core 侧同形重建）：`Signal` / `Order` / `Position` / `Portfolio` / `RiskMetrics` /
    `TargetPosition` / `PositionSize` / `StopLossEvent` / `RiskManager` / `OrderType` 族；`BacktestParams`；
    `FactorKernel` / `FactorCaliber` / `FactorDefinition` / `FactorRegistration` / `FactorMiningRecord` / `FactorOutput` /
    `FactorCacheEntry` / `FactorReturn` / `ICEntry` / `BuiltinFactorDefinitions`；
    `ScreenFilters` / `ScreenRequest` / `ScreenResult` / `SWIndustryNode` / `RelationIndustryMember` / `SupplyChainEdge` /
    `IndustryVitality` / `ETFConstituent` / `FundamentalsDetailRow` / `ExecutionConfig` / `DefaultExecutionConfig`。
  - **删除且不重建**：`StrategyConfig`；`Sector` / `StockSector` / `StockIndustry`（`domain/sector.go` 整文件）；
    `IndexConstituent` / `FundamentalData` 两个**根别名**（实体仍在 `domain/market/`）。
  - **保留**：6 个 market 别名（`Fundamental` 仍被 `tellurion-data` 直接使用）；`FactorMarketActivity`（`FactorType` 名册成员）；
    `MacroSeries` / `MarketActivity`（core + `tellurion-data`）；`BuiltinFactorHypotheses`（core + agent）。
  - **实测**：`go run ./tools/check_api_surface` → 391 符号 / 登记 391（改前 457）；冻结基线
    `git diff --numstat docs/api-surface.txt` → 0 / −63（447 → 384 行）；`docs/API.md` §4 删 66 条失效登记行。
    本仓 `domain/` 删 14 个文件（13 个迁出 + `sector.go`），`types.go` / `backtest.go` / `factor.go` 收缩。
  - **下游影响**：core 侧调用点已同批切至 `tellurion/pkg/domain`（core 仓同笔 commit）。core 当前仍 pin `v0.9.0`，
    在 bump 到含本条目的版本前，core `internal/repoguard/contract_identity_test.go` 的腿 C 以「版本滞后」白名单
    临时登记这批同名类型；**bump 后必须删除，反向棘轮会判红**。0.x 版本规则允许 breaking（删符号 = MINOR），
    本条随下一个 MINOR 发版。

### 破坏性（CT-06-A，2026-11-09）

- **破坏性：`domain.BuiltinHypothesisFor`（func）— 删除 → 由唯一使用方 `tellurion-agent` 本地私有实现 `builtinHypothesisFor` 承接。**
  依据：11 仓重跑 + Go AST 复核，该符号只有 `tellurion-agent` 一个编译使用方，不满足准入规则 C4（≥2 使用方）。
  内置表 `domain.BuiltinFactorHypotheses` 仍有 core + agent 两个使用方，**保留不迁**。
  0.x 版本规则允许 breaking（删符号 = MINOR），本条随下一个 MINOR 发版；`docs/API.md` §4/§6.1 与
  `docs/api-surface.txt` 已同步。

### 新增（risk 契约，2026-11-11）

- **`risk/` 同步 scope 契约（纯加法，既有符号与编码一字未动）**：新增
  `Scope.ProcessLevel()` 与 `Snapshot.ValidateScope(scope)`，把 2026-11-11 主任裁决的正式口径
  编码成纯函数 / 纯结构校验：
  - `global → strategy` 是**正常优先级回退**（global 组合级基线，strategy 单策略覆盖），
    `account` / `symbol` 不参与进程级策略源；
  - 选中 scope 内**四类规则各恰好一条**，缺一类或同类多条 ⇒ 拒（`ErrInvalidParams`），
    未知 scope / type ⇒ `ErrInvalidScope` / `ErrInvalidRuleType`。
- **边界不变**：`ForScope` 的 fallback 选择、`params` 的冻结词典严格解码仍属服务策略（CT-03 ④）；
  本包只做结构判定，不做 decode / 默认值 / 进程内索引。
- **下游影响**：纯加法、无编译破坏，**不需要即时升 pin**；`tellurion-risk` / core 可在下一个
  MINOR（建议 `v0.10.0`）发版后选用。`docs/API.md` §4 与 `docs/api-surface.txt` 已同步。

## [v0.9.0] - 2026-11-07

MINOR 发版：CT-05 新增 `broker/` 包（**只增不改**），导出面 427 → 458
（+31：12 个 `const` + 12 个 struct + 3 个枚举类型 + 4 个谓词方法），冻结基线 417 → 448 行
（`git diff --numstat docs/api-surface.txt` 实测 +31 / −0）。既有六个包一字未动。

### 新增（MINOR 级：新包）

- **CT-05（2026-11-07）：`broker/` —— 把 broker 的订单 / 成交 / 账户 wire 形状收成唯一实现。**
  - 13 个类型（`Side` / `OrderType` / `OrderStatus` + 12 个 struct）与 12 个常量、4 个谓词方法
    （`Side.Valid` / `OrderType.Valid` / `OrderStatus.Valid` / `OrderStatus.Terminal`）。
  - **字节来源**：期望值由 `tellurion-broker:internal/gateway` 的生产类型 `json.Marshal` 打出后
    逐字节抄入 `broker/wire_test.go` 与 `conformance/broker_cases.json`（10 例）；不是按契约侧反推。
  - **不迁**：XTP 适配器内部类型、`OrderRequest.Fingerprint`（仓储幂等规则）、12 个契约哨兵
    （进程内错误值）；`quantity` 整股 `int` / `Side` 两值 / `OrderType` 两值与 `domain/` 的差异
    写进 `docs/api-schema.md` §3.2，**禁止**类型别名抹平。
  - **未随本单落（如实）**：`tellurion-broker` 删本地形状改引本包**未执行**——tag 未 push，
    `go get` 解析不到；解除前置后由 broker 一笔做，判据是报文快照逐字节一致。

## [v0.8.0] - 2026-11-07

MINOR 发版：CT-04 新增 `run/` 包（**只增不改**）＋ CT-06 回收 `domain/` 根包 0 使用方导出面
（**破坏性**，见本节末）。导出面 418 → 427（+9：2 个 `const` + 6 个 struct + 1 个纯函数），
冻结基线 408 → 417 行（`git diff --numstat docs/api-surface.txt` 实测 +9 / −0）。
既有六个包除 CT-06 的删除外一字未动。

### 新增（MINOR 级：新包）

- **CT-04 契约腿（2026-11-07）**：`run/` —— 把编排仓已在三处**手写复刻**的 wire 形状收成唯一实现：
  - `ClaimView` 七字段（`claimView` 镜像；`ParentJobID` 带 `omitempty`，普通 run 不出现该键）。
  - `ProvenanceView` 十五字段（`provenanceView` 镜像；六个可空字段**出 `null` 而不是省略键**）。
    core `docs/openapi.yaml` 的 `BacktestProvenance` 原文 `must not drift`，本包把它变成机械判据。
  - `ParentSubmitRequest` / `ParentBatchRequest` / `ParentWalkForwardRequest` / `ParentWindows`
    （`parentSubmitRequest` 族镜像；`kind` 决定 batch / walkforward 二选一，契约不猜）。
  - `ChildKeyPrefix = "parent"` / `ChildKeySeparator = ":"` / `ChildIdempotencyKey(parentID, childIndex)`
    —— 纯函数 `parent:<parent_job_id>:<child_index>`，逐字节对齐编排仓。
    **⚠️ 改此派生规则 = MAJOR**（全部子 run 的幂等身份改变）。
  - **对拍**：三份形状的 wire 字节由生产侧形状独立重建（字面量），`run/run_test.go` 逐字节断言；
    `conformance/run_cases.json` 是单独一份语料（与 `cases.json` 同形，但 core / agent 的加载器
    只认 `contracts` / `core` / `agent` 三侧，混入新侧会卡住下游 pin 升级），期望值由本仓实现产出后写死。
  - **不迁**（按 CT-04 ④）：`RunRecord` / `ParentRecord` / `WorkerRecord`（队列与租约状态机）、
    `jobView` / `parentView` / `parentChildView` / `parentCountsView`（编排读面投影）、
    `ArtifactStore`（带写语义）、`run_artifacts` DDL。包内零 I/O、零时钟、零 decode。
  - **未随本单落（如实）**：orchestration 改别名、core `orchestrationworker` 改用契约类型、
    analytics 提交报文改用契约类型，都是各自仓的单（CT-04 ⑥）；本单落地时三处仍各自持有副本。
    下游 pin 升级未执行（`go get` 要求 tag 已在远端）。

### 破坏性（CT-06，已随本版发出）

- **CT-06（2026-11-07）：回收 `domain/` 根包 0 使用方导出面，导出符号 487 → 418（−69）。**
  - 取证：`domain/` 根包 176 符号折 119 个根名，逐一在 11 个同级仓（`tellurion*`，排除本仓，
    `*.go` 含测试，仅统计导入 `contracts/domain` 的文件）grep `\bdomain\.<Name>\b`。
    逐名出处见 `docs/API.md` §6.1。
  - **破坏性：删除公司行为族 68 符号** —— `CorporateAction`、`ActionType`、
    `CashDividend` / `BonusShare` / `CorporateActionSplit` / `RightsIssue` / `Placement`、
    `ActionEngine`、`ApplyResult` / `ApplyOutcome` 及其构造函数 / 方法（`domain/corporate_action.go`
    整文件与测试同删）。11 仓实测 **0 使用方**。
  - **破坏性：删除 `domain.MarketDataProvider` 别名 1 符号**（实体 `domain/market.Provider` 保留；
    11 仓实测 0 使用方）。
  - 版本规则口径：README 明确 0.x 阶段**没有**「deprecated 至少保留 1 个 MINOR 周期」的义务，
    故本单直接删净，不先标 Deprecated；发版时按 0.x breaking → MINOR 记入下一 MINOR。
  - **未随本单落（如实）**：61 个根名只有 1 个使用方（`tellurion` 60 + `tellurion-agent` 1），
    按 CT-06 ② 应迁回该使用方仓；迁回须改其它仓并删本仓公共符号，超出本单「只改本仓」边界，
    已在 §6.1 登记「迁回（待跨仓单）」，本轮不执行。`CaliberRank` / `CaliberRaw` 虽 0 外部引用，
    但属 `FactorCaliber` 封闭词表且被 `Valid()` 依赖，保留不改语义。

## [v0.7.0] - 2026-11-07

MINOR 发版：CT-03 新增 `risk/` 包（**只增不改**），导出面 465 → 487
（+22：1 个 `const` 域标签 + 1 个 `const` 体积上限 + 2 个枚举 `type` + 8 个枚举 `const`
+ 2 个 struct + 4 个 `var` 哨兵 + 4 个方法），冻结基线 448 → 470 行（`git diff --numstat` 实测 +22 / −0）。
无破坏性变更；既有五个包一字未动。

### 新增（MINOR 级：新包）

- **CT-03 契约腿（2026-11-07）**：`risk/` —— 把风控规则快照的**形状与 hash 算法**收成跨仓唯一实现，
  任何消费方 `import` 一次就能算出与 risk 逐字节相同的 `risk_rules_hash`。
  - `Snapshot` 三字段（`RiskRulesHash` / `Rules` / `CreatedAt`）+ `Clone()` / `Hash()` / `Verify()`；
    `SnapshotRule` 六字段（`RuleID` / `Revision` / `Scope` / `Type` / `Params` / `VersionHash`）。
  - 冻结词表：`RuleType`（`position_sizing` / `stop_loss` / `regime` / `take_profit`）、
    `Scope`（`global` / `account` / `strategy` / `symbol`）+ `Scope.Valid()`。
  - `DomainSnapshotHash = "tellurion.risk_snapshot_hash.v1"`（逐字节不动；改 = 全历史 hash 作废）；
    `MaxParamsBytes = 65536`（⚠️ 与 `engine/` 的 4096 **同名不同值不同义**，禁止「统一」）。
  - 4 个跨仓哨兵：`ErrInvalidRuleType` / `ErrInvalidScope` / `ErrInvalidParams` / `ErrSnapshotMismatch`，
    **文案与 risk 逐字一致**；服务私有 7 哨兵不迁。
  - **编码逐字节搬运**：域标签裸写一行 + `\n`、长度前缀帧 `len(key):key=len(value):value\n`、
    按 `rule_id` 升序且重复即拒、参数 `UseNumber` 解码后 `json.Marshal`。**禁止**改用
    `engine.NewCanonical`（帧结构 / 域标签 / 参数规范化三处都不同；换 = `quant.run_provenance`
    与 `run.run_artifacts` 两处锚点同时失效）。
  - **对拍**：固定向量 `14ecab865bebaef54052354b63caa5d351f533ab112e1a29a9dcad3707955c5b`
    （由 `tellurion-risk` 内真实调用 `rule.SnapshotHash` 取得，2026-11-07 退出码 0）+
    `TestSnapshotHash_MatchesRiskRepoEncoding` 逐行重建 risk 侧编码并断言同 hex。
  - **不迁**（按 CT-03 ④）：`Bundle` / `ParseBundle` / `ForScope` 带 decode、默认值填充与进程内索引；
    它们属服务策略，进契约就把实现冻进底座。
  - 随单护栏：`tools/check_contract_shape` 的 `contractRoots` 四根 → 五根（+`risk`），
    夹具新增「risk 命名红」用例、缺根红期望补 `risk/`。
  - **未随本单落（如实）**：risk 侧改别名与 `tellurion-risk:docs/ARCHITECTURE.md:85` 判据订正归 risk 仓；
    core 的 `bundle → RiskPolicy` 接线归 `tellurion:SVC-CORE-RISK-RETIRE`；core / risk / agent 的 pin 升级
    **未执行**（tag 未 push，`go get` 拉不到）。⇒ 本包今天**零 import**，有意取舍见 `docs/API.md` §5.6。

## [v0.6.0] - 2026-10-24

MINOR 发版：QT-40 X1a 新增 `dataprotocol/` 包（**只增不改**），导出面 447 → 465
（+18：1 个 `type` + 2 个方法 + 15 个 `const`），冻结基线 430 → 448 行（`git diff --numstat` 实测 +18 / −0）。
无破坏性变更；既有四个包一字未动。

### 新增（MINOR 级：新包）

- **QT-40 X1a 契约腿（2026-10-24）**：`dataprotocol/` —— 把三处**已在漂**的同名形状收成一份共享定义
  （`tellurion-data:internal/ingest/point.go` 带 `VisibleAt`、`tellurion:pkg/data/source/unified.go`
  有 `DeduplicateKey()`、`tellurion:pkg/storage/bulk_insert.go` 因 import 成环另抄一份）。
  - `UnifiedDataPoint` **七字段**：`Symbol` / `TradeTime` / `Source` / `DataType` / `Data` / `IngestTime` / `VisibleAt`
    （字段顺序与 data 镜像一致，`VisibleAt` 末位）。`Data` 是 `map[string]interface{}`（JSONB 逃生阀），
    `DataType` 仍是裸 `string`（具名化会碰 wire，等 X1b 的 schema 冻结）。
  - `ToJSON()` —— `json.Marshal(Data)`；`DeduplicateKey()` —— `Symbol|DataType|TradeTime`（UTC，RFC3339Nano），
    **Source 不入键**（fallback 链可能多源同一数据，留哪条由调用方裁决）。
  - 15 个 `DataType*` 常量：与 core / data 的现取值**逐字一致**；字面量就是表路由，改一个值 = 数据静默落错表。
  - **不提供构造器、不读时钟**：`IngestTime` 必须由调用方注入（ADR-0004；core 版构造器读 `time.Now()` 是这次
    不迁构造器的直接原因）。**不重复 `TimeSemantics`**（它住 `domain/market`）；**不落 `Aggregator`**
    （今天无实现、无消费方，随 X1b 定形）。
  - 裁决：`docs/API.md` §8 #4 定 **(a) 先迁**；判据是「第二份独立实现已经存在」，不是「第二个 import」。
    形状与本包语义见 §4 / §5.5。向量由 `dataprotocol_test.go` 钉死（含 core 的既有向量
    `AAPL|global_ohlcv|2026-05-14T09:30:00Z`）。
  - **未随本单落（如实）**：core 删两份副本、data 镜像改 import、storage 列布局的对表测试，
    都归各自仓的单；在那之前三份形状仍并存。本仓只交付类型、向量与登记。

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
