---
status: active
last-verified: 2026-10-24
verified-by: G0D 裁决材料（只出材料，不含实现；迁入时机由各服务仓协调者拍板）
---

# 服务 DTO 迁入裁决材料 —— QT-45 / QT-43 / QT-42

> **本文只出材料，不是决策**。三组 DTO 今天都住在各自服务仓的 `internal/` 下——Go 的 `internal` 规则本身就是一道跨仓边界，
> 因此「第二消费方」不能按「有没有人 import」单独判；还要看**将来边界是 Go 调用还是 HTTP**。
> 证据口径：`repo:path:line` 为 2026-10-24 实测。跨仓引用写成纯文本。

## 0. 结论先行

| 单 | 第二消费方 | 今天能迁吗 | 触发判据 | 一句话 |
|---|---|---|---|---|
| **QT-45**（broker） | **没有**（仓内 `ackstore` / `reconcile` 消费；core 与 UI 都走各自的路） | **不能** | core 在**进程内**接 `BrokerGateway`（X6a）且需要同一 Go 类型 | 现在迁 = 给一个「零外部消费」的 port 提前发版；contracts 已登记 `BrokerGateway` **未定形**（`tellurion-contracts:docs/API.md:651`） |
| **QT-43**（analytics） | **没有**（只有本仓 `internal/upstream` 实现端口） | **不能** | R2 的 **owner + wire 形状**先裁，再出现第二个 Go 消费方 | 形状与 owner 都没定（`tellurion-analytics:docs/TASKS.md:48`、`:67`、`:79`），现在迁 = 把临时镜像升格成契约 |
| **QT-42**（risk） | **设计上已指名 core，代码未接**（`ParseBundle` 的第二使用方） | **不能（但最接近）** | core 的 run 前消费接线（F3）开工，core 需要 `ParseBundle` 的 Go 类型 | 本仓台账已把触发条件写成一句话（`tellurion-risk:docs/TASKS.md:94-96`），是全组里唯一「判据已成文」的 |

**共同的硬约束**：`tellurion-contracts:docs/API.md:719` 写明跨进程 wire DTO 是 core `openapi.yaml` 的**生成物**，
契约仓只收生成物、禁手改。所以「服务 DTO」要分两种处置：

- **共享的 Go port/类型**（core 与服务在同一进程或同一 Go 依赖图里用）→ 才进 contracts；
- **纯 HTTP 报文**（UI / 别的进程只按 JSON 消费）→ 归 core 的 openapi，**不进 contracts**，也不该手写镜像。

---

## 1. QT-45 · `tellurion-broker` 的 Gateway DTO

### 1.1 当前落点

| 内容 | 路径 |
|---|---|
| 四个 port | `tellurion-broker:internal/gateway/ports.go:11`（`BrokerGateway`）、`:34`（`SettlementQuerier`）、`:42`（`KillSwitch`）、`:49`（`HealthChecker`） |
| 订单 DTO | `internal/gateway/model.go:80`（`OrderRequest`）、`:106`（`OrderAck`）、`:118`（`CancelRequest`）、`:124`（`OrderState`）、`:137`（`Fill`） |
| 账户 / 持仓 | `internal/gateway/model.go:153`（`Position`）、`:164`（`Cash`）、`:178`（`AccountSnapshot`） |
| 清仓结果 | `internal/gateway/model.go:188`（`FlattenedPosition`）、`:199`（`SkippedPosition`）、`:208`（`EmergencyFlattenResult`） |
| 幂等映射 | `internal/gateway/ports.go:58`（`IdempotencyRecord`）、`:68`（`AckStore`） |

### 1.2 第二消费方：没有

- 仓内消费：`internal/ackstore/pg.go:14`、`internal/reconcile/reconcile.go:17`（都 import `tellurion-broker/internal/gateway`）。
- 跨仓 import：在 `tellurion` / `tellurion-agent` / `tellurion-analytics` / `tellurion-orchestration` / `tellurion-strategy` /
  `tellurion-data` / `tellurion-risk` 全仓 grep `tellurion-broker`（`.go` + `go.mod`）**0 命中**。
- core 有**独立镜像**，不是消费方：`tellurion:pkg/live/trader.go:42`（`LiveTrader`）、`:141`（`OrderResult`）、
  `:168`（`AccountInfo`）、`:178`（`PositionInfo`）；`tellurion:pkg/live/reconciliation/reconciliation.go:131`（`AccountSnapshot`）。
- UI 消费的是 **core 的** `/api/execution/*`，不是 broker 的类型（`tellurion-ui:web/src/api/execution.ts:6-8`）。

### 1.3 迁入后谁改 import

- 立刻：`tellurion-broker:internal/gateway`（`mock.go` / `xtp.go` / `mock_test.go`）、`internal/ackstore`、`internal/reconcile`。
- 将来：core 的 `pkg/live` 一侧——**前提是 X6a 把 port 接在 core 进程内**（`tellurion-broker:docs/ARCHITECTURE.md:81`：
  「X6a/X6b 只在 core 进程内接 port」）。如果最终选择独立进程 + HTTP，那 core 只需要 wire DTO，**不该** import Go 类型，
  这时迁 contracts 就是错的落点。

### 1.4 判据与禁区

- **触发**：core 决定在进程内实现/消费 `BrokerGateway`，且两边要同一个 Go 类型。届时只迁**通道中立**的 DTO；
  `XTPConfig` / `XTPDeps`（`tellurion-broker:internal/gateway/xtp.go:22`、`:42`）是 SDK 适配细节，**永不进契约**。
- **禁区**：不要把 `OrderStatus` 和 K9 的 run 三值混用（broker 自己的 `ARCHITECTURE.md` F7 已点名）；
  `AccountSnapshot` 里的费用列属于柜台口径，进契约前必须写清「不是本地记账」。
- **同仓先例**：risk 对 `OrderIntent` / `Verdict` 已用同一条判据裁过——「port 留 core；X3b 冻结的是 **wire JSON**，
  不是 Go 类型。只有出现第二个 Go 使用方才按 ADR-0001 ① 迁 contracts」（`tellurion-risk:docs/ARCHITECTURE.md:115`）。
  broker 这组今天连第一个外部 Go 使用方都没有，比它更早。
- **今天不动**：contracts 的 `BrokerGateway` 仍按 `docs/API.md:651` 标「未定形 / S1 不落」。

---

## 2. QT-43 · `tellurion-analytics` 的 `HoldingsProvider` / `FactorExposureProvider`

### 2.1 当前落点

| 内容 | 路径 |
|---|---|
| 两个 port | `tellurion-analytics:internal/analysis/ports.go:30`（`HoldingsProvider`）、`:35`（`FactorExposureProvider`） |
| 持仓 DTO | `internal/analysis/model.go:78`（`Holding`）、`:89`（`PortfolioSnapshot`） |
| 暴露 DTO | `internal/analysis/model.go:98`（`FactorExposure`）、`:105`（`FactorExposures`） |
| 唯一实现 | `tellurion-analytics:internal/upstream/upstream.go:247-255`（`ArtifactReader` 的编译期断言）、`:263`（`GetPortfolioSnapshot`）、`:275`（`GetFactorExposures`） |

### 2.2 第二消费方：没有

- Go 消费方只有本仓：端口由消费端定义、`upstream` 实现（`internal/upstream/upstream.go:10-13` 的包注释写明来源）。
- **生产者都还没有**：core 全仓 grep `portfolio-snapshot` / `factor-exposures` / `factor_exposures`（`.go`）**0 命中**；
  今天读的是一个 interim artifact 面，形状与 owner 都未裁（`tellurion-analytics:docs/TASKS.md:48`、`:67`、`:79`）。
- HTTP 面是本仓自己的端点；没有第二个 Go 模块 import 这些类型（`internal/` 也不允许）。

### 2.3 迁入后谁改 import

- 立刻：`tellurion-analytics:internal/analysis`（定义处）、`internal/upstream`（实现处）、`internal/httpapi`（响应组装与 e2e 夹具）。
- 将来：artifact 的**生产者**（core 或 orchestration，取决于 R2 裁决）若要在 Go 层产出同一形状，才会 import。
  若生产者只写 JSON 文件 / HTTP，则同样只需 wire schema，**不迁 Go 类型**。

### 2.4 判据与禁区

- **触发（缺一不可）**：① R2 的 owner 与 wire 形状先裁（`tellurion-analytics:docs/TASKS.md:79`）；
  ② 出现第二个 **Go** 消费方（生产者用 Go 类型产出）；③ contracts 有对应的包（今天没有 `analysis` 包，
  `tellurion-analytics:docs/TASKS.md:82`）。
- **禁区**：不要把 interim artifact 形状当冻结形状迁进去；`PortfolioSnapshot` 目前混了 Brinson 需要的基准分行业权重与收益
  （`internal/analysis/model.go:85-102`），其中哪些是 artifact、哪些是请求参数，必须先拆清。

---

## 3. QT-42 · `tellurion-risk` 的 `Snapshot` / `SnapshotRule`

### 3.1 当前落点

| 内容 | 路径 |
|---|---|
| `SnapshotRule` | `tellurion-risk:internal/rule/snapshot.go:24-31`（`RuleID` / `Revision` / `Scope` / `Type` / `Params json.RawMessage` / `VersionHash`） |
| `Snapshot` | `tellurion-risk:internal/rule/snapshot.go:37-41`（`RiskRulesHash` / `Rules` / `CreatedAt`） |
| 域标签 + hash | `internal/rule/snapshot.go:15-20`（`DomainSnapshotHash = "tellurion.risk_snapshot_hash.v1"`）、`:57`（`SnapshotHash`）、`:108`（`verifyHash`） |
| 进程内边界 | `internal/rule/snapshot.go:120-131`（`ParseBundle` / `Bundle.ForScope`） |
| wire 出口 | `internal/httpapi/api.go:86`（`resolveBody.Snapshot rule.Snapshot`）；`GET /snapshots/{hash}/bundle` body 就是 canonical snapshot JSON（`tellurion-risk:docs/api-schema.md:43`） |

### 3.2 第二消费方：设计已指名，代码未接

- 今天：只有 risk 仓自己（`internal/pgstore/pgstore.go:491` 读快照行、`internal/rule/service.go:405-425` 组快照）。
- 指名：core 的**进程内** `RiskPolicy` 必须吃 `ParseBundle`，并回写 `risk_rules_hash`
  （`tellurion-risk:docs/TASKS.md:94-96` 原文：「`ParseBundle` 的 Go 类型出现第二个使用方（core）后，……迁 contracts」；
  `:131`、`:165`、`:200` 三处施工记录重复同一判断）。
- core 侧今天只有列，没有类型：`tellurion:pkg/storage/run_provenance.go:38`、`:61` 的 `risk_rules_hash`；
  全仓没有 `ParseBundle` / `SnapshotRule` 命中。

### 3.3 迁入后谁改 import

- 立刻：`tellurion-risk:internal/rule`（定义与 hash）、`internal/httpapi`（wire 组装）、`internal/pgstore`（读行）、
  `internal/rule/*_test.go`。
- 触发后：core 新增的 bundle 读取/`RiskPolicy` 组装代码（F3 的施工点）从 contracts import。

### 3.4 判据与禁区

- **触发**：core 的消费接线（F3）真的开工，且选择**同进程 Go 类型**（设计就是这么裁的）。若 core 改走「HTTP bundle + 在 core
  自建一份解析类型」，那就等于主动造第三份形状——必须在开始时就说清，不能默默发生。
- **迁什么**：`Snapshot` / `SnapshotRule` 是 DTO；但 `SnapshotHash` 是 canonical 算法，按「谁算 hash 谁唯一」
  （`tellurion:docs/design/qt-decoupling/w1-repo-ownership.md` J4）应随域标签进 `engine/`（或至少由契约仓提供唯一实现），
  **不能**让 core 另写一份。
- **形状坑**：`CreatedAt` 不参与 hash（`internal/rule/snapshot.go:54-56` 注释），迁移时要把「身份字段 vs 元数据」写死；
  `Params json.RawMessage` 是 JSONB 逃生阀，canonical 编码规则必须先冻结（键排序 / 空白 / 数字），否则两仓算出不同 hash。
- **今天不动**：`RiskPolicy` 的同步裁决留 core（`tellurion-risk:docs/ARCHITECTURE.md:13` 的边界判断），本组不碰 core。

---

## 4. 三组共性判据（给协调者的一张表）

| 判据 | QT-45 | QT-43 | QT-42 |
|---|---|---|---|
| 第二份**独立实现**已存在 | ❌（core 是另一套旧面，不是同一 port） | ❌ | ❌（但已指名） |
| 跨进程 / 跨模块的 **Go** 消费方 | ❌ | ❌ | 设计上 ✅，代码 ⏳ |
| wire 形状已冻结 | ❌（port 未定形） | ❌（R2 未裁） | ✅（`docs/api-schema.md:43` + canonical JSON） |
| 迁入触发已写成一句话 | ❌ | ❌ | ✅（`docs/TASKS.md:96`） |
| **今天建议** | 不迁 | 不迁 | 等 F3 开工那批一起迁 |

**共同禁区**（与 G0D 派工一致）：本组不先落 `StrategyVersionInput`、`dataprotocol` 或任何服务 DTO；
不改 `tellurion-broker` / `tellurion-analytics` / `tellurion-risk` 的代码；不替这些仓决定 port 形状。

## 5. 分发与边界

- **只出材料**：三组 DTO 一个都没迁，`git status` 里本组只应有文档。
- **上报对象**：按派工要求报 G1 / G4 / G5 / G6 / G7 协调者（仓与组的对应以派工表为准）；受影响仓：
  `tellurion-broker`、`tellurion-analytics`、`tellurion-risk`，以及作为消费侧的 `tellurion`。
- **材料复用**：QT-43 / QT-42 都引用了「接口对表」模板（`tellurion:docs/design/qt-decoupling/data-service-contract-and-ownership.md:151-155`），
  正式迁入前建议各自补齐该仓 §9 的对表行，再按本文判据执行。
