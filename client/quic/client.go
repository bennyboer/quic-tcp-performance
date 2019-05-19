package quic

import (
	"bytes"
	"github.com/bennyboer/quic-tcp-performance/util/cli"
	"github.com/bennyboer/quic-tcp-performance/util/connection_type"
	"github.com/lucas-clemente/quic-go"
	"log"
	"time"
)

// QUIC Client implementation.
type Client struct {
	// Current QUIC session.
	session quic.Session
}

// Create new QUIC client.
func NewClient(options *cli.Options) (*Client, error) {
	log.Println("Connecting via QUIC")
	session, err := quic.DialAddr(options.Address, &options.TlsConfiguration, nil)
	if err != nil {
		return nil, err
	}

	client := Client{
		session: session,
	}

	return &client, nil
}

func (c *Client) GetType() connection_type.ConnectionType {
	return connection_type.QUIC
}

func (c *Client) SendDuration(duration time.Duration, bufferSize int) (int64, error) {
	stream, err := c.session.OpenStreamSync()
	if err != nil {
		return -1, err
	}

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
		_, err = stream.Write(b)
		if err != nil {
			return -1, err
		}

		if time.Now().After(end) {
			break
		}

		i++
	}

	err = stream.Close()
	if err != nil {
		return -1, err
	}

	return i * int64(bufferSize), nil
}

func (c *Client) SendBytes(numBytes int64) (time.Duration, error) {
	stream, err := c.session.OpenStreamSync()
	if err != nil {
		return -1, err
	}

	var buffer bytes.Buffer
	var i int64 = 0
	for ; i < numBytes; i++ {
		buffer.WriteByte(0) // Write zero byte
	}
	b := buffer.Bytes()

	start := time.Now()

	_, err = stream.Write(b)
	if err != nil {
		return -1, err
	}

	err = stream.Close()
	if err != nil {
		return -1, err
	}

	return time.Since(start), nil
}

func (c *Client) Cleanup() error {
	return c.session.Close()
}
