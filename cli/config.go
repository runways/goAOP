package main

import "github.com/BurntSushi/toml"

type Config struct {
	MiddWare    []middleWare `toml:"middleware"`
	MiddWareMap map[string]middleWare
}

type middleWare struct {
	ID        string   `toml:"id"`
	FuncStmt  []string `toml:"funcStmt"`
	DeferStmt []string `toml:"deferStmt"`
	Package   pack     `toml:"package"`
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

	mwm := make(map[string]middleWare, len(c.MiddWare))
	for _, m := range c.MiddWare {
		mwm[m.ID] = m
	}

	c.MiddWareMap = mwm
	return
}
