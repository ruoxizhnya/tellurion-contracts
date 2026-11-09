// Package domain contains cross-repo domain types and interfaces for
// the quant trading system.
//
// S7-P3-4 (ODR-043 D3): Market-data types (OHLCV, Stock, Fundamental,
// etc.) now live in this repo's `domain/market` (2026-10-22 自 core 的
// `tellurion:pkg/domain/market` 逐字迁入). This package re-exports them as
// type aliases for backward compatibility. New code SHOULD import
// `domain/market` directly.
//
// CT-06-B (2026-11-10): 只有单一使用方（core）的根名已迁回
// `tellurion:pkg/domain`；本包只保留 ≥2 使用方的形状（C4）。迁出的
// 具体名单见 docs/API.md §6.1。
package domain

import (
	"time"

	"github.com/ruoxizhnya/tellurion-contracts/domain/market"
)

// Direction represents trading direction
type Direction string

const (
	DirectionLong  Direction = "long"
	DirectionShort Direction = "short"
	DirectionClose Direction = "close"
	DirectionHold  Direction = "hold"
)

// S7-P3-4 (ODR-043 D3): Market-data types migrated to `domain/market`（本仓；2026-10-22 随 QT-05 ③-① 迁入）.
// These aliases preserve backward compatibility: domain.OHLCV and
// market.OHLCV are the SAME type (not two different types), so all
// existing consumer code compiles unchanged. New code SHOULD prefer
// the `market.` package path.
//
// To migrate a file: change `domain.OHLCV` → `market.OHLCV` and the
// import. No type conversion is needed.
//
// CT-06-B (2026-11-10): `IndexConstituent` / `FundamentalData` 两个别名只有
// core 一个使用方，已删；需要时直接 import `domain/market` 用同名类型。
type (
	OHLCV       = market.OHLCV
	Stock       = market.Stock
	IndexDaily  = market.IndexDaily
	Split       = market.Split
	Dividend    = market.Dividend
	Fundamental = market.Fundamental
)

// MarketRegime represents detected market regime
type MarketRegime struct {
	Trend      string    `json:"trend"`      // "bull", "bear", "sideways"
	Volatility string    `json:"volatility"` // "low", "medium", "high"
	Sentiment  float64   `json:"sentiment"`  // -1.0 to 1.0
	Timestamp  time.Time `json:"timestamp"`
}
