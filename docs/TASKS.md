---
status: active
last-verified: 2026-11-10
verified-by: CT-06 收口（2026-11-10：CT-06-A/B 已执行；余下 6 个 market 别名按「别名不是重复真相」豁免留仓，**代裁·待主任追认**；跨 11 仓实测引用：`OHLCV` core+agent、`Fundamental` core+data，`Stock`/`IndexDaily`/`Split`/`Dividend` 仅 core；本仓 `go build/vet/test ./... -count=1` 与 `check_api_surface` / `check_surface_freeze` / `check_contract_shape` / `check_release_matrix` 全 0）；CT-05 收编完成（2026-11-09：`tellurion-broker` 删本地 DTO 改引 `contracts/broker` 类型别名 + 13 组报文改前/改后逐字节 golden；broker build/vet/test 全 0）+ CT-04 编排 wire 契约腿（2026-11-07：新建 `run/` 包 —— `ClaimView` 七字段 / `ProvenanceView` 十五字段 / 父作业提交请求族 + `ChildIdempotencyKey` 纯函数；三份形状逐字节对拍编排仓 + `conformance/run_cases.json` 七用例；`contractRoots` 五根 → 六根；导出面 418 → 427、基线 408 → 417；本地 tag `v0.8.0` **未 push**，三仓改类型与 pin **未执行**；§4 口径代裁为 A、第六轮 #9 解除；施工记录见「一、本仓任务」）+ CT-03 契约腿（2026-11-07：新建 `risk/` 包 —— 快照形状 + 冻结词表 + `Snapshot.Hash()`/`Verify()`/`Clone()` + 4 哨兵；编码逐字节搬运、固定向量 `14ecab86…` + 跨仓对拍腿；`contractRoots` 四根 → 五根；导出面 465 → 487、基线 448 → 470；本地 tag `v0.7.0` **未 push**，三仓 pin **未执行**；施工记录见「一、本仓任务」）+ A 组三个形状包决策材料（2026-10-26：只出材料、0 代码、三单未立；六条派工前提逐条复现 + 破坏验证实跑；施工记录见「一、本仓任务」）+ QT-40 X1a 契约腿（2026-10-24：新增 `dataprotocol/` 包 —— `UnifiedDataPoint` 七字段（含 `VisibleAt`）+ `ToJSON`/`DeduplicateKey` + 15 个 `DataType*` 常量；无构造器、不读时钟；`v0.6.0` 已打 tag **且已 push**（2026-10-26 实测 `git ls-remote --tags origin` 有 `refs/tags/v0.6.0` → peel `5ee7fae`；此前记的「待 push」是 2026-10-24/25 的旧账。⚠️ 但本仓 `main` 领先 `origin/main` **1 笔**（`3055b6e` 未推）——tag 在远端 ≠ 分支已同步）；施工记录见「一、本仓任务」）+ G4 执行组（2026-10-24：core / agent / data / strategy 四仓升 pin `v0.5.0`；`check_contract_pin` GREEN 6 pin、conformance 2/2 零跳过；施工记录见「一、本仓任务」）+ G2 执行组（2026-10-24：QT-22 契约腿 `StrategyVersionInput` 六项 + `DomainStrategyVersionHash`，`v0.5.0` 已发版；施工记录见「一、本仓任务」）+ G0 执行组（2026-10-24：QT-20 契约腿 + QT-41b ④ 错误族，`v0.4.0` 已发；施工记录见「一、本仓任务」）+ G0D 裁决材料（QT-22 / QT-40 X1a / 服务 DTO 三份，只出材料、0 代码）；QT-09 本仓侧交接（2026-10-23：② examples/ + ③ conformance 语料与契约侧自检腿 + ④ PR 模板/check_pr_body + ⑤ check_contract_shape + surface 冻结 check_surface_freeze/api-surface.txt，全部接进 CI；三条工具另补表驱动单测）；QT-05 ③-① 本仓侧收官（四包 429/429 + §8 #1~#3 裁决 + 下游腿交接单）+ 自审订正（ADR 文件名死链、方法类别、EngineRunner 返回值、fees 注释失真；新增 ADR-0004）；CT-02 结案归档（版本规则表 + 兼容矩阵 + CHANGELOG + check_release_matrix 进 CI）；推送回写（2026-10-22：6 个 commit 已进 `origin/main`，该批 HEAD = `fbe9a2d`，`go get` 实测可解析，CI 429/429；§8 #4 登记为 X1a 待办）；**2026-10-25 B6-1**：QT-07 ⑥ S3 结案回写——agent 删掉 `replace tellurion => ../tellurion`（零 import / 零 require / CI 不 checkout core），core 侧同批清空 `check_contract_pin.py` 的白名单与 `check_repo_autonomy.py` 的 R1 死登记；本仓 AGENTS §3 的「唯一过渡例外」改为无豁免禁令、ADR-0001 附「例外收回记录」、QT-05 在册余项清零转 ✅
---

# TASKS — tellurion-contracts 未完成项

> **归档说明（2026-11-11）**：本文件已瘦身为「只装未做的事」。2026-11-10 及以前已收口（✅）的任务原文、施工记录、验收输出与破坏验证，已逐字移入 [`docs/archive/tasks-archive-2026-11.md`](archive/tasks-archive-2026-11.md)；该文件顶部有「检索索引」表（ID / 标题 / 收口日期 / commit / 证据链接）。检索：`grep -n "<任务ID或commit>" docs/archive/tasks-archive-2026-11.md`。
>
> 状态：`⬜ 待做` / `🔶 进行中` / `⛔ 阻塞（有未解除前置）` / `⏸ 挂起（外部依赖或已裁决暂不实施；**不计入开放项**，重开须另立单并满足该条「重开条件」）` / `✅ 完成（附日期，随即移入归档）`

## 一、未做 / 挂起 / 外部依赖

### CT-06-ALIAS-EXIT · 6 个 `market` 别名的删除前置 [P2] ⏸

- **① 目标**：给「6 个别名是否删」一个可执行的前置与退出动作，避免长期悬空。
- **② 上下文**：`OHLCV` / `Stock` / `IndexDaily` / `Split` / `Dividend` / `Fundamental` 六个 `domain` 根名是 `domain/market` 的类型别名，**不是重复真相**；11 仓逐名 grep（2026-11-11）仍有 core / agent / data 的文本层引用，各仓 pin 的是已发布版本（core/agent `v0.9.0`、data `v0.6.0`，无 `replace`）⇒ 本仓本地删除**不会立刻弄红下游构建**，破坏发生在下一次 pin bump。按「删符号 = MAJOR」必须先做下游迁移。裁决与证据见 `CT-06` 收口块（已归档）。
- **③ 要求**：① 下游（core / agent / data）先把引用迁到 `domain/market.*`；② 三仓同时 pin bump 到删除别名后的 MINOR/MAJOR 版本；③ 本仓再删别名并更新 `docs/API.md`、跑四道护栏。
- **④ 约束与非目标**：本仓不先删；不为「凑数」改接口；不做兼容 shim 之外的新机制。
- **⑤ 验收**：11 仓 grep 对六名的引用为 0；`check_api_surface` / `check_surface_freeze` / `check_contract_shape` / `check_release_matrix` 全绿。
- **⑥ 边界**：形状与版本纪律归本仓；下游迁移与 pin 归各仓。
- **重开条件**：core / agent / data 三仓完成 `domain/market.*` 迁移并同批 pin bump。

#### 执行手册（前置就绪后才做；**本仓今天不执行**）

**A. 前置条件（缺一不动，逐条留证）**

1. 下列三仓的 Go 代码对六名的引用为 0（只数 `.go`，排除本仓 `docs/`、`CHANGELOG` 与归档文本）：
   ```bash
   grep -rn --include='*.go' -E 'domain\.(OHLCV|Stock|IndexDaily|Split|Dividend|Fundamental)\b' <repo>
   ```
   - `tellurion`（core）
   - `tellurion-agent`
   - `tellurion-data`
   每仓附 `git rev-parse HEAD` + 上面命令的**空输出**。另外 8 仓同样跑一遍；未引入 contracts 的仓附 `go.mod` 无依赖行即可（`grep tellurion-contracts go.mod` 为空）。
2. 三仓 `go.mod` 的 contracts pin 已升到删除别名后的版本（见 D 的顺序），且 `go build ./...` 绿。
3. 另一条会话已明确回报「三仓已迁完」；本条由用户转达，未收到前保持 ⏸。

**B. 删除清单（本仓，按序执行，一次提交）**

| # | 文件 | 动作 |
|---|---|---|
| 1 | `domain/types.go` | 删 `type ( OHLCV … Fundamental )` 别名块及上方 S7-P3-4 / CT-06-B 注释；保留 `MarketRegime`；改写包注释（不再说「本包 re-export market 类型」） |
| 2 | `domain/market_alias_test.go` | **整文件删除**（它断言 6 个别名与 `market.*` 同一类型，留着必红） |
| 3 | `domain/market/doc.go` | 删「`domain.OHLCV` — type alias」兼容段，只留 canonical 说明 |
| 4 | `docs/API.md` §4 登记区 | 删 6 行：`OHLCV` / `Stock` / `Fundamental` / `IndexDaily` / `Split` / `Dividend`（别名行，非 `market/` 实体行） |
| 5 | `docs/API.md` §6 / §6.1 | 把「留本仓」处置改为「已删除（CT-06-ALIAS-EXIT，v0.10.0），改用 `domain/market.*`」；同步 §6.1 的命中数表 |
| 6 | `docs/api-surface.txt` | `go run ./tools/check_surface_freeze -update` 重生成（预期比基线少 6 行） |
| 7 | `CHANGELOG.md` | 新增 `## [v0.10.0] - <日期>`，条目写明「**破坏性**：删除 `domain.OHLCV` / `Stock` / `IndexDaily` / `Split` / `Dividend` / `Fundamental` 六个根别名 → 改用 `domain/market.*`」 |
| 8 | `README.md` 兼容矩阵 | 新增 `v0.10` 行；版本规则表不用改（0.x 删符号 = bump MINOR 的既有规则） |

> 版本目标：`v0.9.0` → **`v0.10.0`**（0.x 阶段删符号 = MINOR，见 README「0.x 阶段」）。不跳 `v1.0.0`。

**C. 四道护栏（删除后逐条跑，全绿才提交）**

```bash
go run ./tools/check_api_surface       # 期望：385 符号 / 登记 385（今天 391，-6）；无「未登记」
go run ./tools/check_surface_freeze    # 期望：与第 6 步重生成后的基线逐行一致
go run ./tools/check_contract_shape    # 期望：绿（删符号不触发形状/依赖规则）
go run ./tools/check_release_matrix    # 期望：CHANGELOG v0.10.0 ↔ 矩阵 v0.10 ↔ tag 三件套一致
```

- `check_api_surface`：先跑 `-list` 取当场数字，**不要把 385 抄进文档**（会漂）；红 = 登记区与导出面不一致。
- `check_surface_freeze`：必须先 `-update` 再校验；只删基线不删登记（或反之）两边都会红。
- `check_contract_shape`：不因本次删除放宽任何白名单。
- `check_release_matrix`：tag 只有 push 后才被下游 `go get`；本地 tag 不算发版。

**D. 跨仓 pin 升级顺序（谁先升、为什么）**

1. **本仓先发** `v0.10.0` 并 push annotated tag（不 push 的 tag 无法 `go get`）。
2. 消费仓按「消费面从窄到宽」升，每仓独立 PR（`go get …@v0.10.0` → `go build ./...` → `go vet ./...`）：
   - **`tellurion-data` 先升**：只直接使用 `Fundamental`，改动面最小，先暴露问题。
   - **`tellurion-agent` 次之**：命中集中在 `pkg/ai/prompts.go` 的提示词字符串（非编译使用方），改动面小。
   - **`tellurion`（core）最后升**：AST 复核 nontest 262 处、test 533 处，成本最高；它升完才把矩阵的「已切」标注改为 `v0.10.0`。
3. 其余仓（analytics / strategy / risk / broker / orchestration / kit / ui）无 contracts 依赖或只经上面三仓间接消费，**不需要单独动作**；逐个 `grep tellurion-contracts go.mod` 留证即可。
4. 回滚策略：任一仓 `go build` 红 ⇒ **该仓留在 `v0.9.0`**，本仓不回退已发布版本（只前进）。红的原因写进该仓 TASKS，不改本仓导出面。

**E. 非目标（避免做成新债）**

- 不留 `// Deprecated` 兼容 shim，不留第二份 `domain.OHLCV = market.OHLCV`。
- 不在 core 重建「指回本仓」的别名：core 的 `pkg/domain` 是 CT-06-B 已定的自有同形类型，不随本单改。
- 不为「少改下游」把删除拆成多批 —— 6 个是一次性同一语义（同一类型别名），拆开只会让版本矩阵多出无意义的中间态。

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
