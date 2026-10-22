---
status: active
last-verified: 2026-10-22
verified-by: QT-05 ③-① 本仓侧收官（四包 429/429 + §8 #1~#3 裁决 + 下游腿交接单）+ 自审订正（ADR 文件名死链、方法类别、EngineRunner 返回值、fees 注释失真；新增 ADR-0004）；CT-02 结案归档（版本规则表 + 兼容矩阵 + CHANGELOG + check_release_matrix 进 CI）；A2 推送回写（6 个 commit 已推 `origin/main`，该批 HEAD = `fbe9a2d`，`go get` 实测可解析，CI 429/429；§8 #4 登记为 X1a 待办）
---

# TASKS — tellurion-contracts 未完成项

> 本仓**唯一**任务台账。跨仓任务只在 §二 留指针，**不镜像正文**（双写必然漂移）。
> 状态：`⬜ 待做` / `🔶 进行中` / `⛔ 阻塞（有未解除前置）` / `✅ 完成（附日期）`
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

### QT-05 · S1 契约层与工具层落地（代码直接写进两个新仓）[P0] 🔶 本仓侧已完成并已推（`fbe9a2d`）；下游腿（③-③④⑤⑦⑨）待 core/agent 执行

> **迁入说明（2026-10-21）**：本单自 core 仓 `tellurion:docs/TASKS.md`「S1」迁入（任务台账批量重分布：S1 双仓落点，主单落本仓）。六字段正文逐字保留、不改写；**kit 腿**拆分为 `tellurion-kit:docs/TASKS.md` 的 **QT-05-B**（③-② 的执行单，判据以本单为准）；core 侧留同号指针，core 侧改动（③-③/④/⑥/⑦/⑨）随本单执行。执行顺序（W0 全绿后的 S1）见 core `tellurion:docs/TASKS.md` 执行顺序总览。

> **进度（2026-10-22，本仓侧已完成）**：`fees/` ✅（13 符号）、`domain/market/` ✅（46）、`domain/` ✅（176）、`engine/` ✅（194）四包全落，登记进 `docs/API.md` §4，两条随仓护栏绿（导出面 **429/429**）。§8 待裁 **#1 已裁决并写成 [`adr/0003`](adr/0003-哨兵与-abi-不独立成包.md)**（errors//abi/ 随 engine/，不拆包；61 文件的旧估计实测证伪）、**#2 已裁决（[`adr/0002`](adr/0002-config-族的边界处置.md)）**、**#3 已裁决**（fees 整包）；只剩 **#4**（dataprotocol 落点时机，已登记为 X1a 待办）。
> ✅ **push 阻塞已解除（A2，2026-10-22）**：本仓 6 个 commit 已推到 `origin/main`（该批 HEAD = `fbe9a2d`）；仓外空模块里 `GOPROXY=https://proxy.golang.org,direct go get github.com/ruoxizhnya/tellurion-contracts@fbe9a2d` 实测解析为 `v0.0.0-20261022191510-fbe9a2d0df17`。CI 全绿（导出面 429/429）。剩下的全部是下游腿（③-③④⑤⑦⑨），由 core / agent 执行；③-④ 仍禁 `replace`。
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

1. `go get` 两个仓（contracts + kit）；**它仍保留 `replace tellurion => ../tellurion`**（S3 才删，见该仓 ADR-0001），但**不得**为 contracts/kit 加 replace。
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
| QT-40（T2） | `dataprotocol` 要不要独立成包、何时迁 | 单在 `tellurion-data:docs/TASKS.md`（QT-40 的接口对表待办 T2）；**落点与两条读法登记在本仓** `docs/API.md` §8 待裁 **#4**。⇒ **已登记为 X1a 待办（2026-10-22，A2）**：X1a 开工前必须在 §8 #4 挑定 (a)/(b) 并写回 §8 #4 与设计文档 §4.1 注 ②；X1a 的施工序在 `tellurion-data:docs/TASKS.md` 施工序表 |

---

## 维护本文件

- **新增单**：写全六字段；六字段里必须显式引用**冻结契约的路径与节号**（契约仓的包路径、设计文档的节号），不许写"参考相关实现"。
- **派工前先证伪前提**：读码核对工单写的前提；不符就停下回报 + 给拆法，不硬做、不扩范围。
- **跨仓引用写 `repo:path`**（如 `tellurion:docs/TASKS.md`），不做可点击相对链接。
- **顺序纪律**：同一波次内撞同一批文件的单按单号串行，一单一 commit。
