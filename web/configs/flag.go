package configs

import "flag"

var ConfigPath = flag.String("configPath", "config.yml", "config path")
var PluginRootPath = flag.String("PluginRootPath", "./plugin", "go plugin path")
var Gcflags = flag.String("gcflags", "all=-N -l", "gcflags")

func FlagInit() {
	flag.Parse()
}
