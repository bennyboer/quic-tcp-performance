package main

import (
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

		if opt.Bytes > -1 {
			// Send the set amount of bytes to the server
			time, err := c.SendBytes(opt.Bytes)
			if err != nil {
				log.Fatalf("Encountered error when trying to send %d bytes to the server. Error: %s", opt.Bytes, err.Error())
			}

			log.Printf("Sent %d bytes in %d nanoseconds", opt.Bytes, time.Nanoseconds());
		} else if opt.Duration > -1 {
			// Send for the set duration to the server
			sentBytes, err := c.SendDuration(opt.Duration, opt.BufferSize)
			if err != nil {
				log.Fatalf("Encountered error when trying to send for %d ns to the server. Error: %s", opt.Duration.Nanoseconds(), err.Error())
			}

			log.Printf("Sent %d bytes in %d nanoseconds", sentBytes, opt.Duration.Nanoseconds());
		} else {
			log.Fatalf("You need to either set --bytes or --duration to measure throughput")
		}

		err = c.Cleanup()
		if err != nil {
			log.Fatalf("Could not cleanup client properly. Error: %s", err.Error())
		}
	}
}
