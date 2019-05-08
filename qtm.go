package main

import (
	"fmt"
	"github.com/bennyboer/quic-tcp-performance/client"
	"github.com/bennyboer/quic-tcp-performance/server"
	"github.com/bennyboer/quic-tcp-performance/util/cli"
	"log"
)

func main() {
	opt := cli.ParseOptions()

	if opt.IsServerMode {
		log.Println("Tool started in SERVER mode")

		s, err := server.NewServer(opt)
		if err != nil {
			log.Fatalln(err.Error())
		}

		wg, err := s.Listen(&opt.Address)
		if err != nil {
			log.Fatalf("Server could not start listening: %s", err.Error())
		}

		wg.Wait() // Wait for server termination
	} else {
		log.Println("Tool started in CLIENT mode")

		c, err := client.NewClient(opt)
		if err != nil {
			log.Fatalln(err.Error())
		}

		message := "Hello World"
		byteMessage := []byte(message)
		response, e := c.SendSync(&byteMessage)
		if e != nil {
			panic(e)
		}

		fmt.Printf("Server responded with %s\n", string(*response))
	}
}
