package configs

import "flag"

var ConfigPath = flag.String("configPath", "config.yml", "config path")

func FlagInit()  {
	flag.Parse()
}