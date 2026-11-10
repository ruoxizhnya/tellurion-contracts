// Package market contains market-data domain types for the quant
// trading system.
//
// This sub-package was introduced by S7-P3-4 (ODR-043 D3 "soft layering")
// as the canonical home for market-data types. From 2026-10-22 to
// 2026-11-11 the root `domain` package re-exported them as type aliases
// (`domain.OHLCV` = `market.OHLCV`). **CT-06-ALIAS-EXIT (2026-11-11)
// deleted those aliases** (unreleased; ships in the next minor v0.10.0):
// `OHLCV` / `Stock` / `IndexDaily` / `Split` / `Dividend` /
// `Fundamental` are now spelled `market.X` only, and consumers must
// import this package directly.
package market
