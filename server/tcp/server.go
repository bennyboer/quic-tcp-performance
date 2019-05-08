package tcp

import (
	"github.com/bennyboer/quic-tcp-performance/util/cli"
	"github.com/bennyboer/quic-tcp-performance/util/connection_type"
	"sync"
)

// TCP server implementation.
type Server struct {
	wg sync.WaitGroup
}

// Create TCP server.
func NewServer(options *cli.Options) (*Server, error) {
	server := Server{}

	return &server, nil
}

func (s *Server) GetType() connection_type.ConnectionType {
	return connection_type.TCP
}

func (s *Server) Listen(addr *string) (*sync.WaitGroup, error) {
	panic("implement me")
}
