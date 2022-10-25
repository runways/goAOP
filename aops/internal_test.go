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
		id   string
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "valid position",
			args: struct {
				pkgs map[string]*ast.Package
				id   string
			}{pkgs: pkg, id: "@middleware-a"},
			want: map[string][]string{
				"../unitTests/test.go": []string{
					"InvokeFirstFunction",
				},
			},
		},
		{
			name: "invalid position",
			args: struct {
				pkgs map[string]*ast.Package
				id   string
			}{pkgs: nil, id: "@middleware-a"},
			want: map[string][]string{},
		},
		{
			name: "invalid position",
			args: struct {
				pkgs map[string]*ast.Package
				id   string
			}{pkgs: pkg, id: "@middleware-b"},
			want: map[string][]string{},
		},
		{
			name: "invalid position",
			args: struct {
				pkgs map[string]*ast.Package
				id   string
			}{pkgs: nil, id: "@middleware-b"},
			want: map[string][]string{},
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
	
	pkgs := Position(pkg, "@middleware-a")
	
	type args struct {
		pkgs      map[string][]string
		funStmt   []string
		deferStmt []string
		replace   bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal add code",
			args: struct {
				pkgs      map[string][]string
				funStmt   []string
				deferStmt []string
				replace   bool
			}{pkgs: pkgs, funStmt: []string{
				`func(){fmt.Println("add by addCode")}()`,
			}, deferStmt: []string{
				`defer func(){fmt.Println("add by addCode")}()`,
			}, replace: false},
			wantErr: false,
		},
		{
			name: "only add fun code",
			args: struct {
				pkgs      map[string][]string
				funStmt   []string
				deferStmt []string
				replace   bool
			}{pkgs: pkgs, funStmt: []string{
				`func(){fmt.Println("add by addCode")}()`,
			}, deferStmt: nil, replace: false},
			wantErr: false,
		},
		{
			name: "only add defer code",
			args: struct {
				pkgs      map[string][]string
				funStmt   []string
				deferStmt []string
				replace   bool
			}{pkgs: pkgs, funStmt: nil, deferStmt: []string{
				`defer func(){fmt.Println("add by addCode")}()`,
			}, replace: false},
			wantErr: false,
		},
		{
			name: "add nothing codes",
			args: struct {
				pkgs      map[string][]string
				funStmt   []string
				deferStmt []string
				replace   bool
			}{pkgs: pkgs, funStmt: nil, deferStmt: nil, replace: false},
			wantErr: false,
		},
		{
			name: "add two fun codes and one defer code",
			args: struct {
				pkgs      map[string][]string
				funStmt   []string
				deferStmt []string
				replace   bool
			}{pkgs: pkgs, funStmt: []string{
				`func(){fmt.Println("add by addCode once")}()`,
				`func(){fmt.Println("add by addCode twice")}()`,
			}, deferStmt: []string{
				`defer func(){fmt.Println("add by addCode")}()`,
			}, replace: false},
			wantErr: false,
		},
		{
			name: "add wrong defer code",
			args: struct {
				pkgs      map[string][]string
				funStmt   []string
				deferStmt []string
				replace   bool
			}{pkgs: pkgs, funStmt: []string{
				`func(){fmt.Println("add by addCode once")}()`,
				`func(){fmt.Println("add by addCode twice")}()`,
			}, deferStmt: []string{
				`fmt.Println("add by addCode"`,
			}, replace: false},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddCode(tt.args.pkgs, tt.args.funStmt, tt.args.deferStmt, tt.args.replace); (err != nil) != tt.wantErr {
				t.Errorf("addCode() error = %v, wantErr %v", err, tt.wantErr)
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
