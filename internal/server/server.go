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
	wg       sync.WaitGroup
	quit     chan interface{}
}

func New(address string, handler *handler.Handler, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	quit := make(chan interface{})

	return &Server{address: address, handler: handler, logger: logger, quit: quit}

}

func (s *Server) Start() error {
	op := "server.Start"
	var err error
	s.listener, err = net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal(err)
	}
	s.logger.Info(op, "Listening on ", s.address)
	for {
		listener, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return nil
			default:
				s.logger.Error(op, "Error accepting:", err.Error())
				continue
			}
		}
		s.logger.Info(op, "Accepting connections on", listener.RemoteAddr().String())
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.handleConn(listener)
		}()
	}
}

func (s *Server) handleConn(conn net.Conn) {
	op := "server.handleConn"
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			s.logger.Error(op, "Error closing connection:", err.Error())
		}
	}(conn)
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				s.logger.Error(op, "Error reading from connection:", err.Error())
			}
			s.logger.Info(op, "Connection closed", conn.RemoteAddr().String())
			break
		}
		resp, err := s.handler.Process(line)
		if err != nil {
			s.logger.Error(op, "Error processing:", err.Error())

		} else {
			s.logger.Info(op, conn.RemoteAddr().String(), resp)
			_, err := conn.Write([]byte(resp))
			if err != nil {
				s.logger.Error(op, "unable to write in conn", err)
			}
		}
	}
}

func (s *Server) Stop() {
	op := "server.Stop"
	close(s.quit)
	err := s.listener.Close()
	if err != nil {
		s.logger.Error(op, "Error closing listener:", err.Error())
	}
	s.wg.Wait()
}
