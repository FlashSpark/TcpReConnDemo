package server

import (
	"fmt"
	"log"
	"net"
)

func Start() {

	addr := "0.0.0.0:8080"

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)

	if err != nil {
		log.Fatalf("net.ResovleTCPAddr fail:%s", addr)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("listen %s fail: %s", addr, err)
	} else {

		log.Println("rpc listening", addr)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept error:", err)
			continue
		}

		fmt.Println("client connect in: ", conn.RemoteAddr())

		go handleConnection(conn)

	}

}

func handleConnection(conn net.Conn) {

	var buffer = []byte("You are welcome. I'm server.")

	n, err := conn.Write(buffer)

	if err != nil {
		fmt.Println("Write error:", err)
	}

	fmt.Println("send bytes:", n)
}
