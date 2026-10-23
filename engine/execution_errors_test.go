package contracts

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

// QT-10 ③-① 的钉子：三类哨兵必须互斥、可判、可包裹；`CodeDataGapSkipped`
// 走的是既有 invalid_reasons 通道，不是第四种状态。

func TestExecutionErrorSentinelsAreMutuallyExclusive(t *testing.T) {
	classes := []struct {
		name string
		err  error
	}{
		{"ErrDataGap", ErrDataGap},
		{"ErrStrategyFault", ErrStrategyFault},
		{"ErrContractBreach", ErrContractBreach},
	}
	for _, a := range classes {
		for _, b := range classes {
			if a.name == b.name {
				if !errors.Is(a.err, b.err) {
					t.Fatalf("%s 必须匹配自身", a.name)
				}
				continue
			}
			if errors.Is(a.err, b.err) {
				t.Fatalf("%s 不能匹配 %s —— 三类互斥是分级的前提，串了就等于没分级", a.name, b.name)
			}
		}
	}
}

// 包装保链（docs/API.md §3.2）：provider 侧一律是 `fmt.Errorf("%w: …", 类)`，
// 引擎侧在若干层之后才判类。任何一环用 %v 断链，调用方就只能退回字符串匹配。
func TestExecutionErrorClassesSurviveWrapping(t *testing.T) {
	cases := []struct {
		name  string
		build func() error
		want  error
	}{
		{
			name: "单层 Wrapf",
			build: func() error {
				return fmt.Errorf("%w: 2024-03-15 ohlcv missing", ErrDataGap)
			},
			want: ErrDataGap,
		},
		{
			name: "三层 Wrapf 后仍判得出",
			build: func() error {
				inner := fmt.Errorf("%w: wasm trap", ErrStrategyFault)
				mid := fmt.Errorf("on_bar session: %w", inner)
				return fmt.Errorf("getSignals: %w", mid)
			},
			want: ErrStrategyFault,
		},
		{
			name: "包装链命中兜底类",
			build: func() error {
				inner := fmt.Errorf("%w: unknown param %q", ErrContractBreach, "lookahead_window")
				return fmt.Errorf("provider: %w", inner)
			},
			want: ErrContractBreach,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.build()
			if !errors.Is(err, tc.want) {
				t.Fatalf("包装后 errors.Is 必须命中 %v，实际 err=%v", tc.want, err)
			}
			// 反证腿：命中自己那类，不能同时命中别类。
			for _, other := range []error{ErrDataGap, ErrStrategyFault, ErrContractBreach} {
				if other == tc.want {
					continue
				}
				if errors.Is(err, other) {
					t.Fatalf("%s 不该同时命中 %v（实际 err=%v）", tc.name, other, err)
				}
			}
		})
	}
}

// 判类只认哨兵身份，不认文案（docs/API.md §3.1）：措辞随时可改，
// 按 `err.Error()` 匹配的实现会在别人改注释文案那天静默失效。
func TestExecutionErrorClassesAreNotMatchedByMessage(t *testing.T) {
	lookalikes := map[string]error{
		"strategy fault":    ErrStrategyFault,
		"contract breach":   ErrContractBreach,
		"data gap: x":       ErrDataGap,
		"DataGap":           ErrDataGap,
		"abi: not this one": ErrContractBreach,
	}
	for text, want := range lookalikes {
		if errors.Is(errors.New(text), want) {
			t.Fatalf("同名异源的 error 不该命中哨兵 %v —— 命中说明判定退成了字符串匹配", want)
		}
	}
}

// Code 的字面量即契约：自动化层（loop 的失败计数、验证器维度、下游 SQL）按它判断。
func TestCodeDataGapSkippedLiteralIsPinned(t *testing.T) {
	if CodeDataGapSkipped != "data_gap_skipped" {
		t.Fatalf("Code 字面量变了 = 契约变更（走版本 + 同步 core SPEC）：实际 %q", CodeDataGapSkipped)
	}
}

// **禁止第四值**（QT-10 ④ / QT-63 ④）：本单只加 Code，冻结三值一字未动；
// 且新 Code 不能同时是一个 status —— 那会让「状态」与「原因」两条通道混起来。
func TestResultSpecStatusSetStaysThreeValued(t *testing.T) {
	statuses := map[string]bool{
		StatusCompleted: true,
		StatusInvalid:   true,
		StatusFailed:    true,
	}
	if len(statuses) != 3 {
		t.Fatalf("status 冻结面必须是三值，实测去重后 %d 值：%v", len(statuses), statuses)
	}
	if StatusCompleted != "completed" || StatusInvalid != "invalid" || StatusFailed != "failed" {
		t.Fatalf("三值取值是 frozen 面，一个字母都不许动：%q %q %q",
			StatusCompleted, StatusInvalid, StatusFailed)
	}
	if statuses[CodeDataGapSkipped] {
		t.Fatalf("%q 不能是第四种状态 —— 它只能进 invalid_reasons", CodeDataGapSkipped)
	}
}

// 新 Code 走的是既有通道：`CheckValidity` 只产出它自己那批「无效运行」Code，
// 数据缺口跳天由引擎在 `ErrDataGap` 分支上挂，不经过这条纯函数。
func TestCheckValidityNeverEmitsDataGapSkipped(t *testing.T) {
	res := validResult()
	res.UniverseMaxSize = intPtr(0)
	res.TotalTrades = 0
	res.StartDate = time.Time{}
	res.SharpeRatio = 1e13

	got := CheckValidity(&res)
	if len(got) == 0 {
		t.Fatal("空票池 + 0 成交 + 零值日期 + 超上限指标必须至少命中一条")
	}
	for _, r := range got {
		if r.Code == CodeDataGapSkipped {
			t.Fatalf("CheckValidity 不该产出 %q —— 那条 Code 的来由是 ErrDataGap 跳天", CodeDataGapSkipped)
		}
	}
}

// QT-41b ④：编译 / 资产管线的跨仓错误族。
//
// 四条是**不同的消费决策**，身份必须互不命中：编译失败 → 反馈给模型；
// 服务没配编译器 → 报配置错；源码超限 / 字节超限 → 不可重试，改输入或改上限。
// 再用一条包装腿证明跨层 `%w` 之后仍认得出来（消费方分属 agent / core /
// orchestration，不能要求它们 import 资产服务的实现包）。
func TestCompileFamilySentinelsStayDistinctAndSurviveWrapping(t *testing.T) {
	family := []struct {
		name string
		err  error
	}{
		{"ErrCompileFailed", ErrCompileFailed},
		{"ErrCompileUnavailable", ErrCompileUnavailable},
		{"ErrSourceTooLarge", ErrSourceTooLarge},
		{"ErrArtifactTooLarge", ErrArtifactTooLarge},
	}
	for _, a := range family {
		for _, b := range family {
			if a.name == b.name {
				if !errors.Is(a.err, b.err) {
					t.Fatalf("%s 必须匹配自身", a.name)
				}
				continue
			}
			if errors.Is(a.err, b.err) {
				t.Fatalf("%s 不能匹配 %s —— 四种消费决策不同，串了就分不出该重试还是该改输入", a.name, b.name)
			}
		}
		wrapped := fmt.Errorf("asset service: %w", a.err)
		if !errors.Is(wrapped, a.err) {
			t.Fatalf("%s 包装后必须仍可 errors.Is 判定", a.name)
		}
		// 反证：包装后不得命中同族别条。
		for _, b := range family {
			if b.name != a.name && errors.Is(wrapped, b.err) {
				t.Fatalf("%s 包装后误命中 %s", a.name, b.name)
			}
		}
	}
}
