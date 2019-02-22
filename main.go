package main

import (
	"TcpReConnDemo/client"
	"TcpReConnDemo/server"
	"fmt"
	"os"
)

const (
	typeServer = "server"
	typeClient = "client"
)

// set client or server
func main() {
	switch cmd() {
	case typeServer:
		server.Start()
	case typeClient:
		client.Start()
	default:
		fmt.Println("error of cmd. no params ")
	}
}

// get cmd from os
func cmd() string {
	p := os.Args
	if len(p) < 2 {
		return ""
	}

	return p[1]
}
