---
status: active
last-verified: 2026-11-11
verified-by: CT-06-ALIAS-EXIT（2026-11-11：删除 `domain/` 根包 6 个 market 类型别名 —— `OHLCV` / `Stock` / `IndexDaily` / `Split` / `Dividend` / `Fundamental`，改用 `domain/market.*`；§4 登记区相应删 6 行、`api-surface.txt` 重生成；本仓 `go build/vet/test ./...` 与四道护栏全绿）；代码-文档对齐审计（2026-11-11：补 CT-06-B 执行结论；§6 因子族「待回收」改为已裁决；导出面 387 / 冻结基线 380 经 `check_api_surface` + `check_surface_freeze` 复核）；CT-06-A（2026-11-09：`domain.BuiltinHypothesisFor` 迁回唯一使用方 `tellurion-agent` 并删本仓导出；`domain/` 根包 107 → 106 符号 / 97 → 96 根名，迁回待办剩 60 个 core 根名；`api-surface.txt` 重生成）；CT-06 契约腿（2026-11-07：`domain/` 根名重跑取证 —— 119 根名逐一在 11 个同级仓 grep；删除 22 个 0 使用方根名 / 69 个符号（公司行为族 68 + `MarketDataProvider` 别名 1），61 个单使用方根名登记「迁回（待跨仓单）」，`CaliberRank`/`CaliberRaw` 作为封闭词表成员保留；导出面 487 → 418，冻结基线同步；§6.1 为逐名 grep 出处）；CT-03 契约腿（2026-11-07：新建 `risk/` 包 —— `Snapshot`/`SnapshotRule`/`Scope`/`RuleType` + 冻结词表 + `Snapshot.Hash()`/`Verify()`/`Clone()` + 4 个跨仓哨兵 + `DomainSnapshotHash`/`MaxParamsBytes`；编码逐字节搬运 risk 的 `SnapshotHash`，禁止换 `NewCanonical`；`Bundle`/`ParseBundle`/`ForScope` 与服务私有 7 哨兵不迁；导出面 465 → 487，冻结基线同步；固定向量与跨仓对拍腿在 `risk/snapshot_test.go`）；QT-40 X1a 契约腿（2026-10-24：新增 `dataprotocol/` 包 —— `UnifiedDataPoint` 七字段（含 `VisibleAt`）+ `ToJSON`/`DeduplicateKey` + 15 个 `DataType*` 常量；无构造器、不读时钟、不含 `Aggregator`；导出面 447 → 465，冻结基线 +18 行；§8 #4 按 (a) 先迁落定）；QT-22 契约腿（2026-10-24：`engine/provenance.go` 增 `StrategyVersionInput` 六字段 + `Hash()`/`MustHash()`，`engine/canonical.go` 增 `DomainStrategyVersionHash` = `tellurion.strategy_version_hash.v1`；集合 = `strategy_id + carrier + abi_version + scope + entry + source`，**不含** `capability`、**不含** `params_defaults`（后者若进必须升 v2）；导出面 443 → 447，冻结基线同步 +4 行）；QT-20 契约腿（2026-10-24：`BacktestRequest` 增可选 `StrategyRef`（wire 形态 = `ParseStrategyRef`/`String()` 的 `id`/`id@hash`）与 `StrategyParams`（浅层对象，嵌套键在 `ValidateRunSpec` 拒绝）；`domain.BacktestParams` 增生效参数 `Params`（对齐 core `RunMeta.Params`）；未知键清单不复制进本仓，归 core manifest（J2）；导出符号数不变 439，冻结基线 +2 行字段形状）；QT-41b ④ 错误族契约腿（2026-10-24：`ErrCompileFailed`/`ErrCompileUnavailable`/`ErrSourceTooLarge`/`ErrArtifactTooLarge` 四哨兵迁入 `engine/execution_errors.go`，名字语义不改，资产服务私有哨兵不迁；导出面 439 → 443）；QT-10 ③-① 契约腿（2026-10-23：三类执行错误哨兵 `ErrDataGap`/`ErrStrategyFault`/`ErrContractBreach` 落 `engine/execution_errors.go` + `CodeDataGapSkipped` 落 `engine/validity.go`）；QT-05 ③-① 本仓侧落码完成（`fees` 13 + `domain/market` 46 + `domain` 176 + `engine` 194 = 429 个符号当时已登记）；§8 待裁 #1（errors//abi/，ADR-0003）、#2（`*Config` 族，ADR-0002）、#3（fees 整包）均已裁决；#4（dataprotocol 时机）已登记为 X1a 待办（2026-10-22 推送回写：6 个 commit 已进 `origin/main`，该批 HEAD = `fbe9a2d`）
---

# API — 导出面规格（本仓有哪些符号、各自稳不稳定、错了怎么判）

> **与 [`api-schema.md`](api-schema.md) 的分工**：那份定**规则**（错误分类、幂等、版本兼容、不提供什么）；
> 本文件定**清单**（导出符号逐个回答「它在不在、稳不稳定、失败时是什么」）。
> 符号级冲突以本文件为准（它由 CI 核对，`api-schema.md` 不核对符号名）；规则级冲突以 `api-schema.md` 为准。
>
> **现状口径（开工前必读）**：本仓的**库代码已全部落地**（QT-05 ③-①，2026-10-22）：`fees/`、`domain/market/`、`domain/`、`engine/` 四个包，
> 加 QT-40 X1a 的 `dataprotocol/`（2026-10-24）共五个包，465 个导出符号全部登记在 §4（QT-40 后，2026-10-24）。`tools/` 下的命令是 `package main`，不计入导出面。
> **§5~§7 已从「目标面快照」变成「现行规格」**——它们的符号今天都在 §4 的登记表里，本节只是把冻结面 / 语义讲透。
> ✅ **下游今天已可 `import` 本仓**：6 个 commit 已于 2026-10-22 推到 `origin/main`（该批 HEAD `fbe9a2d`；仓外空模块里 `go get github.com/ruoxizhnya/tellurion-contracts@fbe9a2d` 实测解析为 `v0.0.0-20261022191510-fbe9a2d0df17`）。
> ⚠️ core / agent 切依赖那条腿（③-③④⑤⑦⑨）由下游仓执行；在它落地前，core 仍用自己那份 `pkg/domain` / `pkg/backtest/contracts`（双份真相窗口，见 `docs/ARCHITECTURE.md` §8）。

---

## 1. 本文件怎么被核对（CI 护栏）

```bash
go run ./tools/check_api_surface            # CI 同款；不一致即退出码 1
go run ./tools/check_api_surface -list      # 只打印实际导出面，用来生成登记行
go run ./tools/check_surface_freeze         # 形状指纹 ↔ docs/api-surface.txt；任何一行增删改即退出码 1
go run ./tools/check_surface_freeze -update # 重生成冻结基线（改签名/字段后必须与登记表同一 PR）
```

| 项 | 口径 |
|---|---|
| 核对粒度 | **包级导出符号名**：`func` / `type` / `const` / `var` 的名字，加上**导出类型的导出方法**（写作 `Type.Method`；接口的内联方法同样计入——端口的形状就是方法集） |
| 不参与核对 | `package main`（所以 `tools/` 天然不参与）、`*_test.go`、未导出符号、结构体字段与 tag |
| 红向 ① | **未登记**：代码里有导出符号、§4 登记表没有 ⇒ 红（新增导出符号必须先登记，这就是工单要的「未登记即红」） |
| 红向 ② | **已失效**：§4 登记了、代码里没有 ⇒ 红（改名/搬走没同步登记表，比缺登记更危险，因为它让下游按不存在的符号编码） |
| 红向 ③ | **结构红**：`docs/API.md` 找不到 `<!-- surface:begin/end -->` 标记对 ⇒ 红。登记表整段删不掉，护栏不能靠「删掉检查对象」变绿 |
| 红向 ④ | 解析失败（语法错误）与一个目录里多个包 ⇒ 红。否则解析失败的包会「凭空消失」，护栏静默空过 |
| 登记表形状 | 包小节标题必须带反引号包目录（形状见下面的示例）；**根包写 `.`**——不认根包就会让「根包新增导出符号」无路可登记，那是不可满足的护栏。标题不带目录时，那一节的表行**不算登记**：写错标题的人被「未登记」正面挡住，而不是悄悄登记进一个不存在的包 |

> ⚠️ **覆盖面限制（如实登记，别把它当全知护栏）**：内嵌接口（`type P interface { Base; Extra() }`）的**继承来的**方法名不在扫描粒度内——
> 那些方法本来就属于已登记的 `Base`，不是新符号。真要新增方法，新增的那个会被扫到。
> 结构体字段同理：字段增删**不触发**本护栏，它们由 core 侧的 SPEC 合规测试与 `resultspec_test.go` 那类字段级钉子管（见 §5）。

> **第二道门（surface 冻结）**：`tools/check_surface_freeze` 比对 `docs/api-surface.txt` 的形状指纹
> （函数签名 / 结构体字段 / 接口方法 / 底层类型 / const·var）。任何一行增删改都红——**结构体字段的增删在这里也会咬人**，
> 正好补上上一段说的覆盖面空白。新增符号必须**先登记本表、再 `-update` 基线**，两道门都过才算一次契约变更。

**登记表的形状**（§4 区内，一个包一个小节，小节标题里的反引号目录名就是包路径）：

```markdown
### `.` — 根包（今天只有占位的 doc.go）

| 符号 | 类别 | 稳定性 | 错误语义 / 判据 |
|---|---|---|---|
| `SomeRootSymbol` | var | stable | — |

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
| **frozen** | 冻结面：字段集 / 取值集 / 语义本身就是契约。改前必须先**破坏 core 的护栏测试**（`resultspec_test.go`、`registry_test.go` 等），再走 `tellurion:docs/TASKS.md`「维护本文件」节的**冻结 schema 变更评审（硬规则）** | 最贵：契约变更 = 一次显式发布 + 下游同步 |
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

> 第 3、4 条说的三类（`ErrDataGap` / `ErrStrategyFault` / `ErrContractBreach`）**已于 2026-10-23 落在
> `engine/execution_errors.go`**（QT-10 ③-①，§5.4 已改标 ✅）。
> ⚠️ 但**分类包装（QT-10 ③-②）在 core**：provider 侧把底层错误归入哪一类，是本仓之外的改动。
> 在它落地前，跨边界透出的仍是未分类 error——第 3、4 条对**生产者**依然是验收标准，只有对**消费方**的判据（`errors.Is` + 兜底类）今天可用。

---

## 4. 机器核对的导出面登记表（CI 读这一区）

<!-- surface:begin -->

> **登记口径**：一个包一个小节，标题里的反引号目录名就是包路径（根包写 `.`）。只有第一列参与核对，
> 其余三列给人读——所以**稳定性与错误语义靠评审逐行读，护栏核不到**（§1 末的覆盖面限制）。
> 新增导出符号：先跑 `go run ./tools/check_api_surface -list` 拿到符号名，在本区补行，再落代码；
> 顺序反了 CI 就红（红向 ①）。

### `.` — 根包（只有 `doc.go`，零导出符号）

| 符号 | 类别 | 稳定性 | 错误语义 / 判据 |
|---|---|---|---|

> 本区没有行是事实：根包只放文档（`package contracts`），代码一律住子包。留这个小节是为了
> 「根包新增导出符号」有路可登记——没有它，那条路是不存在的（CT-01 自审修的正是这一处）。

### `dataprotocol/` — 包 `dataprotocol`（数据面共享协议形状；2026-10-24 QT-40 X1a 新建）

> **不是冻结面**：七字段是今天三份副本（data 镜像 / core source / core storage）的交集再加 `VisibleAt`，语义见 §5.5。
> 迁移期 core / data 仍各留副本，**删副本与改 import 归各自单**；契约侧的字段集从本表起是唯一共享定义。
> **零 I/O、零时钟**：包内无构造器，`IngestTime` 必须由调用方注入（ADR-0004）。
> `DataType*` 的**字面量就是表路由**：改一个值 = 数据静默落错表，不是改名问题。

| 符号 | 类别 | 稳定性 | 错误语义 / 判据 |
|---|---|---|---|
| `UnifiedDataPoint` | type | stable | 七字段：`Symbol` / `TradeTime` / `Source` / `DataType` / `Data` / `IngestTime` / `VisibleAt`。字段顺序与 data 镜像一致（位置化字面量按此书写）；`VisibleAt` 零值 = 闸门未填，**不得**当成 `TradeTime` 的替代 |
| `UnifiedDataPoint.ToJSON` | method | stable | `json.Marshal(Data)`；`nil` Data → `null`；不可序列化的值原样返错（不吞、不替换成空对象） |
| `UnifiedDataPoint.DeduplicateKey` | method | stable | `Symbol\|DataType\|TradeTime`（UTC，RFC3339Nano）。**Source 刻意不在键里**（fallback 链可能多源同一数据，留哪条由调用方裁决）；时区先归 UTC，同一时刻不同时区同键 |
| `DataTypeOHLCDaily` | const | stable | `"ohlcv_daily"`：日线 OHLCV |
| `DataTypeOHLCMinute` | const | stable | `"ohlcv_minute"`：分钟线 |
| `DataTypeRealtime` | const | stable | `"realtime_quote"`：实时快照 |
| `DataTypeCapitalFlow` | const | stable | `"capital_flow"`：资金流（同表多源口径靠 `Source` 区分，DBT-09） |
| `DataTypeFundamental` | const | stable | `"fundamentals"`：基本面 |
| `DataTypeSectors` | const | stable | `"sectors"`：板块列表 |
| `DataTypeStockSector` | const | stable | `"stock_sector"`：个股 ↔ 板块归属 |
| `DataTypeTopList` | const | stable | `"top_list"`：榜单（龙虎榜） |
| `DataTypeLimitUpPool` | const | stable | `"limit_up_pool"`：涨停池 |
| `DataTypeAnnounce` | const | stable | `"announcements"`：公告 |
| `DataTypeNews` | const | stable | `"news"`：新闻 |
| `DataTypeHotSearch` | const | stable | `"hot_search"`：热榜 |
| `DataTypeGlobalOHLCV` | const | stable | `"global_ohlcv"`：海外行情 |
| `DataTypeMacro` | const | stable | `"macro_series"`：宏观序列（落 `quant.macro_series`） |
| `DataTypeDailyBasic` | const | stable | `"daily_basic"`：每日行情指标（落 `market_activity`；声明 `after_close`） |

### `fees/` — 包 `fees`（A 股费率模型的数据形状；2026-10-22 自 `tellurion:pkg/fees` 迁入，代码逐字、注释有订正）

| 符号 | 类别 | 稳定性 | 错误语义 / 判据 |
|---|---|---|---|
| `DefaultCommissionRate` | const | **frozen** | `0.0003`（双边佣金，监管上限/历史默认）。**取值即契约**：改一位数，core 每一笔回测成交的成本跟着变，且不会报错 |
| `DefaultStampTaxRate` | const | **frozen** | `0.0005`（仅卖出，2023-08-28 起减半后的档）。AUD-06 修过 2 倍错值——历史教训是「常量错了曲线照样漂亮」 |
| `DefaultStampTaxRateBefore` | const | **frozen** | `0.001`（2008-09-19~2023-08-27 的卖出档）。跨削减日的窗口必须按日取档，见 `StampTaxRateFor` |
| `DefaultTransferFeeRate` | const | **frozen** | `0.00001`（过户费，双边，2022 起） |
| `DefaultMinCommission` | const | **frozen** | `5.0`（**绝对元**，不是比率；只约束佣金那一条线） |
| `DefaultSlippageRate` | const | **frozen** | `0.0001`。**非监管值**，是回测真实性旋钮；置 0 = 确定性测试里关掉滑点 |
| `FixedSlippageRate` | const | **frozen** | `0.001`（"fixed" 模型专用，刻意悲观 10 倍：该模型没有逐单波动率反馈） |
| `StampTaxCutDate` | var | **frozen** | 印花税减半生效日 `2023-08-28`（UTC 零点，只比日期不比墙钟）。⚠️ **它是导出 `var` 不是 `const`**（`time.Time` 做不了 const），进程内可被任何导入方改写；本仓按 QT-05 ④ 没动它，下游**只读**。要 const 化就得换成函数入口 = 破坏性改名，另立单 |
| `StampTaxRateFor` | func | **frozen** | 「按成交日取档」的**唯一**比较点：`asOf < StampTaxCutDate → before`，否则 `current`；`asOf` 零值 → `current`（给不出日期的调用方拿今天的规则，不静默拿历史档）；`before`/`current` 非正 → 回落包常量。判据由 `TestStampTaxRateFor_DateSegmented` / `_ZeroAndNegativeFallBack` 钉住 |
| `AShareFees` | type | stable | 五个费率旋钮的聚合（四个比率 + `MinCommission` 绝对元）。⚠️ **加字段必须先定零值语义**：零值在费率路径上不是「未设置」而是「**免费**」——实测 `tellurion:pkg/portfolio/portfolio.go:45` 的 `ComputeFees` 对零值 `AShareFees` 返回全 0 三项，且**不调** `ApplyDefaults`；两个生产消费方（`tellurion:pkg/backtest/tracker/tracker.go:141` 的 `feeSchedule(asOf)`、`tellurion:pkg/live/mock_trader.go:180` 的 `feeSchedule()`）各自按自己的配置逐字段构造，也不调它。所以加一个默认为 0 的新字段 = 让成交静默变便宜，**没有任何一层会报错** |
| `DefaultAShareFees` | func | stable | 返回监管默认档；等价于五个 `Default*` 常量装进结构体（`TestConstantsAreWiredToAShareFees` 逐位对拍） |
| `AShareFees.ApplyDefaults` | method | stable | **逐字段**填零值，不是整体替换；**幂等**（三条测试钉住：填零 / 保留覆写 / 全填时是 no-op）。**负值不清零**：负费率是调用方 bug，要它显形为错 P&L，不要静默回落。<br>⚠️ **今天零生产调用方**（全仓 grep：只有它自己的 `_test.go`）——它是 `MockTrader` 旧内联默认值（`if cfg.CommissionRate <= 0 {…}`）的迁移 shim。别因此当它可删：删掉就把"给不出默认值的构造方"逼回各自内联。<br>同型教训的正确出处是 **AUD-37**，但那条修的是 `tellurion:pkg/backtest/contracts` 的 `TradingConfig.WithDefaults`（整体替换会犯「我只改了一行 yaml，其余全被重置」的错），**不是本函数**——两者形状相同、来历不同，引用时别混 |
| `AShareFees.Validate` | method | stable | 区间即判据：三个比率 ∈ **[0, 0.05)**、`MinCommission` ∈ [0, 1000)、`SlippageRate` ∈ [0, 0.1)；越界返回 error（**不是** `AppError`，本包不依赖 kit）。**0 是刻意合法值**（免佣 VIP / 确定性测试关滑点），由 `TestValidate_AcceptsBoundary` 钉住——⚠️ 源码原注释曾写 `(0, 0.05)` 且声称「被 `OrderManager` 在下单前调用」，**两句都与实现不符**（实测零生产调用方，只有本包测试）；本仓已订正注释为 `[0, 0.05)` + 「尚未接线」（2026-10-22，仅注释、代码行零差异；core 旧包随 sed 删除）。另：**零值 `AShareFees{}` 也能通过 `Validate`**（全 0 在区间内），所以"先 Validate 再 ApplyDefaults"拦不住"整组费率忘配"——那是 `ApplyDefaults` 的活，不是它的。**放宽区间属语义变更**（原本被拒的配置会通过），按版本规则记账 |

### `domain/market/` — 包 `market`（行情值对象 + 频率/时间语义词表；2026-10-22 自 `tellurion:pkg/domain/market` 逐字迁入）

> **两条纪律决定这里的稳定性标注**：
> ① **字面量契约（跨语言）**：`Frequency` / `FrequencySpec` / `TimeSemantics` 的取值与 nt-sidecar `sidecar/mapping.py` 的常量逐字相同。**改一个已有字面量 = 同时改 Python 侧 = breaking**；**加一个词 = MINOR**（旧值照常解析）。
> ② **SPEC 互锁**：`OHLCV` / `Fundamental` 的字段集与 `Frequency` 词表由 core 的 `pkg/domain/spec_drift_test.go` 与 `tellurion:docs/SPEC.md` 双向比对（Go 有而 SPEC 没写 = 红；SPEC 写了而 Go 没有 = 红）。**那条测试刻意留 core**（它读 core 的 `docs/SPEC.md`），core 侧改依赖后它改为 import 本包 ⇒ **在本包给这两个类型加字段，core 侧会红**，这是设计意图，不是漏配。
> **单位口径（B 类规范化数据）**：`Volume` = **股**、`Turnover` = **CNY**；「缺失 ≠ 0」用指针表达（nil = 上游未发布）。
> **零值口径**：`Frequency` / `VisibleAt` 的零值都表示「**未声明**」，不是「默认」——未声明的 PIT 数据按**不可见**处理（fail-closed）。

| 符号 | 类别 | 稳定性 | 错误语义 / 判据 |
|---|---|---|---|
| `Frequency` | type | **frozen** | 记录自带频率字段（ADR-031 §4.1）。字面量 `"1m"` / `"15m"` / `"1d"`；零值 `""` = **未声明**，不是合法频率。扩展 = 常量 + 路由分支 + 聚合口径 + SPEC 词表**四件套**同步（见 `tellurion:docs/SPEC.md` 规范层 Schema 节） |
| `Freq1m` | const | **frozen** | `"1m"`（ohlcv_minute，mootdx）。与 nt-sidecar 的 `FREQ_1D` 等同族常量逐字对齐 |
| `Freq15m` | const | **frozen** | `"15m"`（由 1m 聚合产出，**不物化**） |
| `Freq1d` | const | **frozen** | `"1d"`（ohlcv_daily_qfq，Tushare 前复权——唯一进回测的口径） |
| `Frequency.Valid` | method | stable | 是否在词表内；零值 → false |
| `Frequency.String` | method | stable | 原样返回字面量，供日志/错误信息 |
| `Frequency.IsIntraday` | method | stable | `1m` / `15m` → true，`1d` 与零值 → false |
| `Frequency.IntradayMinutes` | method | stable | `(1,true)` / `(15,true)`；日频或未声明 → `(0,false)`（**不猜默认周期**） |
| `ParseFrequency` | func | stable | 接受集 = `Frequency` 词表；未知 ⇒ `fmt.Errorf("market: unknown frequency …")`。**不静默回退到"默认"频率**（回测里的"默认"就是泄漏） |
| `KnownFrequencies` | func | stable | 返回词表**拷贝**（防调用方改内部状态）。消费方（频率枚举 / 校验 / 防漂移测试）只读它，不各自维护清单 |
| `FrequencySpec` | type | **frozen** | connector 注册契约（D2）的「源口径」枚举，与记录自带的 `Frequency` **正交**：前者说"这个源怎么产出"，后者说"这条记录是什么频率"。零值 = 未声明 |
| `FreqSpecDaily` | const | **frozen** | `"daily"`（日线 / 资金流 / 龙虎榜等每日一条） |
| `FreqSpecIntraday` | const | **frozen** | `"intraday"`（1m bar / 实时快照） |
| `FreqSpecEvent` | const | **frozen** | `"event"`（披露 / 发布驱动的非等间隔记录；可见时间由记录自带的锚点承载） |
| `FrequencySpec.Valid` | method | stable | 是否在词表内；零值 → false |
| `FrequencySpec.String` | method | stable | 原样返回字面量 |
| `ParseFrequencySpec` | func | stable | 接受集 = 词表；未知 ⇒ error，**不回退默认口径** |
| `KnownFrequencySpecs` | func | stable | 词表**拷贝** |
| `TimeSemantics` | type | **frozen** | 「该数据源的数据什么时候才可见」的注册契约字段（ADR-031 §4.2 红线：每条数据必须自带可见时间，否则即泄漏）。字面量与 nt-sidecar 对齐 |
| `VisibleOnAnnouncementDay` | const | **frozen** | `"announcement_day"`（披露日 ≤ 回测日才可用：财报 / 公告 / 宏观发布） |
| `VisibleAfterClose` | const | **frozen** | `"after_close"`（交易日 15:00 后：日线 / 资金流 / 龙虎榜）。**保守**口径，不早于真实可见时刻 |
| `VisibleNextOpen` | const | **frozen** | `"next_open"`（次日开盘；用于收盘后仍拿不到的数据） |
| `VisibleCustom` | const | **frozen** | `"custom"`（必须随声明给出 `VisibleExpr`；实时快照 = 观察时刻，分钟 bar = 周期结束） |
| `TimeSemantics.Valid` | method | stable | 是否在词表内；零值 → false |
| `TimeSemantics.RequiresExpression` | method | stable | 仅 `VisibleCustom` → true；别的语义给了表达式也不报错（**由注册侧判**，本方法只回答"必需吗"） |
| `TimeSemantics.String` | method | stable | 原样返回字面量 |
| `ParseTimeSemantics` | func | stable | 接受集 = 词表；未知 ⇒ error（时间语义的"默认"就是数据泄漏，**不回退**） |
| `KnownTimeSemantics` | func | stable | 词表**拷贝**（扩展 = 常量 + 校验 + SPEC 三件套同步） |
| `AsiaShanghai` | **var** | **frozen** | A 股交易时区 = `time.FixedZone("Asia/Shanghai", +8h)`。**故意不用 `LoadLocation`**：1991 年后无夏令时，固定偏移语义等价，且 Windows 开发机 / 精简 CI 容器默认没有 tzdata（`LoadLocation` 会直接失败）。⚠️ 它是**导出 `var`**（`*time.Location` 做不了 const），进程内可被重新赋值 ⇒ 下游只读；要改成不可变入口 = 破坏性改名 |
| `DailyVisibleAt` | func | **frozen** | 日线 bar 可见点 = 交易日 **15:00 Asia/Shanghai**。入参按 UTC 零点解释（PG `DATE` 经 pgx 扫出的形态），**只取日历日期分量**，时段与位置另建——不要把入参的钟点带进来 |
| `IntradayVisibleAt` | func | **frozen** | 日内 bar 可见点 = **bar 周期结束**（1m stamped 10:00 → 10:01；15m bucket 09:30 → 09:45）。可见 ≠ 周期起点，提前一分钟就是前视 |
| `OHLCV` | type | **frozen**（字段集 SPEC 互锁） | 任意频率的 bar。`Frequency` / `VisibleAt` 零值 = **未声明**（遗留取数路径不填）。单位见本节抬头。`LimitUp` / `LimitDown` 是当日涨跌停标记 |
| `OHLCV.VisibleBy` | method | **frozen** | PIT 闸门判定原语（D6）：`VisibleAt` 非零且 ≤ asOf 才算可见；**零值 VisibleAt 恒不可见**（未声明 = 泄漏，不是"永远可见"）。闸门判据与 `VisibleAt` 的填充口径必须同文件同语义——拆开就会漂移 |
| `Stock` | type | stable | 证券主数据。⚠️ `DelistDate` 用指针：**nil = 仍在市**（零值 `time.Time` 是 0001-01-01，会让"今天还在市"被算成"一万年前就退市"）。没有它修不了幸存者偏差（P2-4） |
| `Split` | type | stable | 拆股 / 送股 / 转增事件，供前复权校验。字段名与 JSON tag 是形状契约；**落库列名由 storage 侧映射，本包不持有** |
| `Dividend` | type | stable | 分红事件（公告日 / 登记日 / 派息日 / 现金与股票股利） |
| `Fundamental` | type | **frozen**（字段集 SPEC 互锁） | 基本面快照。**数值字段一律 `*float64`**：缺失 ≠ 0（PE=0 会被读成"白送的股票"）。`Date` = **可用日**（读侧 `COALESCE(ann_date, trade_date)`），`AnnDate` = 真实披露日（nil = 不知道）——三季报 9/30 截止、10/25 披露，不可在 9/30 就可见（P0-1 / P1-4） |
| `FundamentalData` | type | stable | **持久化视图**（ODR-043 D3 soft layering 的 "view" 半边）：在 `Fundamental` 之上带 `ID` / `CreatedAt` / `TsCode` / `EndDate`。指针语义同上 |
| `IndexDaily` | type | stable | **指数**日线，与 `OHLCV` **刻意分开**（DATA-12）：身份列是 `Code`（→ `market.index_daily`），混用会让指数被当股票进票池。除 `Close`/`Volume` 外的数值列是指针（nil = 上游未发布） |
| `IndexConstituent` | type | stable | 指数成分（入 / 出日与权重）。`OutDate` 零值 = 仍在成分内 |
| `Provider` | interface | experimental | 行情数据访问接口（原 `domain.MarketDataProvider`，S7-P3-4 改名；旧别名已按 CT-06 删除）。⚠️ **实测今天 0 实现 / 0 消费**：core 全仓除定义外无引用（`pkg/marketdata` 用的是它**自己的** `Provider`，与这个不是同一个类型），agent 侧 0 命中。所以它标 experimental 而不是 stable——一个没人实现的 port，稳定性无从谈起。X1a（data 服务拆分 / `DataPort`）会重新裁这个形状，届时走 ADR，不许两边各留一份 |
| `Provider.GetOHLCV` | method | experimental | 按 symbol + 区间取 bar，返回 `[]OHLCV`；**今天返回未分类 error**，QT-10 ③-② 落地后缺数据须包成 `ErrDataGap`（§5.4） |
| `Provider.GetFundamental` | method | experimental | 按 symbol + 日期取基本面（PIT：用可用日，不是报告期）；错误分类同上 |
| `Provider.GetStocks` | method | experimental | 按交易所取证券列表（含退市，见 `Stock.DelistDate`）；错误分类同上 |
| `Provider.GetLatestPrice` | method | experimental | 最新价（只给价格，不给时间戳——需要 PIT 判定的场景不要用它）；错误分类同上 |
| `Provider.GetIndexConstituents` | method | experimental | 取指数成分代码；错误分类同上 |

### `domain/` — 包 `domain`（信号 / 订单 / 组合 / 公司行为 / 因子名册；2026-10-22 自 `tellurion:pkg/domain` 迁入）

> **迁入形状**：23 个文件（16 非测试 + 7 测试）自 core 拷入，刻意的改动只有三类：① import 换成本仓路径（`domain/market`、`fees`）；② **删掉 5 个零引用的部署配置类型**（`Config` / `DatabaseConfig` / `RedisConfig` / `ServiceConfig` / `TushareConfig`——core 与 agent 实测 0 引用；`Config` 还带 `Password` 字段，不该进公开库；裁决见 §8 #2 与 [`adr/0002`](adr/0002-config-族的边界处置.md)）；③ 随之删掉只为它们存在、且断言恒真的 `TestConfig_Types`。
> **等价性证据**：剥注释 + 归一化 import 后，20/23 文件与 core 完全相同；`execution_test.go` 只差 gofmt 的 import 顺序；`types.go` / `types_test.go` 只差被删的那两段（−31 / −24 行）。
> **稳定性口径（本节）**：`frozen` = SPEC 互锁（`market.OHLCV` / `market.Fundamental`）；`stable` = 引擎 / 策略 / 风控的领域词汇（改 = MINOR，破坏 = MAJOR）；`experimental` = §8 #2 的遗留形状与跨仓数据族（core + `tellurion-data`）。core-only 的族按 §6 的功能分组标 stable，并在备注里写明「core-only，③-⑨」——两处判断都留痕。
> **别名已删除（CT-06-ALIAS-EXIT，2026-11-11；随下一个 MINOR `v0.10.0`，尚未发版）**：`OHLCV` / `Stock` / `IndexDaily` / `Split` / `Dividend` / `Fundamental` 六个根别名已从本包删除，实体仍在 `domain/market`——需要时直接 `import "github.com/ruoxizhnya/tellurion-contracts/domain/market"` 用同名类型。`FundamentalData` / `IndexConstituent` 两个根别名早前已按 CT-06-B 删除。删除前 11 仓重跑：六名经**根包别名**的真实残留为 0（core / agent 的命中经各仓自有 `pkg/domain`，data 已改 `market.Fundamental`）；证据见 §6.1 的「CT-06-ALIAS-EXIT 已执行」块。
> **方法体边界**：本包里的 `Validate` / `VersionHash` 是**纯函数行为**，随值对象迁入；判据（无 I/O、依赖受限、无跨调用持久状态、不替代服务）见 [`adr/0004`](adr/0004-值对象上的纯函数行为随类型走.md)。
>
> **CT-06-B（2026-11-10，补裁 #12）**：`domain/` 根包里**只有 core 一个编译使用方**的根名已物理迁回 `tellurion:pkg/domain`（core 新建该包承载，market 别名在 core 侧同形重建）。本仓本包因此只剩 ≥2 使用方的形状（当时留的 6 个 market 别名已于 2026-11-11 由 `CT-06-ALIAS-EXIT` 删除）。迁出清单：`Signal` / `Order` / `Position` / `Portfolio` / `RiskMetrics` / `TargetPosition` / `PositionSize` / `StopLossEvent` / `RiskManager` / `OrderType` 族；`BacktestParams`；因子族 `FactorKernel` / `FactorCaliber` / `FactorDefinition` / `FactorRegistration` / `FactorMiningRecord` / `FactorOutput` / `FactorCacheEntry` / `FactorReturn` / `ICEntry` / `BuiltinFactorDefinitions`；core-only 数据族 `Screen*` / `SWIndustryNode` / `RelationIndustryMember` / `SupplyChainEdge` / `IndustryVitality` / `ETFConstituent` / `FundamentalsDetailRow` / `ExecutionConfig`。**0 使用方直接删除、不在 core 重建**：`StrategyConfig`、`Sector` / `StockSector` / `StockIndustry`（11 仓重跑 0 命中）。逐名处置见 §6.1。

| 符号 | 类别 | 稳定性 | 错误语义 / 判据 |
|---|---|---|---|
| `Direction` | type | stable | 交易方向枚举；JSON 字面量是契约（`long` / `short` / `close` / `hold`） |
| `DirectionLong` | const | stable | `"long"`（开多 / 加多） |
| `DirectionShort` | const | stable | `"short"` |
| `DirectionClose` | const | stable | `"close"`（平仓；卖出方向判定印花税时按它） |
| `DirectionHold` | const | stable | `"hold"`（不动作） |
| `MarketRegime` | type | stable | 市场状态（趋势 / 波动 / 情绪分 / 时间戳）。`Trend` / `Volatility` 也是裸字符串 |

> **CT-06 回收（2026-11-07）**：公司行为族（`CorporateAction` / `ActionType` / 5 个实现 / `ActionEngine` / `ApplyResult` / `ApplyOutcome` 共 68 个符号）经 11 仓重跑实测 **0 使用方**，已按 CT-06 ② 删除；旧名别名 `MarketDataProvider` 同批删除。完整清单与 grep 出处见 §6。

**因子名册与注册（§6「因子族」；core-only：agent 侧 0 命中。CT-06-B（2026-11-10）已裁决：本节符号留在本仓，见 §6 执行块）**

| 符号 | 类别 | 稳定性 | 错误语义 / 判据 |
|---|---|---|---|
| `FactorType` | type | stable | 因子名册名（**冻结词表**：`FactorDefinition.Validate` 用它判名，未登记的名字当场失败）。扩展 = 加常量 + 加定义 + 同步 SPEC |
| `FactorMomentum` | const | stable | `"momentum"`（价量类；有确定性默认表达式） |
| `FactorValue` | const | stable | `"value"`（价值类；**给不出**默认表达式——ADR-024 明确失败，不套无关价格表达式） |
| `FactorQuality` | const | stable | `"quality"` |
| `FactorSize` | const | stable | `"size"` |
| `FactorVolatility` | const | stable | `"volatility"` |
| `FactorGrowth` | const | stable | `"growth"` |
| `FactorGrossMarginTrend` | const | stable | `"gross_margin_trend"`（桥 B1 纵向基本面因子） |
| `FactorContractLiabilityRatio` | const | stable | `"contract_liability_ratio"`（桥 B1） |
| `FactorOCFToNetProfit` | const | stable | `"ocf_to_net_profit"`（桥 B1） |
| `FactorROEDuPontLeverage` | const | stable | `"roe_dupont_leverage"`（桥 B1） |
| `FactorInventoryTurnoverDelta` | const | stable | `"inventory_turnover_delta"`（桥 B1） |
| `FactorMarketActivity` | const | stable | `"market_activity"`（D5 新数据源样板：相对自身近 20 日均值的换手率异常度） |
| `ParseFactorType` | func | stable | 词表查找；未命中 → `("", false)`（**不回退默认因子**） |
| `FactorHypothesis` | type | stable | 因子假设（可证伪预测 / 依据 / 来源等），供 AI 挖掘链路记录 |
| `BuiltinFactorHypotheses` | var | stable | 内置因子假设种子（同上） |
| `FactorSourceLiterature` | const | stable | `"literature"`（有公开文献支撑） |
| `FactorSourceSupplyChain` | const | stable | `"supply_chain"`（产业链逻辑推导） |
| `FactorSourceAI` | const | stable | `"ai_hypothesis"`（AI 提出；**必须过可证伪检验**才可用） |
| `FactorSourceAdHoc` | const | stable | `"ad_hoc"`（临时 / 未记录来源；能用但别太当真） |

**回测输入输出（§6「回测输入输出」）**

| 符号 | 类别 | 稳定性 | 错误语义 / 判据 |
|---|---|---|---|
| `BacktestResult` | type | stable | 单次回测结果（指标 + 净值曲线 + 成交列表）。指标口径受 `tellurion:docs/SPEC.md` 冻结节约束。⚠️ **它没有也不该有 `Status`**（三值属响应层 `BacktestResponse`）。`UniverseMaxSize` 是**指针**：nil = 未评估/未知，0 = 明确为空票池（OBS-01：把未知编码成 0 会造出假红灯）；`InvalidReasons` 由引擎用 `CheckValidity` 算一次后挂上，消费方**只读不重算** |
| `PortfolioValue` | type | stable | 净值曲线的一个点（日期 / 总市值 / 现金 / 持仓市值） |
| `Trade` | type | stable | 成交记录（方向 / 数量 / 价格 / 三项费用 / 时间戳 / 未成交与已成交数量） |
| `WalkForwardParams` | type | stable | 滚动窗口参数（训练 / 测试 / 步长 / 最小训练窗） |
| `WalkForwardResult` | type | stable | 单窗口结果（两段回测 + Sharpe + 收益 + 回撤）。⚠️ `OOSvsTrain` 名实不符（REV-20）：当前没有参数优化步骤，它只反映「B 段相对 A 段的衰减」，**不是**过拟合度量；json tag 是历史遗留，改名 = 破坏 API |
| `WalkForwardReport` | type | stable | 多窗口汇总（均值 / 标准差 / 通过率 / 总体是否通过）。⚠️ `OverfitScore` / `ProbNoOverfit` 同样是历史名（REV-20），当前语义是分段衰减的分档映射 |

**跨仓数据族（core + `tellurion-data` 两个使用方 ⇒ 留本仓；标 experimental，下游不要依赖）**

| 符号 | 类别 | 稳定性 | 错误语义 / 判据 |
|---|---|---|---|
| `MacroSeries` | type | experimental | 宏观序列行（指标 / 频率 / 观测日 / 值 / 单位 / 来源 / as-of） |
| `MarketActivity` | type | experimental | 市场活跃度行（换手率 / 自由流通换手率 / 量比，全指针：nil = 上游未给） |

### `engine/` — 包 `contracts`（四契约 / 内核路由 / 成交口径 / canonical 身份 / abi；2026-10-22 自 `tellurion:pkg/backtest/contracts` 逐字迁入）

> **目录叫 `engine/`、包名保持 `contracts`，是刻意的**（QT-05 ③-⑧）：调用点的 `contracts.X` 前缀一行不改，只有 import 路径变；改成 `package engine` 会波及 core 50 处 + agent 5 处。
> **迁入形状**：20 个文件（10 非测试 + 10 测试）逐字搬入，刻意的代码改动**只有 import 路径**（`domain` / `fees` / 测试包自引用）——剥注释 + 归一化 import 后，17/20 文件完全相同，3 个测试文件只差 gofmt 的 import 顺序，**0 处真实代码差异**。
> **留在 core**：`resultspec_test.go`（它读 core 的 `docs/SPEC.md` 做文档互锁；core 侧切换依赖后改为 import 本包）。
> **§8 #1 已裁决（[`adr/0003`](adr/0003-哨兵与-abi-不独立成包.md)）**：`errors/` 与 `abi/` **不独立成包**，随 `engine/` 走。理由不是「改动太大」——实测 `contracts.Err*` 只有 10 个消费文件、`contracts.Host*` 只有 4 个，**而且每个文件都在同时用同包的其它符号**（≥2 个），拆出去没有任何调用点能少 import 一个包，只会多 import 一个。
> **稳定性**：四契约 §5.1 frozen；status 与 invalid code §5.3 frozen；abi 常量族 §7 frozen；错误族 §5.4 stable；其余工具函数 stable。`TradingConfig` / `PriceLimitConfig` 按 [`adr/0002`](adr/0002-config-族的边界处置.md) ③「冻结面伴生名」保留原名。

| 符号 | 类别 | 稳定性 | 错误语义 / 判据 |
|---|---|---|---|
| `EngineRunner` | interface | **frozen** | 跑一次回测的接口（§5.1）：`RunBacktest(ctx, req) (*BacktestResponse, error)`。签名与 fail-closed 语义冻结 |
| `EngineRunner.RunBacktest` | method | **frozen** | status 是唯一成功判据；错误分类见 §5.4 |
| `BacktestRunner` | interface | stable | 另一条**跨仓**回测口（返回 `*domain.BacktestResult`，签名与 `EngineRunner` 不同）：被 agent 的 pipeline / copilot 与 core 的 strategy adapter 实现（agent 实测 5 个文件 import 本包）。加方法 = 破坏所有实现方 |
| `BacktestRunner.RunBacktest` | method | stable | `(ctx, strategyName, stockPool, startDate, endDate) → (*domain.BacktestResult, error)`；字符串日期是这条口的既有形状，别改成 time.Time 而不升版本 |
| `BacktestRequest` | type | **frozen** | `RunSpec` 的承载物（§5.1）：策略 / 票池 / 区间 / 资金 / 频率 / 内核 / 引擎参数 / 信号行 / 成交口径 / **策略引用与参数覆写（QT-20，可选）**。未知名拒绝（`ValidateRunSpec`） |
| `BacktestResponse` | type | **frozen** | `ResultSpec` 的承载物（§5.1）：`status` 三值冻结 + `invalid_reasons` + 指标 + 净值 + 成交 |
| `SignalWeight` | type | stable | `SignalSpec` 的权重形状（§5.1；K15 目标权重流）：`Symbol` / `Weight` / `AsOf`（ISO 日期字符串） |
| `TradingConfig` | type | **frozen** | A 股交易规则集合（费率 / 涨跌停 / 新股天数 / 融券费率）。**名字按 [`adr/0002`](adr/0002-config-族的边界处置.md) ③ 保留**（SPEC 冻结文本的伴生符号） |
| `TradingConfig.WithDefaults` | method | **frozen** | **逐字段**填默认（不是整体替换；AUD-37 的教训：整体替换会犯「只改了一行 yaml，其余全被重置」） |
| `PriceLimitConfig` | type | **frozen** | 涨跌停档位（主板 / ST / 新股）——名字按 ADR-0002 ③ 保留 |
| `DefaultTradingConfig` | func | stable | 默认交易规则（五个费率项来自 `fees.Default*`，涨跌停 / 新股窗口 / 融券费率是本包字面量） |
| `StatusCompleted` | const | **frozen** | `"completed"`：`invalid_reasons` 空、`error` 空（§5.3） |
| `StatusInvalid` | const | **frozen** | `"invalid"`：`invalid_reasons` 必填、`error` 必空；指标照常给出但不得当成功观测（§5.3） |
| `StatusFailed` | const | **frozen** | `"failed"`：`error` 必填、`invalid_reasons` 空（§5.3）。**禁止第四值** |
| `DefaultStampTaxRate` | const | **frozen** | `= fees.DefaultStampTaxRate`（0.0005；2023-08-28 起）。别名而非新字面量——费率单一来源 |
| `DefaultStampTaxRateBefore` | const | **frozen** | `= fees.DefaultStampTaxRateBefore`（0.001；削减日前）。与上一个配对，跨削减日的窗口按日取档 |
| `DefaultMinCommission` | const | **frozen** | `= fees.DefaultMinCommission`（¥5 绝对额，不是比率） |
| `DefaultTransferFeeRate` | const | **frozen** | `= fees.DefaultTransferFeeRate`（0.00001，双边） |
| `DefaultPriceLimitNormal` | const | **frozen** | `0.10`（主板 ±10%）。创业板 / 科创板 ±20%、北交所 ±30% 由 symbol 在决策时派生，**不是**改这个常量 |
| `DefaultPriceLimitST` | const | **frozen** | `0.10`：主板 ST / *ST **2026-07-06 起**的涨跌停（AUD-07：原值 0.05 已过时；创业板/科创板风险警示股仍按板块 ±20%） |
| `DefaultPriceLimitSTBefore` | const | **frozen** | `0.05`：2026-07-06 **之前**的主板 ST 档位。跨规则变更的窗口按日取档 |
| `DefaultPriceLimitNew` | const | **frozen** | `0.20`：新股上市首日涨跌停（创业板 / 科创板口径） |
| `DefaultNewStockDays` | const | **frozen** | `60`：上市后多少天算「新股」 |
| `DefaultShortSellingRate` | const | **frozen** | `0.106`：融券年化费率，按持仓市值逐日计提 |
| `TradingDaysPerYear` | const | **frozen** | `252`：年化 → 日的除数（回测按交易日推进；实盘 margin 模块用 365，**别混**） |
| `FillTiming` | type | **frozen** | 成交时点口径（QT-02）：`next_open` / `same_close`。字面量进 `run_spec_hash`，改名 = 改所有历史 hash |
| `FillTimingNextOpen` | const | **frozen** | `"next_open"`：**默认口径**（信号日 D 的决策在 D+1 开盘成交） |
| `FillTimingSameClose` | const | **frozen** | `"same_close"`：D 日收盘成交（含 T 日信息参与成交的已知偏差，必须用 `KnownBiasSameCloseFill` 标注） |
| `DefaultFillTiming` | const | **frozen** | `= FillTimingNextOpen` |
| `FillTiming.Valid` | method | stable | 是否在两个合法字面量内；零值 → false（零值是**未声明**，不是默认口径） |
| `FillTiming.String` | method | stable | 原样返回字面量 |
| `FillTiming.KnownBiases` | method | stable | `same_close` → `["same_close_fill"]`；`next_open` → nil。**结果必须自带已知偏差**，别让读者自己猜 |
| `KnownBiasSameCloseFill` | const | stable | `"same_close_fill"`：同收盘成交的已知偏差标签 |
| `ResolveFillTiming` | func | stable | 严格解析：空 → `ErrFillTimingMissing`；不认识 → `ErrFillTimingUnknown`；**绝不回落默认值**（默认只在 `EffectiveFillTiming` 那一处决定） |
| `BacktestRequest.EffectiveFillTiming` | method | stable | 空值 → `DefaultFillTiming`；**不写回**请求（回填会改变存量 JSON 序列化） |
| `BacktestRequest.EffectiveKernel` | method | stable | 空值 → `DefaultKernel`；同样**不写回**（存量请求的 JSON 不多出 `"kernel":"go"`） |
| `BacktestRequest.ValidateRunSpec` | method | stable | 返回非 nil = 请求非法（API 层映射 400）：引擎参数未知名、signals 行非法值、**显式声明**但不认识的 fill_timing。空值**不拒**（走默认语义） |
| `DefaultKernel` | const | stable | `"go"` |
| `KnownEngineParams` | func | stable | 该内核已登记的引擎参数名（排序副本）；未登记的内核返回空切片 = 「不接受任何引擎特有参数」 |
| `KernelCaliber` | type | stable | 内核口径对（内核名 + **实际成交**口径），K15 双跑的前置断言输入 |
| `ErrFillTimingMissing` | var | stable | 见 §5.4：必填面缺失（声明为空）——与「不认识」「两侧不符」分开，调用方才知道该改哪边 |
| `ErrFillTimingUnknown` | var | stable | 见 §5.4：给了但不认识（拼错 / 上游新增口径而本内核没实现） |
| `ErrFillTimingMismatch` | var | stable | 见 §5.4：两内核各自声明的口径不一致（K15 对拍前置门；口径不符 ≠ 结果超差） |
| `CompareKernelCalibers` | func | stable | 两侧都取自「实际成交」一侧；任一侧为空 / 不认识 → 错误（`ErrFillTiming*` 包装内核名）。**口径不符 ≠ 结果超差**：前者根本没有可比性 |
| `EngineRegistry` | type | stable | 内核注册表（并发安全）。**零静默回退**：未注册内核必须点名报错 |
| `NewEngineRegistry` | func | stable | 构造空注册表 |
| `EngineRegistry.Register` | method | stable | 装配期登记（**无热插拔**）。Must 风格：**空名 / nil runner / 重名 → panic**（装配期编程错误，fail-fast；panic 前先解锁，残留状态仍可用）。注意 typed-nil 陷阱：`Register(name, (*Engine)(nil))` 过不了接口 nil 检查，会在首次 `RunBacktest` 才崩 |
| `EngineRegistry.Resolve` | method | stable | 取内核；未注册 → `ErrKernelNotRegistered`（不回落 `go`） |
| `EngineRegistry.Names` | method | stable | 已注册内核名（排序副本），供诊断 / 诊断端点 |
| `ErrKernelNotRegistered` | var | stable | 见 §5.4：请求期 400；零静默回退 |
| `MarketRouter` | type | stable | 市场 → 内核路由器（配置驱动）。**零值不可用**（nil 表一律拒绝，请用构造器）。市场名构造时归一（去空白 + 大写：`US` 与 `us` 是同一条路由） |
| `NewMarketRouter` | func | stable | 由映射构造；**空表合法**（存量部署没有 `kernel_routing` 配置时，不带市场的请求行为必须逐字节不变）。拒绝并报错：市场名为空、内核名为空、大小写归一并冲突（同一条路由两个答案——挑一个就是静默回退） |
| `NewMarketRouterFromConfig` | func | stable | 配置读取包装；nil / 缺配置节 = 空表（不报错） |
| `MarketRouter.ResolveKernel` | method | stable | 按市场（归一后）取内核；未命中 / 空名 / nil 表 → `ErrMarketNotRegistered`，错误里点名「被拒的市场 + 已登记集合 + 不回落默认内核」——**绝不猜默认内核** |
| `MarketRouter.Markets` | method | stable | 已路由的市场名（排序副本） |
| `ErrMarketNotRegistered` | var | stable | 见 §5.4 |
| `ResolveRunTarget` | func | stable | 市场 + 显式内核 → 生效内核名。四种组合：都不给 → `""`（走默认语义，**不回填**）；只给内核 → 原样；只给市场 → 查路由；两者都给 → **必须一致**，不一致 → 错误（拒绝替调用方挑一个）。**不校验内核是否已注册**（那是 `EngineRegistry.Resolve` 的单一职责） |
| `CheckValidity` | func | stable | 「无效运行」的**唯一**判定实现（OBS-01）：返回 `InvalidReason` 列表；消费方读结果里的 `InvalidReasons`，不各自重算 |
| `InvalidReason` | type | **frozen** | 无效原因（Code + Detail）；**Code 是稳定集，Detail 才可改** |
| `InvalidReasonCodes` | func | **frozen** | 全部合法 Code（排序副本），供文档 / 前端枚举 |
| `InvalidEmptyUniverse` | const | **frozen** | `empty_universe` |
| `InvalidZeroTrades` | const | **frozen** | `zero_trades` |
| `InvalidZeroStartDate` | const | **frozen** | `zero_start_date` |
| `InvalidGarbageMetric` | const | **frozen** | `garbage_metric`（NaN / Inf 等） |
| `InvalidExecAlgoUnknownStyle` | const | **frozen** | `exec_algo_unknown_style` |
| `InvalidExecAlgoAttachFailed` | const | **frozen** | `exec_algo_attach_failed` |
| `CodeDataGapSkipped` | const | **frozen** | QT-10 ③-④：`data_gap_skipped`——`ErrDataGap` 跳天的落账 Code，走**既有** `invalid_reasons` 通道。⚠️ 它是 Code 不是状态：`status` 三值仍是 frozen（§5.3） |
| `Canonical` | type | stable | 确定性编码器（链式 API，错误在内部累积，`Hex` 时一次性返回）。输出是 hash 的字节基础：**编码顺序写死，改顺序 = 所有历史 hash 失效** |
| `NewCanonical` | func | stable | 开启一段编码，先写域标签（`Domain*Hash`） |
| `Canonical.Str` | method | stable | 字符串字段 |
| `Canonical.Int` | method | stable | int64 字段 |
| `Canonical.Float` | method | stable | 浮点字段：**NaN / Inf 不可 canonical 编码** → `ErrCanonicalUnsupported` |
| `Canonical.Time` | method | stable | 时间字段（UTC 日历日编码） |
| `Canonical.Sequence` | method | stable | **保序**序列（如交易日轴） |
| `Canonical.SymbolSet` | method | stable | **集合**语义（排序 + 去重；票池用它）。与 `Sequence` 的区别是刻意的，别互换 |
| `Canonical.Value` | method | stable | 任意值（走 canonicalValue）；不支持的类型 → `ErrCanonicalUnsupported` |
| `Canonical.Hex` | method | stable | 返回 sha256 hex；内部有错误则返回该错误 |
| `Canonical.MustHex` | method | stable | 只用于**已证明可编码**的调用点；内部错误 → panic（不是丢字段） |
| `ErrCanonicalUnsupported` | var | stable | 见 §5.4：不可编码 → 当场失败，禁止退化成「猜一个 hash」 |
| `DomainRunSpecHash` | const | **frozen** | `"tellurion.run_spec_hash.v1"` 域标签（升版本 = 所有历史 hash 失效） |
| `DomainDatasetSnapshotHash` | const | **frozen** | `"tellurion.dataset_snapshot_hash.v1"` |
| `DomainUniverseHash` | const | **frozen** | `"tellurion.universe_hash.v1"` |
| `DomainCalendarHash` | const | **frozen** | `"tellurion.calendar_hash.v1"` |
| `DomainAnchorDigestsHash` | const | **frozen** | `"tellurion.adjustment_anchor_digests.v1"` |
| `DomainParamsHash` | const | **frozen** | `"tellurion.params_hash.v1"`：**生效参数**的独立指纹，与 `run_spec_hash` 里那份**不同域**（只比参数时别误命中换了区间/票池的 run） |
| `DomainStrategyVersionHash` | const | **frozen** | `"tellurion.strategy_version_hash.v1"`（QT-22）：`version_hash` 的域标签，覆盖**六项** `strategy_id` + `carrier` + `abi_version` + `scope` + `entry` + `source`。与 `DomainParamsHash` **不同域，禁止复用**。改标签 = 所有历史 version_hash 失效（它同时是策略资产四张表的键与 `run_spec_hash` 的分量）⇒ 加项必须升 `v2` 并走双标签过渡；`capability` / `params_defaults` **不在集合内** |
| `DomainRunIdentityHash` | const | **frozen** | `"tellurion.run_identity_hash.v1"`（QT-16）：检查点身份复合域，覆盖 `{run_universe_hash, scope, shard_key, param_effective_hash}` 四分量；计算函数在 core `pkg/strategy`（消费方只有 core，不进本仓） |
| `ArtifactHash` | func | stable | 字节码 sha256（hex），**刻意不加域标签**：要能被编译服务 / 缓存 / 外部工具独立算出且逐字节相同（内容寻址）。执行与结果以它为真相（N7） |
| `UniverseHash` | func | stable | 票池**集合**指纹（排序 + 去重）：`["A","B"]` 与 `["B","A"]` 同 hash |
| `CalendarHash` | func | stable | 交易日轴指纹，**保序**（`[d1,d2]` ≠ `[d2,d1]`）。Axis 必须进 dataset hash，否则解释不了「同 hash 不同结果」 |
| `AdjustmentAnchorDigestsHash` | func | stable | 复权锚点 digest 集合指纹（按 `(symbol,start_date)` 排序）：属 **run** 分量，不是 dataset 分量 |
| `DatasetSnapshotInput` | type | stable | dataset hash 的全部输入（票池 / 区间 / 频率 / 数据版本 / 交易日轴 / 复权口径） |
| `DatasetSnapshotInput.Hash` | method | stable | 产出 `dataset_snapshot_hash` |
| `DateRange` | type | stable | 闭区间日期对（UTC 日历日参与编码） |
| `RunSpecInput` | type | **frozen** | 一次回测的**全部决定性输入**（N6-1）：**影响结果的就必须是字段**（加一个不进 hash 的输入 = 造出「同 hash 不同结果」）。空串/零值是合法输入（哈希保持全函数），必填性由各契约自己管 |
| `RunSpecInput.Hash` | method | **frozen** | 产出 `run_spec_hash`；不可编码的值返回错误而**不是跳过字段**（跳过的字段改了也不变 hash = 会撒谎的留痕） |
| `RunSpecInput.MustHash` | method | stable | 已证明可编码时用；错误 → panic |
| `ParamsHash` | func | stable | 生效参数指纹（QT-63 的 `params_hash` 列） |
| `StrategyVersionInput` | type | stable | QT-22：`version_hash` 的**全部决定性输入**，六字段 `StrategyID` / `Carrier` / `ABIVersion` / `Scope`（`StrategyScope`）/ `Entry` / `Source`。与 `tellurion-strategy:internal/asset/hash.go` 的 `VersionInput` **逐字段同序**（两仓同输入必须同 hex）。形状稳定；**加字段 = 换语义 = 升 `v2`**（0.x 阶段按破坏性留账） |
| `StrategyVersionInput.Hash` | method | stable | 产出 `version_hash`（`DomainStrategyVersionHash` 域）；全字段字符串，编码不会失败，但仍返回 error 以统一调用形状 |
| `StrategyVersionInput.MustHash` | method | stable | 已证明可编码时用；错误 → panic（与 `RunSpecInput.MustHash` 同型） |
| `StrategyRef` | type | stable | 策略引用（`name@hash` 形态；`Hash=""` = 最新） |
| `StrategyRef.IsLatest` | method | stable | `Hash == ""` → true |
| `StrategyRef.String` | method | stable | 还原 `name@hash` |
| `ParseStrategyRef` | func | stable | 解析 `name@hash`；格式不对 → `ErrStrategyRefMalformed`（artifact hex 段必须 64 位且**必须小写**——大写会被拒，避免同一 hash 出现两个引用） |
| `StrategyRefSeparator` | const | **frozen** | `"@"`：`strategy_ref` 的分隔符，改它 = 破坏所有已落库的 ref |
| `ErrStrategyRefMalformed` | var | stable | 见 §5.4 |

**abi 面（`engine/abi.go`；host 函数名 / 版本 / 内存区所有权按 §7 冻结）**

| 符号 | 类别 | 稳定性 | 错误语义 / 判据 |
|---|---|---|---|
| `HostFnGetBar` | const | **frozen** | v1：`(i32 bar_idx, i32 field_id) → f64`。见 §7 |
| `HostFnGetSeriesLen` | const | **frozen** | v1：`(i32) → i32` |
| `HostFnGetSymbolCount` | const | **frozen** | v1：`() → i32` |
| `HostFnGetCrossSection` | const | **frozen** | v1：`(i32,i32,i32) → f64` |
| `HostFnGetState` | const | **frozen** | v1：`(i32 key, i32 len) → i64`；traps 见 §7 |
| `HostFnSetState` | const | **frozen** | v1：`(i32,i32,i32,i32)`；traps 见 §7 |
| `HostFnEmitSignal` | const | **frozen** | v1：`(i32 symbol_idx, i32 action, f64 strength)`；traps 见 §7 |
| `HostFnLog` | const | **frozen** | v1：`(i32,i32,i32)` |
| `HostFnGetBar2` | const | **frozen** | v2 新增：`(i32,i32,i32) → f64`（带频率 / 口径参数） |
| `HostFnGetSeriesLen2` | const | **frozen** | v2 新增：`(i32) → i32` |
| `HostFnGetDate2` | const | **frozen** | v2 新增：`() → i32`（yyyymmdd） |
| `HostFnEmitTarget2` | const | **frozen** | v2 新增：目标权重流；与 `emit_signal` **同 bar 同 symbol 混用即错**（`ErrEmitSourcesMixed`） |
| `HostFnCount` | func | **frozen** | 返回 `len(hostFnContracts)`（枚举结果），**不是**文档字面数；今天是 **12** |
| `HostFnContract` | type | **frozen** | 单条 host 函数契约（名 / Since / Params / Results / Traps / Doc） |
| `HostFnContract.Signature` | method | **frozen** | 签名字符串（导入表比对用）；artifact 导入表不符 → `ErrImportSignatureMismatch` |
| `HostABIContracts` | func | **frozen** | 全部契约（深拷贝，防调用方改内部表） |
| `HostABIContract` | func | **frozen** | 按名字取一条；未知 → `(零值, false)` |
| `HostAbiVersion` | type | **frozen** | abi 版本（1 / 2） |
| `HostAbiV1` | const | **frozen** | `1` |
| `HostAbiV2` | const | **frozen** | `2` |
| `HostAbiVersion.Valid` | method | **frozen** | 只认 1 / 2 |
| `HostAbiVersion.String` | method | **frozen** | `"v1"` / `"v2"` / `"unknown"` |
| `HostAbiVersionFromName` | func | **frozen** | 解析 `on_bar/N` 的 N；失败 → `ErrAbiVersionUnknown` |
| `HostValueType` | type | **frozen** | host 值类型（i32 / i64 / f64） |
| `HostI32` | const | **frozen** | 枚举值 1 |
| `HostI64` | const | **frozen** | 枚举值 2 |
| `HostF64` | const | **frozen** | 枚举值 3 |
| `HostValueType.String` | method | **frozen** | `"i32"` / `"i64"` / `"f64"` |
| `FieldID` | type | **frozen** | bar 字段索引（abi 里的整数） |
| `FieldClose` | const | **frozen** | `0` ——**索引即契约**，换序 = 所有 artifact 静默读错列 |
| `FieldOpen` | const | **frozen** | `1` |
| `FieldHigh` | const | **frozen** | `2` |
| `FieldLow` | const | **frozen** | `3` |
| `FieldVolume` | const | **frozen** | `4` |
| `FieldID.Valid` | method | **frozen** | 是否在 0..4 |
| `FieldIDs` | func | **frozen** | 全部字段（按索引升序） |
| `Action` | type | **frozen** | 信号动作（abi 里的整数值） |
| `ActionHold` | const | **frozen** | `0` |
| `ActionLong` | const | **frozen** | `1` |
| `ActionClose` | const | **frozen** | `2` |
| `ActionShort` | const | **frozen** | `3` |
| `Action.Valid` | method | **frozen** | 是否在 0..3 |
| `Action.Direction` | method | **frozen** | 映射到 `domain.Direction`；非法值返回 **error**（调用方 trap，**不回落 hold**） |
| `Actions` | func | **frozen** | 全部动作（按值升序） |
| `StrategyScope` | type | **frozen** | 信号作用域（market / symbol） |
| `ScopeMarket` | const | **frozen** | `"market"`：1 实例/run，横截面，Axis 为日期轴 |
| `ScopeSymbol` | const | **frozen** | `"symbol"`：N 实例/run |
| `StrategyScope.Valid` | method | **frozen** | 是否在两个字面量内 |
| `ResolveScope` | func | **frozen** | 按 abi 版本裁定生效作用域：**v1 只能是 symbol**（v1 协议表达不了横截面）——`""` 按 symbol 补（兼容期唯一允许的「空值有意义」），声明 market ⇒ `ErrScopeMarketV1`；**v2 必须显式声明**且 ∈ {market, symbol}，否则 `ErrScopeUnknown`；版本不认识 ⇒ `ErrAbiVersionUnknown` |
| `EmitSource` | type | stable | 发射来源（signal / target）——用于同 bar 混用检测 |
| `EmitSourceSignal` | const | stable | `"signal"` |
| `EmitSourceTarget` | const | stable | `"target"` |
| `EmitSource.String` | method | stable | 原样返回字面量 |
| `BarEmission` | type | stable | 一根 bar 上的一次发射（来源 / symbol / 值）；预算与混用检查的输入 |
| `CheckEmitSources` | func | stable | 同 bar 同 symbol 混用两种来源 → `ErrEmitSourcesMixed`；未知来源 → `ErrEmitSourceUnknown` |
| `CheckSignalBudget` | func | stable | 每 bar 发射数超 `MaxSignalsPerBar(symbolCount)` → `ErrSignalBudgetExceeded` |
| `CheckStateBudget` | func | stable | 输入是「**本次写入之后**」的 key 数 / 总字节（调用方负责先算再写）：超 `MaxStateKeys` / `MaxInstanceStateBytes` → `ErrStateBudgetExceeded` |
| `ValidateStrength` | func | stable | 非 finite 或越界 → `ErrStrengthOutOfRange`（非法信号 = 策略故障） |
| `ValidateTargetWeight` | func | stable | 非 finite 或越界 → `ErrTargetWeightOutOfRange` |
| `MaxParamsBytes` | const | **frozen** | `4096`：params 区上限（截断**不允许** → `ErrParamsTooLarge`） |
| `MaxStateValueBytes` | const | **frozen** | `64 << 10`：单 state 值上限 |
| `MaxInstanceStateBytes` | const | **frozen** | `4 << 20`：单实例 state 总字节上限 |
| `MaxStateKeys` | const | **frozen** | `4096`：state key 数上限 |
| `MaxScratchPerBar` | const | **frozen** | `256 << 10`：单 bar scratch 上限 |
| `MaxSignalsPerBar` | func | **frozen** | 与 symbolCount 相关的每 bar 发射上限（函数而非常量） |
| `HostRegionAlignment` | const | **frozen** | `8`：内存区指针对齐 |
| `MinHostParamsPtr` | const | **frozen** | `= MaxParamsBytes`：params 区最小起始地址 |
| `MinHostScratchSize` | const | **frozen** | `= MaxStateValueBytes`：scratch 区最小尺寸 |
| `ExportHostParamsPtr` | const | **frozen** | `"__host_params_ptr"` 导出名（artifact 必须导出） |
| `ExportHostScratchPtr` | const | **frozen** | `"__host_scratch_ptr"` |
| `ExportHostScratchSize` | const | **frozen** | `"__host_scratch_size"` |
| `HostRegionExports` | func | **frozen** | 三个必须存在的导出名；缺失 → `ErrHostRegionMissing` |
| `HostRegionDeclaration` | type | **frozen** | 宿主内存区声明（params ptr / scratch ptr / scratch size） |
| `HostRegionDeclaration.Validate` | method | **frozen** | 对齐、最小地址、最小尺寸逐条校验；不合法 → `ErrHostRegionInvalid` |
| `HostRegionDeclaration.CheckAgainstStaticData` | method | stable | QT-62：宿主区 × 模块静态数据（Data section 活动段并集）重叠判定；重叠 → `ErrHostRegionOverlapsData`，图像不可证 → `ErrHostRegionDataUnprovable`。零长段不算重叠；区间算术 uint64 不绕回 |
| `StaticDataSegment` | type | stable | QT-62：模块 Data section 里一个活动段的常量区间（`Offset` / `Length`，uint32） |
| `StaticDataImage` | type | stable | QT-62：从 artifact 字节解出的静态数据布局；`Unprovable` 非空 = 存在无法证明不重叠的段，调用方必须 fail-closed，不许当空图放行 |
| `ValidateParamsWrite` | func | stable | 写 params 前的边界校验：`paramsLen < 0` 或 `paramsLen > min(MaxParamsBytes, scratch_ptr−params_ptr)` → `ErrParamsTooLarge`。**截断是明确禁止的**（截断等于让策略读到半个 JSON）；阈值本身是 frozen 常量 |
| `ErrAbiVersionUnknown` | var | stable | 见 §5.4（abi / artifact 校验族） |
| `ErrAbiVersionMismatch` | var | stable | 见 §5.4 |
| `ErrAbiFamilyMissing` | var | stable | 见 §5.4 |
| `ErrAbiFamilyMixed` | var | stable | 见 §5.4 |
| `ErrImportNotWhitelisted` | var | stable | 见 §5.4 |
| `ErrImportSignatureMismatch` | var | stable | 见 §5.4 |
| `ErrHostRegionMissing` | var | stable | 见 §5.4 |
| `ErrHostRegionInvalid` | var | stable | 见 §5.4 |
| `ErrHostRegionOverlapsData` | var | stable | QT-62：宿主区落在模块自己的静态数据上 |
| `ErrHostRegionDataUnprovable` | var | stable | QT-62：无法从 artifact 字节证明宿主区不踩模块数据（被动段 / 非常量 offset / 非零 memory index）——读不懂即拒绝 |
| `ErrParamsTooLarge` | var | stable | 见 §5.4 |
| `ErrScopeUnknown` | var | stable | 见 §5.4 |
| `ErrScopeMarketV1` | var | stable | 见 §5.4；解禁要立 ADR（D25） |
| `ErrActionIllegal` | var | stable | 见 §5.4（信号校验族，归策略故障） |
| `ErrStrengthOutOfRange` | var | stable | 见 §5.4 |
| `ErrTargetWeightOutOfRange` | var | stable | 见 §5.4 |
| `ErrEmitSourcesMixed` | var | stable | 见 §5.4 |
| `ErrEmitSourceUnknown` | var | stable | 见 §5.4 |
| `ErrSignalBudgetExceeded` | var | stable | 见 §5.4 |
| `ErrStateBudgetExceeded` | var | stable | 见 §5.4 |
| `ErrDataGap` | var | stable | QT-10 ③-① 三分法之首：数据缺口 → **跳过当天** + `invalid_reasons` 写 `CodeDataGapSkipped`，run 继续、终态 `invalid`。三类里唯一不终止 run 的（§5.4） |
| `ErrStrategyFault` | var | stable | QT-10：策略故障（wasm trap / 非法信号 / 预算超限 / artifact 校验失败）→ **终止 run** `status=failed` + `error`；禁止降级成「跳过这天继续」（§5.4） |
| `ErrContractBreach` | var | stable | QT-10：契约违规（前视 / 未知参数 / abi 不符）→ **终止 run** `status=failed`；同时是「无法归类」的兜底类（§3.4） |
| `ErrCompileFailed` | var | stable | QT-41b ④：源码编译失败（跨仓可识别；包装须带 toolchain + 日志）。与 `ErrCompileUnavailable` 分开：一个改源码，一个改配置（§5.4） |
| `ErrCompileUnavailable` | var | stable | QT-41b ④：资产服务未配置编译端口（可重试 / 需报警，不是源码问题） |
| `ErrSourceTooLarge` | var | stable | QT-41b ④：上传源码超过 1 MiB（不可重试，改输入） |
| `ErrArtifactTooLarge` | var | stable | QT-41b ④：artifact 超过 8 MiB blob 上限（fail-loud，不截断；上限本身属资产服务配置） |

### `risk/` — 包 `risk`（风控规则快照的形状 + 冻结词表 + hash 纯函数；2026-11-07 CT-03 新建，`v0.7.0`）

> **编码逐字节冻结**：`Snapshot.Hash()` 与 `tellurion-risk:internal/rule/snapshot.go` 的 `SnapshotHash`
> 必须同 hex；算法细节见 §5.6。**禁止**改用 `engine.NewCanonical`（帧结构、域标签、参数规范化三处都不同，
> 换编码 = 全历史 `risk_rules_hash` 作废）。
> **不迁**（按 CT-03 ④）：`Bundle` / `ParseBundle` / `ForScope` 带 decode、默认值填充与进程内索引，属服务策略；
> 服务私有哨兵（`ErrRuleNotFound` 等 7 个）只经 wire code 消费，也不进契约。
> **零 I/O、零时钟**：包内无构造器；`CreatedAt` 只是被承载的字段，不参与 hash（ADR-0004）。
> **scope 语义（2026-11-11 主任裁决，正式口径）**：进程级策略源的选择顺序是 `global` → `strategy`
> （global 提供组合级基线，strategy 提供单策略覆盖；回退是**正常优先级语义**，不是兜底）；
> `account` / `symbol` 不参与进程级选择。选中 scope 内**四类规则各恰好一条**，缺一类或同类多条 ⇒ 拒。
> `params` 的严格解码（未知键拒收）在消费方，冻结词典归 `tellurion-risk`；本包只做规范化与体积上限。

| 符号 | 类别 | 稳定性 | 错误语义 / 判据 |
|---|---|---|---|
| `DomainSnapshotHash` | const | **frozen** | `"tellurion.risk_snapshot_hash.v1"`。改它 = 全历史 `risk_rules_hash` 作废（`tellurion-risk:docs/db-schema.md` §5）；要改必须 ADR 级裁决并升 v2 |
| `MaxParamsBytes` | const | **frozen** | `65536`（规则参数 JSON 上限，超限 fail-loud 不截断）。⚠️ **与 `engine/` 的 `MaxParamsBytes`（4096，wasm `initialize` 的 params 区）同名不同值不同义**，禁止为「统一」改成 4096——会拒掉今天所有合法快照 |
| `RuleType` | type | **frozen** | 冻结枚举（`string` 底层）。旧读者遇未知类型必须 fail-closed 拒收；新增取值 = MINOR，改既有语义 = breaking |
| `RuleTypePositionSizing` | const | **frozen** | `"position_sizing"`，与 risk 仓逐字一致 |
| `RuleTypeStopLoss` | const | **frozen** | `"stop_loss"` |
| `RuleTypeRegime` | const | **frozen** | `"regime"` |
| `RuleTypeTakeProfit` | const | **frozen** | `"take_profit"` |
| `Scope` | type | **frozen** | 规则适用范围（`string` 底层）。⚠️ 与 `engine/` 的 `ScopeSymbol`（`StrategyScope` 的取值 `"symbol"`，指每实例绑定一个 symbol）**不是同一个概念**，禁止用类型别名互转 |
| `ScopeGlobal` | const | **frozen** | `"global"` |
| `ScopeAccount` | const | **frozen** | `"account"` |
| `ScopeStrategy` | const | **frozen** | `"strategy"` |
| `ScopeSymbol` | const | **frozen** | `"symbol"` |
| `Scope.Valid` | method | stable | 纯函数：是否在冻结四值内（ADR-0004）。零值 `""` → `false` |
| `Scope.ProcessLevel` | method | stable | 纯函数：该 scope 是否参与进程级策略源选择。`global` / `strategy` → `true`；`account` / `symbol` 及未知值 → `false`（2026-11-11 裁决） |
| `SnapshotRule` | type | stable | 六字段：`RuleID` / `Revision` / `Scope` / `Type` / `Params`（`json.RawMessage`，须为 JSON 对象）/ `VersionHash`。**字段顺序即位置化字面量的顺序** |
| `Snapshot` | type | stable | 三字段：`RiskRulesHash` / `Rules` / `CreatedAt`。内容寻址：同 hash 同字节，原地改内容禁止 |
| `Snapshot.Clone` | method | stable | 深拷贝 `Rules` 与每条 `Params`，避免与 Store / 来源共享可变内存 |
| `Snapshot.Hash` | method | **frozen** | 由 `Rules` 派生 `risk_rules_hash`：按 `rule_id` 升序、长度前缀帧、逐条六字段。**不含** `CreatedAt` 与 `RiskRulesHash`（后者只参与 `Verify`）。错误族：空快照 / 重复 `rule_id` / 缺 `rule_id` / `Revision <= 0` / 缺 `version_hash` / 非法 params → `ErrInvalidParams`；未知 `Type` → `ErrInvalidRuleType`；未知 `Scope` → `ErrInvalidScope` |
| `Snapshot.Verify` | method | stable | 重新 derive 并与 `RiskRulesHash` 比对；不符 → `ErrSnapshotMismatch`（fail-loud，不吞、不回退） |
| `Snapshot.ValidateScope` | method | stable | 纯结构校验：指定 scope 下**四类规则各恰好一条**；缺一类或同类多条 → `ErrInvalidParams`，未知 scope → `ErrInvalidScope`，未知 type → `ErrInvalidRuleType`。不做 decode / 默认值 / 进程内索引；`global → strategy` 的 fallback 选择仍属服务策略（2026-11-11 裁决） |
| `ErrInvalidRuleType` | var | stable | `"rule type is not in the frozen enum"`（文案与 risk 逐字一致） |
| `ErrInvalidScope` | var | stable | `"scope is not in the frozen enum"` |
| `ErrInvalidParams` | var | stable | `"rule params are invalid"` |
| `ErrSnapshotMismatch` | var | stable | `"snapshot hash does not match content"` |

### `run/` — 包 `run`（编排 wire 面三份形状；2026-11-07 CT-04 新建，`v0.8.0`）

> **字段名 / 可空性 / 声明顺序逐字冻结**：三份形状分别对应
> `tellurion-orchestration:internal/httpapi/workers.go` 的 `claimView`、
> `internal/httpapi/provenance.go` 的 `provenanceView`、`internal/httpapi/parents.go` 的提交请求族。
> 改任一处 = breaking（已上 wire 的字节改变）。语义与对拍判据见 §5.7。
> **不迁**（按 CT-04 ④）：`RunRecord` / `ParentRecord` / `WorkerRecord`（队列与租约状态机）、
> `jobView` / `parentView` / `parentChildView` / `parentCountsView`（编排读面投影）、
> `ArtifactStore`（带写语义）、`run_artifacts` 的 DDL。包内**零 I/O、零时钟、零 decode**。

| 符号 | 类别 | 稳定性 | 错误语义 / 判据 |
|---|---|---|---|
| `ClaimView` | type | stable | 7 字段：`RunID` / `RunSpec`（`json.RawMessage`，完整 RunSpec 原样字节）/ `Attempts` / `MaxAttempts` / `LeaseUntil`（`*time.Time`，出 `null`）/ `CancelRequested` / `ParentJobID`（**`omitempty`**：普通 run 不出现该键） |
| `ProvenanceView` | type | stable | 15 字段（`run_id` … `created_at`）。`ArtifactURI` / `ArtifactHash` / `VersionHash` / `Seed` / `ResultStatus` / `CreatedAt` 为指针，**出 `null` 而不是省略键**；`RiskRulesHash` 空串 = 未登记，不是「无规则」 |
| `ParentSubmitRequest` | type | stable | `IdempotencyKey` / `Kind` / `Batch` / `WalkForward`；两个载荷都非 `nil` 或都 `nil` = 歧义请求，由服务侧 400（契约不猜） |
| `ParentBatchRequest` | type | stable | `BaseRunSpec`（`json.RawMessage`）+ `Grid map[string][]json.RawMessage` |
| `ParentWalkForwardRequest` | type | stable | `BaseRunSpec` + `Windows` + `TradingDays []string`（调用方传入的交易日历，契约不自造） |
| `ParentWindows` | type | stable | `TrainDays` / `TestDays` / `StepDays` / `MinTrainDays`（int） |
| `ChildKeyPrefix` | const | **frozen** | `"parent"`：子 run 幂等键的保留前缀；调用方自造键不得使用 |
| `ChildKeySeparator` | const | **frozen** | `":"`：改它 = 全部子 run 的幂等身份改变（MAJOR） |
| `ChildIdempotencyKey` | func | **frozen** | 纯函数 `parent:<parent_job_id>:<child_index>`；逐字节对齐 `tellurion-orchestration:internal/parents/service.go`。改规则 = MAJOR |

### `broker/` — 包 `broker`（订单 / 成交 / 账户 wire 面；2026-11-07 CT-05 新建，`v0.9.0`）

> **字段名 / 可空性 / 声明顺序逐字冻结**：形状逐字取自
> `tellurion-broker:internal/gateway/model.go`，期望字节由生产类型 `json.Marshal` 打出后逐字节抄入
> `broker/wire_test.go` 与 `conformance/broker_cases.json`。改任一处 = breaking（已上 wire 的字节改变）。
> **不迁**（按 CT-05 ④）：XTP 适配器内部类型（`BusinessType` / `XTPConfig` / `XTPDeps`）、
> `OrderRequest.Fingerprint`（幂等指纹是 broker 的仓储规则，不是 wire 形状）、12 个契约哨兵
> （进程内错误值，HTTP 消费方只看 wire code）。包内**零 I/O、零时钟、零 decode**。
> **三处与 `domain` 的形状差异**见 `api-schema.md` §3.2，**禁止**用类型别名把两侧等同。

| 符号 | 类别 | 稳定性 | 错误语义 / 判据 |
|---|---|---|---|
| `Side` | type | stable | **两值** `buy` / `sell`；不是 `domain.Direction` 四值 |
| `SideBuy` | const | stable | `"buy"` |
| `SideSell` | const | stable | `"sell"` |
| `Side.Valid` | method | stable | 是否两值之一；零值 → false |
| `OrderType` | type | stable | **两值** `market` / `limit`；broker 刻意不收 `stop` / `trailing` |
| `OrderMarket` | const | stable | `"market"` |
| `OrderLimit` | const | stable | `"limit"` |
| `OrderType.Valid` | method | stable | 是否两值之一；零值 → false |
| `OrderStatus` | type | stable | **八值**（`pending` / `accepted` / `partially_filled` / `filled` / `cancelled` / `rejected` / `expired` / `unknown`）；不是 run 三值 |
| `OrderPending` | const | stable | `"pending"` |
| `OrderAccepted` | const | stable | `"accepted"`（**受理 ≠ 成交**） |
| `OrderPartiallyFilled` | const | stable | `"partially_filled"` |
| `OrderFilled` | const | stable | `"filled"` |
| `OrderCancelled` | const | stable | `"cancelled"` |
| `OrderRejected` | const | stable | `"rejected"` |
| `OrderExpired` | const | stable | `"expired"` |
| `OrderUnknown` | const | stable | `"unknown"`：**未决**状态，不是拒绝；调用方只能查单 / 对账，禁止重发 |
| `OrderStatus.Valid` | method | stable | 是否八值之一；零值 → false |
| `OrderStatus.Terminal` | method | stable | `filled` / `cancelled` / `rejected` / `expired` → true；`unknown` **不是**终态 |
| `OrderRequest` | type | stable | 7 字段；`quantity` 是**整股 `int`**；`limit_price` **`omitempty`**（市价单不出现该键） |
| `OrderAck` | type | stable | 7 字段；`broker_order_id` / `code` / `reason` **`omitempty`**；`accepted_at` 恒出现（非指针） |
| `CancelRequest` | type | stable | 2 字段，都 **`omitempty`**：至少一个非空，优先按 `client_order_id` 幂等 |
| `OrderState` | type | stable | 9 字段；`filled_quantity` / `avg_fill_price` / `remaining_quantity` **无 omitempty**（零也出现） |
| `Fill` | type | stable | 11 字段；`broker_order_id` **`omitempty`**；`stamp_tax` 零值也出现 |
| `Position` | type | stable | 7 字段；`quantity` / `available_quantity` 是**整股 `int`**（不是 `domain.Position.Quantity` 的 `float64`） |
| `Cash` | type | stable | `total` / `available`（float64） |
| `Fees` | type | stable | `commission` / `stamp_tax` / `transfer_fee`（float64，零值也出现） |
| `AccountSnapshot` | type | stable | 6 字段；`positions` 空时出 `[]` **不是 `null`**（生产侧恒非 nil） |
| `FlattenedPosition` | type | stable | 7 字段；只放**已真实提交并确认**的平仓成交 |
| `SkippedPosition` | type | stable | 4 字段；未提交 / 被拒 / 状态未知的一律进这里，不得计入 `sold` |
| `EmergencyFlattenResult` | type | stable | 6 字段；`sold` / `skipped` 空时出 `[]`；`sold_total` 是整股 |

<!-- surface:end -->

---

## 5. 冻结面逐符号规格（S1 已迁入 §4；本节是权威语义口径）

> 每个冻结面只**指回** `tellurion:docs/SPEC.md` 的节标题（SPEC.md 用命名节，不用数字编号），不复制正文——
> 复制的那一份一定会漂移。改这些节本身走 `tellurion:docs/TASKS.md`「维护本文件」的**冻结 schema 变更评审（硬规则）**。

### 5.1 四契约（ADR-031 D1 —— `tellurion-docs:docs/adr/adr-031-kernel-pluggable-nt-ready.md`）

| 契约术语 | 承载符号（✅ 已在本仓 `engine/`） | 稳定性 | 冻结文本指回 |
|---|---|---|---|
| `RunSpec` | `BacktestRequest` + `BacktestRequest.ValidateRunSpec` / `EffectiveKernel` / `EffectiveFillTiming`；伴生：`SignalWeight`、`TradingConfig` + `TradingConfig.WithDefaults`、`PriceLimitConfig`、`FillTiming` | **frozen** | `tellurion:docs/SPEC.md` §「RunSpec 内核路由（K10，2026-10-12）：EngineRunner 可插拔」 |
| `SignalSpec` | `SignalWeight`（K15 目标权重流）；信号本体 `domain.Signal` 随 `domain/` 迁入；供给契约 `SignalProvider`（今天 `tellurion:pkg/strategy`，见 §5.2） | stable（**SPEC.md 无 SignalSpec 节**） | `tellurion-contracts:docs/api-schema.md` §1.2 + `tellurion:docs/design/qt-decoupling/kernel-interface-map.md` §一.4 |
| `ResultSpec` | `BacktestResponse` + `StatusCompleted` / `StatusInvalid` / `StatusFailed`；`InvalidReason` + `CheckValidity` + `InvalidReasonCodes` | **frozen** | `tellurion:docs/SPEC.md` §「ResultSpec 冻结（K9，2026-10-12）：BacktestResponse canonical schema」→ 子节「status 三值语义」「字段集（27 项，冻结）」「指标口径」 |
| `EngineRunner` | `EngineRunner` interface（+ `EngineRegistry` / `NewEngineRegistry`、`MarketRouter` / `ResolveRunTarget`） | **frozen**（签名与 fail-closed 语义） | 同 RunSpec 行：§「RunSpec 内核路由（K10）」的「接口边界」与「fail-closed（零静默回退）」两条 |

> ⚠️ **`RunSpec` / `ResultSpec` 这两个 Go 类型名今天不存在**（承载物是 `BacktestRequest` / `BacktestResponse`）。
> W0（QT-01~04）已结束、S1 按「只搬不改名」落地（2026-10-22）：改名未发生；今后若改名属破坏性变更（MAJOR）。
> 见 `docs/api-schema.md` §1.2 的提示行。看到本表左边是术语、右边是符号名，不是笔误。

### 5.2 port 接口（方法集即形状）

| 接口 | 方法（导出面计入 `Iface.Method`） | 实现方 | 稳定性 | 错误语义 / 判据 |
|---|---|---|---|---|
| `EngineRunner` | `RunBacktest(ctx, BacktestRequest) (*BacktestResponse, error)` | core | **frozen** | 返回值里的 `Status` 是唯一成功判据；error 只装真实故障（§3.5） |
| `EngineRegistry` | `Register` / `Resolve` / `Names` | core（装配期） | stable | 未注册内核 → `ErrKernelNotRegistered`，**绝不静默回退 go**；`Resolve` 可并发，`Register` 仅装配期 |
| `MarketRouter` | `ResolveKernel` / `Markets` | core | stable | 未登记市场 → `ErrMarketNotRegistered`；路由表未装配 → 创建即拒绝、job 不落库 |
| `BacktestRunner` | `RunBacktest(ctx, strategyName, stockPool, startDate, endDate) (*domain.BacktestResult, error)` | core / agent 适配器 | stable | 与 `EngineRunner` **刻意分离**（五参 + 返回 domain 结果，AI 管线端口）；签名不动是 K10 明文约束 |
| `SignalProvider` | `Signals(ctx, Scope, asOf) ([]Signal, error)` | core（`tellurion:pkg/strategy`） | stable（**不住本仓**：QT-05 只搬 `domain` / `engine` / `fees`，接口留 core 至策略资产拆分） | 返回 `as_of > asOf` 的信号即违约：`RejectFutureSignals` fail-loud，不静默跳过 |
| `Provider`（`domain/market`） | `GetOHLCV` / `GetFundamental` / `GetStocks` / `GetLatestPrice` / `GetIndexConstituents` | **无**（实测 0 实现 / 0 消费：core 除定义外 0 引用，`pkg/marketdata` 用的是它自己的同名接口） | **experimental**（原标 stable，2026-10-22 落码时按实测下调——没人实现的 port 没有稳定性可言；§4 已同步） | **今天返回未分类 error**；QT-10 ③-② 落地后缺数据须包成 `ErrDataGap`（跳过当天 + 写 `invalid_reasons`），不是 fault —— 见 §5.4 的「⬜ 目标」行 |
| `RiskManager`（`domain/`） | `CheckStopLoss` / `CalculatePosition` / `DetectRegime` | core | stable | 裁决失败必须可归因；禁止把拒绝伪装成空信号 |
| `BrokerGateway` | **未定形**（形态见 `tellurion:docs/design/qt-decoupling/broker-gateway-design.md`） | `tellurion-broker` | ⬜ 设计中，S1 不落 | 见 `api-schema.md` §1.1「状态」列——今天源码 0 命中，**别当已存在接口去 import** |
| `DataPort` / `ExecutionPort` | **未定形** | `tellurion-data` / core | ⬜ 设计中，S1 不落 | 同上 |

### 5.3 状态与无效 Code（frozen 的取值集）

| 符号 | 取值 | 不变量（判据） | 指回 |
|---|---|---|---|
| `StatusCompleted` | `completed` | `invalid_reasons` 空、`error` 空 | SPEC §「status 三值语义」 |
| `StatusInvalid` | `invalid` | `invalid_reasons` **必填**、`error` 必空；指标照常给出，但**不得**当成功观测 | 同上 |
| `StatusFailed` | `failed` | `error` **必填**、`invalid_reasons` 空 | 同上 |
| `InvalidEmptyUniverse` 等 `Invalid*` 常量 | `empty_universe` / `zero_trades` / `zero_start_date` / `garbage_metric` / `exec_algo_unknown_style` / `exec_algo_attach_failed` | 稳定 Code；`InvalidReason.Detail` 才可改。新增 Code = 契约变更 | SPEC §「status 三值语义」的 Code 列表 |
| `CodeDataGapSkipped` | `data_gap_skipped` | ✅ **已落（2026-10-23，QT-10 ③-①）**。生产者：core 在 `ErrDataGap` 跳天分支上挂进 `invalid_reasons`；映射是 `status=invalid`，**不是第四值** | `tellurion:docs/TASKS.md` QT-10 ③-④ |
| （QT-15 的 Code，**未迁本仓**）| `session_terminated` | ⬜ 本仓源码里没有。今天它定义在 **core** `pkg/sandbox/wasm/replay.go:42`（`CodeSessionTerminated`，QT-15 落沙箱层时自带），映射同为 `status=invalid`。要不要收进本仓与 `CodeDataGapSkipped` 同处，属 core 侧收编时的裁决，**不在 QT-10 ③-① 范围内** | `tellurion:docs/TASKS.md` QT-15 ⑤、`tellurion-agent:docs/TASKS.md` REV-01-A |

> **禁止第四值**已两次明文（QT-10 ③-④、QT-63 ④）。`Failures[]` 只作诊断明细，**不是状态来源**——status 的唯一透传点是 `Engine.ResponseFromState`。

### 5.4 哨兵错误族

| 族 | 符号 | 稳定性 | 消费方行为 |
|---|---|---|---|
| 三分法（**QT-10 ③-①，✅ 2026-10-23 已落在 `engine/execution_errors.go`**） | `ErrDataGap` | stable | 跳过当天 + 写 `invalid_reasons`（`CodeDataGapSkipped`），run 继续、终态 `invalid`。三类里唯一不终止 run 的 |
| 三分法 | `ErrStrategyFault` | stable | **终止 run** → `status=failed` + `error`；不得当可重试的成功 |
| 三分法 | `ErrContractBreach` | stable | **终止 run** → `status=failed`；同时是「无法归类」的兜底类（§3.4） |
| ↑ 上面两行的**归类**（`ErrDataGap` 与两类的分界、provider 侧往哪一类包）在 core 执行（QT-10 ③-②/③-③）。本仓只保证：三类互斥、`errors.Is` 可判、`%w` 保链——钉子是 `execution_errors_test.go` 的六条 | — | — | 归类今天还没人做：core 的 provider 仍可能返回未分类 error（§3 的注） |
| 内核路由 | `ErrKernelNotRegistered`、`ErrMarketNotRegistered` | stable（✅ 已在 `engine/`） | 请求期 400；异步 job 终态 `failed`。**零静默回退** |
| 成交口径 | `ErrFillTimingMissing`、`ErrFillTimingUnknown`、`ErrFillTimingMismatch` | stable | 引擎永远不猜口径；口径不符是错误不是结果差异 |
| abi / artifact 校验 | `ErrAbiVersionUnknown`、`ErrAbiVersionMismatch`、`ErrAbiFamilyMissing`、`ErrAbiFamilyMixed`、`ErrImportNotWhitelisted`、`ErrImportSignatureMismatch`、`ErrHostRegionMissing`、`ErrHostRegionInvalid`、`ErrParamsTooLarge`、`ErrScopeUnknown`、`ErrScopeMarketV1` | stable（✅ 已在 `engine/`；`errors/` 不独立成包见 §8 #1） | 归 `ErrContractBreach` 语义：装载期拦下，不进 run |
| 信号校验 | `ErrActionIllegal`、`ErrStrengthOutOfRange`、`ErrTargetWeightOutOfRange`、`ErrEmitSourcesMixed`、`ErrEmitSourceUnknown`、`ErrSignalBudgetExceeded`、`ErrStateBudgetExceeded` | stable | 归 `ErrStrategyFault`（非法信号 / 预算超限 = 策略故障，终止 run） |
| provenance / hash | `ErrCanonicalUnsupported`、`ErrStrategyRefMalformed` | stable | 输入身份不可 canonical 编码 → 当场失败，禁止退化成「猜一个 hash」 |
| 策略编译 / 资产（**QT-41b ④，✅ 2026-10-24 已落在 `engine/execution_errors.go`**） | `ErrCompileFailed`、`ErrCompileUnavailable`、`ErrSourceTooLarge`、`ErrArtifactTooLarge` | stable | **跨仓消费方必须能识别**（agent 生成闭环 / core 取 artifact / orchestration 终态）：编译失败 → 反馈给模型；编译端口未配置 → 配置错 / 可重试；超限 → 不可重试。四条身份互斥，`errors.Is` 判定、`%w` 保链。**本仓（资产服务）私有哨兵不迁**（`ErrStrategyNotFound` / `ErrInvalidTransition` / `ErrJobNotFound` 等只看 wire code） |

### 5.5 数据协议面（`dataprotocol/`，QT-40 X1a——✅ 2026-10-24 新建，**stable 不是 frozen**）

| 判据 | 口径 | 证据 |
|---|---|---|
| 形状来源 | 三份已在漂的副本收成一份：data 镜像（带 `VisibleAt`、注入时钟）、core `source`（有 `DeduplicateKey`、构造器读墙钟）、core `storage`（因 import 成环另抄）。字段 = 三者并集，不多不少 | `tellurion-contracts:docs/design/qt40-x1a-dataprotocol-decision-material.md` §1 |
| `VisibleAt` 必进 | 它是 data 写入路径依赖的字段；不进类型，data 就得继续留镜像，迁移等于没做 | 同上 §3 判据 2 |
| 零时钟 / 零 I/O | 不提供构造器：`IngestTime` 由调用方注入，保证「同输入同字节」 | ADR-0004；同上 §3 判据 3 |
| 不重复 `TimeSemantics` | 时间语义词表已住 `domain/market`，本包只承载点位本身 | `docs/API.md` §8 #4 |
| 不落 `Aggregator` | 今天无实现、无消费方；接口随 X1b 定形，不许在契约仓造空 port | 同上 |
| 去重坐标 | `Symbol\|DataType\|TradeTime(UTC, RFC3339Nano)`，`Source` 不入键；向量由 `dataprotocol_test.go` 钉死（与 core `pkg/data/source/unified.go` 同式） | `dataprotocol/dataprotocol_test.go` |

> **迁移边界（如实）**：本仓只交付类型与登记。core 删两份副本、data 镜像改 import 是**各自仓的单**；
> 在那之前三份形状仍并存，且 `storage` 的列布局与共享字段的**对表测试**尚未补（`bulk_insert_ddl_test.go` 只核表名）——
> 这是 core 侧迁移单的验收项，不是本仓可以代做的。

### 5.6 风控快照面（`risk/`，CT-03——✅ 2026-11-07 新建，**hash 算法 frozen**）

| 判据 | 口径 | 证据 |
|---|---|---|
| 编码必须逐字节搬运 | `Snapshot.Hash()` 与 risk 的 `SnapshotHash` **产出同一字节流**：域标签裸写一行 + `\n`、长度前缀帧 `len(key):key=len(value):value\n`、按 `rule_id` 升序且重复即拒、参数 `UseNumber` 解码后 `json.Marshal` | `tellurion-risk:internal/rule/{snapshot.go:57-104,hash.go:52-55,model.go:150-174}` |
| **禁止**换 `engine.NewCanonical` | 两套编码器三处不同：帧结构（文本长度前缀 vs 二进制 uvarint + 类型标签）、域标签（带换行 vs 无换行首帧）、参数规范化（`UseNumber` 全收 vs 显式拒 `json.Number`）。换 = 全历史 `risk_rules_hash` 作废，且 `quant.run_provenance` 与 `run.run_artifacts` 两处锚点同时失效 | `engine/canonical.go:70-74`、`:83-91`、`:185-197`；`tellurion-contracts:docs/design/shape-packages-decision-material.md` §1.2 |
| 固定向量 | `14ecab865bebaef54052354b63caa5d351f533ab112e1a29a9dcad3707955c5b`（两条规则：`zz_guard` 乱序在前 + `aa_stop`；首条 params 带多余空白）。该 hex 由 `tellurion-risk` 内真实调用 `rule.SnapshotHash` 取得（2026-11-07，退出码 0），**不是按本仓实现反推** | `risk/snapshot_test.go` 的 `TestGoldenVector_SnapshotHash` |
| 跨仓对拍腿 | 测试里逐行重建 risk 侧编码器（独立实现，不调本包编码辅助函数），断言两侧同 hex；排序或帧格式任一漂移先红 | 同上 `TestSnapshotHash_MatchesRiskRepoEncoding` |
| 两个撞名坑 | `MaxParamsBytes`（本包 65536 vs `engine/` 4096）与 `ScopeSymbol`（`Scope` 取值 vs `StrategyScope` 取值）**同名不同义**，禁类型别名互转 | §4 `risk/` 表；材料 §1.3 |
| 不迁清单 | `Bundle` / `ParseBundle` / `ForScope` + 7 个服务私有哨兵 | §4 `risk/` 表注；CT-03 ③ |

> **迁移边界（如实）**：本单只交付契约形状、向量与登记。risk 侧改别名与 `tellurion-risk:docs/ARCHITECTURE.md:85`
> 的判据订正**归 risk 仓**；core 的 `bundle → RiskPolicy` 接线归 `tellurion:SVC-CORE-RISK-RETIRE`。
> ⇒ **本节落地的真实消费方今天仍为零**（risk 的 F3 未接、core 全仓 `ParseBundle` 零调用点）。
> 这是「先冻结形状、后续线」的有意取舍（与 `dataprotocol` 落库时同型）；代价是这段时间契约里多一个没人 `import` 的包。

### 5.7 编排 wire 面（`run/`，CT-04——✅ 2026-11-07 新建，**stable 不是 frozen**）

三份形状的来源与判据（证据为 2026-11-07 实测）：

| 判据 | 口径 | 证据 |
|---|---|---|
| claim 形状 | 7 字段，`parent_job_id` 带 `omitempty`；`run_spec` 是完整 RunSpec 原样字节 | `tellurion-orchestration:internal/httpapi/workers.go`；core 镜像 `tellurion:internal/orchestrationworker/client.go` |
| provenance 形状 | 15 字段；6 个可空字段**显式 `null`**，键不得省略；字段顺序即字节顺序 | `tellurion-orchestration:internal/httpapi/provenance.go`；core `tellurion:docs/openapi.yaml` 的 `BacktestProvenance`（原文 `must not drift`） |
| 子 run 入队形状 | `parentSubmitRequest` 族四个类型；`kind` 决定 `batch` / `walkforward` 二选一 | `tellurion-orchestration:internal/httpapi/parents.go`；analytics 镜像 `tellurion-analytics:internal/upstream/sensitivity.go` |
| 子 run 幂等键 | `parent:<parent_job_id>:<child_index>`；父 id 是计划身份的行 id，重放派生出同一批键 | `tellurion-orchestration:internal/parents/service.go` 的 `ChildIdempotencyKey` |
| 对拍腿 | 三份形状的 wire 字节由生产侧形状**独立重建**（字面量），与本包实现逐字节比对；改 `omitempty` 或分隔符立即红 | `run/run_test.go` 的三个 `Test*_MatchesOrchestration*` |
| 跨仓语料 | `conformance/run_cases.json`（与 `cases.json` 同形，但**单独文件**：core / agent 的加载器只认三侧，混入新侧会卡住下游 pin 升级）；期望值由本仓实现产出后写死，生产仓断言相同 | `conformance/run_conformance_test.go` |
| 不迁清单 | `RunRecord` / `ParentRecord` / `WorkerRecord` / `jobView` / `parentView` / `parentChildView` / `parentCountsView` / `ArtifactStore` / `run_artifacts` DDL | §4 `run/` 表注；材料 §2.2 |

> **迁移边界（如实）**：本单只交付契约形状、向量与语料。orchestration 改别名、core `orchestrationworker`
> 改用契约类型、analytics 提交报文改用契约类型，**都是各自仓的单**（CT-04 ⑥）；本单落地时三处仍各自持有副本。
> `docs/ARCHITECTURE.md` 与 README 兼容矩阵同步登记本包。

### 5.8 交易通道 wire 面（`broker/`，CT-05——✅ 2026-11-07 新建，**stable 不是 frozen**）

形状来源与判据（证据为 2026-11-07 实测）：

| 判据 | 口径 | 证据 |
|---|---|---|
| 订单三件套 | `OrderRequest`（7）/ `OrderAck`（7）/ `OrderState`（9）逐字段对齐；`limit_price` / `broker_order_id` / `code` / `reason` 带 `omitempty`，其余零值也出现 | `tellurion-broker:internal/gateway/model.go`；`tellurion-broker:docs/api-schema.md` §1.5 逐字段表 |
| 账户族 | `Fill`（11）/ `Position`（7）/ `Cash`（2）/ `Fees`（3）/ `AccountSnapshot`（6）；空持仓出 `[]` 不是 `null` | 同上；`internal/httpapi` 的键集合用例 |
| 紧急清仓 | `EmergencyFlattenResult`（6）+ `FlattenedPosition`（7）/ `SkippedPosition`（4）；未确认成交一律进 `skipped`，不得计入 `sold` | 同上；`SVC-BROKER-G1` 端点清单 |
| 与 `domain` 的三处差异 | `quantity` 整股 `int`（vs `float64`）、`Side` 两值（vs `Direction` 四值）、`OrderType` 两值（vs 四值）；**禁止**类型别名等同 | `docs/api-schema.md` §3.2；`domain/types.go` 对照 |
| 对拍腿 | 期望字节由 `tellurion-broker` 生产类型 `json.Marshal` 打出后逐字节抄入；改字段名 / tag / 顺序立即红 | `broker/wire_test.go` |
| 跨仓语料 | `conformance/broker_cases.json`（10 例，单独文件，理由同 §5.7） | `conformance/broker_conformance_test.go` |
| 不迁清单 | XTP 适配器内部类型、`OrderRequest.Fingerprint`、12 个契约哨兵、`broker.order_acks` DDL | §4 `broker/` 表注；材料 §3.3 / CT-05 ④ |

> **迁移边界（如实，2026-11-07）**：本单交付契约形状、向量与语料；`tellurion-broker` 删本地形状改引本包
> （CT-05 ⑤③ 的后半条）**未执行**——`v0.9.0` tag 未推远端（本机纪律只 commit 不 push），
> `go get …@v0.9.0` 解析不到，broker 切 `go.mod` 会破坏单仓可构建。解除前置后由 broker 一笔做，
> 判据是「同一条报文快照改前 / 改后逐字节一致」。

---

## 6. 值对象面（`domain/`，S1 整包迁入——✅ 2026-10-22 已落）

来源与再生成：**本仓内**跑 `go run ./tools/check_api_surface -list`（逐符号清单以该命令当场输出为准，本节不抄数字——一抄就会漂）。`engine` 的符号见 §4/§5。
**逐符号清单不在这里手抄**——抄一份就会漂；落码时用上面这条命令生成 §4 的表，再按下面分组填稳定性。

| 组 | 代表符号 | 稳定性 | 备注 |
|---|---|---|---|
| 回测输入输出 | `BacktestResult`、`PortfolioValue`、`Trade`、`WalkForward*` | stable（`BacktestResult` 的**指标口径**受 SPEC 冻结节约束） | `BacktestResult` **没有也不该有 `Status`**——三值属响应层，见 `tellurion-agent:docs/TASKS.md` OBS-01-R2-A 的 core 裁决 (b)。`BacktestParams` / `RiskMetrics` 只有 core 使用，已按 CT-06-B 迁回 `tellurion:pkg/domain` |
| 仪器域行情（已移出本包） | ~~`OHLCV`、`Stock`、`Fundamental`、`IndexDaily`、`Dividend`、`Split`~~（原 `market/` 的类型别名） | — | **CT-06-ALIAS-EXIT（2026-11-11；随下一个 MINOR `v0.10.0`，尚未发版）已删除这 6 个根别名**：实体一直在 `domain/market`，请直接 `import ".../domain/market"`。`FundamentalData` / `IndexConstituent` 根别名早前已按 CT-06-B 删除 |
| 时间/频率语义 | `Frequency`、`FrequencySpec`、`TimeSemantics`、`AsiaShanghai`、`Visible*`、`DailyVisibleAt` | **frozen** | PIT 红线：`tellurion:docs/SPEC.md` §「时间语义（PIT 红线）」。⚠️ 说清来源——core「冻结 schema 变更评审」点名的名单是三契约（RunSpec/ResultSpec/SignalSpec），这一行的 frozen 是**本文件按 §2 定义给的标注**，依据是它已有注册期闸门（D2）与逐源钉子（`tellurion:pkg/data/source/registry_rev05_test.go` 等把每个 connector 的 `TimeSemantics` 钉死）；改语义同样要先让那些测试红 |
| 信号与订单 | `Direction`（4 个常量）、`MarketRegime` | stable | `Signal` / `Order` / `OrderType` 族 / `Position` / `Portfolio` / `RiskMetrics` / `TargetPosition` / `PositionSize` / `StopLossEvent` / `RiskManager` 只有 core 使用，已按 CT-06-B 迁回 `tellurion:pkg/domain` |
| 因子族 | `FactorType` + 名册常量、`ParseFactorType`、`FactorSource*`、`FactorHypothesis`、`BuiltinFactorHypotheses` | stable（名册词表扩展 = 加名 + 定义 + 同步 SPEC） | 其余因子形状（`FactorKernel` / `FactorCaliber` / `FactorDefinition` / `FactorRegistration` / `FactorMiningRecord` / `FactorOutput` / `FactorCacheEntry` / `FactorReturn` / `ICEntry` / `BuiltinFactorDefinitions`）只有 core 使用，已按 CT-06-B 迁回 `tellurion:pkg/domain`；`BuiltinFactorHypotheses` 因 core + agent 两个使用方留本仓 |
| 配置形状 | ~~`Config`、`ServiceConfig`、`DatabaseConfig`、`RedisConfig`、`TushareConfig`~~（**已按 §8 #2 裁决删除：零引用 + 含 `Password`**）、~~`StrategyConfig`~~（CT-06-B 删除：0 使用方）、`ExecutionConfig`（CT-06-B 迁回 core） | — | 本节已无 `*Config` 登记行；`engine/` 的 `TradingConfig` / `PriceLimitConfig` 另按「冻结面伴生名」保留（[`adr/0002`](adr/0002-config-族的边界处置.md) ③） |
| 单使用方残留（已清零） | ~~5 个 market 别名（`OHLCV` / `Stock` / `IndexDaily` / `Split` / `Dividend`）+ `Fundamental`~~ | — | **CT-06-ALIAS-EXIT（2026-11-11；随下一个 MINOR `v0.10.0`，尚未发版）已删除全部 6 个根别名**，改用 `domain/market.*`；删除前 11 仓重跑实测真实残留 0。逐名处置见 §6.1 |

> **公司行为族已删除（CT-06，2026-11-07）**：`CorporateAction` / `ActionType` / `CashDividend` / `BonusShare` / `CorporateActionSplit` / `RightsIssue` / `Placement` / `ActionEngine` / `ApplyResult` / `ApplyOutcome` 及构造函数共 68 个符号，连同旧名别名 `MarketDataProvider`（1 个符号），11 仓重跑均 0 使用方，已从 §4 与代码中删除。破坏性记录见 `CHANGELOG.md`。

### 6.1 CT-06 重跑处置清单（2026-11-07，119 个根名）

> **取证命令**（PowerShell，11 个同级仓 `C:\Users\ruoxi\Workspace\tellurion*`，排除本仓；`*.go` 含测试；只统计导入 `github.com/ruoxizhnya/tellurion-contracts/domain` 的文件里的选择器）：
> `Get-ChildItem <repo> -Recurse -Filter *.go | ForEach-Object { Select-String -Path $_ -Pattern '\bdomain\.<Name>\b' }`
> 下表每条的证据列 = 代表性 `repo:path:line` + 该仓命中数（`nontest` / `test` 单列）。「0」= 11 仓全部 0 命中。基线 `go run ./tools/check_api_surface -list` 当日 `domain/` 根包 176 符号 / 119 根名。
>
> **复核口径**：先 grep 初筛，再用 Go AST（`go/parser`）只数 `SelectorExpr`，排除字符串 / 注释里的 `domain.X`。AST 复核后只有 5 个根名从「2 仓」降为「1 仓」：`OHLCV` / `OrderType` / `OrderTypeLimit` / `OrderTypeMarket` / `Portfolio`——agent 侧命中全部位于 `tellurion-agent:pkg/ai/prompts.go` 的提示词代码示例字符串，不是编译使用方；处置已改判为迁回。
>
> **处置统计**：删除 22 个根名（69 个符号，含 `MarketDataProvider`）；保留 36 个根名（34 个 ≥2 使用方 + `CaliberRank` / `CaliberRaw` 两个封闭词表成员）；迁回建议 61 个根名（均只有 1 个使用方：`tellurion` 60 + `tellurion-agent` 1，本轮因跨仓边界不执行）。删除后 `domain/` 根包 107 符号 / 97 根名。
>
> **CT-06-A 已执行（2026-11-09）**：`tellurion-agent` 侧唯一使用方 `BuiltinHypothesisFor` 1 个符号已迁回该仓（`pkg/ai/tools/builtin` 内私有 `builtinHypothesisFor`），本仓同步删除；`domain/` 根包变为 **106 符号 / 96 根名**，迁回待办剩 **60 个根名（全部属 core）**。内置表 `BuiltinFactorHypotheses` 因 core + agent 两个使用方，留在本仓不迁。
>
> **CT-06-B 已执行（2026-11-10，补裁 #12）**：60 个 core-only 根名处理完毕。下表「处置」列保留 2026-11-07 判定原貌（历史证据不追改），实际落地以本块为准，分三类：
> 1. **迁回 core**：本仓已删同名导出，实类型迁入 `tellurion:pkg/domain`（core 新建包承载；market 别名在 core 侧同形重建，不重导出共享符号）。本仓对应删掉 `domain/` 13 个文件，并收缩 `types.go` / `backtest.go` / `factor.go`。
> 2. **删除且不在 core 重建（11 仓实测 0 编译使用方）**：`StrategyConfig`；`Sector` / `StockSector` / `StockIndustry`（`domain/sector.go` 整文件删除）；`IndexConstituent` / `FundamentalData` 两个**根别名**（实体仍在 `domain/market/`，需要时直接 import 该包）。
> 3. **留本仓**：6 个 market 别名（`Fundamental` 仍被 `tellurion-data` 直接使用；其余是同一类型别名）；`FactorMarketActivity`（`FactorType` 名册成员，随 `ParseFactorType` 留）；`MacroSeries` / `MarketActivity`（core + `tellurion-data`）；`BuiltinFactorHypotheses`（core + agent）。
>
> 破坏性记录见 `CHANGELOG.md`；core 侧调用点已同批切至 `tellurion/pkg/domain`（见 core 仓同笔 commit 与 `internal/repoguard/contract_identity_test.go` 的腿 A/C 口径）。
>
> **CT-06-ALIAS-EXIT 已执行（2026-11-11；随下一个 MINOR `v0.10.0`，尚未发版）**：上表 CT-06-B 块第 3 类「留本仓」的 **6 个 market 别名自本块起作废**——`OHLCV` / `Stock` / `IndexDaily` / `Split` / `Dividend` / `Fundamental` 已从 `domain/` 根包删除（`market_alias_test.go` 整文件删除；`types.go` 只剩 `Direction` 族与 `MarketRegime`）。删除前置与实测：
> 1. **11 仓重跑（2026-11-11）**：六名经**根包别名**的真实编译残留 = **0**（core 188 / agent 48 / data 3 / contracts 7 个文件 import 根 `domain`，但均未使用六名）；dot-import = 0；本地再别名 = 0。字符串字面量命中（`"OHLCV"` 等，用于错误分类与测试数据）不依赖类型别名，不构成阻塞。
> 2. **先修后删**：`tellurion-agent` 的 6 处非归档引用已于同日改指 `market.*` / core 的 `domain.Portfolio`（`tellurion-agent:a6e7d02`）；`tellurion-data` 更早已改 `market.Fundamental`。
> 3. **护栏**：本仓 `go build ./...`、`go vet ./...`、`go test ./...` 与 `check_api_surface` / `check_surface_freeze` / `check_contract_shape` / `check_release_matrix` 全部通过（实测输出记在 `docs/TASKS.md` 的收口记录）。
> 4. **下游 pin**：`tellurion` / `tellurion-agent` / `tellurion-data` 仍 pin `v0.9.0`；因已不使用六名，升 pin 前仍可编译。升 pin 是发版后的独立 PR（需远端 tag）。

| 根名 | 使用方仓数 | 处置 | grep 出处（repo:path:line / 命令 + 命中数） |
|---|---|---|---|
| `ActionBonusShare` | 0 | 删除 | 11 仓 grep `\bdomain\.ActionBonusShare\b` → 0 命中（含测试） |
| `ActionCashDividend` | 0 | 删除 | 11 仓 grep `\bdomain\.ActionCashDividend\b` → 0 命中（含测试） |
| `ActionEngine` | 0 | 删除 | 11 仓 grep `\bdomain\.ActionEngine\b` → 0 命中（含测试） |
| `ActionPlacement` | 0 | 删除 | 11 仓 grep `\bdomain\.ActionPlacement\b` → 0 命中（含测试） |
| `ActionRegister` | 1 | 迁回（待跨仓单） | tellurion:pkg/data/factor_governance.go:129 (6 hits, nontest=2, test=4) |
| `ActionRightsIssue` | 0 | 删除 | 11 仓 grep `\bdomain\.ActionRightsIssue\b` → 0 命中（含测试） |
| `ActionRollback` | 1 | 迁回（待跨仓单） | tellurion:pkg/data/factor_governance.go:162 (5 hits, nontest=2, test=3) |
| `ActionSplit` | 0 | 删除 | 11 仓 grep `\bdomain\.ActionSplit\b` → 0 命中（含测试） |
| `ActionType` | 0 | 删除 | 11 仓 grep `\bdomain\.ActionType\b` → 0 命中（含测试） |
| `ApplyOutcome` | 0 | 删除 | 11 仓 grep `\bdomain\.ApplyOutcome\b` → 0 命中（含测试） |
| `ApplyResult` | 0 | 删除 | 11 仓 grep `\bdomain\.ApplyResult\b` → 0 命中（含测试） |
| `BacktestParams` | 1 | 迁回（待跨仓单） | tellurion:cmd/analysis/handlers_backtest_readpath_test.go:56 (34 hits, nontest=13, test=21) |
| `BacktestResult` | 2 | 保留 | tellurion:cmd/analysis/handlers_backtest_readpath_test.go:57 (72 hits, nontest=23, test=49) ; tellurion-agent:cmd/ai/handlers_explore.go:381 (99 hits, nontest=38, test=61) |
| `BonusShare` | 0 | 删除 | 11 仓 grep `\bdomain\.BonusShare\b` → 0 命中（含测试） |
| `BuiltinFactorDefinitions` | 1 | 迁回（待跨仓单） | tellurion:cmd/analysis/handlers_factor.go:195 (27 hits, nontest=3, test=24) |
| `BuiltinFactorHypotheses` | 2 | 保留 | tellurion:pkg/storage/factor_hypothesis.go:55 (2 hits, nontest=2, test=0) ; tellurion-agent:pkg/ai/tools/builtin/factor_hypothesis_tool.go:191 (9 hits, nontest=6, test=3) |
| `BuiltinHypothesisFor` | 1 | ✅ 已迁回 agent（CT-06-A，2026-11-09） | 原 `tellurion-agent:pkg/ai/tools/builtin/factor_hypothesis_tool.go:156`；现为本仓外的 `builtinHypothesisFor`，本仓已删除 |
| `CaliberRank` | 0 | 保留（封闭词表） | 11 仓 grep → 0 命中；`domain/factor_output.go:19` 词表成员，`FactorCaliber.Valid()` 依赖 |
| `CaliberRaw` | 0 | 保留（封闭词表） | 11 仓 grep → 0 命中；`domain/factor_output.go:21` 词表成员，`FactorCaliber.Valid()` 依赖 |
| `CaliberZScore` | 1 | 迁回（待跨仓单） | tellurion:cmd/data/mining_cmd_test.go:105 (14 hits, nontest=1, test=13) |
| `CashDividend` | 0 | 删除 | 11 仓 grep `\bdomain\.CashDividend\b` → 0 命中（含测试） |
| `CorporateAction` | 0 | 删除 | 11 仓 grep `\bdomain\.CorporateAction\b` → 0 命中（含测试） |
| `CorporateActionSplit` | 0 | 删除 | 11 仓 grep `\bdomain\.CorporateActionSplit\b` → 0 命中（含测试） |
| `DefaultExecutionConfig` | 1 | 迁回（待跨仓单） | tellurion:pkg/backtest/execution/execution_test.go:13 (17 hits, nontest=2, test=15) |
| `Direction` | 2 | 保留 | tellurion:cmd/analysis/alert_loop_test.go:71 (62 hits, nontest=40, test=22) ; tellurion-agent:pkg/ai/yaml/config.go:59 (3 hits, nontest=3, test=0) |
| `DirectionClose` | 2 | 保留 | tellurion:cmd/analysis/handlers_export_test.go:51 (126 hits, nontest=32, test=94) ; tellurion-agent:pkg/ai/yaml/config.go:62 (1 hits, nontest=1, test=0) |
| `DirectionHold` | 2 | 保留 | tellurion:pkg/backtest/engine.go:1671 (13 hits, nontest=9, test=4) ; tellurion-agent:pkg/ai/yaml/config.go:62 (1 hits, nontest=1, test=0) |
| `DirectionLong` | 2 | 保留 | tellurion:cmd/analysis/handlers_execution.go:184 (414 hits, nontest=50, test=364) ; tellurion-agent:cmd/ai/handlers_explore_test.go:59 (5 hits, nontest=1, test=4) |
| `DirectionShort` | 2 | 保留 | tellurion:cmd/dualrun/diff.go:63 (75 hits, nontest=33, test=42) ; tellurion-agent:pkg/ai/yaml/config.go:62 (1 hits, nontest=1, test=0) |
| `Dividend` | 1 | 迁回（待跨仓单） | tellurion:cmd/data/sync_dates_rev24_test.go:42 (26 hits, nontest=18, test=8) |
| `ETFConstituent` | 1 | 迁回（待跨仓单） | tellurion:cmd/analysis/handlers_industry_test.go:280 (19 hits, nontest=12, test=7) |
| `ExecutionConfig` | 1 | 迁回（待跨仓单） | tellurion:pkg/backtest/aliases.go:155 (41 hits, nontest=24, test=17) |
| `FactorCacheEntry` | 1 | 迁回（待跨仓单） | tellurion:cmd/data/handlers_factor.go:36 (39 hits, nontest=21, test=18) |
| `FactorCaliber` | 1 | 迁回（待跨仓单） | tellurion:pkg/mining/mining.go:83 (1 hits, nontest=1, test=0) |
| `FactorContractLiabilityRatio` | 2 | 保留 | tellurion:pkg/data/factor_equitydeep.go:476 (7 hits, nontest=1, test=6) ; tellurion-agent:pkg/ai/tools/builtin/factor_hypothesis_tool_test.go:151 (1 hits, nontest=0, test=1) |
| `FactorDefinition` | 1 | 迁回（待跨仓单） | tellurion:cmd/analysis/handlers_factor_registration.go:68 (62 hits, nontest=30, test=32) |
| `FactorGrossMarginTrend` | 2 | 保留 | tellurion:cmd/data/handlers_factor_test.go:86 (13 hits, nontest=1, test=12) ; tellurion-agent:pkg/ai/tools/builtin/factor_hypothesis_tool_test.go:151 (1 hits, nontest=0, test=1) |
| `FactorGrowth` | 2 | 保留 | tellurion:pkg/data/factor_registry_test.go:27 (3 hits, nontest=0, test=3) ; tellurion-agent:pkg/ai/tools/builtin/factor_hypothesis_tool_test.go:150 (1 hits, nontest=0, test=1) |
| `FactorHypothesis` | 2 | 保留 | tellurion:pkg/storage/factor_hypothesis.go:23 (6 hits, nontest=6, test=0) ; tellurion-agent:pkg/ai/tools/builtin/factor_hypothesis_tool.go:34 (24 hits, nontest=13, test=11) |
| `FactorInventoryTurnoverDelta` | 2 | 保留 | tellurion:pkg/data/factor_equitydeep.go:568 (6 hits, nontest=1, test=5) ; tellurion-agent:pkg/ai/tools/builtin/factor_hypothesis_tool_test.go:153 (1 hits, nontest=0, test=1) |
| `FactorKernel` | 1 | 迁回（待跨仓单） | tellurion:pkg/data/factor_registry_test.go:45 (1 hits, nontest=0, test=1) |
| `FactorMarketActivity` | 1 | 迁回（待跨仓单） | tellurion:pkg/data/factor_market_activity_governance_test.go:22 (16 hits, nontest=0, test=16) |
| `FactorMiningRecord` | 1 | 迁回（待跨仓单） | tellurion:cmd/data/mining_cmd.go:63 (29 hits, nontest=10, test=19) |
| `FactorMomentum` | 2 | 保留 | tellurion:cmd/data/handlers_factor_test.go:97 (101 hits, nontest=3, test=98) ; tellurion-agent:pkg/ai/prompts.go:59 (2 hits, nontest=1, test=1) |
| `FactorOCFToNetProfit` | 2 | 保留 | tellurion:pkg/data/factor_equitydeep.go:504 (6 hits, nontest=1, test=5) ; tellurion-agent:pkg/ai/tools/builtin/factor_hypothesis_tool_test.go:152 (1 hits, nontest=0, test=1) |
| `FactorOutput` | 1 | 迁回（待跨仓单） | tellurion:e2e/guard/spec_drift_test.go:218 (2 hits, nontest=0, test=2) |
| `FactorQuality` | 2 | 保留 | tellurion:pkg/backtest/cache/factor_cache.go:229 (14 hits, nontest=3, test=11) ; tellurion-agent:pkg/ai/prompts.go:59 (2 hits, nontest=1, test=1) |
| `FactorRegistration` | 1 | 迁回（待跨仓单） | tellurion:cmd/analysis/handlers_factor_registration.go:12 (43 hits, nontest=22, test=21) |
| `FactorReturn` | 1 | 迁回（待跨仓单） | tellurion:pkg/data/factor_attribution.go:23 (11 hits, nontest=9, test=2) |
| `FactorROEDuPontLeverage` | 2 | 保留 | tellurion:pkg/data/factor_equitydeep.go:532 (6 hits, nontest=1, test=5) ; tellurion-agent:pkg/ai/tools/builtin/factor_hypothesis_tool_test.go:152 (1 hits, nontest=0, test=1) |
| `FactorSize` | 2 | 保留 | tellurion:pkg/data/factor_governance_test.go:149 (11 hits, nontest=0, test=11) ; tellurion-agent:pkg/ai/tools/builtin/factor_hypothesis_tool_test.go:150 (1 hits, nontest=0, test=1) |
| `FactorSourceAdHoc` | 2 | 保留 | tellurion:pkg/mining/mining.go:61 (2 hits, nontest=2, test=0) ; tellurion-agent:pkg/ai/tools/builtin/factor_hypothesis_tool.go:170 (4 hits, nontest=2, test=2) |
| `FactorSourceAI` | 2 | 保留 | tellurion:cmd/data/mining_cmd_test.go:110 (4 hits, nontest=1, test=3) ; tellurion-agent:pkg/ai/tools/builtin/factor_hypothesis_tool_test.go:82 (5 hits, nontest=0, test=5) |
| `FactorSourceLiterature` | 2 | 保留 | tellurion:pkg/mining/mining.go:60 (1 hits, nontest=1, test=0) ; tellurion-agent:pkg/ai/tools/builtin/factor_hypothesis_tool_test.go:54 (2 hits, nontest=0, test=2) |
| `FactorSourceSupplyChain` | 2 | 保留 | tellurion:pkg/mining/mining.go:61 (1 hits, nontest=1, test=0) ; tellurion-agent:pkg/ai/tools/builtin/factor_hypothesis_tool_test.go:68 (2 hits, nontest=0, test=2) |
| `FactorType` | 2 | 保留 | tellurion:cmd/analysis/handlers_factor_registration.go:78 (98 hits, nontest=59, test=39) ; tellurion-agent:pkg/ai/prompts.go:57 (14 hits, nontest=7, test=7) |
| `FactorValue` | 2 | 保留 | tellurion:pkg/backtest/engine_accessors_test.go:155 (23 hits, nontest=3, test=20) ; tellurion-agent:pkg/ai/prompts.go:59 (3 hits, nontest=1, test=2) |
| `FactorVolatility` | 2 | 保留 | tellurion:pkg/data/factor_registry_test.go:27 (3 hits, nontest=0, test=3) ; tellurion-agent:pkg/ai/tools/builtin/factor_hypothesis_tool_test.go:150 (1 hits, nontest=0, test=1) |
| `Fundamental` | 2 | 保留 | tellurion:pkg/backtest/engine_accessors_test.go:25 (82 hits, nontest=57, test=25) ; tellurion-data:internal/ingest/tushare_adapter.go:187 (3 hits, nontest=3, test=0) |
| `FundamentalData` | 1 | 迁回（待跨仓单） | tellurion:cmd/data/handlers_factor_test.go:40 (45 hits, nontest=21, test=24) |
| `FundamentalsDetailRow` | 1 | 迁回（待跨仓单） | tellurion:cmd/data/handlers_equitydeep_ingest.go:104 (37 hits, nontest=9, test=28) |
| `ICEntry` | 1 | 迁回（待跨仓单） | tellurion:pkg/data/factor_attribution.go:24 (12 hits, nontest=10, test=2) |
| `IndexConstituent` | 1 | 迁回（待跨仓单） | tellurion:cmd/data/sync_index_constituents_persist_test.go:43 (28 hits, nontest=13, test=15) |
| `IndexDaily` | 1 | 迁回（待跨仓单） | tellurion:e2e/guard/spec_drift_test.go:228 (17 hits, nontest=11, test=6) |
| `IndustryVitality` | 1 | 迁回（待跨仓单） | tellurion:cmd/analysis/handlers_industry_test.go:232 (17 hits, nontest=8, test=9) |
| `KernelCode` | 1 | 迁回（待跨仓单） | tellurion:pkg/data/factor.go:369 (10 hits, nontest=4, test=6) |
| `KernelFormula` | 1 | 迁回（待跨仓单） | tellurion:cmd/analysis/handlers_factor_registration_test.go:49 (16 hits, nontest=5, test=11) |
| `KernelLeaf` | 1 | 迁回（待跨仓单） | tellurion:cmd/analysis/handlers_factor_registration_test.go:49 (7 hits, nontest=4, test=3) |
| `KernelPending` | 1 | 迁回（待跨仓单） | tellurion:pkg/data/factor.go:389 (9 hits, nontest=4, test=5) |
| `MacroSeries` | 2 | 保留 | tellurion:cmd/analysis/handlers_macro_test.go:122 (33 hits, nontest=14, test=19) ; tellurion-data:internal/ingest/transport.go:37 (1 hits, nontest=1, test=0) |
| `MarketActivity` | 2 | 保留 | tellurion:cmd/data/handlers_factor_test.go:49 (49 hits, nontest=16, test=33) ; tellurion-data:internal/ingest/transport.go:42 (1 hits, nontest=1, test=0) |
| `MarketDataProvider` | 0 | 删除 | 11 仓 grep `\bdomain\.MarketDataProvider\b` → 0 命中（含测试） |
| `MarketRegime` | 2 | 保留 | tellurion:cmd/analysis/handlers_risk.go:92 (73 hits, nontest=35, test=38) ; tellurion-agent:pkg/ai/client/risk_client.go:62 (20 hits, nontest=5, test=15) |
| `NewActionEngine` | 0 | 删除 | 11 仓 grep `\bdomain\.NewActionEngine\b` → 0 命中（含测试） |
| `NewBonusShare` | 0 | 删除 | 11 仓 grep `\bdomain\.NewBonusShare\b` → 0 命中（含测试） |
| `NewCashDividend` | 0 | 删除 | 11 仓 grep `\bdomain\.NewCashDividend\b` → 0 命中（含测试） |
| `NewPlacement` | 0 | 删除 | 11 仓 grep `\bdomain\.NewPlacement\b` → 0 命中（含测试） |
| `NewRightsIssue` | 0 | 删除 | 11 仓 grep `\bdomain\.NewRightsIssue\b` → 0 命中（含测试） |
| `NewSplit` | 0 | 删除 | 11 仓 grep `\bdomain\.NewSplit\b` → 0 命中（含测试） |
| `OHLCV` | 1 | 迁回（待跨仓单） | tellurion:cmd/analysis/handlers_risk.go:43 (795 hits, AST 复核 nontest=262, test=533)；agent 侧仅 `pkg/ai/prompts.go:43` 的提示词代码示例字符串，**非编译使用方** |
| `Order` | 1 | 迁回（待跨仓单） | tellurion:pkg/backtest/engine_daily.go:836 (147 hits, nontest=51, test=96) |
| `OrderType` | 1 | 迁回（待跨仓单） | tellurion:cmd/analysis/alert_loop_test.go:71 (23 hits, AST 复核 nontest=16, test=7)；agent 侧仅 `pkg/ai/prompts.go:34` 的提示词代码示例字符串，**非编译使用方** |
| `OrderTypeLimit` | 1 | 迁回（待跨仓单） | tellurion:cmd/analysis/wiring/live_trader_port_test.go:80 (42 hits, AST 复核 nontest=12, test=30)；agent 侧仅 `pkg/ai/prompts.go:34-35` 的提示词代码示例字符串，**非编译使用方** |
| `OrderTypeMarket` | 1 | 迁回（待跨仓单） | tellurion:cmd/analysis/handlers_execution.go:181 (96 hits, AST 复核 nontest=19, test=77)；agent 侧仅 `pkg/ai/prompts.go:34` 的提示词代码示例字符串，**非编译使用方** |
| `OrderTypeStop` | 1 | 迁回（待跨仓单） | tellurion:pkg/live/engine.go:338 (8 hits, nontest=3, test=5) |
| `OrderTypeTrailing` | 1 | 迁回（待跨仓单） | tellurion:pkg/live/engine.go:281 (13 hits, nontest=4, test=9) |
| `ParseFactorType` | 2 | 保留 | tellurion:cmd/analysis/handlers_factor.go:19 (7 hits, nontest=6, test=1) ; tellurion-agent:pkg/ai/tools/builtin/factor_hypothesis_tool.go:127 (2 hits, nontest=1, test=1) |
| `Placement` | 0 | 删除 | 11 仓 grep `\bdomain\.Placement\b` → 0 命中（含测试） |
| `Portfolio` | 1 | 迁回（待跨仓单） | tellurion:cmd/analysis/handlers_risk.go:91 (125 hits, AST 复核 nontest=42, test=83)；agent 侧仅 `pkg/ai/prompts.go:43` 的提示词代码示例字符串，**非编译使用方** |
| `PortfolioValue` | 2 | 保留 | tellurion:cmd/analysis/handlers_compare_test.go:79 (31 hits, nontest=21, test=10) ; tellurion-agent:cmd/ai/handlers_explore_test.go:51 (24 hits, nontest=5, test=19) |
| `Position` | 1 | 迁回（待跨仓单） | tellurion:cmd/analysis/handlers_risk.go:235 (117 hits, nontest=54, test=63) |
| `PositionSize` | 1 | 迁回（待跨仓单） | tellurion:pkg/backtest/engine.go:1859 (30 hits, nontest=24, test=6) |
| `RegistrationAction` | 1 | 迁回（待跨仓单） | tellurion:pkg/data/factor_governance_test.go:102 (1 hits, nontest=0, test=1) |
| `RegistrationActive` | 1 | 迁回（待跨仓单） | tellurion:cmd/analysis/handlers_factor_registration.go:108 (11 hits, nontest=3, test=8) |
| `RegistrationAuditEntry` | 1 | 迁回（待跨仓单） | tellurion:pkg/data/factor_governance.go:56 (11 hits, nontest=11, test=0) |
| `RegistrationRolledBack` | 1 | 迁回（待跨仓单） | tellurion:pkg/data/factor_governance.go:149 (3 hits, nontest=2, test=1) |
| `RegistrationStatus` | 1 | 迁回（待跨仓单） | tellurion:cmd/analysis/handlers_factor_registration.go:81 (1 hits, nontest=1, test=0) |
| `RelationIndustryMember` | 1 | 迁回（待跨仓单） | tellurion:cmd/analysis/handlers_industry_test.go:196 (10 hits, nontest=1, test=9) |
| `RightsIssue` | 0 | 删除 | 11 仓 grep `\bdomain\.RightsIssue\b` → 0 命中（含测试） |
| `RiskManager` | 1 | 迁回（待跨仓单） | tellurion:pkg/risk/interfaces.go:9 (15 hits, nontest=9, test=6) |
| `RiskMetrics` | 1 | 迁回（待跨仓单） | tellurion:pkg/backtest/metrics/performance.go:13 (4 hits, nontest=4, test=0) |
| `ScreenFilters` | 1 | 迁回（待跨仓单） | tellurion:pkg/storage/fundamentals.go:504 (19 hits, nontest=4, test=15) |
| `ScreenRequest` | 1 | 迁回（待跨仓单） | tellurion:cmd/data/handlers_stocks.go:316 (5 hits, nontest=4, test=1) |
| `ScreenResult` | 1 | 迁回（待跨仓单） | tellurion:pkg/storage/fundamentals.go:504 (34 hits, nontest=11, test=23) |
| `Sector` | 1 | 迁回（待跨仓单） | tellurion:pkg/data/sector_sync.go:19 (4 hits, nontest=4, test=0) |
| `Signal` | 1 | 迁回（待跨仓单） | tellurion:cmd/analysis/handlers_risk.go:90 (171 hits, nontest=74, test=97) |
| `Split` | 1 | 迁回（待跨仓单） | tellurion:cmd/data/sync_dates_rev24_test.go:56 (28 hits, nontest=18, test=10) |
| `Stock` | 1 | 迁回（待跨仓单） | tellurion:cmd/analysis/signals_e2e_test.go:104 (110 hits, nontest=60, test=50) |
| `StockIndustry` | 1 | 迁回（待跨仓单） | tellurion:pkg/data/sector_sync.go:18 (4 hits, nontest=4, test=0) |
| `StockSector` | 1 | 迁回（待跨仓单） | tellurion:pkg/data/sector_sync.go:20 (4 hits, nontest=4, test=0) |
| `StopLossEvent` | 1 | 迁回（待跨仓单） | tellurion:cmd/analysis/handlers_risk_execution_test.go:302 (16 hits, nontest=15, test=1) |
| `StrategyConfig` | 1 | 迁回（待跨仓单） | tellurion:cmd/analysis/handlers_strategy.go:51 (18 hits, nontest=12, test=6) |
| `SupplyChainEdge` | 1 | 迁回（待跨仓单） | tellurion:cmd/analysis/handlers_industry_test.go:195 (20 hits, nontest=13, test=7) |
| `SWIndustryNode` | 1 | 迁回（待跨仓单） | tellurion:cmd/analysis/handlers_industry_test.go:129 (24 hits, nontest=13, test=11) |
| `TargetPosition` | 1 | 迁回（待跨仓单） | tellurion:pkg/backtest/engine.go:842 (12 hits, nontest=7, test=5) |
| `Trade` | 2 | 保留 | tellurion:cmd/analysis/handlers_export_test.go:49 (120 hits, nontest=47, test=73) ; tellurion-agent:cmd/ai/handlers_explore_test.go:58 (14 hits, nontest=3, test=11) |
| `WalkForwardParams` | 2 | 保留 | tellurion:pkg/backtest/batch/batch.go:260 (18 hits, nontest=3, test=15) ; tellurion-agent:cmd/ai/http_runner.go:189 (11 hits, nontest=5, test=6) |
| `WalkForwardReport` | 2 | 保留 | tellurion:pkg/backtest/batch/batch.go:44 (35 hits, nontest=16, test=19) ; tellurion-agent:cmd/ai/http_runner.go:190 (19 hits, nontest=7, test=12) |
| `WalkForwardResult` | 2 | 保留 | tellurion:pkg/backtest/walkforward/walkforward.go:156 (31 hits, nontest=6, test=25) ; tellurion-agent:pkg/ai/tools/builtin/walkforward_tool_test.go:56 (4 hits, nontest=0, test=4) |

---

## 7. abi 常量族（`engine/` 内，✅ QT-05 已迁入）

来源：`tellurion:pkg/backtest/contracts/abi.go`（v2 修订见 `tellurion:docs/design/qt-decoupling/adr-029-amendment-host-api-v2.md`）。

| 组 | 符号 | 稳定性 | 判据 |
|---|---|---|---|
| host 函数名 | **v1 家族 8 个**：`HostFnGetBar` / `HostFnGetSeriesLen` / `HostFnGetSymbolCount` / `HostFnGetCrossSection` / `HostFnGetState` / `HostFnSetState` / `HostFnEmitSignal` / `HostFnLog`<br>**v2 家族 4 个**：`HostFnGetBar2` / `HostFnGetSeriesLen2` / `HostFnGetDate2` / `HostFnEmitTarget2`<br>合计 **12**（`HostFnCount()` = `len(hostFnContracts)`，返回枚举结果而非文档里的字面数） | **frozen** | 名字与签名成对冻结，`HostABIContract` / `HostABIContracts` 是唯一名单源（每条带 `Since: HostAbiV1/V2`）。artifact 导入表不符 ⇒ `ErrImportSignatureMismatch`；两族都不导 ⇒ `ErrAbiFamilyMissing`；混用 ⇒ `ErrAbiFamilyMixed`。<br>⚠️ 旧文档的「13 个」是**已被 core 订正的错数**（2026-10-21，QT-01 落代码时实测枚举为 12；见 `tellurion:docs/design/qt-decoupling/adr-029-amendment-host-api-v2.md` §2.1 的订正条）。真要加第 13 个：先改那张表、走一次 abi 冻结 |
| abi 版本 | `HostAbiV1` / `HostAbiV2` / `HostAbiVersion` / `HostAbiVersionFromName` / `HostFnCount` | **frozen** | `HostAbiVersion.Valid()` 只认 1/2；`HostAbiVersionFromName` 解析 `on_bar/N` 失败 ⇒ `ErrAbiVersionUnknown`；声明与导入段矛盾 ⇒ `ErrAbiVersionMismatch`；两族混用 ⇒ `ErrAbiFamilyMixed` |
| 内存区所有权 | `ExportHostParamsPtr` / `ExportHostScratchPtr` / `ExportHostScratchSize` / `MinHostParamsPtr` / `MinHostScratchSize` / `HostRegionExports` / `HostRegionAlignment` / `HostRegionDeclaration` / `ValidateParamsWrite` | **frozen** | 所有权写死在常量里；缺失 ⇒ `ErrHostRegionMissing`，非法 ⇒ `ErrHostRegionInvalid` |
| 预算上限 | `MaxParamsBytes` / `MaxScratchPerBar` / `MaxSignalsPerBar` / `MaxStateKeys` / `MaxStateValueBytes` / `MaxInstanceStateBytes` | **frozen** | 超限是 `Err*BudgetExceeded` / `ErrParamsTooLarge`，不是静默截断 |
| 值类型 / 字段号 | `HostI32` / `HostI64` / `HostF64` / `HostValueType` / `FieldClose`…`FieldVolume` / `FieldID` / `Action`…`ActionClose` / `Actions` / `FieldIDs` | **frozen** | **数字位就是线上协议**：错一位是「把 volume 当 close 用」——不报错、结果还看着合理（`abi.go` 里 FieldID 注释的原话）。非法值判据分三种，别混：`FieldID.Valid()` / `Action.Valid()` 只返回 **bool**（无哨兵）；`Action` 经 `String()` / `Direction()` 越界才给 `ErrActionIllegal`；`HostValueType` **没有** `Valid()`，它越不越界由导入签名比对拦（`ErrImportSignatureMismatch`） |
| scope | `ScopeMarket` / `ScopeSymbol` / `StrategyScope` / `ResolveScope` | **frozen** | v1 给 market scope ⇒ `ErrScopeMarketV1`（**目标** D25「M0 市场级先行」要解开这一约束，但今天代码是拒的） |

---

## 8. 不登记的东西（边界与非目标）

- **`dto/`**：跨进程 DTO 是 `tellurion:docs/openapi.yaml` 的**生成物**（`oapi-codegen`），真相在 core，本仓只收生成物且**禁手改**（CI 重生成 diff 必须为空）。见 `tellurion:docs/design/qt-decoupling/repo-inventory.md` §3。所以本文件不列 DTO 字段。
- **`schemas/`**：数据契约资产，不是 Go 包，不产生导出符号。
- **内部函数 / 未导出符号 / 结构体字段**：不是导出面（§1 粒度）。
- **默认实现、mock、测试替身**：`api-schema.md` §5 已禁（契约仓被测试代码污染）。

### 待裁决（#1~#3 已于落码时裁决；#4 已登记为 X1a 待办）

| # | 事项 | 现状 |
|---|---|---|
| 1 | `errors/` 与 `abi/` 是否独立成包 | ✅ **已裁决（2026-10-22，落码时实测）：不独立成包，随 `engine/` 走**（[`adr/0003`](adr/0003-哨兵与-abi-不独立成包.md)）。<br>**据实测而非估计**：`contracts.Err*` 消费方 10 个文件（非测试 4 + 测试 6），`contracts.Host*` 消费方 4 个文件（全在 `tellurion:pkg/sandbox/wasm`）；**没有任何一个消费方只想要其中一组**——每个文件同时用到 ≥2 个其它 `contracts.*` 符号（最少 2、最多 35）。⇒ 拆包不会让任何调用点少 import 一个包。<br>**代价**：`engine/` 成为 10 文件 / 2542 行（非测试，本仓拷贝实测）、194 个导出符号的包（占本仓 429 的 45%）。判据是「包内是否有真实的独立演化需求」——今天没有；等有（例如 abi 要独立发给 artifact 加载器）再拆，拆包本身是可排期的独立变更。<br>⚠️ **旧估计已作废**：前期笔记里的「约 61 个 core 文件」没有来源，实测是 4+6；结论不变，理由从「改动太大」换成「改了没人受益」 |
| 2 | `*Config` 命名族要不要进本仓 | ✅ **已裁决（2026-10-22，落码时实测）：[`adr/0002`](adr/0002-config-族的边界处置.md) 按「类型在边界上的角色」分三条，不按名字分**。<br>① **零引用的部署配置不迁**：`Config` / `DatabaseConfig` / `RedisConfig` / `ServiceConfig` / `TushareConfig` —— core 与 agent 实测 **0 引用**（连 core 自己都不用），且 `Config` 带 `Password`，描述的是 core 的部署 YAML。留在 core，由切换依赖那次 commit 一并删除（0 调用点）。本仓**已删除**它们（`domain/types.go`），支撑它们的恒真测试 `TestConfig_Types` 同删。<br>② **core-only 但随整包迁的两条按原名保留，标 experimental**：`StrategyConfig`（core 4 文件）、`ExecutionConfig`（core 21 文件）——进 ③-⑨ 回收清单；不就地改名的理由是它们本就要搬走（`StrategyConfig` → 最可能 `tellurion-strategy`；`ExecutionConfig` → 大概率回 core，因为执行留 core），改名是给即将搬走的类型付一次破坏性变更的账。<br>③ **冻结面伴生名保留**：`TradingConfig` / `PriceLimitConfig`（随 `engine/` 落地）已写在 SPEC 的 RunSpec 冻结文本里，改名 = 改冻结文本 + 弄红 core 的 `spec_drift_test.go`——冻结面优先于命名规范。<br>判据（以后照此判）：名字是否已在冻结文本 / 跨语言契约里 → 是则保留；否则看它是否 core-only 搭车 → 是则不迁或迁后回收；只有「长期留在本仓 + 名字泄漏实现细节」才改名（破坏性，走 ADR） |
| 3 | `fees/` 的取舍 | ✅ **已裁决（2026-10-22，落码时实测）：整包迁入，`StampTaxRateFor` / `ApplyDefaults` / `Validate` 一起进**。<br>**裁决依据是「④ 那句在今天的源码里没有对应物」**：`tellurion:pkg/fees` 一共 2 个文件 533 行（非测试 1 文件 250 行），**零算术**——导出面只有 7 个常量 + 1 个日期 var + 1 个结构体 + 4 个纯函数（取档选择器 / 默认值构造 / 零值填充 / 区间校验），全是 ③-① 说的「费率/滑点/成本模型的定义与枚举」。真正算钱的两处都在 core：`tellurion:pkg/backtest/tracker/tracker.go:145`（成交额 × 费率 + 最低佣金兜底，并且是 `StampTaxRateFor` 的**唯一生产调用方**）与 `tellurion:pkg/live/mock_trader.go`。<br>**为什么不把 `StampTaxRateFor` 留 core**：它与 `StampTaxCutDate` 是一对——源码注释自己写着「This is the ONLY place the cut date may be compared」（削减日只许在这一个地方比较）。把常量放本仓、比较点放 core，等于把这条不变量拆到两个仓，以后改档期的人只会看到其中一半。<br>**代价**：本仓因此有一个「看起来像行为」的函数（`StampTaxRateFor`）。判据是它**不读任何状态、不做算术、只按日期在两个给定值里选一个**，属"取值集的定义"而非"算钱"；`ApplyDefaults`/`Validate` 同理（填零值 / 判区间）。真出现按费率算成本的函数，那才越界，按 ADR-0001 ② 打回 |
| 4 | `dataprotocol` 是否独立成包（QT-40 对表 T2，2026-10-22 登记） | ✅ **已裁决（2026-10-24，QT-40 X1a 落码）：走 (a) 先迁** —— 包已落地（七字段 + `ToJSON` / `DeduplicateKey` + 15 个 `DataType*` 常量；无构造器 / 无时钟 / 无 `Aggregator`，见 §4 + §5.5）。判据是「第二份独立实现已经存在」，不是「第二个 import」。**未随本单落（如实）**：core 删两份副本 / data 镜像改 import / storage 列布局对表测试，归各自仓的单。**以下是裁决前的材料（保留来历）**：**落点已定名**：`contracts/dataprotocol`（包名 `dataprotocol`）——不叫 `data/`（包名要能自证「这是协议不是实现」，与 #1 的 `errors/`/`abi/` 之争同型），不叫 `dataset`（那是 X1b 的 aggregator 面，形状还没定）。依据 `tellurion:docs/design/qt-decoupling/data-service-contract-and-ownership.md` §4.1 A 组注 ①。**要裁的是迁入时机，而且两份上游文档给的判据并不一致**：「只有单个仓使用的类型不许提前迁」对应的是 ①「本仓只放被 ≥2 个仓使用的类型」（[AGENTS.md](../AGENTS.md) §3 把它写成「只有单个仓使用的类型属泄漏」）；对表注 ② 与本仓下游的 `tellurion-data:docs/api-schema.md` §4 **原先都把它引成 ③**（③ 管的是 `Config`/`Options`/`Impl` 命名），两处引用已于 2026-10-22 同步订正。**但同一份 ADR 的「代价」节明写 S1 整包迁时接受 core-only 类型临时进来、输出一份清单、S2 后回收**，而本仓台账已把这条落成硬要求（`docs/TASKS.md` QT-05 ③-⑨：`screen`/`sector`/`supply_chain`/`macro`/`mining` 在 agent 侧各 0 命中，仍整包迁）。于是有两条都自洽的走法：(a) `dataprotocol` 不适用那条容忍（它是**名单外的新建包**，不是整包迁的附带品）⇒ 挂 X1b 第一 commit，`tellurion-data` 一 import，两个使用方立刻成立；(b) 把它当同类容忍处理，随 S1 一起落、进那份 core-only 清单（③-⑨ 那份名单今天不含它，选 (b) 就要把它加进去）。**X1a 待办已由 QT-40 X1a 收口（2026-10-24）**：选 (a)；本仓 §4 / §5.5 已写回，`tellurion:docs/design/qt-decoupling/data-service-contract-and-ownership.md` §4.1 注 ② 的回写归 core 侧协调者（本仓不改 core 文档）；X1a 的施工序在 `tellurion-data:docs/TASKS.md` 施工序表。冻结前要先定形状的三处实测：① `UnifiedDataPoint` 今天有**两份同名同字段副本**（`tellurion:pkg/data/source/unified.go:16` 与 `pkg/storage/bulk_insert.go:17`，后者注释自述「storage 不能 import source（会成环）」）⇒ 进 contracts 后成环理由消失，**必须删副本改引用**，否则协议加一列时 `BulkInsert` 的列布局会静默少写一列（`bulk_insert_ddl_test.go` 只核「目标表在 DDL 里」，核不到字段）；② **没有 `DataType` 类型**——`ConnectorDeclaration.DataType` 与 `UnifiedDataPoint.DataType` 都是裸 `string`，取值集靠 `DataType*` 常量；③ **没有 `Aggregator` 接口**（§5.2 的 port 表里只有未定形的 `DataPort`，标着「⬜ 设计中，S1 不落」）：aggregator 今天只有 HTTP 形状，要么随 X1b 实现定、要么先不落 port，**不许在文档里造一个无实现的接口**（§8 已禁默认实现/mock 进本仓）。另：`TimeSemantics` 枚举定义在 `tellurion:pkg/domain/market/visibility.go:61`，**已随 QT-05 ③-① 整包迁入 `domain/market/`（2026-10-22，见 §4）**，`dataprotocol` 里**不再重复定义**——同一次迁移出现两份 `TimeSemantics`，等于把 #1 与 #4 合成一个真事故 |

---

## 9. 与其他文档的互锁

| 文档 | 关系 |
|---|---|
| `api-schema.md` | 规则源（错误分类 / 幂等 / 兼容 / 不提供什么）；§1.1 的 port「状态」列与本文件 §5.2 一致，§1.3 的三分法与 §1.4 的 abi 计数都由本文件 §5.4/§7 逐条回源码核对（2026-10-21 自审时改掉过那两处的「写得像事实」：哨兵零定义、host 函数 13→12） |
| `ARCHITECTURE.md` §2 | 包布局与本文件 §8 待裁决 #1 有关 |
| `tellurion:docs/SPEC.md` | **冻结文本的唯一源**；本文件只指节标题。改冻结面 = 先破坏 core 合规测试 → 更新 SPEC → 变更评审 → 才动代码（`tellurion:docs/TASKS.md`「维护本文件」硬规则） |
| `tellurion:docs/design/qt-decoupling/kernel-interface-map.md` | 现状接口面清单的来源（M6 的「散在几处」即指它）；本文件不是它的替代，是**随仓可核对**的那一份 |
