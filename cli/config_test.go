package main

import (
	"github.com/runways/goAOP/aops"
	"os"
	"reflect"
	"testing"
)

func Test_parseConfig(t *testing.T) {
	
	f, _ := os.CreateTemp("", "")
	os.WriteFile(f.Name(), []byte(`
[[middleware]]
    id="@middleware-a"
	[[middleware.Stmt]]
	kind="add-func-without-depends"
    code=["""func(){
                log.Println("before")
            }()"""]
	depend=["err"]
	[[middleware.Stmt]]
	kind="add-func-with-var-depend"
    code =["""func(){
                log.Println("before")
            }()"""]
	depend=["str"]
    [[middleware.package]]
        name = "log"
        path = """"github.com/sirupsen/logrus""""

`), 0777)
	
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		wantC   Config
		wantErr bool
	}{
		{
			name: "Normal parse",
			args: struct{ file string }{file: f.Name()},
			wantC: Config{
				MiddWare: []middleWare{
					{
						ID: "@middleware-a",
						Package: []pack{
							{Name: "log", Path: "\"github.com/sirupsen/logrus\""},
						},
						Stmt: []Stmt{
							{
								Kind: "add-func-without-depends",
								Code: []string{`func(){
                log.Println("before")
            }()`},
								Depend: []string{"err"},
							},
							{
								Kind: "add-func-with-var-depend",
								Code: []string{`func(){
                log.Println("before")
            }()`},
								Depend: []string{"str"},
							},
						},
					},
				},
				MiddWareMap: map[string]aops.StmtParams{
					"@middleware-a": aops.StmtParams{
						DeclStmt: nil,
						Stmts: []aops.StmtParam{
							aops.StmtParam{
								Kind: aops.AddFuncWithoutDepends,
								Stmt: []string{
									`func(){
                log.Println("before")
            }()`,
								},
								Depends: nil,
							},
							aops.StmtParam{
								Kind: aops.AddFuncWithVarStmt,
								Stmt: []string{
									`func(){
                log.Println("before")
            }()`,
								},
								Depends: []string{"str"},
							},
						},
						Packs: []aops.Pack{
							aops.Pack{
								Name: "log",
								Path: "\"github.com/sirupsen/logrus\"",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := parseConfig(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("parseConfig() gotC = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}
