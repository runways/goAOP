package main

import "flag"

var (
	// main operation modes
	dir     = flag.String("dir", "", "The source code file dir path")
	replace = flag.Bool("replace", true, "Replace source code file or not, default is true")
	conf    = flag.String("config", "aop.toml", "The runtime config, default is aop.toml")
	// debug operation mode
	debug = flag.Bool("debug", false, "Enable / Disable debug output")
)
