---
status: active
last-verified: 2026-10-24
verified-by: G0D 裁决材料（只出材料，不含实现；X1a 开工前由协调者二选一）
---

# QT-40 · X1a 裁决材料 —— `dataprotocol` 的两种读法（先迁 / 维持 core-only）

> **本文只出材料，不是决策**。`tellurion-contracts` 今日**没有** `dataprotocol` 包；本仓 v0.4.0 的包清单是
> `fees` / `domain` / `domain/market` / `engine`（+ 零导出符号的 `conformance`）。
> 证据口径：`repo:path:line` 为 2026-10-24 实测。跨仓引用写成纯文本。

## 0. 结论先行（含两处对派工前提的订正）

1. **`UnifiedDataPoint` 不是「两份副本」，是「三份形状」**：`tellurion-data:internal/ingest/point.go:18-35`、
   `tellurion:pkg/data/source/unified.go:16-48`、`tellurion:pkg/storage/bulk_insert.go:17-24`。
   派工与 `tellurion-contracts:docs/API.md:731` 只数了 core 的两份，漏了 data 仓这份**已经带 `VisibleAt` 的镜像**。
2. **三份已经在漂**，不是「将来可能漂」：
   - data 镜像有 `VisibleAt`（`point.go:33-35`），core 与 storage 没有；
   - data 镜像**没有** `DeduplicateKey()`（`point.go:50-53` 只有 `ToJSON`），core 有（`unified.go:70-72`）；
   - 构造器时钟策略相反：data 由调用方注入 `ingestTime`（`point.go:37-49`），core **读墙钟**（`unified.go:50-58`）。
   ⇒ 「维持 core-only 直到出现第二个消费方」在今天会得出「第二个消费方永远不会出现（因为 data 用的是自己的副本，
   不是 import）」，从而把协议**永久写死在数据服务里**——正是设计文档警告的代价
   （`tellurion:docs/design/qt-decoupling/data-service-contract-and-ownership.md:163`）。
3. **原读法 (b) 的措辞已过期**：「随 S1 一起落、进 core-only 清单」是 QT-05 施工期写的；
   S1 已结案（本仓 `docs/TASKS.md` QT-05 标 🔶 但本仓侧完成）。今天再选 (b)，只能作为**一次独立 MINOR 发版**
   并把新包补进 ③-⑨ 的 core-only 回收清单，不能假装它「随 S1 顺带进来」。
4. **建议**：X1a 选 **(a) 先迁**，但附带 §3 的四个形状前置——尤其 **`VisibleAt` 必须在类型里**（否则 data 仍得留镜像，
   迁移等于没做），且**契约包不得带读墙钟的构造器**。判据不是「有没有第二个 import」，而是「有没有第二份独立实现」；
   今天两处独立实现已经存在。

---

## 1. `UnifiedDataPoint` 的字段来源（逐字段）

| 字段 | `tellurion-data:internal/ingest/point.go` | `tellurion:pkg/data/source/unified.go` | `tellurion:pkg/storage/bulk_insert.go` | 语义（冻结时必须写清） |
|---|---|---|---|---|
| `Symbol` | `:20` | `:19` | `:18` | 规范代码（`600519.SH` / `AAPL`）；**不做交易所后缀归一**（归一在 connector） |
| `TradeTime` | `:23` | `:25` | `:20` | 数据点归属时刻。日频 = 交易日零点；实时 = 快照时刻；财报 = 报告期截止日（PIT 另走 `VisibleAt`） |
| `Source` | `:25` | `:28` | `:21` | connector 名，落 `source` 列；**同表多源口径靠它区分**（DBT-09 的教训） |
| `DataType` | `:27` | `:31` | `:22` | 裸 `string`（三份都没有具名类型），取值集由 `DataType*` 常量给出；决定写哪张表 |
| `Data` | `:30` | `:34` | `:23` | `map[string]interface{}`，JSONB 形状的逃生阀；**必含键由 connector 声明与 normalizer 决定**，不在类型上表达 |
| `IngestTime` | `:31-32` | `:37-38` | `:24` | 落库时刻；**必须由调用方注入**（可复现：同输入同字节） |
| `VisibleAt` | `:33-35`（**仅 data 有**） | — | — | 按声明的可见锚点导出的 PIT 可见时刻，由 `GatePoints` 填（`tellurion-data:internal/ingest/gate.go:58-79`）。X1d 给表加列后落库（`tellurion-data:docs/TASKS.md:35` 的设计前提） |

行为面（不是字段，但同样是形状）：

| 行为 | data | core source | core storage | 冻结建议 |
|---|---|---|---|---|
| `ToJSON()` | ✅ `point.go:51-53` | ✅ `unified.go:59-61` | ❌ | 保留（纯函数，ADR-0004 允许随值对象走） |
| `DeduplicateKey()` | ❌ | ✅ `unified.go:70-72`（`Symbol|DataType|TradeTime(RFC3339Nano,UTC)`） | ❌ | **必须进共享类型**：这是幂等坐标的唯一实现；留在 core 会让 data 服务复制第二份幂等逻辑 |
| 构造器 | `NewUnifiedDataPoint(..., ingestTime)`（`:37-49`） | `NewUnifiedDataPoint(...)` **读 `time.Now()`**（`:50-58`） | — | **不进 contracts**：契约包不读墙钟；两个仓各留自己的构造/归一入口。data 版注入时钟的做法是对的，core 版应随迁移改为注入 |

## 2. 两条读法的影响面

### 2.1 (a) 先作为 `tellurion-contracts` 新包迁入

**做什么**

1. 本仓新增 `dataprotocol/`（包名 `dataprotocol`，落点已定名见 `tellurion-contracts:docs/API.md:731`）：
   `UnifiedDataPoint` 七字段（含 `VisibleAt`）+ `ToJSON` + `DeduplicateKey`。**不带构造器、不带时钟**。
2. `tellurion-data`：删 `internal/ingest/point.go:18-53` 的类型与两个方法，改 import 契约包；
   `gate.go` 继续填 `VisibleAt`（类型字段已在契约里）；构造仍走本仓一个注入时钟的 helper（可在 `ingest` 内留薄包装）。
3. `tellurion`：删 `pkg/data/source/unified.go:16-61` 与 `pkg/storage/bulk_insert.go:12-24` 两份，
   同时删 `etl.go:14-25` 的 `toStoragePoints` 投影（成环理由消失）；所有 `source.` / `storage.` 前缀改契约包。
4. **同一批**里补齐 `ConnectorDeclaration` 家族（`tellurion:pkg/data/source/contract.go:41-59` 等，见设计 A 组
   `data-service-contract-and-ownership.md:172-175`）——否则只有点位进契约、声明仍留 core，X1b 的 aggregator 仍要
   反向理解 core 的声明类型。若协调者只批 `UnifiedDataPoint`，第二项要显式登记为 X1b 前置，不能默认它「随之解决」。

**代价 / 风险**

- 这是**新包**（MINOR），不是既有包的字段：`check_api_surface` 要登记、`check_surface_freeze` 要 `-update`。
- 形状冻结前必须先拍三件事（§3）：`VisibleAt` 是否必进、`DataType` 是否具名、`DeduplicateKey` 的时钟格式。
- core 删两份副本会碰到 `pkg/data` / `pkg/storage` 大量调用点；`bulk_insert_ddl_test.go` 只核「目标表在 DDL 里」，
  **核不到字段顺序**——迁移必须补一条「共享类型字段 ↔ `TableMapper` 列布局」的对表测试，否则是假安全。
- 版本：`tellurion-data` 今日 pin `v0.1.0`（`tellurion-data:go.mod:7`），切包要升 pin 并重跑单仓绿。

**收益**

- 三份形状归一；`VisibleAt` 的 PIT 语义只定义一次；`DeduplicateKey` 不再有第二份实现的空间。
- X1b 的 aggregator 从第一天就 import 契约，而不是「先把协议写死在服务里，之后再搬」（设计文档点名的返工）。

### 2.2 (b) 维持 core-only 直到出现第二个消费方

**做什么**：不动 contracts；X1b 在 `tellurion-data` 内用自己的 `internal/ingest` 形状实现 aggregator；core 继续双副本。

**代价 / 风险**

- **判据永远不会触发**：data 已经有独立实现，它不需要 import core；「第二个 import 方」在可预见期内不会出现。
  于是 X1b 落地后，协议的事实标准变成 data 服务的内部类型——正是 `data-service-contract-and-ownership.md:163`
  说的「协议没冻结就写 aggregator = 把协议写死在一个服务里」。
- **已存在的漂移会继续扩大**：`VisibleAt` / `DeduplicateKey` / 时钟策略三处分叉，任何一处补字段都可能只改一边；
  storage 的列布局漏写一列不会红（见上）。
- **将来合并更贵**：X1d 给 12 张表加 `visible_at` 之后，再迁就要同时改 data、core 的 connector/ETL、storage 三处，
  且要处理已落库数据的 PIT 语义。
- **流程上也不便宜**：S1 已结束，(b) 只能靠一次独立 MINOR 进入本仓，再补 ③-⑨ 清单；不是「免费等待」。

**收益**

- 今天不动 contracts、不动 core，短期 diff 最小。
- 如果协调者判断「协议形状还会大改（例如 `DataType` 要变成具名枚举、`Data` 要去 JSONB 化）」，等形状稳定再迁，
  可以少一次 MINOR。

## 3. 迁入判据（若选 (a)，四个前置缺一不可）

| # | 判据 | 为什么不能省 | 证据 |
|---|---|---|---|
| 1 | **第二份独立实现已存在**，不是「第二个 import」 | 否则规则会奖励「各写一份副本」；今天 data 镜像就是第二份实现 | `tellurion-data:internal/ingest/point.go:8-17`（注释自述「X1a 冻结前以本地镜像承载写入面」） |
| 2 | **`VisibleAt` 进类型** | 不进，data 就必须继续留镜像（它的写入路径依赖该字段），迁移等于没做 | `point.go:33-35`、`gate.go:58-79` |
| 3 | **契约包零 I/O、零时钟** | 构造器读墙钟会让「同输入同字节」失效；ADR-0004 只允许纯函数 | `unified.go:50-58`（core 读 `time.Now()`） vs `point.go:37-49`（data 注入） |
| 4 | **同批删副本 + 补字段对表测试** | 只迁不删 = 四份形状；无对表测试 = 漏列静默 | `bulk_insert.go:12-24`、`etl.go:14-25`、`tellurion:pkg/storage/bulk_insert_ddl_test.go`（只核表名） |

补充口径：

- `TimeSemantics` **不在 `dataprotocol` 里重复定义**——它已随 QT-05 整包住 `tellurion-contracts:domain/market`，
  重复定义等于把 `docs/API.md:731` 的 #1/#4 合成一个真事故。
- `DataType` 今天三份都是裸 `string`（`point.go:27`、`unified.go:31`、`bulk_insert.go:22`）。具名化是**加法但会碰 wire**；
  建议 X1a 先保持 `string` + 常量，等 X1b 的 schema 冻结再评估（设计文档 `:174` 的「加字段优先落 `Data`」同理）。
- `Aggregator` 接口**不随 X1a 落**：今天只有 HTTP 形状，没有实现（`tellurion-contracts:docs/API.md:652` 的 `DataPort`
  仍标「未定形 / S1 不落」）；不许在契约仓造一个无消费方的 port。

## 4. 分发与边界

- **只出材料**：本组不创建 `dataprotocol`、不改 `tellurion` / `tellurion-data`、不打 tag。
- **上报对象**：按派工要求报 G1 / G4 / G5 / G6 / G7 协调者（仓与组的对应以派工表为准）；受影响仓：
  `tellurion-contracts`（新包）、`tellurion`（删两份副本）、`tellurion-data`（镜像改 import + `VisibleAt` 落点）。
- **禁区复核**：未先落 `dataprotocol` 或任何 `UnifiedDataPoint` 代码；未把 `TimeSemantics` 复制进来；未替 data 服务裁决 aggregator 形状。
