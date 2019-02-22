package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	server = "192.168.31.20"
)

type Client struct {
}

func (c *Client) conn() {
	addr := server + ":8080"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("dial failed:", err)
		os.Exit(1)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("close error. ")
		}
	}()

	buffer := make([]byte, 512)

	n, err2 := conn.Read(buffer)
	if err2 != nil {
		fmt.Println("Read failed:", err2)
		return
	}

	fmt.Println("count:", n, "msg:", string(buffer))
}

// conn to server
func (c *Client) sendMsg(msg string) {

}

func Start() {
	c := Client{}
	c.conn()

	inputReader := bufio.NewReader(os.Stdin)
	fmt.Printf("send to server: ")
	input, err := inputReader.ReadString('\n')
	if err == nil {
		fmt.Printf("The input was: %s", input)
	}

	c.sendMsg(input)
}
