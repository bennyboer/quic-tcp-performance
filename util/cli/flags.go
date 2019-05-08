package cli

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"github.com/bennyboer/quic-tcp-performance/util/connection_type"
	"log"
	"math/big"
)

const (
	rsaKeySize                = 4096
	defaultServerAddress      = "localhost:19191"
	defaultConnectionTypeName = "QUIC"
)

// Options for the CLI.
type Options struct {
	// Whether the tool should run as server or client
	IsServerMode bool

	// Client Mode: Address of the corresponding server to connect to.
	// Server Mode: Address to listen for incoming connections.
	Address string

	// Type of the connection to test
	ConnectionType connection_type.ConnectionType

	// Whether TLS should be enabled (Cannot be enforced if the protocol to use is not supporting it or without it)
	// For example QUIC is enforcing TLS use, thus it cannot be disabled.
	TlsEnabled bool

	// TLS configuration to use
	TlsConfiguration tls.Config
}

// Parse the CLI options.
func ParseOptions() *Options {
	isServerMode := flag.Bool("server", false, "Whether the measurement tool should be started in server mode")
	address := flag.String("address", defaultServerAddress, "Address at which to bind the server (if started in server mode) or at which to connect to (if started in client mode)")
	connectionTypeName := flag.String("type", defaultConnectionTypeName, "type of the connection (Either 'QUIC' or 'TCP')")
	tlsEnabled := flag.Bool("tls", true, "Whether TLS is enabled (Cannot be turned of for some protocols (e. g. QUIC)")

	// TODO Parse certificate and private key flags for custom TLS configuration

	flag.Parse()

	// Parse connection type
	var connectionType connection_type.ConnectionType
	switch *connectionTypeName {
	case "QUIC":
		connectionType = connection_type.QUIC
	case "TCP":
		connectionType = connection_type.TCP
	default:
		log.Fatalf("Could not understand connection type '%s', use either 'QUIC' or 'TCP'", *connectionTypeName)
	}

	// Generate tls configuration
	var tlsConfiguration tls.Config
	if *isServerMode {
		generatedTlsConfig, err := generateDefaultTLSConfig()
		if err != nil {
			log.Fatalf("Could not generate TLS configuration")
		}

		tlsConfiguration = *generatedTlsConfig
	} else {
		tlsConfiguration = tls.Config{
			InsecureSkipVerify: true,
		}
	}

	return &Options{
		IsServerMode:     *isServerMode,
		Address:          *address,
		ConnectionType:   connectionType,
		TlsEnabled:       *tlsEnabled,
		TlsConfiguration: tlsConfiguration,
	}
}

// Setup default TLS configuration.
func generateDefaultTLSConfig() (*tls.Config, error) {
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
