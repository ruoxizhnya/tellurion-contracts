# TASKS 归档 — tellurion-contracts（2026-10-22）

> 规则见 [`README.md`](README.md)：完成即在本仓台账标 ✅，随后整条搬进本目录，**旧代正文不改写**（要订正只在文件头加说明）。
> 本文件是历史证据，不作为当前依据引用；引用前先确认它没被新文档取代。

---

## CT-02 · README 版本规则表 + 兼容矩阵 + CHANGELOG + `GOPRIVATE` 前提 [P1] ✅ 完成（2026-10-22）

- **① 目标**：本仓作为被依赖库有成文的版本纪律与拉取前提，breaking 有账可查。
- **② 上下文**：refactor-plan-review **L4/L5**；规则口径冻结在 `tellurion:docs/TASKS.md` QT-19 ③（**规则内容以彼处为准，本单只落本仓文档**）；`tellurion:docs/design/qt-decoupling/repo-inventory.md` §1.1/§5 的可见性与 GOPRIVATE 风险注。
- **③ 要求**：① README 写**版本规则表**（`MAJOR` = 删/改签名或语义；`MINOR` = 新增兼容能力；`PATCH` = 文档/无行为修正）；② **兼容矩阵**：`vX.Y` ← 支持 core / agent 的哪些版本，每次发版更新；③ CHANGELOG 至少一条记录；④ AGENTS.md 写一行 **`GOPRIVATE=github.com/ruoxizhnya/*` + git 凭证**前提（本仓转私有时 CI/构建必须配，否则 `go get` 走 proxy 必然失败）。
- **④ 约束与非目标**：不引 release 工具链（goreleaser 等）；内容寻址标识（`strategy_ref`/`artifact_hash`/`run_spec_hash`）不参与 semver；0.x 阶段允许 breaking 但须 bump MINOR + CHANGELOG 条目（细则以 QT-19 为准）。
- **⑤ 验收**：README 存在版本规则表与兼容矩阵；CHANGELOG ≥ 1 条；AGENTS.md 含 GOPRIVATE 行；**破坏验证**：发一版不更新兼容矩阵 → 发版 checklist 项不可勾（或检查脚本红）。
- **⑥ 边界**：只定本仓策略与文档；kit 的同名交付是 `tellurion-kit:docs/TASKS.md` KT-02；「禁删 deprecated」机械护栏属 `tellurion:docs/TASKS.md` QT-09 ⑤。

### 结案记录（2026-10-22）

**交付物**

| 文件 | 内容 |
|---|---|
| `README.md` | 新增「版本规则」（三行表 + 0.x 条款 + 两条 Go 特有判定补充 + 内容寻址不参与 + tag 规范）、「兼容矩阵」（今天只有一行基线，见下面的现状说明）、「发版流程」（五条 checklist，其中两条有 CI 护栏）、「拉取前提（私有仓）」（两条：`GOPRIVATE` + git 凭证，含现状 12 仓全 public 的时点说明）；`## 目录` 与「快速开始」补两份新交付 |
| `CHANGELOG.md` | 新建。文件头写明标题约定 `## [vX.Y.Z] - YYYY-MM-DD` 与 breaking 的写法要求；两条 PATCH 级条目（CT-02 本次、CT-01 2026-10-21），全在 `[Unreleased]` 下 —— **本仓还没发版，所以没有版本节** |
| `AGENTS.md` | §2 加 `check_release_matrix` 命令 + **私有化前提**那一行（含 `repo-inventory.md` §1.1/§5 的风险注出处）；§6 登记「第五份 = 版本纪律三件套」并把「发版 → 改矩阵那一行 + CHANGELOG 那一节」写进变更触发串 |
| `tools/check_release_matrix/` | 新建护栏（CT-02 ⑤ 的「检查脚本红」腿）。五个红向：结构红 / 发版没更新矩阵 / 矩阵行没账 / tag 没账 / GOPRIVATE 前提缺失 / 空账 |
| `tools/check_release_matrix/main_test.go` | 16 条用例（夹具 `testdata/ok/` = 已发版且齐；`testdata/baseline/` = 本仓今天的零发版形状，必须绿） |
| `.github/workflows/ci.yml` | 加「版本纪律核对」步骤；checkout 改 `fetch-depth: 0`（默认只取 HEAD 时 `git tag` 恒空，tag 那条腿会变成假绿） |
| `docs/README.md` | §1 加「发版 / 改版本纪律该读哪份」、§2 文档地图登记 `CHANGELOG.md` + 新护栏、§3 变更触发表加「发版 → 矩阵那一行 + CHANGELOG 那一节」 |
| `docs/TASKS.md` | CT-02 标 ✅、六字段正文搬进本文件，台账只留指针 |

**③ 要求逐条落点**

| 要求 | 落在哪 |
|---|---|
| ① 版本规则表 | `README.md`「版本规则」，`<!-- rules:begin/end -->` 之间，MAJOR/MINOR/PATCH 三行 |
| ② 兼容矩阵 | `README.md`「兼容矩阵」，`<!-- matrix:begin/end -->` 之间，按**版本序列** `vX.Y` 建行 |
| ③ CHANGELOG ≥ 1 条 | `CHANGELOG.md` 两条 PATCH 级条目（在 `[Unreleased]` 下，理由见「现状说明」） |
| ④ GOPRIVATE 前提 | `AGENTS.md` §2 的「私有化前提」那一行（护栏按 `GOPRIVATE=github.com/ruoxizhnya/` 子串核对）+ `README.md`「拉取前提」两条细则 |

**⑤ 验收逐条（真退出码）**

1. **README 存在版本规则表与兼容矩阵**：`go run ./tools/check_release_matrix` → 退出码 **0**，
   输出「规则表 3 行 ｜ 矩阵序列 0 ｜ CHANGELOG 版本节 0（条目 2）｜ tag 0 ✅」。（同一判据由 CI 步骤「版本纪律核对」执行。）
2. **CHANGELOG ≥ 1 条**：条目数由脚本数 `- ` 行，输出里的「条目 2」即判据；删空 → 红（夹具腿 `TestChangelogWithoutEntriesIsRed`）。
3. **AGENTS.md 含 GOPRIVATE 行**：脚本核对子串存在；删掉 → 红（`TestGoprivatePrerequisiteMissingIsRed`）。
4. **破坏验证（现场，真退出码）**：往 `CHANGELOG.md` 追加一节 `## [v0.2.0] - 2026-10-22`（矩阵不动）→
   同一条命令退出码 **1**，输出点名「CHANGELOG 有 v0.2.0 这一版，矩阵却没有 v0.2 那一行」→ 按 sha256 逐字节还原 → 退出码回到 **0**。
5. **破坏验证（夹具腿，长期在案）**：`TestReleaseWithoutMatrixRowIsRed`（本单主判据）、
   `TestMatrixRowWithoutChangelogIsRed`（反向：矩阵留了没账的序列）、`TestTagWithoutChangelogSectionIsRed` 与
   `TestTagWithoutMatrixRowIsRed`（tag 有版本没账）、`TestMissingMatrixMarkersIsFatal`（删标记让检查通过）、
   `TestEmptyMatrixDataRowsIsFatal` 与 `TestMissingTableSeparatorIsFatal`（清空矩阵 / 删表头分隔行让「跳过表头」的逻辑变松）、
   `TestMissingRuleRowIsRed`、`TestMalformedMatrixCellIsRed`（打错的 `v0.1.0-beta` 不能被当信息行放过）、
   `TestDuplicateMatrixSeriesIsRed`（`v0.1` 与 `v0.1.0` 同序列两行）、`TestReleaseTitleWithoutDateIsRed`、
   `TestListTagsFailureIsFatal`（读不到 tag 就红，不许静默按零 tag 走下去）。

**现状说明（写进文档的口径，不是缺陷）**

- 本仓**零符号、零 tag**，所以兼容矩阵今天只有一行「`main`（未发版）」的基线行；首个版本 `v0.1.0` 由 QT-07 ③-S2-a 在 S2 打。
  「矩阵里有没有我的版本」对当前所有使用方的答案都是「没有，你在用 `main`」。
- 护栏的 **tag 那条腿今天是空跑**（0 个 tag）。它不是靠这趟空跑产生价值的，价值在 `TestBaselineWithoutReleaseIsGreen`
  （证明零发版不被误判成违规）与 4 条 tag 夹具腿；S2 打上第一个 tag 那天它才开始真正咬人。
  CI 侧的前提（`fetch-depth: 0`）已一并改掉，否则那条腿在 CI 里永远是绿的假象。
- 跨仓那半（「三个仓的 require 是否都落在矩阵声明范围内」）**没做也做不了**：本仓 CI 不 checkout 别的仓（R2），
  已在 README 写明归属 QT-09。

**执行中如实登记的偏差**

- QT-19 ③-③ 的矩阵口径写的是 `vX.Y`，而 tag/CHANGELOG 的形状是 `vX.Y.Z`。本单按**序列建行**落地（一行覆盖该序列的补丁版），
  脚本用「序列对齐」而不是「字符串相等」核对，并把两种形状都收进归一化 —— 否则 S2 打第一个 `v0.1.0` 时就会因为
  「矩阵写 `v0.1`、账写 `v0.1.0`」而红，那是格式摩擦不是纪律。
- QT-19 ⑤ 的「CHANGELOG 至少有一条记录」在「还没发版」的现实下只能是 `[Unreleased]` 下的条目；
  本单没有为了凑一条「版本节」而虚构一个没发过的版本。
- CT-02 ③-④ 只要求 AGENTS.md 一行；本单同时在 README 写了「拉取前提」两条细则（proxy 绕过之外还要 git 凭证），
  因为只配 `GOPRIVATE` 不配凭证同样是拉不到 —— 出处是 `repo-inventory.md` §1.1 的可见性风险注。

**边界**：没引 release 工具链（goreleaser 等一个没加，`go.mod` 零变化）；没写任何契约符号；没碰 `docs/{db-schema,api-schema,dataflow}.md`；
「禁删 deprecated」的机械护栏留给 QT-09 ⑤；kit 的同名交付留给 KT-02（本单只把**版本规则表**的口径在 README 里说清了出处，没替 kit 写文档地图）。
