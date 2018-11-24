package client

import (
	"fmt"
	"net"
	"os"
)

func Conn() {

	conn, err := net.Dial("tcp", "192.168.5.6:8080")
	if err != nil {
		fmt.Println("dial failed:", err)
		os.Exit(1)
	}
	defer conn.Close()

	buffer := make([]byte, 512)

	n, err2 := conn.Read(buffer)
	if err2 != nil {
		fmt.Println("Read failed:", err2)
		return
	}

	fmt.Println("count:", n, "msg:", string(buffer))
}
