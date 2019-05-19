package client

import (
	"errors"
	"fmt"
	"github.com/bennyboer/quic-tcp-performance/client/quic"
	"github.com/bennyboer/quic-tcp-performance/client/tcp"
	"github.com/bennyboer/quic-tcp-performance/util/cli"
	"github.com/bennyboer/quic-tcp-performance/util/connection_type"
	"time"
)

// Client of the QUIC / TCP measurement tool.
type Client interface {
	// Send bytes for the passed time and return the sent bytes
	SendDuration(duration time.Duration, bufferSize int) (int64, error)
	// Send the passed amount of bytes and measure the time it took
	SendBytes(numBytes int64) (time.Duration, error)
	// Get the type of the connection
	GetType() connection_type.ConnectionType
	// Clean up the connection
	Cleanup() error
}

// Create new client of the passed type which connects to the passed address.
func NewClient(options *cli.Options) (Client, error) {
	switch options.ConnectionType {
	case connection_type.TCP:
		return tcp.NewClient(options)
	case connection_type.QUIC:
		return quic.NewClient(options)
	default:
		return nil, errors.New(fmt.Sprintf("Connection type %d unknown", options.ConnectionType))
	}
}
