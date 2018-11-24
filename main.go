package main

import (
	"TcpReConnDemo/client"
	"TcpReConnDemo/server"
)

const isServer = true

func main() {
	if isServer {
		server.Start()
	} else {
		client.Conn()
	}

}
