package main

import (
	"flag"
	"fmt"
	"github.com/runways/goAOP/aops"
	"os"
)

func main() {

	flag.Parse()

	check()

	c, err := parseConfig(*conf)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	if *debug {
		outputConfig(c)
	}

	pkgs, err := aops.ParseDir(*dir, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	pkgMap := aops.Position(pkgs, "@middleware-a")

	aops.AddCode(pkgMap, c.MiddWareMap["@middleware-a"].FuncStmt, c.MiddWareMap["@middleware-a"].DeferStmt, *replace)
}

func check() {
	if *dir == "" {
		fmt.Println("dir is null, please specify source code dir path with -dir")
		os.Exit(-1)
	}

}

func outputConfig(c Config) {
	fmt.Printf("%+v \n", c.MiddWare)
}
