package main

import (
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// fixtureRoot 是 testdata/ok：导出面与登记表逐符号对上的最小样本仓。
const fixtureRoot = "testdata/ok"

func TestFixtureIsRegisteredExactly(t *testing.T) {
	if err := checkFixture(fixtureRoot); err != nil {
		t.Fatalf("一致的夹具必须绿，got %v", err)
	}
}

// TestMainAndTestFilesAreNotSurface 钉住「不参与核对」那条：
// 夹具里的 cmd/tool/main.go 与 engine/face_test.go 都有导出符号，
// 一旦被扫进来，第一条用例就会因为「未登记」而红。
func TestMainAndTestFilesAreNotSurface(t *testing.T) {
	code, err := scanExports(fixtureRoot)
	if err != nil {
		t.Fatal(err)
	}
	got := flatten(code)
	for _, sym := range got {
		if strings.Contains(sym, "ExportedFromMainIsNotSurface") || strings.Contains(sym, "ExportedFromTestFileIsNotSurface") {
			t.Errorf("package main / _test.go 的符号不该进导出面：%s", sym)
		}
	}
	if syms := code["cmd/tool"]; len(syms) != 0 {
		t.Errorf("package main 目录应扫出 0 个符号，got %v", syms)
	}
	// 反面腿：夹具里确实有该被扫到的符号，否则上面的缺席断言是空的。
	for _, want := range []string{"engine/New", "engine/EngineRunner.Close", "domain/OHLCV"} {
		if !hasSym(got, want) {
			t.Errorf("夹具的导出面缺少 %s（扫描本身失效了）", want)
		}
	}
	// 未导出符号、未导出类型的方法必须缺席。
	for _, unwanted := range []string{"engine/statusInternal", "engine/Runner.run", "engine/internal.DoIt"} {
		if hasSym(got, unwanted) {
			t.Errorf("%s 不属于导出面，却被扫出来了", unwanted)
		}
	}
}

func TestUnregisteredExportedSymbolIsRed(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, filepath.Join(root, "engine", "extra.go"),
		"package contracts\n\n// Scratch 是「新增导出符号没登记」的形状。\nfunc Scratch() error { return nil }\n")

	unregistered, stale := runCompare(t, root)
	if len(stale) != 0 {
		t.Fatalf("不该有失效项，got %v", stale)
	}
	if want := []string{"engine/Scratch"}; !reflect.DeepEqual(unregistered, want) {
		t.Fatalf("未登记项应为 %v，got %v", want, unregistered)
	}
	if err := checkFixture(root); err == nil {
		t.Fatal("新增未登记导出符号时 run 必须返回错误（CI 才会红）")
	}
}

// TestUnregisteredInterfaceMethodIsRed 钉住端口形状：给接口加一个方法，
// 类型名没变、方法集变了，同样必须红。
func TestUnregisteredInterfaceMethodIsRed(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, filepath.Join(root, "engine", "port.go"),
		"package contracts\n\ntype Port interface {\n\tDial() error\n}\n")

	unregistered, _ := runCompare(t, root)
	for _, want := range []string{"engine/Port", "engine/Port.Dial"} {
		if !hasSym(unregistered, want) {
			t.Errorf("未登记项缺少 %s，got %v", want, unregistered)
		}
	}
}

// TestStaleRegistrationIsRed 是另一半：表里有、代码里没有。
// 只查「缺登记」的护栏会让文档长期挂着不存在的符号，下游照样踩。
func TestStaleRegistrationIsRed(t *testing.T) {
	root := copyFixture(t)
	api := filepath.Join(root, "API.md")
	raw, err := os.ReadFile(api)
	if err != nil {
		t.Fatal(err)
	}
	hook := "| `Ghost` | type | stable | 代码里不存在 |"
	withGhost := strings.Replace(string(raw), surfaceEnd, hook+"\n"+surfaceEnd, 1)
	if withGhost == string(raw) {
		t.Fatal("夹具 API.md 没有登记区结束标记，改不动")
	}
	writeFile(t, api, withGhost)

	unregistered, stale := runCompare(t, root)
	if len(unregistered) != 0 {
		t.Fatalf("不该有未登记项，got %v", unregistered)
	}
	if want := []string{"domain/Ghost"}; !reflect.DeepEqual(stale, want) {
		t.Fatalf("失效项应为 %v，got %v", want, stale)
	}
	if err := checkFixture(root); err == nil {
		t.Fatal("登记表挂着不存在的符号时 run 必须返回错误")
	}
}

// TestEmptyRegisterIsNotAVacuousPass 防最省事的那种假护栏：把登记区的行删空，
// 检查器就「什么都不缺」。代码有符号 + 表为空 ⇒ 必须全量红。
func TestEmptyRegisterIsNotAVacuousPass(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, filepath.Join(root, "API.md"), "# 空登记表\n\n"+surfaceBegin+"\n\n"+surfaceEnd+"\n")

	unregistered, _ := runCompare(t, root)
	if len(unregistered) < 10 {
		t.Fatalf("空登记表必须把每个导出符号都报成未登记，got %d 项：%v", len(unregistered), unregistered)
	}
	for _, want := range []string{"engine/EngineRunner", "engine/StatusCompleted", "domain/OHLCV"} {
		if !hasSym(unregistered, want) {
			t.Errorf("缺少 %s", want)
		}
	}
}

// TestMissingMarkersIsFatal 防「删掉整段登记区让护栏空过」。
func TestMissingMarkersIsFatal(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, filepath.Join(root, "API.md"), "# 没有核对区的接口文档\n")

	err := checkFixture(root)
	if err == nil {
		t.Fatal("缺少 surface:begin/end 标记必须报错，不能当作「零符号 = 通过」")
	}
	if !strings.Contains(err.Error(), surfaceBegin) {
		t.Errorf("错误信息要点名缺哪个标记，got %v", err)
	}
}

// TestSyntaxErrorIsFatal 防「解析失败的包凭空消失、护栏静默通过」。
func TestSyntaxErrorIsFatal(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, filepath.Join(root, "engine", "broken.go"), "package contracts\n\nfunc Broken( {\n")

	if err := checkFixture(root); err == nil {
		t.Fatal("语法错误必须当场红")
	}
}

// TestTwoPackagesInOneDirIsFatal 防「一个目录两个包」时只扫到其中一个包，
// 另一个包的符号整体消失。
func TestTwoPackagesInOneDirIsFatal(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, filepath.Join(root, "engine", "other.go"), "package other\n\nfunc AlsoExported() {}\n")

	if err := checkFixture(root); err == nil {
		t.Fatal("同一目录多包必须报错")
	}
}

// TestHeadingWithoutPackageDirRegistersNothing 钉住登记表形状：
// 小节标题里没有 `某目录/` 时，那张表不算登记（写错标题不能悄悄漏登记）。
func TestHeadingWithoutPackageDirRegistersNothing(t *testing.T) {
	doc, err := parseRegister(filepath.Join(fixtureRoot, "API.md"))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := doc["NotRegistered"]; ok {
		t.Error("标记区之外的表行不该被登记")
	}
	if got := len(doc); got != 3 {
		t.Errorf("夹具登记了 3 个包（根包 + engine + domain），got %d（%v）", got, doc)
	}
}

// TestRootPackageIsRegistrable 钉住护栏对根包也有效：小节标题写 `.` 就能登记，
// 去掉那一节它立刻变成「未登记」。不认根包的检查器会让根包一新增导出符号就永久红
// ——那是不可满足的护栏，和假护栏一样糟。
func TestRootPackageIsRegistrable(t *testing.T) {
	doc, err := parseRegister(filepath.Join(fixtureRoot, "API.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !hasSym(doc["."], "RootValue") {
		t.Fatalf("`.` 小节没有登记 RootValue，got %v", doc)
	}

	root := copyFixture(t)
	api := filepath.Join(root, "API.md")
	raw, err := os.ReadFile(api)
	if err != nil {
		t.Fatal(err)
	}
	dotted := strings.ReplaceAll(string(raw), "### `.` — 根包 `contracts`", "### 根包（标题里没有反引号目录）")
	if dotted == string(raw) {
		t.Fatal("夹具 API.md 没有根包小节，改不动")
	}
	writeFile(t, api, dotted)

	unregistered, stale := runCompare(t, root)
	if want := []string{"./RootValue"}; !reflect.DeepEqual(unregistered, want) {
		t.Fatalf("根包符号应报成未登记 %v，got %v", want, unregistered)
	}
	// 标题不带目录 ⇒ 那一节的表行根本不算登记（不是「登记了但失效」），
	// 写错标题的人会被「未登记」正面挡住，而不是悄悄登记到一个不存在的包上。
	if len(stale) != 0 {
		t.Errorf("非包标题下的行不该被当成登记，got %v", stale)
	}
}

// TestRealRepoSurfaceMatchesAPImd 是对 CT-01 ⑤ 的现场自证：
// 真实仓（今天只有占位包）必须绿，且核对区必须存在。
func TestRealRepoSurfaceMatchesAPImd(t *testing.T) {
	const repoRoot = "../.."
	if err := run(options{root: repoRoot, apiPath: "docs/API.md", out: io.Discard}); err != nil {
		t.Errorf("本仓导出面与 docs/API.md 不一致：%v", err)
	}
}

// TestListFlagOutputsSurface 钉住 §4 说明里的生成命令：S1 落码时靠它生成登记行。
func TestListFlagOutputsSurface(t *testing.T) {
	var buf strings.Builder
	if err := run(options{root: fixtureRoot, apiPath: "API.md", listOnly: true, out: &buf}); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"engine/EngineRunner.RunBacktest", "domain/OHLCV"} {
		if !strings.Contains(buf.String(), want) {
			t.Errorf("-list 输出缺少 %s，got:\n%s", want, buf.String())
		}
	}
}

func checkFixture(root string) error {
	return run(options{root: root, apiPath: "API.md", out: io.Discard})
}

func runCompare(t *testing.T, root string) (unregistered, stale []string) {
	t.Helper()
	code, err := scanExports(root)
	if err != nil {
		t.Fatalf("扫描 %s: %v", root, err)
	}
	doc, err := parseRegister(filepath.Join(root, "API.md"))
	if err != nil {
		t.Fatalf("读登记表: %v", err)
	}
	return compare(code, doc)
}

// copyFixture 把夹具复制进临时目录，供「改坏 → 应当红」的用例动手。
func copyFixture(t *testing.T) string {
	t.Helper()
	dst := t.TempDir()
	if err := copyTree(fixtureRoot, dst); err != nil {
		t.Fatal(err)
	}
	return dst
}

func copyTree(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, raw, 0o600)
	})
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func hasSym(list []string, want string) bool {
	for _, item := range list {
		if item == want {
			return true
		}
	}
	return false
}
