package server

import (
	"TcpReConnDemo/message"
	"fmt"
	"log"
	"net"
)

const (
	addr = "0.0.0.0"
	port = "8080"
	size = 10000

	inBuffSize = 20
)

const (
	NormalMsg  = 0x01
)

var (
	welcome = "You are welcome. I'm server."
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

type Pool struct {
	set []*message.DataRW
}

func poolIns() Pool {
	p := Pool{}
	p.set = make([]*message.DataRW, 0, size)

	return p
}

func (p *Pool) add(rw *message.DataRW) {
	p.set = append(p.set, rw)
}

type Server struct {
	listenAddr string
	port       string

	pool  Pool
	actor Actor

	// connection go routine
	in chan net.Conn
}

func InsOfServer() Server {
	s := Server{}
	s.listenAddr = addr
	s.port = port

	s.pool = poolIns()
	s.actor = Actor{}

	s.in = make(chan net.Conn, inBuffSize)

	return s
}

func (s *Server) addr() string {
	return addr + ":" + port
}

func (s *Server) Start() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", s.addr())
	if err != nil {
		log.Fatalf("net.ResovleTCPAddr fail:%s", addr)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("listen %s fail: %s", addr, err)
	} else {

		log.Println("rpc listening", addr)
	}

	go s.acceptLoop()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept error:", err)
			continue
		}

		fmt.Println("client connect in: ", conn.RemoteAddr())
		s.in <- conn
	}
}

// server to accept a new connection
func (s *Server) acceptLoop() {
	// Wait for an error or disconnect.
loop:
	for {
		select {
		case c := <-s.in:
			s.handleConnection(c)
			break loop
		}
	}
}

func (s *Server) handleConnection(c net.Conn) {
	// init read-writer
	rw := message.DataRWIns(c)

	err := message.Send(rw, NormalMsg, welcome)
	if err != nil {
		fmt.Println("error send welcome info. ")
		return
	}

	go s.readLoop(rw)
}

func (s *Server) readLoop(rw *message.DataRW) {
	for {
		m, err := rw.ReadMsg()
		if err != nil {
			fmt.Println("error reading msg. ", err.Error())
			return
		}

		fmt.Println("msg code:", m.Code)
		fmt.Println("info: ", m.Payload)

		resp := s.actor.resp("hello")
		fmt.Println("resp message :", resp)

		err = message.Send(rw, NormalMsg, resp)
		if err != nil {
			fmt.Println("error response: ", err.Error())
			return
		}
	}
}

func Start() {
	s := InsOfServer()
	s.Start()
}
