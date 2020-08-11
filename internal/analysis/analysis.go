package analysis

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"

	"golang.org/x/tools/go/packages"

	"github.com/110y/go-test-name/internal/analysis/internal/visitor"
)

var pkgConfigMode = packages.NeedName |
	packages.NeedFiles |
	packages.NeedImports |
	packages.NeedSyntax |
	packages.NeedTypesInfo |
	packages.NeedTypes |
	packages.NeedTypesSizes

func GetTestInfo(ctx context.Context, path string, pos int) (*TestInfo, error) {
	fs := token.NewFileSet()

	fpath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get abs file path: %w", err)
	}

	cfg := &packages.Config{
		Context: ctx,
		Fset:    fs,
		Dir:     filepath.Dir(fpath),
		Mode:    pkgConfigMode,
		Tests:   true,
	}
	pkgs, err := packages.Load(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load package: %w", err)
	}

	var pkgIdx int
	var fIdx int
	for i, pkg := range pkgs {
		for j, f := range pkg.GoFiles {
			if fpath == f {
				fIdx = j
				pkgIdx = i
			}
		}
	}

	pkg := pkgs[pkgIdx]
	f := pkg.Syntax[fIdx]

	v := visitor.New(pos, fs, pkg.TypesInfo)
	ast.Walk(v, f)

	return &TestInfo{
		TestFuncName: v.GetFuncName(),
		SubTestNames: v.GetSubTestNames(),
	}, nil
}
