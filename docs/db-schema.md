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
| `BacktestResult` / `Trade` / `Factor*` / `Direction` | `run`（结果） | **owner = core（2026-10-21 裁决，单一写者；orchestration 只读）** |
| `RunSpec` / `SignalSpec` / `ResultSpec` | `run`（job 元数据） | `tellurion-orchestration` |
| `strategy_ref`（`id@version_hash`） | `strategy` | `tellurion-strategy` |
| `risk_rules_hash` / 规则快照 | `risk` | `tellurion-risk` |

> ⚠️ **`ResultSpec` 的形状是冻结契约**（`tellurion:docs/SPEC.md`「ResultSpec 冻结（K9，2026-10-12）」节：字段集 27 项、指标口径、`status` 三值语义三者都冻结）：改形状属**改冻结契约**，必须先破护栏测试再走正式契约变更评审，**不得顺手改**。

> 判断一条类型该不该进本仓：**被 ≥2 个仓使用**才算；只有一个仓用的，留在那个仓（否则这里会变成公共代码堆放处）。

## 2. 只读输入（不属于本仓，禁止写入）

不适用（本仓不访问数据）。

## 3. 红线

1. **零实现**：不许出现函数体里做业务判断的代码；接口 + 数据 + 常量，仅此三类。
2. **零 DB**：不许 `import` 任何数据库驱动或 storage 包。
3. **禁 `Config` / `Options` / `Impl` 命名**：那是实现细节漏进契约的形状。
4. **禁反向依赖**：不许 import core / agent / 任何服务仓（编译器会拦，但要写清楚为什么）。
5. 结构级变更 = 一次版本发布：改这里必须能回答"哪个下游需要它"。

## 4. 数据治理（本仓无表，但影响别人的表）

- 本仓**不建表、不迁移、不连库**，因此没有 `migrations/`，将来也不该有。
- 但本仓的形状**决定下游的表**：`UnifiedDataPoint` 多一个字段，`tellurion-data` 的表就要多一列。所以：
  1. 形状变更 = 一次契约版本发布（见 [`api-schema.md`](api-schema.md) §4）；
  2. **谁的表谁迁**：本仓不发 DDL、不写迁移说明，只在契约定版时列出"哪些下游需要动表"；
  3. 数据库约束（主键 / 唯一 / 索引 / 分区）**不进本仓** —— 那是实现的选择，不是形状。
- 判据：一个类型一旦开始带 `NOT NULL` / 索引 / 分区这类词，它就属于某个仓的 `db-schema.md`，不属于这里。

## 5. 待定

- 版本规则表与兼容矩阵：QT-19（未落地前不要凭印象发版）。
