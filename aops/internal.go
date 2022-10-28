package aops

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/scanner"
	"go/token"
	"io/fs"
	"os"
	"sort"
	"strings"
)

type fun struct {
	owner  string
	name   string
	aopIds []string
}

const noReceiver = ""

// ParseDir use `token.ParserDir` parser specify dir. And use `filter` for filter flile info at the same time
// If filter is nil, it will pass all files as default.
// When parse success, `ParseDir` will return a map save file name and package pointer. If failed, return a error
func ParseDir(dir string, filter func(info fs.FileInfo) bool) (map[string]*ast.Package, error) {
	if filter == nil {
		filter = func(info fs.FileInfo) bool {
			return true
		}
	}
	
	return parser.ParseDir(token.NewFileSet(), dir, filter, parser.ParseComments)
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

// Position Get all functions that need add AOP.
//
// Pkgs should generate by `parser.ParseDir` and `id` is the AOP middleware name, e.g. @trace.
//
// This function will ignore *_test.go. It will return a map(map[string][]string), key is file
// name, value is a function name array.
func Position(pkgs map[string]*ast.Package, ids map[string]struct{}) map[string][]fun {
	result := make(map[string][]fun)
	
	for _, pack := range pkgs {
		for name, f := range pack.Files {
			if strings.HasSuffix(name, "_test.go") {
				continue
			}
			var functions []fun
			for _, funDecl := range f.Decls {
				switch t := funDecl.(type) {
				case *ast.FuncDecl:
					if t.Doc != nil {
						for _, c := range t.Doc.List {
							// get all valid AOP ids from comment
							_ids := extractIdFromComment(c.Text)
							
							validId := getIntersection(_ids, ids)
							if len(validId) > 0 {
								if t.Recv == nil ||
									len(t.Recv.List) == 0 {
									functions = append(functions, fun{
										owner:  noReceiver,
										name:   t.Name.String(),
										aopIds: validId,
									})
								} else {
									owner := t.Recv.List[0].Type.(*ast.Ident)
									functions = append(functions, fun{
										owner:  owner.Name,
										name:   t.Name.String(),
										aopIds: validId,
									})
								}
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

// AddImport Add import package for build.
// pkgs is generated by `position` function.
//
// stmt is standard data, usually pass by user. The stmt save aop id and it's params.
//
// modify is generated by AddCode. It save the files that have modified.
//
// If replace origin file, then set replace true, otherwise, set false.
func AddImport(pkgs map[string][]fun, stmt map[string]StmtParams, modify map[string][]string, replace bool) error {
	for name := range pkgs {
		aopIds, exist := modify[name]
		if !exist {
			continue
		}
		
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, name, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		
		decls := make([]ast.Decl, 0, len(f.Decls))
		for _, decl := range f.Decls {
			switch t := decl.(type) {
			case *ast.GenDecl:
				var stats []ast.Spec
				_, ok := t.Specs[0].(*ast.ImportSpec)
				if !ok {
					decls = append(decls, t)
					continue
				}
				
				for _, i := range aopIds {
					for _, p := range stmt[i].Packs {
						impor, err := parserImport(p)
						if err != nil {
							return err
						}
						stats = append(stats, impor...)
					}
				}
				
				stats = append(stats, t.Specs...)
				t.Specs = stats
				decls = append(decls, t)
			default:
				decls = append(decls, t)
			}
		}
		f.Decls = decls
		
		cfg := printer.Config{
			Mode: printer.UseSpaces,
		}
		var buf bytes.Buffer
		
		cfg.Fprint(&buf, fset, f)
		
		dest, err := format.Source(buf.Bytes())
		if err != nil {
			return err
		}
		
		if replace {
			os.WriteFile(name, dest, 0777)
		} else {
			fmt.Println(string(dest))
		}
	}
	return nil
}

// AddCode Insert AOP code to source code files.
// `pkgs` is map that save file name and function names.
// `pkgs` is generated by `position` function.
//
// stmt is a map, that save function declare, key is aop id.
// stmtParams save the AOP stmt params. There are some things need attentions.
//
// First,in the FunStmt. A simple `fun` declare does not have any effect, e.g. `func(){ fmt.Println() }`.
// If wants to execute this function, invoke function at the end, like `func(){}()`
//
// Second, DeferStmt also is a string slice, that save defer function declare, e.g. `defer func(){}()`.
// When use deferStmt, please keep deferStmt valid.
//
// At last, Packs save the import data. Maybe user has import the same package, so named a unique name
// for repeat is a good idea.
//
// When adding code, it will add funStmt first, and add deferStmt by follow.
//
// Replace used to indicate replace source file or not. If replace == true, it replaces at the end.
// Otherwise, it will not.
func AddCode(pkgs map[string][]fun, stmt map[string]StmtParams, replace bool) (map[string][]string, error) {
	modify := make(map[string][]string)
	for name, funs := range pkgs {
		
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, name, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		
		fm := make(map[string][]fun)
		var addId []string
		
		for _, n := range funs {
			ns, exist := fm[fmt.Sprintf("%s-%s", n.name, n.owner)]
			if exist {
				ns = append(ns, n)
				fm[fmt.Sprintf("%s-%s", n.name, n.owner)] = ns
			} else {
				fm[fmt.Sprintf("%s-%s", n.name, n.owner)] = []fun{n}
			}
			
		}
		
		decls := make([]ast.Decl, 0, len(f.Decls))
		for _, decl := range f.Decls {
			switch t := decl.(type) {
			case *ast.FuncDecl:
				if _fn, exist := fm[fullId(t)]; exist {
					for _, fn := range _fn {
						if isEqual(t, fn) {
							var stats []ast.Stmt
							for _, id := range fn.aopIds {
								
								funStmt := stmt[id].FunStmt
								deferStmt := stmt[id].DeferStmt
								funcStmt := stmt[id].FunVarStmt
								
								exprs := make([]ast.Expr, len(funStmt))
								stmts := make([]ast.Stmt, len(deferStmt))
								funcs := make([]ast.Stmt, len(funcStmt))
								
								for idx, c := range funStmt {
									exprInsert, err := parser.ParseExpr(c)
									if err != nil {
										return nil, err
									}
									exprs[idx] = exprInsert
								}
								
								for idx, c := range deferStmt {
									stmtInsert, err := parserStmt(c)
									if err != nil {
										return nil, err
									}
									stmts[idx] = stmtInsert
								}
								
								for idx, c := range funcStmt {
									funcInsert, err := parserStmt(c)
									if err != nil {
										return nil, err
									}
									funcs[idx] = funcInsert
								}
								
								for _, e := range exprs {
									stats = append(stats, &ast.ExprStmt{
										X: e,
									})
								}
								
								for _, e := range stmts {
									stats = append(stats, e)
								}
								
								// Invoke addStmtAdReturn to complete return code.
								err = addStmtAsReturn(t, funcs)
								if err != nil {
									return nil, err
								}
								
								stats = append(stats, t.Body.List...)
								if len(stats) > 0 {
									addId = append(addId, id)
								}
							}
							
							t.Body.List = stats
						}
					}
					
				}
				decls = append(decls, t)
			default:
				decls = append(decls, t)
			}
		}
		f.Decls = decls
		
		cfg := printer.Config{
			Mode: printer.UseSpaces,
		}
		var buf bytes.Buffer
		
		cfg.Fprint(&buf, fset, f)
		
		dest, err := format.Source(buf.Bytes())
		if err != nil {
			return nil, err
		}
		if replace {
			os.WriteFile(name, dest, 0777)
		} else {
			fmt.Println(string(dest))
		}
		
		modify[name] = addId
	}
	
	return removeDuplicate(modify), nil
}

// addStmtAsReturn check whether this function has a func variable as return data.
// If it has function as return, then add pre-defined code. Otherwise, fallthrough.
func addStmtAsReturn(t *ast.FuncDecl, fun []ast.Stmt) error {
	if len(t.Body.List) > 0 {
		for _, _f := range t.Body.List {
			rf, ok := _f.(*ast.ReturnStmt)
			if ok {
				//	One function only has one return stmt block. So no need check twice.
				return _addFuncCode(rf, fun)
			}
		}
	}
	
	return nil
}

func _addFuncCode(t *ast.ReturnStmt, exprs []ast.Stmt) error {
	
	for _, returnFunc := range t.Results {
		rf, ok := returnFunc.(*ast.FuncLit)
		if ok && len(rf.Body.List) > 0 {
			stats := make([]ast.Stmt, 0, len(rf.Body.List)+len(exprs))
			for _, e := range exprs {
				stats = append(stats, e)
			}
			
			stats = append(stats, rf.Body.List...)
			rf.Body.List = stats
		}
	}
	
	return nil
}
