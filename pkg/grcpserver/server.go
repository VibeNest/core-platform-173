package grpcserver

import (
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	server *grpc.Server
	notify chan error
	port   string
}

func New(server *grpc.Server, port string) *Server {
	s := &Server{
		server: server,
		notify: make(chan error, 1),
		port:   port,
	}
	s.start()
	return s
}

func (s *Server) start() {
	go func() {
		l, err := net.Listen("tcp", ":"+s.port)
		if err != nil {
			s.notify <- err
			return
		}
		s.notify <- s.server.Serve(l)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() {
	s.server.GracefulStop()
}
