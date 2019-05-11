package tcp

import (
	"bytes"
	"crypto/tls"
	"github.com/bennyboer/quic-tcp-performance/util/cli"
	"github.com/bennyboer/quic-tcp-performance/util/connection_type"
	"io"
	"log"
	"net"
)

const protocol = "tcp"

// TCP Client implementation
type Client struct {
	// Current TCP connection
	conn net.Conn
}

func NewClient(options *cli.Options) (*Client, error) {
	var conn net.Conn
	var err error
	if options.TlsEnabled {
		log.Println("Connecting via TCP over TLS")
		conn, err = tls.Dial(protocol, options.Address, &options.TlsConfiguration)
	} else {
		log.Println("Connection via TCP")
		conn, err = net.Dial(protocol, options.Address)
	}
	if err != nil {
		return nil, err
	}

	client := Client{
		conn: conn,
	}

	return &client, nil
}

func (c *Client) GetType() connection_type.ConnectionType {
	return connection_type.TCP
}

func (c *Client) SendSync(message *[]byte) (*[]byte, error) {
	_, err := c.conn.Write(*message)
	if err != nil {
		return nil, nil
	}

	buffer := bytes.Buffer{}
	_, err = io.Copy(c.conn, &buffer)
	if err != nil {
		return nil, nil
	}

	responseBytes := buffer.Bytes()

	return &responseBytes, nil
}
