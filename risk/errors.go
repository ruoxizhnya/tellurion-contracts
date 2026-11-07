package risk

import "errors"

// 跨仓哨兵：只有这四条能 `errors.Is` 判定，文案与 `tellurion-risk:internal/rule/errors.go`
// **逐字一致**（改文案 = 改下游依赖的判定文本）。
//
// 其余七个哨兵（ErrRuleNotFound / ErrRevisionNotFound / ErrRevisionConflict /
// ErrRevisionNotDraft / ErrIdempotencyConflict / ErrSnapshotNotFound /
// ErrNoApprovedRules）是服务私有的：跨仓只经 wire code 消费，**不进契约**。
var (
	// ErrInvalidRuleType 规则类型不在冻结枚举。
	ErrInvalidRuleType = errors.New("rule type is not in the frozen enum")
	// ErrInvalidScope 适用范围不在冻结枚举。
	ErrInvalidScope = errors.New("scope is not in the frozen enum")
	// ErrInvalidParams 参数不是合法 JSON 对象、为空或超限。
	ErrInvalidParams = errors.New("rule params are invalid")
	// ErrSnapshotMismatch 内容寻址被破坏：derive hash ≠ 声称 hash。fail-loud，不吞。
	ErrSnapshotMismatch = errors.New("snapshot hash does not match content")
)
