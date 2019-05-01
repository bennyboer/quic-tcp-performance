package client

import (
	"crypto/tls"
	"github.com/lucas-clemente/quic-go"
	"io"
)

// Client of the QUIC / TCP measurement tool.
type Client struct {
	session quic.Session
}

// Create new client.
func NewClient(addr *string) (*Client, error) {
	session, err := quic.DialAddr(*addr, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		return nil, err
	}

	client := Client{
		session: session,
	}

	return &client, nil
}

func (c *Client) Send(message *string) (*string, error) {
	stream, err := c.session.OpenStreamSync()
	if err != nil {
		return nil, err
	}

	_, err = stream.Write([]byte(*message))
	if err != nil {
		return nil, err
	}

	buf := make([]byte, len(*message))
	_, err = io.ReadFull(stream, buf)
	if err != nil {
		return nil, err
	}

	response := string(buf)

	return &response, nil
}
