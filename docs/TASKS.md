---
status: active
last-verified: 2026-11-13
verified-by: 2026-11-13 `CT-COVERAGE-FLOOR` 收口 ✅——本机全量实测 **79.3%**（不带 DSN 的全套件），CI 门取「实测 - 5 点」= **74%**，执法步骤 `go test ./... -coverprofile=coverage.out` + `go tool cover -func`（语句加权）；破坏验证：删 `tools/check_release_matrix/main_test.go` ⇒ 实测 69.3% **门红 EXIT=1**，还原 ⇒ 79.3% EXIT=0；`AGENTS.md` §2 已写明阈值与来源；原文与施工记录移入归档。此前： 2026-11-13 五角审查档 4 立单（只立单不施工）——新立 `CT-COVERAGE-FLOOR`；`X-OPENAPI-INGEST` 原文逐字移入归档；此前：2026-11-12 `X-OPENAPI-INGEST` 作废登记（主任裁决 A：维持 core 唯一真相、本仓零 spec；不计完成，理由见该条）；2026-11-12 `X-LOW-CLEANUP` contracts 腿——`fees/ashare.go` 三处指向已删 `tellurion:pkg/live` 的死示例注释清理（`pkg/live` 与 `mock_trader.go` 实测不存在），只改注释、不动字段与导出面；2026-11-12 跨仓伞两条腿——`X-CI-RACE-STATICCHECK`（CI 补 `-race` 与 `staticcheck` v0.8.2；存量 1 条 `domain/market/types_test.go:325` SA4023 恒假比较改为断言动态类型；`-race` 用 Docker golang:1.26 全量跑通 EXIT=0）、`X-DEP-BUMP-CONSUMERS`（本仓为契约生产方，无 kit/pgx pin 可抬；只补 `govulncheck` 门并把 CI `go-version` 升到 1.26.11，`go` directive 保持 1.25）；2026-11-11 DDD 边界审查确认 `contracts/risk` 为风控快照形状的**唯一真相**（ADR-040）：本仓只承载形状/编码/四个跨仓哨兵，`Bundle`/`ParseBundle`/`ForScope` 属 risk 服务策略、明确不迁入；采用方迁移登记在 risk 仓；CT-06-ALIAS-EXIT 收口（2026-11-11：删除 `domain/` 根包 6 个 market 类型别名，11 仓重跑实测经根包别名的真实残留 0；本仓 `go build/vet/test ./... -count=1` + `gofmt -l` 全绿，`check_api_surface` 387/387、`check_surface_freeze` 380 = 380、`check_contract_shape` / `check_release_matrix` 0；**已发版 `v0.10.0`**（2026-11-11，release commit `3a3b9de`，annotated tag `v0.10.0` → `3a3b9de` 已 push；tag 后四护栏复跑 `10 / 11 / 11` 全绿），下游 pin 链 data → agent → core 与 broker / orchestration / strategy 六仓已升 `v0.10.0` 并 push；原文与实测输出归档在 `docs/archive/tasks-archive-2026-11.md` 末节）；CT-06 收口（2026-11-10：CT-06-A/B 已执行；余下 6 个 market 别名按「别名不是重复真相」豁免留仓，**代裁·待主任追认**；跨 11 仓实测引用：`OHLCV` core+agent、`Fundamental` core+data，`Stock`/`IndexDaily`/`Split`/`Dividend` 仅 core；本仓 `go build/vet/test ./... -count=1` 与 `check_api_surface` / `check_surface_freeze` / `check_contract_shape` / `check_release_matrix` 全 0）；CT-05 收编完成（2026-11-09：`tellurion-broker` 删本地 DTO 改引 `contracts/broker` 类型别名 + 13 组报文改前/改后逐字节 golden；broker build/vet/test 全 0）+ CT-04 编排 wire 契约腿（2026-11-07：新建 `run/` 包 —— `ClaimView` 七字段 / `ProvenanceView` 十五字段 / 父作业提交请求族 + `ChildIdempotencyKey` 纯函数；三份形状逐字节对拍编排仓 + `conformance/run_cases.json` 七用例；`contractRoots` 五根 → 六根；导出面 418 → 427、基线 408 → 417；本地 tag `v0.8.0` **未 push**，三仓改类型与 pin **未执行**；§4 口径代裁为 A、第六轮 #9 解除；施工记录见「一、本仓任务」）+ CT-03 契约腿（2026-11-07：新建 `risk/` 包 —— 快照形状 + 冻结词表 + `Snapshot.Hash()`/`Verify()`/`Clone()` + 4 哨兵；编码逐字节搬运、固定向量 `14ecab86…` + 跨仓对拍腿；`contractRoots` 四根 → 五根；导出面 465 → 487、基线 448 → 470；本地 tag `v0.7.0` **未 push**，三仓 pin **未执行**；施工记录见「一、本仓任务」）+ A 组三个形状包决策材料（2026-10-26：只出材料、0 代码、三单未立；六条派工前提逐条复现 + 破坏验证实跑；施工记录见「一、本仓任务」）+ QT-40 X1a 契约腿（2026-10-24：新增 `dataprotocol/` 包 —— `UnifiedDataPoint` 七字段（含 `VisibleAt`）+ `ToJSON`/`DeduplicateKey` + 15 个 `DataType*` 常量；无构造器、不读时钟；`v0.6.0` 已打 tag **且已 push**（2026-10-26 实测 `git ls-remote --tags origin` 有 `refs/tags/v0.6.0` → peel `5ee7fae`；此前记的「待 push」是 2026-10-24/25 的旧账。⚠️ 但本仓 `main` 领先 `origin/main` **1 笔**（`3055b6e` 未推）——tag 在远端 ≠ 分支已同步）；施工记录见「一、本仓任务」）+ G4 执行组（2026-10-24：core / agent / data / strategy 四仓升 pin `v0.5.0`；`check_contract_pin` GREEN 6 pin、conformance 2/2 零跳过；施工记录见「一、本仓任务」）+ G2 执行组（2026-10-24：QT-22 契约腿 `StrategyVersionInput` 六项 + `DomainStrategyVersionHash`，`v0.5.0` 已发版；施工记录见「一、本仓任务」）+ G0 执行组（2026-10-24：QT-20 契约腿 + QT-41b ④ 错误族，`v0.4.0` 已发；施工记录见「一、本仓任务」）+ G0D 裁决材料（QT-22 / QT-40 X1a / 服务 DTO 三份，只出材料、0 代码）；QT-09 本仓侧交接（2026-10-23：② examples/ + ③ conformance 语料与契约侧自检腿 + ④ PR 模板/check_pr_body + ⑤ check_contract_shape + surface 冻结 check_surface_freeze/api-surface.txt，全部接进 CI；三条工具另补表驱动单测）；QT-05 ③-① 本仓侧收官（四包 429/429 + §8 #1~#3 裁决 + 下游腿交接单）+ 自审订正（ADR 文件名死链、方法类别、EngineRunner 返回值、fees 注释失真；新增 ADR-0004）；CT-02 结案归档（版本规则表 + 兼容矩阵 + CHANGELOG + check_release_matrix 进 CI）；推送回写（2026-10-22：6 个 commit 已进 `origin/main`，该批 HEAD = `fbe9a2d`，`go get` 实测可解析，CI 429/429；§8 #4 登记为 X1a 待办）；**2026-10-25 B6-1**：QT-07 ⑥ S3 结案回写——agent 删掉 `replace tellurion => ../tellurion`（零 import / 零 require / CI 不 checkout core），core 侧同批清空 `check_contract_pin.py` 的白名单与 `check_repo_autonomy.py` 的 R1 死登记；本仓 AGENTS §3 的「唯一过渡例外」改为无豁免禁令、ADR-0001 附「例外收回记录」、QT-05 在册余项清零转 ✅
---

# TASKS — tellurion-contracts 未完成项

> **归档说明（2026-11-11）**：本文件已瘦身为「只装未做的事」。2026-11-10 及以前已收口（✅）的任务原文、施工记录、验收输出与破坏验证，已逐字移入 [`docs/archive/tasks-archive-2026-11.md`](archive/tasks-archive-2026-11.md)；该文件顶部有「检索索引」表（ID / 标题 / 收口日期 / commit / 证据链接）。检索：`grep -n "<任务ID或commit>" docs/archive/tasks-archive-2026-11.md`。
>
> 状态：`⬜ 待做` / `🔶 进行中` / `⛔ 阻塞（有未解除前置）` / `⏸ 挂起（外部依赖或已裁决暂不实施；**不计入开放项**，重开须另立单并满足该条「重开条件」）` / `❌ 作废（主任裁决不做或前提证伪；**不计完成、不计入开放项**，留裁决记录）` / `✅ 完成（附日期，随即移入归档）`

## 一、未做 / 挂起 / 外部依赖

> **活跃层无开放项**（2026-11-13：`CT-COVERAGE-FLOOR` 已收口 ✅，原文与施工记录见归档「2026-11-13 收口」节）。`X-OPENAPI-INGEST` 已按主任裁决 A **作废（不计完成）并逐字移入归档**（索引见 [`docs/archive/tasks-archive-2026-11.md`](archive/tasks-archive-2026-11.md) 顶部表，原文在末节）；`CT-06-ALIAS-EXIT`（2026-11-11：删除 6 个 `domain` 根别名，随 `v0.10.0` 发版）的原文与实测输出亦在归档末节。

## 二、跨仓指针（不在本仓，勿双写）

| 事项 | 去哪 |
|---|---|
| core 侧执行闸门与消费面 | `tellurion:docs/TASKS.md` |
| data 侧参考数据形状与零写者表 | `tellurion-data:docs/TASKS.md` |
| agent / strategy 侧契约腿 | 各自 `docs/TASKS.md` |
| HTTP wire 契约（OpenAPI） | `tellurion:docs/openapi.yaml`（唯一真相，2026-11-12 主任裁决 A）；agent 组以 `tellurion-agent:docs/api-schema.md` 为准；本仓不收编（`X-OPENAPI-INGEST` 已作废） |
| `contracts/risk` 的采用方迁移（risk 仓改别名引用；`Bundle`/`ParseBundle` **不进本仓**） | `tellurion-risk:docs/TASKS.md` 的 `RISK-CONTRACTS-SNAPSHOT-MIGRATION`；裁决见跨仓 ADR-040 |

## 维护本文件

- **新增单**：写全六字段；六字段里必须显式引用**冻结契约的路径与节号**（契约仓的包路径、设计文档的节号），不许写"参考相关实现"。
- **派工前先证伪前提**：读码核对工单写的前提；不符就停下回报 + 给拆法，不硬做、不扩范围。
- **跨仓引用写 `repo:path`**（如 `tellurion:docs/TASKS.md`），不做可点击相对链接。
- **顺序纪律**：同一波次内撞同一批文件的单按单号串行，一单一 commit。
