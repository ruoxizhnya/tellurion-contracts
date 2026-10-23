// Command check_contract_shape 扫描契约仓的形状红线（QT-09 ⑤）。
//
// 四条红向，缺一条就是假护栏：
//
//	命名红 —— domain/ / engine/ / fees/ 里出现导出的 `*Config` / `*Options` / `*Impl` 形状名
//	          （精确、前缀或后缀，如 TradingConfig / FillOptions / XxxImpl）。裁决见
//	          docs/adr/0002-config-族的边界处置.md：这两类名字会把部署/实现细节写进契约。
//	DB 红  —— 任何非测试源文件 import 数据库/驱动包。契约仓没有实现、没有 DB 访问，
//	          一旦契约包 import 驱动，下游就被迫带上存储依赖。
//	空扫红 —— 一个 Go 包都没扫到（目录被删空 / 改名）也算红。删掉检查对象让护栏变绿，
//	          是这类扫描器最常见的失效方式。
//	缺根红 —— domain/ / engine/ / fees/ 任一目录扫不到非测试 Go 文件同样红：整棵树被搬走
//	          时「没扫到命中」与「什么都没扫」必须被区分开。
//
// 白名单是显式的、逐符号登记理由的：例外必须先在 ADR 里成立，再写在这里，不许在 commit
// message 里口头豁免。只做 go/parser + go/ast 的语法级扫描（不用 go/types），
// 因此 CI 不需要先编译成功、也不解析依赖图。
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// contractRoots 是三条红线的作用域（目录名，含其子包，如 domain/market）。
var contractRoots = []string{"domain", "engine", "fees"}

// ignoredDirs 与 check_api_surface 同口径：夹具与 vendored 代码不是本仓的承诺。
var ignoredDirs = map[string]bool{
	"testdata": true,
	"vendor":   true,
}

// shapeKeywords 是契约层禁止出现的实现细节形状名。
var shapeKeywords = []string{"Config", "Options", "Impl"}

// allowedNames 是显式白名单：键为 `包目录/符号名`（目录用 slash 相对路径）。
//
// 每一行都必须写清「为什么它可以留」——白名单不是逃生门，是 ADR 的执行形态。
// ⚠️ 与 QT-09 工单「只允许 TradingConfig / PriceLimitConfig 两个名字」的冲突：
// 本仓落码时（QT-05 ③-①）按 ADR-0002 ② 已迁入 StrategyConfig / ExecutionConfig，
// 按 ADR-0002 ③ 保留 TradingConfig / PriceLimitConfig；改名 = 改已冻结的 429 导出面，
// 本单明令禁止。因此白名单只能如实登记这 7 个已存在的符号（2 个冻结面伴生名 +
// 2 个 ADR-0002 ② 的 core-only 搭车类型 + 3 个它们的派生构造函数），
// 新增任何一个 `*Config` / `*Options` / `*Impl` 仍然当场红。
var allowedNames = map[string]string{
	"engine/TradingConfig":             "ADR-0002 ③ 冻结面伴生名：名字已写进 tellurion:docs/SPEC.md 的 RunSpec 冻结文本，改名 = 改冻结文本",
	"engine/PriceLimitConfig":          "ADR-0002 ③ 冻结面伴生名：同上，是 TradingConfig.PriceLimit 字段的冻结类型",
	"engine/DefaultTradingConfig":      "随冻结点 TradingConfig 派生的默认值构造函数；改名即改冻结面（docs/API.md §5.1）",
	"engine/NewMarketRouterFromConfig": "市场路由装配入口，名字里的 Config 指调用方传入的路由表参数，不承载契约配置形状（docs/API.md §5.2 MarketRouter）",
	"domain/StrategyConfig":            "ADR-0002 ②：core-only 搭车类型，按原名迁入并进 QT-05 ③-⑨ 回收清单，不就地改名（改名 = 给即将搬走的类型付一次破坏性变更）",
	"domain/ExecutionConfig":           "ADR-0002 ②：同上，core-only 搭车类型，S2 后回收",
	"domain/DefaultExecutionConfig":    "随 StrategyConfig / ExecutionConfig 族一起迁入的默认值构造函数（ADR-0002 ② 的在册例外）",
}

// dbImportPrefixes 是「契约包不得出现」的数据库/驱动依赖。
// 至少拦这几类：标准库 database/sql（含子包）、pgx（v4/v5 前缀相同）、lib/pq、
// go-redis（v8/v9）、sqlx、gorm（gorm.io/ 整个域）。新增驱动按前缀补进来。
var dbImportPrefixes = []string{
	"database/sql",
	"github.com/jackc/pgx",
	"github.com/lib/pq",
	"github.com/redis/go-redis",
	"github.com/jmoiron/sqlx",
	"gorm.io/",
}

type options struct {
	root string
	out  *os.File
}

type sourceFile struct {
	dir  string // slash 相对目录，根包记为 "."
	rel  string // slash 相对路径（给人读的 文件:行号）
	file *ast.File
}

func main() {
	var opts options
	flag.StringVar(&opts.root, "root", ".", "仓库根目录（包目录相对于它计算）")
	flag.Parse()

	opts.out = os.Stdout
	if err := run(opts); err != nil {
		fmt.Fprintln(os.Stderr, "check_contract_shape:", err)
		os.Exit(1)
	}
}

func run(opts options) error {
	files, fset, err := scanGoFiles(opts.root)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("一个 Go 包都没扫到（%s 下没有非测试 .go 文件）—— 护栏拒绝空绿", opts.root)
	}

	pkgs := map[string]bool{}
	var red []string
	for _, src := range files {
		pkgs[src.dir] = true
		for _, imp := range src.file.Imports {
			path := strings.Trim(imp.Path.Value, `"`)
			if !isDBImport(path) {
				continue
			}
			pos := fset.Position(imp.Pos())
			red = append(red, fmt.Sprintf("%s:%d: import %q —— 契约仓禁止数据库/驱动依赖", src.rel, pos.Line, path))
		}
		if !inContractScope(src.dir) {
			continue
		}
		for _, decl := range src.file.Decls {
			red = append(red, shapeViolations(src.rel, src.dir, fset, decl)...)
		}
	}

	var missing []string
	for _, root := range contractRoots {
		if !hasRoot(pkgs, root) {
			missing = append(missing, root+"/")
		}
	}
	if len(missing) > 0 {
		red = append(red, fmt.Sprintf("契约根目录扫不到非测试 Go 文件：%s —— 搬家/删空后「没命中」不等于「没问题」", strings.Join(missing, "、")))
	}

	sort.Strings(red)
	fmt.Fprintf(opts.out, "契约形状核对：%d 个包 / %d 个非测试文件 ｜ 白名单例外 %d 个\n", len(pkgs), len(files), len(allowedNames))
	if len(red) == 0 {
		fmt.Fprintln(opts.out, "✅ 无 *Config/*Options/*Impl 形状名，无 DB/驱动 import，三条契约根目录都在")
		return nil
	}
	for _, r := range red {
		fmt.Fprintln(opts.out, "❌", r)
	}
	return fmt.Errorf("契约形状违规（%d 项）", len(red))
}

// scanGoFiles 解析 root 下所有非测试 .go 文件；语法错误当场红——
// 解析失败的文件里所有声明会「凭空消失」，护栏随之静默空过。
func scanGoFiles(root string) ([]sourceFile, *token.FileSet, error) {
	fset := token.NewFileSet()
	var files []sourceFile
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
		name := d.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			return nil
		}
		file, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			return fmt.Errorf("解析 %s: %w", path, err)
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		dir, err := filepath.Rel(root, filepath.Dir(path))
		if err != nil {
			return err
		}
		files = append(files, sourceFile{
			dir:  filepath.ToSlash(dir),
			rel:  filepath.ToSlash(rel),
			file: file,
		})
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return files, fset, nil
}

// inContractScope 判断目录是否落在 domain/ engine/ fees/ 及其子包里。
func inContractScope(dir string) bool {
	for _, root := range contractRoots {
		if dir == root || strings.HasPrefix(dir, root+"/") {
			return true
		}
	}
	return false
}

func hasRoot(pkgs map[string]bool, root string) bool {
	for dir := range pkgs {
		if dir == root || strings.HasPrefix(dir, root+"/") {
			return true
		}
	}
	return false
}

// shapeViolations 收集一个声明里的命中：包级导出类型名与非方法导出函数名。
// 方法不算：方法名挂在已登记的类型上，泄漏的是类型名（`Type.Method` 由
// check_api_surface 管名字，`check_surface_freeze` 管签名）。
func shapeViolations(rel, dir string, fset *token.FileSet, decl ast.Decl) []string {
	var hits []string
	report := func(pos token.Pos, kind, name string) {
		fam := shapeFamily(name)
		if fam == "" {
			return
		}
		key := dir + "/" + name
		if _, ok := allowedNames[key]; ok {
			// 例外只在盘点里出现，不打印；理由写在 allowedNames 表上。
			return
		}
		hits = append(hits, fmt.Sprintf("%s:%d: 导出%s %s 命中 %s 族（ADR-0002 只放行已登记例外）",
			rel, fset.Position(pos).Line, kind, name, fam))
	}
	switch d := decl.(type) {
	case *ast.FuncDecl:
		if d.Recv != nil || d.Name == nil || !d.Name.IsExported() {
			return nil
		}
		report(d.Name.Pos(), "函数", d.Name.Name)
	case *ast.GenDecl:
		if d.Tok != token.TYPE {
			return nil
		}
		for _, spec := range d.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok || !ts.Name.IsExported() {
				continue
			}
			report(ts.Name.Pos(), "类型", ts.Name.Name)
		}
	}
	return hits
}

// shapeFamily 返回命中族名；精确、前缀或后缀都算（TradingConfig / Config / Configs 前缀式同罪）。
func shapeFamily(name string) string {
	for _, kw := range shapeKeywords {
		if name == kw || strings.HasPrefix(name, kw) || strings.HasSuffix(name, kw) {
			return kw
		}
	}
	return ""
}

// isDBImport 判断 import 路径是否属于数据库/驱动依赖。
func isDBImport(path string) bool {
	for _, prefix := range dbImportPrefixes {
		if strings.HasSuffix(prefix, "/") {
			if strings.HasPrefix(path, prefix) {
				return true
			}
			continue
		}
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			return true
		}
	}
	return false
}
