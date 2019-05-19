package util

import (
	"fmt"
	"io"
)

type LoggingWriter struct {
	Writer io.Writer
}

func (w LoggingWriter) Write(b []byte) (int, error) {
	fmt.Printf("Received %d bytes\n", len(b))
	return len(b), nil
}
