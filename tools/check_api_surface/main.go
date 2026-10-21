// Command check_api_surface 核对 docs/API.md 的导出面登记表与本仓实际导出符号（CT-01 ③-③）。
//
// 三个方向都算红，缺一个就是假护栏：
//
//	unregistered —— 代码里有导出符号、登记表没有（新增导出符号未登记即红）
//	stale        —— 登记表里有、代码里没有（改名/搬走没同步登记表即红）
//	结构红       —— docs/API.md 找不到 surface:begin/end 标记（删掉登记表来「让检查通过」即红）
//
// 核对粒度 = go doc 的导出面：包级 func / type / const / var 的名字，
// 加上导出类型的导出方法，登记为 `Type.Method`。
// 不参与核对：package main（工具与命令目录）、*_test.go、未导出符号、结构体字段。
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	surfaceBegin = "<!-- surface:begin -->"
	surfaceEnd   = "<!-- surface:end -->"
)

// ignoredDirs 是不参与导出面扫描的目录名（精确匹配，任何层级都生效）。
var ignoredDirs = map[string]bool{
	"testdata": true, // 夹具不是导出面
	"vendor":   true, // vendored 依赖的符号不是本仓的承诺
}

type options struct {
	root     string
	apiPath  string
	listOnly bool
	out      io.Writer
}

func main() {
	var opts options
	flag.StringVar(&opts.root, "root", ".", "仓库根目录（包目录相对于它计算）")
	flag.StringVar(&opts.apiPath, "api", "docs/API.md", "接口规格文档路径（相对仓库根）")
	flag.BoolVar(&opts.listOnly, "list", false, "只打印实际导出面（不核对登记表），用于生成登记行")
	flag.Parse()

	opts.out = os.Stdout
	if err := run(opts); err != nil {
		fmt.Fprintln(os.Stderr, "check_api_surface:", err)
		os.Exit(1)
	}
}

func run(opts options) error {
	out := opts.out
	if out == nil {
		out = os.Stdout
	}
	code, err := scanExports(opts.root)
	if err != nil {
		return err
	}
	if opts.listOnly {
		for _, sym := range flatten(code) {
			fmt.Fprintln(out, sym)
		}
		fmt.Fprintf(out, "（共 %d 个导出符号，分布在 %d 个包）\n", countSym(code), len(code))
		return nil
	}

	doc, err := parseRegister(filepath.Join(opts.root, filepath.FromSlash(opts.apiPath)))
	if err != nil {
		return err
	}
	unregistered, stale := compare(code, doc)
	fmt.Fprintf(out, "导出面核对：%d 个代码包 / %d 个登记包 ｜ 代码符号 %d ｜ 登记符号 %d\n",
		len(code), len(doc), countSym(code), countSym(doc))
	if len(unregistered) == 0 && len(stale) == 0 {
		fmt.Fprintln(out, "✅ docs/API.md 与导出面一一对应")
		return nil
	}
	if len(unregistered) > 0 {
		fmt.Fprintf(out, "❌ 未登记的导出符号（%d）—— 在 %s 的机器核对区补行：\n", len(unregistered), opts.apiPath)
		for _, sym := range unregistered {
			fmt.Fprintln(out, "   +", sym)
		}
	}
	if len(stale) > 0 {
		fmt.Fprintf(out, "❌ 登记表里有、代码里没有（%d）—— 删行，或把符号补回来：\n", len(stale))
		for _, sym := range stale {
			fmt.Fprintln(out, "   -", sym)
		}
	}
	return fmt.Errorf("导出面与登记表不一致：%d 未登记 + %d 已失效", len(unregistered), len(stale))
}

// scanExports 返回 包目录（slash 相对路径，"." 表示根包）→ 该包导出符号名（升序）。
func scanExports(root string) (map[string][]string, error) {
	dirs, err := goDirs(root)
	if err != nil {
		return nil, err
	}
	pkgs := map[string][]string{}
	for _, dir := range dirs {
		syms, err := scanPackageDir(root, dir)
		if err != nil {
			return nil, err
		}
		pkgs[dir] = syms
	}
	return pkgs, nil
}

// goDirs 列出 root 下含非测试 .go 文件的目录（相对 root，slash 形式；根包记为 "."）。
func goDirs(root string) ([]string, error) {
	set := map[string]bool{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if path != root && (ignoredDirs[name] || strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_")) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".go") || strings.HasSuffix(d.Name(), "_test.go") {
			return nil
		}
		key, err := filepath.Rel(root, filepath.Dir(path))
		if err != nil {
			return err
		}
		set[filepath.ToSlash(key)] = true
		return nil
	})
	if err != nil {
		return nil, err
	}
	dirs := make([]string, 0, len(set))
	for key := range set {
		dirs = append(dirs, key)
	}
	sort.Strings(dirs)
	return dirs, nil
}

// scanPackageDir 解析一个目录里的非测试 Go 文件，返回该目录的导出面符号名。
// package main 返回空集：命令与工具不属于库的导出面。
func scanPackageDir(root, dir string) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(dir)))
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	var syms []string
	pkgNames := map[string]bool{}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		// 语法错误必须当场红：解析失败的文件里所有符号会「凭空消失」，护栏静默空过。
		file, err := parser.ParseFile(fset, filepath.Join(root, filepath.FromSlash(dir), name), nil, parser.SkipObjectResolution)
		if err != nil {
			return nil, fmt.Errorf("解析 %s/%s: %w", dir, name, err)
		}
		pkgNames[file.Name.Name] = true
		if file.Name.Name == "main" {
			continue
		}
		syms = append(syms, fileExports(file)...)
	}
	if len(pkgNames) > 1 {
		return nil, fmt.Errorf("%s 里有多个 package（%s）：导出面按包登记，一个目录只能有一个包", dir, joinSet(pkgNames))
	}
	return dedupe(syms), nil
}

func joinSet(set map[string]bool) string {
	names := make([]string, 0, len(set))
	for name := range set {
		names = append(names, name)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

// fileExports 收集一个文件里的导出符号：包级声明的名字 + 导出类型的导出方法。
func fileExports(file *ast.File) []string {
	var syms []string
	for _, decl := range file.Decls {
		switch spec := decl.(type) {
		case *ast.FuncDecl:
			if spec.Name == nil || !spec.Name.IsExported() {
				continue
			}
			if spec.Recv == nil || len(spec.Recv.List) == 0 {
				syms = append(syms, spec.Name.Name)
				continue
			}
			// 方法只有挂在导出的具名类型上才进导出面；未导出类型的方法外部拿不到。
			if recv := receiverTypeName(spec.Recv.List[0].Type); ast.IsExported(recv) {
				syms = append(syms, recv+"."+spec.Name.Name)
			}
		case *ast.GenDecl:
			syms = append(syms, genDeclExports(spec)...)
		}
	}
	return syms
}

func genDeclExports(decl *ast.GenDecl) []string {
	var syms []string
	for _, spec := range decl.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			if !s.Name.IsExported() {
				continue
			}
			syms = append(syms, s.Name.Name)
			syms = append(syms, interfaceMethods(s.Name.Name, s.Type)...)
		case *ast.ValueSpec:
			for _, name := range s.Names {
				if name.IsExported() {
					syms = append(syms, name.Name)
				}
			}
		}
	}
	return syms
}

// interfaceMethods 把接口的导出方法登记为 `Iface.Method`：端口的形状就是方法集，
// 给接口加一个方法对实现方是破坏性变更，必须过登记表。
func interfaceMethods(ifaceName string, expr ast.Expr) []string {
	it, ok := expr.(*ast.InterfaceType)
	if !ok {
		return nil
	}
	var syms []string
	for _, field := range it.Methods.List {
		for _, name := range field.Names {
			if name.IsExported() {
				syms = append(syms, ifaceName+"."+name.Name)
			}
		}
	}
	return syms
}

func receiverTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return receiverTypeName(t.X)
	case *ast.Ident:
		return t.Name
	case *ast.IndexExpr: // 泛型接收者 T[X]
		return receiverTypeName(t.X)
	case *ast.IndexListExpr: // 泛型接收者 T[X, Y]
		return receiverTypeName(t.X)
	}
	return ""
}

var (
	headingRe = regexp.MustCompile(`^#{2,4}\s+(.*)$`)
	// 小节标题里用反引号写包目录：`engine/`；根包写 `.`（不认它的话，根包新增导出符号就无路可登记）。
	dirInTick  = regexp.MustCompile("\x60([^\x60]+/|\\.)\x60")
	cellSymbol = regexp.MustCompile("^\x60([^\x60]+)\x60$")
)

// parseRegister 读 docs/API.md 的机器核对区，返回 包目录 → 登记符号名。
// 找不到 surface:begin/end 标记即报错 —— 把登记表整段删掉时护栏必须红，不能空过。
func parseRegister(path string) (map[string][]string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读登记表 %s: %w", path, err)
	}
	lines := strings.Split(strings.ReplaceAll(string(raw), "\r\n", "\n"), "\n")
	begin, end := -1, -1
	for i, line := range lines {
		switch strings.TrimSpace(line) {
		case surfaceBegin:
			if begin >= 0 {
				return nil, fmt.Errorf("%s 里 %s 出现了多次", path, surfaceBegin)
			}
			begin = i
		case surfaceEnd:
			end = i
		}
	}
	if begin < 0 || end <= begin {
		return nil, fmt.Errorf("%s 缺少 %s / %s 标记对（机器核对区必须存在，可以为空）", path, surfaceBegin, surfaceEnd)
	}

	register := map[string][]string{}
	dir := ""
	for _, line := range lines[begin+1 : end] {
		if m := headingRe.FindStringSubmatch(strings.TrimSpace(line)); m != nil {
			dm := dirInTick.FindStringSubmatch(m[1])
			if dm == nil {
				dir = "" // 小节不属于任何包：它下面的表行不算登记
				continue
			}
			dir = strings.Trim(dm[1], "/")
			if _, ok := register[dir]; !ok {
				register[dir] = nil
			}
			continue
		}
		if dir == "" || !strings.HasPrefix(strings.TrimSpace(line), "|") {
			continue
		}
		cells := strings.Split(strings.TrimSpace(line), "|")
		if len(cells) < 2 {
			continue
		}
		m := cellSymbol.FindStringSubmatch(strings.TrimSpace(cells[1]))
		if m == nil {
			continue // 表头、分隔行、非符号列
		}
		register[dir] = append(register[dir], m[1])
	}
	for dir := range register {
		register[dir] = dedupe(register[dir])
	}
	return register, nil
}

// compare 返回未登记（代码有、表没有）与已失效（表有、代码没有）两个升序列表，
// 每项写作 包目录/符号。
func compare(code, doc map[string][]string) (unregistered, stale []string) {
	for dir, syms := range code {
		have := toSet(doc[dir])
		for _, sym := range syms {
			if !have[sym] {
				unregistered = append(unregistered, dir+"/"+sym)
			}
		}
	}
	for dir, syms := range doc {
		have := toSet(code[dir])
		for _, sym := range syms {
			if !have[sym] {
				stale = append(stale, dir+"/"+sym)
			}
		}
	}
	sort.Strings(unregistered)
	sort.Strings(stale)
	return unregistered, stale
}

func toSet(syms []string) map[string]bool {
	set := make(map[string]bool, len(syms))
	for _, sym := range syms {
		set[sym] = true
	}
	return set
}

func flatten(pkgs map[string][]string) []string {
	var all []string
	for dir, syms := range pkgs {
		for _, sym := range syms {
			all = append(all, dir+"/"+sym)
		}
	}
	sort.Strings(all)
	return all
}

func countSym(pkgs map[string][]string) int {
	n := 0
	for _, syms := range pkgs {
		n += len(syms)
	}
	return n
}

func dedupe(syms []string) []string {
	seen := make(map[string]bool, len(syms))
	out := make([]string, 0, len(syms))
	for _, sym := range syms {
		if seen[sym] {
			continue
		}
		seen[sym] = true
		out = append(out, sym)
	}
	sort.Strings(out)
	return out
}
