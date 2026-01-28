package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/uncomfyhalomacro/httpfromtcp/internal/request"
	"github.com/uncomfyhalomacro/httpfromtcp/internal/response"
)

type Server struct {
	port       int
	addr       string
	listener   net.Listener
	connection net.Conn
	connected  atomic.Bool
	handler    func(target string) Handler
}

func Serve(addr string, port int, handler func(target string) Handler) (*Server, error) {
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
		handler:    handler,
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
		h := response.GetDefaultHeaders(0)
		req, err := request.RequestFromReader(conn)
		if err != nil {
			response.WriteStatusLine(conn, response.BadRequest)
			response.WriteHeaders(conn, h)
			conn.Write([]byte("Woopsie, my bad\n"))
			conn.Close()
		}
		hErr := s.handler(req.RequestLine.RequestTarget)(conn, req)
		if hErr != nil {
			response.WriteStatusLine(conn, response.StatusCode(hErr.StatusCode))
			hlen := len(hErr.Message)
			h["content-length"] = fmt.Sprintf("%d", hlen)
			response.WriteHeaders(conn, h)
		}
		conn.Close()
	} else {
		conn.Close()
	}
}
