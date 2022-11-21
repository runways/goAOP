package aops

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"reflect"
	"testing"
)

func Test_parserStmt(t *testing.T) {
	type args struct {
		stmt string
	}
	tests := []struct {
		name    string
		args    args
		want    ast.Stmt
		wantErr bool
	}{
		{
			name: "valid stmt",
			args: struct{ stmt string }{stmt: `var b = 11`},
			want: &ast.DeclStmt{Decl: &ast.GenDecl{
				Doc:    nil,
				TokPos: 8,
				Tok:    token.VAR,
				Lparen: token.NoPos,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Doc: nil,
						Names: []*ast.Ident{&ast.Ident{
							NamePos: token.Pos(12),
							Name:    "b",
							Obj:     nil,
						}},
						Type: nil,
						Values: []ast.Expr{
							&ast.BasicLit{
								ValuePos: token.Pos(16),
								Kind:     5,
								Value:    "11",
							},
						},
						Comment: nil,
					},
				},
				Rparen: token.NoPos,
			}},
			wantErr: false,
		},
		{
			name:    "inValid stmt",
			args:    struct{ stmt string }{stmt: `var b=`},
			wantErr: true,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parserStmt(tt.args.stmt)
			if (err != nil) != tt.wantErr {
				t.Errorf("parserStmt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parserStmt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Position(t *testing.T) {
	
	pkg, _ := parser.ParseDir(token.NewFileSet(), "../unitTests", func(info fs.FileInfo) bool {
		//	ignore all logic check
		return true
	}, parser.ParseComments)
	
	type args struct {
		pkgs map[string]*ast.Package
		id   map[string]struct{}
	}
	tests := []struct {
		name string
		args args
		want map[string][]fun
	}{
		{
			name: "valid position",
			args: struct {
				pkgs map[string]*ast.Package
				id   map[string]struct{}
			}{
				pkgs: pkg, id: map[string]struct{}{"@middleware-a": struct{}{}},
			},
			want: map[string][]fun{
				"../unitTests/test.go": []fun{
					{
						owner: "FirstStruct",
						name:  "InvokeFirstFunction",
						aopIds: []string{
							"@middleware-a",
						},
					},
					{
						owner: "",
						name:  "InvokeFirstFunction",
						aopIds: []string{
							"@middleware-a",
						},
					},
				},
			},
		},
		{
			name: "nil position",
			args: struct {
				pkgs map[string]*ast.Package
				id   map[string]struct{}
			}{pkgs: nil, id: map[string]struct{}{"@middleware-a": struct{}{}}},
			want: map[string][]fun{},
		},
		{
			name: "no exist position",
			args: struct {
				pkgs map[string]*ast.Package
				id   map[string]struct{}
			}{pkgs: pkg, id: map[string]struct{}{"@middleware-c": struct{}{}}},
			want: map[string][]fun{},
		},
		{
			name: "nil and no exist position",
			args: struct {
				pkgs map[string]*ast.Package
				id   map[string]struct{}
			}{pkgs: nil, id: map[string]struct{}{"@middleware-b": struct{}{}}},
			want: map[string][]fun{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Position(tt.args.pkgs, tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("position() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_AddCode(t *testing.T) {
	pkg, _ := parser.ParseDir(token.NewFileSet(), "../unitTests", func(info fs.FileInfo) bool {
		//	ignore all logic check
		return true
	}, parser.ParseComments)
	
	pkgs := Position(pkg, map[string]struct{}{"@middleware-a": {}, "@middleware-b": {}, "@middleware-return": {}})
	
	type args struct {
		pkgs    map[string][]fun
		stmt    map[string]StmtParams
		replace bool
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantModify map[string][]string
	}{
		{
			name: "normal add return code",
			args: struct {
				pkgs    map[string][]fun
				stmt    map[string]StmtParams
				replace bool
			}{pkgs: pkgs, stmt: map[string]StmtParams{
				"@middleware-a": {
					Stmts: []StmtParam{
						{
							Kind: AddFuncWithoutDepends,
							Stmt: []string{
								`func(){fmt.Println("add funcStmt by addCode")}()`,
							},
							Depends: nil,
						},
						{
							Kind: AddDeferFuncStmt,
							Stmt: []string{
								`defer func(){fmt.Println("add defer func by addCode")}()`,
							},
							Depends: nil,
						},
						//{
						//	Kind: AddFuncWithVarStmt,
						//	Stmt: []string{
						//		`if 1>0 {
						//			fmt.Println("add a var func by addCode")
						//		}`,
						//	},
						//	Depends: []string{
						//		"x",
						//	},
						//},
						{
							Kind: AddFuncWithVarStmt,
							Stmt: []string{
								`if 1>0 {
									fmt.Println("add a func depend by addCode")
								}`,
							},
							FuncDepends: []string{
								"m.Round",
							},
						},
					},
					Packs: nil,
				},
				"@middleware-return": {
					Stmts: []StmtParam{
						{
							Kind: AddReturnFuncWithoutVarStmt,
							Stmt: []string{
								`if 1>0 {
									fmt.Println("add a return func by addCode")
								}`,
							},
						},
					},
					Packs: nil,
				},
			}, replace: false},
			wantErr: false,
			wantModify: map[string][]string{
				"../unitTests/test.go": []string{
					"@middleware-a",
					"@middleware-b",
					"@middleware-return",
				},
			},
		},
		{
			name: "normal add code",
			args: struct {
				pkgs    map[string][]fun
				stmt    map[string]StmtParams
				replace bool
			}{pkgs: pkgs, stmt: map[string]StmtParams{
				"@middleware-a": {
					Stmts: []StmtParam{
						{
							Kind: AddFuncWithoutDepends,
							Stmt: []string{
								`func(){fmt.Println("add by addCode")}()`,
							},
							Depends: nil,
						},
						{
							Kind: AddDeferFuncStmt,
							Stmt: []string{
								`defer func(){fmt.Println("add by addCode")}()`,
							},
							Depends: nil,
						},
					},
					Packs: nil,
				},
			}, replace: false},
			wantErr: false,
			wantModify: map[string][]string{
				"../unitTests/test.go": []string{
					"@middleware-a",
					"@middleware-b",
					"@middleware-return",
				},
			},
		},
		{
			name: "only add fun code",
			args: struct {
				pkgs    map[string][]fun
				stmt    map[string]StmtParams
				replace bool
			}{pkgs: pkgs, stmt: map[string]StmtParams{
				"@middleware-a": {
					Stmts: []StmtParam{
						{
							Kind: AddFuncWithoutDepends,
							Stmt: []string{
								`func(){fmt.Println("add by addCode")}()`,
							},
							Depends: nil,
						},
					},
					Packs: nil,
				},
			}, replace: false},
			wantErr: false,
			wantModify: map[string][]string{
				"../unitTests/test.go": []string{
					"@middleware-a",
					"@middleware-b",
					"@middleware-return",
				},
			},
		},
		{
			name: "only add defer code",
			args: struct {
				pkgs    map[string][]fun
				stmt    map[string]StmtParams
				replace bool
			}{pkgs: pkgs, stmt: map[string]StmtParams{
				"@middleware-a": {
					Stmts: []StmtParam{
						{
							Kind: AddDeferFuncStmt,
							Stmt: []string{
								`defer func(){fmt.Println("add by addCode")}()`,
							},
							Depends: nil,
						},
					},
					Packs: nil,
				},
			}, replace: false},
			wantErr: false,
			wantModify: map[string][]string{
				"../unitTests/test.go": []string{
					"@middleware-a",
					"@middleware-b",
					"@middleware-return",
				},
			},
		},
		{
			name: "add nothing codes",
			args: struct {
				pkgs    map[string][]fun
				stmt    map[string]StmtParams
				replace bool
			}{pkgs: pkgs, stmt: map[string]StmtParams{
				"@middleware-a": {},
			}, replace: false},
			wantErr: false,
			wantModify: map[string][]string{
				"../unitTests/test.go": []string{
					"@middleware-a",
					"@middleware-b",
					"@middleware-return",
				},
			},
		},
		{
			name: "add two fun codes and one defer code",
			args: struct {
				pkgs    map[string][]fun
				stmt    map[string]StmtParams
				replace bool
			}{pkgs: pkgs, stmt: map[string]StmtParams{
				"@middleware-a": {
					Stmts: []StmtParam{
						{
							Kind: AddFuncWithoutDepends,
							Stmt: []string{
								`func(){fmt.Println("add by addCode once")}()`,
								`func(){fmt.Println("add by addCode twice")}()`,
							},
						},
						{
							Kind: AddDeferFuncStmt,
							Stmt: []string{
								`defer func(){fmt.Println("add  defer by addCode")}()`,
							},
							Depends: nil,
						},
					},
					Packs: nil,
				},
				"@middleware-b": {
					Stmts: []StmtParam{
						{
							Kind: AddFuncWithoutDepends,
							Stmt: []string{
								`func(){fmt.Println("middleware-b add middleware-b by addCode once")}()`,
								`func(){fmt.Println("middleware-b add by addCode twice")}()`,
							},
						},
						{
							Kind: AddDeferFuncStmt,
							Stmt: []string{
								`defer func(){fmt.Println("middleware-b add by addCode")}()`,
							},
							Depends: nil,
						},
					},
					Packs: nil,
				},
			}, replace: false},
			wantErr: false,
			wantModify: map[string][]string{
				"../unitTests/test.go": []string{
					"@middleware-a",
					"@middleware-b",
					"@middleware-return",
				},
			},
		},
		{
			name: "add wrong defer code",
			args: struct {
				pkgs    map[string][]fun
				stmt    map[string]StmtParams
				replace bool
			}{pkgs: pkgs, stmt: map[string]StmtParams{
				"@middleware-a": {
					Stmts: []StmtParam{
						{
							Kind: AddFuncWithoutDepends,
							Stmt: []string{
								`func(){fmt.Println("add by addCode once")}()`,
								`func(){fmt.Println("add by addCode twice")}()`,
							},
						},
						{
							Kind: AddDeferFuncStmt,
							Stmt: []string{
								`fmt.Println("add by addCode"`,
							},
							Depends: nil,
						},
					},
					Packs: nil,
				},
			}, replace: false},
			wantErr:    true,
			wantModify: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := AddCode(tt.args.pkgs, tt.args.stmt, tt.args.replace)
			if (err != nil) != tt.wantErr {
				t.Errorf("addCode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(m, tt.wantModify) {
				t.Errorf("addCode() modify = %v, wantErr %v", m, tt.wantModify)
			}
		})
	}
}

func Test_extractIdFromComment(t *testing.T) {
	type args struct {
		comment string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "a valid id",
			args: struct{ comment string }{comment: `// a function comment
// @middleware-a
`},
			want: []string{"@middleware-a"},
		},
		{
			name: "two valid ids",
			args: struct{ comment string }{comment: `// a function comment
// @middleware-a
// @Middleware
`},
			want: []string{"@middleware-a", "@Middleware"},
		},
		{
			name: "inValid id",
			args: struct{ comment string }{comment: `// a function comment
// @ middleware
`},
			want: nil,
		},
		{
			name: "inValid id",
			args: struct{ comment string }{comment: `// a function comment
// @@ middleware
`},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractIdFromComment(tt.args.comment); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractIdFromComment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getIntersection(t *testing.T) {
	type args struct {
		arr []string
		ids map[string]struct{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "path test",
			args: struct {
				arr []string
				ids map[string]struct{}
			}{arr: []string{
				"@middleware-a(path:\"xx\")",
				"@middleware-b()",
				"@middleware-c",
			}, ids: map[string]struct{}{
				"@middleware-a": struct {
				}{},
			}},
			want: []string{"@middleware-a"},
		},
		{
			name: "nil test",
			args: struct {
				arr []string
				ids map[string]struct{}
			}{arr: []string{
				"@middleware-a",
				"@middleware-b",
				"@middleware-c",
			}, ids: map[string]struct{}{
				"@middleware": struct {
				}{},
			}},
			want: nil,
		},
		{
			name: "half intersection test",
			args: struct {
				arr []string
				ids map[string]struct{}
			}{arr: []string{
				"@middleware-a",
				"@middleware-b",
				"@middleware-c",
			}, ids: map[string]struct{}{
				"@middleware": struct {
				}{},
				"@middleware-a": struct {
				}{},
			}},
			want: []string{"@middleware-a"},
		},
		{
			name: "full intersection test",
			args: struct {
				arr []string
				ids map[string]struct{}
			}{arr: []string{
				"@middleware-a",
				"@middleware-b",
				"@middleware-c",
			}, ids: map[string]struct{}{
				"@middleware": struct {
				}{},
				"@middleware-a": struct {
				}{},
				"@middleware-c": struct {
				}{},
				"@middleware-b": struct {
				}{},
			}},
			want: []string{
				"@middleware-a",
				"@middleware-b",
				"@middleware-c",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getIntersection(tt.args.arr, tt.args.ids); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getIntersection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddImport(t *testing.T) {
	pkg, _ := parser.ParseDir(token.NewFileSet(), "../unitTests", func(info fs.FileInfo) bool {
		//	ignore all logic check
		return true
	}, parser.ParseComments)
	
	pkgs := Position(pkg, map[string]struct{}{"@middleware-a": {}, "@middleware-b": {}})
	
	type args struct {
		pkgs    map[string][]fun
		stmt    map[string]StmtParams
		modify  map[string][]string
		replace bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal test",
			args: struct {
				pkgs    map[string][]fun
				stmt    map[string]StmtParams
				modify  map[string][]string
				replace bool
			}{pkgs: pkgs, stmt: map[string]StmtParams{
				"@middleware-a": {
					Packs: []Pack{
						{
							Name: "f",
							Path: "\"fmt\"",
						},
					},
				},
			}, modify: map[string][]string{
				"../unitTests/test.go": []string{
					"@middleware-a",
				},
			}, replace: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddImport(tt.args.pkgs, tt.args.stmt, tt.args.modify, tt.args.replace); (err != nil) != tt.wantErr {
				t.Errorf("AddImport() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddCode(t *testing.T) {
	
	pkg, _ := parser.ParseDir(token.NewFileSet(), "../cases/insert-code-behind-variable", func(info fs.FileInfo) bool {
		//	ignore all logic check
		return true
	}, parser.ParseComments)
	
	pkgs := Position(pkg, map[string]struct{}{"@middleware-err": {}})
	
	type args struct {
		pkgs    map[string][]fun
		stmt    map[string]StmtParams
		replace bool
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]string
		wantErr bool
	}{
		{
			name: "Normal test",
			args: struct {
				pkgs    map[string][]fun
				stmt    map[string]StmtParams
				replace bool
			}{pkgs: pkgs, stmt: map[string]StmtParams{
				"@middleware-err": {
					DeclStmt: []DeclParams{
						{
							VarName: "err",
							Stmt: []string{
								`
	defer func(err error) {
		if err != nil {
			fmt.Println("err not nil")
		} else {
			fmt.Println("err is nil")
		}
	}(err)
`,
							},
						},
					},
				},
			}, replace: false},
			want: map[string][]string{
				"../cases/insert-code-behind-variable/test.go": []string{
					"@middleware-err",
				},
			},
			wantErr: false,
		},
		{
			name: "A ID with default value param",
			args: struct {
				pkgs    map[string][]fun
				stmt    map[string]StmtParams
				replace bool
			}{
				pkgs: map[string][]fun{
					"../cases/case-01.md/code.go": []fun{
						{
							originIds: []string{
								"@middleware-injection(path:\"\")",
							},
							owner: "FirstStruct",
							name:  "invokeSecondFunctionWithInjection",
							aopIds: []string{
								"@middleware-injection",
							},
						},
					},
				},
				stmt: map[string]StmtParams{
					"@middleware-injection": {
						Stmts: []StmtParam{
							{
								Kind: AddFuncWithoutDependsWithInject,
								Stmt: []string{
									`
	fmt.Println("path: %s",path)
`,
								},
								Depends: nil,
							},
						},
					},
				},
				replace: false,
			},
			want: map[string][]string{
				"../cases/case-01.md/code.go": []string{
					"@middleware-injection",
				},
			},
			wantErr: false,
		},
		{
			name: "A ID with int param",
			args: struct {
				pkgs    map[string][]fun
				stmt    map[string]StmtParams
				replace bool
			}{
				pkgs: map[string][]fun{
					"../cases/case-01.md/code.go": []fun{
						{
							originIds: []string{
								"@middleware-injection(path:100)",
							},
							owner: "FirstStruct",
							name:  "invokeThirdFunctionWithInjection",
							aopIds: []string{
								"@middleware-injection",
							},
						},
					},
				},
				stmt: map[string]StmtParams{
					"@middleware-injection": {
						Stmts: []StmtParam{
							{
								Kind: AddFuncWithoutDependsWithInject,
								Stmt: []string{
									`
	fmt.Println("path: %v",path)
`,
								},
								Depends: nil,
							},
						},
					},
				},
				replace: false,
			},
			want: map[string][]string{
				"../cases/case-01.md/code.go": []string{
					"@middleware-injection",
				},
			},
			wantErr: false,
		},
		{
			name: "A ID with two params",
			args: struct {
				pkgs    map[string][]fun
				stmt    map[string]StmtParams
				replace bool
			}{
				pkgs: map[string][]fun{
					"../cases/case-01.md/code.go": []fun{
						{
							originIds: []string{
								"@middleware-injection(path:100, name:\"a string param\")",
							},
							owner: "FirstStruct",
							name:  "invokeFourFunctionWithInjection",
							aopIds: []string{
								"@middleware-injection",
							},
						},
					},
				},
				stmt: map[string]StmtParams{
					"@middleware-injection": {
						Stmts: []StmtParam{
							{
								Kind: AddFuncWithoutDependsWithInject,
								Stmt: []string{
									`
	fmt.Println("path: %v, name: %v",path, name)
`,
								},
								Depends: nil,
							},
						},
					},
				},
				replace: false,
			},
			want: map[string][]string{
				"../cases/case-01.md/code.go": []string{
					"@middleware-injection",
				},
			},
			wantErr: false,
		},
		{
			name: "A ID with param",
			args: struct {
				pkgs    map[string][]fun
				stmt    map[string]StmtParams
				replace bool
			}{
				pkgs: map[string][]fun{
					"../cases/case-01.md/code.go": []fun{
						{
							originIds: []string{
								"@middleware-injection(path:\"/user/id\")",
							},
							owner: "FirstStruct",
							name:  "invokeFunctionWithInjection",
							aopIds: []string{
								"@middleware-injection",
							},
						},
					},
				},
				stmt: map[string]StmtParams{
					"@middleware-injection": {
						Stmts: []StmtParam{
							{
								Kind: AddFuncWithoutDependsWithInject,
								Stmt: []string{
									`
	fmt.Println("path: %s",path)
`,
								},
								Depends: nil,
							},
						},
					},
				},
				replace: false,
			},
			want: map[string][]string{
				"../cases/case-01.md/code.go": []string{
					"@middleware-injection",
				},
			},
			wantErr: false,
		},
		{
			name: "A ID with inject param",
			args: struct {
				pkgs    map[string][]fun
				stmt    map[string]StmtParams
				replace bool
			}{
				pkgs: map[string][]fun{
					"../cases/case-01.md/code.go": []fun{
						{
							originIds: []string{
								"@middleware-injection(path:100, name:\"a string param\", f:@inject)",
							},
							owner: "FirstStruct",
							name:  "invokeFiveFunctionWithInjection",
							aopIds: []string{
								"@middleware-injection",
							},
						},
					},
				},
				stmt: map[string]StmtParams{
					"@middleware-injection": {
						Stmts: []StmtParam{
							{
								Kind: AddFuncWithoutDependsWithInject,
								Stmt: []string{
									`
	fmt.Println("path: %v, name: %v inject: %v",path, name, f)
`,
								},
								Depends: nil,
							},
						},
					},
				},
				replace: false,
			},
			want: map[string][]string{
				"../cases/case-01.md/code.go": []string{
					"@middleware-injection",
				},
			},
			wantErr: false,
		},
		{
			name: "Has same variables, insert correct position test",
			args: struct {
				pkgs    map[string][]fun
				stmt    map[string]StmtParams
				replace bool
			}{pkgs: map[string][]fun{
				"../cases/insert-code-behind-variable/test.go": []fun{
					{
						owner: "FirstStruct",
						name:  "invokeSecondFunction",
						aopIds: []string{
							"@middleware-err",
						},
					},
				},
			}, stmt: map[string]StmtParams{
				"@middleware-err": {
					DeclStmt: []DeclParams{
						{
							VarName: "err",
							Stmt: []string{
								`
	defer func(err error) {
		if err != nil {
			fmt.Println("err not nil")
		} else {
			fmt.Println("err is nil")
		}
	}(err)
`,
							},
						},
					},
				},
			}, replace: false},
			want: map[string][]string{
				"../cases/insert-code-behind-variable/test.go": []string{
					"@middleware-err",
				},
			},
			wantErr: false,
		},
		{
			name: "Has two kinds function insert correct position in no receiver fun",
			args: struct {
				pkgs    map[string][]fun
				stmt    map[string]StmtParams
				replace bool
			}{pkgs: map[string][]fun{
				"../cases/insert-code-behind-variable/test.go": []fun{
					{
						owner: "",
						name:  "invokeThreeFunction",
						aopIds: []string{
							"@middleware-err",
						},
					},
				},
			}, stmt: map[string]StmtParams{
				"@middleware-err": {
					DeclStmt: []DeclParams{
						{
							VarName: "err",
							Stmt: []string{
								`
	defer func(err error) {
		if err != nil {
			fmt.Println("err not nil")
		} else {
			fmt.Println("err is nil")
		}
	}(err)
`,
							},
						},
					},
				},
			}, replace: false},
			want: map[string][]string{
				"../cases/insert-code-behind-variable/test.go": []string{
					"@middleware-err",
				},
			},
			wantErr: false,
		},
		{
			name: "Insert correct position in no receiver fun as return type",
			args: struct {
				pkgs    map[string][]fun
				stmt    map[string]StmtParams
				replace bool
			}{pkgs: map[string][]fun{
				"../cases/insert-code-behind-variable/test.go": []fun{
					{
						owner: "",
						name:  "invokeFourFunction",
						aopIds: []string{
							"@middleware-err",
						},
					},
				},
			}, stmt: map[string]StmtParams{
				"@middleware-err": {
					DeclStmt: []DeclParams{
						{
							VarName: "err",
							Stmt: []string{
								`
	defer func(err error) {
		if err != nil {
			fmt.Println("err not nil")
		} else {
			fmt.Println("err is nil")
		}
	}(err)
`,
							},
						},
					},
				},
			}, replace: false},
			want: map[string][]string{
				"../cases/insert-code-behind-variable/test.go": []string{
					"@middleware-err",
				},
			},
			wantErr: false,
		},
		{
			name: "Insert correct position in no receiver fun with return fun",
			args: struct {
				pkgs    map[string][]fun
				stmt    map[string]StmtParams
				replace bool
			}{pkgs: map[string][]fun{
				"../cases/insert-code-behind-variable/test.go": []fun{
					{
						owner: "",
						name:  "invokeFiveFunction",
						aopIds: []string{
							"@middleware-err",
						},
					},
				},
			}, stmt: map[string]StmtParams{
				"@middleware-err": {
					DeclStmt: []DeclParams{
						{
							VarName: "err",
							Stmt: []string{
								`
	defer func(err error) {
		if err != nil {
			fmt.Println("err not nil")
		} else {
			fmt.Println("err is nil")
		}
	}(err)
`,
							},
						},
					},
				},
			}, replace: false},
			want: map[string][]string{
				"../cases/insert-code-behind-variable/test.go": []string{
					"@middleware-err",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddCode(tt.args.pkgs, tt.args.stmt, tt.args.replace)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddCode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReturnWithVar(t *testing.T) {
	type args struct {
		pkgs    map[string][]fun
		stmt    map[string]StmtParams
		replace bool
	}
	
	tests := []struct {
		name    string
		args    args
		want    map[string][]string
		wantErr bool
	}{
		{
			name: "Insert code in return func with specify variables in params list",
			args: struct {
				pkgs    map[string][]fun
				stmt    map[string]StmtParams
				replace bool
			}{pkgs: map[string][]fun{
				"../cases/insert-return-func-with-var/test.go": []fun{
					{
						owner: "",
						name:  "invokeFourFunction",
						aopIds: []string{
							"@middleware-return-err",
						},
					},
				},
			}, stmt: map[string]StmtParams{
				"@middleware-return-err": {
					Stmts: []StmtParam{
						
						{
							Kind: AddReturnFuncWithVarStmt,
							Stmt: []string{
								`if err != nil {
										fmt.Println("err in this return func not nil")
									} else {
										fmt.Println("err in this return func is nil")
									}
								`,
							},
							Depends: []string{"err"},
						},
					},
				},
			}, replace: false},
			want: map[string][]string{
				"../cases/insert-return-func-with-var/test.go": []string{
					"@middleware-return-err",
				},
			},
			wantErr: false,
		},
		{
			name: "Insert code in return func with specify variables in body",
			args: struct {
				pkgs    map[string][]fun
				stmt    map[string]StmtParams
				replace bool
			}{pkgs: map[string][]fun{
				"../cases/insert-return-func-with-var/test.go": []fun{
					{
						owner: "",
						name:  "invokeFiveFunction",
						aopIds: []string{
							"@middleware-return-err-body",
						},
					},
				},
			}, stmt: map[string]StmtParams{
				"@middleware-return-err-body": {
					Stmts: []StmtParam{
						
						{
							Kind: AddReturnFuncWithVarStmt,
							Stmt: []string{
								`if err != nil {
										fmt.Println("err in this return func not nil")
									} else {
										fmt.Println("err in this return func is nil")
									}
								`,
							},
							Depends: []string{"err"},
						},
					},
				},
			}, replace: false},
			want: map[string][]string{
				"../cases/insert-return-func-with-var/test.go": []string{
					"@middleware-return-err-body",
				},
			},
			wantErr: false,
		},
		{
			name: "Use func as variable",
			args: struct {
				pkgs    map[string][]fun
				stmt    map[string]StmtParams
				replace bool
			}{pkgs: map[string][]fun{
				"../cases/insert-return-func-with-var/test.go": []fun{
					{
						owner: "",
						name:  "invokeSevenFunction",
						aopIds: []string{
							"@middleware-func-var-err",
						},
					},
				},
			}, stmt: map[string]StmtParams{
				"@middleware-func-var-err": {
					Stmts: []StmtParam{
						
						{
							Kind: AddReturnFuncWithVarStmt,
							Stmt: []string{
								`if e != nil {
										fmt.Println("e func in params")
									} else {
										fmt.Println("e func in params is nil")
									}
								`,
							},
							Depends: []string{"e"},
						},
					},
				},
			}, replace: false},
			want: map[string][]string{
				"../cases/insert-return-func-with-var/test.go": []string{
					"@middleware-func-var-err",
				},
			},
			wantErr: false,
		},
		{
			name: "Mixed scene",
			args: struct {
				pkgs    map[string][]fun
				stmt    map[string]StmtParams
				replace bool
			}{pkgs: map[string][]fun{
				"../cases/insert-return-func-with-var/test.go": []fun{
					{
						originIds: []string{
							"@middleware-err",
						},
						owner: "",
						name:  "invokeSixFunction",
						aopIds: []string{
							"@middleware-err",
						},
					},
				},
			}, stmt: map[string]StmtParams{
				"@middleware-err": {
					Stmts: []StmtParam{
						
						{
							Kind: AddReturnFuncWithVarStmt,
							Stmt: []string{
								`if err != nil {
										fmt.Println("err in params")
									} else {
										fmt.Println("err in params is nil")
									}
								`,
							},
							Depends: []string{"err"},
						},
					},
				},
			}, replace: false},
			want: map[string][]string{
				"../cases/insert-return-func-with-var/test.go": []string{
					"@middleware-err",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddCode(tt.args.pkgs, tt.args.stmt, tt.args.replace)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestReturnWithVar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestReturnWithVar() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInjectionWithFuncDepend(t *testing.T) {
	type args struct {
		pkgs    map[string][]fun
		stmt    map[string]StmtParams
		replace bool
	}
	
	tests := []struct {
		name    string
		args    args
		want    map[string][]string
		wantErr bool
	}{
		{
			name: "normal test",
			args: struct {
				pkgs    map[string][]fun
				stmt    map[string]StmtParams
				replace bool
			}{pkgs: map[string][]fun{
				"../cases/case-02/code.go": []fun{
					{
						originIds: []string{
							"@middleware-func-depend",
						},
						owner: "FirstStruct",
						name:  "addWithFuncDependWithInjection",
						aopIds: []string{
							"@middleware-func-depend",
						},
					},
				},
			}, stmt: map[string]StmtParams{
				"@middleware-func-depend": {
					Stmts: []StmtParam{
						
						{
							Kind: AddFuncWithVarStmt,
							Stmt: []string{
								`go fmt.Println(__varName__)`,
							},
							FuncDepends: []string{"math.Round"},
						},
					},
				},
			}, replace: false},
			want: map[string][]string{
				"../cases/case-02/code.go": []string{
					"@middleware-func-depend",
				},
			},
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddCode(tt.args.pkgs, tt.args.stmt, tt.args.replace)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestReturnWithVar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestReturnWithVar() got = %v, want %v", got, tt.want)
			}
		})
	}
}
