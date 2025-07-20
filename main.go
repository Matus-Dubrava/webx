package main

import (
	"fmt"
	"log"
	"webx/core"
)

func main() {
	confPath := "conf.toml"
	conf, err := core.ParseConfig(confPath)
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("config: OK")
	}

	conf.Print()
}
