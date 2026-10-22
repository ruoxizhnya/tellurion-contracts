---
status: active
last-verified: 2026-10-22
verified-by: QT-05 ③-① 落码完成本仓侧（§2 订正「只依赖 kit」为实测零依赖 + 四个包均标已落 + errors//abi/ 按 adr/0003 并入 engine；§3 四行按实测收紧；§5 依赖图同步订正；§8 S1 标本仓侧已完成、下游腿阻塞在 push）
---

# 架构 — tellurion-contracts

## 1. 定位

跨仓唯一缝隙：**port 接口 · DTO/值对象（含 `domain`）· 四个引擎契约 · 哨兵错误 · abi 常量**。

它的存在意义是让"改契约"成为一次**显式发布**，而不是某天某人顺手改了 core 里的一个 struct。

## 2. 分层

```
（被所有 Go 仓依赖）

  domain/     值对象：Signal / Order / Position / Portfolio / 公司行为 / 因子名册 / BacktestResult（✅ 2026-10-22 已落地，176 个导出符号）
  domain/market/  行情值对象 + 频率/时间语义词表 + Provider（✅ 2026-10-22 已落地，46 个导出符号）
  engine/     四契约：RunSpec / SignalSpec / ResultSpec / EngineRunner + MarketRouter（✅ 2026-10-22 已落地，194 个导出符号）
              + 三类哨兵与 abi 常量：**已裁决不独立成包**，随本包走（[`adr/0003`](adr/0003-哨兵与-abi-不独立成包.md)）
  fees/       费率 / 滑点 / 成本模型的**定义**（算钱留在 core）
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
| `domain/` | 跨边界值对象（信号 / 订单 / 组合 / 公司行为 / 因子名册 / 回测结果）+ `domain/market` 的兼容别名 | 不做计算、不读配置；**部署配置形状不进来**（[`adr/0002`](adr/0002-config-族的边界处置.md)：零引用的 5 个 `*Config` 已删，core-only 的 2 个标 experimental 待回收） |
| `domain/market/` | 行情值对象 + 频率/时间语义词表 + `Provider` 接口 | 不做计算；`OHLCV`/`Fundamental` 的字段集与 core `docs/SPEC.md` 互锁（加字段 = core 侧红），见 `docs/API.md` §4 |
| `engine/` | 回测四契约 + 路由接口 + **三类哨兵 + abi 常量**（同一包、`package contracts`） | 不含引擎实现；不把哨兵 / abi 拆成独立包（[`adr/0003`](adr/0003-哨兵与-abi-不独立成包.md)：拆了没有调用点能少 import 一个包） |
| `fees/` | 成本模型的形状 | 不算钱（成交额 × 费率、最低佣金兜底在 core 的 tracker / mock_trader；本包实测零算术，见 `docs/API.md` §8 #3 的裁决依据） |

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
| S1（QT-05） | `pkg/domain` + `pkg/backtest/contracts` 搬入本仓，core/agent 改 `require` 伪版本 | 🔶 **本仓侧已完成（2026-10-22）**：`fees/` ✅（13）、`domain/market/` ✅（46）、`domain/` ✅（176，5 个零引用部署配置按 [`adr/0002`](adr/0002-config-族的边界处置.md) 不迁）、`engine/` ✅（194，哨兵/abi 按 [`adr/0003`](adr/0003-哨兵与-abi-不独立成包.md) 不拆包）——导出面 429/429，两条护栏绿。**剩下的是下游那条腿（core/agent 切依赖）阻塞在 push**：`go get <本仓>@<sha>` 只能从远端解析，而 ③-④ 禁止 `replace`——本地提交不推上去，下游就没法 `require`。窗口期内 core 的 `pkg/domain` 与本仓的 `domain/` **并存（双份真相）**，两边都不许改，直到 core 侧 sed 落地并删掉旧包 |
| S2（QT-07） | 打 `v0.1.0`，下游换 tag | ⬜ |
| QT-19 | 版本规则表 + 兼容矩阵 + CHANGELOG | ✅ 已随 CT-02 落地（2026-10-22，护栏 `tools/check_release_matrix` 进 CI） |
