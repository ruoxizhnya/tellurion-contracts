package conformance

// 契约侧的契约一致性腿（QT-09 ③）：sides 含 "contracts" 的用例直接对契约函数实算。
//
// 这一侧是「真相源自检」：语料里写死的期望值必须与契约仓实现逐条对上。
// 它同时也钉住语料本身的可用性 —— 语料读不出来、为空、某条 sides 声明了
// contracts 却没人实现，都在这里当场红，而不是在 core/agent 侧变成静默跳过。

import (
	"bytes"
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/ruoxizhnya/tellurion-contracts/domain/market"
	contracts "github.com/ruoxizhnya/tellurion-contracts/engine"
	"github.com/ruoxizhnya/tellurion-contracts/fees"
)

// TestContractConformance 是契约仓的入口。
func TestContractConformance(t *testing.T) {
	cases, err := loadCases(".")
	if err != nil {
		t.Fatalf("加载契约语料失败: %v", err)
	}

	handlers := map[string]func(*testing.T, caseSpec){
		"stamp_tax_rate_for":  checkStampTaxRateFor,
		"fee_defaults":        checkFeeDefaults,
		"daily_visible_at":    checkDailyVisibleAt,
		"intraday_visible_at": checkIntradayVisibleAt,
		"visible_by":          checkVisibleBy,
		"time_semantics":      checkTimeSemantics,
		"run_spec_hash":       checkRunSpecHash,
	}

	ran, skipped := 0, 0
	for _, c := range cases {
		c := c
		if !declaresSide(c.Sides, "contracts") {
			skipped++
			t.Logf("跳过 %s（sides=%v，本侧不跑）", c.ID, c.Sides)
			continue
		}
		handler, ok := handlers[c.Kind]
		if !ok {
			t.Fatalf("用例 %q 声明了 contracts 侧，但 kind %q 在契约侧没有实现（未知 kind 或漏接线）", c.ID, c.Kind)
		}
		ran++
		t.Run(c.ID, func(t *testing.T) { handler(t, c) })
	}

	t.Logf("契约侧：执行 %d 条 / 跳过 %d 条 / 语料共 %d 条", ran, skipped, len(cases))
	if ran == 0 {
		t.Fatalf("契约侧一条用例都没执行 —— 语料与 sides 声明脱节")
	}
}

// ─── 各 kind 的断言 ──────────────────────────────────────────

func checkStampTaxRateFor(t *testing.T, c caseSpec) {
	t.Helper()
	var in struct {
		AsOf    string  `json:"as_of"`
		Current float64 `json:"current"`
		Before  float64 `json:"before"`
	}
	var want struct {
		Rate float64 `json:"rate"`
	}
	decodePair(t, c, &in, &want)

	got := fees.StampTaxRateFor(parseTimeArg(t, c, "as_of", in.AsOf), in.Current, in.Before)
	if !almostEqual(got, want.Rate) {
		t.Errorf("用例 %q: stamp_tax_rate_for 期望 %v，实际 %v", c.ID, want.Rate, got)
	}
}

// feeFields 是语料里五个费率字段的公共形状（input.fees 与 expect 共用）。
type feeFields struct {
	CommissionRate  float64 `json:"commission_rate"`
	StampTaxRate    float64 `json:"stamp_tax_rate"`
	TransferFeeRate float64 `json:"transfer_fee_rate"`
	MinCommission   float64 `json:"min_commission"`
	SlippageRate    float64 `json:"slippage_rate"`
}

func checkFeeDefaults(t *testing.T, c caseSpec) {
	t.Helper()
	var in struct {
		Base          string     `json:"base"`
		ApplyDefaults bool       `json:"apply_defaults"`
		Fees          *feeFields `json:"fees"`
	}
	var want feeFields
	decodePair(t, c, &in, &want)

	var got fees.AShareFees
	switch in.Base {
	case "", "defaults":
		got = fees.DefaultAShareFees()
	case "zero":
		got = fees.AShareFees{}
	default:
		t.Fatalf("用例 %q: 未知 base %q（合法：defaults / zero）", c.ID, in.Base)
	}
	if in.Fees != nil {
		got.CommissionRate = in.Fees.CommissionRate
		got.StampTaxRate = in.Fees.StampTaxRate
		got.TransferFeeRate = in.Fees.TransferFeeRate
		got.MinCommission = in.Fees.MinCommission
		got.SlippageRate = in.Fees.SlippageRate
	}
	if in.ApplyDefaults {
		got.ApplyDefaults()
	}

	compareFees(t, c, got, want)
}

func checkDailyVisibleAt(t *testing.T, c caseSpec) {
	t.Helper()
	var in struct {
		TradeDate string `json:"trade_date"`
	}
	var want struct {
		VisibleAt string `json:"visible_at"`
	}
	decodePair(t, c, &in, &want)

	got := market.DailyVisibleAt(parseTimeArg(t, c, "trade_date", in.TradeDate))
	if s := got.Format(time.RFC3339); s != want.VisibleAt {
		t.Errorf("用例 %q: daily_visible_at 期望 %s，实际 %s（零点值 %v）", c.ID, want.VisibleAt, s, got)
	}
}

func checkIntradayVisibleAt(t *testing.T, c caseSpec) {
	t.Helper()
	var in struct {
		BarStart string `json:"bar_start"`
		Minutes  int    `json:"minutes"`
	}
	var want struct {
		VisibleAt string `json:"visible_at"`
	}
	decodePair(t, c, &in, &want)

	got := market.IntradayVisibleAt(parseTimeArg(t, c, "bar_start", in.BarStart), in.Minutes)
	if s := got.Format(time.RFC3339); s != want.VisibleAt {
		t.Errorf("用例 %q: intraday_visible_at(%d) 期望 %s，实际 %s", c.ID, in.Minutes, want.VisibleAt, s)
	}
}

func checkVisibleBy(t *testing.T, c caseSpec) {
	t.Helper()
	var in struct {
		VisibleAt string `json:"visible_at"`
		AsOf      string `json:"as_of"`
	}
	var want struct {
		Visible bool `json:"visible"`
	}
	decodePair(t, c, &in, &want)

	bar := market.OHLCV{VisibleAt: parseTimeArg(t, c, "visible_at", in.VisibleAt)}
	if got := bar.VisibleBy(parseTimeArg(t, c, "as_of", in.AsOf)); got != want.Visible {
		t.Errorf("用例 %q: VisibleBy 期望 %v，实际 %v", c.ID, want.Visible, got)
	}
}

func checkTimeSemantics(t *testing.T, c caseSpec) {
	t.Helper()
	var in struct {
		Value string `json:"value"`
	}
	var want struct {
		OK        bool   `json:"ok"`
		Semantics string `json:"semantics"`
	}
	decodePair(t, c, &in, &want)

	got, err := market.ParseTimeSemantics(in.Value)
	if !want.OK {
		if err == nil {
			t.Errorf("用例 %q: ParseTimeSemantics(%q) 期望报错，实际返回 %q", c.ID, in.Value, got)
		}
		return
	}
	if err != nil {
		t.Fatalf("用例 %q: ParseTimeSemantics(%q) 意外报错: %v", c.ID, in.Value, err)
	}
	if string(got) != want.Semantics {
		t.Errorf("用例 %q: ParseTimeSemantics(%q) 期望 %q，实际 %q", c.ID, in.Value, want.Semantics, got)
	}
}

// runSpecCase 是语料里的 run_spec_hash 入参（与 contracts.RunSpecInput 字段一一对应）。
//
// 刻意与 contracts.RunSpecInput 分开声明：这一层是**语料 JSON 形状**，
// 转换写在 toRunSpec 里。这样契约仓改了字段名/类型时，这里会编译失败或解析失败，
// 而不是悄悄用零值算出一个「碰巧相同」的 hash。
type runSpecCase struct {
	Start                   string             `json:"start"`
	End                     string             `json:"end"`
	InitialCapital          float64            `json:"initial_capital"`
	RiskFreeRate            float64            `json:"risk_free_rate"`
	StrategyArtifactHash    string             `json:"strategy_artifact_hash"`
	StrategyVersionHash     string             `json:"strategy_version_hash"`
	StrategyParams          map[string]any     `json:"strategy_params"`
	StockPool               []string           `json:"stock_pool"`
	IndexCode               string             `json:"index_code"`
	Universe                []string           `json:"universe"`
	Frequency               string             `json:"frequency"`
	FillTiming              string             `json:"fill_timing"`
	Kernel                  string             `json:"kernel"`
	EngineVersion           string             `json:"engine_version"`
	Fees                    map[string]float64 `json:"fees"`
	Slippage                float64            `json:"slippage"`
	CostModel               string             `json:"cost_model"`
	LotSize                 int64              `json:"lot_size"`
	TPlus                   string             `json:"t_plus"`
	PriceLimitRule          string             `json:"price_limit_rule"`
	RiskRulesHash           string             `json:"risk_rules_hash"`
	DatasetSnapshotHash     string             `json:"dataset_snapshot_hash"`
	AdjustmentAnchorDigests []string           `json:"adjustment_anchor_digests"`
	Seed                    int64              `json:"seed"`
}

func checkRunSpecHash(t *testing.T, c caseSpec) {
	t.Helper()
	var in runSpecCase
	var want struct {
		Hash string `json:"hash"`
	}
	decodePair(t, c, &in, &want)

	spec := contracts.RunSpecInput{
		Start:                   parseDateArg(t, c, "start", in.Start),
		End:                     parseDateArg(t, c, "end", in.End),
		InitialCapital:          in.InitialCapital,
		RiskFreeRate:            in.RiskFreeRate,
		StrategyArtifactHash:    in.StrategyArtifactHash,
		StrategyVersionHash:     in.StrategyVersionHash,
		StrategyParamsEffective: in.StrategyParams,
		StockPool:               in.StockPool,
		IndexCode:               in.IndexCode,
		Universe:                in.Universe,
		Frequency:               in.Frequency,
		FillTiming:              contracts.FillTiming(in.FillTiming),
		Kernel:                  in.Kernel,
		EngineVersion:           in.EngineVersion,
		Fees:                    in.Fees,
		Slippage:                in.Slippage,
		CostModel:               in.CostModel,
		LotSize:                 in.LotSize,
		TPlus:                   in.TPlus,
		PriceLimitRule:          in.PriceLimitRule,
		RiskRulesHash:           in.RiskRulesHash,
		DatasetSnapshotHash:     in.DatasetSnapshotHash,
		AdjustmentAnchorDigests: in.AdjustmentAnchorDigests,
		Seed:                    in.Seed,
	}

	got, err := spec.Hash()
	if err != nil {
		t.Fatalf("用例 %q: RunSpecInput.Hash 报错: %v", c.ID, err)
	}
	if got != want.Hash {
		t.Errorf("用例 %q: run_spec_hash 期望 %s，实际 %s（编码漂移会让全部历史留痕失效）", c.ID, want.Hash, got)
	}
}

// ─── 公共工具 ────────────────────────────────────────────────

func declaresSide(sides []string, side string) bool {
	for _, s := range sides {
		if s == side {
			return true
		}
	}
	return false
}

// decodePair 严格解析 input/expect：未知字段一律报错。
//
// 宽松解析会让语料里的键名拼写错误静默变成零值，于是一条本该红的用例
// 拿着「零值输入 + 零值期望」变绿 —— 这正是语料最需要防的失效方式。
func decodePair(t *testing.T, c caseSpec, input, expect any) {
	t.Helper()
	if err := decodeStrict(c.Input, input); err != nil {
		t.Fatalf("用例 %q 的 input 解析失败: %v", c.ID, err)
	}
	if err := decodeStrict(c.Expect, expect); err != nil {
		t.Fatalf("用例 %q 的 expect 解析失败: %v", c.ID, err)
	}
}

func decodeStrict(raw json.RawMessage, dst any) error {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

// parseTimeArg 解析 RFC3339 时间；空串表示零值时间（语料用它表达「调用方拿不到日期」）。
func parseTimeArg(t *testing.T, c caseSpec, field, value string) time.Time {
	t.Helper()
	if value == "" {
		return time.Time{}
	}
	ts, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatalf("用例 %q 的 %s=%q 不是 RFC3339: %v", c.ID, field, value, err)
	}
	return ts
}

// parseDateArg 解析 UTC 日历日（run 区间是交易日，不是时刻）。
func parseDateArg(t *testing.T, c caseSpec, field, value string) time.Time {
	t.Helper()
	ts, err := time.Parse("2006-01-02", value)
	if err != nil {
		t.Fatalf("用例 %q 的 %s=%q 不是 YYYY-MM-DD: %v", c.ID, field, value, err)
	}
	return ts
}

// almostEqual 是费率这类小数的比较口径：1e-12 远小于任何有意义的费率差，
// 又足以容纳浮点乘除的舍入噪声。
func almostEqual(got, want float64) bool {
	return math.Abs(got-want) <= 1e-12
}

func compareFees(t *testing.T, c caseSpec, got fees.AShareFees, want feeFields) {
	t.Helper()
	for _, f := range []struct {
		name     string
		got      float64
		expected float64
	}{
		{"commission_rate", got.CommissionRate, want.CommissionRate},
		{"stamp_tax_rate", got.StampTaxRate, want.StampTaxRate},
		{"transfer_fee_rate", got.TransferFeeRate, want.TransferFeeRate},
		{"min_commission", got.MinCommission, want.MinCommission},
		{"slippage_rate", got.SlippageRate, want.SlippageRate},
	} {
		if !almostEqual(f.got, f.expected) {
			t.Errorf("用例 %q: %s 期望 %v，实际 %v", c.ID, f.name, f.expected, f.got)
		}
	}
}
