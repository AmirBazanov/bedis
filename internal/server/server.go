package server

import (
	"bedis/internal/handler"
	"bufio"
	"fmt"
	"net"
)

type Server struct {
	address string
	handler *handler.Handler
}

func New(address string, handler *handler.Handler) *Server {
	return &Server{address, handler}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	fmt.Println("Listening on " + s.address)
	for {
		listener, err := listener.Accept()
		fmt.Println("Accepting connections on " + s.address)
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		s.handleConn(listener)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	reader := bufio.NewScanner(conn)
	for reader.Scan() {
		line := reader.Text()
		resp, err := s.handler.Process(line)
		if err != nil {
			fmt.Println("Error processing: ", err.Error())
		}
		fmt.Println(conn, resp)
	}
}
