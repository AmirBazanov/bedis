package server

import (
	"bedis/internal/handler"
	"bufio"
	"io"
	"log"
	"log/slog"
	"net"
	"sync"
)

type Server struct {
	address  string
	handler  *handler.Handler
	logger   *slog.Logger
	listener net.Listener
	wg       *sync.WaitGroup
	quit     chan interface{}
}

func New(address string, handler *handler.Handler, logger *slog.Logger) *Server {
	op := "server.New"

	if logger == nil {
		logger = slog.New(
			slog.NewTextHandler(io.Discard, nil),
		)
	}
	wait := &sync.WaitGroup{}
	quit := make(chan interface{})
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info(op, "Listening on ", address)
	return &Server{address, handler, logger, listener, wait, quit}

}

func (s *Server) Start() error {
	op := "server.Start"
	defer s.wg.Done()
	for {
		listener, err := s.listener.Accept()
		s.logger.Info(op, "Accepting connections on ", s.address)
		if err != nil {
			select {
			case <-s.quit:
				return nil
			default:
				s.logger.Error(op, "Error accepting: ", err.Error())
				continue
			}
		}
		s.wg.Add(1)
		go func() {

			s.handleConn(listener)
			s.wg.Done()
		}()
	}
}

func (s *Server) handleConn(conn net.Conn) {
	op := "server.handleConn"
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			s.logger.Error(op, "Error closing connection: ", err.Error())
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

func (s *Server) Stop() {
	op := "server.Stop"
	close(s.quit)
	err := s.listener.Close()
	if err != nil {
		s.logger.Error(op, "Error closing listener: ", err.Error())
	}
	s.wg.Wait()
}
