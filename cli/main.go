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
	
	if *debug {
		fmt.Println("Will handler bellow files:")
		for key := range pkgs {
			fmt.Println(key)
		}
		fmt.Println("=======>")
	}
	aopMap := make(map[string]struct{})
	for name := range c.MiddWareMap {
		aopMap[name] = struct{}{}
	}
	
	pkgMap := aops.Position(pkgs, aopMap)
	
	modify, err := aops.AddCode(pkgMap, c.MiddWareMap, *replace)
	
	err = aops.AddImport(pkgMap, c.MiddWareMap, modify, *replace)
	if err != nil {
		fmt.Println("FAILED")
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	
	fmt.Println("// SUCCESS")
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
