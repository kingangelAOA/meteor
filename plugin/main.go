package main

import (
	"flag"
	"fmt"
	"plugin/shared"
)

func main() {
	var name string
	var path string
	flag.StringVar(&name, "name", "demo", "Just for demo")
	flag.StringVar(&path, "path", "demo", "Just for demo")
	flag.Parse()
	err := shared.BuildPlugin(name, path)
	if err != nil {
		fmt.Println(err.Error())
	}
}
