---
service-card: 1
repo: tellurion-contracts
kind: library
lifecycle: production
owner: ruoxizhnya
impact: compile-time

requires: []
required-by:
  - tellurion
  - tellurion-agent
  - tellurion-data
  - tellurion-strategy
  - tellurion-orchestration
  - tellurion-broker
calls: []
called-by: []

truth:
  readme: README.md
  agents: AGENTS.md
  architecture: docs/ARCHITECTURE.md
  api: docs/API.md
  db: docs/db-schema.md
  dataflow: docs/dataflow.md
  runbook: null
  tasks: docs/TASKS.md

status: evergreen
last-verified: 2026-11-11
verified-by: 2026-11-11 实读 README / AGENTS / API.md / api-schema / db-schema / dataflow / go.mod，并按工作区各仓 go.mod 反向扫描消费方；未写版本号与导出符号数量
---

# 服务卡：tellurion-contracts

> **一句话**：跨仓唯一缝——port 接口、DTO / 值对象、四个引擎契约、哨兵错误、abi 常量；零实现、零 DB、零 HTTP。
>
> `library` · `production` · `compile-time` · `ruoxizhnya`

## 1. 这是什么

它是 12 仓依赖图的根：所有跨边界交换的**形状**只住在这里——port 接口、DTO 与值对象、四个引擎契约、跨仓哨兵错误、abi 常量，以及跨仓一致性语料。它不实现任何东西；谁落库、谁开端口、谁做业务判断，都由下游仓决定。

**它拥有**：`fees` / `domain/market` / `domain` / `engine` / `dataprotocol` / `risk` / `run` / `broker` 等契约包；`docs/API.md` 的导出面登记；`docs/api-surface.txt` 的签名冻结基线；`conformance/cases.json` 的跨仓语料。
**它不拥有**：任何数据库表、HTTP 服务、业务规则、实现代码。

## 2. 它不做什么

- 不放任何实现、业务规则、DB 访问；值对象自带的纯函数方法按 ADR-0004 的判据例外处理。
- 不放只有单个仓使用的类型；只收**被 ≥2 个仓使用**的跨边界形状。
- 禁止 `Config` / `Options` / `Impl` 命名（例外只认 ADR-0002 登记的项）。
- 禁止反向 import core / agent / 任何服务仓。
- 不提供 HTTP / gRPC / 消息队列等任何运行时入口。

## 3. 与谁相连

**运行期没有边**：本仓不被任何进程加载为服务，也不调用任何服务。它的「连接」是编译期依赖。

### 谁 require 它（编译期）

| 消费方 | 用法 | 它坏了会怎样 |
|---|---|---|
| `tellurion` | 合约 DTO、四契约、哨兵、abi、risk / run / broker wire 面 | 编译失败；core 无法构建 / 发版 |
| `tellurion-agent` | 合约 DTO、conformance 语料 | 编译失败 |
| `tellurion-data` | `dataprotocol` 等数据面形状 | 编译失败 |
| `tellurion-strategy` | 策略资产与回测 wire 面 | 编译失败 |
| `tellurion-orchestration` | `run` 编排 wire 面 | 编译失败 |
| `tellurion-broker` | `broker` 订单 / 成交 / 账户 wire 面 | 编译失败 |

> 当前未直接 require：`tellurion-risk`、`tellurion-analytics`、`tellurion-kit`、`tellurion-docs`、`tellurion-ui`。`contracts/risk` 包由 core 消费，不代表 risk 仓直接依赖本仓。
> 消费方清单以各仓 `go.mod` 为准；本卡只记录方向，不复制 pin 版本。版本矩阵见 core `docs/REPO-MAP.md`。

### 它调用（出边）

| 目标 | 协议 | 用来做什么 | 目标挂了会怎样 |
|---|---|---|---|
| 无 | — | — | — |

## 4. 运行时身份

- **常驻进程**：无，纯库。
- **存活/就绪语义**：无端口、无探针。
- **构建入口**：`README.md` 快速开始；单独 checkout 必须能跑通 `go build ./...` / `go vet ./...` / `go test ./...` 与四条导出面护栏。
- **端口**：不适用。

## 5. 数据与契约归属

| 维度 | 归属 | 真相源 |
|---|---|---|
| DB schema | 无表、不连库、不迁移 | `docs/db-schema.md` |
| 跨仓 DTO / 契约 | 本仓是唯一 owner；导出面逐符号登记 | `docs/API.md`、`docs/api-surface.txt` |
| 对外 HTTP | 无 | — |

> 只写「拥有哪块」，不写导出符号数量、包清单的完整枚举。逐符号稳定性词表见 `docs/API.md` §2。
> **未裁事项**：`contracts/risk` 已存在，但 `tellurion-risk` 仍自持一份 `Snapshot` 形状、尚未采用本包；两者合并或择一为唯一真相**待裁决**，见 `tellurion-docs:docs/design/ddd-boundary-review-2026-11-11.md` C5/D3。

## 6. 坏了怎么办

- **影响面**：消费方**编译失败**，不是运行期 503；已经构建好的二进制不受影响。破坏性变更会让下游在升级 pin 时才暴露，见 `CHANGELOG.md` 的 breaking 条目。
- **先确认（按顺序）**：
  1. 消费方 `go build ./...` 是否因合约签名 / 字段 / 哨兵变化失败；对照 `CHANGELOG.md` 与兼容矩阵。
  2. 本仓四条导出面护栏是否红：`check_api_surface` / `check_surface_freeze` / `check_contract_shape` / `check_release_matrix`。
  3. 下游 pin 是否落在兼容矩阵声明的范围内；跨仓 pin 核对归 core `docs/TASKS.md` 的 QT-09。
- **升级路径**：`owner` → 跨仓争议 → core `docs/TASKS.md`
- **runbook**：无（纯库不需要起停手册）；构建与发版见 `README.md`「快速开始」「版本规则」。

## 7. 5 分钟后去哪

| 我想…… | 读这个 |
|---|---|
| 知道每个导出符号的稳定性与错误语义 | `docs/API.md` |
| 知道对外 API 的形状与稳定性词表 | `docs/api-schema.md` |
| 知道数据形状怎么流到下游 | `docs/dataflow.md` |
| 知道为什么本仓不建表 | `docs/db-schema.md` |
| 知道版本规则与兼容矩阵 | `README.md` / `CHANGELOG.md` |
| 当前在做什么 | `docs/TASKS.md` |
| 架构与决策 | `docs/ARCHITECTURE.md` / `docs/adr/` |
