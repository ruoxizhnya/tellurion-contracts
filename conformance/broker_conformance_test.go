package conformance

// CT-05 ⑤③：broker wire 语料的自检腿（契约侧）。
//
// 与 run_cases.json 同因单独成文件：core / agent 的加载器只认
// contracts/core/agent 三侧，把新侧塞进 cases.json 会让下游升 pin 当场红。
// 期望字节全部来自 `tellurion-broker:internal/gateway` 的生产类型序列化。

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/ruoxizhnya/tellurion-contracts/broker"
)

type brokerCaseSpec struct {
	ID     string          `json:"id"`
	Kind   string          `json:"kind"`
	Sides  []string        `json:"sides"`
	Input  json.RawMessage `json:"input"`
	Expect json.RawMessage `json:"expect"`
}

func TestBrokerCases_ContractSide(t *testing.T) {
	raw, err := os.ReadFile("broker_cases.json")
	if err != nil {
		t.Fatalf("读取 broker_cases.json: %v", err)
	}
	var cases []brokerCaseSpec
	if err := json.Unmarshal(raw, &cases); err != nil {
		t.Fatalf("解析 broker_cases.json: %v", err)
	}
	if len(cases) == 0 {
		t.Fatal("broker_cases.json 为空（空语料 = 假绿）")
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
			case "order_request_json":
				var v broker.OrderRequest
				mustDecodeStrict(t, c.Input, &v)
				assertJSON(t, v, c.Expect)

			case "order_ack_json":
				var v broker.OrderAck
				mustDecodeStrict(t, c.Input, &v)
				assertJSON(t, v, c.Expect)

			case "order_state_json":
				var v broker.OrderState
				mustDecodeStrict(t, c.Input, &v)
				assertJSON(t, v, c.Expect)

			case "fill_json":
				var v broker.Fill
				mustDecodeStrict(t, c.Input, &v)
				assertJSON(t, v, c.Expect)

			case "account_snapshot_json":
				var v broker.AccountSnapshot
				mustDecodeStrict(t, c.Input, &v)
				assertJSON(t, v, c.Expect)

			case "emergency_flatten_json":
				var v broker.EmergencyFlattenResult
				mustDecodeStrict(t, c.Input, &v)
				assertJSON(t, v, c.Expect)

			default:
				t.Fatalf("未知 kind %q（新增 kind 必须同时加断言）", c.Kind)
			}
		})
	}
}
