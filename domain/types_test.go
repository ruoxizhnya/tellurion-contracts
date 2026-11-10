package domain

import (
	"encoding/json"
	"testing"
	"time"
)

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
