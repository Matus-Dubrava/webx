package main

import (
	"fmt"
	"log"
	"net"
	"strings"
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

type ConnInfo struct {
	Status   Status
	Nretries int
	Attempt  int
	Address  string
	Path     string
	HcPath   string
}

type Connections map[string]*ConnInfo

func PingConnections(conns *Connections, conf *core.Config) {
	var pingInterval int
	if pingInterval = conf.Global.PingInterval; pingInterval <= 0 {
		pingInterval = 10
	}

	ticker := time.NewTicker(time.Duration(pingInterval) * time.Second)

	for {
		<-ticker.C
		for _, connInfo := range *conns {
			conn, err := net.Dial("tcp", connInfo.Address)
			if err != nil {
				connInfo.Status = Offline
				fmt.Printf("conn %s: offline (%v)\n", connInfo.Address, err)
				continue
			}

			req := fmt.Sprintf("GET %s HTTP/1.1\r\nContent-Length: 0\r\n\r\n", fmt.Sprintf("%s%s", connInfo.Path, connInfo.HcPath))
			req = strings.ReplaceAll(req, "//", "/")
			// fmt.Printf("sending: %s\n", req)
			_, err = conn.Write([]byte(req))
			if err != nil {
				connInfo.Status = Offline
				fmt.Printf("conn %s: offline (%v)\n", connInfo.Address, err)
				continue
			}
			// fmt.Printf("ping: sent %d bytes to upstream\n", nsent)

			buf := make([]byte, 1024)
			_, err = conn.Read(buf)
			if err != nil {
				connInfo.Status = Offline
				fmt.Printf("conn %s: offline (%v)\n", connInfo.Address, err)
				continue
			}
			// fmt.Printf("ping: recevied (%d bytes) %s\n", nread, buf)

			resp, err := core.ParseHTTPResponse(string(buf))
			if err != nil {
				connInfo.Status = Offline
				fmt.Printf("conn %s: offline (%v)\n", connInfo.Address, err)
				continue
			}
			// fmt.Printf("ping: parsed resp: %s\n", resp.ToString())

			if resp.StatusCode == 200 {
				connInfo.Status = Online
				fmt.Printf("conn %s: online\n", connInfo.Address)
			} else {
				connInfo.Status = Offline
				fmt.Printf("conn %s: offline (status code %d)\n", connInfo.Address, resp.StatusCode)
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
		conns[rule.Tpath] = &ConnInfo{
			Status:   NotStarted,
			Nretries: 5,
			Attempt:  0,
			Address:  fmt.Sprintf("%s:%d", rule.Thost, rule.Tport),
			HcPath:   rule.HcPath,
			Path:     rule.Tpath,
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go PingConnections(&conns, conf)

	listenerAddr := fmt.Sprintf("%s:%d", conf.Global.Listener_host, conf.Global.Listener_port)
	listener, err := net.Listen("tcp", listenerAddr)
	if err != nil {
		log.Fatalf("err: failed to start listening; %v\n", err)
	} else {
		fmt.Println("waiting for connections...")
	}

	_, err = listener.Accept()
	if err != nil {
		fmt.Println(err)
	}

	conf.Print()
	wg.Wait()
}
