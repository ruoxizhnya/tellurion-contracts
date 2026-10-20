---
status: active
last-verified: 2026-10-20
verified-by: 功能边界固化（三份契约：db / api / dataflow）
---

# DB Schema — tellurion-contracts

> **所有权声明**：本仓**不持有任何数据库表，也不提供 HTTP 服务**。它是依赖图的根：只定义跨边界交换的**形状**（类型 / 接口 / 错误 / 常量）。谁落库、谁开端口，都由下游仓决定。
> 上级依据：`tellurion:docs/design/qt-decoupling/data-service-contract-and-ownership.md` §1 —— **DDL 只有一个 owner**，同一个 PG 实例、每服务一个 schema，将来拆库是"按 schema 搬走"而不是重写。
> 本文件是本仓表结构的**唯一真相源**：代码里的建表语句与它不一致时，**改代码**；要改表结构，先改本文件（结构级变更还要立 ADR）。

## 1. 归属表

> **无表**。本仓不建表、不迁移、不连接数据库。
>
> 它定义的是"别人要落库的数据长什么样"：

| 形状 | 落在哪个 schema | owner |
|---|---|---|
| `OHLCV` / `UnifiedDataPoint` / `FrequencySpec` / `TimeSemantics` | `market`、`ingest` | `tellurion-data` |
| `BacktestResult` / `Trade` / `Factor*` / `Direction` | `run`（结果） | `tellurion-orchestration` / core |
| `RunSpec` / `SignalSpec` / `ResultSpec` | `run`（job 元数据） | `tellurion-orchestration` |
| `strategy_ref`（`id@version_hash`） | `strategy` | `tellurion-strategy` |
| `risk_rules_hash` / 规则快照 | `risk` | `tellurion-risk` |

> 判断一条类型该不该进本仓：**被 ≥2 个仓使用**才算；只有一个仓用的，留在那个仓（否则这里会变成公共代码堆放处）。

## 2. 只读输入（不属于本仓，禁止写入）

不适用（本仓不访问数据）。

## 3. 红线

1. **零实现**：不许出现函数体里做业务判断的代码；接口 + 数据 + 常量，仅此三类。
2. **零 DB**：不许 `import` 任何数据库驱动或 storage 包。
3. **禁 `Config` / `Options` / `Impl` 命名**：那是实现细节漏进契约的形状。
4. **禁反向依赖**：不许 import core / agent / 任何服务仓（编译器会拦，但要写清楚为什么）。
5. 结构级变更 = 一次版本发布：改这里必须能回答"哪个下游需要它"。

## 4. 迁移策略

- 迁移脚本落在 `migrations/`，命名 `<序号>_<动作>.sql`（全仓一个递增序列，不复用号）。
- **只前滚**：已合并的迁移文件不改；修错靠新增一条。
- 每条迁移必须可重入（`IF NOT EXISTS` / `IF EXISTS`），且能在空库上一次跑完。
- 跨 schema 取数只走**只读视图或 API**；出现 `join` 别人的表就是边界腐化的第一信号。
- 表结构变更的同步义务：本文件 + `docs/ARCHITECTURE.md` 的「数据归属」节必须同一次提交里更新。

## 5. 待定

- 版本规则表与兼容矩阵：QT-19（未落地前不要凭印象发版）。
