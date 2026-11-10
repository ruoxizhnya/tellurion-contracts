package risk

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

// fourTypes 造「四类各恰好一条」的最小快照规则集（params 空对象即可，
// ValidateScope 是纯结构校验，不看 params 内容）。
func fourTypes(scope Scope) []SnapshotRule {
	types := []RuleType{
		RuleTypePositionSizing,
		RuleTypeStopLoss,
		RuleTypeRegime,
		RuleTypeTakeProfit,
	}
	out := make([]SnapshotRule, 0, len(types))
	for _, typ := range types {
		out = append(out, SnapshotRule{
			RuleID:      string(scope) + "_" + string(typ),
			Revision:    1,
			Scope:       scope,
			Type:        typ,
			Params:      json.RawMessage(`{}`),
			VersionHash: "x",
		})
	}
	return out
}

// TestScope_ProcessLevel 钉住 2026-11-11 主任裁决：只有 global / strategy
// 参与进程级策略源选择，account / symbol 不参与。
func TestScope_ProcessLevel(t *testing.T) {
	want := map[Scope]bool{
		ScopeGlobal:        true,
		ScopeStrategy:      true,
		ScopeAccount:       false,
		ScopeSymbol:        false,
		Scope(""):          false,
		Scope("portfolio"): false,
	}
	for scope, expected := range want {
		if got := scope.ProcessLevel(); got != expected {
			t.Errorf("Scope(%q).ProcessLevel() = %v, want %v", scope, got, expected)
		}
	}
}

// TestSnapshot_ValidateScope 钉住「选中 scope 内四类各恰好一条」：
// 缺一类或同类多条都拒（歧义不猜），未知 scope / type fail-closed。
func TestSnapshot_ValidateScope(t *testing.T) {
	t.Run("四类各一条", func(t *testing.T) {
		snap := Snapshot{Rules: fourTypes(ScopeGlobal)}
		if err := snap.ValidateScope(ScopeGlobal); err != nil {
			t.Fatalf("四类各一条应通过，got %v", err)
		}
	})

	t.Run("缺一类", func(t *testing.T) {
		rules := fourTypes(ScopeGlobal)[:3] // 缺 take_profit
		err := Snapshot{Rules: rules}.ValidateScope(ScopeGlobal)
		if !errors.Is(err, ErrInvalidParams) {
			t.Fatalf("errors.Is(%v, ErrInvalidParams) = false", err)
		}
		if !strings.Contains(err.Error(), `missing rule type "take_profit"`) {
			t.Fatalf("错误文本 %q 没点名缺的规则类型", err)
		}
	})

	t.Run("同类多条", func(t *testing.T) {
		rules := fourTypes(ScopeGlobal)
		rules = append(rules, SnapshotRule{
			RuleID: "extra_stop", Revision: 1, Scope: ScopeGlobal,
			Type: RuleTypeStopLoss, Params: json.RawMessage(`{}`), VersionHash: "x",
		})
		err := Snapshot{Rules: rules}.ValidateScope(ScopeGlobal)
		if !errors.Is(err, ErrInvalidParams) {
			t.Fatalf("errors.Is(%v, ErrInvalidParams) = false", err)
		}
		if !strings.Contains(err.Error(), `has 2 rules of type "stop_loss"`) {
			t.Fatalf("错误文本 %q 没点名重复的规则类型", err)
		}
	})

	t.Run("未知 scope", func(t *testing.T) {
		err := Snapshot{Rules: fourTypes(ScopeGlobal)}.ValidateScope(Scope("portfolio"))
		if !errors.Is(err, ErrInvalidScope) {
			t.Fatalf("errors.Is(%v, ErrInvalidScope) = false", err)
		}
	})

	t.Run("未知 type", func(t *testing.T) {
		rules := fourTypes(ScopeGlobal)
		rules[0].Type = RuleType("vol_target")
		err := Snapshot{Rules: rules}.ValidateScope(ScopeGlobal)
		if !errors.Is(err, ErrInvalidRuleType) {
			t.Fatalf("errors.Is(%v, ErrInvalidRuleType) = false", err)
		}
	})

	t.Run("别的 scope 不顶数", func(t *testing.T) {
		global := fourTypes(ScopeGlobal)[:3]   // 缺 take_profit
		account := fourTypes(ScopeAccount)[3:] // 恰好只有 take_profit
		snap := Snapshot{Rules: append(global, account...)}
		if err := snap.ValidateScope(ScopeGlobal); !errors.Is(err, ErrInvalidParams) {
			t.Fatalf("account 的 take_profit 不能补 global 的缺口，got %v", err)
		}
	})
}
