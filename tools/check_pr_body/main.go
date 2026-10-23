// Command check_pr_body 核对 PR 正文是否填了契约仓模板的四个必填段落（QT-09 ④）。
//
// 它只在 `pull_request` 事件里有意义：CI 用
// `if: github.event_name == 'pull_request'` 控制，push / workflow_dispatch 显式跳过
// （GitHub 会把那一步标成 skipped，不是 passed —— 跳过不许装成通过）。
// 本地直接跑而没有正文时，本工具报错并说明原因，而不是静默退出 0。
//
// 四个红向，缺一条就是假护栏：
//
//	缺段红   —— 任一段落标题不存在（删掉模板段落来「让检查通过」= 红）
//	留白红   —— 段落标题在、正文只有标签/空行（占位符没替换 = 红）
//	类型红   —— 「变更类型与版本」段里没有 MINOR / PATCH / MAJOR(braking) 任一取值
//	登记红   —— 「登记区与冻结基线」段里没有同时写明 docs/API.md 与 docs/api-surface.txt
//
// 正文来源：--body-file 指定文件，或 stdin（CI 用 printf '%s' "$PR_BODY" | ...）。
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// requiredSections 是模板里必须保留的段落（用前缀匹配二级标题，避免因尾注改动而误红）。
var requiredSections = []struct {
	name    string
	heading string
}{
	{"使用方", "## 使用方"},
	{"变更类型与版本", "## 变更类型与版本"},
	{"登记区与冻结基线", "## 登记区与冻结基线"},
	{"Config / Options / Impl 与 DB import", "## Config / Options / Impl 与 DB import"},
}

var (
	// 多行 HTML 注释是模板里的填写说明，检查正文时先剥掉。
	commentRe = regexp.MustCompile(`(?s)<!--.*?-->`)
	// 模板里的占位符：`<填写>` / `<填写：不涉及>`。留着没替换即红。
	placeholderRe = regexp.MustCompile(`<\s*填写[^>]*>`)
	// 版本类型取值：大小写宽松，breaking 与 MAJOR 视为同义（README 版本规则用 breaking 描述 MAJOR 级变更）。
	versionTypeRe = regexp.MustCompile(`(?i)\b(MINOR|PATCH|MAJOR|breaking)\b`)
)

type options struct {
	bodyFile string
	out      *os.File
}

func main() {
	var opts options
	flag.StringVar(&opts.bodyFile, "body-file", "", "PR 正文文件；留空则读 stdin")
	flag.Parse()

	opts.out = os.Stdout
	if err := run(opts); err != nil {
		fmt.Fprintln(os.Stderr, "check_pr_body:", err)
		os.Exit(1)
	}
}

func run(opts options) error {
	body, err := readBody(opts.bodyFile)
	if err != nil {
		return err
	}
	if strings.TrimSpace(body) == "" {
		return fmt.Errorf("PR 正文为空：本工具只核对 pull_request 正文，请在 CI 的 pull_request 步骤里经 stdin 传入（其他事件应显式跳过，不要跑成假绿）")
	}

	sections, err := splitSections(body)
	if err != nil {
		return err
	}

	var red []string
	for _, want := range requiredSections {
		text, ok := findSection(sections, want.heading)
		if !ok {
			red = append(red, fmt.Sprintf("缺少必填段落 %q（模板段落不能删，没内容也要写「无」）", want.heading))
			continue
		}
		if placeholderRe.MatchString(text) {
			red = append(red, fmt.Sprintf("段落 %q 的占位符 <填写…> 没有替换", want.heading))
			continue
		}
		if !hasContent(text) {
			red = append(red, fmt.Sprintf("段落 %q 仍是空标签/空行，没有实际内容", want.heading))
			continue
		}
		switch want.name {
		case "变更类型与版本":
			if !versionTypeRe.MatchString(text) {
				red = append(red, "「变更类型与版本」段没有写明 MINOR / PATCH / MAJOR(breaking) 任一类型")
			}
		case "登记区与冻结基线":
			if !strings.Contains(text, "docs/API.md") || !strings.Contains(text, "docs/api-surface.txt") {
				red = append(red, "「登记区与冻结基线」段必须同时写明 docs/API.md 与 docs/api-surface.txt 的处理（无变化也要写「无导出面变化」并点名两个文件）")
			}
		}
	}

	fmt.Fprintf(opts.out, "PR 正文核对：必填段落 %d 个\n", len(requiredSections))
	if len(red) == 0 {
		fmt.Fprintln(opts.out, "✅ 四个必填段落都在且已填写")
		return nil
	}
	for _, r := range red {
		fmt.Fprintln(opts.out, "❌", r)
	}
	return fmt.Errorf("PR 正文不合模板（%d 项）", len(red))
}

func readBody(path string) (string, error) {
	if path != "" {
		raw, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("读 --body-file %s: %w", path, err)
		}
		return string(raw), nil
	}
	raw, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("读 stdin: %w", err)
	}
	return string(raw), nil
}

// splitSections 把正文按二级标题切开，返回 标题行 → 段内正文。
// 只有 `## ` 开头才算段落边界；`###` 与正文里的列表不破坏段落。
func splitSections(body string) (map[string]string, error) {
	sections := map[string]string{}
	var heading string
	var buf []string
	flush := func() {
		if heading != "" {
			sections[heading] = strings.Join(buf, "\n")
		}
		buf = nil
	}
	for _, line := range strings.Split(strings.ReplaceAll(body, "\r\n", "\n"), "\n") {
		if strings.HasPrefix(line, "## ") {
			flush()
			heading = strings.TrimSpace(line)
			continue
		}
		if heading != "" {
			buf = append(buf, line)
		}
	}
	flush()
	return sections, nil
}

// findSection 按标题前缀匹配段落：模板标题带「（必填）」尾注，用户也可能补说明，
// 所以只认 `## 使用方` 这个前缀，不要求整行相等。
func findSection(sections map[string]string, heading string) (string, bool) {
	for key, text := range sections {
		if key == heading || strings.HasPrefix(key, heading+"（") {
			return text, true
		}
	}
	return "", false
}

// hasContent 判断段落里是否有实内容：剥掉 HTML 注释后，跳过空行、
// 「`- 标签：`」这种只有标签没有值的行，以及只有破折号占位的行。
// 模板里的占位符是 `<填写…>` 形状，由 placeholderRe 单独判红。
func hasContent(text string) bool {
	for _, line := range strings.Split(commentRe.ReplaceAllString(text, ""), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasSuffix(line, "：") || strings.HasSuffix(line, ":") {
			continue // 标签后面没有值
		}
		if strings.HasSuffix(line, "——") {
			continue // 只有破折号，值没填
		}
		return true
	}
	return false
}
