package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const freezeBaseline = "docs/api-surface.txt"

// freezeSrc 是最小库包：一个结构体、一个构造函数、一个导出方法。
const freezeSrc = `package domain

// Quote 是最小报价。
type Quote struct {
	Symbol string
	Price  float64
}

// NewQuote 构造报价。
func NewQuote(symbol string, price float64) *Quote {
	return &Quote{Symbol: symbol, Price: price}
}

// Validate 校验报价。
func (q *Quote) Validate() error { return nil }
`

// newFreezeRepo 建夹具包并跑 -update 生成基线，返回仓库根。
// 用例切到临时目录里跑，root 用 "."：与 CLI 的默认用法一致
// （扫描器按 root 拼文件路径，root 只能相对仓库根为 "."）。
func newFreezeRepo(t *testing.T) string {
	t.Helper()
	abs := t.TempDir()
	t.Chdir(abs)
	writeRepoFile(t, "domain/quote.go", freezeSrc)
	if err := os.MkdirAll("docs", 0o755); err != nil {
		t.Fatal(err)
	}
	if err := run(options{root: ".", baseline: freezeBaseline, update: true, out: io.Discard}); err != nil {
		t.Fatalf("生成基线失败：%v", err)
	}
	return "."
}

func checkFreeze(root string) (string, error) {
	var buf bytes.Buffer
	err := run(options{root: root, baseline: freezeBaseline, out: &buf})
	return buf.String(), err
}

func TestFreezeUpdateThenDefaultIsGreen(t *testing.T) {
	root := newFreezeRepo(t)

	raw, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(freezeBaseline)))
	if err != nil {
		t.Fatalf("基线没落盘：%v", err)
	}
	for _, want := range []string{
		"domain/NewQuote(string, float64) *Quote",
		"domain/Quote struct{Symbol string; Price float64}",
		"domain/(*Quote).Validate() error",
	} {
		if !strings.Contains(string(raw), want) {
			t.Errorf("-update 生成的基线缺少 %q，got：\n%s", want, raw)
		}
	}

	out, err := checkFreeze(root)
	if err != nil {
		t.Fatalf("刚生成的基线再跑默认模式必须绿，got %v：\n%s", err, out)
	}
	if !strings.Contains(out, "✅") {
		t.Errorf("绿时缺少结论行，got：\n%s", out)
	}
}

type driftCase struct {
	name      string
	mutate    func(t *testing.T, root string)
	wantErr   string
	wantOut   []string
	rejectOut []string
}

func TestFreezeDriftLegs(t *testing.T) {
	cases := []driftCase{
		{
			name: "形状漂移/改函数签名",
			mutate: func(t *testing.T, root string) {
				writeRepoFile(t, filepath.Join(root, "domain", "quote.go"),
					strings.Replace(freezeSrc,
						"func NewQuote(symbol string, price float64) *Quote",
						"func NewQuote(symbol string) *Quote", 1))
			},
			wantErr: "导出面形状漂移：1 行消失 + 1 行新增",
			wantOut: []string{
				"   - domain/NewQuote(string, float64) *Quote",
				"   + domain/NewQuote(string) *Quote",
				"基线里有、代码里没有（1）",
				"代码里有、基线里没有（1）",
			},
		},
		{
			name: "形状漂移/改结构体字段类型",
			mutate: func(t *testing.T, root string) {
				writeRepoFile(t, filepath.Join(root, "domain", "quote.go"),
					strings.Replace(freezeSrc, "Price  float64", "Price  int64", 1))
			},
			wantErr: "导出面形状漂移：1 行消失 + 1 行新增",
			wantOut: []string{
				"   - domain/Quote struct{Symbol string; Price float64}",
				"   + domain/Quote struct{Symbol string; Price int64}",
			},
		},
		{
			name: "形状漂移/新增导出符号",
			mutate: func(t *testing.T, root string) {
				writeRepoFile(t, filepath.Join(root, "domain", "quote.go"),
					freezeSrc+"\n// Extra 是新导出符号。\nfunc Extra() {}\n")
			},
			wantErr:   "导出面形状漂移：0 行消失 + 1 行新增",
			wantOut:   []string{"   + domain/Extra()", "代码里有、基线里没有（1）"},
			rejectOut: []string{"   -"},
		},
		{
			name: "形状漂移/删导出方法",
			mutate: func(t *testing.T, root string) {
				writeRepoFile(t, filepath.Join(root, "domain", "quote.go"),
					strings.Replace(freezeSrc,
						"\n// Validate 校验报价。\nfunc (q *Quote) Validate() error { return nil }\n", "\n", 1))
			},
			wantErr:   "导出面形状漂移：1 行消失 + 0 行新增",
			wantOut:   []string{"   - domain/(*Quote).Validate() error", "基线里有、代码里没有（1）"},
			rejectOut: []string{"   +"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := newFreezeRepo(t)
			tc.mutate(t, root)

			out, err := checkFreeze(root)
			if err == nil {
				t.Fatalf("形状漂移必须红，got 绿：\n%s", out)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("错误文本缺少 %q，got %v：\n%s", tc.wantErr, err, out)
			}
			for _, want := range tc.wantOut {
				if !strings.Contains(out, want) {
					t.Errorf("输出没有报出 %q，got：\n%s", want, out)
				}
			}
			for _, bad := range tc.rejectOut {
				if strings.Contains(out, bad) {
					t.Errorf("这条漂移只该出现在一个方向，却出现了 %q 行：\n%s", bad, out)
				}
			}
		})
	}
}

func TestFreezeBaselineLegs(t *testing.T) {
	cases := []struct {
		name     string
		breakIt  func(t *testing.T, root string)
		wantErr  string
		wantSubs []string
	}{
		{
			name: "基线红/文件在但一行指纹都没有",
			breakIt: func(t *testing.T, root string) {
				writeRepoFile(t, filepath.Join(root, filepath.FromSlash(freezeBaseline)), "\n   \n")
			},
			wantErr:  "一行指纹都没有",
			wantSubs: []string{"空基线等于没冻结", "-update"},
		},
		{
			name: "基线红/文件缺失",
			breakIt: func(t *testing.T, root string) {
				if err := os.Remove(filepath.Join(root, filepath.FromSlash(freezeBaseline))); err != nil {
					t.Fatal(err)
				}
			},
			wantErr:  "读冻结基线",
			wantSubs: []string{"缺失即红", "-update"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := newFreezeRepo(t)
			tc.breakIt(t, root)

			out, err := checkFreeze(root)
			if err == nil {
				t.Fatalf("基线这条腿必须红，got 绿：\n%s", out)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("错误文本缺少 %q，got %v", tc.wantErr, err)
			}
			for _, want := range tc.wantSubs {
				if !strings.Contains(err.Error(), want) && !strings.Contains(out, want) {
					t.Errorf("错误/输出缺少 %q，got err=%v：\n%s", want, err, out)
				}
			}
			if strings.Contains(out, "❌") {
				t.Errorf("基线腿在读基线前返回，不该有 ❌ 输出：\n%s", out)
			}
		})
	}
}

func TestFreezeEmptyScanLegs(t *testing.T) {
	cases := []struct {
		name     string
		files    map[string]string
		wantErr  string
		wantSubs []string
	}{
		{
			name:     "空扫红/没有 Go 文件",
			files:    map[string]string{},
			wantErr:  "一个 Go 包都没扫到",
			wantSubs: []string{"护栏拒绝空绿"},
		},
		{
			name:     "空扫红/只有 package main",
			files:    map[string]string{"tools/x/main.go": "package main\n\nfunc main() {}\n"},
			wantErr:  "没有任何库包",
			wantSubs: []string{"全是 package main"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			t.Chdir(root)
			for rel, body := range tc.files {
				writeRepoFile(t, filepath.Join(root, filepath.FromSlash(rel)), body)
			}

			out, err := checkFreeze(".")
			if err == nil {
				t.Fatalf("空扫必须红（不许静默绿），got：\n%s", out)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("错误文本缺少 %q，got %v", tc.wantErr, err)
			}
			for _, want := range tc.wantSubs {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("错误缺少 %q，got %v", want, err)
				}
			}
		})
	}
}

func writeRepoFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}
