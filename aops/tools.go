package aops

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"sort"
	"strings"
)

func getAddFuncWithoutDependsStmt(sp StmtParams) (expr []ast.Expr, err error) {
	for _, s := range sp.Stmts {
		switch s.Kind {
		case AddFuncWithoutDepends:
			return getExprsFromStmt(s.Stmt)
		}
	}
	
	return
}

func getDeferFuncStmt(sp StmtParams) (stmts []ast.Stmt, err error) {
	return getStmt(sp.Stmts, AddDeferFuncStmt)
}

func getDeferWithVarFuncStmt(sp StmtParams) (stmts []ast.Stmt, err error) {
	return getStmt(sp.Stmts, AddDeferFuncWithVarStmt)
}

func getReturnFuncWithoutVarStmt(sp StmtParams) (stmts []ast.Stmt, err error) {
	return getStmt(sp.Stmts, AddReturnFuncWithoutVarStmt)
}

func getReturnFuncWithVarStmt(sp StmtParams) (stmts []ast.Stmt, depends []string, err error) {
	for _, s := range sp.Stmts {
		if s.Kind == AddReturnFuncWithVarStmt {
			stmts, err := getStmt(sp.Stmts, AddReturnFuncWithVarStmt)
			return stmts, s.Depends, err
		}
	}
	
	return
}
func getFuncStmt(sp StmtParams) (stmts []ast.Stmt, depends []string, err error) {
	for _, s := range sp.Stmts {
		if s.Kind == AddFuncWithVarStmt {
			stmts, err := getStmt(sp.Stmts, s.Kind)
			return stmts, s.Depends, err
		}
	}
	
	return
}

func getStmt(stmt []StmtParam, id OperationKind) (stmts []ast.Stmt, err error) {
	for _, s := range stmt {
		if s.Kind == id {
			return getStmtsFromStmt(s.Stmt)
		}
	}
	
	return nil, nil
}

// getStmtsFromStmt parser stmt string to ast.Stmt
func getStmtsFromStmt(stmt []string) (stmts []ast.Stmt, err error) {
	for _, e := range stmt {
		s, err := parserStmt(e)
		if err != nil {
			return nil, err
		}
		
		stmts = append(stmts, s)
	}
	
	return stmts, nil
}

// getExprsFromStmt parser expr string to ast.Expr
func getExprsFromStmt(stmt []string) (expr []ast.Expr, err error) {
	for _, e := range stmt {
		exp, err := parser.ParseExpr(e)
		if err != nil {
			return nil, err
		}
		
		expr = append(expr, exp)
	}
	
	return
}

// getIntersection get intersection from string array and AOP id map.
// arr generated by `extractIdFromComment`. ids usually is a const params.
// If these have intersection, then return a slice that contains intersection id.
// Otherwise, return nil.
func getIntersection(arr []string, ids map[string]struct{}) []string {
	if len(arr) == 0 {
		return nil
	}
	
	var result []string
	for _, a := range arr {
		if _, exist := ids[a]; exist {
			result = append(result, a)
		}
	}
	
	return result
}

// extractIdFromComment get AOP ids from comment.
// Valid AOP id starts with '@' and follows a word.
// like `@Middleware`, `@middleware-a`, `@middleWare`
// are valid
// The id can not end with '@', likes that '@','@middleware@'
// are invalid
func extractIdFromComment(comment string) []string {
	var result []string
	comments := strings.Split(comment, "//")
	for _, _comment := range comments {
		comments := strings.Split(_comment, " ")
		for _, c := range comments {
			if strings.HasPrefix(c, "@") {
				_c := strings.TrimSpace(c)
				if len(_c) > 1 &&
					!strings.HasSuffix(_c, "@") {
					result = append(result, strings.TrimSpace(c))
				}
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

func isEqual(fd *ast.FuncDecl, fn fun) bool {
	if fd.Recv != nil && len(fd.Recv.List) > 0 && fn.owner != "" {
		return fd.Recv.List[0].Type.(*ast.Ident).Name == fn.owner && fd.Name.String() == fn.name
	}
	
	if fd.Recv == nil && fn.owner == "" {
		return fd.Name.String() == fn.name
	}
	
	if fd.Recv != nil && fn.owner == "" {
		return false
	}
	
	return false
}

func fullId(t *ast.FuncDecl) string {
	r := ""
	name := t.Name.String()
	if t.Recv != nil && len(t.Recv.List) > 0 {
		r = t.Recv.List[0].Type.(*ast.Ident).Name
	}
	
	return fmt.Sprintf("%s-%s", name, r)
}

// removeDuplicate m is generated by AddCode, save the file path and AOP ids.
// Maybe AOP ids will duplicate, so use this funciton remove surplus ids.
func removeDuplicate(m map[string][]string) map[string][]string {
	for key, val := range m {
		_m := make(map[string]struct{})
		var _ms []string
		for _, v := range val {
			_m[v] = struct{}{}
		}
		
		for v := range _m {
			_ms = append(_ms, v)
		}
		
		sort.Strings(_ms)
		m[key] = _ms
	}
	return m
}

func parserImport(p Pack) (impor []ast.Spec, err error) {
	comment := fmt.Sprintf(`package main
import %s %s
`, p.Name, p.Path)
	
	f, err := parser.ParseFile(token.NewFileSet(), "", comment, parser.ImportsOnly)
	if err != nil {
		return
	}
	
	if len(f.Decls) == 0 {
		return nil, fmt.Errorf("import parser failed. please verify Pack[%s %s] whether valid. ", p.Name, p.Path)
	}
	
	return f.Decls[0].(*ast.GenDecl).Specs, nil
}
