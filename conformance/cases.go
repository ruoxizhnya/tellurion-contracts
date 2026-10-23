// Package conformance 承载跨仓契约一致性语料的形状与加载器（QT-09 ③）。
//
// cases.json 是唯一真相：同一组用例分别喂给契约仓、core 与 agent 的实现，
// 任一侧语义漂移即红。本文件只做「读 + 校验」，不含任何断言逻辑 ——
// 每侧的断言住在该侧自己的测试里（契约仓是 conformance_test.go，core 是
// tellurion:pkg/conformance，agent 是 tellurion-agent:internal/conformance）。
//
// ⚠️ 本包**刻意不导出任何符号**：
//   - 契约仓的导出面被 docs/API.md（check_api_surface）与 docs/api-surface.txt
//     （check_surface_freeze）双冻结，新增导出符号必须同批改文档并重冻基线；
//     本单的边界是「不改任何文档」，导出 `Case`/`Load` 会当场把两道护栏变红。
//   - 下游两侧根本不需要 import 这里的加载器：它们各自读同一份 cases.json，
//     这正是「同一组用例、两套实现」的形态 —— 共享加载器反而会让两侧共用
//     一份解析逻辑，把「解析漂移」也一起藏起来。
//
// 因此加载函数是 loadCases，用例结构是 caseSpec。语料的字段形状由本文件与
// 两侧测试共同钉住；任何一侧读不出某个字段都必须 fail-loud。
package conformance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// casesFileName 是语料文件名（相对语料目录）。
const casesFileName = "cases.json"

// knownSides 是语料里允许声明的实现侧，与三仓一一对应。
//
// 未登记的侧名一律报错而不是忽略：一个拼错的 "cores" 会让那条用例在 core
// 侧悄悄不跑，而错误恰恰是「悄悄不跑」这类问题里最难发现的一种。
var knownSides = map[string]bool{
	"contracts": true,
	"core":      true,
	"agent":     true,
}

// caseSpec 是 cases.json 里的一条用例。
type caseSpec struct {
	// ID 是全局唯一用例名，失败信息里必须能唯一定位到它。
	ID string `json:"id"`
	// Kind 决定由哪一侧的哪个断言函数执行。
	Kind string `json:"kind"`
	// Sides 显式声明必须跑这条用例的侧。写进 sides 的侧没跑 = 红。
	Sides []string `json:"sides"`
	// Input 是该 kind 的入参（各侧按同名字段自行解码）。
	Input json.RawMessage `json:"input"`
	// Expect 是期望结果（各侧按同名字段自行解码）。
	Expect json.RawMessage `json:"expect"`
}

// loadCases 读取并校验语料目录里的 cases.json。
//
// 所有失败都返回错误（fail-loud），不返回空切片让调用方「跑 0 条然后绿」：
// 语料为空、id 缺失/重复、kind 缺失、sides 为空或含未知侧、input/expect
// 缺失，全部是语料本身的缺陷，必须在读入阶段就暴露。
func loadCases(dir string) ([]caseSpec, error) {
	path := filepath.Join(dir, casesFileName)
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("conformance: 读取语料 %s: %w", path, err)
	}

	var cases []caseSpec
	if err := json.Unmarshal(raw, &cases); err != nil {
		return nil, fmt.Errorf("conformance: 解析语料 %s: %w", path, err)
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("conformance: 语料 %s 为空（空语料 = 假绿）", path)
	}

	seen := make(map[string]bool, len(cases))
	for i, c := range cases {
		if c.ID == "" {
			return nil, fmt.Errorf("conformance: 语料 %s 第 %d 条缺少 id", path, i)
		}
		if seen[c.ID] {
			return nil, fmt.Errorf("conformance: 语料 %s 里 id %q 重复（失败信息必须能唯一定位）", path, c.ID)
		}
		seen[c.ID] = true

		if c.Kind == "" {
			return nil, fmt.Errorf("conformance: 用例 %q 缺少 kind", c.ID)
		}
		if len(c.Sides) == 0 {
			return nil, fmt.Errorf("conformance: 用例 %q 的 sides 为空（不再明侧的用例永远不会被执行）", c.ID)
		}
		for _, s := range c.Sides {
			if !knownSides[s] {
				return nil, fmt.Errorf("conformance: 用例 %q 声明了未知侧 %q（已知：agent/contracts/core）", c.ID, s)
			}
		}
		if len(c.Input) == 0 {
			return nil, fmt.Errorf("conformance: 用例 %q 缺少 input", c.ID)
		}
		if len(c.Expect) == 0 {
			return nil, fmt.Errorf("conformance: 用例 %q 缺少 expect", c.ID)
		}
	}

	return cases, nil
}
