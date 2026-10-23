package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const validBody = `# 契约变更说明

## 使用方（必填）
下游 tellurion-engine 通过 RunSpec 调用，无新增依赖。

## 变更类型与版本（必填）
PATCH：仅措辞订正。

## 登记区与冻结基线（必填）
无导出面变化；docs/API.md 与 docs/api-surface.txt 均无需改动。

## Config / Options / Impl 与 DB import（必填）
无新增形状名，无 DB import。
`

type bodyCase struct {
	name        string
	body        string
	missingFile bool
	wantRed     bool
	wantErr     string
	wantSubs    []string
	redLines    int // 输出里的 ❌ 行数；默认 1，读正文前就失败的腿为 0
}

func TestPRBodyLegs(t *testing.T) {
	mustReplace := func(old, new string) string {
		t.Helper()
		got := strings.Replace(validBody, old, new, 1)
		if got == validBody {
			t.Fatalf("夹具改不动：%q 不在正文里", old)
		}
		return got
	}

	cases := []bodyCase{
		{
			name: "绿/四段全填",
			body: validBody,
		},
		{
			name: "类型绿/写 MINOR",
			body: mustReplace("PATCH：仅措辞订正。", "MINOR：新增一个导出符号。"),
		},
		{
			name: "类型绿/写 breaking",
			body: mustReplace("PATCH：仅措辞订正。", "breaking：删掉一个导出方法。"),
		},
		{
			name:     "缺段红/使用方",
			body:     mustReplace("## 使用方（必填）\n下游 tellurion-engine 通过 RunSpec 调用，无新增依赖。\n", ""),
			wantRed:  true,
			wantErr:  "PR 正文不合模板（1 项）",
			wantSubs: []string{`缺少必填段落 "## 使用方"`},
		},
		{
			name:     "缺段红/变更类型与版本",
			body:     mustReplace("## 变更类型与版本（必填）\nPATCH：仅措辞订正。\n", ""),
			wantRed:  true,
			wantErr:  "PR 正文不合模板（1 项）",
			wantSubs: []string{`缺少必填段落 "## 变更类型与版本"`},
		},
		{
			name:     "缺段红/登记区与冻结基线",
			body:     mustReplace("## 登记区与冻结基线（必填）\n无导出面变化；docs/API.md 与 docs/api-surface.txt 均无需改动。\n", ""),
			wantRed:  true,
			wantErr:  "PR 正文不合模板（1 项）",
			wantSubs: []string{`缺少必填段落 "## 登记区与冻结基线"`},
		},
		{
			name:     "缺段红/Config Options Impl 与 DB import",
			body:     mustReplace("## Config / Options / Impl 与 DB import（必填）\n无新增形状名，无 DB import。\n", ""),
			wantRed:  true,
			wantErr:  "PR 正文不合模板（1 项）",
			wantSubs: []string{`缺少必填段落 "## Config / Options / Impl 与 DB import"`},
		},
		{
			name: "占位符红/段落留着 <填写…> 没替换",
			body: mustReplace(
				"下游 tellurion-engine 通过 RunSpec 调用，无新增依赖。",
				"<!-- 模板说明：写下游包名 -->\n<填写：下游包名>"),
			wantRed:  true,
			wantErr:  "PR 正文不合模板（1 项）",
			wantSubs: []string{`段落 "## 使用方" 的占位符 <填写…> 没有替换`},
		},
		{
			name: "留白红/段落只剩标签与注释",
			body: mustReplace(
				"下游 tellurion-engine 通过 RunSpec 调用，无新增依赖。",
				"<!-- 模板说明：写下游包名 -->\n- 使用方："),
			wantRed:  true,
			wantErr:  "PR 正文不合模板（1 项）",
			wantSubs: []string{`段落 "## 使用方" 仍是空标签/空行`},
		},
		{
			name:     "类型红/没写 MINOR PATCH breaking",
			body:     mustReplace("PATCH：仅措辞订正。", "本次只调整段落措辞，无行为变化。"),
			wantRed:  true,
			wantErr:  "PR 正文不合模板（1 项）",
			wantSubs: []string{"「变更类型与版本」段没有写明 MINOR / PATCH / MAJOR(breaking) 任一类型"},
		},
		{
			name: "登记红/只点名 docs/API.md",
			body: mustReplace(
				"无导出面变化；docs/API.md 与 docs/api-surface.txt 均无需改动。",
				"无导出面变化；已更新 docs/API.md。"),
			wantRed:  true,
			wantErr:  "PR 正文不合模板（1 项）",
			wantSubs: []string{"必须同时写明 docs/API.md 与 docs/api-surface.txt"},
		},
		{
			name: "登记红/两个登记文件都没点名",
			body: mustReplace(
				"无导出面变化；docs/API.md 与 docs/api-surface.txt 均无需改动。",
				"无导出面变化，未涉及登记文件。"),
			wantRed:  true,
			wantErr:  "PR 正文不合模板（1 项）",
			wantSubs: []string{"必须同时写明 docs/API.md 与 docs/api-surface.txt"},
		},
		{
			name:     "空正文红/只有空白",
			body:     "   \n\n",
			wantRed:  true,
			wantErr:  "PR 正文为空",
			wantSubs: []string{"只核对 pull_request 正文"},
			redLines: -1,
		},
		{
			name:        "读正文红/--body-file 不存在",
			missingFile: true,
			wantRed:     true,
			wantErr:     "读 --body-file",
			wantSubs:    []string{"no-such-pr-body.md"},
			redLines:    -1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "PR_BODY.md")
			if tc.missingFile {
				path = filepath.Join(t.TempDir(), "no-such-pr-body.md")
			} else if err := os.WriteFile(path, []byte(tc.body), 0o600); err != nil {
				t.Fatal(err)
			}

			var buf bytes.Buffer
			err := run(options{bodyFile: path, out: &buf})
			out := buf.String()

			if tc.wantRed {
				if err == nil {
					t.Fatalf("这条红腿必须报错，却绿了：\n%s", out)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("错误文本缺少 %q，got %v：\n%s", tc.wantErr, err, out)
				}
			} else {
				if err != nil {
					t.Fatalf("这条腿必须绿，got %v：\n%s", err, out)
				}
				if !strings.Contains(out, "✅ 四个必填段落都在且已填写") {
					t.Errorf("绿时缺少结论行，got：\n%s", out)
				}
			}

			for _, want := range tc.wantSubs {
				if !strings.Contains(out, want) && !strings.Contains(err.Error(), want) {
					t.Errorf("输出/错误没有咬住红因 %q，got err=%v：\n%s", want, err, out)
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
