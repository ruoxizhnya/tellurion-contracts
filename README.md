# tellurion-contracts

> module: `github.com/ruoxizhnya/tellurion-contracts` ｜ 状态：**S1 落代码（QT-05）／S2 打 tag（QT-07）——当前为占位仓，尚无实现**

跨仓唯一缝隙：**port 接口 · DTO/值对象（含 domain）· 四个引擎契约 · 哨兵错误 · abi 常量**。所有跨边界的东西只能住这里。

## 目录

| 路径 | 内容 |
|---|---|
| `docs/` | 本仓文档：台账 / 边界 / 决策 / 归档 —— 从这里进：`docs/README.md` |
| `docs/API.md` | 导出面清单（符号 × 稳定性 × 错误语义）；登记区由 CI 核对 |
| `tools/check_api_surface/` | 导出面 AST 扫描护栏（`go run ./tools/check_api_surface`） |
| `AGENTS.md` | 本仓入口（构建、边界、禁止、依赖版本） |
| `doc.go` | 占位包（让空仓 `go build ./...` 从一开始就是绿的） |

## 快速开始

```bash
go build ./...
go vet ./...
go test ./...
go run ./tools/check_api_surface   # 导出面 ↔ docs/API.md 登记区一一对应
```

本仓**单独 checkout 就能跑通**上面四条。需要别的仓源码才过 → 边界错了，不是加 `replace` 能解决的问题。

## 边界

**可以**：被 core / agent / 各服务 `require`；只放**被 ≥2 个仓使用**的跨边界类型与接口

**不可以**：任何实现、业务规则、DB 访问；`Config` / `Options` / `Impl` 命名（契约泄漏实现细节的形状）；只有单个仓使用的类型（那条属于泄漏，应留在使用方）

## 文档

- 文档导航：`docs/README.md`
- 本仓台账：`docs/TASKS.md`（唯一；跨仓任务只留指针）
- 边界与依赖方向：`docs/ARCHITECTURE.md`
- 决策记录：`docs/adr/`

## 许可

MIT（见 `LICENSE`）。
