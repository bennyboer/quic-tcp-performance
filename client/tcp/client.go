package tcp

import (
	"bytes"
	"crypto/tls"
	"github.com/bennyboer/quic-tcp-performance/util/cli"
	"github.com/bennyboer/quic-tcp-performance/util/connection_type"
	"log"
	"net"
	"time"
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

func (c *Client) SendDuration(duration time.Duration, bufferSize int) (int64, error) {
	if bufferSize <= 0 {
		bufferSize = 1 // As small as possible
	}

	// Write buffer first
	var buffer bytes.Buffer
	for i := 0; i < bufferSize; i++ {
		buffer.WriteByte(0)
	}
	b := buffer.Bytes()

	end := time.Now().Add(duration)

	var i int64 = 0
	for {
		_, err := c.conn.Write(b)
		if err != nil {
			return -1, err
		}

		if time.Now().After(end) {
			break
		}

		i++
	}

	return i * int64(bufferSize), nil
}

func (c *Client) SendBytes(numBytes int64) (time.Duration, error) {
	var buffer bytes.Buffer
	var i int64 = 0
	for ; i < numBytes; i++ {
		buffer.WriteByte(0) // Write zero byte
	}
	b := buffer.Bytes()

	start := time.Now()

	_, err := c.conn.Write(b)
	if err != nil {
		return -1, err
	}

	return time.Since(start), nil
}

func (c *Client) Cleanup() error {
	return c.conn.Close()
}
