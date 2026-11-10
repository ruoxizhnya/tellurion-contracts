// Package domain contains cross-repo domain types and interfaces for
// the quant trading system.
//
// S7-P3-4 (ODR-043 D3): market-data types (OHLCV, Stock, Fundamental,
// etc.) live in this repo's `domain/market`. This package used to
// re-export them as type aliases for backward compatibility; **those
// six aliases were deleted by CT-06-ALIAS-EXIT (2026-11-11; unreleased,
// ships in the next minor v0.10.0)**. Import `domain/market` directly.
//
// CT-06-B (2026-11-10): 只有单一使用方（core）的根名已迁回
// `tellurion:pkg/domain`；本包只保留 ≥2 使用方的形状（C4）。迁出的
// 具体名单见 docs/API.md §6.1。
package domain

import "time"

// Direction represents trading direction
type Direction string

const (
	DirectionLong  Direction = "long"
	DirectionShort Direction = "short"
	DirectionClose Direction = "close"
	DirectionHold  Direction = "hold"
)

// MarketRegime represents detected market regime
type MarketRegime struct {
	Trend      string    `json:"trend"`      // "bull", "bear", "sideways"
	Volatility string    `json:"volatility"` // "low", "medium", "high"
	Sentiment  float64   `json:"sentiment"`  // -1.0 to 1.0
	Timestamp  time.Time `json:"timestamp"`
}
