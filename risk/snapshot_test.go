package risk

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

// ct03Rules 是固定向量的输入：**刻意乱序**（zz 在 aa 前）且首条 params 带多余空白，
// 同时咬住「按 rule_id 排序」与「params 规范化」两条编码细节。
func ct03Rules() []SnapshotRule {
	return []SnapshotRule{
		{
			RuleID:      "zz_guard",
			Revision:    2,
			Scope:       ScopeAccount,
			Type:        RuleTypeRegime,
			Params:      json.RawMessage(`{"max_drawdown": 0.15}`),
			VersionHash: "1111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			RuleID:      "aa_stop",
			Revision:    1,
			Scope:       ScopeSymbol,
			Type:        RuleTypeStopLoss,
			Params:      json.RawMessage(`{"pct":0.05}`),
			VersionHash: "2222222222222222222222222222222222222222222222222222222222222222",
		},
	}
}

// TestGoldenVector_SnapshotHash 钉死 risk 仓真实实现算出的 hex
// （2026-11-07 在 `tellurion-risk` 内调 `rule.SnapshotHash(ct03Rules())` 取得，
// 该次命令退出码 0）。向量变 = 编码漂了。
func TestGoldenVector_SnapshotHash(t *testing.T) {
	const want = "14ecab865bebaef54052354b63caa5d351f533ab112e1a29a9dcad3707955c5b"
	got, err := Snapshot{Rules: ct03Rules()}.Hash()
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}
	if got != want {
		t.Fatalf("Hash() = %s, want %s", got, want)
	}
}

// TestSnapshotHash_MatchesRiskRepoEncoding 逐行重建 `tellurion-risk:internal/rule/snapshot.go`
// 的字节流（长度前缀帧 + 域标签 + 裸换行 + 按 rule_id 排序 + UseNumber→Marshal），
// 断言与本包实现同 hex。这是本单的跨仓对拍腿：任一细节漂移，这条先红。
func TestSnapshotHash_MatchesRiskRepoEncoding(t *testing.T) {
	rules := ct03Rules()
	want, err := riskRepoEncoding(rules)
	if err != nil {
		t.Fatalf("risk 侧重建编码失败：%v", err)
	}
	got, err := Snapshot{Rules: rules}.Hash()
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}
	if got != want {
		t.Fatalf("两侧 hex 不一致：contracts %s / risk %s", got, want)
	}
}

// riskRepoEncoding 是 risk 侧编码器的独立重建（不调用本包任何编码辅助函数）。
func riskRepoEncoding(rules []SnapshotRule) (string, error) {
	sorted := make([]SnapshotRule, len(rules))
	copy(sorted, rules)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].RuleID < sorted[j].RuleID })

	var b strings.Builder
	b.WriteString(DomainSnapshotHash)
	b.WriteByte('\n')
	fmt.Fprintf(&b, "%d:%s=%d:%s\n", len("rule_count"), "rule_count", len(strconv.Itoa(len(sorted))), strconv.Itoa(len(sorted)))
	for _, r := range sorted {
		dec := json.NewDecoder(bytes.NewReader(r.Params))
		dec.UseNumber()
		var v any
		if err := dec.Decode(&v); err != nil {
			return "", err
		}
		params, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		for _, kv := range [][2]string{
			{"rule_id", r.RuleID},
			{"revision", strconv.Itoa(r.Revision)},
			{"scope", string(r.Scope)},
			{"type", string(r.Type)},
			{"params", string(params)},
			{"version_hash", r.VersionHash},
		} {
			fmt.Fprintf(&b, "%d:%s=%d:%s\n", len(kv[0]), kv[0], len(kv[1]), kv[1])
		}
	}
	sum := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(sum[:]), nil
}

// TestSnapshotHash_OrderInsensitive 同内容不同写入顺序必须同 hash（排序在编码内完成）。
func TestSnapshotHash_OrderInsensitive(t *testing.T) {
	rules := ct03Rules()
	reversed := []SnapshotRule{rules[1], rules[0]}
	a, err := Snapshot{Rules: rules}.Hash()
	if err != nil {
		t.Fatal(err)
	}
	b, err := Snapshot{Rules: reversed}.Hash()
	if err != nil {
		t.Fatal(err)
	}
	if a != b {
		t.Fatalf("顺序不同导致 hash 不同：%s / %s", a, b)
	}
}

// TestSnapshotHash_IgnoresCreatedAtAndClaimedHash hash 只覆盖内容：两个元数据字段都不进编码。
func TestSnapshotHash_IgnoresCreatedAtAndClaimedHash(t *testing.T) {
	base := Snapshot{Rules: ct03Rules()}
	one, err := base.Hash()
	if err != nil {
		t.Fatal(err)
	}
	meta := Snapshot{
		RiskRulesHash: "deadbeef",
		Rules:         ct03Rules(),
		CreatedAt:     time.Date(2026, 11, 7, 0, 0, 0, 0, time.UTC),
	}
	two, err := meta.Hash()
	if err != nil {
		t.Fatal(err)
	}
	if one != two {
		t.Fatalf("CreatedAt / RiskRulesHash 混进了 hash：%s / %s", one, two)
	}
}

// TestSnapshotHash_Rejections 逐条钉住拒收原因与哨兵身份。
func TestSnapshotHash_Rejections(t *testing.T) {
	valid := ct03Rules()[1]
	cases := []struct {
		name   string
		rules  []SnapshotRule
		wantIs error
		want   string
	}{
		{name: "空快照", rules: nil, wantIs: ErrInvalidParams, want: "at least one active rule"},
		{
			name:   "重复 rule_id",
			rules:  []SnapshotRule{valid, valid},
			wantIs: ErrInvalidParams,
			want:   `duplicate rule_id "aa_stop"`,
		},
		{
			name:   "空 rule_id",
			rules:  []SnapshotRule{{Revision: 1, Scope: ScopeGlobal, Type: RuleTypeRegime, Params: json.RawMessage(`{}`), VersionHash: "x"}},
			wantIs: ErrInvalidParams,
			want:   "rule_id is required",
		},
		{
			name: "revision 非正",
			rules: []SnapshotRule{{
				RuleID: "r", Revision: 0, Scope: ScopeGlobal, Type: RuleTypeRegime,
				Params: json.RawMessage(`{}`), VersionHash: "x",
			}},
			wantIs: ErrInvalidParams,
			want:   "r revision must be > 0",
		},
		{
			name: "未知 type",
			rules: []SnapshotRule{{
				RuleID: "r", Revision: 1, Scope: ScopeGlobal, Type: RuleType("vol_target"),
				Params: json.RawMessage(`{}`), VersionHash: "x",
			}},
			wantIs: ErrInvalidRuleType,
			want:   `"vol_target"`,
		},
		{
			name: "未知 scope",
			rules: []SnapshotRule{{
				RuleID: "r", Revision: 1, Scope: Scope("portfolio"), Type: RuleTypeRegime,
				Params: json.RawMessage(`{}`), VersionHash: "x",
			}},
			wantIs: ErrInvalidScope,
			want:   `"portfolio"`,
		},
		{
			name: "缺 version_hash",
			rules: []SnapshotRule{{
				RuleID: "r", Revision: 1, Scope: ScopeGlobal, Type: RuleTypeRegime,
				Params: json.RawMessage(`{}`),
			}},
			wantIs: ErrInvalidParams,
			want:   "r version_hash is required",
		},
		{
			name: "params 不是对象",
			rules: []SnapshotRule{{
				RuleID: "r", Revision: 1, Scope: ScopeGlobal, Type: RuleTypeRegime,
				Params: json.RawMessage(`[1,2]`), VersionHash: "x",
			}},
			wantIs: ErrInvalidParams,
			want:   "params must be a JSON object",
		},
		{
			name: "params 为空",
			rules: []SnapshotRule{{
				RuleID: "r", Revision: 1, Scope: ScopeGlobal, Type: RuleTypeRegime,
				Params: json.RawMessage(`  `), VersionHash: "x",
			}},
			wantIs: ErrInvalidParams,
			want:   "params is required",
		},
		{
			name: "params 有尾随内容",
			rules: []SnapshotRule{{
				RuleID: "r", Revision: 1, Scope: ScopeGlobal, Type: RuleTypeRegime,
				Params: json.RawMessage(`{} {}`), VersionHash: "x",
			}},
			wantIs: ErrInvalidParams,
			want:   "params has trailing content",
		},
		{
			name: "params 超限",
			rules: []SnapshotRule{{
				RuleID: "r", Revision: 1, Scope: ScopeGlobal, Type: RuleTypeRegime,
				Params:      json.RawMessage(fmt.Sprintf(`{"blob":%q}`, strings.Repeat("x", MaxParamsBytes))),
				VersionHash: "x",
			}},
			wantIs: ErrInvalidParams,
			want:   "params exceed 65536 bytes",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Snapshot{Rules: tc.rules}.Hash()
			if err == nil {
				t.Fatalf("必须报错，却绿了")
			}
			if !errors.Is(err, tc.wantIs) {
				t.Fatalf("errors.Is(%v, %v) = false", err, tc.wantIs)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("错误文本 %q 不含 %q", err.Error(), tc.want)
			}
		})
	}
}

// TestSnapshot_Verify 声称 hash 与内容不符必须 fail-loud。
func TestSnapshot_Verify(t *testing.T) {
	snap := Snapshot{Rules: ct03Rules()}
	h, err := snap.Hash()
	if err != nil {
		t.Fatal(err)
	}
	snap.RiskRulesHash = h
	if err := snap.Verify(); err != nil {
		t.Fatalf("正确 hash 应通过，got %v", err)
	}
	snap.RiskRulesHash = strings.Repeat("0", 64)
	err = snap.Verify()
	if !errors.Is(err, ErrSnapshotMismatch) {
		t.Fatalf("errors.Is(%v, ErrSnapshotMismatch) = false", err)
	}
	if !strings.Contains(err.Error(), "derived "+h) {
		t.Fatalf("错误文本没带 derive 值：%v", err)
	}
}

// TestSnapshot_Clone_DeepCopiesParams 改副本的 Params 不得影响原件。
func TestSnapshot_Clone_DeepCopiesParams(t *testing.T) {
	orig := Snapshot{Rules: ct03Rules()}
	clone := orig.Clone()
	clone.Rules[0].Params[0] = '!'
	clone.Rules[1].RuleID = "mutated"
	if bytes.Equal(orig.Rules[0].Params, clone.Rules[0].Params) {
		t.Fatal("Params 未深拷贝")
	}
	if orig.Rules[1].RuleID == "mutated" {
		t.Fatal("Rules 切片未深拷贝")
	}
}

// TestScope_Valid 词表是冻结的四个取值，其余一律 false。
func TestScope_Valid(t *testing.T) {
	for _, s := range []Scope{ScopeGlobal, ScopeAccount, ScopeStrategy, ScopeSymbol} {
		if !s.Valid() {
			t.Errorf("%q 应在冻结枚举里", s)
		}
	}
	for _, s := range []Scope{"", "portfolio", "GLOBAL"} {
		if s.Valid() {
			t.Errorf("%q 不该在冻结枚举里", s)
		}
	}
}

// TestFrozenLiterals 取值即契约：本轮搬的是字面量，改任何一个都是跨仓 breaking。
func TestFrozenLiterals(t *testing.T) {
	if DomainSnapshotHash != "tellurion.risk_snapshot_hash.v1" {
		t.Fatalf("DomainSnapshotHash = %q", DomainSnapshotHash)
	}
	if MaxParamsBytes != 65536 {
		t.Fatalf("MaxParamsBytes = %d", MaxParamsBytes)
	}
	wantTypes := map[RuleType]string{
		RuleTypePositionSizing: "position_sizing",
		RuleTypeStopLoss:       "stop_loss",
		RuleTypeRegime:         "regime",
		RuleTypeTakeProfit:     "take_profit",
	}
	for typ, want := range wantTypes {
		if string(typ) != want {
			t.Errorf("RuleType %q != %q", typ, want)
		}
	}
	wantScopes := map[Scope]string{
		ScopeGlobal:   "global",
		ScopeAccount:  "account",
		ScopeStrategy: "strategy",
		ScopeSymbol:   "symbol",
	}
	for sc, want := range wantScopes {
		if string(sc) != want {
			t.Errorf("Scope %q != %q", sc, want)
		}
	}
}

// TestSentinelsStayDistinct 四条哨兵互不命中（`errors.Is` 身份唯一）。
func TestSentinelsStayDistinct(t *testing.T) {
	all := []error{ErrInvalidRuleType, ErrInvalidScope, ErrInvalidParams, ErrSnapshotMismatch}
	for i, a := range all {
		for j, b := range all {
			if i != j && errors.Is(a, b) {
				t.Errorf("哨兵 %d 与 %d 混了：%v", i, j, a)
			}
		}
	}
}
