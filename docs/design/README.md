---
status: active
last-verified: 2026-10-26
verified-by: A 组三个形状包决策材料落地（只出材料，不含实现）
---

# 设计 / 裁决材料

> 本目录放**跨仓落点的裁决材料**：只列证据与影响面，**不代表决策已执行**。
> 拍板后请改对应的契约 / 台账，再把这里的状态改成「已裁决 → 落点见 …」。
> 跨仓引用写成 `repo:path` 纯文本，不做相对链接。

| 材料 | 回答什么 | 状态 |
|---|---|---|
| [qt22-version-hash-decision-material.md](qt22-version-hash-decision-material.md) | `version_hash` canonical 集合：冻结 7 项 vs 实测 6 项；`params_defaults` 与 v2 域标签的代价 | ✅ 已裁 6 项（2026-10-24，随 `v0.5.0` 落地） |
| [qt40-x1a-dataprotocol-decision-material.md](qt40-x1a-dataprotocol-decision-material.md) | `UnifiedDataPoint` 三份形状；`dataprotocol` 先迁 (a) vs 维持 core-only (b) 的影响面与迁入判据 | ✅ 已裁 (a)（2026-10-24，随 `v0.6.0` 落地） |
| [service-dto-migration-decision-material.md](service-dto-migration-decision-material.md) | QT-45 / QT-43 / QT-42 三组服务 DTO 的落点、第二消费方、迁入后谁改 import | 三组均「今天不迁」，等各自触发 |
| [shape-packages-decision-material.md](shape-packages-decision-material.md) | A 组三个形状包（`contracts/risk` / 编排三份 wire 形状 / broker DTO）的顺序、前置、禁区与三张候选六字段单；并逮到「wire DTO = openapi 生成物、禁手改」与「服务仓零 openapi」的结构性冲突 | ⬜ 待裁（三单均未立；① 等 risk `QT-74`，② 等 §4 口径，③ 等 `SVC-BROKER-G1`） |