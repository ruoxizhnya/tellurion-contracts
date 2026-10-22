# tellurion-contracts

> module: `github.com/ruoxizhnya/tellurion-contracts` ｜ 状态：**S1 落代码本仓侧已完成（QT-05 ③-①，2026-10-22：`fees/` + `domain/market/` + `domain/` + `engine/` 四个包，429 个导出符号）／S2 打 tag（QT-07）待做——今天零 tag，使用方只能 pin `main`；下游切依赖那一步另需一次 push**

跨仓唯一缝隙：**port 接口 · DTO/值对象（含 domain）· 四个引擎契约 · 哨兵错误 · abi 常量**。所有跨边界的东西只能住这里。

## 目录

| 路径 | 内容 |
|---|---|
| `fees/` | **已落地的库包**：A 股费率模型的数据形状（常量 / 档位 / 结构体 / 纯校验）。导出面见 `docs/API.md` §4，边界裁决见 §8 #3 |
| `domain/market/` | **已落地的库包**：行情值对象（`OHLCV` / `Stock` / `Fundamental` / `IndexDaily` …）+ 频率 / 时间语义词表（PIT 红线）+ `Provider` 接口。46 个导出符号见 `docs/API.md` §4；`OHLCV`/`Fundamental` 的字段集与 core 的 `docs/SPEC.md` 互锁 |
| `domain/` | **已落地的库包**：信号 / 订单 / 持仓 / 组合 / 公司行为引擎 / 因子名册 / 回测结果 176 个导出符号（`docs/API.md` §4）。含 `domain/market` 的兼容别名；5 个零引用部署配置按 [`adr/0002`](docs/adr/0002-config-族的边界处置.md) 不迁 |
| `engine/` | **已落地的库包**：四契约（`BacktestRequest` / `BacktestResponse` / `SignalWeight` / `EngineRunner`）+ 内核路由 + 成交口径 + canonical/provenance hash + 三类哨兵 + abi 常量，共 194 个导出符号。**目录名 `engine/`、包名仍是 `contracts`**（调用点前缀不变） |
| `docs/` | 本仓文档：台账 / 边界 / 决策 / 归档 —— 从这里进：`docs/README.md` |
| `docs/API.md` | 导出面清单（符号 × 稳定性 × 错误语义）；登记区由 CI 核对 |
| `CHANGELOG.md` | 版本账本：每个发版一节，breaking 逐条写明「旧 → 新」 |
| `tools/check_api_surface/` | 导出面 AST 扫描护栏（`go run ./tools/check_api_surface`） |
| `tools/check_release_matrix/` | 版本纪律护栏：矩阵 ↔ CHANGELOG ↔ tag 三件套不一致即红（`go run ./tools/check_release_matrix`） |
| `AGENTS.md` | 本仓入口（构建、边界、禁止、依赖版本、拉取前提） |
| `doc.go` | 占位包（让空仓 `go build ./...` 从一开始就是绿的） |

## 快速开始

```bash
go build ./...
go vet ./...
go test ./...
go run ./tools/check_api_surface      # 导出面 ↔ docs/API.md 登记区一一对应
go run ./tools/check_release_matrix   # 版本规则表 + 兼容矩阵 ↔ CHANGELOG ↔ tag
```

本仓**单独 checkout 就能跑通**上面五条。需要别的仓源码才过 → 边界错了，不是加 `replace` 能解决的问题。

## 边界

**可以**：被 core / agent / 各服务 `require`；只放**被 ≥2 个仓使用**的跨边界类型与接口

**不可以**：任何实现、业务规则、DB 访问；`Config` / `Options` / `Impl` 命名（契约泄漏实现细节的形状）；只有单个仓使用的类型（那条属于泄漏，应留在使用方）

## 版本规则

规则口径的唯一源头是 `tellurion:docs/TASKS.md` **QT-19 ③**（跨仓决策）；本节是它落在本仓的执行形态，改口径先去那张单，不要在这里另起一套。

<!-- rules:begin -->
| 位 | 何时 bump | 例子（符号名指设计文档，本仓今天还没有符号） |
|---|---|---|
| `MAJOR` | **删符号**、**改签名**、**改语义**（默认值、错误语义、单位与口径都算语义） | 去掉 `RunSpec` 的字段；`CheckOrder` 的返回含义翻转；哨兵错误的判定条件变松 |
| `MINOR` | **新增向后兼容的能力**：新包、新类型、新导出符号、struct 新增字段 | 新增 `domain.Interval`；给 DTO 加一个可选字段 |
| `PATCH` | **文档 / 注释 / 无行为修正** | `docs/API.md` 措辞订正、CHANGELOG 补条目、README 加前提说明 |
<!-- rules:end -->

**0.x 阶段（现在到 1.0 前）**：允许 breaking，但**必须 bump MINOR + 写 CHANGELOG 条目**，条目里写明「破坏性：<旧 → 新>」。
**1.0 之后**才启用「`deprecated` 的旧接口至少保留 1 个 MINOR 周期」——0.x 阶段没有这条义务，但也没有「随便删」的豁免：删一样要留账。

**判定补充（Go 特有的两个坑，容易误判成 MINOR）**：

- 给已发布的 **interface 加方法 = breaking**（1.0 之后是 MAJOR；0.x 期间按上面 bump MINOR，但 CHANGELOG 条目必须写成「破坏性」）：所有实现方立刻编不过。
- 给 struct **加字段**通常是 MINOR，但**只要使用方存在无字段名的位置化字面量**（`T{a, b}`）就是 MAJOR —— 不确定就当 breaking 记。

**不参与 semver 的东西**：内容寻址标识 `strategy_ref` / `artifact_hash` / `run_spec_hash` 自带版本语义（同一个 hash = 同一份内容，变了就是新对象），**不映射到版本号**，也不因它们的变化而 bump。

**tag 规范**：一律 annotated tag（`git tag -a`），message 首行 `vX.Y.Z: <类型>`，类型取 `feat` / `fix` / `breaking`；tag 必须打在 `main` 上并 push —— 不 push 的 tag 无法 `go get`，等于没发版。

## 兼容矩阵

每一行 = 一个**版本序列** `vX.Y`（同序列的补丁版共用一行，口径来自 QT-19 ③-③ 的 `vX.Y`），写清该序列同时被哪些上游版本验证过。
**每次发版更新这张表**：新的 `vX.Y` 加一行，同序列的 `vX.Y.Z` 更新那一行 —— 漏更新 → `go run ./tools/check_release_matrix` 当场红（不是提醒）；矩阵里留一个没有账的序列同样红。
「三个仓的 require 版本是否都落在矩阵声明的范围内」是跨仓核对，归 `tellurion:docs/TASKS.md` QT-09 的 CI 步骤（本仓 checkout 不到别的仓，脚本只核对**仓内**三件套的一致性）。

<!-- matrix:begin -->
| 契约版本 | 状态 | 支持的 core（`tellurion`） | 支持的 agent（`tellurion-agent`） |
|---|---|---|---|
| `main`（未发版） | 开发基线，**不做兼容承诺** | 尚无 require（QT-05 落码后接入） | 尚无 require（同上） |
<!-- matrix:end -->

⚠️ 本仓今天**零符号、零 tag**：首个版本将是 `v0.1.0`，由 `tellurion:docs/TASKS.md` QT-07 ③-S2-a 在 S2 打。在那之前，「矩阵里有没有我的版本」这个问题对任何使用方的答案都是「没有，你在用 `main`，随时可能 breaking」。

## 发版流程

```bash
# 1. 先在 CHANGELOG.md 把本次改动记进对应版本节（标题格式 ## [vX.Y.Z] - YYYY-MM-DD）
# 2. 再更新 README 兼容矩阵：新的 vX.Y 序列加一行，同序列的补丁版改那一行
# 3. 导出面有增删改 → 同步 docs/API.md §4 登记区
go run ./tools/check_api_surface && go run ./tools/check_release_matrix   # 两条都绿才继续
git tag -a v0.1.0 -m "v0.1.0: feat <一句话>"
git push origin main && git push origin v0.1.0
```

使用方升级只需 `go get github.com/ruoxizhnya/tellurion-contracts@vX.Y.Z` + `go mod tidy`；契约变更**不走「顺手改一下」**，走这张表 + CHANGELOG 的账。

## 拉取前提（私有仓）

本仓一旦被设为 **private**，`go get` 默认走公共 proxy，**必然失败**——这是 `tellurion:docs/design/qt-decoupling/repo-inventory.md` §1.1/§5 登记的已知风险（现状 12 仓全为 public，转私有前必须先配好下面两条）：

1. 环境与 CI 设 `GOPRIVATE=github.com/ruoxizhnya/*`（让 module 直连 git，不经 proxy/checksum 服务）；
2. 配 git 凭证（HTTPS PAT 或 SSH key）——`GOPRIVATE` 只让 module 直连 git，**认证是另一件事**：只配前者不配后者，`go get` 会在 `git ls-remote https://github.com/ruoxizhnya/tellurion-contracts` 那一步 403。

## 文档

- 文档导航：`docs/README.md`
- 本仓台账：`docs/TASKS.md`（唯一；跨仓任务只留指针）
- 边界与依赖方向：`docs/ARCHITECTURE.md`
- 导出面清单：`docs/API.md`（CI 核对，不登记即红）
- 版本账本：`CHANGELOG.md`（与兼容矩阵、tag 由 CI 三方对拍）
- 决策记录：`docs/adr/`

## 许可

MIT（见 `LICENSE`）。
