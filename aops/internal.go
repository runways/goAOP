package aops

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type fun struct {
	originIds []string
	owner     string
	name      string
	aopIds    []string
}

const noReceiver = ""

// fetchDir get all dirs which contains *.go files.
func fetchDir(root string) (dirs []string, err error) {
	dirFilter := make(map[string]interface{})
	
	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".go") {
			p := filepath.Dir(path)
			if _, exist := dirFilter[p]; !exist {
				dirFilter[p] = struct{}{}
			}
		}
		
		return nil
	})
	
	for key := range dirFilter {
		dirs = append(dirs, key)
	}
	sort.Strings(dirs)
	return dirs, nil
}

// ParseDir use `token.ParserDir` parser specify dir. And use `filter` for filter flile info at the same time
// If filter is nil, it will pass all files as default.
// When parse success, `ParseDir` will return a map save file name and package pointer. If failed, return a error
func ParseDir(dir string, filter func(info fs.FileInfo) bool) (map[string]*ast.Package, error) {
	if filter == nil {
		filter = func(info fs.FileInfo) bool {
			return true
		}
	}
	
	dirs, err := fetchDir(dir)
	if err != nil {
		return nil, err
	}
	
	m := make(map[string]*ast.Package)
	for _, d := range dirs {
		_m, err := parser.ParseDir(token.NewFileSet(), d, filter, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		for key, value := range _m {
			m[key] = value
		}
	}
	
	return m, nil
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
										originIds: _ids,
										owner:     noReceiver,
										name:      t.Name.String(),
										aopIds:    validId,
									})
								} else {
									switch owner := t.Recv.List[0].Type.(type) {
									case *ast.Ident:
										functions = append(functions, fun{
											originIds: _ids,
											owner:     owner.Name,
											name:      t.Name.String(),
											aopIds:    validId,
										})
									case *ast.StarExpr:
										functions = append(functions, fun{
											originIds: _ids,
											owner:     owner.X.(*ast.Ident).Name,
											name:      t.Name.String(),
											aopIds:    validId,
										})
									}
									//owner := t.Recv.List[0].Type.(*ast.Ident)
									//functions = append(functions, fun{
									//	originIds: _ids,
									//	owner:     owner.Name,
									//	name:      t.Name.String(),
									//	aopIds:    validId,
									//})
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
// Typically, in the StmtParams, `DeclStmt` and `Pack` can nil. But Stmts cannot be empty.
// Stmts save different OperationKind stmts. More detail info please reference StmtParam usage in types.go.
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
							for idx, id := range fn.aopIds {
								
								ij := injectDetail{
									owner: fn.owner,
									name:  fn.name,
								}
								
								funcVarStmt, exprs, err := ij.getAddFuncWithoutDependsStmt(stmt[id], fn.originIds[idx])
								if err != nil {
									return nil, err
								}
								
								stmts, err := ij.getDeferFuncStmt(stmt[id])
								if err != nil {
									return nil, err
								}
								
								funcs, depends, funcDepends, stmtStr, err := ij.getFuncStmt(stmt[id])
								if err != nil {
									return nil, err
								}
								
								rets, err := ij.getReturnFuncWithoutVarStmt(stmt[id])
								if err != nil {
									return nil, err
								}
								
								retVars, retDepends, err := ij.getReturnFuncWithVarStmt(stmt[id])
								if err != nil {
									return nil, err
								}
								
								err = addFuncWithoutDependsOperator(t, exprs)
								if err != nil {
									return nil, err
								}
								
								err = addStmtAsFuncWithoutVarOperator(t, funcVarStmt)
								if err != nil {
									return nil, err
								}
								
								err = addDeferWithoutVarOperator(t, stmts)
								if err != nil {
									return nil, err
								}
								err = addStmtAsFuncWithVarOperator(t, funcs, depends, funcDepends, stmtStr)
								if err != nil {
									return nil, err
								}
								err = addStmtAsReturnOperator(t, rets)
								if err != nil {
									return nil, err
								}
								
								err = addReturnWithBindVarOperator(t, retVars, retDepends)
								if err != nil {
									return nil, err
								}
								
								err = addStmtBindVarOperator(t, stmt[id].DeclStmt)
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
