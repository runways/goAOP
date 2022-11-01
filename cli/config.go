package main

import (
	"github.com/BurntSushi/toml"
	"github.com/runways/goAOP/aops"
	"strings"
)

type Config struct {
	MiddWare    []middleWare `toml:"middleware"`
	MiddWareMap map[string]aops.StmtParams
}

type middleWare struct {
	ID      string `toml:"id"`
	Stmt    []Stmt `toml:"Stmt"`
	Package []pack `toml:"package"`
}

type pack struct {
	Name string `toml:"name"`
	Path string `toml:"path"`
}

// Stmt save all the code will injection to source code
// ID is the middleware id, should match with comment in function.
// Kind is the middleware type.Valid values declare in the `aops/const.go`.
// Code is a string array, save the code will injection to source code.
// Depend is a string array, save the injection conditions. Now only support
// signal variable. No need type variable type.
type Stmt struct {
	//ID     string   `toml:"id"`
	Kind   string   `toml:"kind"`
	Code   []string `toml:"code,omitempty"`
	Depend []string `toml:"depend"`
}

func parseConfig(file string) (c Config, err error) {
	_, err = toml.DecodeFile(file, &c)
	if err != nil {
		return
	}
	
	mwm := make(map[string]aops.StmtParams, len(c.MiddWare))
	for _, m := range c.MiddWare {
		var p []aops.Pack
		var stmtBlock []aops.StmtParam
		for _, _p := range m.Package {
			p = append(p, aops.Pack{
				Name: strings.TrimSpace(_p.Name),
				Path: strings.TrimSpace(_p.Path),
			})
		}
		
		for _, s := range m.Stmt {
			switch strings.TrimSpace(strings.ToLower(s.Kind)) {
			case aops.AddFuncWithoutDependsStr:
				stmtBlock = append(stmtBlock, aops.StmtParam{
					Kind:    aops.AddFuncWithoutDepends,
					Stmt:    s.Code,
					Depends: nil,
				})
			case aops.AddFuncWithVarStmtStr:
				stmtBlock = append(stmtBlock, aops.StmtParam{
					Kind:    aops.AddFuncWithVarStmt,
					Stmt:    s.Code,
					Depends: s.Depend,
				})
			case aops.AddDeferFuncStmtStr:
				stmtBlock = append(stmtBlock, aops.StmtParam{
					Kind:    aops.AddDeferFuncStmt,
					Stmt:    s.Code,
					Depends: nil,
				})
			case aops.AddDeferFuncWithVarStmtStr:
				stmtBlock = append(stmtBlock, aops.StmtParam{
					Kind:    aops.AddDeferFuncWithVarStmt,
					Stmt:    s.Code,
					Depends: s.Depend,
				})
			case aops.AddReturnFuncWithoutVarStmtStr:
				stmtBlock = append(stmtBlock, aops.StmtParam{
					Kind:    aops.AddReturnFuncWithoutVarStmt,
					Stmt:    s.Code,
					Depends: nil,
				})
			case aops.AddReturnFuncWithVarStmtStr:
				stmtBlock = append(stmtBlock, aops.StmtParam{
					Kind:    aops.AddReturnFuncWithVarStmt,
					Stmt:    s.Code,
					Depends: s.Depend,
				})
			}
		}
		
		mwm[m.ID] = aops.StmtParams{
			Stmts: stmtBlock,
			Packs: p,
		}
	}
	
	c.MiddWareMap = mwm
	return
}
