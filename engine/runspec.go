package contracts

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"
	"time"
)

// K9（RunSpec 冻结）：内核路由字段（Kernel / EngineParams）的默认语义与
// fail-closed 校验。落点在 contracts 的理由与 validity.go 相同 —— 叶子包，
// API 层（cmd/analysis）与编排层（job / batch / walkforward）都能引用，
// 不引入 import 环，校验口径只有一处。

// DefaultKernel 是未指定 Kernel 时的默认内核（FR-1.1：现有行为不变）。
const DefaultKernel = "go"

// engineParamsByKernel 登记每个内核「认识」的引擎参数名。
//
// **未登记 = 该内核不接受任何引擎特有参数**（fail-closed 的另一面）：
// 请求里出现未知名一律拒绝，绝不静默忽略（FR-4.2）。目前 Go 内核没有引擎
// 特有参数，空集是刻意的 —— 新内核（K13 nautilus 等）接入时在自己的 runner
// 落地任务里登记自己的参数名，别提前为想象中的参数开门。
var engineParamsByKernel = map[string][]string{
	DefaultKernel: {},
}

// KnownEngineParams 返回该内核已登记的引擎参数名（排序后的副本）。
// 未登记的内核返回空切片 —— 与「不接受任何引擎特有参数」同义。
func KnownEngineParams(kernel string) []string {
	names := append([]string(nil), engineParamsByKernel[kernel]...)
	sort.Strings(names)
	return names
}

// EffectiveKernel 返回本次运行实际生效的内核名：显式指定则原样透传
// （路由与未知内核裁定是 K10 的事），空值回落 DefaultKernel。
//
// 默认值**不**写回 r.Kernel：BacktestRequest 会被序列化（job 落库等），
// 回填会让存量请求的 JSON 多出 `"kernel":"go"`，破坏「存量请求序列化行为
// 不变」的契约腿。
func (r BacktestRequest) EffectiveKernel() string {
	if r.Kernel == "" {
		return DefaultKernel
	}
	return r.Kernel
}

// ValidateRunSpec 校验 RunSpec 的内核路由字段，返回非 nil 即请求非法
// （API 层映射 400）。当前拒绝项：engine_params 里的未知名、signals 行的
// 非法值（K15）、显式声明但不认识的 fill_timing（QT-02）、strategy_ref 的
// 非 `id` / `id@hash` 形状与 strategy_params 的嵌套键（QT-20）。
//
// 为什么是「拒绝而非忽略」：静默吞掉一个拼错的参数名，调用方会以为参数
// 生效了并拿这次结果下结论 —— 与 FR-1.3「不得静默回退 go 内核」同一条
// fail-closed 原则，只是发生在参数维度。signals 同理：一行 weight=1.5 的
// 脏数据如果被钳到 1.0，双跑对齐报告里就多了一个「查无实据」的差异。
//
// fill_timing 只拒**显式声明且不认识**的值：空值走 EffectiveFillTiming 的
// 默认语义（QT-02 ④：报告归档前默认口径不动，存量请求必须继续可用），
// 而引擎日循环入口那道 ResolveFillTiming 才是「必填」的执行点 —— 它管的
// 是所有入口（HTTP / job / batch / walkforward / 直接调 Engine）。
//
// 边界（QT-20 · J2）：本函数只做**形状 / 类型**校验。strategy_params 的
// 未知键（清单来自策略 manifest）、默认值合并、参数取值范围都属于执行侧
// 的业务规则，不在这里，也不在本包复制一份清单。
func (r BacktestRequest) ValidateRunSpec() error {
	if err := r.validateSignals(); err != nil {
		return err
	}
	if err := r.validateStrategyRef(); err != nil {
		return err
	}
	if err := r.validateStrategyParams(); err != nil {
		return err
	}
	if r.FillTiming != "" && !r.EffectiveFillTiming().Valid() {
		return fmt.Errorf("%w (got %q for kernel %q)", ErrFillTimingUnknown,
			r.FillTiming, r.EffectiveKernel())
	}

	if len(r.EngineParams) == 0 {
		return nil
	}

	known := KnownEngineParams(r.EffectiveKernel())
	knownSet := make(map[string]struct{}, len(known))
	for _, k := range known {
		knownSet[k] = struct{}{}
	}

	var unknown []string
	for k := range r.EngineParams {
		if _, ok := knownSet[k]; !ok {
			unknown = append(unknown, k)
		}
	}
	if len(unknown) == 0 {
		return nil
	}
	sort.Strings(unknown)

	knownDesc := "(none)"
	if len(known) > 0 {
		knownDesc = strings.Join(known, ", ")
	}
	return fmt.Errorf(
		"engine_params contains unknown key(s) [%s] for kernel %q: unknown engine params are rejected, not ignored (known: %s)",
		strings.Join(unknown, ", "), r.EffectiveKernel(), knownDesc)
}

// validateSignals 校验可选的目标权重信号流（K15）：symbol 非空、weight 为
// [0,1] 内的有限小数、as_of 为 ISO 日期。任一违例整条请求拒绝（400），
// 不钳制、不跳过 —— 与侧车 runspec._parse_signals 逐字同口径，两侧校验
// 不一致会让同一个请求「Go 收、NT 拒」或反之，双跑从入口就不可比。
func (r BacktestRequest) validateSignals() error {
	if len(r.Signals) == 0 {
		return nil
	}
	for i, sig := range r.Signals {
		if strings.TrimSpace(sig.Symbol) == "" {
			return fmt.Errorf("signals[%d].symbol must be a non-blank string", i)
		}
		if math.IsNaN(sig.Weight) || math.IsInf(sig.Weight, 0) || sig.Weight < 0 || sig.Weight > 1 {
			return fmt.Errorf(
				"signals[%d].weight must be a finite fraction in [0, 1], got %v", i, sig.Weight)
		}
		if _, err := time.Parse("2006-01-02", sig.AsOf); err != nil {
			return fmt.Errorf(
				"signals[%d].as_of must be an ISO date (YYYY-MM-DD), got %q", i, sig.AsOf)
		}
	}
	return nil
}

// validateStrategyRef 校验可选的策略引用串（QT-20）：形状必须能被
// ParseStrategyRef 解析 —— `id` 或 `id@64位小写hex`，不做宽松解析、不认
// 第二种引用形状。空值 = 未指定（走既有 Strategy 注册表键），不拒。
//
// 引用**指向什么**不在这里判：artifact_hash 与 version_hash 的形态不可自证，
// 解析成哪个是资产服务查表的事（provenance.go 的 N7-3 注释）。
func (r BacktestRequest) validateStrategyRef() error {
	if r.StrategyRef == "" {
		return nil
	}
	if _, err := ParseStrategyRef(r.StrategyRef); err != nil {
		return fmt.Errorf("strategy_ref: %w", err)
	}
	return nil
}

// validateStrategyParams 校验参数覆写的**形状**（QT-20）：只接受浅层对象，
// 值的类型必须是标量（string / number / bool / null）。
//
// 拒嵌套不是洁癖：浅合并与 params_hash 的编码都要求「一个键一个值」，
// 嵌套会逼每个消费方自己定义遍历顺序，同参复现立刻失去唯一含义。
//
// 未知键**不在这里拒**——清单来自策略 manifest（J2），契约层若复制一份
// 就成了第二个真相源，manifest 加参数时契约会先红。
func (r BacktestRequest) validateStrategyParams() error {
	if len(r.StrategyParams) == 0 {
		return nil
	}
	keys := make([]string, 0, len(r.StrategyParams))
	for k := range r.StrategyParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if strings.TrimSpace(k) == "" {
			return fmt.Errorf("strategy_params key must be a non-blank string")
		}
		v := r.StrategyParams[k]
		if v == nil {
			continue
		}
		if !isScalarValue(v) {
			return fmt.Errorf(
				"strategy_params[%q] has type %T: only scalar values are allowed (nested objects/arrays are rejected)",
				k, v)
		}
	}
	return nil
}

// isScalarValue 报告 v 是否为可进 strategy_params 的标量。
// 数值囊括全部整数 / 浮点宽度：HTTP JSON 只产出 float64，进程内调用方
// 可能直接给 int，两者都应被接受；nil 表示「该键无值」，不算嵌套。
func isScalarValue(v any) bool {
	if v == nil {
		return true
	}
	switch reflect.ValueOf(v).Kind() {
	case reflect.Bool, reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}
