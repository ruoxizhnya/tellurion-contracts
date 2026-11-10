---
status: active
last-verified: 2026-11-11
verified-by: CT-06-ALIAS-EXIT 收口（2026-11-11：删除 `domain/` 根包 6 个 market 类型别名，11 仓重跑实测经根包别名的真实残留 0；本仓 `go build/vet/test ./... -count=1` + `gofmt -l` 全绿，`check_api_surface` 387/387、`check_surface_freeze` 380 = 380、`check_contract_shape` / `check_release_matrix` 0；**未发版**——只记 `CHANGELOG [Unreleased]`，无 tag、不 push，下游仍 pin `v0.9.0`；原文与实测输出归档在 `docs/archive/tasks-archive-2026-11.md` 末节）；CT-06 收口（2026-11-10：CT-06-A/B 已执行；余下 6 个 market 别名按「别名不是重复真相」豁免留仓，**代裁·待主任追认**；跨 11 仓实测引用：`OHLCV` core+agent、`Fundamental` core+data，`Stock`/`IndexDaily`/`Split`/`Dividend` 仅 core；本仓 `go build/vet/test ./... -count=1` 与 `check_api_surface` / `check_surface_freeze` / `check_contract_shape` / `check_release_matrix` 全 0）；CT-05 收编完成（2026-11-09：`tellurion-broker` 删本地 DTO 改引 `contracts/broker` 类型别名 + 13 组报文改前/改后逐字节 golden；broker build/vet/test 全 0）+ CT-04 编排 wire 契约腿（2026-11-07：新建 `run/` 包 —— `ClaimView` 七字段 / `ProvenanceView` 十五字段 / 父作业提交请求族 + `ChildIdempotencyKey` 纯函数；三份形状逐字节对拍编排仓 + `conformance/run_cases.json` 七用例；`contractRoots` 五根 → 六根；导出面 418 → 427、基线 408 → 417；本地 tag `v0.8.0` **未 push**，三仓改类型与 pin **未执行**；§4 口径代裁为 A、第六轮 #9 解除；施工记录见「一、本仓任务」）+ CT-03 契约腿（2026-11-07：新建 `risk/` 包 —— 快照形状 + 冻结词表 + `Snapshot.Hash()`/`Verify()`/`Clone()` + 4 哨兵；编码逐字节搬运、固定向量 `14ecab86…` + 跨仓对拍腿；`contractRoots` 四根 → 五根；导出面 465 → 487、基线 448 → 470；本地 tag `v0.7.0` **未 push**，三仓 pin **未执行**；施工记录见「一、本仓任务」）+ A 组三个形状包决策材料（2026-10-26：只出材料、0 代码、三单未立；六条派工前提逐条复现 + 破坏验证实跑；施工记录见「一、本仓任务」）+ QT-40 X1a 契约腿（2026-10-24：新增 `dataprotocol/` 包 —— `UnifiedDataPoint` 七字段（含 `VisibleAt`）+ `ToJSON`/`DeduplicateKey` + 15 个 `DataType*` 常量；无构造器、不读时钟；`v0.6.0` 已打 tag **且已 push**（2026-10-26 实测 `git ls-remote --tags origin` 有 `refs/tags/v0.6.0` → peel `5ee7fae`；此前记的「待 push」是 2026-10-24/25 的旧账。⚠️ 但本仓 `main` 领先 `origin/main` **1 笔**（`3055b6e` 未推）——tag 在远端 ≠ 分支已同步）；施工记录见「一、本仓任务」）+ G4 执行组（2026-10-24：core / agent / data / strategy 四仓升 pin `v0.5.0`；`check_contract_pin` GREEN 6 pin、conformance 2/2 零跳过；施工记录见「一、本仓任务」）+ G2 执行组（2026-10-24：QT-22 契约腿 `StrategyVersionInput` 六项 + `DomainStrategyVersionHash`，`v0.5.0` 已发版；施工记录见「一、本仓任务」）+ G0 执行组（2026-10-24：QT-20 契约腿 + QT-41b ④ 错误族，`v0.4.0` 已发；施工记录见「一、本仓任务」）+ G0D 裁决材料（QT-22 / QT-40 X1a / 服务 DTO 三份，只出材料、0 代码）；QT-09 本仓侧交接（2026-10-23：② examples/ + ③ conformance 语料与契约侧自检腿 + ④ PR 模板/check_pr_body + ⑤ check_contract_shape + surface 冻结 check_surface_freeze/api-surface.txt，全部接进 CI；三条工具另补表驱动单测）；QT-05 ③-① 本仓侧收官（四包 429/429 + §8 #1~#3 裁决 + 下游腿交接单）+ 自审订正（ADR 文件名死链、方法类别、EngineRunner 返回值、fees 注释失真；新增 ADR-0004）；CT-02 结案归档（版本规则表 + 兼容矩阵 + CHANGELOG + check_release_matrix 进 CI）；推送回写（2026-10-22：6 个 commit 已进 `origin/main`，该批 HEAD = `fbe9a2d`，`go get` 实测可解析，CI 429/429；§8 #4 登记为 X1a 待办）；**2026-10-25 B6-1**：QT-07 ⑥ S3 结案回写——agent 删掉 `replace tellurion => ../tellurion`（零 import / 零 require / CI 不 checkout core），core 侧同批清空 `check_contract_pin.py` 的白名单与 `check_repo_autonomy.py` 的 R1 死登记；本仓 AGENTS §3 的「唯一过渡例外」改为无豁免禁令、ADR-0001 附「例外收回记录」、QT-05 在册余项清零转 ✅
---

# TASKS — tellurion-contracts 未完成项

> **归档说明（2026-11-11）**：本文件已瘦身为「只装未做的事」。2026-11-10 及以前已收口（✅）的任务原文、施工记录、验收输出与破坏验证，已逐字移入 [`docs/archive/tasks-archive-2026-11.md`](archive/tasks-archive-2026-11.md)；该文件顶部有「检索索引」表（ID / 标题 / 收口日期 / commit / 证据链接）。检索：`grep -n "<任务ID或commit>" docs/archive/tasks-archive-2026-11.md`。
>
> 状态：`⬜ 待做` / `🔶 进行中` / `⛔ 阻塞（有未解除前置）` / `⏸ 挂起（外部依赖或已裁决暂不实施；**不计入开放项**，重开须另立单并满足该条「重开条件」）` / `✅ 完成（附日期，随即移入归档）`

## 一、未做 / 挂起 / 外部依赖

> 本仓当前**无开放项**。最近一笔收口 `CT-06-ALIAS-EXIT`（2026-11-11：删除 6 个 `domain` 根别名，未发版）的原文与实测输出在 [`docs/archive/tasks-archive-2026-11.md`](archive/tasks-archive-2026-11.md) 末节，索引见该文件顶部表。

## 二、跨仓指针（不在本仓，勿双写）

| 事项 | 去哪 |
|---|---|
| core 侧执行闸门与消费面 | `tellurion:docs/TASKS.md` |
| data 侧参考数据形状与零写者表 | `tellurion-data:docs/TASKS.md` |
| agent / strategy 侧契约腿 | 各自 `docs/TASKS.md` |

## 维护本文件

- **新增单**：写全六字段；六字段里必须显式引用**冻结契约的路径与节号**（契约仓的包路径、设计文档的节号），不许写"参考相关实现"。
- **派工前先证伪前提**：读码核对工单写的前提；不符就停下回报 + 给拆法，不硬做、不扩范围。
- **跨仓引用写 `repo:path`**（如 `tellurion:docs/TASKS.md`），不做可点击相对链接。
- **顺序纪律**：同一波次内撞同一批文件的单按单号串行，一单一 commit。
