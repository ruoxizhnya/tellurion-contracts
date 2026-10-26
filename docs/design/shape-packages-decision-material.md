---
status: active
last-verified: 2026-10-26
verified-by: A 组决策材料（只出材料，0 个 .go diff；三张候选六字段单待主任按序立单）
---

# A 组 · 三个跨仓形状包决策材料（risk / 编排 / broker DTO）

> **本文只出材料，不是决策，也不开工**。contracts 与 kit 仓当前**没有任何已立单**（`tellurion-contracts:docs/TASKS.md`
> 的未完成标题只剩 QT-05 的下游腿指针，`tellurion-kit:docs/TASKS.md` 同理），按「没有单号别开工」不得动手。
> 本文回答的是：**主任若要立这三单，前提是什么、顺序为什么不能换、禁区在哪、验收怎么机械判定。**
> 证据口径：`repo:path:line` 为 2026-10-26 实测。跨仓引用写成纯文本，不做相对链接。

## 0. 结论先行

| # | 候选单 | 落点 | 今天能开工吗 | 硬前置 | 建议版本 |
|---|---|---|---|---|---|
| **①** | `contracts/risk` 包：`Snapshot` / `SnapshotRule` / `Scope` / `RuleType` + 冻结词表 + hash 纯函数 + 4 哨兵 | `tellurion-contracts/risk/`（新包） | 🔶 **能开工**（2026-10-26 复测：原判「⛔ 不能」的前提已失效） | ~~risk `QT-74` 参数词典定稿（今天 ⬜）~~ → **词典已落**：`tellurion-risk:internal/rule/params.go` 246 行 / 14 标量，写侧 `model.go:147` `validateParams`、读侧 `snapshot.go:163` `rejectUnknownParamKeys`；risk 台账 `docs/TASKS.md:508` 记 **🔶**（前两步已落，未落项是「迁 contracts」本身 + core 消费 + 数值换算口径）。剩的前置只有**没人立单** | MINOR（`v0.7.0`） |
| **②** | 编排形状：claim / `run_artifacts` / 子 run 入队 | `tellurion-contracts/run/`（新包，仅**消费方需要的**三个 wire 形状） | ⛔ **不能** | 先裁 §4 的「wire DTO 归属」冲突（下面第 1 条） | MINOR（`v0.8.0`） |
| **③** | broker DTO：order intent / fill / account | `tellurion-contracts/broker/`（新包） | ⛔ **不能**（顺序已裁定，不再是死锁） | `SVC-BROKER-G1` 先把订单 HTTP 面定形（`tellurion-broker:docs/TASKS.md:37` 现为 **⛔**）——⚠️ 该单原文把「等 A 组 broker DTO」当前置，与本行构成**互为前置的死锁**；2026-10-26 已按 §3.1.1 裁定**G1 先行**，③ 的反向依赖已解除 | MINOR（`v0.9.0`） |

**三条判据支撑上面的排序**（每条都有实测出处，见 §1~§3）：

1. **① 的顺序是「词典先于形状」——顺序判据不变，但前置已于 2026-10-26 满足**：本行原文写「risk 的 `params` 今天**没有词典**」，
   那是 `QT-74` 落码前的状态。现状：`canonicalParams`（`tellurion-risk:internal/rule/model.go:150-174`）确实**仍只做规范化**
   （键有序 / 无多余空白 / ≤64 KiB，不看键名），但键名与单位已有词典并在两侧执法——`internal/rule/params.go`（246 行、
   14 个标量的字段名 / 单位 / 取值域），写侧 `model.go:147` 调 `validateParams`（未知键 / 缺必填 / 越域 / NaN 当场拒），
   读侧 `snapshot.go:163` 调 `rejectUnknownParamKeys`。所以「迁了等于没迁（只是个装任意 JSON 的信封）」这条**不再成立**。
   `QT-74` 现在在 `tellurion-risk:docs/TASKS.md:508`，状态 **🔶**（前两步已落；未落的三项里只有「数值换算口径」
   与单位语义有关，且 risk 已把权重冻结为**小数**，换算差异属 core 侧映射、不改形状）。⇒ ① 缺的是**立单**，不是前置。
2. **② 的顺序是「先裁口径，再迁形状」**：本文逮到一条**结构性冲突**——本仓 `docs/API.md:766` 明写
   「跨进程 DTO 是 `tellurion:docs/openapi.yaml` 的**生成物**，本仓只收生成物且**禁手改**」，
   而 ② ③ 要迁的**恰恰是 wire DTO**，且六个服务仓**一个 openapi 文件都没有**（实测：12 仓 glob
   `*openapi*` / `*swagger*` 只命中 core 的 `docs/openapi.yaml` 与 `cmd/analysis/handlers_openapi.go`）。
   ⇒ 不先裁这条，② ③ 无论怎么写都会被本仓自己的规则判成违规。**这是本组最重要的发现**，详见 §4。
3. **③ 的顺序是「有 HTTP 面才有 wire 形状」**：broker 今天**只有两条路由**（`GET /healthz`、`GET /version`，
   `tellurion-broker:cmd/brokergateway/main.go:45`、`:55`），订单 / 成交 / 账户 DTO **一次都没被 JSON 序列化上 wire**
   （全仓 `json.Marshal` 只有 3 处：`internal/ackstore/pg.go:96` 写 `ack` JSONB、`:53` 读回、`model.go:96` 算指纹）。
   现在迁 = **在 HTTP 面存在之前冻结报文形状**，而 `SVC-BROKER-G1` 还没开工（`tellurion-broker:docs/TASKS.md:37`，现标 **⛔**）。
   ⚠️ 但那条 ⛔ 写的理由是「等 A 组 broker DTO」——与本判据互为前置，构成死锁；顺序见 **§3.1.1**。

**一次只发一包**：三单**不得同批**。理由不只是纪律——① 要改 `tools/check_contract_shape` 的 `contractRoots`
（见 §1.5），② ③ 各要改 `docs/API.md` §4 登记区与 `api-surface.txt` 基线；三单同批会让「这次红是谁造成的」失去答案。

## 0.1 派工前提复核（本轮实测，可直接引用）

| 断言 | 实测命令 | 实测输出 | 结论 |
|---|---|---|---|
| contracts 顶层无 `risk/` | `git -C tellurion-contracts ls-tree --name-only HEAD` | `.gitattributes .github .gitignore AGENTS.md CHANGELOG.md LICENSE README.md conformance dataprotocol doc.go docs domain engine examples fees go.mod go.sum tools` | ✅ **成立** |
| contracts 最新 tag = `v0.6.0` | `git -C tellurion-contracts tag \| tail -1` | `v0.6.0` | ✅ **成立**（并已于 2026-10-26 实测**远端也有**：`git ls-remote --tags origin` → `refs/tags/v0.6.0` = `ec00ef9` / peel `5ee7fae`；kit 同法核得 `v0.1.4` = `f1d24df` / peel `8077c36`。台账 frontmatter 原文写的「待 push」已订正） |
| kit 最新 tag = `v0.1.4` | `git -C tellurion-kit tag \| tail -1` | `v0.1.4` | ✅ **成立** |
| 两仓工作树干净 | `git -C <repo> status --porcelain` | 空（两仓皆是） | ✅ **成立** |
| 两仓 main 与 origin 同步 | `rev-list --count origin/main..HEAD` / `HEAD..origin/main` | ~~`0` / `0`（两仓皆是）~~ → **2026-10-26 复测不成立**：contracts 领先 **1** 笔（`3055b6e`）、kit 领先 **1** 笔（`ba42add`），两边 behind 均为 0。**本批订正提交后 contracts 会再 +1**（推送由主任执行，本材料不代推） | ⚠️ **已过期**：tag 在远端不等于分支已同步；判据现跑 `git -C <repo> rev-list --count origin/main..HEAD` |
| contracts 无 `risk` 形状 | `git -C tellurion-contracts ls-tree -r --name-only HEAD \| grep '(^\|/)risk(/$\|$)'` | 零命中 | ✅ **「risk 形状已进契约」被证伪** |

> **验收判据的原文照收**（派工给的）：`git -C tellurion-contracts tag | tail -1` = `v0.6.0`；`ls tellurion-contracts`
> 输出无 `risk` ⇒ 证伪「risk 形状已进契约」。上面两行就是它的复现。
> **破坏验证**：派工原文写「不适用（无实现可破坏）」。本文仍做了一条**前提守卫的破坏验证**（改坏 → 红 → 逐字节还原 → 绿），
> 因为「无实现」不等于「无可证伪」——详见 §6。

## 1. ① `contracts/risk` 包

### 1.1 要迁什么（逐项，附出处）

| 内容 | 今天在哪 | 迁不迁 | 判据 |
|---|---|---|---|
| `RuleType` + 4 个取值 | `tellurion-risk:internal/rule/model.go:12`、`:14-19`（`position_sizing` / `stop_loss` / `regime` / `take_profit`） | ✅ 迁 | 冻结枚举，跨仓消费方必须按同一词表 fail-closed |
| `Scope` + 4 个取值 + `Valid()` | `model.go:29`、`:31-36`（`global` / `account` / `strategy` / `symbol`）、`:47-50`（`Valid()`） | ✅ 迁 | 同上；`Valid()` 是纯函数（ADR-0004 允许） |
| `SnapshotRule`（6 字段） | `internal/rule/snapshot.go:24-31`（`RuleID` / `Revision` / `Scope` / `Type` / `Params json.RawMessage` / `VersionHash`） | ✅ 迁 | 快照内容身份的唯一形状 |
| `Snapshot`（3 字段）+ `Clone()` | `snapshot.go:37-41`、`:43-52` | ✅ 迁 | 同上；`Clone()` 是纯函数（深拷贝，无 I/O） |
| `DomainSnapshotHash` | `snapshot.go:20`（`"tellurion.risk_snapshot_hash.v1"`） | ✅ 迁，**逐字节不动** | 改它 = 全历史 `risk_rules_hash` 作废 |
| hash 纯函数 | `snapshot.go:57-104`（`SnapshotHash`）、`:108-117`（`verifyHash`）、`hash.go:52-55`（`writeField`）、`model.go:150-174`（`canonicalParams`） | ✅ 迁，建议做成 `Snapshot.Hash()` | 与 `engine/provenance.go:78-87` 的 `StrategyVersionInput.Hash()` 同形 |
| 4 个哨兵 | `internal/rule/errors.go:18`（`ErrInvalidRuleType`）、`:20`（`ErrInvalidScope`）、`:22`（`ErrInvalidParams`）、`:28`（`ErrSnapshotMismatch`） | ✅ 迁，**文案不改** | 跨仓消费方要能 `errors.Is` |
| `MaxParamsBytes = 64 << 10` | `model.go:64` | ✅ 迁（见 §1.3 的撞名风险） | 体积上限是形状的一部分 |
| `Bundle` / `ForScope` / `ParseBundle` | `snapshot.go:121-137`、`:142-164` | ❌ **不迁** | `Bundle` 有未导出索引字段（`rules` / `byScope`）；`ParseBundle` 是 decode + 默认值填充，属服务策略 |
| 其余 7 个哨兵 | `errors.go:8/10/12/14/16/24/26`（`ErrRuleNotFound` / `ErrRevisionNotFound` / `ErrRevisionConflict` / `ErrRevisionNotDraft` / `ErrIdempotencyConflict` / `ErrSnapshotNotFound` / `ErrNoApprovedRules`） | ❌ **不迁** | 服务私有，跨仓只经 wire code 消费 |

> 与第三轮裁决 ② 逐条对齐（`tellurion:docs/design/qt-decoupling/cross-repo-decision-record-2026-10-25.md:372`）：
> 迁 `Snapshot` / `SnapshotRule` / `Scope` / `RuleType` + 冻结词表 + hash 纯函数 + 4 哨兵；`Bundle` / `ParseBundle` /
> `ForScope` 不迁；编码逐字节搬运、禁止换 Canonical；带固定向量 + 对拍腿；前置 = QT-74；发 MINOR。
> 本表把它落成了可核对的清单，**无一处扩大范围**。

### 1.2 编码必须逐字节搬运（最容易做错的一条）

两个编码器**产出不同字节**，换成 contracts 的 `NewCanonical` 会让全部历史 hash 作废：

| | risk 现有编码 | contracts `engine/canonical.go` |
|---|---|---|
| 帧结构 | 文本长度前缀：`fmt.Fprintf(b, "%d:%s=%d:%s\n", len(key), key, len(value), value)`（`risk:internal/rule/hash.go:52-55`） | 二进制 uvarint 长度前缀（`engine/canonical.go:83-91`）+ 值类型标签 `"s;"` / `"n;"`（`:101-134`） |
| 域标签 | 裸写一行 `DomainSnapshotHash` + `\n`（`snapshot.go:71-72`） | `NewCanonical` 首帧即域标签，**无换行**（`:70-74`） |
| 排序 | 按 `rule_id` 升序（`snapshot.go:63`），重复即拒（`:64-68`） | 只对集合语义字段排序（`SymbolSet`，`:167-179`） |
| 参数规范化 | `json.Decoder.UseNumber()` 解码 → `json.Marshal`（map 键自动排序、无多余空白）（`model.go:150-174`） | `Value` 只认受支持类型，**显式拒绝 `json.Number`**（`:185-197`） |

⇒ 结论：**照抄 `snapshot.go:57-104` 的字节流，不要「顺便换成更整洁的编码器」**。这不是风格问题：
`risk_rules_hash` 是回测结果的可复现锚点，落在 `quant.run_provenance`（`tellurion:pkg/storage/run_provenance.go:46`）
与编排的 `run.run_artifacts`（`tellurion-orchestration:internal/queue/ddl.go:70`）两处，换编码 = 两处锚点全部对不上。

### 1.3 三个撞名坑（迁进本仓后才会暴露）

| 撞名 | 本仓现状 | risk 侧 | 后果 / 处置 |
|---|---|---|---|
| `MaxParamsBytes` | `engine/abi.go:429` = **4096**（wasm `initialize` 的 params 区上限） | `internal/rule/model.go:64` = **65536**（规则参数 JSON 上限） | 不同包不冲突，但**同名不同值不同义**。迁入后 `docs/API.md` 必须把两者的语义分列，**禁止**为了「统一」把 65536 改成 4096（会拒掉今天的合法快照） |
| `ScopeSymbol` | `engine/abi.go:389`（`StrategyScope` 的取值，`"symbol"`，指「每实例绑定一个 symbol」） | `internal/rule/model.go:35`（`Scope` 的取值，`"symbol"`，指「规则适用范围」） | 不同包不冲突，但同一个消费方同时 import 两包时**极易混用**。`docs/API.md` 要给两者各写一句「不是同一个概念」，并禁止用类型别名互转 |
| `Scope` | 本仓无裸 `Scope`（只有 `StrategyScope`） | `internal/rule/model.go:29` 裸 `Scope` | 新包内叫 `risk.Scope` 可接受；**不要**并入 `engine/`（会与 `StrategyScope` 并排，加剧上一条） |

### 1.4 与既有判据的冲突（必须由主任拍，本文不代裁）

本仓与 risk 仓各有一条**成文判据**，指向相反的结论：

- **risk 仓的判据**：`tellurion-risk:docs/ARCHITECTURE.md:85` —「契约（`RuleSnapshot` / `RuleEvent` / `BreakerCommand`）
  要等 X3b/X3c 出现**第二个 Go 使用方**后才按 ADR-0001 ① 迁 contracts」。
  实测今天**没有**第二个 Go 使用方：risk `go.mod` 只有 `pgx` + `kit v0.1.4`（**零 contracts 依赖**）；
  core 全仓 `pkg/risk` 对 `Bundle` / `ParseBundle` / `Snapshot` **零命中**（只命中一个与本议题无关的
  `TestChecker_RulesReturnsSnapshot`，`tellurion:pkg/risk/take_profit_test.go:347`）。
- **第三轮裁决 ② 的判据**：`tellurion:docs/design/qt-decoupling/cross-repo-decision-record-2026-10-25.md:372` —
  「前置 = QT-74 先落」，**没有**要求「先有第二个 Go 使用方」。

⇒ 两条判据在今天会给出不同答案。**建议**：以第三轮裁决 ② 为准（它是更晚的、针对本单的裁决），
但**同批把 risk 仓 `ARCHITECTURE.md:85` 那句话订正**为「按 ② 裁决，QT-74 后即迁；第二使用方（core F3）随迁移后接线」——
否则两个仓的文档会长期互相打脸。这条订正归 risk 仓的单，**不在本组**。

> 另需如实登记：risk 的 F3（`bundle` → `ParseBundle` → `RiskPolicy`）**今天仍未接**，core 零调用点
> （`tellurion-risk:docs/ARCHITECTURE.md:108`）。所以 ① 迁完后的**真实消费方仍然为零**——
> 这是「先冻结形状、后续线」的有意取舍，代价是这段时间契约里多一个没人 import 的包（与 `dataprotocol`
> 当时的处境同型，见本仓 `docs/API.md` §8 #4）。**接受这个代价才立单**。

### 1.5 本仓随单必改的两处（否则护栏假绿）

1. `tools/check_contract_shape/main.go:36` 的 `contractRoots = []string{"domain", "engine", "fees", "dataprotocol"}` —
   **不加 `risk` 就是假护栏**：不在作用域里的包，「没命中」不等于「干净」（工具头部注释自己写着这句，
   `main.go:34-36`）。同批要补夹具（现有四条根各有一条红腿）。
2. `docs/API.md` §4 登记区 + `docs/api-surface.txt` 基线（`-update` 重生成）+ `README.md` 兼容矩阵加 `v0.7` 行 +
   `CHANGELOG.md` 一节 —— 三条是 CI 硬门（`check_api_surface` / `check_surface_freeze` / `check_release_matrix`）。

## 2. ② 编排形状（claim / `run_artifacts` / 子 run 入队）

### 2.1 只迁「消费方需要的三个 wire 形状」，不迁队列内部记录

orchestration 的 `RunRecord`（17 字段，`internal/queue/model.go:105-124`）是**队列内部记录**（含
`Attempts` / `MaxAttempts` / `WorkerID` / `LeaseUntil` / `CancelRequested`），把它整份搬进契约等于把租约与重试
状态机搬进契约仓 —— 违反「只放被 ≥2 个仓使用的**跨边界**类型」。要迁的是**消费方真正按字节对齐**的三份：

| 形状 | 生产侧（orchestration） | 消费侧（已有第二份独立实现） | 是否满足「≥2 仓」 |
|---|---|---|---|
| **claim** | `internal/httpapi/workers.go:103-111` 的 `claimView`（7 字段：`run_id` / `run_spec` / `attempts` / `max_attempts` / `lease_until` / `cancel_requested` / `parent_job_id`） | **core 已逐字段复刻**：`tellurion:internal/orchestrationworker/client.go:92-101` 的 `ClaimedJob` | ✅ **是**（两份手写副本，今天靠人肉对齐） |
| **`run_artifacts` 读面** | `internal/httpapi/provenance.go:32-48` 的 `provenanceView`（15 字段） | **core 有两份**：`tellurion:docs/openapi.yaml:3802-3823` 的 `BacktestProvenance`（openapi 里**明文写着**「mirror tellurion-orchestration's `provenanceView` exactly … must not drift」）+ `tellurion:pkg/storage/run_provenance.go:33-47` 的 `RunProvenanceRecord` | ✅ **是**（**三份**形状，且「不许漂」是文档承诺、无机械护栏） |
| **子 run 入队** | `internal/parents/service.go:130-132` 的 `ChildIdempotencyKey`（`"parent:" + parentID + ":" + strconv.Itoa(childIndex)`，注释标 **stable；改规则 = MAJOR**）+ `internal/httpapi/parents.go:23-48` 的提交请求形状 | analytics 已手写复刻：`tellurion-analytics:internal/upstream/sensitivity.go:38-47`（`idempotency_key` / `kind` / `batch{base_run_spec, grid}`）+`:67-68` 打 `POST /api/runs/parents` | ✅ **是** |

**最强的一条证据**是 `provenanceView`：「must not drift」不是评论，是 core openapi 里的**成文承诺**
（`tellurion:docs/openapi.yaml:3806-3807`），而今天没有任何机械手段能判它是否兑现。
这正是契约仓该收的形状——**把一句口头承诺换成一份可 import 的唯一实现**。

**不建议迁的**（留各自仓）：`jobView` / `parentView` / `parentChildView` / `parentCountsView`
（编排自己的读面投影，只有它一个消费者）、`RunRecord` / `ParentRecord` / `WorkerRecord`（队列与租约状态机）、
`ListFilter`（实现细节）。

### 2.2 「不迁」的部分要写清，否则会被下一次误当成漏项

- `ArtifactStore` 接口（`internal/queue/artifacts.go:62-65`）不迁：它带**写语义 + 身份漂移拒绝**，
  属编排的服务策略；契约只收 DTO。
- `run_artifacts` 的 **DDL** 不迁（本仓不持有表，`docs/ARCHITECTURE.md` §7「不持有任何表」）。
- `ChildIdempotencyKey` 是**函数**不是 DTO：按 ADR-0004，若它只是「给定 parent_id 与 index 拼串」的纯函数，
  可以随值对象走；但它决定**幂等身份**，建议**只迁常量/形状、把拼串留在编排**——
  或整段迁入并在 CHANGELOG 标注「改规则 = MAJOR」。**二选一，立单时定死**。

### 2.3 与已有文档的冲突（同样是「先裁再迁」）

orchestration 自己的对表**早就写着**契约位 = contracts：

- `tellurion-orchestration:docs/ARCHITECTURE.md:100`「**A-1 目标 `JobQueue`（contracts，先落 orchestration 内部的同名接口，S1 迁）**」
- `:174`「`RunRecord` … 契约位 = contracts」
- `:185`「`ProvenanceResponse` … 契约位 = 本仓 HTTP」

⇒ orchestration 侧的口径是**混的**：port 与 `RunRecord` 说迁 contracts，`ProvenanceResponse` 说留本仓 HTTP。
而 §2.1 的实测结论恰好相反——**`ProvenanceResponse` 才是三份形状在漂的那个，`RunRecord` 反而只有一份**。
这个矛盾必须由主任在立单时裁定，本文建议取实测结论（迁 `provenanceView`，不迁 `RunRecord`），
并同批订正 orchestration `ARCHITECTURE.md` §9.3 的「契约位」两行（归 orchestration 仓的单）。

## 3. ③ broker DTO（order intent / fill / account）

### 3.1 今天一个都没上 wire

| DTO | 定义处 | 字段数 | 今天被 JSON 序列化吗 |
|---|---|---|---|
| `OrderRequest`（= order intent 的承载物，**broker 侧没有叫 `OrderIntent` 的类型**） | `tellurion-broker:internal/gateway/model.go:80-89` | 7 | ❌ |
| `Fill` | `model.go:137-149` | 11 | ❌ |
| `AccountSnapshot` | `model.go:178-186` | 6 | ❌ |
| `Position` / `Cash` / `Fees` | `model.go:153-162` / `:164-168` / `:170-174` | 7 / 2 / 3 | ❌ |

全仓 `json.Marshal` 只有 3 处，没有一处在 HTTP 路径上（`internal/ackstore/pg.go:96` / `:53` 只读写
`broker.order_acks.ack` 这个 JSONB 列；`model.go:96` 只对字段值数组算指纹）。
`tellurion-broker:docs/api-schema.md:9` 自己写：「**契约地位：Go port（进程内），订单面今天不是 HTTP**」。

⇒ **建议：③ 排最后，前置 = `SVC-BROKER-G1` 落定**（`tellurion-broker:docs/TASKS.md:37` ⛔）。
理由：现在迁 = 把一份「尚未存在的报文」冻进契约，之后 HTTP 面一落地就会发现字段不对，再改就是 MAJOR。
这与 ① 的顺序逻辑同型（词典先于形状），只是这里的前置是「报文先于形状」。

#### 3.1.1 死锁裁定：`CT-05` ③ ↔ `SVC-BROKER-G1` 谁先（2026-10-26 执行侧代裁，待主任追认）

**循环原文**（两侧都指对方）：本材料 §0 行 ③ 与 §3.1 写「③ 的前置 = G1 先把订单 HTTP 面定形」；
`tellurion-broker:docs/TASKS.md:37`（标题）与 `:46`（前置复核）写「G1 ⛔ = 等 A 组 broker DTO」。
两条并读 ⇒ 谁都不许开工，且不会自动解开。材料 §4.3 原本给的破环建议是「③ 只出对表、不落包」——
**那不算产出**（对表 §3.3 已经存在，再抄一遍不产生新形状），所以必须裁顺序而不是绕开。

**裁定：G1 先行，③（`CT-05`）收编在后。** 三条理由：

1. **契约仓自己的准入判据反对先落包**：`tellurion-contracts:AGENTS.md` §3 禁「只有单个仓使用的类型」。
   今天 broker DTO 的真实消费方只有 broker 一侧（core 没有订单 HTTP 客户端；`core:pkg/live` 是被退役的旧实现）。
   在 G1 之前把它冻进契约 = 把单仓形状塞进底座，违反的是 contracts 自己的规则，不是偏好。
2. **HTTP 面是形状唯一的可执行验证场**：G1 落地后，报文被真 `json.Marshal` 过一次、有测试钉住字段名，
   ③ 才是「收编已被验证的形状」；反过来的话 contracts 先出的那份包从未经过上 wire 的检验，
   而 §3.3 那三条差异（`int` vs `float64`、`Side` 两值 vs `Direction` 四值、`OrderType` 两值 vs 四值）恰恰只有真报文才能证伪。
3. **反向顺序会把 G1 变成「等发版」**：③ 落包要走 contracts 的 MINOR + 兼容矩阵 + `check_surface_freeze` + 三仓同测，
   G1 因此要等一次跨仓发版；而 ① ② 也在争同一批发版批次（`v0.7.0` / `v0.8.0`），批次耦合只会更重。

**代价（如实）**：G1 期间订单 DTO 暂住 broker（`internal/gateway/model.go:80-186`，今天已经在那儿），
存在一个「HTTP 报文先定形、契约后收编」的漂移窗口；若 core 在 ③ 之前就要消费，会出现一份手写镜像——
这正是 ③ 想消掉的东西，窗口期内它是**已知且被接受**的复刻，不是新缺陷。

**两条机械判据锁住这个窗口**（缺一即 G1 不算落定，③ 不许开工）：
① `tellurion-broker:docs/api-schema.md` 必须写出订单 / 成交 / 账户三族的**逐字段报文表**（字段名 + 可空性 + 类型），
   并有 httpapi 的用例断言序列化后的键集合；② ③ 收编时的验收 = broker 删掉本地形状改引契约 +
   同一条**报文快照逐字节一致**（改前 / 改后各序列化一次比对），不允许「先并存两份」。

**剩余前置（不是循环）**：③ 落包仍要 §4 的 wire DTO 口径先拍（建议 A：区分 core 的 openapi 生成物与服务间手写形状）。
那是一条**单点裁决**，拍完 ③ 就能开工；G1 不再等它，因此循环已断。

### 3.2 order intent 到底住在谁那儿（三个候选名，选错就是造第三份）

| 候选 | 位置 | 字段 | 归属判定 |
|---|---|---|---|
| `risk.OrderIntent` | `tellurion:pkg/risk/interfaces.go:29-38`（8 字段，`OrderID` / `RunID` / `Symbol` / `Side domain.Direction` / `OrderType domain.OrderType` / `Qty float64` / `Price float64` / `Ts`） | 风控裁决的输入 | **core，不迁**：`pkg/live/interfaces.go:34` 只是 core 内的类型别名；`tellurion-risk:docs/ARCHITECTURE.md:124`（F2）已裁「port 留 core；X3b 冻结的是 **wire JSON**，不是 Go 类型」 |
| `broker.OrderRequest` | `tellurion-broker:internal/gateway/model.go:80-89`（7 字段，`int` 数量、`Side` 两值、`OrderType` 两值） | 通道适配的请求 | **建议的迁移对象**（但等 §3.1 前置） |
| `domain.Order` | `tellurion-contracts:domain/types.go:154-189` | 策略产出的订单 | **已在契约里**，不是本议题 |

⇒ **`order intent` 这个名字在三个仓指三件不同的事**（`tellurion-broker:docs/db-schema.md:59` 还把
「订单意图（`OrderIntent`）」的 owner 标成 **core**）。立单时必须在文档里把「谁的 intent」写死，
否则迁进契约的那份会被两个消费方按各自理解使用。

### 3.3 三个必须写进契约文档的形状差异（今天的真陷阱）

| 概念 | contracts / core 侧 | broker 侧 | 差异 |
|---|---|---|---|
| 数量 | `domain.Position.Quantity` = **`float64`**（`domain/types.go:90`）、`domain.Order.Quantity` = `float64`（`:159`） | `Position.Quantity` / `OrderRequest.Quantity` / `Fill.Quantity` = **`int`**（`model.go:87` / `:143` / `:157`） | 类型不同、语义也不同（broker 是整股）⇒ **禁止**在契约里用类型别名把两者等同 |
| 方向 | `domain.Direction` 四值：`long` / `short` / `close` / `hold`（`domain/types.go:22-26`） | `Side` **两值**：`buy` / `sell`（`model.go:21-24`） | 不是同一词表；映射必须显式（`close` → `sell` 是**业务判断**，不能藏在 JSON 字段名里） |
| 委托类型 | `domain.OrderType` 四值：`market` / `limit` / `stop` / `trailing`（`domain/types.go:38-42`） | `OrderType` **两值**：`market` / `limit`（`model.go:32-35`） | broker 刻意不收 `stop` / `trailing`（`model.go:76-79` 注释：「券商私有字段留在适配器内部，执行层语义留在 core」）⇒ 契约里**只能**收 broker 那两值，不能为了「统一」把四值搬过来 |

`broker:docs/ARCHITECTURE.md:156-164` 其实已经把这三条写成了对表（「`Quantity` 是柜台口径」「`Side(buy/sell)`」
「不含 `BusinessType` / `ExecAlgo` / `StopPrice` / `Trail*`」）—— 本节的贡献是把**与 contracts 现有 `domain` 类型
的冲突**点名，因为迁入后它们在同一个仓里。

### 3.4 kit 机械上被排除（这条让「落点」不用再讨论）

kit 的 `tools/check_scope` 把 `order` / `market` / `portfolio` 等词写进了 `businessWords`
（`tellurion-kit:tools/check_scope/main.go:45-60`，词表出处 = `docs/SCOPE.md:34-49`），
且切词会命中驼峰（`OrderRequest` → `[order request]`，`main.go:431-442`）。
⇒ **任何 broker DTO 放进 kit 都会当场红**。kit 的准入判据第二条还要求「被 ≥2 个部署单元使用」
（`docs/SCOPE.md:18-27`）。所以 ③ 的落点只有两个：contracts 或「各服务自己的 openapi 生成物」。

## 4. 阻塞性裁决：wire DTO 到底能不能手写进 contracts

这是 ② ③ 的共同前置，也是本组唯一需要主任拍的新问题。

### 4.1 冲突的原文

- `tellurion-contracts:docs/API.md:766`：「**`dto/`**：跨进程 DTO 是 `tellurion:docs/openapi.yaml` 的**生成物**
  （`oapi-codegen`），真相在 core，本仓只收生成物且**禁手改**（CI 重生成 diff 必须为空）」
- 而 ② ③ 的 DTO 生产方**不是 core**（是 orchestration / broker），且**六个服务仓全都没有 openapi 文件**
  （实测 glob：`*openapi*` / `*swagger*` 只命中 core 的 `docs/openapi.yaml`、`cmd/analysis/handlers_openapi.go`、
  agent 的一个测试、ui 的一条 ADR）。

### 4.2 两个选项与代价

| 选项 | 做法 | 代价 | 适用 |
|---|---|---|---|
| **A. 区分两类跨进程形状，允许手写** | 一条 ADR 修订 `docs/API.md` §8 的口径：把「跨进程 DTO」拆成 **(a) core 的 HTTP 报文**（仍走 openapi 生成物，禁手改）与 **(b) 服务间 Go 共享形状**（可手写进 contracts，配 `conformance/cases.json` 做跨仓对拍） | 要维护两份口径的边界；手写的形状**没有**「CI 重生成 diff 为空」这道自动护栏，只能靠 conformance 向量 | **建议 ② 走这条**（它是目前唯一能让「must not drift」变成机械判据的路径） |
| **B. 服务先出 openapi，本仓收生成物** | orchestration / broker 各建 `docs/openapi.yaml` + `oapi-codegen` 管线，contracts 只收生成物 | 6 个服务仓今天 0 条 openapi 管线；每个仓要新增一个代码生成步骤 + 一个 CI 校验步骤；orchestration 的 wire 形状现在写在 `ARCHITECTURE.md` §9.3 与 `api-schema.md` 里，要先迁成 openapi 才能生成 | 只有在「所有服务都必须有版本化 API 文档」被裁为硬要求时才值得 |

### 4.3 建议（附代价）

- **① 走 A 但几乎不触发争议**：它的主体是**值对象 + 纯函数 + 哨兵**，只有 `Snapshot` 兼作 `/bundle` 的响应体。
  按 ADR-0004 的「值对象上的纯函数行为随类型走」，hash 函数进契约是既有判据的直接适用。
- **② 走 A，且只迁 §2.1 那三份**。这是本组**收益最高**的一单：它把一句 openapi 里的口头承诺
  （`must not drift`）换成唯一实现，同时消掉 core / analytics 的手写副本。
- **③ 先不做**，等 `SVC-BROKER-G1`。若主任希望 ③ 与 ② 同批裁口径，则 ③ 只出**对表**（把
  `broker:docs/ARCHITECTURE.md:156-164` 的形状搬成 §3.3 的差异表），不落包。
- **不建议**「三包同批发一个 MINOR」：`api-surface.txt` 基线一次加三包，`check_surface_freeze` 的 diff
  会失去可读性；且三单的前置完全不同（① 等 QT-74，② 等口径裁决，③ 等 G1）。

### 4.4 四条组级禁区（对三单同时生效，逐条给机械判据）

| 禁区 | 为什么 | 机械判据（可复跑） |
|---|---|---|
| **不得同批发三个形状包** | 三单前置互不相同；同批会让护栏的 diff 失去可归因性 | `git tag` 每发一包只多一个 tag；`git show --stat` 每笔只含**一个**包的目录 |
| **不得改既有 canonical 域标签** | `engine/canonical.go:44-66` 的 8 个 `Domain*` 是历史 hash 的域分隔；改任一个 = 对应家族全历史 hash 作废 | `git diff engine/canonical.go` 在三单里**必须为空**；`go test ./engine -run 'TestGoldenVector|Test.*Hash_ChangesOnEveryDecidingInput'` 保持绿 |
| **不得引入 `contracts → core` 反向依赖** | 一行反向 import 就让本仓无法单独构建（R2 判据） | ① **编译期已由 Go 强制**：core 的 module path 是单段 `tellurion`（`tellurion:go.mod:1`），Go 要求 module path 首段含点，故**任何仓都无法 `require` 它**（`tellurion:AGENTS.md:33` 原文「任何仓都无法 `require` 它」）；② 静态复核 `go list -deps ./... \| grep -c '^tellurion$'` = 0；③ `go build ./...` 在**只 checkout 本仓**时仍绿 |
| **不得把 `Bundle` / `ParseBundle` 或服务私有哨兵顺手带进来** | 它们带服务策略（decode、默认值填充、wire code 映射），进契约就把服务实现冻进底座 | `docs/API.md` §4 登记区里不得出现对应符号；`go run ./tools/check_api_surface` 的符号数与单据声明的增量一致 |

> **跨仓腿（发版后由协调者统一通知）**：本组三单的产物会被 **E / H / I 三个组的 G1 / G3** 消费
> （risk 的 `SVC-RISK-G3`、orchestration 的 `SVC-ORCH-G1/G3`、broker 的 `SVC-BROKER-G1/G3`）。
> 发版本身**不在本组**——本组只出材料；tag、兼容矩阵、pin 升版由协调者在对应单立后统一安排。

## 5. 三张候选六字段单（可直接派发，按序一次一包）

> 编号建议 `CT-03` / `CT-04` / `CT-05`（`CT-01` / `CT-02` 已被占用，实测）。**主任确认后**再把
> 「编号建议」改成正式编号并移入 `docs/TASKS.md` 的「一、本仓任务」。
> 三张单的 ⑥ 都显式排除彼此：**一次只发一包**。

### CT-03（建议）· `contracts/risk` 包：快照形状 + 冻结词表 + hash 纯函数 + 4 哨兵 [P0] ⬜

- **① 目标**：把 risk 的规则快照**形状与 hash 算法**变成跨仓唯一实现，使任何消费方 `import` 一次就能算出与 risk 逐字节相同的 `risk_rules_hash`。
- **② 上下文**：目标形状与算法全文在 `tellurion-risk:internal/rule/snapshot.go:15-117`、`hash.go:15-55`、`model.go:12-64`、`:150-174`、`errors.go:18/20/22/28`；冻结判据 = `tellurion:docs/design/qt-decoupling/cross-repo-decision-record-2026-10-25.md:372`（第三轮裁决 ②）；同形参照 `tellurion-contracts:engine/provenance.go:78-87`（`StrategyVersionInput.Hash()`）；本仓护栏四处 = `tools/check_contract_shape/main.go:36`、`docs/API.md` §4、`docs/api-surface.txt`、`README.md` 兼容矩阵。
- **③ 要求**：① 新建 `risk/`（包名 `risk`），按本文 §1.1 的清单迁入，**`Bundle` / `ParseBundle` / `ForScope` 与其余 7 哨兵不迁**；② 编码**逐字节搬运**（`writeField` 的文本长度前缀 + 按 `rule_id` 排序 + `canonicalParams` 的 `UseNumber`→`Marshal`），**禁止**改用 `NewCanonical`；③ 固定向量测试（输入 → 期望 hex 硬编码）；④ 与 risk 的**对拍腿**（同一组 `SnapshotRule` 两侧各算一次，断言 hex 相同）；⑤ `contractRoots` 加 `risk` 并补夹具；⑥ `docs/API.md` §4 / `api-surface.txt` / 兼容矩阵 `v0.7` / `CHANGELOG` 四件套同步；⑦ 发 MINOR `v0.7.0`。
- **④ 约束与非目标**：**不改** hash 算法与 `DomainSnapshotHash` 字面量（改 = 全历史 hash 作废，需 ADR 级裁决）；**不改** 4 个哨兵的文案；**不把 `risk` 并入 `engine`**；**不在契约里读时钟 / 连库 / 读环境变量**；**不动** risk 的 Store / HTTP / `Bundle`；**不代做** core 的 F3 消费接线；**不同批**发 ② ③ 的形状包。
- **⑤ 验收**：① `go build ./... && go vet ./... && go test ./... -count=1` 三绿（贴真实退出码）；② `go run ./tools/check_api_surface` / `check_surface_freeze` / `check_contract_shape` / `check_release_matrix` 四绿；③ 对拍腿绿（两侧 hex 相同）；④ **破坏验证**：把 `SnapshotHash` 里的 `sort.Slice`（`risk:snapshot.go:63`）去掉 ⇒ 固定向量与对拍腿必须红；把 `writeField` 的格式串改成明文分隔（`%s=%s\n`）⇒ 对拍腿必须红；两处**逐字节还原**后复绿（`sha256` 与改坏前一致）。
- **⑥ 边界**：只做「契约落形状 + 对拍 + 登记 + 发版」。core 的 `bundle → RiskPolicy` 接线、risk 侧改别名、risk `ARCHITECTURE.md:85` 的判据订正，都在各自仓的单。

### CT-04（建议）· 编排三份 wire 形状（claim / `run_artifacts` / 子 run 入队）[P0] ⬜

- **① 目标**：让「领取到的 run」「留痕读面」「子 run 幂等键」三份**已在多仓手写复刻**的 wire 形状，收成契约里的唯一实现，把 core openapi 里那句 `must not drift` 变成机械判据。
- **② 上下文**：生产侧 `tellurion-orchestration:internal/httpapi/workers.go:103-111`（`claimView`）、`internal/httpapi/provenance.go:32-48`（`provenanceView`）、`internal/parents/service.go:130-132`（`ChildIdempotencyKey`）、`internal/httpapi/parents.go:23-48`（提交请求形状）；消费侧副本 `tellurion:internal/orchestrationworker/client.go:92-101`（`ClaimedJob`）、`tellurion:docs/openapi.yaml:3802-3823`（`BacktestProvenance`，含 `must not drift` 原文 `:3806-3807`）、`tellurion:pkg/storage/run_provenance.go:33-47`（`RunProvenanceRecord`）、`tellurion-analytics:internal/upstream/sensitivity.go:38-47`（提交报文）；口径前置 = 本文 §4。
- **③ 要求**：① 新建包（建议 `run/`，包名待定，**不要**叫 `orchestration`——仓名不是形状名），只收**三份消费方形状**；② **不迁** `RunRecord` / `ParentRecord` / `WorkerRecord` / `jobView` / `parentView` / `ArtifactStore`；③ 每份形状配一条**跨仓对拍腿**（本仓一侧 + 生产仓一侧各算/各序列化一次，断言字节相同），形态照抄 `engine/provenance_vector_test.go`；④ `ChildIdempotencyKey` 的二选一（迁函数 vs 只迁形状）在单内定死并写进 `docs/api-schema.md`；⑤ 四件套登记 + MINOR（`v0.8.0`）。
- **④ 约束与非目标**：**不改**任何 wire 字段名 / 可空性 / 顺序（`provenanceView` 的 `*string` 出 `null` 而非省略键，`provenance.go:38-42` 是冻结语义）；**不迁**队列状态机与租约语义；**不在契约里**做 decode / 默认值填充 / 幂等判定；**不动** orchestration 的 HTTP / PG / DDL；**不代做** orchestration `ARCHITECTURE.md:100/174/185` 三处「契约位」口径的订正；**不同批**发 ① ③。
- **⑤ 验收**：① 三条对拍腿全绿；② `go build/vet/test ./... -count=1` + 四条护栏绿（贴真实退出码）；③ **破坏验证**：把 `claimView` 的 `parent_job_id` 的 `omitempty` 去掉 ⇒ 对拍腿必须红；把 `ChildIdempotencyKey` 的分隔符从 `:` 改成 `-` ⇒ 对拍腿必须红；逐字节还原后复绿。
- **⑥ 边界**：只做「契约落形状 + 对拍 + 登记」。orchestration 侧改别名、core 的 `orchestrationworker` 改用契约类型、analytics 的提交报文改用契约类型，都是各自仓的单。

### CT-05（建议）· broker DTO（order intent / fill / account）——**排在 `SVC-BROKER-G1` 之后落包**（顺序裁定见 §3.1.1）[P1] ⛔（等 G1）

- **① 目标**：把 broker 的订单 / 成交 / 账户记账形状在**订单 HTTP 面定形之后**收进契约，避免两个消费方各写一份镜像。
- **② 上下文**：DTO 定义 `tellurion-broker:internal/gateway/model.go:80-186`；形状对表已在 `tellurion-broker:docs/ARCHITECTURE.md:152-167`；契约位今天标 contracts（同表「契约位」列）；HTTP 面未落 `tellurion-broker:docs/TASKS.md:37`（`SVC-BROKER-G1` **⛔**，且它原文「等 A 组 broker DTO」与本单互为前置——**已由 §3.1.1 裁定 G1 先行、本单反向依赖解除**）；与 `domain` 的三处差异见本文 §3.3。
- **③ 要求**：① **前置 = `SVC-BROKER-G1` 落定**（订单面 HTTP 契约冻结 + §3.1.1 那两条机械判据到手）之后再迁；本单**不再等自己**：G1 的 DTO 形状 owner 暂留 broker，③ 是**收编**而不是先立法；② 迁入时按 §3.3 的三条差异写进 `docs/api-schema.md`（数量 `int` vs `float64`、`Side` 两值 vs `Direction` 四值、`OrderType` 两值 vs 四值），**禁止**用类型别名把两侧等同；③ 12 个哨兵**是否随迁**单内另判（`internal/gateway/errors.go:12-36`，它们与 12 个 wire code 一一对应，属「跨仓消费方要能识别」的候选，比 DTO 更接近 ① 的判据）。
- **④ 约束与非目标**：**不迁** XTP 适配器内部类型（`BusinessType` / `XTPConfig` / `XTPDeps`）——`broker:docs/db-schema.md:69` 与 `ARCHITECTURE.md:165` 都写死了「券商特有字段不许泄漏进契约」；**不在** HTTP 面定形前冻结报文；**不动** `broker.order_acks` 的 DDL 与幂等语义（保留期归 ㉒R10 / ⑨ 那两张单）。
- **⑤ 验收**：① 对表逐行可查（DTO × 字段 × 契约位 × 与 `domain` 的差异）；② `SVC-BROKER-G1` 已标 ✅ 且有可指的端点清单，**并满足 §3.1.1 的两条窗口判据**（`docs/api-schema.md` 有逐字段报文表 + httpapi 用例钉住键集合；③ 收编后 broker 本地形状已删、报文快照逐字节一致）；③ 若落包：对拍腿 + 四条护栏绿 + 破坏验证（改一个字段名 ⇒ 对拍腿红 → 还原 → 绿）。
- **⑥ 边界**：只做 broker 形状。真实 XTP（X6c）与 core `pkg/live` 的退役（`SVC-CORE-BROKER-RETIRE`）都不在本单。

## 6. 破坏验证（前提守卫；本轮实跑）

派工写「不适用（无实现可破坏）」。本文仍对**派工自己给的判据**做了一次改坏-还原，因为它可证伪：

```
# 基线：无 risk/
(Get-ChildItem -Name tellurion-contracts) | Select-String -Pattern '^risk$'
#   → 空（「risk 形状已进契约」被证伪）

# 改坏：临时塞一个 risk/ 目录
mkdir tellurion-contracts/risk ; 'package risk' > tellurion-contracts/risk/doc.go
(Get-ChildItem -Name tellurion-contracts) | Select-String -Pattern '^risk$'
#   → risk（判据翻转 ⇒ 证明这条判据真的在判东西，不是恒空）

# 还原
Remove-Item -Recurse -Force tellurion-contracts/risk
git -C tellurion-contracts status --porcelain
#   → 空（逐字节还原；工作树与 HEAD 一致）
```

真跑结果见 §7。**未做**：kit 侧没有同型判据可破坏（它的断言是「tag = v0.1.4」，改 tag 属不可逆动作，本轮不做、也不该做）。

## 7. 本轮实测输出（可复跑）

见文末「施工记录」。此处只声明口径：**本文不含任何 `.go` 改动**，`git show --stat` 应只出现 `docs/`。

## 8. 未做与如实登记

- **三单全部未立、未开工**（`docs/TASKS.md` 里没有对应的正式条目）——本文只给候选六字段。
- **① 的真实消费方今天仍是 0**（core 的 F3 未接，`tellurion-risk:docs/ARCHITECTURE.md:108`）。迁完后会有一段
  「契约里有包、没人 import」的窗口，**接受这个代价才立单**。
- **② 的对拍腿没有现成模板可直接抄**：`engine/provenance_vector_test.go` 是**同仓**向量，
  跨仓对拍要么复制算法到测试里（会漂），要么读对方仓的源码（违反单仓可自建）。**建议**用
  `conformance/cases.json` 的形态扩展一份新语料（输入 + 期望 hex 都由**契约仓实现**产出，生产仓侧断言相同），
  这条在 CT-04 里要写死，否则对拍腿会退化成「两边都从自己实现取期望值」的假绿。
- **③ 的 `order intent` 三义**（core `risk.OrderIntent` / broker `OrderRequest` / contracts `domain.Order`）
  本文只点名、不裁定；立单时必须显式选一个。
- **kit 侧零改动**：三个形状都不进 kit（§3.4 有机械证据），本组在 kit 只留一条跨仓指针。