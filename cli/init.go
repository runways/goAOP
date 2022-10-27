package main

import "flag"

var (
	// main operation modes
	dir     = flag.String("dir", "", "The source code file dir path")
	replace = flag.Bool("replace", true, "Replace source code file or not")
	conf    = flag.String("config", "aop.toml", "The runtime config")
	// debug operation mode
	debug = flag.Bool("debug", false, "Enable / Disable debug output")
)
