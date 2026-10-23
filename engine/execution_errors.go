package contracts

import "errors"

// QT-10 ③-①（N1）：策略执行错误的三类哨兵。
//
// 三类存在的理由是一个具体的假象：在它们落地之前，引擎对 `getSignals` 的错误
// 只 `Warn` + 跳过当天（`tellurion:pkg/backtest/engine.go`），于是「策略每 50 天
// trap 一次」与「那天真的缺数据」在曲线上长得一模一样 —— 观测是假的，曲线却
// 照样漂亮。分级之后只有数据缺口能跳天，且必须留痕；另外两类终止 run。
//
// 落点在本仓而不是 core：分类的**生产者**（provider / wasm 宿主）与**消费者**
// （引擎、job、AI 管线）分属不同模块，哨兵必须在两侧之间可 import 且不引入环。
// 本包是叶子包，与 `ErrFillTiming*`、abi 族同处 `engine/`——`errors/` 不独立成包
// 的裁决见 docs/API.md §8 #1（adr/0003）。
//
// 用法约束（docs/API.md §3）：判定只用 `errors.Is`，`err.Error()` 的措辞不属于
// 契约；跨层包装只用 `%w`，`%v` 断链等于把调用方逼回字符串匹配。三类互斥：一个
// 错误只归一类，归类看**这条错误的来由**，不看它有多严重。
var (
	// ErrDataGap 是「这次运行需要的输入取不到」——外部数据缺口，既不是策略写坏了，
	// 也不是调用方配错了什么。
	//
	// 消费行为：**跳过当天** + 往 `invalid_reasons` 写 `CodeDataGapSkipped`，run
	// 继续跑完；最终 `status=invalid`（被跳过过交易日的运行不是一次干净的观测）。
	// 它是三类里唯一不终止 run 的；把它当成可重试的成功同样错。
	ErrDataGap = errors.New("data gap: required input is unavailable")

	// ErrStrategyFault 是「策略自己跑坏了」：wasm trap、非法信号、预算超限、
	// artifact 校验失败。
	//
	// 消费行为：**终止 run** → `status=failed` + `error` 必填。禁止「跳过这一天
	// 继续」——那正是本单要堵的洞：策略一直在坏、曲线一直在涨。
	ErrStrategyFault = errors.New("strategy fault")

	// ErrContractBreach 是「跨边界的约定被违反」：前视信号（`as_of` 晚于查询的
	// asOf）、未知参数名、abi 不符。
	//
	// 消费行为：**终止 run** → `status=failed` + `error` 必填。它同时是「无法
	// 归类」的兜底类（docs/API.md §3.4）：分不清是数据缺口还是策略故障时按本类
	// 处理——宁可停，不要跑出漂亮但错的曲线。
	ErrContractBreach = errors.New("contract breach")
)
