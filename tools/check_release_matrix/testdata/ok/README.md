# tellurion-contracts（夹具：一个已发版的最小样本仓）

> 这是 `tools/check_release_matrix` 的夹具，不是本仓文档。真实三件套在仓库根的 README.md / CHANGELOG.md / AGENTS.md。

## 版本规则

<!-- rules:begin -->
| 位 | 何时 bump | 例子 |
|---|---|---|
| `MAJOR` | 删符号、改签名、改语义 | 去掉 DTO 字段 |
| `MINOR` | 新增向后兼容的能力 | 新增一个值对象 |
| `PATCH` | 文档 / 注释 / 无行为修正 | 措辞订正 |
<!-- rules:end -->

## 兼容矩阵

<!-- matrix:begin -->
| 契约版本 | 状态 | 支持的 core | 支持的 agent |
|---|---|---|---|
| `main`（未发版） | 开发基线，不做兼容承诺 | 任意 | 任意 |
| `v0.1` | 已发版 | `v0.1.0` | `v0.1.0` |
<!-- matrix:end -->
