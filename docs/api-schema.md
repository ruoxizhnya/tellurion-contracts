---
status: active
last-verified: 2026-10-20
verified-by: 功能边界固化（三份契约：db / api / dataflow）
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
| `EngineRunner` | 跑一次回测 | core（`pkg/backtest/contracts/contracts.go:366`） | ✅ **已存在** |
| `SignalProvider` | 逐 bar 出信号（v1 per-symbol / v2 混合 scope） | core（`pkg/strategy/signal_provider.go:41`） | ✅ **已存在** |
| `BrokerGateway` | 券商通道：下单 / 撤单 / 查单 / 持仓 / 成交 | `tellurion-broker` | ⬜ 设计中（S1 落 port；形态见 `broker-gateway-design.md` §3） |
| `DataPort` | 数据集解析与取数 | `tellurion-data` | ⬜ 设计中（S1 落 port） |
| `ExecutionPort` | 撮合与成交回传 | core | ⬜ 设计中（S1 落 port） |

> **命名冲突已统一**：D29 §3.2 把券商通道写成 `BrokerPort`，而 `broker-gateway-design.md` 用 `BrokerGateway`。
> **以 `BrokerGateway` 为准**——设计文档里出现 `BrokerPort` 时按同一物的旧写法理解。

### 1.2 四契约（ADR-031）

`RunSpec` · `SignalSpec` · `ResultSpec` · `EngineRunner` —— 回测的**输入 / 信号 / 输出 / 执行者**。

> ⚠️ **这四个是设计术语，不是四个同名 Go 类型。** 今天 core 里的承载物是：
>
> | 契约术语 | 今天的承载 | 依据 |
> |---|---|---|
> | `RunSpec` | `contracts.BacktestRequest`（含 `ValidateRunSpec()`；未知名拒绝） | `pkg/backtest/contracts/contracts.go:171` |
> | `ResultSpec` | `contracts.BacktestResponse`（`status` 三值冻结） | 同上 `:202` + `resultspec_test.go` |
> | `SignalSpec` | 信号形状（`SignalWeight` 等）+ 信号提供者契约 | `contracts.go:160` |
> | `EngineRunner` | `contracts.EngineRunner`（interface） | `contracts.go:366` |
>
> **S1 只搬不改名**：是否把承载类型改名为 `RunSpec` / `ResultSpec`，是 **W0（QT-01~04）** 的决定。
> 看到本仓里没有 `type RunSpec struct` 不是缺漏——这是当前刻意的状态。

### 1.3 哨兵错误（QT-10 的三分法）

| 错误 | 含义 | 消费方行为 |
|---|---|---|
| `ErrDataGap` | 数据缺口 | 跳过当天 + 写 `invalid_reasons`，run 继续 |
| `ErrStrategyFault` | 策略故障（trap / 非法信号 / 预算超限） | **终止 run**，`status=failed` |
| `ErrContractBreach` | 契约违规（前视 / 未知参数 / abi 不符） | **终止 run**，`status=failed` |

### 1.4 abi 常量

host API 的 abi 版本、内存区所有权、枚举上限（abi v2 = 13 个 host 函数，见 QT-01 与 `adr-029-amendment-host-api-v2.md`）。

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
