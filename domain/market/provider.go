package market

import (
	"context"
	"time"
)

// Provider defines the interface for accessing market data.
//
// Renamed from domain.MarketDataProvider during S7-P3-4 migration.
// The legacy `domain.MarketDataProvider` alias was removed by CT-06
// (2026-11-07) after the rerun found 0 consumers.
type Provider interface {
	GetOHLCV(ctx context.Context, symbol string, start, end time.Time) ([]OHLCV, error)
	GetFundamental(ctx context.Context, symbol string, date time.Time) (*Fundamental, error)
	GetStocks(ctx context.Context, exchange string) ([]Stock, error)
	GetLatestPrice(ctx context.Context, symbol string) (float64, error)
	GetIndexConstituents(ctx context.Context, indexCode string) ([]string, error)
}
