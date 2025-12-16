package server

import (
	"bedis/internal/handler"
	"bufio"
	"log/slog"
	"net"
)

type Server struct {
	address string
	handler *handler.Handler
	logger  *slog.Logger
}

func New(address string, handler *handler.Handler, logger *slog.Logger) *Server {
	return &Server{address, handler, logger}

}

func (s *Server) Start() error {
	op := "server.Start"
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	s.logger.Info(op, "Listening on ", s.address)
	for {
		listener, err := listener.Accept()
		s.logger.Info(op, "Accepting connections on ", s.address)
		if err != nil {
			s.logger.Error(op, "Error accepting: ", err.Error())
			continue
		}
		s.handleConn(listener)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	op := "server.handleConn"
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
			s.logger.Error(op, "Error processing: ", err.Error())

		} else {

			s.logger.Info(op, conn.RemoteAddr().String(), resp)
		}
	}
}
