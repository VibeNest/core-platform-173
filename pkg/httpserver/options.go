package httpserver

import "net"

// Option позволяет настраивать сервер
type Option func(*Server)

// Port задает порт для прослушивания
func Port(port string) Option {
	return func(s *Server) {
		s.server.Addr = net.JoinHostPort("", port)
	}
}
