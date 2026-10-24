package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// 最小绿基座：四条契约根各有一个无害的包。任何红腿都在这上面只动一处。
const (
	okDomain       = "package domain\n\n// Price 是价格。\ntype Price struct{}\n"
	okEngine       = "package engine\n\n// Run 启动回测。\nfunc Run() {}\n"
	okFees         = "package fees\n\n// Fee 是费率。\ntype Fee struct{}\n"
	okDataProtocol = "package dataprotocol\n\n// Point 是统一数据点。\ntype Point struct{}\n"
)

func greenFiles() map[string]string {
	return map[string]string{
		"domain/base.go":       okDomain,
		"engine/base.go":       okEngine,
		"fees/base.go":         okFees,
		"dataprotocol/base.go": okDataProtocol,
	}
}

type shapeCase struct {
	name string
	// files 是夹具仓库的全部 .go 文件（slash 相对路径 → 内容）。
	files map[string]string
	// wantRed 表示这一条必须红；wantErr 是错误文本必须含的片段。
	wantRed    bool
	wantErr    string
	wantSubs   []string // 输出里必须出现的片段（咬住红因）
	rejectSubs []string // 输出里不得出现的片段（证明别的红腿没被触发）
	// redLines 是输出里 ❌ 行数；默认 1（只红一条腿）。
	redLines int
}

func TestContractShapeLegs(t *testing.T) {
	cases := []shapeCase{
		{
			name:  "绿/四条根都在且没有形状名",
			files: greenFiles(),
		},
		{
			name:     "命名红/domain 导出类型 *Config",
			files:    withFile(greenFiles(), "domain/config.go", "package domain\n\n// FooConfig 是命中形状名的类型。\ntype FooConfig struct{}\n"),
			wantRed:  true,
			wantErr:  "契约形状违规（1 项）",
			wantSubs: []string{"导出类型 FooConfig 命中 Config 族"},
			rejectSubs: []string{
				"禁止数据库/驱动依赖",
				"契约根目录扫不到",
				"一个 Go 包都没扫到",
			},
		},
		{
			name:     "命名红/engine 导出函数 *Options",
			files:    withFile(greenFiles(), "engine/options.go", "package engine\n\n// FooOptions 是命中形状名的函数。\nfunc FooOptions() {}\n"),
			wantRed:  true,
			wantErr:  "契约形状违规（1 项）",
			wantSubs: []string{"导出函数 FooOptions 命中 Options 族"},
			rejectSubs: []string{
				"禁止数据库/驱动依赖",
				"契约根目录扫不到",
				"一个 Go 包都没扫到",
			},
		},
		{
			name:     "命名红/fees 导出结构体 *Impl",
			files:    withFile(greenFiles(), "fees/impl.go", "package fees\n\n// FooImpl 是命中形状名的结构体。\ntype FooImpl struct{}\n"),
			wantRed:  true,
			wantErr:  "契约形状违规（1 项）",
			wantSubs: []string{"导出类型 FooImpl 命中 Impl 族"},
			rejectSubs: []string{
				"禁止数据库/驱动依赖",
				"契约根目录扫不到",
				"一个 Go 包都没扫到",
			},
		},
		{
			name:     "命名红/dataprotocol 导出结构体 *Options（新增契约根同样在红线内）",
			files:    withFile(greenFiles(), "dataprotocol/options.go", "package dataprotocol\n\n// PointOptions 命中形状名。\ntype PointOptions struct{}\n"),
			wantRed:  true,
			wantErr:  "契约形状违规（1 项）",
			wantSubs: []string{"导出类型 PointOptions 命中 Options 族"},
			rejectSubs: []string{
				"禁止数据库/驱动依赖",
				"契约根目录扫不到",
				"一个 Go 包都没扫到",
			},
		},
		{
			name: "白名单/7 个登记符号各就各位 → 绿",
			files: map[string]string{
				"domain/base.go": okDomain,
				"domain/whitelist.go": "package domain\n\n" +
					"type StrategyConfig struct{}\n" +
					"type ExecutionConfig struct{}\n" +
					"func DefaultExecutionConfig() ExecutionConfig { return ExecutionConfig{} }\n",
				"engine/whitelist.go": "package engine\n\n" +
					"type TradingConfig struct{}\n" +
					"type PriceLimitConfig struct{}\n" +
					"func DefaultTradingConfig() TradingConfig { return TradingConfig{} }\n" +
					"func NewMarketRouterFromConfig() {}\n",
				"fees/base.go":         okFees,
				"dataprotocol/base.go": okDataProtocol,
			},
		},
		{
			name:     "命名红/白名单键按目录精确匹配（TradingConfig 放进 domain 仍红）",
			files:    withFile(greenFiles(), "domain/drift.go", "package domain\n\n// TradingConfig 登记在 engine/ 下，放错目录不算例外。\ntype TradingConfig struct{}\n"),
			wantRed:  true,
			wantErr:  "契约形状违规（1 项）",
			wantSubs: []string{"导出类型 TradingConfig 命中 Config 族"},
			rejectSubs: []string{
				"禁止数据库/驱动依赖",
				"契约根目录扫不到",
			},
		},
		{
			name:     "DB 红/标准库 database/sql",
			files:    withFile(greenFiles(), "domain/db.go", "package domain\n\nimport _ \"database/sql\"\n"),
			wantRed:  true,
			wantErr:  "契约形状违规（1 项）",
			wantSubs: []string{`import "database/sql"`, "契约仓禁止数据库/驱动依赖"},
			rejectSubs: []string{
				"命中 Config 族",
				"命中 Options 族",
				"命中 Impl 族",
				"契约根目录扫不到",
				"一个 Go 包都没扫到",
			},
		},
		{
			name:     "DB 红/gorm.io/gorm",
			files:    withFile(greenFiles(), "fees/dao.go", "package fees\n\nimport _ \"gorm.io/gorm\"\n"),
			wantRed:  true,
			wantErr:  "契约形状违规（1 项）",
			wantSubs: []string{`import "gorm.io/gorm"`, "契约仓禁止数据库/驱动依赖"},
			rejectSubs: []string{
				"命中 Config 族",
				"命中 Options 族",
				"命中 Impl 族",
				"契约根目录扫不到",
				"一个 Go 包都没扫到",
			},
		},
		{
			name:     "空扫红/仓库里一个 Go 文件都没有",
			files:    map[string]string{},
			wantRed:  true,
			wantErr:  "一个 Go 包都没扫到",
			wantSubs: []string{"护栏拒绝空绿"},
			redLines: -1, // 这条腿在打印前就返回错误，输出里没有 ❌
		},
		{
			name:     "缺根红/只有 domain 与 engine，fees/ 与 dataprotocol/ 整棵搬走",
			files:    withFile(withFile(map[string]string{}, "domain/base.go", okDomain), "engine/base.go", okEngine),
			wantRed:  true,
			wantErr:  "契约形状违规（1 项）",
			wantSubs: []string{"契约根目录扫不到非测试 Go 文件：fees/、dataprotocol/"},
			rejectSubs: []string{
				"：domain/",
				"：engine/",
				"禁止数据库/驱动依赖",
				"一个 Go 包都没扫到",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			for rel, body := range tc.files {
				writeFixtureFile(t, filepath.Join(root, filepath.FromSlash(rel)), body)
			}

			var buf bytes.Buffer
			err := run(options{root: root, out: &buf})
			out := buf.String()

			if tc.wantRed {
				if err == nil {
					t.Fatalf("这条红腿必须报错，却绿了：\n%s", out)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("错误文本缺少 %q，got %v\n%s", tc.wantErr, err, out)
				}
			} else {
				if err != nil {
					t.Fatalf("这条腿必须绿，got %v：\n%s", err, out)
				}
				if !strings.Contains(out, "✅") {
					t.Errorf("绿时缺少结论行，got：\n%s", out)
				}
			}

			for _, want := range tc.wantSubs {
				// 空扫这类腿在打印前就返回，红因只出现在 error 里。
				if !strings.Contains(out, want) && !strings.Contains(err.Error(), want) {
					t.Errorf("输出/错误没有咬住红因 %q，got err=%v：\n%s", want, err, out)
				}
			}
			for _, bad := range tc.rejectSubs {
				if strings.Contains(out, bad) {
					t.Errorf("夹具只该触发一条腿，却出现了 %q：\n%s", bad, out)
				}
			}

			wantLines := 0
			if tc.wantRed {
				wantLines = 1
			}
			if tc.redLines < 0 {
				wantLines = 0
			}
			if got := strings.Count(out, "❌"); got != wantLines {
				t.Errorf("❌ 行数应为 %d（证明只红一条腿），got %d：\n%s", wantLines, got, out)
			}
		})
	}
}

func withFile(files map[string]string, rel, body string) map[string]string {
	out := make(map[string]string, len(files)+1)
	for k, v := range files {
		out[k] = v
	}
	out[rel] = body
	return out
}

func writeFixtureFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}
