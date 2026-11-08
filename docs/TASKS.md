---
status: active
last-verified: 2026-11-07
verified-by: CT-04 编排 wire 契约腿（2026-11-07：新建 `run/` 包 —— `ClaimView` 七字段 / `ProvenanceView` 十五字段 / 父作业提交请求族 + `ChildIdempotencyKey` 纯函数；三份形状逐字节对拍编排仓 + `conformance/run_cases.json` 七用例；`contractRoots` 五根 → 六根；导出面 418 → 427、基线 408 → 417；本地 tag `v0.8.0` **未 push**，三仓改类型与 pin **未执行**；§4 口径代裁为 A、第六轮 #9 解除；施工记录见「一、本仓任务」）+ CT-03 契约腿（2026-11-07：新建 `risk/` 包 —— 快照形状 + 冻结词表 + `Snapshot.Hash()`/`Verify()`/`Clone()` + 4 哨兵；编码逐字节搬运、固定向量 `14ecab86…` + 跨仓对拍腿；`contractRoots` 四根 → 五根；导出面 465 → 487、基线 448 → 470；本地 tag `v0.7.0` **未 push**，三仓 pin **未执行**；施工记录见「一、本仓任务」）+ A 组三个形状包决策材料（2026-10-26：只出材料、0 代码、三单未立；六条派工前提逐条复现 + 破坏验证实跑；施工记录见「一、本仓任务」）+ QT-40 X1a 契约腿（2026-10-24：新增 `dataprotocol/` 包 —— `UnifiedDataPoint` 七字段（含 `VisibleAt`）+ `ToJSON`/`DeduplicateKey` + 15 个 `DataType*` 常量；无构造器、不读时钟；`v0.6.0` 已打 tag **且已 push**（2026-10-26 实测 `git ls-remote --tags origin` 有 `refs/tags/v0.6.0` → peel `5ee7fae`；此前记的「待 push」是 2026-10-24/25 的旧账。⚠️ 但本仓 `main` 领先 `origin/main` **1 笔**（`3055b6e` 未推）——tag 在远端 ≠ 分支已同步）；施工记录见「一、本仓任务」）+ G4 执行组（2026-10-24：core / agent / data / strategy 四仓升 pin `v0.5.0`；`check_contract_pin` GREEN 6 pin、conformance 2/2 零跳过；施工记录见「一、本仓任务」）+ G2 执行组（2026-10-24：QT-22 契约腿 `StrategyVersionInput` 六项 + `DomainStrategyVersionHash`，`v0.5.0` 已发版；施工记录见「一、本仓任务」）+ G0 执行组（2026-10-24：QT-20 契约腿 + QT-41b ④ 错误族，`v0.4.0` 已发；施工记录见「一、本仓任务」）+ G0D 裁决材料（QT-22 / QT-40 X1a / 服务 DTO 三份，只出材料、0 代码）；QT-09 本仓侧交接（2026-10-23：② examples/ + ③ conformance 语料与契约侧自检腿 + ④ PR 模板/check_pr_body + ⑤ check_contract_shape + surface 冻结 check_surface_freeze/api-surface.txt，全部接进 CI；三条工具另补表驱动单测）；QT-05 ③-① 本仓侧收官（四包 429/429 + §8 #1~#3 裁决 + 下游腿交接单）+ 自审订正（ADR 文件名死链、方法类别、EngineRunner 返回值、fees 注释失真；新增 ADR-0004）；CT-02 结案归档（版本规则表 + 兼容矩阵 + CHANGELOG + check_release_matrix 进 CI）；推送回写（2026-10-22：6 个 commit 已进 `origin/main`，该批 HEAD = `fbe9a2d`，`go get` 实测可解析，CI 429/429；§8 #4 登记为 X1a 待办）；**2026-10-25 B6-1**：QT-07 ⑥ S3 结案回写——agent 删掉 `replace tellurion => ../tellurion`（零 import / 零 require / CI 不 checkout core），core 侧同批清空 `check_contract_pin.py` 的白名单与 `check_repo_autonomy.py` 的 R1 死登记；本仓 AGENTS §3 的「唯一过渡例外」改为无豁免禁令、ADR-0001 附「例外收回记录」、QT-05 在册余项清零转 ✅
---

# TASKS — tellurion-contracts 未完成项

> 本仓**唯一**任务台账。跨仓任务只在 §二 留指针，**不镜像正文**（双写必然漂移）。
> 状态：`⬜ 待做` / `🔶 进行中` / `⛔ 阻塞（有未解除前置）` / `✅ 完成（附日期）`

> **第三轮裁决指针（2026-10-25）**：② risk 快照迁 contracts = **新建 `risk` 包**，迁 `Snapshot` / `SnapshotRule` / `Scope` / `RuleType` + 冻结词表 + hash 纯函数 + 4 个哨兵；`Bundle` / `ParseBundle` / `ForScope` 不迁；**编码逐字节搬运，禁止换 Canonical**；必须带固定向量 + 与 risk 仓的对拍腿。建议发 MINOR（`v0.7.0`）。正文见 `tellurion:docs/design/qt-decoupling/cross-repo-decision-record-2026-10-25.md` §五。
> 每张单六字段（① 目标 ② 上下文 ③ 要求 ④ 约束与非目标 ⑤ 验收 ⑥ 边界），判据见 `tellurion:AGENTS.md` §8.5。
> 完成即标 ✅（附日期）并移入 `docs/archive/TASKS-YYYY-MM-DD-archived.md`；本文件只留未完成项。

---

## 一、本仓任务

> **不再是占位仓（2026-10-22 起）**：QT-05 ③-① **本仓侧已完成**——`fees/`（13）、`domain/market/`（46）、`domain/`（176）、`engine/`（194）四包落地，`docs/API.md` §4 已登记、两条护栏绿（429/429）；CT-01/CT-02 两张**随仓文档交付**（refactor-plan-review M6/L4/L5 修复传播，2026-10-20 登记）均已结案。
>
> ✅ **CT-01 已结案（2026-10-21）**：`docs/API.md` + `tools/check_api_surface`（CI 步骤「导出面核对」）落地，整条搬进
> [`archive/TASKS-2026-10-21-archived.md`](archive/TASKS-2026-10-21-archived.md)（含验收逐条与实测偏差）。
> **对 QT-05 的直接后果**：本仓落第一个符号那天起，`docs/API.md` §4 的登记区就必须同步填，否则
> `go run ./tools/check_api_surface` 红——登记形状与生成命令都在那一节。
>
> ✅ **CT-02 已结案（2026-10-22）**：README 的版本规则表 + 兼容矩阵 + `CHANGELOG.md` + AGENTS 的 `GOPRIVATE` 前提落地，
> 护栏 `tools/check_release_matrix` 进 CI 步骤「版本纪律核对」，整条搬进
> [`archive/TASKS-2026-10-22-archived.md`](archive/TASKS-2026-10-22-archived.md)。
> **对 QT-05 / QT-07 的直接后果**：S2 打第一个 tag（`v0.1.0`）那天起，README 兼容矩阵必须有 `v0.1` 那一行、
> CHANGELOG 必须有 `## [v0.1.0] - 日期` 那一节，漏一条 `go run ./tools/check_release_matrix` 红（CI 同款）。
> tag 那条腿今天空跑（本仓 0 tag），从第一个 tag 开始才咬人。
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

### CT-01 · 随仓接口规格载体 `docs/API.md` ✅ 完成（2026-10-21，已归档）

> 六字段正文与结案记录（验收逐条、破坏验证的真退出码、五条实测偏差登记）在
> [`archive/TASKS-2026-10-21-archived.md`](archive/TASKS-2026-10-21-archived.md)。本文件按规则只留未完成项，这里只留指针。

### CT-02 · README 版本规则表 + 兼容矩阵 + CHANGELOG + `GOPRIVATE` 前提 [P1] ✅ 完成（2026-10-22，已归档）

> 六字段正文与结案记录（③ 逐条落点、⑤ 验收真退出码、现场破坏验证的还原凭据、三条口径偏差登记）在
> [`archive/TASKS-2026-10-22-archived.md`](archive/TASKS-2026-10-22-archived.md)。本文件按规则只留未完成项，这里只留指针。
> **对 QT-05 / QT-07 的直接后果**：S2 打第一个 tag（`v0.1.0`）那天起，`README.md` 的兼容矩阵必须加 `v0.1` 那一行、
> `CHANGELOG.md` 必须把 `[Unreleased]` 的内容搬进 `## [v0.1.0] - 日期` 节，否则 `go run ./tools/check_release_matrix` 红
> （CI「版本纪律核对」同款）。tag 那条腿今天空跑（0 tag），从第一个 tag 开始才咬人。

### QT-20 契约腿 · `BacktestRequest.StrategyRef/StrategyParams` + `BacktestParams.Params` [P0] ✅ 完成（2026-10-24）

> 由 G0 执行组派工，正文与验收逐条见本节的施工记录；core 侧消费腿（G1-QT-20）不在本仓。

- **① 目标**：core 的 QT-20 参数投递有**入口字段**可消费——请求能指定策略引用与参数覆写，生效参数能落到 `BacktestParams`；形状错误 400，不把 manifest 业务规则搬进契约仓。
- **② 上下文**：`engine/contracts.go:170-215`（v0.3.0 无这两字段）、`engine/provenance.go:249-312`（`StrategyRef` / `ParseStrategyRef` / `String()` 是引用形状的唯一实现）、core `pkg/strategy/signal_provider.go:52-65`（`RunMeta.Params` 是字段名对齐点）、`tellurion:docs/design/qt-decoupling/p0-remediation-design.md` A3 / D26、`w1-repo-ownership.md` J2。
- **③ 要求（全部落地）**：① `BacktestRequest.StrategyRef string`，wire 形态与 `ParseStrategyRef`/`String()` 互逆（`id` / `id@hash`），不新增第二种引用形状；② `BacktestRequest.StrategyParams map[string]any` 只收浅层对象，嵌套键拒绝；③ `domain.BacktestParams.Params map[string]any` 承载生效参数；④ `ValidateRunSpec` 只做形状 / 类型校验，**未知键清单不由本仓提供**（core manifest）。
- **④ 约束与非目标**：不改 `SignalWeight` / `ResultSpec` 冻结面；不新增 canonical 编码器；不在本仓实现 manifest / 合并 / 取值范围；不发 major tag。**未做**：core 的入口接线与 `RunMeta.Params` 投递（G1-QT-20）。
- **⑤ 验收（真退出码）**：`go build ./... && go vet ./... && go test ./... -count=1` → **EXIT=0**（12 个包全绿）；`go run ./tools/check_api_surface` → 443/443；`go run ./tools/check_surface_freeze` → 426 行一致；`go run ./tools/check_contract_shape` → 绿；`go run ./tools/check_release_matrix` → **EXIT=0**。
  - **破坏验证**：删掉 `validateStrategyParams` 里的 `isScalarValue` 拒收分支 → `go test ./engine -run TestRunSpec_ValidateRejectsNestedStrategyParams` **EXIT=1**（4 条子用例全部点名失败）；逐字节还原后 `sha256=4d190be7b2b2809c67871f6c052bc560f362f6a109e1aa7aed1522b1b44e7f3e` 与改坏前一致，全套测试回到 EXIT=0。
- **⑥ 边界**：只改契约形状与形状校验；请求字段的消费、未知键拒绝、`RunMeta.Params` 接线归 core G1-QT-20。版本取 `v0.4.0`（新增导出符号 = MINOR，README 规则表）。

### QT-41b ④ 错误族契约腿 · 编译 / 资产四哨兵 [P1] ✅ 完成（2026-10-24）

- **① 目标**：跨仓消费方（agent 生成闭环 / core 取 artifact / orchestration 终态）能 `errors.Is` 识别编译与资产管线的四类失败，不必 import 资产服务实现包。
- **② 上下文**：实测来源 `tellurion-strategy:internal/asset/errors.go:29-60`（`ErrSourceTooLarge` / `ErrArtifactTooLarge` / `ErrCompileUnavailable` / `ErrCompileFailed` 已存在）；同类先例 `engine/execution_errors.go`（QT-10）。
- **③ 要求（全部落地）**：四哨兵迁入 `engine/execution_errors.go`，**名字与语义（含错误文案）不改**；迁入判据 = 跨仓消费方必须能识别；**资产服务私有哨兵不迁**（`ErrStrategyNotFound` / `ErrStrategyExists` / `ErrVersionNotFound` / `ErrInvalidTransition` / `ErrJobNotFound` 等）。
- **④ 约束与非目标**：不迁私有哨兵；不改四条的文案 / 身份；不在本仓实现编译或资产逻辑；不改 strategy 仓（其别名切换归 G13 通知腿）。
- **⑤ 验收**：`TestCompileFamilySentinelsStayDistinctAndSurviveWrapping` 钉住四条互不命中 + 包装后仍可判定；全套 `go test ./... -count=1` EXIT=0；导出面 439 → 443。
- **⑥ 边界**：只落哨兵；compile 服务实现、asset 状态机与 wire code 映射留在 `tellurion-strategy`。

### QT-40 X1a 契约腿 · `dataprotocol` 包（`UnifiedDataPoint` 收编）[P1] ✅ 完成（2026-10-24，`v0.6.0`）

- **① 目标**：把已在三处漂移的 `UnifiedDataPoint` 收成唯一共享形状，让 data 与 core 都只 import 一份。
- **② 上下文**：裁决材料 [`design/qt40-x1a-dataprotocol-decision-material.md`](design/qt40-x1a-dataprotocol-decision-material.md) 实测**三份形状已在漂**（data 镜像有 `VisibleAt`、core source 有 `DeduplicateKey()`、core 构造器读墙钟而 data 注入时钟）。**2026-10-24 用户裁决走 (a) 先迁**。
- **③ 要求**：① 新增 `dataprotocol/` 包：`UnifiedDataPoint` 七字段（含 `VisibleAt`）+ `ToJSON` + `DeduplicateKey`；**不带构造器、不读时钟**；② `TimeSemantics` 不重复定义（已随 QT-05 住 `domain/market`）；③ `DataType` 先保持 `string` + 常量，具名化等 X1b schema 冻结；④ `Aggregator` 接口**不落**（今天无实现、无消费方）；⑤ 同批登记 `docs/API.md` §4 与 `docs/api-surface.txt`，按 MINOR 发版。
- **④ 约束与非目标**：不在契约包放 I/O 或时钟；不改 wire 语义；不替 data 服务裁决 aggregator 形状；不删 core / data 的副本（那属下游单）。
- **⑤ 验收**：`check_api_surface` / `check_surface_freeze` / `check_contract_shape` 全绿；`DeduplicateKey()` 与 core 现实现（`tellurion:pkg/data/source/unified.go:70-72`）向量同结果；`VisibleAt` 在类型里（缺则 data 迁移无意义）。
- **⑥ 边界**：只出契约类型与登记；core 删两份副本归 core 单，data 镜像改 import 归 data 单。
- **施工记录（2026-10-24）**：
  - 落码：`dataprotocol/dataprotocol.go` —— `UnifiedDataPoint` 七字段（字段顺序与 data 镜像一致，`VisibleAt` 末位）+ `ToJSON()` + `DeduplicateKey()`；15 个 `DataType*` 常量（与 core / data 现取值逐字一致）；**无构造器、无时钟**；`TimeSemantics` 未重复定义；`Aggregator` 未落。包注释写明三条边界（零 I/O / 不重复时间语义 / 不造空 port）。
  - 向量：`dataprotocol_test.go` 钉死 core 的既有向量（`AAPL|global_ohlcv|2026-05-14T09:30:00Z`）+ 纳秒精度 + 时区归 UTC + Source 不入键 + DataType 区分；`TestVisibleAt_RoundTrips` 证明 `VisibleAt` 在类型里且不与 `TradeTime` 混淆；`ToJSON` 三条（正常 / nil→`null` / 不可序列化返错）。
  - 登记：`docs/API.md` §4 新增 `dataprotocol/` 小节（18 行）、§5.5 记语义与迁移边界、§8 #4 落定为 (a)；`docs/api-surface.txt` 重生成（430 → 448 行）；`docs/ARCHITECTURE.md` §2/§3/§8 同步。
  - 验收（真退出码）：`go build ./...` / `go vet ./...` / `go test ./... -count=1` 全 **0**；`check_api_surface` → **465/465**（14 个代码包 / 6 个登记包）；`check_surface_freeze` → 448 行一致；`check_contract_shape` → 绿；`check_release_matrix` → 绿。
  - 破坏验证：把 `DeduplicateKey` 的 `TradeTime.UTC()` 改成 `TradeTime`（去掉时区归一）→ `go test ./dataprotocol -run TestDeduplicateKey_NormalizesToUTC` **EXIT=1**；逐字节还原后 `git diff` 为空、`go test ./dataprotocol` 回 **EXIT=0**。
  - **未随本单落（如实）**：core 删 `pkg/data/source/unified.go` / `pkg/storage/bulk_insert.go` 两份副本、data 镜像改 import、storage 列布局的对表测试 —— 都归各自仓的单；`tellurion:docs/design/qt-decoupling/data-service-contract-and-ownership.md` §4.1 注 ② 的回写归 core 侧协调者（本仓不改 core 文档）。在那之前三份形状仍并存。
  - ~~**已知护栏缺口（如实）**：`tools/check_contract_shape` 的 `contractRoots` 仍是 `domain/` / `engine/` / `fees/` 三根，**`dataprotocol/` 今天不在它的命名 / DB 红线上**。~~ **已关闭（2026-10-25，独立工具单）**：`contractRoots` 加入 `dataprotocol`（四根），夹具同步 —— 新增「dataprotocol 命名红」用例、「缺根红」期望改为 `fees/、dataprotocol/` 同时缺失，`greenFiles()` 与白名单用例补上第四根。`go test ./tools/check_contract_shape` 全绿；破坏验证见该单施工记录。
  - **护栏补洞施工记录（2026-10-25，独立工具单）**：`contractRoots` 三根 → 四根（+ `dataprotocol`），工具头部注释 / 成功文案 / `inContractScope` 注释与 `README.md` 的工具说明同步。验收（真退出码）：`gofmt -l tools/check_contract_shape` → 空；`go test ./tools/check_contract_shape -count=1` → **EXIT=0**（11 个子用例全 PASS）；实仓 `go run ./tools/check_contract_shape` → `14 个包 / 43 个非测试文件`、`✅ 四条契约根目录都在`、**EXIT=0**。破坏验证：把 `dataprotocol` 从 `contractRoots` 临时删掉 → 新增的「dataprotocol 命名红」用例报「这条红腿必须报错，却绿了」+「缺根红」咬不到 `dataprotocol/`，**EXIT=1**；逐字节还原后 **EXIT=0**，`git diff` 只剩本单四处改动。

### G0D · QT-22 / QT-40 X1a / 服务 DTO 三份裁决材料 [P0] ✅ 完成（2026-10-24）

> **只出材料，不写代码、不打 tag**。落点：`docs/design/`（导航见 [`design/README.md`](design/README.md)）。本组不直接改 broker / analytics / risk / data / core。

- **① 目标**：把三个跨仓落点问题的**证据与影响面**写成可拍板的材料，每条结论带 `repo:path:line`。
- **② 上下文**：`tellurion:docs/design/qt-decoupling/w1-repo-ownership.md:101`（version_hash 冻结 7 项）、
  `tellurion-contracts:docs/API.md:731`（dataprotocol 待裁 #4）、`tellurion:docs/design/qt-decoupling/data-service-contract-and-ownership.md:151-179`（A 组对表）、
  `tellurion-risk:docs/TASKS.md:94-96`（ParseBundle 第二使用方判据）。
- **③ 要求（全部落地）**：
  - [QT-22 材料](design/qt22-version-hash-decision-material.md)：逐项对表 7 vs 6，给出保留 / 删除 / 新增与 `v2` 代价。
  - [QT-40 X1a 材料](design/qt40-x1a-dataprotocol-decision-material.md)：`UnifiedDataPoint` 字段来源 + (a)/(b) 两份影响面 + 迁入判据。
  - [服务 DTO 材料](design/service-dto-migration-decision-material.md)：QT-45 / QT-43 / QT-42 的当前落点、第二消费方、迁入后谁改 import。
- **④ 约束与非目标**：不落 `StrategyVersionInput` / `dataprotocol` / 任何服务 DTO；不复用 `DomainParamsHash`；不改其他仓。
- **⑤ 验收**：三份材料均在 `docs/design/`；每条结论带 `repo:path:line`；**本单 0 个 `.go` diff**（提交只含 `docs/`，可用 `git show --stat` 复核）。
- **⑥ 边界**：只出材料。拍板与代码落点分别归 G1（core）、G4（strategy）、X1a 协调者、各服务仓协调者；分发对象 = G1 / G4 / G5 / G6 / G7。

### A 组 · 三个跨仓形状包决策材料（risk / 编排 / broker DTO）[P0] ✅ 完成（2026-10-26）

> **只出材料，不写代码、不打 tag、不立正式单**。落点：[`design/shape-packages-decision-material.md`](design/shape-packages-decision-material.md)（导航见 [`design/README.md`](design/README.md)）。

- **① 目标**：把主任要立的三单（`contracts/risk` / 编排三份 wire 形状 / broker DTO）的**前置、顺序、禁区、验收**写成可直接派发的材料，并核实派工给的六条事实。
- **② 上下文**：第三轮裁决 ②（`tellurion:docs/design/qt-decoupling/cross-repo-decision-record-2026-10-25.md:372`）、risk `QT-74`（`tellurion-risk:docs/TASKS.md:465`，⬜）、本仓 `docs/API.md:766`（wire DTO = openapi 生成物、禁手改）、`tellurion-orchestration:docs/ARCHITECTURE.md:100/174/185`（契约位）、`tellurion-broker:docs/ARCHITECTURE.md:152-167`（DTO 对表）。
- **③ 要求（全部落地）**：① §0.1 六条派工前提逐条复现（tag `v0.6.0` / 无 `risk/` / kit `v0.1.4` / 两仓干净且同步）；② §1~§3 逐项对表「迁什么 / 不迁什么」，每条带 `repo:path:line`；③ §4 点出**阻塞性冲突**、给两个选项与代价，并在 §4.4 落**四条组级禁区**（不得同批发三包 / 不得改 canonical 域标签 / 不得 `contracts → core` 反向依赖 / 不得顺手带 `Bundle`、`ParseBundle` 与私有哨兵，逐条带机械判据）；④ §5 出三张候选六字段单（建议编号 `CT-03`~`CT-05`）+ 逐单禁区；⑤ §8 如实登记未做项。
- **④ 约束与非目标**：**不**落任何包、**不**改 `check_contract_shape` 的 `contractRoots`、**不**动 `docs/API.md` §4 / `api-surface.txt` / 兼容矩阵 / CHANGELOG；**不**改 risk / orchestration / broker / core / analytics 任何文件；**不**代裁 §4 的 openapi 口径，也**不**替主任定单号。
- **⑤ 验收（真输出）**：`git -C tellurion-contracts tag | tail -1` → **`v0.6.0`**；`ls tellurion-contracts` → 无 `risk`（⇒ 证伪「risk 形状已进契约」）；两仓 `rev-list --count origin/main..HEAD` / `HEAD..origin/main` → **`0` / `0`**；**本单 0 个 `.go` diff**（`git show --stat` 只含 `docs/`）。
  - **破坏验证（前提守卫，实跑）**：临时建 `tellurion-contracts/risk/doc.go` ⇒ `ls` 命中 `risk`（判据翻转）→ 删除后 `git status --porcelain` **为空**（逐字节还原）。派工原文的「不适用」指无实现可破坏，本文仍验证了它给的判据本身可证伪。
- **⑥ 边界**：只出材料。三单的立单与编号归主任；① 的消费接线归 core F3、② 的别名切换归 orchestration / core / analytics 各自单、③ 的前置归 `SVC-BROKER-G1`；risk `ARCHITECTURE.md:85` 与 orchestration `ARCHITECTURE.md:100/174/185` 的判据订正各归其仓。

### QT-22 契约腿 · `version_hash` canonical 六项 + `v0.5.0` [P0] ✅ 完成（2026-10-24）

> 派工来自 G2；本文只记本仓（contracts）侧。strategy 仓的别名切换与 core 的消费不在本单。

- **① 目标**：把 `version_hash` 的 canonical 集合收进契约并冻结，新增 `StrategyVersionInput` / `DomainStrategyVersionHash`，使 contracts 与 strategy 两仓对同一输入算出同一个 hex。
- **② 上下文**：`engine/provenance.go`（QT-03 三 hash 的唯一正文）、`engine/canonical.go`（域标签）、`tellurion-strategy:internal/asset/hash.go:32-79`（本地 `DomainStrategyVersionHash` + `VersionInput` + `VersionHash` 的六项编码）、裁决材料 [`design/qt22-version-hash-decision-material.md`](design/qt22-version-hash-decision-material.md)。
- **③ 要求（全部落地）**：① `StrategyVersionInput` 六字段 `StrategyID` / `Carrier` / `ABIVersion` / `Scope`（`StrategyScope`）/ `Entry` / `Source` + `Hash()` / `MustHash()`；② `DomainStrategyVersionHash = "tellurion.strategy_version_hash.v1"`；③ 集合**不含** `capability`、**不含** `params_defaults`；④ 与 strategy 仓 `VersionInput` 逐字段同序并钉固定向量；⑤ 登记 `docs/API.md` §4 + `docs/api-surface.txt` + `CHANGELOG.md` + README 兼容矩阵 + `v0.5.0` tag。
- **④ 约束与非目标**：不复用 `DomainParamsHash`；不改 core / strategy 源码；不发 PATCH 级 tag；`params_defaults` 若要进必须升 `v2`（双标签过渡），本单不落。
- **⑤ 验收（真退出码）**：`go build ./... && go vet ./... && go test ./... -count=1` → **EXIT=0**（13 个包全绿，含 `go test ./conformance`）；`go run ./tools/check_api_surface` → 447/447；`go run ./tools/check_surface_freeze` → 430 行一致（426 → 430，本单 +4）；`go run ./tools/check_release_matrix` → EXIT=0；`go run ./tools/check_contract_shape` → 绿。
  - **对拍向量**：输入 `s1` / `wasm` / `on_bar/1` / `symbol` / `OnBar` / `x` → `0e6ad53637ce43742804fe325d97b6799372fe03319f1020fc26f0d748d68412`；`TestStrategyVersionInput_MatchesStrategyRepoEncoding` 逐行重建 strategy 侧编码断言同 hex。
  - **破坏验证**：删掉 `Hash()` 里的 `c.Str("source", in.Source)` → `go test ./engine -run 'TestGoldenVector_StrategyVersionHash|TestStrategyVersionInput_MatchesStrategyRepoEncoding|TestStrategyVersionHash_ChangesOnEveryDecidingInput'` **EXIT=1**（三条腿全部点名失败，实际 hex 变 `3112a2e7…`）；逐字节还原后 `engine/provenance.go` 的 `sha256=74E3A4BCEBADB8B002B4C055B8136D2EAF067E1E0AFAA259F79A7A86C9C4A90F` 与改坏前一致，`go test ./engine` 回到 EXIT=0。
- **⑥ 边界**：只落契约类型、域标签、向量与登记；strategy 仓的本地常量改为本仓别名归 G4 / G13 通知腿；`v0.5.0` 的下游切换与 core 消费不在本仓。版本取 `v0.5.0`（新增导出符号 = MINOR，README 规则表）。

### G4 · 跨仓发版协调：四仓升 pin `v0.5.0` [P0] ✅ 完成（2026-10-24）

> 本单是**跨仓通知腿**：改动落在 core / agent / data / strategy 四仓的 `go.mod` / `go.sum` / `AGENTS.md` §4；本仓只回写兼容矩阵与这条指针，不镜像四仓台账正文。

- **① 目标**：`v0.5.0` 发版后，core / agent / data / strategy 四仓把 `tellurion-contracts` pin 从 `v0.4.0` 升到 `v0.5.0`，并证明 pin 已合入、三侧 conformance 不漂移、矩阵三件套一致。
- **② 上下文**：本仓 `v0.5.0`（`70d9040`，已 push `origin/main` + tag）；四仓 `go.mod` 实测原 pin 均为 `v0.4.0`；agent 的 `replace tellurion => ../tellurion` 在册（ADR-D31），必须与 core 同批升。
- **③ 要求（全部落地）**：① 四仓 `go.mod` 升 `v0.5.0` + `go.sum` 同步（`GOPRIVATE=github.com/ruoxizhnya/*` + `GOPROXY=direct`，因 sumdb 对新 tag 尚未收录）；② 四仓 `AGENTS.md` §4 回写版本与说明；③ 跑 `tools/check_contract_pin.py` 与 `tools/check_contract_conformance.py`；④ 本仓 README 兼容矩阵 `v0.5` 行改「已切」。
- **④ 约束与非目标**：不改四仓业务代码；不删 agent 的 `replace tellurion`（S3 才删）；不把四个仓的本地 commit push 出去（各仓另有未推提交）；不消费 QT-22 新符号（core / strategy 的收编腿另行立项）。
- **⑤ 验收（真输出）**：
  - `python tools/check_contract_pin.py` → **GREEN**：6 个 pin 全部是目标仓 `origin/main` 祖先；四仓 pin 实测 `v0.5.0` → `70d90404df46`（2026-10-24）。
  - `python tools/check_contract_conformance.py` → core 通过 / agent 通过，**通过 2 / 失败 0 / 跳过 0**（语料 33 条 / 9 kind）。
  - 四仓 `go build ./... && go vet ./... && go test ./... -count=1` → **BUILD/VET/TEST 全 0**（core 78 包、agent 22 包、data 8 包、strategy 6 包）。
  - `go run ./tools/check_release_matrix` → 矩阵 ↔ CHANGELOG ↔ tag 一致（6 tag）。
- **⑥ 边界**：只做版本 pin 与登记回写；QT-22 的 core 消费与 strategy 别名收编不在本单。

### QT-05 · S1 契约层与工具层落地（代码直接写进两个新仓）[P0] ✅ 本仓侧已完成并已推（`fbe9a2d`）；**下游腿在册余项已于 2026-10-25 清零**（agent 删掉 core replace，QT-07 ⑥ S3 结案）

> **迁入说明（2026-10-21）**：本单自 core 仓 `tellurion:docs/TASKS.md`「S1」迁入（任务台账批量重分布：S1 双仓落点，主单落本仓）。六字段正文逐字保留、不改写；**kit 腿**拆分为 `tellurion-kit:docs/TASKS.md` 的 **QT-05-B**（③-② 的执行单，判据以本单为准）；core 侧留同号指针，core 侧改动（③-③/④/⑥/⑦/⑨）随本单执行。执行顺序（W0 全绿后的 S1）见 core `tellurion:docs/TASKS.md` 执行顺序总览。

> **进度（2026-10-22，本仓侧已完成）**：`fees/` ✅（13 符号）、`domain/market/` ✅（46）、`domain/` ✅（176）、`engine/` ✅（194）四包全落，登记进 `docs/API.md` §4，两条随仓护栏绿（导出面 **429/429**）。§8 待裁 **#1 已裁决并写成 [`adr/0003`](adr/0003-哨兵与-abi-不独立成包.md)**（errors//abi/ 随 engine/，不拆包；61 文件的旧估计实测证伪）、**#2 已裁决（[`adr/0002`](adr/0002-config-族的边界处置.md)）**、**#3 已裁决**（fees 整包）；只剩 **#4**（dataprotocol 落点时机，已登记为 X1a 待办）。
> ✅ **push 阻塞已解除（2026-10-22）**：本仓 6 个 commit 已推到 `origin/main`（该批 HEAD = `fbe9a2d`）；仓外空模块里 `GOPROXY=https://proxy.golang.org,direct go get github.com/ruoxizhnya/tellurion-contracts@fbe9a2d` 实测解析为 `v0.0.0-20261022191510-fbe9a2d0df17`。CI 全绿（导出面 429/429）。**下游腿收口状态（2026-10-25 实测）**：③-③ core 已切 pin（`go.mod` = `v0.6.0`）且 `pkg/domain` / `pkg/backtest/contracts` / `pkg/fees` 三目录已删、SPEC 互锁测试已迁 `e2e/guard/`；④ core/agent 的 contracts/kit 均无 `replace`；⑤ agent `go.mod` 已 require `v0.6.0` + `kit v0.1.0`；⑦ `tools/check_contract_pin.py` + `make contracts-pin-check` 在册；⑨ core-only 清单在 `docs/API.md` §6。**唯一余项** = ~~agent 仍保留 `replace tellurion => ../tellurion`（成功删它的前置是 agent 去掉 core 15 个实现包 import，归 QT-07 ⑥ S3，不属本单）~~ → **已结案（2026-10-25，QT-07 ⑥ S3 / B6-1）**：agent 对 core 实现包的 import 现为 **0 处**（11 条 `tellurion/pkg` 命中全是注释），`go.mod` 既无 `require tellurion` 也无任何 `=>`，CI 不再 checkout core；core 侧同批清空为它开的两条过期容忍（`check_contract_pin.py` 的 `LOCAL_REPLACE_WHITELIST` 现为空 dict、`check_repo_autonomy.py` 的 `KNOWN_GAPS` 删掉 `("tellurion-agent","R1")`）。取证见 `tellurion:docs/TASKS.md` QT-07 ⑥ 结案两条。**⇒ 本单在册余项清零。**
> ⚠️ **`domain/` 落地后，core 侧切换依赖那次 commit 多出两条待办**：① **5 个零引用的部署配置类型不必单独删**——它们就在 `pkg/domain/types.go` 里，随整目录删除自动消失（实测 0 调用点，不会牵连任何调用）；② `go.mod` 需要 `github.com/stretchr/testify v1.11.1`（域测试用到；同版本已在 core，不会引入第二个大版本）。
> ⚠️ **③-⑨ 的 core-only 清单已产出（2026-10-22 重测）**：`domain/` 根包 176 个导出符号折 **119 个根名**，agent import 了 **38** 个，**0 命中 81** 个；族与判据、再生成命令都在 `docs/API.md` §6 的 core-only 行——**清单以命令输出为准，不手抄**。S2 后按它回收。
> ⚠️ **`domain/market/` 落地后多出一条跨仓互锁**：core 的 `pkg/domain/spec_drift_test.go`（`OHLCV` / `Fundamental` 字段集 + `Frequency` 词表 ↔ `tellurion:docs/SPEC.md`）**刻意留 core**，core 侧 sed 后改为 import 本包。⇒ 在本仓给这两个类型加字段，core 侧 CI 会红，这是设计意图；那条测试也是本仓这两个类型字段集的**唯一机器判据**，别当它在迁移里被漏掉了。
> ⚠️ **② 的实测数字已滞后，且没写测试/非测试的拆分**（W0 期间 core 侧新增了 filltiming / provenance / canonical）。2026-10-22 重新取（`wc -l`）：
> `pkg/backtest/contracts` **非测试 10 文件 2543 行 ｜ 测试 11 文件 2455 行**（② 写「6 文件 2113 行」）；`pkg/domain` **非测试 16 文件 2208 行 ｜ 测试 8 文件 1795 行**（② 写「17 文件 3983 行」）；`pkg/domain/market` 非测试 5 文件 480 行 ｜ 测试 2 文件 430 行；`pkg/fees` 非测试 1 文件 250 行 ｜ 测试 1 文件 283 行。
> **别照抄 ② 的数字做工作量估计**。另：两个测试文件是**SPEC 互锁守卫，刻意留 core 不迁**（`pkg/domain/spec_drift_test.go`、`pkg/backtest/contracts/resultspec_test.go`——它们读 core 的 SPEC 文本），所以实际迁入的测试文件数是 7 与 10。
> ⚠️ **core/agent 那条腿（③-③④⑤⑦⑨）的 push 前置已解除**（该批 HEAD `fbe9a2d` 已在 `origin/main`），但**双份真相窗口仍在**：core 的 `pkg/domain` / `pkg/backtest/contracts` 与本仓同名包并存，两边都不许改，直到 core 侧 sed 落地并删旧包。

#### 下游腿交接单（2026-10-22，本仓侧交棒；**不在本仓执行**）

> 本仓侧四个代码 commit（`c9e2d6d` fees → `a044f80` domain/market → `f1d0e2d` domain → `cfad6c8` engine）与台账/自审 commit 已推到 **`origin/main`（2026-10-22，该批 HEAD `fbe9a2d`）**，下面每一步直接可执行。取版本用 `@fbe9a2d`（实测解析为 `v0.0.0-20261022191510-fbe9a2d0df17`），不要写 `@main`。

**A. core（`tellurion` 仓）**

1. `go get github.com/ruoxizhnya/tellurion-contracts@<推上去之后的 sha>`，再 `go mod tidy`。**不许**加 `replace`（③-④）。
2. import 全替（2026-10-22 实测文件数，**含测试**）：`"tellurion/pkg/domain"` **314 文件** → `github.com/ruoxizhnya/tellurion-contracts/domain`；`"tellurion/pkg/domain/market"` **19 文件** → `.../domain/market`；`"tellurion/pkg/backtest/contracts"` **61 文件** → `.../engine`（**包名仍是 `contracts`，调用点前缀不改**）；`"tellurion/pkg/fees"` **32 文件** → `.../fees`。
3. **先把两个 SPEC 互锁测试搬到 `e2e/guard/` 再删旧包**（它们读 core 的文档，不能随包走）：`pkg/domain/spec_drift_test.go` → `e2e/guard/spec_drift_test.go`（package 改为 `guard`，SPEC 路径 `../../docs/SPEC.md`，import 改指 `.../domain/market`；文件里的类型已是 `market.` 限定）；`pkg/backtest/contracts/resultspec_test.go` → `e2e/guard/resultspec_test.go`（SPEC 路径 `../../docs/SPEC.md`，契约从 `.../engine` 取并加 `contracts.` 前缀——原文件是 `package contracts`，用的是非限定名）。**先搬后删**，否则两个守卫随目录一起消失。
4. 删三个旧包目录：`pkg/domain/`（含 `market/`，5 个零引用部署配置也随它一起消失——`Config` / `DatabaseConfig` / `RedisConfig` / `ServiceConfig` / `TushareConfig` **不用单独删**）、`pkg/fees/`、`pkg/backtest/contracts/`（都只有 `.go` 文件，无资产）。**保留** `pkg/backtest/aliases.go`（core 自己的门面，35 个符号），只改它的 import。
5. ③-⑦：新增 `tools/check_contract_pin.py`（或在既有 `tools/check_repo_autonomy.py` 的 R5 上加一条腿——R5 现在只核**跨仓版本一致**，核不到「pin 是不是 contracts `main` 的祖先」），判据：pin 必须是该仓 `main` 的祖先，并打印 `repo@sha → 提交日期`；接进 `make contracts-pin-check` 与 CI。
6. **测试新增依赖**：core 的 `go.mod` 已有 testify，无需动；`go.sum` 会因 `go get` 增加 contracts 的条目。
7. 文档触发器（core 的 Rule 1）：`AGENTS.md` §3 目录树（`pkg/domain` / `pkg/fees` / `pkg/backtest/contracts` 三行）、`docs/ARCHITECTURE.md` 的包布局、`docs/SPEC.md` 里指 `pkg/domain` 的路径。
8. 验收（真退出码）：`go build ./...`、`go vet ./...`、`go test ./... -count=1`、`git grep -nE '"tellurion/(pkg/domain|pkg/fees|pkg/backtest/contracts)"'` **0 命中**、`git grep -nE '=> *\.\./tellurion-(contracts|kit)'` **0 命中**、`make contracts-pin-check` 绿、删掉 `go.work` 后仍能 build（证明走 `go.sum`）。

**B. agent（`tellurion-agent` 仓）**

1. `go get` 两个仓（contracts + kit）；~~**它仍保留 `replace tellurion => ../tellurion`**（S3 才删，见该仓 ADR-0001）~~ → **2026-10-25 订正：那条已删**（QT-07 ⑥ S3 结案；该仓 ADR-0001 已附「例外收回记录」，core 侧 `check_contract_pin.py` 的白名单同批清空为空 dict），但**不得**为 contracts/kit 加 replace —— 例外收回后这条是**无豁免的绝对禁令**。
2. import 全替（实测）：`"tellurion/pkg/domain"` **48 文件** → `.../domain`；`"tellurion/pkg/backtest/contracts"` **5 文件** → `.../engine`；`"tellurion/pkg/fees"` **3 文件** → `.../fees`；无 `domain/market` 直接 import。
3. `go build ./... && go vet ./... && go test ./...`；本地多仓联调用仓根 `go.work`（**不入库**，`.gitignore` 已挡）。

**C. 陷阱清单（都是这次实测踩到的）**

- `go get` 在 push 之前必然失败——**这是阻塞，不是 bug**；别用 `replace` 绕过。
- 先搬 SPEC 测试再删目录（见 A-3），否则两个唯一机器判据静默消失。
- `contracts/domain` 自带 `domain/market` 的别名，所以 sed 后 `domain.OHLCV` 这类写法**照常编译**——别顺手把它们改成 `market.OHLCV`，那会扩大 diff 且没有必要。
- `engine/` 目录里是 `package contracts`：sed 只换 import 路径，**不要**把 `contracts.X` 改成 `engine.X`。
- core-only 清单（`docs/API.md` §6）里那 81 个根名**先别删**：S1 整包迁的容忍条款到 S2 后才回收。

**D. 之后**

- S2（QT-07）：打 `v0.1.0` → 下游把伪版本换 tag；打 tag 前先把 README 兼容矩阵补 `v0.1` 那一行、CHANGELOG 把 `[Unreleased]` 搬进 `## [v0.1.0] - 日期`（漏一条 `check_release_matrix` 红）。
- QT-09：跨仓护栏（pin 祖先检查、`go.work` 不入库、surface 冻结）在 core/agent 侧的落地与自证。

> **形态已裁决（2026-10-20，用户选 Plan B → 落地形态见 ADR-D31）**：两仓已存在（各含最终 module path 的 `go.mod` + `AGENTS.md` + `doc.go` 占位），所以**代码直接写进新仓**，core/agent 用 `require <伪版本>` 取用，**全程不出现 `replace`**。原 Plan A（先在 core 内建 `contracts/` `kit/` 子 module，S2 再搬）作废——它会搬两次，且与 core 顶层 `contracts/` 数据目录撞名。

- **① 目标**：契约（含值对象）与横切工具**物理上住进两个独立仓**，依赖方向由编译器 + `go.sum` 强制；core 与 agent 通过**版本依赖**引用它们。
- **② 上下文**：`repo-topology-and-naming.md` §5/§7-S1（D29）+ ADR-D31；实测：`pkg/domain` 17 文件 3983 行，被 **core 306 文件、agent 48 文件**使用，自身仅依赖 `domain/market` 与 `fees`（2 处）；`pkg/backtest/contracts` 6 文件 2113 行，仅依赖 `domain` + `fees`，被 50 处（含 agent 5 处）使用；ADR-031 四契约（`RunSpec`/`SignalSpec`/`ResultSpec`/`EngineRunner`）就在这个包里。两仓现状：`github.com/ruoxizhnya/tellurion-{contracts,kit}`、`go 1.25`、**public**（`GET /repos/...` = 200）；`tellurion-contracts@main` 实测可 `go get` 到 `v0.0.0-20261020112435-a55208cc2412`（Go 1.25.0）。
- **③ 要求**：
  ① `tellurion-contracts`：`domain/`（17 文件，`package domain`）、`engine/`（6 文件，**保持 `package contracts`**）、`fees/`（**只放数据形状**：费率/滑点/成本模型的定义与枚举）；
  ② `tellurion-kit`：迁 12 个包 + 新增 **`db`**（仅 connection/migrate/DSN 解析，**零表结构**，K 裁决落地）；
  ③ core 侧：`go get github.com/ruoxizhnya/tellurion-contracts@<sha>`、`@...-kit@<sha>` 写入 `require` 与 `go.sum`；再一次性 sed 全替 `tellurion/pkg/domain` → `github.com/ruoxizhnya/tellurion-contracts/domain`（306 文件）与 `tellurion/pkg/backtest/contracts` → `.../contracts/engine`（50 处）；
  ④ **禁止 `replace`**：core / agent 的 `go.mod` 里不得出现 `=>`。本地要同时改三个仓时用 **`go.work`**（仓根，`use ( . ../tellurion-contracts ../tellurion-kit )`），**加入 `.gitignore`、不入库**——实测：有 `go.work` 时直接读本地目录、无 `go.sum` 也能构建；没有时按 `go.sum` 走真实版本（**这正是 CI 与 Docker 的路径**）；
  ⑤ **agent 侧同步**：agent 仓根 `go.work` 写 `use ( . ../tellurion ../tellurion-contracts ../tellurion-kit )`（它仍 replace core，到 S3 才删）；**`require` 同样不传递**，agent 必须自己 `require` 两个仓——漏一条它立刻编不过；
  ⑥ **core 顶层 `contracts/` 的 5 个数据文件留在原地**：`field_dictionary.yaml` + 4 个 `*.schema.{json,sql}` 是**运行时按 cwd 相对路径自动发现**的部署资产（`cmd/data/setup.go:223-225` 依次试 `./contracts/`、`../contracts/`、`../../contracts/`），不是 Go 包。跟 module 搬走 = 部署后字典**静默加载失败**（编译期查不出）。本单**只搬 Go 代码**；
  ⑦ 新增 `tools/check_contract_pin.py` + `make contracts-pin-check`：core / agent 的 contracts/kit pin **必须是该仓 `main` 的祖先**（禁止指向未合入的提交），并打印 `repo@sha → 提交日期` 供评审；
  ⑧ 包名策略：`contracts/engine/` 内保持 `package contracts`、`contracts/domain/` 内 `package domain`；调用点前缀一行不改，只在 import 处处理 alias。若改成 `package engine`，本仓 50 处 + agent 5 处都要改前缀；
  ⑨ 输出「**core-only 类型清单**」：`screen`/`sector`/`supply_chain`/`macro`/`mining` 在 agent 侧各 0 命中，属 core 专用；S1 仍整包迁，清单留给 S2 后回收（D29 §5.2）。
- **④ 约束与非目标**：**不改任何业务逻辑**；不搬 `sandbox/wasm`（必须与引擎同进程）；不搬 `fees` 的计算部分；**不动那 5 个数据契约文件**；不在本单打 tag（S2 = QT-07）；**不为省事给 core/agent 加 replace**；一次只搬一个包，每搬一次跑 `go build ./...` + `go vet ./...`。
- **⑤ 验收**：`go list -m all` 显示 3 个 module，且 contracts/kit 两行**带真实伪版本**（不是 `=> ./...`）；`go build ./...` + `go test ./...` 绿；旧路径 grep 零命中；`git grep -E '=> *\.\./tellurion-(contracts|kit)'` 在 core / agent 均零命中；在两仓内 import `tellurion/...` **编译失败**；`kit` 业务词零命中；agent 与 core 均能 build；`make contracts-pin-check` 绿；**删掉 `go.work` 后 `go build ./...` 仍绿**（证明走的是 `go.sum` 而非本地目录）；`equitydeep` 字典在 `cmd/data` 与部署镜像里都能加载（回归）。
- **⑥ 边界**：只做代码搬迁与依赖切换；tag、版本策略、跨仓护栏属 QT-07 / QT-09 / QT-19。

---

### CT-03 · `contracts/risk` 包：快照形状 + 冻结词表 + hash 纯函数 + 4 哨兵 [P0] ✅ 完成（2026-11-07）

> 立单说明：正文原以「候选单（建议编号 `CT-03`）」形态落在 [`design/shape-packages-decision-material.md`](design/shape-packages-decision-material.md) §5（A 组单的 ⑥ 把编号权交回主任）。本轮按第三轮裁决 **②**（contracts 新建 `risk` 包）与六仓 G1/G3 的实际等待关系落成正式条目，编号沿用建议值。**逐字搬运清单与四项组级禁区以材料 §1.1 与 §4.4 为准，本条不复制正文**（防两处漂移）。

- **① 目标**：把 risk 的规则快照**形状与 hash 算法**变成跨仓唯一实现，任何消费方 `import` 一次就能算出与 risk 逐字节相同的 `risk_rules_hash`。
- **② 上下文**：材料 §1.1（迁入清单）/ §1.2（编码必须逐字节搬运）/ §1.3（三个撞名坑）/ §1.4（与本仓既有判据的冲突，**需主任拍**）/ §4.4（四条组级禁区）；目标形状与算法在 `tellurion-risk:internal/rule/{snapshot.go:15-117,hash.go:15-55,model.go:12-64,model.go:150-174,errors.go}`；同形参照 `engine/provenance.go:78-87`；本仓随单必改四处护栏 = `tools/check_contract_shape/main.go:36` 的 `contractRoots`、`docs/API.md` §4、`docs/api-surface.txt`、`README.md` 兼容矩阵（材料 §1.5：不改则**假绿**）。
- **③ 要求**：① 新建 `risk/`（包名 `risk`），按 §1.1 迁入，**`Bundle` / `ParseBundle` / `ForScope` 与其余 7 个哨兵不迁**；② 编码逐字节搬运（`writeField` 的文本长度前缀 + 按 `rule_id` 排序 + `canonicalParams` 的 `UseNumber`→`Marshal`），**禁止**改用 `NewCanonical`；③ 固定向量测试（期望 hex 硬编码）+ 与 risk 的**对拍腿**；④ `contractRoots` 加 `risk` 并补夹具；⑤ 四件套（`API.md` §4 / `api-surface.txt` / 兼容矩阵 / `CHANGELOG`）同笔登记；⑥ 发 MINOR `v0.7.0`。
- **④ 约束与非目标**：不改 hash 算法与 `DomainSnapshotHash` 字面量（改 = 全历史 hash 作废，需 ADR 级裁决）；不改 4 个哨兵文案；不把 `risk` 并入 `engine`；不在契约里读时钟/连库/读环境变量；不动 risk 的 Store/HTTP/`Bundle`；不代做 core 的 F3 消费接线；**不同批**发 ② ③ 的形状包。另：材料 §1.4 那条「`Snapshot` 兼作 `/bundle` 响应体，是否违反『契约不放 decode』」的口径**需主任拍**，拍板前不落 `Snapshot` 之外的方法。
- **⑤ 验收**：① `go build ./... && go vet ./... && go test ./... -count=1` 三绿（贴真实退出码）；② `check_api_surface` / `check_surface_freeze` / `check_contract_shape` / `check_release_matrix` 四绿；③ 对拍腿绿（两侧 hex 相同）；④ 破坏验证两条：去掉 `risk:snapshot.go:63` 的 `sort.Slice` ⇒ 固定向量与对拍腿必须红；把 `writeField` 格式串改成明文分隔 `%s=%s\n` ⇒ 对拍腿必须红；两处**逐字节还原**后复绿（`sha256` 与改坏前一致）。
- **⑥ 边界**：只做「契约落形状 + 对拍 + 登记 + 发版」。core 的 `bundle → RiskPolicy` 接线归 `tellurion:SVC-CORE-RISK-RETIRE`（该单 ⑤ 已把本包列为前置），risk 侧改别名与 `ARCHITECTURE.md:85` 的判据订正归 risk 仓。
- **施工记录（2026-11-07）**：
  - 落码：`risk/snapshot.go` —— `DomainSnapshotHash`（字面量不动）、`MaxParamsBytes = 64 << 10`、`RuleType` + 4 常量、`Scope` + 4 常量 + `Valid()`、`SnapshotRule`（6 字段）、`Snapshot`（3 字段）+ `Clone()` / `Hash()` / `Verify()`；`risk/errors.go` —— 4 个跨仓哨兵，**文案与 `tellurion-risk:internal/rule/errors.go` 逐字一致**。编码照抄 risk：域标签裸写一行 + `\n`、`writeField` 文本长度前缀、按 `rule_id` 升序 + 重复即拒、`canonicalParams` 的 `UseNumber`→`Marshal` + 64 KiB 上限。**未用** `engine.NewCanonical`；`Bundle` / `ParseBundle` / `ForScope` 与 7 个服务私有哨兵**未迁**；包内零 I/O、零时钟。
  - 对拍：固定向量 `14ecab865bebaef54052354b63caa5d351f533ab112e1a29a9dcad3707955c5b` —— 2026-11-07 在 `tellurion-risk` 内临时调 `rule.SnapshotHash`（同一组两条规则）取得，命令退出码 0；临时测试文件跑完即删，`tellurion-risk` `git status --porcelain` 为空（实测）。另落 `TestSnapshotHash_MatchesRiskRepoEncoding`：测试内**独立重建** risk 侧编码器（不调本包编码辅助函数），断言两侧同 hex。
  - 护栏：`tools/check_contract_shape/main.go` 的 `contractRoots` 四根 → **五根**（+`risk`），头部注释 / 成功文案 / `inContractScope` 注释同步；夹具 `greenFiles()` 补 `risk/base.go`，新增「risk 命名红（`SnapshotConfig`）」用例，缺根红期望改为 `fees/、dataprotocol/、risk/`。
  - 登记：`docs/API.md` §4 新增 `risk/` 小节（22 行）+ §5.6「风控快照面」（编码三处差异、向量出处、两个撞名坑、不迁清单、真实消费方为零的取舍）；`docs/api-surface.txt` 重生成（448 → 470 行，`git diff --numstat` 实测 `+22 / −0`）；`README.md` 兼容矩阵加 `v0.7` 行；`CHANGELOG.md` 加 `[v0.7.0]` 节；`AGENTS.md` §1 状态串与 `docs/ARCHITECTURE.md` §2/§3/§8 同步。
  - 验收（真退出码）：`gofmt -l .` 空；`go build ./...` / `go vet ./...` / `go test ./... -count=1` 全 **0**（12 个包全绿）；`check_api_surface` → **487/487**（15 个代码包 / 7 个登记包）；`check_surface_freeze` → **470 = 470**；`check_contract_shape` → 绿（15 个包 / 45 个非测试文件 / 7 个白名单例外 / **五条根都在**）；`check_release_matrix` → 绿（矩阵序列 7 / CHANGELOG 版本节 8 / tag 7）。
  - 破坏验证（两条，均真做真还原）：① 把 `sort.Slice` 的比较器换成恒 `false`（排序摘掉）→ `go test ./risk` **EXIT=1**，`TestGoldenVector_SnapshotHash`（`730cf02e…` ≠ 向量）与 `TestSnapshotHash_MatchesRiskRepoEncoding`（`730cf02e…` vs `14ecab86…`）双红，另 `TestSnapshotHash_OrderInsensitive` 也红；逐字节还原后 `sha256=4B9B1CF63EBDD43DCFA6F05102CC2E22CC7F2E5BB46AB4C373514D8B85F76CE3` 与改坏前一致，`go test ./risk` 回 **EXIT=0**。② 把 `writeField` 格式串改成明文分隔 `%s=%s\n` → `go test ./risk` **EXIT=1**（向量与对拍腿双红，`e2a65ab9…` vs `14ecab86…`）；同样逐字节还原，`sha256` 同上、复绿。
  - 发版：本地 annotated tag `v0.7.0`。**未 push**（本轮纪律：只 commit 不 push）。
  - **未执行（如实）**：① core / risk / agent 三仓升 pin `v0.7.0` —— `go get` 要求 tag 已在远端，tag 未 push 前升 pin 会让三仓 `go build` 直接拉不到模块；该腿归「tag 推远端后的协调者步骤」，本单不发。② risk 侧改别名与 `tellurion-risk:docs/ARCHITECTURE.md:85` 判据订正（归 risk 仓）。③ core 的 `bundle → RiskPolicy` 接线（归 `tellurion:SVC-CORE-RISK-RETIRE`）。⇒ 本包今天**零 import**（与 `dataprotocol` 落库时同型的有意取舍，已写进 `docs/API.md` §5.6）。

### CT-04 · 编排三份 wire 形状（claim / `run_artifacts` / 子 run 入队） [P0] ✅ 完成（2026-11-07）

> **前置两条已代裁（2026-11-07）**：依据用户长期授权「待裁点按建议代裁并落文档」——
> ① 材料 §4 的 **wire DTO 口径 = A**：区分「core 的 HTTP 报文」（`docs/openapi.yaml` 生成物，仍禁手改）与
> 「服务间 Go 共享形状」（可手写进 contracts + conformance 对拍）。依据材料 §4.2/§4.3 与「六个服务仓今天 0 条 openapi 管线」的实测。
> ② 第六轮裁决 **#9**（orchestration / broker 的 contracts 本轮不做）**已解除**：本单开工即视为解除；
> broker 侧 CT-05 仍等 `SVC-BROKER-G1`（材料 §3.1.1 已裁 G1 先行），不在本单。
> 正文 = 材料 §5 的 CT-04（含 §2.1 的「只迁三份消费方形状」清单、§2.2 的「不迁」清单、§2.3 的文档冲突）。

- **① 目标**：把已在多仓**手写复刻**的三份 wire 形状收成契约里的唯一实现，让 core `docs/openapi.yaml` 里那句 `must not drift` 变成机械判据。
- **② 上下文**：生产侧 `tellurion-orchestration:internal/httpapi/{workers.go:103-111,provenance.go:32-48,parents.go:23-48}` + `internal/parents/service.go:130-132`；消费侧副本 `tellurion:internal/orchestrationworker/client.go:92-101`、`tellurion:docs/openapi.yaml:3802-3823`（`must not drift` 原文 `:3806-3807`）、`tellurion:pkg/storage/run_provenance.go:33-47`、`tellurion-analytics:internal/upstream/sensitivity.go:38-47`；向量测试形态照抄 `engine/provenance_vector_test.go`。
- **③ 要求**：① 新建包（建议 `run/`，**不要**叫 `orchestration`——仓名不是形状名），只收 §2.1 三份；② **不迁** `RunRecord` / `ParentRecord` / `WorkerRecord` / `jobView` / `parentView` / `ArtifactStore`；③ 每份一条跨仓对拍腿；④ `ChildIdempotencyKey` 的二选一（迁函数 vs 只迁形状）在单内定死并写进 `docs/api-schema.md`；⑤ 四件套登记 + MINOR `v0.8.0`。
- **④ 约束与非目标**：不改任何 wire 字段名/可空性/顺序（`provenanceView` 的 `*string` 出 `null` 而非省略键是**冻结语义**）；不迁队列状态机与租约语义；不在契约里做 decode / 默认值填充 / 幂等判定；不动 orchestration 的 HTTP / PG / DDL；不代做 orchestration `ARCHITECTURE.md:100/174/185` 的口径订正；**不同批**发 ① ③。
- **⑤ 验收**：① 三条对拍腿全绿；② 三条命令 + 四条护栏绿（贴真实退出码）；③ 破坏验证两条：去掉 `claimView` 的 `parent_job_id` 的 `omitempty` ⇒ 对拍腿必须红；把 `ChildIdempotencyKey` 的分隔符 `:` 改成 `-` ⇒ 必须红；逐字节还原后复绿。
- **⑥ 边界**：只做契约落形状 + 对拍 + 登记；orchestration 改别名、core 的 `orchestrationworker` 改用契约类型、analytics 的提交报文改用契约类型，各自立单。

- **施工记录（2026-11-07）**：
  - **落码**：`run/run.go` —— `ClaimView`（7 字段，`ParentJobID` 带 `omitempty`）、`ProvenanceView`（15 字段，6 个指针字段出 `null` 不省略键）、`ParentSubmitRequest` / `ParentBatchRequest` / `ParentWalkForwardRequest` / `ParentWindows`、`ChildKeyPrefix = "parent"` / `ChildKeySeparator = ":"` / `ChildIdempotencyKey` 纯函数。包注释写明冻结语义、不迁清单（`RunRecord` / `ParentRecord` / `WorkerRecord` / `jobView` / `parentView` / `ArtifactStore` / DDL）与「零 I/O、零时钟、零 decode」。
  - **`ChildIdempotencyKey` 二选一结论 = 整段迁入契约**（不是只迁形状）：它是纯函数（ADR-0004 四条），且 ⑤ 的破坏验证要求期望值由契约侧实现产出。判据写进 `docs/api-schema.md` §3.1（含「改规则 = MAJOR」与父 id 语义）。
  - **三条对拍腿**：`run/run_test.go` —— `TestClaimView_MatchesOrchestrationEncoding`（带/不带父 id 两组，钉 `omitempty`）、`TestProvenanceView_MatchesOrchestrationEncoding`（全字段 + 全空两组，钉 15 键与 `null`）、`TestChildIdempotencyKey_MatchesOrchestrationRule`（3 个金向量 + 前缀/分隔符常量）。另加 `TestParentSubmitRequest_MatchesOrchestrationEncoding`（batch / walkforward 两组）。期望值全是**生产侧形状的字面量**，不从本包反推。
  - **跨仓语料**：`conformance/run_cases.json`（7 用例：claim 2 / provenance 2 / child key 3）+ `conformance/run_conformance_test.go`。**单独一份文件而不是塞进 `cases.json`**——core / agent 的加载器只认 `contracts` / `core` / `agent` 三侧，混入 `analytics` / `orchestration` 侧会让下游升级 pin 到本版本时当场红。语料 `sides` 只用 `["contracts"]`，期望值由本仓实现产出后写死。
  - **护栏**：`tools/check_contract_shape/main.go` 的 `contractRoots` 五根 → **六根**（+`run`），头部注释 / 成功文案 / `inContractScope` 注释同步；夹具 `greenFiles()` 补 `run/base.go`，新增「run 命名红（`ClaimOptions`）」用例（`risk` 命名红用例保留），缺根红期望改为 `fees/、dataprotocol/、risk/、run/`，绿腿标题改「六条根」。
  - **登记**：`docs/API.md` §4 新增 `run/` 表（9 个导出符号：6 struct + 2 const + 1 func）、§5 新增 §5.7「编排 wire 面」（六条判据 + 迁移边界）；`docs/api-surface.txt` 用 `check_surface_freeze -update` 重生成（`git diff --numstat` 实测 **+9 / −0**，408 → 417 行）；`docs/api-schema.md` 新增 §3.1；`README.md` 兼容矩阵加 `v0.8` 行；`CHANGELOG.md` 加 `[v0.8.0]` 节；`AGENTS.md` §1 状态串与 `docs/ARCHITECTURE.md`（包图 / §3 表 / §8 表）同步。
  - **验收（真退出码）**：`gofmt -l .` 空（exit **0**）；`go build ./...` **0**；`go vet ./...` **0**；`go test ./... -count=1` **0**（`run` / `conformance` / 四个工具包全绿）；`check_api_surface` **427/427**、exit **0**（16 个代码包 / 8 个登记包）；`check_surface_freeze` **417 = 417**、exit **0**（10 个库包）；`check_contract_shape` exit **0**（16 个包 / 45 个非测试文件 / 白名单 7 个）；`check_release_matrix` exit **0**（矩阵序列 8 / CHANGELOG 9 节 / tag 8）。
  - **破坏验证（两条，均真做真还原）**：① 去掉 `ClaimView.ParentJobID` 的 `omitempty` ⇒ `go test ./run -run TestClaimView_MatchesOrchestrationEncoding` **EXIT=1**（`without_parent_job` 子用例：多出 `"parent_job_id":""`），`go test ./conformance -run TestRunCases_ContractSide/run-claim/without-parent` 同 **EXIT=1**；逐字节还原后 `sha256=63AC5CEDC148056294F4B17E86BD010A4EE158B225309929461E5D7119D0D398` 与改坏前一致，复绿。② `ChildKeySeparator` 由 `":"` 改 `"-"` ⇒ `go test ./run -run TestChildIdempotencyKey_MatchesOrchestrationRule` **EXIT=1**（常量守卫直接点名 `sep="-"`），`conformance` 三条 child-key 用例全红（`parent-parent-abc-0` vs `parent:parent-abc:0`）；同样逐字节还原，`sha256` 同上、复绿。
  - 发版：本地 annotated tag `v0.8.0`。**未 push**（本轮纪律：只 commit 不 push）。
  - **未执行（如实）**：① orchestration 改别名、core `internal/orchestrationworker` 改用契约类型、analytics `internal/upstream` 提交报文改用契约类型 —— 都是各自仓的单（⑥），本单落地时三处仍各自持有副本。② 三仓升 pin `v0.8.0` —— `go get` 要求 tag 已在远端，未 push 前升 pin 会让 `go build` 拉不到模块。⇒ `run/` 包今天**零 import**（与 `dataprotocol` / `risk` 落库时同型的有意取舍）。
  - **订正（2026-11-08）**：`v0.8.0` / `v0.9.0` tag 与 `main` 已推远端、core 已 pin `v0.9.0`（补裁 #13-2）⇒ 上面「未 push / 升 pin 会拉不到」的前提**已失效**；`run/` 仍零 import 是**未施工**，不是被 push 卡住。

### CT-05 · broker DTO（order intent / fill / account） [P1] 🔶（2026-11-07：契约侧 `broker/` 已落；broker 删本地形状改引本包 = 未执行；tag 已于 2026-11-08 推远端，前置全解除）

> 正文 = 材料 §5 的 CT-05，材料 §4.3 对 ③ 的建议是**先不做**，只出对表。⛔ 的两条前置是**互锁**的：`SVC-BROKER-G1`（订单 HTTP 面）当时的阻塞理由写在 `tellurion-broker:docs/TASKS.md:37`（现 `:37` 起，行号会漂；等 A 组 broker DTO），而本单的阻塞理由又写「等 `SVC-BROKER-G1` 落定」——**这条死锁已于 2026-10-26 裁定，不是待裁项**：见材料新增 **§3.1.1**（三条理由 + 代价 + 两条窗口判据），结论 = **G1 先行、本单事后收编**，broker 侧已据此把 G1 从 ⛔ 改判 ⬜（`tellurion-broker:docs/TASKS.md:37`）。
> 本单剩下的前置因此**只剩 `SVC-BROKER-G1`**：§4 的 wire DTO 口径已于 2026-11-07 在 CT-04 代裁为 **A**
> （服务间 Go 共享形状可手写进 contracts + conformance 对拍；core 的 openapi 生成物仍禁手改），不再是待裁项；G1 不等它。
> ⚠️ 上面那句「这条死锁需主任裁定先后（建议：先裁 §4 口径 → broker 自持 DTO 落 G1 → 本单再收编）」是同批施工方的原文，方向与本裁定一致；**「本单不自行解开」已被执行侧代裁取代，待主任追认**。

- **① 目标**：把 broker 的订单 / 成交 / 账户记账形状在**订单 HTTP 面定形之后**收进契约，避免两个消费方各写一份镜像。
- **② 上下文**：材料 §3.1（今天一个都没上 wire）/ §3.2（order intent 三个候选名，选错=造第三份）/ §3.3（三处必须写进契约文档的形状差异）/ §3.4（kit 机械上被排除）；DTO 定义 `tellurion-broker:internal/gateway/model.go:80-186`；形状对表 `tellurion-broker:docs/ARCHITECTURE.md:152-167`；12 个哨兵 `internal/gateway/errors.go:12-36`。
- **③ 要求**：① 前置满足后才迁；② 迁入时把 §3.3 三条差异写进 `docs/api-schema.md`（数量 `int` vs `float64`、`Side` 两值 vs `Direction` 四值、`OrderType` 两值 vs 四值），**禁止**用类型别名把两侧等同；③ 12 个哨兵是否随迁在单内另判（它们更接近 ① 的判据）。
- **④ 约束与非目标**：不迁 XTP 适配器内部类型（`BusinessType` / `XTPConfig` / `XTPDeps`）；不在 HTTP 面定形前冻结报文；不动 `broker.order_acks` 的 DDL 与幂等语义（保留期归 `tellurion-broker:SVC-BROKER-G2`，R10 归 `SVC-BROKER-G4`）。
- **⑤ 验收**：① 对表逐行可查（DTO × 字段 × 契约位 × 与 `domain` 的差异）；② `SVC-BROKER-G1` 已 ✅ 且有可指的端点清单；③ 若落包：对拍腿 + 四条护栏绿 + 破坏验证（改一个字段名 ⇒ 红 → 还原 → 绿）。
- **⑥ 边界**：只做 broker 形状；真实 XTP（X6c ⛔）与 core `pkg/live` 退役（`tellurion:SVC-CORE-BROKER-RETIRE`）不在本单。

> **施工记录（2026-11-07 · 契约侧落包）**
>
> **前置复核（当场实测，不照抄文档）**：① `SVC-BROKER-G1` 已 ✅（`tellurion-broker:docs/TASKS.md` 头部与正文，含 12 端点清单与
> 逐字段报文表）；② 窗口期两条义务已满足——`tellurion-broker:docs/api-schema.md` §1.5 有订单 / 成交 / 账户三族逐字段表
> （字段名 + 可空性 + 类型），`internal/httpapi` 有序列化键集合用例；③ 材料 §3.3 三处差异已对照 `internal/gateway/model.go:80-186` 逐条核过。
>
> **落点**：新建 `broker/`（`enums.go` / `orders.go` / `account.go`）——31 个导出符号（12 `const` + 3 枚举类型 + 12 struct + 4 谓词方法）。
> **不迁**（④ + 本单另判）：XTP 适配器内部类型、`OrderRequest.Fingerprint`（幂等指纹是 broker 的仓储规则）、12 个契约哨兵
> （进程内错误值；HTTP 消费方按键值码判定）、`broker.order_acks` DDL。⇒ ③ 的「哨兵是否随迁」裁定 = **不迁**。
>
> **⑦ 字节来源（关键判据）**：期望值**不是**按契约侧反推的——在 `tellurion-broker` 加一次性取证用例
> （`internal/gateway/zz_capture_tmp_test.go`，跑完即删），把 13 条代表值（含空数组、`omitempty` 缺席、零时间）
> 由生产类型 `json.Marshal` 打出 `CAP|<id>|<bytes>`，逐字节抄进 `broker/wire_test.go` 与 `conformance/broker_cases.json`（10 例）。
>
> **登记**：`docs/API.md` §4 `broker/` 表（31 行）+ §5.8；`docs/api-schema.md` §3.2（三处差异 + 哨兵不迁）；
> `check_contract_shape` 加 `broker` 根（六条 → 七条，含绿基座、白名单夹具与命名红用例）；`docs/api-surface.txt` 用 `-update` 重生成。
> `README.md` 兼容矩阵加 `v0.9` 行；`CHANGELOG.md` 新建 `## [v0.9.0] - 2026-11-07`（并把 CT-06 的破坏性条目归位到它实际随发的 `v0.8.0` 节）；
> `AGENTS.md` §1 / `docs/ARCHITECTURE.md`（包图、§3 表、§8 表）同步。
>
> **验收（真实输出，2026-11-07）**：`gofmt -l broker conformance` 空；`go test ./broker ./conformance -count=1` 全绿；
> `go run ./tools/check_api_surface` → **458 代码符号 / 458 登记符号，17 个代码包 / 9 个登记包，EXIT=0**；
> `go run ./tools/check_surface_freeze -update` → 基线 **417 → 448 行（+31 / −0）**，复跑 `448 = 448` EXIT=0；
> `go run ./tools/check_contract_shape` → **17 个包 / 48 个非测试文件 / 白名单 7 个，七条契约根目录都在，EXIT=0**；
> `go run ./tools/check_release_matrix` 见下一条。
>
> **破坏验证（两条，真做真还原）**：① 把 `OrderRequest.LimitPrice` 的 `,omitempty` 去掉 ⇒ `winner` 腿红
> （`broker/wire_test.go` 的 `order_request_market` 多出 `"limit_price":0`，`conformance` 同名用例同红）；逐字节还原后复绿。
> ② 把 `OrderStatus.Terminal` 的 `OrderExpired` 从终态集合删掉 ⇒ `TestEnums_BrokerVocabularyIsClosed` 红；还原后复绿。
>
> **未执行（如实，本单保持 🔶）**：`tellurion-broker` 删本地形状改引 `contracts/broker`——**前置已全部解除**（2026-11-08 补裁 #13-3）：`v0.9.0` tag 与 `main` 已推远端、`SVC-BROKER-G1` 已 ✅，`go get …@v0.9.0` 可解析；由 broker 一笔做，判据 = 同一条报文快照改前 / 改后逐字节一致。原文「tag 未 push ⇒ 会不可构建」已被证伪。

### CT-06 · `domain/` core-only 根名的回收（重跑：119 根名） [P3] 🔶（0 使用方已删；61 个单使用方迁回待跨仓单）

- **① 目标**：把 QT-05 ③-⑨ 产出、被两句「先别删 / S2 后回收」押下的 core-only 导出面回收掉，让本仓的导出面只剩**真有 ≥2 使用方**的形状（准入规则 C4）。
- **② 上下文**：本仓 `docs/TASKS.md:159`（QT-05 ③-⑨ 的实测：`domain/` 根包 176 个导出符号折 **119 个根名**，agent import 38，**0 命中 81**；清单以命令输出为准，**不手抄**）与 `:193`（「那 81 个根名先别删：S1 整包迁的容忍条款到 S2 后才回收」）；S2 / S3 均已 ✅（core `docs/TASKS.md` 阶梯表），登记的容忍条款已到期；清单落点 `docs/API.md` §6 的 core-only 行。
- **③ 要求**：① **先重跑再生成**（`81` 是 2026-10-22 的读数，core 已切五包到 kit、agent 已零 import，数字必变）；② 逐名核当前使用方数量：仅 1 个使用方 → 迁回该使用方仓；0 使用方 → 删除；③ 删除即 **MAJOR**（本仓 README 版本规则：删符号属 breaking），按兼容矩阵走「保留一个版本周期」的条款，先在 `CHANGELOG` 标 Deprecated 再删；④ `docs/API.md` §4 登记区与 `api-surface.txt` 同步（`check_surface_freeze -update`），不许只改基线。
- **④ 约束与非目标**：不改任何仍在用的形状；不动 `engine/canonical.go:44-66` 的 8 个 `Domain*`（禁区，材料 §4.4）；不借回收改字段名或语义；不在本单引入 kit/依赖；不代删 core 侧的实现包。
- **⑤ 验收**：① 一份「根名 × 使用方仓 × 处置（迁回/删/保留）」清单，每行有 grep 出处；② `go run ./tools/check_api_surface` 与 `check_surface_freeze` 绿，且导出符号数的**减少量**与清单逐条对得上；③ 消费方（core / agent / 各服务仓）在下一版 pin 上 `go build ./...` 绿；④ 破坏验证：随便删一个仍被 agent 用的名字 ⇒ agent 编译失败（证明清单的「0 使用方」判据不是空话）。
- **⑥ 边界**：本仓导出面回收；core 侧 `pkg/*` 的包级退役归各 `SVC-CORE-*-RETIRE` 单。
- **跨仓裁决（2026-11-08 补裁 #12，**待主任追认**）**：61 个「仅 1 使用方」根名**按原建议迁回**，不放宽 C4；拆两笔——`CT-06-A`（agent 的 `BuiltinHypothesisFor`，单仓单点，先做）与 `CT-06-B`（core 60 个，按「兼容别名 / 实类型」分类：别名类只在 core 建同形别名、身份不变；实类型类才物理搬家并同批改签）。每笔 = 本仓删符号 + 使用方本地定义 + `CHANGELOG` 记 breaking + `api-surface.txt -update`。**本轮只裁决与拆分，不施工**；正文见 `tellurion:docs/design/qt-decoupling/cross-repo-decision-record-2026-10-26.md` §四.12。
- **施工记录（2026-11-07）**：
  - **重跑取证（取代 2026-10-22 的「81」）**：`go run ./tools/check_api_surface -list` 实测 `domain/` 根包 **176 符号 / 119 根名**（方法并回类型）。11 个同级仓（`C:\Users\ruoxi\Workspace\tellurion*`，排除本仓；`*.go` 含测试；只统计导入 `contracts/domain` 的文件）逐名 grep `\bdomain\.<Name>\b`，命令与逐名 `repo:path:line` 出处见 `docs/API.md` §6.1。
  - **处置结果**：**删除 22 根名 / 69 符号**（公司行为族 68：`CorporateAction` / `ActionType` / `CashDividend` / `BonusShare` / `CorporateActionSplit` / `RightsIssue` / `Placement` / `ActionEngine` / `ApplyResult` / `ApplyOutcome` 及构造函数与方法，整文件 `domain/corporate_action.go` 与测试同删；旧名别名 `MarketDataProvider` 1）——11 仓全部 0 命中；**保留 34 根名**（≥2 使用方，含 core+agent / core+data）；**保留 2 根名**（`CaliberRank` / `CaliberRaw`，0 外部引用但属 `FactorCaliber` 封闭词表且 `Valid()` 依赖，删会改语义）；**迁回建议 61 根名**（均只有 1 个使用方：`tellurion` 60 + `tellurion-agent` 1 `BuiltinHypothesisFor`），本轮**不执行**——迁回须改 core / agent 并删本仓公共符号，超出「只改本仓」边界。grep 初筛后用 Go AST（`go/parser` 只数 `SelectorExpr`）复核，排除字符串/注释：`OHLCV` / `OrderType` / `OrderTypeLimit` / `OrderTypeMarket` / `Portfolio` 的 agent 命中在 `pkg/ai/prompts.go` 提示词代码示例字符串里，从「2 仓」改判为「1 仓（core）」。
  - **版本规则裁决**：`README.md:63-64` 明确 0.x **允许 breaking，bump MINOR + 写 CHANGELOG 即可**；「deprecated 旧接口保留 1 个 MINOR 周期」是 1.0 之后的义务，0.x **没有**。故本轮直接删净，未先标 Deprecated；`CHANGELOG.md` 的 `[Unreleased]` 已按「破坏性：<旧 → 新>」留账，下一 MINOR 发版。
  - **登记与基线**：`docs/API.md` §4 删除公司行为族行与 `MarketDataProvider` 行、§5.2 删 `CorporateAction` 纯函数行、§6 分组表同步、新增 §6.1 逐名清单；`docs/api-surface.txt` 用 `go run ./tools/check_surface_freeze -update` 重生成，`git diff --numstat` 实测 **−62 / +0**（69 个 API 符号中 7 个 `CorporateAction` 接口方法在冻结基线里折进 interface 指纹行，故基线少 62 行）。
  - **验收（真退出码）**：`gofmt -l .` 空（exit 0）；`go build ./...` **0**；`go vet ./...` **0**；`go test ./... -count=1` **0**（12 个包全绿）；`check_api_surface` **418/418**、exit **0**；`check_surface_freeze` **408 = 408**、exit **0**；`check_contract_shape` exit **0**；`check_release_matrix` exit **0**。
  - **破坏验证（真做真还原）**：在仓外临时模块 `%TEMP%\ct06-destru` 用 `replace` 指向本仓，引用 agent 在用的 `domain.WalkForwardParams`（agent 侧 11 命中，`pkg/ai/http_runner.go:189`）→ `go build ./...` **exit 0**；临时删掉 `domain/backtest.go` 的该类型声明 → 同一模块 **exit 1**（`undefined: domain.WalkForwardParams`）；逐字节还原后 **exit 0**，本仓 `go build ./...` **exit 0**，`git diff domain/backtest.go` 为空。
  - **未执行（如实）**：① 消费方（core / agent / data）在「下一版 pin」上 `go build ./...` —— 本轮不打 tag / 不 push，`go get` 拉不到新版本，未执行；② 56 个单使用方根名的迁回（需跨仓单，且会删本仓公共符号）；③ tag / push（纪律禁止）。
  - 状态保持 **🔶** 的原因：0 使用方已删净，但 56 个「仅 1 使用方」尚未迁回，导出面还未达到 C4 的「只剩 ≥2 使用方」。

---

## 施工序（候选单，开工前必须展开成六字段）

> 这张表是**施工顺序**，不是任务台账：每条在执行前要展开成六字段单（目标/上下文/要求/约束/验收/边界），并写明冻结契约的路径与节号。
> 顺序有依赖，**不要跳级**：跳序的代价是返工（例：X1a 的协议没冻结就写 aggregator，等于把协议写死在一个服务里）。

| 阶段 | 内容 | 上游单号 | 前置 |
|---|---|---|---|
| S1 | 把 `pkg/domain` + `pkg/backtest/contracts` 搬入本仓（**文件数/行数以 QT-05 进度注里 2026-10-22 的实测为准**，本行原先抄的「17 文件 / 6 文件」已滞后）；core/agent 改 `require` 伪版本（正文见「一、本仓任务」QT-05） | `QT-05` | W0 全绿 |
| S2 | 打 `v0.1.0`；下游把伪版本换 tag | `QT-07` | S1 全绿 |
| S3 | 版本规则表 + 兼容矩阵 + CHANGELOG | `QT-19` | S2 打 tag 后 —— ⚠️ **已提前落地（CT-02，2026-10-22）**：README 规则表 + 兼容矩阵 + `CHANGELOG.md` + 护栏 `tools/check_release_matrix` 都在册。S2 打 tag 时只剩「矩阵补 `v0.1` 那一行 + CHANGELOG 把 `[Unreleased]` 搬进 `## [v0.1.0] - 日期`」，别重做一遍 |

## 二、跨仓指针（不在本仓，勿双写）

| 单号 | 主题 | 在哪 |
|---|---|---|
| QT-05 | S1：把 `pkg/domain` + `pkg/backtest/contracts` 搬进本仓 | **已迁入本仓「一、本仓任务」**（2026-10-21；core `tellurion:docs/TASKS.md` 留同号指针；kit 腿见 `tellurion-kit:docs/TASKS.md` QT-05-B） |
| QT-07 | S2：打 `v0.1.0` 并把下游伪版本换成 tag | `tellurion:docs/TASKS.md` |
| QT-19 | 版本策略 / 兼容矩阵 / CHANGELOG 规则 | `tellurion:docs/TASKS.md` |
| QT-09 | 跨仓护栏（surface 冻结 / 契约形状 / PR 模板 / examples / 契约测试） | `tellurion:docs/TASKS.md`（主单）；**本仓侧交接（2026-10-23）**：② `examples/`（`market.Provider` / `EngineRunner` 最小用法）、③ `conformance/`（`cases.json` 33 条 / 9 kind + 契约侧自检腿；core `pkg/conformance` 与 agent `internal/conformance` 读同一份语料，core 侧编排 `make contracts-conformance`，未 checkout 的侧显式 skip 不算通过）、④ `.github/PULL_REQUEST_TEMPLATE.md` + `tools/check_pr_body`（仅 PR 事件）、⑤ `tools/check_contract_shape`（`*Config`/`*Options`/`*Impl` 与 DB import 即红）、surface 冻结 `tools/check_surface_freeze` + `docs/api-surface.txt`；四条判据见各工具头部注释，全部接进 CI |
| QT-40（T2） | `dataprotocol` 要不要独立成包、何时迁 | 单在 `tellurion-data:docs/TASKS.md`（QT-40 的接口对表待办 T2）；本仓侧 **X1a 契约腿已落（2026-10-24，`v0.6.0`）**：§8 #4 定 (a)、(a) 的形状与本仓落点见 `docs/API.md` §4/§5.5。`tellurion:docs/design/qt-decoupling/data-service-contract-and-ownership.md` §4.1 注 ② 的回写归 core 侧协调者；core / data 的删副本与改 import 归各自单 |

---

## 维护本文件

- **新增单**：写全六字段；六字段里必须显式引用**冻结契约的路径与节号**（契约仓的包路径、设计文档的节号），不许写"参考相关实现"。
- **派工前先证伪前提**：读码核对工单写的前提；不符就停下回报 + 给拆法，不硬做、不扩范围。
- **跨仓引用写 `repo:path`**（如 `tellurion:docs/TASKS.md`），不做可点击相对链接。
- **顺序纪律**：同一波次内撞同一批文件的单按单号串行，一单一 commit。
