package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"webx/core"
)

type Status int

const (
	Online Status = iota
	Offline
	Retrying
	NotStarted
	Unknown
)

func (s Status) ToString() string {
	switch s {
	case Online:
		return "online"
	case Offline:
		return "offline"
	case Retrying:
		return "retrying"
	case NotStarted:
		return "NotStarted"
	default:
		return "Unknown"
	}
}

type Connection struct {
	Status   Status
	Nretries int
	Attempt  int
	Address  string
	HcPath   string
}

type Connections map[string]*Connection

func PingConnections(conns *Connections, conf *core.Config) {
	var pingInterval int
	if pingInterval = conf.Global.PingInterval; pingInterval <= 0 {
		pingInterval = 10
	}

	ticker := time.NewTicker(time.Duration(pingInterval) * time.Second)

	for {
		<-ticker.C
		for _, conn := range *conns {
			path := fmt.Sprintf("http://%s%s", conn.Address, conn.HcPath)
			resp, err := http.Get(path)
			if err != nil || resp.StatusCode != 200 {
				conn.Status = Offline
				fmt.Printf("service: %s - %s\n", path, Offline.ToString())
			} else {
				conn.Status = Online
				fmt.Printf("service: %s - %s\n", path, Online.ToString())
			}
		}
	}
}

func main() {
	confPath := "conf.toml"
	conf, err := core.ParseConfig(confPath)
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("config: OK")
	}

	conns := make(Connections)
	for _, rule := range conf.PassRules {
		conns[rule.Tpath] = &Connection{
			Status:   NotStarted,
			Nretries: 5,
			Attempt:  0,
			Address:  fmt.Sprintf("%s:%d", rule.Thost, rule.Tport),
			HcPath:   rule.HcPath,
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go PingConnections(&conns, conf)

	conf.Print()
	wg.Wait()
}
