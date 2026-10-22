---
status: active
last-verified: 2026-10-22
verified-by: QT-05 ③-① 落码进行中（`fees/` 13 个符号 + `domain/market/` 46 个符号已登记，护栏绿 59/59）；§8 待裁 #3 已裁决（fees 整包迁入，实测本包零算术）
---

# API — 导出面规格（本仓有哪些符号、各自稳不稳定、错了怎么判）

> **与 [`api-schema.md`](api-schema.md) 的分工**：那份定**规则**（错误分类、幂等、版本兼容、不提供什么）；
> 本文件定**清单**（导出符号逐个回答「它在不在、稳不稳定、失败时是什么」）。
> 符号级冲突以本文件为准（它由 CI 核对，`api-schema.md` 不核对符号名）；规则级冲突以 `api-schema.md` 为准。
>
> **现状口径（开工前必读）**：本仓**正在落码中**（QT-05 ③-①，2026-10-22 起）——今天已落地的库包是 `fees/` 与 `domain/market/`，
> `domain/`、`engine/` 按「一次只搬一个包」依次落，每落一个就在 §4 补一节并跑护栏。
> `tools/` 下的命令是 `package main`，不计入导出面。
> **§5~§7 仍是目标面快照**（2026-10-21 用 §1 的命令从 core 实测抓的），只有已经出现在 §4 的包才是「现在能 import 的东西」；
> §4 没登记的包，就算 §5~§7 描述了它，也**还不能 import**。

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
| 登记表形状 | 包小节标题必须带反引号包目录（形状见下面的示例）；**根包写 `.`**——不认根包就会让「根包新增导出符号」无路可登记，那是不可满足的护栏。标题不带目录时，那一节的表行**不算登记**：写错标题的人被「未登记」正面挡住，而不是悄悄登记进一个不存在的包 |

> ⚠️ **覆盖面限制（如实登记，别把它当全知护栏）**：内嵌接口（`type P interface { Base; Extra() }`）的**继承来的**方法名不在扫描粒度内——
> 那些方法本来就属于已登记的 `Base`，不是新符号。真要新增方法，新增的那个会被扫到。
> 结构体字段同理：字段增删**不触发**本护栏，它们由 core 侧的 SPEC 合规测试与 `resultspec_test.go` 那类字段级钉子管（见 §5）。

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

> 第 3、4 条说的三类（`ErrDataGap` / `ErrStrategyFault` / `ErrContractBreach`）**今天还不存在**（§5.4 标 ⬜，QT-10 ③-① 才落）。
> 落码前 provider 侧返回的是未分类 error——所以这两条现在是**验收标准**，不是现状描述。

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

### `fees/` — 包 `fees`（A 股费率模型的数据形状；2026-10-22 自 `tellurion:pkg/fees` 逐字迁入）

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
| `AShareFees.Validate` | method | stable | 区间即判据：三个比率 ∈ **[0, 0.05)**、`MinCommission` ∈ [0, 1000)、`SlippageRate` ∈ [0, 0.1)；越界返回 error（**不是** `AppError`，本包不依赖 kit）。**0 是刻意合法值**（免佣 VIP / 确定性测试关滑点），由 `TestValidate_AcceptsBoundary` 钉住——⚠️ 源码里那条 doc 注释写成 `(0, 0.05)` 与实现不符，**以实现与本行为准**（QT-05 ④ 不改业务逻辑，注释订正留给 core 侧同一次改动）。另：**零值 `AShareFees{}` 也能通过 `Validate`**（全 0 在区间内），所以"先 Validate 再 ApplyDefaults"拦不住"整组费率忘配"——那是 `ApplyDefaults` 的活，不是它的。**放宽区间属语义变更**（原本被拒的配置会通过），按版本规则记账 |

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
| `Provider` | interface | experimental | 行情数据访问接口（原 `domain.MarketDataProvider`，S7-P3-4 改名，旧名以别名保留）。⚠️ **实测今天 0 实现 / 0 消费**：core 全仓除定义、别名与别名测试外无引用（`pkg/marketdata` 用的是它**自己的** `Provider`，与这个不是同一个类型），agent 侧 0 命中。所以它标 experimental 而不是 stable——一个没人实现的 port，稳定性无从谈起。X1a（data 服务拆分 / `DataPort`）会重新裁这个形状，届时走 ADR，不许两边各留一份 |
| `Provider.GetOHLCV` | method | experimental | 按 symbol + 区间取 bar，返回 `[]OHLCV`；**今天返回未分类 error**，QT-10 ③-② 落地后缺数据须包成 `ErrDataGap`（§5.4） |
| `Provider.GetFundamental` | method | experimental | 按 symbol + 日期取基本面（PIT：用可用日，不是报告期）；错误分类同上 |
| `Provider.GetStocks` | method | experimental | 按交易所取证券列表（含退市，见 `Stock.DelistDate`）；错误分类同上 |
| `Provider.GetLatestPrice` | method | experimental | 最新价（只给价格，不给时间戳——需要 PIT 判定的场景不要用它）；错误分类同上 |
| `Provider.GetIndexConstituents` | method | experimental | 取指数成分代码；错误分类同上 |

<!-- surface:end -->

---

## 5. 冻结面逐符号规格（S1 迁入 §4 时的权威口径）

> 每个冻结面只**指回** `tellurion:docs/SPEC.md` 的节标题（SPEC.md 用命名节，不用数字编号），不复制正文——
> 复制的那一份一定会漂移。改这些节本身走 `tellurion:docs/TASKS.md`「维护本文件」的**冻结 schema 变更评审（硬规则）**。

### 5.1 四契约（ADR-031 D1 —— `tellurion-docs:docs/adr/adr-031-kernel-pluggable-nt-ready.md`）

| 契约术语 | 承载符号（今天住 core，S1 → `engine/`） | 稳定性 | 冻结文本指回 |
|---|---|---|---|
| `RunSpec` | `BacktestRequest` + `BacktestRequest.ValidateRunSpec` / `EffectiveKernel` / `EffectiveFillTiming`；伴生：`SignalWeight`、`TradingConfig` + `TradingConfig.WithDefaults`、`PriceLimitConfig`、`FillTiming` | **frozen** | `tellurion:docs/SPEC.md` §「RunSpec 内核路由（K10，2026-10-12）：EngineRunner 可插拔」 |
| `SignalSpec` | `SignalWeight`（K15 目标权重流）；信号本体 `domain.Signal` 随 `domain/` 迁入；供给契约 `SignalProvider`（今天 `tellurion:pkg/strategy`，见 §5.2） | stable（**SPEC.md 无 SignalSpec 节**） | `tellurion-contracts:docs/api-schema.md` §1.2 + `tellurion:docs/design/qt-decoupling/kernel-interface-map.md` §一.4 |
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
| `SignalProvider` | `Signals(ctx, Scope, asOf) ([]Signal, error)` | core（`tellurion:pkg/strategy`） | stable（S1 落点待 QT-05） | 返回 `as_of > asOf` 的信号即违约：`RejectFutureSignals` fail-loud，不静默跳过 |
| `Provider`（`domain/market`，别名 `MarketDataProvider`） | `GetOHLCV` / `GetFundamental` / `GetStocks` / `GetLatestPrice` / `GetIndexConstituents` | **无**（实测 0 实现 / 0 消费：core 除定义/别名/别名测试外 0 引用，`pkg/marketdata` 用的是它自己的同名接口） | **experimental**（原标 stable，2026-10-22 落码时按实测下调——没人实现的 port 没有稳定性可言；§4 已同步） | **今天返回未分类 error**；QT-10 ③-② 落地后缺数据须包成 `ErrDataGap`（跳过当天 + 写 `invalid_reasons`），不是 fault —— 见 §5.4 的「⬜ 目标」行 |
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
| 内核路由 | `ErrKernelNotRegistered`、`ErrMarketNotRegistered` | stable（今天在 `tellurion:pkg/backtest/contracts`，S1 → `engine/`） | 请求期 400；异步 job 终态 `failed`。**零静默回退** |
| 成交口径 | `ErrFillTimingMissing`、`ErrFillTimingUnknown`、`ErrFillTimingMismatch` | stable | 引擎永远不猜口径；口径不符是错误不是结果差异 |
| abi / artifact 校验 | `ErrAbiVersionUnknown`、`ErrAbiVersionMismatch`、`ErrAbiFamilyMissing`、`ErrAbiFamilyMixed`、`ErrImportNotWhitelisted`、`ErrImportSignatureMismatch`、`ErrHostRegionMissing`、`ErrHostRegionInvalid`、`ErrParamsTooLarge`、`ErrScopeUnknown`、`ErrScopeMarketV1` | stable（S1 → `engine/`；`errors/` 是否独立成包见 §8 待裁决） | 归 `ErrContractBreach` 语义：装载期拦下，不进 run |
| 信号校验 | `ErrActionIllegal`、`ErrStrengthOutOfRange`、`ErrTargetWeightOutOfRange`、`ErrEmitSourcesMixed`、`ErrEmitSourceUnknown`、`ErrSignalBudgetExceeded`、`ErrStateBudgetExceeded` | stable | 归 `ErrStrategyFault`（非法信号 / 预算超限 = 策略故障，终止 run） |
| provenance / hash | `ErrCanonicalUnsupported`、`ErrStrategyRefMalformed` | stable | 输入身份不可 canonical 编码 → 当场失败，禁止退化成「猜一个 hash」 |

---

## 6. 值对象面（`domain/`，S1 整包迁入）

来源与再生成：**本仓内**跑 `go run ./tools/check_api_surface -list`（已落地的包）；尚未落地的部分在 core 侧用 `-root <tellurion 仓>/pkg/domain` 取（2026-10-21 实测 227 个导出符号，含 `market/` 子包 ⇒ 其中 **46 个已于 2026-10-22 落地并登记在 §4**，余下是 `domain/` 根包）。
**逐符号清单不在这里手抄**——抄一份就会漂；落码时用上面这条命令生成 §4 的表，再按下面分组填稳定性。

| 组 | 代表符号 | 稳定性 | 备注 |
|---|---|---|---|
| 回测输入输出 | `BacktestResult`、`BacktestParams`、`PortfolioValue`、`Trade`、`RiskMetrics`、`WalkForward*` | stable（`BacktestResult` 的**指标口径**受 SPEC 冻结节约束） | `BacktestResult` **没有也不该有 `Status`**——三值属响应层，见 `tellurion-agent:docs/TASKS.md` OBS-01-R2-A 的 core 裁决 (b) |
| 仪器域行情 | `OHLCV`、`Stock`、`Fundamental(Data)`、`IndexDaily`、`IndexConstituent`、`Dividend`、`Split`（多为 `market/` 的类型别名） | stable | 别名不是重复真相：实体在 `domain/market`，根包别名保兼容 |
| 时间/频率语义 | `Frequency`、`FrequencySpec`、`TimeSemantics`、`AsiaShanghai`、`Visible*`、`DailyVisibleAt` | **frozen** | PIT 红线：`tellurion:docs/SPEC.md` §「时间语义（PIT 红线）」。⚠️ 说清来源——core「冻结 schema 变更评审」点名的名单是三契约（RunSpec/ResultSpec/SignalSpec），这一行的 frozen 是**本文件按 §2 定义给的标注**，依据是它已有注册期闸门（D2）与逐源钉子（`tellurion:pkg/data/source/registry_rev05_test.go` 等把每个 connector 的 `TimeSemantics` 钉死）；改语义同样要先让那些测试红 |
| 信号与订单 | `Signal`、`Direction`、`Order`、`OrderType`、`Position(Size)`、`Portfolio`、`TargetPosition`、`StopLossEvent`、`MarketRegime` | stable | — |
| 因子族 | `FactorDefinition`、`FactorHypothesis`、`FactorType`、`FactorKernel`、`FactorCaliber`、`FactorRegistration`、`FactorMiningRecord`、`FactorOutput`、`ICEntry`、`FactorValue` + 内置定义/假设种子变量 | stable（`VersionHash` 语义受「版本哈希纪律」节约束） | 封闭词表扩展 = 加名 + 加实现 + 同步 SPEC 对应节 |
| 公司行为 | `CorporateAction`、`CashDividend`、`BonusShare`、`CorporateActionSplit`、`RightsIssue`、`Placement`、`ActionEngine`、`ApplyResult`/`ApplyOutcome` | stable | — |
| 配置形状 | `Config`、`ServiceConfig`、`DatabaseConfig`、`RedisConfig`、`TushareConfig`、`StrategyConfig`、`ExecutionConfig` | ⚠️ **experimental，且命名待评审** | 全部命中本仓 [`adr/0001`](adr/0001-契约只增不破且零实现.md) ③ 的 `Config` 命名禁令（`Config` 是精确同名，其余是同一种形状：契约位上写的是实现细节）。**S1 整包迁会把它们带进来**——落码前定「改名 / 留 core / 立 ADR 豁免」，见 §8 待裁决 #2，别默默收下 |
| core-only（agent 侧 0 命中） | QT-05 ③-⑨ 点名的 5 族：`Screen*`、`Sector`/`StockSector`/`SWIndustryNode`、`SupplyChainEdge`、`MacroSeries`、`mining` 族（`FactorMiningRecord`）；实测另含 `IndustryVitality`、`MarketActivity`、`Factor*`（定义/登记/输出）、`CorporateAction` 族、`RiskManager`、`Order`/`Signal`/`Stock` 等 | **experimental + 待回收** | QT-05 ③-⑨ 与 ADR-0001 的代价条款：S1 仍整包迁，S2 后按清单回收。**下游不要依赖这一组**。<br>实测口径（2026-10-21）：`tellurion:pkg/domain` 根包 124 个导出符号里 **86 个在 agent 侧 0 命中**（`grep -rl "domain[.]<名>\b" --include='*.go'` 逐名数文件）。「0 命中」≠「该删」——agent 走 HTTP/JSON 消费的那几族（`Signal`/`Stock`）本来就不摸 Go 类型。完整清单是 **QT-05 ③-⑨ 的产出物**，本文件只登记判据与命令，不替它下结论 |

---

## 7. abi 常量族（`engine/` 内，随 QT-05 迁入）

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

### 待裁决（S1 落码前要有结论，别默默带过）

| # | 事项 | 现状 |
|---|---|---|
| 1 | `errors/` 与 `abi/` 是否独立成包 | [`ARCHITECTURE.md`](ARCHITECTURE.md) §2 列了这两个包，但 QT-05 ③-① 只说迁 `domain/`、`engine/`、`fees/`（`engine/` 内保持 `package contracts`），abi 常量与哨兵今天都随 `engine/` 走。二者口径不一致，本文件按 QT-05 落点写（§5/§7 归 `engine/`） |
| 2 | `*Config` 命名族要不要进本仓 | ADR-0001 ③ 禁 `Config`/`Options`/`Impl` 命名，但 S1 整包迁会带进来 9 个：`domain/` 的 `Config`（精确同名）、`ServiceConfig`、`DatabaseConfig`、`RedisConfig`、`TushareConfig`、`StrategyConfig`、`ExecutionConfig` + `engine/` 的 `TradingConfig`、`PriceLimitConfig`。后两个还是 §5.1 的 frozen 伴生符号——禁令与冻结面在这儿撞车，落码前要结论，不能两边都当没看见 |
| 3 | `fees/` 的取舍 | ✅ **已裁决（2026-10-22，落码时实测）：整包迁入，`StampTaxRateFor` / `ApplyDefaults` / `Validate` 一起进**。<br>**裁决依据是「④ 那句在今天的源码里没有对应物」**：`tellurion:pkg/fees` 一共 2 个文件 533 行（非测试 1 文件 250 行），**零算术**——导出面只有 7 个常量 + 1 个日期 var + 1 个结构体 + 4 个纯函数（取档选择器 / 默认值构造 / 零值填充 / 区间校验），全是 ③-① 说的「费率/滑点/成本模型的定义与枚举」。真正算钱的两处都在 core：`tellurion:pkg/backtest/tracker/tracker.go:145`（成交额 × 费率 + 最低佣金兜底，并且是 `StampTaxRateFor` 的**唯一生产调用方**）与 `tellurion:pkg/live/mock_trader.go`。<br>**为什么不把 `StampTaxRateFor` 留 core**：它与 `StampTaxCutDate` 是一对——源码注释自己写着「This is the ONLY place the cut date may be compared」（削减日只许在这一个地方比较）。把常量放本仓、比较点放 core，等于把这条不变量拆到两个仓，以后改档期的人只会看到其中一半。<br>**代价**：本仓因此有一个「看起来像行为」的函数（`StampTaxRateFor`）。判据是它**不读任何状态、不做算术、只按日期在两个给定值里选一个**，属"取值集的定义"而非"算钱"；`ApplyDefaults`/`Validate` 同理（填零值 / 判区间）。真出现按费率算成本的函数，那才越界，按 ADR-0001 ② 打回 |
| 4 | `dataprotocol` 是否独立成包（QT-40 对表 T2，2026-10-22 登记） | **落点已定名**：`contracts/dataprotocol`（包名 `dataprotocol`）——不叫 `data/`（包名要能自证「这是协议不是实现」，与 #1 的 `errors/`/`abi/` 之争同型），不叫 `dataset`（那是 X1b 的 aggregator 面，形状还没定）。依据 `tellurion:docs/design/qt-decoupling/data-service-contract-and-ownership.md` §4.1 A 组注 ①。**要裁的是迁入时机，而且两份上游文档给的判据并不一致**：「只有单个仓使用的类型不许提前迁」对应的是 ①「本仓只放被 ≥2 个仓使用的类型」（[AGENTS.md](AGENTS.md) §3 把它写成「只有单个仓使用的类型属泄漏」）；对表注 ② 与本仓下游的 `tellurion-data:docs/api-schema.md` §4 **原先都把它引成 ③**（③ 管的是 `Config`/`Options`/`Impl` 命名），两处引用已于 2026-10-22 同步订正。**但同一份 ADR 的「代价」节明写 S1 整包迁时接受 core-only 类型临时进来、输出一份清单、S2 后回收**，而本仓台账已把这条落成硬要求（`docs/TASKS.md` QT-05 ③-⑨：`screen`/`sector`/`supply_chain`/`macro`/`mining` 在 agent 侧各 0 命中，仍整包迁）。于是有两条都自洽的走法：(a) `dataprotocol` 不适用那条容忍（它是**名单外的新建包**，不是整包迁的附带品）⇒ 挂 X1b 第一 commit，`tellurion-data` 一 import，两个使用方立刻成立；(b) 把它当同类容忍处理，随 S1 一起落、进那份 core-only 清单（③-⑨ 那份名单今天不含它，选 (b) 就要把它加进去）。X1a 之前要挑一条并写进本节，**别两份文档各按一条**（现在的写法是注 ② 选 (a)、而它引的条款其实同时允许 (b)）。X1a 冻结前要先定形状的三处实测：① `UnifiedDataPoint` 今天有**两份同名同字段副本**（`tellurion:pkg/data/source/unified.go:16` 与 `pkg/storage/bulk_insert.go:17`，后者注释自述「storage 不能 import source（会成环）」）⇒ 进 contracts 后成环理由消失，**必须删副本改引用**，否则协议加一列时 `BulkInsert` 的列布局会静默少写一列（`bulk_insert_ddl_test.go` 只核「目标表在 DDL 里」，核不到字段）；② **没有 `DataType` 类型**——`ConnectorDeclaration.DataType` 与 `UnifiedDataPoint.DataType` 都是裸 `string`，取值集靠 `DataType*` 常量；③ **没有 `Aggregator` 接口**（§5.2 的 port 表里只有未定形的 `DataPort`，标着「⬜ 设计中，S1 不落」）：aggregator 今天只有 HTTP 形状，要么随 X1b 实现定、要么先不落 port，**不许在文档里造一个无实现的接口**（§8 已禁默认实现/mock 进本仓）。另：`TimeSemantics` 枚举定义在 `tellurion:pkg/domain/market/visibility.go:61`，**已随 QT-05 ③-① 整包迁入 `domain/market/`（2026-10-22，见 §4）**，`dataprotocol` 里**不再重复定义**——同一次迁移出现两份 `TimeSemantics`，等于把 #1 与 #4 合成一个真事故 |

---

## 9. 与其他文档的互锁

| 文档 | 关系 |
|---|---|
| `api-schema.md` | 规则源（错误分类 / 幂等 / 兼容 / 不提供什么）；§1.1 的 port「状态」列与本文件 §5.2 一致，§1.3 的三分法与 §1.4 的 abi 计数都由本文件 §5.4/§7 逐条回源码核对（2026-10-21 自审时改掉过那两处的「写得像事实」：哨兵零定义、host 函数 13→12） |
| `ARCHITECTURE.md` §2 | 包布局与本文件 §8 待裁决 #1 有关 |
| `tellurion:docs/SPEC.md` | **冻结文本的唯一源**；本文件只指节标题。改冻结面 = 先破坏 core 合规测试 → 更新 SPEC → 变更评审 → 才动代码（`tellurion:docs/TASKS.md`「维护本文件」硬规则） |
| `tellurion:docs/design/qt-decoupling/kernel-interface-map.md` | 现状接口面清单的来源（M6 的「散在几处」即指它）；本文件不是它的替代，是**随仓可核对**的那一份 |
