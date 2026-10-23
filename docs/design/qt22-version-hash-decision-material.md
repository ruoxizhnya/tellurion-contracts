---
status: active
last-verified: 2026-10-24
verified-by: G0D 裁决材料（只出材料，不含实现；G1 第二批发版前由协调者拍板）
---

# QT-22 裁决材料 —— `version_hash` canonical 集合：冻结 7 项 vs 实测 6 项

> **本文只出材料，不是决策**。由 G1 协调者拍板后才进代码；本仓今日**没有** `StrategyVersionInput` / `DomainStrategyVersionHash`（全仓 grep 零命中，见 §3）。
> 证据口径：`repo:path:line` 为 2026-10-24 实测（`HEAD` 见各仓 `git log -1`）。跨仓引用写成纯文本，不做相对链接。

## 0. 结论先行（含一处对派工前提的订正）

1. **派工写的「实测 6 项（含 `capability`、无 `params_defaults`）」不成立**。`tellurion-strategy` 的实测输入是
   `strategy_id + carrier + abi_version + scope + entry + source` 六项，**不含 `capability`**，也**没有 `params_defaults`**
   （`tellurion-strategy:internal/asset/hash.go:34-46`；`capability` 被明确排除，理由写在 `:37-38`）。
   `capability` 出现在**定义表与 resolve 回包**里（`tellurion-strategy:internal/asset/model.go:94-103`、`:235`；
   `internal/asset/model.go:228-239` 的 `ResolveResult`），但**从未进入 version hash**。
   派工大概把「定义行的身份字段」与「hash 输入」混成了同一张表。
2. **与冻结 7 项相比，真实缺口只有一项**：`params_defaults`。`tellurion-strategy` 全仓
   `grep -rn params_defaults --exclude-dir=.git` **0 命中**——不是「实现漏了」，是这个字段在资产服务里**从未存在**。
3. **建议**：G1 第一批只落**与实测一致的 6 项**（`v1`），把 `params_defaults` 登记为**需要 `v2` 域标签的显式增量**；
   在资产服务定义长出 `params_defaults` 之前，不把空默认值塞进 canonical（理由见 §2）。
4. **`v2` 的代价不是「重算一下」**：hash 是四张策略表的**主键 / 唯一键 / 外键**，还被 core 的
   `quant.run_provenance` 与 `run_spec_hash` 折入；改标签 = 历史行必须双标签共存或不可逆重写（§2.3 逐条）。

---

## 1. 两套形状逐项对表

| 项 | 冻结 7 项（`tellurion:docs/design/qt-decoupling/w1-repo-ownership.md:101`） | 实测 6 项（`tellurion-strategy:internal/asset/hash.go:39-46`） | 判定 | 判据 |
|---|---|---|---|---|
| `strategy_id` | ✅ | ✅ | **保留** | 内容寻址的第一维；改 id = 新资产 |
| `carrier` | ✅ | ✅ | **保留** | `wasm` / `so` / `builtin` 决定装载路径，换载体语义必变 |
| `abi_version` | ✅ | ✅ | **保留** | host ABI 是执行语义的一部分（`on_bar/1` vs `/2`） |
| `scope` | ✅ | ✅ | **保留** | 同源码 market / symbol 行为不同；策略侧已把它写进校验（`hash.go:59-61`） |
| `entry` | ✅ | ✅ | **保留** | 入口函数名换 = 不同字节语义 |
| `source` | ✅ | ✅ | **保留** | 源码本体；`ArtifactHash` 只覆盖编译产物，源码版本身份靠它 |
| `params_defaults` | ✅ | ❌（全仓 0 命中） | **新增（但必须走 `v2`）** | 冻结设计说「默认值改变执行结果，不进 hash 会让缓存命中错版本」（`tellurion:docs/design/qt-decoupling/p0-remediation-design.md:115`、`:147`）。但今天没有生产者：`VersionInput` 不收、`Definition` 没有该字段、`POST /api/strategies` 的 body 也没有（`tellurion-strategy:docs/api-schema.md:27`、`:126`）。**零生产者的必填输入 = 所有调用方传空对象，字节流仍会变** ⇒ 不是 v1 的兼容增量 |
| `capability`（冻结里没有） | — | ❌（`hash.go:37-38` 明确排除） | **删除（保持排除）** | 它只参与 MarketRouter 路由（`tellurion-strategy:internal/asset/model.go:71-84` 的注释），不改变源码语义。把它放进 hash 会让「给同一份源码补一个能力声明」变成新版本，产物字节没变却要重新 promote——错误方向 |
| `toolchain` | ❌（冻结明确不含） | ❌ | **不加** | 工具链差异由 `artifact_hash` 区分（`tellurion-contracts:engine/provenance.go:37-46` 的 N7 裁决） |

> **另一个易混项**：`active_version_hash` / `active_artifact_hash` 是**指针**（`tellurion-strategy:internal/asset/model.go:94-103`），
> 不进 hash——它们指向被 hash 的对象，进 hash 会自指。

## 2. `v2` 域标签的代价（逐处）

### 2.1 先看「同标签、两套字段」会怎样

`tellurion-strategy` 今天用 `tellurion.strategy_version_hash.v1` 算 6 项（`tellurion-strategy:internal/asset/hash.go:29-32`、`:72-78`）。
如果 contracts 以**同一个 `v1` 标签**实现 7 项（多一个 `params_defaults`），同一份输入会在两个仓得到**两个 hash**，
而字符串标签还相同——这正是「两份真源」的最坏形态：`errors.Is` 那类判定失效，缓存与台账各自命中一半。

⇒ 只要加项，就必须换标签。**没有第三种选择**。

### 2.2 今天没有 contracts 侧实现，所以「污染」尚未发生

`tellurion-contracts` 全仓 grep `StrategyVersionInput` / `DomainStrategyVersionHash` **零命中**（`engine/canonical.go` 只有
`run_spec_hash` / `dataset_snapshot_hash` / `universe_hash` / `calendar_hash` / `adjustment_anchor_digests` / `params_hash` 六个域）。
所以现在**新增 `v2` 是干净起点**：没有下游按 contracts 算过 v2，向量常量可以从零钉。

### 2.3 但策略侧的历史行是真实成本

`version_hash` 在资产服务里同时是**四种角色**：

| 落点 | 角色 | 证据 |
|---|---|---|
| `strategy.strategy_versions` | **主键**的一部分 | `tellurion-strategy:docs/db-schema.md:18`、`:25` |
| `strategy.strategy_artifacts` | `UNIQUE (strategy_id, version_hash, toolchain)` + 外键引用版本行 | `tellurion-strategy:docs/db-schema.md:19`、`:26`、`:30` |
| `strategy.strategy_definitions` | `active_version_hash`（生效指针，带 `CHECK` 与 artifact 指针同空同非空） | `tellurion-strategy:docs/db-schema.md:17`、`:70` |
| `strategy.strategy_build_jobs` | 外键 `(strategy_id, version_hash)` | `tellurion-strategy:docs/db-schema.md:21`、`:30` |
| `quant.run_provenance.version_hash` | 运行留痕；同一 `run_id` 写入不同值会 fail-loud | `tellurion:pkg/storage/run_provenance.go:38`、`:61`、`:100-108` |
| `run_spec_hash` 的分量 | `RunSpecInput.StrategyVersionHash` 折进 run 身份 | `tellurion-contracts:engine/provenance.go:137-146`（字段）、`:188-196`（折进 canonical） |

改标签的概率后果：

- **可重算的行**：`strategy_versions.source` 已存（`tellurion-strategy:docs/db-schema.md:18`），若 `params_defaults` 也能在当时的状态下唯一还原，可离线重算并批量改写。
- **不可重算的行**：默认值是**可变配置**。历史时刻的 defaults 若没留快照，今天无论怎么算都是「猜」。⇒ 旧行只能**保留 v1 原值**。
- **连带失效**：`run_spec_hash` 含 version_hash；同一逻辑 run 在新标签下会得到不同 `run_spec_hash`，历史上的「同参复现」判定立刻失去可比性（`run_spec_hash` 的用途就是这条判定）。
- **主键重写**：`strategy_artifacts` / `build_jobs` 的外键与唯一约束都要跟着改；`active_version_hash` 的 `CHECK` 与 `active_artifact_hash` 成对约束（`db-schema.md:70`）也在同一事务里。

**成本结论**：如果资产服务已有真实行，`v2` 的成本是**一次双标签过渡**（旧行保持 v1、新行写 v2、读侧按标签分派），
不是一次 SQL `UPDATE`。若表还是空的（今天代码已落、生产库未灌数据），成本接近零——**拍板前必须先查目标库行数**，
这个数本文不能替它编。

### 2.4 `params_defaults` 的第二种解法（不建议）

把 defaults 放进**另一个 hash**（例如复用 `DomainParamsHash` 或新开 `params_defaults_hash`），再与 version_hash 组合进
`strategy_ref`。**禁止**：`strategy_ref` 只有 `id@hash` 一种形状（`tellurion-contracts:engine/provenance.go:264-275`），
再拼一段就是第三种引用形状；而 G0D 派工也明令不得用 `DomainParamsHash` 代替新域标签。

## 3. 给 G1 的落地建议（材料，不是已执行）

1. **第一批（可立即做）**：在 `engine/provenance.go` 增 `StrategyVersionInput` 六字段 + `Hash()`，
   域标签 `DomainStrategyVersionHash = "tellurion.strategy_version_hash.v1"`；向量测试钉固定常量。
   与 `tellurion-strategy:internal/asset/hash.go:52-79` **逐字段、逐顺序一致**，先让两个仓算出同一个值（对拍腿）。
2. **`params_defaults` 单独立项**：等资产服务定义表真的长出该字段（今天没有：`tellurion-strategy:internal/asset/model.go:92-103`），
   以 `DomainStrategyVersionHashV2` 发布；过渡期双标签，读侧按标签分派。**不要**在 v1 里塞空对象。
3. **对拍腿**：契约仓与策略仓各跑同一组输入（含 6 个必变 + 1 个不必变的反例），比较 hex 相等。
   策略侧已有 6 个必变用例（`tellurion-strategy:internal/asset/asset_test.go:188-207`），**没有固定向量**（只有确定性 + 敏感性），
   所以今天把向量常量钉进契约仓不会与任何旧断言冲突。
4. **冻结面词表**：`w1-repo-ownership.md:101` 是「冻结 7 项」的正文。若最终选 6 项 + 延后，必须在 §2.4 留一行
   「`params_defaults` 延后到 v2，理由与触发条件」——不能让正文继续写一个不存在的字段。

## 4. 分发与边界

- **只出材料**：本组不落 `StrategyVersionInput`、不改 `tellurion-strategy`、不打 tag。
- **上报对象**：按派工要求报 G1 / G4 / G5 / G6 / G7 协调者（仓与组的对应以派工表为准）；受影响仓：
  `tellurion-contracts`（新类型 + 域标签）、`tellurion-strategy`（别名切换 + 双标签读侧）、`tellurion`（`run_spec_hash` 消费与对拍）。
- **禁区复核**：未复用 `DomainParamsHash`；未把 `capability` 塞进 hash；未提前落任何 `StrategyVersionInput` 代码。
