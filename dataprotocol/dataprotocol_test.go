package dataprotocol

import (
	"strings"
	"testing"
	"time"
)

// TestDeduplicateKey_MatchesCoreVector pins the exact key produced by
// `tellurion:pkg/data/source/unified.go` (its test uses the same inputs and
// asserts only the prefix; this one pins the whole string).
func TestDeduplicateKey_MatchesCoreVector(t *testing.T) {
	p := UnifiedDataPoint{
		Symbol:    "AAPL",
		TradeTime: time.Date(2026, 5, 14, 9, 30, 0, 0, time.UTC),
		DataType:  DataTypeGlobalOHLCV,
	}
	const want = "AAPL|global_ohlcv|2026-05-14T09:30:00Z"
	if got := p.DeduplicateKey(); got != want {
		t.Fatalf("DeduplicateKey() = %q, want %q", got, want)
	}
}

func TestDeduplicateKey_KeepsNanosecondPrecision(t *testing.T) {
	p := UnifiedDataPoint{
		Symbol:    "600519.SH",
		TradeTime: time.Date(2026, 5, 14, 9, 30, 0, 123456789, time.UTC),
		DataType:  DataTypeOHLCDaily,
	}
	const want = "600519.SH|ohlcv_daily|2026-05-14T09:30:00.123456789Z"
	if got := p.DeduplicateKey(); got != want {
		t.Fatalf("DeduplicateKey() = %q, want %q", got, want)
	}
}

func TestDeduplicateKey_ExcludesSource(t *testing.T) {
	base := UnifiedDataPoint{
		Symbol:    "600519.SH",
		TradeTime: time.Date(2026, 5, 14, 0, 0, 0, 0, time.UTC),
		DataType:  DataTypeOHLCDaily,
	}
	a, b := base, base
	a.Source = "tushare"
	b.Source = "mootdx"
	if a.DeduplicateKey() != b.DeduplicateKey() {
		t.Fatalf("source must not affect the dedup key: %q vs %q", a.DeduplicateKey(), b.DeduplicateKey())
	}
}

func TestDeduplicateKey_DistinguishesDataType(t *testing.T) {
	base := UnifiedDataPoint{
		Symbol:    "600519.SH",
		TradeTime: time.Date(2026, 5, 14, 0, 0, 0, 0, time.UTC),
	}
	a, b := base, base
	a.DataType = DataTypeOHLCDaily
	b.DataType = DataTypeCapitalFlow
	if a.DeduplicateKey() == b.DeduplicateKey() {
		t.Fatalf("data types write different tables and must not collapse: %q", a.DeduplicateKey())
	}
}

func TestDeduplicateKey_NormalizesToUTC(t *testing.T) {
	instant := time.Date(2026, 5, 14, 9, 30, 0, 0, time.UTC)
	shanghai := instant.In(time.FixedZone("CST", 8*60*60))
	utcPoint := UnifiedDataPoint{Symbol: "AAPL", TradeTime: instant, DataType: DataTypeGlobalOHLCV}
	localPoint := UnifiedDataPoint{Symbol: "AAPL", TradeTime: shanghai, DataType: DataTypeGlobalOHLCV}
	if localPoint.DeduplicateKey() != utcPoint.DeduplicateKey() {
		t.Fatalf("same instant in different zones must share a key: %q vs %q",
			localPoint.DeduplicateKey(), utcPoint.DeduplicateKey())
	}
}

func TestToJSON_EncodesData(t *testing.T) {
	p := UnifiedDataPoint{Data: map[string]interface{}{"close": 12.5}}
	got, err := p.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() error = %v", err)
	}
	if string(got) != `{"close":12.5}` {
		t.Fatalf("ToJSON() = %s, want {\"close\":12.5}", got)
	}
}

func TestToJSON_NilDataIsNull(t *testing.T) {
	got, err := (UnifiedDataPoint{}).ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() error = %v", err)
	}
	if string(got) != "null" {
		t.Fatalf("ToJSON() = %s, want null", got)
	}
}

func TestToJSON_PropagatesUnsupportedValue(t *testing.T) {
	p := UnifiedDataPoint{Data: map[string]interface{}{"ch": make(chan int)}}
	if _, err := p.ToJSON(); err == nil || !strings.Contains(err.Error(), "unsupported type") {
		t.Fatalf("ToJSON() error = %v, want unsupported type", err)
	}
}

// TestVisibleAt_RoundTrips 的存在是硬要求：VisibleAt 不从类型里进，data 的
// 写入路径就得继续留一份本地镜像，X1a 的迁移等于没做。
func TestVisibleAt_RoundTrips(t *testing.T) {
	trade := time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC)
	visible := time.Date(2026, 4, 28, 0, 0, 0, 0, time.UTC)
	p := UnifiedDataPoint{Symbol: "600519.SH", TradeTime: trade, DataType: DataTypeFundamental, VisibleAt: visible}
	if p.VisibleAt != visible {
		t.Fatalf("VisibleAt = %v, want %v", p.VisibleAt, visible)
	}
	if p.VisibleAt.Equal(p.TradeTime) {
		t.Fatalf("VisibleAt must not collapse into TradeTime")
	}
}
