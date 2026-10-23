# examples — 跨仓 port 的最小用法

> 目的：下游**不读本仓实现**（本仓也没有实现）就能看到调用形状。
> 每个示例都是 `package main`，只做编译期演示（桩恒返回零值/空实现，不算业务），
> 因此**不进导出面**（`check_api_surface` / `check_surface_freeze` 都跳过 `package main`）。
> 判据：`go build ./...` 必须把全部示例编过；示例里的 `var _ Iface = stub{}` 是方法集的编译期钉子。

| 示例 | 契约符号 | 对应 `docs/API.md` 词条 |
|---|---|---|
| [`marketprovider/`](marketprovider/main.go) | `market.Provider`（5 方法：`GetOHLCV` / `GetFundamental` / `GetStocks` / `GetLatestPrice` / `GetIndexConstituents`） | §4 `domain/market/` 的 `Provider` 及五行 `Provider.*`；语义在 §5.2 |
| [`enginerunner/`](enginerunner/main.go) | `contracts.EngineRunner`（`RunBacktest(ctx, BacktestRequest) (*BacktestResponse, error)`），入参/出参 `BacktestRequest` / `BacktestResponse` | §4 `engine/` 的 `EngineRunner` / `EngineRunner.RunBacktest` / `BacktestRequest` / `BacktestResponse`；冻结面见 §5.1、§5.2 |

## 新增示例的规矩

1. 一个跨仓 port 一个目录；目录名 = port 的小写连写（`marketprovider`）。
2. 必须是 `package main` + `func main()`，用 `var _ <端口> = stub{}` 钉住方法集。
3. 桩不许做业务计算（不做算术、不读写任何存储），只返回零值或 `panic("example")`。
4. 在本表补一行：示例 ↔ 契约符号 ↔ `docs/API.md` 词条；`go build ./...` 绿才算完。
