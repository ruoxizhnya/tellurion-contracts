# 0002 · `*Config` 族在边界上的三条处置

- **状态**：已裁决（2026-10-22，QT-05 ③-① 落 `domain/` 时实测；修订 [0001](0001-契约只增不破且零实现.md) ③ 的适用范围，不废除它）
- **场景**：0001 ③ 禁 `Config` / `Options` / `Impl` 命名，理由是这类名字会把实现细节写进契约。但 S1 是**整包迁** `pkg/domain` 与 `pkg/backtest/contracts`，一刀切会同时撞上两堵墙：① 迁进来一堆只有 core 用的配置形状；② `engine/` 的 `TradingConfig` / `PriceLimitConfig` 已经写进 `tellurion:docs/SPEC.md` 的 RunSpec 冻结文本——**改名就是改冻结面**。禁令与冻结面在同一个类型上撞车，必须给出可执行的判据，而不是两边都当没看见。
- **实测（2026-10-22，core + agent 全仓 grep `domain.<T>`）**：
  - `Config` / `DatabaseConfig` / `RedisConfig` / `ServiceConfig` / `TushareConfig`：**0 引用**（连 core 自己都不用；`Config` 还带 `Password` 字段，描述的是部署 YAML）；
  - `StrategyConfig`：core 4 个文件；`ExecutionConfig`：core 21 个文件；agent 两者皆 0；
  - `TradingConfig` / `PriceLimitConfig`：RunSpec 冻结文本里的伴生符号名。
- **决策（三条，按类型在边界上的角色分，不按名字分）**：
  1. **零引用的部署配置不迁入**：`Config` / `DatabaseConfig` / `RedisConfig` / `ServiceConfig` / `TushareConfig` 留在 core，由 core 侧切换依赖那次 commit 直接删除（0 调用点，无需改写）。理由：它们描述的是 core 进程的部署形态，不是任何跨仓缝隙；`Password` 字段尤其不该进公开库。
  2. **core-only 的两个 `*Config` 迁入但按原名保留，标 experimental**：`StrategyConfig` / `ExecutionConfig` 随整包迁入、进 QT-05 ③-⑨ 回收清单。不就地改名的理由：它们本就要在 S2 后搬走（`StrategyConfig` 最可能归 `tellurion-strategy`，`ExecutionConfig` 大概率回 core——执行留在 core），改名等于给一个即将搬走的类型付一次破坏性变更的账，还要在 core 侧多出 25 个文件的 diff。
  3. **冻结面伴生名保留**：`TradingConfig` / `PriceLimitConfig`（`engine/`，QT-05 ③-① 落地）原名进入本仓。理由：名字本身已出现在 SPEC 的 RunSpec 冻结文本里，改名 = 改冻结文本 + 破坏 core 的 `spec_drift_test.go` —— 冻结面优先于命名规范。
- **判据（以后遇到同名族照这条判）**：先问「这个名字是否已经出现在冻结文本 / 跨语言契约里」——是 ⇒ 保留；否 ⇒ 再看「它是不是 core-only 的搭车类型」——是 ⇒ 不迁（①）或迁入后进 ③-⑨ 清单（②）；只有「要长期留在本仓、且名字泄漏实现细节」时才改名，且那是一次破坏性变更。
- **代价 / 后果**：本仓在 S2 前短暂持有 2 个 core-only 的 `*Config`（experimental 标注 + 清单在册，回收有据）；5 个部署配置随 core 侧 `pkg/domain` 整目录删除一起消失（**不单独删**，与 sed 同一次 commit 做，避免「删了没人知道为什么」）。0001 ③ 因此从「绝对禁令」变成「默认规则 + 上述三条例外判据」，例外必须在本条登记才成立。
- **替代关系**：无（修订 0001 ③，不替代）
