---
status: active
last-verified: 2026-10-22
verified-by: CT-02 结案（版本纪律三件套进导航与变更触发表）
---

# 文档导航 — tellurion-contracts

> 本目录只讲**本仓**。跨仓拓扑、阶梯与全局决策读 `tellurion:docs/`（见文末）。
> 跨仓引用一律写成 `repo:path` 纯文本，**不做可点击相对链接**——独立 checkout 时 `../../` 必断。

## 1. 先读哪份

| 你想干什么 | 读这份 |
|---|---|
| 第一次进本仓 | [`../AGENTS.md`](../AGENTS.md) |
| 开工前找活 | [`TASKS.md`](TASKS.md)（**唯一台账**） |
| 确认"这个功能到底归不归我" | [`ARCHITECTURE.md`](ARCHITECTURE.md) 的定位与「三份契约」 |
| 改表 / 加列 | [`db-schema.md`](db-schema.md) |
| 加 / 改接口 | [`api-schema.md`](api-schema.md) |
| 加 / 改导出符号 | [`API.md`](API.md)（**CI 核对，不登记即红**） |
| 发版 / 改版本纪律 | [`../README.md`](../README.md) 的「版本规则 / 兼容矩阵 / 发版流程」+ [`../CHANGELOG.md`](../CHANGELOG.md)（**CI 核对，矩阵与账对不上即红**） |
| 改数据流向或顺序 | [`dataflow.md`](dataflow.md) |
| 想知道"为什么这样设计" | [`adr/`](adr/README.md) |
| 要做删档 / 迁移 / 流程变更 | [`odr/`](odr/README.md) |
| 查历史 | [`archive/`](archive/README.md) |

## 2. 文档地图

| 文件 | 职责 | 谁维护 |
|---|---|---|
| [`ARCHITECTURE.md`](ARCHITECTURE.md) | 高层设计：分层、组件、依赖方向、非目标 | 改结构的人 |
| [`db-schema.md`](db-schema.md) | **表结构唯一真相源** + 所有权声明 | 改表的人 |
| [`api-schema.md`](api-schema.md) | **对外接口唯一真相源** + 错误码 / 幂等 / 兼容 | 改接口的人 |
| [`API.md`](API.md) | **导出面清单**：逐个符号的稳定性与错误语义；§4 的登记区由 `tools/check_api_surface` 在 CI 核对 | 加/改导出符号的人 |
| [`../README.md`](../README.md) §版本规则/兼容矩阵 + [`../CHANGELOG.md`](../CHANGELOG.md) | **版本纪律三件套**：三行 bump 表、按 `vX.Y` 序列建行的矩阵、每个发版一节的账本；由 `tools/check_release_matrix` 在 CI 与 tag 三方对拍 | 发版的人 |
| [`dataflow.md`](dataflow.md) | 数据流向、时序、fail-closed 点、不变量 | 改链路的人 |
| [`design/`](design/README.md) | **跨仓落点的裁决材料**（只出证据与影响面，不含实现；拍板后落成代码/台账） | 提变更的人 |
| [`TASKS.md`](TASKS.md) | 可执行任务（六字段） | 每个人 |
| [`adr/`](adr/README.md) | 架构决策（只增不删） | 提变更的人 |
| [`odr/`](odr/README.md) | 运行 / 治理决策（删档、审计、迁移、工具） | 做事的人 |
| [`archive/`](archive/README.md) | 已完成任务与历史文档 | 收尾的人 |

## 3. 变更触发（代码动了，文档必须跟着动）

| 代码变更 | 必须更新的文档 |
|---|---|
| 建表 / 加列 / 加索引 / 改主键 | `db-schema.md`（+ 结构级变更立 ADR） |
| 新增 / 修改 / 废弃接口 | `api-schema.md`（错误码与幂等两节要一起看） |
| 新增 / 改名 / 删除**导出符号** | `API.md` §4 登记区（`go run ./tools/check_api_surface` 当场红，不是提醒） |
| 增删改 `conformance/cases.json` 的用例或期望值 | 三仓同测：本仓 `go test ./conformance`、core `pkg/conformance`、agent `internal/conformance`（core 侧编排 `make contracts-conformance`）。期望值必须来自契约仓实现，不许按某一侧实现反推 |
| 改导出符号的**签名 / 字段 / 底层类型** | `docs/api-surface.txt` 冻结基线（`go run ./tools/check_surface_freeze -update` 重生成；不 update 就红）。与上一行是两道门：新增符号要同时过登记表与基线 |
| **发版**（打 tag） | `README.md` 兼容矩阵那一行 + `CHANGELOG.md` 那一节（`go run ./tools/check_release_matrix` 当场红；口径源头在 `tellurion:docs/TASKS.md` QT-19 ③） |
| 新增数据源、改聚合口径、改顺序 | `dataflow.md` |
| 改分层、依赖方向、模块职责 | `ARCHITECTURE.md` + 一条 ADR |
| 删档 / 迁移 / 换 CI / 换工具链 | 一条 ODR |
| 完成一张单 | `TASKS.md` 标 ✅，整条移入 `archive/` |

> 这条表的作用不是提醒，是**评审时的检查项**：PR 里改了上表左列，右列没动，就该被打回。

## 4. 跨仓只读资料（不在本仓维护）

| 主题 | 位置 |
|---|---|
| 仓清单与命名规则 | `tellurion:docs/design/qt-decoupling/repo-inventory.md` |
| 分仓拓扑与 S0~S3 阶梯 | `tellurion:docs/design/qt-decoupling/repo-topology-and-naming.md` |
| 全局决策 / ADR-D 索引 | `tellurion:docs/design/qt-decoupling/project-engineering-file.md` |
| 数据归属总表（DDD 分库） | `tellurion:docs/design/qt-decoupling/data-service-contract-and-ownership.md` |
| 跨仓任务总台账 | `tellurion:docs/TASKS.md` |
