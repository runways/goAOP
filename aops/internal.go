package aops

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/scanner"
	"strings"
)

// position Get all functions that need add AOP.
//
// Pkgs should generate by `parser.ParseDir` and `id` is the AOP middleware name, e.g. @trace.
//
// This function will ignore *_test.go. It will return a map(map[string][]string), key is file
// name, value is a function name array.
func position(pkgs map[string]*ast.Package, id string) map[string][]string {
	result := make(map[string][]string)

	for _, pack := range pkgs {
		for name, f := range pack.Files {
			if strings.HasSuffix(name, "_test.go") {
				continue
			}
			var functions []string
			for _, funDecl := range f.Decls {
				switch t := funDecl.(type) {
				case *ast.FuncDecl:
					if t.Doc != nil {
						for _, c := range t.Doc.List {
							if strings.Contains(c.Text, id) {
								functions = append(functions, t.Name.String())
							}
						}
					}
				}
			}
			if len(functions) > 0 {
				result[name] = functions
			}
		}
	}
	
	return result
}

// parserStmt Convert stmt to ast.Stmt
// If the stmt is a valid stmt, then return ast.Stmt.
// Otherwise, return a error.
func parserStmt(stmt string) (ast.Stmt, error) {
	expr := "func(){" + stmt + ";}"
	if e, err := parser.ParseExpr(expr); err != nil {
		if e, err := parser.ParseExpr(stmt); err == nil {
			return &ast.ExprStmt{X: e}, nil
		}
		errs := err.(scanner.ErrorList)
		for i := range errs {
			errs[i].Pos.Offset -= 7
			errs[i].Pos.Column -= 7
		}
		return nil, errs
	} else {
		node := e.(*ast.FuncLit).Body.List[0]
		if stmt, ok := node.(ast.Stmt); !ok {
			return nil, fmt.Errorf("%T not supported", node)
		} else {
			return stmt, nil
		}
	}
}
