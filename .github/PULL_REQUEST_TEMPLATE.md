# 变更说明

<!--
本模板的四个必填段落由 CI 的 `tools/check_pr_body` 核对（仅 pull_request 事件）。
删段落 / 只留标签不填值 / 占位符未替换 = CI 红。确无内容的项写「无」，不要删标题。
契约仓的边界见 README「边界」与 AGENTS.md §3：只放跨仓共享的形状，零实现、零 DB 访问。
-->

## 使用方（必填）

<!-- 哪个使用方（core / agent / 哪个服务仓）需要这条变更？写 repo 名与被使用的符号；不要只写「下游」。 -->

- 使用方：<填写：repo 名>
- 用到的符号：<填写>

## 变更类型与版本（必填）

<!-- MINOR / PATCH / breaking（MAJOR）三选一；写明是否已按 README「版本规则」bump -->

- 变更类型：<填写：MINOR / PATCH / breaking>
- 版本 bump（CHANGELOG / 兼容矩阵）：<填写>

## 登记区与冻结基线（必填）

<!-- 导出面有增删改才需要更新；没有变化也要写明「无导出面变化」并点名两个文件 -->

- `docs/API.md` 登记区：<填写>
- `docs/api-surface.txt` 冻结基线：<填写>
- 核对命令：`go run ./tools/check_api_surface && go run ./tools/check_surface_freeze`

## Config / Options / Impl 与 DB import（必填）

<!-- 契约仓不允许 `*Config` / `*Options` / `*Impl` 命名（ADR-0002 在册例外除外），也不允许任何数据库/驱动 import。确无则写「不涉及」，并给出 grep 结论。 -->

- 是否涉及 `Config` / `Options` / `Impl` 或 DB import：<填写：不涉及>
- 说明 / grep 结论：<填写>
