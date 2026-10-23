---
status: active
last-verified: 2026-10-22
verified-by: QT-05 ③-① 本仓侧落码完成（`fees` 13 + `domain/market` 46 + `domain` 176 + `engine` 194 = 429 个符号已登记，护栏 429/429 绿）；§8 待裁 #1（errors//abi/，ADR-0003）、#2（`*Config` 族，ADR-0002）、#3（fees 整包）均已裁决；#4（dataprotocol 时机）已登记为 X1a 待办（2026-10-22 推送回写：6 个 commit 已进 `origin/main`，该批 HEAD = `fbe9a2d`）
---

# API — 导出面规格（本仓有哪些符号、各自稳不稳定、错了怎么判）

> **与 [`api-schema.md`](api-schema.md) 的分工**：那份定**规则**（错误分类、幂等、版本兼容、不提供什么）；
> 本文件定**清单**（导出符号逐个回答「它在不在、稳不稳定、失败时是什么」）。
> 符号级冲突以本文件为准（它由 CI 核对，`api-schema.md` 不核对符号名）；规则级冲突以 `api-schema.md` 为准。
>
> **现状口径（开工前必读）**：本仓的**库代码已全部落地**（QT-05 ③-①，2026-10-22）：`fees/`、`domain/market/`、`domain/`、`engine/` 四个包，
> 429 个导出符号全部登记在 §4。`tools/` 下的命令是 `package main`，不计入导出面。
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
| `Provider` | interface | experimental | 行情数据访问接口（原 `domain.MarketDataProvider`，S7-P3-4 改名，旧名以别名保留）。⚠️ **实测今天 0 实现 / 0 消费**：core 全仓除定义、别名与别名测试外无引用（`pkg/marketdata` 用的是它**自己的** `Provider`，与这个不是同一个类型），agent 侧 0 命中。所以它标 experimental 而不是 stable——一个没人实现的 port，稳定性无从谈起。X1a（data 服务拆分 / `DataPort`）会重新裁这个形状，届时走 ADR，不许两边各留一份 |
| `Provider.GetOHLCV` | method | experimental | 按 symbol + 区间取 bar，返回 `[]OHLCV`；**今天返回未分类 error**，QT-10 ③-② 落地后缺数据须包成 `ErrDataGap`（§5.4） |
| `Provider.GetFundamental` | method | experimental | 按 symbol + 日期取基本面（PIT：用可用日，不是报告期）；错误分类同上 |
| `Provider.GetStocks` | method | experimental | 按交易所取证券列表（含退市，见 `Stock.DelistDate`）；错误分类同上 |
| `Provider.GetLatestPrice` | method | experimental | 最新价（只给价格，不给时间戳——需要 PIT 判定的场景不要用它）；错误分类同上 |
| `Provider.GetIndexConstituents` | method | experimental | 取指数成分代码；错误分类同上 |

### `domain/` — 包 `domain`（信号 / 订单 / 组合 / 公司行为 / 因子名册；2026-10-22 自 `tellurion:pkg/domain` 迁入）

> **迁入形状**：23 个文件（16 非测试 + 7 测试）自 core 拷入，刻意的改动只有三类：① import 换成本仓路径（`domain/market`、`fees`）；② **删掉 5 个零引用的部署配置类型**（`Config` / `DatabaseConfig` / `RedisConfig` / `ServiceConfig` / `TushareConfig`——core 与 agent 实测 0 引用；`Config` 还带 `Password` 字段，不该进公开库；裁决见 §8 #2 与 [`adr/0002`](adr/0002-config-族的边界处置.md)）；③ 随之删掉只为它们存在、且断言恒真的 `TestConfig_Types`。
> **等价性证据**：剥注释 + 归一化 import 后，20/23 文件与 core 完全相同；`execution_test.go` 只差 gofmt 的 import 顺序；`types.go` / `types_test.go` 只差被删的那两段（−31 / −24 行）。
> **稳定性口径（本节）**：`frozen` = SPEC 互锁（`OHLCV` / `Fundamental` 别名）；`stable` = 引擎 / 策略 / 风控的领域词汇（改 = MINOR，破坏 = MAJOR）；`experimental` = 配置形状（§8 #2）与 core-only 数据族（QT-05 ③-⑨ 回收清单）。core-only 的族按 §6 的功能分组标 stable，并在备注里写明「core-only，③-⑨」——两处判断都留痕。
> **别名不是重复真相**：`OHLCV` / `Stock` / … / `MarketDataProvider` 的实体在 `domain/market`，这里只为兼容旧拼写（`domain.X` 与 `market.X` 是**同一个类型**）。
> **方法体边界**：本包里的 `Validate` / `VersionHash` / `CorporateAction.Apply` / `ActionEngine.ApplyAll` 是**纯函数行为**，随值对象迁入；判据（无 I/O、依赖受限、无跨调用持久状态、不替代服务）见 [`adr/0004`](adr/0004-值对象上的纯函数行为随类型走.md)。

| 符号 | 类别 | 稳定性 | 错误语义 / 判据 |
|---|---|---|---|
| `Direction` | type | stable | 交易方向枚举；JSON 字面量是契约（`long` / `short` / `close` / `hold`） |
| `DirectionLong` | const | stable | `"long"`（开多 / 加多） |
| `DirectionShort` | const | stable | `"short"` |
| `DirectionClose` | const | stable | `"close"`（平仓；卖出方向判定印花税时按它） |
| `DirectionHold` | const | stable | `"hold"`（不动作） |
| `OrderType` | type | stable | 订单类型枚举（P1-3 / ODR-016）；JSON 字面量是契约 |
| `OrderTypeMarket` | const | stable | `"market"`：立即按当前价成交 |
| `OrderTypeLimit` | const | stable | `"limit"`：仅价格穿越限价时成交 |
| `OrderTypeStop` | const | stable | `"stop"`：触及 `StopPrice` 后转市价 |
| `OrderTypeTrailing` | const | stable | `"trailing"`：跟踪 HWM 动态调整止损价 |
| `Signal` | type | stable | 策略信号（v1 强度流）：`Symbol` / `Date` / `Direction` / `Strength` / `CompositeScore` / `Factors` / `Metadata` / `LimitPrice` / `OrderType` + K17 可选执行算法三字段（`ExecAlgo` / `ExecSlices` / `ExecProfile`）。零值 `ExecAlgo` = 不拆单（默认路径逐字节不变） |
| `Order` | type | stable | 订单（含 K17 执行算法字段）。⚠️ 改字段要看 SPEC 的 RunSpec/ResultSpec 互锁，别只看本表 |
| `Position` | type | stable | 持仓：含 T+1 结算三字段（`BuyDate` / `QuantityToday` / `QuantityYesterday`）与 `Metadata` 开放槽。⚠️ T+1 字段是 A 股规则载体，删/改 = 破坏回测正确性 |
| `Portfolio` | type | stable | 组合快照（`Cash` / `Positions` / `TotalValue` / `DailyReturn` / `UpdatedAt`） |
| `TargetPosition` | type | stable | 目标仓位跟踪（目标 / 实际 / 未成交数量），供对账 |
| `PositionSize` | type | stable | 风控给出的仓位建议（大小 / 权重 / 止损 / 止盈 / 风险分） |
| `StopLossEvent` | type | stable | 止损 / 止盈触发事件（`Type` 是**裸字符串** `"stop_loss"` / `"take_profit"`，不是枚举——改字面量 = 破坏下游判定） |
| `MarketRegime` | type | stable | 市场状态（趋势 / 波动 / 情绪分 / 时间戳）。`Trend` / `Volatility` 也是裸字符串 |
| `RiskMetrics` | type | stable | 组合风险指标（波动 / Beta / Sharpe / Sortino / 最大回撤 / VaR …） |
| `RiskManager` | interface | stable | 风控**同步裁决**接口（零网络，留在进程内——ADR-023 / 分仓边界）。加方法对实现方是破坏性变更 |
| `RiskManager.CalculatePosition` | method | stable | 给定信号 + 组合 + 状态 + 现价 + bar，返回 `PositionSize`；错误即拒绝下单，**不许静默降级为 0 仓位以外的默认** |
| `RiskManager.DetectRegime` | method | stable | 由 bar 序列识别市场状态 |
| `RiskManager.CheckStopLoss` | method | stable | 扫描持仓 + 价格表，返回触发事件（空切片 = 无触发） |
| `OHLCV` | type | **frozen** | `= market.OHLCV`（别名）。字段集由 core 的 `spec_drift_test.go` 与 `tellurion:docs/SPEC.md` 互锁——加字段 core 侧红 |
| `Stock` | type | stable | `= market.Stock`（别名）。`DelistDate` nil = 仍在市（幸存者偏差） |
| `Fundamental` | type | **frozen** | `= market.Fundamental`（别名）。数值字段全指针（缺失 ≠ 0）；SPEC 互锁 |
| `FundamentalData` | type | stable | `= market.FundamentalData`（别名，持久化视图） |
| `IndexDaily` | type | stable | `= market.IndexDaily`（别名，指数日线，身份列 `Code`） |
| `IndexConstituent` | type | stable | `= market.IndexConstituent`（别名） |
| `Split` | type | stable | `= market.Split`（别名；Tushare 风格的拆股记录，**不是**公司行为引擎里的 `CorporateActionSplit`） |
| `Dividend` | type | stable | `= market.Dividend`（别名） |
| `MarketDataProvider` | type | experimental | `= market.Provider`（旧名别名）。与实体同命：**今天 0 实现 / 0 消费**，X1a 定 `DataPort` 时一并裁（见 §4 `domain/market/`） |
| `ActionType` | type | stable | 公司行为类型枚举（5 类，字面量即契约） |
| `ActionCashDividend` | const | stable | `"cash_dividend"`：每股派现；持仓不变、现金增加 `qty × cash_per_share` |
| `ActionBonusShare` | const | stable | `"bonus_share"`：送股；`qty × (1 + bonus_per_10/10)`，均价同比例下调 |
| `ActionSplit` | const | stable | `"split"`：拆股 / 并股；`qty × ratio`，均价 `/ ratio` |
| `ActionRightsIssue` | const | stable | `"rights_issue"`：配股；默认按**放弃**处理，缴款走 `ApplyPaid` |
| `ActionPlacement` | const | stable | `"placement"`：增发；持仓与现金都不变（只调整参考价 / 记事件） |
| `CorporateAction` | interface | stable | 公司行为统一接口（`Type` / `Symbol` / `ExDate` / `RecordDate` / `PayDate` / `Apply` / `Description`）。`Apply` 是**纯函数**：不改入参，返回新持仓 + 现金增量（正=收到） |
| `CorporateAction.Type` | method | stable | 返回判别子（`ActionType`） |
| `CorporateAction.Symbol` | method | stable | 返回 ts_code |
| `CorporateAction.ExDate` | method | stable | 除权除息日（此日开盘股价已调整、持仓与现金已结算） |
| `CorporateAction.RecordDate` | method | stable | 股权登记日（此日收盘后持股者享有本次行为） |
| `CorporateAction.PayDate` | method | stable | 实际到账 / 缴款日；零值 = 不适用 |
| `CorporateAction.Apply` | method | stable | `(pos, prevClose) → (newPos, cashDelta)`，纯函数 |
| `CorporateAction.Description` | method | stable | 审计用的人类可读描述（含日期与数量，格式非契约、内容可读即可） |
| `CashDividend` | type | stable | 现金分红实现（税**不在**引擎内处理，留 tax 模块） |
| `NewCashDividend` | func | stable | 构造函数（`symbol, exDate, recordDate, payDate, cashPerShare`） |
| `CashDividend.Apply` | method | stable | `qty ≤ 0` 或 `cashPerShare ≤ 0` → 原样返回、现金增量 0（**不报错**，由引擎按无持仓/无效事件跳过） |
| `CashDividend.Description` | method | stable | 见接口同名方法 |
| `CashDividend.ExDate` | method | stable | 见接口同名方法 |
| `CashDividend.PayDate` | method | stable | 见接口同名方法 |
| `CashDividend.RecordDate` | method | stable | 见接口同名方法 |
| `CashDividend.Symbol` | method | stable | 见接口同名方法 |
| `CashDividend.Type` | method | stable | `ActionCashDividend` |
| `BonusShare` | type | stable | 送股实现（`BonusPer10`，如 10 送 5 → 5） |
| `NewBonusShare` | func | stable | 构造函数 |
| `BonusShare.Apply` | method | stable | `qty × (1 + bonus/10)`；`AvgCost > 0` 时按同比例下调；`CurrentPrice > 0` 时顺带重算市值与浮盈（调用方通常会用最新价再算一遍） |
| `BonusShare.Description` | method | stable | 见接口同名方法 |
| `BonusShare.ExDate` | method | stable | 见接口同名方法 |
| `BonusShare.PayDate` | method | stable | 见接口同名方法 |
| `BonusShare.RecordDate` | method | stable | 见接口同名方法 |
| `BonusShare.Symbol` | method | stable | 见接口同名方法 |
| `BonusShare.Type` | method | stable | `ActionBonusShare` |
| `CorporateActionSplit` | type | stable | 拆股 / 并股实现。**名字刻意带前缀**：与 `types.go` 的 `Split`（Tushare 记录，别名到 `market.Split`）区分 |
| `NewSplit` | func | stable | 构造函数；`splitRatio` 必须 > 0（≤0 时 `Apply` 原样返回） |
| `CorporateActionSplit.Apply` | method | stable | `qty × ratio`、`AvgCost / ratio`；`ratio ≤ 0` → 原样返回 |
| `CorporateActionSplit.Description` | method | stable | 见接口同名方法 |
| `CorporateActionSplit.ExDate` | method | stable | 见接口同名方法 |
| `CorporateActionSplit.PayDate` | method | stable | 见接口同名方法 |
| `CorporateActionSplit.RecordDate` | method | stable | 见接口同名方法 |
| `CorporateActionSplit.Symbol` | method | stable | 见接口同名方法 |
| `CorporateActionSplit.Type` | method | stable | `ActionSplit` |
| `RightsIssue` | type | stable | 配股实现（`RightsPer10` + 配股价） |
| `NewRightsIssue` | func | stable | 构造函数 |
| `RightsIssue.Apply` | method | stable | **默认按放弃处理**：持仓与现金都不变（缴款走 `ApplyPaid`，由券商客户端在 pay_date 前提示） |
| `RightsIssue.ApplyPaid` | method | stable | 已缴款：qty 增加 `qty × rightsPer10/10`，均价按 (原成本+新股成本)/新股数 重算，现金增量 = `-cost` |
| `RightsIssue.Description` | method | stable | 见接口同名方法 |
| `RightsIssue.ExDate` | method | stable | 见接口同名方法 |
| `RightsIssue.ExRefPrice` | method | stable | 除权参考价 `(prevClose + 配股价×ratio) / (1+ratio)`；`rightsPer10 ≤ 0` 或 `prevClose ≤ 0` → 原样返回 prevClose |
| `RightsIssue.PayDate` | method | stable | 见接口同名方法 |
| `RightsIssue.RecordDate` | method | stable | 见接口同名方法 |
| `RightsIssue.Symbol` | method | stable | 见接口同名方法 |
| `RightsIssue.Type` | method | stable | `ActionRightsIssue` |
| `Placement` | type | stable | 增发实现（新股数 + 增发价），供稀释提示 |
| `NewPlacement` | func | stable | 构造函数 |
| `Placement.Apply` | method | stable | 持仓与现金都不变（引擎只记录事件） |
| `Placement.Description` | method | stable | 见接口同名方法 |
| `Placement.ExDate` | method | stable | 见接口同名方法 |
| `Placement.PayDate` | method | stable | 见接口同名方法 |
| `Placement.RecordDate` | method | stable | 见接口同名方法 |
| `Placement.Symbol` | method | stable | 见接口同名方法 |
| `Placement.Type` | method | stable | `ActionPlacement` |
| `ActionEngine` | type | stable | 多 action 调度器，**无状态承诺**：已应用历史由调用方持有（生产应落库 / Redis）。⚠️ 它内部有 `appliedLog` 去重表——这是进程内便利，不是持久化真相 |
| `NewActionEngine` | func | stable | 构造空引擎（去重表已初始化） |
| `ActionEngine.ApplyAll` | method | stable | 按「除权日升序 + A 股标准顺序（split → bonus → rights → cash → placement）」应用 `ex_date ≤ asOf` 且未 applied 的 action。返回新持仓、现金增量合计、逐条审计结果。**不修改入参**；无持仓的 action 也标记 applied（附 `SkipReason="no_position"`）；未到除权日 → `ex_date_in_future`；已应用 → `already_applied` |
| `ActionEngine.IsApplied` | method | stable | 该 action（key = `symbol:ex_date:type`）是否已应用 |
| `ActionEngine.MarkApplied` | method | stable | 标记已应用（幂等）。手工回放 / 对账时用 |
| `ApplyResult` | type | stable | 单个 action 对单个持仓的结果（旧/新持仓、现金增量、是否应用、跳过原因） |
| `ApplyOutcome` | type | stable | `ApplyAll` 的逐条审计项（action + `ApplyResult`） |

**因子名册与注册（§6「因子族」；core-only：agent 侧 0 命中，QT-05 ③-⑨ 待回收）**

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
| `FactorKernel` | type | stable | 因子内核类型（四种：leaf / formula / code / pending） |
| `KernelLeaf` | const | stable | `"leaf"`：直接取字段，无公式 |
| `KernelFormula` | const | stable | `"formula"`：公式文本，由 `tellurion:pkg/expression` 求值 |
| `KernelCode` | const | stable | `"code"`：代码引用（`CodeRef`），必须能过编译校验 |
| `KernelPending` | const | stable | `"pending"`：未定型；**不可注册为生效**（`FactorRegistration.Validate` 拒） |
| `FactorKernel.Valid` | method | stable | 是否在四种词表内；零值 → false |
| `FactorCaliber` | type | stable | 因子口径（z-score / rank / raw）——**未知口径 fail-closed，不默认成 z-score**（默认恰是要显式化的隐含假设） |
| `CaliberZScore` | const | stable | `"z-score"`：截面标准化 |
| `CaliberRank` | const | stable | `"rank"`：截面分位 |
| `CaliberRaw` | const | stable | `"raw"`：原始值 |
| `FactorCaliber.Valid` | method | stable | 是否在三种词表内 |
| `FactorDefinition` | type | stable | 因子定义（名 / 内核 / 公式 / 校验 / 代码引用 / 口径 / 输入字段 / 生效日 / 前视规则 / 回看窗口）。`Lookahead` **必填**（蓝图 §5.2 硬要求） |
| `FactorDefinition.Validate` | method | stable | 结构自洽 fail-closed：名字必须在冻结词表、kernel/caliber 合法、`Lookahead` 非空、`EffectiveFrom` 必须是 `YYYY-MM-DD`、非 pending 的 `InputFields` 至少一个且不得重复；leaf/formula 必须有公式且**不得**带 `code_ref`，code 必须带 `code_ref` 且**不得**带公式，pending 的公式 / `code_ref` / `InputFields` 都必须为空。**跨层校验（公式可解析 / 字段在注册表）由 data 包加载时做**——domain 不 import expression |
| `FactorDefinition.VersionHash` | method | stable | canonical JSON（输入字段**排序**后）的 SHA-256 hex。改公式 → 哈希天然变；**方法而非字段**：没有可篡改的存储入口，库里的哈希加载时重新 derive 校验 |
| `FactorHypothesis` | type | stable | 因子假设（可证伪预测 / 依据 / 来源等），供 AI 挖掘链路记录 |
| `BuiltinHypothesisFor` | func | stable | 按因子名取内置假设；未登记 → `(零值, false)` |
| `BuiltinFactorDefinitions` | var | stable | 内置因子定义种子（**内容是数据**：加条目 = MINOR；别当冻结清单照抄） |
| `BuiltinFactorHypotheses` | var | stable | 内置因子假设种子（同上） |
| `FactorSourceLiterature` | const | stable | `"literature"`（有公开文献支撑） |
| `FactorSourceSupplyChain` | const | stable | `"supply_chain"`（产业链逻辑推导） |
| `FactorSourceAI` | const | stable | `"ai_hypothesis"`（AI 提出；**必须过可证伪检验**才可用） |
| `FactorSourceAdHoc` | const | stable | `"ad_hoc"`（临时 / 未记录来源；能用但别太当真） |
| `FactorRegistration` | type | stable | 因子登记（定义 + 版本哈希 + 登记人 + 状态 + 回滚留痕）。四样齐备才算登记：公式 / 哈希 / 登记人 / 生效日 |
| `FactorRegistration.Validate` | method | stable | 四样 + 状态校验：定义合法且**不是 `KernelPending`**、`VersionHash` 必须等于 `Definition.VersionHash()` 派生值、`Registrar` 非空白、状态只能是 active / rolled_back |
| `RegistrationStatus` | type | stable | 登记状态枚举（active / rolled_back） |
| `RegistrationActive` | const | stable | `"active"` |
| `RegistrationRolledBack` | const | stable | `"rolled_back"` |
| `RegistrationStatus.Valid` | method | stable | 是否在词表内；零值 → false |
| `RegistrationAction` | type | stable | 审计动作（register / rollback） |
| `ActionRegister` | const | stable | `"register"`（注意：这是**因子登记审计**的动作，与公司行为的 `Action*` 常量无关） |
| `ActionRollback` | const | stable | `"rollback"` |
| `RegistrationAuditEntry` | type | stable | 登记审计条目（谁 / 何时 / 对哪个版本 / 做了什么 / 为什么） |
| `FactorMiningRecord` | type | stable | 因子挖掘记录（来源 / 生成器 / 代数 / 试验次数与失败数 / 沙箱证明 / 过拟合与 OOS 指标 / 采纳与否） |
| `FactorMiningRecord.Validate` | method | stable | 名字在冻结词表、`VersionHash` 长度必须是 64（sha-256 hex 的长度；**不做字符集校验**）、`Generator` 非空、`SourceKind` 合法、`NumTrials > 0`；**采纳的必须带与版本绑定的沙箱哈希**（`SandboxHash == VersionHash`），**未采纳的必须写拒绝理由** |
| `FactorOutput` | type | stable | 因子产出 schema（「所有因子一个样」；D1 冻结字段形状）：值 + 口径 + 频率 + `VisibleAt`。与 `OHLCV` / `Fundamental` 并列构成规范层三类字段 |
| `FactorCacheEntry` | type | stable | 因子缓存行（原始值 / z-score / 分位 + `Citation`）。⚠️ `Citation` 是 `json.RawMessage`：存储形态是 hash 链、输出形态是五元组，空引用必须是 `'[]'` 而**不是 null**（ODR-061） |
| `FactorReturn` | type | stable | 分层回测结果（五分位平均 / 累计收益 + top-minus-bottom） |
| `ICEntry` | type | stable | IC 检验结果（IC / p 值 / 头部 IC） |

**回测输入输出（§6「回测输入输出」）**

| 符号 | 类别 | 稳定性 | 错误语义 / 判据 |
|---|---|---|---|
| `BacktestResult` | type | stable | 单次回测结果（指标 + 净值曲线 + 成交列表）。指标口径受 `tellurion:docs/SPEC.md` 冻结节约束。⚠️ **它没有也不该有 `Status`**（三值属响应层 `BacktestResponse`）。`UniverseMaxSize` 是**指针**：nil = 未评估/未知，0 = 明确为空票池（OBS-01：把未知编码成 0 会造出假红灯）；`InvalidReasons` 由引擎用 `CheckValidity` 算一次后挂上，消费方**只读不重算** |
| `PortfolioValue` | type | stable | 净值曲线的一个点（日期 / 总市值 / 现金 / 持仓市值） |
| `Trade` | type | stable | 成交记录（方向 / 数量 / 价格 / 三项费用 / 时间戳 / 未成交与已成交数量） |
| `BacktestParams` | type | stable | 回测参数（策略 / 票池 / 区间 / 初始资金 / 无风险利率 / 回看窗）。⚠️ `FillTiming` 与 `Kernel` 是**裸 string**（不能引 contracts 的类型——contracts import domain，反向会成环）：由请求解析点写入**实际生效值**，空值 = **口径未知**而不是「默认收盘成交」；`IndexCode` 保留票池来源，供 `run_spec_hash` 记录 |
| `WalkForwardParams` | type | stable | 滚动窗口参数（训练 / 测试 / 步长 / 最小训练窗） |
| `WalkForwardResult` | type | stable | 单窗口结果（两段回测 + Sharpe + 收益 + 回撤）。⚠️ `OOSvsTrain` 名实不符（REV-20）：当前没有参数优化步骤，它只反映「B 段相对 A 段的衰减」，**不是**过拟合度量；json tag 是历史遗留，改名 = 破坏 API |
| `WalkForwardReport` | type | stable | 多窗口汇总（均值 / 标准差 / 通过率 / 总体是否通过）。⚠️ `OverfitScore` / `ProbNoOverfit` 同样是历史名（REV-20），当前语义是分段衰减的分档映射 |

**core-only 数据族（QT-05 ③-⑨ 回收清单；agent 侧 0 命中 ⇒ 标 experimental，下游不要依赖）**

| 符号 | 类别 | 稳定性 | 错误语义 / 判据 |
|---|---|---|---|
| `ScreenFilters` | type | experimental | 多因子筛选条件（PE / PB / PS / ROE / ROA / 负债率 / 毛利率 / 净利率 / 市值，全部**指针**：nil = 不限，不是 0） |
| `ScreenRequest` | type | experimental | 筛选请求（条件 + 日期 + 上限） |
| `ScreenResult` | type | experimental | 筛选命中行（指标全指针，语义同 `ScreenFilters`） |
| `Sector` | type | experimental | 板块行情行（代码 / 名称 / 分类 / 涨跌幅 / 领涨股 / 来源） |
| `StockSector` | type | experimental | 股票 ↔ 板块归属 |
| `StockIndustry` | type | experimental | 股票 ↔ 行业归属（**无 json tag**：字段名即形状，改字段名 = 破坏兼容） |
| `SWIndustryNode` | type | experimental | 申万行业树节点（代码 / 父节点 / 层级 / 是否叶子 / 来源 / as-of） |
| `RelationIndustryMember` | const | experimental | `"industry_member"`（产业链关系字面量） |
| `SupplyChainEdge` | type | experimental | 产业链边（上游 → 下游 / 关系 / 权重 / 来源 / as-of） |
| `IndustryVitality` | type | experimental | 行业景气度指标行（指标 / 观测日 / 值 / 单位 / 来源 / as-of） |
| `ETFConstituent` | type | experimental | ETF 成分与权重（ETF 代码 / 成分 / 权重 / 观测日 / 来源） |
| `MacroSeries` | type | experimental | 宏观序列行（指标 / 频率 / 观测日 / 值 / 单位 / 来源 / as-of） |
| `MarketActivity` | type | experimental | 市场活跃度行（换手率 / 自由流通换手率 / 量比，全指针：nil = 上游未给） |
| `FundamentalsDetailRow` | type | experimental | 深财务明细行（PIT 键是 `AnnDate`；`Value` 指针：nil = 源侧缺失 ≠ 0；`SnapshotURI` 回指归档原始响应，`RawFieldName` 保留源字段名以便追溯） |
| `StrategyConfig` | type | experimental | 策略配置的**存储行**（定义见 §8 #2：ADR-0001 ③ 禁 `Config` 命名，本类型按豁免保留原名——core 4 个文件在用；③-⑨ 待回收，**最可能的家是 `tellurion-strategy`**，不是本仓） |
| `ExecutionConfig` | type | experimental | 执行配置（订单类型 / 滑点模型 / 费率 / 冲击模型参数）。§8 #2：名字按豁免保留（core 21 个文件在用）；③-⑨ 待回收，**执行留在 core**，所以它大概率也该回 core |
| `DefaultExecutionConfig` | func | experimental | 默认执行配置；费率字段取自 `fees.Default*`（单一来源，S7-P1-4 修过 0.00025 与 0.0003 漂移） |

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
| `BacktestRequest` | type | **frozen** | `RunSpec` 的承载物（§5.1）：策略 / 票池 / 区间 / 资金 / 频率 / 内核 / 引擎参数 / 信号行 / 成交口径。未知名拒绝（`ValidateRunSpec`） |
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
| `ValidateParamsWrite` | func | stable | 写 params 前的边界校验：`paramsLen < 0` 或 `paramsLen > min(MaxParamsBytes, scratch_ptr−params_ptr)` → `ErrParamsTooLarge`。**截断是明确禁止的**（截断等于让策略读到半个 JSON）；阈值本身是 frozen 常量 |
| `ErrAbiVersionUnknown` | var | stable | 见 §5.4（abi / artifact 校验族） |
| `ErrAbiVersionMismatch` | var | stable | 见 §5.4 |
| `ErrAbiFamilyMissing` | var | stable | 见 §5.4 |
| `ErrAbiFamilyMixed` | var | stable | 见 §5.4 |
| `ErrImportNotWhitelisted` | var | stable | 见 §5.4 |
| `ErrImportSignatureMismatch` | var | stable | 见 §5.4 |
| `ErrHostRegionMissing` | var | stable | 见 §5.4 |
| `ErrHostRegionInvalid` | var | stable | 见 §5.4 |
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
| 内核路由 | `ErrKernelNotRegistered`、`ErrMarketNotRegistered` | stable（✅ 已在 `engine/`） | 请求期 400；异步 job 终态 `failed`。**零静默回退** |
| 成交口径 | `ErrFillTimingMissing`、`ErrFillTimingUnknown`、`ErrFillTimingMismatch` | stable | 引擎永远不猜口径；口径不符是错误不是结果差异 |
| abi / artifact 校验 | `ErrAbiVersionUnknown`、`ErrAbiVersionMismatch`、`ErrAbiFamilyMissing`、`ErrAbiFamilyMixed`、`ErrImportNotWhitelisted`、`ErrImportSignatureMismatch`、`ErrHostRegionMissing`、`ErrHostRegionInvalid`、`ErrParamsTooLarge`、`ErrScopeUnknown`、`ErrScopeMarketV1` | stable（✅ 已在 `engine/`；`errors/` 不独立成包见 §8 #1） | 归 `ErrContractBreach` 语义：装载期拦下，不进 run |
| 信号校验 | `ErrActionIllegal`、`ErrStrengthOutOfRange`、`ErrTargetWeightOutOfRange`、`ErrEmitSourcesMixed`、`ErrEmitSourceUnknown`、`ErrSignalBudgetExceeded`、`ErrStateBudgetExceeded` | stable | 归 `ErrStrategyFault`（非法信号 / 预算超限 = 策略故障，终止 run） |
| provenance / hash | `ErrCanonicalUnsupported`、`ErrStrategyRefMalformed` | stable | 输入身份不可 canonical 编码 → 当场失败，禁止退化成「猜一个 hash」 |

---

## 6. 值对象面（`domain/`，S1 整包迁入——✅ 2026-10-22 已落）

来源与再生成：**本仓内**跑 `go run ./tools/check_api_surface -list`（今天输出全部 429 个；本节覆盖其中 `domain` 一侧的 235 个 = `fees` 13 + `domain/market` 46 + `domain` 176，根包 0，`engine` 的 194 个见 §4/§5）。落地前从 core 抓的快照是 2026-10-21 的 227 个（`market/` 46 + 根包 181）：`market/` 46 个与根包 176 个**已全部登记在 §4**；差的 5 个是 §8 #2 裁决删掉的零引用部署配置（`Config` / `DatabaseConfig` / `RedisConfig` / `ServiceConfig` / `TushareConfig`，见 [`adr/0002`](adr/0002-config-族的边界处置.md)）。
**逐符号清单不在这里手抄**——抄一份就会漂；落码时用上面这条命令生成 §4 的表，再按下面分组填稳定性。

| 组 | 代表符号 | 稳定性 | 备注 |
|---|---|---|---|
| 回测输入输出 | `BacktestResult`、`BacktestParams`、`PortfolioValue`、`Trade`、`RiskMetrics`、`WalkForward*` | stable（`BacktestResult` 的**指标口径**受 SPEC 冻结节约束） | `BacktestResult` **没有也不该有 `Status`**——三值属响应层，见 `tellurion-agent:docs/TASKS.md` OBS-01-R2-A 的 core 裁决 (b) |
| 仪器域行情 | `OHLCV`、`Stock`、`Fundamental(Data)`、`IndexDaily`、`IndexConstituent`、`Dividend`、`Split`（多为 `market/` 的类型别名） | stable | 别名不是重复真相：实体在 `domain/market`，根包别名保兼容 |
| 时间/频率语义 | `Frequency`、`FrequencySpec`、`TimeSemantics`、`AsiaShanghai`、`Visible*`、`DailyVisibleAt` | **frozen** | PIT 红线：`tellurion:docs/SPEC.md` §「时间语义（PIT 红线）」。⚠️ 说清来源——core「冻结 schema 变更评审」点名的名单是三契约（RunSpec/ResultSpec/SignalSpec），这一行的 frozen 是**本文件按 §2 定义给的标注**，依据是它已有注册期闸门（D2）与逐源钉子（`tellurion:pkg/data/source/registry_rev05_test.go` 等把每个 connector 的 `TimeSemantics` 钉死）；改语义同样要先让那些测试红 |
| 信号与订单 | `Signal`、`Direction`、`Order`、`OrderType`、`Position(Size)`、`Portfolio`、`TargetPosition`、`StopLossEvent`、`MarketRegime` | stable | — |
| 因子族 | `FactorDefinition`、`FactorHypothesis`、`FactorType`、`FactorKernel`、`FactorCaliber`、`FactorRegistration`、`FactorMiningRecord`、`FactorOutput`、`ICEntry`、`FactorValue` + 内置定义/假设种子变量 | stable（`VersionHash` 语义受「版本哈希纪律」节约束） | 封闭词表扩展 = 加名 + 加实现 + 同步 SPEC 对应节 |
| 公司行为 | `CorporateAction`、`CashDividend`、`BonusShare`、`CorporateActionSplit`、`RightsIssue`、`Placement`、`ActionEngine`、`ApplyResult`/`ApplyOutcome` | stable | — |
| 配置形状 | ~~`Config`、`ServiceConfig`、`DatabaseConfig`、`RedisConfig`、`TushareConfig`~~（**已按 §8 #2 裁决删除：零引用 + 含 `Password`**）、`StrategyConfig`、`ExecutionConfig` | ⚠️ **experimental（core-only，③-⑨ 待回收）** | 命中 `adr/0001` ③ 的 `Config` 命名禁令；裁决见 **§8 #2 + [`adr/0002`](adr/0002-config-族的边界处置.md)**：前五个不迁（core 侧随 sed 删除），后两个按原名保留、进回收清单。`engine/` 的 `TradingConfig` / `PriceLimitConfig` 另按「冻结面伴生名」保留（同一条 ADR ③） |
| core-only（agent 侧 0 命中） | QT-05 ③-⑨ 点名的 5 族：`Screen*`、`Sector`/`StockSector`/`SWIndustryNode`、`SupplyChainEdge`、`MacroSeries`、`mining` 族（`FactorMiningRecord`）；实测另含 `CorporateAction` 族、`FactorDefinition`/`FactorRegistration`/`FactorOutput`、`RiskManager`、`ExecutionConfig`/`StrategyConfig`、`Signal`/`Order`、`FundamentalsDetailRow`、`ETFConstituent`、`IndustryVitality`、`MarketActivity` 等 | **experimental + 待回收** | QT-05 ③-⑨ 与 ADR-0001 的代价条款：S1 仍整包迁，S2 后按清单回收。**下游不要依赖这一组**。<br>**实测口径（2026-10-22，落码时重测，取代 2026-10-21 的 86/124）**：`domain/` 根包 176 个导出符号折成 **119 个根名**（方法并回所属类型），agent 侧 import 了其中 **38 个**，**81 个 0 命中**。再生成命令：`grep -rhoE '\bdomain\.[A-Za-z0-9_]+' --include='*.go' <tellurion-agent 仓>` 去重，与 `go run ./tools/check_api_surface -list` 的 `domain/` 行对差。<br>agent 实际在用的 38 个是：`BacktestResult`、`Direction*`、`Factor*` 名册与假设、`MarketRegime`、`OHLCV`、`OrderType*`、`ParseFactorType`、`Portfolio`、`PortfolioValue`、`Stock`、`Trade`、`WalkForward*`。<br>「0 命中」≠「该删」——agent 走 HTTP/JSON 消费的那几族（`Signal`/`Order`）本来就不摸 Go 类型。**完整逐符号清单以命令输出为准**（本表只记判据、族与数字，不手抄会漂的名单） |

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
| 4 | `dataprotocol` 是否独立成包（QT-40 对表 T2，2026-10-22 登记） | **落点已定名**：`contracts/dataprotocol`（包名 `dataprotocol`）——不叫 `data/`（包名要能自证「这是协议不是实现」，与 #1 的 `errors/`/`abi/` 之争同型），不叫 `dataset`（那是 X1b 的 aggregator 面，形状还没定）。依据 `tellurion:docs/design/qt-decoupling/data-service-contract-and-ownership.md` §4.1 A 组注 ①。**要裁的是迁入时机，而且两份上游文档给的判据并不一致**：「只有单个仓使用的类型不许提前迁」对应的是 ①「本仓只放被 ≥2 个仓使用的类型」（[AGENTS.md](../AGENTS.md) §3 把它写成「只有单个仓使用的类型属泄漏」）；对表注 ② 与本仓下游的 `tellurion-data:docs/api-schema.md` §4 **原先都把它引成 ③**（③ 管的是 `Config`/`Options`/`Impl` 命名），两处引用已于 2026-10-22 同步订正。**但同一份 ADR 的「代价」节明写 S1 整包迁时接受 core-only 类型临时进来、输出一份清单、S2 后回收**，而本仓台账已把这条落成硬要求（`docs/TASKS.md` QT-05 ③-⑨：`screen`/`sector`/`supply_chain`/`macro`/`mining` 在 agent 侧各 0 命中，仍整包迁）。于是有两条都自洽的走法：(a) `dataprotocol` 不适用那条容忍（它是**名单外的新建包**，不是整包迁的附带品）⇒ 挂 X1b 第一 commit，`tellurion-data` 一 import，两个使用方立刻成立；(b) 把它当同类容忍处理，随 S1 一起落、进那份 core-only 清单（③-⑨ 那份名单今天不含它，选 (b) 就要把它加进去）。**X1a 待办已登记（2026-10-22）：X1a 开工前必须挑一条并写回本节（#4）+ `tellurion:docs/design/qt-decoupling/data-service-contract-and-ownership.md` §4.1 注 ②**，别两份文档各按一条（现在的写法是注 ② 选 (a)、而它引的条款其实同时允许 (b)）；X1a 的施工序在 `tellurion-data:docs/TASKS.md` 施工序表。冻结前要先定形状的三处实测：① `UnifiedDataPoint` 今天有**两份同名同字段副本**（`tellurion:pkg/data/source/unified.go:16` 与 `pkg/storage/bulk_insert.go:17`，后者注释自述「storage 不能 import source（会成环）」）⇒ 进 contracts 后成环理由消失，**必须删副本改引用**，否则协议加一列时 `BulkInsert` 的列布局会静默少写一列（`bulk_insert_ddl_test.go` 只核「目标表在 DDL 里」，核不到字段）；② **没有 `DataType` 类型**——`ConnectorDeclaration.DataType` 与 `UnifiedDataPoint.DataType` 都是裸 `string`，取值集靠 `DataType*` 常量；③ **没有 `Aggregator` 接口**（§5.2 的 port 表里只有未定形的 `DataPort`，标着「⬜ 设计中，S1 不落」）：aggregator 今天只有 HTTP 形状，要么随 X1b 实现定、要么先不落 port，**不许在文档里造一个无实现的接口**（§8 已禁默认实现/mock 进本仓）。另：`TimeSemantics` 枚举定义在 `tellurion:pkg/domain/market/visibility.go:61`，**已随 QT-05 ③-① 整包迁入 `domain/market/`（2026-10-22，见 §4）**，`dataprotocol` 里**不再重复定义**——同一次迁移出现两份 `TimeSemantics`，等于把 #1 与 #4 合成一个真事故 |

---

## 9. 与其他文档的互锁

| 文档 | 关系 |
|---|---|
| `api-schema.md` | 规则源（错误分类 / 幂等 / 兼容 / 不提供什么）；§1.1 的 port「状态」列与本文件 §5.2 一致，§1.3 的三分法与 §1.4 的 abi 计数都由本文件 §5.4/§7 逐条回源码核对（2026-10-21 自审时改掉过那两处的「写得像事实」：哨兵零定义、host 函数 13→12） |
| `ARCHITECTURE.md` §2 | 包布局与本文件 §8 待裁决 #1 有关 |
| `tellurion:docs/SPEC.md` | **冻结文本的唯一源**；本文件只指节标题。改冻结面 = 先破坏 core 合规测试 → 更新 SPEC → 变更评审 → 才动代码（`tellurion:docs/TASKS.md`「维护本文件」硬规则） |
| `tellurion:docs/design/qt-decoupling/kernel-interface-map.md` | 现状接口面清单的来源（M6 的「散在几处」即指它）；本文件不是它的替代，是**随仓可核对**的那一份 |
