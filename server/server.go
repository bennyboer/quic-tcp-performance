package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"io"
	"math/big"
)

const (
	rsaKeySize = 4096
)

// Server of the QUIC / TCP measurement tool.
type Server struct {
	listener quic.Listener
}

// Create new server.
func NewServer(addr *string) (*Server, error) {
	tlsConfig, err := generateTLSConfig()
	if err != nil {
		return nil, err
	}

	listener, err := quic.ListenAddr(*addr, tlsConfig, nil)
	if err != nil {
		return nil, err
	}

	sess, err := listener.Accept()
	if err != nil {
		return nil, err
	}
	stream, err := sess.AcceptStream()
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(loggingWriter{stream}, stream)

	return &Server{
		listener: listener,
	}, nil
}

// Setup TLS config for the server.
func generateTLSConfig() (*tls.Config, error) {
	key, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}

	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}, nil
}

type loggingWriter struct{ io.Writer }

func (w loggingWriter) Write(b []byte) (int, error) {
	fmt.Printf("Server: Got '%s'\n", string(b))
	return w.Writer.Write(b)
}
