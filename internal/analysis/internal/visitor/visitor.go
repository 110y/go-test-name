package visitor

import (
	"go/ast"
	"go/token"
	"go/types"
	"regexp"
	"strconv"
)

var _ ast.Visitor = (*Visitor)(nil)

func New(pos int, fs *token.FileSet, info *types.Info) *Visitor {
	return &Visitor{
		fileset:      fs,
		cursorPos:    pos,
		info:         info,
		subTestNames: []string{},
	}
}

type Visitor struct {
	cursorPos int
	fileset   *token.FileSet
	info      *types.Info

	testFuncName string
	subTestNames []string
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	startPos := v.getPositionOffset(node.Pos())
	endPos := v.getPositionOffset(node.End())

	if v.cursorPos < startPos || v.cursorPos > endPos {
		return nil
	}

	if v.testFuncName == "" {
		if fd, ok := node.(*ast.FuncDecl); ok {
			v.testFuncName = fd.Name.Name
		}

		return v
	}

	if stmt, ok := node.(*ast.AssignStmt); ok {
		if len(stmt.Rhs) > 0 {
			if clit, ok := stmt.Rhs[0].(*ast.CompositeLit); ok {
				if v.extractSubTestNamesFromCompositeLit(clit) {
					return nil
				}
			}
		}
	}

	if stmt, ok := node.(*ast.RangeStmt); ok {
		if clit, ok := stmt.X.(*ast.CompositeLit); ok {
			if v.extractSubTestNamesFromCompositeLit(clit) {
				return nil
			}
		}
	}

	return v
}

func (v *Visitor) GetFuncName() string {
	return v.testFuncName
}

func (v *Visitor) GetSubTestNames() []string {
	return v.subTestNames
}

func (v *Visitor) extractSubTestNamesFromCompositeLit(clit *ast.CompositeLit) bool {
	if _, ok := clit.Type.(*ast.MapType); ok {
		for _, elt := range clit.Elts {
			if v.cursorPos >= v.getPositionOffset(elt.Pos()) && v.cursorPos <= v.getPositionOffset(elt.End()) {
				if kve, ok := elt.(*ast.KeyValueExpr); ok {
					if bl, ok := kve.Key.(*ast.BasicLit); ok && bl.Kind == token.STRING {
						s, err := strconv.Unquote(bl.Value)
						if err != nil {
							return true
						}
						v.subTestNames = append(v.subTestNames, regexp.QuoteMeta(s))
						return true
					}
				}
			}
		}
	}
	return false
}

func (v *Visitor) getPositionOffset(pos token.Pos) int {
	return v.fileset.Position(pos).Offset
}
