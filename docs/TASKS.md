---
status: active
last-verified: 2026-10-20
verified-by: 首次建立本仓台账（QT-08 十二仓自治三件套）
---

# TASKS — tellurion-contracts 未完成项

> 本仓**唯一**任务台账。跨仓任务只在 §二 留指针，**不镜像正文**（双写必然漂移）。
> 状态：`⬜ 待做` / `🔶 进行中` / `⛔ 阻塞（有未解除前置）` / `✅ 完成（附日期）`
> 每张单六字段（① 目标 ② 上下文 ③ 要求 ④ 约束与非目标 ⑤ 验收 ⑥ 边界），判据见 `tellurion:AGENTS.md` §8.5。
> 完成即标 ✅（附日期）并移入 `docs/archive/TASKS-YYYY-MM-DD-archived.md`；本文件只留未完成项。

---

## 一、本仓任务

> **占位仓**：实现落地前不写业务代码；以下两张单是**随仓文档交付**（refactor-plan-review M6/L4/L5 修复传播，2026-10-20 登记），与 S1 落码（QT-05）并行可做。
>
> 第一张单长什么样（照此写，别自造格式）：
>
> ### XX-01 · <一句话目标> [P?] ⬜
> - **① 目标**：…
> - **② 上下文**：本仓路径/行号 + 冻结契约的路径与节号
> - **③ 要求**：①②③（可执行、可核对）
> - **④ 约束与非目标**：明确不做什么
> - **⑤ 验收**：可复现的命令 + 期望输出；至少一条"破坏验证"（改坏它应当红）
> - **⑥ 边界**：只做什么；相邻的什么归别的单

### CT-01 · 随仓接口规格载体 `docs/API.md`（导出面清单 + 稳定性 + 错误语义）[P1] ⬜

- **① 目标**：本仓导出面有随仓、可自证的接口规格，下游（core/agent）不读 core 实现即可知调用形状与错误语义。
- **② 上下文**：refactor-plan-review **M6**（报告：`tellurion:docs/design/qt-decoupling/refactor-plan-review-2026-10-20.md`）；上游编排单 `tellurion:docs/TASKS.md` QT-07 ③-S2-c；冻结面指回 `tellurion:docs/SPEC.md` 四契约节（改冻结面走 `tellurion:docs/TASKS.md`「维护本文件」的冻结 schema 变更评审硬规则）。现状导出面清单散在 `tellurion:docs/design/qt-decoupling/kernel-interface-map.md`，无随仓载体。
- **③ 要求**：① 建 `docs/API.md`：逐项列导出符号——四契约（`engine/`）+ port + 值对象（`domain/`）+ 哨兵错误族，每项标**稳定性标注**（frozen / stable / experimental）与**错误语义**（`errors.Is` 判定、禁字符串匹配）；② 每个冻结面指回 `tellurion:docs/SPEC.md` 对应**节号**（只指节号，不复制正文）；③ 护栏进 CI（AST 扫描）：新增导出符号未登记进 `docs/API.md` → 红。
- **④ 约束与非目标**：不写内部函数/实现细节；**不做手写 DTO**——跨进程 DTO 是 `tellurion:docs/openapi.yaml` 的生成物、归 core（见 `tellurion:docs/design/qt-decoupling/repo-inventory.md` §3），本仓 `dto/` 只收生成物且禁手改。
- **⑤ 验收**：`go doc ./...` 导出面与 `docs/API.md` 清单一一对应（比对脚本跑通）；**破坏验证**：故意导出一个未登记符号 → CI 红（留证后撤回）。
- **⑥ 边界**：只做接口规格载体与护栏；版本规则/兼容矩阵是 CT-02。

### CT-02 · README 版本规则表 + 兼容矩阵 + CHANGELOG + `GOPRIVATE` 前提 [P1] ⬜

- **① 目标**：本仓作为被依赖库有成文的版本纪律与拉取前提，breaking 有账可查。
- **② 上下文**：refactor-plan-review **L4/L5**；规则口径冻结在 `tellurion:docs/TASKS.md` QT-19 ③（**规则内容以彼处为准，本单只落本仓文档**）；`tellurion:docs/design/qt-decoupling/repo-inventory.md` §1.1/§5 的可见性与 GOPRIVATE 风险注。
- **③ 要求**：① README 写**版本规则表**（`MAJOR` = 删/改签名或语义；`MINOR` = 新增兼容能力；`PATCH` = 文档/无行为修正）；② **兼容矩阵**：`vX.Y` ← 支持 core / agent 的哪些版本，每次发版更新；③ CHANGELOG 至少一条记录；④ AGENTS.md 写一行 **`GOPRIVATE=github.com/ruoxizhnya/*` + git 凭证**前提（本仓转私有时 CI/构建必须配，否则 `go get` 走 proxy 必然失败）。
- **④ 约束与非目标**：不引 release 工具链（goreleaser 等）；内容寻址标识（`strategy_ref`/`artifact_hash`/`run_spec_hash`）不参与 semver；0.x 阶段允许 breaking 但须 bump MINOR + CHANGELOG 条目（细则以 QT-19 为准）。
- **⑤ 验收**：README 存在版本规则表与兼容矩阵；CHANGELOG ≥ 1 条；AGENTS.md 含 GOPRIVATE 行；**破坏验证**：发一版不更新兼容矩阵 → 发版 checklist 项不可勾（或检查脚本红）。
- **⑥ 边界**：只定本仓策略与文档；kit 的同名交付是 `tellurion-kit:docs/TASKS.md` KT-02；「禁删 deprecated」机械护栏属 `tellurion:docs/TASKS.md` QT-09 ⑤。


---

## 施工序（候选单，开工前必须展开成六字段）

> 这张表是**施工顺序**，不是任务台账：每条在执行前要展开成六字段单（目标/上下文/要求/约束/验收/边界），并写明冻结契约的路径与节号。
> 顺序有依赖，**不要跳级**：跳序的代价是返工（例：X1a 的协议没冻结就写 aggregator，等于把协议写死在一个服务里）。

| 阶段 | 内容 | 上游单号 | 前置 |
|---|---|---|---|
| S1 | 把 `pkg/domain`（17 文件）+ `pkg/backtest/contracts`（6 文件）搬入本仓；core/agent 改 `require` 伪版本 | `QT-05` | W0 全绿 |
| S2 | 打 `v0.1.0`；下游把伪版本换 tag | `QT-07` | S1 全绿 |
| S3 | 版本规则表 + 兼容矩阵 + CHANGELOG | `QT-19` | S2 打 tag 后 |

## 二、跨仓指针（不在本仓，勿双写）

| 单号 | 主题 | 在哪 |
|---|---|---|
| QT-05 | S1：把 `pkg/domain` + `pkg/backtest/contracts` 搬进本仓 | `tellurion:docs/TASKS.md` |
| QT-07 | S2：打 `v0.1.0` 并把下游伪版本换成 tag | `tellurion:docs/TASKS.md` |
| QT-19 | 版本策略 / 兼容矩阵 / CHANGELOG 规则 | `tellurion:docs/TASKS.md` |

---

## 维护本文件

- **新增单**：写全六字段；六字段里必须显式引用**冻结契约的路径与节号**（契约仓的包路径、设计文档的节号），不许写"参考相关实现"。
- **派工前先证伪前提**：读码核对工单写的前提；不符就停下回报 + 给拆法，不硬做、不扩范围。
- **跨仓引用写 `repo:path`**（如 `tellurion:docs/TASKS.md`），不做可点击相对链接。
- **顺序纪律**：同一波次内撞同一批文件的单按单号串行，一单一 commit。
