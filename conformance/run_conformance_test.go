package conformance

// CT-04 ⑤：编排 wire 语料的自检腿（契约侧）。
//
// 语料与 cases.json 同形，但**单独一份文件**：core / agent 的加载器只认
// contracts/core/agent 三侧，把 orchestration / analytics 侧名塞进 cases.json
// 会让下游升级 pin 到本版本时当场红。下游各仓的腿是各自仓的单（CT-04 ⑥）。

import (
	"encoding/json"
	"os"
	"testing"

	run "github.com/ruoxizhnya/tellurion-contracts/run"
)

type runCaseSpec struct {
	ID     string          `json:"id"`
	Kind   string          `json:"kind"`
	Sides  []string        `json:"sides"`
	Input  json.RawMessage `json:"input"`
	Expect json.RawMessage `json:"expect"`
}

func mustDecodeStrict(t *testing.T, raw json.RawMessage, dst any) {
	t.Helper()
	if err := decodeStrict(raw, dst); err != nil {
		t.Fatalf("解码用例输入失败: %v", err)
	}
}

func TestRunCases_ContractSide(t *testing.T) {
	raw, err := os.ReadFile("run_cases.json")
	if err != nil {
		t.Fatalf("读取 run_cases.json: %v", err)
	}
	var cases []runCaseSpec
	if err := json.Unmarshal(raw, &cases); err != nil {
		t.Fatalf("解析 run_cases.json: %v", err)
	}
	if len(cases) == 0 {
		t.Fatal("run_cases.json 为空（空语料 = 假绿）")
	}

	for _, c := range cases {
		if len(c.Sides) == 0 {
			t.Fatalf("用例 %q 的 sides 为空", c.ID)
		}
		for _, side := range c.Sides {
			if side != "contracts" {
				t.Fatalf("用例 %q 声明了本语料不执行的侧 %q", c.ID, side)
			}
		}

		t.Run(c.ID, func(t *testing.T) {
			switch c.Kind {
			case "claim_view_json":
				var v run.ClaimView
				mustDecodeStrict(t, c.Input, &v)
				assertJSON(t, v, c.Expect)

			case "provenance_view_json":
				var v run.ProvenanceView
				mustDecodeStrict(t, c.Input, &v)
				assertJSON(t, v, c.Expect)

			case "child_idempotency_key":
				var in struct {
					ParentJobID string `json:"parent_job_id"`
					ChildIndex  int    `json:"child_index"`
				}
				mustDecodeStrict(t, c.Input, &in)
				var exp struct {
					Key string `json:"key"`
				}
				if err := json.Unmarshal(c.Expect, &exp); err != nil {
					t.Fatalf("解析 expect: %v", err)
				}
				if got := run.ChildIdempotencyKey(in.ParentJobID, in.ChildIndex); got != exp.Key {
					t.Fatalf("ChildIdempotencyKey(%q, %d) = %q，期望 %q",
						in.ParentJobID, in.ChildIndex, got, exp.Key)
				}

			default:
				t.Fatalf("未知 kind %q（新增 kind 必须同时加断言）", c.Kind)
			}
		})
	}
}

func assertJSON(t *testing.T, v any, expect json.RawMessage) {
	t.Helper()
	var exp struct {
		JSON string `json:"json"`
	}
	if err := json.Unmarshal(expect, &exp); err != nil {
		t.Fatalf("解析 expect: %v", err)
	}
	got, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	if string(got) != exp.JSON {
		t.Fatalf("wire 字节与语料不一致:\n got: %s\nwant: %s", got, exp.JSON)
	}
}
