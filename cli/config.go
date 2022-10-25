package main

import (
	"github.com/BurntSushi/toml"
	"github.com/runways/goAOP/aops"
)

type Config struct {
	MiddWare    []middleWare `toml:"middleware"`
	MiddWareMap map[string]aops.StmtParams
}

type middleWare struct {
	ID        string   `toml:"id"`
	FuncStmt  []string `toml:"funcStmt"`
	DeferStmt []string `toml:"deferStmt"`
	Package   []pack   `toml:"package"`
}

type pack struct {
	Name string `toml:"name"`
	Path string `toml:"path"`
}

func parseConfig(file string) (c Config, err error) {
	_, err = toml.DecodeFile(file, &c)
	if err != nil {
		return
	}
	
	mwm := make(map[string]aops.StmtParams, len(c.MiddWare))
	for _, m := range c.MiddWare {
		var p []aops.Pack
		for _, _p := range m.Package {
			p = append(p, aops.Pack{
				Name: _p.Name,
				Path: _p.Path,
			})
		}
		mwm[m.ID] = aops.StmtParams{
			FunStmt:   m.FuncStmt,
			DeferStmt: m.DeferStmt,
			Packs:     p,
		}
	}
	
	c.MiddWareMap = mwm
	return
}
