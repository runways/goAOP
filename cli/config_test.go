package main

import (
	"os"
	"reflect"
	"testing"
)

func Test_parseConfig(t *testing.T) {

	f, _ := os.CreateTemp("", "")
	os.WriteFile(f.Name(), []byte(`
[[middleware]]
	id="@middleware-a"
	funcStmt=["""func(){
			fmt.Println("")
		}()
	"""]
	deferStmt=["""defer func(){
			fmt.Println("")
		}()
	"""]
	[middleware.package]
		name = "f"
		path = "fmt"
[[middleware]]
	id="@middleware-b"
	funcStmt=["""func(){
			fmt.Println("b")
		}()
	"""]
	deferStmt=["""defer func(){
			fmt.Println("b")
		}()
	"""]
	[middleware.package]
		name = "f1"
		path = "fmt"


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
			wantC: Config(struct {
				MiddWare    []middleWare
				MiddWareMap map[string]middleWare
			}{MiddWare: []middleWare{
				{
					ID: "@middleware-a",
					FuncStmt: []string{
						`func(){
			fmt.Println("")
		}()
	`,
					},
					DeferStmt: []string{
						`defer func(){
			fmt.Println("")
		}()
	`,
					},
					Package: pack(struct {
						Name string
						Path string
					}{Name: "f", Path: "fmt"}),
				},
				{
					ID: "@middleware-b",
					FuncStmt: []string{
						`func(){
			fmt.Println("b")
		}()
	`,
					},
					DeferStmt: []string{
						`defer func(){
			fmt.Println("b")
		}()
	`,
					},
					Package: pack(struct {
						Name string
						Path string
					}{Name: "f1", Path: "fmt"}),
				},
			},
				MiddWareMap: map[string]middleWare{
					"@middleware-a": {
						ID: "@middleware-a",
						FuncStmt: []string{
							`func(){
			fmt.Println("")
		}()
	`,
						},
						DeferStmt: []string{
							`defer func(){
			fmt.Println("")
		}()
	`,
						},
						Package: pack(struct {
							Name string
							Path string
						}{Name: "f", Path: "fmt"}),
					},
					"@middleware-b": {
						ID: "@middleware-b",
						FuncStmt: []string{
							`func(){
			fmt.Println("b")
		}()
	`,
						},
						DeferStmt: []string{
							`defer func(){
			fmt.Println("b")
		}()
	`,
						},
						Package: pack(struct {
							Name string
							Path string
						}{Name: "f1", Path: "fmt"}),
					},
				}}),
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
