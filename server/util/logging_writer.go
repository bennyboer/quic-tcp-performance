package util

import (
	"fmt"
	"io"
)

type LoggingWriter struct {
	Writer io.Writer
}

func (w LoggingWriter) Write(b []byte) (int, error) {
	fmt.Printf("Server: Got '%s'\n", string(b))
	return w.Writer.Write(b)
}
