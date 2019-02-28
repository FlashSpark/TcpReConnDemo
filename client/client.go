package client

import (
	"TcpReConnDemo/message"
	server2 "TcpReConnDemo/server"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

const (
	server  = "192.168.31.20"
	msgSize = 50
)

type Client struct {
	rw *message.DataRW

	in chan message.Msg

	lock sync.Mutex
}

func InsOfClient() Client {
	c := Client{}
	c.in = make(chan message.Msg, msgSize)

	return c
}

func (c *Client) conn() {
	addr := server + ":8080"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("dial failed:", err)
		os.Exit(1)
	}

	c.rw = message.DataRWIns(conn)

	go c.dataDeal()
	go c.readLoop()
}

func (c *Client) close() {
	defer c.lock.Unlock()
	c.lock.Lock()

	c.rw.Close()
}

// conn to server
func (c *Client) sendMsg(msg string) {
	err := message.Send(c.rw, server2.NormalMsg, msg)
	if err != nil {
		fmt.Println("msg send error :", err.Error())
	}
}

func (c *Client) readLoop() {
	for {
		msg, err := c.rw.ReadMsg()
		if err != nil {
			fmt.Println("client read error. ", err.Error())
			break
		}

		c.in <- msg
	}

	c.close()
}

func (c *Client) dataDeal() {
	for {
		select {
		case msg := <-c.in:
			c.display(msg)
		}
	}
}

func (c *Client) display(msg message.Msg) {
	fmt.Println("server resp code :", msg.Code)
	var content string
	err := msg.Decode(&content)
	if err != nil {
		fmt.Println("msg decode error. ", err.Error())
	}
	fmt.Println("info :", content)
}

func Start() {
	for i:= 0 ;i < 60000;i ++ {
		c := InsOfClient()
		c.conn()
		time.Sleep(time.Second)
	}


	//for {
	//	inputReader := bufio.NewReader(os.Stdin)
	//	fmt.Printf("send to server: ")
	//	input, err := inputReader.ReadString('\n')
	//	if err == nil {
	//		fmt.Printf("The input was: %s", input)
	//	}
	//
	//	if strings.Compare(input, "exit") == 0 {
	//		break
	//	}
	//
	//	c.sendMsg(input)
	//}

	for {
		time.Sleep(time.Second)
	}
}
