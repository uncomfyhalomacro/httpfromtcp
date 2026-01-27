package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	port       int
	addr       string
	listener   net.Listener
	connection net.Conn
	connected  atomic.Bool
}

func Serve(addr string, port int) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return nil, err
	}
	server := Server{
		port:       port,
		addr:       addr,
		listener:   l,
		connection: nil,
		connected:  atomic.Bool{},
	}
	go server.listen()
	return &server, nil
}

func (s *Server) Close() error {
	err := s.listener.Close()

	isAlreadyConnected := s.connected.Load()
	swapped := s.connected.CompareAndSwap(isAlreadyConnected, false)
	if swapped {
		log.Println("Connection Closed")
	}
	return err
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("errors: %v\n", err)
		}
		isAlreadyConnected := s.connected.Load()
		swapped := s.connected.CompareAndSwap(isAlreadyConnected, true)
		if swapped {
			log.Println("Connection Established")
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	isServerConnected := s.connected.Load()
	if isServerConnected {
		response := []byte(`HTTP/1.1 200 OK
Content-Type: text/plain

Hello World!
`)
		conn.Write(response)
		conn.Close()
	} else {
		conn.Close()
	}
}
