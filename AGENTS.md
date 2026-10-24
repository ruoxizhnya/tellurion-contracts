# AGENTS — tellurion-contracts

> 本文件是**本仓**的唯一入口。判据：一个 agent 只 checkout 本仓，就能知道「怎么构建、边界在哪、依赖哪版契约、活在哪、禁止什么」，**不需要读别的仓的源码**。

## 1. 这是什么仓

跨仓唯一缝隙：**port 接口 · DTO/值对象（含 domain）· 四个引擎契约 · 哨兵错误 · abi 常量**。所有跨边界的东西只能住这里。

**状态**：S1 落代码已完成并已推（QT-05 ③-①，2026-10-22：`fees` / `domain/market` / `domain` / `engine` 四包、429 个导出符号）／2026-10-23 一天四版：`v0.1.0`（QT-07 首版）→ `v0.1.1`（把 `conformance/` 语料打进 tag）→ `v0.2.0`（429→435）→ `v0.3.0`（435→439，QT-10 三类执行错误哨兵）／`v0.4.0`（2026-10-24：QT-20 参数投递入口 `StrategyRef`/`StrategyParams` + `BacktestParams.Params`；QT-41b ④ 编译 / 资产四哨兵；439→443）／`v0.5.0`（2026-10-24：QT-22 `version_hash` canonical 六项 `StrategyVersionInput` + `DomainStrategyVersionHash`；443→447）——core / agent 已于 2026-10-24 由 **G13** 切至 `v0.4.0`（strategy 同日更早切）；`v0.5.0` 于同日由 **G4** 回写：core / agent / data / strategy 四仓 `go.mod` 实测已切，矩阵 `v0.5` 行已改「已切」／`v0.6.0`（2026-10-24：QT-40 X1a 新增 `dataprotocol/` 包 —— `UnifiedDataPoint` 七字段 + `ToJSON`/`DeduplicateKey` + 15 个 `DataType*` 常量；447→465，**`main` 与 tag 均已推远端（2026-10-25）**）

## 2. 构建 / 测试（必须单仓可跑）

```bash
go build ./...
go vet ./...
go test ./...
go run ./tools/check_api_surface      # 导出面与 docs/API.md 登记区一一对应（CI 同款，CT-01 护栏）
go run ./tools/check_release_matrix   # 版本规则表 + 兼容矩阵 ↔ CHANGELOG ↔ tag（CI 同款，CT-02 护栏）
go run ./tools/check_contract_shape   # 无 *Config/*Options/*Impl 形状名、无 DB/驱动 import（CI 同款，QT-09 ⑤）
go run ./tools/check_surface_freeze   # 导出面形状 ↔ docs/api-surface.txt 逐行一致（CI 同款，QT-09 冻结）
go test ./conformance                 # 跨仓契约语料的自检腿（QT-09 ③；core/agent 侧读同一份 cases.json）
```

- 前三条是 R1/R2 的判据：**单独 checkout 就能绿**；后五条是本仓的随仓护栏（CI 同款步骤）。
- `conformance/cases.json` 是跨仓契约一致性的**唯一真相**：改用例或期望值要三仓同测，且期望值只能来自本仓实现（不许按某一侧反推）。本包刻意零导出符号，别往里加导出类型/函数。
- 导出面签名改了：先登记 `docs/API.md` §4，再 `go run ./tools/check_surface_freeze -update`；只改基线不登记，`check_api_surface` 会挡住新增符号。
- PR 正文模板（`.github/PULL_REQUEST_TEMPLATE.md`）的四段必填项由 `tools/check_pr_body` 在 `pull_request` 事件里核对，其他事件显式跳过。
- 若某一步需要别的仓的源码目录才过，那是**边界画错了**——停下来修边界，**不要**加 `replace ../xxx` 绕过。
- 本机同时改多个仓时可用 `go.work` 联调，但**它不入库**（`.gitignore` 已挡）。
- **私有化前提**：本仓一旦被设为 private，构建与全部 CI 环境必须配 `GOPRIVATE=github.com/ruoxizhnya/*` **加** git 凭证（PAT / SSH），否则 `go get` 走公共 proxy 必然拉不到（风险注见 `tellurion:docs/design/qt-decoupling/repo-inventory.md` §1.1/§5；现状 12 仓全 public）。细则与发版流程在 `README.md`「拉取前提」「发版流程」。

## 3. 边界（硬约束）

**允许**

- 被 core / agent / 各服务 `require`
- 只放**被 ≥2 个仓使用**的跨边界类型与接口

**禁止**

- 任何实现、业务规则、DB 访问
- `Config` / `Options` / `Impl` 命名（契约泄漏实现细节的形状）
- 只有单个仓使用的类型（那条属于泄漏，应留在使用方）

> **`*Config` 族的例外判据**：不按名字判、按类型在边界上的角色判——见 [`docs/adr/0002`](docs/adr/0002-config-族的边界处置.md)。
> 简版：冻结文本 / 跨语言契约里的名字保留；core-only 的搭车类型不迁或迁后回收；只有「长期留在本仓 + 名字泄漏实现细节」才改名（破坏性）。
> **新备案的例外必须写进那条 ADR**，不许在 commit message 里口头豁免。
>
> **「零实现」的判据**：值对象自带的纯函数方法（校验 / 哈希 / 公司行为 Apply）不属禁止项；四条判据（无 I/O、依赖受限、无跨调用持久状态、不替代服务）见 [`docs/adr/0004`](docs/adr/0004-值对象上的纯函数行为随类型走.md)。回测 / 撮合 / 存储 / 调度 / 服务仍一律禁止。

下面三条对所有仓都成立：

- **禁止** `replace` 指向仓外路径（单仓就不可构建了）。**唯一过渡例外**：`tellurion-agent` 的 `replace tellurion => ../tellurion`（拆仓过渡态，S3 后收回），以 `tellurion-agent:docs/adr/0001-编译期依赖-core-是过渡态.md` 为准。
- **禁止**提交 `go.work` / `go.work.sum`（本地产物；已在 `.gitignore`）。
- 新增第三方依赖前，先在第 4 节登记，再改 `go.mod`；没登记就加依赖 = 绕过评审。

## 4. 依赖的契约版本

| 依赖 | 版本 | 说明 |
|---|---|---|
| `github.com/ruoxizhnya/tellurion-kit` | **今天不 require**（2026-10-22 对 core 侧四个待迁包实测） | 原写「只用它的 decimal / errors」——实测迁移输入（`tellurion:pkg/fees`、`pkg/domain`、`pkg/domain/market`、`pkg/backtest/contracts`）**一个 kit 符号都不碰**：非标准库 import 只有 testify（测试），仓内 import 只在这四个包彼此之间。要用 kit 是一次显式 MINOR + 一条 ADR（例如把 `float64` 换成 `kit/decimal`），不是默认前提 |

第三方依赖（QT-05 迁入时登记；**新增先改这张表，再改 `go.mod`**）：

| 依赖 | 版本 | 用途 | 引入包 |
|---|---|---|---|
| `github.com/stretchr/testify` | v1.11.1 | **仅测试**断言（随 core 的 `pkg/domain` / `pkg/backtest/contracts` 测试一起迁入；`domain` 与 `engine` 均已落地） | `domain` / `engine` 的 `_test.go` |

> testify 的间接依赖随它自动进 `go.mod`（`davecgh/go-spew`、`pmezard/go-difflib`、`gopkg.in/yaml.v3`，均 `// indirect`）；`go.sum` 里另有 `gopkg.in/check.v1`（testify 的测试依赖，不进构建）。它们不是独立决策，但**一并登记**，免得下次 tidy 时以为被谁偷偷加了。
> 版本与 core 当前 `go.mod` 对齐（v1.11.1），避免迁移期同一个库在两个仓里出现两个大版本。
> `fees` / `domain/market` 的测试只用标准库，**不引 testify**。

> 跨仓依赖**只认版本号**（伪版本或 tag），不认源码路径。契约变更走 `tellurion-contracts` 的发布，不走"顺手改一下"。

## 5. 活在哪、文档在哪

- **本仓台账**：`docs/TASKS.md`（唯一）。**没有单号别开工。**
- 文档导航：`docs/README.md`；边界与依赖方向：`docs/ARCHITECTURE.md`；决策：`docs/adr/`。
- 跨仓背景（只读，不在本仓改）：`tellurion:docs/design/qt-decoupling/repo-topology-and-naming.md`、`tellurion:docs/TASKS.md`。

> 跨仓引用一律写 `repo:path` 纯文本，**不做可点击相对链接**——独立 checkout 时 `../../` 必断，R2 会永远红。

## 6. 文档真相源（改功能前先看这里）

本仓的「功能边界」由三份文件钉住，**它们是唯一真相源，不是参考读物**：

| 维度 | 文件 |
|---|---|
| 数据形状（我拥有哪些表、谁能写） | `docs/db-schema.md` |
| 接口形状（对外暴露什么、错误码、幂等） | `docs/api-schema.md` |
| 流向与顺序（数据怎么走、哪里 fail-closed） | `docs/dataflow.md` |

**本仓还有第四份**：`docs/API.md` —— 导出符号的清单（稳定性标注 + 错误语义 + 冻结面指回的 SPEC 节号）。
它与别的仓的三份不同，因为它是**机器核对**的：新增导出符号没登记，`go run ./tools/check_api_surface` 当场红（CT-01）。

**第五份是版本纪律三件套**：`README.md` 的「版本规则 / 兼容矩阵」+ `CHANGELOG.md` + `git tag`。
同样是机器核对的：发版只写 CHANGELOG 不更新矩阵（或反过来、或 tag 打在没账的版本上），`go run ./tools/check_release_matrix` 当场红（CT-02）。

**变更触发是硬约束**（完整表见 `docs/README.md` §3）：改表 → 改 `db-schema.md`；改接口 → 改 `api-schema.md`；加/改/删导出符号 → 改 `API.md` 登记区；发版 → 改兼容矩阵那一行 + `CHANGELOG.md` 那一节；改链路 → 改 `dataflow.md`；改分层 → 改 `ARCHITECTURE.md` + 一条 ADR；删档/迁移/换工具 → 一条 ODR。PR 里代码动了而对应文档没动，应被打回。
