package tcp

import (
	"bytes"
	"github.com/bennyboer/quic-tcp-performance/util/cli"
	"github.com/bennyboer/quic-tcp-performance/util/connection_type"
	"io"
	"net"
)

const protocol = "tcp"

// TCP Client implementation
type Client struct {
	// Current TCP connection
	con net.Conn
}

func NewClient(options *cli.Options) (*Client, error) {
	con, err := net.Dial(protocol, options.Address)
	if err != nil {
		return nil, err
	}

	client := Client{
		con: con,
	}

	return &client, nil
}

func (c *Client) GetType() connection_type.ConnectionType {
	return connection_type.TCP
}

func (c *Client) SendSync(message *[]byte) (*[]byte, error) {
	_, err := c.con.Write(*message)
	if err != nil {
		return nil, nil
	}

	buffer := bytes.Buffer{}
	_, err = io.Copy(c.con, &buffer)
	if err != nil {
		return nil, nil
	}

	responseBytes := buffer.Bytes()

	return &responseBytes, nil
}
