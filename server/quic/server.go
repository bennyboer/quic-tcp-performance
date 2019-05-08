package quic

import (
	"crypto/tls"
	"fmt"
	"github.com/bennyboer/quic-tcp-performance/util/cli"
	"github.com/bennyboer/quic-tcp-performance/util/connection_type"
	"github.com/lucas-clemente/quic-go"
	"io"
	"log"
	"sync"
)

// QUIC Server implementation.
type Server struct {
	// TLS configuration to use
	tlsConfig tls.Config
}

// Create new QUIC server.
func NewServer(options *cli.Options) (*Server, error) {
	server := Server{
		tlsConfig: options.TlsConfiguration,
	}

	return &server, nil
}

func (s *Server) GetType() connection_type.ConnectionType {
	return connection_type.QUIC
}

func (s *Server) Listen(addr *string) (*sync.WaitGroup, error) {
	listener, err := quic.ListenAddr(*addr, &s.tlsConfig, nil)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go s.listen(listener, &wg)

	return &wg, nil
}

// Start listening to incoming connections.
func (s *Server) listen(listener quic.Listener, wg *sync.WaitGroup) {
	for {
		sess, err := listener.Accept()
		if err != nil {
			log.Println("QUIC server failed while listening to incoming connection requests. Cancelling listening.")
			break
		}

		wg.Add(1)
		go s.inSession(&sess, wg)
	}
}

func (s *Server) inSession(session *quic.Session, wg *sync.WaitGroup) {
	stream, err := (*session).AcceptStream()
	if err != nil {
		log.Println("QUIC session failed while trying to accept stream. Cancelling session.")
		wg.Done()
		return
	}

	_, err = io.Copy(loggingWriter{stream}, stream)

	wg.Done()
}

type loggingWriter struct{ io.Writer }

func (w loggingWriter) Write(b []byte) (int, error) {
	fmt.Printf("Server: Got '%s'\n", string(b))
	return w.Writer.Write(b)
}
