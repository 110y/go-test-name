package visitor

import (
	"go/ast"
	"go/token"
	"go/types"
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

	inTestCasesStatement bool
	inMapExpr            bool
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

	if !v.inTestCasesStatement {
		if stmt, ok := node.(*ast.AssignStmt); ok {
			if len(stmt.Rhs) > 0 {
				if clit, ok := stmt.Rhs[0].(*ast.CompositeLit); ok {
					if _, ok := clit.Type.(*ast.MapType); ok {
						for _, elt := range clit.Elts {
							if v.cursorPos >= v.getPositionOffset(elt.Pos()) && v.cursorPos <= v.getPositionOffset(elt.End()) {
								if kve, ok := elt.(*ast.KeyValueExpr); ok {
									if bl, ok := kve.Key.(*ast.BasicLit); ok && bl.Kind == token.STRING {
										s, err := strconv.Unquote(bl.Value)
										if err != nil {
											return nil
										}
										v.subTestNames = append(v.subTestNames, s)
										return nil
									}
								}
							}
						}
					}
				}
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

func (v *Visitor) getPositionOffset(pos token.Pos) int {
	return v.fileset.Position(pos).Offset
}
