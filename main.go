package main

import (
	"fmt"
	"log"
	"net"
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

	conns := make(map[string]net.Conn)

	for _, rule := range conf.PassRules {
		addr := fmt.Sprintf("%s:%d", rule.Thost, rule.Tport)
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Printf("error: failed to establish connection with %s; err: %v\n", addr, err)
			continue
		}
		defer conn.Close()

		conns[rule.Spath] = conn
	}

	conf.Print()
}
