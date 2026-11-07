package run_test

// CT-04 ⑤ 的跨仓对拍腿：三份形状的 **wire 字节**由生产侧（orchestration）的
// 结构体形状独立重建，与本包实现逐字节比对。
//
// 期望值是字面量，不是「两边都从本包取」——那会把两侧一起漂走变成假绿。
// 字形来源：
//   - claim        = tellurion-orchestration:internal/httpapi/workers.go 的 claimView
//   - provenance   = 同仓 internal/httpapi/provenance.go 的 provenanceView
//   - 父作业提交   = 同仓 internal/httpapi/parents.go 的 parentSubmitRequest 族
//   - 子 run 键    = 同仓 internal/parents/service.go 的 ChildIdempotencyKey

import (
	"encoding/json"
	"testing"
	"time"

	run "github.com/ruoxizhnya/tellurion-contracts/run"
)

func mustJSON(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	return string(b)
}

func TestClaimView_MatchesOrchestrationEncoding(t *testing.T) {
	lease := time.Date(2026, 11, 7, 8, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		in   run.ClaimView
		want string
	}{
		{
			name: "with parent job",
			in: run.ClaimView{
				RunID:           "run-01",
				RunSpec:         json.RawMessage(`{"strategy_ref":"momentum","start":"2026-01-02","end":"2026-03-31"}`),
				Attempts:        1,
				MaxAttempts:     3,
				LeaseUntil:      &lease,
				CancelRequested: false,
				ParentJobID:     "parent-abc",
			},
			want: `{"run_id":"run-01","run_spec":{"strategy_ref":"momentum","start":"2026-01-02","end":"2026-03-31"},"attempts":1,"max_attempts":3,"lease_until":"2026-11-07T08:00:00Z","cancel_requested":false,"parent_job_id":"parent-abc"}`,
		},
		{
			// 无父作业时 `parent_job_id` 必须**整个键不出现**（omitempty 冻结）。
			name: "without parent job",
			in: run.ClaimView{
				RunID:           "run-02",
				Attempts:        0,
				MaxAttempts:     5,
				CancelRequested: true,
			},
			want: `{"run_id":"run-02","run_spec":null,"attempts":0,"max_attempts":5,"lease_until":null,"cancel_requested":true}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := mustJSON(t, tc.in); got != tc.want {
				t.Fatalf("claim 字节与编排仓 claimView 不一致:\n got: %s\nwant: %s", got, tc.want)
			}
		})
	}
}

func TestProvenanceView_MatchesOrchestrationEncoding(t *testing.T) {
	created := time.Date(2026, 11, 7, 7, 45, 42, 753900000, time.UTC)
	uri, artifactHash, versionHash, status := "s3://artifacts/run-01.wasm", "artifact-hash", "version-hash", "completed"
	seed := int64(42)

	tests := []struct {
		name string
		in   run.ProvenanceView
		want string
	}{
		{
			name: "all fields set",
			in: run.ProvenanceView{
				RunID:               "run-01",
				RunSpecHash:         "run-spec-hash",
				DatasetSnapshotHash: "dataset-snapshot-hash",
				RiskRulesHash:       "risk-rules-hash",
				StrategyRef:         "momentum",
				ArtifactURI:         &uri,
				ArtifactHash:        &artifactHash,
				VersionHash:         &versionHash,
				ParamsHash:          "params-hash",
				Seed:                &seed,
				Kernel:              "go",
				EngineVersion:       "0.1.0",
				FillTiming:          "next_open",
				ResultStatus:        &status,
				CreatedAt:           &created,
			},
			want: `{"run_id":"run-01","run_spec_hash":"run-spec-hash","dataset_snapshot_hash":"dataset-snapshot-hash","risk_rules_hash":"risk-rules-hash","strategy_ref":"momentum","artifact_uri":"s3://artifacts/run-01.wasm","artifact_hash":"artifact-hash","version_hash":"version-hash","params_hash":"params-hash","seed":42,"kernel":"go","engine_version":"0.1.0","fill_timing":"next_open","result_status":"completed","created_at":"2026-11-07T07:45:42.7539Z"}`,
		},
		{
			// 可空字段出 `null` 而**不是**省略键；15 个键一个都不能少（冻结语义）。
			name: "nullable fields stay explicit",
			in:   run.ProvenanceView{RunID: "run-02"},
			want: `{"run_id":"run-02","run_spec_hash":"","dataset_snapshot_hash":"","risk_rules_hash":"","strategy_ref":"","artifact_uri":null,"artifact_hash":null,"version_hash":null,"params_hash":"","seed":null,"kernel":"","engine_version":"","fill_timing":"","result_status":null,"created_at":null}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := mustJSON(t, tc.in); got != tc.want {
				t.Fatalf("provenance 字节与编排仓 provenanceView 不一致:\n got: %s\nwant: %s", got, tc.want)
			}
		})
	}
}

func TestParentSubmitRequest_MatchesOrchestrationEncoding(t *testing.T) {
	tests := []struct {
		name string
		in   run.ParentSubmitRequest
		want string
	}{
		{
			name: "batch",
			in: run.ParentSubmitRequest{
				IdempotencyKey: "analysis-sensitivity:report-1",
				Kind:           "batch",
				Batch: &run.ParentBatchRequest{
					BaseRunSpec: json.RawMessage(`{"run_id":"run-01"}`),
					Grid: map[string][]json.RawMessage{
						"lookback": {json.RawMessage("20"), json.RawMessage("30")},
					},
				},
			},
			want: `{"idempotency_key":"analysis-sensitivity:report-1","kind":"batch","batch":{"base_run_spec":{"run_id":"run-01"},"grid":{"lookback":[20,30]}},"walkforward":null}`,
		},
		{
			name: "walkforward",
			in: run.ParentSubmitRequest{
				IdempotencyKey: "wf:1",
				Kind:           "walkforward",
				WalkForward: &run.ParentWalkForwardRequest{
					BaseRunSpec: json.RawMessage(`{"run_id":"run-02"}`),
					Windows:     run.ParentWindows{TrainDays: 60, TestDays: 20, StepDays: 20, MinTrainDays: 40},
					TradingDays: []string{"2026-01-02", "2026-01-05"},
				},
			},
			want: `{"idempotency_key":"wf:1","kind":"walkforward","batch":null,"walkforward":{"base_run_spec":{"run_id":"run-02"},"windows":{"train_days":60,"test_days":20,"step_days":20,"min_train_days":40},"trading_days":["2026-01-02","2026-01-05"]}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := mustJSON(t, tc.in); got != tc.want {
				t.Fatalf("父作业提交字节与编排仓 parentSubmitRequest 不一致:\n got: %s\nwant: %s", got, tc.want)
			}
		})
	}
}

func TestChildIdempotencyKey_MatchesOrchestrationRule(t *testing.T) {
	if run.ChildKeyPrefix != "parent" || run.ChildKeySeparator != ":" {
		t.Fatalf("保留前缀/分隔符被改动：prefix=%q sep=%q（改规则 = MAJOR）",
			run.ChildKeyPrefix, run.ChildKeySeparator)
	}

	for _, tc := range []struct {
		parent string
		index  int
		want   string
	}{
		{"parent-abc", 0, "parent:parent-abc:0"},
		{"parent-abc", 12, "parent:parent-abc:12"},
		{"3866fdc02b693c0a9a6b006669f7900c", 7, "parent:3866fdc02b693c0a9a6b006669f7900c:7"},
	} {
		if got := run.ChildIdempotencyKey(tc.parent, tc.index); got != tc.want {
			t.Errorf("ChildIdempotencyKey(%q, %d) = %q，期望 %q（与编排仓派生规则不一致）",
				tc.parent, tc.index, got, tc.want)
		}
	}
}
