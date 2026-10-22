---
status: active
last-verified: 2026-10-22
verified-by: QT-05 ③-① 落码（§2 订正「只依赖 kit」为实测零依赖 + 标注 errors//abi/ 不独立成包 + `domain/market/` 已落；§3 加 `domain/market/` 行、fees 行按实测收紧；§5 依赖图同步订正；§8 S1 标 🔶 并写明 push 阻塞与双份真相窗口）
---

# 架构 — tellurion-contracts

## 1. 定位

跨仓唯一缝隙：**port 接口 · DTO/值对象（含 `domain`）· 四个引擎契约 · 哨兵错误 · abi 常量**。

它的存在意义是让"改契约"成为一次**显式发布**，而不是某天某人顺手改了 core 里的一个 struct。

## 2. 分层

```
（被所有 Go 仓依赖）

  domain/     值对象：OHLCV / Trade / BacktestResult / Factor* / Direction
  domain/market/  行情值对象 + 频率/时间语义词表 + Provider（✅ 2026-10-22 已落地，46 个导出符号）
  engine/     四契约：RunSpec / SignalSpec / ResultSpec / EngineRunner + MarketRouter
  fees/       费率 / 滑点 / 成本模型的**定义**（算钱留在 core）
  errors/     三类哨兵          ⚠️ 已裁决**不独立成包**（随 engine/），图待 engine 落地时改
  abi/        host 函数版本…    ⚠️ 同上（依据见 docs/API.md §8 待裁 #1）
```

> **依赖 kit 那句已订正（2026-10-22，QT-05 落码时实测）**：本节原先写「自己只依赖 tellurion-kit 的 decimal / errors」。
> 对迁移输入的四个 core 包实测（`tellurion:pkg/fees`、`pkg/domain`、`pkg/domain/market`、`pkg/backtest/contracts` —— 即本仓的 `fees` / `domain` / `domain/market` / `engine`）：
> 非标准库 import **只有 `github.com/stretchr/testify`（仅测试）**，仓内 import **只在这四个包彼此之间**，kit 符号 **0 命中**。
> 证据命令：`go list -f '{{join .Imports "\n"}}' ./pkg/{fees,domain,backtest/contracts}/...`（另取 `.TestImports`），过滤标准库后逐条看过。
> 所以本仓 `go.mod` 落码后**不 require kit**；真要用（例如把 `float64` 换成 `kit/decimal`）是一次显式的 MINOR + 一条 ADR，不是默认前提。
> `AGENTS.md` §4 那张表已按此订正（testify 已登记为**仅测试**依赖）。

## 3. 组件职责

| 包 | 职责 | 不该做什么 |
|---|---|---|
| `domain/` | 跨边界值对象 | 不做计算、不读配置 |
| `domain/market/` | 行情值对象 + 频率/时间语义词表 + `Provider` 接口 | 不做计算；`OHLCV`/`Fundamental` 的字段集与 core `docs/SPEC.md` 互锁（加字段 = core 侧红），见 `docs/API.md` §4 |
| `engine/` | 回测四契约 + 路由接口 | 不含引擎实现 |
| `fees/` | 成本模型的形状 | 不算钱（成交额 × 费率、最低佣金兜底在 core 的 tracker / mock_trader；本包实测零算术，见 `docs/API.md` §8 #3 的裁决依据） |
| `errors/` | 三类哨兵 + 包装约定 | 不做日志 |

## 4. 三份契约（功能的边界就在这三份文件里）

| 维度 | 文件 | 回答什么问题 |
|---|---|---|
| 数据形状 | [`db-schema.md`](db-schema.md) | 本仓拥有哪些表？谁能写？ |
| 接口形状 | [`api-schema.md`](api-schema.md) | 对外暴露什么？错误码与幂等怎么做？ |
| 流向与顺序 | [`dataflow.md`](dataflow.md) | 数据怎么流？哪里 fail-closed？ |

## 5. 依赖方向

```
（今天无内部依赖：本仓不 require tellurion-kit，见 §2 的实测订正；
  未来若要用 kit/decimal，是一次显式 MINOR + 一条 ADR，不是默认前提）

tellurion-contracts
        ↑            ↑            ↑
      core        agent        各服务
```

## 6. 关键决策

- [`adr/0001-契约只增不破且零实现.md`](adr/0001-契约只增不破且零实现.md)

## 7. 非目标（明确不做）

- 不做实现、不做默认实现
- 不做 HTTP / gRPC 服务（跨进程协议是 OpenAPI，归 core；本仓只管进程内形状）
- 不持有任何表
- 不做策略执行（那是 core 的 WASM host）

## 8. 分期与状态

| 阶段 | 内容 | 状态 |
|---|---|---|
| S1（QT-05） | `pkg/domain` + `pkg/backtest/contracts` 搬入本仓，core/agent 改 `require` 伪版本 | 🔶 **进行中（2026-10-22）**：本仓侧 `fees/` ✅（13 符号）、`domain/market/` ✅（46 符号）已落并登记进 `docs/API.md` §4，护栏 59/59 绿；余下 `domain` → `engine` 按「一次只搬一个包」依次落。**core/agent 切依赖那条腿阻塞在 push**：`go get <本仓>@<sha>` 只能从远端解析，而 ③-④ 禁止 `replace`——本地提交不推上去，下游就没法 `require`。窗口期内 core 的 `pkg/domain` 与本仓的 `domain/` **并存（双份真相）**，两边都不许改，直到 core 侧 sed 落地并删掉旧包 |
| S2（QT-07） | 打 `v0.1.0`，下游换 tag | ⬜ |
| QT-19 | 版本规则表 + 兼容矩阵 + CHANGELOG | ✅ 已随 CT-02 落地（2026-10-22，护栏 `tools/check_release_matrix` 进 CI） |
