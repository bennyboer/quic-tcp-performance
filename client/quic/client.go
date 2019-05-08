package quic

import (
	"bytes"
	"encoding/gob"
	"github.com/bennyboer/quic-tcp-performance/util/cli"
	"github.com/bennyboer/quic-tcp-performance/util/connection_type"
	"github.com/lucas-clemente/quic-go"
	"io"
)

// QUIC Client implementation.
type Client struct {
	// Current QUIC session.
	session quic.Session
}

// Create new QUIC client.
func NewClient(options *cli.Options) (*Client, error) {
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

func (c *Client) SendSync(message *[]byte) (*[]byte, error) {
	stream, err := c.session.OpenStreamSync()
	if err != nil {
		return nil, err
	}

	// Convert message to byte array
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	if err := enc.Encode(message); err != nil {
		return nil, err
	}

	_, err = stream.Write(buffer.Bytes())
	if err != nil {
		return nil, err
	}

	response := &bytes.Buffer{}
	_, err = io.Copy(stream, response)
	if err != nil {
		return nil, err
	}

	responseBytes := response.Bytes()

	return &responseBytes, nil
}
