// Package run 承载编排（orchestration）与消费方之间**已在上 wire 的三份形状**：
// 领取到的 run（claim）、`run_artifacts` 读面（provenance）、子 run 入队请求。
//
// CT-04 的判据（`docs/TASKS.md`）只有一句：把已在多仓**手写复刻**的 wire 形状
// 收成契约里的唯一实现，让 `tellurion:docs/openapi.yaml` 里那句
// 「mirror tellurion-orchestration's provenanceView exactly … must not drift」
// 从口头承诺变成机械判据（本包的字节级向量 + conformance 语料）。
//
// # 冻结语义
//
// 字段名 / 可空性 / 声明顺序三件事**逐字对齐生产侧**：
//
//   - `ClaimView` ↔ `tellurion-orchestration:internal/httpapi/workers.go` 的 `claimView`；
//   - `ProvenanceView` ↔ 同仓 `internal/httpapi/provenance.go` 的 `provenanceView`，
//     可空字段出 `null` 而**不是省略键**（读侧不用猜「没有」与「空串」的区别）；
//   - `ParentSubmitRequest` 族 ↔ 同仓 `internal/httpapi/parents.go` 的提交请求形状。
//
// 改字段名 / 类型 / 可空性 / 顺序 = breaking；`ChildIdempotencyKey` 的规则改动
// 同样是 MAJOR（它决定子 run 的幂等身份）。
//
// # 不迁（刻意留在编排仓）
//
// `RunRecord` / `ParentRecord` / `WorkerRecord`（队列与租约状态机）、
// `jobView` / `parentView` / `parentChildView` / `parentCountsView`（编排自己的读面投影）、
// `ArtifactStore`（带写语义与身份漂移拒绝）、`run_artifacts` 的 DDL。
// 契约只收 DTO，不做 decode / 默认值填充 / 幂等判定。
//
// # 零 I/O
//
// 包内无构造器、无时钟、无数据库；`ChildIdempotencyKey` 是纯函数（ADR-0004）。
package run

import (
	"encoding/json"
	"strconv"
	"time"
)

// ClaimView 是 worker 领取 run 的响应体（`claimView` 的契约镜像）。
//
// `RunSpec` 是完整 RunSpec 的原样字节，worker 需要它才能执行，所以不复用
// `GET /api/runs/jobs/{id}`（那里不返回请求体）。
//
// `ParentJobID` 带 `omitempty`：普通 run 没有父作业，键**不出现**；
// 去掉它会改变已上 wire 的字节（CT-04 ⑤ 的破坏验证之一）。
type ClaimView struct {
	RunID           string          `json:"run_id"`
	RunSpec         json.RawMessage `json:"run_spec"`
	Attempts        int             `json:"attempts"`
	MaxAttempts     int             `json:"max_attempts"`
	LeaseUntil      *time.Time      `json:"lease_until"`
	CancelRequested bool            `json:"cancel_requested"`
	ParentJobID     string          `json:"parent_job_id,omitempty"`
}

// ProvenanceView 是编排仓 `GET /api/runs/{id}/provenance` 的 15 字段响应体
// （`provenanceView` 的契约镜像；core openapi 的 `BacktestProvenance` 描述同一个概念）。
//
// `artifact_uri` / `artifact_hash` / `version_hash` 可空：内置策略没有独立
// artifact；`seed` 也可空。三者与 `result_status` / `created_at` 一样，
// **显式出 `null` 而不是省略键**（冻结语义，禁止加 `omitempty`）。
// `risk_rules_hash` 是空串时表示**未登记**，不是「无规则」。
type ProvenanceView struct {
	RunID               string     `json:"run_id"`
	RunSpecHash         string     `json:"run_spec_hash"`
	DatasetSnapshotHash string     `json:"dataset_snapshot_hash"`
	RiskRulesHash       string     `json:"risk_rules_hash"`
	StrategyRef         string     `json:"strategy_ref"`
	ArtifactURI         *string    `json:"artifact_uri"`
	ArtifactHash        *string    `json:"artifact_hash"`
	VersionHash         *string    `json:"version_hash"`
	ParamsHash          string     `json:"params_hash"`
	Seed                *int64     `json:"seed"`
	Kernel              string     `json:"kernel"`
	EngineVersion       string     `json:"engine_version"`
	FillTiming          string     `json:"fill_timing"`
	ResultStatus        *string    `json:"result_status"`
	CreatedAt           *time.Time `json:"created_at"`
}

// ParentSubmitRequest 是 `POST /api/runs/parents` 的请求体（batch / walkforward 共用）。
//
// `Kind` 决定二者取一：同时给出两个载荷是歧义请求，服务侧直接 400
// （不猜哪一个才是意图）。
type ParentSubmitRequest struct {
	IdempotencyKey string                    `json:"idempotency_key"`
	Kind           string                    `json:"kind"`
	Batch          *ParentBatchRequest       `json:"batch"`
	WalkForward    *ParentWalkForwardRequest `json:"walkforward"`
}

// ParentBatchRequest 是 batch 父作业的载荷。
type ParentBatchRequest struct {
	BaseRunSpec json.RawMessage              `json:"base_run_spec"`
	Grid        map[string][]json.RawMessage `json:"grid"`
}

// ParentWalkForwardRequest 是 walkforward 父作业的载荷。
//
// `TradingDays` 是调用方传入的交易日历（`YYYY-MM-DD`，严格递增）；
// 编排仓不自造日历。
type ParentWalkForwardRequest struct {
	BaseRunSpec json.RawMessage `json:"base_run_spec"`
	Windows     ParentWindows   `json:"windows"`
	TradingDays []string        `json:"trading_days"`
}

// ParentWindows 是 walkforward 的窗口切分参数。
type ParentWindows struct {
	TrainDays    int `json:"train_days"`
	TestDays     int `json:"test_days"`
	StepDays     int `json:"step_days"`
	MinTrainDays int `json:"min_train_days"`
}

const (
	// ChildKeyPrefix 是子 run 幂等键的保留前缀。调用方自己的幂等键不应使用它，
	// 否则可能与某个子 run 撞键。
	ChildKeyPrefix = "parent"
	// ChildKeySeparator 是键的分隔符。改动 = 全部子 run 的幂等身份改变（MAJOR）。
	ChildKeySeparator = ":"
)

// ChildIdempotencyKey 是子 run 幂等键的派生规则（stable；改规则 = MAJOR）。
//
// 形状：`parent:<parent_job_id>:<child_index>`。父 id 是计划身份（plan_hash）
// 在库里的行 id——同一份计划（即使换幂等键提交）拿到同一个父 id，因此重放
// 派生出同一批键；不同计划必然不同父 id，不会串键。
//
// 本函数逐字节对齐 `tellurion-orchestration:internal/parents/service.go` 的
// `ChildIdempotencyKey`；两侧必须始终产出同一串（CT-04 的对拍腿与 conformance
// 语料各钉一次）。
func ChildIdempotencyKey(parentID string, childIndex int) string {
	return ChildKeyPrefix + ChildKeySeparator + parentID + ChildKeySeparator + strconv.Itoa(childIndex)
}
