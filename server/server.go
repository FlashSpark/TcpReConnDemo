package server

import (
	"TcpReConnDemo/message"
	"fmt"
	"log"
	"net"
)

const (
	addr     = "0.0.0.0"
	port     = "8080"
)

type Actor struct {
}

func (a *Actor) resp(req string) string {
	switch req {
	case "hello":
	case "hi":
		return "hi, i'm server. "
	default:
		return "i don't quit understand. "
	}

	return "something wrong. "
}

type Server struct {
	listenAddr string
	port       string

	pool  Pool
	actor Actor
	rw    *message.DataRW
}

type Pool struct {
	cs []net.Conn
}

func (p *Pool) add(c net.Conn) {
}

func (s *Server) Start() {
	s.listenAddr = addr
	s.port = port

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr+":"+port)
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
		s.rw = message.DataRWIns(conn)

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(c net.Conn) {

	defer func() {
		err := c.Close()
		if err != nil {
			fmt.Println("connection close error. ")
		}
	}()

	// return hello
	var buffer = []byte("You are welcome. I'm server.")

	n, err := c.Write(buffer)

	if err != nil {
		fmt.Println("Write error:", err)
	}

	fmt.Println("send bytes:", n)

	for {
		m, err := s.rw.ReadMsg()
		if err != nil {
			fmt.Println("error:", err.Error())
			return
		}

		fmt.Println(n, " bytes received. ")

		fmt.Println("message: ", m.Payload)
		resp := s.actor.resp("this")

		fmt.Println("resp message :", resp)

		n, err = c.Write([]byte(resp))
		if err != nil {
			fmt.Println("error:", err.Error())
			return
		}

		fmt.Println("write success. ")
	}
}

func Start() {
	s := Server{}
	s.Start()
}
