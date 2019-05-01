package main

import (
	"flag"
	"fmt"
	"github.com/bennyboer/quic-tcp-performance/client"
	"github.com/bennyboer/quic-tcp-performance/server"
)

const (
	defaultServerAddress = "localhost:19191"
)

func main() {
	isServerMode := flag.Bool("server", false, "Whether the measurement tool should be started in server mode")
	addr := flag.String("address", defaultServerAddress, "Address at which to bind the server (if started in server mode) or at which to connect (if started in client mode)")

	flag.Parse()

	if *isServerMode {
		fmt.Println("Tool started in SERVER mode")

		_, e := server.NewServer(addr)
		if e != nil {
			panic(e)
		}
	} else {
		fmt.Println("Tool started in CLIENT mode")

		c, e := client.NewClient(addr)
		if e != nil {
			panic(e)
		}

		message := "Hello World"
		response, e := c.Send(&message)
		if e != nil {
			panic(e)
		}

		fmt.Printf("Server responded with %s\n", *response)
	}
}
