package server

import (
	"errors"
	"fmt"
	"github.com/bennyboer/quic-tcp-performance/server/quic"
	"github.com/bennyboer/quic-tcp-performance/server/tcp"
	"github.com/bennyboer/quic-tcp-performance/util/cli"
	"github.com/bennyboer/quic-tcp-performance/util/connection_type"
	"sync"
)

// Interface for servers in the tool.
type Server interface {
	GetType() connection_type.ConnectionType
	Listen(addr *string) (*sync.WaitGroup, error)
}

func NewServer(options *cli.Options) (Server, error) {
	switch options.ConnectionType {
	case connection_type.TCP:
		return tcp.NewServer(options)
	case connection_type.QUIC:
		return quic.NewServer(options)
	default:
		return nil, errors.New(fmt.Sprintf("Connection type %d unknown", options.ConnectionType))
	}
}
