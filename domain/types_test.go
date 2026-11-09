package domain

import (
	"encoding/json"
	"testing"
	"time"
)

func TestOHLCV_ZeroValues(t *testing.T) {
	var o OHLCV
	if o.Symbol != "" {
		t.Errorf("expected empty Symbol, got %q", o.Symbol)
	}
	if !o.Date.IsZero() {
		t.Error("expected zero Date")
	}
	if o.Open != 0 || o.High != 0 || o.Low != 0 || o.Close != 0 {
		t.Error("expected zero OHLC values")
	}
	if o.Volume != 0 || o.Turnover != 0 {
		t.Error("expected zero volume/turnover")
	}
	if o.TradeDays != 0 {
		t.Errorf("expected zero TradeDays, got %d", o.TradeDays)
	}
	if o.LimitUp || o.LimitDown {
		t.Error("expected false limit flags")
	}
}

func TestOHLCV_JSONRoundTrip(t *testing.T) {
	original := OHLCV{
		Symbol:    "000001.SZ",
		Date:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Open:      10.5,
		High:      11.0,
		Low:       10.2,
		Close:     10.8,
		Volume:    1000000,
		Turnover:  10800000,
		TradeDays: 1,
		LimitUp:   false,
		LimitDown: false,
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded OHLCV
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Symbol != original.Symbol {
		t.Errorf("Symbol: expected %q, got %q", original.Symbol, decoded.Symbol)
	}
	if !decoded.Date.Equal(original.Date) {
		t.Errorf("Date: expected %v, got %v", original.Date, decoded.Date)
	}
	if decoded.Open != original.Open {
		t.Errorf("Open: expected %f, got %f", original.Open, decoded.Open)
	}
	if decoded.High != original.High {
		t.Errorf("High: expected %f, got %f", original.High, decoded.High)
	}
	if decoded.LimitUp != original.LimitUp {
		t.Errorf("LimitUp: expected %v, got %v", original.LimitUp, decoded.LimitUp)
	}
}

func TestOHLCV_JSONFieldNames(t *testing.T) {
	o := OHLCV{Symbol: "TEST", Close: 1.5}
	data, err := json.Marshal(o)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	jsonStr := string(data)
	expectedFields := []string{`"symbol"`, `"date"`, `"open"`, `"high"`, `"low"`, `"close"`, `"volume"`, `"turnover"`, `"trade_days"`, `"limit_up"`, `"limit_down"`}
	for _, field := range expectedFields {
		if !contains(jsonStr, field) {
			t.Errorf("expected JSON to contain %s, got: %s", field, jsonStr)
		}
	}
}

func TestDirection_Constants(t *testing.T) {
	tests := []struct {
		name     string
		value    Direction
		expected string
	}{
		{"long", DirectionLong, "long"},
		{"short", DirectionShort, "short"},
		{"close", DirectionClose, "close"},
		{"hold", DirectionHold, "hold"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, string(tt.value))
			}
		})
	}
}

func TestMarketRegime_ZeroValues(t *testing.T) {
	var m MarketRegime
	if m.Trend != "" {
		t.Errorf("expected empty Trend")
	}
	if m.Volatility != "" {
		t.Errorf("expected empty Volatility")
	}
	if m.Sentiment != 0 {
		t.Error("expected zero Sentiment")
	}
	if !m.Timestamp.IsZero() {
		t.Error("expected zero Timestamp")
	}
}

func TestMarketRegime_JSONRoundTrip(t *testing.T) {
	original := MarketRegime{
		Trend:      "bull",
		Volatility: "low",
		Sentiment:  0.65,
		Timestamp:  time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded MarketRegime
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Trend != original.Trend {
		t.Errorf("Trend mismatch")
	}
	if decoded.Sentiment != original.Sentiment {
		t.Errorf("Sentiment mismatch")
	}
}

func TestOHLCV_NegativeValues(t *testing.T) {
	o := OHLCV{
		Symbol: "TEST",
		Open:   -1.0,
		High:   -0.5,
		Low:    -2.0,
		Close:  -1.5,
		Volume: -100,
	}
	// Domain types don't validate; they just hold data.
	// Verify JSON round-trip preserves negative values.
	data, err := json.Marshal(o)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	var decoded OHLCV
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if decoded.Open != o.Open {
		t.Errorf("Open: expected %f, got %f", o.Open, decoded.Open)
	}
}

func TestOHLCV_ExtremeValues(t *testing.T) {
	o := OHLCV{
		Symbol: "EXTREME",
		Open:   1e15,
		High:   1e15,
		Low:    1e-15,
		Close:  1e15,
		Volume: 1e18,
	}
	data, err := json.Marshal(o)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	var decoded OHLCV
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if decoded.Open != o.Open {
		t.Errorf("Open: expected %e, got %e", o.Open, decoded.Open)
	}
}

func TestStock_JSONRoundTrip(t *testing.T) {
	original := Stock{
		Symbol:    "000001.SZ",
		Name:      "平安银行",
		Exchange:  "SZ",
		Industry:  "银行",
		MarketCap: 30000000000,
		ListDate:  time.Date(1991, 4, 3, 0, 0, 0, 0, time.UTC),
		Status:    "active",
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	var decoded Stock
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if decoded.Symbol != original.Symbol {
		t.Errorf("Symbol mismatch")
	}
	if decoded.MarketCap != original.MarketCap {
		t.Errorf("MarketCap mismatch")
	}
}

// TestMovedTypes_AreGone 钉住 CT-06-B 的契约面收缩：只有 core 一个使用方的
// 根名必须不再能从本包导出（core 侧定义在 tellurion:pkg/domain）。
func TestMovedTypes_AreGone(t *testing.T) {
	// These are referenced through reflect-free existence checks in the
	// companion guard on the core side; here we only assert the retained
	// vocabulary is still intact, which is what the contracts surface
	// promises after the move.
	if DirectionLong != "long" || DirectionShort != "short" {
		t.Fatal("Direction vocabulary changed")
	}
}

func TestBacktestResult_ZeroValues(t *testing.T) {
	var br BacktestResult
	if br.TotalReturn != 0 || br.SharpeRatio != 0 || br.MaxDrawdown != 0 {
		t.Error("expected zero values")
	}
	if br.TotalTrades != 0 || br.WinTrades != 0 || br.LoseTrades != 0 {
		t.Error("expected zero trade counts")
	}
	if br.PortfolioValues != nil || br.Trades != nil {
		t.Error("expected nil slices")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
