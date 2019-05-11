package tcp

import (
	"crypto/tls"
	"github.com/bennyboer/quic-tcp-performance/server/util"
	"github.com/bennyboer/quic-tcp-performance/util/cli"
	"github.com/bennyboer/quic-tcp-performance/util/connection_type"
	"io"
	"log"
	"net"
	"sync"
)

const network = "tcp"

// TCP server implementation.
type Server struct {
	// Whether TLS should be used
	useTls bool
	// TLS configuration to use
	tlsConfig tls.Config
}

// Create TCP server.
func NewServer(options *cli.Options) (*Server, error) {
	server := Server{
		useTls:    options.TlsEnabled,
		tlsConfig: options.TlsConfiguration,
	}

	return &server, nil
}

func (s *Server) GetType() connection_type.ConnectionType {
	return connection_type.TCP
}

func (s *Server) Listen(addr *string) (*sync.WaitGroup, error) {
	var listener net.Listener
	var err error
	if s.useTls {
		log.Println("Setting up TCP listener over TLS")
		listener, err = tls.Listen(network, *addr, &s.tlsConfig)
	} else {
		log.Println("Setting up TCP listener")
		listener, err = net.Listen(network, *addr)
	}

	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go s.listen(listener, &wg)

	return &wg, nil
}

// Start listening to incoming connections.
func (s *Server) listen(listener net.Listener, wg *sync.WaitGroup) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("TCP server failed while listening to incoming connection requests. Cancelling listening.")
			break
		}

		wg.Add(1)
		go s.inConnection(&conn, wg)
	}
}

func (s *Server) inConnection(conn *net.Conn, wg *sync.WaitGroup) {
	_, err := io.Copy(util.LoggingWriter{
		Writer: *conn,
	}, *conn)
	if err != nil {
		log.Println("Error while reading from TCP socket")
	}

	wg.Done()
}
