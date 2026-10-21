# TASKS 归档 — tellurion-contracts（2026-10-21）

> 规则见 [`README.md`](README.md)：完成即在本仓台账标 ✅，随后整条搬进本目录，**旧代正文不改写**（要订正只在文件头加说明）。
> 本文件是历史证据，不作为当前依据引用；引用前先确认它没被新文档取代。

---

## CT-01 · 随仓接口规格载体 `docs/API.md`（导出面清单 + 稳定性 + 错误语义）[P1] ✅ 完成（2026-10-21）

- **① 目标**：本仓导出面有随仓、可自证的接口规格，下游（core/agent）不读 core 实现即可知调用形状与错误语义。
- **② 上下文**：refactor-plan-review **M6**（报告：`tellurion:docs/design/qt-decoupling/refactor-plan-review-2026-10-20.md`）；上游编排单 `tellurion:docs/TASKS.md` QT-07 ③-S2-c；冻结面指回 `tellurion:docs/SPEC.md` 四契约节（改冻结面走 `tellurion:docs/TASKS.md`「维护本文件」的冻结 schema 变更评审硬规则）。现状导出面清单散在 `tellurion:docs/design/qt-decoupling/kernel-interface-map.md`，无随仓载体。
- **③ 要求**：① 建 `docs/API.md`：逐项列导出符号——四契约（`engine/`）+ port + 值对象（`domain/`）+ 哨兵错误族，每项标**稳定性标注**（frozen / stable / experimental）与**错误语义**（`errors.Is` 判定、禁字符串匹配）；② 每个冻结面指回 `tellurion:docs/SPEC.md` 对应**节号**（只指节号，不复制正文）；③ 护栏进 CI（AST 扫描）：新增导出符号未登记进 `docs/API.md` → 红。
- **④ 约束与非目标**：不写内部函数/实现细节；**不做手写 DTO**——跨进程 DTO 是 `tellurion:docs/openapi.yaml` 的生成物、归 core（见 `tellurion:docs/design/qt-decoupling/repo-inventory.md` §3），本仓 `dto/` 只收生成物且禁手改。
- **⑤ 验收**：`go doc ./...` 导出面与 `docs/API.md` 清单一一对应（比对脚本跑通）；**破坏验证**：故意导出一个未登记符号 → CI 红（留证后撤回）。
- **⑥ 边界**：只做接口规格载体与护栏；版本规则/兼容矩阵是 CT-02。

### 结案记录（2026-10-21）

**交付物**

| 文件 | 内容 |
|---|---|
| `docs/API.md` | 新建。§1 核对判据（粒度 + 四个红向 + **覆盖面限制**）；§2 稳定性词表；§3 错误语义六条；§4 **机器核对的登记区**（今天为空，见下）；§5 冻结面逐符号规格（四契约 / port 方法集 / status 与 invalid Code / 哨兵族，每条指回 SPEC 节标题）；§6 值对象分组面；§7 abi 常量族；§8 不登记的东西 + 三条待裁决；§9 与其他文档的互锁 |
| `tools/check_api_surface/` | 新建。`go/ast` 扫描包级导出符号（func/type/const/var + `Type.Method`），与 `docs/API.md` 的 `<!-- surface:begin/end -->` 区双向对拍；`package main` 与 `_test.go` 不参与 |
| `tools/check_api_surface/main_test.go` | 13 条用例（含夹具 `testdata/ok/`） |
| `.github/workflows/ci.yml` | Vet 与 Test 之间加「导出面核对」步骤 |
| `AGENTS.md` / `README.md` / `docs/README.md` | 登记第四份真相源与变更触发行、命令 |
| `docs/TASKS.md` | CT-01 标 ✅ 并把六字段正文搬进本文件，台账只留指针；§一 引言补「QT-05 落码时必须同步 §4 登记区」 |
| `docs/api-schema.md` | 自审时回源码核对那两处后订正：§1.3 加「三条哨兵零定义」的现状标注、§1.4 的 `13 个 host 函数` 改回 **12**（v1 8 + v2 4），frontmatter `last-verified` 一起更新 |

**⑤ 验收逐条**

1. **一一对应**：`go run ./tools/check_api_surface` → 退出码 **0**，输出「2 个代码包 / 0 个登记包 ｜ 代码符号 0 ｜ 登记符号 0 ✅」。
   ⚠️ **本仓今天没有库**（唯一的库文件是根包 `doc.go` 占位，零导出符号；`tools/` 是 `package main` 不计入），
   所以「一一对应」目前成立的形态是**空对空**。护栏的实质内容在 2/3 两条与 13 条单测里，不在这一行输出里；
   S1（QT-05）落码后 §4 区才有符号可登记。
2. **破坏验证（现场，真退出码）**：临时新建 `scratch_destruction_leg.go`（`package contracts` + `func ScratchUnregistered() {}`）→
   同一条 CI 命令退出码 **1**，输出点名 `+ ./ScratchUnregistered` → 删除该文件 → 退出码回到 **0**，`git status` 只剩本次交付的新文件（无 scratch 残留）。
3. **夹具里的同一条腿长期在案**（不需要每次都手改代码）：`TestUnregisteredExportedSymbolIsRed`（新增导出符号 → 红）、
   `TestUnregisteredInterfaceMethodIsRed`（给接口加方法 → 红）、`TestStaleRegistrationIsRed`（登记了不存在的符号 → 红）、
   `TestEmptyRegisterIsNotAVacuousPass`（把登记区删空 → 全量红）、`TestMissingMarkersIsFatal`（连标记一起删 → 红）、
   `TestSyntaxErrorIsFatal` / `TestTwoPackagesInOneDirIsFatal`（解析失败 / 一目录两包 → 红，防「包凭空消失后静默通过」）、
   `TestRootPackageIsRegistrable`（根包符号必须能登记，见下面的自审补记 #2）。

**执行中如实登记的偏差（不是缺陷，是给后续单的输入）**

- QT-05 ③-① 说 `engine/` = **6 文件**，实测 core `pkg/backtest/contracts` 今天是 **9 个非测试文件 / 190 个导出符号**（`abi.go`、`canonical.go`、`provenance.go`、`filltiming.go` 等是 K9~QT-02 之后长出来的）。QT-05 落码前那个数得重测。
- `ARCHITECTURE.md` §2 列了 `errors/` 与 `abi/` 两个包，QT-05 ③-① 只说迁 `domain/`/`engine/`/`fees/` —— 口径不一致，登记在 `docs/API.md` §8 待裁决 #1。
- SPEC.md 里**没有** `SignalSpec` 节（实测 0 命中），四契约的第三个只能指向 `api-schema.md` §1.2 与 `kernel-interface-map.md`，已在 §5.1 写明。
- 三分法哨兵 `ErrDataGap`/`ErrStrategyFault`/`ErrContractBreach` **源码里今天不存在**（QT-10 ③-① 才落），§5.4 标成「⬜ 目标」而不是现状。
- 实测「agent 侧 0 命中」的 `pkg/domain` 根包符号是 **86 / 124**，比 QT-05 ③-⑨ 点名的 5 族宽得多；命令与读法登记在 §6，完整清单仍属 QT-05 ③-⑨ 的产出物（本单不替它下结论）。

**边界**：没碰版本规则/兼容矩阵/CHANGELOG/GOPRIVATE（那是 CT-02）；没写任何 DTO；没动 `dto/`、`schemas/`。

### 自审补记（同日，`fc6eb57` 之后逐条通读 diff 的结果）

原则：叙述与实现逐句核对，每个数字回源码取。改掉 6 处，其中 2 处是我自己刚写进去的。

| # | 发现 | 处置 |
|---|---|---|
| 1 | **我自己写出的过时声明**：`docs/API.md` 顶部与 §4 都写「本仓唯一的 Go 文件是根包 `doc.go`」——交付 `tools/` 之后这句话当场失效 | 改成「唯一的**库**文件是 `doc.go`；`tools/` 是 `package main` 不计入导出面」；§4 的「零个包对零行登记」同步改成实测的「0 个库导出符号对 0 行登记」 |
| 2 | **不可满足的护栏**：`parseRegister` 只认带 `/` 的包目录标题，根包（键为 `.`）的导出符号**没有任何登记入口**——谁在根包加一个导出符号，CI 就永久红 | `dirInTick` 增认 `` `.` ``；夹具加 `placeholder.go`（根包 `RootValue`）+ `TestRootPackageIsRegistrable` 正反两腿（认它 → 绿；标题不带目录 → 报「未登记」且**不**算失效项）；§1 加「登记表形状」一行把这个写法钉成文 |
| 3 | `api-schema.md` §1.4 的「abi v2 = **13** 个 host 函数」是 core 已在 2026-10-21 订正掉的错数（`adr-029-amendment-host-api-v2.md` §2.1 明写原枚举 8+4=12，代码侧判据是 `HostFnCount()`） | 改回 12 并给出家族拆分，指回那条订正；顺带在 `API.md` §7 逐条列全 12 个名字（v1 8 / v2 4），注明「真要补第 13 个要先改表走 abi 冻结」 |
| 4 | `api-schema.md` §1.3 用现状口吻列三条哨兵，读者会去 import 不存在的符号 | 表下加「1 处命中（`abi.go:646` 的注释，且不含 `ErrDataGap`）/ 定义 0 处」的现状标注，指回 `API.md` §5.4 |
| 5 | `API.md` §7「值类型 / 字段号」那行写「未知 ⇒ `ErrAbiVersionUnknown` 类」——核 `abi.go` 后是错的：`FieldID.Valid()` / `Action.Valid()` 只返回 bool，`HostValueType` 根本没有 `Valid()` | 按三种真实判据重写（bool / `ErrActionIllegal` 只在 `Action.String()`+`Direction()` 路径 / `HostValueType` 由签名比对拦），并在 abi 版本行补上 `HostAbiVersionFromName` ⇒ `ErrAbiVersionUnknown` |
| 6 | `API.md` §5.2 的 `Provider` 行把「缺数据归 `ErrDataGap`」写成今天的判据 | 改成「今天返回未分类 error；QT-10 ③-② 落地后才归 `ErrDataGap`」 |

**核对过、无需改的断言**（列出来是为了下次不必重跑）：`domain.BacktestResult` 确实无 `Status` 字段；6 条 `Invalid*` Code 的字面值与 §5.3 一致；`EngineRunner` 在 `contracts.go:366`、`BacktestRequest` 在 `:171`、`BacktestResponse` 在 `:202`、`SignalWeight` 在 `:160`、`SignalProvider` 在 `signal_provider.go:41`；`market.Provider` 就是 5 个方法（`kernel-interface-map.md` 列的 10 个是 core `pkg/marketdata` 那个同名不同物的 Provider）。

**门禁复跑（干净 checkout，`fc6eb57` + 本次自审）**：`gofmt -l` 空、`go build` / `go vet` / `go test ./... -count=1` 全 0（13 条用例 PASS）、`go run ./tools/check_api_surface` 0。

补记 #2 的通路也做了现场三腿（不只在夹具里）：真实仓根包加 `func ScratchRootLeg()` 且不登记 → **exit 1**；
在 §4 区临时加一个 `` `.` `` 小节登记它 → **exit 0**，输出「1 个登记包 ｜ 代码符号 1 ｜ 登记符号 1」；
撤掉文件与那一节 → **exit 0**，`grep -c ScratchRootLeg docs/API.md` = 0，`git status` 只剩本次自审的改动。
