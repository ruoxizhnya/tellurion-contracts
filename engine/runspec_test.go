package contracts

import (
	"encoding/json"
	"errors"
	"math"
	"strings"
	"testing"
)

// K9 契约腿：RunSpec 的内核路由字段。
//
// 三条腿：
//   1. kernel 透传 —— JSON → struct → JSON 往返不丢 kernel / engine_params，
//      缺省时 EffectiveKernel() == "go"；
//   2. 存量请求序列化行为不变 —— 不带新字段的请求序列化后**逐字节**没有
//      kernel / engine_params 键；
//   3. engine_params 未知名 fail-closed 拒绝（ValidateRunSpec 非 nil），
//      而不是被静默忽略。

func TestRunSpec_KernelPassthrough(t *testing.T) {
	const body = `{"strategy":"momentum","stock_pool":["600000.SH"],"start_date":"2024-01-01","end_date":"2024-06-30","kernel":"nautilus","engine_params":{"timeout_sec":30}}`

	var req BacktestRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if req.Kernel != "nautilus" {
		t.Errorf("Kernel = %q, want %q（kernel 必须随请求透传，不得被丢弃）", req.Kernel, "nautilus")
	}
	if req.EffectiveKernel() != "nautilus" {
		t.Errorf("EffectiveKernel() = %q, want %q", req.EffectiveKernel(), "nautilus")
	}
	if len(req.EngineParams) != 1 || req.EngineParams["timeout_sec"] != float64(30) {
		t.Errorf("EngineParams = %#v, want {timeout_sec:30}（engine_params 必须随请求透传）", req.EngineParams)
	}

	// 往返：序列化再反序列化后仍然一致（透传到底，不被中间层吃掉）。
	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var round BacktestRequest
	if err := json.Unmarshal(raw, &round); err != nil {
		t.Fatalf("round-trip unmarshal: %v", err)
	}
	if round.Kernel != req.Kernel || len(round.EngineParams) != len(req.EngineParams) {
		t.Errorf("round-trip 丢字段: kernel=%q engine_params=%#v", round.Kernel, round.EngineParams)
	}
}

func TestRunSpec_DefaultKernelIsGo(t *testing.T) {
	var req BacktestRequest
	if req.Kernel != "" {
		t.Fatalf("零值 Kernel 应为空串（默认值不回填字段）, got %q", req.Kernel)
	}
	if got := req.EffectiveKernel(); got != DefaultKernel {
		t.Errorf("EffectiveKernel() = %q, want %q（未指定内核时默认 go）", got, DefaultKernel)
	}
}

func TestRunSpec_LegacyRequestSerializationUnchanged(t *testing.T) {
	// 不带新字段的存量请求：序列化结果必须与加字段之前**逐字节相同**——
	// 不得出现 "kernel" / "engine_params" 键。带 omitempty 是为了这一条；
	// 把默认值 "go" 回填进字段会让这条腿立刻红。
	req := BacktestRequest{
		Strategy:       "momentum",
		StockPool:      []string{"600000.SH"},
		StartDate:      "2024-01-01",
		EndDate:        "2024-06-30",
		InitialCapital: 1_000_000,
	}
	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	// golden = 加 kernel/engine_params **之前**该请求的真实序列化（index_code /
	// risk_free_rate 本就无 omitempty，一直如此）。断言逐字节相等 ⇒ 新字段
	// 只要漏进输出这条腿就红。
	const want = `{"strategy":"momentum","stock_pool":["600000.SH"],"index_code":"","start_date":"2024-01-01","end_date":"2024-06-30","initial_capital":1000000,"risk_free_rate":0}`
	if string(raw) != want {
		t.Errorf("存量请求序列化行为变了:\n got: %s\nwant: %s", raw, want)
	}
}

func TestRunSpec_ValidateRejectsUnknownEngineParams(t *testing.T) {
	req := BacktestRequest{
		Strategy:     "momentum",
		Kernel:       DefaultKernel,
		EngineParams: map[string]any{"typo_param": 1},
	}
	err := req.ValidateRunSpec()
	if err == nil {
		t.Fatal("未知名 engine_params 必须被拒绝（fail-closed），不能静默忽略")
	}
	if !strings.Contains(err.Error(), "typo_param") {
		t.Errorf("错误必须点名未知键, got: %v", err)
	}
	if !strings.Contains(err.Error(), "rejected, not ignored") {
		t.Errorf("错误必须说明「拒绝而非忽略」语义, got: %v", err)
	}
}

func TestRunSpec_ValidateRejectsUnknownParamsOnUnregisteredKernel(t *testing.T) {
	// 未登记内核的已知参数集为空 ⇒ 任何 engine_params 都是未知名。
	// （内核本身认不认识是 K10 路由的事，参数维度先 fail-closed。）
	req := BacktestRequest{
		Kernel:       "nautilus",
		EngineParams: map[string]any{"timeout_sec": 30},
	}
	if err := req.ValidateRunSpec(); err == nil {
		t.Fatal("未登记内核 + engine_params 必须被拒绝")
	}
}

func TestRunSpec_ValidateAcceptsNoParams(t *testing.T) {
	for _, req := range []BacktestRequest{
		{Strategy: "momentum"},
		{Strategy: "momentum", Kernel: "nautilus"},
		{Strategy: "momentum", Kernel: DefaultKernel, EngineParams: map[string]any{}},
	} {
		if err := req.ValidateRunSpec(); err != nil {
			t.Errorf("无引擎参数的请求不得被拒: %#v -> %v", req, err)
		}
	}
}

func TestRunSpec_ValidateListsAllUnknownKeys(t *testing.T) {
	req := BacktestRequest{
		EngineParams: map[string]any{"bbb": 1, "aaa": 2},
	}
	err := req.ValidateRunSpec()
	if err == nil {
		t.Fatal("expected rejection")
	}
	// 排序稳定输出：aaa 在前，两个都必须点名。
	msg := err.Error()
	if !strings.Contains(msg, "aaa") || !strings.Contains(msg, "bbb") {
		t.Errorf("所有未知键都要点名, got: %v", err)
	}
	if strings.Index(msg, "aaa") > strings.Index(msg, "bbb") {
		t.Errorf("未知键输出应排序稳定, got: %v", err)
	}
}

// K15 契约腿：signals（目标权重流）的 fail-closed 校验。
//
// 反证腿三例（与侧车 runspec._parse_signals 同口径）：非法值拒绝而非钳制
// /跳过 —— 一行 weight=1.5 的脏数据被钳到 1.0 会让双跑对齐报告多出「查无
// 实据」的差异。
func TestRunSpec_ValidateSignalsRejectsInvalidRows(t *testing.T) {
	for _, tc := range []struct {
		name    string
		signals []SignalWeight
		wantMsg string
	}{
		{
			name:    "weight above one",
			signals: []SignalWeight{{Symbol: "000001.SZ", Weight: 1.01, AsOf: "2024-01-02"}},
			wantMsg: "signals[0].weight",
		},
		{
			name:    "negative weight",
			signals: []SignalWeight{{Symbol: "000001.SZ", Weight: -0.1, AsOf: "2024-01-02"}},
			wantMsg: "signals[0].weight",
		},
		{
			name:    "NaN weight",
			signals: []SignalWeight{{Symbol: "000001.SZ", Weight: math.NaN(), AsOf: "2024-01-02"}},
			wantMsg: "signals[0].weight",
		},
		{
			name:    "blank symbol",
			signals: []SignalWeight{{Symbol: "  ", Weight: 0.5, AsOf: "2024-01-02"}},
			wantMsg: "signals[0].symbol",
		},
		{
			name:    "non ISO date",
			signals: []SignalWeight{{Symbol: "000001.SZ", Weight: 0.5, AsOf: "20240102"}},
			wantMsg: "signals[0].as_of",
		},
		{
			name:    "second row invalid",
			signals: []SignalWeight{{Symbol: "000001.SZ", Weight: 0.5, AsOf: "2024-01-02"}, {Symbol: "600000.SH", Weight: 0.5, AsOf: "2024-01-02T00:00:00Z"}},
			wantMsg: "signals[1].as_of",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := BacktestRequest{Strategy: "weight_replay", Signals: tc.signals}
			err := req.ValidateRunSpec()
			if err == nil {
				t.Fatalf("%s 必须被拒绝（fail-closed）", tc.name)
			}
			if !strings.Contains(err.Error(), tc.wantMsg) {
				t.Fatalf("错误必须点名 %q, got: %v", tc.wantMsg, err)
			}
		})
	}
}

func TestRunSpec_ValidateSignalsAcceptsValidRows(t *testing.T) {
	req := BacktestRequest{
		Strategy: "weight_replay",
		Signals: []SignalWeight{
			{Symbol: "000001.SZ", Weight: 0.5, AsOf: "2024-01-02"},
			{Symbol: "600000.SH", Weight: 0, AsOf: "2024-12-31"},
			{Symbol: "600000.SH", Weight: 1, AsOf: "2024-06-03"},
		},
	}
	if err := req.ValidateRunSpec(); err != nil {
		t.Fatalf("合法信号流不得被拒: %v", err)
	}
	// 边界：weight 恰 0 与 1 合法（0 = 全平，1 = 全仓）。
}

func TestRunSpec_SignalsWireRoundTrip(t *testing.T) {
	// signals 随 JSON 往返不丢（NautilusRunner 是直拷映射，字段名必须对齐）。
	const body = `{"strategy":"weight_replay","start_date":"2024-01-02","end_date":"2024-12-31","signals":[{"symbol":"000001.SZ","weight":0.5,"as_of":"2024-01-02"}]}`
	var req BacktestRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(req.Signals) != 1 || req.Signals[0].Symbol != "000001.SZ" || req.Signals[0].Weight != 0.5 || req.Signals[0].AsOf != "2024-01-02" {
		t.Fatalf("signals 反序列化失真: %#v", req.Signals)
	}
	if err := req.ValidateRunSpec(); err != nil {
		t.Fatalf("合法请求被拒: %v", err)
	}
}

// ─── QT-20：strategy_ref / strategy_params 的契约腿 ───────────────────

// 引用串的 wire 形态必须与 provenance.go 的 ParseStrategyRef / String()
// **互逆**（`id` / `id@hash`）：JSON 往返、解析、还原三步同一个值。
// 这条同时钉住「不得新增第二种引用形状」——一旦引入别的分隔符或前缀，
// 往返就会失真。
func TestRunSpec_StrategyRefWireRoundTrip(t *testing.T) {
	artifact := strings.Repeat("ab", 32)
	body := `{"strategy":"momentum","start_date":"2024-01-01","end_date":"2024-06-30",` +
		`"strategy_ref":"momentum@` + artifact + `","strategy_params":{"fast":5,"slow":20.5,"enabled":true,"note":null}}`

	var req BacktestRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := req.ValidateRunSpec(); err != nil {
		t.Fatalf("合法的 strategy_ref / strategy_params 不得被拒: %v", err)
	}

	ref, err := ParseStrategyRef(req.StrategyRef)
	if err != nil {
		t.Fatalf("ParseStrategyRef(%q): %v", req.StrategyRef, err)
	}
	if ref.ID != "momentum" || ref.Hash != artifact {
		t.Fatalf("解析结果失真: %+v", ref)
	}
	if got := ref.String(); got != req.StrategyRef {
		t.Fatalf("String() 与 wire 形态不互逆: got %q, want %q", got, req.StrategyRef)
	}

	// latest 形态（不带 hash）同样往返。
	latest := BacktestRequest{Strategy: "momentum", StrategyRef: "momentum"}
	if err := latest.ValidateRunSpec(); err != nil {
		t.Fatalf("latest 引用（无 hash）必须合法: %v", err)
	}
	lRef, err := ParseStrategyRef(latest.StrategyRef)
	if err != nil || !lRef.IsLatest() || lRef.String() != "momentum" {
		t.Fatalf("latest 往返失真: ref=%+v err=%v", lRef, err)
	}

	// 序列化再反序列化：两个新字段都要在（有值时不能丢）。
	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var round BacktestRequest
	if err := json.Unmarshal(raw, &round); err != nil {
		t.Fatalf("round-trip unmarshal: %v", err)
	}
	if round.StrategyRef != req.StrategyRef {
		t.Fatalf("往返丢 strategy_ref: %q != %q", round.StrategyRef, req.StrategyRef)
	}
	if len(round.StrategyParams) != len(req.StrategyParams) {
		t.Fatalf("往返丢 strategy_params: %#v", round.StrategyParams)
	}
	if round.StrategyParams["slow"] != 20.5 {
		t.Fatalf("strategy_params 数值失真: %#v", round.StrategyParams["slow"])
	}
}

// 形状不对的引用一律拒绝（fail-closed），且错误必须可 errors.Is 到
// ErrStrategyRefMalformed —— 调用方不该靠字符串匹配认这条。
func TestRunSpec_ValidateRejectsMalformedStrategyRef(t *testing.T) {
	valid64 := strings.Repeat("a", 64)
	for _, tc := range []struct {
		name string
		ref  string
	}{
		{"empty id", "@" + valid64},
		{"empty hash", "momentum@"},
		{"short hash", "momentum@" + strings.Repeat("a", 63)},
		{"long hash", "momentum@" + strings.Repeat("a", 65)},
		{"two separators", "momentum@a@b"},
		{"uppercase hex", "momentum@" + strings.Repeat("A", 64)},
		{"non hex", "momentum@" + strings.Repeat("z", 64)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := BacktestRequest{Strategy: "momentum", StrategyRef: tc.ref}
			err := req.ValidateRunSpec()
			if err == nil {
				t.Fatalf("%s 的引用 %q 必须被拒绝", tc.name, tc.ref)
			}
			if !errors.Is(err, ErrStrategyRefMalformed) {
				t.Fatalf("错误必须可 errors.Is 到 ErrStrategyRefMalformed, got: %v", err)
			}
		})
	}
}

// 嵌套键必须被拒：浅层对象是「一个键一个标量」，嵌套会让浅合并与
// params_hash 的编码失去唯一定义。
func TestRunSpec_ValidateRejectsNestedStrategyParams(t *testing.T) {
	for _, tc := range []struct {
		name   string
		params map[string]any
		key    string
	}{
		{"nested object", map[string]any{"risk": map[string]any{"max": 0.2}}, "risk"},
		{"nested array", map[string]any{"windows": []any{5, 10}}, "windows"},
		{"typed map", map[string]any{"risk": map[string]float64{"max": 0.2}}, "risk"},
		{"typed slice", map[string]any{"windows": []int{5, 10}}, "windows"},
		{"blank key", map[string]any{" ": 1}, " "},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := BacktestRequest{Strategy: "momentum", StrategyParams: tc.params}
			err := req.ValidateRunSpec()
			if err == nil {
				t.Fatalf("%s 必须被拒绝（只接受浅层标量）", tc.name)
			}
			if !strings.Contains(err.Error(), "strategy_params") {
				t.Fatalf("错误必须点名 strategy_params, got: %v", err)
			}
			if tc.key != " " && !strings.Contains(err.Error(), tc.key) {
				t.Fatalf("错误必须点名嵌套键 %q, got: %v", tc.key, err)
			}
		})
	}
}

// 合法浅层对象（含各数值宽度与 nil）不得误拒。
func TestRunSpec_ValidateAcceptsScalarStrategyParams(t *testing.T) {
	req := BacktestRequest{
		Strategy: "momentum",
		StrategyParams: map[string]any{
			"fast":    int(5),
			"slow":    int64(20),
			"ratio":   0.25,
			"factor":  float32(1.5),
			"enabled": true,
			"label":   "v1",
			"unset":   nil,
		},
	}
	if err := req.ValidateRunSpec(); err != nil {
		t.Fatalf("合法浅层参数不得被拒: %v", err)
	}
}
