// Package risk 承载风控规则快照的跨仓唯一形状与 hash 纯函数。
//
// 三条边界（CT-03）：
//
//   - **只落形状与编码**：`Snapshot` / `SnapshotRule` / `Scope` / `RuleType` 的字段集、
//     冻结词表、4 个哨兵、`MaxParamsBytes` 与 `Snapshot.Hash()`。`Bundle` / `ParseBundle` /
//     `ForScope` 刻意**不迁**——它们带 decode、默认值填充与进程内索引，属服务策略。
//   - **零 I/O、零时钟、零环境依赖**：包内不读时间、不连库、不读环境变量（ADR-0004）。
//     `CreatedAt` 只是被承载的字段，hash 不覆盖它。
//   - **编码逐字节搬运**：`Snapshot.Hash()` 必须与 `tellurion-risk:internal/rule/snapshot.go`
//     的 `SnapshotHash` 产出同一个 hex。**禁止**改用 `engine.NewCanonical`：两套编码器的帧
//     结构、域标签写法与参数规范化都不同（见 `docs/API.md` §5.6），换编码 = 全历史
//     `risk_rules_hash` 作废，且 `quant.run_provenance` 与 `run.run_artifacts` 两处锚点同时失效。
package risk

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// DomainSnapshotHash 是 risk_rules_hash 的域标签，**逐字节冻结**。
//
// 改它 = 全历史 risk_rules_hash 作废（`tellurion-risk:docs/db-schema.md` §5、
// `tellurion-risk:docs/api-schema.md` §4）。要改必须走 ADR 级裁决并升域标签到 v2。
const DomainSnapshotHash = "tellurion.risk_snapshot_hash.v1"

// MaxParamsBytes 是规则参数 JSON 的体积上限；超限 fail-loud，不截断。
//
// ⚠️ 与 `engine/abi.go` 的 `MaxParamsBytes`（4096，wasm initialize 的 params 区上限）
// **同名不同值不同义**，别为了「统一」改成 4096——那会拒掉今天所有合法快照。
const MaxParamsBytes = 64 << 10

// RuleType 是冻结的规则类型枚举。旧读者遇未知类型必须 fail-closed 拒收。
// 新增类型 = MINOR；改既有取值的语义 = breaking。
type RuleType string

const (
	// RuleTypePositionSizing 仓位规模规则。
	RuleTypePositionSizing RuleType = "position_sizing"
	// RuleTypeStopLoss 止损规则。
	RuleTypeStopLoss RuleType = "stop_loss"
	// RuleTypeRegime 市场状态规则。
	RuleTypeRegime RuleType = "regime"
	// RuleTypeTakeProfit 止盈规则。
	RuleTypeTakeProfit RuleType = "take_profit"
)

var ruleTypes = map[RuleType]struct{}{
	RuleTypePositionSizing: {},
	RuleTypeStopLoss:       {},
	RuleTypeRegime:         {},
	RuleTypeTakeProfit:     {},
}

// Scope 是规则的适用范围。找不到对应 scope 的规则 = 没有规则，
// **不许**回落到默认限额。
//
// ⚠️ 与 `engine/abi.go` 的 `ScopeSymbol`（`StrategyScope` 的取值 `"symbol"`，指
// 「每实例绑定一个 symbol」）**不是同一个概念**：这里指「规则适用范围」。
// 两者禁止用类型别名互转。
type Scope string

const (
	// ScopeGlobal 全局规则。
	ScopeGlobal Scope = "global"
	// ScopeAccount 账户级规则。
	ScopeAccount Scope = "account"
	// ScopeStrategy 策略级规则。
	ScopeStrategy Scope = "strategy"
	// ScopeSymbol 标的级规则。
	ScopeSymbol Scope = "symbol"
)

var scopes = map[Scope]struct{}{
	ScopeGlobal:   {},
	ScopeAccount:  {},
	ScopeStrategy: {},
	ScopeSymbol:   {},
}

// Valid 报告 scope 是否在冻结枚举里（纯函数，ADR-0004）。
func (s Scope) Valid() bool {
	_, ok := scopes[s]
	return ok
}

// SnapshotRule 是快照里的一条生效规则（active revision 的内容身份）。
// draft 永不进快照：能到达 runtime 的只有 active_revision。
type SnapshotRule struct {
	RuleID      string          `json:"rule_id"`
	Revision    int             `json:"revision"`
	Scope       Scope           `json:"scope"`
	Type        RuleType        `json:"type"`
	Params      json.RawMessage `json:"params"`
	VersionHash string          `json:"version_hash"`
}

// Snapshot 是内容寻址的规则快照：risk_rules_hash 是回测结果里的可复现锚点。
//
// 不变量：同一 hash 的规则包字节一致；内容变 ⇒ 必然是新 hash，
// 原地改内容是绝对禁止项（写入与读取都重新 derive 校验）。
type Snapshot struct {
	RiskRulesHash string         `json:"risk_rules_hash"`
	Rules         []SnapshotRule `json:"rules"`
	CreatedAt     time.Time      `json:"created_at"`
}

// Clone 返回深拷贝（Rules 与每个 Params 都拷），避免调用方与来源共享可变内存。
func (s Snapshot) Clone() Snapshot {
	out := s
	out.Rules = make([]SnapshotRule, len(s.Rules))
	for i, r := range s.Rules {
		r.Params = append(json.RawMessage(nil), r.Params...)
		out.Rules[i] = r
	}
	return out
}

// Hash 由快照内容派生 risk_rules_hash：按 rule_id 升序后逐条写进长度前缀域。
//
// **不含** CreatedAt —— 同一内容永远同一 hash，快照缓存按 hash 永不失效。
// **不读** s.RiskRulesHash —— 声称值只参与 Verify 的比对，不参与派生。
func (s Snapshot) Hash() (string, error) {
	return snapshotHash(s.Rules)
}

// Verify 重新 derive 并与声称的 RiskRulesHash 比对。
// 内容寻址被破坏即 fail-loud 返回 ErrSnapshotMismatch，不吞、不回退。
func (s Snapshot) Verify() error {
	derived, err := s.Hash()
	if err != nil {
		return err
	}
	if s.RiskRulesHash != derived {
		return fmt.Errorf("%w: derived %s, claimed %q", ErrSnapshotMismatch, derived, s.RiskRulesHash)
	}
	return nil
}

// snapshotHash 是 risk `SnapshotHash` 的逐字节搬运（编码细节见包注释）。
func snapshotHash(rules []SnapshotRule) (string, error) {
	if len(rules) == 0 {
		return "", fmt.Errorf("%w: snapshot must contain at least one active rule", ErrInvalidParams)
	}
	sorted := make([]SnapshotRule, len(rules))
	copy(sorted, rules)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].RuleID < sorted[j].RuleID })
	for i := 1; i < len(sorted); i++ {
		if sorted[i].RuleID == sorted[i-1].RuleID {
			return "", fmt.Errorf("%w: duplicate rule_id %q", ErrInvalidParams, sorted[i].RuleID)
		}
	}

	var b strings.Builder
	b.WriteString(DomainSnapshotHash)
	b.WriteByte('\n')
	writeField(&b, "rule_count", strconv.Itoa(len(sorted)))
	for _, r := range sorted {
		if r.RuleID == "" {
			return "", fmt.Errorf("%w: rule_id is required", ErrInvalidParams)
		}
		if r.Revision <= 0 {
			return "", fmt.Errorf("%w: %s revision must be > 0", ErrInvalidParams, r.RuleID)
		}
		if _, ok := ruleTypes[r.Type]; !ok {
			// 旧读者遇未知类型必须 fail-closed 拒收。
			return "", fmt.Errorf("%w: %q", ErrInvalidRuleType, r.Type)
		}
		if _, ok := scopes[r.Scope]; !ok {
			return "", fmt.Errorf("%w: %q", ErrInvalidScope, r.Scope)
		}
		if r.VersionHash == "" {
			return "", fmt.Errorf("%w: %s version_hash is required", ErrInvalidParams, r.RuleID)
		}
		params, err := canonicalParams(r.Params)
		if err != nil {
			return "", err
		}
		writeField(&b, "rule_id", r.RuleID)
		writeField(&b, "revision", strconv.Itoa(r.Revision))
		writeField(&b, "scope", string(r.Scope))
		writeField(&b, "type", string(r.Type))
		writeField(&b, "params", string(params))
		writeField(&b, "version_hash", r.VersionHash)
	}
	sum := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(sum[:]), nil
}

// writeField 写一帧「长度前缀」文本：`len(key):key=len(value):value\n`。
// 长度前缀避免拼接歧义（空串与分隔符不会产生同一 canonical 串）。
func writeField(b *strings.Builder, key, value string) {
	fmt.Fprintf(b, "%d:%s=%d:%s\n", len(key), key, len(value), value)
}

// canonicalParams 把参数规范化成「键有序、无多余空白」的 JSON 对象字节。
// 规范化是 hash 可复现的前提：同一对象的不同书写必须同 hash。
//
// 刻意用 `UseNumber` + `json.Marshal`（与 risk 一致）：`engine` 的 canonical
// 显式拒绝 json.Number，换成它会让浮点参数的字节表示改变。
func canonicalParams(raw json.RawMessage) ([]byte, error) {
	if len(bytes.TrimSpace(raw)) == 0 {
		return nil, fmt.Errorf("%w: params is required", ErrInvalidParams)
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, fmt.Errorf("%w: params is not valid JSON: %v", ErrInvalidParams, err)
	}
	if dec.More() {
		return nil, fmt.Errorf("%w: params has trailing content", ErrInvalidParams)
	}
	if _, ok := v.(map[string]any); !ok {
		return nil, fmt.Errorf("%w: params must be a JSON object", ErrInvalidParams)
	}
	out, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("%w: params cannot be canonicalised: %v", ErrInvalidParams, err)
	}
	if len(out) > MaxParamsBytes {
		return nil, fmt.Errorf("%w: params exceed %d bytes", ErrInvalidParams, MaxParamsBytes)
	}
	return out, nil
}
