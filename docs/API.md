---
status: active
last-verified: 2026-10-21
verified-by: CT-01（建本文件 + tools/check_api_surface 护栏）
---

# API — 导出面规格（本仓有哪些符号、各自稳不稳定、错了怎么判）

> **与本仓另两份文档的分工**：[`api-schema.md`](api-schema.md) 定**规则**（错误分类、幂等、版本兼容、不提供什么）；
> 本文件定**清单**（导出符号逐个回答「它在不在、稳不稳定、失败时是什么」）。
> 符号级冲突以本文件为准（它由 CI 核对，`api-schema.md` 不核对符号名）；规则级冲突以 `api-schema.md` 为准。
>
> **现状口径（开工前必读）**：本仓**今天只有一个占位包**（根 `doc.go`，零导出符号）。下面 §5~§7 全部是
> **目标面快照**——2026-10-21 用 §1 的命令从 core 实测抓出来的，S1（QT-05）落码时逐包迁入 §4 的登记表并接受 CI 核对。
> 把它们当"现在已经能 import 的东西"就错了。

---

## 1. 本文件怎么被核对（CI 护栏）

```bash
go run ./tools/check_api_surface            # CI 同款；不一致即退出码 1
go run ./tools/check_api_surface -list      # 只打印实际导出面，用来生成登记行
```

| 项 | 口径 |
|---|---|
| 核对粒度 | **包级导出符号名**：`func` / `type` / `const` / `var` 的名字，加上**导出类型的导出方法**（写作 `Type.Method`；接口的内联方法同样计入——端口的形状就是方法集） |
| 不参与核对 | `package main`（所以 `tools/` 天然不参与）、`*_test.go`、未导出符号、结构体字段与 tag |
| 红向 ① | **未登记**：代码里有导出符号、§4 登记表没有 ⇒ 红（新增导出符号必须先登记，这就是工单要的「未登记即红」） |
| 红向 ② | **已失效**：§4 登记了、代码里没有 ⇒ 红（改名/搬走没同步登记表，比缺登记更危险，因为它让下游按不存在的符号编码） |
| 红向 ③ | **结构红**：`docs/API.md` 找不到 `<!-- surface:begin/end -->` 标记对 ⇒ 红。登记表整段删不掉，护栏不能靠「删掉检查对象」变绿 |
| 红向 ④ | 解析失败（语法错误）与一个目录里多个包 ⇒ 红。否则解析失败的包会「凭空消失」，护栏静默空过 |

> ⚠️ **覆盖面限制（如实登记，别把它当全知护栏）**：内嵌接口（`type P interface { Base; Extra() }`）的**继承来的**方法名不在扫描粒度内——
> 那些方法本来就属于已登记的 `Base`，不是新符号。真要新增方法，新增的那个会被扫到。
> 结构体字段同理：字段增删**不触发**本护栏，它们由 core 侧的 SPEC 合规测试与 `resultspec_test.go` 那类字段级钉子管（见 §5）。

**登记表的形状**（§4 区内，一个包一个小节，小节标题里的反引号目录名就是包路径）：

```markdown
### `engine/` — 包 `contracts`

| 符号 | 类别 | 稳定性 | 错误语义 / 判据 |
|---|---|---|---|
| `EngineRunner` | interface | frozen | 见 §5.2 |
| `EngineRunner.RunBacktest` | method | frozen | `errors.Is(err, ErrStrategyFault)`；status 是唯一成功判据 |
```

只有**第一列**（单反引号包裹的符号名）参与核对，其余列给人读。列写错了不会红——所以 §5 的冻结面每次变更评审必须人工逐行读这张表。

---

## 2. 稳定性词表

| 标注 | 含义 | 变更代价 |
|---|---|---|
| **frozen** | 冻结面：字段集 / 取值集 / 语义本身就是契约。改前必须先**破坏 core 的护栏测试**（`resultspec_test.go`、`registry_test.go` 等），再走 core `docs/TASKS.md`「维护本文件」节的**冻结 schema 变更评审（硬规则）** | 最贵：契约变更 = 一次显式发布 + 下游同步 |
| **stable** | 已发布、承诺兼容。加字段/加方法可以，破坏性改名要按 CT-02 的版本规则（0.x：bump MINOR + CHANGELOG） | 中：走版本，不走「顺手改一下」 |
| **experimental** | 形状未定，允许随时破坏。下游用之前应在自己仓里包一层适配 | 低：但**不许**被当成已冻结面引用 |
| 未标注 | 不存在——§4 的表要求逐行填写；空标注在评审里按「未登记」处理 | — |

> **为什么要有这一列**：本仓最大的诱惑是「契约层放实现、字段随手加」。稳定性把「这条承诺多重」写进表里，
> 让下游能分辨哪些形状改了要发版、哪些还不能依赖。

---

## 3. 错误语义（全仓统一，§4/§5 的表里只指回来，不重复正文）

1. **判定只用 `errors.Is` / `errors.As` + 本仓哨兵**。禁止字符串匹配（`err.Error()` 的措辞不属于契约，随时可改；拿它判断等于把日志文案当 API）。
2. **包装必须保链**：跨层用 `%w`，禁止 `%v` 断链——断链后 `errors.Is` 失效，调用方只能退化成字符串匹配（就是在违反第 1 条）。
3. **分类即契约**：provider 侧必须把底层错误归入 §5.4 的三类之一，禁止裸 `fmt.Errorf` 透出（`api-schema.md` §2）。
4. **无法归类 = fail-closed**：按 `ErrContractBreach` 处理，宁可停，不要跑出漂亮但错的曲线。
5. **`status` 三值不是错误**：`completed` / `invalid` / `failed` 是**结果属性**，走响应字段不走 error。调用方唯一判据是 `status`——**禁止**用 `id` 是否存在、HTTP 200 之类间接信号推断成功（这条在 core 的 SPEC 冻结节里，本节只指路）。
6. **`invalid_reasons` 的取值是稳定 Code**（见 §5.3）：`Code` 给机器判读、`Detail` 给人读且可改。按 `Detail` 判断 = 违反第 1 条。

---

## 4. 机器核对的导出面登记表（CI 读这一区）

<!-- surface:begin -->

> **本区今天没有条目，这是事实而不是缺漏。** 本仓唯一的 Go 文件是根包 `doc.go`（占位，零导出符号），
> 所以 `go run ./tools/check_api_surface` 现在绿——它核对的是「零个包对零行登记」。
>
> **S1（QT-05）落码时必须做的事**：对每个落地的包跑
> `go run ./tools/check_api_surface -root <core 里的旧包目录> -list`，把输出逐行登记成 §1 的表格形状，
> 逐行填稳定性与错误语义，然后本区的核对才有实质内容。核对强度不靠人保证，靠红向 ①②③④。

<!-- surface:end -->

---

## 5. 冻结面逐符号规格（S1 迁入 §4 时的权威口径）

> 每个冻结面只**指回** `tellurion:docs/SPEC.md` 的节标题（SPEC.md 用命名节，不用数字编号），不复制正文——
> 复制的那一份一定会漂移。改这些节本身走 core `docs/TASKS.md`「维护本文件」的**冻结 schema 变更评审（硬规则）**。

### 5.1 四契约（ADR-031 D1）

| 契约术语 | 承载符号（今天住 core，S1 → `engine/`） | 稳定性 | 冻结文本指回 |
|---|---|---|---|
| `RunSpec` | `BacktestRequest` + `BacktestRequest.ValidateRunSpec` / `EffectiveKernel` / `EffectiveFillTiming`；伴生：`SignalWeight`、`TradingConfig` + `TradingConfig.WithDefaults`、`PriceLimitConfig`、`FillTiming` | **frozen** | `tellurion:docs/SPEC.md` §「RunSpec 内核路由（K10，2026-10-12）：EngineRunner 可插拔」 |
| `SignalSpec` | `SignalWeight`（K15 目标权重流）；信号本体 `domain.Signal` 随 `domain/` 迁入；供给契约 `SignalProvider`（今天 `pkg/strategy`，见 §5.2） | stable（**SPEC.md 无 SignalSpec 节**） | `tellurion-contracts:docs/api-schema.md` §1.2 + `tellurion:docs/design/qt-decoupling/kernel-interface-map.md` §一.4 |
| `ResultSpec` | `BacktestResponse` + `StatusCompleted` / `StatusInvalid` / `StatusFailed`；`InvalidReason` + `CheckValidity` + `InvalidReasonCodes` | **frozen** | `tellurion:docs/SPEC.md` §「ResultSpec 冻结（K9，2026-10-12）：BacktestResponse canonical schema」→ 子节「status 三值语义」「字段集（27 项，冻结）」「指标口径」 |
| `EngineRunner` | `EngineRunner` interface（+ `EngineRegistry` / `NewEngineRegistry`、`MarketRouter` / `ResolveRunTarget`） | **frozen**（签名与 fail-closed 语义） | 同 RunSpec 行：§「RunSpec 内核路由（K10）」的「接口边界」与「fail-closed（零静默回退）」两条 |

> ⚠️ **`RunSpec` / `ResultSpec` 这两个 Go 类型名今天不存在**（承载物是 `BacktestRequest` / `BacktestResponse`）。
> 是否改名是 W0（QT-01~04）的决定，见 `tellurion-contracts:docs/api-schema.md` §1.2 的提示行。
> 看到本表左边是术语、右边是符号名，不是笔误。

### 5.2 port 接口（方法集即形状）

| 接口 | 方法（导出面计入 `Iface.Method`） | 实现方 | 稳定性 | 错误语义 / 判据 |
|---|---|---|---|---|
| `EngineRunner` | `RunBacktest(ctx, BacktestRequest) (*BacktestResponse, error)` | core | **frozen** | 返回值里的 `Status` 是唯一成功判据；error 只装真实故障（§3.5） |
| `EngineRegistry` | `Register` / `Resolve` / `Names` | core（装配期） | stable | 未注册内核 → `ErrKernelNotRegistered`，**绝不静默回退 go**；`Resolve` 可并发，`Register` 仅装配期 |
| `MarketRouter` | `ResolveKernel` / `Markets` | core | stable | 未登记市场 → `ErrMarketNotRegistered`；路由表未装配 → 创建即拒绝、job 不落库 |
| `BacktestRunner` | `RunBacktest(ctx, strategyName, stockPool, startDate, endDate) (*domain.BacktestResult, error)` | core / agent 适配器 | stable | 与 `EngineRunner` **刻意分离**（五参 + 返回 domain 结果，AI 管线端口）；签名不动是 K10 明文约束 |
| `SignalProvider` | `Signals(ctx, Scope, asOf) ([]Signal, error)` | core（`pkg/strategy`） | stable（S1 落点待 QT-05） | 返回 `as_of > asOf` 的信号即违约：`RejectFutureSignals` fail-loud，不静默跳过 |
| `Provider`（`domain/market`，别名 `MarketDataProvider`） | `GetOHLCV` / `GetFundamental` / `GetStocks` / `GetLatestPrice` / `GetIndexConstituents` | core | stable | 缺数据 → 归 `ErrDataGap`（跳过当天 + 写 `invalid_reasons`），不是 fault |
| `RiskManager`（`domain/`） | `CheckStopLoss` / `CalculatePosition` / `DetectRegime` | core | stable | 裁决失败必须可归因；禁止把拒绝伪装成空信号 |
| `CorporateAction`（`domain/`） | `Type` / `Symbol` / `ExDate` / `RecordDate` / `PayDate` / `Description` / `Apply` | core | stable | 公司行为按 ex-date 生效；`Apply` 的现金变动走返回值，不静默改持仓 |
| `BrokerGateway` | **未定形**（形态见 `tellurion:docs/design/qt-decoupling/broker-gateway-design.md`） | `tellurion-broker` | ⬜ 设计中，S1 不落 | 见 `api-schema.md` §1.1「状态」列——今天源码 0 命中，**别当已存在接口去 import** |
| `DataPort` / `ExecutionPort` | **未定形** | `tellurion-data` / core | ⬜ 设计中，S1 不落 | 同上 |

### 5.3 状态与无效 Code（frozen 的取值集）

| 符号 | 取值 | 不变量（判据） | 指回 |
|---|---|---|---|
| `StatusCompleted` | `completed` | `invalid_reasons` 空、`error` 空 | SPEC §「status 三值语义」 |
| `StatusInvalid` | `invalid` | `invalid_reasons` **必填**、`error` 必空；指标照常给出，但**不得**当成功观测 | 同上 |
| `StatusFailed` | `failed` | `error` **必填**、`invalid_reasons` 空 | 同上 |
| `InvalidEmptyUniverse` 等 `Invalid*` 常量 | `empty_universe` / `zero_trades` / `zero_start_date` / `garbage_metric` / `exec_algo_unknown_style` / `exec_algo_attach_failed` | 稳定 Code；`InvalidReason.Detail` 才可改。新增 Code = 契约变更 | SPEC §「status 三值语义」的 Code 列表 |
| （QT-10 待落）| `data_gap_skipped` / `session_terminated` | ⬜ **今天源码里没有**。QT-10 ③-④ 把它俩映射进上面三值（不引入第四值），agent 侧消费口径见 `tellurion-agent:docs/TASKS.md` REV-01-A | `tellurion:docs/TASKS.md` QT-10 ③-④ |

> **禁止第四值**已两次明文（QT-10 ③-④、QT-63 ④）。`Failures[]` 只作诊断明细，**不是状态来源**——status 的唯一透传点是 `Engine.ResponseFromState`。

### 5.4 哨兵错误族

| 族 | 符号 | 稳定性 | 消费方行为 |
|---|---|---|---|
| 三分法（**QT-10 待落，源码里今天没有**） | `ErrDataGap` | ⬜ 目标 | 跳过当天 + 写 `invalid_reasons`（`data_gap_skipped`），run 继续 |
| 三分法 | `ErrStrategyFault` | ⬜ 目标 | **终止 run** → `status=failed` + `error`；不得当可重试的成功 |
| 三分法 | `ErrContractBreach` | ⬜ 目标 | **终止 run** → `status=failed`；同时是「无法归类」的兜底类（§3.4） |
| 内核路由 | `ErrKernelNotRegistered`、`ErrMarketNotRegistered` | stable（今天在 `pkg/backtest/contracts`，S1 → `engine/`） | 请求期 400；异步 job 终态 `failed`。**零静默回退** |
| 成交口径 | `ErrFillTimingMissing`、`ErrFillTimingUnknown`、`ErrFillTimingMismatch` | stable | 引擎永远不猜口径；口径不符是错误不是结果差异 |
| abi / artifact 校验 | `ErrAbiVersionUnknown`、`ErrAbiVersionMismatch`、`ErrAbiFamilyMissing`、`ErrAbiFamilyMixed`、`ErrImportNotWhitelisted`、`ErrImportSignatureMismatch`、`ErrHostRegionMissing`、`ErrHostRegionInvalid`、`ErrParamsTooLarge`、`ErrScopeUnknown`、`ErrScopeMarketV1` | stable（S1 → `engine/`；`errors/` 是否独立成包见 §8 待裁决） | 归 `ErrContractBreach` 语义：装载期拦下，不进 run |
| 信号校验 | `ErrActionIllegal`、`ErrStrengthOutOfRange`、`ErrTargetWeightOutOfRange`、`ErrEmitSourcesMixed`、`ErrEmitSourceUnknown`、`ErrSignalBudgetExceeded`、`ErrStateBudgetExceeded` | stable | 归 `ErrStrategyFault`（非法信号 / 预算超限 = 策略故障，终止 run） |
| provenance / hash | `ErrCanonicalUnsupported`、`ErrStrategyRefMalformed` | stable | 输入身份不可 canonical 编码 → 当场失败，禁止退化成「猜一个 hash」 |

---

## 6. 值对象面（`domain/`，S1 整包迁入）

来源与再生成：`go run ./tools/check_api_surface -root ../tellurion/pkg/domain -list`（2026-10-21 实测 227 个导出符号，含 `market/` 子包）。
**逐符号清单不在这里手抄**——抄一份就会漂；落码时用上面这条命令生成 §4 的表，再按下面分组填稳定性。

| 组 | 代表符号 | 稳定性 | 备注 |
|---|---|---|---|
| 回测输入输出 | `BacktestResult`、`BacktestParams`、`PortfolioValue`、`Trade`、`RiskMetrics`、`WalkForward*` | stable（`BacktestResult` 的**指标口径**受 SPEC 冻结节约束） | `BacktestResult` **没有也不该有 `Status`**——三值属响应层，见 `tellurion-agent:docs/TASKS.md` OBS-01-R2-A 的 core 裁决 (b) |
| 仪器域行情 | `OHLCV`、`Stock`、`Fundamental(Data)`、`IndexDaily`、`IndexConstituent`、`Dividend`、`Split`（多为 `market/` 的类型别名） | stable | 别名不是重复真相：实体在 `domain/market`，根包别名保兼容 |
| 时间/频率语义 | `Frequency`、`FrequencySpec`、`TimeSemantics`、`AsiaShanghai`、`Visible*`、`DailyVisibleAt` | **frozen**（PIT 红线） | 可见性语义是防前视的地基，改它=契约变更 |
| 信号与订单 | `Signal`、`Direction`、`Order`、`OrderType`、`Position(Size)`、`Portfolio`、`TargetPosition`、`StopLossEvent`、`MarketRegime` | stable | — |
| 因子族 | `FactorDefinition`、`FactorHypothesis`、`FactorType`、`FactorKernel`、`FactorCaliber`、`FactorRegistration`、`FactorMiningRecord`、`FactorOutput`、`ICEntry`、`FactorValue` + 内置定义/假设种子变量 | stable（`VersionHash` 语义受「版本哈希纪律」节约束） | 封闭词表扩展 = 加名 + 加实现 + 同步 SPEC 对应节 |
| 公司行为 | `CorporateAction`、`CashDividend`、`BonusShare`、`CorporateActionSplit`、`RightsIssue`、`Placement`、`ActionEngine`、`ApplyResult`/`ApplyOutcome` | stable | — |
| 配置形状 | `Config`、`ServiceConfig`、`DatabaseConfig`、`RedisConfig`、`TushareConfig`、`StrategyConfig`、`ExecutionConfig` | ⚠️ **experimental，且命名待评审** | 全部命中本仓 [`adr/0001`](adr/0001-契约只增不破且零实现.md) ③ 的 `Config` 命名禁令（`Config` 是精确同名，其余是同一种形状：契约位上写的是实现细节）。**S1 整包迁会把它们带进来**——落码前定「改名 / 留 core / 立 ADR 豁免」，见 §8 待裁决 #2，别默默收下 |
| core-only（agent 侧 0 命中） | QT-05 ③-⑨ 点名的 5 族：`Screen*`、`Sector`/`StockSector`/`SWIndustryNode`、`SupplyChainEdge`、`MacroSeries`、`mining` 族（`FactorMiningRecord`）；实测另含 `IndustryVitality`、`MarketActivity`、`Factor*`（定义/登记/输出）、`CorporateAction` 族、`RiskManager`、`Order`/`Signal`/`Stock` 等 | **experimental + 待回收** | QT-05 ③-⑨ 与 ADR-0001 的代价条款：S1 仍整包迁，S2 后按清单回收。**下游不要依赖这一组**。<br>实测口径（2026-10-21）：`pkg/domain` 根包 124 个导出符号里 **86 个在 agent 侧 0 命中**（`grep -rl "domain[.]<名>\b" --include='*.go'` 逐名数文件）。「0 命中」≠「该删」——agent 走 HTTP/JSON 消费的那几族（`Signal`/`Stock`）本来就不摸 Go 类型。完整清单是 **QT-05 ③-⑨ 的产出物**，本文件只登记判据与命令，不替它下结论 |

---

## 7. abi 常量族（`engine/` 内，随 QT-05 迁入）

来源：`tellurion:pkg/backtest/contracts/abi.go`（v2 修订见 `tellurion:docs/design/qt-decoupling/adr-029-amendment-host-api-v2.md`）。

| 组 | 符号 | 稳定性 | 判据 |
|---|---|---|---|
| host 函数名 | `HostFnGetBar` / `HostFnGetBar2` / `HostFnGetDate2` / `HostFnGetSeriesLen` / `HostFnGetSeriesLen2` / `HostFnGetCrossSection` / `HostFnGetSymbolCount` / `HostFnEmitSignal` / `HostFnEmitTarget2` / `HostFnGetState` / `HostFnSetState` / `HostFnLog` | **frozen** | 名字与签名成对冻结；`HostABIContract`/`HostABIContracts` 是唯一名单源。artifact 导入表不符 ⇒ `ErrImportSignatureMismatch` |
| abi 版本 | `HostAbiV1` / `HostAbiV2` / `HostAbiVersion` / `HostAbiVersionFromName` / `HostFnCount` | **frozen** | 声明与导入段矛盾 ⇒ `ErrAbiVersionMismatch`；两族混用 ⇒ `ErrAbiFamilyMixed` |
| 内存区所有权 | `ExportHostParamsPtr` / `ExportHostScratchPtr` / `ExportHostScratchSize` / `MinHostParamsPtr` / `MinHostScratchSize` / `HostRegionExports` / `HostRegionAlignment` / `HostRegionDeclaration` / `ValidateParamsWrite` | **frozen** | 所有权写死在常量里；缺失 ⇒ `ErrHostRegionMissing`，非法 ⇒ `ErrHostRegionInvalid` |
| 预算上限 | `MaxParamsBytes` / `MaxScratchPerBar` / `MaxSignalsPerBar` / `MaxStateKeys` / `MaxStateValueBytes` / `MaxInstanceStateBytes` | **frozen** | 超限是 `Err*BudgetExceeded` / `ErrParamsTooLarge`，不是静默截断 |
| 值类型 / 字段号 | `HostI32` / `HostI64` / `HostF64` / `HostValueType` / `FieldOpen`…`FieldVolume` / `FieldID` | **frozen** | 枚举上限同属 abi；未知 ⇒ `ErrAbiVersionUnknown` 类 |
| scope | `ScopeMarket` / `ScopeSymbol` / `StrategyScope` / `ResolveScope` | **frozen** | v1 给 market scope ⇒ `ErrScopeMarketV1`（**目标** D25「M0 市场级先行」要解开这一约束，但今天代码是拒的） |

---

## 8. 不登记的东西（边界与非目标）

- **`dto/`**：跨进程 DTO 是 `tellurion:docs/openapi.yaml` 的**生成物**（`oapi-codegen`），真相在 core，本仓只收生成物且**禁手改**（CI 重生成 diff 必须为空）。见 `tellurion:docs/design/qt-decoupling/repo-inventory.md` §3。所以本文件不列 DTO 字段。
- **`schemas/`**：数据契约资产，不是 Go 包，不产生导出符号。
- **内部函数 / 未导出符号 / 结构体字段**：不是导出面（§1 粒度）。
- **默认实现、mock、测试替身**：`api-schema.md` §5 已禁（契约仓被测试代码污染）。

### 待裁决（S1 落码前要有结论，别默默带过）

| # | 事项 | 现状 |
|---|---|---|
| 1 | `errors/` 与 `abi/` 是否独立成包 | [`ARCHITECTURE.md`](ARCHITECTURE.md) §2 列了这两个包，但 QT-05 ③-① 只说迁 `domain/`、`engine/`、`fees/`（`engine/` 内保持 `package contracts`），abi 常量与哨兵今天都随 `engine/` 走。二者口径不一致，本文件按 QT-05 落点写（§5/§7 归 `engine/`） |
| 2 | `*Config` 命名族要不要进本仓 | ADR-0001 ③ 禁 `Config`/`Options`/`Impl` 命名，但 S1 整包迁会带进来 9 个：`domain/` 的 `Config`（精确同名）、`ServiceConfig`、`DatabaseConfig`、`RedisConfig`、`TushareConfig`、`StrategyConfig`、`ExecutionConfig` + `engine/` 的 `TradingConfig`、`PriceLimitConfig`。后两个还是 §5.1 的 frozen 伴生符号——禁令与冻结面在这儿撞车，落码前要结论，不能两边都当没看见 |
| 3 | `fees/` 的取舍 | QT-05 ①「只放数据形状」+ ④「不搬 `fees` 的计算部分」：`AShareFees`/`Default*` 常量/`StampTaxCutDate` 进；`StampTaxRateFor`（按日取档的函数）与 `ApplyDefaults`/`Validate` 留 core —— 但 `AShareFees` 的方法在迁入包里，切不开。落码前必须定「方法算不算计算部分」 |

---

## 9. 与其他文档的互锁

| 文档 | 关系 |
|---|---|
| `api-schema.md` | 规则源（错误分类 / 幂等 / 兼容 / 不提供什么）；§1.1 的 port「状态」列与本文件 §5.2 必须一致 |
| `ARCHITECTURE.md` §2 | 包布局与本文件 §8 待裁决 #1 有关 |
| `tellurion:docs/SPEC.md` | **冻结文本的唯一源**；本文件只指节标题。改冻结面 = 先破坏 core 合规测试 → 更新 SPEC → 变更评审 → 才动代码（`tellurion:docs/TASKS.md`「维护本文件」硬规则） |
| `tellurion:docs/design/qt-decoupling/kernel-interface-map.md` | 现状接口面清单的来源（M6 的「散在几处」即指它）；本文件不是它的替代，是**随仓可核对**的那一份 |
