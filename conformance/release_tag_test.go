package conformance

// release tag 必须带语料（C1 / QT-07 ⑤）。
//
// 这条守卫的存在理由是一个已经发生过的真实缺口：`v0.1.0` 打在语料落地之前，
// 所以 `git ls-tree -r v0.1.0` 里 conformance/ 命中 0 个文件，而 core
// （tellurion:pkg/conformance）与 agent（tellurion-agent:internal/conformance）
// 两侧的腿都**按 pin 的版本**从模块缓存里读同一份 cases.json，且约定「定位失败
// 一律 t.Fatalf，不许 skip」。⇒ 下游把这个 tag 升上去之前，任何一侧的 CI 都会红，
// 而红在别人的仓里；本仓当时四条护栏全绿。
//
// 判据取「最新的 release tag」而不是「每个 tag」：v0.1.0 缺语料是历史事实，
// 改写历史（重打 tag）比把它钉在账上更糟。往前看的那条不变量才是我们要的：
// 语料已经被三侧依赖，此后每一次发版都必须带上它。
//
// 失败一律 t.Fatalf，不 t.Skip：git 不可用、不在 checkout 里、一个 tag 都没有，
// 都可能是「语料今天没人管」的伪装，静默跳过就是假绿。

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// corpusTagPath 是语料在仓库内的路径（相对仓库根，即 git ls-tree 输出的形状）。
const corpusTagPath = "conformance/cases.json"

var releaseTagPattern = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)

// TestReleaseTagShipsCorpus 核对最新的 release tag 里确实有 cases.json。
func TestReleaseTagShipsCorpus(t *testing.T) {
	root := repoRootForTest(t)
	tags := versionTagsForTest(t, root)

	newest := newestReleaseTag(tags)
	if newest == "" {
		t.Fatalf("本仓没有任何 vX.Y.Z 形态的 release tag（现有：%v）—— 使用方只能 pin 伪版本或 main，"+
			"跨仓语料就无法随版本发布；先按 README「发版流程」打 tag", tags)
	}

	if !tagHasPath(t, root, newest, corpusTagPath) {
		t.Fatalf("最新 release tag %s 里没有 %s —— 下游 core/agent 的 conformance 腿按这个 tag 取语料必然 t.Fatalf。"+
			"修法：按 README「发版流程」补 CHANGELOG 节 + 兼容矩阵那一行，再打新 tag 并 push（不 push 的 tag 等于没发版）。",
			newest, corpusTagPath)
	}
	t.Logf("最新 release tag %s 含语料 %s", newest, corpusTagPath)
}

// newestReleaseTag 返回 semver 最大的 vX.Y.Z tag；没有任何合格 tag 时返回空串。
//
// 刻意不按字典序取「最后一个」：`v0.10.0` 字典序小于 `v0.1.9`，
// 拿字符串排序挑出来的"最新"会在 0.9 → 0.10 那天挑错对象，守卫从此核的是旧 tag。
func newestReleaseTag(tags []string) string {
	var best string
	var bestParts [3]int
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		parts, ok := parseSemver(tag)
		if !ok {
			continue
		}
		if best == "" || compareTriple(parts, bestParts) > 0 {
			best, bestParts = tag, parts
		}
	}
	return best
}

// parseSemver 解析 vX.Y.Z；任何其它形状（v0.1、v0.4.0-rc1、非 v 前缀）都 ok=false。
func parseSemver(v string) ([3]int, bool) {
	if !releaseTagPattern.MatchString(v) {
		return [3]int{}, false
	}
	var out [3]int
	for i, s := range strings.Split(strings.TrimPrefix(v, "v"), ".") {
		n, err := strconv.Atoi(s)
		if err != nil {
			return [3]int{}, false
		}
		out[i] = n
	}
	return out, true
}

func compareTriple(a, b [3]int) int {
	for i := range a {
		if a[i] != b[i] {
			return a[i] - b[i]
		}
	}
	return 0
}

// ─── git 取数 ────────────────────────────────────────────────

func repoRootForTest(t *testing.T) string {
	t.Helper()
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").CombinedOutput()
	if err != nil {
		t.Fatalf("git rev-parse --show-toplevel 失败（本守卫只在真实 checkout 里成立）：%v\n%s", err, out)
	}
	root := strings.TrimSpace(string(out))
	if root == "" {
		t.Fatal("git rev-parse --show-toplevel 返回空串")
	}
	return root
}

func versionTagsForTest(t *testing.T, root string) []string {
	t.Helper()
	out, err := exec.Command("git", "-C", root, "tag", "-l", "v*").CombinedOutput()
	if err != nil {
		t.Fatalf("git tag -l 失败：%v\n%s", err, out)
	}
	var tags []string
	for _, line := range strings.Split(string(out), "\n") {
		if s := strings.TrimSpace(line); s != "" {
			tags = append(tags, s)
		}
	}
	return tags
}

// tagHasPath 判断该 tag 的树里有没有指定路径。
//
// annotated tag 会被 git 自己剥到 commit，不需要手写 ^{}；但 --exit-code 那条腿
// 不能省：`git ls-tree` 对不存在的 pathspec 也是退出码 0，只能看输出是否为空。
func tagHasPath(t *testing.T, root, tag, path string) bool {
	t.Helper()
	out, err := exec.Command("git", "-C", root, "ls-tree", "-r", "--name-only", tag, "--", path).CombinedOutput()
	if err != nil {
		t.Fatalf("git ls-tree %s -- %s 失败：%v\n%s", tag, path, err, out)
	}
	return strings.TrimSpace(string(out)) != ""
}

// ─── 纯逻辑表驱动：tag 选取 ─────────────────────────────────

func TestNewestReleaseTag(t *testing.T) {
	tests := []struct {
		name string
		tags []string
		want string
	}{
		{"单个 tag", []string{"v0.1.0"}, "v0.1.0"},
		{"补丁版胜", []string{"v0.1.1", "v0.1.0"}, "v0.1.1"},
		{"次版本胜", []string{"v0.1.9", "v0.2.0"}, "v0.2.0"},
		{"字典序会挑错的两位数", []string{"v0.1.9", "v0.10.0"}, "v0.10.0"},
		{"主版本胜", []string{"v0.99.99", "v1.0.0"}, "v1.0.0"},
		{"乱序输入", []string{"v1.0.0", "v0.2.0", "v0.10.0", "v0.9.0"}, "v1.0.0"},
		{"忽略非完整版本 tag", []string{"v0.1", "release-3", "v0.3.1", "v0.2.9"}, "v0.3.1"},
		{"忽略预发布后缀", []string{"v0.4.0-rc1", "v0.3.0"}, "v0.3.0"},
		{"空白与空串", []string{"", "  ", "v0.1.2 "}, "v0.1.2"},
		{"全不合格返回空串", []string{"v0.1", "latest"}, ""},
		{"空列表返回空串", nil, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newestReleaseTag(tt.tags); got != tt.want {
				t.Errorf("newestReleaseTag(%q) = %q，期望 %q", tt.tags, got, tt.want)
			}
		})
	}
}

func TestParseSemverAndCompare(t *testing.T) {
	for _, tt := range []struct {
		in   string
		want [3]int
		ok   bool
	}{
		{"v1.23.4", [3]int{1, 23, 4}, true},
		{"v0.1.0", [3]int{0, 1, 0}, true},
		{"v0.1", [3]int{}, false},
		{"v0.4.0-rc1", [3]int{}, false},
		{"0.1.0", [3]int{}, false},
		{"main", [3]int{}, false},
	} {
		got, ok := parseSemver(tt.in)
		if ok != tt.ok || got != tt.want {
			t.Errorf("parseSemver(%q) = %v, %v；期望 %v, %v", tt.in, got, ok, tt.want, tt.ok)
		}
	}

	// 三档各自都要能分出胜负：只核主次、漏掉补丁，等于放行「同次版本的旧补丁」。
	for _, c := range []struct {
		a, b string
		want int
	}{
		{"v0.1.1", "v0.1.0", 1},
		{"v0.1.0", "v0.1.1", -1},
		{"v0.2.0", "v0.1.9", 1},
		{"v0.1.9", "v0.10.0", -1},
		{"v1.0.0", "v0.99.99", 1},
		{"v0.1.0", "v0.1.0", 0},
	} {
		a, _ := parseSemver(c.a)
		b, _ := parseSemver(c.b)
		if got := sign(compareTriple(a, b)); got != c.want {
			t.Errorf("compareTriple(%s, %s) 符号 = %d，期望 %d", c.a, c.b, got, c.want)
		}
	}
}

func sign(n int) int {
	switch {
	case n > 0:
		return 1
	case n < 0:
		return -1
	}
	return 0
}
