---
status: active
last-verified: 2026-10-20
verified-by: 功能边界固化（三份契约：db / api / dataflow）
---

# 架构 — tellurion-contracts

## 1. 定位

跨仓唯一缝隙：**port 接口 · DTO/值对象（含 `domain`）· 四个引擎契约 · 哨兵错误 · abi 常量**。

它的存在意义是让"改契约"成为一次**显式发布**，而不是某天某人顺手改了 core 里的一个 struct。

## 2. 分层

```
（被所有 Go 仓依赖；自己只依赖 tellurion-kit 的 decimal / errors）

  domain/     值对象：OHLCV / Trade / BacktestResult / Factor* / Direction
  engine/     四契约：RunSpec / SignalSpec / ResultSpec / EngineRunner + MarketRouter
  fees/       费率 / 滑点 / 成本模型的**定义**（计算留在 core）
  errors/     三类哨兵
  abi/        host 函数版本、内存区、枚举上限
```

## 3. 组件职责

| 包 | 职责 | 不该做什么 |
|---|---|---|
| `domain/` | 跨边界值对象 | 不做计算、不读配置 |
| `engine/` | 回测四契约 + 路由接口 | 不含引擎实现 |
| `fees/` | 成本模型的形状 | 不算钱（算法在 core） |
| `errors/` | 三类哨兵 + 包装约定 | 不做日志 |

## 4. 三份契约（功能的边界就在这三份文件里）

| 维度 | 文件 | 回答什么问题 |
|---|---|---|
| 数据形状 | [`db-schema.md`](db-schema.md) | 本仓拥有哪些表？谁能写？ |
| 接口形状 | [`api-schema.md`](api-schema.md) | 对外暴露什么？错误码与幂等怎么做？ |
| 流向与顺序 | [`dataflow.md`](dataflow.md) | 数据怎么流？哪里 fail-closed？ |

## 5. 依赖方向

```
tellurion-kit（decimal / errors）
        ↑
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
| S1（QT-05） | `pkg/domain` + `pkg/backtest/contracts` 搬入本仓，core/agent 改 `require` 伪版本 | ⬜ |
| S2（QT-07） | 打 `v0.1.0`，下游换 tag | ⬜ |
| QT-19 | 版本规则表 + 兼容矩阵 + CHANGELOG | ⬜ |
