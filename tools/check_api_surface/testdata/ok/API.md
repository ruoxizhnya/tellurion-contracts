# fixture — 登记表与导出面完全一致

前言里出现的 `Ghost` 不参与核对（标记区之外）。

<!-- surface:begin -->

### `engine/` — 包 `contracts`

| 符号 | 类别 | 稳定性 | 错误语义 / 判据 |
|---|---|---|---|
| `EngineRunner` | interface | frozen | 见 §5.2 |
| `EngineRunner.RunBacktest` | method | frozen | `errors.Is` |
| `EngineRunner.Close` | method | frozen | — |
| `Status` | type | stable | — |
| `StatusCompleted` | const | frozen | — |
| `ErrFault` | var | stable | — |
| `New` | func | stable | — |
| `Runner` | type | stable | — |
| `Runner.Run` | method | stable | — |

### `domain/` — 包 `domain`

| 符号 | 类别 | 稳定性 | 错误语义 / 判据 |
|---|---|---|---|
| `OHLCV` | type | stable | — |

<!-- surface:end -->

## 本节以下的表格都不算登记

| 符号 | 说明 |
|---|---|
| `NotRegistered` | 在核对区之外，检查器必须看不见它 |
