package client

import (
	"TcpReConnDemo/message"
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	server   = "192.168.31.20"
	buffSize = 512
)

type Client struct {
	rw *message.DataRW
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

	c.rw = message.DataRWIns(conn)

	m, err2 := c.rw.ReadMsg()
	if err2 != nil {
		fmt.Println("Read failed:", err2)
		return
	}

	fmt.Println("count:", m.Size, "message:", m.String())
}

// conn to server
func (c *Client) sendMsg(msg string) {

}

func Start() {
	c := Client{}
	c.conn()

	for {
		inputReader := bufio.NewReader(os.Stdin)
		fmt.Printf("send to server: ")
		input, err := inputReader.ReadString('\n')
		if err == nil {
			fmt.Printf("The input was: %s", input)
		}

		if strings.Compare(input, "exit") == 0 {
			break
		}

		c.sendMsg(input)
	}
}
