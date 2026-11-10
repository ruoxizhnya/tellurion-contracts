package risk

import "fmt"

// ProcessLevel 报告该 scope 是否参与进程级策略源的选择。
//
// 2026-11-11 主任裁决（正式口径）：core 的 pickScope 只读 `global` / `strategy`，
// 顺序是「有 global 用 global，否则退 strategy」；`account` / `symbol` 不参与。
// 纯函数，零 I/O（ADR-0004）。
func (s Scope) ProcessLevel() bool {
	switch s {
	case ScopeGlobal, ScopeStrategy:
		return true
	default:
		return false
	}
}

// ValidateScope 校验某个 scope 下「四类规则各恰好一条」。
//
// 2026-11-11 主任裁决：选中 scope 内四类规则各恰好一条；缺一类或同类多条
// ⇒ 拒（歧义不猜），不得回落默认限额。未知 scope ⇒ ErrInvalidScope，
// 未知 type ⇒ ErrInvalidRuleType（旧读者 fail-closed）。
//
// 纯结构校验：不做 params decode、不填默认值、不建进程内索引（CT-03 ④ 的边界
// 不变；`global → strategy` 的 fallback 选择仍属服务策略）。
func (s Snapshot) ValidateScope(scope Scope) error {
	if !scope.Valid() {
		return fmt.Errorf("%w: %q", ErrInvalidScope, scope)
	}
	counts := make(map[RuleType]int, len(ruleTypes))
	for _, r := range s.Rules {
		if r.Scope != scope {
			continue
		}
		if _, ok := ruleTypes[r.Type]; !ok {
			return fmt.Errorf("%w: %q", ErrInvalidRuleType, r.Type)
		}
		counts[r.Type]++
	}
	for _, typ := range []RuleType{
		RuleTypePositionSizing,
		RuleTypeStopLoss,
		RuleTypeRegime,
		RuleTypeTakeProfit,
	} {
		switch n := counts[typ]; n {
		case 1:
		case 0:
			return fmt.Errorf("%w: scope %q missing rule type %q", ErrInvalidParams, scope, typ)
		default:
			return fmt.Errorf("%w: scope %q has %d rules of type %q", ErrInvalidParams, scope, n, typ)
		}
	}
	return nil
}
