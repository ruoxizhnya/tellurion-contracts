---
status: active
last-verified: 2026-10-24
verified-by: QT-20 / QT-41b ④ 契约腿（§1.2 增 `StrategyRef`/`StrategyParams`/`BacktestParams.Params` 的只增不改说明；§1.3 增编译 / 资产四哨兵族与「私有哨兵不迁」判据）；QT-05 ③-① 落码后订正（§1.1 port 状态与承载路径、§1.2 「只搬不改名」已落地；符号本体住本仓 `engine/`）
---

# API Schema — tellurion-contracts

> **契约地位**：本仓的"API"是 **Go 接口契约**，不是 HTTP。消费方式是 `require` + 版本号，禁止源码路径依赖（`replace`）。
> 本文件是**本仓对外接口的唯一真相源**。代码与它不一致时，先判断谁错：契约错 → 改本文 + 立 ADR；实现错 → 改代码。
> 跨仓消费方（core / agent / ui / 其他服务）按本文件编码，**不读本仓实现**。

## 1. 接口面

### 1.1 port（能力边界）

> **"状态"这一列是必须的**：`EngineRunner` / `SignalProvider` 今天在 core 里是**真实 interface**；
> 而 `DataPort` / `BrokerGateway` / `ExecutionPort` 只是**设计中的名字**，源码里 0 命中——不能当成已存在的接口去 import。

| 接口 | 语义 | 实现方 | 状态 |
|---|---|---|---|
| `EngineRunner` | 跑一次回测 | core（实现）；接口本体已迁本仓 `engine/` | ✅ **已存在** |
| `SignalProvider` | 逐 bar 出信号（v1 per-symbol / v2 混合 scope） | core（`pkg/strategy/signal_provider.go:41`） | ✅ **已存在** |
| `BrokerGateway` | 券商通道：下单 / 撤单 / 查单 / 持仓 / 成交 | `tellurion-broker` | ⬜ 设计中（**S1 未落**：S1 只搬既有面；形态见 `broker-gateway-design.md` §3） |
| `DataPort` | 数据集解析与取数 | `tellurion-data` | ⬜ 设计中（**S1 未落**；归 X1） |
| `ExecutionPort` | 撮合与成交回传 | core | ⬜ 设计中（**S1 未落**；执行留 core） |

> **命名冲突已统一**：D29 §3.3 把券商通道写成 `BrokerPort`，而 `broker-gateway-design.md` 用 `BrokerGateway`。
> **以 `BrokerGateway` 为准**——设计文档里出现 `BrokerPort` 时按同一物的旧写法理解。

### 1.2 四契约（ADR-031）

`RunSpec` · `SignalSpec` · `ResultSpec` · `EngineRunner` —— 回测的**输入 / 信号 / 输出 / 执行者**。

> ⚠️ **这四个是设计术语，不是四个同名 Go 类型。** 承载物是：
>
> | 契约术语 | 承载 | 依据 |
> |---|---|---|
> | `RunSpec` | `contracts.BacktestRequest`（含 `ValidateRunSpec()`；未知名拒绝） | 本仓 `engine/contracts.go`（core 旧路径 `tellurion:pkg/backtest/contracts/contracts.go`） |
> | `ResultSpec` | `contracts.BacktestResponse`（`status` 三值冻结） | 本仓 `engine/contracts.go` + core 的 `resultspec_test.go`（SPEC 互锁，刻意留 core） |
> | `SignalSpec` | 信号形状（`SignalWeight` 等）+ 信号提供者契约 | 本仓 `engine/contracts.go` |
> | `EngineRunner` | `contracts.EngineRunner`（interface） | 本仓 `engine/contracts.go` |
>
> **S1 只搬不改名（已落地，2026-10-22）**：W0 已结束，承载类型未改名——`BacktestRequest` / `BacktestResponse` 就是 `RunSpec` / `ResultSpec` 的 Go 符号，符号本体今天住本仓 `engine/`（core 侧旧包保留到 sed 落地）。
> 看到本仓里没有 `type RunSpec struct` 不是缺漏——这是当前刻意的状态；今后改名属破坏性变更（MAJOR）。
>
> **QT-20 增量（2026-10-24，只增不改）**：`BacktestRequest` 增两个**可选**字段——
> `StrategyRef`（wire 形态 = `ParseStrategyRef` / `StrategyRef.String()` 的 `id` / `id@hash`，
> **不新增第二种引用形状**）与 `StrategyParams`（**浅层对象**，值为标量；嵌套键在
> `ValidateRunSpec` 拒绝）。`domain.BacktestParams` 增 `Params`（**生效参数**，字段名对齐
> core `RunMeta.Params`）。**未知键清单不在这里**：它属于 core 的策略 manifest（J2），
> 契约层只做形状 / 类型校验；合并（defaults ⊕ 覆写）与取值域归执行侧。

### 1.3 哨兵错误（QT-10 的三分法 + QT-41b 编译族）

| 错误 | 含义 | 消费方行为 |
|---|---|---|
| `ErrDataGap` | 数据缺口 | 跳过当天 + 写 `invalid_reasons`，run 继续 |
| `ErrStrategyFault` | 策略故障（trap / 非法信号 / 预算超限） | **终止 run**，`status=failed` |
| `ErrContractBreach` | 契约违规（前视 / 未知参数 / abi 不符） | **终止 run**，`status=failed` |

> ✅ **三条哨兵已落本仓**（2026-10-23，QT-10 ③-①）：`engine/execution_errors.go`，符号级现状见 `API.md` §5.4。
> ⚠️ **但归类（③-②）与状态映射（③-③）在 core**：provider 侧把底层错误包成哪一类、`RunBacktest` 在哪条分支终止 run，
> 都是 core 的改动，今天还没做——所以 §2 前两条对生产者是**验收标准**，对消费方（`errors.Is` + 兜底类）才是现状。
> 与本节含义配对的新稳定 Code：`CodeDataGapSkipped`（`data_gap_skipped`）走既有 `invalid_reasons` 通道，**不是第四种状态**。

**编译 / 资产族（QT-41b ④，✅ 2026-10-24 已落本仓）**：

| 错误 | 含义 | 消费方行为 |
|---|---|---|
| `ErrCompileFailed` | 源码编译失败（诊断可检索） | 反馈给生成方 / 模型；不可当「服务挂了」重试 |
| `ErrCompileUnavailable` | 资产服务未配置编译端口 | 配置错 / 可重试；与上一条分开，否则模型会对着坏配置反复改源码 |
| `ErrSourceTooLarge` | 上传源码超过 1 MiB | 不可重试，改输入 |
| `ErrArtifactTooLarge` | artifact 超过 8 MiB blob 上限 | 不可重试（fail-loud，不截断） |

> 迁入判据 = **跨仓消费方必须能识别**（agent 生成闭环 / core 取 artifact / orchestration 终态）。
> 资产服务的**私有哨兵不迁**：`ErrStrategyNotFound` / `ErrStrategyExists` / `ErrVersionNotFound` /
> `ErrInvalidTransition` / `ErrJobNotFound` 等描述的是服务自己的存储与状态机，跨仓调用方只看
> HTTP status + wire code。四条名字与语义逐字取自 `tellurion-strategy:internal/asset/errors.go`，
> 该仓升 pin 后改为本包别名（过渡态是同一语义的两处定义，不是两份真相）。

### 1.4 abi 常量

host API 的 abi 版本、内存区所有权、枚举上限。冻结集合是 **12 个 host 函数名**（v1 家族 8 个原样保留 + v2 新增 4 个：
`get_bar2` / `get_series_len2` / `get_date2` / `emit_target2`），逐条见 `API.md` §7。

> ⚠️ **本节原写「abi v2 = 13 个 host 函数」，是滞后于 core 的错数**：那个 13 来自设计文档的一处笔误，
> core 在 2026-10-21（QT-01 落代码时）已订正为 12，并写明代码侧判据是 `contracts.HostFnCount()`
> （返回枚举结果，不返回文档字面数），由 `abi_contract_table_test.go` 钉住。
> 依据：`tellurion:docs/design/qt-decoupling/adr-029-amendment-host-api-v2.md` §2.1 的 2026-10-21 订正条。
> 真要补第 13 个（如那条可选的 `is_tradable`）：先改那张表、走一次 abi 冻结，数字才跟着变。

## 2. 错误码与 fail-closed 语义

- **错误分类是契约的一部分**：provider 侧必须把底层错误包成上面三类之一，**禁止裸 `fmt.Errorf` 透出**（否则调用方无法区分"跳过"与"终止"）。
- fail-closed 默认：无法归类 → 按 `ErrContractBreach` 处理（宁可停，不要跑出漂亮但错的曲线）。

## 3. 幂等性与重试

内容寻址三兄弟（改了内容必改 hash，因此天然幂等）：

| 标识 | 覆盖什么 |
|---|---|
| `run_spec_hash` | 策略 + 参数 + 区间 + 票池 + 引擎配置 |
| `dataset_snapshot_hash` | 票池 + 窗口 + 频率 + 数据版本 |
| `artifact_hash` | 策略字节码本体 |

## 4. 版本与兼容

- **0.x 阶段**（现在 → 1.0）：允许 breaking，但必须 **bump MINOR + 写 CHANGELOG 条目**。
- **1.0 之后**：`deprecated` 至少保留 1 个 MINOR 周期（QT-19 定规则，QT-09 ⑤ 机械检查）。
- 下游只认 tag / 伪版本；**任何 `replace` 指向仓外路径都视为破坏单仓可构建**。

## 5. 不提供什么（边界）

- 不提供实现（没有"默认实现"这种东西，那会让下游偷懒）
- 不提供业务规则（风控阈值、策略逻辑都不在这）
- 不提供 DB 访问、不提供 HTTP handler
- 不提供测试替身（mock 由使用方自己写，避免契约仓被测试代码污染）
