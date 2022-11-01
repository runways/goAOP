package main

import (
	"fmt"
	"github.com/runways/goAOP/aops"
	"os"
	"reflect"
	"testing"
)

func Test_parseConfig(t *testing.T) {

	subFile, _ := os.CreateTemp("", "")
	os.WriteFile(subFile.Name(), []byte(`
[[middleware]]
    id="@middleware-b"
	[[middleware.Stmt]]
	kind="add-return-func-with-var"
    code=["""func(){
                log.Println("before")
            }()"""]
	depend=["err"]
	[[middleware.Stmt]]
	kind="add-return-func-without-var"
    code =["""func(){
                log.Println("before")
            }()"""]
	depend=["str"]
`), 0777)

	conf := fmt.Sprintf(`
include=["%s"]
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

`, subFile.Name())
	f, _ := os.CreateTemp("", "")
	os.WriteFile(f.Name(), []byte(conf), 0777)

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
				Include: []string{subFile.Name()},
				MidWare: []middleWare{
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
				MidWareMap: map[string]aops.StmtParams{
					"@middleware-b": aops.StmtParams{
						DeclStmt: nil,
						Stmts: []aops.StmtParam{
							aops.StmtParam{
								Kind: aops.AddReturnFuncWithVarStmt,
								Stmt: []string{
									`func(){
                log.Println("before")
            }()`,
								},
								Depends: []string{"err"},
							},
							aops.StmtParam{
								Kind: aops.AddReturnFuncWithoutVarStmt,
								Stmt: []string{
									`func(){
                log.Println("before")
            }()`,
								},
								Depends: nil,
							},
						},
						Packs: nil,
					},
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
