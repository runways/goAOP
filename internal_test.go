package goAOP

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

func Test_position(t *testing.T) {
	
	pkg, _ := parser.ParseDir(token.NewFileSet(), "./unitTests", func(info fs.FileInfo) bool {
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
				"unitTests/test.go": []string{
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
			if got := position(tt.args.pkgs, tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("position() = %v, want %v", got, tt.want)
			}
		})
	}
}
