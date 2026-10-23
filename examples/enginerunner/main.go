// Command enginerunner 是 `contracts.EngineRunner` 的最小用法示例（QT-09 ②）。
//
// 它只演示调用形状：一个恒成功、不做任何计算的 Runner 桩，加上请求/响应 DTO 的最小构造。
// 真正的回测实现在 core，下游只需要知道「跑一次回测」的入参与判据。
// 对应契约：engine/ 的 EngineRunner（docs/API.md §4 `engine/` → EngineRunner / EngineRunner.RunBacktest，§5.2）
package main

import (
	"context"
	"fmt"

	contracts "github.com/ruoxizhnya/tellurion-contracts/engine"
)

// stub 是编译期可用的空实现：不算业务，只把请求原样回显成 completed。
type stub struct{}

func (stub) RunBacktest(ctx context.Context, req contracts.BacktestRequest) (*contracts.BacktestResponse, error) {
	return &contracts.BacktestResponse{
		Status:   contracts.StatusCompleted,
		Strategy: req.Strategy,
	}, nil
}

// 编译期断言：桩实现了端口。签名变了这里先红。
var _ contracts.EngineRunner = stub{}

func main() {
	var runner contracts.EngineRunner = stub{}
	req := contracts.BacktestRequest{
		Strategy:  "demo",
		StartDate: "2026-01-01",
		EndDate:   "2026-06-30",
	}
	resp, err := runner.RunBacktest(context.Background(), req)
	if err != nil {
		fmt.Printf("engine.EngineRunner 示例：RunBacktest 失败：%v\n", err)
		return
	}
	// status 是唯一成功判据（docs/API.md §3 第 5 条）。
	fmt.Printf("engine.EngineRunner 示例：RunBacktest 返回 status=%s，strategy=%s（桩不做业务）\n", resp.Status, resp.Strategy)
}
